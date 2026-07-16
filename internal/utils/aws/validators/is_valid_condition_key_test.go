package validators

import "testing"

func TestIsValidConditionKey(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		valid bool
	}{
		{"Empty string", "", false},
		{"AWS SourceIp", "aws:SourceIp", true},
		{"AWS SecureTransport", "aws:SecureTransport", true},
		{"AWS PrincipalArn", "aws:PrincipalArn", true},
		{"AWS PrincipalTag/key", "aws:PrincipalTag/department", true},
		{"AWS Tag/key", "aws:Tag/environment", true},
		{"AWS with wildcard", "aws:PrincipalTag/*", true},
		{"AWS unknown", "aws:UnknownKey", false},
		{"env prefix", "env:MY_VAR", true},
		{"env prefix short", "env:ab", true},
		{"http prefix short", "http:ab", true},
		{"s3 prefix short", "s3:ab", true},
		{"ec2 prefix short", "ec2:ab", true},
		{"aws:sms", "aws:sms", false},
		{"aws:sns", "aws:sns", false},
		{"aws:cloudwatch", "aws:cloudwatch", false},
		{"Random key", "random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidConditionKey(tt.key); got != tt.valid {
				t.Errorf("IsValidConditionKey(%q) = %v, want %v", tt.key, got, tt.valid)
			}
		})
	}
}
