// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
	"unicode"
)

const defaultMaxPasswordLength = 128

// PasswordPolicy represents a password policy configuration for IAM users.
type PasswordPolicy struct {
	MinimumLength    int
	MaximumLength    int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSymbols   bool
}

func containsUppercase(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

func containsLowercase(s string) bool {
	for _, c := range s {
		if unicode.IsLower(c) {
			return true
		}
	}
	return false
}

func containsNumber(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

func containsSymbol(s string) bool {
	for _, c := range s {
		if unicode.IsPunct(c) || unicode.IsSymbol(c) {
			return true
		}
	}
	return false
}

// ValidatePasswordPolicy validates a password against an IAM password policy.
//
// Parameters:
//   - password: The password to validate
//   - policy: The password policy to validate against
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidatePasswordPolicy(password string, policy *PasswordPolicy) error {
	if policy == nil {
		return nil
	}

	maxLen := policy.MaximumLength
	if maxLen <= 0 {
		maxLen = defaultMaxPasswordLength
	}

	if len(password) < policy.MinimumLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinimumLength)
	}

	if len(password) > maxLen {
		return fmt.Errorf("password must be at most %d characters long", maxLen)
	}

	if policy.RequireUppercase && !containsUppercase(password) {
		return fmt.Errorf("password must contain uppercase letters")
	}

	if policy.RequireLowercase && !containsLowercase(password) {
		return fmt.Errorf("password must contain lowercase letters")
	}

	if policy.RequireNumbers && !containsNumber(password) {
		return fmt.Errorf("password must contain numbers")
	}

	if policy.RequireSymbols && !containsSymbol(password) {
		return fmt.Errorf("password must contain special characters")
	}

	return nil
}
