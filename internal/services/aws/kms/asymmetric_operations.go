package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"encoding/base64"

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
	return s.signCore(stores, &SignInput{
		KeyID:            s.getKeyID(req.Parameters),
		CallerPrincipal:  s.resolveCallerPrincipal(reqCtx, req),
		MessageB64:       request.GetStringParam(req.Parameters, "Message"),
		MessageType:      request.GetStringParam(req.Parameters, "MessageType"),
		SigningAlgorithm: request.GetStringParam(req.Parameters, "SigningAlgorithm"),
		DryRun:           getDryRunParam(req.Parameters),
	})
}

// Verify verifies a digital signature for the specified message.
func (s *KMSService) Verify(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.verifyCore(stores, &VerifyInput{
		KeyID:            s.getKeyID(req.Parameters),
		CallerPrincipal:  s.resolveCallerPrincipal(reqCtx, req),
		MessageB64:       request.GetStringParam(req.Parameters, "Message"),
		MessageType:      request.GetStringParam(req.Parameters, "MessageType"),
		SigningAlgorithm: request.GetStringParam(req.Parameters, "SigningAlgorithm"),
		SignatureB64:     request.GetStringParam(req.Parameters, "Signature"),
		DryRun:           getDryRunParam(req.Parameters),
	})
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

// validSigningAlgorithmSpec reports whether alg is a member of the Smithy
// SigningAlgorithmSpec enum that types the SigningAlgorithm request member
// of Sign and Verify. SM2DSA belongs to the enum even though this platform
// has no SM2 key material; a request naming it fails the key-level
// algorithm support check.
func validSigningAlgorithmSpec(alg string) bool {
	switch alg {
	case string(hsm.SigningAlgorithmRSAPKCS1SHA256),
		string(hsm.SigningAlgorithmRSAPKCS1SHA384),
		string(hsm.SigningAlgorithmRSAPKCS1SHA512),
		string(hsm.SigningAlgorithmRSAPSSSHA256),
		string(hsm.SigningAlgorithmRSAPSSSHA384),
		string(hsm.SigningAlgorithmRSAPSSSHA512),
		string(hsm.SigningAlgorithmECDSASHA256),
		string(hsm.SigningAlgorithmECDSASHA384),
		string(hsm.SigningAlgorithmECDSASHA512),
		"SM2DSA":
		return true
	}
	return false
}

// resolveSigningAlgorithm applies the SigningAlgorithm member contract: an
// explicit value outside the SigningAlgorithmSpec enum is a shape violation
// rejected with SerializationException, and a value the key does not support
// is rejected with InvalidKeyUsageException, whose documentation covers "the
// signing algorithm specified for the operation is incompatible with the
// type of key material in the KMS key".
func resolveSigningAlgorithm(algorithm string, key *kmsstore.Key) error {
	if !validSigningAlgorithmSpec(algorithm) {
		return ErrSerializationException
	}
	if !algorithmSupported(algorithm, key.SigningAlgorithms) {
		return ErrInvalidKeyUsage
	}
	return nil
}
