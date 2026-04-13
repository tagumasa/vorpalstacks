// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"encoding/json"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/aws/types"
)

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

// parseTagKeys parses tag keys from request parameters.
func parseTagKeys(params map[string]interface{}) map[string]bool {
	result := make(map[string]bool)
	keys := tags.ParseTagKeysWithQueryFallback(params, "TagKeys")
	for _, k := range keys {
		result[k] = true
	}
	return result
}

// tagsToResponse converts tags to response format.
func tagsToResponse(t []types.Tag) []map[string]interface{} {
	return tags.ToResponse(t)
}

// applyTags merges new tags into existing tags.
func applyTags(existing []types.Tag, newTags []types.Tag) []types.Tag {
	return tags.Apply(existing, newTags)
}

// removeTags removes tags with specified keys from a tag slice.
func removeTags(t []types.Tag, keysToRemove map[string]bool) []types.Tag {
	return tags.Remove(t, keysToRemove)
}

// getMaxItems extracts the maximum items parameter from request parameters.
func getMaxItems(params map[string]interface{}) int {
	return pagination.GetMaxItems(params, pagination.DefaultMaxItems)
}

// GetStringParam extracts a string parameter from request parameters.
func GetStringParam(params map[string]interface{}, key string) string {
	return request.GetStringParam(params, key)
}

// GetIntParam extracts an integer parameter from request parameters.
func GetIntParam(params map[string]interface{}, key string) int {
	return request.GetIntParam(params, key)
}

// GetBoolParam extracts a boolean parameter from request parameters.
func GetBoolParam(params map[string]interface{}, key string) bool {
	return request.GetBoolParam(params, key)
}
