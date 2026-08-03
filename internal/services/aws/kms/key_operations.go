package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/kms/hsm"
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

	bypassLockoutCheck := false
	if v, ok := req.Parameters["BypassPolicyLockoutSafetyCheck"]; ok {
		if b, ok := v.(string); ok && b == "true" {
			bypassLockoutCheck = true
		}
	}

	key, err := s.createKeyCore(stores, CreateKeyInput{
		Description:        request.GetStringParam(req.Parameters, "Description"),
		KeyUsage:           kmsstore.KeyUsage(request.GetStringParam(req.Parameters, "KeyUsage")),
		KeySpec:            kmsstore.KeySpec(request.GetStringParam(req.Parameters, "KeySpec")),
		Origin:             kmsstore.OriginType(request.GetStringParam(req.Parameters, "Origin")),
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
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "DescribeKey", nil)
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
	if err := validateMarkerLength(marker); err != nil {
		return nil, err
	}
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
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "EnableKey", nil)
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
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "DisableKey", nil)
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
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "ScheduleKeyDeletion", nil)
	if err != nil {
		return nil, err
	}

	// PendingWindowInDays: AWS min 7, max 30, default 30. Use existence
	// check to distinguish unset from explicit 0. Range validation is
	// performed by scheduleKeyDeletionCore.
	pendingWindowInDays := 30
	if _, ok := req.Parameters["PendingWindowInDays"]; ok {
		pendingWindowInDays = request.GetIntParam(req.Parameters, "PendingWindowInDays")
	}

	updatedKey, days, err := s.scheduleKeyDeletionCore(stores, key.KeyID, pendingWindowInDays)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyId":               key.KeyID,
		"DeletionDate":        updatedKey.DeletionDate.Unix(),
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

	if err := stores.keys.CancelDeletion(key.KeyID); err != nil {
		return nil, err
	}

	updatedKey, err := stores.keys.Get(key.KeyID)
	if err != nil || updatedKey == nil {
		return nil, fmt.Errorf("failed to retrieve key after cancelling deletion: %w", err)
	}

	// AWS returns KeyId and KeyState. KeyState reflects the restored state
	// (Enabled or Disabled depending on the pre-ScheduleKeyDeletion value).
	return map[string]interface{}{
		"KeyId":    key.KeyID,
		"KeyState": updatedKey.KeyState,
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

	description := request.GetStringParam(req.Parameters, "Description")
	if err := validateDescriptionLength(description); err != nil {
		return nil, err
	}

	if err := stores.keys.UpdateDescription(key.KeyID, description); err != nil {
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
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "GetParametersForImport", keyID, nil); err != nil {
		return nil, err
	}

	// AWS: GetParametersForImport requires the key to be in PendingImport
	// state. Calling it on an Enabled key or a key with already-imported
	// material returns KMSInvalidStateException.
	key, err := stores.keys.Get(keyID)
	if err != nil {
		return nil, NewKeyNotFoundError(keyID)
	}
	if key.Origin != kmsstore.OriginTypeExternal {
		return nil, ErrUnsupportedOperation
	}
	if key.KeyState != kmsstore.KeyStatePendingImport {
		return nil, ErrKeyPendingImport
	}

	// AWS: WrappingAlgorithm selects between symmetric and RSA wrapping.
	// The platform implements RSA_* wrapping only; SYMMETRIC_KEY_WRAPPING
	// is rejected as unsupported until the underlying implementation lands.
	wrappingAlgorithm := request.GetStringParam(req.Parameters, "WrappingAlgorithm")
	switch wrappingAlgorithm {
	case "":
		wrappingAlgorithm = "RSAES_OAEP_SHA_256"
	case "RSAES_OAEP_SHA_256", "RSAES_OAEP_SHA_1":
		// supported
	default:
		return nil, ErrUnsupportedOperation
	}

	// AWS: WrappingKeySpec is a required parameter (smithy.api#required
	// trait on the GetParametersForImport input). The previous code
	// silently defaulted to RSA_2048 when omitted, diverging from the
	// documented contract.
	wrappingKeySpec := request.GetStringParam(req.Parameters, "WrappingKeySpec")
	switch wrappingKeySpec {
	case "RSA_2048", "RSA_4096":
		// supported
	default:
		return nil, NewValidationError(fmt.Sprintf("WrappingKeySpec must be RSA_2048 or RSA_4096, got %q", wrappingKeySpec))
	}

	importToken, publicKey, err := stores.keys.GetParametersForImport(keyID, wrappingKeySpec, wrappingAlgorithm)
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

	keyID := s.getKeyID(req.Parameters)
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "ImportKeyMaterial", keyID, nil); err != nil {
		return nil, err
	}

	key, err := stores.keys.Get(keyID)
	if err != nil {
		return nil, err
	}

	importToken := request.GetStringParam(req.Parameters, "ImportToken")
	encryptedKeyMaterialB64 := request.GetStringParam(req.Parameters, "EncryptedKeyMaterial")
	validTo := request.GetIntParam(req.Parameters, "ValidTo")
	expirationModel := request.GetStringParam(req.Parameters, "ExpirationModel")

	// Enforce AWS-documented length constraints (max 6144
	// base64-encoded characters) before decoding to prevent memory
	// exhaustion from arbitrarily large blobs.
	if err := validateImportTokenSize(importToken); err != nil {
		return nil, err
	}
	if err := validateEncryptedKeyMaterialSize(encryptedKeyMaterialB64); err != nil {
		return nil, err
	}

	encryptedKeyMaterial, err := base64.StdEncoding.DecodeString(encryptedKeyMaterialB64)
	if err != nil {
		return nil, fmt.Errorf("invalid EncryptedKeyMaterial: %w", err)
	}

	// Reject expired import tokens. GetParametersForImport issues tokens
	// with a 24-hour validity window; using an expired token returns
	// ExpiredImportTokenException per the AWS contract.
	if key.ImportTokenValidTo != nil && time.Now().After(*key.ImportTokenValidTo) {
		return nil, ErrExpiredImportToken
	}

	// AWS: ExpirationModel controls whether the imported key material
	// expires. KEY_MATERIAL_DOES_NOT_EXPIRE causes ValidTo to be ignored.
	// When ExpirationModel is omitted but ValidTo is supplied, the default
	// is KEY_MATERIAL_EXPIRES. When neither is supplied, the default is
	// KEY_MATERIAL_DOES_NOT_EXPIRE.
	var validToTime *time.Time
	if expirationModel != "KEY_MATERIAL_DOES_NOT_EXPIRE" && validTo > 0 {
		t := time.Unix(int64(validTo), 0)
		validToTime = &t
	}
	if expirationModel == "" {
		if validToTime != nil {
			expirationModel = "KEY_MATERIAL_EXPIRES"
		} else {
			expirationModel = "KEY_MATERIAL_DOES_NOT_EXPIRE"
		}
	}

	rawKeyMaterial, err := stores.keys.ImportKeyMaterial(keyID, importToken, encryptedKeyMaterial, validToTime)
	if err != nil {
		return nil, err
	}

	if err := s.hsmBackend.ImportKey(keyID, rawKeyMaterial, hsm.KeySpec(key.KeySpec)); err != nil {
		return nil, err
	}

	// Re-fetch the key because ImportKeyMaterial's atomicUpdate changed
	// its state (PendingImport → Enabled) and we must not overwrite that
	// with the stale pre-import copy.
	updatedKey, err := stores.keys.Get(keyID)
	if err != nil {
		return nil, err
	}
	updatedKey.ExpirationModel = expirationModel
	if err := stores.keys.Update(updatedKey); err != nil {
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

	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "DeleteImportedKeyMaterial", keyID, nil); err != nil {
		return nil, err
	}

	// Delete from the HSM first so that if it fails, the store still
	// reflects the imported state (Enabled + key material present).
	// The previous ordering (store → HSM) left an inconsistent state
	// when the HSM delete failed: the store showed PendingImport but
	// key material still resided in the HSM.
	if err := s.hsmBackend.DeleteKey(keyID); err != nil {
		logs.Error("DeleteImportedKeyMaterial: HSM DeleteKey failed", logs.String("keyId", keyID), logs.Err(err))
		return nil, ErrKMSInternal
	}
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
	if replicaRegion == "" {
		return nil, NewValidationError("ReplicaRegion is required")
	}
	if err := validateReplicaRegion(replicaRegion); err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "ReplicateKey", keyID, nil); err != nil {
		return nil, err
	}

	primary, err := stores.keys.Get(keyID)
	if err != nil {
		return nil, NewKeyNotFoundError(keyID)
	}
	if !primary.MultiRegion {
		return nil, kmsstore.ErrNotMultiRegionKey
	}

	replicaKeyID, err := kmsstore.GenerateKeyID()
	if err != nil {
		return nil, err
	}

	replicaKey, err := stores.keys.ReplicateKey(keyID, replicaRegion, replicaKeyID)
	if err != nil {
		return nil, err
	}

	// Apply optional caller-supplied parameters to the replica. Per the
	// Smithy ReplicateKeyRequest, the caller may set Policy, Description,
	// BypassPolicyLockoutSafetyCheck, and Tags on the replica.
	replicaPolicy := request.GetStringParam(req.Parameters, "Policy")
	if replicaPolicy == "" {
		// Inherit the primary's key policy when not explicitly supplied.
		primaryPolicy, pErr := stores.keyPolicies.GetDefault(primary.KeyID)
		if pErr == nil {
			replicaPolicy = primaryPolicy.PolicyDocument
		} else {
			replicaPolicy = kmsstore.DefaultKeyPolicy
		}
	}
	bypassCheck := false
	if v, ok := req.Parameters["BypassPolicyLockoutSafetyCheck"]; ok {
		if b, ok := v.(string); ok && b == "true" {
			bypassCheck = true
		}
	}
	if err := validatePolicySize(replicaPolicy); err != nil {
		return nil, err
	}
	if !bypassCheck {
		if err := validatePolicyDoesNotLockOutRoot(replicaPolicy, reqCtx.GetAccountID()); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, replicaKeyID); delErr != nil {
				logs.Error("Failed to cascade-delete replica after policy lockout check", logs.Err(delErr), logs.String("replicaKeyId", replicaKeyID))
			}
			return nil, err
		}
	}
	if desc := request.GetStringParam(req.Parameters, "Description"); desc != "" {
		if err := validateDescriptionLength(desc); err != nil {
			return nil, err
		}
		replicaKey.Description = desc
	}
	replicaKey.BypassPolicyLockoutSafetyCheck = bypassCheck
	if err := stores.keys.Update(replicaKey); err != nil {
		if delErr := stores.CascadeDeleteKey(s.hsmBackend, replicaKeyID); delErr != nil {
			logs.Error("Failed to cascade-delete replica after Update", logs.Err(delErr), logs.String("replicaKeyId", replicaKeyID))
		}
		return nil, err
	}

	// Apply tags from the request.
	replicaTags := tagutil.ParseTagsWithKeyNames(req.Parameters, "Tags", "TagKey", "TagValue")
	if len(replicaTags) > 0 {
		if err := validateKMSTags(replicaTags); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, replicaKeyID); delErr != nil {
				logs.Error("Failed to cascade-delete replica after tag validation failure", logs.Err(delErr), logs.String("replicaKeyId", replicaKeyID))
			}
			return nil, err
		}
		if err := stores.keys.TagStore.Tag(replicaKeyID, tagutil.ToMap(replicaTags)); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, replicaKeyID); delErr != nil {
				logs.Error("Failed to cascade-delete replica after Tag failure", logs.Err(delErr), logs.String("replicaKeyId", replicaKeyID))
			}
			return nil, err
		}
	}

	// Store the replica's key policy.
	if err := stores.keyPolicies.PutDefault(replicaKeyID, replicaPolicy); err != nil {
		if delErr := stores.CascadeDeleteKey(s.hsmBackend, replicaKeyID); delErr != nil {
			logs.Error("Failed to cascade-delete replica after PutDefault policy", logs.Err(delErr), logs.String("replicaKeyId", replicaKeyID))
		}
		return nil, err
	}

	// Provision the replica in the HSM with the primary's key material.
	// Without this, crypto operations on the replica fail with
	// ErrKeyNotFound because the HSM has no material for the new keyID.
	if err := s.hsmBackend.ReplicateKey(primary.KeyID, replicaKeyID); err != nil {
		if delErr := stores.CascadeDeleteKey(s.hsmBackend, replicaKeyID); delErr != nil {
			logs.Error("Failed to cascade-delete replica after HSM ReplicateKey failure", logs.Err(delErr), logs.String("replicaKeyId", replicaKeyID))
		}
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
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "UpdatePrimaryRegion", keyID, nil); err != nil {
		return nil, err
	}

	// Reject empty PrimaryRegion before it reaches the store.
	// Without this check the store persists a broken PrimaryKeyInfo
	// with an empty Region field.
	if err := validatePrimaryRegion(primaryRegion); err != nil {
		return nil, err
	}

	if err := stores.keys.UpdatePrimaryRegion(keyID, primaryRegion); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
