package iam

import (
	"testing"
)

func TestValidatePolicyDocument(t *testing.T) {
	tests := []struct {
		name     string
		document string
		expected bool
	}{
		{
			name:     "valid single statement",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Action":"s3:Get*","Resource":"*"}}`,
			expected: true,
		},
		{
			name:     "valid array of statements",
			document: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"},{"Effect":"Deny","Action":"iam:*","Resource":"*"}]}`,
			expected: true,
		},
		{
			name:     "empty string",
			document: "",
			expected: false,
		},
		{
			name:     "invalid JSON",
			document: `{not json`,
			expected: false,
		},
		{
			name:     "missing Statement",
			document: `{"Version":"2012-10-17"}`,
			expected: false,
		},
		{
			name:     "empty Statement array",
			document: `{"Version":"2012-10-17","Statement":[]}`,
			expected: false,
		},
		{
			name:     "Statement missing Effect",
			document: `{"Version":"2012-10-17","Statement":[{"Action":"*","Resource":"*"}]}`,
			expected: false,
		},
		{
			name:     "Statement invalid Effect",
			document: `{"Version":"2012-10-17","Statement":[{"Effect":"Maybe","Action":"*","Resource":"*"}]}`,
			expected: false,
		},
		{
			name:     "one valid one invalid in array",
			document: `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"},{"Effect":"Maybe","Action":"iam:*","Resource":"*"}]}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validatePolicyDocument(tt.document)
			if got != tt.expected {
				t.Errorf("validatePolicyDocument() = %v, want %v", got, tt.expected)
			}
		})
	}
}
