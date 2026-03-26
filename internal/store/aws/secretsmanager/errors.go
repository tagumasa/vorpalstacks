package secretsmanager

import (
	"errors"
)

var (
	// ErrSecretNotFound is returned when the specified secret does not exist.
	ErrSecretNotFound = errors.New("secret not found")

	// ErrSecretAlreadyExists is returned when attempting to create a secret
	// that already exists.
	ErrSecretAlreadyExists = errors.New("secret already exists")

	// ErrInvalidSecretName is returned when the secret name does not meet
	// Secrets Manager naming requirements.
	ErrInvalidSecretName = errors.New("invalid secret name")

	// ErrInvalidVersionId is returned when the specified version ID is not
	// valid for the secret.
	ErrInvalidVersionId = errors.New("invalid version id")
)
