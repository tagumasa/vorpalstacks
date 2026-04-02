package hsm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/chacha20poly1305"
)

// PersistentBackend provides persistent HSM-backed key storage.
type PersistentBackend struct {
	mu        sync.RWMutex
	keys      map[string]*memoryKey
	masterKey [32]byte
	dataPath  string
}

type serializedKey struct {
	KeySpec      string `json:"key_spec"`
	Symmetric    []byte `json:"symmetric,omitempty"`
	RSAPrivate   []byte `json:"rsa_private,omitempty"`
	ECDSAPrivate []byte `json:"ecdsa_private,omitempty"`
}

// NewPersistentBackend creates a new persistent HSM backend.
func NewPersistentBackend(dataPath string) (*PersistentBackend, error) {
	if err := os.MkdirAll(filepath.Join(dataPath, "hsm"), 0700); err != nil {
		return nil, fmt.Errorf("failed to create hsm directory: %w", err)
	}

	b := &PersistentBackend{
		keys:     make(map[string]*memoryKey),
		dataPath: dataPath,
	}

	masterKeyPath := filepath.Join(dataPath, "hsm", "master.key")
	if _, err := os.Stat(masterKeyPath); os.IsNotExist(err) {
		if _, err := rand.Read(b.masterKey[:]); err != nil {
			return nil, fmt.Errorf("failed to generate master key: %w", err)
		}
		if err := b.saveMasterKey(masterKeyPath); err != nil {
			return nil, fmt.Errorf("failed to save master key: %w", err)
		}
	} else {
		if err := b.loadMasterKey(masterKeyPath); err != nil {
			return nil, fmt.Errorf("failed to load master key: %w", err)
		}
	}

	if err := b.loadKeys(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to load keys: %w", err)
	}

	return b, nil
}

func (b *PersistentBackend) saveMasterKey(path string) error {
	return os.WriteFile(path, b.masterKey[:], 0600)
}

func (b *PersistentBackend) loadMasterKey(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading master key file: %w", err)
	}

	if len(data) != 32 {
		return errors.New("invalid master key file: expected 32 bytes")
	}

	copy(b.masterKey[:], data)
	return nil
}

func (b *PersistentBackend) saveKeys() error {
	aead, err := chacha20poly1305.NewX(b.masterKey[:])
	if err != nil {
		return fmt.Errorf("creating AEAD cipher for key save: %w", err)
	}

	serialized := make(map[string]serializedKey)
	for keyID, mk := range b.keys {
		sk := serializedKey{KeySpec: string(mk.keySpec)}
		if mk.symmetric != nil {
			sk.Symmetric = mk.symmetric
		}
		if mk.rsaKey != nil {
			sk.RSAPrivate = x509.MarshalPKCS1PrivateKey(mk.rsaKey)
		}
		if mk.ecdsaKey != nil {
			ecdsaBytes, err := x509.MarshalECPrivateKey(mk.ecdsaKey)
			if err != nil {
				return fmt.Errorf("failed to marshal ECDSA private key for %s: %w", keyID, err)
			}
			sk.ECDSAPrivate = ecdsaBytes
		}
		serialized[keyID] = sk
	}

	plaintext, err := json.Marshal(serialized)
	if err != nil {
		return fmt.Errorf("marshaling serialized keys: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("generating encryption nonce: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)
	data := append(nonce, ciphertext...)

	keysPath := filepath.Join(b.dataPath, "hsm", "keys.enc")
	return os.WriteFile(keysPath, data, 0600)
}

func (b *PersistentBackend) loadKeys() error {
	keysPath := filepath.Join(b.dataPath, "hsm", "keys.enc")
	data, err := os.ReadFile(keysPath)
	if err != nil {
		return fmt.Errorf("reading encrypted keys file: %w", err)
	}

	if len(data) < chacha20poly1305.NonceSizeX {
		return errors.New("invalid keys file")
	}

	aead, err := chacha20poly1305.NewX(b.masterKey[:])
	if err != nil {
		return fmt.Errorf("creating AEAD cipher for key load: %w", err)
	}

	nonce := data[:aead.NonceSize()]
	ciphertext := data[aead.NonceSize():]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decrypting keys data: %w", err)
	}

	var serialized map[string]serializedKey
	if err := json.Unmarshal(plaintext, &serialized); err != nil {
		return fmt.Errorf("unmarshaling keys JSON: %w", err)
	}

	for keyID, sk := range serialized {
		mk := &memoryKey{keySpec: KeySpec(sk.KeySpec)}
		if sk.Symmetric != nil {
			mk.symmetric = sk.Symmetric
		}
		if sk.RSAPrivate != nil {
			rsaKey, err := x509.ParsePKCS1PrivateKey(sk.RSAPrivate)
			if err != nil {
				return fmt.Errorf("failed to parse RSA private key for %s: %w", keyID, err)
			}
			mk.rsaKey = rsaKey
		}
		if sk.ECDSAPrivate != nil {
			ecdsaKey, err := x509.ParseECPrivateKey(sk.ECDSAPrivate)
			if err != nil {
				return fmt.Errorf("failed to parse ECDSA private key for %s: %w", keyID, err)
			}
			mk.ecdsaKey = ecdsaKey
		}
		b.keys[keyID] = mk
	}

	return nil
}

// GenerateKey generates a new key in the HSM.
func (b *PersistentBackend) GenerateKey(keyID string, keySpec KeySpec) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.keys[keyID]; exists {
		return ErrKeyAlreadyExists
	}

	key := &memoryKey{keySpec: keySpec}

	switch keySpec {
	case KeySpecSymmetricDefault:
		key.symmetric = make([]byte, 32)
		if _, err := rand.Read(key.symmetric); err != nil {
			return fmt.Errorf("generating symmetric key: %w", err)
		}
	case KeySpecHMAC224:
		key.symmetric = make([]byte, 28)
		if _, err := rand.Read(key.symmetric); err != nil {
			return fmt.Errorf("generating HMAC-224 key: %w", err)
		}
	case KeySpecHMAC256:
		key.symmetric = make([]byte, 32)
		if _, err := rand.Read(key.symmetric); err != nil {
			return fmt.Errorf("generating HMAC-256 key: %w", err)
		}
	case KeySpecHMAC384:
		key.symmetric = make([]byte, 48)
		if _, err := rand.Read(key.symmetric); err != nil {
			return fmt.Errorf("generating HMAC-384 key: %w", err)
		}
	case KeySpecHMAC512:
		key.symmetric = make([]byte, 64)
		if _, err := rand.Read(key.symmetric); err != nil {
			return fmt.Errorf("generating HMAC-512 key: %w", err)
		}
	case KeySpecRSA2048:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return fmt.Errorf("generating RSA-2048 key pair: %w", err)
		}
		key.rsaKey = rsaKey
	case KeySpecRSA3072:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 3072)
		if err != nil {
			return fmt.Errorf("generating RSA-3072 key pair: %w", err)
		}
		key.rsaKey = rsaKey
	case KeySpecRSA4096:
		rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return fmt.Errorf("generating RSA-4096 key pair: %w", err)
		}
		key.rsaKey = rsaKey
	case KeySpecECCNISTP256, KeySpecECCSECPP256R1:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return fmt.Errorf("generating ECDSA P-256 key pair: %w", err)
		}
		key.ecdsaKey = ecdsaKey
	case KeySpecECCNISTP384:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return fmt.Errorf("generating ECDSA P-384 key pair: %w", err)
		}
		key.ecdsaKey = ecdsaKey
	case KeySpecECCNISTP521:
		ecdsaKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return fmt.Errorf("generating ECDSA P-521 key pair: %w", err)
		}
		key.ecdsaKey = ecdsaKey
	default:
		return ErrInvalidKeySpec
	}

	b.keys[keyID] = key
	return b.saveKeys()
}

// ImportKey imports an existing key into the HSM.
func (b *PersistentBackend) ImportKey(keyID string, keyMaterial []byte, keySpec KeySpec) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.keys[keyID]; exists {
		return ErrKeyAlreadyExists
	}

	key := &memoryKey{keySpec: keySpec}
	switch keySpec {
	case KeySpecSymmetricDefault, KeySpecHMAC256:
		key.symmetric = make([]byte, len(keyMaterial))
		copy(key.symmetric, keyMaterial)
	default:
		return ErrInvalidKeySpec
	}

	b.keys[keyID] = key
	return b.saveKeys()
}

// DeleteKey deletes a key from the HSM.
func (b *PersistentBackend) DeleteKey(keyID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.keys[keyID]; !exists {
		return ErrKeyNotFound
	}

	delete(b.keys, keyID)
	return b.saveKeys()
}

// Encrypt encrypts plaintext using an HSM key.
func (b *PersistentBackend) Encrypt(keyID string, plaintext []byte, context map[string]string) (*EncryptResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrEncryptFailed
	}

	ciphertext, err := aesEncrypt(key.symmetric, plaintext, serializeEncryptionContext(context))
	if err != nil {
		return nil, fmt.Errorf("encrypting plaintext: %w", err)
	}

	ciphertextWithKeyID := prependKeyIDToCiphertext(keyID, ciphertext)

	return &EncryptResult{
		Ciphertext:      ciphertextWithKeyID,
		CiphertextCRC32: computeCRC32(ciphertextWithKeyID),
	}, nil
}

// Decrypt decrypts ciphertext using an HSM key.
func (b *PersistentBackend) Decrypt(keyID string, ciphertext []byte, context map[string]string) (*DecryptResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	extractedKeyID, actualCiphertext, err := extractKeyIDFromCiphertext(ciphertext)
	if err == nil && extractedKeyID != "" {
		keyID = extractedKeyID
		ciphertext = actualCiphertext
	}

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrDecryptFailed
	}

	plaintext, err := aesDecrypt(key.symmetric, ciphertext, serializeEncryptionContext(context))
	if err != nil {
		return nil, fmt.Errorf("decrypting ciphertext: %w", err)
	}

	return &DecryptResult{
		Plaintext:      plaintext,
		PlaintextCRC32: computeCRC32(plaintext),
	}, nil
}

// DecryptWithoutKeyID decrypts ciphertext without specifying a key ID.
func (b *PersistentBackend) DecryptWithoutKeyID(ciphertext []byte, context map[string]string) (*DecryptResult, string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	keyID, actualCiphertext, err := extractKeyIDFromCiphertext(ciphertext)
	if err != nil {
		return nil, "", fmt.Errorf("extracting key ID from ciphertext: %w", err)
	}

	key, exists := b.keys[keyID]
	if !exists {
		return nil, "", ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, "", ErrDecryptFailed
	}

	plaintext, err := aesDecrypt(key.symmetric, actualCiphertext, serializeEncryptionContext(context))
	if err != nil {
		return nil, "", fmt.Errorf("decrypting ciphertext: %w", err)
	}

	return &DecryptResult{
		Plaintext:      plaintext,
		PlaintextCRC32: computeCRC32(plaintext),
	}, keyID, nil
}

// GenerateDataKey generates a data key for envelope encryption.
func (b *PersistentBackend) GenerateDataKey(keyID string, keySpec string, numberOfBytes int) (*GenerateDataKeyResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

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

	ciphertext, err := aesEncrypt(key.symmetric, plaintext, nil)
	if err != nil {
		return nil, err
	}

	ciphertextWithKeyID := prependKeyIDToCiphertext(keyID, ciphertext)

	return &GenerateDataKeyResult{
		Plaintext:      plaintext,
		Ciphertext:     ciphertextWithKeyID,
		PlaintextCRC32: computeCRC32(plaintext),
		KeyID:          keyID,
	}, nil
}

// KeyExists checks if a key exists in the HSM.
func (b *PersistentBackend) KeyExists(keyID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, exists := b.keys[keyID]
	return exists
}

// IsKeyAvailable checks if a key is available in the HSM.
func (b *PersistentBackend) IsKeyAvailable(keyID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return false
	}

	return key.symmetric != nil || key.rsaKey != nil || key.ecdsaKey != nil
}

// Sign signs a message using an HSM key.
func (b *PersistentBackend) Sign(keyID string, message []byte, algorithm SigningAlgorithm) (*SignResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
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

// Verify verifies a signature using an HSM key.
func (b *PersistentBackend) Verify(keyID string, message, signature []byte, algorithm SigningAlgorithm) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
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

// GenerateMAC generates a message authentication code using an HSM key.
func (b *PersistentBackend) GenerateMAC(keyID string, message []byte, algorithm MACAlgorithm) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
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

// VerifyMAC verifies a message authentication code using an HSM key.
func (b *PersistentBackend) VerifyMAC(keyID string, message, macValue []byte, algorithm MACAlgorithm) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
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

	if len(macValue) != len(expectedMAC) {
		return false, nil
	}

	return hmac.Equal(macValue, expectedMAC), nil
}

// GetPublicKey retrieves the public key from an HSM key.
func (b *PersistentBackend) GetPublicKey(keyID string) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
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
