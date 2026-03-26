// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import "strings"

var validAWSConditionKeys = map[string]bool{
	"CurrentTime":              true,
	"EpochTime":                true,
	"SourceIp":                 true,
	"UserAgent":                true,
	"Referer":                  true,
	"PrincipalArn":             true,
	"PrincipalId":              true,
	"ResourceArn":              true,
	"Action":                   true,
	"Via":                      true,
	"RequestedRegion":          true,
	"PrincipalAccount":         true,
	"PrincipalARN":             true,
	"PrincipalTag":             true,
	"PrincipalIAMUserArn":      true,
	"PrincipalIAMRoleArn":      true,
	"PrincipalServiceName":     true,
	"SecureTransport":          true,
	"SourceVpc":                true,
	"SourceVpce":               true,
	"RequestedProvider":        true,
	"RequestedResourceAccount": true,
	"ResourceAccount":          true,
	"ResourceOwner":            true,
	"ResourceVpc":              true,
	"TagKey":                   true,
	"RequestTime":              true,
	"PrincipalOrgId":           true,
	"PrincipalOrgPaths":        true,
}

// IsValidConditionKey checks if a condition key is valid for AWS IAM policies.
// It validates both AWS-wide condition keys and service-specific keys.
//
// Parameters:
//   - key: The condition key to validate
//
// Returns:
//   - bool: True if the condition key is valid
//
// Valid keys include:
//   - AWS-wide keys (aws:CurrentTime, aws:SourceIp, aws:PrincipalArn, etc.)
//   - Environment variables (env:VAR_NAME)
//   - HTTP headers (http:HeaderName)
//   - EC2 tags (ec2:tag/TagName)
//   - S3 keys (s3:prefix, s3:Delimiter)
//   - Principal tags (PrincipalTag/tagname)
//
// Example:
//
//	IsValidConditionKey("aws:SourceIp")       // true
//	IsValidConditionKey("env:MY_VAR")          // true
//	IsValidConditionKey("ec2:tag/Department") // true
func IsValidConditionKey(key string) bool {
	if key == "" {
		return false
	}

	if strings.HasPrefix(key, "aws:") {
		awsKey := strings.TrimPrefix(key, "aws:")
		if _, ok := validAWSConditionKeys[awsKey]; ok {
			return true
		}
		if strings.Contains(awsKey, "*") {
			return true
		}
		if strings.HasPrefix(awsKey, "PrincipalTag/") {
			return true
		}
		if strings.HasPrefix(awsKey, "Tag/") {
			return true
		}
		if strings.HasPrefix(awsKey, "sms:") || strings.HasPrefix(awsKey, "sns:") || strings.HasPrefix(awsKey, "cloudwatch:") {
			return true
		}
		return false
	}

	if strings.HasPrefix(key, "env:") {
		return len(key) > 4
	}

	if strings.HasPrefix(key, "http:") {
		return len(key) > 5
	}

	if strings.HasPrefix(key, "s3:") {
		return len(key) > 3
	}

	if strings.HasPrefix(key, "ec2:") {
		return len(key) > 4
	}

	return false
}
