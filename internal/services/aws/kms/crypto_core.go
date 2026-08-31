package kms

// crypto_core.go carries the Core functions of the two-step crypto
// operations (Decrypt and ReEncrypt). The single-shot operations (Encrypt,
// GenerateDataKey*, Sign, Verify, GenerateMac, VerifyMac) reach the store
// only through the auth_core resolution/authorisation layer and the HSM
// backend, so they need no Core bundle of their own; Decrypt and ReEncrypt
// additionally resolve keys from ciphertext and fall back to a keyless
// HSM lookup, which requires store reads that belong behind the Core
// boundary.

import (
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// validEncryptionAlgorithmSpec reports whether alg is a member of the
// Smithy EncryptionAlgorithmSpec enum that types the EncryptionAlgorithm
// request member of Decrypt and ReEncrypt. SM2PKE belongs to the enum even
// though this platform has no SM2 key material; a request naming it fails
// the key-level algorithm support check.
func validEncryptionAlgorithmSpec(alg string) bool {
	switch alg {
	case string(hsm.EncryptionAlgorithmSymmetricDefault),
		string(hsm.EncryptionAlgorithmRSAOAEPSHA1),
		string(hsm.EncryptionAlgorithmRSAOAEPSHA256),
		"SM2PKE":
		return true
	}
	return false
}

// defaultEncryptionAlgorithm returns the keyspec-default encryption
// algorithm. For key specs that support no encryption at all (HMAC,
// ECC_SIGNATURE, etc.) it returns an empty string; callers must check for
// this and surface UnsupportedOperationException rather than silently
// selecting an RSA default that would then fail inside the HSM.
func defaultEncryptionAlgorithm(key *kmsstore.Key) string {
	if key.KeySpec == kmsstore.KeySpecSymmetricDefault {
		return "SYMMETRIC_DEFAULT"
	}
	if isRSAKeySpec(key.KeySpec) {
		return "RSAES_OAEP_SHA_256"
	}
	return ""
}

// resolveEncryptionAlgorithm resolves the encryption algorithm for the
// request member named by memberName: the single-member operations carry
// "EncryptionAlgorithm", while ReEncrypt names "SourceEncryptionAlgorithm"
// and "DestinationEncryptionAlgorithm" separately.
//
// An explicitly supplied value must be a member of the EncryptionAlgorithmSpec
// enum — any other string is a shape violation the aws-json-1.1 contract
// rejects with SerializationException — and must be supported by the key:
// the InvalidKeyUsageException documentation covers "the encryption
// algorithm ... specified for the operation is incompatible with the type of
// key material in the KMS key". An omitted member falls back to the keyspec
// default.
func resolveEncryptionAlgorithm(key *kmsstore.Key, params map[string]interface{}, memberName string) (string, error) {
	algorithm := defaultEncryptionAlgorithm(key)
	if explicit := request.GetStringParam(params, memberName); explicit != "" {
		if !validEncryptionAlgorithmSpec(explicit) {
			return "", ErrSerializationException
		}
		algorithm = explicit
	}
	if algorithm != "" && !algorithmSupported(algorithm, key.EncryptionAlgorithms) {
		return "", ErrInvalidKeyUsage
	}
	return algorithm, nil
}

// DecryptInput carries the Decrypt members. Params transports the raw
// wire parameters so the Core can read EncryptionAlgorithm and DryRun in
// their original positions; Ciphertext travels already decoded because
// the wire gates (empty check, base64 decode, length bound) precede every
// store read in the original failure precedence.
type DecryptInput struct {
	KeyID             string
	Ciphertext        []byte
	EncryptionContext map[string]string
	Principal         string
	Params            map[string]interface{}
}

// DecryptCoreResult is the transport-agnostic Decrypt outcome.
type DecryptCoreResult struct {
	Plaintext           []byte
	KeyArn              string
	EncryptionAlgorithm string
}

// decryptCore is the single entry point for Decrypt. With a KeyId it
// resolves and authorises the key up front; without one it falls back to
// the HSM's keyless lookup and authorises the resolved key afterwards.
func (s *KMSService) decryptCore(stores *kmsStores, in DecryptInput) (*DecryptCoreResult, error) {
	var result *hsm.DecryptResult
	var keyArn string
	var resolvedKey *kmsstore.Key
	// usedAlgorithm is the algorithm actually dispatched to the HSM; the
	// response echoes this value, never the raw caller string.
	var usedAlgorithm string

	if in.KeyID != "" {
		key, err := s.resolveKeyByKeyID(stores, in.KeyID)
		if err != nil {
			return nil, err
		}
		if err := s.authorizeOperation(stores, in.Principal, "Decrypt", key.KeyID, in.EncryptionContext); err != nil {
			return nil, err
		}
		if err := s.validateKeyState(key); err != nil {
			return nil, err
		}
		// AWS: Decrypt on a non-ENCRYPT_DECRYPT key returns
		// InvalidKeyUsageException. Encrypt has this guard but Decrypt was
		// historically missing it, allowing HMAC/SignVerify keys to reach
		// the HSM where they fail with the misleading ErrDecryptFailed.
		if key.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
			return nil, ErrInvalidKeyUsage
		}
		decryptionAlgorithm, err := resolveEncryptionAlgorithm(key, in.Params, "EncryptionAlgorithm")
		if err != nil {
			return nil, err
		}
		if err := checkKMSDryRun(in.Params); err != nil {
			return nil, err
		}
		result, err = s.hsmBackend.Decrypt(key.KeyID, in.Ciphertext, hsm.EncryptionAlgorithm(decryptionAlgorithm), in.EncryptionContext)
		if err != nil {
			return nil, s.mapHSMError(err)
		}
		keyArn = key.Arn
		resolvedKey = key
		usedAlgorithm = decryptionAlgorithm
	} else {
		// With the KeyId member omitted, the AWS Decrypt contract resolves
		// the key from the metadata inside the symmetric ciphertext blob,
		// so this resolution path is symmetric-only. An explicit algorithm
		// is still enum-checked up front; after resolution the resolved key
		// must support it — the asymmetric algorithms name ciphertext this
		// path cannot resolve, so the caller must supply KeyId for those.
		wireAlgorithm := request.GetStringParam(in.Params, "EncryptionAlgorithm")
		if wireAlgorithm != "" && !validEncryptionAlgorithmSpec(wireAlgorithm) {
			return nil, ErrSerializationException
		}
		if err := checkKMSDryRun(in.Params); err != nil {
			return nil, err
		}
		var resolvedKeyID string
		var err error
		result, resolvedKeyID, err = s.hsmBackend.DecryptWithoutKeyID(in.Ciphertext, hsm.EncryptionAlgorithmSymmetricDefault, in.EncryptionContext)
		if err != nil {
			return nil, s.mapHSMError(err)
		}
		key, err := stores.keys.Get(resolvedKeyID)
		if err != nil {
			return nil, NewKeyNotFoundError(resolvedKeyID)
		}
		if err := s.authorizeOperation(stores, in.Principal, "Decrypt", key.KeyID, in.EncryptionContext); err != nil {
			return nil, err
		}
		if err := s.validateKeyState(key); err != nil {
			return nil, err
		}
		if wireAlgorithm != "" && !algorithmSupported(wireAlgorithm, key.EncryptionAlgorithms) {
			return nil, ErrInvalidKeyUsage
		}
		keyArn = key.Arn
		resolvedKey = key
		usedAlgorithm = string(hsm.EncryptionAlgorithmSymmetricDefault)
	}

	s.markKeyLastUsed(stores, resolvedKey.KeyID, "Decrypt")

	return &DecryptCoreResult{
		Plaintext:           result.Plaintext,
		KeyArn:              keyArn,
		EncryptionAlgorithm: usedAlgorithm,
	}, nil
}

// ReEncryptInput carries the ReEncrypt members. Params transports the raw
// wire parameters for the EncryptionAlgorithm reads; both encryption
// contexts travel parsed because the parse is pure.
type ReEncryptInput struct {
	SourceKeyID                  string
	DestinationKeyID             string
	Ciphertext                   []byte
	SourceEncryptionContext      map[string]string
	DestinationEncryptionContext map[string]string
	Principal                    string
	Params                       map[string]interface{}
}

// ReEncryptCoreResult is the transport-agnostic ReEncrypt outcome.
type ReEncryptCoreResult struct {
	CiphertextBlob                 []byte
	SourceKeyArn                   string
	DestinationKeyArn              string
	SourceEncryptionAlgorithm      string
	DestinationEncryptionAlgorithm string
}

// reEncryptCore is the single entry point for ReEncrypt: decrypt under the
// source key (explicit or resolved from the ciphertext), re-encrypt under
// the destination key, record usage on both.
func (s *KMSService) reEncryptCore(stores *kmsStores, in ReEncryptInput) (*ReEncryptCoreResult, error) {
	var sourceKey *kmsstore.Key

	var decryptResult *hsm.DecryptResult
	var sourceKeyArn string
	var sourceAlgorithm string
	var err error

	if in.SourceKeyID != "" {
		sourceKey, err = s.resolveKeyByKeyID(stores, in.SourceKeyID)
		if err != nil {
			return nil, err
		}
		if err := s.authorizeOperation(stores, in.Principal, "Decrypt", sourceKey.KeyID, in.SourceEncryptionContext); err != nil {
			return nil, err
		}
		if err := s.validateKeyState(sourceKey); err != nil {
			return nil, err
		}
		if sourceKey.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
			return nil, ErrInvalidKeyUsage
		}
		if err := checkKMSDryRun(in.Params); err != nil {
			return nil, err
		}
		sourceAlgorithm, err = resolveEncryptionAlgorithm(sourceKey, in.Params, "SourceEncryptionAlgorithm")
		if err != nil {
			return nil, err
		}
		decryptResult, err = s.hsmBackend.Decrypt(sourceKey.KeyID, in.Ciphertext, hsm.EncryptionAlgorithm(sourceAlgorithm), in.SourceEncryptionContext)
		if err != nil {
			return nil, err
		}
		sourceKeyArn = sourceKey.Arn
	} else {
		if err := checkKMSDryRun(in.Params); err != nil {
			return nil, err
		}
		var resolvedKeyID string
		decryptResult, resolvedKeyID, err = s.hsmBackend.DecryptWithoutKeyID(in.Ciphertext, hsm.EncryptionAlgorithmSymmetricDefault, in.SourceEncryptionContext)
		if err != nil {
			return nil, err
		}
		sourceKey, err = stores.keys.Get(resolvedKeyID)
		if err != nil {
			return nil, err
		}
		if err := s.authorizeOperation(stores, in.Principal, "Decrypt", sourceKey.KeyID, in.SourceEncryptionContext); err != nil {
			return nil, err
		}
		if err := s.validateKeyState(sourceKey); err != nil {
			return nil, err
		}
		if sourceKey.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
			return nil, ErrInvalidKeyUsage
		}
		sourceKeyArn = sourceKey.Arn
		// The keyless resolution path is symmetric-only, so an explicit
		// SourceEncryptionAlgorithm must still be enum-valid and must be
		// supported by the resolved key — an asymmetric algorithm cannot
		// name ciphertext this path is able to resolve.
		if srcAlg := request.GetStringParam(in.Params, "SourceEncryptionAlgorithm"); srcAlg != "" {
			if !validEncryptionAlgorithmSpec(srcAlg) {
				return nil, ErrSerializationException
			}
			if !algorithmSupported(srcAlg, sourceKey.EncryptionAlgorithms) {
				return nil, ErrInvalidKeyUsage
			}
		}
		// No caller-supplied KeyId ⇒ no caller-supplied algorithm; the
		// DecryptWithoutKeyID path always uses SYMMETRIC_DEFAULT.
		sourceAlgorithm = "SYMMETRIC_DEFAULT"
	}

	destinationKey, err := s.resolveKeyByKeyID(stores, in.DestinationKeyID)
	if err != nil {
		return nil, err
	}
	if err := s.authorizeOperation(stores, in.Principal, "Encrypt", destinationKey.KeyID, in.DestinationEncryptionContext); err != nil {
		return nil, err
	}

	if err := s.validateKeyState(destinationKey); err != nil {
		return nil, err
	}

	if destinationKey.KeyUsage != kmsstore.KeyUsageEncryptDecrypt {
		return nil, ErrInvalidKeyUsage
	}

	destinationAlgorithm, err := resolveEncryptionAlgorithm(destinationKey, in.Params, "DestinationEncryptionAlgorithm")
	if err != nil {
		return nil, err
	}
	encryptResult, err := s.hsmBackend.Encrypt(destinationKey.KeyID, decryptResult.Plaintext, hsm.EncryptionAlgorithm(destinationAlgorithm), in.DestinationEncryptionContext)
	if err != nil {
		return nil, err
	}
	s.markKeyLastUsed(stores, sourceKey.KeyID, "ReEncrypt")
	s.markKeyLastUsed(stores, destinationKey.KeyID, "ReEncrypt")

	return &ReEncryptCoreResult{
		CiphertextBlob:                 encryptResult.Ciphertext,
		SourceKeyArn:                   sourceKeyArn,
		DestinationKeyArn:              destinationKey.Arn,
		SourceEncryptionAlgorithm:      sourceAlgorithm,
		DestinationEncryptionAlgorithm: destinationAlgorithm,
	}, nil
}
