// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"fmt"
)

// ValidateProtocol validates an SNS subscription protocol.
//
// Parameters:
//   - protocol: The protocol to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
//
// Valid protocols: http, https, email, email-json, sqs, sms, lambda, application, firehose
func ValidateProtocol(protocol string) error {
	validProtocols := map[string]bool{
		"http":        true,
		"https":       true,
		"email":       true,
		"email-json":  true,
		"sqs":         true,
		"sms":         true,
		"lambda":      true,
		"application": true,
		"firehose":    true,
	}
	if !validProtocols[protocol] {
		return fmt.Errorf("invalid protocol: %s", protocol)
	}
	return nil
}
