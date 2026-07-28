package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"vorpalstacks/internal/common/request"
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
	if messageB64 == "" {
		// AWS rejects empty Message with ValidationException.
		return nil, NewValidationError("Message is required")
	}
	messageType := request.GetStringParam(req.Parameters, "MessageType")
	if messageType == "" {
		messageType = string(hsm.MessageTypeRaw)
	}
	if messageType != string(hsm.MessageTypeRaw) && messageType != string(hsm.MessageTypeDigest) {
		return nil, NewValidationError(fmt.Sprintf("MessageType must be RAW or DIGEST, got %q", messageType))
	}
	algorithm := request.GetStringParam(req.Parameters, "SigningAlgorithm")
	if algorithm == "" {
		return nil, NewValidationError("SigningAlgorithm is required")
	}
	if !algorithmSupported(algorithm, key.SigningAlgorithms) {
		return nil, ErrInvalidAlgorithm
	}

	message, err := base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		// AWS requires base64-encoded Message; non-base64 input is a
		// validation error, not a signature failure.
		return nil, NewValidationError("Message is not valid base64")
	}

	if err := checkKMSDryRun(req.Parameters); err != nil {
		return nil, err
	}

	result, err := s.hsmBackend.Sign(key.KeyID, message, hsm.SigningAlgorithm(algorithm), hsm.MessageType(messageType))
	if err != nil {
		if errors.Is(err, hsm.ErrInvalidDigestLength) {
			return nil, NewValidationError(fmt.Sprintf("Digest length does not match %s", algorithm))
		}
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
	if messageB64 == "" {
		return nil, NewValidationError("Message is required")
	}
	messageType := request.GetStringParam(req.Parameters, "MessageType")
	if messageType == "" {
		messageType = string(hsm.MessageTypeRaw)
	}
	if messageType != string(hsm.MessageTypeRaw) && messageType != string(hsm.MessageTypeDigest) {
		return nil, NewValidationError(fmt.Sprintf("MessageType must be RAW or DIGEST, got %q", messageType))
	}
	algorithm := request.GetStringParam(req.Parameters, "SigningAlgorithm")
	if algorithm == "" {
		return nil, NewValidationError("SigningAlgorithm is required")
	}
	if !algorithmSupported(algorithm, key.SigningAlgorithms) {
		return nil, ErrInvalidAlgorithm
	}
	signatureB64 := request.GetStringParam(req.Parameters, "Signature")
	if signatureB64 == "" {
		return nil, NewValidationError("Signature is required")
	}

	message, err := base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		// Non-base64 Message is a validation error, not a key-material error.
		return nil, NewValidationError("Message is not valid base64")
	}

	signature, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		// Signature was malformed; surface as ValidationException rather
		// than the misleading ErrInvalidAlgorithm the previous code returned.
		return nil, NewValidationError("Signature is not valid base64")
	}

	if err := checkKMSDryRun(req.Parameters); err != nil {
		return nil, err
	}

	valid, err := s.hsmBackend.Verify(key.KeyID, message, signature, hsm.SigningAlgorithm(algorithm), hsm.MessageType(messageType))
	if err != nil {
		if errors.Is(err, hsm.ErrInvalidDigestLength) {
			return nil, ErrValidation
		}
		return nil, err
	}

	return map[string]interface{}{
		"KeyId":            key.Arn,
		"SignatureValid":   valid,
		"SigningAlgorithm": algorithm,
	}, nil
}

// GetPublicKey returns the public key for the specified KMS key.
// AWS supports GetPublicKey only on asymmetric key specs (RSA, ECC, SM2,
// and HMAC keys expose the public material as a symmetric blob). A
// SYMMETRIC_DEFAULT key has no exportable public component and AWS
// returns UnsupportedOperationException.
func (s *KMSService) GetPublicKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GetPublicKey", nil)
	if err != nil {
		return nil, err
	}

	if key.KeySpec == kmsstore.KeySpecSymmetricDefault {
		return nil, ErrUnsupportedOperation
	}

	publicKey, err := s.hsmBackend.GetPublicKey(key.KeyID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"KeyId":     key.Arn,
		"PublicKey": base64.StdEncoding.EncodeToString(publicKey),
		"KeySpec":   key.KeySpec,
		"KeyUsage":  key.KeyUsage,
	}
	if len(key.EncryptionAlgorithms) > 0 {
		result["EncryptionAlgorithms"] = key.EncryptionAlgorithms
	}
	if len(key.SigningAlgorithms) > 0 {
		result["SigningAlgorithms"] = key.SigningAlgorithms
	}
	if len(key.MacAlgorithms) > 0 {
		result["MacAlgorithms"] = key.MacAlgorithms
	}

	// CustomerMasterKeySpec is the deprecated alias of KeySpec; AWS still
	// returns it for backward compatibility.
	result["CustomerMasterKeySpec"] = key.KeySpec

	return result, nil
}

func algorithmSupported(algorithm string, supported []string) bool {
	for _, a := range supported {
		if a == algorithm {
			return true
		}
	}
	return false
}
