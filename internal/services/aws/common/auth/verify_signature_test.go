package auth

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCredentialsProvider struct {
	creds *Credentials
	err   error
}

func (m *mockCredentialsProvider) GetCredentials() (*Credentials, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.creds, nil
}

func TestSignatureV4Verifier_VerifyRequest_MissingAuthHeader(t *testing.T) {
	provider := &mockCredentialsProvider{
		creds: &Credentials{SecretAccessKey: "secret"},
	}
	verifier := NewSignatureV4Verifier(provider)

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	err := verifier.VerifyRequest(req, "s3", "us-east-1")
	assert.Equal(t, ErrMissingAuthHeader, err)
}

func TestSignatureV4Verifier_VerifyRequest_InvalidAlgorithm(t *testing.T) {
	provider := &mockCredentialsProvider{
		creds: &Credentials{SecretAccessKey: "secret"},
	}
	verifier := NewSignatureV4Verifier(provider)

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Set("Authorization", "AWS4-HMAC-SHA512 Credential=...")
	err := verifier.VerifyRequest(req, "s3", "us-east-1")
	assert.Equal(t, ErrInvalidAlgorithm, err)
}

func TestSignatureV4Verifier_VerifyRequest_MissingAuthFields(t *testing.T) {
	provider := &mockCredentialsProvider{
		creds: &Credentials{SecretAccessKey: "secret"},
	}
	verifier := NewSignatureV4Verifier(provider)

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=")
	err := verifier.VerifyRequest(req, "s3", "us-east-1")
	assert.Equal(t, ErrMissingAuthFields, err)
}

func TestSignatureV4Verifier_VerifyRequest_MissingDateHeader(t *testing.T) {
	provider := &mockCredentialsProvider{
		creds: &Credentials{SecretAccessKey: "secret"},
	}
	verifier := NewSignatureV4Verifier(provider)

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=test/20240101/us-east-1/s3/aws4_request, SignedHeaders=host, Signature=abc123")
	err := verifier.VerifyRequest(req, "s3", "us-east-1")
	assert.Equal(t, ErrMissingDateHeader, err)
}

func TestSignatureV4Verifier_BuildCanonicalQueryString(t *testing.T) {
	verifier := &SignatureV4Verifier{}

	tests := []struct {
		name     string
		query    url.Values
		expected string
	}{
		{
			name:     "empty query",
			query:    url.Values{},
			expected: "",
		},
		{
			name:     "single param",
			query:    url.Values{"key": []string{"value"}},
			expected: "key=value",
		},
		{
			name:     "multiple params",
			query:    url.Values{"b": []string{"2"}, "a": []string{"1"}},
			expected: "a=1&b=2",
		},
		{
			name:     "special chars",
			query:    url.Values{"key": []string{"value with space"}},
			expected: "key=value+with+space",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := verifier.buildCanonicalQueryString(tt.query)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSignatureV4Verifier_BuildCanonicalRequest(t *testing.T) {
	provider := &mockCredentialsProvider{
		creds: &Credentials{SecretAccessKey: "secret"},
	}
	verifier := NewSignatureV4Verifier(provider)

	req, _ := http.NewRequest("GET", "http://example.com/test?foo=bar", nil)
	req.Host = "example.com"

	result, err := verifier.buildCanonicalRequest(req, "host", nil)
	require.NoError(t, err)
	assert.Contains(t, result, "GET")
	assert.Contains(t, result, "/test")
	assert.Contains(t, result, "host:example.com")
}

func TestSignatureV4Verifier_BuildStringToSign(t *testing.T) {
	verifier := &SignatureV4Verifier{}

	credential := "AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request"
	canonicalRequest := "GET\n/\n\nhost:example.com\n\nhost\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	scope, stringToSign, err := verifier.buildStringToSign("20130524T000000Z", credential, canonicalRequest)
	require.NoError(t, err)
	assert.Equal(t, "20130524", scope.Date)
	assert.Equal(t, "us-east-1", scope.Region)
	assert.Equal(t, "s3", scope.Service)
	assert.Contains(t, stringToSign, "AWS4-HMAC-SHA256")
}

func TestSignatureV4Verifier_BuildStringToSign_InvalidCredential(t *testing.T) {
	verifier := &SignatureV4Verifier{}

	_, _, err := verifier.buildStringToSign("20130524T000000Z", "invalid", "request")
	assert.Equal(t, ErrInvalidCredentialFormat, err)
}

func TestSignatureV4Verifier_Sha256Hash(t *testing.T) {
	verifier := &SignatureV4Verifier{}

	result := verifier.sha256Hash("")
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", result)
}

func TestSignatureV4Verifier_HmacSHA256(t *testing.T) {
	verifier := &SignatureV4Verifier{}

	key := []byte("secret")
	data := "test"
	result := verifier.hmacSHA256(key, data)
	assert.Len(t, result, 32)
}

func TestSignatureV4Verifier_CalculateSignature(t *testing.T) {
	verifier := &SignatureV4Verifier{}

	signature, err := verifier.calculateSignature("20130524T000000Z", "us-east-1", "s3", "test", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	require.NoError(t, err)
	assert.Len(t, signature, 64)
}

func TestSignatureV4Verifier_VerifyRequest_WithBody(t *testing.T) {
	provider := &mockCredentialsProvider{
		creds: &Credentials{SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"},
	}
	verifier := NewSignatureV4Verifier(provider)

	body := []byte(`{"test": "data"}`)
	req, _ := http.NewRequest("POST", "http://example.com/test", bytes.NewReader(body))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-date, Signature=abc123")
	req.Header.Set("X-Amz-Date", "20130524T000000Z")
	req.Header.Set("Host", "example.com")

	err := verifier.VerifyRequest(req, "s3", "us-east-1")
	assert.Equal(t, ErrSignatureMismatch, err)
}

func TestNewSignatureV4Verifier(t *testing.T) {
	provider := NewStaticCredentialsProvider("key", "secret", "us-east-1", "s3")
	verifier := NewSignatureV4Verifier(provider)
	assert.NotNil(t, verifier)
}
