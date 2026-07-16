package kms

import (
	"errors"
	"testing"
)

func TestKMSStoreErrors(t *testing.T) {
	t.Run("ErrKeyNotFound", func(t *testing.T) {
		if ErrKeyNotFound.Error() != "key not found" {
			t.Errorf("expected 'key not found', got %s", ErrKeyNotFound.Error())
		}
	})

	t.Run("ErrKeyAlreadyExists", func(t *testing.T) {
		if ErrKeyAlreadyExists.Error() != "key already exists" {
			t.Errorf("expected 'key already exists', got %s", ErrKeyAlreadyExists.Error())
		}
	})

	t.Run("ErrAliasNotFound", func(t *testing.T) {
		if ErrAliasNotFound.Error() != "alias not found" {
			t.Errorf("expected 'alias not found', got %s", ErrAliasNotFound.Error())
		}
	})

	t.Run("ErrAliasAlreadyExists", func(t *testing.T) {
		if ErrAliasAlreadyExists.Error() != "alias already exists" {
			t.Errorf("expected 'alias already exists', got %s", ErrAliasAlreadyExists.Error())
		}
	})

	t.Run("ErrGrantNotFound", func(t *testing.T) {
		if ErrGrantNotFound.Error() != "grant not found" {
			t.Errorf("expected 'grant not found', got %s", ErrGrantNotFound.Error())
		}
	})

	t.Run("ErrGrantAlreadyExists", func(t *testing.T) {
		if ErrGrantAlreadyExists.Error() != "grant already exists" {
			t.Errorf("expected 'grant already exists', got %s", ErrGrantAlreadyExists.Error())
		}
	})

	t.Run("ErrKeyPolicyNotFound", func(t *testing.T) {
		if ErrKeyPolicyNotFound.Error() != "key policy not found" {
			t.Errorf("expected 'key policy not found', got %s", ErrKeyPolicyNotFound.Error())
		}
	})

	t.Run("ErrKeyDisabled", func(t *testing.T) {
		if ErrKeyDisabled.Error() != "key is disabled" {
			t.Errorf("expected 'key is disabled', got %s", ErrKeyDisabled.Error())
		}
	})

	t.Run("ErrKeyPendingDeletion", func(t *testing.T) {
		if ErrKeyPendingDeletion.Error() != "key is pending deletion" {
			t.Errorf("expected 'key is pending deletion', got %s", ErrKeyPendingDeletion.Error())
		}
	})

	t.Run("ErrKeyPendingImport", func(t *testing.T) {
		if ErrKeyPendingImport.Error() != "key is pending import" {
			t.Errorf("expected 'key is pending import', got %s", ErrKeyPendingImport.Error())
		}
	})

	t.Run("ErrInvalidKeyState", func(t *testing.T) {
		if ErrInvalidKeyState.Error() != "invalid key state for operation" {
			t.Errorf("expected 'invalid key state for operation', got %s", ErrInvalidKeyState.Error())
		}
	})

	t.Run("ErrInvalidKeyUsage", func(t *testing.T) {
		if ErrInvalidKeyUsage.Error() != "invalid key usage" {
			t.Errorf("expected 'invalid key usage', got %s", ErrInvalidKeyUsage.Error())
		}
	})

	t.Run("ErrInvalidKeySpec", func(t *testing.T) {
		if ErrInvalidKeySpec.Error() != "invalid key specification" {
			t.Errorf("expected 'invalid key specification', got %s", ErrInvalidKeySpec.Error())
		}
	})

	t.Run("ErrInvalidAliasName", func(t *testing.T) {
		if ErrInvalidAliasName.Error() != "invalid alias name" {
			t.Errorf("expected 'invalid alias name', got %s", ErrInvalidAliasName.Error())
		}
	})

	t.Run("ErrInvalidGrantToken", func(t *testing.T) {
		if ErrInvalidGrantToken.Error() != "invalid grant token" {
			t.Errorf("expected 'invalid grant token', got %s", ErrInvalidGrantToken.Error())
		}
	})

	t.Run("ErrCustomKeyStoreNotFound", func(t *testing.T) {
		if ErrCustomKeyStoreNotFound.Error() != "custom key store not found" {
			t.Errorf("expected 'custom key store not found', got %s", ErrCustomKeyStoreNotFound.Error())
		}
	})

	t.Run("ErrImportKeyMaterialExpired", func(t *testing.T) {
		if ErrImportKeyMaterialExpired.Error() != "import key material expired" {
			t.Errorf("expected 'import key material expired', got %s", ErrImportKeyMaterialExpired.Error())
		}
	})

	t.Run("ErrInvalidKeyOrigin", func(t *testing.T) {
		if ErrInvalidKeyOrigin.Error() != "invalid key origin for this operation" {
			t.Errorf("expected 'invalid key origin for this operation', got %s", ErrInvalidKeyOrigin.Error())
		}
	})

	t.Run("ErrNotMultiRegionKey", func(t *testing.T) {
		if ErrNotMultiRegionKey.Error() != "key is not a multi-region key" {
			t.Errorf("expected 'key is not a multi-region key', got %s", ErrNotMultiRegionKey.Error())
		}
	})
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		found bool
	}{
		{"ErrKeyNotFound", ErrKeyNotFound, true},
		{"ErrAliasNotFound", ErrAliasNotFound, true},
		{"ErrGrantNotFound", ErrGrantNotFound, true},
		{"Nil error", nil, false},
		{"Random error", errors.New("random error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.found {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.found)
			}
		})
	}
}
