package iam

import (
	"testing"

	"vorpalstacks/internal/services/aws/iam/policy"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

func TestGetServicePrincipal(t *testing.T) {
	tests := []struct {
		serviceName string
		expected    ServicePrincipal
	}{
		{"lambda", ServicePrincipalLambda},
		{"states", ServicePrincipalStates},
		{"events", ServicePrincipalEvents},
		{"scheduler", ServicePrincipalScheduler},
		{"cloudtrail", ServicePrincipalCloudTrail},
		{"timestream", ServicePrincipalTimestream},
		{"cognito", ServicePrincipalCognito},
		{"cognito-identity", ServicePrincipalCognito},
		{"unknown", "unknown.amazonaws.com"},
	}

	for _, tt := range tests {
		t.Run(tt.serviceName, func(t *testing.T) {
			result := GetServicePrincipal(tt.serviceName)
			if result != tt.expected {
				t.Errorf("GetServicePrincipal(%q) = %v, want %v", tt.serviceName, result, tt.expected)
			}
		})
	}
}

func TestParseARN(t *testing.T) {
	tests := []struct {
		name      string
		arn       string
		wantErr   bool
		partition string
		service   string
		region    string
		accountID string
		resource  string
		roleName  string
	}{
		{
			name:      "valid role ARN",
			arn:       "arn:aws:iam::123456789012:role/MyRole",
			wantErr:   false,
			partition: "aws",
			service:   "iam",
			region:    "",
			accountID: "123456789012",
			resource:  "role/MyRole",
			roleName:  "MyRole",
		},
		{
			name:      "valid role ARN with path",
			arn:       "arn:aws:iam::123456789012:role/path/MyRole",
			wantErr:   false,
			partition: "aws",
			service:   "iam",
			region:    "",
			accountID: "123456789012",
			resource:  "role/path/MyRole",
			roleName:  "MyRole",
		},
		{
			name:    "invalid ARN format - no prefix",
			arn:     "aws:iam::123456789012:role/MyRole",
			wantErr: true,
		},
		{
			name:    "invalid ARN format - too few parts",
			arn:     "arn:aws:iam::123456789012",
			wantErr: true,
		},
		{
			name:      "non-role ARN",
			arn:       "arn:aws:s3:::mybucket",
			wantErr:   false,
			partition: "aws",
			service:   "s3",
			region:    "",
			accountID: "",
			resource:  "mybucket",
			roleName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := arnutil.ParseARN(tt.arn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseARN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if result.Partition != tt.partition {
					t.Errorf("Partition = %v, want %v", result.Partition, tt.partition)
				}
				if result.Service != tt.service {
					t.Errorf("Service = %v, want %v", result.Service, tt.service)
				}
				if result.Region != tt.region {
					t.Errorf("Region = %v, want %v", result.Region, tt.region)
				}
				if result.AccountID != tt.accountID {
					t.Errorf("AccountID = %v, want %v", result.AccountID, tt.accountID)
				}
				if result.Resource != tt.resource {
					t.Errorf("Resource = %v, want %v", result.Resource, tt.resource)
				}
				if result.RoleName != tt.roleName {
					t.Errorf("RoleName = %v, want %v", result.RoleName, tt.roleName)
				}
			}
		})
	}
}

func TestParsedARN_IsCrossAccount(t *testing.T) {
	tests := []struct {
		name         string
		arn          string
		localAccount string
		expected     bool
	}{
		{
			name:         "cross account",
			arn:          "arn:aws:iam::123456789012:role/MyRole",
			localAccount: "987654321098",
			expected:     true,
		},
		{
			name:         "same account",
			arn:          "arn:aws:iam::123456789012:role/MyRole",
			localAccount: "123456789012",
			expected:     false,
		},
		{
			name:         "empty account",
			arn:          "arn:aws:s3:::mybucket",
			localAccount: "123456789012",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := arnutil.ParseARN(tt.arn)
			if err != nil {
				t.Fatalf("ParseARN() error = %v", err)
			}
			result := parsed.IsCrossAccount(tt.localAccount)
			if result != tt.expected {
				t.Errorf("IsCrossAccount(%q) = %v, want %v", tt.localAccount, result, tt.expected)
			}
		})
	}
}

func TestIsValidRoleARN(t *testing.T) {
	tests := []struct {
		name     string
		arn      string
		expected bool
	}{
		{
			name:     "valid role ARN",
			arn:      "arn:aws:iam::123456789012:role/MyRole",
			expected: true,
		},
		{
			name:     "valid role ARN with path",
			arn:      "arn:aws:iam::123456789012:role/path/MyRole",
			expected: true,
		},
		{
			name:     "invalid - not IAM service",
			arn:      "arn:aws:s3:::mybucket",
			expected: false,
		},
		{
			name:     "invalid - user ARN",
			arn:      "arn:aws:iam::123456789012:user/MyUser",
			expected: false,
		},
		{
			name:     "invalid - empty ARN",
			arn:      "",
			expected: false,
		},
		{
			name:     "invalid - malformed ARN",
			arn:      "not-an-arn",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := arnutil.IsValidRoleARN(tt.arn)
			if result != tt.expected {
				t.Errorf("IsValidRoleARN(%q) = %v, want %v", tt.arn, result, tt.expected)
			}
		})
	}
}

func TestBuildEvaluationContext(t *testing.T) {
	accountID := "123456789012"
	ctx := BuildEvaluationContext(accountID, "")

	if ctx.PrincipalAccount != accountID {
		t.Errorf("PrincipalAccount = %v, want %v", ctx.PrincipalAccount, accountID)
	}

	if ctx.Variables["aws:SourceAccount"] != accountID {
		t.Errorf("Variables[aws:SourceAccount] = %v, want %v", ctx.Variables["aws:SourceAccount"], accountID)
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		value    string
		expected bool
	}{
		{
			name:     "exact match",
			pattern:  "test",
			value:    "test",
			expected: true,
		},
		{
			name:     "exact mismatch",
			pattern:  "test",
			value:    "other",
			expected: false,
		},
		{
			name:     "wildcard - asterisk only",
			pattern:  "*",
			value:    "anything",
			expected: true,
		},
		{
			name:     "wildcard - prefix",
			pattern:  "test*",
			value:    "testing",
			expected: true,
		},
		{
			name:     "wildcard - prefix no match",
			pattern:  "test*",
			value:    "other",
			expected: false,
		},
		{
			name:     "wildcard - suffix",
			pattern:  "*test",
			value:    "mytest",
			expected: true,
		},
		{
			name:     "wildcard - middle",
			pattern:  "foo*bar",
			value:    "foobar",
			expected: true,
		},
		{
			name:     "wildcard - middle with content",
			pattern:  "foo*bar",
			value:    "foobazbar",
			expected: true,
		},
		{
			name:     "wildcard - middle no match",
			pattern:  "foo*bar",
			value:    "foobaz",
			expected: false,
		},
		{
			name:     "empty pattern",
			pattern:  "",
			value:    "test",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPattern(tt.pattern, tt.value)
			if result != tt.expected {
				t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.pattern, tt.value, result, tt.expected)
			}
		})
	}
}

func TestMatchesPrincipal(t *testing.T) {
	tests := []struct {
		name         string
		principal    *policy.Principal
		servicePrinc string
		expected     bool
	}{
		{
			name:         "nil principal",
			principal:    nil,
			servicePrinc: "lambda.amazonaws.com",
			expected:     false,
		},
		{
			name: "everyone",
			principal: &policy.Principal{
				Everyone: true,
			},
			servicePrinc: "lambda.amazonaws.com",
			expected:     true,
		},
		{
			name: "matching service",
			principal: &policy.Principal{
				Service: policy.StringList{"lambda.amazonaws.com"},
			},
			servicePrinc: "lambda.amazonaws.com",
			expected:     true,
		},
		{
			name: "non-matching service",
			principal: &policy.Principal{
				Service: policy.StringList{"ec2.amazonaws.com"},
			},
			servicePrinc: "lambda.amazonaws.com",
			expected:     false,
		},
		{
			name: "wildcard service",
			principal: &policy.Principal{
				Service: policy.StringList{"*"},
			},
			servicePrinc: "lambda.amazonaws.com",
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesPrincipal(tt.principal, tt.servicePrinc)
			if result != tt.expected {
				t.Errorf("matchesPrincipal() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMatchesAction(t *testing.T) {
	tests := []struct {
		name     string
		actions  policy.ActionList
		action   string
		expected bool
	}{
		{
			name:     "empty actions",
			actions:  nil,
			action:   "sts:AssumeRole",
			expected: false,
		},
		{
			name:     "matching action",
			actions:  policy.ActionList{"sts:AssumeRole"},
			action:   "sts:AssumeRole",
			expected: true,
		},
		{
			name:     "non-matching action",
			actions:  policy.ActionList{"s3:GetObject"},
			action:   "sts:AssumeRole",
			expected: false,
		},
		{
			name:     "wildcard action",
			actions:  policy.ActionList{"*"},
			action:   "sts:AssumeRole",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesAction(tt.actions, tt.action)
			if result != tt.expected {
				t.Errorf("matchesAction() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMatchConditionValue(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		value    string
		expected []string
		result   bool
	}{
		{
			name:     "StringEquals match",
			operator: "StringEquals",
			value:    "123456789012",
			expected: []string{"123456789012"},
			result:   true,
		},
		{
			name:     "StringEquals no match",
			operator: "StringEquals",
			value:    "123456789012",
			expected: []string{"987654321098"},
			result:   false,
		},
		{
			name:     "StringLike match with wildcard",
			operator: "StringLike",
			value:    "123456789012",
			expected: []string{"*"},
			result:   true,
		},
		{
			name:     "StringLike prefix match",
			operator: "StringLike",
			value:    "123456789012",
			expected: []string{"123*"},
			result:   true,
		},
		{
			name:     "StringLike no match",
			operator: "StringLike",
			value:    "123456789012",
			expected: []string{"987*"},
			result:   false,
		},
		{
			name:     "unknown operator",
			operator: "UnknownOp",
			value:    "value",
			expected: []string{"value"},
			result:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchConditionValue(tt.operator, tt.value, tt.expected)
			if result != tt.result {
				t.Errorf("matchConditionValue() = %v, want %v", result, tt.result)
			}
		})
	}
}

func TestEvaluateTrustPolicy(t *testing.T) {
	doc := &policy.Document{
		Statement: []policy.Statement{
			{
				Effect: policy.EffectAllow,
				Principal: &policy.Principal{
					Service: policy.StringList{"lambda.amazonaws.com"},
				},
				Action: policy.ActionList{"sts:AssumeRole"},
			},
		},
	}

	evalCtx := &policy.EvaluationContext{
		PrincipalAccount: "123456789012",
		Variables:        map[string]string{"aws:SourceAccount": "123456789012"},
	}

	err := EvaluateTrustPolicy(doc, string(ServicePrincipalLambda), evalCtx)
	if err != nil {
		t.Errorf("EvaluateTrustPolicy() error = %v", err)
	}

	docDeny := &policy.Document{
		Statement: []policy.Statement{
			{
				Effect: policy.EffectDeny,
				Principal: &policy.Principal{
					Service: policy.StringList{"lambda.amazonaws.com"},
				},
				Action: policy.ActionList{"sts:AssumeRole"},
			},
		},
	}

	err = EvaluateTrustPolicy(docDeny, string(ServicePrincipalLambda), evalCtx)
	if err != ErrRoleCannotBeAssumed {
		t.Errorf("EvaluateTrustPolicy() expected ErrRoleCannotBeAssumed, got %v", err)
	}

	docNoMatch := &policy.Document{
		Statement: []policy.Statement{
			{
				Effect: policy.EffectAllow,
				Principal: &policy.Principal{
					Service: policy.StringList{"ec2.amazonaws.com"},
				},
				Action: policy.ActionList{"sts:AssumeRole"},
			},
		},
	}

	err = EvaluateTrustPolicy(docNoMatch, string(ServicePrincipalLambda), evalCtx)
	if err != ErrRoleCannotBeAssumed {
		t.Errorf("EvaluateTrustPolicy() expected ErrRoleCannotBeAssumed, got %v", err)
	}
}

func TestErrorTypes(t *testing.T) {
	t.Run("InvalidParameterValueError", func(t *testing.T) {
		err := &InvalidParameterValueError{
			Code:    "TestCode",
			Message: "Test message",
		}
		if err.Error() != "Test message" {
			t.Errorf("Error() = %v, want Test message", err.Error())
		}
	})

	t.Run("InvalidArnError", func(t *testing.T) {
		err := &InvalidArnError{
			Message: "Invalid ARN",
		}
		if err.Error() != "Invalid ARN" {
			t.Errorf("Error() = %v, want Invalid ARN", err.Error())
		}
	})

	t.Run("ValidationException", func(t *testing.T) {
		err := &ValidationException{
			Message: "Validation failed",
		}
		if err.Error() != "Validation failed" {
			t.Errorf("Error() = %v, want Validation failed", err.Error())
		}
	})

	t.Run("InvalidParameterCombinationError", func(t *testing.T) {
		err := &InvalidParameterCombinationError{
			Message: "Invalid combination",
		}
		if err.Error() != "Invalid combination" {
			t.Errorf("Error() = %v, want Invalid combination", err.Error())
		}
	})
}

func TestErrorConstructors(t *testing.T) {
	roleArn := "arn:aws:iam::123456789012:role/TestRole"

	tests := []struct {
		name string
		fn   func(string) error
		want string
	}{
		{
			name: "NewLambdaRoleNotFoundError",
			fn:   NewLambdaRoleNotFoundError,
			want: "The role defined for the function (arn:aws:iam::123456789012:role/TestRole) does not exist.",
		},
		{
			name: "NewLambdaRoleCannotBeAssumedError",
			fn:   NewLambdaRoleCannotBeAssumedError,
			want: "The role defined for the function cannot be assumed by Lambda.",
		},
		{
			name: "NewInvalidRoleArnError",
			fn:   NewInvalidRoleArnError,
			want: "Invalid RoleArn: arn:aws:iam::123456789012:role/TestRole",
		},
		{
			name: "NewRoleCannotBeAssumedError",
			fn:   NewRoleCannotBeAssumedError,
			want: "Role arn:aws:iam::123456789012:role/TestRole is invalid or cannot be assumed.",
		},
		{
			name: "NewEventsRoleNotFoundError",
			fn:   NewEventsRoleNotFoundError,
			want: "Parameter RoleArn is not valid. Reason: Role arn:aws:iam::123456789012:role/TestRole does not exist.",
		},
		{
			name: "NewCloudTrailRoleError",
			fn:   NewCloudTrailRoleError,
			want: "Role ARN arn:aws:iam::123456789012:role/TestRole does not exist or cannot be assumed by CloudTrail.",
		},
		{
			name: "NewTimestreamRoleError",
			fn:   NewTimestreamRoleError,
			want: "Role arn:aws:iam::123456789012:role/TestRole is invalid or cannot be assumed by Timestream.",
		},
		{
			name: "NewCognitoRoleError",
			fn:   NewCognitoRoleError,
			want: "Role ARN arn:aws:iam::123456789012:role/TestRole is invalid or cannot be assumed.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(roleArn)
			if err.Error() != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, err.Error(), tt.want)
			}
		})
	}
}
