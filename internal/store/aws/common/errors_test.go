package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreError(t *testing.T) {
	err := NewStoreError("s3", "get", ErrNotFound)
	require.NotNil(t, err)
	assert.Equal(t, "s3 store get: not found", err.Error())
}

func TestStoreErrorWithKey(t *testing.T) {
	err := NewStoreErrorWithKey("s3", "get", "my-bucket", ErrNotFound)
	require.NotNil(t, err)
	assert.Equal(t, "s3 store get my-bucket: not found", err.Error())
}

func TestStoreErrorUnwrap(t *testing.T) {
	err := NewStoreError("s3", "get", ErrNotFound)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "StoreError with ErrNotFound",
			err:  NewStoreError("s3", "get", ErrNotFound),
			want: true,
		},
		{
			name: "StoreError with ErrNotFound as wrapped",
			err:  NewStoreErrorWithKey("s3", "get", "key", ErrNotFound),
			want: true,
		},
		{
			name: "StoreError with other error",
			err:  NewStoreError("s3", "get", ErrAlreadyExists),
			want: false,
		},
		{
			name: "Plain ErrNotFound",
			err:  ErrNotFound,
			want: true,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "Other error",
			err:  errors.New("some other error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsNotFound(tt.err))
		})
	}
}

func TestIsAlreadyExists(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "StoreError with ErrAlreadyExists",
			err:  NewStoreError("s3", "create", ErrAlreadyExists),
			want: true,
		},
		{
			name: "StoreError with other error",
			err:  NewStoreError("s3", "get", ErrNotFound),
			want: false,
		},
		{
			name: "Plain ErrAlreadyExists",
			err:  ErrAlreadyExists,
			want: true,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsAlreadyExists(tt.err))
		})
	}
}
