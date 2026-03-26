package kms

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/kms"
	"vorpalstacks/internal/pb/aws/kms/kmsconnect"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// AdminHandler provides administrative operations for KMS keys.
// It implements the KMSServiceHandler interface for gRPC-Web communication.
type AdminHandler struct {
	store *kmsstore.KeyStore
}

var _ kmsconnect.KMSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new AdminHandler with the given key store.
// The store is used for listing KMS keys.
func NewAdminHandler(store *kmsstore.KeyStore) *AdminHandler {
	return &AdminHandler{store: store}
}

// ListKeys retrieves a list of KMS keys from the key store.
// It supports pagination via the Marker and Limit parameters.
// Returns a list of KeyListEntry objects containing KeyId and KeyArn.
func (h *AdminHandler) ListKeys(ctx context.Context, req *connect.Request[pb.ListKeysRequest]) (*connect.Response[pb.ListKeysResponse], error) {
	result, err := h.store.List(req.Msg.Marker, int(req.Msg.Limit))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	keys := make([]*pb.KeyListEntry, len(result.Keys))
	for i, k := range result.Keys {
		keys[i] = &pb.KeyListEntry{
			Keyid:  k.KeyID,
			Keyarn: k.KeyArn,
		}
	}

	return connect.NewResponse(&pb.ListKeysResponse{
		Keys:       keys,
		Nextmarker: result.NextMarker,
		Truncated:  result.IsTruncated,
	}), nil
}

// CancelKeyDeletion cancels the scheduled deletion of a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelKeyDeletion(ctx context.Context, req *connect.Request[pb.CancelKeyDeletionRequest]) (*connect.Response[pb.CancelKeyDeletionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ConnectCustomKeyStore connects to a custom key store in AWS CloudHSM.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ConnectCustomKeyStore(ctx context.Context, req *connect.Request[pb.ConnectCustomKeyStoreRequest]) (*connect.Response[pb.ConnectCustomKeyStoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateAlias creates a display name (alias) for a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateAlias(ctx context.Context, req *connect.Request[pb.CreateAliasRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateCustomKeyStore creates a custom key store backed by AWS CloudHSM or an external key manager.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateCustomKeyStore(ctx context.Context, req *connect.Request[pb.CreateCustomKeyStoreRequest]) (*connect.Response[pb.CreateCustomKeyStoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateGrant adds a grant to a KMS key, permitting the grantee principal to use the key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateGrant(ctx context.Context, req *connect.Request[pb.CreateGrantRequest]) (*connect.Response[pb.CreateGrantResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// CreateKey creates a new KMS key (customer master key) in the key store.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateKey(ctx context.Context, req *connect.Request[pb.CreateKeyRequest]) (*connect.Response[pb.CreateKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Decrypt decrypts ciphertext using a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Decrypt(ctx context.Context, req *connect.Request[pb.DecryptRequest]) (*connect.Response[pb.DecryptResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteAlias removes an alias from a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAlias(ctx context.Context, req *connect.Request[pb.DeleteAliasRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteCustomKeyStore deletes a custom key store.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteCustomKeyStore(ctx context.Context, req *connect.Request[pb.DeleteCustomKeyStoreRequest]) (*connect.Response[pb.DeleteCustomKeyStoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeleteImportedKeyMaterial deletes imported key material from a KMS key.
// After this operation, the key can no longer be used to decrypt data.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteImportedKeyMaterial(ctx context.Context, req *connect.Request[pb.DeleteImportedKeyMaterialRequest]) (*connect.Response[pb.DeleteImportedKeyMaterialResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DeriveSharedSecret derives a shared secret using ECDH key exchange with a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeriveSharedSecret(ctx context.Context, req *connect.Request[pb.DeriveSharedSecretRequest]) (*connect.Response[pb.DeriveSharedSecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeCustomKeyStores returns information about one or more custom key stores.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeCustomKeyStores(ctx context.Context, req *connect.Request[pb.DescribeCustomKeyStoresRequest]) (*connect.Response[pb.DescribeCustomKeyStoresResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DescribeKey returns detailed information about a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeKey(ctx context.Context, req *connect.Request[pb.DescribeKeyRequest]) (*connect.Response[pb.DescribeKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisableKey disables a KMS key, preventing its use for cryptographic operations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisableKey(ctx context.Context, req *connect.Request[pb.DisableKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisableKeyRotation disables automatic key rotation for a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisableKeyRotation(ctx context.Context, req *connect.Request[pb.DisableKeyRotationRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// DisconnectCustomKeyStore disconnects a custom key store from its underlying key store.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisconnectCustomKeyStore(ctx context.Context, req *connect.Request[pb.DisconnectCustomKeyStoreRequest]) (*connect.Response[pb.DisconnectCustomKeyStoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// EnableKey enables a KMS key, permitting its use for cryptographic operations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableKey(ctx context.Context, req *connect.Request[pb.EnableKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// EnableKeyRotation enables automatic annual key rotation for a symmetric KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableKeyRotation(ctx context.Context, req *connect.Request[pb.EnableKeyRotationRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Encrypt encrypts plaintext data using a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Encrypt(ctx context.Context, req *connect.Request[pb.EncryptRequest]) (*connect.Response[pb.EncryptResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateDataKey generates a unique symmetric data key for client-side encryption.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateDataKey(ctx context.Context, req *connect.Request[pb.GenerateDataKeyRequest]) (*connect.Response[pb.GenerateDataKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateDataKeyPair generates a public/private key pair for asymmetric encryption.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateDataKeyPair(ctx context.Context, req *connect.Request[pb.GenerateDataKeyPairRequest]) (*connect.Response[pb.GenerateDataKeyPairResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateDataKeyPairWithoutPlaintext generates a public/private key pair without the plaintext private key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateDataKeyPairWithoutPlaintext(ctx context.Context, req *connect.Request[pb.GenerateDataKeyPairWithoutPlaintextRequest]) (*connect.Response[pb.GenerateDataKeyPairWithoutPlaintextResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateDataKeyWithoutPlaintext generates a data key without the plaintext copy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateDataKeyWithoutPlaintext(ctx context.Context, req *connect.Request[pb.GenerateDataKeyWithoutPlaintextRequest]) (*connect.Response[pb.GenerateDataKeyWithoutPlaintextResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateMac generates a message authentication code (MAC) for a message using a symmetric KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateMac(ctx context.Context, req *connect.Request[pb.GenerateMacRequest]) (*connect.Response[pb.GenerateMacResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GenerateRandom returns a random byte string of the specified length.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateRandom(ctx context.Context, req *connect.Request[pb.GenerateRandomRequest]) (*connect.Response[pb.GenerateRandomResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetKeyPolicy retrieves the key policy attached to a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetKeyPolicy(ctx context.Context, req *connect.Request[pb.GetKeyPolicyRequest]) (*connect.Response[pb.GetKeyPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetKeyRotationStatus returns the status of automatic key rotation for a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetKeyRotationStatus(ctx context.Context, req *connect.Request[pb.GetKeyRotationStatusRequest]) (*connect.Response[pb.GetKeyRotationStatusResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetParametersForImport returns the parameters required to import key material into a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetParametersForImport(ctx context.Context, req *connect.Request[pb.GetParametersForImportRequest]) (*connect.Response[pb.GetParametersForImportResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// GetPublicKey returns the public key of an asymmetric KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetPublicKey(ctx context.Context, req *connect.Request[pb.GetPublicKeyRequest]) (*connect.Response[pb.GetPublicKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ImportKeyMaterial imports key material into a KMS key that was created without key material.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ImportKeyMaterial(ctx context.Context, req *connect.Request[pb.ImportKeyMaterialRequest]) (*connect.Response[pb.ImportKeyMaterialResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListAliases returns a list of all aliases in the account and region.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAliases(ctx context.Context, req *connect.Request[pb.ListAliasesRequest]) (*connect.Response[pb.ListAliasesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListGrants returns a list of grants for a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListGrants(ctx context.Context, req *connect.Request[pb.ListGrantsRequest]) (*connect.Response[pb.ListGrantsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListKeyPolicies returns the names of key policies for a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListKeyPolicies(ctx context.Context, req *connect.Request[pb.ListKeyPoliciesRequest]) (*connect.Response[pb.ListKeyPoliciesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListKeyRotations returns a list of key rotation intervals for a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListKeyRotations(ctx context.Context, req *connect.Request[pb.ListKeyRotationsRequest]) (*connect.Response[pb.ListKeyRotationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListResourceTags returns the tags attached to a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListResourceTags(ctx context.Context, req *connect.Request[pb.ListResourceTagsRequest]) (*connect.Response[pb.ListResourceTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ListRetirableGrants returns grants that can be retired by the specified retiring principal.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRetirableGrants(ctx context.Context, req *connect.Request[pb.ListRetirableGrantsRequest]) (*connect.Response[pb.ListGrantsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// PutKeyPolicy attaches a key policy to a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutKeyPolicy(ctx context.Context, req *connect.Request[pb.PutKeyPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ReEncrypt decrypts ciphertext encrypted with one KMS key and re-encrypts it with a different KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ReEncrypt(ctx context.Context, req *connect.Request[pb.ReEncryptRequest]) (*connect.Response[pb.ReEncryptResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ReplicateKey replicates a multi-Region key to another region.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ReplicateKey(ctx context.Context, req *connect.Request[pb.ReplicateKeyRequest]) (*connect.Response[pb.ReplicateKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RetireGrant retires a grant, removing its permissions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RetireGrant(ctx context.Context, req *connect.Request[pb.RetireGrantRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RevokeGrant revokes a grant, immediately removing its permissions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RevokeGrant(ctx context.Context, req *connect.Request[pb.RevokeGrantRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// RotateKeyOnDemand rotates a symmetric KMS key immediately rather than waiting for the scheduled rotation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RotateKeyOnDemand(ctx context.Context, req *connect.Request[pb.RotateKeyOnDemandRequest]) (*connect.Response[pb.RotateKeyOnDemandResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// ScheduleKeyDeletion schedules the deletion of a KMS key after a waiting period.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ScheduleKeyDeletion(ctx context.Context, req *connect.Request[pb.ScheduleKeyDeletionRequest]) (*connect.Response[pb.ScheduleKeyDeletionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Sign generates a digital signature for a message using a private key from an asymmetric KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Sign(ctx context.Context, req *connect.Request[pb.SignRequest]) (*connect.Response[pb.SignResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// TagResource adds or overwrites one or more tags on a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UntagResource removes one or more tags from a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateAlias changes the alias associated with a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateAlias(ctx context.Context, req *connect.Request[pb.UpdateAliasRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateCustomKeyStore updates the configuration of a custom key store.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateCustomKeyStore(ctx context.Context, req *connect.Request[pb.UpdateCustomKeyStoreRequest]) (*connect.Response[pb.UpdateCustomKeyStoreResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdateKeyDescription updates the description of a KMS key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateKeyDescription(ctx context.Context, req *connect.Request[pb.UpdateKeyDescriptionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// UpdatePrimaryRegion moves a multi-Region primary key to a different region.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdatePrimaryRegion(ctx context.Context, req *connect.Request[pb.UpdatePrimaryRegionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// Verify verifies a digital signature generated by the Sign operation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) Verify(ctx context.Context, req *connect.Request[pb.VerifyRequest]) (*connect.Response[pb.VerifyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}

// VerifyMac verifies a message authentication code (MAC) generated by the GenerateMac operation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) VerifyMac(ctx context.Context, req *connect.Request[pb.VerifyMacRequest]) (*connect.Response[pb.VerifyMacResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for this operation"))
}
