// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	keyIDRegex     = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	aliasNameRegex = regexp.MustCompile(`^alias/([a-zA-Z0-9/_-]+|aws/[a-zA-Z0-9/_-]+)$`)
	arnRegex       = regexp.MustCompile(`^arn:aws:kms:[a-z0-9-]+:[0-9]{12}:key/[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	aliasArnRegex  = regexp.MustCompile(`^arn:aws:kms:[a-z0-9-]+:[0-9]{12}:alias/([a-zA-Z0-9/_-]+|aws/[a-zA-Z0-9/_-]+)$`)
)

// ValidateKeyID validates a KMS key ID.
//
// Valid formats:
//   - Key ID (UUID): 8-4-4-4-12 hexadecimal format
//   - Alias: alias/alias-name
//   - Key ARN: arn:aws:kms:region:account:key/key-id
//   - Alias ARN: arn:aws:kms:region:account:alias/alias-name
//
// Parameters:
//   - keyID: The key ID to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidateKeyID(keyID string) error {
	if keyID == "" {
		return fmt.Errorf("key ID cannot be empty")
	}

	switch {
	case strings.HasPrefix(keyID, "arn:aws:kms:"):
		if arnRegex.MatchString(keyID) || aliasArnRegex.MatchString(keyID) {
			return nil
		}
		return fmt.Errorf("invalid key ARN format: %s", keyID)
	case strings.HasPrefix(keyID, "alias/"):
		if aliasNameRegex.MatchString(keyID) {
			return nil
		}
		return fmt.Errorf("invalid alias name format: %s", keyID)
	case keyIDRegex.MatchString(keyID):
		return nil
	default:
		return fmt.Errorf("invalid key ID format: %s", keyID)
	}
}
