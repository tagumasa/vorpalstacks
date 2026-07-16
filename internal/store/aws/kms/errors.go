package kms

// Package kms provides KMS (Key Management Service) data store implementations
// for vorpalstacks.

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrKeyNotFound is returned when the specified KMS key does not exist.
	ErrKeyNotFound = errors.New("key not found")

	// ErrKeyAlreadyExists is returned when attempting to create a key that
	// already exists.
	ErrKeyAlreadyExists = errors.New("key already exists")

	// ErrAliasNotFound is returned when the specified key alias does not exist.
	ErrAliasNotFound = errors.New("alias not found")

	// ErrAliasAlreadyExists is returned when attempting to create an alias
	// that already exists.
	ErrAliasAlreadyExists = errors.New("alias already exists")

	// ErrGrantNotFound is returned when the specified grant does not exist.
	ErrGrantNotFound = errors.New("grant not found")

	// ErrGrantAlreadyExists is returned when attempting to create a grant
	// that already exists.
	ErrGrantAlreadyExists = errors.New("grant already exists")

	// ErrKeyPolicyNotFound is returned when the key policy cannot be found.
	ErrKeyPolicyNotFound = errors.New("key policy not found")

	// ErrKeyDisabled is returned when the key is disabled and the operation
	// requires an enabled key.
	ErrKeyDisabled = errors.New("key is disabled")

	// ErrKeyPendingDeletion is returned when the key is pending deletion
	// and cannot be used for cryptographic operations.
	ErrKeyPendingDeletion = errors.New("key is pending deletion")

	// ErrKeyPendingImport is returned when the key material is pending import.
	ErrKeyPendingImport = errors.New("key is pending import")

	// ErrInvalidKeyState is returned when the key state does not permit
	// the requested operation.
	ErrInvalidKeyState = errors.New("invalid key state for operation")

	// ErrInvalidKeyUsage is returned when the key usage is not compatible
	// with the requested operation.
	ErrInvalidKeyUsage = errors.New("invalid key usage")

	// ErrInvalidKeySpec is returned when the key specification is not valid.
	ErrInvalidKeySpec = errors.New("invalid key specification")

	// ErrInvalidAliasName is returned when the alias name does not meet
	// KMS naming requirements.
	ErrInvalidAliasName = errors.New("invalid alias name")

	// ErrInvalidGrantToken is returned when the grant token is not valid.
	ErrInvalidGrantToken = errors.New("invalid grant token")

	// ErrCustomKeyStoreNotFound is returned when the specified custom key
	// store does not exist.
	ErrCustomKeyStoreNotFound = errors.New("custom key store not found")

	// ErrImportKeyMaterialExpired is returned when the key material import
	// has expired.
	ErrImportKeyMaterialExpired = errors.New("import key material expired")

	// ErrInvalidKeyOrigin is returned when the key origin is not valid for
	// the requested operation.
	ErrInvalidKeyOrigin = errors.New("invalid key origin for this operation")

	// ErrNotMultiRegionKey is returned when the operation requires a
	// multi-region key but the key is a single-region key.
	ErrNotMultiRegionKey = errors.New("key is not a multi-region key")

	// ErrInvalidPendingWindowDays is returned when PendingWindowInDays is
	// outside the valid range of 7-30 days.
	ErrInvalidPendingWindowDays = errors.New("PendingWindowInDays must be between 7 and 30 days")

	// ErrAliasNameReserved is returned when attempting to create an alias
	// that starts with alias/aws/ which is reserved for AWS managed keys.
	ErrAliasNameReserved = errors.New("alias name cannot begin with alias/aws/")
)

// NewStoreError creates a new KMS store error with the given operation and error.
func NewStoreError(op string, err error) *common.StoreError {
	return common.NewStoreError("kms", op, err)
}

// IsNotFound checks if the given error is a "not found" type error for any
// KMS resource type. This includes key, alias, and grant not found errors.
func IsNotFound(err error) bool {
	return common.IsNotFound(err) ||
		errors.Is(err, ErrKeyNotFound) ||
		errors.Is(err, ErrAliasNotFound) ||
		errors.Is(err, ErrGrantNotFound)
}
