// This is a generated file - do not edit.
//
// Generated from cloudfront.proto.

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

import 'cloudfront.pb.dart' as $0;
import 'common.pb.dart' as $1;

export 'cloudfront.pb.dart';

/// CloudFrontService provides cloudfront API operations.
@$pb.GrpcServiceName('cloudfront.CloudFrontService')
class CloudFrontServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CloudFrontServiceClient(super.channel, {super.options, super.interceptors});

  /// The AssociateAlias API operation only supports standard distributions. To move domains between distribution tenants and/or standard distributions, we recommend that you use the UpdateDomainAssociat...
  /// HTTP: PUT /2020-05-31/distribution/{TargetDistributionId}/associate-alias
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> associateAlias(
    $0.AssociateAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateAlias, request, options: options);
  }

  /// Associates the WAF web ACL with a distribution tenant.
  /// HTTP: PUT /2020-05-31/distribution-tenant/{Id}/associate-web-acl
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.AssociateDistributionTenantWebACLResult>
      associateDistributionTenantWebACL(
    $0.AssociateDistributionTenantWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateDistributionTenantWebACL, request,
        options: options);
  }

  /// Associates the WAF web ACL with a distribution.
  /// HTTP: PUT /2020-05-31/distribution/{Id}/associate-web-acl
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.AssociateDistributionWebACLResult>
      associateDistributionWebACL(
    $0.AssociateDistributionWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateDistributionWebACL, request,
        options: options);
  }

  /// Creates a staging distribution using the configuration of the provided primary distribution. A staging distribution is a copy of an existing distribution (called the primary distribution) that you ...
  /// HTTP: POST /2020-05-31/distribution/{PrimaryDistributionId}/copy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CopyDistributionResult> copyDistribution(
    $0.CopyDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$copyDistribution, request, options: options);
  }

  /// Creates an Anycast static IP list.
  /// HTTP: POST /2020-05-31/anycast-ip-list
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateAnycastIpListResult> createAnycastIpList(
    $0.CreateAnycastIpListRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAnycastIpList, request, options: options);
  }

  /// Creates a cache policy. After you create a cache policy, you can attach it to one or more cache behaviors. When it's attached to a cache behavior, the cache policy determines the following: The val...
  /// HTTP: POST /2020-05-31/cache-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateCachePolicyResult> createCachePolicy(
    $0.CreateCachePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createCachePolicy, request, options: options);
  }

  /// Creates a new origin access identity. If you're using Amazon S3 for your origin, you can use an origin access identity to require users to access your content using a CloudFront URL instead of the ...
  /// HTTP: POST /2020-05-31/origin-access-identity/cloudfront
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateCloudFrontOriginAccessIdentityResult>
      createCloudFrontOriginAccessIdentity(
    $0.CreateCloudFrontOriginAccessIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createCloudFrontOriginAccessIdentity, request,
        options: options);
  }

  /// Creates a connection function.
  /// HTTP: POST /2020-05-31/connection-function
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateConnectionFunctionResult>
      createConnectionFunction(
    $0.CreateConnectionFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createConnectionFunction, request,
        options: options);
  }

  /// Creates a connection group.
  /// HTTP: POST /2020-05-31/connection-group
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateConnectionGroupResult> createConnectionGroup(
    $0.CreateConnectionGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createConnectionGroup, request, options: options);
  }

  /// Creates a continuous deployment policy that distributes traffic for a custom domain name to two different CloudFront distributions. To use a continuous deployment policy, first use CopyDistribution...
  /// HTTP: POST /2020-05-31/continuous-deployment-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateContinuousDeploymentPolicyResult>
      createContinuousDeploymentPolicy(
    $0.CreateContinuousDeploymentPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createContinuousDeploymentPolicy, request,
        options: options);
  }

  /// Creates a CloudFront distribution.
  /// HTTP: POST /2020-05-31/distribution
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateDistributionResult> createDistribution(
    $0.CreateDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDistribution, request, options: options);
  }

  /// Creates a distribution tenant.
  /// HTTP: POST /2020-05-31/distribution-tenant
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateDistributionTenantResult>
      createDistributionTenant(
    $0.CreateDistributionTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDistributionTenant, request,
        options: options);
  }

  /// Create a new distribution with tags. This API operation requires the following IAM permissions: CreateDistribution TagResource
  /// HTTP: POST /2020-05-31/distribution?WithTags
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateDistributionWithTagsResult>
      createDistributionWithTags(
    $0.CreateDistributionWithTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDistributionWithTags, request,
        options: options);
  }

  /// Create a new field-level encryption configuration.
  /// HTTP: POST /2020-05-31/field-level-encryption
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateFieldLevelEncryptionConfigResult>
      createFieldLevelEncryptionConfig(
    $0.CreateFieldLevelEncryptionConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createFieldLevelEncryptionConfig, request,
        options: options);
  }

  /// Create a field-level encryption profile.
  /// HTTP: POST /2020-05-31/field-level-encryption-profile
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateFieldLevelEncryptionProfileResult>
      createFieldLevelEncryptionProfile(
    $0.CreateFieldLevelEncryptionProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createFieldLevelEncryptionProfile, request,
        options: options);
  }

  /// Creates a CloudFront function. To create a function, you provide the function code and some configuration information about the function. The response contains an Amazon Resource Name (ARN) that un...
  /// HTTP: POST /2020-05-31/function
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateFunctionResult> createFunction(
    $0.CreateFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createFunction, request, options: options);
  }

  /// Create a new invalidation. For more information, see Invalidating files in the Amazon CloudFront Developer Guide.
  /// HTTP: POST /2020-05-31/distribution/{DistributionId}/invalidation
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateInvalidationResult> createInvalidation(
    $0.CreateInvalidationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createInvalidation, request, options: options);
  }

  /// Creates an invalidation for a distribution tenant. For more information, see Invalidating files in the Amazon CloudFront Developer Guide.
  /// HTTP: POST /2020-05-31/distribution-tenant/{Id}/invalidation
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateInvalidationForDistributionTenantResult>
      createInvalidationForDistributionTenant(
    $0.CreateInvalidationForDistributionTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createInvalidationForDistributionTenant, request,
        options: options);
  }

  /// Creates a key group that you can use with CloudFront signed URLs and signed cookies. To create a key group, you must specify at least one public key for the key group. After you create a key group,...
  /// HTTP: POST /2020-05-31/key-group
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateKeyGroupResult> createKeyGroup(
    $0.CreateKeyGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createKeyGroup, request, options: options);
  }

  /// Specifies the key value store resource to add to your account. In your account, the key value store names must be unique. You can also import key value store data in JSON format from an S3 bucket b...
  /// HTTP: POST /2020-05-31/key-value-store
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateKeyValueStoreResult> createKeyValueStore(
    $0.CreateKeyValueStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createKeyValueStore, request, options: options);
  }

  /// Enables or disables additional Amazon CloudWatch metrics for the specified CloudFront distribution. The additional metrics incur an additional cost. For more information, see Viewing additional Clo...
  /// HTTP: POST /2020-05-31/distributions/{DistributionId}/monitoring-subscription
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateMonitoringSubscriptionResult>
      createMonitoringSubscription(
    $0.CreateMonitoringSubscriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createMonitoringSubscription, request,
        options: options);
  }

  /// Creates a new origin access control in CloudFront. After you create an origin access control, you can add it to an origin in a CloudFront distribution so that CloudFront sends authenticated (signed...
  /// HTTP: POST /2020-05-31/origin-access-control
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateOriginAccessControlResult>
      createOriginAccessControl(
    $0.CreateOriginAccessControlRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createOriginAccessControl, request,
        options: options);
  }

  /// Creates an origin request policy. After you create an origin request policy, you can attach it to one or more cache behaviors. When it's attached to a cache behavior, the origin request policy dete...
  /// HTTP: POST /2020-05-31/origin-request-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateOriginRequestPolicyResult>
      createOriginRequestPolicy(
    $0.CreateOriginRequestPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createOriginRequestPolicy, request,
        options: options);
  }

  /// Uploads a public key to CloudFront that you can use with signed URLs and signed cookies, or with field-level encryption.
  /// HTTP: POST /2020-05-31/public-key
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreatePublicKeyResult> createPublicKey(
    $0.CreatePublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPublicKey, request, options: options);
  }

  /// Creates a real-time log configuration. After you create a real-time log configuration, you can attach it to one or more cache behaviors to send real-time log data to the specified Amazon Kinesis da...
  /// HTTP: POST /2020-05-31/realtime-log-config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateRealtimeLogConfigResult>
      createRealtimeLogConfig(
    $0.CreateRealtimeLogConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRealtimeLogConfig, request,
        options: options);
  }

  /// Creates a response headers policy. A response headers policy contains information about a set of HTTP headers. To create a response headers policy, you provide some metadata about the policy and a ...
  /// HTTP: POST /2020-05-31/response-headers-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateResponseHeadersPolicyResult>
      createResponseHeadersPolicy(
    $0.CreateResponseHeadersPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createResponseHeadersPolicy, request,
        options: options);
  }

  /// This API is deprecated. Amazon CloudFront is deprecating real-time messaging protocol (RTMP) distributions on December 31, 2020. For more information, read the announcement on the Amazon CloudFront...
  /// HTTP: POST /2020-05-31/streaming-distribution
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateStreamingDistributionResult>
      createStreamingDistribution(
    $0.CreateStreamingDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createStreamingDistribution, request,
        options: options);
  }

  /// This API is deprecated. Amazon CloudFront is deprecating real-time messaging protocol (RTMP) distributions on December 31, 2020. For more information, read the announcement on the Amazon CloudFront...
  /// HTTP: POST /2020-05-31/streaming-distribution?WithTags
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateStreamingDistributionWithTagsResult>
      createStreamingDistributionWithTags(
    $0.CreateStreamingDistributionWithTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createStreamingDistributionWithTags, request,
        options: options);
  }

  /// Creates a trust store.
  /// HTTP: POST /2020-05-31/trust-store
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateTrustStoreResult> createTrustStore(
    $0.CreateTrustStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTrustStore, request, options: options);
  }

  /// Create an Amazon CloudFront VPC origin.
  /// HTTP: POST /2020-05-31/vpc-origin
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateVpcOriginResult> createVpcOrigin(
    $0.CreateVpcOriginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createVpcOrigin, request, options: options);
  }

  /// Deletes an Anycast static IP list.
  /// HTTP: DELETE /2020-05-31/anycast-ip-list/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteAnycastIpList(
    $0.DeleteAnycastIpListRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAnycastIpList, request, options: options);
  }

  /// Deletes a cache policy. You cannot delete a cache policy if it's attached to a cache behavior. First update your distributions to remove the cache policy from all cache behaviors, then delete the c...
  /// HTTP: DELETE /2020-05-31/cache-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteCachePolicy(
    $0.DeleteCachePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteCachePolicy, request, options: options);
  }

  /// Delete an origin access identity.
  /// HTTP: DELETE /2020-05-31/origin-access-identity/cloudfront/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteCloudFrontOriginAccessIdentity(
    $0.DeleteCloudFrontOriginAccessIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteCloudFrontOriginAccessIdentity, request,
        options: options);
  }

  /// Deletes a connection function.
  /// HTTP: DELETE /2020-05-31/connection-function/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteConnectionFunction(
    $0.DeleteConnectionFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteConnectionFunction, request,
        options: options);
  }

  /// Deletes a connection group.
  /// HTTP: DELETE /2020-05-31/connection-group/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteConnectionGroup(
    $0.DeleteConnectionGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteConnectionGroup, request, options: options);
  }

  /// Deletes a continuous deployment policy. You cannot delete a continuous deployment policy that's attached to a primary distribution. First update your distribution to remove the continuous deploymen...
  /// HTTP: DELETE /2020-05-31/continuous-deployment-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteContinuousDeploymentPolicy(
    $0.DeleteContinuousDeploymentPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteContinuousDeploymentPolicy, request,
        options: options);
  }

  /// Delete a distribution. Before you can delete a distribution, you must disable it, which requires permission to update the distribution. Once deleted, a distribution cannot be recovered.
  /// HTTP: DELETE /2020-05-31/distribution/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteDistribution(
    $0.DeleteDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDistribution, request, options: options);
  }

  /// Deletes a distribution tenant. If you use this API operation to delete a distribution tenant that is currently enabled, the request will fail. To delete a distribution tenant, you must first disabl...
  /// HTTP: DELETE /2020-05-31/distribution-tenant/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteDistributionTenant(
    $0.DeleteDistributionTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDistributionTenant, request,
        options: options);
  }

  /// Remove a field-level encryption configuration.
  /// HTTP: DELETE /2020-05-31/field-level-encryption/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteFieldLevelEncryptionConfig(
    $0.DeleteFieldLevelEncryptionConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteFieldLevelEncryptionConfig, request,
        options: options);
  }

  /// Remove a field-level encryption profile.
  /// HTTP: DELETE /2020-05-31/field-level-encryption-profile/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteFieldLevelEncryptionProfile(
    $0.DeleteFieldLevelEncryptionProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteFieldLevelEncryptionProfile, request,
        options: options);
  }

  /// Deletes a CloudFront function. You cannot delete a function if it's associated with a cache behavior. First, update your distributions to remove the function association from all cache behaviors, t...
  /// HTTP: DELETE /2020-05-31/function/{Name}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteFunction(
    $0.DeleteFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteFunction, request, options: options);
  }

  /// Deletes a key group. You cannot delete a key group that is referenced in a cache behavior. First update your distributions to remove the key group from all cache behaviors, then delete the key grou...
  /// HTTP: DELETE /2020-05-31/key-group/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteKeyGroup(
    $0.DeleteKeyGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteKeyGroup, request, options: options);
  }

  /// Specifies the key value store to delete.
  /// HTTP: DELETE /2020-05-31/key-value-store/{Name}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteKeyValueStore(
    $0.DeleteKeyValueStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteKeyValueStore, request, options: options);
  }

  /// Disables additional CloudWatch metrics for the specified CloudFront distribution.
  /// HTTP: DELETE /2020-05-31/distributions/{DistributionId}/monitoring-subscription
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteMonitoringSubscriptionResult>
      deleteMonitoringSubscription(
    $0.DeleteMonitoringSubscriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMonitoringSubscription, request,
        options: options);
  }

  /// Deletes a CloudFront origin access control. You cannot delete an origin access control if it's in use. First, update all distributions to remove the origin access control from all origins, then del...
  /// HTTP: DELETE /2020-05-31/origin-access-control/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteOriginAccessControl(
    $0.DeleteOriginAccessControlRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteOriginAccessControl, request,
        options: options);
  }

  /// Deletes an origin request policy. You cannot delete an origin request policy if it's attached to any cache behaviors. First update your distributions to remove the origin request policy from all ca...
  /// HTTP: DELETE /2020-05-31/origin-request-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteOriginRequestPolicy(
    $0.DeleteOriginRequestPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteOriginRequestPolicy, request,
        options: options);
  }

  /// Remove a public key you previously added to CloudFront.
  /// HTTP: DELETE /2020-05-31/public-key/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deletePublicKey(
    $0.DeletePublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePublicKey, request, options: options);
  }

  /// Deletes a real-time log configuration. You cannot delete a real-time log configuration if it's attached to a cache behavior. First update your distributions to remove the real-time log configuratio...
  /// HTTP: POST /2020-05-31/delete-realtime-log-config
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteRealtimeLogConfig(
    $0.DeleteRealtimeLogConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRealtimeLogConfig, request,
        options: options);
  }

  /// Deletes the resource policy attached to the CloudFront resource.
  /// HTTP: POST /2020-05-31/delete-resource-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteResourcePolicy(
    $0.DeleteResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourcePolicy, request, options: options);
  }

  /// Deletes a response headers policy. You cannot delete a response headers policy if it's attached to a cache behavior. First update your distributions to remove the response headers policy from all c...
  /// HTTP: DELETE /2020-05-31/response-headers-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteResponseHeadersPolicy(
    $0.DeleteResponseHeadersPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResponseHeadersPolicy, request,
        options: options);
  }

  /// Delete a streaming distribution. To delete an RTMP distribution using the CloudFront API, perform the following steps. To delete an RTMP distribution using the CloudFront API: Disable the RTMP dist...
  /// HTTP: DELETE /2020-05-31/streaming-distribution/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteStreamingDistribution(
    $0.DeleteStreamingDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteStreamingDistribution, request,
        options: options);
  }

  /// Deletes a trust store.
  /// HTTP: DELETE /2020-05-31/trust-store/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteTrustStore(
    $0.DeleteTrustStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTrustStore, request, options: options);
  }

  /// Delete an Amazon CloudFront VPC origin.
  /// HTTP: DELETE /2020-05-31/vpc-origin/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteVpcOriginResult> deleteVpcOrigin(
    $0.DeleteVpcOriginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteVpcOrigin, request, options: options);
  }

  /// Describes a connection function.
  /// HTTP: GET /2020-05-31/connection-function/{Identifier}/describe
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DescribeConnectionFunctionResult>
      describeConnectionFunction(
    $0.DescribeConnectionFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeConnectionFunction, request,
        options: options);
  }

  /// Gets configuration information and metadata about a CloudFront function, but not the function's code. To get a function's code, use GetFunction. To get configuration information and metadata about ...
  /// HTTP: GET /2020-05-31/function/{Name}/describe
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DescribeFunctionResult> describeFunction(
    $0.DescribeFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeFunction, request, options: options);
  }

  /// Specifies the key value store and its configuration.
  /// HTTP: GET /2020-05-31/key-value-store/{Name}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DescribeKeyValueStoreResult> describeKeyValueStore(
    $0.DescribeKeyValueStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeKeyValueStore, request, options: options);
  }

  /// Disassociates a distribution tenant from the WAF web ACL.
  /// HTTP: PUT /2020-05-31/distribution-tenant/{Id}/disassociate-web-acl
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DisassociateDistributionTenantWebACLResult>
      disassociateDistributionTenantWebACL(
    $0.DisassociateDistributionTenantWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disassociateDistributionTenantWebACL, request,
        options: options);
  }

  /// Disassociates a distribution from the WAF web ACL.
  /// HTTP: PUT /2020-05-31/distribution/{Id}/disassociate-web-acl
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DisassociateDistributionWebACLResult>
      disassociateDistributionWebACL(
    $0.DisassociateDistributionWebACLRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disassociateDistributionWebACL, request,
        options: options);
  }

  /// Gets an Anycast static IP list.
  /// HTTP: GET /2020-05-31/anycast-ip-list/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetAnycastIpListResult> getAnycastIpList(
    $0.GetAnycastIpListRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAnycastIpList, request, options: options);
  }

  /// Gets a cache policy, including the following metadata: The policy's identifier. The date and time when the policy was last modified. To get a cache policy, you must provide the policy's identifier....
  /// HTTP: GET /2020-05-31/cache-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetCachePolicyResult> getCachePolicy(
    $0.GetCachePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCachePolicy, request, options: options);
  }

  /// Gets a cache policy configuration. To get a cache policy configuration, you must provide the policy's identifier. If the cache policy is attached to a distribution's cache behavior, you can get the...
  /// HTTP: GET /2020-05-31/cache-policy/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetCachePolicyConfigResult> getCachePolicyConfig(
    $0.GetCachePolicyConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCachePolicyConfig, request, options: options);
  }

  /// Get the information about an origin access identity.
  /// HTTP: GET /2020-05-31/origin-access-identity/cloudfront/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetCloudFrontOriginAccessIdentityResult>
      getCloudFrontOriginAccessIdentity(
    $0.GetCloudFrontOriginAccessIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCloudFrontOriginAccessIdentity, request,
        options: options);
  }

  /// Get the configuration information about an origin access identity.
  /// HTTP: GET /2020-05-31/origin-access-identity/cloudfront/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetCloudFrontOriginAccessIdentityConfigResult>
      getCloudFrontOriginAccessIdentityConfig(
    $0.GetCloudFrontOriginAccessIdentityConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCloudFrontOriginAccessIdentityConfig, request,
        options: options);
  }

  /// Gets a connection function.
  /// HTTP: GET /2020-05-31/connection-function/{Identifier}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetConnectionFunctionResult> getConnectionFunction(
    $0.GetConnectionFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConnectionFunction, request, options: options);
  }

  /// Gets information about a connection group.
  /// HTTP: GET /2020-05-31/connection-group/{Identifier}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetConnectionGroupResult> getConnectionGroup(
    $0.GetConnectionGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConnectionGroup, request, options: options);
  }

  /// Gets information about a connection group by using the endpoint that you specify.
  /// HTTP: GET /2020-05-31/connection-group
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetConnectionGroupByRoutingEndpointResult>
      getConnectionGroupByRoutingEndpoint(
    $0.GetConnectionGroupByRoutingEndpointRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConnectionGroupByRoutingEndpoint, request,
        options: options);
  }

  /// Gets a continuous deployment policy, including metadata (the policy's identifier and the date and time when the policy was last modified).
  /// HTTP: GET /2020-05-31/continuous-deployment-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetContinuousDeploymentPolicyResult>
      getContinuousDeploymentPolicy(
    $0.GetContinuousDeploymentPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getContinuousDeploymentPolicy, request,
        options: options);
  }

  /// Gets configuration information about a continuous deployment policy.
  /// HTTP: GET /2020-05-31/continuous-deployment-policy/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetContinuousDeploymentPolicyConfigResult>
      getContinuousDeploymentPolicyConfig(
    $0.GetContinuousDeploymentPolicyConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getContinuousDeploymentPolicyConfig, request,
        options: options);
  }

  /// Get the information about a distribution.
  /// HTTP: GET /2020-05-31/distribution/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetDistributionResult> getDistribution(
    $0.GetDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDistribution, request, options: options);
  }

  /// Get the configuration information about a distribution.
  /// HTTP: GET /2020-05-31/distribution/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetDistributionConfigResult> getDistributionConfig(
    $0.GetDistributionConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDistributionConfig, request, options: options);
  }

  /// Gets information about a distribution tenant.
  /// HTTP: GET /2020-05-31/distribution-tenant/{Identifier}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetDistributionTenantResult> getDistributionTenant(
    $0.GetDistributionTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDistributionTenant, request, options: options);
  }

  /// Gets information about a distribution tenant by the associated domain.
  /// HTTP: GET /2020-05-31/distribution-tenant
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetDistributionTenantByDomainResult>
      getDistributionTenantByDomain(
    $0.GetDistributionTenantByDomainRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDistributionTenantByDomain, request,
        options: options);
  }

  /// Get the field-level encryption configuration information.
  /// HTTP: GET /2020-05-31/field-level-encryption/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetFieldLevelEncryptionResult>
      getFieldLevelEncryption(
    $0.GetFieldLevelEncryptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFieldLevelEncryption, request,
        options: options);
  }

  /// Get the field-level encryption configuration information.
  /// HTTP: GET /2020-05-31/field-level-encryption/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetFieldLevelEncryptionConfigResult>
      getFieldLevelEncryptionConfig(
    $0.GetFieldLevelEncryptionConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFieldLevelEncryptionConfig, request,
        options: options);
  }

  /// Get the field-level encryption profile information.
  /// HTTP: GET /2020-05-31/field-level-encryption-profile/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetFieldLevelEncryptionProfileResult>
      getFieldLevelEncryptionProfile(
    $0.GetFieldLevelEncryptionProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFieldLevelEncryptionProfile, request,
        options: options);
  }

  /// Get the field-level encryption profile configuration information.
  /// HTTP: GET /2020-05-31/field-level-encryption-profile/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetFieldLevelEncryptionProfileConfigResult>
      getFieldLevelEncryptionProfileConfig(
    $0.GetFieldLevelEncryptionProfileConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFieldLevelEncryptionProfileConfig, request,
        options: options);
  }

  /// Gets the code of a CloudFront function. To get configuration information and metadata about a function, use DescribeFunction. To get a function's code, you must provide the function's name and stag...
  /// HTTP: GET /2020-05-31/function/{Name}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetFunctionResult> getFunction(
    $0.GetFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFunction, request, options: options);
  }

  /// Get the information about an invalidation.
  /// HTTP: GET /2020-05-31/distribution/{DistributionId}/invalidation/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetInvalidationResult> getInvalidation(
    $0.GetInvalidationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getInvalidation, request, options: options);
  }

  /// Gets information about a specific invalidation for a distribution tenant.
  /// HTTP: GET /2020-05-31/distribution-tenant/{DistributionTenantId}/invalidation/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetInvalidationForDistributionTenantResult>
      getInvalidationForDistributionTenant(
    $0.GetInvalidationForDistributionTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getInvalidationForDistributionTenant, request,
        options: options);
  }

  /// Gets a key group, including the date and time when the key group was last modified. To get a key group, you must provide the key group's identifier. If the key group is referenced in a distribution...
  /// HTTP: GET /2020-05-31/key-group/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetKeyGroupResult> getKeyGroup(
    $0.GetKeyGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getKeyGroup, request, options: options);
  }

  /// Gets a key group configuration. To get a key group configuration, you must provide the key group's identifier. If the key group is referenced in a distribution's cache behavior, you can get the key...
  /// HTTP: GET /2020-05-31/key-group/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetKeyGroupConfigResult> getKeyGroupConfig(
    $0.GetKeyGroupConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getKeyGroupConfig, request, options: options);
  }

  /// Gets details about the CloudFront managed ACM certificate.
  /// HTTP: GET /2020-05-31/managed-certificate/{Identifier}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetManagedCertificateDetailsResult>
      getManagedCertificateDetails(
    $0.GetManagedCertificateDetailsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getManagedCertificateDetails, request,
        options: options);
  }

  /// Gets information about whether additional CloudWatch metrics are enabled for the specified CloudFront distribution.
  /// HTTP: GET /2020-05-31/distributions/{DistributionId}/monitoring-subscription
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetMonitoringSubscriptionResult>
      getMonitoringSubscription(
    $0.GetMonitoringSubscriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMonitoringSubscription, request,
        options: options);
  }

  /// Gets a CloudFront origin access control, including its unique identifier.
  /// HTTP: GET /2020-05-31/origin-access-control/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetOriginAccessControlResult> getOriginAccessControl(
    $0.GetOriginAccessControlRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOriginAccessControl, request,
        options: options);
  }

  /// Gets a CloudFront origin access control configuration.
  /// HTTP: GET /2020-05-31/origin-access-control/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetOriginAccessControlConfigResult>
      getOriginAccessControlConfig(
    $0.GetOriginAccessControlConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOriginAccessControlConfig, request,
        options: options);
  }

  /// Gets an origin request policy, including the following metadata: The policy's identifier. The date and time when the policy was last modified. To get an origin request policy, you must provide the ...
  /// HTTP: GET /2020-05-31/origin-request-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetOriginRequestPolicyResult> getOriginRequestPolicy(
    $0.GetOriginRequestPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOriginRequestPolicy, request,
        options: options);
  }

  /// Gets an origin request policy configuration. To get an origin request policy configuration, you must provide the policy's identifier. If the origin request policy is attached to a distribution's ca...
  /// HTTP: GET /2020-05-31/origin-request-policy/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetOriginRequestPolicyConfigResult>
      getOriginRequestPolicyConfig(
    $0.GetOriginRequestPolicyConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOriginRequestPolicyConfig, request,
        options: options);
  }

  /// Gets a public key.
  /// HTTP: GET /2020-05-31/public-key/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetPublicKeyResult> getPublicKey(
    $0.GetPublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPublicKey, request, options: options);
  }

  /// Gets a public key configuration.
  /// HTTP: GET /2020-05-31/public-key/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetPublicKeyConfigResult> getPublicKeyConfig(
    $0.GetPublicKeyConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPublicKeyConfig, request, options: options);
  }

  /// Gets a real-time log configuration. To get a real-time log configuration, you can provide the configuration's name or its Amazon Resource Name (ARN). You must provide at least one. If you provide b...
  /// HTTP: POST /2020-05-31/get-realtime-log-config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetRealtimeLogConfigResult> getRealtimeLogConfig(
    $0.GetRealtimeLogConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRealtimeLogConfig, request, options: options);
  }

  /// Retrieves the resource policy for the specified CloudFront resource that you own and have shared.
  /// HTTP: POST /2020-05-31/get-resource-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetResourcePolicyResult> getResourcePolicy(
    $0.GetResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResourcePolicy, request, options: options);
  }

  /// Gets a response headers policy, including metadata (the policy's identifier and the date and time when the policy was last modified). To get a response headers policy, you must provide the policy's...
  /// HTTP: GET /2020-05-31/response-headers-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetResponseHeadersPolicyResult>
      getResponseHeadersPolicy(
    $0.GetResponseHeadersPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResponseHeadersPolicy, request,
        options: options);
  }

  /// Gets a response headers policy configuration. To get a response headers policy configuration, you must provide the policy's identifier. If the response headers policy is attached to a distribution'...
  /// HTTP: GET /2020-05-31/response-headers-policy/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetResponseHeadersPolicyConfigResult>
      getResponseHeadersPolicyConfig(
    $0.GetResponseHeadersPolicyConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResponseHeadersPolicyConfig, request,
        options: options);
  }

  /// Gets information about a specified RTMP distribution, including the distribution configuration.
  /// HTTP: GET /2020-05-31/streaming-distribution/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetStreamingDistributionResult>
      getStreamingDistribution(
    $0.GetStreamingDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getStreamingDistribution, request,
        options: options);
  }

  /// Get the configuration information about a streaming distribution.
  /// HTTP: GET /2020-05-31/streaming-distribution/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetStreamingDistributionConfigResult>
      getStreamingDistributionConfig(
    $0.GetStreamingDistributionConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getStreamingDistributionConfig, request,
        options: options);
  }

  /// Gets a trust store.
  /// HTTP: GET /2020-05-31/trust-store/{Identifier}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetTrustStoreResult> getTrustStore(
    $0.GetTrustStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTrustStore, request, options: options);
  }

  /// Get the details of an Amazon CloudFront VPC origin.
  /// HTTP: GET /2020-05-31/vpc-origin/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetVpcOriginResult> getVpcOrigin(
    $0.GetVpcOriginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getVpcOrigin, request, options: options);
  }

  /// Lists your Anycast static IP lists.
  /// HTTP: GET /2020-05-31/anycast-ip-list
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListAnycastIpListsResult> listAnycastIpLists(
    $0.ListAnycastIpListsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAnycastIpLists, request, options: options);
  }

  /// Gets a list of cache policies. You can optionally apply a filter to return only the managed policies created by Amazon Web Services, or only the custom policies created in your Amazon Web Services ...
  /// HTTP: GET /2020-05-31/cache-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListCachePoliciesResult> listCachePolicies(
    $0.ListCachePoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCachePolicies, request, options: options);
  }

  /// Lists origin access identities.
  /// HTTP: GET /2020-05-31/origin-access-identity/cloudfront
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListCloudFrontOriginAccessIdentitiesResult>
      listCloudFrontOriginAccessIdentities(
    $0.ListCloudFrontOriginAccessIdentitiesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCloudFrontOriginAccessIdentities, request,
        options: options);
  }

  /// The ListConflictingAliases API operation only supports standard distributions. To list domain conflicts for both standard distributions and distribution tenants, we recommend that you use the ListD...
  /// HTTP: GET /2020-05-31/conflicting-alias
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListConflictingAliasesResult> listConflictingAliases(
    $0.ListConflictingAliasesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listConflictingAliases, request,
        options: options);
  }

  /// Lists connection functions.
  /// HTTP: POST /2020-05-31/connection-functions
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListConnectionFunctionsResult>
      listConnectionFunctions(
    $0.ListConnectionFunctionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listConnectionFunctions, request,
        options: options);
  }

  /// Lists the connection groups in your Amazon Web Services account.
  /// HTTP: POST /2020-05-31/connection-groups
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListConnectionGroupsResult> listConnectionGroups(
    $0.ListConnectionGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listConnectionGroups, request, options: options);
  }

  /// Gets a list of the continuous deployment policies in your Amazon Web Services account. You can optionally specify the maximum number of items to receive in the response. If the total number of item...
  /// HTTP: GET /2020-05-31/continuous-deployment-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListContinuousDeploymentPoliciesResult>
      listContinuousDeploymentPolicies(
    $0.ListContinuousDeploymentPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listContinuousDeploymentPolicies, request,
        options: options);
  }

  /// List CloudFront distributions.
  /// HTTP: GET /2020-05-31/distribution
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsResult> listDistributions(
    $0.ListDistributionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributions, request, options: options);
  }

  /// Lists the distributions in your account that are associated with the specified AnycastIpListId.
  /// HTTP: GET /2020-05-31/distributionsByAnycastIpListId/{AnycastIpListId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByAnycastIpListIdResult>
      listDistributionsByAnycastIpListId(
    $0.ListDistributionsByAnycastIpListIdRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByAnycastIpListId, request,
        options: options);
  }

  /// Gets a list of distribution IDs for distributions that have a cache behavior that's associated with the specified cache policy. You can optionally specify the maximum number of items to receive in ...
  /// HTTP: GET /2020-05-31/distributionsByCachePolicyId/{CachePolicyId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByCachePolicyIdResult>
      listDistributionsByCachePolicyId(
    $0.ListDistributionsByCachePolicyIdRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByCachePolicyId, request,
        options: options);
  }

  /// Lists distributions by connection function.
  /// HTTP: GET /2020-05-31/distributionsByConnectionFunction
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByConnectionFunctionResult>
      listDistributionsByConnectionFunction(
    $0.ListDistributionsByConnectionFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByConnectionFunction, request,
        options: options);
  }

  /// Lists the distributions by the connection mode that you specify.
  /// HTTP: GET /2020-05-31/distributionsByConnectionMode/{ConnectionMode}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByConnectionModeResult>
      listDistributionsByConnectionMode(
    $0.ListDistributionsByConnectionModeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByConnectionMode, request,
        options: options);
  }

  /// Gets a list of distribution IDs for distributions that have a cache behavior that references the specified key group. You can optionally specify the maximum number of items to receive in the respon...
  /// HTTP: GET /2020-05-31/distributionsByKeyGroupId/{KeyGroupId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByKeyGroupResult>
      listDistributionsByKeyGroup(
    $0.ListDistributionsByKeyGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByKeyGroup, request,
        options: options);
  }

  /// Gets a list of distribution IDs for distributions that have a cache behavior that's associated with the specified origin request policy. You can optionally specify the maximum number of items to re...
  /// HTTP: GET /2020-05-31/distributionsByOriginRequestPolicyId/{OriginRequestPolicyId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByOriginRequestPolicyIdResult>
      listDistributionsByOriginRequestPolicyId(
    $0.ListDistributionsByOriginRequestPolicyIdRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByOriginRequestPolicyId, request,
        options: options);
  }

  /// Lists the CloudFront distributions that are associated with the specified resource that you own.
  /// HTTP: GET /2020-05-31/distributionsByOwnedResource/{ResourceArn}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByOwnedResourceResult>
      listDistributionsByOwnedResource(
    $0.ListDistributionsByOwnedResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByOwnedResource, request,
        options: options);
  }

  /// Gets a list of distributions that have a cache behavior that's associated with the specified real-time log configuration. You can specify the real-time log configuration by its name or its Amazon R...
  /// HTTP: POST /2020-05-31/distributionsByRealtimeLogConfig
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByRealtimeLogConfigResult>
      listDistributionsByRealtimeLogConfig(
    $0.ListDistributionsByRealtimeLogConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByRealtimeLogConfig, request,
        options: options);
  }

  /// Gets a list of distribution IDs for distributions that have a cache behavior that's associated with the specified response headers policy. You can optionally specify the maximum number of items to ...
  /// HTTP: GET /2020-05-31/distributionsByResponseHeadersPolicyId/{ResponseHeadersPolicyId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByResponseHeadersPolicyIdResult>
      listDistributionsByResponseHeadersPolicyId(
    $0.ListDistributionsByResponseHeadersPolicyIdRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$listDistributionsByResponseHeadersPolicyId, request,
        options: options);
  }

  /// Lists distributions by trust store.
  /// HTTP: GET /2020-05-31/distributionsByTrustStore
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByTrustStoreResult>
      listDistributionsByTrustStore(
    $0.ListDistributionsByTrustStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByTrustStore, request,
        options: options);
  }

  /// List CloudFront distributions by their VPC origin ID.
  /// HTTP: GET /2020-05-31/distributionsByVpcOriginId/{VpcOriginId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByVpcOriginIdResult>
      listDistributionsByVpcOriginId(
    $0.ListDistributionsByVpcOriginIdRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByVpcOriginId, request,
        options: options);
  }

  /// List the distributions that are associated with a specified WAF web ACL.
  /// HTTP: GET /2020-05-31/distributionsByWebACLId/{WebACLId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionsByWebACLIdResult>
      listDistributionsByWebACLId(
    $0.ListDistributionsByWebACLIdRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionsByWebACLId, request,
        options: options);
  }

  /// Lists the distribution tenants in your Amazon Web Services account.
  /// HTTP: POST /2020-05-31/distribution-tenants
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionTenantsResult>
      listDistributionTenants(
    $0.ListDistributionTenantsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionTenants, request,
        options: options);
  }

  /// Lists distribution tenants by the customization that you specify. You must specify either the CertificateArn parameter or WebACLArn parameter, but not both in the same request.
  /// HTTP: POST /2020-05-31/distribution-tenants-by-customization
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDistributionTenantsByCustomizationResult>
      listDistributionTenantsByCustomization(
    $0.ListDistributionTenantsByCustomizationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDistributionTenantsByCustomization, request,
        options: options);
  }

  /// We recommend that you use the ListDomainConflicts API operation to check for domain conflicts, as it supports both standard distributions and distribution tenants. ListConflictingAliases performs s...
  /// HTTP: POST /2020-05-31/domain-conflicts
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDomainConflictsResult> listDomainConflicts(
    $0.ListDomainConflictsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDomainConflicts, request, options: options);
  }

  /// List all field-level encryption configurations that have been created in CloudFront for this account.
  /// HTTP: GET /2020-05-31/field-level-encryption
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListFieldLevelEncryptionConfigsResult>
      listFieldLevelEncryptionConfigs(
    $0.ListFieldLevelEncryptionConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFieldLevelEncryptionConfigs, request,
        options: options);
  }

  /// Request a list of field-level encryption profiles that have been created in CloudFront for this account.
  /// HTTP: GET /2020-05-31/field-level-encryption-profile
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListFieldLevelEncryptionProfilesResult>
      listFieldLevelEncryptionProfiles(
    $0.ListFieldLevelEncryptionProfilesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFieldLevelEncryptionProfiles, request,
        options: options);
  }

  /// Gets a list of all CloudFront functions in your Amazon Web Services account. You can optionally apply a filter to return only the functions that are in the specified stage, either DEVELOPMENT or LI...
  /// HTTP: GET /2020-05-31/function
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListFunctionsResult> listFunctions(
    $0.ListFunctionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFunctions, request, options: options);
  }

  /// Lists invalidation batches.
  /// HTTP: GET /2020-05-31/distribution/{DistributionId}/invalidation
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListInvalidationsResult> listInvalidations(
    $0.ListInvalidationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listInvalidations, request, options: options);
  }

  /// Lists the invalidations for a distribution tenant.
  /// HTTP: GET /2020-05-31/distribution-tenant/{Id}/invalidation
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListInvalidationsForDistributionTenantResult>
      listInvalidationsForDistributionTenant(
    $0.ListInvalidationsForDistributionTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listInvalidationsForDistributionTenant, request,
        options: options);
  }

  /// Gets a list of key groups. You can optionally specify the maximum number of items to receive in the response. If the total number of items in the list exceeds the maximum that you specify, or the d...
  /// HTTP: GET /2020-05-31/key-group
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListKeyGroupsResult> listKeyGroups(
    $0.ListKeyGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listKeyGroups, request, options: options);
  }

  /// Specifies the key value stores to list.
  /// HTTP: GET /2020-05-31/key-value-store
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListKeyValueStoresResult> listKeyValueStores(
    $0.ListKeyValueStoresRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listKeyValueStores, request, options: options);
  }

  /// Gets the list of CloudFront origin access controls (OACs) in this Amazon Web Services account. You can optionally specify the maximum number of items to receive in the response. If the total number...
  /// HTTP: GET /2020-05-31/origin-access-control
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListOriginAccessControlsResult>
      listOriginAccessControls(
    $0.ListOriginAccessControlsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOriginAccessControls, request,
        options: options);
  }

  /// Gets a list of origin request policies. You can optionally apply a filter to return only the managed policies created by Amazon Web Services, or only the custom policies created in your Amazon Web ...
  /// HTTP: GET /2020-05-31/origin-request-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListOriginRequestPoliciesResult>
      listOriginRequestPolicies(
    $0.ListOriginRequestPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOriginRequestPolicies, request,
        options: options);
  }

  /// List all public keys that have been added to CloudFront for this account.
  /// HTTP: GET /2020-05-31/public-key
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListPublicKeysResult> listPublicKeys(
    $0.ListPublicKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPublicKeys, request, options: options);
  }

  /// Gets a list of real-time log configurations. You can optionally specify the maximum number of items to receive in the response. If the total number of items in the list exceeds the maximum that you...
  /// HTTP: GET /2020-05-31/realtime-log-config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListRealtimeLogConfigsResult> listRealtimeLogConfigs(
    $0.ListRealtimeLogConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRealtimeLogConfigs, request,
        options: options);
  }

  /// Gets a list of response headers policies. You can optionally apply a filter to get only the managed policies created by Amazon Web Services, or only the custom policies created in your Amazon Web S...
  /// HTTP: GET /2020-05-31/response-headers-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListResponseHeadersPoliciesResult>
      listResponseHeadersPolicies(
    $0.ListResponseHeadersPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listResponseHeadersPolicies, request,
        options: options);
  }

  /// List streaming distributions.
  /// HTTP: GET /2020-05-31/streaming-distribution
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListStreamingDistributionsResult>
      listStreamingDistributions(
    $0.ListStreamingDistributionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listStreamingDistributions, request,
        options: options);
  }

  /// List tags for a CloudFront resource. For more information, see Tagging a distribution in the Amazon CloudFront Developer Guide.
  /// HTTP: GET /2020-05-31/tagging
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTagsForResourceResult> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Lists trust stores.
  /// HTTP: POST /2020-05-31/trust-stores
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTrustStoresResult> listTrustStores(
    $0.ListTrustStoresRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTrustStores, request, options: options);
  }

  /// List the CloudFront VPC origins in your account.
  /// HTTP: GET /2020-05-31/vpc-origin
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListVpcOriginsResult> listVpcOrigins(
    $0.ListVpcOriginsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listVpcOrigins, request, options: options);
  }

  /// Publishes a connection function.
  /// HTTP: POST /2020-05-31/connection-function/{Id}/publish
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PublishConnectionFunctionResult>
      publishConnectionFunction(
    $0.PublishConnectionFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$publishConnectionFunction, request,
        options: options);
  }

  /// Publishes a CloudFront function by copying the function code from the DEVELOPMENT stage to LIVE. This automatically updates all cache behaviors that are using this function to use the newly publish...
  /// HTTP: POST /2020-05-31/function/{Name}/publish
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PublishFunctionResult> publishFunction(
    $0.PublishFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$publishFunction, request, options: options);
  }

  /// Creates a resource control policy for a given CloudFront resource.
  /// HTTP: POST /2020-05-31/put-resource-policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PutResourcePolicyResult> putResourcePolicy(
    $0.PutResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putResourcePolicy, request, options: options);
  }

  /// Add tags to a CloudFront resource. For more information, see Tagging a distribution in the Amazon CloudFront Developer Guide.
  /// HTTP: POST /2020-05-31/tagging?Operation=Tag
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Tests a connection function.
  /// HTTP: POST /2020-05-31/connection-function/{Id}/test
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.TestConnectionFunctionResult> testConnectionFunction(
    $0.TestConnectionFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testConnectionFunction, request,
        options: options);
  }

  /// Tests a CloudFront function. To test a function, you provide an event object that represents an HTTP request or response that your CloudFront distribution could receive in production. CloudFront ru...
  /// HTTP: POST /2020-05-31/function/{Name}/test
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.TestFunctionResult> testFunction(
    $0.TestFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testFunction, request, options: options);
  }

  /// Remove tags from a CloudFront resource. For more information, see Tagging a distribution in the Amazon CloudFront Developer Guide.
  /// HTTP: POST /2020-05-31/tagging?Operation=Untag
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Updates an Anycast static IP list.
  /// HTTP: PUT /2020-05-31/anycast-ip-list/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateAnycastIpListResult> updateAnycastIpList(
    $0.UpdateAnycastIpListRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAnycastIpList, request, options: options);
  }

  /// Updates a cache policy configuration. When you update a cache policy configuration, all the fields are updated with the values provided in the request. You cannot update some fields independent of ...
  /// HTTP: PUT /2020-05-31/cache-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateCachePolicyResult> updateCachePolicy(
    $0.UpdateCachePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateCachePolicy, request, options: options);
  }

  /// Update an origin access identity.
  /// HTTP: PUT /2020-05-31/origin-access-identity/cloudfront/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateCloudFrontOriginAccessIdentityResult>
      updateCloudFrontOriginAccessIdentity(
    $0.UpdateCloudFrontOriginAccessIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateCloudFrontOriginAccessIdentity, request,
        options: options);
  }

  /// Updates a connection function.
  /// HTTP: PUT /2020-05-31/connection-function/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateConnectionFunctionResult>
      updateConnectionFunction(
    $0.UpdateConnectionFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateConnectionFunction, request,
        options: options);
  }

  /// Updates a connection group.
  /// HTTP: PUT /2020-05-31/connection-group/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateConnectionGroupResult> updateConnectionGroup(
    $0.UpdateConnectionGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateConnectionGroup, request, options: options);
  }

  /// Updates a continuous deployment policy. You can update a continuous deployment policy to enable or disable it, to change the percentage of traffic that it sends to the staging distribution, or to c...
  /// HTTP: PUT /2020-05-31/continuous-deployment-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateContinuousDeploymentPolicyResult>
      updateContinuousDeploymentPolicy(
    $0.UpdateContinuousDeploymentPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateContinuousDeploymentPolicy, request,
        options: options);
  }

  /// Updates the configuration for a CloudFront distribution. The update process includes getting the current distribution configuration, updating it to make your changes, and then submitting an UpdateD...
  /// HTTP: PUT /2020-05-31/distribution/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateDistributionResult> updateDistribution(
    $0.UpdateDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDistribution, request, options: options);
  }

  /// Updates a distribution tenant.
  /// HTTP: PUT /2020-05-31/distribution-tenant/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateDistributionTenantResult>
      updateDistributionTenant(
    $0.UpdateDistributionTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDistributionTenant, request,
        options: options);
  }

  /// Copies the staging distribution's configuration to its corresponding primary distribution. The primary distribution retains its Aliases (also known as alternate domain names or CNAMEs) and Continuo...
  /// HTTP: PUT /2020-05-31/distribution/{Id}/promote-staging-config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateDistributionWithStagingConfigResult>
      updateDistributionWithStagingConfig(
    $0.UpdateDistributionWithStagingConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDistributionWithStagingConfig, request,
        options: options);
  }

  /// We recommend that you use the UpdateDomainAssociation API operation to move a domain association, as it supports both standard distributions and distribution tenants. AssociateAlias performs simila...
  /// HTTP: POST /2020-05-31/domain-association
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateDomainAssociationResult>
      updateDomainAssociation(
    $0.UpdateDomainAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDomainAssociation, request,
        options: options);
  }

  /// Update a field-level encryption configuration.
  /// HTTP: PUT /2020-05-31/field-level-encryption/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateFieldLevelEncryptionConfigResult>
      updateFieldLevelEncryptionConfig(
    $0.UpdateFieldLevelEncryptionConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateFieldLevelEncryptionConfig, request,
        options: options);
  }

  /// Update a field-level encryption profile.
  /// HTTP: PUT /2020-05-31/field-level-encryption-profile/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateFieldLevelEncryptionProfileResult>
      updateFieldLevelEncryptionProfile(
    $0.UpdateFieldLevelEncryptionProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateFieldLevelEncryptionProfile, request,
        options: options);
  }

  /// Updates a CloudFront function. You can update a function's code or the comment that describes the function. You cannot update a function's name. To update a function, you provide the function's nam...
  /// HTTP: PUT /2020-05-31/function/{Name}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateFunctionResult> updateFunction(
    $0.UpdateFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateFunction, request, options: options);
  }

  /// Updates a key group. When you update a key group, all the fields are updated with the values provided in the request. You cannot update some fields independent of others. To update a key group: Get...
  /// HTTP: PUT /2020-05-31/key-group/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateKeyGroupResult> updateKeyGroup(
    $0.UpdateKeyGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateKeyGroup, request, options: options);
  }

  /// Specifies the key value store to update.
  /// HTTP: PUT /2020-05-31/key-value-store/{Name}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateKeyValueStoreResult> updateKeyValueStore(
    $0.UpdateKeyValueStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateKeyValueStore, request, options: options);
  }

  /// Updates a CloudFront origin access control.
  /// HTTP: PUT /2020-05-31/origin-access-control/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateOriginAccessControlResult>
      updateOriginAccessControl(
    $0.UpdateOriginAccessControlRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateOriginAccessControl, request,
        options: options);
  }

  /// Updates an origin request policy configuration. When you update an origin request policy configuration, all the fields are updated with the values provided in the request. You cannot update some fi...
  /// HTTP: PUT /2020-05-31/origin-request-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateOriginRequestPolicyResult>
      updateOriginRequestPolicy(
    $0.UpdateOriginRequestPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateOriginRequestPolicy, request,
        options: options);
  }

  /// Update public key information. Note that the only value you can change is the comment.
  /// HTTP: PUT /2020-05-31/public-key/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdatePublicKeyResult> updatePublicKey(
    $0.UpdatePublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updatePublicKey, request, options: options);
  }

  /// Updates a real-time log configuration. When you update a real-time log configuration, all the parameters are updated with the values provided in the request. You cannot update some parameters indep...
  /// HTTP: PUT /2020-05-31/realtime-log-config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateRealtimeLogConfigResult>
      updateRealtimeLogConfig(
    $0.UpdateRealtimeLogConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRealtimeLogConfig, request,
        options: options);
  }

  /// Updates a response headers policy. When you update a response headers policy, the entire policy is replaced. You cannot update some policy fields independent of others. To update a response headers...
  /// HTTP: PUT /2020-05-31/response-headers-policy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateResponseHeadersPolicyResult>
      updateResponseHeadersPolicy(
    $0.UpdateResponseHeadersPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateResponseHeadersPolicy, request,
        options: options);
  }

  /// Update a streaming distribution.
  /// HTTP: PUT /2020-05-31/streaming-distribution/{Id}/config
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateStreamingDistributionResult>
      updateStreamingDistribution(
    $0.UpdateStreamingDistributionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateStreamingDistribution, request,
        options: options);
  }

  /// Updates a trust store.
  /// HTTP: PUT /2020-05-31/trust-store/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateTrustStoreResult> updateTrustStore(
    $0.UpdateTrustStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTrustStore, request, options: options);
  }

  /// Update an Amazon CloudFront VPC origin in your account.
  /// HTTP: PUT /2020-05-31/vpc-origin/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateVpcOriginResult> updateVpcOrigin(
    $0.UpdateVpcOriginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateVpcOrigin, request, options: options);
  }

  /// Verify the DNS configuration for your domain names. This API operation checks whether your domain name points to the correct routing endpoint of the connection group, such as d111111abcdef8.cloudfr...
  /// HTTP: POST /2020-05-31/verify-dns-configuration
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.VerifyDnsConfigurationResult> verifyDnsConfiguration(
    $0.VerifyDnsConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$verifyDnsConfiguration, request,
        options: options);
  }

  // method descriptors

  static final _$associateAlias =
      $grpc.ClientMethod<$0.AssociateAliasRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/AssociateAlias',
          ($0.AssociateAliasRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$associateDistributionTenantWebACL = $grpc.ClientMethod<
          $0.AssociateDistributionTenantWebACLRequest,
          $0.AssociateDistributionTenantWebACLResult>(
      '/cloudfront.CloudFrontService/AssociateDistributionTenantWebACL',
      ($0.AssociateDistributionTenantWebACLRequest value) =>
          value.writeToBuffer(),
      $0.AssociateDistributionTenantWebACLResult.fromBuffer);
  static final _$associateDistributionWebACL = $grpc.ClientMethod<
          $0.AssociateDistributionWebACLRequest,
          $0.AssociateDistributionWebACLResult>(
      '/cloudfront.CloudFrontService/AssociateDistributionWebACL',
      ($0.AssociateDistributionWebACLRequest value) => value.writeToBuffer(),
      $0.AssociateDistributionWebACLResult.fromBuffer);
  static final _$copyDistribution =
      $grpc.ClientMethod<$0.CopyDistributionRequest, $0.CopyDistributionResult>(
          '/cloudfront.CloudFrontService/CopyDistribution',
          ($0.CopyDistributionRequest value) => value.writeToBuffer(),
          $0.CopyDistributionResult.fromBuffer);
  static final _$createAnycastIpList = $grpc.ClientMethod<
          $0.CreateAnycastIpListRequest, $0.CreateAnycastIpListResult>(
      '/cloudfront.CloudFrontService/CreateAnycastIpList',
      ($0.CreateAnycastIpListRequest value) => value.writeToBuffer(),
      $0.CreateAnycastIpListResult.fromBuffer);
  static final _$createCachePolicy = $grpc.ClientMethod<
          $0.CreateCachePolicyRequest, $0.CreateCachePolicyResult>(
      '/cloudfront.CloudFrontService/CreateCachePolicy',
      ($0.CreateCachePolicyRequest value) => value.writeToBuffer(),
      $0.CreateCachePolicyResult.fromBuffer);
  static final _$createCloudFrontOriginAccessIdentity = $grpc.ClientMethod<
          $0.CreateCloudFrontOriginAccessIdentityRequest,
          $0.CreateCloudFrontOriginAccessIdentityResult>(
      '/cloudfront.CloudFrontService/CreateCloudFrontOriginAccessIdentity',
      ($0.CreateCloudFrontOriginAccessIdentityRequest value) =>
          value.writeToBuffer(),
      $0.CreateCloudFrontOriginAccessIdentityResult.fromBuffer);
  static final _$createConnectionFunction = $grpc.ClientMethod<
          $0.CreateConnectionFunctionRequest,
          $0.CreateConnectionFunctionResult>(
      '/cloudfront.CloudFrontService/CreateConnectionFunction',
      ($0.CreateConnectionFunctionRequest value) => value.writeToBuffer(),
      $0.CreateConnectionFunctionResult.fromBuffer);
  static final _$createConnectionGroup = $grpc.ClientMethod<
          $0.CreateConnectionGroupRequest, $0.CreateConnectionGroupResult>(
      '/cloudfront.CloudFrontService/CreateConnectionGroup',
      ($0.CreateConnectionGroupRequest value) => value.writeToBuffer(),
      $0.CreateConnectionGroupResult.fromBuffer);
  static final _$createContinuousDeploymentPolicy = $grpc.ClientMethod<
          $0.CreateContinuousDeploymentPolicyRequest,
          $0.CreateContinuousDeploymentPolicyResult>(
      '/cloudfront.CloudFrontService/CreateContinuousDeploymentPolicy',
      ($0.CreateContinuousDeploymentPolicyRequest value) =>
          value.writeToBuffer(),
      $0.CreateContinuousDeploymentPolicyResult.fromBuffer);
  static final _$createDistribution = $grpc.ClientMethod<
          $0.CreateDistributionRequest, $0.CreateDistributionResult>(
      '/cloudfront.CloudFrontService/CreateDistribution',
      ($0.CreateDistributionRequest value) => value.writeToBuffer(),
      $0.CreateDistributionResult.fromBuffer);
  static final _$createDistributionTenant = $grpc.ClientMethod<
          $0.CreateDistributionTenantRequest,
          $0.CreateDistributionTenantResult>(
      '/cloudfront.CloudFrontService/CreateDistributionTenant',
      ($0.CreateDistributionTenantRequest value) => value.writeToBuffer(),
      $0.CreateDistributionTenantResult.fromBuffer);
  static final _$createDistributionWithTags = $grpc.ClientMethod<
          $0.CreateDistributionWithTagsRequest,
          $0.CreateDistributionWithTagsResult>(
      '/cloudfront.CloudFrontService/CreateDistributionWithTags',
      ($0.CreateDistributionWithTagsRequest value) => value.writeToBuffer(),
      $0.CreateDistributionWithTagsResult.fromBuffer);
  static final _$createFieldLevelEncryptionConfig = $grpc.ClientMethod<
          $0.CreateFieldLevelEncryptionConfigRequest,
          $0.CreateFieldLevelEncryptionConfigResult>(
      '/cloudfront.CloudFrontService/CreateFieldLevelEncryptionConfig',
      ($0.CreateFieldLevelEncryptionConfigRequest value) =>
          value.writeToBuffer(),
      $0.CreateFieldLevelEncryptionConfigResult.fromBuffer);
  static final _$createFieldLevelEncryptionProfile = $grpc.ClientMethod<
          $0.CreateFieldLevelEncryptionProfileRequest,
          $0.CreateFieldLevelEncryptionProfileResult>(
      '/cloudfront.CloudFrontService/CreateFieldLevelEncryptionProfile',
      ($0.CreateFieldLevelEncryptionProfileRequest value) =>
          value.writeToBuffer(),
      $0.CreateFieldLevelEncryptionProfileResult.fromBuffer);
  static final _$createFunction =
      $grpc.ClientMethod<$0.CreateFunctionRequest, $0.CreateFunctionResult>(
          '/cloudfront.CloudFrontService/CreateFunction',
          ($0.CreateFunctionRequest value) => value.writeToBuffer(),
          $0.CreateFunctionResult.fromBuffer);
  static final _$createInvalidation = $grpc.ClientMethod<
          $0.CreateInvalidationRequest, $0.CreateInvalidationResult>(
      '/cloudfront.CloudFrontService/CreateInvalidation',
      ($0.CreateInvalidationRequest value) => value.writeToBuffer(),
      $0.CreateInvalidationResult.fromBuffer);
  static final _$createInvalidationForDistributionTenant = $grpc.ClientMethod<
          $0.CreateInvalidationForDistributionTenantRequest,
          $0.CreateInvalidationForDistributionTenantResult>(
      '/cloudfront.CloudFrontService/CreateInvalidationForDistributionTenant',
      ($0.CreateInvalidationForDistributionTenantRequest value) =>
          value.writeToBuffer(),
      $0.CreateInvalidationForDistributionTenantResult.fromBuffer);
  static final _$createKeyGroup =
      $grpc.ClientMethod<$0.CreateKeyGroupRequest, $0.CreateKeyGroupResult>(
          '/cloudfront.CloudFrontService/CreateKeyGroup',
          ($0.CreateKeyGroupRequest value) => value.writeToBuffer(),
          $0.CreateKeyGroupResult.fromBuffer);
  static final _$createKeyValueStore = $grpc.ClientMethod<
          $0.CreateKeyValueStoreRequest, $0.CreateKeyValueStoreResult>(
      '/cloudfront.CloudFrontService/CreateKeyValueStore',
      ($0.CreateKeyValueStoreRequest value) => value.writeToBuffer(),
      $0.CreateKeyValueStoreResult.fromBuffer);
  static final _$createMonitoringSubscription = $grpc.ClientMethod<
          $0.CreateMonitoringSubscriptionRequest,
          $0.CreateMonitoringSubscriptionResult>(
      '/cloudfront.CloudFrontService/CreateMonitoringSubscription',
      ($0.CreateMonitoringSubscriptionRequest value) => value.writeToBuffer(),
      $0.CreateMonitoringSubscriptionResult.fromBuffer);
  static final _$createOriginAccessControl = $grpc.ClientMethod<
          $0.CreateOriginAccessControlRequest,
          $0.CreateOriginAccessControlResult>(
      '/cloudfront.CloudFrontService/CreateOriginAccessControl',
      ($0.CreateOriginAccessControlRequest value) => value.writeToBuffer(),
      $0.CreateOriginAccessControlResult.fromBuffer);
  static final _$createOriginRequestPolicy = $grpc.ClientMethod<
          $0.CreateOriginRequestPolicyRequest,
          $0.CreateOriginRequestPolicyResult>(
      '/cloudfront.CloudFrontService/CreateOriginRequestPolicy',
      ($0.CreateOriginRequestPolicyRequest value) => value.writeToBuffer(),
      $0.CreateOriginRequestPolicyResult.fromBuffer);
  static final _$createPublicKey =
      $grpc.ClientMethod<$0.CreatePublicKeyRequest, $0.CreatePublicKeyResult>(
          '/cloudfront.CloudFrontService/CreatePublicKey',
          ($0.CreatePublicKeyRequest value) => value.writeToBuffer(),
          $0.CreatePublicKeyResult.fromBuffer);
  static final _$createRealtimeLogConfig = $grpc.ClientMethod<
          $0.CreateRealtimeLogConfigRequest, $0.CreateRealtimeLogConfigResult>(
      '/cloudfront.CloudFrontService/CreateRealtimeLogConfig',
      ($0.CreateRealtimeLogConfigRequest value) => value.writeToBuffer(),
      $0.CreateRealtimeLogConfigResult.fromBuffer);
  static final _$createResponseHeadersPolicy = $grpc.ClientMethod<
          $0.CreateResponseHeadersPolicyRequest,
          $0.CreateResponseHeadersPolicyResult>(
      '/cloudfront.CloudFrontService/CreateResponseHeadersPolicy',
      ($0.CreateResponseHeadersPolicyRequest value) => value.writeToBuffer(),
      $0.CreateResponseHeadersPolicyResult.fromBuffer);
  static final _$createStreamingDistribution = $grpc.ClientMethod<
          $0.CreateStreamingDistributionRequest,
          $0.CreateStreamingDistributionResult>(
      '/cloudfront.CloudFrontService/CreateStreamingDistribution',
      ($0.CreateStreamingDistributionRequest value) => value.writeToBuffer(),
      $0.CreateStreamingDistributionResult.fromBuffer);
  static final _$createStreamingDistributionWithTags = $grpc.ClientMethod<
          $0.CreateStreamingDistributionWithTagsRequest,
          $0.CreateStreamingDistributionWithTagsResult>(
      '/cloudfront.CloudFrontService/CreateStreamingDistributionWithTags',
      ($0.CreateStreamingDistributionWithTagsRequest value) =>
          value.writeToBuffer(),
      $0.CreateStreamingDistributionWithTagsResult.fromBuffer);
  static final _$createTrustStore =
      $grpc.ClientMethod<$0.CreateTrustStoreRequest, $0.CreateTrustStoreResult>(
          '/cloudfront.CloudFrontService/CreateTrustStore',
          ($0.CreateTrustStoreRequest value) => value.writeToBuffer(),
          $0.CreateTrustStoreResult.fromBuffer);
  static final _$createVpcOrigin =
      $grpc.ClientMethod<$0.CreateVpcOriginRequest, $0.CreateVpcOriginResult>(
          '/cloudfront.CloudFrontService/CreateVpcOrigin',
          ($0.CreateVpcOriginRequest value) => value.writeToBuffer(),
          $0.CreateVpcOriginResult.fromBuffer);
  static final _$deleteAnycastIpList =
      $grpc.ClientMethod<$0.DeleteAnycastIpListRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteAnycastIpList',
          ($0.DeleteAnycastIpListRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteCachePolicy =
      $grpc.ClientMethod<$0.DeleteCachePolicyRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteCachePolicy',
          ($0.DeleteCachePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteCloudFrontOriginAccessIdentity = $grpc.ClientMethod<
          $0.DeleteCloudFrontOriginAccessIdentityRequest, $1.Empty>(
      '/cloudfront.CloudFrontService/DeleteCloudFrontOriginAccessIdentity',
      ($0.DeleteCloudFrontOriginAccessIdentityRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$deleteConnectionFunction =
      $grpc.ClientMethod<$0.DeleteConnectionFunctionRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteConnectionFunction',
          ($0.DeleteConnectionFunctionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteConnectionGroup =
      $grpc.ClientMethod<$0.DeleteConnectionGroupRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteConnectionGroup',
          ($0.DeleteConnectionGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteContinuousDeploymentPolicy =
      $grpc.ClientMethod<$0.DeleteContinuousDeploymentPolicyRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteContinuousDeploymentPolicy',
          ($0.DeleteContinuousDeploymentPolicyRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDistribution =
      $grpc.ClientMethod<$0.DeleteDistributionRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteDistribution',
          ($0.DeleteDistributionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDistributionTenant =
      $grpc.ClientMethod<$0.DeleteDistributionTenantRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteDistributionTenant',
          ($0.DeleteDistributionTenantRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteFieldLevelEncryptionConfig =
      $grpc.ClientMethod<$0.DeleteFieldLevelEncryptionConfigRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteFieldLevelEncryptionConfig',
          ($0.DeleteFieldLevelEncryptionConfigRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteFieldLevelEncryptionProfile =
      $grpc.ClientMethod<$0.DeleteFieldLevelEncryptionProfileRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteFieldLevelEncryptionProfile',
          ($0.DeleteFieldLevelEncryptionProfileRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteFunction =
      $grpc.ClientMethod<$0.DeleteFunctionRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteFunction',
          ($0.DeleteFunctionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteKeyGroup =
      $grpc.ClientMethod<$0.DeleteKeyGroupRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteKeyGroup',
          ($0.DeleteKeyGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteKeyValueStore =
      $grpc.ClientMethod<$0.DeleteKeyValueStoreRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteKeyValueStore',
          ($0.DeleteKeyValueStoreRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteMonitoringSubscription = $grpc.ClientMethod<
          $0.DeleteMonitoringSubscriptionRequest,
          $0.DeleteMonitoringSubscriptionResult>(
      '/cloudfront.CloudFrontService/DeleteMonitoringSubscription',
      ($0.DeleteMonitoringSubscriptionRequest value) => value.writeToBuffer(),
      $0.DeleteMonitoringSubscriptionResult.fromBuffer);
  static final _$deleteOriginAccessControl =
      $grpc.ClientMethod<$0.DeleteOriginAccessControlRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteOriginAccessControl',
          ($0.DeleteOriginAccessControlRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteOriginRequestPolicy =
      $grpc.ClientMethod<$0.DeleteOriginRequestPolicyRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteOriginRequestPolicy',
          ($0.DeleteOriginRequestPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deletePublicKey =
      $grpc.ClientMethod<$0.DeletePublicKeyRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeletePublicKey',
          ($0.DeletePublicKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteRealtimeLogConfig =
      $grpc.ClientMethod<$0.DeleteRealtimeLogConfigRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteRealtimeLogConfig',
          ($0.DeleteRealtimeLogConfigRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteResourcePolicy =
      $grpc.ClientMethod<$0.DeleteResourcePolicyRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteResourcePolicy',
          ($0.DeleteResourcePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteResponseHeadersPolicy =
      $grpc.ClientMethod<$0.DeleteResponseHeadersPolicyRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteResponseHeadersPolicy',
          ($0.DeleteResponseHeadersPolicyRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteStreamingDistribution =
      $grpc.ClientMethod<$0.DeleteStreamingDistributionRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteStreamingDistribution',
          ($0.DeleteStreamingDistributionRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteTrustStore =
      $grpc.ClientMethod<$0.DeleteTrustStoreRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/DeleteTrustStore',
          ($0.DeleteTrustStoreRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteVpcOrigin =
      $grpc.ClientMethod<$0.DeleteVpcOriginRequest, $0.DeleteVpcOriginResult>(
          '/cloudfront.CloudFrontService/DeleteVpcOrigin',
          ($0.DeleteVpcOriginRequest value) => value.writeToBuffer(),
          $0.DeleteVpcOriginResult.fromBuffer);
  static final _$describeConnectionFunction = $grpc.ClientMethod<
          $0.DescribeConnectionFunctionRequest,
          $0.DescribeConnectionFunctionResult>(
      '/cloudfront.CloudFrontService/DescribeConnectionFunction',
      ($0.DescribeConnectionFunctionRequest value) => value.writeToBuffer(),
      $0.DescribeConnectionFunctionResult.fromBuffer);
  static final _$describeFunction =
      $grpc.ClientMethod<$0.DescribeFunctionRequest, $0.DescribeFunctionResult>(
          '/cloudfront.CloudFrontService/DescribeFunction',
          ($0.DescribeFunctionRequest value) => value.writeToBuffer(),
          $0.DescribeFunctionResult.fromBuffer);
  static final _$describeKeyValueStore = $grpc.ClientMethod<
          $0.DescribeKeyValueStoreRequest, $0.DescribeKeyValueStoreResult>(
      '/cloudfront.CloudFrontService/DescribeKeyValueStore',
      ($0.DescribeKeyValueStoreRequest value) => value.writeToBuffer(),
      $0.DescribeKeyValueStoreResult.fromBuffer);
  static final _$disassociateDistributionTenantWebACL = $grpc.ClientMethod<
          $0.DisassociateDistributionTenantWebACLRequest,
          $0.DisassociateDistributionTenantWebACLResult>(
      '/cloudfront.CloudFrontService/DisassociateDistributionTenantWebACL',
      ($0.DisassociateDistributionTenantWebACLRequest value) =>
          value.writeToBuffer(),
      $0.DisassociateDistributionTenantWebACLResult.fromBuffer);
  static final _$disassociateDistributionWebACL = $grpc.ClientMethod<
          $0.DisassociateDistributionWebACLRequest,
          $0.DisassociateDistributionWebACLResult>(
      '/cloudfront.CloudFrontService/DisassociateDistributionWebACL',
      ($0.DisassociateDistributionWebACLRequest value) => value.writeToBuffer(),
      $0.DisassociateDistributionWebACLResult.fromBuffer);
  static final _$getAnycastIpList =
      $grpc.ClientMethod<$0.GetAnycastIpListRequest, $0.GetAnycastIpListResult>(
          '/cloudfront.CloudFrontService/GetAnycastIpList',
          ($0.GetAnycastIpListRequest value) => value.writeToBuffer(),
          $0.GetAnycastIpListResult.fromBuffer);
  static final _$getCachePolicy =
      $grpc.ClientMethod<$0.GetCachePolicyRequest, $0.GetCachePolicyResult>(
          '/cloudfront.CloudFrontService/GetCachePolicy',
          ($0.GetCachePolicyRequest value) => value.writeToBuffer(),
          $0.GetCachePolicyResult.fromBuffer);
  static final _$getCachePolicyConfig = $grpc.ClientMethod<
          $0.GetCachePolicyConfigRequest, $0.GetCachePolicyConfigResult>(
      '/cloudfront.CloudFrontService/GetCachePolicyConfig',
      ($0.GetCachePolicyConfigRequest value) => value.writeToBuffer(),
      $0.GetCachePolicyConfigResult.fromBuffer);
  static final _$getCloudFrontOriginAccessIdentity = $grpc.ClientMethod<
          $0.GetCloudFrontOriginAccessIdentityRequest,
          $0.GetCloudFrontOriginAccessIdentityResult>(
      '/cloudfront.CloudFrontService/GetCloudFrontOriginAccessIdentity',
      ($0.GetCloudFrontOriginAccessIdentityRequest value) =>
          value.writeToBuffer(),
      $0.GetCloudFrontOriginAccessIdentityResult.fromBuffer);
  static final _$getCloudFrontOriginAccessIdentityConfig = $grpc.ClientMethod<
          $0.GetCloudFrontOriginAccessIdentityConfigRequest,
          $0.GetCloudFrontOriginAccessIdentityConfigResult>(
      '/cloudfront.CloudFrontService/GetCloudFrontOriginAccessIdentityConfig',
      ($0.GetCloudFrontOriginAccessIdentityConfigRequest value) =>
          value.writeToBuffer(),
      $0.GetCloudFrontOriginAccessIdentityConfigResult.fromBuffer);
  static final _$getConnectionFunction = $grpc.ClientMethod<
          $0.GetConnectionFunctionRequest, $0.GetConnectionFunctionResult>(
      '/cloudfront.CloudFrontService/GetConnectionFunction',
      ($0.GetConnectionFunctionRequest value) => value.writeToBuffer(),
      $0.GetConnectionFunctionResult.fromBuffer);
  static final _$getConnectionGroup = $grpc.ClientMethod<
          $0.GetConnectionGroupRequest, $0.GetConnectionGroupResult>(
      '/cloudfront.CloudFrontService/GetConnectionGroup',
      ($0.GetConnectionGroupRequest value) => value.writeToBuffer(),
      $0.GetConnectionGroupResult.fromBuffer);
  static final _$getConnectionGroupByRoutingEndpoint = $grpc.ClientMethod<
          $0.GetConnectionGroupByRoutingEndpointRequest,
          $0.GetConnectionGroupByRoutingEndpointResult>(
      '/cloudfront.CloudFrontService/GetConnectionGroupByRoutingEndpoint',
      ($0.GetConnectionGroupByRoutingEndpointRequest value) =>
          value.writeToBuffer(),
      $0.GetConnectionGroupByRoutingEndpointResult.fromBuffer);
  static final _$getContinuousDeploymentPolicy = $grpc.ClientMethod<
          $0.GetContinuousDeploymentPolicyRequest,
          $0.GetContinuousDeploymentPolicyResult>(
      '/cloudfront.CloudFrontService/GetContinuousDeploymentPolicy',
      ($0.GetContinuousDeploymentPolicyRequest value) => value.writeToBuffer(),
      $0.GetContinuousDeploymentPolicyResult.fromBuffer);
  static final _$getContinuousDeploymentPolicyConfig = $grpc.ClientMethod<
          $0.GetContinuousDeploymentPolicyConfigRequest,
          $0.GetContinuousDeploymentPolicyConfigResult>(
      '/cloudfront.CloudFrontService/GetContinuousDeploymentPolicyConfig',
      ($0.GetContinuousDeploymentPolicyConfigRequest value) =>
          value.writeToBuffer(),
      $0.GetContinuousDeploymentPolicyConfigResult.fromBuffer);
  static final _$getDistribution =
      $grpc.ClientMethod<$0.GetDistributionRequest, $0.GetDistributionResult>(
          '/cloudfront.CloudFrontService/GetDistribution',
          ($0.GetDistributionRequest value) => value.writeToBuffer(),
          $0.GetDistributionResult.fromBuffer);
  static final _$getDistributionConfig = $grpc.ClientMethod<
          $0.GetDistributionConfigRequest, $0.GetDistributionConfigResult>(
      '/cloudfront.CloudFrontService/GetDistributionConfig',
      ($0.GetDistributionConfigRequest value) => value.writeToBuffer(),
      $0.GetDistributionConfigResult.fromBuffer);
  static final _$getDistributionTenant = $grpc.ClientMethod<
          $0.GetDistributionTenantRequest, $0.GetDistributionTenantResult>(
      '/cloudfront.CloudFrontService/GetDistributionTenant',
      ($0.GetDistributionTenantRequest value) => value.writeToBuffer(),
      $0.GetDistributionTenantResult.fromBuffer);
  static final _$getDistributionTenantByDomain = $grpc.ClientMethod<
          $0.GetDistributionTenantByDomainRequest,
          $0.GetDistributionTenantByDomainResult>(
      '/cloudfront.CloudFrontService/GetDistributionTenantByDomain',
      ($0.GetDistributionTenantByDomainRequest value) => value.writeToBuffer(),
      $0.GetDistributionTenantByDomainResult.fromBuffer);
  static final _$getFieldLevelEncryption = $grpc.ClientMethod<
          $0.GetFieldLevelEncryptionRequest, $0.GetFieldLevelEncryptionResult>(
      '/cloudfront.CloudFrontService/GetFieldLevelEncryption',
      ($0.GetFieldLevelEncryptionRequest value) => value.writeToBuffer(),
      $0.GetFieldLevelEncryptionResult.fromBuffer);
  static final _$getFieldLevelEncryptionConfig = $grpc.ClientMethod<
          $0.GetFieldLevelEncryptionConfigRequest,
          $0.GetFieldLevelEncryptionConfigResult>(
      '/cloudfront.CloudFrontService/GetFieldLevelEncryptionConfig',
      ($0.GetFieldLevelEncryptionConfigRequest value) => value.writeToBuffer(),
      $0.GetFieldLevelEncryptionConfigResult.fromBuffer);
  static final _$getFieldLevelEncryptionProfile = $grpc.ClientMethod<
          $0.GetFieldLevelEncryptionProfileRequest,
          $0.GetFieldLevelEncryptionProfileResult>(
      '/cloudfront.CloudFrontService/GetFieldLevelEncryptionProfile',
      ($0.GetFieldLevelEncryptionProfileRequest value) => value.writeToBuffer(),
      $0.GetFieldLevelEncryptionProfileResult.fromBuffer);
  static final _$getFieldLevelEncryptionProfileConfig = $grpc.ClientMethod<
          $0.GetFieldLevelEncryptionProfileConfigRequest,
          $0.GetFieldLevelEncryptionProfileConfigResult>(
      '/cloudfront.CloudFrontService/GetFieldLevelEncryptionProfileConfig',
      ($0.GetFieldLevelEncryptionProfileConfigRequest value) =>
          value.writeToBuffer(),
      $0.GetFieldLevelEncryptionProfileConfigResult.fromBuffer);
  static final _$getFunction =
      $grpc.ClientMethod<$0.GetFunctionRequest, $0.GetFunctionResult>(
          '/cloudfront.CloudFrontService/GetFunction',
          ($0.GetFunctionRequest value) => value.writeToBuffer(),
          $0.GetFunctionResult.fromBuffer);
  static final _$getInvalidation =
      $grpc.ClientMethod<$0.GetInvalidationRequest, $0.GetInvalidationResult>(
          '/cloudfront.CloudFrontService/GetInvalidation',
          ($0.GetInvalidationRequest value) => value.writeToBuffer(),
          $0.GetInvalidationResult.fromBuffer);
  static final _$getInvalidationForDistributionTenant = $grpc.ClientMethod<
          $0.GetInvalidationForDistributionTenantRequest,
          $0.GetInvalidationForDistributionTenantResult>(
      '/cloudfront.CloudFrontService/GetInvalidationForDistributionTenant',
      ($0.GetInvalidationForDistributionTenantRequest value) =>
          value.writeToBuffer(),
      $0.GetInvalidationForDistributionTenantResult.fromBuffer);
  static final _$getKeyGroup =
      $grpc.ClientMethod<$0.GetKeyGroupRequest, $0.GetKeyGroupResult>(
          '/cloudfront.CloudFrontService/GetKeyGroup',
          ($0.GetKeyGroupRequest value) => value.writeToBuffer(),
          $0.GetKeyGroupResult.fromBuffer);
  static final _$getKeyGroupConfig = $grpc.ClientMethod<
          $0.GetKeyGroupConfigRequest, $0.GetKeyGroupConfigResult>(
      '/cloudfront.CloudFrontService/GetKeyGroupConfig',
      ($0.GetKeyGroupConfigRequest value) => value.writeToBuffer(),
      $0.GetKeyGroupConfigResult.fromBuffer);
  static final _$getManagedCertificateDetails = $grpc.ClientMethod<
          $0.GetManagedCertificateDetailsRequest,
          $0.GetManagedCertificateDetailsResult>(
      '/cloudfront.CloudFrontService/GetManagedCertificateDetails',
      ($0.GetManagedCertificateDetailsRequest value) => value.writeToBuffer(),
      $0.GetManagedCertificateDetailsResult.fromBuffer);
  static final _$getMonitoringSubscription = $grpc.ClientMethod<
          $0.GetMonitoringSubscriptionRequest,
          $0.GetMonitoringSubscriptionResult>(
      '/cloudfront.CloudFrontService/GetMonitoringSubscription',
      ($0.GetMonitoringSubscriptionRequest value) => value.writeToBuffer(),
      $0.GetMonitoringSubscriptionResult.fromBuffer);
  static final _$getOriginAccessControl = $grpc.ClientMethod<
          $0.GetOriginAccessControlRequest, $0.GetOriginAccessControlResult>(
      '/cloudfront.CloudFrontService/GetOriginAccessControl',
      ($0.GetOriginAccessControlRequest value) => value.writeToBuffer(),
      $0.GetOriginAccessControlResult.fromBuffer);
  static final _$getOriginAccessControlConfig = $grpc.ClientMethod<
          $0.GetOriginAccessControlConfigRequest,
          $0.GetOriginAccessControlConfigResult>(
      '/cloudfront.CloudFrontService/GetOriginAccessControlConfig',
      ($0.GetOriginAccessControlConfigRequest value) => value.writeToBuffer(),
      $0.GetOriginAccessControlConfigResult.fromBuffer);
  static final _$getOriginRequestPolicy = $grpc.ClientMethod<
          $0.GetOriginRequestPolicyRequest, $0.GetOriginRequestPolicyResult>(
      '/cloudfront.CloudFrontService/GetOriginRequestPolicy',
      ($0.GetOriginRequestPolicyRequest value) => value.writeToBuffer(),
      $0.GetOriginRequestPolicyResult.fromBuffer);
  static final _$getOriginRequestPolicyConfig = $grpc.ClientMethod<
          $0.GetOriginRequestPolicyConfigRequest,
          $0.GetOriginRequestPolicyConfigResult>(
      '/cloudfront.CloudFrontService/GetOriginRequestPolicyConfig',
      ($0.GetOriginRequestPolicyConfigRequest value) => value.writeToBuffer(),
      $0.GetOriginRequestPolicyConfigResult.fromBuffer);
  static final _$getPublicKey =
      $grpc.ClientMethod<$0.GetPublicKeyRequest, $0.GetPublicKeyResult>(
          '/cloudfront.CloudFrontService/GetPublicKey',
          ($0.GetPublicKeyRequest value) => value.writeToBuffer(),
          $0.GetPublicKeyResult.fromBuffer);
  static final _$getPublicKeyConfig = $grpc.ClientMethod<
          $0.GetPublicKeyConfigRequest, $0.GetPublicKeyConfigResult>(
      '/cloudfront.CloudFrontService/GetPublicKeyConfig',
      ($0.GetPublicKeyConfigRequest value) => value.writeToBuffer(),
      $0.GetPublicKeyConfigResult.fromBuffer);
  static final _$getRealtimeLogConfig = $grpc.ClientMethod<
          $0.GetRealtimeLogConfigRequest, $0.GetRealtimeLogConfigResult>(
      '/cloudfront.CloudFrontService/GetRealtimeLogConfig',
      ($0.GetRealtimeLogConfigRequest value) => value.writeToBuffer(),
      $0.GetRealtimeLogConfigResult.fromBuffer);
  static final _$getResourcePolicy = $grpc.ClientMethod<
          $0.GetResourcePolicyRequest, $0.GetResourcePolicyResult>(
      '/cloudfront.CloudFrontService/GetResourcePolicy',
      ($0.GetResourcePolicyRequest value) => value.writeToBuffer(),
      $0.GetResourcePolicyResult.fromBuffer);
  static final _$getResponseHeadersPolicy = $grpc.ClientMethod<
          $0.GetResponseHeadersPolicyRequest,
          $0.GetResponseHeadersPolicyResult>(
      '/cloudfront.CloudFrontService/GetResponseHeadersPolicy',
      ($0.GetResponseHeadersPolicyRequest value) => value.writeToBuffer(),
      $0.GetResponseHeadersPolicyResult.fromBuffer);
  static final _$getResponseHeadersPolicyConfig = $grpc.ClientMethod<
          $0.GetResponseHeadersPolicyConfigRequest,
          $0.GetResponseHeadersPolicyConfigResult>(
      '/cloudfront.CloudFrontService/GetResponseHeadersPolicyConfig',
      ($0.GetResponseHeadersPolicyConfigRequest value) => value.writeToBuffer(),
      $0.GetResponseHeadersPolicyConfigResult.fromBuffer);
  static final _$getStreamingDistribution = $grpc.ClientMethod<
          $0.GetStreamingDistributionRequest,
          $0.GetStreamingDistributionResult>(
      '/cloudfront.CloudFrontService/GetStreamingDistribution',
      ($0.GetStreamingDistributionRequest value) => value.writeToBuffer(),
      $0.GetStreamingDistributionResult.fromBuffer);
  static final _$getStreamingDistributionConfig = $grpc.ClientMethod<
          $0.GetStreamingDistributionConfigRequest,
          $0.GetStreamingDistributionConfigResult>(
      '/cloudfront.CloudFrontService/GetStreamingDistributionConfig',
      ($0.GetStreamingDistributionConfigRequest value) => value.writeToBuffer(),
      $0.GetStreamingDistributionConfigResult.fromBuffer);
  static final _$getTrustStore =
      $grpc.ClientMethod<$0.GetTrustStoreRequest, $0.GetTrustStoreResult>(
          '/cloudfront.CloudFrontService/GetTrustStore',
          ($0.GetTrustStoreRequest value) => value.writeToBuffer(),
          $0.GetTrustStoreResult.fromBuffer);
  static final _$getVpcOrigin =
      $grpc.ClientMethod<$0.GetVpcOriginRequest, $0.GetVpcOriginResult>(
          '/cloudfront.CloudFrontService/GetVpcOrigin',
          ($0.GetVpcOriginRequest value) => value.writeToBuffer(),
          $0.GetVpcOriginResult.fromBuffer);
  static final _$listAnycastIpLists = $grpc.ClientMethod<
          $0.ListAnycastIpListsRequest, $0.ListAnycastIpListsResult>(
      '/cloudfront.CloudFrontService/ListAnycastIpLists',
      ($0.ListAnycastIpListsRequest value) => value.writeToBuffer(),
      $0.ListAnycastIpListsResult.fromBuffer);
  static final _$listCachePolicies = $grpc.ClientMethod<
          $0.ListCachePoliciesRequest, $0.ListCachePoliciesResult>(
      '/cloudfront.CloudFrontService/ListCachePolicies',
      ($0.ListCachePoliciesRequest value) => value.writeToBuffer(),
      $0.ListCachePoliciesResult.fromBuffer);
  static final _$listCloudFrontOriginAccessIdentities = $grpc.ClientMethod<
          $0.ListCloudFrontOriginAccessIdentitiesRequest,
          $0.ListCloudFrontOriginAccessIdentitiesResult>(
      '/cloudfront.CloudFrontService/ListCloudFrontOriginAccessIdentities',
      ($0.ListCloudFrontOriginAccessIdentitiesRequest value) =>
          value.writeToBuffer(),
      $0.ListCloudFrontOriginAccessIdentitiesResult.fromBuffer);
  static final _$listConflictingAliases = $grpc.ClientMethod<
          $0.ListConflictingAliasesRequest, $0.ListConflictingAliasesResult>(
      '/cloudfront.CloudFrontService/ListConflictingAliases',
      ($0.ListConflictingAliasesRequest value) => value.writeToBuffer(),
      $0.ListConflictingAliasesResult.fromBuffer);
  static final _$listConnectionFunctions = $grpc.ClientMethod<
          $0.ListConnectionFunctionsRequest, $0.ListConnectionFunctionsResult>(
      '/cloudfront.CloudFrontService/ListConnectionFunctions',
      ($0.ListConnectionFunctionsRequest value) => value.writeToBuffer(),
      $0.ListConnectionFunctionsResult.fromBuffer);
  static final _$listConnectionGroups = $grpc.ClientMethod<
          $0.ListConnectionGroupsRequest, $0.ListConnectionGroupsResult>(
      '/cloudfront.CloudFrontService/ListConnectionGroups',
      ($0.ListConnectionGroupsRequest value) => value.writeToBuffer(),
      $0.ListConnectionGroupsResult.fromBuffer);
  static final _$listContinuousDeploymentPolicies = $grpc.ClientMethod<
          $0.ListContinuousDeploymentPoliciesRequest,
          $0.ListContinuousDeploymentPoliciesResult>(
      '/cloudfront.CloudFrontService/ListContinuousDeploymentPolicies',
      ($0.ListContinuousDeploymentPoliciesRequest value) =>
          value.writeToBuffer(),
      $0.ListContinuousDeploymentPoliciesResult.fromBuffer);
  static final _$listDistributions = $grpc.ClientMethod<
          $0.ListDistributionsRequest, $0.ListDistributionsResult>(
      '/cloudfront.CloudFrontService/ListDistributions',
      ($0.ListDistributionsRequest value) => value.writeToBuffer(),
      $0.ListDistributionsResult.fromBuffer);
  static final _$listDistributionsByAnycastIpListId = $grpc.ClientMethod<
          $0.ListDistributionsByAnycastIpListIdRequest,
          $0.ListDistributionsByAnycastIpListIdResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByAnycastIpListId',
      ($0.ListDistributionsByAnycastIpListIdRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionsByAnycastIpListIdResult.fromBuffer);
  static final _$listDistributionsByCachePolicyId = $grpc.ClientMethod<
          $0.ListDistributionsByCachePolicyIdRequest,
          $0.ListDistributionsByCachePolicyIdResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByCachePolicyId',
      ($0.ListDistributionsByCachePolicyIdRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionsByCachePolicyIdResult.fromBuffer);
  static final _$listDistributionsByConnectionFunction = $grpc.ClientMethod<
          $0.ListDistributionsByConnectionFunctionRequest,
          $0.ListDistributionsByConnectionFunctionResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByConnectionFunction',
      ($0.ListDistributionsByConnectionFunctionRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionsByConnectionFunctionResult.fromBuffer);
  static final _$listDistributionsByConnectionMode = $grpc.ClientMethod<
          $0.ListDistributionsByConnectionModeRequest,
          $0.ListDistributionsByConnectionModeResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByConnectionMode',
      ($0.ListDistributionsByConnectionModeRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionsByConnectionModeResult.fromBuffer);
  static final _$listDistributionsByKeyGroup = $grpc.ClientMethod<
          $0.ListDistributionsByKeyGroupRequest,
          $0.ListDistributionsByKeyGroupResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByKeyGroup',
      ($0.ListDistributionsByKeyGroupRequest value) => value.writeToBuffer(),
      $0.ListDistributionsByKeyGroupResult.fromBuffer);
  static final _$listDistributionsByOriginRequestPolicyId = $grpc.ClientMethod<
          $0.ListDistributionsByOriginRequestPolicyIdRequest,
          $0.ListDistributionsByOriginRequestPolicyIdResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByOriginRequestPolicyId',
      ($0.ListDistributionsByOriginRequestPolicyIdRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionsByOriginRequestPolicyIdResult.fromBuffer);
  static final _$listDistributionsByOwnedResource = $grpc.ClientMethod<
          $0.ListDistributionsByOwnedResourceRequest,
          $0.ListDistributionsByOwnedResourceResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByOwnedResource',
      ($0.ListDistributionsByOwnedResourceRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionsByOwnedResourceResult.fromBuffer);
  static final _$listDistributionsByRealtimeLogConfig = $grpc.ClientMethod<
          $0.ListDistributionsByRealtimeLogConfigRequest,
          $0.ListDistributionsByRealtimeLogConfigResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByRealtimeLogConfig',
      ($0.ListDistributionsByRealtimeLogConfigRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionsByRealtimeLogConfigResult.fromBuffer);
  static final _$listDistributionsByResponseHeadersPolicyId = $grpc.ClientMethod<
          $0.ListDistributionsByResponseHeadersPolicyIdRequest,
          $0.ListDistributionsByResponseHeadersPolicyIdResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByResponseHeadersPolicyId',
      ($0.ListDistributionsByResponseHeadersPolicyIdRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionsByResponseHeadersPolicyIdResult.fromBuffer);
  static final _$listDistributionsByTrustStore = $grpc.ClientMethod<
          $0.ListDistributionsByTrustStoreRequest,
          $0.ListDistributionsByTrustStoreResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByTrustStore',
      ($0.ListDistributionsByTrustStoreRequest value) => value.writeToBuffer(),
      $0.ListDistributionsByTrustStoreResult.fromBuffer);
  static final _$listDistributionsByVpcOriginId = $grpc.ClientMethod<
          $0.ListDistributionsByVpcOriginIdRequest,
          $0.ListDistributionsByVpcOriginIdResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByVpcOriginId',
      ($0.ListDistributionsByVpcOriginIdRequest value) => value.writeToBuffer(),
      $0.ListDistributionsByVpcOriginIdResult.fromBuffer);
  static final _$listDistributionsByWebACLId = $grpc.ClientMethod<
          $0.ListDistributionsByWebACLIdRequest,
          $0.ListDistributionsByWebACLIdResult>(
      '/cloudfront.CloudFrontService/ListDistributionsByWebACLId',
      ($0.ListDistributionsByWebACLIdRequest value) => value.writeToBuffer(),
      $0.ListDistributionsByWebACLIdResult.fromBuffer);
  static final _$listDistributionTenants = $grpc.ClientMethod<
          $0.ListDistributionTenantsRequest, $0.ListDistributionTenantsResult>(
      '/cloudfront.CloudFrontService/ListDistributionTenants',
      ($0.ListDistributionTenantsRequest value) => value.writeToBuffer(),
      $0.ListDistributionTenantsResult.fromBuffer);
  static final _$listDistributionTenantsByCustomization = $grpc.ClientMethod<
          $0.ListDistributionTenantsByCustomizationRequest,
          $0.ListDistributionTenantsByCustomizationResult>(
      '/cloudfront.CloudFrontService/ListDistributionTenantsByCustomization',
      ($0.ListDistributionTenantsByCustomizationRequest value) =>
          value.writeToBuffer(),
      $0.ListDistributionTenantsByCustomizationResult.fromBuffer);
  static final _$listDomainConflicts = $grpc.ClientMethod<
          $0.ListDomainConflictsRequest, $0.ListDomainConflictsResult>(
      '/cloudfront.CloudFrontService/ListDomainConflicts',
      ($0.ListDomainConflictsRequest value) => value.writeToBuffer(),
      $0.ListDomainConflictsResult.fromBuffer);
  static final _$listFieldLevelEncryptionConfigs = $grpc.ClientMethod<
          $0.ListFieldLevelEncryptionConfigsRequest,
          $0.ListFieldLevelEncryptionConfigsResult>(
      '/cloudfront.CloudFrontService/ListFieldLevelEncryptionConfigs',
      ($0.ListFieldLevelEncryptionConfigsRequest value) =>
          value.writeToBuffer(),
      $0.ListFieldLevelEncryptionConfigsResult.fromBuffer);
  static final _$listFieldLevelEncryptionProfiles = $grpc.ClientMethod<
          $0.ListFieldLevelEncryptionProfilesRequest,
          $0.ListFieldLevelEncryptionProfilesResult>(
      '/cloudfront.CloudFrontService/ListFieldLevelEncryptionProfiles',
      ($0.ListFieldLevelEncryptionProfilesRequest value) =>
          value.writeToBuffer(),
      $0.ListFieldLevelEncryptionProfilesResult.fromBuffer);
  static final _$listFunctions =
      $grpc.ClientMethod<$0.ListFunctionsRequest, $0.ListFunctionsResult>(
          '/cloudfront.CloudFrontService/ListFunctions',
          ($0.ListFunctionsRequest value) => value.writeToBuffer(),
          $0.ListFunctionsResult.fromBuffer);
  static final _$listInvalidations = $grpc.ClientMethod<
          $0.ListInvalidationsRequest, $0.ListInvalidationsResult>(
      '/cloudfront.CloudFrontService/ListInvalidations',
      ($0.ListInvalidationsRequest value) => value.writeToBuffer(),
      $0.ListInvalidationsResult.fromBuffer);
  static final _$listInvalidationsForDistributionTenant = $grpc.ClientMethod<
          $0.ListInvalidationsForDistributionTenantRequest,
          $0.ListInvalidationsForDistributionTenantResult>(
      '/cloudfront.CloudFrontService/ListInvalidationsForDistributionTenant',
      ($0.ListInvalidationsForDistributionTenantRequest value) =>
          value.writeToBuffer(),
      $0.ListInvalidationsForDistributionTenantResult.fromBuffer);
  static final _$listKeyGroups =
      $grpc.ClientMethod<$0.ListKeyGroupsRequest, $0.ListKeyGroupsResult>(
          '/cloudfront.CloudFrontService/ListKeyGroups',
          ($0.ListKeyGroupsRequest value) => value.writeToBuffer(),
          $0.ListKeyGroupsResult.fromBuffer);
  static final _$listKeyValueStores = $grpc.ClientMethod<
          $0.ListKeyValueStoresRequest, $0.ListKeyValueStoresResult>(
      '/cloudfront.CloudFrontService/ListKeyValueStores',
      ($0.ListKeyValueStoresRequest value) => value.writeToBuffer(),
      $0.ListKeyValueStoresResult.fromBuffer);
  static final _$listOriginAccessControls = $grpc.ClientMethod<
          $0.ListOriginAccessControlsRequest,
          $0.ListOriginAccessControlsResult>(
      '/cloudfront.CloudFrontService/ListOriginAccessControls',
      ($0.ListOriginAccessControlsRequest value) => value.writeToBuffer(),
      $0.ListOriginAccessControlsResult.fromBuffer);
  static final _$listOriginRequestPolicies = $grpc.ClientMethod<
          $0.ListOriginRequestPoliciesRequest,
          $0.ListOriginRequestPoliciesResult>(
      '/cloudfront.CloudFrontService/ListOriginRequestPolicies',
      ($0.ListOriginRequestPoliciesRequest value) => value.writeToBuffer(),
      $0.ListOriginRequestPoliciesResult.fromBuffer);
  static final _$listPublicKeys =
      $grpc.ClientMethod<$0.ListPublicKeysRequest, $0.ListPublicKeysResult>(
          '/cloudfront.CloudFrontService/ListPublicKeys',
          ($0.ListPublicKeysRequest value) => value.writeToBuffer(),
          $0.ListPublicKeysResult.fromBuffer);
  static final _$listRealtimeLogConfigs = $grpc.ClientMethod<
          $0.ListRealtimeLogConfigsRequest, $0.ListRealtimeLogConfigsResult>(
      '/cloudfront.CloudFrontService/ListRealtimeLogConfigs',
      ($0.ListRealtimeLogConfigsRequest value) => value.writeToBuffer(),
      $0.ListRealtimeLogConfigsResult.fromBuffer);
  static final _$listResponseHeadersPolicies = $grpc.ClientMethod<
          $0.ListResponseHeadersPoliciesRequest,
          $0.ListResponseHeadersPoliciesResult>(
      '/cloudfront.CloudFrontService/ListResponseHeadersPolicies',
      ($0.ListResponseHeadersPoliciesRequest value) => value.writeToBuffer(),
      $0.ListResponseHeadersPoliciesResult.fromBuffer);
  static final _$listStreamingDistributions = $grpc.ClientMethod<
          $0.ListStreamingDistributionsRequest,
          $0.ListStreamingDistributionsResult>(
      '/cloudfront.CloudFrontService/ListStreamingDistributions',
      ($0.ListStreamingDistributionsRequest value) => value.writeToBuffer(),
      $0.ListStreamingDistributionsResult.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResult>(
      '/cloudfront.CloudFrontService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResult.fromBuffer);
  static final _$listTrustStores =
      $grpc.ClientMethod<$0.ListTrustStoresRequest, $0.ListTrustStoresResult>(
          '/cloudfront.CloudFrontService/ListTrustStores',
          ($0.ListTrustStoresRequest value) => value.writeToBuffer(),
          $0.ListTrustStoresResult.fromBuffer);
  static final _$listVpcOrigins =
      $grpc.ClientMethod<$0.ListVpcOriginsRequest, $0.ListVpcOriginsResult>(
          '/cloudfront.CloudFrontService/ListVpcOrigins',
          ($0.ListVpcOriginsRequest value) => value.writeToBuffer(),
          $0.ListVpcOriginsResult.fromBuffer);
  static final _$publishConnectionFunction = $grpc.ClientMethod<
          $0.PublishConnectionFunctionRequest,
          $0.PublishConnectionFunctionResult>(
      '/cloudfront.CloudFrontService/PublishConnectionFunction',
      ($0.PublishConnectionFunctionRequest value) => value.writeToBuffer(),
      $0.PublishConnectionFunctionResult.fromBuffer);
  static final _$publishFunction =
      $grpc.ClientMethod<$0.PublishFunctionRequest, $0.PublishFunctionResult>(
          '/cloudfront.CloudFrontService/PublishFunction',
          ($0.PublishFunctionRequest value) => value.writeToBuffer(),
          $0.PublishFunctionResult.fromBuffer);
  static final _$putResourcePolicy = $grpc.ClientMethod<
          $0.PutResourcePolicyRequest, $0.PutResourcePolicyResult>(
      '/cloudfront.CloudFrontService/PutResourcePolicy',
      ($0.PutResourcePolicyRequest value) => value.writeToBuffer(),
      $0.PutResourcePolicyResult.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$testConnectionFunction = $grpc.ClientMethod<
          $0.TestConnectionFunctionRequest, $0.TestConnectionFunctionResult>(
      '/cloudfront.CloudFrontService/TestConnectionFunction',
      ($0.TestConnectionFunctionRequest value) => value.writeToBuffer(),
      $0.TestConnectionFunctionResult.fromBuffer);
  static final _$testFunction =
      $grpc.ClientMethod<$0.TestFunctionRequest, $0.TestFunctionResult>(
          '/cloudfront.CloudFrontService/TestFunction',
          ($0.TestFunctionRequest value) => value.writeToBuffer(),
          $0.TestFunctionResult.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $1.Empty>(
          '/cloudfront.CloudFrontService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateAnycastIpList = $grpc.ClientMethod<
          $0.UpdateAnycastIpListRequest, $0.UpdateAnycastIpListResult>(
      '/cloudfront.CloudFrontService/UpdateAnycastIpList',
      ($0.UpdateAnycastIpListRequest value) => value.writeToBuffer(),
      $0.UpdateAnycastIpListResult.fromBuffer);
  static final _$updateCachePolicy = $grpc.ClientMethod<
          $0.UpdateCachePolicyRequest, $0.UpdateCachePolicyResult>(
      '/cloudfront.CloudFrontService/UpdateCachePolicy',
      ($0.UpdateCachePolicyRequest value) => value.writeToBuffer(),
      $0.UpdateCachePolicyResult.fromBuffer);
  static final _$updateCloudFrontOriginAccessIdentity = $grpc.ClientMethod<
          $0.UpdateCloudFrontOriginAccessIdentityRequest,
          $0.UpdateCloudFrontOriginAccessIdentityResult>(
      '/cloudfront.CloudFrontService/UpdateCloudFrontOriginAccessIdentity',
      ($0.UpdateCloudFrontOriginAccessIdentityRequest value) =>
          value.writeToBuffer(),
      $0.UpdateCloudFrontOriginAccessIdentityResult.fromBuffer);
  static final _$updateConnectionFunction = $grpc.ClientMethod<
          $0.UpdateConnectionFunctionRequest,
          $0.UpdateConnectionFunctionResult>(
      '/cloudfront.CloudFrontService/UpdateConnectionFunction',
      ($0.UpdateConnectionFunctionRequest value) => value.writeToBuffer(),
      $0.UpdateConnectionFunctionResult.fromBuffer);
  static final _$updateConnectionGroup = $grpc.ClientMethod<
          $0.UpdateConnectionGroupRequest, $0.UpdateConnectionGroupResult>(
      '/cloudfront.CloudFrontService/UpdateConnectionGroup',
      ($0.UpdateConnectionGroupRequest value) => value.writeToBuffer(),
      $0.UpdateConnectionGroupResult.fromBuffer);
  static final _$updateContinuousDeploymentPolicy = $grpc.ClientMethod<
          $0.UpdateContinuousDeploymentPolicyRequest,
          $0.UpdateContinuousDeploymentPolicyResult>(
      '/cloudfront.CloudFrontService/UpdateContinuousDeploymentPolicy',
      ($0.UpdateContinuousDeploymentPolicyRequest value) =>
          value.writeToBuffer(),
      $0.UpdateContinuousDeploymentPolicyResult.fromBuffer);
  static final _$updateDistribution = $grpc.ClientMethod<
          $0.UpdateDistributionRequest, $0.UpdateDistributionResult>(
      '/cloudfront.CloudFrontService/UpdateDistribution',
      ($0.UpdateDistributionRequest value) => value.writeToBuffer(),
      $0.UpdateDistributionResult.fromBuffer);
  static final _$updateDistributionTenant = $grpc.ClientMethod<
          $0.UpdateDistributionTenantRequest,
          $0.UpdateDistributionTenantResult>(
      '/cloudfront.CloudFrontService/UpdateDistributionTenant',
      ($0.UpdateDistributionTenantRequest value) => value.writeToBuffer(),
      $0.UpdateDistributionTenantResult.fromBuffer);
  static final _$updateDistributionWithStagingConfig = $grpc.ClientMethod<
          $0.UpdateDistributionWithStagingConfigRequest,
          $0.UpdateDistributionWithStagingConfigResult>(
      '/cloudfront.CloudFrontService/UpdateDistributionWithStagingConfig',
      ($0.UpdateDistributionWithStagingConfigRequest value) =>
          value.writeToBuffer(),
      $0.UpdateDistributionWithStagingConfigResult.fromBuffer);
  static final _$updateDomainAssociation = $grpc.ClientMethod<
          $0.UpdateDomainAssociationRequest, $0.UpdateDomainAssociationResult>(
      '/cloudfront.CloudFrontService/UpdateDomainAssociation',
      ($0.UpdateDomainAssociationRequest value) => value.writeToBuffer(),
      $0.UpdateDomainAssociationResult.fromBuffer);
  static final _$updateFieldLevelEncryptionConfig = $grpc.ClientMethod<
          $0.UpdateFieldLevelEncryptionConfigRequest,
          $0.UpdateFieldLevelEncryptionConfigResult>(
      '/cloudfront.CloudFrontService/UpdateFieldLevelEncryptionConfig',
      ($0.UpdateFieldLevelEncryptionConfigRequest value) =>
          value.writeToBuffer(),
      $0.UpdateFieldLevelEncryptionConfigResult.fromBuffer);
  static final _$updateFieldLevelEncryptionProfile = $grpc.ClientMethod<
          $0.UpdateFieldLevelEncryptionProfileRequest,
          $0.UpdateFieldLevelEncryptionProfileResult>(
      '/cloudfront.CloudFrontService/UpdateFieldLevelEncryptionProfile',
      ($0.UpdateFieldLevelEncryptionProfileRequest value) =>
          value.writeToBuffer(),
      $0.UpdateFieldLevelEncryptionProfileResult.fromBuffer);
  static final _$updateFunction =
      $grpc.ClientMethod<$0.UpdateFunctionRequest, $0.UpdateFunctionResult>(
          '/cloudfront.CloudFrontService/UpdateFunction',
          ($0.UpdateFunctionRequest value) => value.writeToBuffer(),
          $0.UpdateFunctionResult.fromBuffer);
  static final _$updateKeyGroup =
      $grpc.ClientMethod<$0.UpdateKeyGroupRequest, $0.UpdateKeyGroupResult>(
          '/cloudfront.CloudFrontService/UpdateKeyGroup',
          ($0.UpdateKeyGroupRequest value) => value.writeToBuffer(),
          $0.UpdateKeyGroupResult.fromBuffer);
  static final _$updateKeyValueStore = $grpc.ClientMethod<
          $0.UpdateKeyValueStoreRequest, $0.UpdateKeyValueStoreResult>(
      '/cloudfront.CloudFrontService/UpdateKeyValueStore',
      ($0.UpdateKeyValueStoreRequest value) => value.writeToBuffer(),
      $0.UpdateKeyValueStoreResult.fromBuffer);
  static final _$updateOriginAccessControl = $grpc.ClientMethod<
          $0.UpdateOriginAccessControlRequest,
          $0.UpdateOriginAccessControlResult>(
      '/cloudfront.CloudFrontService/UpdateOriginAccessControl',
      ($0.UpdateOriginAccessControlRequest value) => value.writeToBuffer(),
      $0.UpdateOriginAccessControlResult.fromBuffer);
  static final _$updateOriginRequestPolicy = $grpc.ClientMethod<
          $0.UpdateOriginRequestPolicyRequest,
          $0.UpdateOriginRequestPolicyResult>(
      '/cloudfront.CloudFrontService/UpdateOriginRequestPolicy',
      ($0.UpdateOriginRequestPolicyRequest value) => value.writeToBuffer(),
      $0.UpdateOriginRequestPolicyResult.fromBuffer);
  static final _$updatePublicKey =
      $grpc.ClientMethod<$0.UpdatePublicKeyRequest, $0.UpdatePublicKeyResult>(
          '/cloudfront.CloudFrontService/UpdatePublicKey',
          ($0.UpdatePublicKeyRequest value) => value.writeToBuffer(),
          $0.UpdatePublicKeyResult.fromBuffer);
  static final _$updateRealtimeLogConfig = $grpc.ClientMethod<
          $0.UpdateRealtimeLogConfigRequest, $0.UpdateRealtimeLogConfigResult>(
      '/cloudfront.CloudFrontService/UpdateRealtimeLogConfig',
      ($0.UpdateRealtimeLogConfigRequest value) => value.writeToBuffer(),
      $0.UpdateRealtimeLogConfigResult.fromBuffer);
  static final _$updateResponseHeadersPolicy = $grpc.ClientMethod<
          $0.UpdateResponseHeadersPolicyRequest,
          $0.UpdateResponseHeadersPolicyResult>(
      '/cloudfront.CloudFrontService/UpdateResponseHeadersPolicy',
      ($0.UpdateResponseHeadersPolicyRequest value) => value.writeToBuffer(),
      $0.UpdateResponseHeadersPolicyResult.fromBuffer);
  static final _$updateStreamingDistribution = $grpc.ClientMethod<
          $0.UpdateStreamingDistributionRequest,
          $0.UpdateStreamingDistributionResult>(
      '/cloudfront.CloudFrontService/UpdateStreamingDistribution',
      ($0.UpdateStreamingDistributionRequest value) => value.writeToBuffer(),
      $0.UpdateStreamingDistributionResult.fromBuffer);
  static final _$updateTrustStore =
      $grpc.ClientMethod<$0.UpdateTrustStoreRequest, $0.UpdateTrustStoreResult>(
          '/cloudfront.CloudFrontService/UpdateTrustStore',
          ($0.UpdateTrustStoreRequest value) => value.writeToBuffer(),
          $0.UpdateTrustStoreResult.fromBuffer);
  static final _$updateVpcOrigin =
      $grpc.ClientMethod<$0.UpdateVpcOriginRequest, $0.UpdateVpcOriginResult>(
          '/cloudfront.CloudFrontService/UpdateVpcOrigin',
          ($0.UpdateVpcOriginRequest value) => value.writeToBuffer(),
          $0.UpdateVpcOriginResult.fromBuffer);
  static final _$verifyDnsConfiguration = $grpc.ClientMethod<
          $0.VerifyDnsConfigurationRequest, $0.VerifyDnsConfigurationResult>(
      '/cloudfront.CloudFrontService/VerifyDnsConfiguration',
      ($0.VerifyDnsConfigurationRequest value) => value.writeToBuffer(),
      $0.VerifyDnsConfigurationResult.fromBuffer);
}

@$pb.GrpcServiceName('cloudfront.CloudFrontService')
abstract class CloudFrontServiceBase extends $grpc.Service {
  $core.String get $name => 'cloudfront.CloudFrontService';

  CloudFrontServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AssociateAliasRequest, $1.Empty>(
        'AssociateAlias',
        associateAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateAliasRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssociateDistributionTenantWebACLRequest,
            $0.AssociateDistributionTenantWebACLResult>(
        'AssociateDistributionTenantWebACL',
        associateDistributionTenantWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateDistributionTenantWebACLRequest.fromBuffer(value),
        ($0.AssociateDistributionTenantWebACLResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssociateDistributionWebACLRequest,
            $0.AssociateDistributionWebACLResult>(
        'AssociateDistributionWebACL',
        associateDistributionWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateDistributionWebACLRequest.fromBuffer(value),
        ($0.AssociateDistributionWebACLResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CopyDistributionRequest,
            $0.CopyDistributionResult>(
        'CopyDistribution',
        copyDistribution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CopyDistributionRequest.fromBuffer(value),
        ($0.CopyDistributionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateAnycastIpListRequest,
            $0.CreateAnycastIpListResult>(
        'CreateAnycastIpList',
        createAnycastIpList_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateAnycastIpListRequest.fromBuffer(value),
        ($0.CreateAnycastIpListResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateCachePolicyRequest,
            $0.CreateCachePolicyResult>(
        'CreateCachePolicy',
        createCachePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCachePolicyRequest.fromBuffer(value),
        ($0.CreateCachePolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.CreateCloudFrontOriginAccessIdentityRequest,
            $0.CreateCloudFrontOriginAccessIdentityResult>(
        'CreateCloudFrontOriginAccessIdentity',
        createCloudFrontOriginAccessIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCloudFrontOriginAccessIdentityRequest.fromBuffer(value),
        ($0.CreateCloudFrontOriginAccessIdentityResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateConnectionFunctionRequest,
            $0.CreateConnectionFunctionResult>(
        'CreateConnectionFunction',
        createConnectionFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateConnectionFunctionRequest.fromBuffer(value),
        ($0.CreateConnectionFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateConnectionGroupRequest,
            $0.CreateConnectionGroupResult>(
        'CreateConnectionGroup',
        createConnectionGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateConnectionGroupRequest.fromBuffer(value),
        ($0.CreateConnectionGroupResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateContinuousDeploymentPolicyRequest,
            $0.CreateContinuousDeploymentPolicyResult>(
        'CreateContinuousDeploymentPolicy',
        createContinuousDeploymentPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateContinuousDeploymentPolicyRequest.fromBuffer(value),
        ($0.CreateContinuousDeploymentPolicyResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDistributionRequest,
            $0.CreateDistributionResult>(
        'CreateDistribution',
        createDistribution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDistributionRequest.fromBuffer(value),
        ($0.CreateDistributionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDistributionTenantRequest,
            $0.CreateDistributionTenantResult>(
        'CreateDistributionTenant',
        createDistributionTenant_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDistributionTenantRequest.fromBuffer(value),
        ($0.CreateDistributionTenantResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDistributionWithTagsRequest,
            $0.CreateDistributionWithTagsResult>(
        'CreateDistributionWithTags',
        createDistributionWithTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDistributionWithTagsRequest.fromBuffer(value),
        ($0.CreateDistributionWithTagsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateFieldLevelEncryptionConfigRequest,
            $0.CreateFieldLevelEncryptionConfigResult>(
        'CreateFieldLevelEncryptionConfig',
        createFieldLevelEncryptionConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateFieldLevelEncryptionConfigRequest.fromBuffer(value),
        ($0.CreateFieldLevelEncryptionConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateFieldLevelEncryptionProfileRequest,
            $0.CreateFieldLevelEncryptionProfileResult>(
        'CreateFieldLevelEncryptionProfile',
        createFieldLevelEncryptionProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateFieldLevelEncryptionProfileRequest.fromBuffer(value),
        ($0.CreateFieldLevelEncryptionProfileResult value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateFunctionRequest, $0.CreateFunctionResult>(
            'CreateFunction',
            createFunction_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateFunctionRequest.fromBuffer(value),
            ($0.CreateFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateInvalidationRequest,
            $0.CreateInvalidationResult>(
        'CreateInvalidation',
        createInvalidation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateInvalidationRequest.fromBuffer(value),
        ($0.CreateInvalidationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.CreateInvalidationForDistributionTenantRequest,
            $0.CreateInvalidationForDistributionTenantResult>(
        'CreateInvalidationForDistributionTenant',
        createInvalidationForDistributionTenant_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateInvalidationForDistributionTenantRequest.fromBuffer(value),
        ($0.CreateInvalidationForDistributionTenantResult value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateKeyGroupRequest, $0.CreateKeyGroupResult>(
            'CreateKeyGroup',
            createKeyGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateKeyGroupRequest.fromBuffer(value),
            ($0.CreateKeyGroupResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateKeyValueStoreRequest,
            $0.CreateKeyValueStoreResult>(
        'CreateKeyValueStore',
        createKeyValueStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateKeyValueStoreRequest.fromBuffer(value),
        ($0.CreateKeyValueStoreResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateMonitoringSubscriptionRequest,
            $0.CreateMonitoringSubscriptionResult>(
        'CreateMonitoringSubscription',
        createMonitoringSubscription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateMonitoringSubscriptionRequest.fromBuffer(value),
        ($0.CreateMonitoringSubscriptionResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateOriginAccessControlRequest,
            $0.CreateOriginAccessControlResult>(
        'CreateOriginAccessControl',
        createOriginAccessControl_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateOriginAccessControlRequest.fromBuffer(value),
        ($0.CreateOriginAccessControlResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateOriginRequestPolicyRequest,
            $0.CreateOriginRequestPolicyResult>(
        'CreateOriginRequestPolicy',
        createOriginRequestPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateOriginRequestPolicyRequest.fromBuffer(value),
        ($0.CreateOriginRequestPolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePublicKeyRequest,
            $0.CreatePublicKeyResult>(
        'CreatePublicKey',
        createPublicKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePublicKeyRequest.fromBuffer(value),
        ($0.CreatePublicKeyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRealtimeLogConfigRequest,
            $0.CreateRealtimeLogConfigResult>(
        'CreateRealtimeLogConfig',
        createRealtimeLogConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRealtimeLogConfigRequest.fromBuffer(value),
        ($0.CreateRealtimeLogConfigResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateResponseHeadersPolicyRequest,
            $0.CreateResponseHeadersPolicyResult>(
        'CreateResponseHeadersPolicy',
        createResponseHeadersPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateResponseHeadersPolicyRequest.fromBuffer(value),
        ($0.CreateResponseHeadersPolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateStreamingDistributionRequest,
            $0.CreateStreamingDistributionResult>(
        'CreateStreamingDistribution',
        createStreamingDistribution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateStreamingDistributionRequest.fromBuffer(value),
        ($0.CreateStreamingDistributionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.CreateStreamingDistributionWithTagsRequest,
            $0.CreateStreamingDistributionWithTagsResult>(
        'CreateStreamingDistributionWithTags',
        createStreamingDistributionWithTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateStreamingDistributionWithTagsRequest.fromBuffer(value),
        ($0.CreateStreamingDistributionWithTagsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateTrustStoreRequest,
            $0.CreateTrustStoreResult>(
        'CreateTrustStore',
        createTrustStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateTrustStoreRequest.fromBuffer(value),
        ($0.CreateTrustStoreResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateVpcOriginRequest,
            $0.CreateVpcOriginResult>(
        'CreateVpcOrigin',
        createVpcOrigin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateVpcOriginRequest.fromBuffer(value),
        ($0.CreateVpcOriginResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAnycastIpListRequest, $1.Empty>(
        'DeleteAnycastIpList',
        deleteAnycastIpList_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAnycastIpListRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteCachePolicyRequest, $1.Empty>(
        'DeleteCachePolicy',
        deleteCachePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteCachePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeleteCloudFrontOriginAccessIdentityRequest, $1.Empty>(
        'DeleteCloudFrontOriginAccessIdentity',
        deleteCloudFrontOriginAccessIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteCloudFrontOriginAccessIdentityRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteConnectionFunctionRequest, $1.Empty>(
            'DeleteConnectionFunction',
            deleteConnectionFunction_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteConnectionFunctionRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteConnectionGroupRequest, $1.Empty>(
        'DeleteConnectionGroup',
        deleteConnectionGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteConnectionGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteContinuousDeploymentPolicyRequest,
            $1.Empty>(
        'DeleteContinuousDeploymentPolicy',
        deleteContinuousDeploymentPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteContinuousDeploymentPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDistributionRequest, $1.Empty>(
        'DeleteDistribution',
        deleteDistribution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDistributionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteDistributionTenantRequest, $1.Empty>(
            'DeleteDistributionTenant',
            deleteDistributionTenant_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteDistributionTenantRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteFieldLevelEncryptionConfigRequest,
            $1.Empty>(
        'DeleteFieldLevelEncryptionConfig',
        deleteFieldLevelEncryptionConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteFieldLevelEncryptionConfigRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteFieldLevelEncryptionProfileRequest,
            $1.Empty>(
        'DeleteFieldLevelEncryptionProfile',
        deleteFieldLevelEncryptionProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteFieldLevelEncryptionProfileRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteFunctionRequest, $1.Empty>(
        'DeleteFunction',
        deleteFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteFunctionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteKeyGroupRequest, $1.Empty>(
        'DeleteKeyGroup',
        deleteKeyGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteKeyGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteKeyValueStoreRequest, $1.Empty>(
        'DeleteKeyValueStore',
        deleteKeyValueStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteKeyValueStoreRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMonitoringSubscriptionRequest,
            $0.DeleteMonitoringSubscriptionResult>(
        'DeleteMonitoringSubscription',
        deleteMonitoringSubscription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMonitoringSubscriptionRequest.fromBuffer(value),
        ($0.DeleteMonitoringSubscriptionResult value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteOriginAccessControlRequest, $1.Empty>(
            'DeleteOriginAccessControl',
            deleteOriginAccessControl_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteOriginAccessControlRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteOriginRequestPolicyRequest, $1.Empty>(
            'DeleteOriginRequestPolicy',
            deleteOriginRequestPolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteOriginRequestPolicyRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePublicKeyRequest, $1.Empty>(
        'DeletePublicKey',
        deletePublicKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeletePublicKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRealtimeLogConfigRequest, $1.Empty>(
        'DeleteRealtimeLogConfig',
        deleteRealtimeLogConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRealtimeLogConfigRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourcePolicyRequest, $1.Empty>(
        'DeleteResourcePolicy',
        deleteResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourcePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteResponseHeadersPolicyRequest, $1.Empty>(
            'DeleteResponseHeadersPolicy',
            deleteResponseHeadersPolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteResponseHeadersPolicyRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteStreamingDistributionRequest, $1.Empty>(
            'DeleteStreamingDistribution',
            deleteStreamingDistribution_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteStreamingDistributionRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTrustStoreRequest, $1.Empty>(
        'DeleteTrustStore',
        deleteTrustStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteTrustStoreRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteVpcOriginRequest,
            $0.DeleteVpcOriginResult>(
        'DeleteVpcOrigin',
        deleteVpcOrigin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteVpcOriginRequest.fromBuffer(value),
        ($0.DeleteVpcOriginResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeConnectionFunctionRequest,
            $0.DescribeConnectionFunctionResult>(
        'DescribeConnectionFunction',
        describeConnectionFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeConnectionFunctionRequest.fromBuffer(value),
        ($0.DescribeConnectionFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeFunctionRequest,
            $0.DescribeFunctionResult>(
        'DescribeFunction',
        describeFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeFunctionRequest.fromBuffer(value),
        ($0.DescribeFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeKeyValueStoreRequest,
            $0.DescribeKeyValueStoreResult>(
        'DescribeKeyValueStore',
        describeKeyValueStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeKeyValueStoreRequest.fromBuffer(value),
        ($0.DescribeKeyValueStoreResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DisassociateDistributionTenantWebACLRequest,
            $0.DisassociateDistributionTenantWebACLResult>(
        'DisassociateDistributionTenantWebACL',
        disassociateDistributionTenantWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisassociateDistributionTenantWebACLRequest.fromBuffer(value),
        ($0.DisassociateDistributionTenantWebACLResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisassociateDistributionWebACLRequest,
            $0.DisassociateDistributionWebACLResult>(
        'DisassociateDistributionWebACL',
        disassociateDistributionWebACL_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisassociateDistributionWebACLRequest.fromBuffer(value),
        ($0.DisassociateDistributionWebACLResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAnycastIpListRequest,
            $0.GetAnycastIpListResult>(
        'GetAnycastIpList',
        getAnycastIpList_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAnycastIpListRequest.fromBuffer(value),
        ($0.GetAnycastIpListResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetCachePolicyRequest, $0.GetCachePolicyResult>(
            'GetCachePolicy',
            getCachePolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetCachePolicyRequest.fromBuffer(value),
            ($0.GetCachePolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCachePolicyConfigRequest,
            $0.GetCachePolicyConfigResult>(
        'GetCachePolicyConfig',
        getCachePolicyConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCachePolicyConfigRequest.fromBuffer(value),
        ($0.GetCachePolicyConfigResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCloudFrontOriginAccessIdentityRequest,
            $0.GetCloudFrontOriginAccessIdentityResult>(
        'GetCloudFrontOriginAccessIdentity',
        getCloudFrontOriginAccessIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCloudFrontOriginAccessIdentityRequest.fromBuffer(value),
        ($0.GetCloudFrontOriginAccessIdentityResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetCloudFrontOriginAccessIdentityConfigRequest,
            $0.GetCloudFrontOriginAccessIdentityConfigResult>(
        'GetCloudFrontOriginAccessIdentityConfig',
        getCloudFrontOriginAccessIdentityConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCloudFrontOriginAccessIdentityConfigRequest.fromBuffer(value),
        ($0.GetCloudFrontOriginAccessIdentityConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetConnectionFunctionRequest,
            $0.GetConnectionFunctionResult>(
        'GetConnectionFunction',
        getConnectionFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConnectionFunctionRequest.fromBuffer(value),
        ($0.GetConnectionFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetConnectionGroupRequest,
            $0.GetConnectionGroupResult>(
        'GetConnectionGroup',
        getConnectionGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConnectionGroupRequest.fromBuffer(value),
        ($0.GetConnectionGroupResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetConnectionGroupByRoutingEndpointRequest,
            $0.GetConnectionGroupByRoutingEndpointResult>(
        'GetConnectionGroupByRoutingEndpoint',
        getConnectionGroupByRoutingEndpoint_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConnectionGroupByRoutingEndpointRequest.fromBuffer(value),
        ($0.GetConnectionGroupByRoutingEndpointResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetContinuousDeploymentPolicyRequest,
            $0.GetContinuousDeploymentPolicyResult>(
        'GetContinuousDeploymentPolicy',
        getContinuousDeploymentPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetContinuousDeploymentPolicyRequest.fromBuffer(value),
        ($0.GetContinuousDeploymentPolicyResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetContinuousDeploymentPolicyConfigRequest,
            $0.GetContinuousDeploymentPolicyConfigResult>(
        'GetContinuousDeploymentPolicyConfig',
        getContinuousDeploymentPolicyConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetContinuousDeploymentPolicyConfigRequest.fromBuffer(value),
        ($0.GetContinuousDeploymentPolicyConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDistributionRequest,
            $0.GetDistributionResult>(
        'GetDistribution',
        getDistribution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDistributionRequest.fromBuffer(value),
        ($0.GetDistributionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDistributionConfigRequest,
            $0.GetDistributionConfigResult>(
        'GetDistributionConfig',
        getDistributionConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDistributionConfigRequest.fromBuffer(value),
        ($0.GetDistributionConfigResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDistributionTenantRequest,
            $0.GetDistributionTenantResult>(
        'GetDistributionTenant',
        getDistributionTenant_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDistributionTenantRequest.fromBuffer(value),
        ($0.GetDistributionTenantResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDistributionTenantByDomainRequest,
            $0.GetDistributionTenantByDomainResult>(
        'GetDistributionTenantByDomain',
        getDistributionTenantByDomain_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDistributionTenantByDomainRequest.fromBuffer(value),
        ($0.GetDistributionTenantByDomainResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFieldLevelEncryptionRequest,
            $0.GetFieldLevelEncryptionResult>(
        'GetFieldLevelEncryption',
        getFieldLevelEncryption_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFieldLevelEncryptionRequest.fromBuffer(value),
        ($0.GetFieldLevelEncryptionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFieldLevelEncryptionConfigRequest,
            $0.GetFieldLevelEncryptionConfigResult>(
        'GetFieldLevelEncryptionConfig',
        getFieldLevelEncryptionConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFieldLevelEncryptionConfigRequest.fromBuffer(value),
        ($0.GetFieldLevelEncryptionConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFieldLevelEncryptionProfileRequest,
            $0.GetFieldLevelEncryptionProfileResult>(
        'GetFieldLevelEncryptionProfile',
        getFieldLevelEncryptionProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFieldLevelEncryptionProfileRequest.fromBuffer(value),
        ($0.GetFieldLevelEncryptionProfileResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetFieldLevelEncryptionProfileConfigRequest,
            $0.GetFieldLevelEncryptionProfileConfigResult>(
        'GetFieldLevelEncryptionProfileConfig',
        getFieldLevelEncryptionProfileConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFieldLevelEncryptionProfileConfigRequest.fromBuffer(value),
        ($0.GetFieldLevelEncryptionProfileConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFunctionRequest, $0.GetFunctionResult>(
        'GetFunction',
        getFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFunctionRequest.fromBuffer(value),
        ($0.GetFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetInvalidationRequest,
            $0.GetInvalidationResult>(
        'GetInvalidation',
        getInvalidation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetInvalidationRequest.fromBuffer(value),
        ($0.GetInvalidationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetInvalidationForDistributionTenantRequest,
            $0.GetInvalidationForDistributionTenantResult>(
        'GetInvalidationForDistributionTenant',
        getInvalidationForDistributionTenant_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetInvalidationForDistributionTenantRequest.fromBuffer(value),
        ($0.GetInvalidationForDistributionTenantResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetKeyGroupRequest, $0.GetKeyGroupResult>(
        'GetKeyGroup',
        getKeyGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetKeyGroupRequest.fromBuffer(value),
        ($0.GetKeyGroupResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetKeyGroupConfigRequest,
            $0.GetKeyGroupConfigResult>(
        'GetKeyGroupConfig',
        getKeyGroupConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetKeyGroupConfigRequest.fromBuffer(value),
        ($0.GetKeyGroupConfigResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetManagedCertificateDetailsRequest,
            $0.GetManagedCertificateDetailsResult>(
        'GetManagedCertificateDetails',
        getManagedCertificateDetails_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetManagedCertificateDetailsRequest.fromBuffer(value),
        ($0.GetManagedCertificateDetailsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMonitoringSubscriptionRequest,
            $0.GetMonitoringSubscriptionResult>(
        'GetMonitoringSubscription',
        getMonitoringSubscription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMonitoringSubscriptionRequest.fromBuffer(value),
        ($0.GetMonitoringSubscriptionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOriginAccessControlRequest,
            $0.GetOriginAccessControlResult>(
        'GetOriginAccessControl',
        getOriginAccessControl_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetOriginAccessControlRequest.fromBuffer(value),
        ($0.GetOriginAccessControlResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOriginAccessControlConfigRequest,
            $0.GetOriginAccessControlConfigResult>(
        'GetOriginAccessControlConfig',
        getOriginAccessControlConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetOriginAccessControlConfigRequest.fromBuffer(value),
        ($0.GetOriginAccessControlConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOriginRequestPolicyRequest,
            $0.GetOriginRequestPolicyResult>(
        'GetOriginRequestPolicy',
        getOriginRequestPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetOriginRequestPolicyRequest.fromBuffer(value),
        ($0.GetOriginRequestPolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOriginRequestPolicyConfigRequest,
            $0.GetOriginRequestPolicyConfigResult>(
        'GetOriginRequestPolicyConfig',
        getOriginRequestPolicyConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetOriginRequestPolicyConfigRequest.fromBuffer(value),
        ($0.GetOriginRequestPolicyConfigResult value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetPublicKeyRequest, $0.GetPublicKeyResult>(
            'GetPublicKey',
            getPublicKey_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetPublicKeyRequest.fromBuffer(value),
            ($0.GetPublicKeyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPublicKeyConfigRequest,
            $0.GetPublicKeyConfigResult>(
        'GetPublicKeyConfig',
        getPublicKeyConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPublicKeyConfigRequest.fromBuffer(value),
        ($0.GetPublicKeyConfigResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRealtimeLogConfigRequest,
            $0.GetRealtimeLogConfigResult>(
        'GetRealtimeLogConfig',
        getRealtimeLogConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRealtimeLogConfigRequest.fromBuffer(value),
        ($0.GetRealtimeLogConfigResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourcePolicyRequest,
            $0.GetResourcePolicyResult>(
        'GetResourcePolicy',
        getResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourcePolicyRequest.fromBuffer(value),
        ($0.GetResourcePolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResponseHeadersPolicyRequest,
            $0.GetResponseHeadersPolicyResult>(
        'GetResponseHeadersPolicy',
        getResponseHeadersPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResponseHeadersPolicyRequest.fromBuffer(value),
        ($0.GetResponseHeadersPolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResponseHeadersPolicyConfigRequest,
            $0.GetResponseHeadersPolicyConfigResult>(
        'GetResponseHeadersPolicyConfig',
        getResponseHeadersPolicyConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResponseHeadersPolicyConfigRequest.fromBuffer(value),
        ($0.GetResponseHeadersPolicyConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetStreamingDistributionRequest,
            $0.GetStreamingDistributionResult>(
        'GetStreamingDistribution',
        getStreamingDistribution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetStreamingDistributionRequest.fromBuffer(value),
        ($0.GetStreamingDistributionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetStreamingDistributionConfigRequest,
            $0.GetStreamingDistributionConfigResult>(
        'GetStreamingDistributionConfig',
        getStreamingDistributionConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetStreamingDistributionConfigRequest.fromBuffer(value),
        ($0.GetStreamingDistributionConfigResult value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetTrustStoreRequest, $0.GetTrustStoreResult>(
            'GetTrustStore',
            getTrustStore_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetTrustStoreRequest.fromBuffer(value),
            ($0.GetTrustStoreResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetVpcOriginRequest, $0.GetVpcOriginResult>(
            'GetVpcOrigin',
            getVpcOrigin_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetVpcOriginRequest.fromBuffer(value),
            ($0.GetVpcOriginResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAnycastIpListsRequest,
            $0.ListAnycastIpListsResult>(
        'ListAnycastIpLists',
        listAnycastIpLists_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAnycastIpListsRequest.fromBuffer(value),
        ($0.ListAnycastIpListsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCachePoliciesRequest,
            $0.ListCachePoliciesResult>(
        'ListCachePolicies',
        listCachePolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCachePoliciesRequest.fromBuffer(value),
        ($0.ListCachePoliciesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListCloudFrontOriginAccessIdentitiesRequest,
            $0.ListCloudFrontOriginAccessIdentitiesResult>(
        'ListCloudFrontOriginAccessIdentities',
        listCloudFrontOriginAccessIdentities_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCloudFrontOriginAccessIdentitiesRequest.fromBuffer(value),
        ($0.ListCloudFrontOriginAccessIdentitiesResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListConflictingAliasesRequest,
            $0.ListConflictingAliasesResult>(
        'ListConflictingAliases',
        listConflictingAliases_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListConflictingAliasesRequest.fromBuffer(value),
        ($0.ListConflictingAliasesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListConnectionFunctionsRequest,
            $0.ListConnectionFunctionsResult>(
        'ListConnectionFunctions',
        listConnectionFunctions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListConnectionFunctionsRequest.fromBuffer(value),
        ($0.ListConnectionFunctionsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListConnectionGroupsRequest,
            $0.ListConnectionGroupsResult>(
        'ListConnectionGroups',
        listConnectionGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListConnectionGroupsRequest.fromBuffer(value),
        ($0.ListConnectionGroupsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListContinuousDeploymentPoliciesRequest,
            $0.ListContinuousDeploymentPoliciesResult>(
        'ListContinuousDeploymentPolicies',
        listContinuousDeploymentPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListContinuousDeploymentPoliciesRequest.fromBuffer(value),
        ($0.ListContinuousDeploymentPoliciesResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsRequest,
            $0.ListDistributionsResult>(
        'ListDistributions',
        listDistributions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsRequest.fromBuffer(value),
        ($0.ListDistributionsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsByAnycastIpListIdRequest,
            $0.ListDistributionsByAnycastIpListIdResult>(
        'ListDistributionsByAnycastIpListId',
        listDistributionsByAnycastIpListId_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByAnycastIpListIdRequest.fromBuffer(value),
        ($0.ListDistributionsByAnycastIpListIdResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsByCachePolicyIdRequest,
            $0.ListDistributionsByCachePolicyIdResult>(
        'ListDistributionsByCachePolicyId',
        listDistributionsByCachePolicyId_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByCachePolicyIdRequest.fromBuffer(value),
        ($0.ListDistributionsByCachePolicyIdResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListDistributionsByConnectionFunctionRequest,
            $0.ListDistributionsByConnectionFunctionResult>(
        'ListDistributionsByConnectionFunction',
        listDistributionsByConnectionFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByConnectionFunctionRequest.fromBuffer(value),
        ($0.ListDistributionsByConnectionFunctionResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsByConnectionModeRequest,
            $0.ListDistributionsByConnectionModeResult>(
        'ListDistributionsByConnectionMode',
        listDistributionsByConnectionMode_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByConnectionModeRequest.fromBuffer(value),
        ($0.ListDistributionsByConnectionModeResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsByKeyGroupRequest,
            $0.ListDistributionsByKeyGroupResult>(
        'ListDistributionsByKeyGroup',
        listDistributionsByKeyGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByKeyGroupRequest.fromBuffer(value),
        ($0.ListDistributionsByKeyGroupResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListDistributionsByOriginRequestPolicyIdRequest,
            $0.ListDistributionsByOriginRequestPolicyIdResult>(
        'ListDistributionsByOriginRequestPolicyId',
        listDistributionsByOriginRequestPolicyId_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByOriginRequestPolicyIdRequest.fromBuffer(
                value),
        ($0.ListDistributionsByOriginRequestPolicyIdResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsByOwnedResourceRequest,
            $0.ListDistributionsByOwnedResourceResult>(
        'ListDistributionsByOwnedResource',
        listDistributionsByOwnedResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByOwnedResourceRequest.fromBuffer(value),
        ($0.ListDistributionsByOwnedResourceResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListDistributionsByRealtimeLogConfigRequest,
            $0.ListDistributionsByRealtimeLogConfigResult>(
        'ListDistributionsByRealtimeLogConfig',
        listDistributionsByRealtimeLogConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByRealtimeLogConfigRequest.fromBuffer(value),
        ($0.ListDistributionsByRealtimeLogConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListDistributionsByResponseHeadersPolicyIdRequest,
            $0.ListDistributionsByResponseHeadersPolicyIdResult>(
        'ListDistributionsByResponseHeadersPolicyId',
        listDistributionsByResponseHeadersPolicyId_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByResponseHeadersPolicyIdRequest.fromBuffer(
                value),
        ($0.ListDistributionsByResponseHeadersPolicyIdResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsByTrustStoreRequest,
            $0.ListDistributionsByTrustStoreResult>(
        'ListDistributionsByTrustStore',
        listDistributionsByTrustStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByTrustStoreRequest.fromBuffer(value),
        ($0.ListDistributionsByTrustStoreResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsByVpcOriginIdRequest,
            $0.ListDistributionsByVpcOriginIdResult>(
        'ListDistributionsByVpcOriginId',
        listDistributionsByVpcOriginId_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByVpcOriginIdRequest.fromBuffer(value),
        ($0.ListDistributionsByVpcOriginIdResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionsByWebACLIdRequest,
            $0.ListDistributionsByWebACLIdResult>(
        'ListDistributionsByWebACLId',
        listDistributionsByWebACLId_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionsByWebACLIdRequest.fromBuffer(value),
        ($0.ListDistributionsByWebACLIdResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDistributionTenantsRequest,
            $0.ListDistributionTenantsResult>(
        'ListDistributionTenants',
        listDistributionTenants_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionTenantsRequest.fromBuffer(value),
        ($0.ListDistributionTenantsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListDistributionTenantsByCustomizationRequest,
            $0.ListDistributionTenantsByCustomizationResult>(
        'ListDistributionTenantsByCustomization',
        listDistributionTenantsByCustomization_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDistributionTenantsByCustomizationRequest.fromBuffer(value),
        ($0.ListDistributionTenantsByCustomizationResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDomainConflictsRequest,
            $0.ListDomainConflictsResult>(
        'ListDomainConflicts',
        listDomainConflicts_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDomainConflictsRequest.fromBuffer(value),
        ($0.ListDomainConflictsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListFieldLevelEncryptionConfigsRequest,
            $0.ListFieldLevelEncryptionConfigsResult>(
        'ListFieldLevelEncryptionConfigs',
        listFieldLevelEncryptionConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListFieldLevelEncryptionConfigsRequest.fromBuffer(value),
        ($0.ListFieldLevelEncryptionConfigsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListFieldLevelEncryptionProfilesRequest,
            $0.ListFieldLevelEncryptionProfilesResult>(
        'ListFieldLevelEncryptionProfiles',
        listFieldLevelEncryptionProfiles_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListFieldLevelEncryptionProfilesRequest.fromBuffer(value),
        ($0.ListFieldLevelEncryptionProfilesResult value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListFunctionsRequest, $0.ListFunctionsResult>(
            'ListFunctions',
            listFunctions_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListFunctionsRequest.fromBuffer(value),
            ($0.ListFunctionsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListInvalidationsRequest,
            $0.ListInvalidationsResult>(
        'ListInvalidations',
        listInvalidations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListInvalidationsRequest.fromBuffer(value),
        ($0.ListInvalidationsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListInvalidationsForDistributionTenantRequest,
            $0.ListInvalidationsForDistributionTenantResult>(
        'ListInvalidationsForDistributionTenant',
        listInvalidationsForDistributionTenant_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListInvalidationsForDistributionTenantRequest.fromBuffer(value),
        ($0.ListInvalidationsForDistributionTenantResult value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListKeyGroupsRequest, $0.ListKeyGroupsResult>(
            'ListKeyGroups',
            listKeyGroups_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListKeyGroupsRequest.fromBuffer(value),
            ($0.ListKeyGroupsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListKeyValueStoresRequest,
            $0.ListKeyValueStoresResult>(
        'ListKeyValueStores',
        listKeyValueStores_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListKeyValueStoresRequest.fromBuffer(value),
        ($0.ListKeyValueStoresResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOriginAccessControlsRequest,
            $0.ListOriginAccessControlsResult>(
        'ListOriginAccessControls',
        listOriginAccessControls_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOriginAccessControlsRequest.fromBuffer(value),
        ($0.ListOriginAccessControlsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOriginRequestPoliciesRequest,
            $0.ListOriginRequestPoliciesResult>(
        'ListOriginRequestPolicies',
        listOriginRequestPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOriginRequestPoliciesRequest.fromBuffer(value),
        ($0.ListOriginRequestPoliciesResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListPublicKeysRequest, $0.ListPublicKeysResult>(
            'ListPublicKeys',
            listPublicKeys_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListPublicKeysRequest.fromBuffer(value),
            ($0.ListPublicKeysResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRealtimeLogConfigsRequest,
            $0.ListRealtimeLogConfigsResult>(
        'ListRealtimeLogConfigs',
        listRealtimeLogConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRealtimeLogConfigsRequest.fromBuffer(value),
        ($0.ListRealtimeLogConfigsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListResponseHeadersPoliciesRequest,
            $0.ListResponseHeadersPoliciesResult>(
        'ListResponseHeadersPolicies',
        listResponseHeadersPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListResponseHeadersPoliciesRequest.fromBuffer(value),
        ($0.ListResponseHeadersPoliciesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListStreamingDistributionsRequest,
            $0.ListStreamingDistributionsResult>(
        'ListStreamingDistributions',
        listStreamingDistributions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListStreamingDistributionsRequest.fromBuffer(value),
        ($0.ListStreamingDistributionsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResult>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTrustStoresRequest,
            $0.ListTrustStoresResult>(
        'ListTrustStores',
        listTrustStores_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTrustStoresRequest.fromBuffer(value),
        ($0.ListTrustStoresResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListVpcOriginsRequest, $0.ListVpcOriginsResult>(
            'ListVpcOrigins',
            listVpcOrigins_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListVpcOriginsRequest.fromBuffer(value),
            ($0.ListVpcOriginsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PublishConnectionFunctionRequest,
            $0.PublishConnectionFunctionResult>(
        'PublishConnectionFunction',
        publishConnectionFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PublishConnectionFunctionRequest.fromBuffer(value),
        ($0.PublishConnectionFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PublishFunctionRequest,
            $0.PublishFunctionResult>(
        'PublishFunction',
        publishFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PublishFunctionRequest.fromBuffer(value),
        ($0.PublishFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutResourcePolicyRequest,
            $0.PutResourcePolicyResult>(
        'PutResourcePolicy',
        putResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutResourcePolicyRequest.fromBuffer(value),
        ($0.PutResourcePolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceRequest, $1.Empty>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestConnectionFunctionRequest,
            $0.TestConnectionFunctionResult>(
        'TestConnectionFunction',
        testConnectionFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TestConnectionFunctionRequest.fromBuffer(value),
        ($0.TestConnectionFunctionResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TestFunctionRequest, $0.TestFunctionResult>(
            'TestFunction',
            testFunction_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TestFunctionRequest.fromBuffer(value),
            ($0.TestFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagResourceRequest, $1.Empty>(
        'UntagResource',
        untagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAnycastIpListRequest,
            $0.UpdateAnycastIpListResult>(
        'UpdateAnycastIpList',
        updateAnycastIpList_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAnycastIpListRequest.fromBuffer(value),
        ($0.UpdateAnycastIpListResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateCachePolicyRequest,
            $0.UpdateCachePolicyResult>(
        'UpdateCachePolicy',
        updateCachePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateCachePolicyRequest.fromBuffer(value),
        ($0.UpdateCachePolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateCloudFrontOriginAccessIdentityRequest,
            $0.UpdateCloudFrontOriginAccessIdentityResult>(
        'UpdateCloudFrontOriginAccessIdentity',
        updateCloudFrontOriginAccessIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateCloudFrontOriginAccessIdentityRequest.fromBuffer(value),
        ($0.UpdateCloudFrontOriginAccessIdentityResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateConnectionFunctionRequest,
            $0.UpdateConnectionFunctionResult>(
        'UpdateConnectionFunction',
        updateConnectionFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateConnectionFunctionRequest.fromBuffer(value),
        ($0.UpdateConnectionFunctionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateConnectionGroupRequest,
            $0.UpdateConnectionGroupResult>(
        'UpdateConnectionGroup',
        updateConnectionGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateConnectionGroupRequest.fromBuffer(value),
        ($0.UpdateConnectionGroupResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateContinuousDeploymentPolicyRequest,
            $0.UpdateContinuousDeploymentPolicyResult>(
        'UpdateContinuousDeploymentPolicy',
        updateContinuousDeploymentPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateContinuousDeploymentPolicyRequest.fromBuffer(value),
        ($0.UpdateContinuousDeploymentPolicyResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDistributionRequest,
            $0.UpdateDistributionResult>(
        'UpdateDistribution',
        updateDistribution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDistributionRequest.fromBuffer(value),
        ($0.UpdateDistributionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDistributionTenantRequest,
            $0.UpdateDistributionTenantResult>(
        'UpdateDistributionTenant',
        updateDistributionTenant_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDistributionTenantRequest.fromBuffer(value),
        ($0.UpdateDistributionTenantResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateDistributionWithStagingConfigRequest,
            $0.UpdateDistributionWithStagingConfigResult>(
        'UpdateDistributionWithStagingConfig',
        updateDistributionWithStagingConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDistributionWithStagingConfigRequest.fromBuffer(value),
        ($0.UpdateDistributionWithStagingConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDomainAssociationRequest,
            $0.UpdateDomainAssociationResult>(
        'UpdateDomainAssociation',
        updateDomainAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDomainAssociationRequest.fromBuffer(value),
        ($0.UpdateDomainAssociationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateFieldLevelEncryptionConfigRequest,
            $0.UpdateFieldLevelEncryptionConfigResult>(
        'UpdateFieldLevelEncryptionConfig',
        updateFieldLevelEncryptionConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateFieldLevelEncryptionConfigRequest.fromBuffer(value),
        ($0.UpdateFieldLevelEncryptionConfigResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateFieldLevelEncryptionProfileRequest,
            $0.UpdateFieldLevelEncryptionProfileResult>(
        'UpdateFieldLevelEncryptionProfile',
        updateFieldLevelEncryptionProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateFieldLevelEncryptionProfileRequest.fromBuffer(value),
        ($0.UpdateFieldLevelEncryptionProfileResult value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateFunctionRequest, $0.UpdateFunctionResult>(
            'UpdateFunction',
            updateFunction_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateFunctionRequest.fromBuffer(value),
            ($0.UpdateFunctionResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateKeyGroupRequest, $0.UpdateKeyGroupResult>(
            'UpdateKeyGroup',
            updateKeyGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateKeyGroupRequest.fromBuffer(value),
            ($0.UpdateKeyGroupResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateKeyValueStoreRequest,
            $0.UpdateKeyValueStoreResult>(
        'UpdateKeyValueStore',
        updateKeyValueStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateKeyValueStoreRequest.fromBuffer(value),
        ($0.UpdateKeyValueStoreResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateOriginAccessControlRequest,
            $0.UpdateOriginAccessControlResult>(
        'UpdateOriginAccessControl',
        updateOriginAccessControl_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateOriginAccessControlRequest.fromBuffer(value),
        ($0.UpdateOriginAccessControlResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateOriginRequestPolicyRequest,
            $0.UpdateOriginRequestPolicyResult>(
        'UpdateOriginRequestPolicy',
        updateOriginRequestPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateOriginRequestPolicyRequest.fromBuffer(value),
        ($0.UpdateOriginRequestPolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdatePublicKeyRequest,
            $0.UpdatePublicKeyResult>(
        'UpdatePublicKey',
        updatePublicKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdatePublicKeyRequest.fromBuffer(value),
        ($0.UpdatePublicKeyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRealtimeLogConfigRequest,
            $0.UpdateRealtimeLogConfigResult>(
        'UpdateRealtimeLogConfig',
        updateRealtimeLogConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRealtimeLogConfigRequest.fromBuffer(value),
        ($0.UpdateRealtimeLogConfigResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateResponseHeadersPolicyRequest,
            $0.UpdateResponseHeadersPolicyResult>(
        'UpdateResponseHeadersPolicy',
        updateResponseHeadersPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateResponseHeadersPolicyRequest.fromBuffer(value),
        ($0.UpdateResponseHeadersPolicyResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateStreamingDistributionRequest,
            $0.UpdateStreamingDistributionResult>(
        'UpdateStreamingDistribution',
        updateStreamingDistribution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateStreamingDistributionRequest.fromBuffer(value),
        ($0.UpdateStreamingDistributionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateTrustStoreRequest,
            $0.UpdateTrustStoreResult>(
        'UpdateTrustStore',
        updateTrustStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateTrustStoreRequest.fromBuffer(value),
        ($0.UpdateTrustStoreResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateVpcOriginRequest,
            $0.UpdateVpcOriginResult>(
        'UpdateVpcOrigin',
        updateVpcOrigin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateVpcOriginRequest.fromBuffer(value),
        ($0.UpdateVpcOriginResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.VerifyDnsConfigurationRequest,
            $0.VerifyDnsConfigurationResult>(
        'VerifyDnsConfiguration',
        verifyDnsConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.VerifyDnsConfigurationRequest.fromBuffer(value),
        ($0.VerifyDnsConfigurationResult value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> associateAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AssociateAliasRequest> $request) async {
    return associateAlias($call, await $request);
  }

  $async.Future<$1.Empty> associateAlias(
      $grpc.ServiceCall call, $0.AssociateAliasRequest request);

  $async.Future<$0.AssociateDistributionTenantWebACLResult>
      associateDistributionTenantWebACL_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.AssociateDistributionTenantWebACLRequest>
              $request) async {
    return associateDistributionTenantWebACL($call, await $request);
  }

  $async.Future<$0.AssociateDistributionTenantWebACLResult>
      associateDistributionTenantWebACL($grpc.ServiceCall call,
          $0.AssociateDistributionTenantWebACLRequest request);

  $async.Future<$0.AssociateDistributionWebACLResult>
      associateDistributionWebACL_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AssociateDistributionWebACLRequest> $request) async {
    return associateDistributionWebACL($call, await $request);
  }

  $async.Future<$0.AssociateDistributionWebACLResult>
      associateDistributionWebACL($grpc.ServiceCall call,
          $0.AssociateDistributionWebACLRequest request);

  $async.Future<$0.CopyDistributionResult> copyDistribution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CopyDistributionRequest> $request) async {
    return copyDistribution($call, await $request);
  }

  $async.Future<$0.CopyDistributionResult> copyDistribution(
      $grpc.ServiceCall call, $0.CopyDistributionRequest request);

  $async.Future<$0.CreateAnycastIpListResult> createAnycastIpList_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateAnycastIpListRequest> $request) async {
    return createAnycastIpList($call, await $request);
  }

  $async.Future<$0.CreateAnycastIpListResult> createAnycastIpList(
      $grpc.ServiceCall call, $0.CreateAnycastIpListRequest request);

  $async.Future<$0.CreateCachePolicyResult> createCachePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateCachePolicyRequest> $request) async {
    return createCachePolicy($call, await $request);
  }

  $async.Future<$0.CreateCachePolicyResult> createCachePolicy(
      $grpc.ServiceCall call, $0.CreateCachePolicyRequest request);

  $async.Future<$0.CreateCloudFrontOriginAccessIdentityResult>
      createCloudFrontOriginAccessIdentity_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateCloudFrontOriginAccessIdentityRequest>
              $request) async {
    return createCloudFrontOriginAccessIdentity($call, await $request);
  }

  $async.Future<$0.CreateCloudFrontOriginAccessIdentityResult>
      createCloudFrontOriginAccessIdentity($grpc.ServiceCall call,
          $0.CreateCloudFrontOriginAccessIdentityRequest request);

  $async.Future<$0.CreateConnectionFunctionResult> createConnectionFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateConnectionFunctionRequest> $request) async {
    return createConnectionFunction($call, await $request);
  }

  $async.Future<$0.CreateConnectionFunctionResult> createConnectionFunction(
      $grpc.ServiceCall call, $0.CreateConnectionFunctionRequest request);

  $async.Future<$0.CreateConnectionGroupResult> createConnectionGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateConnectionGroupRequest> $request) async {
    return createConnectionGroup($call, await $request);
  }

  $async.Future<$0.CreateConnectionGroupResult> createConnectionGroup(
      $grpc.ServiceCall call, $0.CreateConnectionGroupRequest request);

  $async.Future<$0.CreateContinuousDeploymentPolicyResult>
      createContinuousDeploymentPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateContinuousDeploymentPolicyRequest>
              $request) async {
    return createContinuousDeploymentPolicy($call, await $request);
  }

  $async.Future<$0.CreateContinuousDeploymentPolicyResult>
      createContinuousDeploymentPolicy($grpc.ServiceCall call,
          $0.CreateContinuousDeploymentPolicyRequest request);

  $async.Future<$0.CreateDistributionResult> createDistribution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDistributionRequest> $request) async {
    return createDistribution($call, await $request);
  }

  $async.Future<$0.CreateDistributionResult> createDistribution(
      $grpc.ServiceCall call, $0.CreateDistributionRequest request);

  $async.Future<$0.CreateDistributionTenantResult> createDistributionTenant_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDistributionTenantRequest> $request) async {
    return createDistributionTenant($call, await $request);
  }

  $async.Future<$0.CreateDistributionTenantResult> createDistributionTenant(
      $grpc.ServiceCall call, $0.CreateDistributionTenantRequest request);

  $async.Future<$0.CreateDistributionWithTagsResult>
      createDistributionWithTags_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateDistributionWithTagsRequest> $request) async {
    return createDistributionWithTags($call, await $request);
  }

  $async.Future<$0.CreateDistributionWithTagsResult> createDistributionWithTags(
      $grpc.ServiceCall call, $0.CreateDistributionWithTagsRequest request);

  $async.Future<$0.CreateFieldLevelEncryptionConfigResult>
      createFieldLevelEncryptionConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateFieldLevelEncryptionConfigRequest>
              $request) async {
    return createFieldLevelEncryptionConfig($call, await $request);
  }

  $async.Future<$0.CreateFieldLevelEncryptionConfigResult>
      createFieldLevelEncryptionConfig($grpc.ServiceCall call,
          $0.CreateFieldLevelEncryptionConfigRequest request);

  $async.Future<$0.CreateFieldLevelEncryptionProfileResult>
      createFieldLevelEncryptionProfile_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateFieldLevelEncryptionProfileRequest>
              $request) async {
    return createFieldLevelEncryptionProfile($call, await $request);
  }

  $async.Future<$0.CreateFieldLevelEncryptionProfileResult>
      createFieldLevelEncryptionProfile($grpc.ServiceCall call,
          $0.CreateFieldLevelEncryptionProfileRequest request);

  $async.Future<$0.CreateFunctionResult> createFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateFunctionRequest> $request) async {
    return createFunction($call, await $request);
  }

  $async.Future<$0.CreateFunctionResult> createFunction(
      $grpc.ServiceCall call, $0.CreateFunctionRequest request);

  $async.Future<$0.CreateInvalidationResult> createInvalidation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateInvalidationRequest> $request) async {
    return createInvalidation($call, await $request);
  }

  $async.Future<$0.CreateInvalidationResult> createInvalidation(
      $grpc.ServiceCall call, $0.CreateInvalidationRequest request);

  $async.Future<$0.CreateInvalidationForDistributionTenantResult>
      createInvalidationForDistributionTenant_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateInvalidationForDistributionTenantRequest>
              $request) async {
    return createInvalidationForDistributionTenant($call, await $request);
  }

  $async.Future<$0.CreateInvalidationForDistributionTenantResult>
      createInvalidationForDistributionTenant($grpc.ServiceCall call,
          $0.CreateInvalidationForDistributionTenantRequest request);

  $async.Future<$0.CreateKeyGroupResult> createKeyGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateKeyGroupRequest> $request) async {
    return createKeyGroup($call, await $request);
  }

  $async.Future<$0.CreateKeyGroupResult> createKeyGroup(
      $grpc.ServiceCall call, $0.CreateKeyGroupRequest request);

  $async.Future<$0.CreateKeyValueStoreResult> createKeyValueStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateKeyValueStoreRequest> $request) async {
    return createKeyValueStore($call, await $request);
  }

  $async.Future<$0.CreateKeyValueStoreResult> createKeyValueStore(
      $grpc.ServiceCall call, $0.CreateKeyValueStoreRequest request);

  $async.Future<$0.CreateMonitoringSubscriptionResult>
      createMonitoringSubscription_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateMonitoringSubscriptionRequest>
              $request) async {
    return createMonitoringSubscription($call, await $request);
  }

  $async.Future<$0.CreateMonitoringSubscriptionResult>
      createMonitoringSubscription($grpc.ServiceCall call,
          $0.CreateMonitoringSubscriptionRequest request);

  $async.Future<$0.CreateOriginAccessControlResult>
      createOriginAccessControl_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateOriginAccessControlRequest> $request) async {
    return createOriginAccessControl($call, await $request);
  }

  $async.Future<$0.CreateOriginAccessControlResult> createOriginAccessControl(
      $grpc.ServiceCall call, $0.CreateOriginAccessControlRequest request);

  $async.Future<$0.CreateOriginRequestPolicyResult>
      createOriginRequestPolicy_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateOriginRequestPolicyRequest> $request) async {
    return createOriginRequestPolicy($call, await $request);
  }

  $async.Future<$0.CreateOriginRequestPolicyResult> createOriginRequestPolicy(
      $grpc.ServiceCall call, $0.CreateOriginRequestPolicyRequest request);

  $async.Future<$0.CreatePublicKeyResult> createPublicKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreatePublicKeyRequest> $request) async {
    return createPublicKey($call, await $request);
  }

  $async.Future<$0.CreatePublicKeyResult> createPublicKey(
      $grpc.ServiceCall call, $0.CreatePublicKeyRequest request);

  $async.Future<$0.CreateRealtimeLogConfigResult> createRealtimeLogConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateRealtimeLogConfigRequest> $request) async {
    return createRealtimeLogConfig($call, await $request);
  }

  $async.Future<$0.CreateRealtimeLogConfigResult> createRealtimeLogConfig(
      $grpc.ServiceCall call, $0.CreateRealtimeLogConfigRequest request);

  $async.Future<$0.CreateResponseHeadersPolicyResult>
      createResponseHeadersPolicy_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateResponseHeadersPolicyRequest> $request) async {
    return createResponseHeadersPolicy($call, await $request);
  }

  $async.Future<$0.CreateResponseHeadersPolicyResult>
      createResponseHeadersPolicy($grpc.ServiceCall call,
          $0.CreateResponseHeadersPolicyRequest request);

  $async.Future<$0.CreateStreamingDistributionResult>
      createStreamingDistribution_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateStreamingDistributionRequest> $request) async {
    return createStreamingDistribution($call, await $request);
  }

  $async.Future<$0.CreateStreamingDistributionResult>
      createStreamingDistribution($grpc.ServiceCall call,
          $0.CreateStreamingDistributionRequest request);

  $async.Future<$0.CreateStreamingDistributionWithTagsResult>
      createStreamingDistributionWithTags_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateStreamingDistributionWithTagsRequest>
              $request) async {
    return createStreamingDistributionWithTags($call, await $request);
  }

  $async.Future<$0.CreateStreamingDistributionWithTagsResult>
      createStreamingDistributionWithTags($grpc.ServiceCall call,
          $0.CreateStreamingDistributionWithTagsRequest request);

  $async.Future<$0.CreateTrustStoreResult> createTrustStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateTrustStoreRequest> $request) async {
    return createTrustStore($call, await $request);
  }

  $async.Future<$0.CreateTrustStoreResult> createTrustStore(
      $grpc.ServiceCall call, $0.CreateTrustStoreRequest request);

  $async.Future<$0.CreateVpcOriginResult> createVpcOrigin_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateVpcOriginRequest> $request) async {
    return createVpcOrigin($call, await $request);
  }

  $async.Future<$0.CreateVpcOriginResult> createVpcOrigin(
      $grpc.ServiceCall call, $0.CreateVpcOriginRequest request);

  $async.Future<$1.Empty> deleteAnycastIpList_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAnycastIpListRequest> $request) async {
    return deleteAnycastIpList($call, await $request);
  }

  $async.Future<$1.Empty> deleteAnycastIpList(
      $grpc.ServiceCall call, $0.DeleteAnycastIpListRequest request);

  $async.Future<$1.Empty> deleteCachePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteCachePolicyRequest> $request) async {
    return deleteCachePolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteCachePolicy(
      $grpc.ServiceCall call, $0.DeleteCachePolicyRequest request);

  $async.Future<$1.Empty> deleteCloudFrontOriginAccessIdentity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteCloudFrontOriginAccessIdentityRequest>
          $request) async {
    return deleteCloudFrontOriginAccessIdentity($call, await $request);
  }

  $async.Future<$1.Empty> deleteCloudFrontOriginAccessIdentity(
      $grpc.ServiceCall call,
      $0.DeleteCloudFrontOriginAccessIdentityRequest request);

  $async.Future<$1.Empty> deleteConnectionFunction_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteConnectionFunctionRequest> $request) async {
    return deleteConnectionFunction($call, await $request);
  }

  $async.Future<$1.Empty> deleteConnectionFunction(
      $grpc.ServiceCall call, $0.DeleteConnectionFunctionRequest request);

  $async.Future<$1.Empty> deleteConnectionGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteConnectionGroupRequest> $request) async {
    return deleteConnectionGroup($call, await $request);
  }

  $async.Future<$1.Empty> deleteConnectionGroup(
      $grpc.ServiceCall call, $0.DeleteConnectionGroupRequest request);

  $async.Future<$1.Empty> deleteContinuousDeploymentPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteContinuousDeploymentPolicyRequest>
          $request) async {
    return deleteContinuousDeploymentPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteContinuousDeploymentPolicy(
      $grpc.ServiceCall call,
      $0.DeleteContinuousDeploymentPolicyRequest request);

  $async.Future<$1.Empty> deleteDistribution_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDistributionRequest> $request) async {
    return deleteDistribution($call, await $request);
  }

  $async.Future<$1.Empty> deleteDistribution(
      $grpc.ServiceCall call, $0.DeleteDistributionRequest request);

  $async.Future<$1.Empty> deleteDistributionTenant_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDistributionTenantRequest> $request) async {
    return deleteDistributionTenant($call, await $request);
  }

  $async.Future<$1.Empty> deleteDistributionTenant(
      $grpc.ServiceCall call, $0.DeleteDistributionTenantRequest request);

  $async.Future<$1.Empty> deleteFieldLevelEncryptionConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteFieldLevelEncryptionConfigRequest>
          $request) async {
    return deleteFieldLevelEncryptionConfig($call, await $request);
  }

  $async.Future<$1.Empty> deleteFieldLevelEncryptionConfig(
      $grpc.ServiceCall call,
      $0.DeleteFieldLevelEncryptionConfigRequest request);

  $async.Future<$1.Empty> deleteFieldLevelEncryptionProfile_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteFieldLevelEncryptionProfileRequest>
          $request) async {
    return deleteFieldLevelEncryptionProfile($call, await $request);
  }

  $async.Future<$1.Empty> deleteFieldLevelEncryptionProfile(
      $grpc.ServiceCall call,
      $0.DeleteFieldLevelEncryptionProfileRequest request);

  $async.Future<$1.Empty> deleteFunction_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteFunctionRequest> $request) async {
    return deleteFunction($call, await $request);
  }

  $async.Future<$1.Empty> deleteFunction(
      $grpc.ServiceCall call, $0.DeleteFunctionRequest request);

  $async.Future<$1.Empty> deleteKeyGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteKeyGroupRequest> $request) async {
    return deleteKeyGroup($call, await $request);
  }

  $async.Future<$1.Empty> deleteKeyGroup(
      $grpc.ServiceCall call, $0.DeleteKeyGroupRequest request);

  $async.Future<$1.Empty> deleteKeyValueStore_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteKeyValueStoreRequest> $request) async {
    return deleteKeyValueStore($call, await $request);
  }

  $async.Future<$1.Empty> deleteKeyValueStore(
      $grpc.ServiceCall call, $0.DeleteKeyValueStoreRequest request);

  $async.Future<$0.DeleteMonitoringSubscriptionResult>
      deleteMonitoringSubscription_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeleteMonitoringSubscriptionRequest>
              $request) async {
    return deleteMonitoringSubscription($call, await $request);
  }

  $async.Future<$0.DeleteMonitoringSubscriptionResult>
      deleteMonitoringSubscription($grpc.ServiceCall call,
          $0.DeleteMonitoringSubscriptionRequest request);

  $async.Future<$1.Empty> deleteOriginAccessControl_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteOriginAccessControlRequest> $request) async {
    return deleteOriginAccessControl($call, await $request);
  }

  $async.Future<$1.Empty> deleteOriginAccessControl(
      $grpc.ServiceCall call, $0.DeleteOriginAccessControlRequest request);

  $async.Future<$1.Empty> deleteOriginRequestPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteOriginRequestPolicyRequest> $request) async {
    return deleteOriginRequestPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteOriginRequestPolicy(
      $grpc.ServiceCall call, $0.DeleteOriginRequestPolicyRequest request);

  $async.Future<$1.Empty> deletePublicKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeletePublicKeyRequest> $request) async {
    return deletePublicKey($call, await $request);
  }

  $async.Future<$1.Empty> deletePublicKey(
      $grpc.ServiceCall call, $0.DeletePublicKeyRequest request);

  $async.Future<$1.Empty> deleteRealtimeLogConfig_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteRealtimeLogConfigRequest> $request) async {
    return deleteRealtimeLogConfig($call, await $request);
  }

  $async.Future<$1.Empty> deleteRealtimeLogConfig(
      $grpc.ServiceCall call, $0.DeleteRealtimeLogConfigRequest request);

  $async.Future<$1.Empty> deleteResourcePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourcePolicyRequest> $request) async {
    return deleteResourcePolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteResourcePolicy(
      $grpc.ServiceCall call, $0.DeleteResourcePolicyRequest request);

  $async.Future<$1.Empty> deleteResponseHeadersPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteResponseHeadersPolicyRequest> $request) async {
    return deleteResponseHeadersPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteResponseHeadersPolicy(
      $grpc.ServiceCall call, $0.DeleteResponseHeadersPolicyRequest request);

  $async.Future<$1.Empty> deleteStreamingDistribution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteStreamingDistributionRequest> $request) async {
    return deleteStreamingDistribution($call, await $request);
  }

  $async.Future<$1.Empty> deleteStreamingDistribution(
      $grpc.ServiceCall call, $0.DeleteStreamingDistributionRequest request);

  $async.Future<$1.Empty> deleteTrustStore_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteTrustStoreRequest> $request) async {
    return deleteTrustStore($call, await $request);
  }

  $async.Future<$1.Empty> deleteTrustStore(
      $grpc.ServiceCall call, $0.DeleteTrustStoreRequest request);

  $async.Future<$0.DeleteVpcOriginResult> deleteVpcOrigin_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteVpcOriginRequest> $request) async {
    return deleteVpcOrigin($call, await $request);
  }

  $async.Future<$0.DeleteVpcOriginResult> deleteVpcOrigin(
      $grpc.ServiceCall call, $0.DeleteVpcOriginRequest request);

  $async.Future<$0.DescribeConnectionFunctionResult>
      describeConnectionFunction_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeConnectionFunctionRequest> $request) async {
    return describeConnectionFunction($call, await $request);
  }

  $async.Future<$0.DescribeConnectionFunctionResult> describeConnectionFunction(
      $grpc.ServiceCall call, $0.DescribeConnectionFunctionRequest request);

  $async.Future<$0.DescribeFunctionResult> describeFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeFunctionRequest> $request) async {
    return describeFunction($call, await $request);
  }

  $async.Future<$0.DescribeFunctionResult> describeFunction(
      $grpc.ServiceCall call, $0.DescribeFunctionRequest request);

  $async.Future<$0.DescribeKeyValueStoreResult> describeKeyValueStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeKeyValueStoreRequest> $request) async {
    return describeKeyValueStore($call, await $request);
  }

  $async.Future<$0.DescribeKeyValueStoreResult> describeKeyValueStore(
      $grpc.ServiceCall call, $0.DescribeKeyValueStoreRequest request);

  $async.Future<$0.DisassociateDistributionTenantWebACLResult>
      disassociateDistributionTenantWebACL_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DisassociateDistributionTenantWebACLRequest>
              $request) async {
    return disassociateDistributionTenantWebACL($call, await $request);
  }

  $async.Future<$0.DisassociateDistributionTenantWebACLResult>
      disassociateDistributionTenantWebACL($grpc.ServiceCall call,
          $0.DisassociateDistributionTenantWebACLRequest request);

  $async.Future<$0.DisassociateDistributionWebACLResult>
      disassociateDistributionWebACL_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DisassociateDistributionWebACLRequest>
              $request) async {
    return disassociateDistributionWebACL($call, await $request);
  }

  $async.Future<$0.DisassociateDistributionWebACLResult>
      disassociateDistributionWebACL($grpc.ServiceCall call,
          $0.DisassociateDistributionWebACLRequest request);

  $async.Future<$0.GetAnycastIpListResult> getAnycastIpList_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAnycastIpListRequest> $request) async {
    return getAnycastIpList($call, await $request);
  }

  $async.Future<$0.GetAnycastIpListResult> getAnycastIpList(
      $grpc.ServiceCall call, $0.GetAnycastIpListRequest request);

  $async.Future<$0.GetCachePolicyResult> getCachePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCachePolicyRequest> $request) async {
    return getCachePolicy($call, await $request);
  }

  $async.Future<$0.GetCachePolicyResult> getCachePolicy(
      $grpc.ServiceCall call, $0.GetCachePolicyRequest request);

  $async.Future<$0.GetCachePolicyConfigResult> getCachePolicyConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCachePolicyConfigRequest> $request) async {
    return getCachePolicyConfig($call, await $request);
  }

  $async.Future<$0.GetCachePolicyConfigResult> getCachePolicyConfig(
      $grpc.ServiceCall call, $0.GetCachePolicyConfigRequest request);

  $async.Future<$0.GetCloudFrontOriginAccessIdentityResult>
      getCloudFrontOriginAccessIdentity_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetCloudFrontOriginAccessIdentityRequest>
              $request) async {
    return getCloudFrontOriginAccessIdentity($call, await $request);
  }

  $async.Future<$0.GetCloudFrontOriginAccessIdentityResult>
      getCloudFrontOriginAccessIdentity($grpc.ServiceCall call,
          $0.GetCloudFrontOriginAccessIdentityRequest request);

  $async.Future<$0.GetCloudFrontOriginAccessIdentityConfigResult>
      getCloudFrontOriginAccessIdentityConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetCloudFrontOriginAccessIdentityConfigRequest>
              $request) async {
    return getCloudFrontOriginAccessIdentityConfig($call, await $request);
  }

  $async.Future<$0.GetCloudFrontOriginAccessIdentityConfigResult>
      getCloudFrontOriginAccessIdentityConfig($grpc.ServiceCall call,
          $0.GetCloudFrontOriginAccessIdentityConfigRequest request);

  $async.Future<$0.GetConnectionFunctionResult> getConnectionFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetConnectionFunctionRequest> $request) async {
    return getConnectionFunction($call, await $request);
  }

  $async.Future<$0.GetConnectionFunctionResult> getConnectionFunction(
      $grpc.ServiceCall call, $0.GetConnectionFunctionRequest request);

  $async.Future<$0.GetConnectionGroupResult> getConnectionGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetConnectionGroupRequest> $request) async {
    return getConnectionGroup($call, await $request);
  }

  $async.Future<$0.GetConnectionGroupResult> getConnectionGroup(
      $grpc.ServiceCall call, $0.GetConnectionGroupRequest request);

  $async.Future<$0.GetConnectionGroupByRoutingEndpointResult>
      getConnectionGroupByRoutingEndpoint_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetConnectionGroupByRoutingEndpointRequest>
              $request) async {
    return getConnectionGroupByRoutingEndpoint($call, await $request);
  }

  $async.Future<$0.GetConnectionGroupByRoutingEndpointResult>
      getConnectionGroupByRoutingEndpoint($grpc.ServiceCall call,
          $0.GetConnectionGroupByRoutingEndpointRequest request);

  $async.Future<$0.GetContinuousDeploymentPolicyResult>
      getContinuousDeploymentPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetContinuousDeploymentPolicyRequest>
              $request) async {
    return getContinuousDeploymentPolicy($call, await $request);
  }

  $async.Future<$0.GetContinuousDeploymentPolicyResult>
      getContinuousDeploymentPolicy($grpc.ServiceCall call,
          $0.GetContinuousDeploymentPolicyRequest request);

  $async.Future<$0.GetContinuousDeploymentPolicyConfigResult>
      getContinuousDeploymentPolicyConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetContinuousDeploymentPolicyConfigRequest>
              $request) async {
    return getContinuousDeploymentPolicyConfig($call, await $request);
  }

  $async.Future<$0.GetContinuousDeploymentPolicyConfigResult>
      getContinuousDeploymentPolicyConfig($grpc.ServiceCall call,
          $0.GetContinuousDeploymentPolicyConfigRequest request);

  $async.Future<$0.GetDistributionResult> getDistribution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDistributionRequest> $request) async {
    return getDistribution($call, await $request);
  }

  $async.Future<$0.GetDistributionResult> getDistribution(
      $grpc.ServiceCall call, $0.GetDistributionRequest request);

  $async.Future<$0.GetDistributionConfigResult> getDistributionConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDistributionConfigRequest> $request) async {
    return getDistributionConfig($call, await $request);
  }

  $async.Future<$0.GetDistributionConfigResult> getDistributionConfig(
      $grpc.ServiceCall call, $0.GetDistributionConfigRequest request);

  $async.Future<$0.GetDistributionTenantResult> getDistributionTenant_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDistributionTenantRequest> $request) async {
    return getDistributionTenant($call, await $request);
  }

  $async.Future<$0.GetDistributionTenantResult> getDistributionTenant(
      $grpc.ServiceCall call, $0.GetDistributionTenantRequest request);

  $async.Future<$0.GetDistributionTenantByDomainResult>
      getDistributionTenantByDomain_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetDistributionTenantByDomainRequest>
              $request) async {
    return getDistributionTenantByDomain($call, await $request);
  }

  $async.Future<$0.GetDistributionTenantByDomainResult>
      getDistributionTenantByDomain($grpc.ServiceCall call,
          $0.GetDistributionTenantByDomainRequest request);

  $async.Future<$0.GetFieldLevelEncryptionResult> getFieldLevelEncryption_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetFieldLevelEncryptionRequest> $request) async {
    return getFieldLevelEncryption($call, await $request);
  }

  $async.Future<$0.GetFieldLevelEncryptionResult> getFieldLevelEncryption(
      $grpc.ServiceCall call, $0.GetFieldLevelEncryptionRequest request);

  $async.Future<$0.GetFieldLevelEncryptionConfigResult>
      getFieldLevelEncryptionConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetFieldLevelEncryptionConfigRequest>
              $request) async {
    return getFieldLevelEncryptionConfig($call, await $request);
  }

  $async.Future<$0.GetFieldLevelEncryptionConfigResult>
      getFieldLevelEncryptionConfig($grpc.ServiceCall call,
          $0.GetFieldLevelEncryptionConfigRequest request);

  $async.Future<$0.GetFieldLevelEncryptionProfileResult>
      getFieldLevelEncryptionProfile_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetFieldLevelEncryptionProfileRequest>
              $request) async {
    return getFieldLevelEncryptionProfile($call, await $request);
  }

  $async.Future<$0.GetFieldLevelEncryptionProfileResult>
      getFieldLevelEncryptionProfile($grpc.ServiceCall call,
          $0.GetFieldLevelEncryptionProfileRequest request);

  $async.Future<$0.GetFieldLevelEncryptionProfileConfigResult>
      getFieldLevelEncryptionProfileConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetFieldLevelEncryptionProfileConfigRequest>
              $request) async {
    return getFieldLevelEncryptionProfileConfig($call, await $request);
  }

  $async.Future<$0.GetFieldLevelEncryptionProfileConfigResult>
      getFieldLevelEncryptionProfileConfig($grpc.ServiceCall call,
          $0.GetFieldLevelEncryptionProfileConfigRequest request);

  $async.Future<$0.GetFunctionResult> getFunction_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetFunctionRequest> $request) async {
    return getFunction($call, await $request);
  }

  $async.Future<$0.GetFunctionResult> getFunction(
      $grpc.ServiceCall call, $0.GetFunctionRequest request);

  $async.Future<$0.GetInvalidationResult> getInvalidation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetInvalidationRequest> $request) async {
    return getInvalidation($call, await $request);
  }

  $async.Future<$0.GetInvalidationResult> getInvalidation(
      $grpc.ServiceCall call, $0.GetInvalidationRequest request);

  $async.Future<$0.GetInvalidationForDistributionTenantResult>
      getInvalidationForDistributionTenant_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetInvalidationForDistributionTenantRequest>
              $request) async {
    return getInvalidationForDistributionTenant($call, await $request);
  }

  $async.Future<$0.GetInvalidationForDistributionTenantResult>
      getInvalidationForDistributionTenant($grpc.ServiceCall call,
          $0.GetInvalidationForDistributionTenantRequest request);

  $async.Future<$0.GetKeyGroupResult> getKeyGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetKeyGroupRequest> $request) async {
    return getKeyGroup($call, await $request);
  }

  $async.Future<$0.GetKeyGroupResult> getKeyGroup(
      $grpc.ServiceCall call, $0.GetKeyGroupRequest request);

  $async.Future<$0.GetKeyGroupConfigResult> getKeyGroupConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetKeyGroupConfigRequest> $request) async {
    return getKeyGroupConfig($call, await $request);
  }

  $async.Future<$0.GetKeyGroupConfigResult> getKeyGroupConfig(
      $grpc.ServiceCall call, $0.GetKeyGroupConfigRequest request);

  $async.Future<$0.GetManagedCertificateDetailsResult>
      getManagedCertificateDetails_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetManagedCertificateDetailsRequest>
              $request) async {
    return getManagedCertificateDetails($call, await $request);
  }

  $async.Future<$0.GetManagedCertificateDetailsResult>
      getManagedCertificateDetails($grpc.ServiceCall call,
          $0.GetManagedCertificateDetailsRequest request);

  $async.Future<$0.GetMonitoringSubscriptionResult>
      getMonitoringSubscription_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetMonitoringSubscriptionRequest> $request) async {
    return getMonitoringSubscription($call, await $request);
  }

  $async.Future<$0.GetMonitoringSubscriptionResult> getMonitoringSubscription(
      $grpc.ServiceCall call, $0.GetMonitoringSubscriptionRequest request);

  $async.Future<$0.GetOriginAccessControlResult> getOriginAccessControl_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetOriginAccessControlRequest> $request) async {
    return getOriginAccessControl($call, await $request);
  }

  $async.Future<$0.GetOriginAccessControlResult> getOriginAccessControl(
      $grpc.ServiceCall call, $0.GetOriginAccessControlRequest request);

  $async.Future<$0.GetOriginAccessControlConfigResult>
      getOriginAccessControlConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetOriginAccessControlConfigRequest>
              $request) async {
    return getOriginAccessControlConfig($call, await $request);
  }

  $async.Future<$0.GetOriginAccessControlConfigResult>
      getOriginAccessControlConfig($grpc.ServiceCall call,
          $0.GetOriginAccessControlConfigRequest request);

  $async.Future<$0.GetOriginRequestPolicyResult> getOriginRequestPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetOriginRequestPolicyRequest> $request) async {
    return getOriginRequestPolicy($call, await $request);
  }

  $async.Future<$0.GetOriginRequestPolicyResult> getOriginRequestPolicy(
      $grpc.ServiceCall call, $0.GetOriginRequestPolicyRequest request);

  $async.Future<$0.GetOriginRequestPolicyConfigResult>
      getOriginRequestPolicyConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetOriginRequestPolicyConfigRequest>
              $request) async {
    return getOriginRequestPolicyConfig($call, await $request);
  }

  $async.Future<$0.GetOriginRequestPolicyConfigResult>
      getOriginRequestPolicyConfig($grpc.ServiceCall call,
          $0.GetOriginRequestPolicyConfigRequest request);

  $async.Future<$0.GetPublicKeyResult> getPublicKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetPublicKeyRequest> $request) async {
    return getPublicKey($call, await $request);
  }

  $async.Future<$0.GetPublicKeyResult> getPublicKey(
      $grpc.ServiceCall call, $0.GetPublicKeyRequest request);

  $async.Future<$0.GetPublicKeyConfigResult> getPublicKeyConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPublicKeyConfigRequest> $request) async {
    return getPublicKeyConfig($call, await $request);
  }

  $async.Future<$0.GetPublicKeyConfigResult> getPublicKeyConfig(
      $grpc.ServiceCall call, $0.GetPublicKeyConfigRequest request);

  $async.Future<$0.GetRealtimeLogConfigResult> getRealtimeLogConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRealtimeLogConfigRequest> $request) async {
    return getRealtimeLogConfig($call, await $request);
  }

  $async.Future<$0.GetRealtimeLogConfigResult> getRealtimeLogConfig(
      $grpc.ServiceCall call, $0.GetRealtimeLogConfigRequest request);

  $async.Future<$0.GetResourcePolicyResult> getResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResourcePolicyRequest> $request) async {
    return getResourcePolicy($call, await $request);
  }

  $async.Future<$0.GetResourcePolicyResult> getResourcePolicy(
      $grpc.ServiceCall call, $0.GetResourcePolicyRequest request);

  $async.Future<$0.GetResponseHeadersPolicyResult> getResponseHeadersPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResponseHeadersPolicyRequest> $request) async {
    return getResponseHeadersPolicy($call, await $request);
  }

  $async.Future<$0.GetResponseHeadersPolicyResult> getResponseHeadersPolicy(
      $grpc.ServiceCall call, $0.GetResponseHeadersPolicyRequest request);

  $async.Future<$0.GetResponseHeadersPolicyConfigResult>
      getResponseHeadersPolicyConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetResponseHeadersPolicyConfigRequest>
              $request) async {
    return getResponseHeadersPolicyConfig($call, await $request);
  }

  $async.Future<$0.GetResponseHeadersPolicyConfigResult>
      getResponseHeadersPolicyConfig($grpc.ServiceCall call,
          $0.GetResponseHeadersPolicyConfigRequest request);

  $async.Future<$0.GetStreamingDistributionResult> getStreamingDistribution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetStreamingDistributionRequest> $request) async {
    return getStreamingDistribution($call, await $request);
  }

  $async.Future<$0.GetStreamingDistributionResult> getStreamingDistribution(
      $grpc.ServiceCall call, $0.GetStreamingDistributionRequest request);

  $async.Future<$0.GetStreamingDistributionConfigResult>
      getStreamingDistributionConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetStreamingDistributionConfigRequest>
              $request) async {
    return getStreamingDistributionConfig($call, await $request);
  }

  $async.Future<$0.GetStreamingDistributionConfigResult>
      getStreamingDistributionConfig($grpc.ServiceCall call,
          $0.GetStreamingDistributionConfigRequest request);

  $async.Future<$0.GetTrustStoreResult> getTrustStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetTrustStoreRequest> $request) async {
    return getTrustStore($call, await $request);
  }

  $async.Future<$0.GetTrustStoreResult> getTrustStore(
      $grpc.ServiceCall call, $0.GetTrustStoreRequest request);

  $async.Future<$0.GetVpcOriginResult> getVpcOrigin_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetVpcOriginRequest> $request) async {
    return getVpcOrigin($call, await $request);
  }

  $async.Future<$0.GetVpcOriginResult> getVpcOrigin(
      $grpc.ServiceCall call, $0.GetVpcOriginRequest request);

  $async.Future<$0.ListAnycastIpListsResult> listAnycastIpLists_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAnycastIpListsRequest> $request) async {
    return listAnycastIpLists($call, await $request);
  }

  $async.Future<$0.ListAnycastIpListsResult> listAnycastIpLists(
      $grpc.ServiceCall call, $0.ListAnycastIpListsRequest request);

  $async.Future<$0.ListCachePoliciesResult> listCachePolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCachePoliciesRequest> $request) async {
    return listCachePolicies($call, await $request);
  }

  $async.Future<$0.ListCachePoliciesResult> listCachePolicies(
      $grpc.ServiceCall call, $0.ListCachePoliciesRequest request);

  $async.Future<$0.ListCloudFrontOriginAccessIdentitiesResult>
      listCloudFrontOriginAccessIdentities_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListCloudFrontOriginAccessIdentitiesRequest>
              $request) async {
    return listCloudFrontOriginAccessIdentities($call, await $request);
  }

  $async.Future<$0.ListCloudFrontOriginAccessIdentitiesResult>
      listCloudFrontOriginAccessIdentities($grpc.ServiceCall call,
          $0.ListCloudFrontOriginAccessIdentitiesRequest request);

  $async.Future<$0.ListConflictingAliasesResult> listConflictingAliases_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListConflictingAliasesRequest> $request) async {
    return listConflictingAliases($call, await $request);
  }

  $async.Future<$0.ListConflictingAliasesResult> listConflictingAliases(
      $grpc.ServiceCall call, $0.ListConflictingAliasesRequest request);

  $async.Future<$0.ListConnectionFunctionsResult> listConnectionFunctions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListConnectionFunctionsRequest> $request) async {
    return listConnectionFunctions($call, await $request);
  }

  $async.Future<$0.ListConnectionFunctionsResult> listConnectionFunctions(
      $grpc.ServiceCall call, $0.ListConnectionFunctionsRequest request);

  $async.Future<$0.ListConnectionGroupsResult> listConnectionGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListConnectionGroupsRequest> $request) async {
    return listConnectionGroups($call, await $request);
  }

  $async.Future<$0.ListConnectionGroupsResult> listConnectionGroups(
      $grpc.ServiceCall call, $0.ListConnectionGroupsRequest request);

  $async.Future<$0.ListContinuousDeploymentPoliciesResult>
      listContinuousDeploymentPolicies_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListContinuousDeploymentPoliciesRequest>
              $request) async {
    return listContinuousDeploymentPolicies($call, await $request);
  }

  $async.Future<$0.ListContinuousDeploymentPoliciesResult>
      listContinuousDeploymentPolicies($grpc.ServiceCall call,
          $0.ListContinuousDeploymentPoliciesRequest request);

  $async.Future<$0.ListDistributionsResult> listDistributions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDistributionsRequest> $request) async {
    return listDistributions($call, await $request);
  }

  $async.Future<$0.ListDistributionsResult> listDistributions(
      $grpc.ServiceCall call, $0.ListDistributionsRequest request);

  $async.Future<$0.ListDistributionsByAnycastIpListIdResult>
      listDistributionsByAnycastIpListId_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByAnycastIpListIdRequest>
              $request) async {
    return listDistributionsByAnycastIpListId($call, await $request);
  }

  $async.Future<$0.ListDistributionsByAnycastIpListIdResult>
      listDistributionsByAnycastIpListId($grpc.ServiceCall call,
          $0.ListDistributionsByAnycastIpListIdRequest request);

  $async.Future<$0.ListDistributionsByCachePolicyIdResult>
      listDistributionsByCachePolicyId_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByCachePolicyIdRequest>
              $request) async {
    return listDistributionsByCachePolicyId($call, await $request);
  }

  $async.Future<$0.ListDistributionsByCachePolicyIdResult>
      listDistributionsByCachePolicyId($grpc.ServiceCall call,
          $0.ListDistributionsByCachePolicyIdRequest request);

  $async.Future<$0.ListDistributionsByConnectionFunctionResult>
      listDistributionsByConnectionFunction_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByConnectionFunctionRequest>
              $request) async {
    return listDistributionsByConnectionFunction($call, await $request);
  }

  $async.Future<$0.ListDistributionsByConnectionFunctionResult>
      listDistributionsByConnectionFunction($grpc.ServiceCall call,
          $0.ListDistributionsByConnectionFunctionRequest request);

  $async.Future<$0.ListDistributionsByConnectionModeResult>
      listDistributionsByConnectionMode_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByConnectionModeRequest>
              $request) async {
    return listDistributionsByConnectionMode($call, await $request);
  }

  $async.Future<$0.ListDistributionsByConnectionModeResult>
      listDistributionsByConnectionMode($grpc.ServiceCall call,
          $0.ListDistributionsByConnectionModeRequest request);

  $async.Future<$0.ListDistributionsByKeyGroupResult>
      listDistributionsByKeyGroup_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByKeyGroupRequest> $request) async {
    return listDistributionsByKeyGroup($call, await $request);
  }

  $async.Future<$0.ListDistributionsByKeyGroupResult>
      listDistributionsByKeyGroup($grpc.ServiceCall call,
          $0.ListDistributionsByKeyGroupRequest request);

  $async.Future<$0.ListDistributionsByOriginRequestPolicyIdResult>
      listDistributionsByOriginRequestPolicyId_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByOriginRequestPolicyIdRequest>
              $request) async {
    return listDistributionsByOriginRequestPolicyId($call, await $request);
  }

  $async.Future<$0.ListDistributionsByOriginRequestPolicyIdResult>
      listDistributionsByOriginRequestPolicyId($grpc.ServiceCall call,
          $0.ListDistributionsByOriginRequestPolicyIdRequest request);

  $async.Future<$0.ListDistributionsByOwnedResourceResult>
      listDistributionsByOwnedResource_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByOwnedResourceRequest>
              $request) async {
    return listDistributionsByOwnedResource($call, await $request);
  }

  $async.Future<$0.ListDistributionsByOwnedResourceResult>
      listDistributionsByOwnedResource($grpc.ServiceCall call,
          $0.ListDistributionsByOwnedResourceRequest request);

  $async.Future<$0.ListDistributionsByRealtimeLogConfigResult>
      listDistributionsByRealtimeLogConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByRealtimeLogConfigRequest>
              $request) async {
    return listDistributionsByRealtimeLogConfig($call, await $request);
  }

  $async.Future<$0.ListDistributionsByRealtimeLogConfigResult>
      listDistributionsByRealtimeLogConfig($grpc.ServiceCall call,
          $0.ListDistributionsByRealtimeLogConfigRequest request);

  $async.Future<$0.ListDistributionsByResponseHeadersPolicyIdResult>
      listDistributionsByResponseHeadersPolicyId_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByResponseHeadersPolicyIdRequest>
              $request) async {
    return listDistributionsByResponseHeadersPolicyId($call, await $request);
  }

  $async.Future<$0.ListDistributionsByResponseHeadersPolicyIdResult>
      listDistributionsByResponseHeadersPolicyId($grpc.ServiceCall call,
          $0.ListDistributionsByResponseHeadersPolicyIdRequest request);

  $async.Future<$0.ListDistributionsByTrustStoreResult>
      listDistributionsByTrustStore_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByTrustStoreRequest>
              $request) async {
    return listDistributionsByTrustStore($call, await $request);
  }

  $async.Future<$0.ListDistributionsByTrustStoreResult>
      listDistributionsByTrustStore($grpc.ServiceCall call,
          $0.ListDistributionsByTrustStoreRequest request);

  $async.Future<$0.ListDistributionsByVpcOriginIdResult>
      listDistributionsByVpcOriginId_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByVpcOriginIdRequest>
              $request) async {
    return listDistributionsByVpcOriginId($call, await $request);
  }

  $async.Future<$0.ListDistributionsByVpcOriginIdResult>
      listDistributionsByVpcOriginId($grpc.ServiceCall call,
          $0.ListDistributionsByVpcOriginIdRequest request);

  $async.Future<$0.ListDistributionsByWebACLIdResult>
      listDistributionsByWebACLId_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionsByWebACLIdRequest> $request) async {
    return listDistributionsByWebACLId($call, await $request);
  }

  $async.Future<$0.ListDistributionsByWebACLIdResult>
      listDistributionsByWebACLId($grpc.ServiceCall call,
          $0.ListDistributionsByWebACLIdRequest request);

  $async.Future<$0.ListDistributionTenantsResult> listDistributionTenants_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDistributionTenantsRequest> $request) async {
    return listDistributionTenants($call, await $request);
  }

  $async.Future<$0.ListDistributionTenantsResult> listDistributionTenants(
      $grpc.ServiceCall call, $0.ListDistributionTenantsRequest request);

  $async.Future<$0.ListDistributionTenantsByCustomizationResult>
      listDistributionTenantsByCustomization_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDistributionTenantsByCustomizationRequest>
              $request) async {
    return listDistributionTenantsByCustomization($call, await $request);
  }

  $async.Future<$0.ListDistributionTenantsByCustomizationResult>
      listDistributionTenantsByCustomization($grpc.ServiceCall call,
          $0.ListDistributionTenantsByCustomizationRequest request);

  $async.Future<$0.ListDomainConflictsResult> listDomainConflicts_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDomainConflictsRequest> $request) async {
    return listDomainConflicts($call, await $request);
  }

  $async.Future<$0.ListDomainConflictsResult> listDomainConflicts(
      $grpc.ServiceCall call, $0.ListDomainConflictsRequest request);

  $async.Future<$0.ListFieldLevelEncryptionConfigsResult>
      listFieldLevelEncryptionConfigs_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListFieldLevelEncryptionConfigsRequest>
              $request) async {
    return listFieldLevelEncryptionConfigs($call, await $request);
  }

  $async.Future<$0.ListFieldLevelEncryptionConfigsResult>
      listFieldLevelEncryptionConfigs($grpc.ServiceCall call,
          $0.ListFieldLevelEncryptionConfigsRequest request);

  $async.Future<$0.ListFieldLevelEncryptionProfilesResult>
      listFieldLevelEncryptionProfiles_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListFieldLevelEncryptionProfilesRequest>
              $request) async {
    return listFieldLevelEncryptionProfiles($call, await $request);
  }

  $async.Future<$0.ListFieldLevelEncryptionProfilesResult>
      listFieldLevelEncryptionProfiles($grpc.ServiceCall call,
          $0.ListFieldLevelEncryptionProfilesRequest request);

  $async.Future<$0.ListFunctionsResult> listFunctions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListFunctionsRequest> $request) async {
    return listFunctions($call, await $request);
  }

  $async.Future<$0.ListFunctionsResult> listFunctions(
      $grpc.ServiceCall call, $0.ListFunctionsRequest request);

  $async.Future<$0.ListInvalidationsResult> listInvalidations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListInvalidationsRequest> $request) async {
    return listInvalidations($call, await $request);
  }

  $async.Future<$0.ListInvalidationsResult> listInvalidations(
      $grpc.ServiceCall call, $0.ListInvalidationsRequest request);

  $async.Future<$0.ListInvalidationsForDistributionTenantResult>
      listInvalidationsForDistributionTenant_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListInvalidationsForDistributionTenantRequest>
              $request) async {
    return listInvalidationsForDistributionTenant($call, await $request);
  }

  $async.Future<$0.ListInvalidationsForDistributionTenantResult>
      listInvalidationsForDistributionTenant($grpc.ServiceCall call,
          $0.ListInvalidationsForDistributionTenantRequest request);

  $async.Future<$0.ListKeyGroupsResult> listKeyGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListKeyGroupsRequest> $request) async {
    return listKeyGroups($call, await $request);
  }

  $async.Future<$0.ListKeyGroupsResult> listKeyGroups(
      $grpc.ServiceCall call, $0.ListKeyGroupsRequest request);

  $async.Future<$0.ListKeyValueStoresResult> listKeyValueStores_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListKeyValueStoresRequest> $request) async {
    return listKeyValueStores($call, await $request);
  }

  $async.Future<$0.ListKeyValueStoresResult> listKeyValueStores(
      $grpc.ServiceCall call, $0.ListKeyValueStoresRequest request);

  $async.Future<$0.ListOriginAccessControlsResult> listOriginAccessControls_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListOriginAccessControlsRequest> $request) async {
    return listOriginAccessControls($call, await $request);
  }

  $async.Future<$0.ListOriginAccessControlsResult> listOriginAccessControls(
      $grpc.ServiceCall call, $0.ListOriginAccessControlsRequest request);

  $async.Future<$0.ListOriginRequestPoliciesResult>
      listOriginRequestPolicies_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListOriginRequestPoliciesRequest> $request) async {
    return listOriginRequestPolicies($call, await $request);
  }

  $async.Future<$0.ListOriginRequestPoliciesResult> listOriginRequestPolicies(
      $grpc.ServiceCall call, $0.ListOriginRequestPoliciesRequest request);

  $async.Future<$0.ListPublicKeysResult> listPublicKeys_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPublicKeysRequest> $request) async {
    return listPublicKeys($call, await $request);
  }

  $async.Future<$0.ListPublicKeysResult> listPublicKeys(
      $grpc.ServiceCall call, $0.ListPublicKeysRequest request);

  $async.Future<$0.ListRealtimeLogConfigsResult> listRealtimeLogConfigs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRealtimeLogConfigsRequest> $request) async {
    return listRealtimeLogConfigs($call, await $request);
  }

  $async.Future<$0.ListRealtimeLogConfigsResult> listRealtimeLogConfigs(
      $grpc.ServiceCall call, $0.ListRealtimeLogConfigsRequest request);

  $async.Future<$0.ListResponseHeadersPoliciesResult>
      listResponseHeadersPolicies_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListResponseHeadersPoliciesRequest> $request) async {
    return listResponseHeadersPolicies($call, await $request);
  }

  $async.Future<$0.ListResponseHeadersPoliciesResult>
      listResponseHeadersPolicies($grpc.ServiceCall call,
          $0.ListResponseHeadersPoliciesRequest request);

  $async.Future<$0.ListStreamingDistributionsResult>
      listStreamingDistributions_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListStreamingDistributionsRequest> $request) async {
    return listStreamingDistributions($call, await $request);
  }

  $async.Future<$0.ListStreamingDistributionsResult> listStreamingDistributions(
      $grpc.ServiceCall call, $0.ListStreamingDistributionsRequest request);

  $async.Future<$0.ListTagsForResourceResult> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResult> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ListTrustStoresResult> listTrustStores_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTrustStoresRequest> $request) async {
    return listTrustStores($call, await $request);
  }

  $async.Future<$0.ListTrustStoresResult> listTrustStores(
      $grpc.ServiceCall call, $0.ListTrustStoresRequest request);

  $async.Future<$0.ListVpcOriginsResult> listVpcOrigins_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListVpcOriginsRequest> $request) async {
    return listVpcOrigins($call, await $request);
  }

  $async.Future<$0.ListVpcOriginsResult> listVpcOrigins(
      $grpc.ServiceCall call, $0.ListVpcOriginsRequest request);

  $async.Future<$0.PublishConnectionFunctionResult>
      publishConnectionFunction_Pre($grpc.ServiceCall $call,
          $async.Future<$0.PublishConnectionFunctionRequest> $request) async {
    return publishConnectionFunction($call, await $request);
  }

  $async.Future<$0.PublishConnectionFunctionResult> publishConnectionFunction(
      $grpc.ServiceCall call, $0.PublishConnectionFunctionRequest request);

  $async.Future<$0.PublishFunctionResult> publishFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PublishFunctionRequest> $request) async {
    return publishFunction($call, await $request);
  }

  $async.Future<$0.PublishFunctionResult> publishFunction(
      $grpc.ServiceCall call, $0.PublishFunctionRequest request);

  $async.Future<$0.PutResourcePolicyResult> putResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutResourcePolicyRequest> $request) async {
    return putResourcePolicy($call, await $request);
  }

  $async.Future<$0.PutResourcePolicyResult> putResourcePolicy(
      $grpc.ServiceCall call, $0.PutResourcePolicyRequest request);

  $async.Future<$1.Empty> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$1.Empty> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$0.TestConnectionFunctionResult> testConnectionFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TestConnectionFunctionRequest> $request) async {
    return testConnectionFunction($call, await $request);
  }

  $async.Future<$0.TestConnectionFunctionResult> testConnectionFunction(
      $grpc.ServiceCall call, $0.TestConnectionFunctionRequest request);

  $async.Future<$0.TestFunctionResult> testFunction_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TestFunctionRequest> $request) async {
    return testFunction($call, await $request);
  }

  $async.Future<$0.TestFunctionResult> testFunction(
      $grpc.ServiceCall call, $0.TestFunctionRequest request);

  $async.Future<$1.Empty> untagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$1.Empty> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.UpdateAnycastIpListResult> updateAnycastIpList_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateAnycastIpListRequest> $request) async {
    return updateAnycastIpList($call, await $request);
  }

  $async.Future<$0.UpdateAnycastIpListResult> updateAnycastIpList(
      $grpc.ServiceCall call, $0.UpdateAnycastIpListRequest request);

  $async.Future<$0.UpdateCachePolicyResult> updateCachePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateCachePolicyRequest> $request) async {
    return updateCachePolicy($call, await $request);
  }

  $async.Future<$0.UpdateCachePolicyResult> updateCachePolicy(
      $grpc.ServiceCall call, $0.UpdateCachePolicyRequest request);

  $async.Future<$0.UpdateCloudFrontOriginAccessIdentityResult>
      updateCloudFrontOriginAccessIdentity_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateCloudFrontOriginAccessIdentityRequest>
              $request) async {
    return updateCloudFrontOriginAccessIdentity($call, await $request);
  }

  $async.Future<$0.UpdateCloudFrontOriginAccessIdentityResult>
      updateCloudFrontOriginAccessIdentity($grpc.ServiceCall call,
          $0.UpdateCloudFrontOriginAccessIdentityRequest request);

  $async.Future<$0.UpdateConnectionFunctionResult> updateConnectionFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateConnectionFunctionRequest> $request) async {
    return updateConnectionFunction($call, await $request);
  }

  $async.Future<$0.UpdateConnectionFunctionResult> updateConnectionFunction(
      $grpc.ServiceCall call, $0.UpdateConnectionFunctionRequest request);

  $async.Future<$0.UpdateConnectionGroupResult> updateConnectionGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateConnectionGroupRequest> $request) async {
    return updateConnectionGroup($call, await $request);
  }

  $async.Future<$0.UpdateConnectionGroupResult> updateConnectionGroup(
      $grpc.ServiceCall call, $0.UpdateConnectionGroupRequest request);

  $async.Future<$0.UpdateContinuousDeploymentPolicyResult>
      updateContinuousDeploymentPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateContinuousDeploymentPolicyRequest>
              $request) async {
    return updateContinuousDeploymentPolicy($call, await $request);
  }

  $async.Future<$0.UpdateContinuousDeploymentPolicyResult>
      updateContinuousDeploymentPolicy($grpc.ServiceCall call,
          $0.UpdateContinuousDeploymentPolicyRequest request);

  $async.Future<$0.UpdateDistributionResult> updateDistribution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDistributionRequest> $request) async {
    return updateDistribution($call, await $request);
  }

  $async.Future<$0.UpdateDistributionResult> updateDistribution(
      $grpc.ServiceCall call, $0.UpdateDistributionRequest request);

  $async.Future<$0.UpdateDistributionTenantResult> updateDistributionTenant_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDistributionTenantRequest> $request) async {
    return updateDistributionTenant($call, await $request);
  }

  $async.Future<$0.UpdateDistributionTenantResult> updateDistributionTenant(
      $grpc.ServiceCall call, $0.UpdateDistributionTenantRequest request);

  $async.Future<$0.UpdateDistributionWithStagingConfigResult>
      updateDistributionWithStagingConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateDistributionWithStagingConfigRequest>
              $request) async {
    return updateDistributionWithStagingConfig($call, await $request);
  }

  $async.Future<$0.UpdateDistributionWithStagingConfigResult>
      updateDistributionWithStagingConfig($grpc.ServiceCall call,
          $0.UpdateDistributionWithStagingConfigRequest request);

  $async.Future<$0.UpdateDomainAssociationResult> updateDomainAssociation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDomainAssociationRequest> $request) async {
    return updateDomainAssociation($call, await $request);
  }

  $async.Future<$0.UpdateDomainAssociationResult> updateDomainAssociation(
      $grpc.ServiceCall call, $0.UpdateDomainAssociationRequest request);

  $async.Future<$0.UpdateFieldLevelEncryptionConfigResult>
      updateFieldLevelEncryptionConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateFieldLevelEncryptionConfigRequest>
              $request) async {
    return updateFieldLevelEncryptionConfig($call, await $request);
  }

  $async.Future<$0.UpdateFieldLevelEncryptionConfigResult>
      updateFieldLevelEncryptionConfig($grpc.ServiceCall call,
          $0.UpdateFieldLevelEncryptionConfigRequest request);

  $async.Future<$0.UpdateFieldLevelEncryptionProfileResult>
      updateFieldLevelEncryptionProfile_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateFieldLevelEncryptionProfileRequest>
              $request) async {
    return updateFieldLevelEncryptionProfile($call, await $request);
  }

  $async.Future<$0.UpdateFieldLevelEncryptionProfileResult>
      updateFieldLevelEncryptionProfile($grpc.ServiceCall call,
          $0.UpdateFieldLevelEncryptionProfileRequest request);

  $async.Future<$0.UpdateFunctionResult> updateFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateFunctionRequest> $request) async {
    return updateFunction($call, await $request);
  }

  $async.Future<$0.UpdateFunctionResult> updateFunction(
      $grpc.ServiceCall call, $0.UpdateFunctionRequest request);

  $async.Future<$0.UpdateKeyGroupResult> updateKeyGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateKeyGroupRequest> $request) async {
    return updateKeyGroup($call, await $request);
  }

  $async.Future<$0.UpdateKeyGroupResult> updateKeyGroup(
      $grpc.ServiceCall call, $0.UpdateKeyGroupRequest request);

  $async.Future<$0.UpdateKeyValueStoreResult> updateKeyValueStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateKeyValueStoreRequest> $request) async {
    return updateKeyValueStore($call, await $request);
  }

  $async.Future<$0.UpdateKeyValueStoreResult> updateKeyValueStore(
      $grpc.ServiceCall call, $0.UpdateKeyValueStoreRequest request);

  $async.Future<$0.UpdateOriginAccessControlResult>
      updateOriginAccessControl_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateOriginAccessControlRequest> $request) async {
    return updateOriginAccessControl($call, await $request);
  }

  $async.Future<$0.UpdateOriginAccessControlResult> updateOriginAccessControl(
      $grpc.ServiceCall call, $0.UpdateOriginAccessControlRequest request);

  $async.Future<$0.UpdateOriginRequestPolicyResult>
      updateOriginRequestPolicy_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateOriginRequestPolicyRequest> $request) async {
    return updateOriginRequestPolicy($call, await $request);
  }

  $async.Future<$0.UpdateOriginRequestPolicyResult> updateOriginRequestPolicy(
      $grpc.ServiceCall call, $0.UpdateOriginRequestPolicyRequest request);

  $async.Future<$0.UpdatePublicKeyResult> updatePublicKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdatePublicKeyRequest> $request) async {
    return updatePublicKey($call, await $request);
  }

  $async.Future<$0.UpdatePublicKeyResult> updatePublicKey(
      $grpc.ServiceCall call, $0.UpdatePublicKeyRequest request);

  $async.Future<$0.UpdateRealtimeLogConfigResult> updateRealtimeLogConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRealtimeLogConfigRequest> $request) async {
    return updateRealtimeLogConfig($call, await $request);
  }

  $async.Future<$0.UpdateRealtimeLogConfigResult> updateRealtimeLogConfig(
      $grpc.ServiceCall call, $0.UpdateRealtimeLogConfigRequest request);

  $async.Future<$0.UpdateResponseHeadersPolicyResult>
      updateResponseHeadersPolicy_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateResponseHeadersPolicyRequest> $request) async {
    return updateResponseHeadersPolicy($call, await $request);
  }

  $async.Future<$0.UpdateResponseHeadersPolicyResult>
      updateResponseHeadersPolicy($grpc.ServiceCall call,
          $0.UpdateResponseHeadersPolicyRequest request);

  $async.Future<$0.UpdateStreamingDistributionResult>
      updateStreamingDistribution_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateStreamingDistributionRequest> $request) async {
    return updateStreamingDistribution($call, await $request);
  }

  $async.Future<$0.UpdateStreamingDistributionResult>
      updateStreamingDistribution($grpc.ServiceCall call,
          $0.UpdateStreamingDistributionRequest request);

  $async.Future<$0.UpdateTrustStoreResult> updateTrustStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateTrustStoreRequest> $request) async {
    return updateTrustStore($call, await $request);
  }

  $async.Future<$0.UpdateTrustStoreResult> updateTrustStore(
      $grpc.ServiceCall call, $0.UpdateTrustStoreRequest request);

  $async.Future<$0.UpdateVpcOriginResult> updateVpcOrigin_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateVpcOriginRequest> $request) async {
    return updateVpcOrigin($call, await $request);
  }

  $async.Future<$0.UpdateVpcOriginResult> updateVpcOrigin(
      $grpc.ServiceCall call, $0.UpdateVpcOriginRequest request);

  $async.Future<$0.VerifyDnsConfigurationResult> verifyDnsConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.VerifyDnsConfigurationRequest> $request) async {
    return verifyDnsConfiguration($call, await $request);
  }

  $async.Future<$0.VerifyDnsConfigurationResult> verifyDnsConfiguration(
      $grpc.ServiceCall call, $0.VerifyDnsConfigurationRequest request);
}
