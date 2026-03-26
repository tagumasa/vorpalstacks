package crypto

// Package crypto provides cryptographic utilities for vorpalstacks, including
// AWS signature version 4 signature calculation.

import (
	"encoding/hex"
	"errors"
	"time"
)

// ErrInvalidDateFormat is returned when the date format is invalid.
var ErrInvalidDateFormat = errors.New("failed to parse date")

// CalculateSignature calculates the AWS Signature Version 4 signature.
// This is used for signing API requests to AWS services.
//
// Parameters:
//   - amzDate: The date in ISO 8601 format (YYYYMMDDTHHMMSSZ)
//   - region: The AWS region
//   - service: The AWS service name
//   - stringToSign: The string to sign
//   - secretKey: The AWS secret access key
//
// Returns:
//   - string: The hexadecimal-encoded signature
//   - error: An error if signature calculation fails
//
// Example:
//
//	signature, err := CalculateSignature("20240101T120000Z", "us-east-1", "s3", stringToSign, secretKey)
func CalculateSignature(amzDate, region, service, stringToSign, secretKey string) (string, error) {
	date, err := time.Parse("20060102T150405Z", amzDate)
	if err != nil {
		return "", ErrInvalidDateFormat
	}
	dateStr := date.Format("20060102")

	kDate := HMACSHA256String([]byte("AWS4"+secretKey), dateStr)
	kRegion := HMACSHA256String(kDate, region)
	kService := HMACSHA256String(kRegion, service)
	kSigning := HMACSHA256String(kService, ServicePrefix)

	signature := HMACSHA256String(kSigning, stringToSign)

	return hex.EncodeToString(signature), nil
}

// DeriveSigningKey derives the signing key for AWS Signature Version 4.
// The signing key is derived from the secret key, date, region, and service.
//
// Parameters:
//   - secretKey: The AWS secret access key
//   - dateStr: The date in YYYYMMDD format
//   - region: The AWS region
//   - service: The AWS service name
//
// Returns:
//   - []byte: The derived signing key
//
// Example:
//
//	signingKey := DeriveSigningKey(secretKey, "20240101", "us-east-1", "s3")
func DeriveSigningKey(secretKey, dateStr, region, service string) []byte {
	kDate := HMACSHA256String([]byte("AWS4"+secretKey), dateStr)
	kRegion := HMACSHA256String(kDate, region)
	kService := HMACSHA256String(kRegion, service)
	kSigning := HMACSHA256String(kService, ServicePrefix)
	return kSigning
}
