package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"encoding/base64"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// GenerateMac generates a MAC (Message Authentication Code) for the specified message.

func macAlgorithmSupported(algorithm string, supported []string) bool {
	for _, a := range supported {
		if a == algorithm {
			return true
		}
	}
	return false
}
func (s *KMSService) GenerateMac(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GenerateMac", nil)
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
	if messageB64 == "" {
		// AWS rejects empty Message with ValidationException.
		return nil, ErrValidation
	}
	// AWS requires base64-encoded Message. The previous code fell back to
	// []byte(messageB64) which silently accepted non-base64 input and
	// computed a MAC over the raw string, diverging from AWS.
	message, err := base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		return nil, ErrValidation
	}

	algorithm := request.GetStringParam(req.Parameters, "MacAlgorithm")
	if algorithm == "" {
		return nil, ErrValidation
	}
	if !macAlgorithmSupported(algorithm, key.MacAlgorithms) {
		return nil, ErrInvalidAlgorithm
	}

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
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "VerifyMac", nil)
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
	if messageB64 == "" {
		return nil, ErrValidation
	}
	macB64 := request.GetStringParam(req.Parameters, "Mac")
	if macB64 == "" {
		return nil, ErrValidation
	}

	message, err := base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		// AWS requires base64 Message; non-base64 is a validation error.
		return nil, ErrValidation
	}
	macValue, err := base64.StdEncoding.DecodeString(macB64)
	if err != nil {
		// Previous code returned ErrInvalidAlgorithm for malformed Mac
		// which conflated the input validation failure with the
		// algorithm check. ValidationException is the correct AWS error.
		return nil, ErrValidation
	}

	algorithm := request.GetStringParam(req.Parameters, "MacAlgorithm")
	if algorithm == "" {
		return nil, ErrValidation
	}
	if !macAlgorithmSupported(algorithm, key.MacAlgorithms) {
		return nil, ErrInvalidAlgorithm
	}

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
