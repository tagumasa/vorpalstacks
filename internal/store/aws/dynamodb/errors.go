// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"errors"
)

var (
	// ErrTableNotFound is returned when the specified DynamoDB table
	// does not exist.
	ErrTableNotFound = errors.New("table not found")

	// ErrTableAlreadyExists is returned when attempting to create a table
	// that already exists.
	ErrTableAlreadyExists = errors.New("table already exists")

	// ErrTableNotActive is returned when the table is not in the ACTIVE state.
	ErrTableNotActive = errors.New("table not active")

	// ErrItemNotFound is returned when the specified item does not exist
	// in the table.
	ErrItemNotFound = errors.New("item not found")

	// ErrItemAlreadyExists is returned when attempting to create an item
	// that already exists.
	ErrItemAlreadyExists = errors.New("item already exists")

	// ErrInvalidKey is returned when the key is not valid (e.g., missing
	// partition key or sort key).
	ErrInvalidKey = errors.New("invalid key")

	// ErrInvalidAttributeType is returned when the attribute type is not valid
	// (e.g., S, N, B for string, number, binary).
	ErrInvalidAttributeType = errors.New("invalid attribute type")

	// ErrMissingKeyAttribute is returned when a required key attribute
	// is missing from the item.
	ErrMissingKeyAttribute = errors.New("missing key attribute")

	// ErrIndexNotFound is returned when the specified index does not exist.
	ErrIndexNotFound = errors.New("index not found")

	// ErrIndexAlreadyExists is returned when attempting to create an index
	// that already exists.
	ErrIndexAlreadyExists = errors.New("index already exists")

	// ErrTTLNotFound is returned when the TTL attribute is not configured
	// for the table.
	ErrTTLNotFound = errors.New("ttl not found")

	// ErrBackupNotFound is returned when the specified backup does not exist.
	ErrBackupNotFound = errors.New("backup not found")
)

// IsTableNotFound checks if the error indicates that a DynamoDB table
// was not found.
func IsTableNotFound(err error) bool {
	return errors.Is(err, ErrTableNotFound)
}

// IsTableAlreadyExists checks if the error indicates that a DynamoDB table
// already exists.
func IsTableAlreadyExists(err error) bool {
	return errors.Is(err, ErrTableAlreadyExists)
}

// IsItemNotFound checks if the error indicates that a DynamoDB item
// was not found.
func IsItemNotFound(err error) bool {
	return errors.Is(err, ErrItemNotFound)
}

// IsItemAlreadyExists checks if the error indicates that a DynamoDB item
// already exists.
func IsItemAlreadyExists(err error) bool {
	return errors.Is(err, ErrItemAlreadyExists)
}
