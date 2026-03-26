// This is a generated file - do not edit.
//
// Generated from sns.proto.

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
import 'sns.pb.dart' as $0;

export 'sns.pb.dart';

/// SNSService provides sns API operations.
@$pb.GrpcServiceName('sns.SNSService')
class SNSServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  SNSServiceClient(super.channel, {super.options, super.interceptors});

  /// Adds a statement to a topic's access control policy, granting access for the specified Amazon Web Services accounts to the specified actions. To remove the ability to change topic permissions, you ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> addPermission(
    $0.AddPermissionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addPermission, request, options: options);
  }

  /// Accepts a phone number and indicates whether the phone holder has opted out of receiving SMS messages from your Amazon Web Services account. You cannot send SMS messages to a number that is opted o...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CheckIfPhoneNumberIsOptedOutResponse>
      checkIfPhoneNumberIsOptedOut(
    $0.CheckIfPhoneNumberIsOptedOutInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkIfPhoneNumberIsOptedOut, request,
        options: options);
  }

  /// Verifies an endpoint owner's intent to receive messages by validating the token sent to the endpoint by an earlier Subscribe action. If the token is valid, the action creates a new subscription and...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ConfirmSubscriptionResponse> confirmSubscription(
    $0.ConfirmSubscriptionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$confirmSubscription, request, options: options);
  }

  /// Creates a platform application object for one of the supported push notification services, such as APNS and GCM (Firebase Cloud Messaging), to which devices and mobile apps may register. You must s...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreatePlatformApplicationResponse>
      createPlatformApplication(
    $0.CreatePlatformApplicationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPlatformApplication, request,
        options: options);
  }

  /// Creates an endpoint for a device and mobile app on one of the supported push notification services, such as GCM (Firebase Cloud Messaging) and APNS. CreatePlatformEndpoint requires the PlatformAppl...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateEndpointResponse> createPlatformEndpoint(
    $0.CreatePlatformEndpointInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPlatformEndpoint, request,
        options: options);
  }

  /// Adds a destination phone number to an Amazon Web Services account in the SMS sandbox and sends a one-time password (OTP) to that phone number. When you start using Amazon SNS to send SMS messages, ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateSMSSandboxPhoneNumberResult>
      createSMSSandboxPhoneNumber(
    $0.CreateSMSSandboxPhoneNumberInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSMSSandboxPhoneNumber, request,
        options: options);
  }

  /// Creates a topic to which notifications can be published. Users can create at most 100,000 standard topics (at most 1,000 FIFO topics). For more information, see Creating an Amazon SNS topic in the ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateTopicResponse> createTopic(
    $0.CreateTopicInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTopic, request, options: options);
  }

  /// Deletes the endpoint for a device and mobile app from Amazon SNS. This action is idempotent. For more information, see Using Amazon SNS Mobile Push Notifications. When you delete an endpoint that i...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteEndpoint(
    $0.DeleteEndpointInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteEndpoint, request, options: options);
  }

  /// Deletes a platform application object for one of the supported push notification services, such as APNS and GCM (Firebase Cloud Messaging). For more information, see Using Amazon SNS Mobile Push No...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deletePlatformApplication(
    $0.DeletePlatformApplicationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePlatformApplication, request,
        options: options);
  }

  /// Deletes an Amazon Web Services account's verified or pending phone number from the SMS sandbox. When you start using Amazon SNS to send SMS messages, your Amazon Web Services account is in the SMS ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.DeleteSMSSandboxPhoneNumberResult>
      deleteSMSSandboxPhoneNumber(
    $0.DeleteSMSSandboxPhoneNumberInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSMSSandboxPhoneNumber, request,
        options: options);
  }

  /// Deletes a topic and all its subscriptions. Deleting a topic might prevent some messages previously sent to the topic from being delivered to subscribers. This action is idempotent, so deleting a to...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteTopic(
    $0.DeleteTopicInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTopic, request, options: options);
  }

  /// Retrieves the specified inline DataProtectionPolicy document that is stored in the specified Amazon SNS topic.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetDataProtectionPolicyResponse>
      getDataProtectionPolicy(
    $0.GetDataProtectionPolicyInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDataProtectionPolicy, request,
        options: options);
  }

  /// Retrieves the endpoint attributes for a device on one of the supported push notification services, such as GCM (Firebase Cloud Messaging) and APNS. For more information, see Using Amazon SNS Mobile...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetEndpointAttributesResponse> getEndpointAttributes(
    $0.GetEndpointAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEndpointAttributes, request, options: options);
  }

  /// Retrieves the attributes of the platform application object for the supported push notification services, such as APNS and GCM (Firebase Cloud Messaging). For more information, see Using Amazon SNS...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetPlatformApplicationAttributesResponse>
      getPlatformApplicationAttributes(
    $0.GetPlatformApplicationAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPlatformApplicationAttributes, request,
        options: options);
  }

  /// Returns the settings for sending SMS messages from your Amazon Web Services account. These settings are set with the SetSMSAttributes action.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetSMSAttributesResponse> getSMSAttributes(
    $0.GetSMSAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSMSAttributes, request, options: options);
  }

  /// Retrieves the SMS sandbox status for the calling Amazon Web Services account in the target Amazon Web Services Region. When you start using Amazon SNS to send SMS messages, your Amazon Web Services...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetSMSSandboxAccountStatusResult>
      getSMSSandboxAccountStatus(
    $0.GetSMSSandboxAccountStatusInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSMSSandboxAccountStatus, request,
        options: options);
  }

  /// Returns all of the properties of a subscription.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetSubscriptionAttributesResponse>
      getSubscriptionAttributes(
    $0.GetSubscriptionAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSubscriptionAttributes, request,
        options: options);
  }

  /// Returns all of the properties of a topic. Topic properties returned might differ based on the authorization of the user.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetTopicAttributesResponse> getTopicAttributes(
    $0.GetTopicAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTopicAttributes, request, options: options);
  }

  /// Lists the endpoints and endpoint attributes for devices in a supported push notification service, such as GCM (Firebase Cloud Messaging) and APNS. The results for ListEndpointsByPlatformApplication...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListEndpointsByPlatformApplicationResponse>
      listEndpointsByPlatformApplication(
    $0.ListEndpointsByPlatformApplicationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEndpointsByPlatformApplication, request,
        options: options);
  }

  /// Lists the calling Amazon Web Services account's dedicated origination numbers and their metadata. For more information about origination numbers, see Origination numbers in the Amazon SNS Developer...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListOriginationNumbersResult> listOriginationNumbers(
    $0.ListOriginationNumbersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOriginationNumbers, request,
        options: options);
  }

  /// Returns a list of phone numbers that are opted out, meaning you cannot send SMS messages to them. The results for ListPhoneNumbersOptedOut are paginated, and each page returns up to 100 phone numbe...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListPhoneNumbersOptedOutResponse>
      listPhoneNumbersOptedOut(
    $0.ListPhoneNumbersOptedOutInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPhoneNumbersOptedOut, request,
        options: options);
  }

  /// Lists the platform application objects for the supported push notification services, such as APNS and GCM (Firebase Cloud Messaging). The results for ListPlatformApplications are paginated and retu...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListPlatformApplicationsResponse>
      listPlatformApplications(
    $0.ListPlatformApplicationsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPlatformApplications, request,
        options: options);
  }

  /// Lists the calling Amazon Web Services account's current verified and pending destination phone numbers in the SMS sandbox. When you start using Amazon SNS to send SMS messages, your Amazon Web Serv...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListSMSSandboxPhoneNumbersResult>
      listSMSSandboxPhoneNumbers(
    $0.ListSMSSandboxPhoneNumbersInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSMSSandboxPhoneNumbers, request,
        options: options);
  }

  /// Returns a list of the requester's subscriptions. Each call returns a limited list of subscriptions, up to 100. If there are more subscriptions, a NextToken is also returned. Use the NextToken param...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListSubscriptionsResponse> listSubscriptions(
    $0.ListSubscriptionsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSubscriptions, request, options: options);
  }

  /// Returns a list of the subscriptions to a specific topic. Each call returns a limited list of subscriptions, up to 100. If there are more subscriptions, a NextToken is also returned. Use the NextTok...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListSubscriptionsByTopicResponse>
      listSubscriptionsByTopic(
    $0.ListSubscriptionsByTopicInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSubscriptionsByTopic, request,
        options: options);
  }

  /// List all tags added to the specified Amazon SNS topic. For an overview, see Amazon SNS Tags in the Amazon Simple Notification Service Developer Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Returns a list of the requester's topics. Each call returns a limited list of topics, up to 100. If there are more topics, a NextToken is also returned. Use the NextToken parameter in a new ListTop...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListTopicsResponse> listTopics(
    $0.ListTopicsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTopics, request, options: options);
  }

  /// Use this request to opt in a phone number that is opted out, which enables you to resume sending SMS messages to the number. You can opt in a phone number only once every 30 days.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.OptInPhoneNumberResponse> optInPhoneNumber(
    $0.OptInPhoneNumberInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$optInPhoneNumber, request, options: options);
  }

  /// Sends a message to an Amazon SNS topic, a text message (SMS message) directly to a phone number, or a message to a mobile platform endpoint (when you specify the TargetArn). If you send a message t...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.PublishResponse> publish(
    $0.PublishInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$publish, request, options: options);
  }

  /// Publishes up to 10 messages to the specified topic in a single batch. This is a batch version of the Publish API. If you try to send more than 10 messages in a single batch request, you will receiv...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.PublishBatchResponse> publishBatch(
    $0.PublishBatchInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$publishBatch, request, options: options);
  }

  /// Adds or updates an inline policy document that is stored in the specified Amazon SNS topic.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> putDataProtectionPolicy(
    $0.PutDataProtectionPolicyInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDataProtectionPolicy, request,
        options: options);
  }

  /// Removes a statement from a topic's access control policy. To remove the ability to change topic permissions, you must deny permissions to the AddPermission, RemovePermission, and SetTopicAttributes...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> removePermission(
    $0.RemovePermissionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removePermission, request, options: options);
  }

  /// Sets the attributes for an endpoint for a device on one of the supported push notification services, such as GCM (Firebase Cloud Messaging) and APNS. For more information, see Using Amazon SNS Mobi...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> setEndpointAttributes(
    $0.SetEndpointAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setEndpointAttributes, request, options: options);
  }

  /// Sets the attributes of the platform application object for the supported push notification services, such as APNS and GCM (Firebase Cloud Messaging). For more information, see Using Amazon SNS Mobi...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> setPlatformApplicationAttributes(
    $0.SetPlatformApplicationAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setPlatformApplicationAttributes, request,
        options: options);
  }

  /// Use this request to set the default settings for sending SMS messages and receiving daily SMS usage reports. You can override some of these settings for a single message when you use the Publish ac...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.SetSMSAttributesResponse> setSMSAttributes(
    $0.SetSMSAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setSMSAttributes, request, options: options);
  }

  /// Allows a subscription owner to set an attribute of the subscription to a new value.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> setSubscriptionAttributes(
    $0.SetSubscriptionAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setSubscriptionAttributes, request,
        options: options);
  }

  /// Allows a topic owner to set an attribute of the topic to a new value. To remove the ability to change topic permissions, you must deny permissions to the AddPermission, RemovePermission, and SetTop...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> setTopicAttributes(
    $0.SetTopicAttributesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setTopicAttributes, request, options: options);
  }

  /// Subscribes an endpoint to an Amazon SNS topic. If the endpoint type is HTTP/S or email, or if the endpoint and the topic are not in the same Amazon Web Services account, the endpoint owner must run...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.SubscribeResponse> subscribe(
    $0.SubscribeInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$subscribe, request, options: options);
  }

  /// Add tags to the specified Amazon SNS topic. For an overview, see Amazon SNS Tags in the Amazon SNS Developer Guide. When you use topic tags, keep the following guidelines in mind: Adding more than ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Deletes a subscription. If the subscription requires authentication for deletion, only the owner of the subscription or the topic's owner can unsubscribe, and an Amazon Web Services signature is re...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> unsubscribe(
    $0.UnsubscribeInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$unsubscribe, request, options: options);
  }

  /// Remove tags from the specified Amazon SNS topic. For an overview, see Amazon SNS Tags in the Amazon SNS Developer Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Verifies a destination phone number with a one-time password (OTP) for the calling Amazon Web Services account. When you start using Amazon SNS to send SMS messages, your Amazon Web Services accoun...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.VerifySMSSandboxPhoneNumberResult>
      verifySMSSandboxPhoneNumber(
    $0.VerifySMSSandboxPhoneNumberInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$verifySMSSandboxPhoneNumber, request,
        options: options);
  }

  // method descriptors

  static final _$addPermission =
      $grpc.ClientMethod<$0.AddPermissionInput, $1.Empty>(
          '/sns.SNSService/AddPermission',
          ($0.AddPermissionInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$checkIfPhoneNumberIsOptedOut = $grpc.ClientMethod<
          $0.CheckIfPhoneNumberIsOptedOutInput,
          $0.CheckIfPhoneNumberIsOptedOutResponse>(
      '/sns.SNSService/CheckIfPhoneNumberIsOptedOut',
      ($0.CheckIfPhoneNumberIsOptedOutInput value) => value.writeToBuffer(),
      $0.CheckIfPhoneNumberIsOptedOutResponse.fromBuffer);
  static final _$confirmSubscription = $grpc.ClientMethod<
          $0.ConfirmSubscriptionInput, $0.ConfirmSubscriptionResponse>(
      '/sns.SNSService/ConfirmSubscription',
      ($0.ConfirmSubscriptionInput value) => value.writeToBuffer(),
      $0.ConfirmSubscriptionResponse.fromBuffer);
  static final _$createPlatformApplication = $grpc.ClientMethod<
          $0.CreatePlatformApplicationInput,
          $0.CreatePlatformApplicationResponse>(
      '/sns.SNSService/CreatePlatformApplication',
      ($0.CreatePlatformApplicationInput value) => value.writeToBuffer(),
      $0.CreatePlatformApplicationResponse.fromBuffer);
  static final _$createPlatformEndpoint = $grpc.ClientMethod<
          $0.CreatePlatformEndpointInput, $0.CreateEndpointResponse>(
      '/sns.SNSService/CreatePlatformEndpoint',
      ($0.CreatePlatformEndpointInput value) => value.writeToBuffer(),
      $0.CreateEndpointResponse.fromBuffer);
  static final _$createSMSSandboxPhoneNumber = $grpc.ClientMethod<
          $0.CreateSMSSandboxPhoneNumberInput,
          $0.CreateSMSSandboxPhoneNumberResult>(
      '/sns.SNSService/CreateSMSSandboxPhoneNumber',
      ($0.CreateSMSSandboxPhoneNumberInput value) => value.writeToBuffer(),
      $0.CreateSMSSandboxPhoneNumberResult.fromBuffer);
  static final _$createTopic =
      $grpc.ClientMethod<$0.CreateTopicInput, $0.CreateTopicResponse>(
          '/sns.SNSService/CreateTopic',
          ($0.CreateTopicInput value) => value.writeToBuffer(),
          $0.CreateTopicResponse.fromBuffer);
  static final _$deleteEndpoint =
      $grpc.ClientMethod<$0.DeleteEndpointInput, $1.Empty>(
          '/sns.SNSService/DeleteEndpoint',
          ($0.DeleteEndpointInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deletePlatformApplication =
      $grpc.ClientMethod<$0.DeletePlatformApplicationInput, $1.Empty>(
          '/sns.SNSService/DeletePlatformApplication',
          ($0.DeletePlatformApplicationInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteSMSSandboxPhoneNumber = $grpc.ClientMethod<
          $0.DeleteSMSSandboxPhoneNumberInput,
          $0.DeleteSMSSandboxPhoneNumberResult>(
      '/sns.SNSService/DeleteSMSSandboxPhoneNumber',
      ($0.DeleteSMSSandboxPhoneNumberInput value) => value.writeToBuffer(),
      $0.DeleteSMSSandboxPhoneNumberResult.fromBuffer);
  static final _$deleteTopic =
      $grpc.ClientMethod<$0.DeleteTopicInput, $1.Empty>(
          '/sns.SNSService/DeleteTopic',
          ($0.DeleteTopicInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$getDataProtectionPolicy = $grpc.ClientMethod<
          $0.GetDataProtectionPolicyInput, $0.GetDataProtectionPolicyResponse>(
      '/sns.SNSService/GetDataProtectionPolicy',
      ($0.GetDataProtectionPolicyInput value) => value.writeToBuffer(),
      $0.GetDataProtectionPolicyResponse.fromBuffer);
  static final _$getEndpointAttributes = $grpc.ClientMethod<
          $0.GetEndpointAttributesInput, $0.GetEndpointAttributesResponse>(
      '/sns.SNSService/GetEndpointAttributes',
      ($0.GetEndpointAttributesInput value) => value.writeToBuffer(),
      $0.GetEndpointAttributesResponse.fromBuffer);
  static final _$getPlatformApplicationAttributes = $grpc.ClientMethod<
          $0.GetPlatformApplicationAttributesInput,
          $0.GetPlatformApplicationAttributesResponse>(
      '/sns.SNSService/GetPlatformApplicationAttributes',
      ($0.GetPlatformApplicationAttributesInput value) => value.writeToBuffer(),
      $0.GetPlatformApplicationAttributesResponse.fromBuffer);
  static final _$getSMSAttributes =
      $grpc.ClientMethod<$0.GetSMSAttributesInput, $0.GetSMSAttributesResponse>(
          '/sns.SNSService/GetSMSAttributes',
          ($0.GetSMSAttributesInput value) => value.writeToBuffer(),
          $0.GetSMSAttributesResponse.fromBuffer);
  static final _$getSMSSandboxAccountStatus = $grpc.ClientMethod<
          $0.GetSMSSandboxAccountStatusInput,
          $0.GetSMSSandboxAccountStatusResult>(
      '/sns.SNSService/GetSMSSandboxAccountStatus',
      ($0.GetSMSSandboxAccountStatusInput value) => value.writeToBuffer(),
      $0.GetSMSSandboxAccountStatusResult.fromBuffer);
  static final _$getSubscriptionAttributes = $grpc.ClientMethod<
          $0.GetSubscriptionAttributesInput,
          $0.GetSubscriptionAttributesResponse>(
      '/sns.SNSService/GetSubscriptionAttributes',
      ($0.GetSubscriptionAttributesInput value) => value.writeToBuffer(),
      $0.GetSubscriptionAttributesResponse.fromBuffer);
  static final _$getTopicAttributes = $grpc.ClientMethod<
          $0.GetTopicAttributesInput, $0.GetTopicAttributesResponse>(
      '/sns.SNSService/GetTopicAttributes',
      ($0.GetTopicAttributesInput value) => value.writeToBuffer(),
      $0.GetTopicAttributesResponse.fromBuffer);
  static final _$listEndpointsByPlatformApplication = $grpc.ClientMethod<
          $0.ListEndpointsByPlatformApplicationInput,
          $0.ListEndpointsByPlatformApplicationResponse>(
      '/sns.SNSService/ListEndpointsByPlatformApplication',
      ($0.ListEndpointsByPlatformApplicationInput value) =>
          value.writeToBuffer(),
      $0.ListEndpointsByPlatformApplicationResponse.fromBuffer);
  static final _$listOriginationNumbers = $grpc.ClientMethod<
          $0.ListOriginationNumbersRequest, $0.ListOriginationNumbersResult>(
      '/sns.SNSService/ListOriginationNumbers',
      ($0.ListOriginationNumbersRequest value) => value.writeToBuffer(),
      $0.ListOriginationNumbersResult.fromBuffer);
  static final _$listPhoneNumbersOptedOut = $grpc.ClientMethod<
          $0.ListPhoneNumbersOptedOutInput,
          $0.ListPhoneNumbersOptedOutResponse>(
      '/sns.SNSService/ListPhoneNumbersOptedOut',
      ($0.ListPhoneNumbersOptedOutInput value) => value.writeToBuffer(),
      $0.ListPhoneNumbersOptedOutResponse.fromBuffer);
  static final _$listPlatformApplications = $grpc.ClientMethod<
          $0.ListPlatformApplicationsInput,
          $0.ListPlatformApplicationsResponse>(
      '/sns.SNSService/ListPlatformApplications',
      ($0.ListPlatformApplicationsInput value) => value.writeToBuffer(),
      $0.ListPlatformApplicationsResponse.fromBuffer);
  static final _$listSMSSandboxPhoneNumbers = $grpc.ClientMethod<
          $0.ListSMSSandboxPhoneNumbersInput,
          $0.ListSMSSandboxPhoneNumbersResult>(
      '/sns.SNSService/ListSMSSandboxPhoneNumbers',
      ($0.ListSMSSandboxPhoneNumbersInput value) => value.writeToBuffer(),
      $0.ListSMSSandboxPhoneNumbersResult.fromBuffer);
  static final _$listSubscriptions = $grpc.ClientMethod<
          $0.ListSubscriptionsInput, $0.ListSubscriptionsResponse>(
      '/sns.SNSService/ListSubscriptions',
      ($0.ListSubscriptionsInput value) => value.writeToBuffer(),
      $0.ListSubscriptionsResponse.fromBuffer);
  static final _$listSubscriptionsByTopic = $grpc.ClientMethod<
          $0.ListSubscriptionsByTopicInput,
          $0.ListSubscriptionsByTopicResponse>(
      '/sns.SNSService/ListSubscriptionsByTopic',
      ($0.ListSubscriptionsByTopicInput value) => value.writeToBuffer(),
      $0.ListSubscriptionsByTopicResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/sns.SNSService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$listTopics =
      $grpc.ClientMethod<$0.ListTopicsInput, $0.ListTopicsResponse>(
          '/sns.SNSService/ListTopics',
          ($0.ListTopicsInput value) => value.writeToBuffer(),
          $0.ListTopicsResponse.fromBuffer);
  static final _$optInPhoneNumber =
      $grpc.ClientMethod<$0.OptInPhoneNumberInput, $0.OptInPhoneNumberResponse>(
          '/sns.SNSService/OptInPhoneNumber',
          ($0.OptInPhoneNumberInput value) => value.writeToBuffer(),
          $0.OptInPhoneNumberResponse.fromBuffer);
  static final _$publish =
      $grpc.ClientMethod<$0.PublishInput, $0.PublishResponse>(
          '/sns.SNSService/Publish',
          ($0.PublishInput value) => value.writeToBuffer(),
          $0.PublishResponse.fromBuffer);
  static final _$publishBatch =
      $grpc.ClientMethod<$0.PublishBatchInput, $0.PublishBatchResponse>(
          '/sns.SNSService/PublishBatch',
          ($0.PublishBatchInput value) => value.writeToBuffer(),
          $0.PublishBatchResponse.fromBuffer);
  static final _$putDataProtectionPolicy =
      $grpc.ClientMethod<$0.PutDataProtectionPolicyInput, $1.Empty>(
          '/sns.SNSService/PutDataProtectionPolicy',
          ($0.PutDataProtectionPolicyInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$removePermission =
      $grpc.ClientMethod<$0.RemovePermissionInput, $1.Empty>(
          '/sns.SNSService/RemovePermission',
          ($0.RemovePermissionInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$setEndpointAttributes =
      $grpc.ClientMethod<$0.SetEndpointAttributesInput, $1.Empty>(
          '/sns.SNSService/SetEndpointAttributes',
          ($0.SetEndpointAttributesInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$setPlatformApplicationAttributes =
      $grpc.ClientMethod<$0.SetPlatformApplicationAttributesInput, $1.Empty>(
          '/sns.SNSService/SetPlatformApplicationAttributes',
          ($0.SetPlatformApplicationAttributesInput value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$setSMSAttributes =
      $grpc.ClientMethod<$0.SetSMSAttributesInput, $0.SetSMSAttributesResponse>(
          '/sns.SNSService/SetSMSAttributes',
          ($0.SetSMSAttributesInput value) => value.writeToBuffer(),
          $0.SetSMSAttributesResponse.fromBuffer);
  static final _$setSubscriptionAttributes =
      $grpc.ClientMethod<$0.SetSubscriptionAttributesInput, $1.Empty>(
          '/sns.SNSService/SetSubscriptionAttributes',
          ($0.SetSubscriptionAttributesInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$setTopicAttributes =
      $grpc.ClientMethod<$0.SetTopicAttributesInput, $1.Empty>(
          '/sns.SNSService/SetTopicAttributes',
          ($0.SetTopicAttributesInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$subscribe =
      $grpc.ClientMethod<$0.SubscribeInput, $0.SubscribeResponse>(
          '/sns.SNSService/Subscribe',
          ($0.SubscribeInput value) => value.writeToBuffer(),
          $0.SubscribeResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
          '/sns.SNSService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$unsubscribe =
      $grpc.ClientMethod<$0.UnsubscribeInput, $1.Empty>(
          '/sns.SNSService/Unsubscribe',
          ($0.UnsubscribeInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
          '/sns.SNSService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$verifySMSSandboxPhoneNumber = $grpc.ClientMethod<
          $0.VerifySMSSandboxPhoneNumberInput,
          $0.VerifySMSSandboxPhoneNumberResult>(
      '/sns.SNSService/VerifySMSSandboxPhoneNumber',
      ($0.VerifySMSSandboxPhoneNumberInput value) => value.writeToBuffer(),
      $0.VerifySMSSandboxPhoneNumberResult.fromBuffer);
}

@$pb.GrpcServiceName('sns.SNSService')
abstract class SNSServiceBase extends $grpc.Service {
  $core.String get $name => 'sns.SNSService';

  SNSServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AddPermissionInput, $1.Empty>(
        'AddPermission',
        addPermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddPermissionInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CheckIfPhoneNumberIsOptedOutInput,
            $0.CheckIfPhoneNumberIsOptedOutResponse>(
        'CheckIfPhoneNumberIsOptedOut',
        checkIfPhoneNumberIsOptedOut_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CheckIfPhoneNumberIsOptedOutInput.fromBuffer(value),
        ($0.CheckIfPhoneNumberIsOptedOutResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ConfirmSubscriptionInput,
            $0.ConfirmSubscriptionResponse>(
        'ConfirmSubscription',
        confirmSubscription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ConfirmSubscriptionInput.fromBuffer(value),
        ($0.ConfirmSubscriptionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePlatformApplicationInput,
            $0.CreatePlatformApplicationResponse>(
        'CreatePlatformApplication',
        createPlatformApplication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePlatformApplicationInput.fromBuffer(value),
        ($0.CreatePlatformApplicationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePlatformEndpointInput,
            $0.CreateEndpointResponse>(
        'CreatePlatformEndpoint',
        createPlatformEndpoint_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePlatformEndpointInput.fromBuffer(value),
        ($0.CreateEndpointResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateSMSSandboxPhoneNumberInput,
            $0.CreateSMSSandboxPhoneNumberResult>(
        'CreateSMSSandboxPhoneNumber',
        createSMSSandboxPhoneNumber_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateSMSSandboxPhoneNumberInput.fromBuffer(value),
        ($0.CreateSMSSandboxPhoneNumberResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateTopicInput, $0.CreateTopicResponse>(
        'CreateTopic',
        createTopic_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateTopicInput.fromBuffer(value),
        ($0.CreateTopicResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteEndpointInput, $1.Empty>(
        'DeleteEndpoint',
        deleteEndpoint_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteEndpointInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePlatformApplicationInput, $1.Empty>(
        'DeletePlatformApplication',
        deletePlatformApplication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeletePlatformApplicationInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteSMSSandboxPhoneNumberInput,
            $0.DeleteSMSSandboxPhoneNumberResult>(
        'DeleteSMSSandboxPhoneNumber',
        deleteSMSSandboxPhoneNumber_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteSMSSandboxPhoneNumberInput.fromBuffer(value),
        ($0.DeleteSMSSandboxPhoneNumberResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTopicInput, $1.Empty>(
        'DeleteTopic',
        deleteTopic_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteTopicInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDataProtectionPolicyInput,
            $0.GetDataProtectionPolicyResponse>(
        'GetDataProtectionPolicy',
        getDataProtectionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDataProtectionPolicyInput.fromBuffer(value),
        ($0.GetDataProtectionPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEndpointAttributesInput,
            $0.GetEndpointAttributesResponse>(
        'GetEndpointAttributes',
        getEndpointAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEndpointAttributesInput.fromBuffer(value),
        ($0.GetEndpointAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPlatformApplicationAttributesInput,
            $0.GetPlatformApplicationAttributesResponse>(
        'GetPlatformApplicationAttributes',
        getPlatformApplicationAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPlatformApplicationAttributesInput.fromBuffer(value),
        ($0.GetPlatformApplicationAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSMSAttributesInput,
            $0.GetSMSAttributesResponse>(
        'GetSMSAttributes',
        getSMSAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSMSAttributesInput.fromBuffer(value),
        ($0.GetSMSAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSMSSandboxAccountStatusInput,
            $0.GetSMSSandboxAccountStatusResult>(
        'GetSMSSandboxAccountStatus',
        getSMSSandboxAccountStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSMSSandboxAccountStatusInput.fromBuffer(value),
        ($0.GetSMSSandboxAccountStatusResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSubscriptionAttributesInput,
            $0.GetSubscriptionAttributesResponse>(
        'GetSubscriptionAttributes',
        getSubscriptionAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSubscriptionAttributesInput.fromBuffer(value),
        ($0.GetSubscriptionAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTopicAttributesInput,
            $0.GetTopicAttributesResponse>(
        'GetTopicAttributes',
        getTopicAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTopicAttributesInput.fromBuffer(value),
        ($0.GetTopicAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEndpointsByPlatformApplicationInput,
            $0.ListEndpointsByPlatformApplicationResponse>(
        'ListEndpointsByPlatformApplication',
        listEndpointsByPlatformApplication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEndpointsByPlatformApplicationInput.fromBuffer(value),
        ($0.ListEndpointsByPlatformApplicationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOriginationNumbersRequest,
            $0.ListOriginationNumbersResult>(
        'ListOriginationNumbers',
        listOriginationNumbers_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOriginationNumbersRequest.fromBuffer(value),
        ($0.ListOriginationNumbersResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPhoneNumbersOptedOutInput,
            $0.ListPhoneNumbersOptedOutResponse>(
        'ListPhoneNumbersOptedOut',
        listPhoneNumbersOptedOut_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPhoneNumbersOptedOutInput.fromBuffer(value),
        ($0.ListPhoneNumbersOptedOutResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPlatformApplicationsInput,
            $0.ListPlatformApplicationsResponse>(
        'ListPlatformApplications',
        listPlatformApplications_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPlatformApplicationsInput.fromBuffer(value),
        ($0.ListPlatformApplicationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSMSSandboxPhoneNumbersInput,
            $0.ListSMSSandboxPhoneNumbersResult>(
        'ListSMSSandboxPhoneNumbers',
        listSMSSandboxPhoneNumbers_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSMSSandboxPhoneNumbersInput.fromBuffer(value),
        ($0.ListSMSSandboxPhoneNumbersResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSubscriptionsInput,
            $0.ListSubscriptionsResponse>(
        'ListSubscriptions',
        listSubscriptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSubscriptionsInput.fromBuffer(value),
        ($0.ListSubscriptionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSubscriptionsByTopicInput,
            $0.ListSubscriptionsByTopicResponse>(
        'ListSubscriptionsByTopic',
        listSubscriptionsByTopic_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSubscriptionsByTopicInput.fromBuffer(value),
        ($0.ListSubscriptionsByTopicResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTopicsInput, $0.ListTopicsResponse>(
        'ListTopics',
        listTopics_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListTopicsInput.fromBuffer(value),
        ($0.ListTopicsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.OptInPhoneNumberInput,
            $0.OptInPhoneNumberResponse>(
        'OptInPhoneNumber',
        optInPhoneNumber_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.OptInPhoneNumberInput.fromBuffer(value),
        ($0.OptInPhoneNumberResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PublishInput, $0.PublishResponse>(
        'Publish',
        publish_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PublishInput.fromBuffer(value),
        ($0.PublishResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PublishBatchInput, $0.PublishBatchResponse>(
            'PublishBatch',
            publishBatch_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PublishBatchInput.fromBuffer(value),
            ($0.PublishBatchResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDataProtectionPolicyInput, $1.Empty>(
        'PutDataProtectionPolicy',
        putDataProtectionPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDataProtectionPolicyInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemovePermissionInput, $1.Empty>(
        'RemovePermission',
        removePermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemovePermissionInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetEndpointAttributesInput, $1.Empty>(
        'SetEndpointAttributes',
        setEndpointAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetEndpointAttributesInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.SetPlatformApplicationAttributesInput, $1.Empty>(
            'SetPlatformApplicationAttributes',
            setPlatformApplicationAttributes_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SetPlatformApplicationAttributesInput.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetSMSAttributesInput,
            $0.SetSMSAttributesResponse>(
        'SetSMSAttributes',
        setSMSAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetSMSAttributesInput.fromBuffer(value),
        ($0.SetSMSAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetSubscriptionAttributesInput, $1.Empty>(
        'SetSubscriptionAttributes',
        setSubscriptionAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetSubscriptionAttributesInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetTopicAttributesInput, $1.Empty>(
        'SetTopicAttributes',
        setTopicAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetTopicAttributesInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SubscribeInput, $0.SubscribeResponse>(
        'Subscribe',
        subscribe_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SubscribeInput.fromBuffer(value),
        ($0.SubscribeResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
            'TagResource',
            tagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TagResourceRequest.fromBuffer(value),
            ($0.TagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UnsubscribeInput, $1.Empty>(
        'Unsubscribe',
        unsubscribe_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UnsubscribeInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceRequest.fromBuffer(value),
            ($0.UntagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.VerifySMSSandboxPhoneNumberInput,
            $0.VerifySMSSandboxPhoneNumberResult>(
        'VerifySMSSandboxPhoneNumber',
        verifySMSSandboxPhoneNumber_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.VerifySMSSandboxPhoneNumberInput.fromBuffer(value),
        ($0.VerifySMSSandboxPhoneNumberResult value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> addPermission_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddPermissionInput> $request) async {
    return addPermission($call, await $request);
  }

  $async.Future<$1.Empty> addPermission(
      $grpc.ServiceCall call, $0.AddPermissionInput request);

  $async.Future<$0.CheckIfPhoneNumberIsOptedOutResponse>
      checkIfPhoneNumberIsOptedOut_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CheckIfPhoneNumberIsOptedOutInput> $request) async {
    return checkIfPhoneNumberIsOptedOut($call, await $request);
  }

  $async.Future<$0.CheckIfPhoneNumberIsOptedOutResponse>
      checkIfPhoneNumberIsOptedOut(
          $grpc.ServiceCall call, $0.CheckIfPhoneNumberIsOptedOutInput request);

  $async.Future<$0.ConfirmSubscriptionResponse> confirmSubscription_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ConfirmSubscriptionInput> $request) async {
    return confirmSubscription($call, await $request);
  }

  $async.Future<$0.ConfirmSubscriptionResponse> confirmSubscription(
      $grpc.ServiceCall call, $0.ConfirmSubscriptionInput request);

  $async.Future<$0.CreatePlatformApplicationResponse>
      createPlatformApplication_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreatePlatformApplicationInput> $request) async {
    return createPlatformApplication($call, await $request);
  }

  $async.Future<$0.CreatePlatformApplicationResponse> createPlatformApplication(
      $grpc.ServiceCall call, $0.CreatePlatformApplicationInput request);

  $async.Future<$0.CreateEndpointResponse> createPlatformEndpoint_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreatePlatformEndpointInput> $request) async {
    return createPlatformEndpoint($call, await $request);
  }

  $async.Future<$0.CreateEndpointResponse> createPlatformEndpoint(
      $grpc.ServiceCall call, $0.CreatePlatformEndpointInput request);

  $async.Future<$0.CreateSMSSandboxPhoneNumberResult>
      createSMSSandboxPhoneNumber_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateSMSSandboxPhoneNumberInput> $request) async {
    return createSMSSandboxPhoneNumber($call, await $request);
  }

  $async.Future<$0.CreateSMSSandboxPhoneNumberResult>
      createSMSSandboxPhoneNumber(
          $grpc.ServiceCall call, $0.CreateSMSSandboxPhoneNumberInput request);

  $async.Future<$0.CreateTopicResponse> createTopic_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateTopicInput> $request) async {
    return createTopic($call, await $request);
  }

  $async.Future<$0.CreateTopicResponse> createTopic(
      $grpc.ServiceCall call, $0.CreateTopicInput request);

  $async.Future<$1.Empty> deleteEndpoint_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteEndpointInput> $request) async {
    return deleteEndpoint($call, await $request);
  }

  $async.Future<$1.Empty> deleteEndpoint(
      $grpc.ServiceCall call, $0.DeleteEndpointInput request);

  $async.Future<$1.Empty> deletePlatformApplication_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeletePlatformApplicationInput> $request) async {
    return deletePlatformApplication($call, await $request);
  }

  $async.Future<$1.Empty> deletePlatformApplication(
      $grpc.ServiceCall call, $0.DeletePlatformApplicationInput request);

  $async.Future<$0.DeleteSMSSandboxPhoneNumberResult>
      deleteSMSSandboxPhoneNumber_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteSMSSandboxPhoneNumberInput> $request) async {
    return deleteSMSSandboxPhoneNumber($call, await $request);
  }

  $async.Future<$0.DeleteSMSSandboxPhoneNumberResult>
      deleteSMSSandboxPhoneNumber(
          $grpc.ServiceCall call, $0.DeleteSMSSandboxPhoneNumberInput request);

  $async.Future<$1.Empty> deleteTopic_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteTopicInput> $request) async {
    return deleteTopic($call, await $request);
  }

  $async.Future<$1.Empty> deleteTopic(
      $grpc.ServiceCall call, $0.DeleteTopicInput request);

  $async.Future<$0.GetDataProtectionPolicyResponse> getDataProtectionPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDataProtectionPolicyInput> $request) async {
    return getDataProtectionPolicy($call, await $request);
  }

  $async.Future<$0.GetDataProtectionPolicyResponse> getDataProtectionPolicy(
      $grpc.ServiceCall call, $0.GetDataProtectionPolicyInput request);

  $async.Future<$0.GetEndpointAttributesResponse> getEndpointAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetEndpointAttributesInput> $request) async {
    return getEndpointAttributes($call, await $request);
  }

  $async.Future<$0.GetEndpointAttributesResponse> getEndpointAttributes(
      $grpc.ServiceCall call, $0.GetEndpointAttributesInput request);

  $async.Future<$0.GetPlatformApplicationAttributesResponse>
      getPlatformApplicationAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetPlatformApplicationAttributesInput>
              $request) async {
    return getPlatformApplicationAttributes($call, await $request);
  }

  $async.Future<$0.GetPlatformApplicationAttributesResponse>
      getPlatformApplicationAttributes($grpc.ServiceCall call,
          $0.GetPlatformApplicationAttributesInput request);

  $async.Future<$0.GetSMSAttributesResponse> getSMSAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSMSAttributesInput> $request) async {
    return getSMSAttributes($call, await $request);
  }

  $async.Future<$0.GetSMSAttributesResponse> getSMSAttributes(
      $grpc.ServiceCall call, $0.GetSMSAttributesInput request);

  $async.Future<$0.GetSMSSandboxAccountStatusResult>
      getSMSSandboxAccountStatus_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetSMSSandboxAccountStatusInput> $request) async {
    return getSMSSandboxAccountStatus($call, await $request);
  }

  $async.Future<$0.GetSMSSandboxAccountStatusResult> getSMSSandboxAccountStatus(
      $grpc.ServiceCall call, $0.GetSMSSandboxAccountStatusInput request);

  $async.Future<$0.GetSubscriptionAttributesResponse>
      getSubscriptionAttributes_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetSubscriptionAttributesInput> $request) async {
    return getSubscriptionAttributes($call, await $request);
  }

  $async.Future<$0.GetSubscriptionAttributesResponse> getSubscriptionAttributes(
      $grpc.ServiceCall call, $0.GetSubscriptionAttributesInput request);

  $async.Future<$0.GetTopicAttributesResponse> getTopicAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetTopicAttributesInput> $request) async {
    return getTopicAttributes($call, await $request);
  }

  $async.Future<$0.GetTopicAttributesResponse> getTopicAttributes(
      $grpc.ServiceCall call, $0.GetTopicAttributesInput request);

  $async.Future<$0.ListEndpointsByPlatformApplicationResponse>
      listEndpointsByPlatformApplication_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListEndpointsByPlatformApplicationInput>
              $request) async {
    return listEndpointsByPlatformApplication($call, await $request);
  }

  $async.Future<$0.ListEndpointsByPlatformApplicationResponse>
      listEndpointsByPlatformApplication($grpc.ServiceCall call,
          $0.ListEndpointsByPlatformApplicationInput request);

  $async.Future<$0.ListOriginationNumbersResult> listOriginationNumbers_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListOriginationNumbersRequest> $request) async {
    return listOriginationNumbers($call, await $request);
  }

  $async.Future<$0.ListOriginationNumbersResult> listOriginationNumbers(
      $grpc.ServiceCall call, $0.ListOriginationNumbersRequest request);

  $async.Future<$0.ListPhoneNumbersOptedOutResponse>
      listPhoneNumbersOptedOut_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListPhoneNumbersOptedOutInput> $request) async {
    return listPhoneNumbersOptedOut($call, await $request);
  }

  $async.Future<$0.ListPhoneNumbersOptedOutResponse> listPhoneNumbersOptedOut(
      $grpc.ServiceCall call, $0.ListPhoneNumbersOptedOutInput request);

  $async.Future<$0.ListPlatformApplicationsResponse>
      listPlatformApplications_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListPlatformApplicationsInput> $request) async {
    return listPlatformApplications($call, await $request);
  }

  $async.Future<$0.ListPlatformApplicationsResponse> listPlatformApplications(
      $grpc.ServiceCall call, $0.ListPlatformApplicationsInput request);

  $async.Future<$0.ListSMSSandboxPhoneNumbersResult>
      listSMSSandboxPhoneNumbers_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListSMSSandboxPhoneNumbersInput> $request) async {
    return listSMSSandboxPhoneNumbers($call, await $request);
  }

  $async.Future<$0.ListSMSSandboxPhoneNumbersResult> listSMSSandboxPhoneNumbers(
      $grpc.ServiceCall call, $0.ListSMSSandboxPhoneNumbersInput request);

  $async.Future<$0.ListSubscriptionsResponse> listSubscriptions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSubscriptionsInput> $request) async {
    return listSubscriptions($call, await $request);
  }

  $async.Future<$0.ListSubscriptionsResponse> listSubscriptions(
      $grpc.ServiceCall call, $0.ListSubscriptionsInput request);

  $async.Future<$0.ListSubscriptionsByTopicResponse>
      listSubscriptionsByTopic_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListSubscriptionsByTopicInput> $request) async {
    return listSubscriptionsByTopic($call, await $request);
  }

  $async.Future<$0.ListSubscriptionsByTopicResponse> listSubscriptionsByTopic(
      $grpc.ServiceCall call, $0.ListSubscriptionsByTopicInput request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ListTopicsResponse> listTopics_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListTopicsInput> $request) async {
    return listTopics($call, await $request);
  }

  $async.Future<$0.ListTopicsResponse> listTopics(
      $grpc.ServiceCall call, $0.ListTopicsInput request);

  $async.Future<$0.OptInPhoneNumberResponse> optInPhoneNumber_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.OptInPhoneNumberInput> $request) async {
    return optInPhoneNumber($call, await $request);
  }

  $async.Future<$0.OptInPhoneNumberResponse> optInPhoneNumber(
      $grpc.ServiceCall call, $0.OptInPhoneNumberInput request);

  $async.Future<$0.PublishResponse> publish_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.PublishInput> $request) async {
    return publish($call, await $request);
  }

  $async.Future<$0.PublishResponse> publish(
      $grpc.ServiceCall call, $0.PublishInput request);

  $async.Future<$0.PublishBatchResponse> publishBatch_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PublishBatchInput> $request) async {
    return publishBatch($call, await $request);
  }

  $async.Future<$0.PublishBatchResponse> publishBatch(
      $grpc.ServiceCall call, $0.PublishBatchInput request);

  $async.Future<$1.Empty> putDataProtectionPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutDataProtectionPolicyInput> $request) async {
    return putDataProtectionPolicy($call, await $request);
  }

  $async.Future<$1.Empty> putDataProtectionPolicy(
      $grpc.ServiceCall call, $0.PutDataProtectionPolicyInput request);

  $async.Future<$1.Empty> removePermission_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemovePermissionInput> $request) async {
    return removePermission($call, await $request);
  }

  $async.Future<$1.Empty> removePermission(
      $grpc.ServiceCall call, $0.RemovePermissionInput request);

  $async.Future<$1.Empty> setEndpointAttributes_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetEndpointAttributesInput> $request) async {
    return setEndpointAttributes($call, await $request);
  }

  $async.Future<$1.Empty> setEndpointAttributes(
      $grpc.ServiceCall call, $0.SetEndpointAttributesInput request);

  $async.Future<$1.Empty> setPlatformApplicationAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetPlatformApplicationAttributesInput> $request) async {
    return setPlatformApplicationAttributes($call, await $request);
  }

  $async.Future<$1.Empty> setPlatformApplicationAttributes(
      $grpc.ServiceCall call, $0.SetPlatformApplicationAttributesInput request);

  $async.Future<$0.SetSMSAttributesResponse> setSMSAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetSMSAttributesInput> $request) async {
    return setSMSAttributes($call, await $request);
  }

  $async.Future<$0.SetSMSAttributesResponse> setSMSAttributes(
      $grpc.ServiceCall call, $0.SetSMSAttributesInput request);

  $async.Future<$1.Empty> setSubscriptionAttributes_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetSubscriptionAttributesInput> $request) async {
    return setSubscriptionAttributes($call, await $request);
  }

  $async.Future<$1.Empty> setSubscriptionAttributes(
      $grpc.ServiceCall call, $0.SetSubscriptionAttributesInput request);

  $async.Future<$1.Empty> setTopicAttributes_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetTopicAttributesInput> $request) async {
    return setTopicAttributes($call, await $request);
  }

  $async.Future<$1.Empty> setTopicAttributes(
      $grpc.ServiceCall call, $0.SetTopicAttributesInput request);

  $async.Future<$0.SubscribeResponse> subscribe_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SubscribeInput> $request) async {
    return subscribe($call, await $request);
  }

  $async.Future<$0.SubscribeResponse> subscribe(
      $grpc.ServiceCall call, $0.SubscribeInput request);

  $async.Future<$0.TagResourceResponse> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceResponse> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$1.Empty> unsubscribe_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UnsubscribeInput> $request) async {
    return unsubscribe($call, await $request);
  }

  $async.Future<$1.Empty> unsubscribe(
      $grpc.ServiceCall call, $0.UnsubscribeInput request);

  $async.Future<$0.UntagResourceResponse> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceResponse> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.VerifySMSSandboxPhoneNumberResult>
      verifySMSSandboxPhoneNumber_Pre($grpc.ServiceCall $call,
          $async.Future<$0.VerifySMSSandboxPhoneNumberInput> $request) async {
    return verifySMSSandboxPhoneNumber($call, await $request);
  }

  $async.Future<$0.VerifySMSSandboxPhoneNumberResult>
      verifySMSSandboxPhoneNumber(
          $grpc.ServiceCall call, $0.VerifySMSSandboxPhoneNumberInput request);
}
