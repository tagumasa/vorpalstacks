package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionalCheckFailedError(t *testing.T) {
	err := NewConditionalCheckFailedError([]byte("key"), []byte("value"))
	assert.Equal(t, "conditional check failed", err.Error())
	assert.Nil(t, err.Unwrap())
}

func TestConditionalCheckFailedError_WithWrappedError(t *testing.T) {
	inner := errors.New("inner error")
	err := NewConditionalCheckFailedError([]byte("key"), []byte("value"))
	err.Err = inner

	assert.Equal(t, "conditional check failed: inner error", err.Error())
	assert.Equal(t, inner, err.Unwrap())
}

func TestVersionConflictError(t *testing.T) {
	err := &VersionConflictError{Key: []byte("key")}
	assert.Equal(t, "version conflict: failed to acquire unique version after retries", err.Error())
	assert.Nil(t, err.Unwrap())
}

func TestVersionConflictError_WithWrappedError(t *testing.T) {
	inner := errors.New("inner error")
	err := &VersionConflictError{Key: []byte("key"), Err: inner}

	assert.Equal(t, "version conflict: inner error", err.Error())
	assert.Equal(t, inner, err.Unwrap())
}
