package policy

import (
	"encoding/json"
	"testing"
)

func TestParseDocument(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		wantErr bool
	}{
		{
			name:    "invalid json",
			jsonStr: `{"Version":`,
			wantErr: true,
		},
		{
			name:    "document without statement",
			jsonStr: `{"Version":"2012-10-17"}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := ParseDocument(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && doc == nil {
				t.Error("expected non-nil document")
			}
		})
	}
}

func TestEffect(t *testing.T) {
	if EffectAllow != "Allow" {
		t.Errorf("expected Allow, got %s", EffectAllow)
	}
	if EffectDeny != "Deny" {
		t.Errorf("expected Deny, got %s", EffectDeny)
	}
}

func TestPrincipalUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		data  string
		check func(*Principal) bool
	}{
		{
			name:  "wildcard principal",
			data:  `"*"`,
			check: func(p *Principal) bool { return p.Everyone },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Principal
			err := json.Unmarshal([]byte(tt.data), &p)
			if err != nil {
				t.Errorf("unmarshal error = %v", err)
			}
			if !tt.check(&p) {
				t.Errorf("check failed for %s", tt.name)
			}
		})
	}
}

func TestPrincipalMatches(t *testing.T) {
	tests := []struct {
		name      string
		principal *Principal
		arn       string
		expected  bool
	}{
		{
			name:      "nil principal",
			principal: nil,
			arn:       "arn:aws:iam::123456789012:user/test",
			expected:  true,
		},
		{
			name:      "everyone",
			principal: &Principal{Everyone: true},
			arn:       "arn:aws:iam::123456789012:user/test",
			expected:  true,
		},
		{
			name:      "exact match",
			principal: &Principal{AWS: StringList{"arn:aws:iam::123456789012:user/test"}},
			arn:       "arn:aws:iam::123456789012:user/test",
			expected:  true,
		},
		{
			name:      "wildcard match",
			principal: &Principal{AWS: StringList{"*"}},
			arn:       "arn:aws:iam::123456789012:user/test",
			expected:  true,
		},
		{
			name:      "no match",
			principal: &Principal{AWS: StringList{"arn:aws:iam::111111111111:user/other"}},
			arn:       "arn:aws:iam::123456789012:user/test",
			expected:  false,
		},
		{
			name:      "prefix match",
			principal: &Principal{AWS: StringList{"arn:aws:iam::123456789012:user/*"}},
			arn:       "arn:aws:iam::123456789012:user/test",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.principal.Matches(tt.arn)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestActionListMatches(t *testing.T) {
	tests := []struct {
		name     string
		actions  ActionList
		action   string
		expected bool
	}{
		{
			name:     "empty action list",
			actions:  ActionList{},
			action:   "s3:GetObject",
			expected: true,
		},
		{
			name:     "wildcard action",
			actions:  ActionList{"*"},
			action:   "s3:GetObject",
			expected: true,
		},
		{
			name:     "exact match",
			actions:  ActionList{"s3:GetObject"},
			action:   "s3:GetObject",
			expected: true,
		},
		{
			name:     "no match",
			actions:  ActionList{"s3:PutObject"},
			action:   "s3:GetObject",
			expected: false,
		},
		{
			name:     "prefix match",
			actions:  ActionList{"s3:*"},
			action:   "s3:GetObject",
			expected: true,
		},
		{
			name:     "prefix no match",
			actions:  ActionList{"ec2:*"},
			action:   "s3:GetObject",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.actions.Matches(tt.action)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestResourceListMatches(t *testing.T) {
	tests := []struct {
		name      string
		resources ResourceList
		resource  string
		expected  bool
	}{
		{
			name:      "empty resource list",
			resources: ResourceList{},
			resource:  "arn:aws:s3:::bucket/object",
			expected:  true,
		},
		{
			name:      "wildcard resource",
			resources: ResourceList{"*"},
			resource:  "arn:aws:s3:::bucket/object",
			expected:  true,
		},
		{
			name:      "exact match",
			resources: ResourceList{"arn:aws:s3:::bucket/object"},
			resource:  "arn:aws:s3:::bucket/object",
			expected:  true,
		},
		{
			name:      "no match",
			resources: ResourceList{"arn:aws:s3:::other-bucket/object"},
			resource:  "arn:aws:s3:::bucket/object",
			expected:  false,
		},
		{
			name:      "prefix match",
			resources: ResourceList{"arn:aws:s3:::bucket/*"},
			resource:  "arn:aws:s3:::bucket/object",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.resources.Matches(tt.resource)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchActionPattern(t *testing.T) {
	tests := []struct {
		pattern  string
		action   string
		expected bool
	}{
		{"*", "s3:GetObject", true},
		{"s3:*", "s3:GetObject", true},
		{"s3:GetObject", "s3:GetObject", true},
		{"s3:Get*", "s3:GetObject", true},
		{"s3:Get*", "s3:PutObject", false},
		{"s3:PutObject", "s3:GetObject", false},
		{"ec2:*", "s3:GetObject", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.action, func(t *testing.T) {
			result := matchActionPattern(tt.pattern, tt.action)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchArnPattern(t *testing.T) {
	tests := []struct {
		pattern  string
		arn      string
		expected bool
	}{
		{"*", "arn:aws:s3:::bucket/object", true},
		{"arn:aws:s3:::*", "arn:aws:s3:::bucket/object", true},
		{"arn:aws:s3:::bucket/*", "arn:aws:s3:::bucket/object", true},
		{"arn:aws:s3:::bucket/object", "arn:aws:s3:::bucket/object", true},
		{"arn:aws:s3:::other-bucket/*", "arn:aws:s3:::bucket/object", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.arn, func(t *testing.T) {
			result := matchArnPattern(tt.pattern, tt.arn)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
