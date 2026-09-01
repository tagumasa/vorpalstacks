package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	kmsstore "vorpalstacks/internal/store/aws/kms"
	"vorpalstacks/internal/utils/aws/arn"
)

var keyUsageKeySpecConstraints = map[kmsstore.KeyUsage][]kmsstore.KeySpec{
	kmsstore.KeyUsageEncryptDecrypt: {
		kmsstore.KeySpecSymmetricDefault,
		kmsstore.KeySpecRSA2048, kmsstore.KeySpecRSA3072, kmsstore.KeySpecRSA4096,
		kmsstore.KeySpecSM2,
	},
	kmsstore.KeyUsageSignVerify: {
		kmsstore.KeySpecRSA2048, kmsstore.KeySpecRSA3072, kmsstore.KeySpecRSA4096,
		kmsstore.KeySpecECCNISTP256, kmsstore.KeySpecECCNISTP384, kmsstore.KeySpecECCNISTP521,
		kmsstore.KeySpecECCSECGP256K1, kmsstore.KeySpecSM2,
	},
	kmsstore.KeyUsageGenerateVerifyMAC: {
		kmsstore.KeySpecHMAC224, kmsstore.KeySpecHMAC256, kmsstore.KeySpecHMAC384, kmsstore.KeySpecHMAC512,
	},
}

// CreateKey creates a new customer master key (CMK) in KMS.
// You can specify the key usage, key spec, origin, and tags for the new key.
// If no key usage is specified, ENCRYPT_DECRYPT is used by default.
// If no key spec is specified, SYMMETRIC_DEFAULT is used by default.
func (s *KMSService) CreateKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	bypassLockoutCheck := request.GetBoolParam(req.Parameters, "BypassPolicyLockoutSafetyCheck")

	keyMeta, err := s.createKeyCore(stores, CreateKeyInput{
		Description:        request.GetStringParam(req.Parameters, "Description"),
		KeyUsage:           request.GetStringParam(req.Parameters, "KeyUsage"),
		KeySpec:            request.GetStringParam(req.Parameters, "KeySpec"),
		Origin:             request.GetStringParam(req.Parameters, "Origin"),
		MultiRegion:        request.GetBoolParam(req.Parameters, "MultiRegion"),
		CustomKeyStoreID:   request.GetStringParam(req.Parameters, "CustomKeyStoreId"),
		XksKeyID:           request.GetStringParam(req.Parameters, "XksKeyId"),
		Policy:             request.GetStringParam(req.Parameters, "Policy"),
		BypassLockoutCheck: bypassLockoutCheck,
		Tags:               tagutil.ParseTagsWithKeyNames(req.Parameters, "Tags", "TagKey", "TagValue"),
		AccountID:          reqCtx.GetAccountID(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyMetadata": keyMetadataResultToResponse(keyMeta),
	}, nil
}

// DescribeKey retrieves detailed information about a specified key.
// You can identify the key by its key ID, key ARN, or alias.
func (s *KMSService) DescribeKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "DescribeKey", nil)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyMetadata": s.buildKeyMetadata(key),
	}, nil
}

func (s *KMSService) ListKeys(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	marker := pagination.GetMarker(req.Parameters)
	if err := validateMarkerLength(marker); err != nil {
		return nil, err
	}
	maxItems := pagination.GetMaxItems(req.Parameters, 100)

	result, err := s.listKeysCore(stores, marker, maxItems)
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
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "EnableKey", nil)
	if err != nil {
		return nil, err
	}

	if err := s.enableKeyCore(stores, key); err != nil {
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
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "DisableKey", nil)
	if err != nil {
		return nil, err
	}

	if err := s.disableKeyCore(stores, key); err != nil {
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
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "ScheduleKeyDeletion", nil)
	if err != nil {
		return nil, err
	}

	// PendingWindowInDays: AWS min 7, max 30, default 30. Use existence
	// check to distinguish unset from explicit 0 — an explicit 0 stays 0
	// and is rejected by the core range validation. Range validation and
	// key resolution are performed by scheduleKeyDeletionCore.
	pendingWindowInDays := defaultPendingWindowDays
	if _, ok := req.Parameters["PendingWindowInDays"]; ok {
		pendingWindowInDays = request.GetIntParam(req.Parameters, "PendingWindowInDays")
	}

	updatedKey, days, err := s.scheduleKeyDeletionCore(stores, key.KeyID, pendingWindowInDays)
	if err != nil {
		return nil, err
	}

	deletionDate := int64(0)
	if updatedKey.DeletionDate != nil {
		deletionDate = updatedKey.DeletionDate.Unix()
	}

	return map[string]interface{}{
		"KeyId":               updatedKey.KeyID,
		"DeletionDate":        deletionDate,
		"KeyState":            updatedKey.KeyState,
		"PendingWindowInDays": days,
	}, nil
}

// CancelKeyDeletion cancels the scheduled deletion of a key.
// This reverts the key to its state before ScheduleKeyDeletion was called.
func (s *KMSService) CancelKeyDeletion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "CancelKeyDeletion", nil)
	if err != nil {
		return nil, err
	}

	keyID, keyState, err := s.cancelKeyDeletionCore(stores, key)
	if err != nil {
		return nil, err
	}

	// AWS returns KeyId and KeyState. KeyState reflects the restored state
	// (Enabled or Disabled depending on the pre-ScheduleKeyDeletion value).
	return map[string]interface{}{
		"KeyId":    keyID,
		"KeyState": keyState,
	}, nil
}

// UpdateKeyDescription updates the description of a specified key.
func (s *KMSService) UpdateKeyDescription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "UpdateKeyDescription", nil)
	if err != nil {
		return nil, err
	}

	if err := s.updateKeyDescriptionCore(stores, key, request.GetStringParam(req.Parameters, "Description")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// buildKeyMetadata constructs a KeyMetadata map from a Key store object.
func (s *KMSService) buildKeyMetadata(key *kmsstore.Key) map[string]interface{} {
	_, _, _, accountID, _ := arn.SplitARN(key.Arn)
	metadata := map[string]interface{}{
		"AWSAccountId":          accountID,
		"KeyId":                 key.KeyID,
		"Arn":                   key.Arn,
		"KeyState":              key.KeyState,
		"Enabled":               key.Enabled,
		"KeyUsage":              key.KeyUsage,
		"KeySpec":               key.KeySpec,
		"CustomerMasterKeySpec": key.KeySpec,
		"CreationDate":          key.CreationDate.Unix(),
		"Origin":                key.Origin,
		"KeyManager":            key.KeyManager,
		"MultiRegion":           key.MultiRegion,
	}

	if key.Description != "" {
		metadata["Description"] = key.Description
	}
	if key.DeletionDate != nil {
		metadata["DeletionDate"] = key.DeletionDate.Unix()
	}
	// AWS: PendingDeletionWindowInDays appears in KeyMetadata when the
	// key is scheduled for deletion. The field name in the Smithy model is
	// "PendingDeletionWindowInDays".
	if key.KeyState == kmsstore.KeyStatePendingDeletion && key.PendingWindowInDays > 0 {
		metadata["PendingDeletionWindowInDays"] = key.PendingWindowInDays
	}
	if key.ValidTo != nil {
		metadata["ValidTo"] = key.ValidTo.Unix()
	}
	// AWS: ExpirationModel is present only for keys with imported key
	// material (Origin EXTERNAL). Values: KEY_MATERIAL_EXPIRES or
	// KEY_MATERIAL_DOES_NOT_EXPIRE.
	if key.Origin == kmsstore.OriginTypeExternal && key.ExpirationModel != "" {
		metadata["ExpirationModel"] = key.ExpirationModel
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
	if len(key.MacAlgorithms) > 0 {
		metadata["MacAlgorithms"] = key.MacAlgorithms
	}
	if key.MultiRegion {
		if key.MultiRegionConfiguration != nil {
			mrc := map[string]interface{}{
				"MultiRegionKeyType": key.MultiRegionConfiguration.MultiRegionKeyType,
			}
			if key.MultiRegionConfiguration.PrimaryKey != nil {
				mrc["PrimaryKey"] = map[string]interface{}{
					"Arn":    key.MultiRegionConfiguration.PrimaryKey.Arn,
					"Region": key.MultiRegionConfiguration.PrimaryKey.Region,
				}
			}
			replicas := make([]interface{}, 0, len(key.MultiRegionConfiguration.ReplicaKeys))
			for _, r := range key.MultiRegionConfiguration.ReplicaKeys {
				replicas = append(replicas, map[string]interface{}{
					"Arn":    r.Arn,
					"Region": r.Region,
				})
			}
			mrc["ReplicaKeys"] = replicas
			metadata["MultiRegionConfiguration"] = mrc
		} else {
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
	importToken, publicKey, err := s.getParametersForImportCore(stores, s.resolveCallerPrincipal(reqCtx, req),
		keyID,
		request.GetStringParam(req.Parameters, "WrappingAlgorithm"),
		request.GetStringParam(req.Parameters, "WrappingKeySpec"))
	if err != nil {
		return nil, err
	}

	validTo := time.Now().Add(24 * time.Hour)

	return map[string]interface{}{
		"KeyId":             keyID,
		"ImportToken":       importToken,
		"PublicKey":         publicKey,
		"ParametersValidTo": float64(validTo.Unix()),
	}, nil
}

// ImportKeyMaterial imports key material into a CMK that was created with external origin.
// The key material is encrypted using the public key from GetParametersForImport.
func (s *KMSService) ImportKeyMaterial(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.importKeyMaterialCore(stores, ImportKeyMaterialInput{
		KeyID:                   s.getKeyID(req.Parameters),
		ImportToken:             request.GetStringParam(req.Parameters, "ImportToken"),
		EncryptedKeyMaterialB64: request.GetStringParam(req.Parameters, "EncryptedKeyMaterial"),
		ValidTo:                 request.GetIntParam(req.Parameters, "ValidTo"),
		ExpirationModel:         request.GetStringParam(req.Parameters, "ExpirationModel"),
		Principal:               s.resolveCallerPrincipal(reqCtx, req),
	}); err != nil {
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

	if err := s.deleteImportedKeyMaterialCore(stores, s.resolveCallerPrincipal(reqCtx, req), s.getKeyID(req.Parameters)); err != nil {
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

	replicaRegion := request.GetStringParam(req.Parameters, "ReplicaRegion")

	replicaKey, err := s.replicateKeyCore(stores, ReplicateKeyInput{
		KeyID:         s.getKeyID(req.Parameters),
		ReplicaRegion: replicaRegion,
		Policy:        request.GetStringParam(req.Parameters, "Policy"),
		Description:   request.GetStringParam(req.Parameters, "Description"),
		BypassCheck:   request.GetBoolParam(req.Parameters, "BypassPolicyLockoutSafetyCheck"),
		Tags:          tagutil.ParseTagsWithKeyNames(req.Parameters, "Tags", "TagKey", "TagValue"),
		Principal:     s.resolveCallerPrincipal(reqCtx, req),
		AccountID:     reqCtx.GetAccountID(),
		Region:        reqCtx.GetRegion(),
	})
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

	if err := s.updatePrimaryRegionCore(stores, s.resolveCallerPrincipal(reqCtx, req),
		s.getKeyID(req.Parameters),
		request.GetStringParam(req.Parameters, "PrimaryRegion")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
