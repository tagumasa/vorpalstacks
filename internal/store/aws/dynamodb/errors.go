// Package dynamodb provides DynamoDB storage functionality for vorpalstacks.
package dynamodb

import (
	"errors"
	"fmt"

	"vorpalstacks/internal/store/aws/common"
)

var (
	// ErrTableNotFound is returned when the specified DynamoDB table
	// does not exist. It wraps the common store not-found class so the
	// table's two read paths — the direct store read and the
	// transactional read — answer absence with one error that both
	// sentinel families recognise: errors.Is(err, ErrTableNotFound) and
	// the common IsNotFound.
	ErrTableNotFound = fmt.Errorf("table not found: %w", common.ErrNotFound)

	// ErrTableAlreadyExists is returned when attempting to create a table
	// that already exists.
	ErrTableAlreadyExists = errors.New("table already exists")

	// ErrTableNotActive is returned when the table is not in the ACTIVE state.
	ErrTableNotActive = errors.New("table not active")

	// errTableGenerationGone marks a metric flush whose carrying write
	// belonged to a superseded generation of the table name: a same-name
	// successor created between the write's commit and the flush must not
	// inherit the old generation's counters.
	errTableGenerationGone = errors.New("table generation replaced")

	// ErrItemNotFound is returned when the specified item does not exist
	// in the table.
	ErrItemNotFound = errors.New("item not found")

	// ErrItemAlreadyExists is returned when attempting to create an item
	// that already exists.
	ErrItemAlreadyExists = errors.New("item already exists")

	// ErrInvalidKey is returned when the key is not valid (e.g., missing
	// partition key or sort key).
	ErrInvalidKey = errors.New("invalid key")

	// ErrIndexNotFound is returned when the specified index does not exist.
	ErrIndexNotFound = errors.New("index not found")

	// ErrIndexAlreadyExists is returned when attempting to create an index
	// that already exists.
	ErrIndexAlreadyExists = errors.New("index already exists")

	// ErrBackupNotFound is returned when the specified backup does not exist.
	ErrBackupNotFound = errors.New("backup not found")

	// ErrPolicyRevisionMismatch is returned by the revision-checked
	// resource policy writes when the record's current revision no longer
	// equals the expected revision — the optimistic-lock loss the
	// ExpectedRevisionId parameter exists to detect.
	ErrPolicyRevisionMismatch = errors.New("resource policy revision mismatch")
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
