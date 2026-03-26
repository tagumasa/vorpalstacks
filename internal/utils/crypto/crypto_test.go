package crypto

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestSHA256Hash(t *testing.T) {
	data := []byte("hello world")
	hash := SHA256Hash(data)

	if hash == "" {
		t.Error("SHA256Hash returned empty string")
	}

	if len(hash) != 64 {
		t.Errorf("SHA256Hash returned wrong length: got %d, want 64", len(hash))
	}

	sameHash := SHA256Hash(data)
	if hash != sameHash {
		t.Error("SHA256Hash is not deterministic")
	}
}

func TestHMACSHA256(t *testing.T) {
	key := []byte("secret-key")
	data := []byte("message")

	mac := HMACSHA256(key, data)
	if len(mac) != 32 {
		t.Errorf("HMACSHA256 returned wrong length: got %d, want 32", len(mac))
	}
}

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

func TestRSAEncryptDecrypt(t *testing.T) {
	privKey, err := GenerateRSAKey(2048)
	if err != nil {
		t.Fatalf("GenerateRSAKey failed: %v", err)
	}

	plaintext := []byte("secret message")

	ciphertext, err := RSAEncrypt(&privKey.PublicKey, plaintext)
	if err != nil {
		t.Fatalf("RSAEncrypt failed: %v", err)
	}

	decrypted, err := RSADecrypt(privKey, ciphertext)
	if err != nil {
		t.Fatalf("RSADecrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("RSADecrypt returned wrong value: got %s, want %s", decrypted, plaintext)
	}
}

func TestRSASignVerify(t *testing.T) {
	privKey, err := GenerateRSAKey(2048)
	if err != nil {
		t.Fatalf("GenerateRSAKey failed: %v", err)
	}

	message := []byte("message to sign")

	signature, err := RSASign(privKey, message)
	if err != nil {
		t.Fatalf("RSASign failed: %v", err)
	}

	err = RSAVerify(&privKey.PublicKey, message, signature)
	if err != nil {
		t.Fatalf("RSAVerify failed: %v", err)
	}
}

func TestRSAPemRoundTrip(t *testing.T) {
	privKey, err := GenerateRSAKey(2048)
	if err != nil {
		t.Fatalf("GenerateRSAKey failed: %v", err)
	}

	privPEM := RSAPrivateKeyToPEM(privKey)
	pubPEM := RSAPublicKeyToPEM(&privKey.PublicKey)

	parsedPriv, err := PEMToRSAPrivateKey(privPEM)
	if err != nil {
		t.Fatalf("PEMToRSAPrivateKey failed: %v", err)
	}

	parsedPub, err := PEMToRSAPublicKey(pubPEM)
	if err != nil {
		t.Fatalf("PEMToRSAPublicKey failed: %v", err)
	}

	if parsedPriv.N.Cmp(privKey.N) != 0 {
		t.Error("Private key N mismatch after PEM round trip")
	}

	if parsedPub.N.Cmp(privKey.PublicKey.N) != 0 {
		t.Error("Public key N mismatch after PEM round trip")
	}
}

func TestBuildCanonicalRequest(t *testing.T) {
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme:   "https",
			Host:     "example.com",
			Path:     "/test",
			RawQuery: "param=value",
		},
		Header: http.Header{
			"Host":       []string{"example.com"},
			"X-Amz-Date": []string{"20230101T000000Z"},
		},
	}

	signedHeaders := "host;x-amz-date"
	body := []byte("")

	canonical := BuildCanonicalRequest(req, signedHeaders, body)

	if canonical == "" {
		t.Error("BuildCanonicalRequest returned empty string")
	}

	if !bytes.Contains([]byte(canonical), []byte("GET\n/test\nparam=value")) {
		t.Errorf("BuildCanonicalRequest missing expected parts: %q", canonical)
	}
}

func TestBuildStringToSign(t *testing.T) {
	amzDate := "20230101T000000Z"
	credential := "AKIAIOSFODNN7EXAMPLE/20230101/us-east-1/s3/aws4_request"
	canonicalRequest := "GET\n/\n\nhost:example.com\n\nhost\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	stringToSign, err := BuildStringToSign(amzDate, credential, canonicalRequest)
	if err != nil {
		t.Fatalf("BuildStringToSign failed: %v", err)
	}

	if stringToSign == "" {
		t.Error("BuildStringToSign returned empty string")
	}
}

func TestCalculateSignature(t *testing.T) {
	amzDate := "20230101T000000Z"
	region := "us-east-1"
	service := "s3"
	stringToSign := "AWS4-HMAC-SHA256\n20230101T000000Z\n20230101/us-east-1/s3/aws4_request\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	secretKey := "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"

	signature, err := CalculateSignature(amzDate, region, service, stringToSign, secretKey)
	if err != nil {
		t.Fatalf("CalculateSignature failed: %v", err)
	}

	if signature == "" {
		t.Error("CalculateSignature returned empty string")
	}

	if len(signature) != 64 {
		t.Errorf("CalculateSignature returned wrong length: got %d, want 64", len(signature))
	}
}

func TestCredentialsProvider(t *testing.T) {
	provider := NewStaticCredentialsProvider(
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
