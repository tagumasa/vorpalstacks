package hsm

// Package hsm provides Hardware Security Module (HSM) interface definitions
// for vorpalstacks KMS operations.

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
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
	// KeySpecECCNISTEdwards25519 specifies an Edwards25519 curve key
	// (Ed25519 signature algorithm). Smithy com.amazonaws.kms#DataKeyPairSpec
	// includes ECC_NIST_EDWARDS25519 as a native member.
	KeySpecECCNISTEdwards25519 KeySpec = "ECC_NIST_EDWARDS25519"
	// KeySpecECCSECGP256K1 specifies an ECC SECG P-256K1 curve key
	// (secp256k1, the Bitcoin curve). Distinct from KeySpecECCNISTP256
	// which is the NIST/SECG P-256 (secp256r1) curve.
	KeySpecECCSECGP256K1 KeySpec = "ECC_SECG_P256K1"
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

// MessageType defines whether a message to be signed is raw data or a pre-computed digest.
type MessageType string

const (
	// MessageTypeRaw indicates that the message is raw data that the HSM must hash before signing.
	MessageTypeRaw MessageType = "RAW"
	// MessageTypeDigest indicates that the message is a pre-computed digest that the HSM signs directly.
	MessageTypeDigest MessageType = "DIGEST"
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
	// ReplicateKey copies the cryptographic material of an existing key into
	// a new key ID. Used by ReplicateKey to provision the replica key in the
	// HSM with the primary's material so that crypto operations on the
	// replica succeed. Returns ErrKeyAlreadyExists if destKeyID is already
	// populated, or ErrKeyNotFound if sourceKeyID does not exist.
	ReplicateKey(sourceKeyID, destKeyID string) error
	DeleteKey(keyID string) error

	Encrypt(keyID string, plaintext []byte, algorithm EncryptionAlgorithm, context map[string]string) (*EncryptResult, error)
	Decrypt(keyID string, ciphertext []byte, algorithm EncryptionAlgorithm, context map[string]string) (*DecryptResult, error)
	DecryptWithoutKeyID(ciphertext []byte, algorithm EncryptionAlgorithm, context map[string]string) (*DecryptResult, string, error)

	Sign(keyID string, message []byte, algorithm SigningAlgorithm, messageType MessageType) (*SignResult, error)
	Verify(keyID string, message, signature []byte, algorithm SigningAlgorithm, messageType MessageType) (bool, error)

	GenerateMAC(keyID string, message []byte, algorithm MACAlgorithm) ([]byte, error)
	VerifyMAC(keyID string, message, mac []byte, algorithm MACAlgorithm) (bool, error)

	GenerateDataKey(keyID string, keySpec string, numberOfBytes int, encryptionContext map[string]string) (*GenerateDataKeyResult, error)

	GetPublicKey(keyID string) ([]byte, error)
	GenerateKeyPair(keySpec KeySpec) (privateKeyDER []byte, publicKeyDER []byte, err error)
	// RotateKey generates new cryptographic material for the key,
	// preserving previous versions so that ciphertexts encrypted before
	// rotation remain decryptable. Used by RotateKeyOnDemand.
	RotateKey(keyID string) error
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
	// ErrInvalidCiphertext is returned when the ciphertext cannot be
	// parsed (missing or truncated keyID prefix). Surfacing this from the
	// HSM lets the service layer map it to InvalidCiphertextException
	// rather than the generic KMSInternalException that decrypt-failure
	// mapping produces.
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	// ErrSignFailed is returned when signing fails.
	ErrSignFailed = errors.New("signing failed")
	// ErrVerifyFailed is returned when signature verification fails.
	ErrVerifyFailed = errors.New("signature verification failed")
	// ErrInvalidAlgorithm is returned when an invalid algorithm is specified.
	ErrInvalidAlgorithm = errors.New("invalid algorithm")
	// ErrInvalidDigestLength is returned when a DIGEST message type has an unexpected length for the signing algorithm.
	ErrInvalidDigestLength = errors.New("invalid digest length for the specified signing algorithm")
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

// EncryptionAlgorithm defines the asymmetric or symmetric algorithm used
// to encrypt plaintext under a KMS key.
type EncryptionAlgorithm string

const (
	// EncryptionAlgorithmSymmetricDefault is the only valid algorithm for
	// SYMMETRIC_DEFAULT keys.
	EncryptionAlgorithmSymmetricDefault EncryptionAlgorithm = "SYMMETRIC_DEFAULT"
	// EncryptionAlgorithmRSAOAEPSHA1 is RSAES-OAEP with SHA-1.
	EncryptionAlgorithmRSAOAEPSHA1 EncryptionAlgorithm = "RSAES_OAEP_SHA_1"
	// EncryptionAlgorithmRSAOAEPSHA256 is RSAES-OAEP with SHA-256.
	EncryptionAlgorithmRSAOAEPSHA256 EncryptionAlgorithm = "RSAES_OAEP_SHA_256"
)

// resolveRSAHash returns the crypto.Hash to use for the requested OAEP
// algorithm. SHA-1 is intentionally supported because AWS KMS still
// exposes RSAES_OAEP_SHA_1 on RSA keys.
func resolveRSAHash(algorithm EncryptionAlgorithm) crypto.Hash {
	switch algorithm {
	case EncryptionAlgorithmRSAOAEPSHA1:
		return crypto.SHA1
	case EncryptionAlgorithmRSAOAEPSHA256:
		return crypto.SHA256
	default:
		return crypto.SHA256
	}
}

// rsaEncrypt encrypts plaintext under an RSA public key using the supplied
// OAEP algorithm. AWS KMS only supports OAEP padding for RSA encryption;
// PKCS#1 v1.5 is not exposed.
func rsaEncrypt(pub *rsa.PublicKey, plaintext []byte, algorithm EncryptionAlgorithm) ([]byte, error) {
	hash := resolveRSAHash(algorithm)
	if err := validateRSAPlaintextLength(pub, plaintext, hash); err != nil {
		return nil, err
	}
	return rsa.EncryptOAEP(hash.New(), rand.Reader, pub, plaintext, nil)
}

// rsaDecrypt decrypts ciphertext that was produced by rsaEncrypt.
func rsaDecrypt(priv *rsa.PrivateKey, ciphertext []byte, algorithm EncryptionAlgorithm) ([]byte, error) {
	hash := resolveRSAHash(algorithm)
	return rsa.DecryptOAEP(hash.New(), rand.Reader, priv, ciphertext, nil)
}

// validateRSAPlaintextLength mirrors the AWS KMS plaintext length limits
// for RSA encryption. RSAES_OAEP_SHA_256 supports up to 190 (RSA_2048),
// 318 (RSA_3072) and 446 (RSA_4096) bytes; RSAES_OAEP_SHA_1 supports up
// to 214/342/470 bytes respectively. AWS rejects oversized plaintext
// with ValidationException, which we surface as ErrEncryptFailed.
func validateRSAPlaintextLength(pub *rsa.PublicKey, plaintext []byte, hash crypto.Hash) error {
	maxLen := pub.Size() - 2*hash.Size() - 2
	if maxLen < 0 || len(plaintext) > maxLen {
		return fmt.Errorf("%w: plaintext exceeds %d bytes for %d-bit RSA with %s", ErrEncryptFailed, maxLen, pub.Size()*8, hash)
	}
	return nil
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
