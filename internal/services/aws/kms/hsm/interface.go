package hsm

// Package hsm provides Hardware Security Module (HSM) interface definitions
// for vorpalstacks KMS operations.

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

// KeySpec defines the type of cryptographic key.
type KeySpec string

const (
	// KeySpecSymmetricDefault is the default symmetric key specification.
	KeySpecSymmetricDefault KeySpec = "SYMMETRIC_DEFAULT"
	// KeySpecHMAC224 specifies an HMAC key with 224-bit strength.
	KeySpecHMAC224 KeySpec = "HMAC_224"
	// KeySpecHMAC256 specifies an HMAC key with 256-bit strength.
	KeySpecHMAC256 KeySpec = "HMAC_256"
	// KeySpecHMAC384 specifies an HMAC key with 384-bit strength.
	KeySpecHMAC384 KeySpec = "HMAC_384"
	// KeySpecHMAC512 specifies an HMAC key with 512-bit strength.
	KeySpecHMAC512 KeySpec = "HMAC_512"
	// KeySpecSM2 specifies an SM2 key (Chinese national standard).
	KeySpecSM2 KeySpec = "SM2"
	// KeySpecRSA2048 specifies an RSA key with 2048-bit strength.
	KeySpecRSA2048 KeySpec = "RSA_2048"
	// KeySpecRSA3072 specifies an RSA key with 3072-bit strength.
	KeySpecRSA3072 KeySpec = "RSA_3072"
	// KeySpecRSA4096 specifies an RSA key with 4096-bit strength.
	KeySpecRSA4096 KeySpec = "RSA_4096"
	// KeySpecECCNISTP256 specifies an ECC NIST P-256 curve key.
	KeySpecECCNISTP256 KeySpec = "ECC_NIST_P256"
	// KeySpecECCNISTP384 specifies an ECC NIST P-384 curve key.
	KeySpecECCNISTP384 KeySpec = "ECC_NIST_P384"
	// KeySpecECCNISTP521 specifies an ECC NIST P-521 curve key.
	KeySpecECCNISTP521 KeySpec = "ECC_NIST_P521"
	// KeySpecECCSECPP256R1 specifies an ECC SECG P-256K1 curve key.
	KeySpecECCSECPP256R1 KeySpec = "ECC_SECG_P256K1"
	// KeySpecSM2ECC is an alias for SM2 (Chinese national standard).
	KeySpecSM2ECC KeySpec = "SM2"
)

// SigningAlgorithm defines the algorithm used for digital signatures.
type SigningAlgorithm string

const (
	// SigningAlgorithmRSAPKCS1SHA256 specifies RSA PKCS#1 v1.5 with SHA-256.
	SigningAlgorithmRSAPKCS1SHA256 SigningAlgorithm = "RSASSA_PKCS1_V1_5_SHA_256"
	// SigningAlgorithmRSAPKCS1SHA384 specifies RSA PKCS#1 v1.5 with SHA-384.
	SigningAlgorithmRSAPKCS1SHA384 SigningAlgorithm = "RSASSA_PKCS1_V1_5_SHA_384"
	// SigningAlgorithmRSAPKCS1SHA512 specifies RSA PKCS#1 v1.5 with SHA-512.
	SigningAlgorithmRSAPKCS1SHA512 SigningAlgorithm = "RSASSA_PKCS1_V1_5_SHA_512"
	// SigningAlgorithmRSAPSSSHA256 specifies RSA PSS with SHA-256.
	SigningAlgorithmRSAPSSSHA256 SigningAlgorithm = "RSASSA_PSS_SHA_256"
	// SigningAlgorithmRSAPSSSHA384 specifies RSA PSS with SHA-384.
	SigningAlgorithmRSAPSSSHA384 SigningAlgorithm = "RSASSA_PSS_SHA_384"
	// SigningAlgorithmRSAPSSSHA512 specifies RSA PSS with SHA-512.
	SigningAlgorithmRSAPSSSHA512 SigningAlgorithm = "RSASSA_PSS_SHA_512"
	// SigningAlgorithmECDSASHA256 specifies ECDSA with SHA-256.
	SigningAlgorithmECDSASHA256 SigningAlgorithm = "ECDSA_SHA_256"
	// SigningAlgorithmECDSASHA384 specifies ECDSA with SHA-384.
	SigningAlgorithmECDSASHA384 SigningAlgorithm = "ECDSA_SHA_384"
	// SigningAlgorithmECDSASHA512 specifies ECDSA with SHA-512.
	SigningAlgorithmECDSASHA512 SigningAlgorithm = "ECDSA_SHA_512"
	// SigningAlgorithmSM2 specifies the SM2 signature algorithm.
	SigningAlgorithmSM2 SigningAlgorithm = "SM2"
)

// MACAlgorithm defines the algorithm used for message authentication codes.
type MACAlgorithm string

const (
	// MACAlgorithmHMACSHA224 specifies HMAC with SHA-224.
	MACAlgorithmHMACSHA224 MACAlgorithm = "HMAC_SHA_224"
	// MACAlgorithmHMACSHA256 specifies HMAC with SHA-256.
	MACAlgorithmHMACSHA256 MACAlgorithm = "HMAC_SHA_256"
	// MACAlgorithmHMACSHA384 specifies HMAC with SHA-384.
	MACAlgorithmHMACSHA384 MACAlgorithm = "HMAC_SHA_384"
	// MACAlgorithmHMACSHA512 specifies HMAC with SHA-512.
	MACAlgorithmHMACSHA512 MACAlgorithm = "HMAC_SHA_512"
)

// EncryptResult contains the result of an encryption operation.
type EncryptResult struct {
	Ciphertext      []byte
	CiphertextCRC32 uint32
}

// DecryptResult contains the result of a decryption operation.
type DecryptResult struct {
	Plaintext      []byte
	PlaintextCRC32 uint32
}

// SignResult contains the result of a signing operation.
type SignResult struct {
	Signature      []byte
	SignatureCRC32 uint32
}

// GenerateDataKeyResult contains the result of generating a data key.
type GenerateDataKeyResult struct {
	Plaintext      []byte
	Ciphertext     []byte
	PlaintextCRC32 uint32
	KeyID          string
}

// Backend defines the interface for HSM operations.
type Backend interface {
	GenerateKey(keyID string, keySpec KeySpec) error
	ImportKey(keyID string, keyMaterial []byte, keySpec KeySpec) error
	DeleteKey(keyID string) error

	Encrypt(keyID string, plaintext []byte, context map[string]string) (*EncryptResult, error)
	Decrypt(keyID string, ciphertext []byte, context map[string]string) (*DecryptResult, error)
	DecryptWithoutKeyID(ciphertext []byte, context map[string]string) (*DecryptResult, string, error)

	Sign(keyID string, message []byte, algorithm SigningAlgorithm) (*SignResult, error)
	Verify(keyID string, message, signature []byte, algorithm SigningAlgorithm) (bool, error)

	GenerateMAC(keyID string, message []byte, algorithm MACAlgorithm) ([]byte, error)
	VerifyMAC(keyID string, message, mac []byte, algorithm MACAlgorithm) (bool, error)

	GenerateDataKey(keyID string, keySpec string, numberOfBytes int) (*GenerateDataKeyResult, error)

	GetPublicKey(keyID string) ([]byte, error)
	IsKeyAvailable(keyID string) bool
	KeyExists(keyID string) bool
}

// Config holds the configuration for creating an HSM backend.
type Config struct {
	BackendType string
	LibraryPath string
	SlotID      int
	Pin         string
}

// ErrKeyNotFound is returned when a key is not found in the HSM.
var (
	ErrKeyNotFound = errors.New("key not found in HSM")
	// ErrKeyAlreadyExists is returned when a key already exists in the HSM.
	ErrKeyAlreadyExists = errors.New("key already exists in HSM")
	// ErrInvalidKeySpec is returned when an invalid key specification is provided.
	ErrInvalidKeySpec = errors.New("invalid key specification")
	// ErrEncryptFailed is returned when encryption fails.
	ErrEncryptFailed = errors.New("encryption failed")
	// ErrDecryptFailed is returned when decryption fails.
	ErrDecryptFailed = errors.New("decryption failed")
	// ErrSignFailed is returned when signing fails.
	ErrSignFailed = errors.New("signing failed")
	// ErrVerifyFailed is returned when signature verification fails.
	ErrVerifyFailed = errors.New("signature verification failed")
	// ErrInvalidAlgorithm is returned when an invalid algorithm is specified.
	ErrInvalidAlgorithm = errors.New("invalid algorithm")
)

// NewBackend creates a new HSM backend based on the configuration.
// Supported backend types are "memory" (default) and "pkcs11" (requires external library).
func NewBackend(cfg *Config) (Backend, error) {
	switch cfg.BackendType {
	case "memory", "":
		return NewMemoryBackend()
	case "softhsm", "pkcs11":
		return nil, fmt.Errorf("PKCS#11 backend requires external library, use memory backend for development")
	default:
		return nil, fmt.Errorf("unknown HSM backend type: %s", cfg.BackendType)
	}
}

func computeCRC32(data []byte) uint32 {
	var crc uint32 = 0xFFFFFFFF
	for _, b := range data {
		crc ^= uint32(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xEDB88320
			} else {
				crc >>= 1
			}
		}
	}
	return crc ^ 0xFFFFFFFF
}

func serializeEncryptionContext(context map[string]string) []byte {
	if len(context) == 0 {
		return nil
	}
	keys := make([]string, 0, len(context))
	for k := range context {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+context[k])
	}
	return []byte(strings.Join(parts, ";"))
}

func aesEncrypt(key, plaintext []byte, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, aad)
	return ciphertext, nil
}

func aesDecrypt(key, ciphertext []byte, aad []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func encodePublicKey(pubKey interface{}) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pubKey)
}
