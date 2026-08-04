package kms

import (
	stderrors "errors"
	"fmt"
	"time"

	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

// ---------------------------------------------------------------------------
// Transport-agnostic DTOs for CreateKey
// ---------------------------------------------------------------------------

// CreateKeyInput carries every field that CreateKey needs, in a format
// independent of the wire protocol. Both the HTTP API handler and the
// admin gRPC handler build this struct and delegate to createKeyCore,
// ensuring that validation, key creation, and policy handling follow a
// single code path.
type CreateKeyInput struct {
	Description        string
	KeyUsage           string
	KeySpec            string
	Origin             string
	MultiRegion        bool
	CustomKeyStoreID   string
	XksKeyID           string
	Policy             string
	BypassLockoutCheck bool
	Tags               []types.Tag
	AccountID          string
}

// ---------------------------------------------------------------------------
// Service-layer result DTOs (no store types)
// ---------------------------------------------------------------------------

// KeyMetadataResult is the transport-agnostic representation of KMS key
// metadata. Core functions return this type so that admin handlers never
// need to import the store package.
type KeyMetadataResult struct {
	KeyID                    string
	Arn                      string
	KeyState                 string
	KeyUsage                 string
	KeySpec                  string
	Description              string
	Enabled                  bool
	Origin                   string
	KeyManager               string
	MultiRegion              bool
	CreationDate             time.Time
	DeletionDate             *time.Time
	PendingWindowInDays      int
	ValidTo                  *time.Time
	ExpirationModel          string
	BypassPolicyLockoutCheck bool
	CustomKeyStoreID         string
	EncryptionAlgorithms     []string
	SigningAlgorithms        []string
	MacAlgorithms            []string
	MultiRegionConfig        *MultiRegionConfigResult
}

// MultiRegionConfigResult is the service-layer representation of
// multi-region key configuration.
type MultiRegionConfigResult struct {
	MultiRegionKeyType string
	PrimaryKey         *MultiRegionKeyResult
	ReplicaKeys        []MultiRegionKeyResult
}

// MultiRegionKeyResult represents a primary or replica key in a
// multi-region configuration.
type MultiRegionKeyResult struct {
	Arn    string
	Region string
}

// KeyListEntry is a single entry in a ListKeys result.
type KeyListEntry struct {
	KeyID  string
	KeyArn string
}

// ListKeysResult is the paginated result of ListKeys.
type ListKeysResult struct {
	Keys        []KeyListEntry
	NextMarker  string
	IsTruncated bool
}

// keyToMetadataResult converts a store-layer Key into the service-layer
// KeyMetadataResult DTO.
func keyToMetadataResult(key *kmsstore.Key) *KeyMetadataResult {
	result := &KeyMetadataResult{
		KeyID:                    key.KeyID,
		Arn:                      key.Arn,
		KeyState:                 string(key.KeyState),
		KeyUsage:                 string(key.KeyUsage),
		KeySpec:                  string(key.KeySpec),
		Description:              key.Description,
		Enabled:                  key.Enabled,
		Origin:                   string(key.Origin),
		KeyManager:               string(key.KeyManager),
		MultiRegion:              key.MultiRegion,
		CreationDate:             key.CreationDate,
		DeletionDate:             key.DeletionDate,
		PendingWindowInDays:      key.PendingWindowInDays,
		ValidTo:                  key.ValidTo,
		ExpirationModel:          key.ExpirationModel,
		BypassPolicyLockoutCheck: key.BypassPolicyLockoutSafetyCheck,
		CustomKeyStoreID:         key.CustomKeyStoreID,
		EncryptionAlgorithms:     key.EncryptionAlgorithms,
		SigningAlgorithms:        key.SigningAlgorithms,
		MacAlgorithms:            key.MacAlgorithms,
	}
	if key.MultiRegionConfiguration != nil {
		mrc := &MultiRegionConfigResult{
			MultiRegionKeyType: key.MultiRegionConfiguration.MultiRegionKeyType,
		}
		if key.MultiRegionConfiguration.PrimaryKey != nil {
			mrc.PrimaryKey = &MultiRegionKeyResult{
				Arn:    key.MultiRegionConfiguration.PrimaryKey.Arn,
				Region: key.MultiRegionConfiguration.PrimaryKey.Region,
			}
		}
		for _, r := range key.MultiRegionConfiguration.ReplicaKeys {
			mrc.ReplicaKeys = append(mrc.ReplicaKeys, MultiRegionKeyResult{
				Arn:    r.Arn,
				Region: r.Region,
			})
		}
		result.MultiRegionConfig = mrc
	}
	return result
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createKeyCore is the single entry point for key creation shared by the
// HTTP API and the admin gRPC handler. It applies defaults, performs all
// validation, generates the key material (HSM or import-pending), persists
// the key policy and tags, and returns the fully-created key metadata. On
// any failure after key creation it cleans up partial state via
// CascadeDeleteKey.
func (s *KMSService) createKeyCore(stores *kmsStores, in CreateKeyInput) (*KeyMetadataResult, error) {
	// Convert string fields to store enum types.
	keyUsage := kmsstore.KeyUsage(in.KeyUsage)
	keySpec := kmsstore.KeySpec(in.KeySpec)
	origin := kmsstore.OriginType(in.Origin)

	// 1. Apply defaults (Smithy: KeyUsage/KeySpec/Origin are OPTIONAL).
	if keyUsage == "" {
		keyUsage = kmsstore.KeyUsageEncryptDecrypt
	}
	if keySpec == "" {
		keySpec = kmsstore.KeySpecSymmetricDefault
	}
	if origin == "" {
		origin = kmsstore.OriginTypeAWSKMS
	}

	// 2. AWS_CLOUDHSM origin requires a CustomKeyStoreId pointing at a real
	// CloudHSM custom key store. CloudHSM integration is not implemented,
	// so reject AWS_CLOUDHSM explicitly rather than silently storing a key
	// that no HSM backend can serve.
	if origin == kmsstore.OriginTypeAWSCloudHSM {
		return nil, ErrUnsupportedOperation
	}

	// 3. CustomKeyStoreId and XksKeyId must not be silently dropped.
	if err := validateOriginParams(in.CustomKeyStoreID, in.XksKeyID); err != nil {
		return nil, err
	}

	// 4. KeyUsage/KeySpec combination validation (Smithy matrix).
	if err := validateKeyUsageSpecCombo(keyUsage, keySpec); err != nil {
		return nil, err
	}

	// 5. Description length (Smithy DescriptionType: 0-8192).
	if err := validateDescriptionLength(in.Description); err != nil {
		return nil, err
	}

	// 6. Validate tags BEFORE creating any state so that an invalid tag
	// list cannot leave a half-created key behind.
	if len(in.Tags) > 0 {
		if err := validateKMSTags(in.Tags); err != nil {
			return nil, err
		}
	}

	// 7. Policy defaults and size validation (Smithy PolicyType: 1-131072).
	policyStr := in.Policy
	if policyStr == "" {
		policyStr = kmsstore.DefaultKeyPolicy
	}
	if err := validatePolicySize(policyStr); err != nil {
		return nil, err
	}

	// 8. Validate the supplied policy does not lock out the root principal
	// unless BypassPolicyLockoutSafetyCheck is explicitly true.
	if !in.BypassLockoutCheck {
		if err := validatePolicyDoesNotLockOutRoot(policyStr, in.AccountID); err != nil {
			return nil, err
		}
	}

	// 9. Generate key ID.
	keyID, err := kmsstore.GenerateKeyID()
	if err != nil {
		return nil, err
	}

	// 10. Create key in store.
	key, err := stores.keys.Create(
		keyID,
		keyUsage,
		keySpec,
		in.Description,
		origin,
		in.MultiRegion,
	)
	if err != nil {
		return nil, err
	}

	// 11. Apply tags.
	if len(in.Tags) > 0 {
		if err := stores.keys.TagStore.Tag(keyID, tagutil.ToMap(in.Tags)); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after Tag failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, err
		}
	}

	// 12. Set up key material.
	if origin == kmsstore.OriginTypeExternal {
		if err := stores.keys.SetPendingImport(keyID); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after SetPendingImport failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, err
		}
		key.KeyState = kmsstore.KeyStatePendingImport
		key.Enabled = false
	} else {
		if err := s.hsmBackend.GenerateKey(keyID, hsm.KeySpec(keySpec)); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after HSM GenerateKey failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, err
		}
	}

	// 13. Set bypass flag and persist.
	key.BypassPolicyLockoutSafetyCheck = in.BypassLockoutCheck
	if err := stores.keys.Update(key); err != nil {
		if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
			logs.Error("Failed to cascade-delete key after Update", logs.Err(delErr), logs.String("keyId", keyID))
		}
		return nil, err
	}

	// 14. Persist key policy.
	if err := stores.keyPolicies.PutDefault(keyID, policyStr); err != nil {
		if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
			logs.Error("Failed to cascade-delete key after PutDefault policy failure", logs.Err(delErr), logs.String("keyId", keyID))
		}
		return nil, err
	}

	return keyToMetadataResult(key), nil
}

// scheduleKeyDeletionCore is the single entry point for scheduling key
// deletion shared by the HTTP API and the admin gRPC handler. It
// validates the pending window (AWS range 7-30), schedules deletion in
// the store, and returns the updated key metadata.
func (s *KMSService) scheduleKeyDeletionCore(stores *kmsStores, keyID string, pendingWindowInDays int) (*KeyMetadataResult, int, error) {
	// Validate pending window (AWS: 7-30).
	if pendingWindowInDays < 7 || pendingWindowInDays > 30 {
		return nil, 0, NewValidationError(fmt.Sprintf("PendingWindowInDays %d does not fit range 7-30", pendingWindowInDays))
	}

	// Schedule deletion. Translate the store-layer ErrInvalidKeyState
	// (key is already PendingDeletion) into the service-layer
	// KMSInvalidStateException so that both the HTTP path and the admin
	// gRPC path surface the correct HTTP 400 / CodeInvalidArgument.
	if err := stores.keys.ScheduleDeletion(keyID, pendingWindowInDays); err != nil {
		if stderrors.Is(err, kmsstore.ErrInvalidKeyState) {
			return nil, 0, ErrKeyPendingDeletion
		}
		return nil, 0, err
	}

	// Retrieve updated key.
	updatedKey, err := stores.keys.Get(keyID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve key after scheduling deletion: %w", err)
	}
	if updatedKey == nil {
		return nil, 0, fmt.Errorf("failed to retrieve key after scheduling deletion: store returned nil")
	}

	return keyToMetadataResult(updatedKey), pendingWindowInDays, nil
}

// listKeysCore is the single entry point for ListKeys shared by the HTTP
// API and the admin gRPC handler. It applies the Smithy @range(1-1000)
// limit clamp and delegates to the store.
func (s *KMSService) listKeysCore(stores *kmsStores, marker string, limit int) (*ListKeysResult, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	result, err := stores.keys.List(marker, limit)
	if err != nil {
		return nil, err
	}

	entries := make([]KeyListEntry, len(result.Keys))
	for i, k := range result.Keys {
		entries[i] = KeyListEntry{
			KeyID:  k.KeyID,
			KeyArn: k.KeyArn,
		}
	}

	return &ListKeysResult{
		Keys:        entries,
		NextMarker:  result.NextMarker,
		IsTruncated: result.IsTruncated,
	}, nil
}

// resolveKeyID resolves a key identifier (key ID, key ARN, or alias) to
// the canonical key ID string. This is used by the admin handler to
// avoid leaking the store-layer Key type.
func (s *KMSService) resolveKeyID(stores *kmsStores, keyID string) (string, error) {
	if keyID == "" {
		return "", ErrKeyNotFound
	}
	if err := validateKeyIdLength(keyID); err != nil {
		return "", err
	}

	if stores.keys.ARNBuilder().IsAlias(keyID) {
		alias, err := stores.aliases.Get(keyID)
		if err != nil {
			return "", NewAliasNotFoundError(keyID)
		}
		keyID = alias.TargetKeyID
	}

	key, err := stores.keys.Get(keyID)
	if err != nil {
		return "", NewKeyNotFoundError(keyID)
	}
	return key.KeyID, nil
}

// keyMetadataResultToResponse converts a KeyMetadataResult into the
// HTTP API response map, mirroring the output of buildKeyMetadata.
func keyMetadataResultToResponse(meta *KeyMetadataResult) map[string]interface{} {
	_, _, _, accountID, _ := arn.SplitARN(meta.Arn)
	metadata := map[string]interface{}{
		"AWSAccountId":          accountID,
		"KeyId":                 meta.KeyID,
		"Arn":                   meta.Arn,
		"KeyState":              meta.KeyState,
		"Enabled":               meta.Enabled,
		"KeyUsage":              meta.KeyUsage,
		"KeySpec":               meta.KeySpec,
		"CustomerMasterKeySpec": meta.KeySpec,
		"CreationDate":          meta.CreationDate.Unix(),
		"Origin":                meta.Origin,
		"KeyManager":            meta.KeyManager,
		"MultiRegion":           meta.MultiRegion,
	}

	if meta.Description != "" {
		metadata["Description"] = meta.Description
	}
	if meta.DeletionDate != nil {
		metadata["DeletionDate"] = meta.DeletionDate.Unix()
	}
	if meta.KeyState == "PendingDeletion" && meta.PendingWindowInDays > 0 {
		metadata["PendingDeletionWindowInDays"] = meta.PendingWindowInDays
	}
	if meta.ValidTo != nil {
		metadata["ValidTo"] = meta.ValidTo.Unix()
	}
	if meta.Origin == "EXTERNAL" && meta.ExpirationModel != "" {
		metadata["ExpirationModel"] = meta.ExpirationModel
	}
	if meta.CustomKeyStoreID != "" {
		metadata["CustomKeyStoreId"] = meta.CustomKeyStoreID
	}
	if len(meta.EncryptionAlgorithms) > 0 {
		metadata["EncryptionAlgorithms"] = meta.EncryptionAlgorithms
	}
	if len(meta.SigningAlgorithms) > 0 {
		metadata["SigningAlgorithms"] = meta.SigningAlgorithms
	}
	if len(meta.MacAlgorithms) > 0 {
		metadata["MacAlgorithms"] = meta.MacAlgorithms
	}
	if meta.MultiRegion {
		if meta.MultiRegionConfig != nil {
			mrc := map[string]interface{}{
				"MultiRegionKeyType": meta.MultiRegionConfig.MultiRegionKeyType,
			}
			if meta.MultiRegionConfig.PrimaryKey != nil {
				mrc["PrimaryKey"] = map[string]interface{}{
					"Arn":    meta.MultiRegionConfig.PrimaryKey.Arn,
					"Region": meta.MultiRegionConfig.PrimaryKey.Region,
				}
			}
			replicas := make([]interface{}, 0, len(meta.MultiRegionConfig.ReplicaKeys))
			for _, r := range meta.MultiRegionConfig.ReplicaKeys {
				replicas = append(replicas, map[string]interface{}{
					"Arn":    r.Arn,
					"Region": r.Region,
				})
			}
			mrc["ReplicaKeys"] = replicas
			metadata["MultiRegionConfiguration"] = mrc
		} else {
			_, _, region, _, _ := arn.SplitARN(meta.Arn)
			metadata["MultiRegionConfiguration"] = map[string]interface{}{
				"MultiRegionKeyType": "PRIMARY",
				"PrimaryKey": map[string]interface{}{
					"Arn":    meta.Arn,
					"Region": region,
				},
				"ReplicaKeys": []interface{}{},
			}
		}
	}

	return metadata
}
