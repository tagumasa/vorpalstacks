package crypto

// Package crypto provides cryptographic utilities for vorpalstacks, including
// AWS signature version 4 verification functions.

import (
	"crypto/subtle"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Constants for AWS Signature Version 4.
const (
	// Algorithm is the AWS Signature Version 4 algorithm identifier.
	Algorithm = "AWS4-HMAC-SHA256"
	// ServicePrefix is the AWS Signature Version 4 service prefix.
	ServicePrefix = "aws4_request"
)

// Error variables for signature verification.
var (
	// ErrMissingAuthorizationHeader is returned when the Authorization header is missing.
	ErrMissingAuthorizationHeader = errors.New("missing Authorization header")
	// ErrInvalidAlgorithm is returned when the algorithm in the header doesn't match.
	ErrInvalidAlgorithm = errors.New("invalid Authorization header: algorithm mismatch")
	// ErrMissingRequiredFields is returned when required fields are missing from the header.
	ErrMissingRequiredFields = errors.New("missing required fields in Authorization header")
	// ErrMissingAmzDate is returned when the X-Amz-Date header is missing.
	ErrMissingAmzDate = errors.New("missing X-Amz-Date header")
	// ErrSignatureMismatch is returned when the signature doesn't match.
	ErrSignatureMismatch = errors.New("signature mismatch")
)

// SignatureV4Verifier verifies AWS Signature Version 4 requests.
type SignatureV4Verifier struct {
	credentialsProvider CredentialsProvider
}

// NewSignatureV4Verifier creates a new Signature V4 verifier.
//
// Parameters:
//   - provider: The credentials provider
//
// Returns:
//   - *SignatureV4Verifier: A new verifier instance
//
// Example:
//
//	verifier := NewSignatureV4Verifier(credentialsProvider)
func NewSignatureV4Verifier(provider CredentialsProvider) *SignatureV4Verifier {
	return &SignatureV4Verifier{
		credentialsProvider: provider,
	}
}

// VerifyRequest verifies an incoming HTTP request using AWS Signature Version 4.
// It validates the Authorization header against the computed signature.
//
// Parameters:
//   - r: The HTTP request to verify
//   - service: The AWS service name (e.g., "s3", "lambda")
//   - region: The AWS region
//
// Returns:
//   - error: An error if verification fails, nil if successful
//
// Example:
//
//	err := verifier.VerifyRequest(r, "s3", "us-east-1")
func (v *SignatureV4Verifier) VerifyRequest(r *http.Request, service, region string) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ErrMissingAuthorizationHeader
	}

	if !strings.HasPrefix(authHeader, Algorithm+" ") {
		return ErrInvalidAlgorithm
	}

	authHeader = strings.TrimPrefix(authHeader, Algorithm+" ")

	parsed, err := parseAuthorizationHeader(authHeader)
	if err != nil {
		return err
	}

	if parsed.Credential == "" || parsed.SignedHeaders == "" || parsed.Signature == "" {
		return ErrMissingRequiredFields
	}

	amzDate := r.Header.Get("X-Amz-Date")
	if amzDate == "" {
		return ErrMissingAmzDate
	}

	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}

	canonicalRequest := BuildCanonicalRequest(r, parsed.SignedHeaders, bodyBytes)

	stringToSign, err := BuildStringToSign(amzDate, parsed.Credential, canonicalRequest)
	if err != nil {
		return err
	}

	credentials, err := v.credentialsProvider.GetCredentials()
	if err != nil {
		return err
	}

	calculatedSignature, err := CalculateSignature(amzDate, credentials.Region, service, stringToSign, credentials.SecretAccessKey)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare([]byte(calculatedSignature), []byte(parsed.Signature)) != 1 {
		return ErrSignatureMismatch
	}

	return nil
}

func parseAuthorizationHeader(authHeader string) (*AuthorizationHeader, error) {
	result := &AuthorizationHeader{
		Algorithm: Algorithm,
	}

	parts := strings.Split(authHeader, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "Credential=") {
			result.Credential = strings.TrimPrefix(part, "Credential=")
		} else if strings.HasPrefix(part, "SignedHeaders=") {
			result.SignedHeaders = strings.TrimPrefix(part, "SignedHeaders=")
		} else if strings.HasPrefix(part, "Signature=") {
			result.Signature = strings.TrimPrefix(part, "Signature=")
		}
	}

	return result, nil
}
