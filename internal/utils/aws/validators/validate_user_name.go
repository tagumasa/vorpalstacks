// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
)

func validateIAMName(name string, maxLength int, nameType string) error {
	if name == "" {
		return fmt.Errorf("%s cannot be empty", nameType)
	}
	if len(name) > maxLength {
		return fmt.Errorf("%s cannot exceed %d characters", nameType, maxLength)
	}
	for _, r := range name {
		if !isAllowedIAMNameChar(r) {
			return fmt.Errorf("%s contains invalid character: %c", nameType, r)
		}
	}
	return nil
}

func isAllowedIAMNameChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '+' || r == '=' || r == '.' || r == '@' || r == '-' || r == '_'
}

// ValidateUserName validates an IAM user name.
//
// Parameters:
//   - userName: The user name to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidateUserName(userName string) error {
	return validateIAMName(userName, 64, "user name")
}
