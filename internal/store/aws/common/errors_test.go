package common

import (
	"errors"
	"testing"
)

func TestStoreError(t *testing.T) {
	err := NewStoreError("s3", "get", ErrNotFound)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	expected := "s3 store get: not found"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestStoreErrorWithKey(t *testing.T) {
	err := NewStoreErrorWithKey("s3", "get", "my-bucket", ErrNotFound)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	expected := "s3 store get my-bucket: not found"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestStoreErrorUnwrap(t *testing.T) {
	err := NewStoreError("s3", "get", ErrNotFound)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Unwrap() failed to unwrap ErrNotFound")
	}
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
			if got := IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
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
			if got := IsAlreadyExists(tt.err); got != tt.want {
				t.Errorf("IsAlreadyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
