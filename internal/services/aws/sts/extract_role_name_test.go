package sts

import (
	"testing"

	arnutil "vorpalstacks/internal/utils/aws/arn"

	"github.com/stretchr/testify/assert"
)

func TestExtractRoleNameFromArn(t *testing.T) {
	tests := []struct {
		name     string
		arn      string
		expected string
	}{
		{
			name:     "simple role ARN",
			arn:      "arn:aws:iam::123456789012:role/MyRole",
			expected: "MyRole",
		},
		{
			name:     "role ARN with path",
			arn:      "arn:aws:iam::123456789012:role/application/backend/MyRole",
			expected: "MyRole",
		},
		{
			name:     "role ARN with multiple path segments",
			arn:      "arn:aws:iam::123456789012:role/path/to/role/MyRole",
			expected: "MyRole",
		},
		{
			name:     "role ARN with leading slash",
			arn:      "arn:aws:iam::123456789012:role/",
			expected: "",
		},
		{
			name:     "arn without slash separator",
			arn:      "arn:aws:iam::123456789012:role",
			expected: "role",
		},
		{
			name:     "empty ARN",
			arn:      "",
			expected: "",
		},
		{
			name:     "role ARN with special characters",
			arn:      "arn:aws:iam::123456789012:role/My-Role_123",
			expected: "My-Role_123",
		},
		{
			name:     "role ARN with dots in path",
			arn:      "arn:aws:iam::123456789012:role/com.example/MyRole",
			expected: "MyRole",
		},
		{
			name:     "IAM ARN without role prefix",
			arn:      "arn:aws:iam::123456789012:user/myuser",
			expected: "myuser",
		},
		{
			name:     "user ARN with path",
			arn:      "arn:aws:iam::123456789012:user/admins/sysadmin",
			expected: "sysadmin",
		},
		{
			name:     "assumed role ARN",
			arn:      "arn:aws:sts::123456789012:assumed-role/MyRole/session",
			expected: "session",
		},
		{
			name:     "role ARN no partition",
			arn:      "arn:aws:iam:::role/MyRole",
			expected: "MyRole",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := arnutil.ExtractRoleNameFromARN(tt.arn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractRoleNameFromArn_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		arn      string
		expected string
	}{
		{
			name:     "only slashes and colons",
			arn:      "arn:aws:iam::123456789012:role/",
			expected: "",
		},
		{
			name:     "role at start after colon",
			arn:      "arn:aws:iam::123456789012:role",
			expected: "role",
		},
		{
			name:     "very long role name",
			arn:      "arn:aws:iam::123456789012:role/VeryLongRoleNameThatExceedsNormalLengthForTestingPurposes",
			expected: "VeryLongRoleNameThatExceedsNormalLengthForTestingPurposes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := arnutil.ExtractRoleNameFromARN(tt.arn)
			assert.Equal(t, tt.expected, result)
		})
	}
}
