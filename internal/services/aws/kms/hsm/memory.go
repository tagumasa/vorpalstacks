// Package hsm provides Hardware Security Module (HSM) in-memory implementations
// for vorpalstacks KMS operations.
package hsm

// MemoryBackend is an in-memory implementation of a Hardware Security Module (HSM)
// backend for key storage and cryptographic operations.
type MemoryBackend struct {
	cryptoOps
}

// NewMemoryBackend creates a new in-memory HSM backend for key storage.
// This is useful for development and testing purposes.
func NewMemoryBackend() (*MemoryBackend, error) {
	b := &MemoryBackend{}
	b.initKeys()
	return b, nil
}

// GenerateKey generates a new key with the specified key ID and key specification.
// Supports symmetric, HMAC, RSA, and ECDSA key types.
func (b *MemoryBackend) GenerateKey(keyID string, keySpec KeySpec) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.keys[keyID]; exists {
		return ErrKeyAlreadyExists
	}

	key, err := b.generateKeyMaterial(keySpec)
	if err != nil {
		return err
	}

	b.keys[keyID] = key
	return nil
}

// ImportKey imports an existing key into the in-memory HSM backend.
// Only supports symmetric and HMAC key specifications.
func (b *MemoryBackend) ImportKey(keyID string, keyMaterial []byte, keySpec KeySpec) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	key, err := b.importKeyMaterial(keyID, keyMaterial, keySpec)
	if err != nil {
		return err
	}
	b.keys[keyID] = key
	return nil
}

// ReplicateKey copies the cryptographic material of an existing key into a
// new key ID. The replica shares the primary's secret material so crypto
// operations on the replica succeed. Required because ImportKey only
// supports symmetric/HMAC material; RSA/ECDSA private keys cannot be
// round-tripped through ImportKey.
func (b *MemoryBackend) ReplicateKey(sourceKeyID, destKeyID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	src, exists := b.keys[sourceKeyID]
	if !exists {
		return ErrKeyNotFound
	}
	if _, exists := b.keys[destKeyID]; exists {
		return ErrKeyAlreadyExists
	}

	replica := &memoryKey{keySpec: src.keySpec}
	if src.symmetric != nil {
		replica.symmetric = make([]byte, len(src.symmetric))
		copy(replica.symmetric, src.symmetric)
	}
	if src.rsaKey != nil {
		replica.rsaKey = src.rsaKey
	}
	if src.ecdsaKey != nil {
		replica.ecdsaKey = src.ecdsaKey
	}
	b.keys[destKeyID] = replica
	return nil
}

// DeleteKey removes a key from the in-memory HSM backend.
// Returns an error if the key does not exist.
func (b *MemoryBackend) DeleteKey(keyID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.keys[keyID]; !exists {
		return ErrKeyNotFound
	}

	delete(b.keys, keyID)
	return nil
}

// Encrypt encrypts plaintext using the specified symmetric or RSA key.
// Returns the ciphertext with key ID prepended and CRC32 checksum.
func (b *MemoryBackend) Encrypt(keyID string, plaintext []byte, algorithm EncryptionAlgorithm, context map[string]string) (*EncryptResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.encrypt(keyID, plaintext, algorithm, context)
}

// Decrypt decrypts ciphertext using the specified symmetric or RSA key.
// Extracts the key ID from the ciphertext if present.
func (b *MemoryBackend) Decrypt(keyID string, ciphertext []byte, algorithm EncryptionAlgorithm, context map[string]string) (*DecryptResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.decrypt(keyID, ciphertext, algorithm, context)
}

// DecryptWithoutKeyID decrypts ciphertext and extracts the key ID from the ciphertext.
// Returns the plaintext, key ID used, and any error encountered.
func (b *MemoryBackend) DecryptWithoutKeyID(ciphertext []byte, algorithm EncryptionAlgorithm, context map[string]string) (*DecryptResult, string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.decryptWithoutKeyID(ciphertext, algorithm, context)
}

// Sign signs a message using the specified key and signing algorithm.
// Supports RSA and ECDSA signing algorithms.
func (b *MemoryBackend) Sign(keyID string, message []byte, algorithm SigningAlgorithm, messageType MessageType) (*SignResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.sign(keyID, message, algorithm, messageType)
}

// Verify verifies a digital signature for a message using the specified key.
// Returns true if the signature is valid, false otherwise.
func (b *MemoryBackend) Verify(keyID string, message, signature []byte, algorithm SigningAlgorithm, messageType MessageType) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.verify(keyID, message, signature, algorithm, messageType)
}

// GenerateMAC generates a Message Authentication Code (MAC) for a message
// using the specified symmetric key and MAC algorithm.
func (b *MemoryBackend) GenerateMAC(keyID string, message []byte, algorithm MACAlgorithm) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.generateMAC(keyID, message, algorithm)
}

// VerifyMAC verifies a Message Authentication Code (MAC) for a message
// using the specified symmetric key and MAC algorithm.
func (b *MemoryBackend) VerifyMAC(keyID string, message, macValue []byte, algorithm MACAlgorithm) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.verifyMAC(keyID, message, macValue, algorithm)
}

// GenerateDataKey generates a new data key for envelope encryption.
// Returns both plaintext and encrypted versions of the data key.
func (b *MemoryBackend) GenerateDataKey(keyID string, keySpec string, numberOfBytes int, encryptionContext map[string]string) (*GenerateDataKeyResult, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.cryptoOps.generateDataKey(keyID, keySpec, numberOfBytes, encryptionContext)
}

// GetPublicKey retrieves the public key for an RSA or ECDSA key.
// Returns the encoded public key bytes.
func (b *MemoryBackend) GetPublicKey(keyID string) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.getPublicKey(keyID)
}

// GenerateKeyPair generates an asymmetric key pair and returns DER-encoded keys.
// The private key is PKCS8 and the public key is PKIX.
func (b *MemoryBackend) GenerateKeyPair(keySpec KeySpec) ([]byte, []byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.generateKeyPairDER(keySpec)
}

// IsKeyAvailable checks whether a key is available for cryptographic operations.
// A key is considered available if it has key material (symmetric, RSA, or ECDSA).
func (b *MemoryBackend) IsKeyAvailable(keyID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isKeyAvailable(keyID)
}

// KeyExists checks whether a key with the specified ID exists in the backend.
func (b *MemoryBackend) KeyExists(keyID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.keyExists(keyID)
}

// RotateKey generates new key material, preserving the old version.
func (b *MemoryBackend) RotateKey(keyID string) error {
	return b.cryptoOps.rotateKey(keyID)
}
