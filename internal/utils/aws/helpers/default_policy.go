// Package helpers provides AWS-specific utility functions for vorpalstacks.
package helpers

import (
	"fmt"
	"regexp"

	"vorpalstacks/internal/config"
)

var awsAccountIDRegex = regexp.MustCompile(`^\d{12}$`)

// DefaultPolicy returns a default KMS key policy JSON string.
//
// The policy grants the root user of the specified AWS account permissions
// to use the KMS key for encrypt, decrypt, re-encrypt, generate data key,
// describe key, create grant, list grants, and revoke grant operations.
//
// Parameters:
//   - accountID: The AWS account ID (12 digits). If invalid, defaults to 000000000000
//
// Returns:
//   - string: The default policy JSON string
func DefaultPolicy(accountID string) string {
	if !awsAccountIDRegex.MatchString(accountID) {
		accountID = config.AWSAccountID()
		if !awsAccountIDRegex.MatchString(accountID) {
			accountID = "000000000000"
		}
	}
	return fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Id": "key-default-1",
		"Statement": [
			{
				"Sid": "Enable IAM User Permissions",
				"Effect": "Allow",
				"Principal": {
					"AWS": "arn:aws:iam::%s:root"
				},
				"Action": [
					"kms:Encrypt",
					"kms:Decrypt",
					"kms:ReEncrypt*",
					"kms:GenerateDataKey*",
					"kms:DescribeKey",
					"kms:CreateGrant",
					"kms:ListGrants",
					"kms:RevokeGrant"
				],
				"Resource": "*"
			}
		]
	}`, accountID)
}
