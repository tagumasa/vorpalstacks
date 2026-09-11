package policy

import (
	"testing"
)

// TestConditionEvaluator_ResolvedAWSVariables exercises the regression:
// conditions with aws:* keys must be evaluated against the resolved
// value, not re-resolved as if the value were a fresh key.
func TestConditionEvaluator_ResolvedAWSVariables(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name       string
		conditions ConditionMap
		ctx        *EvaluationContext
		expected   bool
	}{
		{
			name: "aws:SourceIp equals matches",
			conditions: ConditionMap{
				"StringEquals": {
					"aws:SourceIp": []string{"192.168.1.1"},
				},
			},
			ctx:      &EvaluationContext{SourceIP: "192.168.1.1"},
			expected: true,
		},
		{
			name: "aws:SourceIp equals mismatch",
			conditions: ConditionMap{
				"StringEquals": {
					"aws:SourceIp": []string{"10.0.0.1"},
				},
			},
			ctx:      &EvaluationContext{SourceIP: "192.168.1.1"},
			expected: false,
		},
		{
			name: "aws:username equals matches",
			conditions: ConditionMap{
				"StringEquals": {
					"aws:username": []string{"alice"},
				},
			},
			ctx:      &EvaluationContext{UserName: "alice"},
			expected: true,
		},
		{
			name: "aws:MultiFactorAuthPresent true matches when present",
			conditions: ConditionMap{
				"Bool": {
					"aws:MultiFactorAuthPresent": []string{"true"},
				},
			},
			ctx:      &EvaluationContext{MultiFactorAuthPresent: true},
			expected: true,
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
			name: "string equals - action context matches",
			conditions: ConditionMap{
				"StringEquals": {
					"action": []string{"s3:GetObject"},
				},
			},
			ctx:      &EvaluationContext{Action: "s3:GetObject"},
			expected: true,
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
			name: "string not equals - different value satisfies condition",
			conditions: ConditionMap{
				"StringNotEquals": {
					"action": []string{"s3:PutObject"},
				},
			},
			ctx:      &EvaluationContext{Action: "s3:GetObject"},
			expected: true,
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

func TestConditionEvaluator_Null(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name       string
		conditions ConditionMap
		ctx        *EvaluationContext
		expected   bool
	}{
		{
			name: "null true - key absent",
			conditions: ConditionMap{
				"Null": {
					"aws:referer": []string{"true"},
				},
			},
			ctx:      &EvaluationContext{},
			expected: true,
		},
		{
			name: "null true - key present",
			conditions: ConditionMap{
				"Null": {
					"aws:referer": []string{"true"},
				},
			},
			ctx:      &EvaluationContext{Referer: "https://example.com"},
			expected: false,
		},
		{
			name: "null false - key present",
			conditions: ConditionMap{
				"Null": {
					"aws:referer": []string{"false"},
				},
			},
			ctx:      &EvaluationContext{Referer: "https://example.com"},
			expected: true,
		},
		{
			name: "null false - key absent",
			conditions: ConditionMap{
				"Null": {
					"aws:referer": []string{"false"},
				},
			},
			ctx:      &EvaluationContext{},
			expected: false,
		},
		{
			name: "null true case insensitive",
			conditions: ConditionMap{
				"Null": {
					"aws:referer": []string{"TRUE"},
				},
			},
			ctx:      &EvaluationContext{},
			expected: true,
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

func TestConditionEvaluator_BinaryEquals(t *testing.T) {
	ce := NewConditionEvaluator()

	tests := []struct {
		name         string
		op           ConditionOperator
		contextValue string
		policyValue  string
		expected     bool
	}{
		{
			name:         "binary equals match",
			op:           ConditionBinaryEquals,
			contextValue: "dGVzdA==",
			policyValue:  "dGVzdA==",
			expected:     true,
		},
		{
			name:         "binary equals no match",
			op:           ConditionBinaryEquals,
			contextValue: "dGVzdA==",
			policyValue:  "b3RoZXI=",
			expected:     false,
		},
		{
			name:         "binary equals empty",
			op:           ConditionBinaryEquals,
			contextValue: "",
			policyValue:  "dGVzdA==",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ce.matchValue(tt.op, tt.contextValue, tt.policyValue)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestConditionEvaluator_SetOperatorsStructuredValues pins the multivalued
// context-key semantics: a key's SessionContext entry is its verbatim value
// list, so a value containing a comma stays one value and an absent key has
// no values at all — ForAnyValue is then false, ForAllValues vacuously true.
func TestConditionEvaluator_SetOperatorsStructuredValues(t *testing.T) {
	evaluator := NewConditionEvaluator()

	tests := []struct {
		name       string
		operator   ConditionOperator
		policyVals []string
		session    map[string][]string
		expected   bool
	}{
		{
			name:       "comma inside a value matches that whole value only",
			operator:   ConditionForAnyValueStringEquals,
			policyVals: []string{"a,b"},
			session:    map[string][]string{"aws:principaltag/team": {"a,b"}},
			expected:   true,
		},
		{
			name:       "comma inside a value does not become two values",
			operator:   ConditionForAnyValueStringEquals,
			policyVals: []string{"a"},
			session:    map[string][]string{"aws:principaltag/team": {"a,b"}},
			expected:   false,
		},
		{
			name:       "each list element participates",
			operator:   ConditionForAnyValueStringEquals,
			policyVals: []string{"dev"},
			session:    map[string][]string{"aws:tagkeys": {"prod", "dev"}},
			expected:   true,
		},
		{
			name:       "for all values requires every element",
			operator:   ConditionForAllValuesStringEquals,
			policyVals: []string{"prod"},
			session:    map[string][]string{"aws:tagkeys": {"prod", "dev"}},
			expected:   false,
		},
		{
			name:       "for all values satisfied by both elements",
			operator:   ConditionForAllValuesStringEquals,
			policyVals: []string{"prod", "dev"},
			session:    map[string][]string{"aws:tagkeys": {"prod", "dev"}},
			expected:   true,
		},
		{
			name:       "absent key fails for any value",
			operator:   ConditionForAnyValueStringEquals,
			policyVals: []string{"prod"},
			session:    nil,
			expected:   false,
		},
		{
			name:       "absent key fails for any value not equals",
			operator:   ConditionForAnyValueStringNotEquals,
			policyVals: []string{"prod"},
			session:    nil,
			expected:   false,
		},
		{
			name:       "absent key is vacuously true for all values",
			operator:   ConditionForAllValuesStringEquals,
			policyVals: []string{"prod"},
			session:    nil,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions := ConditionMap{string(tt.operator): {"aws:TagKeys": tt.policyVals}}
			if tt.session != nil {
				// Each case exercises the key its session map carries.
				for k := range tt.session {
					conditions = ConditionMap{string(tt.operator): {k: tt.policyVals}}
					break
				}
			}
			ctx := &EvaluationContext{SessionContext: tt.session}
			if got := evaluator.Evaluate(conditions, ctx); got != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.expected, got)
			}
		})
	}
}

// TestEvaluationContextGetContextValueJoinsMultivaluedKey pins the scalar
// resolution of a multivalued key: GetContextValue joins the list, and
// ContextValues returns it verbatim.
func TestEvaluationContextGetContextValueJoinsMultivaluedKey(t *testing.T) {
	ctx := &EvaluationContext{SessionContext: map[string][]string{
		"aws:tagkeys": {"prod", "dev"},
	}}
	if got := ctx.GetContextValue("aws:TagKeys"); got != "prod,dev" {
		t.Errorf("GetContextValue: got %q, want %q", got, "prod,dev")
	}
	values := ctx.ContextValues("aws:TagKeys")
	if len(values) != 2 || values[0] != "prod" || values[1] != "dev" {
		t.Errorf("ContextValues: got %v, want [prod dev]", values)
	}
	if got := ctx.ContextValues("aws:Unknown"); got != nil {
		t.Errorf("ContextValues on an absent key: got %v, want nil", got)
	}
}
