package policy

import (
	"testing"
)

func TestConditionEvaluator_Evaluate(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name       string
		conditions ConditionMap
		ctx        *EvaluationContext
		expected   bool
	}{
		{
			name:       "empty conditions",
			conditions: nil,
			ctx:        &EvaluationContext{},
			expected:   true,
		},
		{
			name:       "empty condition map",
			conditions: ConditionMap{},
			ctx:        &EvaluationContext{},
			expected:   true,
		},
		{
			name: "string equals - missing context",
			conditions: ConditionMap{
				"StringEquals": {
					"action": []string{"s3:GetObject"},
				},
			},
			ctx:      &EvaluationContext{Action: "s3:GetObject"},
			expected: false,
		},
		{
			name: "string equals false",
			conditions: ConditionMap{
				"StringEquals": {
					"action": []string{"s3:PutObject"},
				},
			},
			ctx:      &EvaluationContext{Action: "s3:GetObject"},
			expected: false,
		},
		{
			name: "string not equals - different value returns false",
			conditions: ConditionMap{
				"StringNotEquals": {
					"action": []string{"s3:PutObject"},
				},
			},
			ctx:      &EvaluationContext{Action: "s3:GetObject"},
			expected: false,
		},
		{
			name: "missing context value",
			conditions: ConditionMap{
				"StringEquals": {
					"missing": []string{"value"},
				},
			},
			ctx:      &EvaluationContext{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.Evaluate(tt.conditions, tt.ctx)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchLikePattern(t *testing.T) {
	tests := []struct {
		pattern  string
		value    string
		expected bool
	}{
		{"*", "anything", true},
		{"test*", "testvalue", true},
		{"test*", "test", true},
		{"test*", "othervalue", false},
		{"*test*", "prefixtestsuffix", true},
		{"prefix*suffix", "prefixvaluesuffix", true},
		{"prefix*suffix", "prefixsuffix", true},
		{"prefix*suffix", "prefixtesuffix", true},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.value, func(t *testing.T) {
			result := matchLikePattern(tt.pattern, tt.value)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestWildcardMatch(t *testing.T) {
	tests := []struct {
		pattern  string
		value    string
		expected bool
	}{
		{"*", "anything", true},
		{"test*", "testvalue", true},
		{"test*", "test", true},
		{"test*", "other", false},
		{"test", "test", true},
		{"test", "other", false},
		{"prefix*middle*suffix", "prefix123middle456suffix", true},
		{"prefix*middle*suffix", "prefix123othersuffix", false},
		{"*suffix", "prefixsuffix", true},
		{"prefix*", "prefix", true},
		{"prefix*", "prefixvalue", true},
		{"a*b*c", "axbxc", true},
		{"a*b*c", "axbyc", true},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.value, func(t *testing.T) {
			result := wildcardMatch(tt.pattern, tt.value)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchArnLike(t *testing.T) {
	tests := []struct {
		pattern  string
		arn      string
		expected bool
	}{
		{"arn:aws:s3:::bucket/*", "arn:aws:s3:::bucket/object", true},
		{"arn:aws:s3:::bucket/*", "arn:aws:s3:::other/object", false},
		{"arn:aws:s3:::bucket", "arn:aws:s3:::bucket", true},
		{"arn:aws:s3:::bucket", "arn:aws:s3:::other", false},
		{"arn:aws:ec2:*:*:instance/*", "arn:aws:ec2:us-east-1:123456789012:instance/i-12345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.arn, func(t *testing.T) {
			result := matchArnLike(tt.pattern, tt.arn)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchIPAddress(t *testing.T) {
	tests := []struct {
		value    string
		cidr     string
		expected bool
	}{
		{"192.168.1.1", "192.168.1.0/24", true},
		{"192.168.1.1", "192.168.2.0/24", false},
		{"192.168.1.1", "192.168.1.1", true},
		{"invalid", "192.168.1.0/24", false},
		{"192.168.1.1", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.value+"_"+tt.cidr, func(t *testing.T) {
			result := matchIPAddress(tt.value, tt.cidr)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCompareNumeric(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected int
	}{
		{"10", "20", -1},
		{"20", "10", 1},
		{"10", "10", 0},
		{"10.5", "10.6", -1},
		{"-5", "5", -1},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			result := compareNumeric(tt.a, tt.b)
			if (result < 0) != (tt.expected < 0) || (result > 0) != (tt.expected > 0) || (result == 0) != (tt.expected == 0) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"10", 10},
		{"-10", -10},
		{"10.5", 10.5},
		{"-10.5", -10.5},
		{"123", 123},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseFloat(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCompareDates(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected int
	}{
		{"2023-01-01T00:00:00Z", "2023-01-02T00:00:00Z", -1},
		{"2023-01-02T00:00:00Z", "2023-01-01T00:00:00Z", 1},
		{"2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z", 0},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			result := compareDates(tt.a, tt.b)
			if (result < 0) != (tt.expected < 0) || (result > 0) != (tt.expected > 0) || (result == 0) != (tt.expected == 0) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFilepathMatch(t *testing.T) {
	tests := []struct {
		pattern  string
		name     string
		expected bool
	}{
		{"*.txt", "file.txt", true},
		{"*.txt", "file.md", false},
		{"test*", "testfile", true},
		{"test*", "otherfile", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.name, func(t *testing.T) {
			result := filepathMatch(tt.pattern, tt.name)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
