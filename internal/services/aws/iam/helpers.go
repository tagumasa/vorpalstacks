// Package iam provides IAM service operations for vorpalstacks.
package iam

import "encoding/json"

const (
	// MaxAccessKeysPerUser is the maximum number of access keys a user can have.
	MaxAccessKeysPerUser = 2
	// MaxPolicyVersions is the maximum number of policy versions allowed.
	MaxPolicyVersions = 5
)

// validatePolicyDocument checks if a policy document is valid JSON.
func validatePolicyDocument(document string) bool {
	if document == "" {
		return false
	}
	var js interface{}
	return json.Unmarshal([]byte(document), &js) == nil
}
