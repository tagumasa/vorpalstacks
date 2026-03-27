package policy

import (
	"testing"
	"time"
)

func TestEvaluationContext_ResolveVariable(t *testing.T) {
	tests := []struct {
		name string
		ctx  *EvaluationContext
		key  string
		want string
	}{
		{
			name: "custom variable",
			ctx: &EvaluationContext{
				Variables: map[string]string{"key": "value"},
			},
			key:  "key",
			want: "value",
		},
		{
			name: "aws username",
			ctx: &EvaluationContext{
				UserName: "testuser",
			},
			key:  "aws:username",
			want: "testuser",
		},
		{
			name: "aws userid",
			ctx: &EvaluationContext{
				UserID: "user123",
			},
			key:  "aws:userid",
			want: "user123",
		},
		{
			name: "aws principalarn",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
			},
			key:  "aws:principalarn",
			want: "arn:aws:iam::123456789012:user/test",
		},
		{
			name: "aws sourceip",
			ctx: &EvaluationContext{
				SourceIP: "192.168.1.1",
			},
			key:  "aws:sourceip",
			want: "192.168.1.1",
		},
		{
			name: "aws useragent",
			ctx: &EvaluationContext{
				UserAgent: "test-agent",
			},
			key:  "aws:useragent",
			want: "test-agent",
		},
		{
			name: "kms encryption context",
			ctx: &EvaluationContext{
				EncryptionContext: map[string]string{"key": "value"},
			},
			key:  "kms:EncryptionContext:key",
			want: "value",
		},
		{
			name: "unknown variable",
			ctx:  &EvaluationContext{},
			key:  "unknown",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ctx.ResolveVariable(tt.key)
			if got != tt.want {
				t.Errorf("ResolveVariable(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestEvaluationContext_GetContextValue(t *testing.T) {
	tests := []struct {
		name string
		ctx  *EvaluationContext
		key  string
		want string
	}{
		{
			name: "principal",
			ctx:  &EvaluationContext{Principal: "arn:aws:iam::123456789012:user/test"},
			key:  "principal",
			want: "arn:aws:iam::123456789012:user/test",
		},
		{
			name: "action",
			ctx:  &EvaluationContext{Action: "s3:GetObject"},
			key:  "action",
			want: "s3:GetObject",
		},
		{
			name: "resource",
			ctx:  &EvaluationContext{Resource: "arn:aws:s3:::bucket/object"},
			key:  "resource",
			want: "arn:aws:s3:::bucket/object",
		},
		{
			name: "sourceip",
			ctx:  &EvaluationContext{SourceIP: "192.168.1.1"},
			key:  "sourceip",
			want: "192.168.1.1",
		},
		{
			name: "unknown key",
			ctx:  &EvaluationContext{},
			key:  "unknown",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ctx.GetContextValue(tt.key)
			if got != tt.want {
				t.Errorf("GetContextValue(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestDecisionEffect(t *testing.T) {
	if DecisionEffectAllow != "Allow" {
		t.Errorf("expected Allow, got %s", DecisionEffectAllow)
	}
	if DecisionEffectDeny != "Deny" {
		t.Errorf("expected Deny, got %s", DecisionEffectDeny)
	}
	if DecisionEffectDefaultDeny != "DefaultDeny" {
		t.Errorf("expected DefaultDeny, got %s", DecisionEffectDefaultDeny)
	}
}

func TestPolicyEvaluator_New(t *testing.T) {
	evaluator := NewPolicyEvaluator()
	if evaluator == nil {
		t.Fatal("expected non-nil evaluator")
	}
	if evaluator.conditionEvaluator == nil {
		t.Error("expected non-nil condition evaluator")
	}
}

func TestPolicyEvaluator_Evaluate(t *testing.T) {
	evaluator := NewPolicyEvaluator()

	tests := []struct {
		name       string
		ctx        *EvaluationContext
		policies   []*Document
		wantEffect DecisionEffect
	}{
		{
			name: "explicit deny",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			policies: []*Document{
				{
					Version: "2012-10-17",
					Statement: []Statement{
						{
							Effect:   EffectDeny,
							Action:   ActionList{"s3:GetObject"},
							Resource: ResourceList{"arn:aws:s3:::bucket/object"},
						},
					},
				},
			},
			wantEffect: DecisionEffectDeny,
		},
		{
			name: "allow",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			policies: []*Document{
				{
					Version: "2012-10-17",
					Statement: []Statement{
						{
							Effect:   EffectAllow,
							Action:   ActionList{"s3:GetObject"},
							Resource: ResourceList{"arn:aws:s3:::bucket/object"},
						},
					},
				},
			},
			wantEffect: DecisionEffectAllow,
		},
		{
			name: "no matching statement",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			policies: []*Document{
				{
					Version: "2012-10-17",
					Statement: []Statement{
						{
							Effect:   EffectAllow,
							Action:   ActionList{"s3:PutObject"},
							Resource: ResourceList{"arn:aws:s3:::bucket/object"},
						},
					},
				},
			},
			wantEffect: DecisionEffectDefaultDeny,
		},
		{
			name: "empty policies",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			policies:   []*Document{},
			wantEffect: DecisionEffectDefaultDeny,
		},
		{
			name: "deny overrides allow",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			policies: []*Document{
				{
					Version: "2012-10-17",
					Statement: []Statement{
						{
							Effect:   EffectAllow,
							Action:   ActionList{"s3:GetObject"},
							Resource: ResourceList{"arn:aws:s3:::bucket/object"},
						},
					},
				},
				{
					Version: "2012-10-17",
					Statement: []Statement{
						{
							Effect:   EffectDeny,
							Action:   ActionList{"s3:GetObject"},
							Resource: ResourceList{"arn:aws:s3:::bucket/object"},
						},
					},
				},
			},
			wantEffect: DecisionEffectDeny,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := evaluator.Evaluate(tt.ctx, tt.policies)
			if decision.Effect != tt.wantEffect {
				t.Errorf("expected %v, got %v", tt.wantEffect, decision.Effect)
			}
		})
	}
}

func TestPolicyEvaluator_StatementMatches(t *testing.T) {
	evaluator := NewPolicyEvaluator()

	tests := []struct {
		name string
		ctx  *EvaluationContext
		stmt *Statement
		want bool
	}{
		{
			name: "principal match",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			stmt: &Statement{
				Effect:   EffectAllow,
				Action:   ActionList{"s3:GetObject"},
				Resource: ResourceList{"arn:aws:s3:::bucket/object"},
			},
			want: true,
		},
		{
			name: "principal no match",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			stmt: &Statement{
				Effect:    EffectAllow,
				Principal: &Principal{AWS: StringList{"arn:aws:iam::999999999999:user/other"}},
				Action:    ActionList{"s3:GetObject"},
				Resource:  ResourceList{"arn:aws:s3:::bucket/object"},
			},
			want: false,
		},
		{
			name: "action match",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			stmt: &Statement{
				Effect:   EffectAllow,
				Action:   ActionList{"s3:GetObject"},
				Resource: ResourceList{"*"},
			},
			want: true,
		},
		{
			name: "action no match",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			stmt: &Statement{
				Effect:   EffectAllow,
				Action:   ActionList{"s3:PutObject"},
				Resource: ResourceList{"*"},
			},
			want: false,
		},
		{
			name: "resource match",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			stmt: &Statement{
				Effect:   EffectAllow,
				Action:   ActionList{"*"},
				Resource: ResourceList{"arn:aws:s3:::bucket/object"},
			},
			want: true,
		},
		{
			name: "resource no match",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			stmt: &Statement{
				Effect:   EffectAllow,
				Action:   ActionList{"*"},
				Resource: ResourceList{"arn:aws:s3:::other-bucket/object"},
			},
			want: false,
		},
		{
			name: "wildcard action and resource",
			ctx: &EvaluationContext{
				Principal: "arn:aws:iam::123456789012:user/test",
				Action:    "s3:GetObject",
				Resource:  "arn:aws:s3:::bucket/object",
			},
			stmt: &Statement{
				Effect:   EffectAllow,
				Action:   ActionList{"*"},
				Resource: ResourceList{"*"},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evaluator.statementMatches(tt.ctx, tt.stmt)
			if got != tt.want {
				t.Errorf("statementMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicyEvaluator_EvaluateKeyPolicy(t *testing.T) {
	evaluator := NewPolicyEvaluator()

	invalidPolicy := `{"Version":`

	tests := []struct {
		name       string
		ctx        *EvaluationContext
		policyJSON string
		wantEffect DecisionEffect
	}{
		{
			name: "invalid policy",
			ctx: &EvaluationContext{
				Action: "s3:GetObject",
			},
			policyJSON: invalidPolicy,
			wantEffect: DecisionEffectDefaultDeny,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := evaluator.EvaluateKeyPolicy(tt.ctx, tt.policyJSON)
			if decision.Effect != tt.wantEffect {
				t.Errorf("expected %v, got %v", tt.wantEffect, decision.Effect)
			}
		})
	}
}

func TestEvaluationContext_resolveAWSVariable(t *testing.T) {
	tests := []struct {
		name string
		ctx  *EvaluationContext
		key  string
		want string
	}{
		{
			name: "aws:username",
			ctx:  &EvaluationContext{UserName: "testuser"},
			key:  "aws:username",
			want: "testuser",
		},
		{
			name: "aws:userid",
			ctx:  &EvaluationContext{UserID: "user123"},
			key:  "aws:userid",
			want: "user123",
		},
		{
			name: "aws:principalarn lowercase",
			ctx:  &EvaluationContext{Principal: "arn:aws:iam::123456789012:user/test"},
			key:  "aws:principalarn",
			want: "arn:aws:iam::123456789012:user/test",
		},
		{
			name: "aws:PrincipalArn mixed case",
			ctx:  &EvaluationContext{Principal: "arn:aws:iam::123456789012:user/test"},
			key:  "aws:PrincipalArn",
			want: "arn:aws:iam::123456789012:user/test",
		},
		{
			name: "aws:principalaccount",
			ctx:  &EvaluationContext{PrincipalAccount: "123456789012"},
			key:  "aws:principalaccount",
			want: "123456789012",
		},
		{
			name: "aws:sourceip",
			ctx:  &EvaluationContext{SourceIP: "192.168.1.1"},
			key:  "aws:sourceip",
			want: "192.168.1.1",
		},
		{
			name: "aws:useragent",
			ctx:  &EvaluationContext{UserAgent: "test-agent"},
			key:  "aws:useragent",
			want: "test-agent",
		},
		{
			name: "aws:epochtime",
			ctx:  &EvaluationContext{RequestTime: time.Unix(1672531200, 0)},
			key:  "aws:epochtime",
			want: "1672531200",
		},
		{
			name: "aws:RequestedRegion with service context",
			ctx: &EvaluationContext{
				ServiceContext: map[string]string{"region": "us-east-1"},
			},
			key:  "aws:RequestedRegion",
			want: "us-east-1",
		},
		{
			name: "aws:RequestedRegion without service context",
			ctx:  &EvaluationContext{},
			key:  "aws:RequestedRegion",
			want: "",
		},
		{
			name: "aws:RequestedAccount",
			ctx:  &EvaluationContext{PrincipalAccount: "123456789012"},
			key:  "aws:RequestedAccount",
			want: "123456789012",
		},
		{
			name: "unknown aws variable",
			ctx:  &EvaluationContext{},
			key:  "aws:unknown",
			want: "",
		},
		{
			name: "aws:referer lowercase",
			ctx:  &EvaluationContext{Referer: "https://example.com"},
			key:  "aws:referer",
			want: "https://example.com",
		},
		{
			name: "aws:Referer mixed case",
			ctx:  &EvaluationContext{Referer: "https://example.com"},
			key:  "aws:Referer",
			want: "https://example.com",
		},
		{
			name: "aws:referer empty",
			ctx:  &EvaluationContext{},
			key:  "aws:referer",
			want: "",
		},
		{
			name: "aws:securetransport true",
			ctx:  &EvaluationContext{SecureTransport: true},
			key:  "aws:securetransport",
			want: "true",
		},
		{
			name: "aws:SecureTransport false",
			ctx:  &EvaluationContext{SecureTransport: false},
			key:  "aws:SecureTransport",
			want: "false",
		},
		{
			name: "aws:tokenissuetime set",
			ctx:  &EvaluationContext{TokenIssueTime: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
			key:  "aws:tokenissuetime",
			want: "2026-01-01T00:00:00Z",
		},
		{
			name: "aws:tokenissuetime zero",
			ctx:  &EvaluationContext{},
			key:  "aws:tokenissuetime",
			want: "",
		},
		{
			name: "aws:multifactorauthpresent true",
			ctx:  &EvaluationContext{MultiFactorAuthPresent: true},
			key:  "aws:multifactorauthpresent",
			want: "true",
		},
		{
			name: "aws:MultiFactorAuthPresent false",
			ctx:  &EvaluationContext{MultiFactorAuthPresent: false},
			key:  "aws:MultiFactorAuthPresent",
			want: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ctx.resolveAWSVariable(tt.key)
			if got != tt.want {
				t.Errorf("resolveAWSVariable(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestEvaluationContext_resolveKMSVariable(t *testing.T) {
	tests := []struct {
		name string
		ctx  *EvaluationContext
		key  string
		want string
	}{
		{
			name: "kms encryption context",
			ctx: &EvaluationContext{
				EncryptionContext: map[string]string{"key": "value"},
			},
			key:  "kms:EncryptionContext:key",
			want: "value",
		},
		{
			name: "kms encryption context empty",
			ctx:  &EvaluationContext{},
			key:  "kms:EncryptionContext:key",
			want: "",
		},
		{
			name: "not kms variable",
			ctx:  &EvaluationContext{},
			key:  "kms:Other:key",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ctx.resolveKMSVariable(tt.key)
			if got != tt.want {
				t.Errorf("resolveKMSVariable(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}
