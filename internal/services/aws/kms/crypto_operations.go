package kms

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"

	"vorpalstacks/internal/services/aws/common/pagination"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// Encrypt encrypts plaintext using the specified KMS key.
func (s *KMSService) Encrypt(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	if key.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
		return nil, ErrInvalidKeyUsage
	}

	plaintextB64 := request.GetStringParam(req.Parameters, "Plaintext")
	plaintext, err := base64.StdEncoding.DecodeString(plaintextB64)
	if err != nil {
		plaintext, err = base64.RawStdEncoding.DecodeString(plaintextB64)
		if err != nil {
			return nil, ErrValidation
		}
	}

	encryptionContext := parseEncryptionContext(req.Parameters)
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "Encrypt", key.KeyID, encryptionContext); err != nil {
		return nil, err
	}

	result, err := s.hsmBackend.Encrypt(key.KeyID, plaintext, encryptionContext)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CiphertextBlob":      base64.StdEncoding.EncodeToString(result.Ciphertext),
		"KeyId":               key.Arn,
		"EncryptionAlgorithm": "SYMMETRIC_DEFAULT",
	}, nil
}

// Decrypt decrypts ciphertext using the specified KMS key.
func (s *KMSService) Decrypt(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ciphertextB64 := request.GetStringParam(req.Parameters, "CiphertextBlob")
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	encryptionContext := parseEncryptionContext(req.Parameters)

	keyID := s.getKeyID(req.Parameters)
	var result *hsm.DecryptResult
	var keyArn string

	if keyID != "" {
		key, err := s.resolveKey(stores, req.Parameters)
		if err != nil {
			return nil, err
		}
		if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "Decrypt", key.KeyID, encryptionContext); err != nil {
			return nil, err
		}
		if err := s.validateKeyState(key); err != nil {
			return nil, err
		}
		result, err = s.hsmBackend.Decrypt(key.KeyID, ciphertext, encryptionContext)
		if err != nil {
			return nil, s.mapHSMError(err)
		}
		keyArn = key.Arn
	} else {
		var resolvedKeyID string
		result, resolvedKeyID, err = s.hsmBackend.DecryptWithoutKeyID(ciphertext, encryptionContext)
		if err != nil {
			return nil, s.mapHSMError(err)
		}
		key, err := stores.keys.Get(resolvedKeyID)
		if err != nil {
			return nil, err
		}
		if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "Decrypt", key.KeyID, encryptionContext); err != nil {
			return nil, err
		}
		keyArn = key.Arn
	}

	return map[string]interface{}{
		"Plaintext":           base64.StdEncoding.EncodeToString(result.Plaintext),
		"KeyId":               keyArn,
		"EncryptionAlgorithm": "SYMMETRIC_DEFAULT",
	}, nil
}

// ReEncrypt decrypts ciphertext using the source KMS key and then re-encrypts it using the destination KMS key.
func (s *KMSService) ReEncrypt(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sourceKeyID := request.GetStringParam(req.Parameters, "SourceKeyId")
	var sourceKey *kmsstore.Key

	ciphertextB64 := request.GetStringParam(req.Parameters, "CiphertextBlob")
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	sourceEncryptionContext := parseEncryptionContextForPrefix(req.Parameters, "SourceEncryptionContext")
	destinationEncryptionContext := parseEncryptionContextForPrefix(req.Parameters, "DestinationEncryptionContext")

	var decryptResult *hsm.DecryptResult
	var sourceKeyArn string

	if sourceKeyID != "" {
		sourceKey, err = s.resolveKeyByKeyID(stores, sourceKeyID)
		if err != nil {
			return nil, err
		}
		if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "Decrypt", sourceKey.KeyID, sourceEncryptionContext); err != nil {
			return nil, err
		}
		if err := s.validateKeyState(sourceKey); err != nil {
			return nil, err
		}
		if sourceKey.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
			return nil, ErrInvalidKeyUsage
		}
		decryptResult, err = s.hsmBackend.Decrypt(sourceKey.KeyID, ciphertext, sourceEncryptionContext)
		if err != nil {
			return nil, err
		}
		sourceKeyArn = sourceKey.Arn
	} else {
		var resolvedKeyID string
		decryptResult, resolvedKeyID, err = s.hsmBackend.DecryptWithoutKeyID(ciphertext, sourceEncryptionContext)
		if err != nil {
			return nil, err
		}
		sourceKey, err = stores.keys.Get(resolvedKeyID)
		if err != nil {
			return nil, err
		}
		if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "Decrypt", sourceKey.KeyID, sourceEncryptionContext); err != nil {
			return nil, err
		}
		if err := s.validateKeyState(sourceKey); err != nil {
			return nil, err
		}
		if sourceKey.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
			return nil, ErrInvalidKeyUsage
		}
		sourceKeyArn = sourceKey.Arn
	}

	destinationKeyID := request.GetStringParam(req.Parameters, "DestinationKeyId")
	destinationKey, err := s.resolveKeyByKeyID(stores, destinationKeyID)
	if err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "Encrypt", destinationKey.KeyID, destinationEncryptionContext); err != nil {
		return nil, err
	}

	if err := s.validateKeyState(destinationKey); err != nil {
		return nil, err
	}

	if destinationKey.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
		return nil, ErrInvalidKeyUsage
	}

	encryptResult, err := s.hsmBackend.Encrypt(destinationKey.KeyID, decryptResult.Plaintext, destinationEncryptionContext)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CiphertextBlob":            base64.StdEncoding.EncodeToString(encryptResult.Ciphertext),
		"SourceKeyId":               sourceKeyArn,
		"KeyId":                     destinationKey.Arn,
		"EncryptionAlgorithm":       "SYMMETRIC_DEFAULT",
		"SourceEncryptionAlgorithm": "SYMMETRIC_DEFAULT",
	}, nil
}

// GenerateDataKey generates a unique data key for encrypting data outside of KMS.
func (s *KMSService) GenerateDataKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	encryptionContext := parseEncryptionContext(req.Parameters)
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GenerateDataKey", encryptionContext)
	if err != nil {
		return nil, err
	}

	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}

	keySpec := request.GetStringParam(req.Parameters, "KeySpec")
	numberOfBytes := request.GetIntParam(req.Parameters, "NumberOfBytes")

	if keySpec == "" && numberOfBytes == 0 {
		keySpec = "AES_256"
	}

	result, err := s.hsmBackend.GenerateDataKey(key.KeyID, keySpec, numberOfBytes)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CiphertextBlob": base64.StdEncoding.EncodeToString(result.Ciphertext),
		"Plaintext":      base64.StdEncoding.EncodeToString(result.Plaintext),
		"KeyId":          key.Arn,
	}, nil
}

// GenerateDataKeyWithoutPlaintext generates a unique data key but returns only the ciphertext.
func (s *KMSService) GenerateDataKeyWithoutPlaintext(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	encryptionContext := parseEncryptionContext(req.Parameters)
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GenerateDataKeyWithoutPlaintext", encryptionContext)
	if err != nil {
		return nil, err
	}

	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}

	keySpec := request.GetStringParam(req.Parameters, "KeySpec")
	numberOfBytes := request.GetIntParam(req.Parameters, "NumberOfBytes")

	if keySpec == "" && numberOfBytes == 0 {
		keySpec = "AES_256"
	}

	result, err := s.hsmBackend.GenerateDataKey(key.KeyID, keySpec, numberOfBytes)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CiphertextBlob": base64.StdEncoding.EncodeToString(result.Ciphertext),
		"KeyId":          key.Arn,
	}, nil
}

// GenerateRandom returns a random byte string for use in cryptographic operations.
func (s *KMSService) GenerateRandom(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	numberOfBytes := request.GetIntParam(req.Parameters, "NumberOfBytes")
	if numberOfBytes == 0 {
		return nil, ErrValidation
	}
	if numberOfBytes < 1 || numberOfBytes > 1024 {
		return nil, ErrValidation
	}

	randomBytes := make([]byte, numberOfBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Plaintext": base64.StdEncoding.EncodeToString(randomBytes),
	}, nil
}

func (s *KMSService) resolveKeyByKeyID(stores *kmsStores, keyID string) (*kmsstore.Key, error) {
	if keyID == "" {
		return nil, ErrKeyNotFound
	}

	if stores.keys.ARNBuilder().IsAlias(keyID) {
		alias, err := stores.aliases.Get(keyID)
		if err != nil {
			return nil, err
		}
		keyID = alias.TargetKeyID
	}

	return stores.keys.Get(keyID)
}

func parseEncryptionContextForPrefix(params map[string]interface{}, prefix string) map[string]string {
	ctx := make(map[string]string)
	if ec, ok := params[prefix]; ok {
		if ecMap, ok := ec.(map[string]interface{}); ok {
			for k, v := range ecMap {
				if vs, ok := v.(string); ok {
					ctx[k] = vs
				}
			}
		}
	}
	return ctx
}

// GenerateDataKeyPair generates an asymmetric key pair and encrypts the private key with the KMS key.
func (s *KMSService) GenerateDataKeyPair(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	encryptionContext := parseEncryptionContext(req.Parameters)
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GenerateDataKeyPair", encryptionContext)
	if err != nil {
		return nil, err
	}

	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}

	if key.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
		return nil, ErrInvalidKeyUsage
	}

	keyPairSpec := hsm.KeySpec(request.GetStringParam(req.Parameters, "KeyPairSpec"))
	if keyPairSpec == "" {
		return nil, ErrValidation
	}

	privKey, pubKey, err := generateAsymmetricKeyPair(keyPairSpec)
	if err != nil {
		return nil, err
	}

	privKeyDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	pubKeyDER, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	encryptedResult, err := s.hsmBackend.Encrypt(key.KeyID, privKeyDER, encryptionContext)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PrivateKeyCiphertextBlob": base64.StdEncoding.EncodeToString(encryptedResult.Ciphertext),
		"PrivateKeyPlaintext":      base64.StdEncoding.EncodeToString(privKeyDER),
		"PublicKey":                base64.StdEncoding.EncodeToString(pubKeyDER),
		"KeyId":                    key.Arn,
		"KeyPairSpec":              keyPairSpec,
	}, nil
}

// GenerateDataKeyPairWithoutPlaintext generates an asymmetric key pair but returns only the encrypted private key.
func (s *KMSService) GenerateDataKeyPairWithoutPlaintext(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	encryptionContext := parseEncryptionContext(req.Parameters)
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GenerateDataKeyPairWithoutPlaintext", encryptionContext)
	if err != nil {
		return nil, err
	}

	if err := s.validateKeyState(key); err != nil {
		return nil, err
	}

	if key.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
		return nil, ErrInvalidKeyUsage
	}

	keyPairSpec := hsm.KeySpec(request.GetStringParam(req.Parameters, "KeyPairSpec"))
	if keyPairSpec == "" {
		return nil, ErrValidation
	}

	privKey, pubKey, err := generateAsymmetricKeyPair(keyPairSpec)
	if err != nil {
		return nil, err
	}

	privKeyDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	pubKeyDER, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	encryptedResult, err := s.hsmBackend.Encrypt(key.KeyID, privKeyDER, encryptionContext)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PrivateKeyCiphertextBlob": base64.StdEncoding.EncodeToString(encryptedResult.Ciphertext),
		"PublicKey":                base64.StdEncoding.EncodeToString(pubKeyDER),
		"KeyId":                    key.Arn,
		"KeyPairSpec":              keyPairSpec,
	}, nil
}

// ListKeyRotations returns the key rotation history for a KMS key.
func (s *KMSService) ListKeyRotations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "ListKeyRotations", nil)
	if err != nil {
		return nil, err
	}

	var rotations []map[string]interface{}

	if key.KeyRotationEnabled && !key.CreationDate.IsZero() {
		rotations = append(rotations, map[string]interface{}{
			"KeyId":        key.Arn,
			"RotationDate": key.CreationDate.Unix(),
			"RotationType": "AUTOMATIC",
		})
	}

	marker := pagination.GetMarker(req.Parameters)
	maxItems := pagination.GetMaxItems(req.Parameters, 100)

	type rotationEntry struct {
		data map[string]interface{}
	}
	entries := make([]rotationEntry, len(rotations))
	for i, r := range rotations {
		entries[i] = rotationEntry{data: r}
	}

	result := pagination.PaginateSlice(entries, marker, maxItems, func(e rotationEntry) string {
		return fmt.Sprintf("%v", e.data["RotationDate"])
	})

	paginatedRotations := make([]map[string]interface{}, len(result.Items))
	for i, e := range result.Items {
		paginatedRotations[i] = e.data
	}

	response := map[string]interface{}{
		"Rotations": paginatedRotations,
		"Truncated": result.IsTruncated,
	}
	if result.NextMarker != "" {
		response["NextMarker"] = result.NextMarker
	}

	return response, nil
}

func generateAsymmetricKeyPair(spec hsm.KeySpec) (interface{}, interface{}, error) {
	switch spec {
	case hsm.KeySpecRSA2048:
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	case hsm.KeySpecRSA3072:
		key, err := rsa.GenerateKey(rand.Reader, 3072)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	case hsm.KeySpecRSA4096:
		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	case hsm.KeySpecECCNISTP256:
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	case hsm.KeySpecECCNISTP384:
		key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	case hsm.KeySpecECCNISTP521:
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	case hsm.KeySpecECCSECPP256R1:
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	default:
		return nil, nil, fmt.Errorf("unsupported KeyPairSpec: %s", spec)
	}
}

func (s *KMSService) mapHSMError(err error) error {
	if errors.Is(err, hsm.ErrKeyNotFound) {
		return NewKeyNotFoundError("")
	}
	if errors.Is(err, hsm.ErrDecryptFailed) {
		return ErrInvalidCiphertext
	}
	return ErrKMSInternal
}
