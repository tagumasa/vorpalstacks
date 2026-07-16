// Package arn provides utilities for parsing and constructing Amazon Resource Names (ARNs).
package arn

import (
	"fmt"
	"strings"
)

// ParseARN parses an Amazon Resource Name string and returns a ParsedARN structure.
// Returns an error if the ARN format is invalid.
func ParseARN(arn string) (*ParsedARN, error) {
	if !strings.HasPrefix(arn, "arn:") {
		return nil, fmt.Errorf("invalid ARN format: %s", arn)
	}
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid ARN format: %s", arn)
	}
	parsed := &ParsedARN{
		Partition: parts[1],
		Service:   parts[2],
		Region:    parts[3],
		AccountID: parts[4],
		Resource:  parts[5],
	}
	if strings.HasPrefix(parsed.Resource, "role/") {
		rolePath := strings.TrimPrefix(parsed.Resource, "role/")
		if idx := strings.LastIndex(rolePath, "/"); idx >= 0 {
			parsed.RoleName = rolePath[idx+1:]
		} else {
			parsed.RoleName = rolePath
		}
	}
	return parsed, nil
}

// SplitARN splits an ARN into its constituent parts.
// Returns partition, service, region, account ID, and resource.
func SplitARN(arn string) (partition, service, region, accountID, resource string) {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) < 6 || parts[0] != "arn" {
		return "", "", "", "", ""
	}
	return parts[1], parts[2], parts[3], parts[4], parts[5]
}

// GetServiceFromARN extracts the service name from an ARN.
func GetServiceFromARN(arn string) string {
	_, service, _, _, _ := SplitARN(arn)
	return service
}

// IsValidRoleARN returns true if the ARN represents a valid IAM role.
func IsValidRoleARN(arn string) bool {
	parsed, err := ParseARN(arn)
	if err != nil {
		return false
	}
	return parsed.Service == "iam" && parsed.RoleName != ""
}
