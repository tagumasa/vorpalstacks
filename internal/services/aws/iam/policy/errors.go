// Package policy provides IAM policy evaluation functionality for vorpalstacks.
package policy

import (
	"fmt"
)

// PolicyError represents an error that occurs during IAM policy evaluation.
type PolicyError struct {
	Code    string
	Message string
}

// Error returns the string representation of the policy error.
func (e *PolicyError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewPolicyError creates a new PolicyError with the specified code and message.
func NewPolicyError(code, message string) *PolicyError {
	return &PolicyError{Code: code, Message: message}
}

var (
	// ErrInvalidPolicyDocument is returned when the policy document is not valid.
	ErrInvalidPolicyDocument = NewPolicyError("InvalidPolicyDocument", "The policy document is invalid")

	// ErrMalformedPolicy is returned when the policy contains malformed JSON.
	ErrMalformedPolicy = NewPolicyError("MalformedPolicy", "The policy contains malformed JSON")

	// ErrInvalidEffect is returned when the policy effect is not Allow or Deny.
	ErrInvalidEffect = NewPolicyError("InvalidEffect", "Policy effect must be Allow or Deny")

	// ErrMissingAction is returned when the policy statement does not specify
	// Action or NotAction.
	ErrMissingAction = NewPolicyError("MissingAction", "Policy statement must specify Action or NotAction")

	// ErrMissingResource is returned when the policy statement does not specify
	// Resource or NotResource.
	ErrMissingResource = NewPolicyError("MissingResource", "Policy statement must specify Resource or NotResource")

	// ErrAccessDenied is returned when the user is not authorized to perform
	// the requested action.
	ErrAccessDenied = NewPolicyError("AccessDenied", "User is not authorized to perform this action")
)

// IsAccessDenied checks if the error indicates an access denied error.
func IsAccessDenied(err error) bool {
	if policyErr, ok := err.(*PolicyError); ok {
		return policyErr.Code == "AccessDenied"
	}
	return false
}
