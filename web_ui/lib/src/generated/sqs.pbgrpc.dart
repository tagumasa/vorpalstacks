// This is a generated file - do not edit.
//
// Generated from sqs.proto.

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
import 'sqs.pb.dart' as $0;

export 'sqs.pb.dart';

/// SQSService provides sqs API operations.
@$pb.GrpcServiceName('sqs.SQSService')
class SQSServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  SQSServiceClient(super.channel, {super.options, super.interceptors});

  /// Adds a permission to a queue for a specific principal. This allows sharing access to the queue. When you create a queue, you have full control access rights for the queue. Only you, the owner of th...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> addPermission(
    $0.AddPermissionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addPermission, request, options: options);
  }

  /// Cancels a specified message movement task. A message movement can only be cancelled when the current status is RUNNING. Cancelling a message movement task does not revert the messages that have alr...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CancelMessageMoveTaskResult> cancelMessageMoveTask(
    $0.CancelMessageMoveTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelMessageMoveTask, request, options: options);
  }

  /// Changes the visibility timeout of a specified message in a queue to a new value. The default visibility timeout for a message is 30 seconds. The minimum is 0 seconds. The maximum is 12 hours. For m...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> changeMessageVisibility(
    $0.ChangeMessageVisibilityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changeMessageVisibility, request,
        options: options);
  }

  /// Changes the visibility timeout of multiple messages. This is a batch version of ChangeMessageVisibility. The result of the action on each message is reported individually in the response. You can s...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ChangeMessageVisibilityBatchResult>
      changeMessageVisibilityBatch(
    $0.ChangeMessageVisibilityBatchRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changeMessageVisibilityBatch, request,
        options: options);
  }

  /// Creates a new standard or FIFO queue. You can pass one or more attributes in the request. Keep the following in mind: If you don't specify the FifoQueue attribute, Amazon SQS creates a standard que...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateQueueResult> createQueue(
    $0.CreateQueueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createQueue, request, options: options);
  }

  /// Deletes the specified message from the specified queue. To select the message to delete, use the ReceiptHandle of the message (not the MessageId which you receive when you send the message). Amazon...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> deleteMessage(
    $0.DeleteMessageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMessage, request, options: options);
  }

  /// Deletes up to ten messages from the specified queue. This is a batch version of DeleteMessage. The result of the action on each message is reported individually in the response. Because the batch r...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteMessageBatchResult> deleteMessageBatch(
    $0.DeleteMessageBatchRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMessageBatch, request, options: options);
  }

  /// Deletes the queue specified by the QueueUrl, regardless of the queue's contents. Be careful with the DeleteQueue action: When you delete a queue, any messages in the queue are no longer available. ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> deleteQueue(
    $0.DeleteQueueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteQueue, request, options: options);
  }

  /// Gets attributes for the specified queue. To determine whether a queue is FIFO, you can check whether QueueName ends with the .fifo suffix.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetQueueAttributesResult> getQueueAttributes(
    $0.GetQueueAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getQueueAttributes, request, options: options);
  }

  /// The GetQueueUrl API returns the URL of an existing Amazon SQS queue. This is useful when you know the queue's name but need to retrieve its URL for further operations. To access a queue owned by an...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetQueueUrlResult> getQueueUrl(
    $0.GetQueueUrlRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getQueueUrl, request, options: options);
  }

  /// Returns a list of your queues that have the RedrivePolicy queue attribute configured with a dead-letter queue. The ListDeadLetterSourceQueues methods supports pagination. Set parameter MaxResults i...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListDeadLetterSourceQueuesResult>
      listDeadLetterSourceQueues(
    $0.ListDeadLetterSourceQueuesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDeadLetterSourceQueues, request,
        options: options);
  }

  /// Gets the most recent message movement tasks (up to 10) under a specific source queue. This action is currently limited to supporting message redrive from dead-letter queues (DLQs) only. In this con...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListMessageMoveTasksResult> listMessageMoveTasks(
    $0.ListMessageMoveTasksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMessageMoveTasks, request, options: options);
  }

  /// Returns a list of your queues in the current region. The response includes a maximum of 1,000 results. If you specify a value for the optional QueueNamePrefix parameter, only queues with a name tha...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListQueuesResult> listQueues(
    $0.ListQueuesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listQueues, request, options: options);
  }

  /// List all cost allocation tags added to the specified Amazon SQS queue. For an overview, see Tagging Your Amazon SQS Queues in the Amazon SQS Developer Guide. Cross-account permissions don't apply t...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListQueueTagsResult> listQueueTags(
    $0.ListQueueTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listQueueTags, request, options: options);
  }

  /// Deletes available messages in a queue (including in-flight messages) specified by the QueueURL parameter. When you use the PurgeQueue action, you can't retrieve any messages deleted from a queue. T...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> purgeQueue(
    $0.PurgeQueueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$purgeQueue, request, options: options);
  }

  /// Retrieves one or more messages (up to 10), from the specified queue. Using the WaitTimeSeconds parameter enables long-poll support. For more information, see Amazon SQS Long Polling in the Amazon S...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ReceiveMessageResult> receiveMessage(
    $0.ReceiveMessageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$receiveMessage, request, options: options);
  }

  /// Revokes any permissions in the queue policy that matches the specified Label parameter. Only the owner of a queue can remove permissions from it. Cross-account permissions don't apply to this actio...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> removePermission(
    $0.RemovePermissionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removePermission, request, options: options);
  }

  /// Delivers a message to the specified queue. A message can include only XML, JSON, and unformatted text. The following Unicode characters are allowed. For more information, see the W3C specification ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.SendMessageResult> sendMessage(
    $0.SendMessageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendMessage, request, options: options);
  }

  /// You can use SendMessageBatch to send up to 10 messages to the specified queue by assigning either identical or different values to each message (or by not assigning values at all). This is a batch ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.SendMessageBatchResult> sendMessageBatch(
    $0.SendMessageBatchRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendMessageBatch, request, options: options);
  }

  /// Sets the value of one or more queue attributes, like a policy. When you change a queue's attributes, the change can take up to 60 seconds for most of the attributes to propagate throughout the Amaz...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> setQueueAttributes(
    $0.SetQueueAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setQueueAttributes, request, options: options);
  }

  /// Starts an asynchronous task to move messages from a specified source queue to a specified destination queue. This action is currently limited to supporting message redrive from queues that are conf...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.StartMessageMoveTaskResult> startMessageMoveTask(
    $0.StartMessageMoveTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startMessageMoveTask, request, options: options);
  }

  /// Add cost allocation tags to the specified Amazon SQS queue. For an overview, see Tagging Your Amazon SQS Queues in the Amazon SQS Developer Guide. When you use queue tags, keep the following guidel...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> tagQueue(
    $0.TagQueueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagQueue, request, options: options);
  }

  /// Remove cost allocation tags from the specified Amazon SQS queue. For an overview, see Tagging Your Amazon SQS Queues in the Amazon SQS Developer Guide. Cross-account permissions don't apply to this...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> untagQueue(
    $0.UntagQueueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagQueue, request, options: options);
  }

  // method descriptors

  static final _$addPermission =
      $grpc.ClientMethod<$0.AddPermissionRequest, $1.Empty>(
          '/sqs.SQSService/AddPermission',
          ($0.AddPermissionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$cancelMessageMoveTask = $grpc.ClientMethod<
          $0.CancelMessageMoveTaskRequest, $0.CancelMessageMoveTaskResult>(
      '/sqs.SQSService/CancelMessageMoveTask',
      ($0.CancelMessageMoveTaskRequest value) => value.writeToBuffer(),
      $0.CancelMessageMoveTaskResult.fromBuffer);
  static final _$changeMessageVisibility =
      $grpc.ClientMethod<$0.ChangeMessageVisibilityRequest, $1.Empty>(
          '/sqs.SQSService/ChangeMessageVisibility',
          ($0.ChangeMessageVisibilityRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$changeMessageVisibilityBatch = $grpc.ClientMethod<
          $0.ChangeMessageVisibilityBatchRequest,
          $0.ChangeMessageVisibilityBatchResult>(
      '/sqs.SQSService/ChangeMessageVisibilityBatch',
      ($0.ChangeMessageVisibilityBatchRequest value) => value.writeToBuffer(),
      $0.ChangeMessageVisibilityBatchResult.fromBuffer);
  static final _$createQueue =
      $grpc.ClientMethod<$0.CreateQueueRequest, $0.CreateQueueResult>(
          '/sqs.SQSService/CreateQueue',
          ($0.CreateQueueRequest value) => value.writeToBuffer(),
          $0.CreateQueueResult.fromBuffer);
  static final _$deleteMessage =
      $grpc.ClientMethod<$0.DeleteMessageRequest, $1.Empty>(
          '/sqs.SQSService/DeleteMessage',
          ($0.DeleteMessageRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteMessageBatch = $grpc.ClientMethod<
          $0.DeleteMessageBatchRequest, $0.DeleteMessageBatchResult>(
      '/sqs.SQSService/DeleteMessageBatch',
      ($0.DeleteMessageBatchRequest value) => value.writeToBuffer(),
      $0.DeleteMessageBatchResult.fromBuffer);
  static final _$deleteQueue =
      $grpc.ClientMethod<$0.DeleteQueueRequest, $1.Empty>(
          '/sqs.SQSService/DeleteQueue',
          ($0.DeleteQueueRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$getQueueAttributes = $grpc.ClientMethod<
          $0.GetQueueAttributesRequest, $0.GetQueueAttributesResult>(
      '/sqs.SQSService/GetQueueAttributes',
      ($0.GetQueueAttributesRequest value) => value.writeToBuffer(),
      $0.GetQueueAttributesResult.fromBuffer);
  static final _$getQueueUrl =
      $grpc.ClientMethod<$0.GetQueueUrlRequest, $0.GetQueueUrlResult>(
          '/sqs.SQSService/GetQueueUrl',
          ($0.GetQueueUrlRequest value) => value.writeToBuffer(),
          $0.GetQueueUrlResult.fromBuffer);
  static final _$listDeadLetterSourceQueues = $grpc.ClientMethod<
          $0.ListDeadLetterSourceQueuesRequest,
          $0.ListDeadLetterSourceQueuesResult>(
      '/sqs.SQSService/ListDeadLetterSourceQueues',
      ($0.ListDeadLetterSourceQueuesRequest value) => value.writeToBuffer(),
      $0.ListDeadLetterSourceQueuesResult.fromBuffer);
  static final _$listMessageMoveTasks = $grpc.ClientMethod<
          $0.ListMessageMoveTasksRequest, $0.ListMessageMoveTasksResult>(
      '/sqs.SQSService/ListMessageMoveTasks',
      ($0.ListMessageMoveTasksRequest value) => value.writeToBuffer(),
      $0.ListMessageMoveTasksResult.fromBuffer);
  static final _$listQueues =
      $grpc.ClientMethod<$0.ListQueuesRequest, $0.ListQueuesResult>(
          '/sqs.SQSService/ListQueues',
          ($0.ListQueuesRequest value) => value.writeToBuffer(),
          $0.ListQueuesResult.fromBuffer);
  static final _$listQueueTags =
      $grpc.ClientMethod<$0.ListQueueTagsRequest, $0.ListQueueTagsResult>(
          '/sqs.SQSService/ListQueueTags',
          ($0.ListQueueTagsRequest value) => value.writeToBuffer(),
          $0.ListQueueTagsResult.fromBuffer);
  static final _$purgeQueue =
      $grpc.ClientMethod<$0.PurgeQueueRequest, $1.Empty>(
          '/sqs.SQSService/PurgeQueue',
          ($0.PurgeQueueRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$receiveMessage =
      $grpc.ClientMethod<$0.ReceiveMessageRequest, $0.ReceiveMessageResult>(
          '/sqs.SQSService/ReceiveMessage',
          ($0.ReceiveMessageRequest value) => value.writeToBuffer(),
          $0.ReceiveMessageResult.fromBuffer);
  static final _$removePermission =
      $grpc.ClientMethod<$0.RemovePermissionRequest, $1.Empty>(
          '/sqs.SQSService/RemovePermission',
          ($0.RemovePermissionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$sendMessage =
      $grpc.ClientMethod<$0.SendMessageRequest, $0.SendMessageResult>(
          '/sqs.SQSService/SendMessage',
          ($0.SendMessageRequest value) => value.writeToBuffer(),
          $0.SendMessageResult.fromBuffer);
  static final _$sendMessageBatch =
      $grpc.ClientMethod<$0.SendMessageBatchRequest, $0.SendMessageBatchResult>(
          '/sqs.SQSService/SendMessageBatch',
          ($0.SendMessageBatchRequest value) => value.writeToBuffer(),
          $0.SendMessageBatchResult.fromBuffer);
  static final _$setQueueAttributes =
      $grpc.ClientMethod<$0.SetQueueAttributesRequest, $1.Empty>(
          '/sqs.SQSService/SetQueueAttributes',
          ($0.SetQueueAttributesRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$startMessageMoveTask = $grpc.ClientMethod<
          $0.StartMessageMoveTaskRequest, $0.StartMessageMoveTaskResult>(
      '/sqs.SQSService/StartMessageMoveTask',
      ($0.StartMessageMoveTaskRequest value) => value.writeToBuffer(),
      $0.StartMessageMoveTaskResult.fromBuffer);
  static final _$tagQueue = $grpc.ClientMethod<$0.TagQueueRequest, $1.Empty>(
      '/sqs.SQSService/TagQueue',
      ($0.TagQueueRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$untagQueue =
      $grpc.ClientMethod<$0.UntagQueueRequest, $1.Empty>(
          '/sqs.SQSService/UntagQueue',
          ($0.UntagQueueRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
}

@$pb.GrpcServiceName('sqs.SQSService')
abstract class SQSServiceBase extends $grpc.Service {
  $core.String get $name => 'sqs.SQSService';

  SQSServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AddPermissionRequest, $1.Empty>(
        'AddPermission',
        addPermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddPermissionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelMessageMoveTaskRequest,
            $0.CancelMessageMoveTaskResult>(
        'CancelMessageMoveTask',
        cancelMessageMoveTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelMessageMoveTaskRequest.fromBuffer(value),
        ($0.CancelMessageMoveTaskResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangeMessageVisibilityRequest, $1.Empty>(
        'ChangeMessageVisibility',
        changeMessageVisibility_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ChangeMessageVisibilityRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangeMessageVisibilityBatchRequest,
            $0.ChangeMessageVisibilityBatchResult>(
        'ChangeMessageVisibilityBatch',
        changeMessageVisibilityBatch_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ChangeMessageVisibilityBatchRequest.fromBuffer(value),
        ($0.ChangeMessageVisibilityBatchResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateQueueRequest, $0.CreateQueueResult>(
        'CreateQueue',
        createQueue_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateQueueRequest.fromBuffer(value),
        ($0.CreateQueueResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMessageRequest, $1.Empty>(
        'DeleteMessage',
        deleteMessage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMessageRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMessageBatchRequest,
            $0.DeleteMessageBatchResult>(
        'DeleteMessageBatch',
        deleteMessageBatch_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMessageBatchRequest.fromBuffer(value),
        ($0.DeleteMessageBatchResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteQueueRequest, $1.Empty>(
        'DeleteQueue',
        deleteQueue_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteQueueRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetQueueAttributesRequest,
            $0.GetQueueAttributesResult>(
        'GetQueueAttributes',
        getQueueAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetQueueAttributesRequest.fromBuffer(value),
        ($0.GetQueueAttributesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetQueueUrlRequest, $0.GetQueueUrlResult>(
        'GetQueueUrl',
        getQueueUrl_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetQueueUrlRequest.fromBuffer(value),
        ($0.GetQueueUrlResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDeadLetterSourceQueuesRequest,
            $0.ListDeadLetterSourceQueuesResult>(
        'ListDeadLetterSourceQueues',
        listDeadLetterSourceQueues_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDeadLetterSourceQueuesRequest.fromBuffer(value),
        ($0.ListDeadLetterSourceQueuesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMessageMoveTasksRequest,
            $0.ListMessageMoveTasksResult>(
        'ListMessageMoveTasks',
        listMessageMoveTasks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListMessageMoveTasksRequest.fromBuffer(value),
        ($0.ListMessageMoveTasksResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListQueuesRequest, $0.ListQueuesResult>(
        'ListQueues',
        listQueues_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListQueuesRequest.fromBuffer(value),
        ($0.ListQueuesResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListQueueTagsRequest, $0.ListQueueTagsResult>(
            'ListQueueTags',
            listQueueTags_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListQueueTagsRequest.fromBuffer(value),
            ($0.ListQueueTagsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PurgeQueueRequest, $1.Empty>(
        'PurgeQueue',
        purgeQueue_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PurgeQueueRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ReceiveMessageRequest, $0.ReceiveMessageResult>(
            'ReceiveMessage',
            receiveMessage_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ReceiveMessageRequest.fromBuffer(value),
            ($0.ReceiveMessageResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemovePermissionRequest, $1.Empty>(
        'RemovePermission',
        removePermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemovePermissionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendMessageRequest, $0.SendMessageResult>(
        'SendMessage',
        sendMessage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendMessageRequest.fromBuffer(value),
        ($0.SendMessageResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendMessageBatchRequest,
            $0.SendMessageBatchResult>(
        'SendMessageBatch',
        sendMessageBatch_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendMessageBatchRequest.fromBuffer(value),
        ($0.SendMessageBatchResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetQueueAttributesRequest, $1.Empty>(
        'SetQueueAttributes',
        setQueueAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetQueueAttributesRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartMessageMoveTaskRequest,
            $0.StartMessageMoveTaskResult>(
        'StartMessageMoveTask',
        startMessageMoveTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartMessageMoveTaskRequest.fromBuffer(value),
        ($0.StartMessageMoveTaskResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagQueueRequest, $1.Empty>(
        'TagQueue',
        tagQueue_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagQueueRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagQueueRequest, $1.Empty>(
        'UntagQueue',
        untagQueue_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UntagQueueRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> addPermission_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddPermissionRequest> $request) async {
    return addPermission($call, await $request);
  }

  $async.Future<$1.Empty> addPermission(
      $grpc.ServiceCall call, $0.AddPermissionRequest request);

  $async.Future<$0.CancelMessageMoveTaskResult> cancelMessageMoveTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelMessageMoveTaskRequest> $request) async {
    return cancelMessageMoveTask($call, await $request);
  }

  $async.Future<$0.CancelMessageMoveTaskResult> cancelMessageMoveTask(
      $grpc.ServiceCall call, $0.CancelMessageMoveTaskRequest request);

  $async.Future<$1.Empty> changeMessageVisibility_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ChangeMessageVisibilityRequest> $request) async {
    return changeMessageVisibility($call, await $request);
  }

  $async.Future<$1.Empty> changeMessageVisibility(
      $grpc.ServiceCall call, $0.ChangeMessageVisibilityRequest request);

  $async.Future<$0.ChangeMessageVisibilityBatchResult>
      changeMessageVisibilityBatch_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ChangeMessageVisibilityBatchRequest>
              $request) async {
    return changeMessageVisibilityBatch($call, await $request);
  }

  $async.Future<$0.ChangeMessageVisibilityBatchResult>
      changeMessageVisibilityBatch($grpc.ServiceCall call,
          $0.ChangeMessageVisibilityBatchRequest request);

  $async.Future<$0.CreateQueueResult> createQueue_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateQueueRequest> $request) async {
    return createQueue($call, await $request);
  }

  $async.Future<$0.CreateQueueResult> createQueue(
      $grpc.ServiceCall call, $0.CreateQueueRequest request);

  $async.Future<$1.Empty> deleteMessage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteMessageRequest> $request) async {
    return deleteMessage($call, await $request);
  }

  $async.Future<$1.Empty> deleteMessage(
      $grpc.ServiceCall call, $0.DeleteMessageRequest request);

  $async.Future<$0.DeleteMessageBatchResult> deleteMessageBatch_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteMessageBatchRequest> $request) async {
    return deleteMessageBatch($call, await $request);
  }

  $async.Future<$0.DeleteMessageBatchResult> deleteMessageBatch(
      $grpc.ServiceCall call, $0.DeleteMessageBatchRequest request);

  $async.Future<$1.Empty> deleteQueue_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteQueueRequest> $request) async {
    return deleteQueue($call, await $request);
  }

  $async.Future<$1.Empty> deleteQueue(
      $grpc.ServiceCall call, $0.DeleteQueueRequest request);

  $async.Future<$0.GetQueueAttributesResult> getQueueAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetQueueAttributesRequest> $request) async {
    return getQueueAttributes($call, await $request);
  }

  $async.Future<$0.GetQueueAttributesResult> getQueueAttributes(
      $grpc.ServiceCall call, $0.GetQueueAttributesRequest request);

  $async.Future<$0.GetQueueUrlResult> getQueueUrl_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetQueueUrlRequest> $request) async {
    return getQueueUrl($call, await $request);
  }

  $async.Future<$0.GetQueueUrlResult> getQueueUrl(
      $grpc.ServiceCall call, $0.GetQueueUrlRequest request);

  $async.Future<$0.ListDeadLetterSourceQueuesResult>
      listDeadLetterSourceQueues_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListDeadLetterSourceQueuesRequest> $request) async {
    return listDeadLetterSourceQueues($call, await $request);
  }

  $async.Future<$0.ListDeadLetterSourceQueuesResult> listDeadLetterSourceQueues(
      $grpc.ServiceCall call, $0.ListDeadLetterSourceQueuesRequest request);

  $async.Future<$0.ListMessageMoveTasksResult> listMessageMoveTasks_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListMessageMoveTasksRequest> $request) async {
    return listMessageMoveTasks($call, await $request);
  }

  $async.Future<$0.ListMessageMoveTasksResult> listMessageMoveTasks(
      $grpc.ServiceCall call, $0.ListMessageMoveTasksRequest request);

  $async.Future<$0.ListQueuesResult> listQueues_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListQueuesRequest> $request) async {
    return listQueues($call, await $request);
  }

  $async.Future<$0.ListQueuesResult> listQueues(
      $grpc.ServiceCall call, $0.ListQueuesRequest request);

  $async.Future<$0.ListQueueTagsResult> listQueueTags_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListQueueTagsRequest> $request) async {
    return listQueueTags($call, await $request);
  }

  $async.Future<$0.ListQueueTagsResult> listQueueTags(
      $grpc.ServiceCall call, $0.ListQueueTagsRequest request);

  $async.Future<$1.Empty> purgeQueue_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PurgeQueueRequest> $request) async {
    return purgeQueue($call, await $request);
  }

  $async.Future<$1.Empty> purgeQueue(
      $grpc.ServiceCall call, $0.PurgeQueueRequest request);

  $async.Future<$0.ReceiveMessageResult> receiveMessage_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ReceiveMessageRequest> $request) async {
    return receiveMessage($call, await $request);
  }

  $async.Future<$0.ReceiveMessageResult> receiveMessage(
      $grpc.ServiceCall call, $0.ReceiveMessageRequest request);

  $async.Future<$1.Empty> removePermission_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemovePermissionRequest> $request) async {
    return removePermission($call, await $request);
  }

  $async.Future<$1.Empty> removePermission(
      $grpc.ServiceCall call, $0.RemovePermissionRequest request);

  $async.Future<$0.SendMessageResult> sendMessage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendMessageRequest> $request) async {
    return sendMessage($call, await $request);
  }

  $async.Future<$0.SendMessageResult> sendMessage(
      $grpc.ServiceCall call, $0.SendMessageRequest request);

  $async.Future<$0.SendMessageBatchResult> sendMessageBatch_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SendMessageBatchRequest> $request) async {
    return sendMessageBatch($call, await $request);
  }

  $async.Future<$0.SendMessageBatchResult> sendMessageBatch(
      $grpc.ServiceCall call, $0.SendMessageBatchRequest request);

  $async.Future<$1.Empty> setQueueAttributes_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetQueueAttributesRequest> $request) async {
    return setQueueAttributes($call, await $request);
  }

  $async.Future<$1.Empty> setQueueAttributes(
      $grpc.ServiceCall call, $0.SetQueueAttributesRequest request);

  $async.Future<$0.StartMessageMoveTaskResult> startMessageMoveTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartMessageMoveTaskRequest> $request) async {
    return startMessageMoveTask($call, await $request);
  }

  $async.Future<$0.StartMessageMoveTaskResult> startMessageMoveTask(
      $grpc.ServiceCall call, $0.StartMessageMoveTaskRequest request);

  $async.Future<$1.Empty> tagQueue_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagQueueRequest> $request) async {
    return tagQueue($call, await $request);
  }

  $async.Future<$1.Empty> tagQueue(
      $grpc.ServiceCall call, $0.TagQueueRequest request);

  $async.Future<$1.Empty> untagQueue_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagQueueRequest> $request) async {
    return untagQueue($call, await $request);
  }

  $async.Future<$1.Empty> untagQueue(
      $grpc.ServiceCall call, $0.UntagQueueRequest request);
}
