package crypto

// Package crypto provides cryptographic utilities for vorpalstacks, including
// AWS signature version 4 string-to-sign generation.

import (
	"errors"
	"strings"
)

// ErrInvalidCredentialFormat is returned when the credential format is invalid.
var ErrInvalidCredentialFormat = errors.New("invalid credential format")

// BuildStringToSign builds the string to sign for AWS Signature Version 4.
func BuildStringToSign(amzDate, credential, canonicalRequest string) (string, error) {
	parts := strings.Split(credential, "/")
	if len(parts) < 5 {
		return "", ErrInvalidCredentialFormat
	}

	date := parts[1]
	region := parts[2]
	service := parts[3]

	canonicalRequestHash := SHA256HashString(canonicalRequest)

	return Algorithm + "\n" +
		amzDate + "\n" +
		date + "/" + region + "/" + service + "/" + ServicePrefix + "\n" +
		canonicalRequestHash, nil
}
