package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// EnableKeyRotation enables automatic key rotation for a symmetric key.
// Key rotation automatically generates new cryptographic material every year.
func (s *KMSService) EnableKeyRotation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "EnableKeyRotation", nil)
	if err != nil {
		return nil, err
	}

	if err := stores.keys.SetKeyRotation(key.KeyID, true); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DisableKeyRotation disables automatic key rotation for a symmetric key.
func (s *KMSService) DisableKeyRotation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "DisableKeyRotation", nil)
	if err != nil {
		return nil, err
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
