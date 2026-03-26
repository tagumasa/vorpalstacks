// This is a generated file - do not edit.
//
// Generated from kinesis.proto.

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
import 'kinesis.pb.dart' as $0;

export 'kinesis.pb.dart';

/// KinesisService provides kinesis API operations.
@$pb.GrpcServiceName('kinesis.KinesisService')
class KinesisServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  KinesisServiceClient(super.channel, {super.options, super.interceptors});

  /// Adds or updates tags for the specified Kinesis data stream. You can assign up to 50 tags to a data stream. When invoking this API, you must use either the StreamARN or the StreamName parameter, or ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> addTagsToStream(
    $0.AddTagsToStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addTagsToStream, request, options: options);
  }

  /// Creates a Kinesis data stream. A stream captures and transports data records that are continuously emitted from different data sources or producers. Scale-out within a stream is explicitly supporte...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> createStream(
    $0.CreateStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createStream, request, options: options);
  }

  /// Decreases the Kinesis data stream's retention period, which is the length of time data records are accessible after they are added to the stream. The minimum value of a stream's retention period is...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> decreaseStreamRetentionPeriod(
    $0.DecreaseStreamRetentionPeriodInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$decreaseStreamRetentionPeriod, request,
        options: options);
  }

  /// Delete a policy for the specified data stream or consumer. Request patterns can be one of the following: Data stream pattern: arn:aws.*:kinesis:.*:\d{12}:.*stream/\S+ Consumer pattern: ^(arn):aws.*...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteResourcePolicy(
    $0.DeleteResourcePolicyInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourcePolicy, request, options: options);
  }

  /// Deletes a Kinesis data stream and all its shards and data. You must shut down any applications that are operating on the stream before you delete the stream. If an application attempts to operate o...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteStream(
    $0.DeleteStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteStream, request, options: options);
  }

  /// To deregister a consumer, provide its ARN. Alternatively, you can provide the ARN of the data stream and the name you gave the consumer when you registered it. You may also provide all three parame...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deregisterStreamConsumer(
    $0.DeregisterStreamConsumerInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deregisterStreamConsumer, request,
        options: options);
  }

  /// Describes the account-level settings for Amazon Kinesis Data Streams. This operation returns information about the minimum throughput billing commitments and other account-level configurations. Thi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAccountSettingsOutput>
      describeAccountSettings(
    $0.DescribeAccountSettingsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAccountSettings, request,
        options: options);
  }

  /// Describes the shard limits and usage for the account. If you update your account limits, the old limits might be returned for a few minutes. This operation has a limit of one transaction per second...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeLimitsOutput> describeLimits(
    $0.DescribeLimitsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeLimits, request, options: options);
  }

  /// Describes the specified Kinesis data stream. This API has been revised. It's highly recommended that you use the DescribeStreamSummary API to get a summarized description of the specified Kinesis d...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeStreamOutput> describeStream(
    $0.DescribeStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeStream, request, options: options);
  }

  /// To get the description of a registered consumer, provide the ARN of the consumer. Alternatively, you can provide the ARN of the data stream and the name you gave the consumer when you registered it...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeStreamConsumerOutput> describeStreamConsumer(
    $0.DescribeStreamConsumerInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeStreamConsumer, request,
        options: options);
  }

  /// Provides a summarized description of the specified Kinesis data stream without the shard list. When invoking this API, you must use either the StreamARN or the StreamName parameter, or both. It is ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeStreamSummaryOutput> describeStreamSummary(
    $0.DescribeStreamSummaryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeStreamSummary, request, options: options);
  }

  /// Disables enhanced monitoring. When invoking this API, you must use either the StreamARN or the StreamName parameter, or both. It is recommended that you use the StreamARN input parameter when you i...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.EnhancedMonitoringOutput> disableEnhancedMonitoring(
    $0.DisableEnhancedMonitoringInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableEnhancedMonitoring, request,
        options: options);
  }

  /// Enables enhanced Kinesis data stream monitoring for shard-level metrics. When invoking this API, you must use either the StreamARN or the StreamName parameter, or both. It is recommended that you u...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.EnhancedMonitoringOutput> enableEnhancedMonitoring(
    $0.EnableEnhancedMonitoringInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableEnhancedMonitoring, request,
        options: options);
  }

  /// Gets data records from a Kinesis data stream's shard. When invoking this API, you must use either the StreamARN or the StreamName parameter, or both. It is recommended that you use the StreamARN in...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRecordsOutput> getRecords(
    $0.GetRecordsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRecords, request, options: options);
  }

  /// Returns a policy attached to the specified data stream or consumer. Request patterns can be one of the following: Data stream pattern: arn:aws.*:kinesis:.*:\d{12}:.*stream/\S+ Consumer pattern: ^(a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetResourcePolicyOutput> getResourcePolicy(
    $0.GetResourcePolicyInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResourcePolicy, request, options: options);
  }

  /// Gets an Amazon Kinesis shard iterator. A shard iterator expires 5 minutes after it is returned to the requester. When invoking this API, you must use either the StreamARN or the StreamName paramete...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetShardIteratorOutput> getShardIterator(
    $0.GetShardIteratorInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getShardIterator, request, options: options);
  }

  /// Increases the Kinesis data stream's retention period, which is the length of time data records are accessible after they are added to the stream. The maximum value of a stream's retention period is...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> increaseStreamRetentionPeriod(
    $0.IncreaseStreamRetentionPeriodInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$increaseStreamRetentionPeriod, request,
        options: options);
  }

  /// Lists the shards in a stream and provides information about each shard. This operation has a limit of 1000 transactions per second per data stream. When invoking this API, you must use either the S...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListShardsOutput> listShards(
    $0.ListShardsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listShards, request, options: options);
  }

  /// Lists the consumers registered to receive data from a stream using enhanced fan-out, and provides information about each consumer. This operation has a limit of 5 transactions per second per stream.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListStreamConsumersOutput> listStreamConsumers(
    $0.ListStreamConsumersInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listStreamConsumers, request, options: options);
  }

  /// Lists your Kinesis data streams. The number of streams may be too large to return from a single call to ListStreams. You can limit the number of returned streams using the Limit parameter. If you d...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListStreamsOutput> listStreams(
    $0.ListStreamsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listStreams, request, options: options);
  }

  /// List all tags added to the specified Kinesis resource. Each tag is a label consisting of a user-defined key and value. Tags can help you manage, identify, organize, search for, and filter resources...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceOutput> listTagsForResource(
    $0.ListTagsForResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Lists the tags for the specified Kinesis data stream. This operation has a limit of five transactions per second per account. When invoking this API, you must use either the StreamARN or the Stream...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForStreamOutput> listTagsForStream(
    $0.ListTagsForStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForStream, request, options: options);
  }

  /// Merges two adjacent shards in a Kinesis data stream and combines them into a single shard to reduce the stream's capacity to ingest and transport data. This API is only supported for the data strea...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> mergeShards(
    $0.MergeShardsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$mergeShards, request, options: options);
  }

  /// Writes a single data record into an Amazon Kinesis data stream. Call PutRecord to send data into the stream for real-time ingestion and subsequent processing, one record at a time. Each shard can s...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutRecordOutput> putRecord(
    $0.PutRecordInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putRecord, request, options: options);
  }

  /// Writes multiple data records into a Kinesis data stream in a single call (also referred to as a PutRecords request). Use this operation to send data into the stream for data ingestion and processin...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutRecordsOutput> putRecords(
    $0.PutRecordsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putRecords, request, options: options);
  }

  /// Attaches a resource-based policy to a data stream or registered consumer. If you are using an identity other than the root user of the Amazon Web Services account that owns the resource, the callin...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putResourcePolicy(
    $0.PutResourcePolicyInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putResourcePolicy, request, options: options);
  }

  /// Registers a consumer with a Kinesis data stream. When you use this operation, the consumer you register can then call SubscribeToShard to receive data from the stream using enhanced fan-out, at a r...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RegisterStreamConsumerOutput> registerStreamConsumer(
    $0.RegisterStreamConsumerInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerStreamConsumer, request,
        options: options);
  }

  /// Removes tags from the specified Kinesis data stream. Removed tags are deleted and cannot be recovered after this operation successfully completes. When invoking this API, you must use either the St...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> removeTagsFromStream(
    $0.RemoveTagsFromStreamInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeTagsFromStream, request, options: options);
  }

  /// Splits a shard into two new shards in the Kinesis data stream, to increase the stream's capacity to ingest and transport data. SplitShard is called when there is a need to increase the overall capa...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> splitShard(
    $0.SplitShardInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$splitShard, request, options: options);
  }

  /// Enables or updates server-side encryption using an Amazon Web Services KMS key for a specified stream. When invoking this API, you must use either the StreamARN or the StreamName parameter, or both...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> startStreamEncryption(
    $0.StartStreamEncryptionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startStreamEncryption, request, options: options);
  }

  /// Disables server-side encryption for a specified stream. When invoking this API, you must use either the StreamARN or the StreamName parameter, or both. It is recommended that you use the StreamARN ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> stopStreamEncryption(
    $0.StopStreamEncryptionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopStreamEncryption, request, options: options);
  }

  /// This operation establishes an HTTP/2 connection between the consumer you specify in the ConsumerARN parameter and the shard you specify in the ShardId parameter. After the connection is successfull...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SubscribeToShardOutput> subscribeToShard(
    $0.SubscribeToShardInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$subscribeToShard, request, options: options);
  }

  /// Adds or updates tags for the specified Kinesis resource. Each tag is a label consisting of a user-defined key and value. Tags can help you manage, identify, organize, search for, and filter resourc...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> tagResource(
    $0.TagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Removes tags from the specified Kinesis resource. Removed tags are deleted and can't be recovered after this operation completes successfully.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> untagResource(
    $0.UntagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Updates the account-level settings for Amazon Kinesis Data Streams. Updating account settings is a synchronous operation. Upon receiving the request, Kinesis Data Streams will return immediately wi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateAccountSettingsOutput> updateAccountSettings(
    $0.UpdateAccountSettingsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAccountSettings, request, options: options);
  }

  /// This allows you to update the MaxRecordSize of a single record that you can write to, and read from a stream. You can ingest and digest single records up to 10240 KiB.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> updateMaxRecordSize(
    $0.UpdateMaxRecordSizeInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateMaxRecordSize, request, options: options);
  }

  /// Updates the shard count of the specified stream to the specified number of shards. This API is only supported for the data streams with the provisioned capacity mode. When invoking this API, you mu...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateShardCountOutput> updateShardCount(
    $0.UpdateShardCountInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateShardCount, request, options: options);
  }

  /// Updates the capacity mode of the data stream. Currently, in Kinesis Data Streams, you can choose between an on-demand capacity mode and a provisioned capacity mode for your data stream. If you'd st...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> updateStreamMode(
    $0.UpdateStreamModeInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateStreamMode, request, options: options);
  }

  /// Updates the warm throughput configuration for the specified Amazon Kinesis Data Streams on-demand data stream. This operation allows you to proactively scale your on-demand data stream to a specifi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateStreamWarmThroughputOutput>
      updateStreamWarmThroughput(
    $0.UpdateStreamWarmThroughputInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateStreamWarmThroughput, request,
        options: options);
  }

  // method descriptors

  static final _$addTagsToStream =
      $grpc.ClientMethod<$0.AddTagsToStreamInput, $1.Empty>(
          '/kinesis.KinesisService/AddTagsToStream',
          ($0.AddTagsToStreamInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$createStream =
      $grpc.ClientMethod<$0.CreateStreamInput, $1.Empty>(
          '/kinesis.KinesisService/CreateStream',
          ($0.CreateStreamInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$decreaseStreamRetentionPeriod =
      $grpc.ClientMethod<$0.DecreaseStreamRetentionPeriodInput, $1.Empty>(
          '/kinesis.KinesisService/DecreaseStreamRetentionPeriod',
          ($0.DecreaseStreamRetentionPeriodInput value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteResourcePolicy =
      $grpc.ClientMethod<$0.DeleteResourcePolicyInput, $1.Empty>(
          '/kinesis.KinesisService/DeleteResourcePolicy',
          ($0.DeleteResourcePolicyInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteStream =
      $grpc.ClientMethod<$0.DeleteStreamInput, $1.Empty>(
          '/kinesis.KinesisService/DeleteStream',
          ($0.DeleteStreamInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deregisterStreamConsumer =
      $grpc.ClientMethod<$0.DeregisterStreamConsumerInput, $1.Empty>(
          '/kinesis.KinesisService/DeregisterStreamConsumer',
          ($0.DeregisterStreamConsumerInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$describeAccountSettings = $grpc.ClientMethod<
          $0.DescribeAccountSettingsInput, $0.DescribeAccountSettingsOutput>(
      '/kinesis.KinesisService/DescribeAccountSettings',
      ($0.DescribeAccountSettingsInput value) => value.writeToBuffer(),
      $0.DescribeAccountSettingsOutput.fromBuffer);
  static final _$describeLimits =
      $grpc.ClientMethod<$0.DescribeLimitsInput, $0.DescribeLimitsOutput>(
          '/kinesis.KinesisService/DescribeLimits',
          ($0.DescribeLimitsInput value) => value.writeToBuffer(),
          $0.DescribeLimitsOutput.fromBuffer);
  static final _$describeStream =
      $grpc.ClientMethod<$0.DescribeStreamInput, $0.DescribeStreamOutput>(
          '/kinesis.KinesisService/DescribeStream',
          ($0.DescribeStreamInput value) => value.writeToBuffer(),
          $0.DescribeStreamOutput.fromBuffer);
  static final _$describeStreamConsumer = $grpc.ClientMethod<
          $0.DescribeStreamConsumerInput, $0.DescribeStreamConsumerOutput>(
      '/kinesis.KinesisService/DescribeStreamConsumer',
      ($0.DescribeStreamConsumerInput value) => value.writeToBuffer(),
      $0.DescribeStreamConsumerOutput.fromBuffer);
  static final _$describeStreamSummary = $grpc.ClientMethod<
          $0.DescribeStreamSummaryInput, $0.DescribeStreamSummaryOutput>(
      '/kinesis.KinesisService/DescribeStreamSummary',
      ($0.DescribeStreamSummaryInput value) => value.writeToBuffer(),
      $0.DescribeStreamSummaryOutput.fromBuffer);
  static final _$disableEnhancedMonitoring = $grpc.ClientMethod<
          $0.DisableEnhancedMonitoringInput, $0.EnhancedMonitoringOutput>(
      '/kinesis.KinesisService/DisableEnhancedMonitoring',
      ($0.DisableEnhancedMonitoringInput value) => value.writeToBuffer(),
      $0.EnhancedMonitoringOutput.fromBuffer);
  static final _$enableEnhancedMonitoring = $grpc.ClientMethod<
          $0.EnableEnhancedMonitoringInput, $0.EnhancedMonitoringOutput>(
      '/kinesis.KinesisService/EnableEnhancedMonitoring',
      ($0.EnableEnhancedMonitoringInput value) => value.writeToBuffer(),
      $0.EnhancedMonitoringOutput.fromBuffer);
  static final _$getRecords =
      $grpc.ClientMethod<$0.GetRecordsInput, $0.GetRecordsOutput>(
          '/kinesis.KinesisService/GetRecords',
          ($0.GetRecordsInput value) => value.writeToBuffer(),
          $0.GetRecordsOutput.fromBuffer);
  static final _$getResourcePolicy =
      $grpc.ClientMethod<$0.GetResourcePolicyInput, $0.GetResourcePolicyOutput>(
          '/kinesis.KinesisService/GetResourcePolicy',
          ($0.GetResourcePolicyInput value) => value.writeToBuffer(),
          $0.GetResourcePolicyOutput.fromBuffer);
  static final _$getShardIterator =
      $grpc.ClientMethod<$0.GetShardIteratorInput, $0.GetShardIteratorOutput>(
          '/kinesis.KinesisService/GetShardIterator',
          ($0.GetShardIteratorInput value) => value.writeToBuffer(),
          $0.GetShardIteratorOutput.fromBuffer);
  static final _$increaseStreamRetentionPeriod =
      $grpc.ClientMethod<$0.IncreaseStreamRetentionPeriodInput, $1.Empty>(
          '/kinesis.KinesisService/IncreaseStreamRetentionPeriod',
          ($0.IncreaseStreamRetentionPeriodInput value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$listShards =
      $grpc.ClientMethod<$0.ListShardsInput, $0.ListShardsOutput>(
          '/kinesis.KinesisService/ListShards',
          ($0.ListShardsInput value) => value.writeToBuffer(),
          $0.ListShardsOutput.fromBuffer);
  static final _$listStreamConsumers = $grpc.ClientMethod<
          $0.ListStreamConsumersInput, $0.ListStreamConsumersOutput>(
      '/kinesis.KinesisService/ListStreamConsumers',
      ($0.ListStreamConsumersInput value) => value.writeToBuffer(),
      $0.ListStreamConsumersOutput.fromBuffer);
  static final _$listStreams =
      $grpc.ClientMethod<$0.ListStreamsInput, $0.ListStreamsOutput>(
          '/kinesis.KinesisService/ListStreams',
          ($0.ListStreamsInput value) => value.writeToBuffer(),
          $0.ListStreamsOutput.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceInput, $0.ListTagsForResourceOutput>(
      '/kinesis.KinesisService/ListTagsForResource',
      ($0.ListTagsForResourceInput value) => value.writeToBuffer(),
      $0.ListTagsForResourceOutput.fromBuffer);
  static final _$listTagsForStream =
      $grpc.ClientMethod<$0.ListTagsForStreamInput, $0.ListTagsForStreamOutput>(
          '/kinesis.KinesisService/ListTagsForStream',
          ($0.ListTagsForStreamInput value) => value.writeToBuffer(),
          $0.ListTagsForStreamOutput.fromBuffer);
  static final _$mergeShards =
      $grpc.ClientMethod<$0.MergeShardsInput, $1.Empty>(
          '/kinesis.KinesisService/MergeShards',
          ($0.MergeShardsInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putRecord =
      $grpc.ClientMethod<$0.PutRecordInput, $0.PutRecordOutput>(
          '/kinesis.KinesisService/PutRecord',
          ($0.PutRecordInput value) => value.writeToBuffer(),
          $0.PutRecordOutput.fromBuffer);
  static final _$putRecords =
      $grpc.ClientMethod<$0.PutRecordsInput, $0.PutRecordsOutput>(
          '/kinesis.KinesisService/PutRecords',
          ($0.PutRecordsInput value) => value.writeToBuffer(),
          $0.PutRecordsOutput.fromBuffer);
  static final _$putResourcePolicy =
      $grpc.ClientMethod<$0.PutResourcePolicyInput, $1.Empty>(
          '/kinesis.KinesisService/PutResourcePolicy',
          ($0.PutResourcePolicyInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$registerStreamConsumer = $grpc.ClientMethod<
          $0.RegisterStreamConsumerInput, $0.RegisterStreamConsumerOutput>(
      '/kinesis.KinesisService/RegisterStreamConsumer',
      ($0.RegisterStreamConsumerInput value) => value.writeToBuffer(),
      $0.RegisterStreamConsumerOutput.fromBuffer);
  static final _$removeTagsFromStream =
      $grpc.ClientMethod<$0.RemoveTagsFromStreamInput, $1.Empty>(
          '/kinesis.KinesisService/RemoveTagsFromStream',
          ($0.RemoveTagsFromStreamInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$splitShard = $grpc.ClientMethod<$0.SplitShardInput, $1.Empty>(
      '/kinesis.KinesisService/SplitShard',
      ($0.SplitShardInput value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$startStreamEncryption =
      $grpc.ClientMethod<$0.StartStreamEncryptionInput, $1.Empty>(
          '/kinesis.KinesisService/StartStreamEncryption',
          ($0.StartStreamEncryptionInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$stopStreamEncryption =
      $grpc.ClientMethod<$0.StopStreamEncryptionInput, $1.Empty>(
          '/kinesis.KinesisService/StopStreamEncryption',
          ($0.StopStreamEncryptionInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$subscribeToShard =
      $grpc.ClientMethod<$0.SubscribeToShardInput, $0.SubscribeToShardOutput>(
          '/kinesis.KinesisService/SubscribeToShard',
          ($0.SubscribeToShardInput value) => value.writeToBuffer(),
          $0.SubscribeToShardOutput.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceInput, $1.Empty>(
          '/kinesis.KinesisService/TagResource',
          ($0.TagResourceInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceInput, $1.Empty>(
          '/kinesis.KinesisService/UntagResource',
          ($0.UntagResourceInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateAccountSettings = $grpc.ClientMethod<
          $0.UpdateAccountSettingsInput, $0.UpdateAccountSettingsOutput>(
      '/kinesis.KinesisService/UpdateAccountSettings',
      ($0.UpdateAccountSettingsInput value) => value.writeToBuffer(),
      $0.UpdateAccountSettingsOutput.fromBuffer);
  static final _$updateMaxRecordSize =
      $grpc.ClientMethod<$0.UpdateMaxRecordSizeInput, $1.Empty>(
          '/kinesis.KinesisService/UpdateMaxRecordSize',
          ($0.UpdateMaxRecordSizeInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateShardCount =
      $grpc.ClientMethod<$0.UpdateShardCountInput, $0.UpdateShardCountOutput>(
          '/kinesis.KinesisService/UpdateShardCount',
          ($0.UpdateShardCountInput value) => value.writeToBuffer(),
          $0.UpdateShardCountOutput.fromBuffer);
  static final _$updateStreamMode =
      $grpc.ClientMethod<$0.UpdateStreamModeInput, $1.Empty>(
          '/kinesis.KinesisService/UpdateStreamMode',
          ($0.UpdateStreamModeInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateStreamWarmThroughput = $grpc.ClientMethod<
          $0.UpdateStreamWarmThroughputInput,
          $0.UpdateStreamWarmThroughputOutput>(
      '/kinesis.KinesisService/UpdateStreamWarmThroughput',
      ($0.UpdateStreamWarmThroughputInput value) => value.writeToBuffer(),
      $0.UpdateStreamWarmThroughputOutput.fromBuffer);
}

@$pb.GrpcServiceName('kinesis.KinesisService')
abstract class KinesisServiceBase extends $grpc.Service {
  $core.String get $name => 'kinesis.KinesisService';

  KinesisServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AddTagsToStreamInput, $1.Empty>(
        'AddTagsToStream',
        addTagsToStream_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddTagsToStreamInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateStreamInput, $1.Empty>(
        'CreateStream',
        createStream_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateStreamInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DecreaseStreamRetentionPeriodInput, $1.Empty>(
            'DecreaseStreamRetentionPeriod',
            decreaseStreamRetentionPeriod_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DecreaseStreamRetentionPeriodInput.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourcePolicyInput, $1.Empty>(
        'DeleteResourcePolicy',
        deleteResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourcePolicyInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteStreamInput, $1.Empty>(
        'DeleteStream',
        deleteStream_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteStreamInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeregisterStreamConsumerInput, $1.Empty>(
        'DeregisterStreamConsumer',
        deregisterStreamConsumer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeregisterStreamConsumerInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAccountSettingsInput,
            $0.DescribeAccountSettingsOutput>(
        'DescribeAccountSettings',
        describeAccountSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAccountSettingsInput.fromBuffer(value),
        ($0.DescribeAccountSettingsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeLimitsInput, $0.DescribeLimitsOutput>(
            'DescribeLimits',
            describeLimits_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeLimitsInput.fromBuffer(value),
            ($0.DescribeLimitsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeStreamInput, $0.DescribeStreamOutput>(
            'DescribeStream',
            describeStream_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeStreamInput.fromBuffer(value),
            ($0.DescribeStreamOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeStreamConsumerInput,
            $0.DescribeStreamConsumerOutput>(
        'DescribeStreamConsumer',
        describeStreamConsumer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeStreamConsumerInput.fromBuffer(value),
        ($0.DescribeStreamConsumerOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeStreamSummaryInput,
            $0.DescribeStreamSummaryOutput>(
        'DescribeStreamSummary',
        describeStreamSummary_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeStreamSummaryInput.fromBuffer(value),
        ($0.DescribeStreamSummaryOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableEnhancedMonitoringInput,
            $0.EnhancedMonitoringOutput>(
        'DisableEnhancedMonitoring',
        disableEnhancedMonitoring_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableEnhancedMonitoringInput.fromBuffer(value),
        ($0.EnhancedMonitoringOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableEnhancedMonitoringInput,
            $0.EnhancedMonitoringOutput>(
        'EnableEnhancedMonitoring',
        enableEnhancedMonitoring_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableEnhancedMonitoringInput.fromBuffer(value),
        ($0.EnhancedMonitoringOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRecordsInput, $0.GetRecordsOutput>(
        'GetRecords',
        getRecords_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetRecordsInput.fromBuffer(value),
        ($0.GetRecordsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourcePolicyInput,
            $0.GetResourcePolicyOutput>(
        'GetResourcePolicy',
        getResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourcePolicyInput.fromBuffer(value),
        ($0.GetResourcePolicyOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetShardIteratorInput,
            $0.GetShardIteratorOutput>(
        'GetShardIterator',
        getShardIterator_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetShardIteratorInput.fromBuffer(value),
        ($0.GetShardIteratorOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.IncreaseStreamRetentionPeriodInput, $1.Empty>(
            'IncreaseStreamRetentionPeriod',
            increaseStreamRetentionPeriod_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.IncreaseStreamRetentionPeriodInput.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListShardsInput, $0.ListShardsOutput>(
        'ListShards',
        listShards_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListShardsInput.fromBuffer(value),
        ($0.ListShardsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListStreamConsumersInput,
            $0.ListStreamConsumersOutput>(
        'ListStreamConsumers',
        listStreamConsumers_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListStreamConsumersInput.fromBuffer(value),
        ($0.ListStreamConsumersOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListStreamsInput, $0.ListStreamsOutput>(
        'ListStreams',
        listStreams_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListStreamsInput.fromBuffer(value),
        ($0.ListStreamsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceInput,
            $0.ListTagsForResourceOutput>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceInput.fromBuffer(value),
        ($0.ListTagsForResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForStreamInput,
            $0.ListTagsForStreamOutput>(
        'ListTagsForStream',
        listTagsForStream_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForStreamInput.fromBuffer(value),
        ($0.ListTagsForStreamOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.MergeShardsInput, $1.Empty>(
        'MergeShards',
        mergeShards_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.MergeShardsInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutRecordInput, $0.PutRecordOutput>(
        'PutRecord',
        putRecord_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutRecordInput.fromBuffer(value),
        ($0.PutRecordOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutRecordsInput, $0.PutRecordsOutput>(
        'PutRecords',
        putRecords_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutRecordsInput.fromBuffer(value),
        ($0.PutRecordsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutResourcePolicyInput, $1.Empty>(
        'PutResourcePolicy',
        putResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutResourcePolicyInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RegisterStreamConsumerInput,
            $0.RegisterStreamConsumerOutput>(
        'RegisterStreamConsumer',
        registerStreamConsumer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterStreamConsumerInput.fromBuffer(value),
        ($0.RegisterStreamConsumerOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemoveTagsFromStreamInput, $1.Empty>(
        'RemoveTagsFromStream',
        removeTagsFromStream_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemoveTagsFromStreamInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SplitShardInput, $1.Empty>(
        'SplitShard',
        splitShard_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SplitShardInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartStreamEncryptionInput, $1.Empty>(
        'StartStreamEncryption',
        startStreamEncryption_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartStreamEncryptionInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopStreamEncryptionInput, $1.Empty>(
        'StopStreamEncryption',
        stopStreamEncryption_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopStreamEncryptionInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SubscribeToShardInput,
            $0.SubscribeToShardOutput>(
        'SubscribeToShard',
        subscribeToShard_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SubscribeToShardInput.fromBuffer(value),
        ($0.SubscribeToShardOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceInput, $1.Empty>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagResourceInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagResourceInput, $1.Empty>(
        'UntagResource',
        untagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagResourceInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAccountSettingsInput,
            $0.UpdateAccountSettingsOutput>(
        'UpdateAccountSettings',
        updateAccountSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAccountSettingsInput.fromBuffer(value),
        ($0.UpdateAccountSettingsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateMaxRecordSizeInput, $1.Empty>(
        'UpdateMaxRecordSize',
        updateMaxRecordSize_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateMaxRecordSizeInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateShardCountInput,
            $0.UpdateShardCountOutput>(
        'UpdateShardCount',
        updateShardCount_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateShardCountInput.fromBuffer(value),
        ($0.UpdateShardCountOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateStreamModeInput, $1.Empty>(
        'UpdateStreamMode',
        updateStreamMode_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateStreamModeInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateStreamWarmThroughputInput,
            $0.UpdateStreamWarmThroughputOutput>(
        'UpdateStreamWarmThroughput',
        updateStreamWarmThroughput_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateStreamWarmThroughputInput.fromBuffer(value),
        ($0.UpdateStreamWarmThroughputOutput value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> addTagsToStream_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddTagsToStreamInput> $request) async {
    return addTagsToStream($call, await $request);
  }

  $async.Future<$1.Empty> addTagsToStream(
      $grpc.ServiceCall call, $0.AddTagsToStreamInput request);

  $async.Future<$1.Empty> createStream_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateStreamInput> $request) async {
    return createStream($call, await $request);
  }

  $async.Future<$1.Empty> createStream(
      $grpc.ServiceCall call, $0.CreateStreamInput request);

  $async.Future<$1.Empty> decreaseStreamRetentionPeriod_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DecreaseStreamRetentionPeriodInput> $request) async {
    return decreaseStreamRetentionPeriod($call, await $request);
  }

  $async.Future<$1.Empty> decreaseStreamRetentionPeriod(
      $grpc.ServiceCall call, $0.DecreaseStreamRetentionPeriodInput request);

  $async.Future<$1.Empty> deleteResourcePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourcePolicyInput> $request) async {
    return deleteResourcePolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteResourcePolicy(
      $grpc.ServiceCall call, $0.DeleteResourcePolicyInput request);

  $async.Future<$1.Empty> deleteStream_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteStreamInput> $request) async {
    return deleteStream($call, await $request);
  }

  $async.Future<$1.Empty> deleteStream(
      $grpc.ServiceCall call, $0.DeleteStreamInput request);

  $async.Future<$1.Empty> deregisterStreamConsumer_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeregisterStreamConsumerInput> $request) async {
    return deregisterStreamConsumer($call, await $request);
  }

  $async.Future<$1.Empty> deregisterStreamConsumer(
      $grpc.ServiceCall call, $0.DeregisterStreamConsumerInput request);

  $async.Future<$0.DescribeAccountSettingsOutput> describeAccountSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAccountSettingsInput> $request) async {
    return describeAccountSettings($call, await $request);
  }

  $async.Future<$0.DescribeAccountSettingsOutput> describeAccountSettings(
      $grpc.ServiceCall call, $0.DescribeAccountSettingsInput request);

  $async.Future<$0.DescribeLimitsOutput> describeLimits_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeLimitsInput> $request) async {
    return describeLimits($call, await $request);
  }

  $async.Future<$0.DescribeLimitsOutput> describeLimits(
      $grpc.ServiceCall call, $0.DescribeLimitsInput request);

  $async.Future<$0.DescribeStreamOutput> describeStream_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeStreamInput> $request) async {
    return describeStream($call, await $request);
  }

  $async.Future<$0.DescribeStreamOutput> describeStream(
      $grpc.ServiceCall call, $0.DescribeStreamInput request);

  $async.Future<$0.DescribeStreamConsumerOutput> describeStreamConsumer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeStreamConsumerInput> $request) async {
    return describeStreamConsumer($call, await $request);
  }

  $async.Future<$0.DescribeStreamConsumerOutput> describeStreamConsumer(
      $grpc.ServiceCall call, $0.DescribeStreamConsumerInput request);

  $async.Future<$0.DescribeStreamSummaryOutput> describeStreamSummary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeStreamSummaryInput> $request) async {
    return describeStreamSummary($call, await $request);
  }

  $async.Future<$0.DescribeStreamSummaryOutput> describeStreamSummary(
      $grpc.ServiceCall call, $0.DescribeStreamSummaryInput request);

  $async.Future<$0.EnhancedMonitoringOutput> disableEnhancedMonitoring_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DisableEnhancedMonitoringInput> $request) async {
    return disableEnhancedMonitoring($call, await $request);
  }

  $async.Future<$0.EnhancedMonitoringOutput> disableEnhancedMonitoring(
      $grpc.ServiceCall call, $0.DisableEnhancedMonitoringInput request);

  $async.Future<$0.EnhancedMonitoringOutput> enableEnhancedMonitoring_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.EnableEnhancedMonitoringInput> $request) async {
    return enableEnhancedMonitoring($call, await $request);
  }

  $async.Future<$0.EnhancedMonitoringOutput> enableEnhancedMonitoring(
      $grpc.ServiceCall call, $0.EnableEnhancedMonitoringInput request);

  $async.Future<$0.GetRecordsOutput> getRecords_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetRecordsInput> $request) async {
    return getRecords($call, await $request);
  }

  $async.Future<$0.GetRecordsOutput> getRecords(
      $grpc.ServiceCall call, $0.GetRecordsInput request);

  $async.Future<$0.GetResourcePolicyOutput> getResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResourcePolicyInput> $request) async {
    return getResourcePolicy($call, await $request);
  }

  $async.Future<$0.GetResourcePolicyOutput> getResourcePolicy(
      $grpc.ServiceCall call, $0.GetResourcePolicyInput request);

  $async.Future<$0.GetShardIteratorOutput> getShardIterator_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetShardIteratorInput> $request) async {
    return getShardIterator($call, await $request);
  }

  $async.Future<$0.GetShardIteratorOutput> getShardIterator(
      $grpc.ServiceCall call, $0.GetShardIteratorInput request);

  $async.Future<$1.Empty> increaseStreamRetentionPeriod_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.IncreaseStreamRetentionPeriodInput> $request) async {
    return increaseStreamRetentionPeriod($call, await $request);
  }

  $async.Future<$1.Empty> increaseStreamRetentionPeriod(
      $grpc.ServiceCall call, $0.IncreaseStreamRetentionPeriodInput request);

  $async.Future<$0.ListShardsOutput> listShards_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListShardsInput> $request) async {
    return listShards($call, await $request);
  }

  $async.Future<$0.ListShardsOutput> listShards(
      $grpc.ServiceCall call, $0.ListShardsInput request);

  $async.Future<$0.ListStreamConsumersOutput> listStreamConsumers_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListStreamConsumersInput> $request) async {
    return listStreamConsumers($call, await $request);
  }

  $async.Future<$0.ListStreamConsumersOutput> listStreamConsumers(
      $grpc.ServiceCall call, $0.ListStreamConsumersInput request);

  $async.Future<$0.ListStreamsOutput> listStreams_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListStreamsInput> $request) async {
    return listStreams($call, await $request);
  }

  $async.Future<$0.ListStreamsOutput> listStreams(
      $grpc.ServiceCall call, $0.ListStreamsInput request);

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceInput> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceInput request);

  $async.Future<$0.ListTagsForStreamOutput> listTagsForStream_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForStreamInput> $request) async {
    return listTagsForStream($call, await $request);
  }

  $async.Future<$0.ListTagsForStreamOutput> listTagsForStream(
      $grpc.ServiceCall call, $0.ListTagsForStreamInput request);

  $async.Future<$1.Empty> mergeShards_Pre($grpc.ServiceCall $call,
      $async.Future<$0.MergeShardsInput> $request) async {
    return mergeShards($call, await $request);
  }

  $async.Future<$1.Empty> mergeShards(
      $grpc.ServiceCall call, $0.MergeShardsInput request);

  $async.Future<$0.PutRecordOutput> putRecord_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutRecordInput> $request) async {
    return putRecord($call, await $request);
  }

  $async.Future<$0.PutRecordOutput> putRecord(
      $grpc.ServiceCall call, $0.PutRecordInput request);

  $async.Future<$0.PutRecordsOutput> putRecords_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutRecordsInput> $request) async {
    return putRecords($call, await $request);
  }

  $async.Future<$0.PutRecordsOutput> putRecords(
      $grpc.ServiceCall call, $0.PutRecordsInput request);

  $async.Future<$1.Empty> putResourcePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutResourcePolicyInput> $request) async {
    return putResourcePolicy($call, await $request);
  }

  $async.Future<$1.Empty> putResourcePolicy(
      $grpc.ServiceCall call, $0.PutResourcePolicyInput request);

  $async.Future<$0.RegisterStreamConsumerOutput> registerStreamConsumer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RegisterStreamConsumerInput> $request) async {
    return registerStreamConsumer($call, await $request);
  }

  $async.Future<$0.RegisterStreamConsumerOutput> registerStreamConsumer(
      $grpc.ServiceCall call, $0.RegisterStreamConsumerInput request);

  $async.Future<$1.Empty> removeTagsFromStream_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemoveTagsFromStreamInput> $request) async {
    return removeTagsFromStream($call, await $request);
  }

  $async.Future<$1.Empty> removeTagsFromStream(
      $grpc.ServiceCall call, $0.RemoveTagsFromStreamInput request);

  $async.Future<$1.Empty> splitShard_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SplitShardInput> $request) async {
    return splitShard($call, await $request);
  }

  $async.Future<$1.Empty> splitShard(
      $grpc.ServiceCall call, $0.SplitShardInput request);

  $async.Future<$1.Empty> startStreamEncryption_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StartStreamEncryptionInput> $request) async {
    return startStreamEncryption($call, await $request);
  }

  $async.Future<$1.Empty> startStreamEncryption(
      $grpc.ServiceCall call, $0.StartStreamEncryptionInput request);

  $async.Future<$1.Empty> stopStreamEncryption_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StopStreamEncryptionInput> $request) async {
    return stopStreamEncryption($call, await $request);
  }

  $async.Future<$1.Empty> stopStreamEncryption(
      $grpc.ServiceCall call, $0.StopStreamEncryptionInput request);

  $async.Future<$0.SubscribeToShardOutput> subscribeToShard_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SubscribeToShardInput> $request) async {
    return subscribeToShard($call, await $request);
  }

  $async.Future<$0.SubscribeToShardOutput> subscribeToShard(
      $grpc.ServiceCall call, $0.SubscribeToShardInput request);

  $async.Future<$1.Empty> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceInput> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$1.Empty> tagResource(
      $grpc.ServiceCall call, $0.TagResourceInput request);

  $async.Future<$1.Empty> untagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceInput> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$1.Empty> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceInput request);

  $async.Future<$0.UpdateAccountSettingsOutput> updateAccountSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateAccountSettingsInput> $request) async {
    return updateAccountSettings($call, await $request);
  }

  $async.Future<$0.UpdateAccountSettingsOutput> updateAccountSettings(
      $grpc.ServiceCall call, $0.UpdateAccountSettingsInput request);

  $async.Future<$1.Empty> updateMaxRecordSize_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateMaxRecordSizeInput> $request) async {
    return updateMaxRecordSize($call, await $request);
  }

  $async.Future<$1.Empty> updateMaxRecordSize(
      $grpc.ServiceCall call, $0.UpdateMaxRecordSizeInput request);

  $async.Future<$0.UpdateShardCountOutput> updateShardCount_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateShardCountInput> $request) async {
    return updateShardCount($call, await $request);
  }

  $async.Future<$0.UpdateShardCountOutput> updateShardCount(
      $grpc.ServiceCall call, $0.UpdateShardCountInput request);

  $async.Future<$1.Empty> updateStreamMode_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateStreamModeInput> $request) async {
    return updateStreamMode($call, await $request);
  }

  $async.Future<$1.Empty> updateStreamMode(
      $grpc.ServiceCall call, $0.UpdateStreamModeInput request);

  $async.Future<$0.UpdateStreamWarmThroughputOutput>
      updateStreamWarmThroughput_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateStreamWarmThroughputInput> $request) async {
    return updateStreamWarmThroughput($call, await $request);
  }

  $async.Future<$0.UpdateStreamWarmThroughputOutput> updateStreamWarmThroughput(
      $grpc.ServiceCall call, $0.UpdateStreamWarmThroughputInput request);
}
