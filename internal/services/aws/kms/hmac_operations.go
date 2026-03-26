package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"encoding/base64"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// GenerateMac generates a MAC (Message Authentication Code) for the specified message.
func (s *KMSService) GenerateMac(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}

	if key.KeyUsage != kmsstore.KeyUsageGenerateVerifyMAC {
		return nil, ErrInvalidKeyUsage
	}

	messageB64 := request.GetStringParam(req.Parameters, "Message")
	message, err := base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		message = []byte(messageB64)
	}

	algorithm := request.GetStringParam(req.Parameters, "MacAlgorithm")

	macValue, err := s.hsmBackend.GenerateMAC(key.KeyID, message, hsm.MACAlgorithm(algorithm))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyId":        key.Arn,
		"Mac":          base64.StdEncoding.EncodeToString(macValue),
		"MacAlgorithm": algorithm,
	}, nil
}

// VerifyMac verifies a MAC (Message Authentication Code) for the specified message.
func (s *KMSService) VerifyMac(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveKey(stores, req.Parameters)
	if err != nil {
		return nil, err
	}

	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}

	if key.KeyUsage != kmsstore.KeyUsageGenerateVerifyMAC {
		return nil, ErrInvalidKeyUsage
	}

	messageB64 := request.GetStringParam(req.Parameters, "Message")
	message, err := base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		message = []byte(messageB64)
	}

	macB64 := request.GetStringParam(req.Parameters, "Mac")
	macValue, err := base64.StdEncoding.DecodeString(macB64)
	if err != nil {
		return nil, ErrInvalidAlgorithm
	}

	algorithm := request.GetStringParam(req.Parameters, "MacAlgorithm")

	valid, err := s.hsmBackend.VerifyMAC(key.KeyID, message, macValue, hsm.MACAlgorithm(algorithm))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"KeyId":        key.Arn,
		"MacValid":     valid,
		"MacAlgorithm": algorithm,
	}, nil
}
