// Package cognitoidentity provides Cognito Identity Pool storage functionality for vorpalstacks.
package cognitoidentity

import (
	"errors"
)

// Error definitions for Cognito Identity Pool operations.
var (
	// ErrIdentityPoolNotFound is returned when an Identity Pool is not found.
	ErrIdentityPoolNotFound = errors.New("identity pool not found")
	// ErrIdentityPoolAlreadyExists is returned when attempting to create an Identity Pool that already exists.
	ErrIdentityPoolAlreadyExists = errors.New("identity pool already exists")
	// ErrInvalidIdentityPoolName is returned when an invalid Identity Pool name is provided.
	ErrInvalidIdentityPoolName = errors.New("invalid identity pool name")
	// ErrInvalidIdentityPoolID is returned when an invalid Identity Pool ID is provided.
	ErrInvalidIdentityPoolID = errors.New("invalid identity pool id")
	// ErrIdentityNotFound is returned when an Identity is not found.
	ErrIdentityNotFound = errors.New("identity not found")
	// ErrIdentityAlreadyExists is returned when attempting to create an Identity that already exists.
	ErrIdentityAlreadyExists = errors.New("identity already exists")
)
