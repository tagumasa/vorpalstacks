// This is a generated file - do not edit.
//
// Generated from s3.proto.

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
import 's3.pb.dart' as $0;

export 's3.pb.dart';

/// S3Service provides s3 API operations.
@$pb.GrpcServiceName('s3.S3Service')
class S3ServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  S3ServiceClient(super.channel, {super.options, super.interceptors});

  /// This operation aborts a multipart upload. After a multipart upload is aborted, no additional parts can be uploaded using that upload ID. The storage consumed by any previously uploaded parts will b...
  /// HTTP: DELETE /{Bucket}/{Key+}?x-id=AbortMultipartUpload
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.AbortMultipartUploadOutput> abortMultipartUpload(
    $0.AbortMultipartUploadRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$abortMultipartUpload, request, options: options);
  }

  /// Completes a multipart upload by assembling previously uploaded parts. You first initiate the multipart upload and then upload all parts using the UploadPart operation or the UploadPartCopy operatio...
  /// HTTP: POST /{Bucket}/{Key+}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CompleteMultipartUploadOutput>
      completeMultipartUpload(
    $0.CompleteMultipartUploadRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$completeMultipartUpload, request,
        options: options);
  }

  /// Creates a copy of an object that is already stored in Amazon S3. End of support notice: As of October 1, 2025, Amazon S3 has discontinued support for Email Grantee Access Control Lists (ACLs). If y...
  /// HTTP: PUT /{Bucket}/{Key+}?x-id=CopyObject
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CopyObjectOutput> copyObject(
    $0.CopyObjectRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$copyObject, request, options: options);
  }

  /// This action creates an Amazon S3 bucket. To create an Amazon S3 on Outposts bucket, see CreateBucket . Creates a new S3 bucket. To create a bucket, you must set up Amazon S3 and have a valid Amazon...
  /// HTTP: PUT /{Bucket}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateBucketOutput> createBucket(
    $0.CreateBucketRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createBucket, request, options: options);
  }

  /// Creates an S3 Metadata V2 metadata configuration for a general purpose bucket. For more information, see Accelerating data discovery with S3 Metadata in the Amazon S3 User Guide. Permissions To use...
  /// HTTP: POST /{Bucket}?metadataConfiguration
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> createBucketMetadataConfiguration(
    $0.CreateBucketMetadataConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createBucketMetadataConfiguration, request,
        options: options);
  }

  /// We recommend that you create your S3 Metadata configurations by using the V2 CreateBucketMetadataConfiguration API operation. We no longer recommend using the V1 CreateBucketMetadataTableConfigurat...
  /// HTTP: POST /{Bucket}?metadataTable
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> createBucketMetadataTableConfiguration(
    $0.CreateBucketMetadataTableConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createBucketMetadataTableConfiguration, request,
        options: options);
  }

  /// End of support notice: As of October 1, 2025, Amazon S3 has discontinued support for Email Grantee Access Control Lists (ACLs). If you attempt to use an Email Grantee ACL in a request after October...
  /// HTTP: POST /{Bucket}/{Key+}?uploads
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateMultipartUploadOutput> createMultipartUpload(
    $0.CreateMultipartUploadRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createMultipartUpload, request, options: options);
  }

  /// Creates a session that establishes temporary security credentials to support fast authentication and authorization for the Zonal endpoint API operations on directory buckets. For more information a...
  /// HTTP: GET /{Bucket}?session
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateSessionOutput> createSession(
    $0.CreateSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSession, request, options: options);
  }

  /// Deletes the S3 bucket. All objects (including all object versions and delete markers) in the bucket must be deleted before the bucket itself can be deleted. Directory buckets - If multipart uploads...
  /// HTTP: DELETE /{Bucket}
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucket(
    $0.DeleteBucketRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucket, request, options: options);
  }

  /// This operation is not supported for directory buckets. Deletes an analytics configuration for the bucket (specified by the analytics configuration ID). To use this operation, you must have permissi...
  /// HTTP: DELETE /{Bucket}?analytics
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketAnalyticsConfiguration(
    $0.DeleteBucketAnalyticsConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketAnalyticsConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Deletes the cors configuration information set for the bucket. To use this operation, you must have permission to perform the s3:PutBucketCORS...
  /// HTTP: DELETE /{Bucket}?cors
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketCors(
    $0.DeleteBucketCorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketCors, request, options: options);
  }

  /// This implementation of the DELETE action resets the default encryption for the bucket as server-side encryption with Amazon S3 managed keys (SSE-S3). General purpose buckets - For information about...
  /// HTTP: DELETE /{Bucket}?encryption
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketEncryption(
    $0.DeleteBucketEncryptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketEncryption, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Deletes the S3 Intelligent-Tiering configuration from the specified bucket. The S3 Intelligent-Tiering storage class is designed to optimize s...
  /// HTTP: DELETE /{Bucket}?intelligent-tiering
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketIntelligentTieringConfiguration(
    $0.DeleteBucketIntelligentTieringConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$deleteBucketIntelligentTieringConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Deletes an S3 Inventory configuration (identified by the inventory ID) from the bucket. To use this operation, you must have permissions to pe...
  /// HTTP: DELETE /{Bucket}?inventory
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketInventoryConfiguration(
    $0.DeleteBucketInventoryConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketInventoryConfiguration, request,
        options: options);
  }

  /// Deletes the lifecycle configuration from the specified bucket. Amazon S3 removes all the lifecycle configuration rules in the lifecycle subresource associated with the bucket. Your objects never ex...
  /// HTTP: DELETE /{Bucket}?lifecycle
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketLifecycle(
    $0.DeleteBucketLifecycleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketLifecycle, request, options: options);
  }

  /// Deletes an S3 Metadata configuration from a general purpose bucket. For more information, see Accelerating data discovery with S3 Metadata in the Amazon S3 User Guide. You can use the V2 DeleteBuck...
  /// HTTP: DELETE /{Bucket}?metadataConfiguration
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketMetadataConfiguration(
    $0.DeleteBucketMetadataConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketMetadataConfiguration, request,
        options: options);
  }

  /// We recommend that you delete your S3 Metadata configurations by using the V2 DeleteBucketMetadataTableConfiguration API operation. We no longer recommend using the V1 DeleteBucketMetadataTableConfi...
  /// HTTP: DELETE /{Bucket}?metadataTable
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketMetadataTableConfiguration(
    $0.DeleteBucketMetadataTableConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketMetadataTableConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Deletes a metrics configuration for the Amazon CloudWatch request metrics (specified by the metrics configuration ID) from the bucket. Note th...
  /// HTTP: DELETE /{Bucket}?metrics
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketMetricsConfiguration(
    $0.DeleteBucketMetricsConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketMetricsConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Removes OwnershipControls for an Amazon S3 bucket. To use this operation, you must have the s3:PutBucketOwnershipControls permission. For more...
  /// HTTP: DELETE /{Bucket}?ownershipControls
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketOwnershipControls(
    $0.DeleteBucketOwnershipControlsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketOwnershipControls, request,
        options: options);
  }

  /// Deletes the policy of a specified bucket. Directory buckets - For directory buckets, you must make requests for this API operation to the Regional endpoint. These endpoints support path-style reque...
  /// HTTP: DELETE /{Bucket}?policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketPolicy(
    $0.DeleteBucketPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketPolicy, request, options: options);
  }

  /// This operation is not supported for directory buckets. Deletes the replication configuration from the bucket. To use this operation, you must have permissions to perform the s3:PutReplicationConfig...
  /// HTTP: DELETE /{Bucket}?replication
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketReplication(
    $0.DeleteBucketReplicationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketReplication, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Deletes tags from the general purpose bucket if attribute based access control (ABAC) is not enabled for the bucket. When you enable ABAC for ...
  /// HTTP: DELETE /{Bucket}?tagging
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketTagging(
    $0.DeleteBucketTaggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketTagging, request, options: options);
  }

  /// This operation is not supported for directory buckets. This action removes the website configuration for a bucket. Amazon S3 returns a 200 OK response upon successfully deleting a website configura...
  /// HTTP: DELETE /{Bucket}?website
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deleteBucketWebsite(
    $0.DeleteBucketWebsiteRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBucketWebsite, request, options: options);
  }

  /// Removes an object from a bucket. The behavior depends on the bucket's versioning state: If bucket versioning is not enabled, the operation permanently deletes the object. If bucket versioning is en...
  /// HTTP: DELETE /{Bucket}/{Key+}?x-id=DeleteObject
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteObjectOutput> deleteObject(
    $0.DeleteObjectRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteObject, request, options: options);
  }

  /// This operation enables you to delete multiple objects from a bucket using a single HTTP request. If you know the object keys that you want to delete, then this operation provides a suitable alterna...
  /// HTTP: POST /{Bucket}?delete
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteObjectsOutput> deleteObjects(
    $0.DeleteObjectsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteObjects, request, options: options);
  }

  /// This operation is not supported for directory buckets. Removes the entire tag set from the specified object. For more information about managing object tags, see Object Tagging. To use this operati...
  /// HTTP: DELETE /{Bucket}/{Key+}?tagging
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteObjectTaggingOutput> deleteObjectTagging(
    $0.DeleteObjectTaggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteObjectTagging, request, options: options);
  }

  /// This operation is not supported for directory buckets. Removes the PublicAccessBlock configuration for an Amazon S3 bucket. This operation removes the bucket-level configuration only. The effective...
  /// HTTP: DELETE /{Bucket}?publicAccessBlock
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> deletePublicAccessBlock(
    $0.DeletePublicAccessBlockRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePublicAccessBlock, request,
        options: options);
  }

  /// Returns the attribute-based access control (ABAC) property of the general purpose bucket. If ABAC is enabled on your bucket, you can use tags on the bucket for access control. For more information,...
  /// HTTP: GET /{Bucket}?abac
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketAbacOutput> getBucketAbac(
    $0.GetBucketAbacRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketAbac, request, options: options);
  }

  /// This operation is not supported for directory buckets. This implementation of the GET action uses the accelerate subresource to return the Transfer Acceleration state of a bucket, which is either E...
  /// HTTP: GET /{Bucket}?accelerate
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketAccelerateConfigurationOutput>
      getBucketAccelerateConfiguration(
    $0.GetBucketAccelerateConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketAccelerateConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. This implementation of the GET action uses the acl subresource to return the access control list (ACL) of a bucket. To use GET to return the A...
  /// HTTP: GET /{Bucket}?acl
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketAclOutput> getBucketAcl(
    $0.GetBucketAclRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketAcl, request, options: options);
  }

  /// This operation is not supported for directory buckets. This implementation of the GET action returns an analytics configuration (identified by the analytics configuration ID) from the bucket. To us...
  /// HTTP: GET /{Bucket}?analytics&x-id=GetBucketAnalyticsConfiguration
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketAnalyticsConfigurationOutput>
      getBucketAnalyticsConfiguration(
    $0.GetBucketAnalyticsConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketAnalyticsConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Returns the Cross-Origin Resource Sharing (CORS) configuration information set for the bucket. To use this operation, you must have permission...
  /// HTTP: GET /{Bucket}?cors
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketCorsOutput> getBucketCors(
    $0.GetBucketCorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketCors, request, options: options);
  }

  /// Returns the default encryption configuration for an Amazon S3 bucket. By default, all buckets have a default encryption configuration that uses server-side encryption with Amazon S3 managed keys (S...
  /// HTTP: GET /{Bucket}?encryption
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketEncryptionOutput> getBucketEncryption(
    $0.GetBucketEncryptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketEncryption, request, options: options);
  }

  /// This operation is not supported for directory buckets. Gets the S3 Intelligent-Tiering configuration from the specified bucket. The S3 Intelligent-Tiering storage class is designed to optimize stor...
  /// HTTP: GET /{Bucket}?intelligent-tiering&x-id=GetBucketIntelligentTieringConfiguration
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketIntelligentTieringConfigurationOutput>
      getBucketIntelligentTieringConfiguration(
    $0.GetBucketIntelligentTieringConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketIntelligentTieringConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Returns an S3 Inventory configuration (identified by the inventory configuration ID) from the bucket. To use this operation, you must have per...
  /// HTTP: GET /{Bucket}?inventory&x-id=GetBucketInventoryConfiguration
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketInventoryConfigurationOutput>
      getBucketInventoryConfiguration(
    $0.GetBucketInventoryConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketInventoryConfiguration, request,
        options: options);
  }

  /// Returns the lifecycle configuration information set on the bucket. For information about lifecycle configuration, see Object Lifecycle Management. Bucket lifecycle configuration now supports specif...
  /// HTTP: GET /{Bucket}?lifecycle
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketLifecycleConfigurationOutput>
      getBucketLifecycleConfiguration(
    $0.GetBucketLifecycleConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketLifecycleConfiguration, request,
        options: options);
  }

  /// Using the GetBucketLocation operation is no longer a best practice. To return the Region that a bucket resides in, we recommend that you use the HeadBucket operation instead. For backward compatibi...
  /// HTTP: GET /{Bucket}?location
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketLocationOutput> getBucketLocation(
    $0.GetBucketLocationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketLocation, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns the logging status of a bucket and the permissions users have to view and modify that status. The following operations are related to ...
  /// HTTP: GET /{Bucket}?logging
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketLoggingOutput> getBucketLogging(
    $0.GetBucketLoggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketLogging, request, options: options);
  }

  /// Retrieves the S3 Metadata configuration for a general purpose bucket. For more information, see Accelerating data discovery with S3 Metadata in the Amazon S3 User Guide. You can use the V2 GetBucke...
  /// HTTP: GET /{Bucket}?metadataConfiguration
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketMetadataConfigurationOutput>
      getBucketMetadataConfiguration(
    $0.GetBucketMetadataConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketMetadataConfiguration, request,
        options: options);
  }

  /// We recommend that you retrieve your S3 Metadata configurations by using the V2 GetBucketMetadataTableConfiguration API operation. We no longer recommend using the V1 GetBucketMetadataTableConfigura...
  /// HTTP: GET /{Bucket}?metadataTable
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketMetadataTableConfigurationOutput>
      getBucketMetadataTableConfiguration(
    $0.GetBucketMetadataTableConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketMetadataTableConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Gets a metrics configuration (specified by the metrics configuration ID) from the bucket. Note that this doesn't include the daily storage met...
  /// HTTP: GET /{Bucket}?metrics&x-id=GetBucketMetricsConfiguration
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketMetricsConfigurationOutput>
      getBucketMetricsConfiguration(
    $0.GetBucketMetricsConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketMetricsConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Returns the notification configuration of a bucket. If notifications are not enabled on the bucket, the action returns an empty NotificationCo...
  /// HTTP: GET /{Bucket}?notification
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.NotificationConfiguration>
      getBucketNotificationConfiguration(
    $0.GetBucketNotificationConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketNotificationConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Retrieves OwnershipControls for an Amazon S3 bucket. To use this operation, you must have the s3:GetBucketOwnershipControls permission. For mo...
  /// HTTP: GET /{Bucket}?ownershipControls
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketOwnershipControlsOutput>
      getBucketOwnershipControls(
    $0.GetBucketOwnershipControlsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketOwnershipControls, request,
        options: options);
  }

  /// Returns the policy of a specified bucket. Directory buckets - For directory buckets, you must make requests for this API operation to the Regional endpoint. These endpoints support path-style reque...
  /// HTTP: GET /{Bucket}?policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketPolicyOutput> getBucketPolicy(
    $0.GetBucketPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketPolicy, request, options: options);
  }

  /// This operation is not supported for directory buckets. Retrieves the policy status for an Amazon S3 bucket, indicating whether the bucket is public. In order to use this operation, you must have th...
  /// HTTP: GET /{Bucket}?policyStatus
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketPolicyStatusOutput> getBucketPolicyStatus(
    $0.GetBucketPolicyStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketPolicyStatus, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns the replication configuration of a bucket. It can take a while to propagate the put or delete a replication configuration to all Amazo...
  /// HTTP: GET /{Bucket}?replication
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketReplicationOutput> getBucketReplication(
    $0.GetBucketReplicationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketReplication, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns the request payment configuration of a bucket. To use this version of the operation, you must be the bucket owner. For more informatio...
  /// HTTP: GET /{Bucket}?requestPayment
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketRequestPaymentOutput>
      getBucketRequestPayment(
    $0.GetBucketRequestPaymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketRequestPayment, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Returns the tag set associated with the general purpose bucket. if ABAC is not enabled for the bucket. When you enable ABAC for a general purp...
  /// HTTP: GET /{Bucket}?tagging
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketTaggingOutput> getBucketTagging(
    $0.GetBucketTaggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketTagging, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns the versioning state of a bucket. To retrieve the versioning state of a bucket, you must be the bucket owner. This implementation also...
  /// HTTP: GET /{Bucket}?versioning
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketVersioningOutput> getBucketVersioning(
    $0.GetBucketVersioningRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketVersioning, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns the website configuration for a bucket. To host website on Amazon S3, you can configure a bucket as website by adding a website config...
  /// HTTP: GET /{Bucket}?website
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetBucketWebsiteOutput> getBucketWebsite(
    $0.GetBucketWebsiteRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBucketWebsite, request, options: options);
  }

  /// Retrieves an object from Amazon S3. In the GetObject request, specify the full key name for the object. General purpose buckets - Both the virtual-hosted-style requests and the path-style requests ...
  /// HTTP: GET /{Bucket}/{Key+}?x-id=GetObject
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetObjectOutput> getObject(
    $0.GetObjectRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getObject, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns the access control list (ACL) of an object. To use this operation, you must have s3:GetObjectAcl permissions or READ_ACP access to the...
  /// HTTP: GET /{Bucket}/{Key+}?acl
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetObjectAclOutput> getObjectAcl(
    $0.GetObjectAclRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getObjectAcl, request, options: options);
  }

  /// Retrieves all of the metadata from an object without returning the object itself. This operation is useful if you're interested only in an object's metadata. GetObjectAttributes combines the functi...
  /// HTTP: GET /{Bucket}/{Key+}?attributes
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetObjectAttributesOutput> getObjectAttributes(
    $0.GetObjectAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getObjectAttributes, request, options: options);
  }

  /// This operation is not supported for directory buckets. Gets an object's current legal hold status. For more information, see Locking Objects. This functionality is not supported for Amazon S3 on Ou...
  /// HTTP: GET /{Bucket}/{Key+}?legal-hold
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetObjectLegalHoldOutput> getObjectLegalHold(
    $0.GetObjectLegalHoldRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getObjectLegalHold, request, options: options);
  }

  /// This operation is not supported for directory buckets. Gets the Object Lock configuration for a bucket. The rule specified in the Object Lock configuration will be applied by default to every new o...
  /// HTTP: GET /{Bucket}?object-lock
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetObjectLockConfigurationOutput>
      getObjectLockConfiguration(
    $0.GetObjectLockConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getObjectLockConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Retrieves an object's retention settings. For more information, see Locking Objects. This functionality is not supported for Amazon S3 on Outp...
  /// HTTP: GET /{Bucket}/{Key+}?retention
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetObjectRetentionOutput> getObjectRetention(
    $0.GetObjectRetentionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getObjectRetention, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns the tag-set of an object. You send the GET request against the tagging subresource associated with the object. To use this operation, ...
  /// HTTP: GET /{Bucket}/{Key+}?tagging
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetObjectTaggingOutput> getObjectTagging(
    $0.GetObjectTaggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getObjectTagging, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns torrent files from a bucket. BitTorrent can save you bandwidth when you're distributing large files. You can get torrent only for obje...
  /// HTTP: GET /{Bucket}/{Key+}?torrent
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetObjectTorrentOutput> getObjectTorrent(
    $0.GetObjectTorrentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getObjectTorrent, request, options: options);
  }

  /// This operation is not supported for directory buckets. Retrieves the PublicAccessBlock configuration for an Amazon S3 bucket. This operation returns the bucket-level configuration only. To understa...
  /// HTTP: GET /{Bucket}?publicAccessBlock
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetPublicAccessBlockOutput> getPublicAccessBlock(
    $0.GetPublicAccessBlockRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPublicAccessBlock, request, options: options);
  }

  /// You can use this operation to determine if a bucket exists and if you have permission to access it. The action returns a 200 OK HTTP status code if the bucket exists and you have permission to acce...
  /// HTTP: HEAD /{Bucket}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.HeadBucketOutput> headBucket(
    $0.HeadBucketRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$headBucket, request, options: options);
  }

  /// The HEAD operation retrieves metadata from an object without returning the object itself. This operation is useful if you're interested only in an object's metadata. A HEAD request has the same opt...
  /// HTTP: HEAD /{Bucket}/{Key+}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.HeadObjectOutput> headObject(
    $0.HeadObjectRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$headObject, request, options: options);
  }

  /// This operation is not supported for directory buckets. Lists the analytics configurations for the bucket. You can have up to 1,000 analytics configurations per bucket. This action supports list pag...
  /// HTTP: GET /{Bucket}?analytics&x-id=ListBucketAnalyticsConfigurations
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListBucketAnalyticsConfigurationsOutput>
      listBucketAnalyticsConfigurations(
    $0.ListBucketAnalyticsConfigurationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listBucketAnalyticsConfigurations, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Lists the S3 Intelligent-Tiering configuration from the specified bucket. The S3 Intelligent-Tiering storage class is designed to optimize sto...
  /// HTTP: GET /{Bucket}?intelligent-tiering&x-id=ListBucketIntelligentTieringConfigurations
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListBucketIntelligentTieringConfigurationsOutput>
      listBucketIntelligentTieringConfigurations(
    $0.ListBucketIntelligentTieringConfigurationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$listBucketIntelligentTieringConfigurations, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Returns a list of S3 Inventory configurations for the bucket. You can have up to 1,000 inventory configurations per bucket. This action suppor...
  /// HTTP: GET /{Bucket}?inventory&x-id=ListBucketInventoryConfigurations
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListBucketInventoryConfigurationsOutput>
      listBucketInventoryConfigurations(
    $0.ListBucketInventoryConfigurationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listBucketInventoryConfigurations, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Lists the metrics configurations for the bucket. The metrics configurations are only for the request metrics of the bucket and do not provide ...
  /// HTTP: GET /{Bucket}?metrics&x-id=ListBucketMetricsConfigurations
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListBucketMetricsConfigurationsOutput>
      listBucketMetricsConfigurations(
    $0.ListBucketMetricsConfigurationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listBucketMetricsConfigurations, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Returns a list of all buckets owned by the authenticated sender of the request. To grant IAM permission to use this operation, you must add th...
  /// HTTP: GET /?x-id=ListBuckets
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListBucketsOutput> listBuckets(
    $0.ListBucketsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listBuckets, request, options: options);
  }

  /// Returns a list of all Amazon S3 directory buckets owned by the authenticated sender of the request. For more information about directory buckets, see Directory buckets in the Amazon S3 User Guide. ...
  /// HTTP: GET /?x-id=ListDirectoryBuckets
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListDirectoryBucketsOutput> listDirectoryBuckets(
    $0.ListDirectoryBucketsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDirectoryBuckets, request, options: options);
  }

  /// This operation lists in-progress multipart uploads in a bucket. An in-progress multipart upload is a multipart upload that has been initiated by the CreateMultipartUpload request, but has not yet b...
  /// HTTP: GET /{Bucket}?uploads
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListMultipartUploadsOutput> listMultipartUploads(
    $0.ListMultipartUploadsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMultipartUploads, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns some or all (up to 1,000) of the objects in a bucket. You can use the request parameters as selection criteria to return a subset of t...
  /// HTTP: GET /{Bucket}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListObjectsOutput> listObjects(
    $0.ListObjectsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listObjects, request, options: options);
  }

  /// Returns some or all (up to 1,000) of the objects in a bucket with each request. You can use the request parameters as selection criteria to return a subset of the objects in a bucket. A 200 OK resp...
  /// HTTP: GET /{Bucket}?list-type=2
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListObjectsV2Output> listObjectsV2(
    $0.ListObjectsV2Request request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listObjectsV2, request, options: options);
  }

  /// This operation is not supported for directory buckets. Returns metadata about all versions of the objects in a bucket. You can also use request parameters as selection criteria to return metadata a...
  /// HTTP: GET /{Bucket}?versions
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListObjectVersionsOutput> listObjectVersions(
    $0.ListObjectVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listObjectVersions, request, options: options);
  }

  /// Lists the parts that have been uploaded for a specific multipart upload. To use this operation, you must provide the upload ID in the request. You obtain this uploadID by sending the initiate multi...
  /// HTTP: GET /{Bucket}/{Key+}?x-id=ListParts
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListPartsOutput> listParts(
    $0.ListPartsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listParts, request, options: options);
  }

  /// Sets the attribute-based access control (ABAC) property of the general purpose bucket. You must have s3:PutBucketABAC permission to perform this action. When you enable ABAC, you can use tags for a...
  /// HTTP: PUT /{Bucket}?abac
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketAbac(
    $0.PutBucketAbacRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketAbac, request, options: options);
  }

  /// This operation is not supported for directory buckets. Sets the accelerate configuration of an existing bucket. Amazon S3 Transfer Acceleration is a bucket-level feature that enables you to perform...
  /// HTTP: PUT /{Bucket}?accelerate
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketAccelerateConfiguration(
    $0.PutBucketAccelerateConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketAccelerateConfiguration, request,
        options: options);
  }

  /// End of support notice: As of October 1, 2025, Amazon S3 has discontinued support for Email Grantee Access Control Lists (ACLs). If you attempt to use an Email Grantee ACL in a request after October...
  /// HTTP: PUT /{Bucket}?acl
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketAcl(
    $0.PutBucketAclRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketAcl, request, options: options);
  }

  /// This operation is not supported for directory buckets. Sets an analytics configuration for the bucket (specified by the analytics configuration ID). You can have up to 1,000 analytics configuration...
  /// HTTP: PUT /{Bucket}?analytics
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketAnalyticsConfiguration(
    $0.PutBucketAnalyticsConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketAnalyticsConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Sets the cors configuration for your bucket. If the configuration exists, Amazon S3 replaces it. To use this operation, you must be allowed to...
  /// HTTP: PUT /{Bucket}?cors
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketCors(
    $0.PutBucketCorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketCors, request, options: options);
  }

  /// This operation configures default encryption and Amazon S3 Bucket Keys for an existing bucket. You can also block encryption types using this operation. Directory buckets - For directory buckets, y...
  /// HTTP: PUT /{Bucket}?encryption
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketEncryption(
    $0.PutBucketEncryptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketEncryption, request, options: options);
  }

  /// This operation is not supported for directory buckets. Puts a S3 Intelligent-Tiering configuration to the specified bucket. You can have up to 1,000 S3 Intelligent-Tiering configurations per bucket...
  /// HTTP: PUT /{Bucket}?intelligent-tiering
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketIntelligentTieringConfiguration(
    $0.PutBucketIntelligentTieringConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketIntelligentTieringConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. This implementation of the PUT action adds an S3 Inventory configuration (identified by the inventory ID) to the bucket. You can have up to 1,...
  /// HTTP: PUT /{Bucket}?inventory
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketInventoryConfiguration(
    $0.PutBucketInventoryConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketInventoryConfiguration, request,
        options: options);
  }

  /// Creates a new lifecycle configuration for the bucket or replaces an existing lifecycle configuration. Keep in mind that this will overwrite an existing lifecycle configuration, so if you want to re...
  /// HTTP: PUT /{Bucket}?lifecycle
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PutBucketLifecycleConfigurationOutput>
      putBucketLifecycleConfiguration(
    $0.PutBucketLifecycleConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketLifecycleConfiguration, request,
        options: options);
  }

  /// End of support notice: As of October 1, 2025, Amazon S3 has discontinued support for Email Grantee Access Control Lists (ACLs). If you attempt to use an Email Grantee ACL in a request after October...
  /// HTTP: PUT /{Bucket}?logging
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketLogging(
    $0.PutBucketLoggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketLogging, request, options: options);
  }

  /// This operation is not supported for directory buckets. Sets a metrics configuration (specified by the metrics configuration ID) for the bucket. You can have up to 1,000 metrics configurations per b...
  /// HTTP: PUT /{Bucket}?metrics
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketMetricsConfiguration(
    $0.PutBucketMetricsConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketMetricsConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Enables notifications of specified events for a bucket. For more information about event notifications, see Configuring Event Notifications. U...
  /// HTTP: PUT /{Bucket}?notification
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketNotificationConfiguration(
    $0.PutBucketNotificationConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketNotificationConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Creates or modifies OwnershipControls for an Amazon S3 bucket. To use this operation, you must have the s3:PutBucketOwnershipControls permissi...
  /// HTTP: PUT /{Bucket}?ownershipControls
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketOwnershipControls(
    $0.PutBucketOwnershipControlsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketOwnershipControls, request,
        options: options);
  }

  /// Applies an Amazon S3 bucket policy to an Amazon S3 bucket. Directory buckets - For directory buckets, you must make requests for this API operation to the Regional endpoint. These endpoints support...
  /// HTTP: PUT /{Bucket}?policy
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketPolicy(
    $0.PutBucketPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketPolicy, request, options: options);
  }

  /// This operation is not supported for directory buckets. Creates a replication configuration or replaces an existing one. For more information, see Replication in the Amazon S3 User Guide. Specify th...
  /// HTTP: PUT /{Bucket}?replication
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketReplication(
    $0.PutBucketReplicationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketReplication, request, options: options);
  }

  /// This operation is not supported for directory buckets. Sets the request payment configuration for a bucket. By default, the bucket owner pays for downloads from the bucket. This configuration param...
  /// HTTP: PUT /{Bucket}?requestPayment
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketRequestPayment(
    $0.PutBucketRequestPaymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketRequestPayment, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Sets the tags for a general purpose bucket if attribute based access control (ABAC) is not enabled for the bucket. When you enable ABAC for a ...
  /// HTTP: PUT /{Bucket}?tagging
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketTagging(
    $0.PutBucketTaggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketTagging, request, options: options);
  }

  /// This operation is not supported for directory buckets. When you enable versioning on a bucket for the first time, it might take a short amount of time for the change to be fully propagated. While t...
  /// HTTP: PUT /{Bucket}?versioning
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketVersioning(
    $0.PutBucketVersioningRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketVersioning, request, options: options);
  }

  /// This operation is not supported for directory buckets. Sets the configuration of the website that is specified in the website subresource. To configure a bucket as a website, you can add this subre...
  /// HTTP: PUT /{Bucket}?website
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putBucketWebsite(
    $0.PutBucketWebsiteRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putBucketWebsite, request, options: options);
  }

  /// End of support notice: As of October 1, 2025, Amazon S3 has discontinued support for Email Grantee Access Control Lists (ACLs). If you attempt to use an Email Grantee ACL in a request after October...
  /// HTTP: PUT /{Bucket}/{Key+}?x-id=PutObject
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PutObjectOutput> putObject(
    $0.PutObjectRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putObject, request, options: options);
  }

  /// End of support notice: As of October 1, 2025, Amazon S3 has discontinued support for Email Grantee Access Control Lists (ACLs). If you attempt to use an Email Grantee ACL in a request after October...
  /// HTTP: PUT /{Bucket}/{Key+}?acl
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PutObjectAclOutput> putObjectAcl(
    $0.PutObjectAclRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putObjectAcl, request, options: options);
  }

  /// This operation is not supported for directory buckets. Applies a legal hold configuration to the specified object. For more information, see Locking Objects. This functionality is not supported for...
  /// HTTP: PUT /{Bucket}/{Key+}?legal-hold
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PutObjectLegalHoldOutput> putObjectLegalHold(
    $0.PutObjectLegalHoldRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putObjectLegalHold, request, options: options);
  }

  /// This operation is not supported for directory buckets. Places an Object Lock configuration on the specified bucket. The rule specified in the Object Lock configuration will be applied by default to...
  /// HTTP: PUT /{Bucket}?object-lock
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PutObjectLockConfigurationOutput>
      putObjectLockConfiguration(
    $0.PutObjectLockConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putObjectLockConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets. Places an Object Retention configuration on an object. For more information, see Locking Objects. Users or accounts require the s3:PutObjectRe...
  /// HTTP: PUT /{Bucket}/{Key+}?retention
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PutObjectRetentionOutput> putObjectRetention(
    $0.PutObjectRetentionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putObjectRetention, request, options: options);
  }

  /// This operation is not supported for directory buckets. Sets the supplied tag-set to an object that already exists in a bucket. A tag is a key-value pair. For more information, see Object Tagging. Y...
  /// HTTP: PUT /{Bucket}/{Key+}?tagging
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.PutObjectTaggingOutput> putObjectTagging(
    $0.PutObjectTaggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putObjectTagging, request, options: options);
  }

  /// This operation is not supported for directory buckets. Creates or modifies the PublicAccessBlock configuration for an Amazon S3 bucket. To use this operation, you must have the s3:PutBucketPublicAc...
  /// HTTP: PUT /{Bucket}?publicAccessBlock
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> putPublicAccessBlock(
    $0.PutPublicAccessBlockRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putPublicAccessBlock, request, options: options);
  }

  /// Renames an existing object in a directory bucket that uses the S3 Express One Zone storage class. You can use RenameObject by specifying an existing object’s name as the source and the new name o...
  /// HTTP: PUT /{Bucket}/{Key+}?renameObject
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.RenameObjectOutput> renameObject(
    $0.RenameObjectRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$renameObject, request, options: options);
  }

  /// This operation is not supported for directory buckets. Restores an archived copy of an object back into Amazon S3 This functionality is not supported for Amazon S3 on Outposts. This action performs...
  /// HTTP: POST /{Bucket}/{Key+}?restore
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.RestoreObjectOutput> restoreObject(
    $0.RestoreObjectRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$restoreObject, request, options: options);
  }

  /// This operation is not supported for directory buckets. This action filters the contents of an Amazon S3 object based on a simple structured query language (SQL) statement. In the request, along wit...
  /// HTTP: POST /{Bucket}/{Key+}?select&select-type=2
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.SelectObjectContentOutput> selectObjectContent(
    $0.SelectObjectContentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$selectObjectContent, request, options: options);
  }

  /// Enables or disables a live inventory table for an S3 Metadata configuration on a general purpose bucket. For more information, see Accelerating data discovery with S3 Metadata in the Amazon S3 User...
  /// HTTP: PUT /{Bucket}?metadataInventoryTable
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty>
      updateBucketMetadataInventoryTableConfiguration(
    $0.UpdateBucketMetadataInventoryTableConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$updateBucketMetadataInventoryTableConfiguration, request,
        options: options);
  }

  /// Enables or disables journal table record expiration for an S3 Metadata configuration on a general purpose bucket. For more information, see Accelerating data discovery with S3 Metadata in the Amazo...
  /// HTTP: PUT /{Bucket}?metadataJournalTable
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> updateBucketMetadataJournalTableConfiguration(
    $0.UpdateBucketMetadataJournalTableConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$updateBucketMetadataJournalTableConfiguration, request,
        options: options);
  }

  /// This operation is not supported for directory buckets or Amazon S3 on Outposts buckets. Updates the server-side encryption type of an existing encrypted object in a general purpose bucket. You can ...
  /// HTTP: PUT /{Bucket}/{Key+}?encryption
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateObjectEncryptionResponse>
      updateObjectEncryption(
    $0.UpdateObjectEncryptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateObjectEncryption, request,
        options: options);
  }

  /// Uploads a part in a multipart upload. In this operation, you provide new data as a part of an object in your request. However, you have an option to specify your existing Amazon S3 object as a data...
  /// HTTP: PUT /{Bucket}/{Key+}?x-id=UploadPart
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UploadPartOutput> uploadPart(
    $0.UploadPartRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$uploadPart, request, options: options);
  }

  /// Uploads a part by copying data from an existing object as data source. To specify the data source, you add the request header x-amz-copy-source in your request. To specify a byte range, you add the...
  /// HTTP: PUT /{Bucket}/{Key+}?x-id=UploadPartCopy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UploadPartCopyOutput> uploadPartCopy(
    $0.UploadPartCopyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$uploadPartCopy, request, options: options);
  }

  /// This operation is not supported for directory buckets. Passes transformed objects to a GetObject operation when using Object Lambda access points. For information about Object Lambda access points,...
  /// HTTP: POST /WriteGetObjectResponse
  /// Protocol: restXml
  $grpc.ResponseFuture<$1.Empty> writeGetObjectResponse(
    $0.WriteGetObjectResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$writeGetObjectResponse, request,
        options: options);
  }

  // method descriptors

  static final _$abortMultipartUpload = $grpc.ClientMethod<
          $0.AbortMultipartUploadRequest, $0.AbortMultipartUploadOutput>(
      '/s3.S3Service/AbortMultipartUpload',
      ($0.AbortMultipartUploadRequest value) => value.writeToBuffer(),
      $0.AbortMultipartUploadOutput.fromBuffer);
  static final _$completeMultipartUpload = $grpc.ClientMethod<
          $0.CompleteMultipartUploadRequest, $0.CompleteMultipartUploadOutput>(
      '/s3.S3Service/CompleteMultipartUpload',
      ($0.CompleteMultipartUploadRequest value) => value.writeToBuffer(),
      $0.CompleteMultipartUploadOutput.fromBuffer);
  static final _$copyObject =
      $grpc.ClientMethod<$0.CopyObjectRequest, $0.CopyObjectOutput>(
          '/s3.S3Service/CopyObject',
          ($0.CopyObjectRequest value) => value.writeToBuffer(),
          $0.CopyObjectOutput.fromBuffer);
  static final _$createBucket =
      $grpc.ClientMethod<$0.CreateBucketRequest, $0.CreateBucketOutput>(
          '/s3.S3Service/CreateBucket',
          ($0.CreateBucketRequest value) => value.writeToBuffer(),
          $0.CreateBucketOutput.fromBuffer);
  static final _$createBucketMetadataConfiguration =
      $grpc.ClientMethod<$0.CreateBucketMetadataConfigurationRequest, $1.Empty>(
          '/s3.S3Service/CreateBucketMetadataConfiguration',
          ($0.CreateBucketMetadataConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$createBucketMetadataTableConfiguration = $grpc.ClientMethod<
          $0.CreateBucketMetadataTableConfigurationRequest, $1.Empty>(
      '/s3.S3Service/CreateBucketMetadataTableConfiguration',
      ($0.CreateBucketMetadataTableConfigurationRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$createMultipartUpload = $grpc.ClientMethod<
          $0.CreateMultipartUploadRequest, $0.CreateMultipartUploadOutput>(
      '/s3.S3Service/CreateMultipartUpload',
      ($0.CreateMultipartUploadRequest value) => value.writeToBuffer(),
      $0.CreateMultipartUploadOutput.fromBuffer);
  static final _$createSession =
      $grpc.ClientMethod<$0.CreateSessionRequest, $0.CreateSessionOutput>(
          '/s3.S3Service/CreateSession',
          ($0.CreateSessionRequest value) => value.writeToBuffer(),
          $0.CreateSessionOutput.fromBuffer);
  static final _$deleteBucket =
      $grpc.ClientMethod<$0.DeleteBucketRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucket',
          ($0.DeleteBucketRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketAnalyticsConfiguration = $grpc.ClientMethod<
          $0.DeleteBucketAnalyticsConfigurationRequest, $1.Empty>(
      '/s3.S3Service/DeleteBucketAnalyticsConfiguration',
      ($0.DeleteBucketAnalyticsConfigurationRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$deleteBucketCors =
      $grpc.ClientMethod<$0.DeleteBucketCorsRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketCors',
          ($0.DeleteBucketCorsRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketEncryption =
      $grpc.ClientMethod<$0.DeleteBucketEncryptionRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketEncryption',
          ($0.DeleteBucketEncryptionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketIntelligentTieringConfiguration =
      $grpc.ClientMethod<$0.DeleteBucketIntelligentTieringConfigurationRequest,
              $1.Empty>(
          '/s3.S3Service/DeleteBucketIntelligentTieringConfiguration',
          ($0.DeleteBucketIntelligentTieringConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketInventoryConfiguration = $grpc.ClientMethod<
          $0.DeleteBucketInventoryConfigurationRequest, $1.Empty>(
      '/s3.S3Service/DeleteBucketInventoryConfiguration',
      ($0.DeleteBucketInventoryConfigurationRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$deleteBucketLifecycle =
      $grpc.ClientMethod<$0.DeleteBucketLifecycleRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketLifecycle',
          ($0.DeleteBucketLifecycleRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketMetadataConfiguration =
      $grpc.ClientMethod<$0.DeleteBucketMetadataConfigurationRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketMetadataConfiguration',
          ($0.DeleteBucketMetadataConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketMetadataTableConfiguration = $grpc.ClientMethod<
          $0.DeleteBucketMetadataTableConfigurationRequest, $1.Empty>(
      '/s3.S3Service/DeleteBucketMetadataTableConfiguration',
      ($0.DeleteBucketMetadataTableConfigurationRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$deleteBucketMetricsConfiguration =
      $grpc.ClientMethod<$0.DeleteBucketMetricsConfigurationRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketMetricsConfiguration',
          ($0.DeleteBucketMetricsConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketOwnershipControls =
      $grpc.ClientMethod<$0.DeleteBucketOwnershipControlsRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketOwnershipControls',
          ($0.DeleteBucketOwnershipControlsRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketPolicy =
      $grpc.ClientMethod<$0.DeleteBucketPolicyRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketPolicy',
          ($0.DeleteBucketPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketReplication =
      $grpc.ClientMethod<$0.DeleteBucketReplicationRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketReplication',
          ($0.DeleteBucketReplicationRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketTagging =
      $grpc.ClientMethod<$0.DeleteBucketTaggingRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketTagging',
          ($0.DeleteBucketTaggingRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBucketWebsite =
      $grpc.ClientMethod<$0.DeleteBucketWebsiteRequest, $1.Empty>(
          '/s3.S3Service/DeleteBucketWebsite',
          ($0.DeleteBucketWebsiteRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteObject =
      $grpc.ClientMethod<$0.DeleteObjectRequest, $0.DeleteObjectOutput>(
          '/s3.S3Service/DeleteObject',
          ($0.DeleteObjectRequest value) => value.writeToBuffer(),
          $0.DeleteObjectOutput.fromBuffer);
  static final _$deleteObjects =
      $grpc.ClientMethod<$0.DeleteObjectsRequest, $0.DeleteObjectsOutput>(
          '/s3.S3Service/DeleteObjects',
          ($0.DeleteObjectsRequest value) => value.writeToBuffer(),
          $0.DeleteObjectsOutput.fromBuffer);
  static final _$deleteObjectTagging = $grpc.ClientMethod<
          $0.DeleteObjectTaggingRequest, $0.DeleteObjectTaggingOutput>(
      '/s3.S3Service/DeleteObjectTagging',
      ($0.DeleteObjectTaggingRequest value) => value.writeToBuffer(),
      $0.DeleteObjectTaggingOutput.fromBuffer);
  static final _$deletePublicAccessBlock =
      $grpc.ClientMethod<$0.DeletePublicAccessBlockRequest, $1.Empty>(
          '/s3.S3Service/DeletePublicAccessBlock',
          ($0.DeletePublicAccessBlockRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$getBucketAbac =
      $grpc.ClientMethod<$0.GetBucketAbacRequest, $0.GetBucketAbacOutput>(
          '/s3.S3Service/GetBucketAbac',
          ($0.GetBucketAbacRequest value) => value.writeToBuffer(),
          $0.GetBucketAbacOutput.fromBuffer);
  static final _$getBucketAccelerateConfiguration = $grpc.ClientMethod<
          $0.GetBucketAccelerateConfigurationRequest,
          $0.GetBucketAccelerateConfigurationOutput>(
      '/s3.S3Service/GetBucketAccelerateConfiguration',
      ($0.GetBucketAccelerateConfigurationRequest value) =>
          value.writeToBuffer(),
      $0.GetBucketAccelerateConfigurationOutput.fromBuffer);
  static final _$getBucketAcl =
      $grpc.ClientMethod<$0.GetBucketAclRequest, $0.GetBucketAclOutput>(
          '/s3.S3Service/GetBucketAcl',
          ($0.GetBucketAclRequest value) => value.writeToBuffer(),
          $0.GetBucketAclOutput.fromBuffer);
  static final _$getBucketAnalyticsConfiguration = $grpc.ClientMethod<
          $0.GetBucketAnalyticsConfigurationRequest,
          $0.GetBucketAnalyticsConfigurationOutput>(
      '/s3.S3Service/GetBucketAnalyticsConfiguration',
      ($0.GetBucketAnalyticsConfigurationRequest value) =>
          value.writeToBuffer(),
      $0.GetBucketAnalyticsConfigurationOutput.fromBuffer);
  static final _$getBucketCors =
      $grpc.ClientMethod<$0.GetBucketCorsRequest, $0.GetBucketCorsOutput>(
          '/s3.S3Service/GetBucketCors',
          ($0.GetBucketCorsRequest value) => value.writeToBuffer(),
          $0.GetBucketCorsOutput.fromBuffer);
  static final _$getBucketEncryption = $grpc.ClientMethod<
          $0.GetBucketEncryptionRequest, $0.GetBucketEncryptionOutput>(
      '/s3.S3Service/GetBucketEncryption',
      ($0.GetBucketEncryptionRequest value) => value.writeToBuffer(),
      $0.GetBucketEncryptionOutput.fromBuffer);
  static final _$getBucketIntelligentTieringConfiguration = $grpc.ClientMethod<
          $0.GetBucketIntelligentTieringConfigurationRequest,
          $0.GetBucketIntelligentTieringConfigurationOutput>(
      '/s3.S3Service/GetBucketIntelligentTieringConfiguration',
      ($0.GetBucketIntelligentTieringConfigurationRequest value) =>
          value.writeToBuffer(),
      $0.GetBucketIntelligentTieringConfigurationOutput.fromBuffer);
  static final _$getBucketInventoryConfiguration = $grpc.ClientMethod<
          $0.GetBucketInventoryConfigurationRequest,
          $0.GetBucketInventoryConfigurationOutput>(
      '/s3.S3Service/GetBucketInventoryConfiguration',
      ($0.GetBucketInventoryConfigurationRequest value) =>
          value.writeToBuffer(),
      $0.GetBucketInventoryConfigurationOutput.fromBuffer);
  static final _$getBucketLifecycleConfiguration = $grpc.ClientMethod<
          $0.GetBucketLifecycleConfigurationRequest,
          $0.GetBucketLifecycleConfigurationOutput>(
      '/s3.S3Service/GetBucketLifecycleConfiguration',
      ($0.GetBucketLifecycleConfigurationRequest value) =>
          value.writeToBuffer(),
      $0.GetBucketLifecycleConfigurationOutput.fromBuffer);
  static final _$getBucketLocation = $grpc.ClientMethod<
          $0.GetBucketLocationRequest, $0.GetBucketLocationOutput>(
      '/s3.S3Service/GetBucketLocation',
      ($0.GetBucketLocationRequest value) => value.writeToBuffer(),
      $0.GetBucketLocationOutput.fromBuffer);
  static final _$getBucketLogging =
      $grpc.ClientMethod<$0.GetBucketLoggingRequest, $0.GetBucketLoggingOutput>(
          '/s3.S3Service/GetBucketLogging',
          ($0.GetBucketLoggingRequest value) => value.writeToBuffer(),
          $0.GetBucketLoggingOutput.fromBuffer);
  static final _$getBucketMetadataConfiguration = $grpc.ClientMethod<
          $0.GetBucketMetadataConfigurationRequest,
          $0.GetBucketMetadataConfigurationOutput>(
      '/s3.S3Service/GetBucketMetadataConfiguration',
      ($0.GetBucketMetadataConfigurationRequest value) => value.writeToBuffer(),
      $0.GetBucketMetadataConfigurationOutput.fromBuffer);
  static final _$getBucketMetadataTableConfiguration = $grpc.ClientMethod<
          $0.GetBucketMetadataTableConfigurationRequest,
          $0.GetBucketMetadataTableConfigurationOutput>(
      '/s3.S3Service/GetBucketMetadataTableConfiguration',
      ($0.GetBucketMetadataTableConfigurationRequest value) =>
          value.writeToBuffer(),
      $0.GetBucketMetadataTableConfigurationOutput.fromBuffer);
  static final _$getBucketMetricsConfiguration = $grpc.ClientMethod<
          $0.GetBucketMetricsConfigurationRequest,
          $0.GetBucketMetricsConfigurationOutput>(
      '/s3.S3Service/GetBucketMetricsConfiguration',
      ($0.GetBucketMetricsConfigurationRequest value) => value.writeToBuffer(),
      $0.GetBucketMetricsConfigurationOutput.fromBuffer);
  static final _$getBucketNotificationConfiguration = $grpc.ClientMethod<
          $0.GetBucketNotificationConfigurationRequest,
          $0.NotificationConfiguration>(
      '/s3.S3Service/GetBucketNotificationConfiguration',
      ($0.GetBucketNotificationConfigurationRequest value) =>
          value.writeToBuffer(),
      $0.NotificationConfiguration.fromBuffer);
  static final _$getBucketOwnershipControls = $grpc.ClientMethod<
          $0.GetBucketOwnershipControlsRequest,
          $0.GetBucketOwnershipControlsOutput>(
      '/s3.S3Service/GetBucketOwnershipControls',
      ($0.GetBucketOwnershipControlsRequest value) => value.writeToBuffer(),
      $0.GetBucketOwnershipControlsOutput.fromBuffer);
  static final _$getBucketPolicy =
      $grpc.ClientMethod<$0.GetBucketPolicyRequest, $0.GetBucketPolicyOutput>(
          '/s3.S3Service/GetBucketPolicy',
          ($0.GetBucketPolicyRequest value) => value.writeToBuffer(),
          $0.GetBucketPolicyOutput.fromBuffer);
  static final _$getBucketPolicyStatus = $grpc.ClientMethod<
          $0.GetBucketPolicyStatusRequest, $0.GetBucketPolicyStatusOutput>(
      '/s3.S3Service/GetBucketPolicyStatus',
      ($0.GetBucketPolicyStatusRequest value) => value.writeToBuffer(),
      $0.GetBucketPolicyStatusOutput.fromBuffer);
  static final _$getBucketReplication = $grpc.ClientMethod<
          $0.GetBucketReplicationRequest, $0.GetBucketReplicationOutput>(
      '/s3.S3Service/GetBucketReplication',
      ($0.GetBucketReplicationRequest value) => value.writeToBuffer(),
      $0.GetBucketReplicationOutput.fromBuffer);
  static final _$getBucketRequestPayment = $grpc.ClientMethod<
          $0.GetBucketRequestPaymentRequest, $0.GetBucketRequestPaymentOutput>(
      '/s3.S3Service/GetBucketRequestPayment',
      ($0.GetBucketRequestPaymentRequest value) => value.writeToBuffer(),
      $0.GetBucketRequestPaymentOutput.fromBuffer);
  static final _$getBucketTagging =
      $grpc.ClientMethod<$0.GetBucketTaggingRequest, $0.GetBucketTaggingOutput>(
          '/s3.S3Service/GetBucketTagging',
          ($0.GetBucketTaggingRequest value) => value.writeToBuffer(),
          $0.GetBucketTaggingOutput.fromBuffer);
  static final _$getBucketVersioning = $grpc.ClientMethod<
          $0.GetBucketVersioningRequest, $0.GetBucketVersioningOutput>(
      '/s3.S3Service/GetBucketVersioning',
      ($0.GetBucketVersioningRequest value) => value.writeToBuffer(),
      $0.GetBucketVersioningOutput.fromBuffer);
  static final _$getBucketWebsite =
      $grpc.ClientMethod<$0.GetBucketWebsiteRequest, $0.GetBucketWebsiteOutput>(
          '/s3.S3Service/GetBucketWebsite',
          ($0.GetBucketWebsiteRequest value) => value.writeToBuffer(),
          $0.GetBucketWebsiteOutput.fromBuffer);
  static final _$getObject =
      $grpc.ClientMethod<$0.GetObjectRequest, $0.GetObjectOutput>(
          '/s3.S3Service/GetObject',
          ($0.GetObjectRequest value) => value.writeToBuffer(),
          $0.GetObjectOutput.fromBuffer);
  static final _$getObjectAcl =
      $grpc.ClientMethod<$0.GetObjectAclRequest, $0.GetObjectAclOutput>(
          '/s3.S3Service/GetObjectAcl',
          ($0.GetObjectAclRequest value) => value.writeToBuffer(),
          $0.GetObjectAclOutput.fromBuffer);
  static final _$getObjectAttributes = $grpc.ClientMethod<
          $0.GetObjectAttributesRequest, $0.GetObjectAttributesOutput>(
      '/s3.S3Service/GetObjectAttributes',
      ($0.GetObjectAttributesRequest value) => value.writeToBuffer(),
      $0.GetObjectAttributesOutput.fromBuffer);
  static final _$getObjectLegalHold = $grpc.ClientMethod<
          $0.GetObjectLegalHoldRequest, $0.GetObjectLegalHoldOutput>(
      '/s3.S3Service/GetObjectLegalHold',
      ($0.GetObjectLegalHoldRequest value) => value.writeToBuffer(),
      $0.GetObjectLegalHoldOutput.fromBuffer);
  static final _$getObjectLockConfiguration = $grpc.ClientMethod<
          $0.GetObjectLockConfigurationRequest,
          $0.GetObjectLockConfigurationOutput>(
      '/s3.S3Service/GetObjectLockConfiguration',
      ($0.GetObjectLockConfigurationRequest value) => value.writeToBuffer(),
      $0.GetObjectLockConfigurationOutput.fromBuffer);
  static final _$getObjectRetention = $grpc.ClientMethod<
          $0.GetObjectRetentionRequest, $0.GetObjectRetentionOutput>(
      '/s3.S3Service/GetObjectRetention',
      ($0.GetObjectRetentionRequest value) => value.writeToBuffer(),
      $0.GetObjectRetentionOutput.fromBuffer);
  static final _$getObjectTagging =
      $grpc.ClientMethod<$0.GetObjectTaggingRequest, $0.GetObjectTaggingOutput>(
          '/s3.S3Service/GetObjectTagging',
          ($0.GetObjectTaggingRequest value) => value.writeToBuffer(),
          $0.GetObjectTaggingOutput.fromBuffer);
  static final _$getObjectTorrent =
      $grpc.ClientMethod<$0.GetObjectTorrentRequest, $0.GetObjectTorrentOutput>(
          '/s3.S3Service/GetObjectTorrent',
          ($0.GetObjectTorrentRequest value) => value.writeToBuffer(),
          $0.GetObjectTorrentOutput.fromBuffer);
  static final _$getPublicAccessBlock = $grpc.ClientMethod<
          $0.GetPublicAccessBlockRequest, $0.GetPublicAccessBlockOutput>(
      '/s3.S3Service/GetPublicAccessBlock',
      ($0.GetPublicAccessBlockRequest value) => value.writeToBuffer(),
      $0.GetPublicAccessBlockOutput.fromBuffer);
  static final _$headBucket =
      $grpc.ClientMethod<$0.HeadBucketRequest, $0.HeadBucketOutput>(
          '/s3.S3Service/HeadBucket',
          ($0.HeadBucketRequest value) => value.writeToBuffer(),
          $0.HeadBucketOutput.fromBuffer);
  static final _$headObject =
      $grpc.ClientMethod<$0.HeadObjectRequest, $0.HeadObjectOutput>(
          '/s3.S3Service/HeadObject',
          ($0.HeadObjectRequest value) => value.writeToBuffer(),
          $0.HeadObjectOutput.fromBuffer);
  static final _$listBucketAnalyticsConfigurations = $grpc.ClientMethod<
          $0.ListBucketAnalyticsConfigurationsRequest,
          $0.ListBucketAnalyticsConfigurationsOutput>(
      '/s3.S3Service/ListBucketAnalyticsConfigurations',
      ($0.ListBucketAnalyticsConfigurationsRequest value) =>
          value.writeToBuffer(),
      $0.ListBucketAnalyticsConfigurationsOutput.fromBuffer);
  static final _$listBucketIntelligentTieringConfigurations =
      $grpc.ClientMethod<$0.ListBucketIntelligentTieringConfigurationsRequest,
              $0.ListBucketIntelligentTieringConfigurationsOutput>(
          '/s3.S3Service/ListBucketIntelligentTieringConfigurations',
          ($0.ListBucketIntelligentTieringConfigurationsRequest value) =>
              value.writeToBuffer(),
          $0.ListBucketIntelligentTieringConfigurationsOutput.fromBuffer);
  static final _$listBucketInventoryConfigurations = $grpc.ClientMethod<
          $0.ListBucketInventoryConfigurationsRequest,
          $0.ListBucketInventoryConfigurationsOutput>(
      '/s3.S3Service/ListBucketInventoryConfigurations',
      ($0.ListBucketInventoryConfigurationsRequest value) =>
          value.writeToBuffer(),
      $0.ListBucketInventoryConfigurationsOutput.fromBuffer);
  static final _$listBucketMetricsConfigurations = $grpc.ClientMethod<
          $0.ListBucketMetricsConfigurationsRequest,
          $0.ListBucketMetricsConfigurationsOutput>(
      '/s3.S3Service/ListBucketMetricsConfigurations',
      ($0.ListBucketMetricsConfigurationsRequest value) =>
          value.writeToBuffer(),
      $0.ListBucketMetricsConfigurationsOutput.fromBuffer);
  static final _$listBuckets =
      $grpc.ClientMethod<$0.ListBucketsRequest, $0.ListBucketsOutput>(
          '/s3.S3Service/ListBuckets',
          ($0.ListBucketsRequest value) => value.writeToBuffer(),
          $0.ListBucketsOutput.fromBuffer);
  static final _$listDirectoryBuckets = $grpc.ClientMethod<
          $0.ListDirectoryBucketsRequest, $0.ListDirectoryBucketsOutput>(
      '/s3.S3Service/ListDirectoryBuckets',
      ($0.ListDirectoryBucketsRequest value) => value.writeToBuffer(),
      $0.ListDirectoryBucketsOutput.fromBuffer);
  static final _$listMultipartUploads = $grpc.ClientMethod<
          $0.ListMultipartUploadsRequest, $0.ListMultipartUploadsOutput>(
      '/s3.S3Service/ListMultipartUploads',
      ($0.ListMultipartUploadsRequest value) => value.writeToBuffer(),
      $0.ListMultipartUploadsOutput.fromBuffer);
  static final _$listObjects =
      $grpc.ClientMethod<$0.ListObjectsRequest, $0.ListObjectsOutput>(
          '/s3.S3Service/ListObjects',
          ($0.ListObjectsRequest value) => value.writeToBuffer(),
          $0.ListObjectsOutput.fromBuffer);
  static final _$listObjectsV2 =
      $grpc.ClientMethod<$0.ListObjectsV2Request, $0.ListObjectsV2Output>(
          '/s3.S3Service/ListObjectsV2',
          ($0.ListObjectsV2Request value) => value.writeToBuffer(),
          $0.ListObjectsV2Output.fromBuffer);
  static final _$listObjectVersions = $grpc.ClientMethod<
          $0.ListObjectVersionsRequest, $0.ListObjectVersionsOutput>(
      '/s3.S3Service/ListObjectVersions',
      ($0.ListObjectVersionsRequest value) => value.writeToBuffer(),
      $0.ListObjectVersionsOutput.fromBuffer);
  static final _$listParts =
      $grpc.ClientMethod<$0.ListPartsRequest, $0.ListPartsOutput>(
          '/s3.S3Service/ListParts',
          ($0.ListPartsRequest value) => value.writeToBuffer(),
          $0.ListPartsOutput.fromBuffer);
  static final _$putBucketAbac =
      $grpc.ClientMethod<$0.PutBucketAbacRequest, $1.Empty>(
          '/s3.S3Service/PutBucketAbac',
          ($0.PutBucketAbacRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketAccelerateConfiguration =
      $grpc.ClientMethod<$0.PutBucketAccelerateConfigurationRequest, $1.Empty>(
          '/s3.S3Service/PutBucketAccelerateConfiguration',
          ($0.PutBucketAccelerateConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketAcl =
      $grpc.ClientMethod<$0.PutBucketAclRequest, $1.Empty>(
          '/s3.S3Service/PutBucketAcl',
          ($0.PutBucketAclRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketAnalyticsConfiguration =
      $grpc.ClientMethod<$0.PutBucketAnalyticsConfigurationRequest, $1.Empty>(
          '/s3.S3Service/PutBucketAnalyticsConfiguration',
          ($0.PutBucketAnalyticsConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketCors =
      $grpc.ClientMethod<$0.PutBucketCorsRequest, $1.Empty>(
          '/s3.S3Service/PutBucketCors',
          ($0.PutBucketCorsRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketEncryption =
      $grpc.ClientMethod<$0.PutBucketEncryptionRequest, $1.Empty>(
          '/s3.S3Service/PutBucketEncryption',
          ($0.PutBucketEncryptionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketIntelligentTieringConfiguration = $grpc.ClientMethod<
          $0.PutBucketIntelligentTieringConfigurationRequest, $1.Empty>(
      '/s3.S3Service/PutBucketIntelligentTieringConfiguration',
      ($0.PutBucketIntelligentTieringConfigurationRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$putBucketInventoryConfiguration =
      $grpc.ClientMethod<$0.PutBucketInventoryConfigurationRequest, $1.Empty>(
          '/s3.S3Service/PutBucketInventoryConfiguration',
          ($0.PutBucketInventoryConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketLifecycleConfiguration = $grpc.ClientMethod<
          $0.PutBucketLifecycleConfigurationRequest,
          $0.PutBucketLifecycleConfigurationOutput>(
      '/s3.S3Service/PutBucketLifecycleConfiguration',
      ($0.PutBucketLifecycleConfigurationRequest value) =>
          value.writeToBuffer(),
      $0.PutBucketLifecycleConfigurationOutput.fromBuffer);
  static final _$putBucketLogging =
      $grpc.ClientMethod<$0.PutBucketLoggingRequest, $1.Empty>(
          '/s3.S3Service/PutBucketLogging',
          ($0.PutBucketLoggingRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketMetricsConfiguration =
      $grpc.ClientMethod<$0.PutBucketMetricsConfigurationRequest, $1.Empty>(
          '/s3.S3Service/PutBucketMetricsConfiguration',
          ($0.PutBucketMetricsConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketNotificationConfiguration = $grpc.ClientMethod<
          $0.PutBucketNotificationConfigurationRequest, $1.Empty>(
      '/s3.S3Service/PutBucketNotificationConfiguration',
      ($0.PutBucketNotificationConfigurationRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$putBucketOwnershipControls =
      $grpc.ClientMethod<$0.PutBucketOwnershipControlsRequest, $1.Empty>(
          '/s3.S3Service/PutBucketOwnershipControls',
          ($0.PutBucketOwnershipControlsRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketPolicy =
      $grpc.ClientMethod<$0.PutBucketPolicyRequest, $1.Empty>(
          '/s3.S3Service/PutBucketPolicy',
          ($0.PutBucketPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketReplication =
      $grpc.ClientMethod<$0.PutBucketReplicationRequest, $1.Empty>(
          '/s3.S3Service/PutBucketReplication',
          ($0.PutBucketReplicationRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketRequestPayment =
      $grpc.ClientMethod<$0.PutBucketRequestPaymentRequest, $1.Empty>(
          '/s3.S3Service/PutBucketRequestPayment',
          ($0.PutBucketRequestPaymentRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketTagging =
      $grpc.ClientMethod<$0.PutBucketTaggingRequest, $1.Empty>(
          '/s3.S3Service/PutBucketTagging',
          ($0.PutBucketTaggingRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketVersioning =
      $grpc.ClientMethod<$0.PutBucketVersioningRequest, $1.Empty>(
          '/s3.S3Service/PutBucketVersioning',
          ($0.PutBucketVersioningRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putBucketWebsite =
      $grpc.ClientMethod<$0.PutBucketWebsiteRequest, $1.Empty>(
          '/s3.S3Service/PutBucketWebsite',
          ($0.PutBucketWebsiteRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putObject =
      $grpc.ClientMethod<$0.PutObjectRequest, $0.PutObjectOutput>(
          '/s3.S3Service/PutObject',
          ($0.PutObjectRequest value) => value.writeToBuffer(),
          $0.PutObjectOutput.fromBuffer);
  static final _$putObjectAcl =
      $grpc.ClientMethod<$0.PutObjectAclRequest, $0.PutObjectAclOutput>(
          '/s3.S3Service/PutObjectAcl',
          ($0.PutObjectAclRequest value) => value.writeToBuffer(),
          $0.PutObjectAclOutput.fromBuffer);
  static final _$putObjectLegalHold = $grpc.ClientMethod<
          $0.PutObjectLegalHoldRequest, $0.PutObjectLegalHoldOutput>(
      '/s3.S3Service/PutObjectLegalHold',
      ($0.PutObjectLegalHoldRequest value) => value.writeToBuffer(),
      $0.PutObjectLegalHoldOutput.fromBuffer);
  static final _$putObjectLockConfiguration = $grpc.ClientMethod<
          $0.PutObjectLockConfigurationRequest,
          $0.PutObjectLockConfigurationOutput>(
      '/s3.S3Service/PutObjectLockConfiguration',
      ($0.PutObjectLockConfigurationRequest value) => value.writeToBuffer(),
      $0.PutObjectLockConfigurationOutput.fromBuffer);
  static final _$putObjectRetention = $grpc.ClientMethod<
          $0.PutObjectRetentionRequest, $0.PutObjectRetentionOutput>(
      '/s3.S3Service/PutObjectRetention',
      ($0.PutObjectRetentionRequest value) => value.writeToBuffer(),
      $0.PutObjectRetentionOutput.fromBuffer);
  static final _$putObjectTagging =
      $grpc.ClientMethod<$0.PutObjectTaggingRequest, $0.PutObjectTaggingOutput>(
          '/s3.S3Service/PutObjectTagging',
          ($0.PutObjectTaggingRequest value) => value.writeToBuffer(),
          $0.PutObjectTaggingOutput.fromBuffer);
  static final _$putPublicAccessBlock =
      $grpc.ClientMethod<$0.PutPublicAccessBlockRequest, $1.Empty>(
          '/s3.S3Service/PutPublicAccessBlock',
          ($0.PutPublicAccessBlockRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$renameObject =
      $grpc.ClientMethod<$0.RenameObjectRequest, $0.RenameObjectOutput>(
          '/s3.S3Service/RenameObject',
          ($0.RenameObjectRequest value) => value.writeToBuffer(),
          $0.RenameObjectOutput.fromBuffer);
  static final _$restoreObject =
      $grpc.ClientMethod<$0.RestoreObjectRequest, $0.RestoreObjectOutput>(
          '/s3.S3Service/RestoreObject',
          ($0.RestoreObjectRequest value) => value.writeToBuffer(),
          $0.RestoreObjectOutput.fromBuffer);
  static final _$selectObjectContent = $grpc.ClientMethod<
          $0.SelectObjectContentRequest, $0.SelectObjectContentOutput>(
      '/s3.S3Service/SelectObjectContent',
      ($0.SelectObjectContentRequest value) => value.writeToBuffer(),
      $0.SelectObjectContentOutput.fromBuffer);
  static final _$updateBucketMetadataInventoryTableConfiguration =
      $grpc.ClientMethod<
              $0.UpdateBucketMetadataInventoryTableConfigurationRequest,
              $1.Empty>(
          '/s3.S3Service/UpdateBucketMetadataInventoryTableConfiguration',
          ($0.UpdateBucketMetadataInventoryTableConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateBucketMetadataJournalTableConfiguration =
      $grpc.ClientMethod<
              $0.UpdateBucketMetadataJournalTableConfigurationRequest,
              $1.Empty>(
          '/s3.S3Service/UpdateBucketMetadataJournalTableConfiguration',
          ($0.UpdateBucketMetadataJournalTableConfigurationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateObjectEncryption = $grpc.ClientMethod<
          $0.UpdateObjectEncryptionRequest, $0.UpdateObjectEncryptionResponse>(
      '/s3.S3Service/UpdateObjectEncryption',
      ($0.UpdateObjectEncryptionRequest value) => value.writeToBuffer(),
      $0.UpdateObjectEncryptionResponse.fromBuffer);
  static final _$uploadPart =
      $grpc.ClientMethod<$0.UploadPartRequest, $0.UploadPartOutput>(
          '/s3.S3Service/UploadPart',
          ($0.UploadPartRequest value) => value.writeToBuffer(),
          $0.UploadPartOutput.fromBuffer);
  static final _$uploadPartCopy =
      $grpc.ClientMethod<$0.UploadPartCopyRequest, $0.UploadPartCopyOutput>(
          '/s3.S3Service/UploadPartCopy',
          ($0.UploadPartCopyRequest value) => value.writeToBuffer(),
          $0.UploadPartCopyOutput.fromBuffer);
  static final _$writeGetObjectResponse =
      $grpc.ClientMethod<$0.WriteGetObjectResponseRequest, $1.Empty>(
          '/s3.S3Service/WriteGetObjectResponse',
          ($0.WriteGetObjectResponseRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
}

@$pb.GrpcServiceName('s3.S3Service')
abstract class S3ServiceBase extends $grpc.Service {
  $core.String get $name => 's3.S3Service';

  S3ServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AbortMultipartUploadRequest,
            $0.AbortMultipartUploadOutput>(
        'AbortMultipartUpload',
        abortMultipartUpload_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AbortMultipartUploadRequest.fromBuffer(value),
        ($0.AbortMultipartUploadOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CompleteMultipartUploadRequest,
            $0.CompleteMultipartUploadOutput>(
        'CompleteMultipartUpload',
        completeMultipartUpload_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CompleteMultipartUploadRequest.fromBuffer(value),
        ($0.CompleteMultipartUploadOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CopyObjectRequest, $0.CopyObjectOutput>(
        'CopyObject',
        copyObject_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CopyObjectRequest.fromBuffer(value),
        ($0.CopyObjectOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateBucketRequest, $0.CreateBucketOutput>(
            'CreateBucket',
            createBucket_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateBucketRequest.fromBuffer(value),
            ($0.CreateBucketOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateBucketMetadataConfigurationRequest,
            $1.Empty>(
        'CreateBucketMetadataConfiguration',
        createBucketMetadataConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateBucketMetadataConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.CreateBucketMetadataTableConfigurationRequest, $1.Empty>(
        'CreateBucketMetadataTableConfiguration',
        createBucketMetadataTableConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateBucketMetadataTableConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateMultipartUploadRequest,
            $0.CreateMultipartUploadOutput>(
        'CreateMultipartUpload',
        createMultipartUpload_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateMultipartUploadRequest.fromBuffer(value),
        ($0.CreateMultipartUploadOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateSessionRequest, $0.CreateSessionOutput>(
            'CreateSession',
            createSession_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateSessionRequest.fromBuffer(value),
            ($0.CreateSessionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketRequest, $1.Empty>(
        'DeleteBucket',
        deleteBucket_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketAnalyticsConfigurationRequest,
            $1.Empty>(
        'DeleteBucketAnalyticsConfiguration',
        deleteBucketAnalyticsConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketAnalyticsConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketCorsRequest, $1.Empty>(
        'DeleteBucketCors',
        deleteBucketCors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketCorsRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketEncryptionRequest, $1.Empty>(
        'DeleteBucketEncryption',
        deleteBucketEncryption_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketEncryptionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeleteBucketIntelligentTieringConfigurationRequest, $1.Empty>(
        'DeleteBucketIntelligentTieringConfiguration',
        deleteBucketIntelligentTieringConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketIntelligentTieringConfigurationRequest.fromBuffer(
                value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketInventoryConfigurationRequest,
            $1.Empty>(
        'DeleteBucketInventoryConfiguration',
        deleteBucketInventoryConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketInventoryConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketLifecycleRequest, $1.Empty>(
        'DeleteBucketLifecycle',
        deleteBucketLifecycle_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketLifecycleRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketMetadataConfigurationRequest,
            $1.Empty>(
        'DeleteBucketMetadataConfiguration',
        deleteBucketMetadataConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketMetadataConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeleteBucketMetadataTableConfigurationRequest, $1.Empty>(
        'DeleteBucketMetadataTableConfiguration',
        deleteBucketMetadataTableConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketMetadataTableConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketMetricsConfigurationRequest,
            $1.Empty>(
        'DeleteBucketMetricsConfiguration',
        deleteBucketMetricsConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketMetricsConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteBucketOwnershipControlsRequest, $1.Empty>(
            'DeleteBucketOwnershipControls',
            deleteBucketOwnershipControls_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteBucketOwnershipControlsRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketPolicyRequest, $1.Empty>(
        'DeleteBucketPolicy',
        deleteBucketPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketReplicationRequest, $1.Empty>(
        'DeleteBucketReplication',
        deleteBucketReplication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketReplicationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketTaggingRequest, $1.Empty>(
        'DeleteBucketTagging',
        deleteBucketTagging_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketTaggingRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBucketWebsiteRequest, $1.Empty>(
        'DeleteBucketWebsite',
        deleteBucketWebsite_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBucketWebsiteRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteObjectRequest, $0.DeleteObjectOutput>(
            'DeleteObject',
            deleteObject_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteObjectRequest.fromBuffer(value),
            ($0.DeleteObjectOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteObjectsRequest, $0.DeleteObjectsOutput>(
            'DeleteObjects',
            deleteObjects_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteObjectsRequest.fromBuffer(value),
            ($0.DeleteObjectsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteObjectTaggingRequest,
            $0.DeleteObjectTaggingOutput>(
        'DeleteObjectTagging',
        deleteObjectTagging_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteObjectTaggingRequest.fromBuffer(value),
        ($0.DeleteObjectTaggingOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePublicAccessBlockRequest, $1.Empty>(
        'DeletePublicAccessBlock',
        deletePublicAccessBlock_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeletePublicAccessBlockRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetBucketAbacRequest, $0.GetBucketAbacOutput>(
            'GetBucketAbac',
            getBucketAbac_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetBucketAbacRequest.fromBuffer(value),
            ($0.GetBucketAbacOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketAccelerateConfigurationRequest,
            $0.GetBucketAccelerateConfigurationOutput>(
        'GetBucketAccelerateConfiguration',
        getBucketAccelerateConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketAccelerateConfigurationRequest.fromBuffer(value),
        ($0.GetBucketAccelerateConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetBucketAclRequest, $0.GetBucketAclOutput>(
            'GetBucketAcl',
            getBucketAcl_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetBucketAclRequest.fromBuffer(value),
            ($0.GetBucketAclOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketAnalyticsConfigurationRequest,
            $0.GetBucketAnalyticsConfigurationOutput>(
        'GetBucketAnalyticsConfiguration',
        getBucketAnalyticsConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketAnalyticsConfigurationRequest.fromBuffer(value),
        ($0.GetBucketAnalyticsConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetBucketCorsRequest, $0.GetBucketCorsOutput>(
            'GetBucketCors',
            getBucketCors_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetBucketCorsRequest.fromBuffer(value),
            ($0.GetBucketCorsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketEncryptionRequest,
            $0.GetBucketEncryptionOutput>(
        'GetBucketEncryption',
        getBucketEncryption_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketEncryptionRequest.fromBuffer(value),
        ($0.GetBucketEncryptionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetBucketIntelligentTieringConfigurationRequest,
            $0.GetBucketIntelligentTieringConfigurationOutput>(
        'GetBucketIntelligentTieringConfiguration',
        getBucketIntelligentTieringConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketIntelligentTieringConfigurationRequest.fromBuffer(
                value),
        ($0.GetBucketIntelligentTieringConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketInventoryConfigurationRequest,
            $0.GetBucketInventoryConfigurationOutput>(
        'GetBucketInventoryConfiguration',
        getBucketInventoryConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketInventoryConfigurationRequest.fromBuffer(value),
        ($0.GetBucketInventoryConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketLifecycleConfigurationRequest,
            $0.GetBucketLifecycleConfigurationOutput>(
        'GetBucketLifecycleConfiguration',
        getBucketLifecycleConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketLifecycleConfigurationRequest.fromBuffer(value),
        ($0.GetBucketLifecycleConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketLocationRequest,
            $0.GetBucketLocationOutput>(
        'GetBucketLocation',
        getBucketLocation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketLocationRequest.fromBuffer(value),
        ($0.GetBucketLocationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketLoggingRequest,
            $0.GetBucketLoggingOutput>(
        'GetBucketLogging',
        getBucketLogging_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketLoggingRequest.fromBuffer(value),
        ($0.GetBucketLoggingOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketMetadataConfigurationRequest,
            $0.GetBucketMetadataConfigurationOutput>(
        'GetBucketMetadataConfiguration',
        getBucketMetadataConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketMetadataConfigurationRequest.fromBuffer(value),
        ($0.GetBucketMetadataConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetBucketMetadataTableConfigurationRequest,
            $0.GetBucketMetadataTableConfigurationOutput>(
        'GetBucketMetadataTableConfiguration',
        getBucketMetadataTableConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketMetadataTableConfigurationRequest.fromBuffer(value),
        ($0.GetBucketMetadataTableConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketMetricsConfigurationRequest,
            $0.GetBucketMetricsConfigurationOutput>(
        'GetBucketMetricsConfiguration',
        getBucketMetricsConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketMetricsConfigurationRequest.fromBuffer(value),
        ($0.GetBucketMetricsConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketNotificationConfigurationRequest,
            $0.NotificationConfiguration>(
        'GetBucketNotificationConfiguration',
        getBucketNotificationConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketNotificationConfigurationRequest.fromBuffer(value),
        ($0.NotificationConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketOwnershipControlsRequest,
            $0.GetBucketOwnershipControlsOutput>(
        'GetBucketOwnershipControls',
        getBucketOwnershipControls_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketOwnershipControlsRequest.fromBuffer(value),
        ($0.GetBucketOwnershipControlsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketPolicyRequest,
            $0.GetBucketPolicyOutput>(
        'GetBucketPolicy',
        getBucketPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketPolicyRequest.fromBuffer(value),
        ($0.GetBucketPolicyOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketPolicyStatusRequest,
            $0.GetBucketPolicyStatusOutput>(
        'GetBucketPolicyStatus',
        getBucketPolicyStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketPolicyStatusRequest.fromBuffer(value),
        ($0.GetBucketPolicyStatusOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketReplicationRequest,
            $0.GetBucketReplicationOutput>(
        'GetBucketReplication',
        getBucketReplication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketReplicationRequest.fromBuffer(value),
        ($0.GetBucketReplicationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketRequestPaymentRequest,
            $0.GetBucketRequestPaymentOutput>(
        'GetBucketRequestPayment',
        getBucketRequestPayment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketRequestPaymentRequest.fromBuffer(value),
        ($0.GetBucketRequestPaymentOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketTaggingRequest,
            $0.GetBucketTaggingOutput>(
        'GetBucketTagging',
        getBucketTagging_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketTaggingRequest.fromBuffer(value),
        ($0.GetBucketTaggingOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketVersioningRequest,
            $0.GetBucketVersioningOutput>(
        'GetBucketVersioning',
        getBucketVersioning_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketVersioningRequest.fromBuffer(value),
        ($0.GetBucketVersioningOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBucketWebsiteRequest,
            $0.GetBucketWebsiteOutput>(
        'GetBucketWebsite',
        getBucketWebsite_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBucketWebsiteRequest.fromBuffer(value),
        ($0.GetBucketWebsiteOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetObjectRequest, $0.GetObjectOutput>(
        'GetObject',
        getObject_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetObjectRequest.fromBuffer(value),
        ($0.GetObjectOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetObjectAclRequest, $0.GetObjectAclOutput>(
            'GetObjectAcl',
            getObjectAcl_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetObjectAclRequest.fromBuffer(value),
            ($0.GetObjectAclOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetObjectAttributesRequest,
            $0.GetObjectAttributesOutput>(
        'GetObjectAttributes',
        getObjectAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetObjectAttributesRequest.fromBuffer(value),
        ($0.GetObjectAttributesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetObjectLegalHoldRequest,
            $0.GetObjectLegalHoldOutput>(
        'GetObjectLegalHold',
        getObjectLegalHold_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetObjectLegalHoldRequest.fromBuffer(value),
        ($0.GetObjectLegalHoldOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetObjectLockConfigurationRequest,
            $0.GetObjectLockConfigurationOutput>(
        'GetObjectLockConfiguration',
        getObjectLockConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetObjectLockConfigurationRequest.fromBuffer(value),
        ($0.GetObjectLockConfigurationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetObjectRetentionRequest,
            $0.GetObjectRetentionOutput>(
        'GetObjectRetention',
        getObjectRetention_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetObjectRetentionRequest.fromBuffer(value),
        ($0.GetObjectRetentionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetObjectTaggingRequest,
            $0.GetObjectTaggingOutput>(
        'GetObjectTagging',
        getObjectTagging_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetObjectTaggingRequest.fromBuffer(value),
        ($0.GetObjectTaggingOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetObjectTorrentRequest,
            $0.GetObjectTorrentOutput>(
        'GetObjectTorrent',
        getObjectTorrent_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetObjectTorrentRequest.fromBuffer(value),
        ($0.GetObjectTorrentOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPublicAccessBlockRequest,
            $0.GetPublicAccessBlockOutput>(
        'GetPublicAccessBlock',
        getPublicAccessBlock_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPublicAccessBlockRequest.fromBuffer(value),
        ($0.GetPublicAccessBlockOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.HeadBucketRequest, $0.HeadBucketOutput>(
        'HeadBucket',
        headBucket_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.HeadBucketRequest.fromBuffer(value),
        ($0.HeadBucketOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.HeadObjectRequest, $0.HeadObjectOutput>(
        'HeadObject',
        headObject_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.HeadObjectRequest.fromBuffer(value),
        ($0.HeadObjectOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListBucketAnalyticsConfigurationsRequest,
            $0.ListBucketAnalyticsConfigurationsOutput>(
        'ListBucketAnalyticsConfigurations',
        listBucketAnalyticsConfigurations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListBucketAnalyticsConfigurationsRequest.fromBuffer(value),
        ($0.ListBucketAnalyticsConfigurationsOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListBucketIntelligentTieringConfigurationsRequest,
            $0.ListBucketIntelligentTieringConfigurationsOutput>(
        'ListBucketIntelligentTieringConfigurations',
        listBucketIntelligentTieringConfigurations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListBucketIntelligentTieringConfigurationsRequest.fromBuffer(
                value),
        ($0.ListBucketIntelligentTieringConfigurationsOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListBucketInventoryConfigurationsRequest,
            $0.ListBucketInventoryConfigurationsOutput>(
        'ListBucketInventoryConfigurations',
        listBucketInventoryConfigurations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListBucketInventoryConfigurationsRequest.fromBuffer(value),
        ($0.ListBucketInventoryConfigurationsOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListBucketMetricsConfigurationsRequest,
            $0.ListBucketMetricsConfigurationsOutput>(
        'ListBucketMetricsConfigurations',
        listBucketMetricsConfigurations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListBucketMetricsConfigurationsRequest.fromBuffer(value),
        ($0.ListBucketMetricsConfigurationsOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListBucketsRequest, $0.ListBucketsOutput>(
        'ListBuckets',
        listBuckets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListBucketsRequest.fromBuffer(value),
        ($0.ListBucketsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDirectoryBucketsRequest,
            $0.ListDirectoryBucketsOutput>(
        'ListDirectoryBuckets',
        listDirectoryBuckets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDirectoryBucketsRequest.fromBuffer(value),
        ($0.ListDirectoryBucketsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMultipartUploadsRequest,
            $0.ListMultipartUploadsOutput>(
        'ListMultipartUploads',
        listMultipartUploads_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListMultipartUploadsRequest.fromBuffer(value),
        ($0.ListMultipartUploadsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListObjectsRequest, $0.ListObjectsOutput>(
        'ListObjects',
        listObjects_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListObjectsRequest.fromBuffer(value),
        ($0.ListObjectsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListObjectsV2Request, $0.ListObjectsV2Output>(
            'ListObjectsV2',
            listObjectsV2_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListObjectsV2Request.fromBuffer(value),
            ($0.ListObjectsV2Output value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListObjectVersionsRequest,
            $0.ListObjectVersionsOutput>(
        'ListObjectVersions',
        listObjectVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListObjectVersionsRequest.fromBuffer(value),
        ($0.ListObjectVersionsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPartsRequest, $0.ListPartsOutput>(
        'ListParts',
        listParts_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListPartsRequest.fromBuffer(value),
        ($0.ListPartsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketAbacRequest, $1.Empty>(
        'PutBucketAbac',
        putBucketAbac_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketAbacRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketAccelerateConfigurationRequest,
            $1.Empty>(
        'PutBucketAccelerateConfiguration',
        putBucketAccelerateConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketAccelerateConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketAclRequest, $1.Empty>(
        'PutBucketAcl',
        putBucketAcl_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketAclRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketAnalyticsConfigurationRequest,
            $1.Empty>(
        'PutBucketAnalyticsConfiguration',
        putBucketAnalyticsConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketAnalyticsConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketCorsRequest, $1.Empty>(
        'PutBucketCors',
        putBucketCors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketCorsRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketEncryptionRequest, $1.Empty>(
        'PutBucketEncryption',
        putBucketEncryption_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketEncryptionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.PutBucketIntelligentTieringConfigurationRequest, $1.Empty>(
        'PutBucketIntelligentTieringConfiguration',
        putBucketIntelligentTieringConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketIntelligentTieringConfigurationRequest.fromBuffer(
                value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketInventoryConfigurationRequest,
            $1.Empty>(
        'PutBucketInventoryConfiguration',
        putBucketInventoryConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketInventoryConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketLifecycleConfigurationRequest,
            $0.PutBucketLifecycleConfigurationOutput>(
        'PutBucketLifecycleConfiguration',
        putBucketLifecycleConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketLifecycleConfigurationRequest.fromBuffer(value),
        ($0.PutBucketLifecycleConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketLoggingRequest, $1.Empty>(
        'PutBucketLogging',
        putBucketLogging_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketLoggingRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutBucketMetricsConfigurationRequest, $1.Empty>(
            'PutBucketMetricsConfiguration',
            putBucketMetricsConfiguration_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutBucketMetricsConfigurationRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketNotificationConfigurationRequest,
            $1.Empty>(
        'PutBucketNotificationConfiguration',
        putBucketNotificationConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketNotificationConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutBucketOwnershipControlsRequest, $1.Empty>(
            'PutBucketOwnershipControls',
            putBucketOwnershipControls_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutBucketOwnershipControlsRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketPolicyRequest, $1.Empty>(
        'PutBucketPolicy',
        putBucketPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketReplicationRequest, $1.Empty>(
        'PutBucketReplication',
        putBucketReplication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketReplicationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketRequestPaymentRequest, $1.Empty>(
        'PutBucketRequestPayment',
        putBucketRequestPayment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketRequestPaymentRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketTaggingRequest, $1.Empty>(
        'PutBucketTagging',
        putBucketTagging_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketTaggingRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketVersioningRequest, $1.Empty>(
        'PutBucketVersioning',
        putBucketVersioning_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketVersioningRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutBucketWebsiteRequest, $1.Empty>(
        'PutBucketWebsite',
        putBucketWebsite_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutBucketWebsiteRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutObjectRequest, $0.PutObjectOutput>(
        'PutObject',
        putObject_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutObjectRequest.fromBuffer(value),
        ($0.PutObjectOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutObjectAclRequest, $0.PutObjectAclOutput>(
            'PutObjectAcl',
            putObjectAcl_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutObjectAclRequest.fromBuffer(value),
            ($0.PutObjectAclOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutObjectLegalHoldRequest,
            $0.PutObjectLegalHoldOutput>(
        'PutObjectLegalHold',
        putObjectLegalHold_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutObjectLegalHoldRequest.fromBuffer(value),
        ($0.PutObjectLegalHoldOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutObjectLockConfigurationRequest,
            $0.PutObjectLockConfigurationOutput>(
        'PutObjectLockConfiguration',
        putObjectLockConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutObjectLockConfigurationRequest.fromBuffer(value),
        ($0.PutObjectLockConfigurationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutObjectRetentionRequest,
            $0.PutObjectRetentionOutput>(
        'PutObjectRetention',
        putObjectRetention_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutObjectRetentionRequest.fromBuffer(value),
        ($0.PutObjectRetentionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutObjectTaggingRequest,
            $0.PutObjectTaggingOutput>(
        'PutObjectTagging',
        putObjectTagging_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutObjectTaggingRequest.fromBuffer(value),
        ($0.PutObjectTaggingOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutPublicAccessBlockRequest, $1.Empty>(
        'PutPublicAccessBlock',
        putPublicAccessBlock_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutPublicAccessBlockRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RenameObjectRequest, $0.RenameObjectOutput>(
            'RenameObject',
            renameObject_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RenameObjectRequest.fromBuffer(value),
            ($0.RenameObjectOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RestoreObjectRequest, $0.RestoreObjectOutput>(
            'RestoreObject',
            restoreObject_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RestoreObjectRequest.fromBuffer(value),
            ($0.RestoreObjectOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SelectObjectContentRequest,
            $0.SelectObjectContentOutput>(
        'SelectObjectContent',
        selectObjectContent_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SelectObjectContentRequest.fromBuffer(value),
        ($0.SelectObjectContentOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateBucketMetadataInventoryTableConfigurationRequest,
            $1.Empty>(
        'UpdateBucketMetadataInventoryTableConfiguration',
        updateBucketMetadataInventoryTableConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateBucketMetadataInventoryTableConfigurationRequest
                .fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateBucketMetadataJournalTableConfigurationRequest, $1.Empty>(
        'UpdateBucketMetadataJournalTableConfiguration',
        updateBucketMetadataJournalTableConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateBucketMetadataJournalTableConfigurationRequest.fromBuffer(
                value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateObjectEncryptionRequest,
            $0.UpdateObjectEncryptionResponse>(
        'UpdateObjectEncryption',
        updateObjectEncryption_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateObjectEncryptionRequest.fromBuffer(value),
        ($0.UpdateObjectEncryptionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UploadPartRequest, $0.UploadPartOutput>(
        'UploadPart',
        uploadPart_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UploadPartRequest.fromBuffer(value),
        ($0.UploadPartOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UploadPartCopyRequest, $0.UploadPartCopyOutput>(
            'UploadPartCopy',
            uploadPartCopy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UploadPartCopyRequest.fromBuffer(value),
            ($0.UploadPartCopyOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.WriteGetObjectResponseRequest, $1.Empty>(
        'WriteGetObjectResponse',
        writeGetObjectResponse_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.WriteGetObjectResponseRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
  }

  $async.Future<$0.AbortMultipartUploadOutput> abortMultipartUpload_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AbortMultipartUploadRequest> $request) async {
    return abortMultipartUpload($call, await $request);
  }

  $async.Future<$0.AbortMultipartUploadOutput> abortMultipartUpload(
      $grpc.ServiceCall call, $0.AbortMultipartUploadRequest request);

  $async.Future<$0.CompleteMultipartUploadOutput> completeMultipartUpload_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CompleteMultipartUploadRequest> $request) async {
    return completeMultipartUpload($call, await $request);
  }

  $async.Future<$0.CompleteMultipartUploadOutput> completeMultipartUpload(
      $grpc.ServiceCall call, $0.CompleteMultipartUploadRequest request);

  $async.Future<$0.CopyObjectOutput> copyObject_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CopyObjectRequest> $request) async {
    return copyObject($call, await $request);
  }

  $async.Future<$0.CopyObjectOutput> copyObject(
      $grpc.ServiceCall call, $0.CopyObjectRequest request);

  $async.Future<$0.CreateBucketOutput> createBucket_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateBucketRequest> $request) async {
    return createBucket($call, await $request);
  }

  $async.Future<$0.CreateBucketOutput> createBucket(
      $grpc.ServiceCall call, $0.CreateBucketRequest request);

  $async.Future<$1.Empty> createBucketMetadataConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateBucketMetadataConfigurationRequest>
          $request) async {
    return createBucketMetadataConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> createBucketMetadataConfiguration(
      $grpc.ServiceCall call,
      $0.CreateBucketMetadataConfigurationRequest request);

  $async.Future<$1.Empty> createBucketMetadataTableConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateBucketMetadataTableConfigurationRequest>
          $request) async {
    return createBucketMetadataTableConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> createBucketMetadataTableConfiguration(
      $grpc.ServiceCall call,
      $0.CreateBucketMetadataTableConfigurationRequest request);

  $async.Future<$0.CreateMultipartUploadOutput> createMultipartUpload_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateMultipartUploadRequest> $request) async {
    return createMultipartUpload($call, await $request);
  }

  $async.Future<$0.CreateMultipartUploadOutput> createMultipartUpload(
      $grpc.ServiceCall call, $0.CreateMultipartUploadRequest request);

  $async.Future<$0.CreateSessionOutput> createSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateSessionRequest> $request) async {
    return createSession($call, await $request);
  }

  $async.Future<$0.CreateSessionOutput> createSession(
      $grpc.ServiceCall call, $0.CreateSessionRequest request);

  $async.Future<$1.Empty> deleteBucket_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketRequest> $request) async {
    return deleteBucket($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucket(
      $grpc.ServiceCall call, $0.DeleteBucketRequest request);

  $async.Future<$1.Empty> deleteBucketAnalyticsConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketAnalyticsConfigurationRequest>
          $request) async {
    return deleteBucketAnalyticsConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketAnalyticsConfiguration(
      $grpc.ServiceCall call,
      $0.DeleteBucketAnalyticsConfigurationRequest request);

  $async.Future<$1.Empty> deleteBucketCors_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketCorsRequest> $request) async {
    return deleteBucketCors($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketCors(
      $grpc.ServiceCall call, $0.DeleteBucketCorsRequest request);

  $async.Future<$1.Empty> deleteBucketEncryption_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketEncryptionRequest> $request) async {
    return deleteBucketEncryption($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketEncryption(
      $grpc.ServiceCall call, $0.DeleteBucketEncryptionRequest request);

  $async.Future<$1.Empty> deleteBucketIntelligentTieringConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketIntelligentTieringConfigurationRequest>
          $request) async {
    return deleteBucketIntelligentTieringConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketIntelligentTieringConfiguration(
      $grpc.ServiceCall call,
      $0.DeleteBucketIntelligentTieringConfigurationRequest request);

  $async.Future<$1.Empty> deleteBucketInventoryConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketInventoryConfigurationRequest>
          $request) async {
    return deleteBucketInventoryConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketInventoryConfiguration(
      $grpc.ServiceCall call,
      $0.DeleteBucketInventoryConfigurationRequest request);

  $async.Future<$1.Empty> deleteBucketLifecycle_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketLifecycleRequest> $request) async {
    return deleteBucketLifecycle($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketLifecycle(
      $grpc.ServiceCall call, $0.DeleteBucketLifecycleRequest request);

  $async.Future<$1.Empty> deleteBucketMetadataConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketMetadataConfigurationRequest>
          $request) async {
    return deleteBucketMetadataConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketMetadataConfiguration(
      $grpc.ServiceCall call,
      $0.DeleteBucketMetadataConfigurationRequest request);

  $async.Future<$1.Empty> deleteBucketMetadataTableConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketMetadataTableConfigurationRequest>
          $request) async {
    return deleteBucketMetadataTableConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketMetadataTableConfiguration(
      $grpc.ServiceCall call,
      $0.DeleteBucketMetadataTableConfigurationRequest request);

  $async.Future<$1.Empty> deleteBucketMetricsConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketMetricsConfigurationRequest>
          $request) async {
    return deleteBucketMetricsConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketMetricsConfiguration(
      $grpc.ServiceCall call,
      $0.DeleteBucketMetricsConfigurationRequest request);

  $async.Future<$1.Empty> deleteBucketOwnershipControls_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketOwnershipControlsRequest> $request) async {
    return deleteBucketOwnershipControls($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketOwnershipControls(
      $grpc.ServiceCall call, $0.DeleteBucketOwnershipControlsRequest request);

  $async.Future<$1.Empty> deleteBucketPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketPolicyRequest> $request) async {
    return deleteBucketPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketPolicy(
      $grpc.ServiceCall call, $0.DeleteBucketPolicyRequest request);

  $async.Future<$1.Empty> deleteBucketReplication_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketReplicationRequest> $request) async {
    return deleteBucketReplication($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketReplication(
      $grpc.ServiceCall call, $0.DeleteBucketReplicationRequest request);

  $async.Future<$1.Empty> deleteBucketTagging_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketTaggingRequest> $request) async {
    return deleteBucketTagging($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketTagging(
      $grpc.ServiceCall call, $0.DeleteBucketTaggingRequest request);

  $async.Future<$1.Empty> deleteBucketWebsite_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBucketWebsiteRequest> $request) async {
    return deleteBucketWebsite($call, await $request);
  }

  $async.Future<$1.Empty> deleteBucketWebsite(
      $grpc.ServiceCall call, $0.DeleteBucketWebsiteRequest request);

  $async.Future<$0.DeleteObjectOutput> deleteObject_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteObjectRequest> $request) async {
    return deleteObject($call, await $request);
  }

  $async.Future<$0.DeleteObjectOutput> deleteObject(
      $grpc.ServiceCall call, $0.DeleteObjectRequest request);

  $async.Future<$0.DeleteObjectsOutput> deleteObjects_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteObjectsRequest> $request) async {
    return deleteObjects($call, await $request);
  }

  $async.Future<$0.DeleteObjectsOutput> deleteObjects(
      $grpc.ServiceCall call, $0.DeleteObjectsRequest request);

  $async.Future<$0.DeleteObjectTaggingOutput> deleteObjectTagging_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteObjectTaggingRequest> $request) async {
    return deleteObjectTagging($call, await $request);
  }

  $async.Future<$0.DeleteObjectTaggingOutput> deleteObjectTagging(
      $grpc.ServiceCall call, $0.DeleteObjectTaggingRequest request);

  $async.Future<$1.Empty> deletePublicAccessBlock_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeletePublicAccessBlockRequest> $request) async {
    return deletePublicAccessBlock($call, await $request);
  }

  $async.Future<$1.Empty> deletePublicAccessBlock(
      $grpc.ServiceCall call, $0.DeletePublicAccessBlockRequest request);

  $async.Future<$0.GetBucketAbacOutput> getBucketAbac_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketAbacRequest> $request) async {
    return getBucketAbac($call, await $request);
  }

  $async.Future<$0.GetBucketAbacOutput> getBucketAbac(
      $grpc.ServiceCall call, $0.GetBucketAbacRequest request);

  $async.Future<$0.GetBucketAccelerateConfigurationOutput>
      getBucketAccelerateConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketAccelerateConfigurationRequest>
              $request) async {
    return getBucketAccelerateConfiguration($call, await $request);
  }

  $async.Future<$0.GetBucketAccelerateConfigurationOutput>
      getBucketAccelerateConfiguration($grpc.ServiceCall call,
          $0.GetBucketAccelerateConfigurationRequest request);

  $async.Future<$0.GetBucketAclOutput> getBucketAcl_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetBucketAclRequest> $request) async {
    return getBucketAcl($call, await $request);
  }

  $async.Future<$0.GetBucketAclOutput> getBucketAcl(
      $grpc.ServiceCall call, $0.GetBucketAclRequest request);

  $async.Future<$0.GetBucketAnalyticsConfigurationOutput>
      getBucketAnalyticsConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketAnalyticsConfigurationRequest>
              $request) async {
    return getBucketAnalyticsConfiguration($call, await $request);
  }

  $async.Future<$0.GetBucketAnalyticsConfigurationOutput>
      getBucketAnalyticsConfiguration($grpc.ServiceCall call,
          $0.GetBucketAnalyticsConfigurationRequest request);

  $async.Future<$0.GetBucketCorsOutput> getBucketCors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketCorsRequest> $request) async {
    return getBucketCors($call, await $request);
  }

  $async.Future<$0.GetBucketCorsOutput> getBucketCors(
      $grpc.ServiceCall call, $0.GetBucketCorsRequest request);

  $async.Future<$0.GetBucketEncryptionOutput> getBucketEncryption_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketEncryptionRequest> $request) async {
    return getBucketEncryption($call, await $request);
  }

  $async.Future<$0.GetBucketEncryptionOutput> getBucketEncryption(
      $grpc.ServiceCall call, $0.GetBucketEncryptionRequest request);

  $async.Future<$0.GetBucketIntelligentTieringConfigurationOutput>
      getBucketIntelligentTieringConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketIntelligentTieringConfigurationRequest>
              $request) async {
    return getBucketIntelligentTieringConfiguration($call, await $request);
  }

  $async.Future<$0.GetBucketIntelligentTieringConfigurationOutput>
      getBucketIntelligentTieringConfiguration($grpc.ServiceCall call,
          $0.GetBucketIntelligentTieringConfigurationRequest request);

  $async.Future<$0.GetBucketInventoryConfigurationOutput>
      getBucketInventoryConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketInventoryConfigurationRequest>
              $request) async {
    return getBucketInventoryConfiguration($call, await $request);
  }

  $async.Future<$0.GetBucketInventoryConfigurationOutput>
      getBucketInventoryConfiguration($grpc.ServiceCall call,
          $0.GetBucketInventoryConfigurationRequest request);

  $async.Future<$0.GetBucketLifecycleConfigurationOutput>
      getBucketLifecycleConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketLifecycleConfigurationRequest>
              $request) async {
    return getBucketLifecycleConfiguration($call, await $request);
  }

  $async.Future<$0.GetBucketLifecycleConfigurationOutput>
      getBucketLifecycleConfiguration($grpc.ServiceCall call,
          $0.GetBucketLifecycleConfigurationRequest request);

  $async.Future<$0.GetBucketLocationOutput> getBucketLocation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketLocationRequest> $request) async {
    return getBucketLocation($call, await $request);
  }

  $async.Future<$0.GetBucketLocationOutput> getBucketLocation(
      $grpc.ServiceCall call, $0.GetBucketLocationRequest request);

  $async.Future<$0.GetBucketLoggingOutput> getBucketLogging_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketLoggingRequest> $request) async {
    return getBucketLogging($call, await $request);
  }

  $async.Future<$0.GetBucketLoggingOutput> getBucketLogging(
      $grpc.ServiceCall call, $0.GetBucketLoggingRequest request);

  $async.Future<$0.GetBucketMetadataConfigurationOutput>
      getBucketMetadataConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketMetadataConfigurationRequest>
              $request) async {
    return getBucketMetadataConfiguration($call, await $request);
  }

  $async.Future<$0.GetBucketMetadataConfigurationOutput>
      getBucketMetadataConfiguration($grpc.ServiceCall call,
          $0.GetBucketMetadataConfigurationRequest request);

  $async.Future<$0.GetBucketMetadataTableConfigurationOutput>
      getBucketMetadataTableConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketMetadataTableConfigurationRequest>
              $request) async {
    return getBucketMetadataTableConfiguration($call, await $request);
  }

  $async.Future<$0.GetBucketMetadataTableConfigurationOutput>
      getBucketMetadataTableConfiguration($grpc.ServiceCall call,
          $0.GetBucketMetadataTableConfigurationRequest request);

  $async.Future<$0.GetBucketMetricsConfigurationOutput>
      getBucketMetricsConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketMetricsConfigurationRequest>
              $request) async {
    return getBucketMetricsConfiguration($call, await $request);
  }

  $async.Future<$0.GetBucketMetricsConfigurationOutput>
      getBucketMetricsConfiguration($grpc.ServiceCall call,
          $0.GetBucketMetricsConfigurationRequest request);

  $async.Future<$0.NotificationConfiguration>
      getBucketNotificationConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetBucketNotificationConfigurationRequest>
              $request) async {
    return getBucketNotificationConfiguration($call, await $request);
  }

  $async.Future<$0.NotificationConfiguration>
      getBucketNotificationConfiguration($grpc.ServiceCall call,
          $0.GetBucketNotificationConfigurationRequest request);

  $async.Future<$0.GetBucketOwnershipControlsOutput>
      getBucketOwnershipControls_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetBucketOwnershipControlsRequest> $request) async {
    return getBucketOwnershipControls($call, await $request);
  }

  $async.Future<$0.GetBucketOwnershipControlsOutput> getBucketOwnershipControls(
      $grpc.ServiceCall call, $0.GetBucketOwnershipControlsRequest request);

  $async.Future<$0.GetBucketPolicyOutput> getBucketPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketPolicyRequest> $request) async {
    return getBucketPolicy($call, await $request);
  }

  $async.Future<$0.GetBucketPolicyOutput> getBucketPolicy(
      $grpc.ServiceCall call, $0.GetBucketPolicyRequest request);

  $async.Future<$0.GetBucketPolicyStatusOutput> getBucketPolicyStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketPolicyStatusRequest> $request) async {
    return getBucketPolicyStatus($call, await $request);
  }

  $async.Future<$0.GetBucketPolicyStatusOutput> getBucketPolicyStatus(
      $grpc.ServiceCall call, $0.GetBucketPolicyStatusRequest request);

  $async.Future<$0.GetBucketReplicationOutput> getBucketReplication_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketReplicationRequest> $request) async {
    return getBucketReplication($call, await $request);
  }

  $async.Future<$0.GetBucketReplicationOutput> getBucketReplication(
      $grpc.ServiceCall call, $0.GetBucketReplicationRequest request);

  $async.Future<$0.GetBucketRequestPaymentOutput> getBucketRequestPayment_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketRequestPaymentRequest> $request) async {
    return getBucketRequestPayment($call, await $request);
  }

  $async.Future<$0.GetBucketRequestPaymentOutput> getBucketRequestPayment(
      $grpc.ServiceCall call, $0.GetBucketRequestPaymentRequest request);

  $async.Future<$0.GetBucketTaggingOutput> getBucketTagging_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketTaggingRequest> $request) async {
    return getBucketTagging($call, await $request);
  }

  $async.Future<$0.GetBucketTaggingOutput> getBucketTagging(
      $grpc.ServiceCall call, $0.GetBucketTaggingRequest request);

  $async.Future<$0.GetBucketVersioningOutput> getBucketVersioning_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketVersioningRequest> $request) async {
    return getBucketVersioning($call, await $request);
  }

  $async.Future<$0.GetBucketVersioningOutput> getBucketVersioning(
      $grpc.ServiceCall call, $0.GetBucketVersioningRequest request);

  $async.Future<$0.GetBucketWebsiteOutput> getBucketWebsite_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBucketWebsiteRequest> $request) async {
    return getBucketWebsite($call, await $request);
  }

  $async.Future<$0.GetBucketWebsiteOutput> getBucketWebsite(
      $grpc.ServiceCall call, $0.GetBucketWebsiteRequest request);

  $async.Future<$0.GetObjectOutput> getObject_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetObjectRequest> $request) async {
    return getObject($call, await $request);
  }

  $async.Future<$0.GetObjectOutput> getObject(
      $grpc.ServiceCall call, $0.GetObjectRequest request);

  $async.Future<$0.GetObjectAclOutput> getObjectAcl_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetObjectAclRequest> $request) async {
    return getObjectAcl($call, await $request);
  }

  $async.Future<$0.GetObjectAclOutput> getObjectAcl(
      $grpc.ServiceCall call, $0.GetObjectAclRequest request);

  $async.Future<$0.GetObjectAttributesOutput> getObjectAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetObjectAttributesRequest> $request) async {
    return getObjectAttributes($call, await $request);
  }

  $async.Future<$0.GetObjectAttributesOutput> getObjectAttributes(
      $grpc.ServiceCall call, $0.GetObjectAttributesRequest request);

  $async.Future<$0.GetObjectLegalHoldOutput> getObjectLegalHold_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetObjectLegalHoldRequest> $request) async {
    return getObjectLegalHold($call, await $request);
  }

  $async.Future<$0.GetObjectLegalHoldOutput> getObjectLegalHold(
      $grpc.ServiceCall call, $0.GetObjectLegalHoldRequest request);

  $async.Future<$0.GetObjectLockConfigurationOutput>
      getObjectLockConfiguration_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetObjectLockConfigurationRequest> $request) async {
    return getObjectLockConfiguration($call, await $request);
  }

  $async.Future<$0.GetObjectLockConfigurationOutput> getObjectLockConfiguration(
      $grpc.ServiceCall call, $0.GetObjectLockConfigurationRequest request);

  $async.Future<$0.GetObjectRetentionOutput> getObjectRetention_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetObjectRetentionRequest> $request) async {
    return getObjectRetention($call, await $request);
  }

  $async.Future<$0.GetObjectRetentionOutput> getObjectRetention(
      $grpc.ServiceCall call, $0.GetObjectRetentionRequest request);

  $async.Future<$0.GetObjectTaggingOutput> getObjectTagging_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetObjectTaggingRequest> $request) async {
    return getObjectTagging($call, await $request);
  }

  $async.Future<$0.GetObjectTaggingOutput> getObjectTagging(
      $grpc.ServiceCall call, $0.GetObjectTaggingRequest request);

  $async.Future<$0.GetObjectTorrentOutput> getObjectTorrent_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetObjectTorrentRequest> $request) async {
    return getObjectTorrent($call, await $request);
  }

  $async.Future<$0.GetObjectTorrentOutput> getObjectTorrent(
      $grpc.ServiceCall call, $0.GetObjectTorrentRequest request);

  $async.Future<$0.GetPublicAccessBlockOutput> getPublicAccessBlock_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPublicAccessBlockRequest> $request) async {
    return getPublicAccessBlock($call, await $request);
  }

  $async.Future<$0.GetPublicAccessBlockOutput> getPublicAccessBlock(
      $grpc.ServiceCall call, $0.GetPublicAccessBlockRequest request);

  $async.Future<$0.HeadBucketOutput> headBucket_Pre($grpc.ServiceCall $call,
      $async.Future<$0.HeadBucketRequest> $request) async {
    return headBucket($call, await $request);
  }

  $async.Future<$0.HeadBucketOutput> headBucket(
      $grpc.ServiceCall call, $0.HeadBucketRequest request);

  $async.Future<$0.HeadObjectOutput> headObject_Pre($grpc.ServiceCall $call,
      $async.Future<$0.HeadObjectRequest> $request) async {
    return headObject($call, await $request);
  }

  $async.Future<$0.HeadObjectOutput> headObject(
      $grpc.ServiceCall call, $0.HeadObjectRequest request);

  $async.Future<$0.ListBucketAnalyticsConfigurationsOutput>
      listBucketAnalyticsConfigurations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListBucketAnalyticsConfigurationsRequest>
              $request) async {
    return listBucketAnalyticsConfigurations($call, await $request);
  }

  $async.Future<$0.ListBucketAnalyticsConfigurationsOutput>
      listBucketAnalyticsConfigurations($grpc.ServiceCall call,
          $0.ListBucketAnalyticsConfigurationsRequest request);

  $async.Future<$0.ListBucketIntelligentTieringConfigurationsOutput>
      listBucketIntelligentTieringConfigurations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListBucketIntelligentTieringConfigurationsRequest>
              $request) async {
    return listBucketIntelligentTieringConfigurations($call, await $request);
  }

  $async.Future<$0.ListBucketIntelligentTieringConfigurationsOutput>
      listBucketIntelligentTieringConfigurations($grpc.ServiceCall call,
          $0.ListBucketIntelligentTieringConfigurationsRequest request);

  $async.Future<$0.ListBucketInventoryConfigurationsOutput>
      listBucketInventoryConfigurations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListBucketInventoryConfigurationsRequest>
              $request) async {
    return listBucketInventoryConfigurations($call, await $request);
  }

  $async.Future<$0.ListBucketInventoryConfigurationsOutput>
      listBucketInventoryConfigurations($grpc.ServiceCall call,
          $0.ListBucketInventoryConfigurationsRequest request);

  $async.Future<$0.ListBucketMetricsConfigurationsOutput>
      listBucketMetricsConfigurations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListBucketMetricsConfigurationsRequest>
              $request) async {
    return listBucketMetricsConfigurations($call, await $request);
  }

  $async.Future<$0.ListBucketMetricsConfigurationsOutput>
      listBucketMetricsConfigurations($grpc.ServiceCall call,
          $0.ListBucketMetricsConfigurationsRequest request);

  $async.Future<$0.ListBucketsOutput> listBuckets_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListBucketsRequest> $request) async {
    return listBuckets($call, await $request);
  }

  $async.Future<$0.ListBucketsOutput> listBuckets(
      $grpc.ServiceCall call, $0.ListBucketsRequest request);

  $async.Future<$0.ListDirectoryBucketsOutput> listDirectoryBuckets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDirectoryBucketsRequest> $request) async {
    return listDirectoryBuckets($call, await $request);
  }

  $async.Future<$0.ListDirectoryBucketsOutput> listDirectoryBuckets(
      $grpc.ServiceCall call, $0.ListDirectoryBucketsRequest request);

  $async.Future<$0.ListMultipartUploadsOutput> listMultipartUploads_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListMultipartUploadsRequest> $request) async {
    return listMultipartUploads($call, await $request);
  }

  $async.Future<$0.ListMultipartUploadsOutput> listMultipartUploads(
      $grpc.ServiceCall call, $0.ListMultipartUploadsRequest request);

  $async.Future<$0.ListObjectsOutput> listObjects_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListObjectsRequest> $request) async {
    return listObjects($call, await $request);
  }

  $async.Future<$0.ListObjectsOutput> listObjects(
      $grpc.ServiceCall call, $0.ListObjectsRequest request);

  $async.Future<$0.ListObjectsV2Output> listObjectsV2_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListObjectsV2Request> $request) async {
    return listObjectsV2($call, await $request);
  }

  $async.Future<$0.ListObjectsV2Output> listObjectsV2(
      $grpc.ServiceCall call, $0.ListObjectsV2Request request);

  $async.Future<$0.ListObjectVersionsOutput> listObjectVersions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListObjectVersionsRequest> $request) async {
    return listObjectVersions($call, await $request);
  }

  $async.Future<$0.ListObjectVersionsOutput> listObjectVersions(
      $grpc.ServiceCall call, $0.ListObjectVersionsRequest request);

  $async.Future<$0.ListPartsOutput> listParts_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListPartsRequest> $request) async {
    return listParts($call, await $request);
  }

  $async.Future<$0.ListPartsOutput> listParts(
      $grpc.ServiceCall call, $0.ListPartsRequest request);

  $async.Future<$1.Empty> putBucketAbac_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketAbacRequest> $request) async {
    return putBucketAbac($call, await $request);
  }

  $async.Future<$1.Empty> putBucketAbac(
      $grpc.ServiceCall call, $0.PutBucketAbacRequest request);

  $async.Future<$1.Empty> putBucketAccelerateConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutBucketAccelerateConfigurationRequest>
          $request) async {
    return putBucketAccelerateConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> putBucketAccelerateConfiguration(
      $grpc.ServiceCall call,
      $0.PutBucketAccelerateConfigurationRequest request);

  $async.Future<$1.Empty> putBucketAcl_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketAclRequest> $request) async {
    return putBucketAcl($call, await $request);
  }

  $async.Future<$1.Empty> putBucketAcl(
      $grpc.ServiceCall call, $0.PutBucketAclRequest request);

  $async.Future<$1.Empty> putBucketAnalyticsConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutBucketAnalyticsConfigurationRequest> $request) async {
    return putBucketAnalyticsConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> putBucketAnalyticsConfiguration(
      $grpc.ServiceCall call,
      $0.PutBucketAnalyticsConfigurationRequest request);

  $async.Future<$1.Empty> putBucketCors_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketCorsRequest> $request) async {
    return putBucketCors($call, await $request);
  }

  $async.Future<$1.Empty> putBucketCors(
      $grpc.ServiceCall call, $0.PutBucketCorsRequest request);

  $async.Future<$1.Empty> putBucketEncryption_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketEncryptionRequest> $request) async {
    return putBucketEncryption($call, await $request);
  }

  $async.Future<$1.Empty> putBucketEncryption(
      $grpc.ServiceCall call, $0.PutBucketEncryptionRequest request);

  $async.Future<$1.Empty> putBucketIntelligentTieringConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutBucketIntelligentTieringConfigurationRequest>
          $request) async {
    return putBucketIntelligentTieringConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> putBucketIntelligentTieringConfiguration(
      $grpc.ServiceCall call,
      $0.PutBucketIntelligentTieringConfigurationRequest request);

  $async.Future<$1.Empty> putBucketInventoryConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutBucketInventoryConfigurationRequest> $request) async {
    return putBucketInventoryConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> putBucketInventoryConfiguration(
      $grpc.ServiceCall call,
      $0.PutBucketInventoryConfigurationRequest request);

  $async.Future<$0.PutBucketLifecycleConfigurationOutput>
      putBucketLifecycleConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutBucketLifecycleConfigurationRequest>
              $request) async {
    return putBucketLifecycleConfiguration($call, await $request);
  }

  $async.Future<$0.PutBucketLifecycleConfigurationOutput>
      putBucketLifecycleConfiguration($grpc.ServiceCall call,
          $0.PutBucketLifecycleConfigurationRequest request);

  $async.Future<$1.Empty> putBucketLogging_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketLoggingRequest> $request) async {
    return putBucketLogging($call, await $request);
  }

  $async.Future<$1.Empty> putBucketLogging(
      $grpc.ServiceCall call, $0.PutBucketLoggingRequest request);

  $async.Future<$1.Empty> putBucketMetricsConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutBucketMetricsConfigurationRequest> $request) async {
    return putBucketMetricsConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> putBucketMetricsConfiguration(
      $grpc.ServiceCall call, $0.PutBucketMetricsConfigurationRequest request);

  $async.Future<$1.Empty> putBucketNotificationConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutBucketNotificationConfigurationRequest>
          $request) async {
    return putBucketNotificationConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> putBucketNotificationConfiguration(
      $grpc.ServiceCall call,
      $0.PutBucketNotificationConfigurationRequest request);

  $async.Future<$1.Empty> putBucketOwnershipControls_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutBucketOwnershipControlsRequest> $request) async {
    return putBucketOwnershipControls($call, await $request);
  }

  $async.Future<$1.Empty> putBucketOwnershipControls(
      $grpc.ServiceCall call, $0.PutBucketOwnershipControlsRequest request);

  $async.Future<$1.Empty> putBucketPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketPolicyRequest> $request) async {
    return putBucketPolicy($call, await $request);
  }

  $async.Future<$1.Empty> putBucketPolicy(
      $grpc.ServiceCall call, $0.PutBucketPolicyRequest request);

  $async.Future<$1.Empty> putBucketReplication_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketReplicationRequest> $request) async {
    return putBucketReplication($call, await $request);
  }

  $async.Future<$1.Empty> putBucketReplication(
      $grpc.ServiceCall call, $0.PutBucketReplicationRequest request);

  $async.Future<$1.Empty> putBucketRequestPayment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketRequestPaymentRequest> $request) async {
    return putBucketRequestPayment($call, await $request);
  }

  $async.Future<$1.Empty> putBucketRequestPayment(
      $grpc.ServiceCall call, $0.PutBucketRequestPaymentRequest request);

  $async.Future<$1.Empty> putBucketTagging_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketTaggingRequest> $request) async {
    return putBucketTagging($call, await $request);
  }

  $async.Future<$1.Empty> putBucketTagging(
      $grpc.ServiceCall call, $0.PutBucketTaggingRequest request);

  $async.Future<$1.Empty> putBucketVersioning_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketVersioningRequest> $request) async {
    return putBucketVersioning($call, await $request);
  }

  $async.Future<$1.Empty> putBucketVersioning(
      $grpc.ServiceCall call, $0.PutBucketVersioningRequest request);

  $async.Future<$1.Empty> putBucketWebsite_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutBucketWebsiteRequest> $request) async {
    return putBucketWebsite($call, await $request);
  }

  $async.Future<$1.Empty> putBucketWebsite(
      $grpc.ServiceCall call, $0.PutBucketWebsiteRequest request);

  $async.Future<$0.PutObjectOutput> putObject_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutObjectRequest> $request) async {
    return putObject($call, await $request);
  }

  $async.Future<$0.PutObjectOutput> putObject(
      $grpc.ServiceCall call, $0.PutObjectRequest request);

  $async.Future<$0.PutObjectAclOutput> putObjectAcl_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutObjectAclRequest> $request) async {
    return putObjectAcl($call, await $request);
  }

  $async.Future<$0.PutObjectAclOutput> putObjectAcl(
      $grpc.ServiceCall call, $0.PutObjectAclRequest request);

  $async.Future<$0.PutObjectLegalHoldOutput> putObjectLegalHold_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutObjectLegalHoldRequest> $request) async {
    return putObjectLegalHold($call, await $request);
  }

  $async.Future<$0.PutObjectLegalHoldOutput> putObjectLegalHold(
      $grpc.ServiceCall call, $0.PutObjectLegalHoldRequest request);

  $async.Future<$0.PutObjectLockConfigurationOutput>
      putObjectLockConfiguration_Pre($grpc.ServiceCall $call,
          $async.Future<$0.PutObjectLockConfigurationRequest> $request) async {
    return putObjectLockConfiguration($call, await $request);
  }

  $async.Future<$0.PutObjectLockConfigurationOutput> putObjectLockConfiguration(
      $grpc.ServiceCall call, $0.PutObjectLockConfigurationRequest request);

  $async.Future<$0.PutObjectRetentionOutput> putObjectRetention_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutObjectRetentionRequest> $request) async {
    return putObjectRetention($call, await $request);
  }

  $async.Future<$0.PutObjectRetentionOutput> putObjectRetention(
      $grpc.ServiceCall call, $0.PutObjectRetentionRequest request);

  $async.Future<$0.PutObjectTaggingOutput> putObjectTagging_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutObjectTaggingRequest> $request) async {
    return putObjectTagging($call, await $request);
  }

  $async.Future<$0.PutObjectTaggingOutput> putObjectTagging(
      $grpc.ServiceCall call, $0.PutObjectTaggingRequest request);

  $async.Future<$1.Empty> putPublicAccessBlock_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutPublicAccessBlockRequest> $request) async {
    return putPublicAccessBlock($call, await $request);
  }

  $async.Future<$1.Empty> putPublicAccessBlock(
      $grpc.ServiceCall call, $0.PutPublicAccessBlockRequest request);

  $async.Future<$0.RenameObjectOutput> renameObject_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RenameObjectRequest> $request) async {
    return renameObject($call, await $request);
  }

  $async.Future<$0.RenameObjectOutput> renameObject(
      $grpc.ServiceCall call, $0.RenameObjectRequest request);

  $async.Future<$0.RestoreObjectOutput> restoreObject_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RestoreObjectRequest> $request) async {
    return restoreObject($call, await $request);
  }

  $async.Future<$0.RestoreObjectOutput> restoreObject(
      $grpc.ServiceCall call, $0.RestoreObjectRequest request);

  $async.Future<$0.SelectObjectContentOutput> selectObjectContent_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SelectObjectContentRequest> $request) async {
    return selectObjectContent($call, await $request);
  }

  $async.Future<$0.SelectObjectContentOutput> selectObjectContent(
      $grpc.ServiceCall call, $0.SelectObjectContentRequest request);

  $async.Future<$1.Empty> updateBucketMetadataInventoryTableConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateBucketMetadataInventoryTableConfigurationRequest>
          $request) async {
    return updateBucketMetadataInventoryTableConfiguration(
        $call, await $request);
  }

  $async.Future<$1.Empty> updateBucketMetadataInventoryTableConfiguration(
      $grpc.ServiceCall call,
      $0.UpdateBucketMetadataInventoryTableConfigurationRequest request);

  $async.Future<$1.Empty> updateBucketMetadataJournalTableConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateBucketMetadataJournalTableConfigurationRequest>
          $request) async {
    return updateBucketMetadataJournalTableConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> updateBucketMetadataJournalTableConfiguration(
      $grpc.ServiceCall call,
      $0.UpdateBucketMetadataJournalTableConfigurationRequest request);

  $async.Future<$0.UpdateObjectEncryptionResponse> updateObjectEncryption_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateObjectEncryptionRequest> $request) async {
    return updateObjectEncryption($call, await $request);
  }

  $async.Future<$0.UpdateObjectEncryptionResponse> updateObjectEncryption(
      $grpc.ServiceCall call, $0.UpdateObjectEncryptionRequest request);

  $async.Future<$0.UploadPartOutput> uploadPart_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UploadPartRequest> $request) async {
    return uploadPart($call, await $request);
  }

  $async.Future<$0.UploadPartOutput> uploadPart(
      $grpc.ServiceCall call, $0.UploadPartRequest request);

  $async.Future<$0.UploadPartCopyOutput> uploadPartCopy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UploadPartCopyRequest> $request) async {
    return uploadPartCopy($call, await $request);
  }

  $async.Future<$0.UploadPartCopyOutput> uploadPartCopy(
      $grpc.ServiceCall call, $0.UploadPartCopyRequest request);

  $async.Future<$1.Empty> writeGetObjectResponse_Pre($grpc.ServiceCall $call,
      $async.Future<$0.WriteGetObjectResponseRequest> $request) async {
    return writeGetObjectResponse($call, await $request);
  }

  $async.Future<$1.Empty> writeGetObjectResponse(
      $grpc.ServiceCall call, $0.WriteGetObjectResponseRequest request);
}
