package kms

// key_lifecycle_core.go carries the Core functions of the KMS key
// lifecycle beyond creation (enable/disable, deletion scheduling, import,
// multi-region replication) plus the thin enable/disable mutations. Every
// Core replays the original validation-before-persistence order.

import (
	"fmt"
	"time"

	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// enableKeyCore is the single entry point for enabling a previously
// disabled key.
func (s *KMSService) enableKeyCore(stores *kmsStores, key *kmsstore.Key) error {
	return stores.keys.Enable(key.KeyID)
}

// disableKeyCore is the single entry point for disabling a key.
func (s *KMSService) disableKeyCore(stores *kmsStores, key *kmsstore.Key) error {
	return stores.keys.Disable(key.KeyID)
}

// cancelKeyDeletionCore is the single entry point for cancelling a
// scheduled deletion. It returns the key ID and the restored key state
// (Enabled or Disabled depending on the pre-ScheduleKeyDeletion value).
func (s *KMSService) cancelKeyDeletionCore(stores *kmsStores, key *kmsstore.Key) (string, kmsstore.KeyState, error) {
	if err := stores.keys.CancelDeletion(key.KeyID); err != nil {
		return "", "", err
	}

	updatedKey, err := stores.keys.Get(key.KeyID)
	if err != nil || updatedKey == nil {
		return "", "", fmt.Errorf("failed to retrieve key after cancelling deletion: %w", err)
	}

	return key.KeyID, updatedKey.KeyState, nil
}

// updateKeyDescriptionCore is the single entry point for updating a key's
// description (Smithy DescriptionType: 0-8192).
func (s *KMSService) updateKeyDescriptionCore(stores *kmsStores, key *kmsstore.Key, description string) error {
	if err := validateDescriptionLength(description); err != nil {
		return err
	}

	return stores.keys.UpdateDescription(key.KeyID, description)
}

// getParametersForImportCore is the single entry point for the import
// parameter exchange. It enforces the PendingImport state requirement,
// the EXTERNAL origin requirement and the wrapping algorithm/spec
// envelopes before issuing the token and wrapping public key.
func (s *KMSService) getParametersForImportCore(stores *kmsStores, principal, keyID, wrappingAlgorithm, wrappingKeySpec string) (string, string, error) {
	if err := s.authorizeOperation(stores, principal, "GetParametersForImport", keyID, nil); err != nil {
		return "", "", err
	}

	// AWS: GetParametersForImport requires the key to be in PendingImport
	// state. Calling it on an Enabled key or a key with already-imported
	// material returns KMSInvalidStateException.
	key, err := stores.keys.Get(keyID)
	if err != nil {
		return "", "", NewKeyNotFoundError(keyID)
	}
	if key.Origin != kmsstore.OriginTypeExternal {
		return "", "", ErrUnsupportedOperation
	}
	if key.KeyState != kmsstore.KeyStatePendingImport {
		return "", "", ErrKeyPendingImport
	}

	// AWS: WrappingAlgorithm selects between symmetric and RSA wrapping.
	// The platform implements RSA_* wrapping only; SYMMETRIC_KEY_WRAPPING
	// is rejected as unsupported until the underlying implementation lands.
	switch wrappingAlgorithm {
	case "":
		wrappingAlgorithm = "RSAES_OAEP_SHA_256"
	case "RSAES_OAEP_SHA_256", "RSAES_OAEP_SHA_1":
		// supported
	default:
		return "", "", ErrUnsupportedOperation
	}

	// AWS: WrappingKeySpec is a required parameter (smithy.api#required
	// trait on the GetParametersForImport input). The previous code
	// silently defaulted to RSA_2048 when omitted, diverging from the
	// documented contract.
	switch wrappingKeySpec {
	case "RSA_2048", "RSA_4096":
		// supported
	default:
		return "", "", NewValidationError(fmt.Sprintf("WrappingKeySpec must be RSA_2048 or RSA_4096, got %q", wrappingKeySpec))
	}

	return stores.keys.GetParametersForImport(keyID, wrappingKeySpec, wrappingAlgorithm)
}

// ImportKeyMaterialInput carries the ImportKeyMaterial members; the
// encrypted key material travels in its wire (base64) form so the Core
// can run the size and decode validations in their original position
// after the key fetch.
type ImportKeyMaterialInput struct {
	KeyID                   string
	ImportToken             string
	EncryptedKeyMaterialB64 string
	ValidTo                 int
	ExpirationModel         string
	Principal               string
}

// importKeyMaterialCore is the single entry point for importing key
// material into an EXTERNAL-origin key, including the expiration-model
// resolution and the post-import state refresh.
func (s *KMSService) importKeyMaterialCore(stores *kmsStores, in ImportKeyMaterialInput) error {
	if err := s.authorizeOperation(stores, in.Principal, "ImportKeyMaterial", in.KeyID, nil); err != nil {
		return err
	}

	key, err := stores.keys.Get(in.KeyID)
	if err != nil {
		return err
	}

	// Enforce the CiphertextType constraints: both blobs must decode to at
	// most 6144 bytes, and oversized encoded input is rejected before the
	// decoder allocates memory.
	if err := validateImportTokenSize(in.ImportToken); err != nil {
		return err
	}
	encryptedKeyMaterial, err := decodeEncryptedKeyMaterial(in.EncryptedKeyMaterialB64)
	if err != nil {
		return err
	}

	// Reject expired import tokens. GetParametersForImport issues tokens
	// with a 24-hour validity window; using an expired token returns
	// ExpiredImportTokenException per the AWS contract.
	if key.ImportTokenValidTo != nil && time.Now().After(*key.ImportTokenValidTo) {
		return ErrExpiredImportToken
	}

	// AWS: ExpirationModel controls whether the imported key material
	// expires. KEY_MATERIAL_DOES_NOT_EXPIRE causes ValidTo to be ignored.
	// When ExpirationModel is omitted but ValidTo is supplied, the default
	// is KEY_MATERIAL_EXPIRES. When neither is supplied, the default is
	// KEY_MATERIAL_DOES_NOT_EXPIRE.
	expirationModel := in.ExpirationModel
	var validToTime *time.Time
	if expirationModel != "KEY_MATERIAL_DOES_NOT_EXPIRE" && in.ValidTo > 0 {
		t := time.Unix(int64(in.ValidTo), 0)
		validToTime = &t
	}
	if expirationModel == "" {
		if validToTime != nil {
			expirationModel = "KEY_MATERIAL_EXPIRES"
		} else {
			expirationModel = "KEY_MATERIAL_DOES_NOT_EXPIRE"
		}
	}

	rawKeyMaterial, err := stores.keys.ImportKeyMaterial(in.KeyID, in.ImportToken, encryptedKeyMaterial, validToTime)
	if err != nil {
		return err
	}

	if err := s.hsmBackend.ImportKey(in.KeyID, rawKeyMaterial, hsm.KeySpec(key.KeySpec)); err != nil {
		return err
	}

	// Re-fetch the key because ImportKeyMaterial's atomicUpdate changed
	// its state (PendingImport → Enabled) and we must not overwrite that
	// with the stale pre-import copy.
	updatedKey, err := stores.keys.Get(in.KeyID)
	if err != nil {
		return err
	}
	updatedKey.ExpirationModel = expirationModel
	return stores.keys.Update(updatedKey)
}

// deleteImportedKeyMaterialCore is the single entry point for deleting
// imported key material. The HSM delete runs first so that on failure the
// store still reflects the imported state (Enabled + key material
// present); the previous ordering (store → HSM) left an inconsistent
// state when the HSM delete failed.
func (s *KMSService) deleteImportedKeyMaterialCore(stores *kmsStores, principal, keyID string) error {
	if err := s.authorizeOperation(stores, principal, "DeleteImportedKeyMaterial", keyID, nil); err != nil {
		return err
	}

	if err := s.hsmBackend.DeleteKey(keyID); err != nil {
		logs.Error("DeleteImportedKeyMaterial: HSM DeleteKey failed", logs.String("keyId", keyID), logs.Err(err))
		return ErrKMSInternal
	}
	return stores.keys.DeleteImportedKeyMaterial(keyID)
}

// ReplicateKeyInput carries the ReplicateKey members. AccountID and Region
// come from the request context because the replica's primary reference
// records the primary's region.
type ReplicateKeyInput struct {
	KeyID         string
	ReplicaRegion string
	Policy        string
	Description   string
	BypassCheck   bool
	Tags          []tagutil.Tag
	Principal     string
	AccountID     string
	Region        string
}

// replicateKeyCore is the single entry point for replicating a
// multi-Region key. The replica is persisted with the replica region's
// ARN builder, registered in the primary's MultiRegionConfiguration,
// inherits the primary's policy and tags, and is provisioned in the HSM
// with the primary's key material; any failure after the replica write
// cascades the partial state away.
func (s *KMSService) replicateKeyCore(stores *kmsStores, in ReplicateKeyInput) (*kmsstore.Key, error) {
	if in.ReplicaRegion == "" {
		return nil, NewValidationError("ReplicaRegion is required")
	}
	if err := validateReplicaRegion(in.ReplicaRegion); err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, in.Principal, "ReplicateKey", in.KeyID, nil); err != nil {
		return nil, err
	}

	primary, err := stores.keys.Get(in.KeyID)
	if err != nil {
		return nil, NewKeyNotFoundError(in.KeyID)
	}
	if !primary.MultiRegion {
		return nil, kmsstore.ErrNotMultiRegionKey
	}

	replicaKeyID, err := kmsstore.GenerateKeyID()
	if err != nil {
		return nil, err
	}

	// Obtain the replica region's stores so the replica key is persisted
	// with the correct region's ARN and bucket. The previous
	// implementation saved the replica to the local (primary) store with
	// a local-region ARN, making it invisible to clients in the replica
	// region.
	replicaStores, err := s.GetStoreForRegion(in.ReplicaRegion)
	if err != nil {
		return nil, err
	}

	// Build the replica key using the replica region's ARN builder.
	replicaArn := replicaStores.keys.ARNBuilder().KeyArn(replicaKeyID)
	replicaKey := &kmsstore.Key{
		KeyID:              replicaKeyID,
		Arn:                replicaArn,
		KeyState:           kmsstore.KeyStateEnabled,
		KeyUsage:           primary.KeyUsage,
		KeySpec:            primary.KeySpec,
		Description:        primary.Description,
		Enabled:            true,
		CreationDate:       time.Now(),
		Origin:             primary.Origin,
		KeyManager:         primary.KeyManager,
		KeyRotationEnabled: primary.KeyRotationEnabled,
		MultiRegion:        true,
		MultiRegionConfiguration: &kmsstore.MultiRegionConfiguration{
			MultiRegionKeyType: "REPLICA",
			PrimaryKey: &kmsstore.PrimaryKeyInfo{
				Arn:    primary.Arn,
				Region: in.Region,
			},
		},
	}

	// Apply optional caller-supplied parameters to the replica. Per the
	// Smithy ReplicateKeyRequest, the caller may set Policy, Description,
	// BypassPolicyLockoutSafetyCheck, and Tags on the replica.
	replicaPolicy := in.Policy
	if replicaPolicy == "" {
		// Inherit the primary's key policy when not explicitly supplied.
		primaryPolicy, pErr := stores.keyPolicies.GetDefault(primary.KeyID)
		if pErr == nil {
			replicaPolicy = primaryPolicy.PolicyDocument
		} else {
			replicaPolicy = kmsstore.DefaultKeyPolicy
		}
	}
	if err := validatePolicySize(replicaPolicy); err != nil {
		return nil, err
	}
	if !in.BypassCheck {
		if err := validatePolicyDoesNotLockOutRoot(replicaPolicy, in.AccountID); err != nil {
			return nil, err
		}
	}
	if in.Description != "" {
		if err := validateDescriptionLength(in.Description); err != nil {
			return nil, err
		}
		replicaKey.Description = in.Description
	}
	replicaKey.BypassPolicyLockoutSafetyCheck = in.BypassCheck

	// Step 1: Save the replica key to the replica region's store.
	if err := replicaStores.keys.CreateReplica(replicaKey); err != nil {
		return nil, err
	}

	// cleanupReplica removes the replica from the replica region store,
	// the primary's MultiRegionConfiguration, and the HSM. Used on any
	// failure after step 1.
	cleanupReplica := func() {
		if delErr := replicaStores.CascadeDeleteKey(s.hsmBackend, replicaKeyID); delErr != nil {
			logs.Error("ReplicateKey: failed to cascade-delete replica from replica store",
				logs.Err(delErr), logs.String("replicaKeyId", replicaKeyID))
		}
		if rmErr := stores.keys.RemoveReplicaFromPrimary(primary.KeyID, in.ReplicaRegion); rmErr != nil {
			logs.Error("ReplicateKey: failed to remove replica from primary config",
				logs.Err(rmErr), logs.String("replicaRegion", in.ReplicaRegion))
		}
	}

	// Step 2: Register the replica in the primary's MultiRegionConfiguration.
	replicaInfo := kmsstore.ReplicaKeyInfo{
		Arn:    replicaArn,
		Region: in.ReplicaRegion,
	}
	if err := stores.keys.AddReplicaToPrimary(primary.KeyID, replicaInfo); err != nil {
		cleanupReplica()
		return nil, err
	}

	// Step 3: Copy tags from the primary to the replica (cross-region).
	sourceTags, err := stores.keys.TagStore.ListAsSlice(primary.KeyID)
	if err != nil {
		cleanupReplica()
		return nil, err
	}
	if len(sourceTags) > 0 {
		if err := replicaStores.keys.TagStore.TagFromSlice(replicaKeyID, sourceTags); err != nil {
			cleanupReplica()
			return nil, err
		}
	}

	// Step 4: Apply caller-supplied tags (in addition to inherited tags).
	if len(in.Tags) > 0 {
		if err := validateKMSTags(in.Tags); err != nil {
			cleanupReplica()
			return nil, err
		}
		if err := replicaStores.keys.TagStore.Tag(replicaKeyID, tagutil.ToMap(in.Tags)); err != nil {
			cleanupReplica()
			return nil, err
		}
	}

	// Step 5: Store the replica's key policy in the replica region.
	if err := replicaStores.keyPolicies.PutDefault(replicaKeyID, replicaPolicy); err != nil {
		cleanupReplica()
		return nil, err
	}

	// Step 6: Provision the replica in the HSM with the primary's key material.
	if err := s.hsmBackend.ReplicateKey(primary.KeyID, replicaKeyID); err != nil {
		cleanupReplica()
		return nil, err
	}

	return replicaKey, nil
}

// updatePrimaryRegionCore is the single entry point for changing a
// multi-Region key's primary region. The empty-region check runs before
// the store call because the store would otherwise persist a broken
// PrimaryKeyInfo with an empty Region field.
func (s *KMSService) updatePrimaryRegionCore(stores *kmsStores, principal, keyID, primaryRegion string) error {
	if err := s.authorizeOperation(stores, principal, "UpdatePrimaryRegion", keyID, nil); err != nil {
		return err
	}

	if err := validatePrimaryRegion(primaryRegion); err != nil {
		return err
	}

	return stores.keys.UpdatePrimaryRegion(keyID, primaryRegion)
}
