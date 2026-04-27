package hsm

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/chacha20poly1305"
)

// PersistentBackend provides persistent HSM-backed key storage.
type PersistentBackend struct {
	cryptoOps
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
		dataPath: dataPath,
	}
	b.initKeys()

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

	key, err := b.generateKeyMaterial(keySpec)
	if err != nil {
		return fmt.Errorf("generating %s key: %w", keySpec, err)
	}

	b.keys[keyID] = key
	return b.saveKeys()
}

// ImportKey imports an existing key into the HSM.
func (b *PersistentBackend) ImportKey(keyID string, keyMaterial []byte, keySpec KeySpec) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	key, err := b.importKeyMaterial(keyID, keyMaterial, keySpec)
	if err != nil {
		return err
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
	result, err := b.encrypt(keyID, plaintext, context)
	if err != nil {
		return nil, fmt.Errorf("encrypting plaintext: %w", err)
	}
	return result, nil
}

// Decrypt decrypts ciphertext using an HSM key.
func (b *PersistentBackend) Decrypt(keyID string, ciphertext []byte, context map[string]string) (*DecryptResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result, err := b.decrypt(keyID, ciphertext, context)
	if err != nil {
		return nil, fmt.Errorf("decrypting ciphertext: %w", err)
	}
	return result, nil
}

// DecryptWithoutKeyID decrypts ciphertext without specifying a key ID.
func (b *PersistentBackend) DecryptWithoutKeyID(ciphertext []byte, context map[string]string) (*DecryptResult, string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result, keyID, err := b.decryptWithoutKeyID(ciphertext, context)
	if err != nil {
		return nil, "", fmt.Errorf("decrypting ciphertext: %w", err)
	}
	return result, keyID, nil
}

// Sign signs a message using an HSM key.
func (b *PersistentBackend) Sign(keyID string, message []byte, algorithm SigningAlgorithm) (*SignResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.sign(keyID, message, algorithm)
}

// Verify verifies a signature using an HSM key.
func (b *PersistentBackend) Verify(keyID string, message, signature []byte, algorithm SigningAlgorithm) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.verify(keyID, message, signature, algorithm)
}

// GenerateMAC generates a message authentication code using an HSM key.
func (b *PersistentBackend) GenerateMAC(keyID string, message []byte, algorithm MACAlgorithm) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.generateMAC(keyID, message, algorithm)
}

// VerifyMAC verifies a message authentication code using an HSM key.
func (b *PersistentBackend) VerifyMAC(keyID string, message, macValue []byte, algorithm MACAlgorithm) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.verifyMAC(keyID, message, macValue, algorithm)
}

// GenerateDataKey generates a data key for envelope encryption.
func (b *PersistentBackend) GenerateDataKey(keyID string, keySpec string, numberOfBytes int, encryptionContext map[string]string) (*GenerateDataKeyResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	key, exists := b.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}

	if key.symmetric == nil {
		return nil, ErrEncryptFailed
	}

	plaintext, err := b.generateDataKeyBytes(keySpec, numberOfBytes)
	if err != nil {
		return nil, err
	}

	ciphertext, err := aesEncrypt(key.symmetric, plaintext, serializeEncryptionContext(encryptionContext))
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

// GetPublicKey retrieves the public key from an HSM key.
func (b *PersistentBackend) GetPublicKey(keyID string) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.getPublicKey(keyID)
}

// IsKeyAvailable checks if a key is available in the HSM.
func (b *PersistentBackend) IsKeyAvailable(keyID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isKeyAvailable(keyID)
}

// KeyExists checks if a key exists in the HSM.
func (b *PersistentBackend) KeyExists(keyID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.keyExists(keyID)
}
