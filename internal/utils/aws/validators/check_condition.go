// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
)

// CheckCondition evaluates a condition and returns an error if it fails.
// This is a simplified version for utility purposes.
//
// Parameters:
//   - conditionResult: The result of condition evaluation
//   - errorMessage: The error message if condition fails
//
// Returns:
//   - error: An error if condition fails, nil otherwise
func CheckCondition(conditionResult bool, errorMessage string) error {
	if !conditionResult {
		return fmt.Errorf("the conditional request failed: %s", errorMessage)
	}
	return nil
}
