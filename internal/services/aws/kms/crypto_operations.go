package kms

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
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

	encryptionAlgorithm := determineEncryptionAlgorithm(key, req.Parameters)
	if !algorithmSupported(encryptionAlgorithm, key.EncryptionAlgorithms) {
		return nil, ErrInvalidAlgorithm
	}

	plaintextB64 := request.GetStringParam(req.Parameters, "Plaintext")
	if plaintextB64 == "" {
		// AWS rejects empty Plaintext with ValidationException. The
		// previous code silently treated empty as 0-byte plaintext.
		return nil, ErrValidation
	}
	plaintext, err := base64.StdEncoding.DecodeString(plaintextB64)
	if err != nil {
		plaintext, err = base64.RawStdEncoding.DecodeString(plaintextB64)
		if err != nil {
			return nil, ErrValidation
		}
	}

	// AWS plaintext length limits:
	//   SYMMETRIC_DEFAULT: 1-4096 bytes
	//   RSA: 190/318/446 bytes for RSA_2048/3072/4096 with SHA-256 OAEP,
	//        214/342/470 bytes for SHA-1 OAEP.
	// The HSM backend enforces the RSA cap (returns ErrEncryptFailed
	// which surfaces as KMSInternalException); enforce the symmetric
	// cap here so callers get a ValidationException instead.
	if encryptionAlgorithm == "SYMMETRIC_DEFAULT" && len(plaintext) > 4096 {
		return nil, ErrValidation
	}

	encryptionContext := parseEncryptionContext(req.Parameters)
	if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "Encrypt", key.KeyID, encryptionContext); err != nil {
		return nil, err
	}

	if err := checkKMSDryRun(req.Parameters); err != nil {
		return nil, err
	}

	result, err := s.hsmBackend.Encrypt(key.KeyID, plaintext, hsm.EncryptionAlgorithm(encryptionAlgorithm), encryptionContext)
	if err != nil {
		return nil, err
	}
	s.markKeyLastUsed(stores, key.KeyID, "Encrypt")

	return map[string]interface{}{
		"CiphertextBlob":      base64.StdEncoding.EncodeToString(result.Ciphertext),
		"KeyId":               key.Arn,
		"EncryptionAlgorithm": encryptionAlgorithm,
	}, nil
}

// errRecipientNotSupported rejects the Recipient parameter (Nitro Enclaves
// envelope encryption) with a message that names the actual gap; the shared
// ErrUnsupportedOperation sentinel talks about key types, which does not
// fit this guard.
var errRecipientNotSupported = awserrors.NewAWSError("UnsupportedOperationException",
	"Recipient is not supported: Nitro Enclaves envelope encryption is not implemented", http.StatusBadRequest)

// rejectRecipient rejects requests carrying the Recipient parameter
// (Decrypt, GenerateDataKey, GenerateDataKeyPair, GenerateRandom,
// DeriveSharedSecret in the Smithy model). Recipient is not implemented;
// silently dropping it would hand normal plaintext to an enclave caller
// that expects CiphertextForRecipient, so the request is rejected
// explicitly.
func rejectRecipient(req *request.ParsedRequest) error {
	if v, ok := req.Parameters["Recipient"]; ok && v != nil {
		return errRecipientNotSupported
	}
	return nil
}

// Decrypt decrypts ciphertext using the specified KMS key.
func (s *KMSService) Decrypt(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectRecipient(req); err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ciphertextB64 := request.GetStringParam(req.Parameters, "CiphertextBlob")
	if ciphertextB64 == "" {
		// AWS rejects empty CiphertextBlob with InvalidCiphertext.
		return nil, ErrInvalidCiphertext
	}
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}
	if err := validateCiphertextLength(len(ciphertext)); err != nil {
		return nil, err
	}

	encryptionContext := parseEncryptionContext(req.Parameters)

	keyID := s.getKeyID(req.Parameters)
	var result *hsm.DecryptResult
	var keyArn string
	var resolvedKey *kmsstore.Key

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
		// AWS: Decrypt on a non-ENCRYPT_DECRYPT key returns
		// InvalidKeyUsageException. Encrypt has this guard at line 30 but
		// Decrypt was missing it, allowing HMAC/SignVerify keys to reach
		// the HSM where they fail with the misleading ErrDecryptFailed.
		if key.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
			return nil, ErrInvalidKeyUsage
		}
		decryptionAlgorithm := determineEncryptionAlgorithm(key, req.Parameters)
		// Mirror the Encrypt guard at line 35: when the caller explicitly
		// supplies an EncryptionAlgorithm that is not in the key's
		// supported list, reject with InvalidAlgorithmException rather
		// than dispatching a request the HSM will refuse.
		if decryptionAlgorithm != "" && !algorithmSupported(decryptionAlgorithm, key.EncryptionAlgorithms) {
			return nil, ErrInvalidAlgorithm
		}
		if err := checkKMSDryRun(req.Parameters); err != nil {
			return nil, err
		}
		result, err = s.hsmBackend.Decrypt(key.KeyID, ciphertext, hsm.EncryptionAlgorithm(decryptionAlgorithm), encryptionContext)
		if err != nil {
			return nil, s.mapHSMError(err)
		}
		keyArn = key.Arn
		resolvedKey = key
	} else {
		var resolvedKeyID string
		// Without a KeyId, the algorithm cannot be requested by the caller
		// (DecryptRequest has no EncryptionAlgorithm member). Fall back to
		// SYMMETRIC_DEFAULT first; if the HSM cannot decrypt (key is RSA)
		// the caller must resubmit with an explicit KeyId.
		if err := checkKMSDryRun(req.Parameters); err != nil {
			return nil, err
		}
		result, resolvedKeyID, err = s.hsmBackend.DecryptWithoutKeyID(ciphertext, hsm.EncryptionAlgorithmSymmetricDefault, encryptionContext)
		if err != nil {
			return nil, s.mapHSMError(err)
		}
		key, err := stores.keys.Get(resolvedKeyID)
		if err != nil {
			return nil, NewKeyNotFoundError(resolvedKeyID)
		}
		if err := s.authorizeOperation(stores, s.resolveCallerPrincipal(reqCtx, req), "Decrypt", key.KeyID, encryptionContext); err != nil {
			return nil, err
		}
		if err := s.validateKeyState(key); err != nil {
			return nil, err
		}
		keyArn = key.Arn
		resolvedKey = key
	}

	s.markKeyLastUsed(stores, resolvedKey.KeyID, "Decrypt")

	return map[string]interface{}{
		"Plaintext":           base64.StdEncoding.EncodeToString(result.Plaintext),
		"KeyId":               keyArn,
		"EncryptionAlgorithm": determineEncryptionAlgorithm(resolvedKey, req.Parameters),
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
	if err := validateCiphertextLength(len(ciphertext)); err != nil {
		return nil, err
	}

	sourceEncryptionContext := parseEncryptionContextForPrefix(req.Parameters, "SourceEncryptionContext")
	destinationEncryptionContext := parseEncryptionContextForPrefix(req.Parameters, "DestinationEncryptionContext")

	var decryptResult *hsm.DecryptResult
	var sourceKeyArn string
	var sourceAlgorithm string

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
		if err := checkKMSDryRun(req.Parameters); err != nil {
			return nil, err
		}
		sourceAlgorithm = determineEncryptionAlgorithm(sourceKey, req.Parameters)
		decryptResult, err = s.hsmBackend.Decrypt(sourceKey.KeyID, ciphertext, hsm.EncryptionAlgorithm(sourceAlgorithm), sourceEncryptionContext)
		if err != nil {
			return nil, err
		}
		sourceKeyArn = sourceKey.Arn
	} else {
		var resolvedKeyID string
		if err := checkKMSDryRun(req.Parameters); err != nil {
			return nil, err
		}
		decryptResult, resolvedKeyID, err = s.hsmBackend.DecryptWithoutKeyID(ciphertext, hsm.EncryptionAlgorithmSymmetricDefault, sourceEncryptionContext)
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
		// No caller-supplied KeyId ⇒ no caller-supplied algorithm; the
		// DecryptWithoutKeyID path always uses SYMMETRIC_DEFAULT.
		sourceAlgorithm = "SYMMETRIC_DEFAULT"
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

	destinationAlgorithm := determineEncryptionAlgorithm(destinationKey, req.Parameters)
	encryptResult, err := s.hsmBackend.Encrypt(destinationKey.KeyID, decryptResult.Plaintext, hsm.EncryptionAlgorithm(destinationAlgorithm), destinationEncryptionContext)
	if err != nil {
		return nil, err
	}
	s.markKeyLastUsed(stores, sourceKey.KeyID, "ReEncrypt")
	s.markKeyLastUsed(stores, destinationKey.KeyID, "ReEncrypt")

	return map[string]interface{}{
		"CiphertextBlob":                 base64.StdEncoding.EncodeToString(encryptResult.Ciphertext),
		"SourceKeyId":                    sourceKeyArn,
		"KeyId":                          destinationKey.Arn,
		"SourceEncryptionAlgorithm":      sourceAlgorithm,
		"DestinationEncryptionAlgorithm": destinationAlgorithm,
	}, nil
}

// GenerateDataKey generates a unique data key for encrypting data outside of KMS.
func (s *KMSService) GenerateDataKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectRecipient(req); err != nil {
		return nil, err
	}

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

	// AWS: GenerateDataKey requires an ENCRYPT_DECRYPT key. Using a
	// SIGN_VERIFY or GENERATE_VERIFY_MAC key is rejected with
	// InvalidKeyUsageException, consistent with Encrypt/Decrypt/ReEncrypt.
	if key.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
		return nil, ErrInvalidKeyUsage
	}

	if err := validateDataKeySpecAndBytes(key, req.Parameters); err != nil {
		return nil, err
	}

	if err := checkKMSDryRun(req.Parameters); err != nil {
		return nil, err
	}

	keySpec := request.GetStringParam(req.Parameters, "KeySpec")
	numberOfBytes := request.GetIntParam(req.Parameters, "NumberOfBytes")

	if keySpec == "" && numberOfBytes == 0 {
		keySpec = "AES_256"
	}

	result, err := s.hsmBackend.GenerateDataKey(key.KeyID, keySpec, numberOfBytes, encryptionContext)
	if err != nil {
		return nil, err
	}
	s.markKeyLastUsed(stores, key.KeyID, "GenerateDataKey")

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

	// AWS: GenerateDataKeyWithoutPlaintext requires an ENCRYPT_DECRYPT key,
	// same as GenerateDataKey. Using a non-encryption key is rejected with
	// InvalidKeyUsageException.
	if key.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
		return nil, ErrInvalidKeyUsage
	}

	if err := validateDataKeySpecAndBytes(key, req.Parameters); err != nil {
		return nil, err
	}

	if err := checkKMSDryRun(req.Parameters); err != nil {
		return nil, err
	}

	keySpec := request.GetStringParam(req.Parameters, "KeySpec")
	numberOfBytes := request.GetIntParam(req.Parameters, "NumberOfBytes")

	if keySpec == "" && numberOfBytes == 0 {
		keySpec = "AES_256"
	}

	result, err := s.hsmBackend.GenerateDataKey(key.KeyID, keySpec, numberOfBytes, encryptionContext)
	if err != nil {
		return nil, err
	}
	s.markKeyLastUsed(stores, key.KeyID, "GenerateDataKeyWithoutPlaintext")

	return map[string]interface{}{
		"CiphertextBlob": base64.StdEncoding.EncodeToString(result.Ciphertext),
		"KeyId":          key.Arn,
	}, nil
}

// validateDataKeySpecAndBytes implements the AWS GenerateDataKey input rules:
//   - KeySpec and NumberOfBytes are mutually exclusive
//   - KeySpec must be AES_128 or AES_256 when set
//   - NumberOfBytes must be 1-1024 when set
//
// The previous code silently accepted both, both empty, or invalid KeySpec
// values (defaulting to AES_256), which violates the AWS contract.
func validateDataKeySpecAndBytes(key *kmsstore.Key, params map[string]interface{}) error {
	keySpec := request.GetStringParam(params, "KeySpec")
	numberOfBytes := request.GetIntParam(params, "NumberOfBytes")
	if keySpec != "" && numberOfBytes != 0 {
		return ErrValidation
	}
	if keySpec != "" && keySpec != "AES_128" && keySpec != "AES_256" {
		return ErrValidation
	}
	if numberOfBytes != 0 && (numberOfBytes < 1 || numberOfBytes > 1024) {
		return ErrValidation
	}
	return nil
}

// GenerateRandom returns a random byte string for use in cryptographic operations.
func (s *KMSService) GenerateRandom(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectRecipient(req); err != nil {
		return nil, err
	}

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
	return s.resolveKey(stores, map[string]interface{}{"KeyId": keyID})
}

func parseEncryptionContextForPrefix(params map[string]interface{}, prefix string) map[string]string {
	if ec, ok := params[prefix]; ok {
		if ecMap, ok := ec.(map[string]interface{}); ok {
			return request.CopyStringMap(ecMap)
		}
	}
	return nil
}

// GenerateDataKeyPair generates an asymmetric key pair and encrypts the private key with the KMS key.
func (s *KMSService) GenerateDataKeyPair(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := rejectRecipient(req); err != nil {
		return nil, err
	}

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
	if !isValidKeyPairSpec(keyPairSpec) {
		return nil, ErrValidation
	}

	if err := checkKMSDryRun(req.Parameters); err != nil {
		return nil, err
	}

	privKeyDER, pubKeyDER, err := s.hsmBackend.GenerateKeyPair(hsm.KeySpec(keyPairSpec))
	if err != nil {
		return nil, err
	}

	encryptedResult, err := s.hsmBackend.Encrypt(key.KeyID, privKeyDER, hsm.EncryptionAlgorithmSymmetricDefault, encryptionContext)
	if err != nil {
		return nil, err
	}
	s.markKeyLastUsed(stores, key.KeyID, "GenerateDataKeyPair")

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
	if !isValidKeyPairSpec(keyPairSpec) {
		return nil, ErrValidation
	}

	if err := checkKMSDryRun(req.Parameters); err != nil {
		return nil, err
	}

	privKeyDER, pubKeyDER, err := s.hsmBackend.GenerateKeyPair(keyPairSpec)
	if err != nil {
		return nil, err
	}

	encryptedResult, err := s.hsmBackend.Encrypt(key.KeyID, privKeyDER, hsm.EncryptionAlgorithmSymmetricDefault, encryptionContext)
	if err != nil {
		return nil, err
	}
	s.markKeyLastUsed(stores, key.KeyID, "GenerateDataKeyPairWithoutPlaintext")

	return map[string]interface{}{
		"PrivateKeyCiphertextBlob": base64.StdEncoding.EncodeToString(encryptedResult.Ciphertext),
		"PublicKey":                base64.StdEncoding.EncodeToString(pubKeyDER),
		"KeyId":                    key.Arn,
		"KeyPairSpec":              keyPairSpec,
	}, nil
}

// ListKeyRotations returns the key rotation history for a KMS key.
// Returns an empty list when no rotations have been performed.
// AWS supports pagination via Marker/Limit on this operation; the previous
// implementation ignored those parameters entirely.
func (s *KMSService) ListKeyRotations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "ListKeyRotations", nil)
	if err != nil {
		return nil, err
	}

	marker := pagination.GetMarker(req.Parameters)
	if err := validateMarkerLength(marker); err != nil {
		return nil, err
	}
	maxItems := pagination.GetMaxItems(req.Parameters, 100)

	// Return the actual rotation history recorded by RotateKeyOnDemand.
	// Each entry includes the rotation date, type, and key material version.
	rotations := make([]map[string]interface{}, 0, len(key.RotationHistory))
	for _, entry := range key.RotationHistory {
		rotations = append(rotations, map[string]interface{}{
			"RotationDate":  entry.RotationDate.Unix(),
			"RotationType":  entry.RotationType,
			"KeyMaterialId": entry.KeyMaterialId,
		})
	}

	result := pagination.PaginateSlice(rotations, marker, maxItems, func(item map[string]interface{}) string {
		if ts, ok := item["RotationDate"].(int64); ok {
			return fmt.Sprintf("%d", ts)
		}
		return ""
	})

	response := map[string]interface{}{
		"Rotations": result.Items,
		"Truncated": result.IsTruncated,
	}
	if result.NextMarker != "" {
		response["NextMarker"] = result.NextMarker
	}
	return response, nil
}

func (s *KMSService) mapHSMError(err error) error {
	if errors.Is(err, hsm.ErrKeyNotFound) {
		return NewKeyNotFoundError("")
	}
	if errors.Is(err, hsm.ErrDecryptFailed) {
		return ErrInvalidCiphertext
	}
	if errors.Is(err, hsm.ErrInvalidCiphertext) {
		return ErrInvalidCiphertext
	}
	return ErrKMSInternal
}

// determineEncryptionAlgorithm returns the encryption algorithm to use for
// the given key. If the caller specifies one via the EncryptionAlgorithm
// parameter, it is used. Otherwise, the default for the key spec is returned.
// For key specs that do not support encryption at all (HMAC, ECC_SIGNATURE,
// etc.), the function returns an empty string; callers must check for this
// and surface UnsupportedOperationException rather than silently selecting
// an RSA default that would then fail inside the HSM.
func determineEncryptionAlgorithm(key *kmsstore.Key, params map[string]interface{}) string {
	if alg := request.GetStringParam(params, "EncryptionAlgorithm"); alg != "" {
		return alg
	}
	if key.KeySpec == kmsstore.KeySpecSymmetricDefault {
		return "SYMMETRIC_DEFAULT"
	}
	if isRSAKeySpec(key.KeySpec) {
		return "RSAES_OAEP_SHA_256"
	}
	return ""
}

// isRSAKeySpec reports whether the key spec is one of the RSA_* specs that
// support Encrypt/Decrypt with RSAES_OAEP_SHA_* algorithms.
func isRSAKeySpec(spec kmsstore.KeySpec) bool {
	switch spec {
	case kmsstore.KeySpecRSA2048,
		kmsstore.KeySpecRSA3072,
		kmsstore.KeySpecRSA4096:
		return true
	}
	return false
}

// isValidKeyPairSpec reports whether spec is one of the AWS-supported
// DataKeyPairSpec values for GenerateDataKeyPair[WithoutPlaintext].
// Per Smithy com.amazonaws.kms#DataKeyPairSpec, the enum has 9 members:
// RSA_2048, RSA_3072, RSA_4096, ECC_NIST_P256, ECC_NIST_P384,
// ECC_NIST_P521, ECC_SECG_P256K1, SM2, ECC_NIST_EDWARDS25519.
func isValidKeyPairSpec(spec hsm.KeySpec) bool {
	// The HSM backend (generateKeyPairDER) supports 7 specs:
	// RSA_2048/3072/4096, ECC_NIST_P256/P384/P521, ECC_SECG_P256K1.
	// SM2 and ECC_NIST_EDWARDS25519 are not implemented; accepting them
	// here caused the HSM to return ErrInvalidKeySpec which
	// mapHSMError translated to ErrKMSInternal (500) instead of the
	// expected ValidationException (400).
	switch spec {
	case hsm.KeySpecRSA2048, hsm.KeySpecRSA3072, hsm.KeySpecRSA4096,
		hsm.KeySpecECCNISTP256, hsm.KeySpecECCNISTP384, hsm.KeySpecECCNISTP521,
		hsm.KeySpecECCSECGP256K1:
		return true
	}
	return false
}
