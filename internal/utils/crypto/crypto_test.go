package crypto

import (
	"bytes"
	"testing"

	"vorpalstacks/internal/common/auth"
)

func TestRandomBytes(t *testing.T) {
	bytes1, err := RandomBytes(16)
	if err != nil {
		t.Fatalf("RandomBytes failed: %v", err)
	}

	if len(bytes1) != 16 {
		t.Errorf("RandomBytes returned wrong length: got %d, want 16", len(bytes1))
	}

	bytes2, err := RandomBytes(16)
	if err != nil {
		t.Fatalf("RandomBytes failed: %v", err)
	}

	if bytes.Equal(bytes1, bytes2) {
		t.Error("RandomBytes returned identical values (unlikely)")
	}
}

func TestAESGCMEncryptDecrypt(t *testing.T) {
	key, err := RandomBytes(32)
	if err != nil {
		t.Fatalf("RandomBytes failed: %v", err)
	}

	plaintext := []byte("secret message")

	ciphertext, err := AESGCMEncrypt(key, plaintext)
	if err != nil {
		t.Fatalf("AESGCMEncrypt failed: %v", err)
	}

	decrypted, err := AESGCMDecrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("AESGCMDecrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("AESGCMDecrypt returned wrong value: got %s, want %s", decrypted, plaintext)
	}
}

func TestAESGCMDecryptWithWrongKey(t *testing.T) {
	key1, _ := RandomBytes(32)
	key2, _ := RandomBytes(32)

	plaintext := []byte("secret message")

	ciphertext, _ := AESGCMEncrypt(key1, plaintext)

	_, err := AESGCMDecrypt(key2, ciphertext)
	if err == nil {
		t.Error("AESGCMDecrypt should fail with wrong key")
	}
}

func TestCredentialsProvider(t *testing.T) {
	provider := auth.NewStaticCredentialsProvider(
		"AKIAIOSFODNN7EXAMPLE",
		"wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY",
		"us-east-1",
		"s3",
	)

	creds, err := provider.GetCredentials()
	if err != nil {
		t.Fatalf("GetCredentials failed: %v", err)
	}

	if creds.AccessKeyID != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("AccessKeyID mismatch: got %s", creds.AccessKeyID)
	}

	if creds.SecretAccessKey != "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY" {
		t.Errorf("SecretAccessKey mismatch: got %s", creds.SecretAccessKey)
	}

	if creds.Region != "us-east-1" {
		t.Errorf("Region mismatch: got %s", creds.Region)
	}

	if creds.Service != "s3" {
		t.Errorf("Service mismatch: got %s", creds.Service)
	}
}
