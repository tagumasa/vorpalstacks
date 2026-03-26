// This is a generated file - do not edit.
//
// Generated from kms.proto.

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
import 'kms.pb.dart' as $0;

export 'kms.pb.dart';

/// KMSService provides kms API operations.
@$pb.GrpcServiceName('kms.KMSService')
class KMSServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  KMSServiceClient(super.channel, {super.options, super.interceptors});

  /// Cancels the deletion of a KMS key. When this operation succeeds, the key state of the KMS key is Disabled. To enable the KMS key, use EnableKey. For more information about scheduling and canceling ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CancelKeyDeletionResponse> cancelKeyDeletion(
    $0.CancelKeyDeletionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelKeyDeletion, request, options: options);
  }

  /// Connects or reconnects a custom key store to its backing key store. For an CloudHSM key store, ConnectCustomKeyStore connects the key store to its associated CloudHSM cluster. For an external key s...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ConnectCustomKeyStoreResponse> connectCustomKeyStore(
    $0.ConnectCustomKeyStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$connectCustomKeyStore, request, options: options);
  }

  /// Creates a friendly name for a KMS key. Adding, deleting, or updating an alias can allow or deny permission to the KMS key. For details, see ABAC for KMS in the Key Management Service Developer Guid...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> createAlias(
    $0.CreateAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAlias, request, options: options);
  }

  /// Creates a custom key store backed by a key store that you own and manage. When you use a KMS key in a custom key store for a cryptographic operation, the cryptographic operation is actually perform...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateCustomKeyStoreResponse> createCustomKeyStore(
    $0.CreateCustomKeyStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createCustomKeyStore, request, options: options);
  }

  /// Adds a grant to a KMS key. A grant is a policy instrument that allows Amazon Web Services principals to use KMS keys in cryptographic operations. It also can allow them to view a KMS key (DescribeK...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateGrantResponse> createGrant(
    $0.CreateGrantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createGrant, request, options: options);
  }

  /// Creates a unique customer managed KMS key in your Amazon Web Services account and Region. You can use a KMS key in cryptographic operations, such as encryption and signing. Some Amazon Web Services...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateKeyResponse> createKey(
    $0.CreateKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createKey, request, options: options);
  }

  /// Decrypts ciphertext that was encrypted by a KMS key using any of the following operations: Encrypt GenerateDataKey GenerateDataKeyPair GenerateDataKeyWithoutPlaintext GenerateDataKeyPairWithoutPlai...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DecryptResponse> decrypt(
    $0.DecryptRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$decrypt, request, options: options);
  }

  /// Deletes the specified alias. Adding, deleting, or updating an alias can allow or deny permission to the KMS key. For details, see ABAC for KMS in the Key Management Service Developer Guide. Because...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteAlias(
    $0.DeleteAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAlias, request, options: options);
  }

  /// Deletes a custom key store. This operation does not affect any backing elements of the custom key store. It does not delete the CloudHSM cluster that is associated with an CloudHSM key store, or af...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteCustomKeyStoreResponse> deleteCustomKeyStore(
    $0.DeleteCustomKeyStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteCustomKeyStore, request, options: options);
  }

  /// Deletes key material that was previously imported. This operation makes the specified KMS key temporarily unusable. To restore the usability of the KMS key, reimport the same key material. For more...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteImportedKeyMaterialResponse>
      deleteImportedKeyMaterial(
    $0.DeleteImportedKeyMaterialRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteImportedKeyMaterial, request,
        options: options);
  }

  /// Derives a shared secret using a key agreement algorithm. You must use an asymmetric NIST-standard elliptic curve (ECC) or SM2 (China Regions only) KMS key pair with a KeyUsage value of KEY_AGREEMEN...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeriveSharedSecretResponse> deriveSharedSecret(
    $0.DeriveSharedSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deriveSharedSecret, request, options: options);
  }

  /// Gets information about custom key stores in the account and Region. This operation is part of the custom key stores feature in KMS, which combines the convenience and extensive integration of KMS w...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeCustomKeyStoresResponse>
      describeCustomKeyStores(
    $0.DescribeCustomKeyStoresRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeCustomKeyStores, request,
        options: options);
  }

  /// Provides detailed information about a KMS key. You can run DescribeKey on a customer managed key or an Amazon Web Services managed key. This detailed information includes the key ARN, creation date...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeKeyResponse> describeKey(
    $0.DescribeKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeKey, request, options: options);
  }

  /// Sets the state of a KMS key to disabled. This change temporarily prevents use of the KMS key for cryptographic operations. The KMS key that you use for this operation must be in a compatible key st...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> disableKey(
    $0.DisableKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableKey, request, options: options);
  }

  /// Disables automatic rotation of the key material of the specified symmetric encryption KMS key. Automatic key rotation is supported only on symmetric encryption KMS keys. You cannot enable automatic...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> disableKeyRotation(
    $0.DisableKeyRotationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableKeyRotation, request, options: options);
  }

  /// Disconnects the custom key store from its backing key store. This operation disconnects an CloudHSM key store from its associated CloudHSM cluster or disconnects an external key store from the exte...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DisconnectCustomKeyStoreResponse>
      disconnectCustomKeyStore(
    $0.DisconnectCustomKeyStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disconnectCustomKeyStore, request,
        options: options);
  }

  /// Sets the key state of a KMS key to enabled. This allows you to use the KMS key for cryptographic operations. The KMS key that you use for this operation must be in a compatible key state. For detai...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> enableKey(
    $0.EnableKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableKey, request, options: options);
  }

  /// Enables automatic rotation of the key material of the specified symmetric encryption KMS key. By default, when you enable automatic rotation of a customer managed KMS key, KMS rotates the key mater...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> enableKeyRotation(
    $0.EnableKeyRotationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableKeyRotation, request, options: options);
  }

  /// Encrypts plaintext of up to 4,096 bytes using a KMS key. You can use a symmetric or asymmetric KMS key with a KeyUsage of ENCRYPT_DECRYPT. You can use this operation to encrypt small amounts of arb...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.EncryptResponse> encrypt(
    $0.EncryptRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$encrypt, request, options: options);
  }

  /// Returns a unique symmetric data key for use outside of KMS. This operation returns a plaintext copy of the data key and a copy that is encrypted under a symmetric encryption KMS key that you specif...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GenerateDataKeyResponse> generateDataKey(
    $0.GenerateDataKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateDataKey, request, options: options);
  }

  /// Returns a unique asymmetric data key pair for use outside of KMS. This operation returns a plaintext public key, a plaintext private key, and a copy of the private key that is encrypted under the s...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GenerateDataKeyPairResponse> generateDataKeyPair(
    $0.GenerateDataKeyPairRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateDataKeyPair, request, options: options);
  }

  /// Returns a unique asymmetric data key pair for use outside of KMS. This operation returns a plaintext public key and a copy of the private key that is encrypted under the symmetric encryption KMS ke...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GenerateDataKeyPairWithoutPlaintextResponse>
      generateDataKeyPairWithoutPlaintext(
    $0.GenerateDataKeyPairWithoutPlaintextRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateDataKeyPairWithoutPlaintext, request,
        options: options);
  }

  /// Returns a unique symmetric data key for use outside of KMS. This operation returns a data key that is encrypted under a symmetric encryption KMS key that you specify. The bytes in the key are rando...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GenerateDataKeyWithoutPlaintextResponse>
      generateDataKeyWithoutPlaintext(
    $0.GenerateDataKeyWithoutPlaintextRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateDataKeyWithoutPlaintext, request,
        options: options);
  }

  /// Generates a hash-based message authentication code (HMAC) for a message using an HMAC KMS key and a MAC algorithm that the key supports. HMAC KMS keys and the HMAC algorithms that KMS uses conform ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GenerateMacResponse> generateMac(
    $0.GenerateMacRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateMac, request, options: options);
  }

  /// Returns a random byte string that is cryptographically secure. You must use the NumberOfBytes parameter to specify the length of the random byte string. There is no default value for string length....
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GenerateRandomResponse> generateRandom(
    $0.GenerateRandomRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateRandom, request, options: options);
  }

  /// Gets a key policy attached to the specified KMS key. Cross-account use: No. You cannot perform this operation on a KMS key in a different Amazon Web Services account. Required permissions: kms:GetK...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetKeyPolicyResponse> getKeyPolicy(
    $0.GetKeyPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getKeyPolicy, request, options: options);
  }

  /// Provides detailed information about the rotation status for a KMS key, including whether automatic rotation of the key material is enabled for the specified KMS key, the rotation period, and the ne...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetKeyRotationStatusResponse> getKeyRotationStatus(
    $0.GetKeyRotationStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getKeyRotationStatus, request, options: options);
  }

  /// Returns the public key and an import token you need to import or reimport key material for a KMS key. By default, KMS keys are created with key material that KMS generates. This operation supports ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetParametersForImportResponse>
      getParametersForImport(
    $0.GetParametersForImportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getParametersForImport, request,
        options: options);
  }

  /// Returns the public key of an asymmetric KMS key. Unlike the private key of a asymmetric KMS key, which never leaves KMS unencrypted, callers with kms:GetPublicKey permission can download the public...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetPublicKeyResponse> getPublicKey(
    $0.GetPublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPublicKey, request, options: options);
  }

  /// Imports or reimports key material into an existing KMS key that was created without key material. You can also use this operation to set or update the expiration model and expiration date of the im...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ImportKeyMaterialResponse> importKeyMaterial(
    $0.ImportKeyMaterialRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$importKeyMaterial, request, options: options);
  }

  /// Gets a list of aliases in the caller's Amazon Web Services account and region. For more information about aliases, see CreateAlias. By default, the ListAliases operation returns all aliases in the ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListAliasesResponse> listAliases(
    $0.ListAliasesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAliases, request, options: options);
  }

  /// Gets a list of all grants for the specified KMS key. You must specify the KMS key in all requests. You can filter the grant list by grant ID or grantee principal. For detailed information about gra...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListGrantsResponse> listGrants(
    $0.ListGrantsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listGrants, request, options: options);
  }

  /// Gets the names of the key policies that are attached to a KMS key. This operation is designed to get policy names that you can use in a GetKeyPolicy operation. However, the only valid policy name i...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListKeyPoliciesResponse> listKeyPolicies(
    $0.ListKeyPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listKeyPolicies, request, options: options);
  }

  /// Returns information about the key materials associated with the specified KMS key. You can use the optional IncludeKeyMaterial parameter to control which key materials are included in the response....
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListKeyRotationsResponse> listKeyRotations(
    $0.ListKeyRotationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listKeyRotations, request, options: options);
  }

  /// Gets a list of all KMS keys in the caller's Amazon Web Services account and Region. Cross-account use: No. You cannot perform this operation on a KMS key in a different Amazon Web Services account....
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListKeysResponse> listKeys(
    $0.ListKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listKeys, request, options: options);
  }

  /// Returns all tags on the specified KMS key. For general information about tags, including the format and syntax, see Tagging Amazon Web Services resources in the Amazon Web Services General Referenc...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListResourceTagsResponse> listResourceTags(
    $0.ListResourceTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listResourceTags, request, options: options);
  }

  /// Returns information about all grants in the Amazon Web Services account and Region that have the specified retiring principal. You can specify any principal in your Amazon Web Services account. The...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListGrantsResponse> listRetirableGrants(
    $0.ListRetirableGrantsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRetirableGrants, request, options: options);
  }

  /// Attaches a key policy to the specified KMS key. For more information about key policies, see Key Policies in the Key Management Service Developer Guide. For help writing and formatting a JSON polic...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putKeyPolicy(
    $0.PutKeyPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putKeyPolicy, request, options: options);
  }

  /// Decrypts ciphertext and then reencrypts it entirely within KMS. You can use this operation to change the KMS key under which data is encrypted, such as when you manually rotate a KMS key or change ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ReEncryptResponse> reEncrypt(
    $0.ReEncryptRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$reEncrypt, request, options: options);
  }

  /// Replicates a multi-Region key into the specified Region. This operation creates a multi-Region replica key based on a multi-Region primary key in a different Region of the same Amazon Web Services ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ReplicateKeyResponse> replicateKey(
    $0.ReplicateKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$replicateKey, request, options: options);
  }

  /// Deletes a grant. Typically, you retire a grant when you no longer need its permissions. To identify the grant to retire, use a grant token, or both the grant ID and a key identifier (key ID or key ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> retireGrant(
    $0.RetireGrantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$retireGrant, request, options: options);
  }

  /// Deletes the specified grant. You revoke a grant to terminate the permissions that the grant allows. For more information, see Retiring and revoking grants in the Key Management Service Developer Gu...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> revokeGrant(
    $0.RevokeGrantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$revokeGrant, request, options: options);
  }

  /// Immediately initiates rotation of the key material of the specified symmetric encryption KMS key. You can perform on-demand rotation of the key material in customer managed KMS keys, regardless of ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RotateKeyOnDemandResponse> rotateKeyOnDemand(
    $0.RotateKeyOnDemandRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$rotateKeyOnDemand, request, options: options);
  }

  /// Schedules the deletion of a KMS key. By default, KMS applies a waiting period of 30 days, but you can specify a waiting period of 7-30 days. When this operation is successful, the key state of the ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ScheduleKeyDeletionResponse> scheduleKeyDeletion(
    $0.ScheduleKeyDeletionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$scheduleKeyDeletion, request, options: options);
  }

  /// Creates a digital signature for a message or message digest by using the private key in an asymmetric signing KMS key. To verify the signature, use the Verify operation, or use the public key in th...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SignResponse> sign(
    $0.SignRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sign, request, options: options);
  }

  /// Adds or edits tags on a customer managed key. Tagging or untagging a KMS key can allow or deny permission to the KMS key. For details, see ABAC for KMS in the Key Management Service Developer Guide...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Deletes tags from a customer managed key. To delete a tag, specify the tag key and the KMS key. Tagging or untagging a KMS key can allow or deny permission to the KMS key. For details, see ABAC for...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Associates an existing KMS alias with a different KMS key. Each alias is associated with only one KMS key at a time, although a KMS key can have multiple aliases. The alias and the KMS key must be ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> updateAlias(
    $0.UpdateAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAlias, request, options: options);
  }

  /// Changes the properties of a custom key store. You can use this operation to change the properties of an CloudHSM key store or an external key store. Use the required CustomKeyStoreId parameter to i...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateCustomKeyStoreResponse> updateCustomKeyStore(
    $0.UpdateCustomKeyStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateCustomKeyStore, request, options: options);
  }

  /// Updates the description of a KMS key. To see the description of a KMS key, use DescribeKey. The KMS key that you use for this operation must be in a compatible key state. For details, see Key state...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> updateKeyDescription(
    $0.UpdateKeyDescriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateKeyDescription, request, options: options);
  }

  /// Changes the primary key of a multi-Region key. This operation changes the replica key in the specified Region to a primary key and changes the former primary key to a replica key. For example, supp...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> updatePrimaryRegion(
    $0.UpdatePrimaryRegionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updatePrimaryRegion, request, options: options);
  }

  /// Verifies a digital signature that was generated by the Sign operation. Verification confirms that an authorized user signed the message with the specified KMS key and signing algorithm, and the mes...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.VerifyResponse> verify(
    $0.VerifyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$verify, request, options: options);
  }

  /// Verifies the hash-based message authentication code (HMAC) for a specified message, HMAC KMS key, and MAC algorithm. To verify the HMAC, VerifyMac computes an HMAC using the message, HMAC KMS key, ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.VerifyMacResponse> verifyMac(
    $0.VerifyMacRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$verifyMac, request, options: options);
  }

  // method descriptors

  static final _$cancelKeyDeletion = $grpc.ClientMethod<
          $0.CancelKeyDeletionRequest, $0.CancelKeyDeletionResponse>(
      '/kms.KMSService/CancelKeyDeletion',
      ($0.CancelKeyDeletionRequest value) => value.writeToBuffer(),
      $0.CancelKeyDeletionResponse.fromBuffer);
  static final _$connectCustomKeyStore = $grpc.ClientMethod<
          $0.ConnectCustomKeyStoreRequest, $0.ConnectCustomKeyStoreResponse>(
      '/kms.KMSService/ConnectCustomKeyStore',
      ($0.ConnectCustomKeyStoreRequest value) => value.writeToBuffer(),
      $0.ConnectCustomKeyStoreResponse.fromBuffer);
  static final _$createAlias =
      $grpc.ClientMethod<$0.CreateAliasRequest, $1.Empty>(
          '/kms.KMSService/CreateAlias',
          ($0.CreateAliasRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$createCustomKeyStore = $grpc.ClientMethod<
          $0.CreateCustomKeyStoreRequest, $0.CreateCustomKeyStoreResponse>(
      '/kms.KMSService/CreateCustomKeyStore',
      ($0.CreateCustomKeyStoreRequest value) => value.writeToBuffer(),
      $0.CreateCustomKeyStoreResponse.fromBuffer);
  static final _$createGrant =
      $grpc.ClientMethod<$0.CreateGrantRequest, $0.CreateGrantResponse>(
          '/kms.KMSService/CreateGrant',
          ($0.CreateGrantRequest value) => value.writeToBuffer(),
          $0.CreateGrantResponse.fromBuffer);
  static final _$createKey =
      $grpc.ClientMethod<$0.CreateKeyRequest, $0.CreateKeyResponse>(
          '/kms.KMSService/CreateKey',
          ($0.CreateKeyRequest value) => value.writeToBuffer(),
          $0.CreateKeyResponse.fromBuffer);
  static final _$decrypt =
      $grpc.ClientMethod<$0.DecryptRequest, $0.DecryptResponse>(
          '/kms.KMSService/Decrypt',
          ($0.DecryptRequest value) => value.writeToBuffer(),
          $0.DecryptResponse.fromBuffer);
  static final _$deleteAlias =
      $grpc.ClientMethod<$0.DeleteAliasRequest, $1.Empty>(
          '/kms.KMSService/DeleteAlias',
          ($0.DeleteAliasRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteCustomKeyStore = $grpc.ClientMethod<
          $0.DeleteCustomKeyStoreRequest, $0.DeleteCustomKeyStoreResponse>(
      '/kms.KMSService/DeleteCustomKeyStore',
      ($0.DeleteCustomKeyStoreRequest value) => value.writeToBuffer(),
      $0.DeleteCustomKeyStoreResponse.fromBuffer);
  static final _$deleteImportedKeyMaterial = $grpc.ClientMethod<
          $0.DeleteImportedKeyMaterialRequest,
          $0.DeleteImportedKeyMaterialResponse>(
      '/kms.KMSService/DeleteImportedKeyMaterial',
      ($0.DeleteImportedKeyMaterialRequest value) => value.writeToBuffer(),
      $0.DeleteImportedKeyMaterialResponse.fromBuffer);
  static final _$deriveSharedSecret = $grpc.ClientMethod<
          $0.DeriveSharedSecretRequest, $0.DeriveSharedSecretResponse>(
      '/kms.KMSService/DeriveSharedSecret',
      ($0.DeriveSharedSecretRequest value) => value.writeToBuffer(),
      $0.DeriveSharedSecretResponse.fromBuffer);
  static final _$describeCustomKeyStores = $grpc.ClientMethod<
          $0.DescribeCustomKeyStoresRequest,
          $0.DescribeCustomKeyStoresResponse>(
      '/kms.KMSService/DescribeCustomKeyStores',
      ($0.DescribeCustomKeyStoresRequest value) => value.writeToBuffer(),
      $0.DescribeCustomKeyStoresResponse.fromBuffer);
  static final _$describeKey =
      $grpc.ClientMethod<$0.DescribeKeyRequest, $0.DescribeKeyResponse>(
          '/kms.KMSService/DescribeKey',
          ($0.DescribeKeyRequest value) => value.writeToBuffer(),
          $0.DescribeKeyResponse.fromBuffer);
  static final _$disableKey =
      $grpc.ClientMethod<$0.DisableKeyRequest, $1.Empty>(
          '/kms.KMSService/DisableKey',
          ($0.DisableKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$disableKeyRotation =
      $grpc.ClientMethod<$0.DisableKeyRotationRequest, $1.Empty>(
          '/kms.KMSService/DisableKeyRotation',
          ($0.DisableKeyRotationRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$disconnectCustomKeyStore = $grpc.ClientMethod<
          $0.DisconnectCustomKeyStoreRequest,
          $0.DisconnectCustomKeyStoreResponse>(
      '/kms.KMSService/DisconnectCustomKeyStore',
      ($0.DisconnectCustomKeyStoreRequest value) => value.writeToBuffer(),
      $0.DisconnectCustomKeyStoreResponse.fromBuffer);
  static final _$enableKey = $grpc.ClientMethod<$0.EnableKeyRequest, $1.Empty>(
      '/kms.KMSService/EnableKey',
      ($0.EnableKeyRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$enableKeyRotation =
      $grpc.ClientMethod<$0.EnableKeyRotationRequest, $1.Empty>(
          '/kms.KMSService/EnableKeyRotation',
          ($0.EnableKeyRotationRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$encrypt =
      $grpc.ClientMethod<$0.EncryptRequest, $0.EncryptResponse>(
          '/kms.KMSService/Encrypt',
          ($0.EncryptRequest value) => value.writeToBuffer(),
          $0.EncryptResponse.fromBuffer);
  static final _$generateDataKey =
      $grpc.ClientMethod<$0.GenerateDataKeyRequest, $0.GenerateDataKeyResponse>(
          '/kms.KMSService/GenerateDataKey',
          ($0.GenerateDataKeyRequest value) => value.writeToBuffer(),
          $0.GenerateDataKeyResponse.fromBuffer);
  static final _$generateDataKeyPair = $grpc.ClientMethod<
          $0.GenerateDataKeyPairRequest, $0.GenerateDataKeyPairResponse>(
      '/kms.KMSService/GenerateDataKeyPair',
      ($0.GenerateDataKeyPairRequest value) => value.writeToBuffer(),
      $0.GenerateDataKeyPairResponse.fromBuffer);
  static final _$generateDataKeyPairWithoutPlaintext = $grpc.ClientMethod<
          $0.GenerateDataKeyPairWithoutPlaintextRequest,
          $0.GenerateDataKeyPairWithoutPlaintextResponse>(
      '/kms.KMSService/GenerateDataKeyPairWithoutPlaintext',
      ($0.GenerateDataKeyPairWithoutPlaintextRequest value) =>
          value.writeToBuffer(),
      $0.GenerateDataKeyPairWithoutPlaintextResponse.fromBuffer);
  static final _$generateDataKeyWithoutPlaintext = $grpc.ClientMethod<
          $0.GenerateDataKeyWithoutPlaintextRequest,
          $0.GenerateDataKeyWithoutPlaintextResponse>(
      '/kms.KMSService/GenerateDataKeyWithoutPlaintext',
      ($0.GenerateDataKeyWithoutPlaintextRequest value) =>
          value.writeToBuffer(),
      $0.GenerateDataKeyWithoutPlaintextResponse.fromBuffer);
  static final _$generateMac =
      $grpc.ClientMethod<$0.GenerateMacRequest, $0.GenerateMacResponse>(
          '/kms.KMSService/GenerateMac',
          ($0.GenerateMacRequest value) => value.writeToBuffer(),
          $0.GenerateMacResponse.fromBuffer);
  static final _$generateRandom =
      $grpc.ClientMethod<$0.GenerateRandomRequest, $0.GenerateRandomResponse>(
          '/kms.KMSService/GenerateRandom',
          ($0.GenerateRandomRequest value) => value.writeToBuffer(),
          $0.GenerateRandomResponse.fromBuffer);
  static final _$getKeyPolicy =
      $grpc.ClientMethod<$0.GetKeyPolicyRequest, $0.GetKeyPolicyResponse>(
          '/kms.KMSService/GetKeyPolicy',
          ($0.GetKeyPolicyRequest value) => value.writeToBuffer(),
          $0.GetKeyPolicyResponse.fromBuffer);
  static final _$getKeyRotationStatus = $grpc.ClientMethod<
          $0.GetKeyRotationStatusRequest, $0.GetKeyRotationStatusResponse>(
      '/kms.KMSService/GetKeyRotationStatus',
      ($0.GetKeyRotationStatusRequest value) => value.writeToBuffer(),
      $0.GetKeyRotationStatusResponse.fromBuffer);
  static final _$getParametersForImport = $grpc.ClientMethod<
          $0.GetParametersForImportRequest, $0.GetParametersForImportResponse>(
      '/kms.KMSService/GetParametersForImport',
      ($0.GetParametersForImportRequest value) => value.writeToBuffer(),
      $0.GetParametersForImportResponse.fromBuffer);
  static final _$getPublicKey =
      $grpc.ClientMethod<$0.GetPublicKeyRequest, $0.GetPublicKeyResponse>(
          '/kms.KMSService/GetPublicKey',
          ($0.GetPublicKeyRequest value) => value.writeToBuffer(),
          $0.GetPublicKeyResponse.fromBuffer);
  static final _$importKeyMaterial = $grpc.ClientMethod<
          $0.ImportKeyMaterialRequest, $0.ImportKeyMaterialResponse>(
      '/kms.KMSService/ImportKeyMaterial',
      ($0.ImportKeyMaterialRequest value) => value.writeToBuffer(),
      $0.ImportKeyMaterialResponse.fromBuffer);
  static final _$listAliases =
      $grpc.ClientMethod<$0.ListAliasesRequest, $0.ListAliasesResponse>(
          '/kms.KMSService/ListAliases',
          ($0.ListAliasesRequest value) => value.writeToBuffer(),
          $0.ListAliasesResponse.fromBuffer);
  static final _$listGrants =
      $grpc.ClientMethod<$0.ListGrantsRequest, $0.ListGrantsResponse>(
          '/kms.KMSService/ListGrants',
          ($0.ListGrantsRequest value) => value.writeToBuffer(),
          $0.ListGrantsResponse.fromBuffer);
  static final _$listKeyPolicies =
      $grpc.ClientMethod<$0.ListKeyPoliciesRequest, $0.ListKeyPoliciesResponse>(
          '/kms.KMSService/ListKeyPolicies',
          ($0.ListKeyPoliciesRequest value) => value.writeToBuffer(),
          $0.ListKeyPoliciesResponse.fromBuffer);
  static final _$listKeyRotations = $grpc.ClientMethod<
          $0.ListKeyRotationsRequest, $0.ListKeyRotationsResponse>(
      '/kms.KMSService/ListKeyRotations',
      ($0.ListKeyRotationsRequest value) => value.writeToBuffer(),
      $0.ListKeyRotationsResponse.fromBuffer);
  static final _$listKeys =
      $grpc.ClientMethod<$0.ListKeysRequest, $0.ListKeysResponse>(
          '/kms.KMSService/ListKeys',
          ($0.ListKeysRequest value) => value.writeToBuffer(),
          $0.ListKeysResponse.fromBuffer);
  static final _$listResourceTags = $grpc.ClientMethod<
          $0.ListResourceTagsRequest, $0.ListResourceTagsResponse>(
      '/kms.KMSService/ListResourceTags',
      ($0.ListResourceTagsRequest value) => value.writeToBuffer(),
      $0.ListResourceTagsResponse.fromBuffer);
  static final _$listRetirableGrants =
      $grpc.ClientMethod<$0.ListRetirableGrantsRequest, $0.ListGrantsResponse>(
          '/kms.KMSService/ListRetirableGrants',
          ($0.ListRetirableGrantsRequest value) => value.writeToBuffer(),
          $0.ListGrantsResponse.fromBuffer);
  static final _$putKeyPolicy =
      $grpc.ClientMethod<$0.PutKeyPolicyRequest, $1.Empty>(
          '/kms.KMSService/PutKeyPolicy',
          ($0.PutKeyPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$reEncrypt =
      $grpc.ClientMethod<$0.ReEncryptRequest, $0.ReEncryptResponse>(
          '/kms.KMSService/ReEncrypt',
          ($0.ReEncryptRequest value) => value.writeToBuffer(),
          $0.ReEncryptResponse.fromBuffer);
  static final _$replicateKey =
      $grpc.ClientMethod<$0.ReplicateKeyRequest, $0.ReplicateKeyResponse>(
          '/kms.KMSService/ReplicateKey',
          ($0.ReplicateKeyRequest value) => value.writeToBuffer(),
          $0.ReplicateKeyResponse.fromBuffer);
  static final _$retireGrant =
      $grpc.ClientMethod<$0.RetireGrantRequest, $1.Empty>(
          '/kms.KMSService/RetireGrant',
          ($0.RetireGrantRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$revokeGrant =
      $grpc.ClientMethod<$0.RevokeGrantRequest, $1.Empty>(
          '/kms.KMSService/RevokeGrant',
          ($0.RevokeGrantRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$rotateKeyOnDemand = $grpc.ClientMethod<
          $0.RotateKeyOnDemandRequest, $0.RotateKeyOnDemandResponse>(
      '/kms.KMSService/RotateKeyOnDemand',
      ($0.RotateKeyOnDemandRequest value) => value.writeToBuffer(),
      $0.RotateKeyOnDemandResponse.fromBuffer);
  static final _$scheduleKeyDeletion = $grpc.ClientMethod<
          $0.ScheduleKeyDeletionRequest, $0.ScheduleKeyDeletionResponse>(
      '/kms.KMSService/ScheduleKeyDeletion',
      ($0.ScheduleKeyDeletionRequest value) => value.writeToBuffer(),
      $0.ScheduleKeyDeletionResponse.fromBuffer);
  static final _$sign = $grpc.ClientMethod<$0.SignRequest, $0.SignResponse>(
      '/kms.KMSService/Sign',
      ($0.SignRequest value) => value.writeToBuffer(),
      $0.SignResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $1.Empty>(
          '/kms.KMSService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $1.Empty>(
          '/kms.KMSService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateAlias =
      $grpc.ClientMethod<$0.UpdateAliasRequest, $1.Empty>(
          '/kms.KMSService/UpdateAlias',
          ($0.UpdateAliasRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateCustomKeyStore = $grpc.ClientMethod<
          $0.UpdateCustomKeyStoreRequest, $0.UpdateCustomKeyStoreResponse>(
      '/kms.KMSService/UpdateCustomKeyStore',
      ($0.UpdateCustomKeyStoreRequest value) => value.writeToBuffer(),
      $0.UpdateCustomKeyStoreResponse.fromBuffer);
  static final _$updateKeyDescription =
      $grpc.ClientMethod<$0.UpdateKeyDescriptionRequest, $1.Empty>(
          '/kms.KMSService/UpdateKeyDescription',
          ($0.UpdateKeyDescriptionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updatePrimaryRegion =
      $grpc.ClientMethod<$0.UpdatePrimaryRegionRequest, $1.Empty>(
          '/kms.KMSService/UpdatePrimaryRegion',
          ($0.UpdatePrimaryRegionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$verify =
      $grpc.ClientMethod<$0.VerifyRequest, $0.VerifyResponse>(
          '/kms.KMSService/Verify',
          ($0.VerifyRequest value) => value.writeToBuffer(),
          $0.VerifyResponse.fromBuffer);
  static final _$verifyMac =
      $grpc.ClientMethod<$0.VerifyMacRequest, $0.VerifyMacResponse>(
          '/kms.KMSService/VerifyMac',
          ($0.VerifyMacRequest value) => value.writeToBuffer(),
          $0.VerifyMacResponse.fromBuffer);
}

@$pb.GrpcServiceName('kms.KMSService')
abstract class KMSServiceBase extends $grpc.Service {
  $core.String get $name => 'kms.KMSService';

  KMSServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CancelKeyDeletionRequest,
            $0.CancelKeyDeletionResponse>(
        'CancelKeyDeletion',
        cancelKeyDeletion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelKeyDeletionRequest.fromBuffer(value),
        ($0.CancelKeyDeletionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ConnectCustomKeyStoreRequest,
            $0.ConnectCustomKeyStoreResponse>(
        'ConnectCustomKeyStore',
        connectCustomKeyStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ConnectCustomKeyStoreRequest.fromBuffer(value),
        ($0.ConnectCustomKeyStoreResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateAliasRequest, $1.Empty>(
        'CreateAlias',
        createAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateAliasRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateCustomKeyStoreRequest,
            $0.CreateCustomKeyStoreResponse>(
        'CreateCustomKeyStore',
        createCustomKeyStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCustomKeyStoreRequest.fromBuffer(value),
        ($0.CreateCustomKeyStoreResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateGrantRequest, $0.CreateGrantResponse>(
            'CreateGrant',
            createGrant_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateGrantRequest.fromBuffer(value),
            ($0.CreateGrantResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateKeyRequest, $0.CreateKeyResponse>(
        'CreateKey',
        createKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateKeyRequest.fromBuffer(value),
        ($0.CreateKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DecryptRequest, $0.DecryptResponse>(
        'Decrypt',
        decrypt_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DecryptRequest.fromBuffer(value),
        ($0.DecryptResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAliasRequest, $1.Empty>(
        'DeleteAlias',
        deleteAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAliasRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteCustomKeyStoreRequest,
            $0.DeleteCustomKeyStoreResponse>(
        'DeleteCustomKeyStore',
        deleteCustomKeyStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteCustomKeyStoreRequest.fromBuffer(value),
        ($0.DeleteCustomKeyStoreResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteImportedKeyMaterialRequest,
            $0.DeleteImportedKeyMaterialResponse>(
        'DeleteImportedKeyMaterial',
        deleteImportedKeyMaterial_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteImportedKeyMaterialRequest.fromBuffer(value),
        ($0.DeleteImportedKeyMaterialResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeriveSharedSecretRequest,
            $0.DeriveSharedSecretResponse>(
        'DeriveSharedSecret',
        deriveSharedSecret_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeriveSharedSecretRequest.fromBuffer(value),
        ($0.DeriveSharedSecretResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeCustomKeyStoresRequest,
            $0.DescribeCustomKeyStoresResponse>(
        'DescribeCustomKeyStores',
        describeCustomKeyStores_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeCustomKeyStoresRequest.fromBuffer(value),
        ($0.DescribeCustomKeyStoresResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeKeyRequest, $0.DescribeKeyResponse>(
            'DescribeKey',
            describeKey_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeKeyRequest.fromBuffer(value),
            ($0.DescribeKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableKeyRequest, $1.Empty>(
        'DisableKey',
        disableKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DisableKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableKeyRotationRequest, $1.Empty>(
        'DisableKeyRotation',
        disableKeyRotation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableKeyRotationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisconnectCustomKeyStoreRequest,
            $0.DisconnectCustomKeyStoreResponse>(
        'DisconnectCustomKeyStore',
        disconnectCustomKeyStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisconnectCustomKeyStoreRequest.fromBuffer(value),
        ($0.DisconnectCustomKeyStoreResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableKeyRequest, $1.Empty>(
        'EnableKey',
        enableKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.EnableKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableKeyRotationRequest, $1.Empty>(
        'EnableKeyRotation',
        enableKeyRotation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableKeyRotationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EncryptRequest, $0.EncryptResponse>(
        'Encrypt',
        encrypt_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.EncryptRequest.fromBuffer(value),
        ($0.EncryptResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateDataKeyRequest,
            $0.GenerateDataKeyResponse>(
        'GenerateDataKey',
        generateDataKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateDataKeyRequest.fromBuffer(value),
        ($0.GenerateDataKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateDataKeyPairRequest,
            $0.GenerateDataKeyPairResponse>(
        'GenerateDataKeyPair',
        generateDataKeyPair_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateDataKeyPairRequest.fromBuffer(value),
        ($0.GenerateDataKeyPairResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GenerateDataKeyPairWithoutPlaintextRequest,
            $0.GenerateDataKeyPairWithoutPlaintextResponse>(
        'GenerateDataKeyPairWithoutPlaintext',
        generateDataKeyPairWithoutPlaintext_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateDataKeyPairWithoutPlaintextRequest.fromBuffer(value),
        ($0.GenerateDataKeyPairWithoutPlaintextResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateDataKeyWithoutPlaintextRequest,
            $0.GenerateDataKeyWithoutPlaintextResponse>(
        'GenerateDataKeyWithoutPlaintext',
        generateDataKeyWithoutPlaintext_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateDataKeyWithoutPlaintextRequest.fromBuffer(value),
        ($0.GenerateDataKeyWithoutPlaintextResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GenerateMacRequest, $0.GenerateMacResponse>(
            'GenerateMac',
            generateMac_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GenerateMacRequest.fromBuffer(value),
            ($0.GenerateMacResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateRandomRequest,
            $0.GenerateRandomResponse>(
        'GenerateRandom',
        generateRandom_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateRandomRequest.fromBuffer(value),
        ($0.GenerateRandomResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetKeyPolicyRequest, $0.GetKeyPolicyResponse>(
            'GetKeyPolicy',
            getKeyPolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetKeyPolicyRequest.fromBuffer(value),
            ($0.GetKeyPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetKeyRotationStatusRequest,
            $0.GetKeyRotationStatusResponse>(
        'GetKeyRotationStatus',
        getKeyRotationStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetKeyRotationStatusRequest.fromBuffer(value),
        ($0.GetKeyRotationStatusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetParametersForImportRequest,
            $0.GetParametersForImportResponse>(
        'GetParametersForImport',
        getParametersForImport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetParametersForImportRequest.fromBuffer(value),
        ($0.GetParametersForImportResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetPublicKeyRequest, $0.GetPublicKeyResponse>(
            'GetPublicKey',
            getPublicKey_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetPublicKeyRequest.fromBuffer(value),
            ($0.GetPublicKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ImportKeyMaterialRequest,
            $0.ImportKeyMaterialResponse>(
        'ImportKeyMaterial',
        importKeyMaterial_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ImportKeyMaterialRequest.fromBuffer(value),
        ($0.ImportKeyMaterialResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListAliasesRequest, $0.ListAliasesResponse>(
            'ListAliases',
            listAliases_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListAliasesRequest.fromBuffer(value),
            ($0.ListAliasesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListGrantsRequest, $0.ListGrantsResponse>(
        'ListGrants',
        listGrants_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListGrantsRequest.fromBuffer(value),
        ($0.ListGrantsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListKeyPoliciesRequest,
            $0.ListKeyPoliciesResponse>(
        'ListKeyPolicies',
        listKeyPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListKeyPoliciesRequest.fromBuffer(value),
        ($0.ListKeyPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListKeyRotationsRequest,
            $0.ListKeyRotationsResponse>(
        'ListKeyRotations',
        listKeyRotations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListKeyRotationsRequest.fromBuffer(value),
        ($0.ListKeyRotationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListKeysRequest, $0.ListKeysResponse>(
        'ListKeys',
        listKeys_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListKeysRequest.fromBuffer(value),
        ($0.ListKeysResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListResourceTagsRequest,
            $0.ListResourceTagsResponse>(
        'ListResourceTags',
        listResourceTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListResourceTagsRequest.fromBuffer(value),
        ($0.ListResourceTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRetirableGrantsRequest,
            $0.ListGrantsResponse>(
        'ListRetirableGrants',
        listRetirableGrants_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRetirableGrantsRequest.fromBuffer(value),
        ($0.ListGrantsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutKeyPolicyRequest, $1.Empty>(
        'PutKeyPolicy',
        putKeyPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutKeyPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ReEncryptRequest, $0.ReEncryptResponse>(
        'ReEncrypt',
        reEncrypt_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ReEncryptRequest.fromBuffer(value),
        ($0.ReEncryptResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ReplicateKeyRequest, $0.ReplicateKeyResponse>(
            'ReplicateKey',
            replicateKey_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ReplicateKeyRequest.fromBuffer(value),
            ($0.ReplicateKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RetireGrantRequest, $1.Empty>(
        'RetireGrant',
        retireGrant_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RetireGrantRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RevokeGrantRequest, $1.Empty>(
        'RevokeGrant',
        revokeGrant_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RevokeGrantRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RotateKeyOnDemandRequest,
            $0.RotateKeyOnDemandResponse>(
        'RotateKeyOnDemand',
        rotateKeyOnDemand_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RotateKeyOnDemandRequest.fromBuffer(value),
        ($0.RotateKeyOnDemandResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ScheduleKeyDeletionRequest,
            $0.ScheduleKeyDeletionResponse>(
        'ScheduleKeyDeletion',
        scheduleKeyDeletion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ScheduleKeyDeletionRequest.fromBuffer(value),
        ($0.ScheduleKeyDeletionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SignRequest, $0.SignResponse>(
        'Sign',
        sign_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SignRequest.fromBuffer(value),
        ($0.SignResponse value) => value.writeToBuffer()));
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
    $addMethod($grpc.ServiceMethod<$0.UpdateAliasRequest, $1.Empty>(
        'UpdateAlias',
        updateAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAliasRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateCustomKeyStoreRequest,
            $0.UpdateCustomKeyStoreResponse>(
        'UpdateCustomKeyStore',
        updateCustomKeyStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateCustomKeyStoreRequest.fromBuffer(value),
        ($0.UpdateCustomKeyStoreResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateKeyDescriptionRequest, $1.Empty>(
        'UpdateKeyDescription',
        updateKeyDescription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateKeyDescriptionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdatePrimaryRegionRequest, $1.Empty>(
        'UpdatePrimaryRegion',
        updatePrimaryRegion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdatePrimaryRegionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.VerifyRequest, $0.VerifyResponse>(
        'Verify',
        verify_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.VerifyRequest.fromBuffer(value),
        ($0.VerifyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.VerifyMacRequest, $0.VerifyMacResponse>(
        'VerifyMac',
        verifyMac_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.VerifyMacRequest.fromBuffer(value),
        ($0.VerifyMacResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.CancelKeyDeletionResponse> cancelKeyDeletion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelKeyDeletionRequest> $request) async {
    return cancelKeyDeletion($call, await $request);
  }

  $async.Future<$0.CancelKeyDeletionResponse> cancelKeyDeletion(
      $grpc.ServiceCall call, $0.CancelKeyDeletionRequest request);

  $async.Future<$0.ConnectCustomKeyStoreResponse> connectCustomKeyStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ConnectCustomKeyStoreRequest> $request) async {
    return connectCustomKeyStore($call, await $request);
  }

  $async.Future<$0.ConnectCustomKeyStoreResponse> connectCustomKeyStore(
      $grpc.ServiceCall call, $0.ConnectCustomKeyStoreRequest request);

  $async.Future<$1.Empty> createAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateAliasRequest> $request) async {
    return createAlias($call, await $request);
  }

  $async.Future<$1.Empty> createAlias(
      $grpc.ServiceCall call, $0.CreateAliasRequest request);

  $async.Future<$0.CreateCustomKeyStoreResponse> createCustomKeyStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateCustomKeyStoreRequest> $request) async {
    return createCustomKeyStore($call, await $request);
  }

  $async.Future<$0.CreateCustomKeyStoreResponse> createCustomKeyStore(
      $grpc.ServiceCall call, $0.CreateCustomKeyStoreRequest request);

  $async.Future<$0.CreateGrantResponse> createGrant_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateGrantRequest> $request) async {
    return createGrant($call, await $request);
  }

  $async.Future<$0.CreateGrantResponse> createGrant(
      $grpc.ServiceCall call, $0.CreateGrantRequest request);

  $async.Future<$0.CreateKeyResponse> createKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateKeyRequest> $request) async {
    return createKey($call, await $request);
  }

  $async.Future<$0.CreateKeyResponse> createKey(
      $grpc.ServiceCall call, $0.CreateKeyRequest request);

  $async.Future<$0.DecryptResponse> decrypt_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DecryptRequest> $request) async {
    return decrypt($call, await $request);
  }

  $async.Future<$0.DecryptResponse> decrypt(
      $grpc.ServiceCall call, $0.DecryptRequest request);

  $async.Future<$1.Empty> deleteAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAliasRequest> $request) async {
    return deleteAlias($call, await $request);
  }

  $async.Future<$1.Empty> deleteAlias(
      $grpc.ServiceCall call, $0.DeleteAliasRequest request);

  $async.Future<$0.DeleteCustomKeyStoreResponse> deleteCustomKeyStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteCustomKeyStoreRequest> $request) async {
    return deleteCustomKeyStore($call, await $request);
  }

  $async.Future<$0.DeleteCustomKeyStoreResponse> deleteCustomKeyStore(
      $grpc.ServiceCall call, $0.DeleteCustomKeyStoreRequest request);

  $async.Future<$0.DeleteImportedKeyMaterialResponse>
      deleteImportedKeyMaterial_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteImportedKeyMaterialRequest> $request) async {
    return deleteImportedKeyMaterial($call, await $request);
  }

  $async.Future<$0.DeleteImportedKeyMaterialResponse> deleteImportedKeyMaterial(
      $grpc.ServiceCall call, $0.DeleteImportedKeyMaterialRequest request);

  $async.Future<$0.DeriveSharedSecretResponse> deriveSharedSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeriveSharedSecretRequest> $request) async {
    return deriveSharedSecret($call, await $request);
  }

  $async.Future<$0.DeriveSharedSecretResponse> deriveSharedSecret(
      $grpc.ServiceCall call, $0.DeriveSharedSecretRequest request);

  $async.Future<$0.DescribeCustomKeyStoresResponse> describeCustomKeyStores_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeCustomKeyStoresRequest> $request) async {
    return describeCustomKeyStores($call, await $request);
  }

  $async.Future<$0.DescribeCustomKeyStoresResponse> describeCustomKeyStores(
      $grpc.ServiceCall call, $0.DescribeCustomKeyStoresRequest request);

  $async.Future<$0.DescribeKeyResponse> describeKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DescribeKeyRequest> $request) async {
    return describeKey($call, await $request);
  }

  $async.Future<$0.DescribeKeyResponse> describeKey(
      $grpc.ServiceCall call, $0.DescribeKeyRequest request);

  $async.Future<$1.Empty> disableKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DisableKeyRequest> $request) async {
    return disableKey($call, await $request);
  }

  $async.Future<$1.Empty> disableKey(
      $grpc.ServiceCall call, $0.DisableKeyRequest request);

  $async.Future<$1.Empty> disableKeyRotation_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DisableKeyRotationRequest> $request) async {
    return disableKeyRotation($call, await $request);
  }

  $async.Future<$1.Empty> disableKeyRotation(
      $grpc.ServiceCall call, $0.DisableKeyRotationRequest request);

  $async.Future<$0.DisconnectCustomKeyStoreResponse>
      disconnectCustomKeyStore_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DisconnectCustomKeyStoreRequest> $request) async {
    return disconnectCustomKeyStore($call, await $request);
  }

  $async.Future<$0.DisconnectCustomKeyStoreResponse> disconnectCustomKeyStore(
      $grpc.ServiceCall call, $0.DisconnectCustomKeyStoreRequest request);

  $async.Future<$1.Empty> enableKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EnableKeyRequest> $request) async {
    return enableKey($call, await $request);
  }

  $async.Future<$1.Empty> enableKey(
      $grpc.ServiceCall call, $0.EnableKeyRequest request);

  $async.Future<$1.Empty> enableKeyRotation_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EnableKeyRotationRequest> $request) async {
    return enableKeyRotation($call, await $request);
  }

  $async.Future<$1.Empty> enableKeyRotation(
      $grpc.ServiceCall call, $0.EnableKeyRotationRequest request);

  $async.Future<$0.EncryptResponse> encrypt_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EncryptRequest> $request) async {
    return encrypt($call, await $request);
  }

  $async.Future<$0.EncryptResponse> encrypt(
      $grpc.ServiceCall call, $0.EncryptRequest request);

  $async.Future<$0.GenerateDataKeyResponse> generateDataKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GenerateDataKeyRequest> $request) async {
    return generateDataKey($call, await $request);
  }

  $async.Future<$0.GenerateDataKeyResponse> generateDataKey(
      $grpc.ServiceCall call, $0.GenerateDataKeyRequest request);

  $async.Future<$0.GenerateDataKeyPairResponse> generateDataKeyPair_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GenerateDataKeyPairRequest> $request) async {
    return generateDataKeyPair($call, await $request);
  }

  $async.Future<$0.GenerateDataKeyPairResponse> generateDataKeyPair(
      $grpc.ServiceCall call, $0.GenerateDataKeyPairRequest request);

  $async.Future<$0.GenerateDataKeyPairWithoutPlaintextResponse>
      generateDataKeyPairWithoutPlaintext_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GenerateDataKeyPairWithoutPlaintextRequest>
              $request) async {
    return generateDataKeyPairWithoutPlaintext($call, await $request);
  }

  $async.Future<$0.GenerateDataKeyPairWithoutPlaintextResponse>
      generateDataKeyPairWithoutPlaintext($grpc.ServiceCall call,
          $0.GenerateDataKeyPairWithoutPlaintextRequest request);

  $async.Future<$0.GenerateDataKeyWithoutPlaintextResponse>
      generateDataKeyWithoutPlaintext_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GenerateDataKeyWithoutPlaintextRequest>
              $request) async {
    return generateDataKeyWithoutPlaintext($call, await $request);
  }

  $async.Future<$0.GenerateDataKeyWithoutPlaintextResponse>
      generateDataKeyWithoutPlaintext($grpc.ServiceCall call,
          $0.GenerateDataKeyWithoutPlaintextRequest request);

  $async.Future<$0.GenerateMacResponse> generateMac_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GenerateMacRequest> $request) async {
    return generateMac($call, await $request);
  }

  $async.Future<$0.GenerateMacResponse> generateMac(
      $grpc.ServiceCall call, $0.GenerateMacRequest request);

  $async.Future<$0.GenerateRandomResponse> generateRandom_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GenerateRandomRequest> $request) async {
    return generateRandom($call, await $request);
  }

  $async.Future<$0.GenerateRandomResponse> generateRandom(
      $grpc.ServiceCall call, $0.GenerateRandomRequest request);

  $async.Future<$0.GetKeyPolicyResponse> getKeyPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetKeyPolicyRequest> $request) async {
    return getKeyPolicy($call, await $request);
  }

  $async.Future<$0.GetKeyPolicyResponse> getKeyPolicy(
      $grpc.ServiceCall call, $0.GetKeyPolicyRequest request);

  $async.Future<$0.GetKeyRotationStatusResponse> getKeyRotationStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetKeyRotationStatusRequest> $request) async {
    return getKeyRotationStatus($call, await $request);
  }

  $async.Future<$0.GetKeyRotationStatusResponse> getKeyRotationStatus(
      $grpc.ServiceCall call, $0.GetKeyRotationStatusRequest request);

  $async.Future<$0.GetParametersForImportResponse> getParametersForImport_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetParametersForImportRequest> $request) async {
    return getParametersForImport($call, await $request);
  }

  $async.Future<$0.GetParametersForImportResponse> getParametersForImport(
      $grpc.ServiceCall call, $0.GetParametersForImportRequest request);

  $async.Future<$0.GetPublicKeyResponse> getPublicKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPublicKeyRequest> $request) async {
    return getPublicKey($call, await $request);
  }

  $async.Future<$0.GetPublicKeyResponse> getPublicKey(
      $grpc.ServiceCall call, $0.GetPublicKeyRequest request);

  $async.Future<$0.ImportKeyMaterialResponse> importKeyMaterial_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ImportKeyMaterialRequest> $request) async {
    return importKeyMaterial($call, await $request);
  }

  $async.Future<$0.ImportKeyMaterialResponse> importKeyMaterial(
      $grpc.ServiceCall call, $0.ImportKeyMaterialRequest request);

  $async.Future<$0.ListAliasesResponse> listAliases_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListAliasesRequest> $request) async {
    return listAliases($call, await $request);
  }

  $async.Future<$0.ListAliasesResponse> listAliases(
      $grpc.ServiceCall call, $0.ListAliasesRequest request);

  $async.Future<$0.ListGrantsResponse> listGrants_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListGrantsRequest> $request) async {
    return listGrants($call, await $request);
  }

  $async.Future<$0.ListGrantsResponse> listGrants(
      $grpc.ServiceCall call, $0.ListGrantsRequest request);

  $async.Future<$0.ListKeyPoliciesResponse> listKeyPolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListKeyPoliciesRequest> $request) async {
    return listKeyPolicies($call, await $request);
  }

  $async.Future<$0.ListKeyPoliciesResponse> listKeyPolicies(
      $grpc.ServiceCall call, $0.ListKeyPoliciesRequest request);

  $async.Future<$0.ListKeyRotationsResponse> listKeyRotations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListKeyRotationsRequest> $request) async {
    return listKeyRotations($call, await $request);
  }

  $async.Future<$0.ListKeyRotationsResponse> listKeyRotations(
      $grpc.ServiceCall call, $0.ListKeyRotationsRequest request);

  $async.Future<$0.ListKeysResponse> listKeys_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListKeysRequest> $request) async {
    return listKeys($call, await $request);
  }

  $async.Future<$0.ListKeysResponse> listKeys(
      $grpc.ServiceCall call, $0.ListKeysRequest request);

  $async.Future<$0.ListResourceTagsResponse> listResourceTags_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListResourceTagsRequest> $request) async {
    return listResourceTags($call, await $request);
  }

  $async.Future<$0.ListResourceTagsResponse> listResourceTags(
      $grpc.ServiceCall call, $0.ListResourceTagsRequest request);

  $async.Future<$0.ListGrantsResponse> listRetirableGrants_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRetirableGrantsRequest> $request) async {
    return listRetirableGrants($call, await $request);
  }

  $async.Future<$0.ListGrantsResponse> listRetirableGrants(
      $grpc.ServiceCall call, $0.ListRetirableGrantsRequest request);

  $async.Future<$1.Empty> putKeyPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutKeyPolicyRequest> $request) async {
    return putKeyPolicy($call, await $request);
  }

  $async.Future<$1.Empty> putKeyPolicy(
      $grpc.ServiceCall call, $0.PutKeyPolicyRequest request);

  $async.Future<$0.ReEncryptResponse> reEncrypt_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ReEncryptRequest> $request) async {
    return reEncrypt($call, await $request);
  }

  $async.Future<$0.ReEncryptResponse> reEncrypt(
      $grpc.ServiceCall call, $0.ReEncryptRequest request);

  $async.Future<$0.ReplicateKeyResponse> replicateKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ReplicateKeyRequest> $request) async {
    return replicateKey($call, await $request);
  }

  $async.Future<$0.ReplicateKeyResponse> replicateKey(
      $grpc.ServiceCall call, $0.ReplicateKeyRequest request);

  $async.Future<$1.Empty> retireGrant_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RetireGrantRequest> $request) async {
    return retireGrant($call, await $request);
  }

  $async.Future<$1.Empty> retireGrant(
      $grpc.ServiceCall call, $0.RetireGrantRequest request);

  $async.Future<$1.Empty> revokeGrant_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RevokeGrantRequest> $request) async {
    return revokeGrant($call, await $request);
  }

  $async.Future<$1.Empty> revokeGrant(
      $grpc.ServiceCall call, $0.RevokeGrantRequest request);

  $async.Future<$0.RotateKeyOnDemandResponse> rotateKeyOnDemand_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RotateKeyOnDemandRequest> $request) async {
    return rotateKeyOnDemand($call, await $request);
  }

  $async.Future<$0.RotateKeyOnDemandResponse> rotateKeyOnDemand(
      $grpc.ServiceCall call, $0.RotateKeyOnDemandRequest request);

  $async.Future<$0.ScheduleKeyDeletionResponse> scheduleKeyDeletion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ScheduleKeyDeletionRequest> $request) async {
    return scheduleKeyDeletion($call, await $request);
  }

  $async.Future<$0.ScheduleKeyDeletionResponse> scheduleKeyDeletion(
      $grpc.ServiceCall call, $0.ScheduleKeyDeletionRequest request);

  $async.Future<$0.SignResponse> sign_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.SignRequest> $request) async {
    return sign($call, await $request);
  }

  $async.Future<$0.SignResponse> sign(
      $grpc.ServiceCall call, $0.SignRequest request);

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

  $async.Future<$1.Empty> updateAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateAliasRequest> $request) async {
    return updateAlias($call, await $request);
  }

  $async.Future<$1.Empty> updateAlias(
      $grpc.ServiceCall call, $0.UpdateAliasRequest request);

  $async.Future<$0.UpdateCustomKeyStoreResponse> updateCustomKeyStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateCustomKeyStoreRequest> $request) async {
    return updateCustomKeyStore($call, await $request);
  }

  $async.Future<$0.UpdateCustomKeyStoreResponse> updateCustomKeyStore(
      $grpc.ServiceCall call, $0.UpdateCustomKeyStoreRequest request);

  $async.Future<$1.Empty> updateKeyDescription_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateKeyDescriptionRequest> $request) async {
    return updateKeyDescription($call, await $request);
  }

  $async.Future<$1.Empty> updateKeyDescription(
      $grpc.ServiceCall call, $0.UpdateKeyDescriptionRequest request);

  $async.Future<$1.Empty> updatePrimaryRegion_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdatePrimaryRegionRequest> $request) async {
    return updatePrimaryRegion($call, await $request);
  }

  $async.Future<$1.Empty> updatePrimaryRegion(
      $grpc.ServiceCall call, $0.UpdatePrimaryRegionRequest request);

  $async.Future<$0.VerifyResponse> verify_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.VerifyRequest> $request) async {
    return verify($call, await $request);
  }

  $async.Future<$0.VerifyResponse> verify(
      $grpc.ServiceCall call, $0.VerifyRequest request);

  $async.Future<$0.VerifyMacResponse> verifyMac_Pre($grpc.ServiceCall $call,
      $async.Future<$0.VerifyMacRequest> $request) async {
    return verifyMac($call, await $request);
  }

  $async.Future<$0.VerifyMacResponse> verifyMac(
      $grpc.ServiceCall call, $0.VerifyMacRequest request);
}
