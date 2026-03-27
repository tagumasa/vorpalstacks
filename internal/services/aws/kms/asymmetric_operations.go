package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"encoding/base64"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// Sign generates a digital signature for the specified message.
func (s *KMSService) Sign(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "Sign", nil)
	if err != nil {
		return nil, err
	}

	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}

	if key.KeyUsage != kmsstore.KeyUsageSignVerify {
		return nil, ErrInvalidKeyUsage
	}

	messageB64 := request.GetStringParam(req.Parameters, "Message")
	_ = request.GetStringParam(req.Parameters, "MessageType")
	algorithm := request.GetStringParam(req.Parameters, "SigningAlgorithm")

	var message []byte
	message, err = base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		message = []byte(messageB64)
	}

	result, err := s.hsmBackend.Sign(key.KeyID, message, hsm.SigningAlgorithm(algorithm))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyId":            key.Arn,
		"Signature":        base64.StdEncoding.EncodeToString(result.Signature),
		"SigningAlgorithm": algorithm,
	}, nil
}

// Verify verifies a digital signature for the specified message.
func (s *KMSService) Verify(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "Verify", nil)
	if err != nil {
		return nil, err
	}

	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}

	if key.KeyUsage != kmsstore.KeyUsageSignVerify {
		return nil, ErrInvalidKeyUsage
	}

	messageB64 := request.GetStringParam(req.Parameters, "Message")
	algorithm := request.GetStringParam(req.Parameters, "SigningAlgorithm")
	signatureB64 := request.GetStringParam(req.Parameters, "Signature")

	var message []byte
	message, err = base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		message = []byte(messageB64)
	}

	signature, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return nil, ErrInvalidAlgorithm
	}

	valid, err := s.hsmBackend.Verify(key.KeyID, message, signature, hsm.SigningAlgorithm(algorithm))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyId":            key.Arn,
		"SignatureValid":   valid,
		"SigningAlgorithm": algorithm,
	}, nil
}

// GetPublicKey returns the public key for the specified KMS key.
func (s *KMSService) GetPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GetPublicKey", nil)
	if err != nil {
		return nil, err
	}

	publicKey, err := s.hsmBackend.GetPublicKey(key.KeyID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyId":     key.Arn,
		"PublicKey": base64.StdEncoding.EncodeToString(publicKey),
		"KeySpec":   key.KeySpec,
		"KeyUsage":  key.KeyUsage,
	}, nil
}
