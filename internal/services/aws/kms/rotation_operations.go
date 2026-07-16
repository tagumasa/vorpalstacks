package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// EnableKeyRotation enables automatic key rotation for a symmetric key.
// Key rotation automatically generates new cryptographic material every year.
// Per AWS spec, automatic rotation is supported only on SYMMETRIC_DEFAULT keys
// with AWS_KMS origin. The key must be enabled and not pending deletion or import.
func (s *KMSService) EnableKeyRotation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "EnableKeyRotation", nil)
	if err != nil {
		return nil, err
	}

	if key.KeyState == kmsstore.KeyStatePendingDeletion || key.KeyState == kmsstore.KeyStatePendingImport {
		return nil, ErrKeyPendingDeletion
	}
	if key.KeyState == kmsstore.KeyStateDisabled || !key.Enabled {
		return nil, ErrKeyDisabled
	}
	if key.KeySpec != kmsstore.KeySpecSymmetricDefault || key.Origin != kmsstore.OriginTypeAWSKMS {
		return nil, ErrUnsupportedOperation
	}

	if err := stores.keys.SetKeyRotation(key.KeyID, true); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DisableKeyRotation disables automatic key rotation for a symmetric key.
// Per AWS spec, rotation status cannot be changed while the key is pending
// deletion or disabled.
func (s *KMSService) DisableKeyRotation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "DisableKeyRotation", nil)
	if err != nil {
		return nil, err
	}

	if key.KeyState == kmsstore.KeyStatePendingDeletion || key.KeyState == kmsstore.KeyStatePendingImport {
		return nil, ErrKeyPendingDeletion
	}
	if key.KeyState == kmsstore.KeyStateDisabled || !key.Enabled {
		return nil, ErrKeyDisabled
	}

	if err := stores.keys.SetKeyRotation(key.KeyID, false); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetKeyRotationStatus returns the status of key rotation for a specified key.
// It indicates whether automatic key rotation is enabled and the next rotation date.
func (s *KMSService) GetKeyRotationStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GetKeyRotationStatus", nil)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"KeyRotationEnabled": key.KeyRotationEnabled,
		"KeyId":              key.Arn,
	}

	if key.KeyRotationEnabled && !key.CreationDate.IsZero() {
		nextRotation := key.CreationDate.AddDate(1, 0, 0)
		response["NextRotationDate"] = nextRotation.Unix()
		response["RotationPeriodInDays"] = int32(365)
	}

	return response, nil
}
