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

	encryptionAlgorithms, signingAlgorithms, keyAgreementAlgorithms := getAlgorithmSets(key.KeySpec, key.KeyUsage)

	result := map[string]interface{}{
		"KeyId":     key.Arn,
		"PublicKey": base64.StdEncoding.EncodeToString(publicKey),
		"KeySpec":   key.KeySpec,
		"KeyUsage":  key.KeyUsage,
	}
	if len(encryptionAlgorithms) > 0 {
		result["EncryptionAlgorithms"] = encryptionAlgorithms
	}
	if len(signingAlgorithms) > 0 {
		result["SigningAlgorithms"] = signingAlgorithms
	}
	if len(keyAgreementAlgorithms) > 0 {
		result["KeyAgreementAlgorithms"] = keyAgreementAlgorithms
	}

	return result, nil
}

// getAlgorithmSets returns the default algorithm sets for a given KeySpec and KeyUsage.
func getAlgorithmSets(keySpec kmsstore.KeySpec, keyUsage kmsstore.KeyUsage) (encryptionAlgorithms []string, signingAlgorithms []string, keyAgreementAlgorithms []string) {
	switch {
	case keyUsage == kmsstore.KeyUsageSignVerify:
		switch keySpec {
		case kmsstore.KeySpecRSA2048, kmsstore.KeySpecRSA3072, kmsstore.KeySpecRSA4096:
			signingAlgorithms = []string{
				"RSASSA_PKCS1_V1_5_SHA_256",
				"RSASSA_PKCS1_V1_5_SHA_384",
				"RSASSA_PKCS1_V1_5_SHA_512",
				"RSASSA_PSS_SHA_256",
				"RSASSA_PSS_SHA_384",
				"RSASSA_PSS_SHA_512",
			}
		case kmsstore.KeySpecECCNISTP256, kmsstore.KeySpecECCSECGP256K1:
			signingAlgorithms = []string{"ECDSA_SHA_256"}
		case kmsstore.KeySpecECCNISTP384:
			signingAlgorithms = []string{"ECDSA_SHA_384"}
		case kmsstore.KeySpecECCNISTP521:
			signingAlgorithms = []string{"ECDSA_SHA_512"}
		case kmsstore.KeySpecSM2:
			signingAlgorithms = []string{"SM2"}
		}
	case keyUsage == kmsstore.KeyUsageEncryptDecrypt && keySpec != kmsstore.KeySpecSymmetricDefault:
		switch keySpec {
		case kmsstore.KeySpecRSA2048, kmsstore.KeySpecRSA3072, kmsstore.KeySpecRSA4096:
			encryptionAlgorithms = []string{
				"RSAES_OAEP_SHA_256",
				"RSAES_OAEP_SHA_1",
			}
		case kmsstore.KeySpecSM2:
			encryptionAlgorithms = []string{"SM2PKE"}
		}
	}
	return
}
