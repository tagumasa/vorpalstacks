// Package sts provides STS (Security Token Service) operations for vorpalstacks.
package sts

import (
	"net/http"

	"vorpalstacks/internal/services/aws/common/errors"
)

// Duration constants for STS sessions.
const (
	// MinDurationSeconds is the minimum allowed session duration in seconds.
	MinDurationSeconds = 900
	// MaxDurationSeconds is the maximum allowed session duration in seconds for standard operations.
	MaxDurationSeconds = 43200
	// MaxDurationSecondsExtended is the maximum allowed session duration in seconds for extended operations.
	MaxDurationSecondsExtended = 129600
	// DefaultDurationSeconds is the default session duration in seconds.
	DefaultDurationSeconds = 3600
)

var (
	// ErrInvalidRoleArn is returned when the specified role ARN is invalid.
	ErrInvalidRoleArn = errors.NewAWSError("InvalidRoleArn", "The specified role ARN is invalid.", http.StatusBadRequest)
	// ErrInvalidPrincipalArn is returned when the specified principal ARN is invalid.
	ErrInvalidPrincipalArn = errors.NewAWSError("InvalidPrincipalArn", "The specified principal ARN is invalid.", http.StatusBadRequest)
	// ErrInvalidSAMLAssertion is returned when the SAML assertion is invalid.
	ErrInvalidSAMLAssertion = errors.NewAWSError("InvalidSAMLAssertion", "The SAML assertion is invalid.", http.StatusBadRequest)
	// ErrInvalidWebIdentityToken is returned when the web identity token is invalid.
	ErrInvalidWebIdentityToken = errors.NewAWSError("InvalidIdentityToken", "The web identity token is invalid.", http.StatusBadRequest)
	// ErrInvalidEncodedMessage is returned when the encoded message is invalid.
	ErrInvalidEncodedMessage = errors.NewAWSError("InvalidEncodedMessage", "The encoded message is invalid.", http.StatusBadRequest)
	// ErrInvalidAccessKeyId is returned when the access key ID is invalid.
	ErrInvalidAccessKeyId = errors.NewAWSError("InvalidAccessKeyId", "The access key ID is invalid.", http.StatusBadRequest)
	// ErrInvalidFederationName is returned when the federation name is invalid.
	ErrInvalidFederationName = errors.NewAWSError("InvalidInput", "The federation name is invalid.", http.StatusBadRequest)
	// ErrInvalidTradeInToken is returned when the trade-in token is invalid.
	ErrInvalidTradeInToken = errors.NewAWSError("InvalidToken", "The trade-in token is invalid.", http.StatusBadRequest)
	// ErrInvalidDuration is returned when the duration is outside the allowed range.
	ErrInvalidDuration = errors.NewAWSError("InvalidDuration", "The duration seconds must be between 900 and 43200 seconds.", http.StatusBadRequest)
	// ErrInvalidDurationExtended is returned when the duration is outside the extended allowed range.
	ErrInvalidDurationExtended = errors.NewAWSError("InvalidDuration", "The duration seconds must be between 900 and 129600 seconds.", http.StatusBadRequest)
	// ErrNoSuchRole is returned when the specified role cannot be found.
	ErrNoSuchRole = errors.NewAWSError("NoSuchEntity", "The role with the specified ARN cannot be found.", http.StatusNotFound)
	// ErrMalformedPolicyDocument is returned when the policy document is invalid.
	ErrMalformedPolicyDocument = errors.NewAWSError("MalformedPolicyDocument", "The policy document is invalid.", http.StatusBadRequest)
	// ErrInvalidParameter is returned when a parameter value is invalid.
	ErrInvalidParameter = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'roleSessionName' failed to satisfy constraint: Member must satisfy regular expression pattern: [\\w+=,.@-]*", http.StatusBadRequest)
)
