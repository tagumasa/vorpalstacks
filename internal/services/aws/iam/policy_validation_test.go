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
			name:     "valid with NotAction and NotResource",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow","NotAction":"iam:*","NotResource":"*"}}`,
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
		{
			name:     "Statement missing Action",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Resource":"*"}}`,
			expected: false,
		},
		{
			name:     "Statement missing Resource",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Action":"*"}}`,
			expected: false,
		},
		{
			name:     "fail-open bypass attempt: only Effect",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow"}}`,
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

func TestValidateTrustPolicyDocument(t *testing.T) {
	tests := []struct {
		name     string
		document string
		expected bool
	}{
		{
			name:     "valid trust policy with Principal",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Action":"sts:AssumeRole","Principal":{"Service":"ec2.amazonaws.com"}}}`,
			expected: true,
		},
		{
			name:     "valid trust policy with NotPrincipal",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Deny","Action":"sts:AssumeRole","NotPrincipal":{"AWS":"arn:aws:iam::123456789012:root"}}}`,
			expected: true,
		},
		{
			name:     "trust policy missing Principal",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Action":"sts:AssumeRole"}}`,
			expected: false,
		},
		{
			name:     "trust policy with Resource instead of Principal (should fail)",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Action":"sts:AssumeRole","Resource":"*"}}`,
			expected: false,
		},
		{
			name:     "trust policy missing Action",
			document: `{"Version":"2012-10-17","Statement":{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"}}}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateTrustPolicyDocument(tt.document)
			if got != tt.expected {
				t.Errorf("validateTrustPolicyDocument() = %v, want %v", got, tt.expected)
			}
		})
	}
}
