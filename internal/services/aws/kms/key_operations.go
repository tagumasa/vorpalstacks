package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"vorpalstacks/internal/services/aws/common/pagination"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
	"vorpalstacks/internal/utils/aws/arn"
)

// CreateKey creates a new customer master key (CMK) in KMS.
// You can specify the key usage, key spec, origin, and tags for the new key.
// If no key usage is specified, ENCRYPT_DECRYPT is used by default.
// If no key spec is specified, SYMMETRIC_DEFAULT is used by default.
func (s *KMSService) CreateKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	description := request.GetStringParam(req.Parameters, "Description")
	keyUsage := request.GetStringParam(req.Parameters, "KeyUsage")
	keySpec := request.GetStringParam(req.Parameters, "KeySpec")
	origin := request.GetStringParam(req.Parameters, "Origin")
	multiRegion := request.GetBoolParam(req.Parameters, "MultiRegion")
	tagList := tagutil.ParseTagsWithKeyNames(req.Parameters, "Tags", "TagKey", "TagValue")

	if keyUsage == "" {
		keyUsage = string(kmsstore.KeyUsageEncryptDecrypt)
	}
	if keySpec == "" {
		keySpec = string(kmsstore.KeySpecSymmetricDefault)
	}
	if origin == "" {
		origin = string(kmsstore.OriginTypeAWSKMS)
	}

	keyID, err := kmsstore.GenerateKeyID()
	if err != nil {
		return nil, err
	}

	key, err := stores.keys.Create(
		keyID,
		kmsstore.KeyUsage(keyUsage),
		kmsstore.KeySpec(keySpec),
		description,
		kmsstore.OriginType(origin),
		multiRegion,
		tagList,
	)
	if err != nil {
		return nil, err
	}

	if err := s.hsmBackend.GenerateKey(keyID, hsm.KeySpec(keySpec)); err != nil {
		_ = stores.keys.Delete(keyID)
		return nil, err
	}

	policyStr := request.GetStringParam(req.Parameters, "Policy")
	if policyStr == "" {
		policyStr = kmsstore.DefaultKeyPolicy
	}
	if err := stores.keyPolicies.PutDefault(keyID, policyStr); err != nil {
		if delErr := stores.keys.Delete(keyID); delErr != nil {
			log.Printf("Failed to delete key during rollback: %v", delErr)
		}
		if delErr := s.hsmBackend.DeleteKey(keyID); delErr != nil {
			log.Printf("Failed to delete HSM key during rollback: %v", delErr)
		}
		return nil, err
	}

	return map[string]interface{}{
		"KeyMetadata": s.buildKeyMetadata(key),
	}, nil
}

// DescribeKey retrieves detailed information about a specified key.
// You can identify the key by its key ID, key ARN, or alias.
func (s *KMSService) DescribeKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyMetadata": s.buildKeyMetadata(key),
	}, nil
}

// ListKeys returns a list of all customer master keys (CMKs) in the account.
// The list includes the key ID and ARN for each key.
// Results can be paginated using the Marker and MaxItems parameters.
func (s *KMSService) ListKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	marker := pagination.GetMarker(req.Parameters)
	maxItems := pagination.GetMaxItems(req.Parameters, 100)

	result, err := stores.keys.List(marker, maxItems)
	if err != nil {
		return nil, err
	}

	keys := make([]map[string]interface{}, len(result.Keys))
	for i, k := range result.Keys {
		keys[i] = map[string]interface{}{
			"KeyId":  k.KeyID,
			"KeyArn": k.KeyArn,
		}
	}

	response := map[string]interface{}{
		"Keys": keys,
	}
	if result.IsTruncated {
		response["NextMarker"] = result.NextMarker
		response["Truncated"] = true
	} else {
		response["Truncated"] = false
	}

	return response, nil
}

// EnableKey enables a previously disabled key.
// A key must be enabled before it can be used for cryptographic operations.
func (s *KMSService) EnableKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := stores.keys.Enable(key.KeyID); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DisableKey disables a key, preventing it from being used for cryptographic operations.
// Disabled keys can be re-enabled using EnableKey.
func (s *KMSService) DisableKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := stores.keys.Disable(key.KeyID); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ScheduleKeyDeletion schedules the deletion of a key.
// The key will be deleted after the specified waiting period (7-30 days).
// You can cancel the deletion using CancelKeyDeletion before the deletion date.
func (s *KMSService) ScheduleKeyDeletion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	pendingWindowInDays := request.GetIntParam(req.Parameters, "PendingWindowInDays")
	if pendingWindowInDays == 0 {
		pendingWindowInDays = 30
	} else if pendingWindowInDays < 7 || pendingWindowInDays > 30 {
		return nil, ErrValidation
	}

	if err := stores.keys.ScheduleDeletion(key.KeyID, pendingWindowInDays); err != nil {
		return nil, err
	}

	updatedKey, err := stores.keys.Get(key.KeyID)
	if err != nil || updatedKey == nil {
		return nil, fmt.Errorf("failed to retrieve key after scheduling deletion: %w", err)
	}

	return map[string]interface{}{
		"KeyId":               key.Arn,
		"DeletionDate":        updatedKey.DeletionDate.Unix(),
		"KeyState":            updatedKey.KeyState,
		"PendingWindowInDays": pendingWindowInDays,
	}, nil
}

// CancelKeyDeletion cancels the scheduled deletion of a key.
// This reverts the key to its state before ScheduleKeyDeletion was called.
func (s *KMSService) CancelKeyDeletion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := stores.keys.CancelDeletion(key.KeyID); err != nil {
		return nil, err
	}

	updatedKey, err := stores.keys.Get(key.KeyID)
	if err != nil || updatedKey == nil {
		return nil, fmt.Errorf("failed to retrieve key after cancelling deletion: %w", err)
	}

	return map[string]interface{}{
		"KeyId": key.Arn,
	}, nil
}

// UpdateKeyDescription updates the description of a specified key.
func (s *KMSService) UpdateKeyDescription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	description := request.GetStringParam(req.Parameters, "Description")

	if err := stores.keys.UpdateDescription(key.KeyID, description); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// buildKeyMetadata constructs a KeyMetadata map from a Key store object.
func (s *KMSService) buildKeyMetadata(key *kmsstore.Key) map[string]interface{} {
	_, _, _, accountID, _ := arn.SplitARN(key.Arn)
	metadata := map[string]interface{}{
		"AWSAccountId":       accountID,
		"KeyId":              key.KeyID,
		"Arn":                key.Arn,
		"KeyState":           key.KeyState,
		"Enabled":            key.Enabled,
		"KeyUsage":           key.KeyUsage,
		"KeySpec":            key.KeySpec,
		"CreationDate":       key.CreationDate.Unix(),
		"Origin":             key.Origin,
		"KeyManager":         key.KeyManager,
		"KeyRotationEnabled": key.KeyRotationEnabled,
		"MultiRegion":        key.MultiRegion,
	}

	if key.Description != "" {
		metadata["Description"] = key.Description
	}
	if key.DeletionDate != nil {
		metadata["DeletionDate"] = key.DeletionDate.Unix()
	}
	if key.ValidTo != nil {
		metadata["ValidTo"] = key.ValidTo.Unix()
	}
	if key.CustomKeyStoreID != "" {
		metadata["CustomKeyStoreId"] = key.CustomKeyStoreID
	}
	if len(key.EncryptionAlgorithms) > 0 {
		metadata["EncryptionAlgorithms"] = key.EncryptionAlgorithms
	}
	if len(key.SigningAlgorithms) > 0 {
		metadata["SigningAlgorithms"] = key.SigningAlgorithms
	}
	if key.MultiRegion {
		_, _, region, _, _ := arn.SplitARN(key.Arn)
		metadata["MultiRegionConfiguration"] = map[string]interface{}{
			"MultiRegionKeyType": "PRIMARY",
			"PrimaryKey": map[string]interface{}{
				"Arn":    key.Arn,
				"Region": region,
			},
			"ReplicaKeys": []interface{}{},
		}
	}

	return metadata
}

// GetParametersForImport returns the items required to import key material into a CMK.
// This includes an import token and the public key that wraps the key material.
// You must import the key material within the validity period of these parameters.
func (s *KMSService) GetParametersForImport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	keyID := s.getKeyID(req.Parameters)
	wrappingKeySpec := request.GetStringParam(req.Parameters, "WrappingKeySpec")
	if wrappingKeySpec == "" {
		wrappingKeySpec = "RSA_2048"
	}

	importToken, publicKey, err := stores.keys.GetParametersForImport(keyID, wrappingKeySpec)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyId":           keyID,
		"ImportToken":     importToken,
		"PublicKey":       publicKey,
		"ParametersValid": true,
	}, nil
}

// ImportKeyMaterial imports key material into a CMK that was created with external origin.
// The key material is encrypted using the public key from GetParametersForImport.
func (s *KMSService) ImportKeyMaterial(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	keyID := s.getKeyID(req.Parameters)
	importToken := request.GetStringParam(req.Parameters, "ImportToken")
	encryptedKeyMaterialB64 := request.GetStringParam(req.Parameters, "EncryptedKeyMaterial")
	validTo := request.GetIntParam(req.Parameters, "ValidTo")

	encryptedKeyMaterial, err := base64.StdEncoding.DecodeString(encryptedKeyMaterialB64)
	if err != nil {
		return nil, fmt.Errorf("invalid EncryptedKeyMaterial: %w", err)
	}

	var validToTime *time.Time
	if validTo > 0 {
		t := time.Unix(int64(validTo), 0)
		validToTime = &t
	}

	if err := stores.keys.ImportKeyMaterial(keyID, importToken, encryptedKeyMaterial, validToTime); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteImportedKeyMaterial deletes the imported key material from a CMK.
// This makes the CMK unusable until new key material is imported.
func (s *KMSService) DeleteImportedKeyMaterial(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	keyID := s.getKeyID(req.Parameters)

	if err := stores.keys.DeleteImportedKeyMaterial(keyID); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ReplicateKey replicates a multi-Region key to another region.
// This creates a replica key in the specified region with the same key material.
func (s *KMSService) ReplicateKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	keyID := s.getKeyID(req.Parameters)
	replicaRegion := request.GetStringParam(req.Parameters, "ReplicaRegion")

	replicaKeyID, err := kmsstore.GenerateKeyID()
	if err != nil {
		return nil, err
	}

	replicaKey, err := stores.keys.ReplicateKey(keyID, replicaRegion, replicaKeyID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ReplicaKeyMetadata": s.buildKeyMetadata(replicaKey),
		"ReplicaKeyArn":      replicaKey.Arn,
	}, nil
}

// UpdatePrimaryRegion changes the primary region of a multi-Region key.
// The primary region is the one where the key is the authoritative copy.
func (s *KMSService) UpdatePrimaryRegion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	keyID := s.getKeyID(req.Parameters)
	primaryRegion := request.GetStringParam(req.Parameters, "PrimaryRegion")

	if err := stores.keys.UpdatePrimaryRegion(keyID, primaryRegion); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
