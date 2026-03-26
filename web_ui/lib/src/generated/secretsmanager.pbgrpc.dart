// This is a generated file - do not edit.
//
// Generated from secretsmanager.proto.

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
import 'secretsmanager.pb.dart' as $0;

export 'secretsmanager.pb.dart';

/// SecretsManagerService provides secretsmanager API operations.
@$pb.GrpcServiceName('secretsmanager.SecretsManagerService')
class SecretsManagerServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  SecretsManagerServiceClient(super.channel,
      {super.options, super.interceptors});

  /// Retrieves the contents of the encrypted fields SecretString or SecretBinary for up to 20 secrets. To retrieve a single secret, call GetSecretValue. To choose which secrets to retrieve, you can spec...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.BatchGetSecretValueResponse> batchGetSecretValue(
    $0.BatchGetSecretValueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$batchGetSecretValue, request, options: options);
  }

  /// Turns off automatic rotation, and if a rotation is currently in progress, cancels the rotation. If you cancel a rotation in progress, it can leave the VersionStage labels in an unexpected state. Yo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CancelRotateSecretResponse> cancelRotateSecret(
    $0.CancelRotateSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelRotateSecret, request, options: options);
  }

  /// Creates a new secret. A secret can be a password, a set of credentials such as a user name and password, an OAuth token, or other secret information that you store in an encrypted form in Secrets M...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateSecretResponse> createSecret(
    $0.CreateSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSecret, request, options: options);
  }

  /// Deletes the resource-based permission policy attached to the secret. To attach a policy to a secret, use PutResourcePolicy. Secrets Manager generates a CloudTrail log entry when you call this actio...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteResourcePolicyResponse> deleteResourcePolicy(
    $0.DeleteResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourcePolicy, request, options: options);
  }

  /// Deletes a secret and all of its versions. You can specify a recovery window during which you can restore the secret. The minimum recovery window is 7 days. The default recovery window is 30 days. S...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteSecretResponse> deleteSecret(
    $0.DeleteSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSecret, request, options: options);
  }

  /// Retrieves the details of a secret. It does not include the encrypted secret value. Secrets Manager only returns fields that have a value in the response. Secrets Manager generates a CloudTrail log ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeSecretResponse> describeSecret(
    $0.DescribeSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeSecret, request, options: options);
  }

  /// Generates a random password. We recommend that you specify the maximum length and include every character type that the system you are generating a password for can support. By default, Secrets Man...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetRandomPasswordResponse> getRandomPassword(
    $0.GetRandomPasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRandomPassword, request, options: options);
  }

  /// Retrieves the JSON text of the resource-based policy document attached to the secret. For more information about permissions policies attached to a secret, see Permissions policies attached to a se...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetResourcePolicyResponse> getResourcePolicy(
    $0.GetResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResourcePolicy, request, options: options);
  }

  /// Retrieves the contents of the encrypted fields SecretString or SecretBinary from the specified version of a secret, whichever contains content. To retrieve the values for a group of secrets, call B...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSecretValueResponse> getSecretValue(
    $0.GetSecretValueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSecretValue, request, options: options);
  }

  /// Lists the secrets that are stored by Secrets Manager in the Amazon Web Services account, not including secrets that are marked for deletion. To see secrets marked for deletion, use the Secrets Mana...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListSecretsResponse> listSecrets(
    $0.ListSecretsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSecrets, request, options: options);
  }

  /// Lists the versions of a secret. Secrets Manager uses staging labels to indicate the different versions of a secret. For more information, see Secrets Manager concepts: Versions. To list the secrets...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListSecretVersionIdsResponse> listSecretVersionIds(
    $0.ListSecretVersionIdsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSecretVersionIds, request, options: options);
  }

  /// Attaches a resource-based permission policy to a secret. A resource-based policy is optional. For more information, see Authentication and access control for Secrets Manager For information about a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutResourcePolicyResponse> putResourcePolicy(
    $0.PutResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putResourcePolicy, request, options: options);
  }

  /// Creates a new version of your secret by creating a new encrypted value and attaching it to the secret. version can contain a new SecretString value or a new SecretBinary value. Do not call PutSecre...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutSecretValueResponse> putSecretValue(
    $0.PutSecretValueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putSecretValue, request, options: options);
  }

  /// For a secret that is replicated to other Regions, deletes the secret replicas from the Regions you specify. Secrets Manager generates a CloudTrail log entry when you call this action. Do not includ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RemoveRegionsFromReplicationResponse>
      removeRegionsFromReplication(
    $0.RemoveRegionsFromReplicationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeRegionsFromReplication, request,
        options: options);
  }

  /// Replicates the secret to a new Regions. See Multi-Region secrets. Secrets Manager generates a CloudTrail log entry when you call this action. Do not include sensitive information in request paramet...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ReplicateSecretToRegionsResponse>
      replicateSecretToRegions(
    $0.ReplicateSecretToRegionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$replicateSecretToRegions, request,
        options: options);
  }

  /// Cancels the scheduled deletion of a secret by removing the DeletedDate time stamp. You can access a secret again after it has been restored. Secrets Manager generates a CloudTrail log entry when yo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RestoreSecretResponse> restoreSecret(
    $0.RestoreSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$restoreSecret, request, options: options);
  }

  /// Configures and starts the asynchronous process of rotating the secret. For information about rotation, see Rotate secrets in the Secrets Manager User Guide. If you include the configuration paramet...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RotateSecretResponse> rotateSecret(
    $0.RotateSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$rotateSecret, request, options: options);
  }

  /// Removes the link between the replica secret and the primary secret and promotes the replica to a primary secret in the replica Region. You must call this operation from the Region in which you want...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopReplicationToReplicaResponse>
      stopReplicationToReplica(
    $0.StopReplicationToReplicaRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopReplicationToReplica, request,
        options: options);
  }

  /// Attaches tags to a secret. Tags consist of a key name and a value. Tags are part of the secret's metadata. They are not associated with specific versions of the secret. This operation appends tags ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Removes specific tags from a secret. This operation is idempotent. If a requested tag is not attached to the secret, no error is returned and the secret metadata is unchanged. If you use tags as pa...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Modifies the details of a secret, including metadata and the secret value. To change the secret value, you can also use PutSecretValue. To change the rotation configuration of a secret, use RotateS...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateSecretResponse> updateSecret(
    $0.UpdateSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSecret, request, options: options);
  }

  /// Modifies the staging labels attached to a version of a secret. Secrets Manager uses staging labels to track a version as it progresses through the secret rotation process. Each staging label can be...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateSecretVersionStageResponse>
      updateSecretVersionStage(
    $0.UpdateSecretVersionStageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSecretVersionStage, request,
        options: options);
  }

  /// Validates that a resource policy does not grant a wide range of principals access to your secret. A resource-based policy is optional for secrets. The API performs three checks when validating the ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ValidateResourcePolicyResponse>
      validateResourcePolicy(
    $0.ValidateResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$validateResourcePolicy, request,
        options: options);
  }

  // method descriptors

  static final _$batchGetSecretValue = $grpc.ClientMethod<
          $0.BatchGetSecretValueRequest, $0.BatchGetSecretValueResponse>(
      '/secretsmanager.SecretsManagerService/BatchGetSecretValue',
      ($0.BatchGetSecretValueRequest value) => value.writeToBuffer(),
      $0.BatchGetSecretValueResponse.fromBuffer);
  static final _$cancelRotateSecret = $grpc.ClientMethod<
          $0.CancelRotateSecretRequest, $0.CancelRotateSecretResponse>(
      '/secretsmanager.SecretsManagerService/CancelRotateSecret',
      ($0.CancelRotateSecretRequest value) => value.writeToBuffer(),
      $0.CancelRotateSecretResponse.fromBuffer);
  static final _$createSecret =
      $grpc.ClientMethod<$0.CreateSecretRequest, $0.CreateSecretResponse>(
          '/secretsmanager.SecretsManagerService/CreateSecret',
          ($0.CreateSecretRequest value) => value.writeToBuffer(),
          $0.CreateSecretResponse.fromBuffer);
  static final _$deleteResourcePolicy = $grpc.ClientMethod<
          $0.DeleteResourcePolicyRequest, $0.DeleteResourcePolicyResponse>(
      '/secretsmanager.SecretsManagerService/DeleteResourcePolicy',
      ($0.DeleteResourcePolicyRequest value) => value.writeToBuffer(),
      $0.DeleteResourcePolicyResponse.fromBuffer);
  static final _$deleteSecret =
      $grpc.ClientMethod<$0.DeleteSecretRequest, $0.DeleteSecretResponse>(
          '/secretsmanager.SecretsManagerService/DeleteSecret',
          ($0.DeleteSecretRequest value) => value.writeToBuffer(),
          $0.DeleteSecretResponse.fromBuffer);
  static final _$describeSecret =
      $grpc.ClientMethod<$0.DescribeSecretRequest, $0.DescribeSecretResponse>(
          '/secretsmanager.SecretsManagerService/DescribeSecret',
          ($0.DescribeSecretRequest value) => value.writeToBuffer(),
          $0.DescribeSecretResponse.fromBuffer);
  static final _$getRandomPassword = $grpc.ClientMethod<
          $0.GetRandomPasswordRequest, $0.GetRandomPasswordResponse>(
      '/secretsmanager.SecretsManagerService/GetRandomPassword',
      ($0.GetRandomPasswordRequest value) => value.writeToBuffer(),
      $0.GetRandomPasswordResponse.fromBuffer);
  static final _$getResourcePolicy = $grpc.ClientMethod<
          $0.GetResourcePolicyRequest, $0.GetResourcePolicyResponse>(
      '/secretsmanager.SecretsManagerService/GetResourcePolicy',
      ($0.GetResourcePolicyRequest value) => value.writeToBuffer(),
      $0.GetResourcePolicyResponse.fromBuffer);
  static final _$getSecretValue =
      $grpc.ClientMethod<$0.GetSecretValueRequest, $0.GetSecretValueResponse>(
          '/secretsmanager.SecretsManagerService/GetSecretValue',
          ($0.GetSecretValueRequest value) => value.writeToBuffer(),
          $0.GetSecretValueResponse.fromBuffer);
  static final _$listSecrets =
      $grpc.ClientMethod<$0.ListSecretsRequest, $0.ListSecretsResponse>(
          '/secretsmanager.SecretsManagerService/ListSecrets',
          ($0.ListSecretsRequest value) => value.writeToBuffer(),
          $0.ListSecretsResponse.fromBuffer);
  static final _$listSecretVersionIds = $grpc.ClientMethod<
          $0.ListSecretVersionIdsRequest, $0.ListSecretVersionIdsResponse>(
      '/secretsmanager.SecretsManagerService/ListSecretVersionIds',
      ($0.ListSecretVersionIdsRequest value) => value.writeToBuffer(),
      $0.ListSecretVersionIdsResponse.fromBuffer);
  static final _$putResourcePolicy = $grpc.ClientMethod<
          $0.PutResourcePolicyRequest, $0.PutResourcePolicyResponse>(
      '/secretsmanager.SecretsManagerService/PutResourcePolicy',
      ($0.PutResourcePolicyRequest value) => value.writeToBuffer(),
      $0.PutResourcePolicyResponse.fromBuffer);
  static final _$putSecretValue =
      $grpc.ClientMethod<$0.PutSecretValueRequest, $0.PutSecretValueResponse>(
          '/secretsmanager.SecretsManagerService/PutSecretValue',
          ($0.PutSecretValueRequest value) => value.writeToBuffer(),
          $0.PutSecretValueResponse.fromBuffer);
  static final _$removeRegionsFromReplication = $grpc.ClientMethod<
          $0.RemoveRegionsFromReplicationRequest,
          $0.RemoveRegionsFromReplicationResponse>(
      '/secretsmanager.SecretsManagerService/RemoveRegionsFromReplication',
      ($0.RemoveRegionsFromReplicationRequest value) => value.writeToBuffer(),
      $0.RemoveRegionsFromReplicationResponse.fromBuffer);
  static final _$replicateSecretToRegions = $grpc.ClientMethod<
          $0.ReplicateSecretToRegionsRequest,
          $0.ReplicateSecretToRegionsResponse>(
      '/secretsmanager.SecretsManagerService/ReplicateSecretToRegions',
      ($0.ReplicateSecretToRegionsRequest value) => value.writeToBuffer(),
      $0.ReplicateSecretToRegionsResponse.fromBuffer);
  static final _$restoreSecret =
      $grpc.ClientMethod<$0.RestoreSecretRequest, $0.RestoreSecretResponse>(
          '/secretsmanager.SecretsManagerService/RestoreSecret',
          ($0.RestoreSecretRequest value) => value.writeToBuffer(),
          $0.RestoreSecretResponse.fromBuffer);
  static final _$rotateSecret =
      $grpc.ClientMethod<$0.RotateSecretRequest, $0.RotateSecretResponse>(
          '/secretsmanager.SecretsManagerService/RotateSecret',
          ($0.RotateSecretRequest value) => value.writeToBuffer(),
          $0.RotateSecretResponse.fromBuffer);
  static final _$stopReplicationToReplica = $grpc.ClientMethod<
          $0.StopReplicationToReplicaRequest,
          $0.StopReplicationToReplicaResponse>(
      '/secretsmanager.SecretsManagerService/StopReplicationToReplica',
      ($0.StopReplicationToReplicaRequest value) => value.writeToBuffer(),
      $0.StopReplicationToReplicaResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $1.Empty>(
          '/secretsmanager.SecretsManagerService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $1.Empty>(
          '/secretsmanager.SecretsManagerService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateSecret =
      $grpc.ClientMethod<$0.UpdateSecretRequest, $0.UpdateSecretResponse>(
          '/secretsmanager.SecretsManagerService/UpdateSecret',
          ($0.UpdateSecretRequest value) => value.writeToBuffer(),
          $0.UpdateSecretResponse.fromBuffer);
  static final _$updateSecretVersionStage = $grpc.ClientMethod<
          $0.UpdateSecretVersionStageRequest,
          $0.UpdateSecretVersionStageResponse>(
      '/secretsmanager.SecretsManagerService/UpdateSecretVersionStage',
      ($0.UpdateSecretVersionStageRequest value) => value.writeToBuffer(),
      $0.UpdateSecretVersionStageResponse.fromBuffer);
  static final _$validateResourcePolicy = $grpc.ClientMethod<
          $0.ValidateResourcePolicyRequest, $0.ValidateResourcePolicyResponse>(
      '/secretsmanager.SecretsManagerService/ValidateResourcePolicy',
      ($0.ValidateResourcePolicyRequest value) => value.writeToBuffer(),
      $0.ValidateResourcePolicyResponse.fromBuffer);
}

@$pb.GrpcServiceName('secretsmanager.SecretsManagerService')
abstract class SecretsManagerServiceBase extends $grpc.Service {
  $core.String get $name => 'secretsmanager.SecretsManagerService';

  SecretsManagerServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.BatchGetSecretValueRequest,
            $0.BatchGetSecretValueResponse>(
        'BatchGetSecretValue',
        batchGetSecretValue_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.BatchGetSecretValueRequest.fromBuffer(value),
        ($0.BatchGetSecretValueResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelRotateSecretRequest,
            $0.CancelRotateSecretResponse>(
        'CancelRotateSecret',
        cancelRotateSecret_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelRotateSecretRequest.fromBuffer(value),
        ($0.CancelRotateSecretResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateSecretRequest, $0.CreateSecretResponse>(
            'CreateSecret',
            createSecret_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateSecretRequest.fromBuffer(value),
            ($0.CreateSecretResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourcePolicyRequest,
            $0.DeleteResourcePolicyResponse>(
        'DeleteResourcePolicy',
        deleteResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourcePolicyRequest.fromBuffer(value),
        ($0.DeleteResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteSecretRequest, $0.DeleteSecretResponse>(
            'DeleteSecret',
            deleteSecret_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteSecretRequest.fromBuffer(value),
            ($0.DeleteSecretResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeSecretRequest,
            $0.DescribeSecretResponse>(
        'DescribeSecret',
        describeSecret_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeSecretRequest.fromBuffer(value),
        ($0.DescribeSecretResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRandomPasswordRequest,
            $0.GetRandomPasswordResponse>(
        'GetRandomPassword',
        getRandomPassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRandomPasswordRequest.fromBuffer(value),
        ($0.GetRandomPasswordResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourcePolicyRequest,
            $0.GetResourcePolicyResponse>(
        'GetResourcePolicy',
        getResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourcePolicyRequest.fromBuffer(value),
        ($0.GetResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSecretValueRequest,
            $0.GetSecretValueResponse>(
        'GetSecretValue',
        getSecretValue_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSecretValueRequest.fromBuffer(value),
        ($0.GetSecretValueResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListSecretsRequest, $0.ListSecretsResponse>(
            'ListSecrets',
            listSecrets_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListSecretsRequest.fromBuffer(value),
            ($0.ListSecretsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSecretVersionIdsRequest,
            $0.ListSecretVersionIdsResponse>(
        'ListSecretVersionIds',
        listSecretVersionIds_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSecretVersionIdsRequest.fromBuffer(value),
        ($0.ListSecretVersionIdsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutResourcePolicyRequest,
            $0.PutResourcePolicyResponse>(
        'PutResourcePolicy',
        putResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutResourcePolicyRequest.fromBuffer(value),
        ($0.PutResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutSecretValueRequest,
            $0.PutSecretValueResponse>(
        'PutSecretValue',
        putSecretValue_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutSecretValueRequest.fromBuffer(value),
        ($0.PutSecretValueResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemoveRegionsFromReplicationRequest,
            $0.RemoveRegionsFromReplicationResponse>(
        'RemoveRegionsFromReplication',
        removeRegionsFromReplication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemoveRegionsFromReplicationRequest.fromBuffer(value),
        ($0.RemoveRegionsFromReplicationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ReplicateSecretToRegionsRequest,
            $0.ReplicateSecretToRegionsResponse>(
        'ReplicateSecretToRegions',
        replicateSecretToRegions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ReplicateSecretToRegionsRequest.fromBuffer(value),
        ($0.ReplicateSecretToRegionsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RestoreSecretRequest, $0.RestoreSecretResponse>(
            'RestoreSecret',
            restoreSecret_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RestoreSecretRequest.fromBuffer(value),
            ($0.RestoreSecretResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RotateSecretRequest, $0.RotateSecretResponse>(
            'RotateSecret',
            rotateSecret_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RotateSecretRequest.fromBuffer(value),
            ($0.RotateSecretResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopReplicationToReplicaRequest,
            $0.StopReplicationToReplicaResponse>(
        'StopReplicationToReplica',
        stopReplicationToReplica_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopReplicationToReplicaRequest.fromBuffer(value),
        ($0.StopReplicationToReplicaResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceRequest, $1.Empty>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagResourceRequest, $1.Empty>(
        'UntagResource',
        untagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateSecretRequest, $0.UpdateSecretResponse>(
            'UpdateSecret',
            updateSecret_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateSecretRequest.fromBuffer(value),
            ($0.UpdateSecretResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateSecretVersionStageRequest,
            $0.UpdateSecretVersionStageResponse>(
        'UpdateSecretVersionStage',
        updateSecretVersionStage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateSecretVersionStageRequest.fromBuffer(value),
        ($0.UpdateSecretVersionStageResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ValidateResourcePolicyRequest,
            $0.ValidateResourcePolicyResponse>(
        'ValidateResourcePolicy',
        validateResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ValidateResourcePolicyRequest.fromBuffer(value),
        ($0.ValidateResourcePolicyResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.BatchGetSecretValueResponse> batchGetSecretValue_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.BatchGetSecretValueRequest> $request) async {
    return batchGetSecretValue($call, await $request);
  }

  $async.Future<$0.BatchGetSecretValueResponse> batchGetSecretValue(
      $grpc.ServiceCall call, $0.BatchGetSecretValueRequest request);

  $async.Future<$0.CancelRotateSecretResponse> cancelRotateSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelRotateSecretRequest> $request) async {
    return cancelRotateSecret($call, await $request);
  }

  $async.Future<$0.CancelRotateSecretResponse> cancelRotateSecret(
      $grpc.ServiceCall call, $0.CancelRotateSecretRequest request);

  $async.Future<$0.CreateSecretResponse> createSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateSecretRequest> $request) async {
    return createSecret($call, await $request);
  }

  $async.Future<$0.CreateSecretResponse> createSecret(
      $grpc.ServiceCall call, $0.CreateSecretRequest request);

  $async.Future<$0.DeleteResourcePolicyResponse> deleteResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourcePolicyRequest> $request) async {
    return deleteResourcePolicy($call, await $request);
  }

  $async.Future<$0.DeleteResourcePolicyResponse> deleteResourcePolicy(
      $grpc.ServiceCall call, $0.DeleteResourcePolicyRequest request);

  $async.Future<$0.DeleteSecretResponse> deleteSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteSecretRequest> $request) async {
    return deleteSecret($call, await $request);
  }

  $async.Future<$0.DeleteSecretResponse> deleteSecret(
      $grpc.ServiceCall call, $0.DeleteSecretRequest request);

  $async.Future<$0.DescribeSecretResponse> describeSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeSecretRequest> $request) async {
    return describeSecret($call, await $request);
  }

  $async.Future<$0.DescribeSecretResponse> describeSecret(
      $grpc.ServiceCall call, $0.DescribeSecretRequest request);

  $async.Future<$0.GetRandomPasswordResponse> getRandomPassword_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRandomPasswordRequest> $request) async {
    return getRandomPassword($call, await $request);
  }

  $async.Future<$0.GetRandomPasswordResponse> getRandomPassword(
      $grpc.ServiceCall call, $0.GetRandomPasswordRequest request);

  $async.Future<$0.GetResourcePolicyResponse> getResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResourcePolicyRequest> $request) async {
    return getResourcePolicy($call, await $request);
  }

  $async.Future<$0.GetResourcePolicyResponse> getResourcePolicy(
      $grpc.ServiceCall call, $0.GetResourcePolicyRequest request);

  $async.Future<$0.GetSecretValueResponse> getSecretValue_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSecretValueRequest> $request) async {
    return getSecretValue($call, await $request);
  }

  $async.Future<$0.GetSecretValueResponse> getSecretValue(
      $grpc.ServiceCall call, $0.GetSecretValueRequest request);

  $async.Future<$0.ListSecretsResponse> listSecrets_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListSecretsRequest> $request) async {
    return listSecrets($call, await $request);
  }

  $async.Future<$0.ListSecretsResponse> listSecrets(
      $grpc.ServiceCall call, $0.ListSecretsRequest request);

  $async.Future<$0.ListSecretVersionIdsResponse> listSecretVersionIds_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSecretVersionIdsRequest> $request) async {
    return listSecretVersionIds($call, await $request);
  }

  $async.Future<$0.ListSecretVersionIdsResponse> listSecretVersionIds(
      $grpc.ServiceCall call, $0.ListSecretVersionIdsRequest request);

  $async.Future<$0.PutResourcePolicyResponse> putResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutResourcePolicyRequest> $request) async {
    return putResourcePolicy($call, await $request);
  }

  $async.Future<$0.PutResourcePolicyResponse> putResourcePolicy(
      $grpc.ServiceCall call, $0.PutResourcePolicyRequest request);

  $async.Future<$0.PutSecretValueResponse> putSecretValue_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutSecretValueRequest> $request) async {
    return putSecretValue($call, await $request);
  }

  $async.Future<$0.PutSecretValueResponse> putSecretValue(
      $grpc.ServiceCall call, $0.PutSecretValueRequest request);

  $async.Future<$0.RemoveRegionsFromReplicationResponse>
      removeRegionsFromReplication_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.RemoveRegionsFromReplicationRequest>
              $request) async {
    return removeRegionsFromReplication($call, await $request);
  }

  $async.Future<$0.RemoveRegionsFromReplicationResponse>
      removeRegionsFromReplication($grpc.ServiceCall call,
          $0.RemoveRegionsFromReplicationRequest request);

  $async.Future<$0.ReplicateSecretToRegionsResponse>
      replicateSecretToRegions_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ReplicateSecretToRegionsRequest> $request) async {
    return replicateSecretToRegions($call, await $request);
  }

  $async.Future<$0.ReplicateSecretToRegionsResponse> replicateSecretToRegions(
      $grpc.ServiceCall call, $0.ReplicateSecretToRegionsRequest request);

  $async.Future<$0.RestoreSecretResponse> restoreSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RestoreSecretRequest> $request) async {
    return restoreSecret($call, await $request);
  }

  $async.Future<$0.RestoreSecretResponse> restoreSecret(
      $grpc.ServiceCall call, $0.RestoreSecretRequest request);

  $async.Future<$0.RotateSecretResponse> rotateSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RotateSecretRequest> $request) async {
    return rotateSecret($call, await $request);
  }

  $async.Future<$0.RotateSecretResponse> rotateSecret(
      $grpc.ServiceCall call, $0.RotateSecretRequest request);

  $async.Future<$0.StopReplicationToReplicaResponse>
      stopReplicationToReplica_Pre($grpc.ServiceCall $call,
          $async.Future<$0.StopReplicationToReplicaRequest> $request) async {
    return stopReplicationToReplica($call, await $request);
  }

  $async.Future<$0.StopReplicationToReplicaResponse> stopReplicationToReplica(
      $grpc.ServiceCall call, $0.StopReplicationToReplicaRequest request);

  $async.Future<$1.Empty> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$1.Empty> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$1.Empty> untagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$1.Empty> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.UpdateSecretResponse> updateSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateSecretRequest> $request) async {
    return updateSecret($call, await $request);
  }

  $async.Future<$0.UpdateSecretResponse> updateSecret(
      $grpc.ServiceCall call, $0.UpdateSecretRequest request);

  $async.Future<$0.UpdateSecretVersionStageResponse>
      updateSecretVersionStage_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateSecretVersionStageRequest> $request) async {
    return updateSecretVersionStage($call, await $request);
  }

  $async.Future<$0.UpdateSecretVersionStageResponse> updateSecretVersionStage(
      $grpc.ServiceCall call, $0.UpdateSecretVersionStageRequest request);

  $async.Future<$0.ValidateResourcePolicyResponse> validateResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ValidateResourcePolicyRequest> $request) async {
    return validateResourcePolicy($call, await $request);
  }

  $async.Future<$0.ValidateResourcePolicyResponse> validateResourcePolicy(
      $grpc.ServiceCall call, $0.ValidateResourcePolicyRequest request);
}
