package kms

// Package kms: asymmetric_core.go holds the Sign/Verify execution path:
// request-contract validation, key authorisation, and signature production
// and verification.

import (
	"encoding/base64"
	"errors"
	"fmt"

	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// SignInput is the transport-agnostic input for Sign.
type SignInput struct {
	KeyID            string
	CallerPrincipal  string
	MessageB64       string
	MessageType      string
	SigningAlgorithm string
	DryRun           bool
}

// VerifyInput is the transport-agnostic input for Verify.
type VerifyInput struct {
	KeyID            string
	CallerPrincipal  string
	MessageB64       string
	MessageType      string
	SigningAlgorithm string
	SignatureB64     string
	DryRun           bool
}

// authorizeSignVerifyKey resolves and authorises the key for a Sign or
// Verify operation and enforces the key state and SIGN_VERIFY usage
// requirements shared by both operations.
func (s *KMSService) authorizeSignVerifyKey(stores *kmsStores, keyID, callerPrincipal, action string) (*kmsstore.Key, error) {
	key, err := s.resolveKeyByKeyID(stores, keyID)
	if err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, callerPrincipal, action, key.KeyID, nil); err != nil {
		return nil, err
	}
	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}
	if key.KeyUsage != kmsstore.KeyUsageSignVerify {
		return nil, ErrInvalidKeyUsage
	}
	return key, nil
}

// signCore executes the Sign operation: key resolution, authorisation,
// request-contract validation, and signature production.
func (s *KMSService) signCore(stores *kmsStores, in *SignInput) (map[string]interface{}, error) {
	key, err := s.authorizeSignVerifyKey(stores, in.KeyID, in.CallerPrincipal, "Sign")
	if err != nil {
		return nil, err
	}

	messageType, err := validateSignInputCore(key, in.MessageB64, in.MessageType, in.SigningAlgorithm)
	if err != nil {
		return nil, err
	}

	message, err := base64.StdEncoding.DecodeString(in.MessageB64)
	if err != nil {
		// AWS requires base64-encoded Message; non-base64 input is a
		// validation error, not a signature failure.
		return nil, NewValidationError("Message is not valid base64")
	}

	if in.DryRun {
		return nil, ErrDryRunOperation
	}

	result, err := s.hsmBackend.Sign(key.KeyID, message, hsm.SigningAlgorithm(in.SigningAlgorithm), hsm.MessageType(messageType))
	if err != nil {
		if errors.Is(err, hsm.ErrInvalidDigestLength) {
			return nil, NewValidationError(fmt.Sprintf("Digest length does not match %s", in.SigningAlgorithm))
		}
		return nil, err
	}
	s.markKeyLastUsed(stores, key.KeyID, "Sign")

	return map[string]interface{}{
		"KeyId":            key.Arn,
		"Signature":        base64.StdEncoding.EncodeToString(result.Signature),
		"SigningAlgorithm": in.SigningAlgorithm,
	}, nil
}

// verifyCore executes the Verify operation: key resolution, authorisation,
// request-contract validation, and signature verification.
func (s *KMSService) verifyCore(stores *kmsStores, in *VerifyInput) (map[string]interface{}, error) {
	key, err := s.authorizeSignVerifyKey(stores, in.KeyID, in.CallerPrincipal, "Verify")
	if err != nil {
		return nil, err
	}

	messageType, err := validateVerifyInputCore(key, in.MessageB64, in.MessageType, in.SigningAlgorithm, in.SignatureB64)
	if err != nil {
		return nil, err
	}

	message, err := base64.StdEncoding.DecodeString(in.MessageB64)
	if err != nil {
		// Non-base64 Message is a validation error, not a key-material error.
		return nil, NewValidationError("Message is not valid base64")
	}

	signature, err := base64.StdEncoding.DecodeString(in.SignatureB64)
	if err != nil {
		// A malformed Signature is an input validation failure, not an
		// algorithm-contract violation; ValidationException is the AWS error.
		return nil, NewValidationError("Signature is not valid base64")
	}

	if in.DryRun {
		return nil, ErrDryRunOperation
	}

	valid, err := s.hsmBackend.Verify(key.KeyID, message, signature, hsm.SigningAlgorithm(in.SigningAlgorithm), hsm.MessageType(messageType))
	if err != nil {
		if errors.Is(err, hsm.ErrInvalidDigestLength) {
			return nil, ErrValidation
		}
		return nil, err
	}
	// AWS: when the signature does not match the message, Verify fails
	// with KMSInvalidSignatureException. The operation never returns a
	// success response with SignatureValid=false.
	if !valid {
		return nil, ErrKMSInvalidSignature
	}
	s.markKeyLastUsed(stores, key.KeyID, "Verify")

	return map[string]interface{}{
		"KeyId":            key.Arn,
		"SignatureValid":   true,
		"SigningAlgorithm": in.SigningAlgorithm,
	}, nil
}

// validateSignInputCore enforces the Sign request contract: Message and
// SigningAlgorithm are required, MessageType defaults to RAW and accepts
// only RAW or DIGEST, and the algorithm must match the key spec. It
// returns the effective MessageType.
func validateSignInputCore(key *kmsstore.Key, messageB64, messageType, algorithm string) (string, error) {
	if messageB64 == "" {
		// AWS rejects empty Message with ValidationException.
		return "", NewValidationError("Message is required")
	}
	if messageType == "" {
		messageType = string(hsm.MessageTypeRaw)
	}
	if messageType != string(hsm.MessageTypeRaw) && messageType != string(hsm.MessageTypeDigest) {
		return "", NewValidationError(fmt.Sprintf("MessageType must be RAW or DIGEST, got %q", messageType))
	}
	if algorithm == "" {
		return "", NewValidationError("SigningAlgorithm is required")
	}
	if err := resolveSigningAlgorithm(algorithm, key); err != nil {
		return "", err
	}
	return messageType, nil
}

// validateVerifyInputCore enforces the Verify request contract on top of
// the Sign contract: Signature is additionally required.
func validateVerifyInputCore(key *kmsstore.Key, messageB64, messageType, algorithm, signatureB64 string) (string, error) {
	effective, err := validateSignInputCore(key, messageB64, messageType, algorithm)
	if err != nil {
		return "", err
	}
	if signatureB64 == "" {
		return "", NewValidationError("Signature is required")
	}
	return effective, nil
}
