package hsm

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"errors"
	"math/big"
	"sync"
)

const (
	aes256KeySize     = 32
	hmacSHA224KeySize = 28
	hmacSHA384KeySize = 48
	hmacSHA512KeySize = 64
)

type memoryKey struct {
	keySpec   KeySpec
	symmetric []byte
	rsaKey    *rsa.PrivateKey
	ecdsaKey  *ecdsa.PrivateKey
}

type cryptoOps struct {
	mu   sync.RWMutex
	keys map[string]*memoryKey
}

func (c *cryptoOps) initKeys() {
	c.keys = make(map[string]*memoryKey)
}

func (c *cryptoOps) generateKeyMaterial(keySpec KeySpec) (*memoryKey, error) {
	key := &memoryKey{keySpec: keySpec}

	switch keySpec {
	case KeySpecSymmetricDefault:
		key.symmetric = make([]byte, aes256KeySize)
		if _, err := rand.Read(key.symmetric); err != nil {
			return nil, err
		}
	case KeySpecHMAC224:
		key.symmetric = make([]byte, hmacSHA224KeySize)
		if _, err := rand.Read(key.symmetric); err != nil {
			return nil, err
		}
	case KeySpecHMAC256:
		key.symmetric = make([]byte, aes256KeySize)
		if _, err := rand.Read(key.symmetric); err != nil {
			return nil, err
		}
	case KeySpecHMAC384:
		key.symmetric = make([]byte, hmacSHA384KeySize)
		if _, err := rand.Read(key.symmetric); err != nil {
			return nil, err
		}
	case KeySpecHMAC512:
		key.symmetric = make([]byte, hmacSHA512KeySize)
		if _, err := rand.Read(key.symmetric); err != nil {
			return nil, err
		}
	case KeySpecRSA2048:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		key.rsaKey = rsaKey
	case KeySpecRSA3072:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 3072)
		if err != nil {
			return nil, err
		}
		key.rsaKey = rsaKey
	case KeySpecRSA4096:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, err
		}
		key.rsaKey = rsaKey
	case KeySpecECCNISTP256, KeySpecECCSECPP256R1:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		key.ecdsaKey = ecdsaKey
	case KeySpecECCNISTP384:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, err
		}
		key.ecdsaKey = ecdsaKey
	case KeySpecECCNISTP521:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, err
		}
		key.ecdsaKey = ecdsaKey
	default:
		return nil, ErrInvalidKeySpec
	}

	return key, nil
}

func (c *cryptoOps) encrypt(keyID string, plaintext []byte, context map[string]string) (*EncryptResult, error) {
	key, exists := c.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrEncryptFailed
	}

	ciphertext, err := aesEncrypt(key.symmetric, plaintext, serializeEncryptionContext(context))
	if err != nil {
		return nil, err
	}

	ciphertextWithKeyID := prependKeyIDToCiphertext(keyID, ciphertext)

	return &EncryptResult{
		Ciphertext:      ciphertextWithKeyID,
		CiphertextCRC32: computeCRC32(ciphertextWithKeyID),
	}, nil
}

func (c *cryptoOps) decrypt(keyID string, ciphertext []byte, context map[string]string) (*DecryptResult, error) {
	extractedKeyID, actualCiphertext, err := extractKeyIDFromCiphertext(ciphertext)
	if err == nil && extractedKeyID != "" {
		keyID = extractedKeyID
		ciphertext = actualCiphertext
	}

	key, exists := c.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrDecryptFailed
	}

	plaintext, err := aesDecrypt(key.symmetric, ciphertext, serializeEncryptionContext(context))
	if err != nil {
		return nil, err
	}

	return &DecryptResult{
		Plaintext:      plaintext,
		PlaintextCRC32: computeCRC32(plaintext),
	}, nil
}

func (c *cryptoOps) decryptWithoutKeyID(ciphertext []byte, context map[string]string) (*DecryptResult, string, error) {
	keyID, actualCiphertext, err := extractKeyIDFromCiphertext(ciphertext)
	if err != nil {
		return nil, "", err
	}

	key, exists := c.keys[keyID]
	if !exists {
		return nil, "", ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, "", ErrDecryptFailed
	}

	plaintext, err := aesDecrypt(key.symmetric, actualCiphertext, serializeEncryptionContext(context))
	if err != nil {
		return nil, "", err
	}

	return &DecryptResult{
		Plaintext:      plaintext,
		PlaintextCRC32: computeCRC32(plaintext),
	}, keyID, nil
}

func (c *cryptoOps) sign(keyID string, message []byte, algorithm SigningAlgorithm) (*SignResult, error) {
	key, exists := c.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	var signature []byte
	var err error

	if key.rsaKey != nil {
		hash := hashForAlgorithm(algorithm, message)
		switch algorithm {
		case SigningAlgorithmRSAPKCS1SHA256, SigningAlgorithmRSAPKCS1SHA384, SigningAlgorithmRSAPKCS1SHA512:
			signature, err = rsa.SignPKCS1v15(rand.Reader, key.rsaKey, hashAlgorithm(algorithm), hash)
		case SigningAlgorithmRSAPSSSHA256, SigningAlgorithmRSAPSSSHA384, SigningAlgorithmRSAPSSSHA512:
			signature, err = rsa.SignPSS(rand.Reader, key.rsaKey, hashAlgorithm(algorithm), hash, nil)
		default:
			return nil, ErrInvalidAlgorithm
		}
	} else if key.ecdsaKey != nil {
		hash := hashForAlgorithm(algorithm, message)
		var r, s *big.Int
		r, s, err = ecdsa.Sign(rand.Reader, key.ecdsaKey, hash)
		if err != nil {
			return nil, ErrSignFailed
		}
		signature, err = asn1.Marshal(struct {
			R, S *big.Int
		}{r, s})
	} else {
		return nil, ErrInvalidKeySpec
	}

	if err != nil {
		return nil, ErrSignFailed
	}

	return &SignResult{
		Signature:      signature,
		SignatureCRC32: computeCRC32(signature),
	}, nil
}

func (c *cryptoOps) verify(keyID string, message, signature []byte, algorithm SigningAlgorithm) (bool, error) {
	key, exists := c.keys[keyID]
	if !exists {
		return false, ErrKeyNotFound
	}

	if key.rsaKey != nil {
		hash := hashForAlgorithm(algorithm, message)
		switch algorithm {
		case SigningAlgorithmRSAPKCS1SHA256, SigningAlgorithmRSAPKCS1SHA384, SigningAlgorithmRSAPKCS1SHA512:
			err := rsa.VerifyPKCS1v15(&key.rsaKey.PublicKey, hashAlgorithm(algorithm), hash, signature)
			return err == nil, nil
		case SigningAlgorithmRSAPSSSHA256, SigningAlgorithmRSAPSSSHA384, SigningAlgorithmRSAPSSSHA512:
			err := rsa.VerifyPSS(&key.rsaKey.PublicKey, hashAlgorithm(algorithm), hash, signature, nil)
			return err == nil, nil
		default:
			return false, ErrInvalidAlgorithm
		}
	} else if key.ecdsaKey != nil {
		hash := hashForAlgorithm(algorithm, message)
		var sig struct {
			R, S *big.Int
		}
		if _, err := asn1.Unmarshal(signature, &sig); err != nil {
			return false, err
		}
		return ecdsa.Verify(&key.ecdsaKey.PublicKey, hash, sig.R, sig.S), nil
	}

	return false, ErrInvalidKeySpec
}

func (c *cryptoOps) generateMAC(keyID string, message []byte, algorithm MACAlgorithm) ([]byte, error) {
	key, exists := c.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrInvalidKeySpec
	}

	var hmacHash []byte
	switch algorithm {
	case MACAlgorithmHMACSHA224:
		h := hmac.New(sha256.New224, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	case MACAlgorithmHMACSHA256:
		h := hmac.New(sha256.New, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	case MACAlgorithmHMACSHA384:
		h := hmac.New(sha512.New384, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	case MACAlgorithmHMACSHA512:
		h := hmac.New(sha512.New, key.symmetric)
		h.Write(message)
		hmacHash = h.Sum(nil)
	default:
		return nil, ErrInvalidAlgorithm
	}

	return hmacHash, nil
}

func (c *cryptoOps) verifyMAC(keyID string, message, macValue []byte, algorithm MACAlgorithm) (bool, error) {
	key, exists := c.keys[keyID]
	if !exists {
		return false, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return false, ErrInvalidKeySpec
	}

	var expectedMAC []byte
	switch algorithm {
	case MACAlgorithmHMACSHA224:
		h := hmac.New(sha256.New224, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	case MACAlgorithmHMACSHA256:
		h := hmac.New(sha256.New, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	case MACAlgorithmHMACSHA384:
		h := hmac.New(sha512.New384, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	case MACAlgorithmHMACSHA512:
		h := hmac.New(sha512.New, key.symmetric)
		h.Write(message)
		expectedMAC = h.Sum(nil)
	default:
		return false, ErrInvalidAlgorithm
	}

	return hmac.Equal(macValue, expectedMAC), nil
}

func (c *cryptoOps) getPublicKey(keyID string) ([]byte, error) {
	key, exists := c.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.rsaKey != nil {
		return encodePublicKey(&key.rsaKey.PublicKey)
	} else if key.ecdsaKey != nil {
		return encodePublicKey(&key.ecdsaKey.PublicKey)
	}

	return nil, ErrInvalidKeySpec
}

func (c *cryptoOps) isKeyAvailable(keyID string) bool {
	key, exists := c.keys[keyID]
	if !exists {
		return false
	}

	return key.symmetric != nil || key.rsaKey != nil || key.ecdsaKey != nil
}

func (c *cryptoOps) keyExists(keyID string) bool {
	_, exists := c.keys[keyID]
	return exists
}

func hashForAlgorithm(algorithm SigningAlgorithm, message []byte) []byte {
	switch algorithm {
	case SigningAlgorithmRSAPKCS1SHA256, SigningAlgorithmRSAPSSSHA256, SigningAlgorithmECDSASHA256:
		h := sha256.Sum256(message)
		return h[:]
	case SigningAlgorithmRSAPKCS1SHA384, SigningAlgorithmRSAPSSSHA384, SigningAlgorithmECDSASHA384:
		h := sha512.Sum384(message)
		return h[:]
	case SigningAlgorithmRSAPKCS1SHA512, SigningAlgorithmRSAPSSSHA512, SigningAlgorithmECDSASHA512:
		h := sha512.Sum512(message)
		return h[:]
	default:
		h := sha256.Sum256(message)
		return h[:]
	}
}

func hashAlgorithm(algorithm SigningAlgorithm) crypto.Hash {
	switch algorithm {
	case SigningAlgorithmRSAPKCS1SHA256, SigningAlgorithmRSAPSSSHA256:
		return crypto.SHA256
	case SigningAlgorithmRSAPKCS1SHA384, SigningAlgorithmRSAPSSSHA384:
		return crypto.SHA384
	case SigningAlgorithmRSAPKCS1SHA512, SigningAlgorithmRSAPSSSHA512:
		return crypto.SHA512
	default:
		return crypto.SHA256
	}
}

func extractKeyIDFromCiphertext(ciphertext []byte) (string, []byte, error) {
	if len(ciphertext) < 2 {
		return "", nil, errors.New("ciphertext too short")
	}
	keyIDLen := int(ciphertext[0])<<8 | int(ciphertext[1])
	if len(ciphertext) < 2+keyIDLen {
		return "", nil, errors.New("ciphertext too short for keyID")
	}
	keyID := string(ciphertext[2 : 2+keyIDLen])
	actualCiphertext := ciphertext[2+keyIDLen:]
	return keyID, actualCiphertext, nil
}

func (o *cryptoOps) importKeyMaterial(keyID string, keyMaterial []byte, keySpec KeySpec) (*memoryKey, error) {
	if _, exists := o.keys[keyID]; exists {
		return nil, ErrKeyAlreadyExists
	}
	key := &memoryKey{keySpec: keySpec}
	switch keySpec {
	case KeySpecSymmetricDefault, KeySpecHMAC224, KeySpecHMAC256, KeySpecHMAC384, KeySpecHMAC512:
		key.symmetric = make([]byte, len(keyMaterial))
		copy(key.symmetric, keyMaterial)
	default:
		return nil, ErrInvalidKeySpec
	}
	return key, nil
}

func (o *cryptoOps) generateDataKeyBytes(keySpec string, numberOfBytes int) ([]byte, error) {
	var keyBytes int
	if numberOfBytes > 0 {
		keyBytes = numberOfBytes
	} else {
		switch keySpec {
		case "AES_128":
			keyBytes = 16
		case "AES_256":
			keyBytes = 32
		default:
			keyBytes = 32
		}
	}
	plaintext := make([]byte, keyBytes)
	if _, err := rand.Read(plaintext); err != nil {
		return nil, err
	}
	return plaintext, nil
}

func prependKeyIDToCiphertext(keyID string, ciphertext []byte) []byte {
	keyIDBytes := []byte(keyID)
	result := make([]byte, 2+len(keyIDBytes)+len(ciphertext))
	result[0] = byte(len(keyIDBytes) >> 8)
	result[1] = byte(len(keyIDBytes))
	copy(result[2:], keyIDBytes)
	copy(result[2+len(keyIDBytes):], ciphertext)
	return result
}
