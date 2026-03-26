// Package parsers provides AWS-specific utility functions for vorpalstacks.
package parsers

import (
	"fmt"
	"strings"
)

// ParseConditionKey parses a condition key and returns its components.
// Format: operator:condition_key
//
// Parameters:
//   - key: The condition key to parse
//
// Returns:
//   - operator: The operator part
//   - conditionKey: The condition key part
//   - error: An error if parsing fails
func ParseConditionKey(key string) (operator, conditionKey string, err error) {
	parts := strings.SplitN(key, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid condition key format: %s", key)
	}
	return parts[0], parts[1], nil
}
