package crypto

import (
	"errors"
	"strings"
)

// Constants for AWS Signature Version 4.
const (
	// Algorithm is the AWS Signature Version 4 algorithm identifier.
	Algorithm = "AWS4-HMAC-SHA256"
	// ServicePrefix is the AWS Signature Version 4 service prefix.
	ServicePrefix = "aws4_request"
)

// AuthorizationHeader represents the components of an AWS authorization header.
type AuthorizationHeader struct {
	Algorithm     string
	Credential    string
	SignedHeaders string
	Signature     string
}

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

// DeriveSigningKey derives the signing key for AWS Signature Version 4.
// The signing key is derived from the secret key, date, region, and service.
func DeriveSigningKey(secretKey, dateStr, region, service string) []byte {
	kDate := HMACSHA256String([]byte("AWS4"+secretKey), dateStr)
	kRegion := HMACSHA256String(kDate, region)
	kService := HMACSHA256String(kRegion, service)
	kSigning := HMACSHA256String(kService, ServicePrefix)
	return kSigning
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
