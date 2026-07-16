package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAWSAlias(t *testing.T) {
	assert.True(t, IsAWSAlias("alias/aws/s3"))
	assert.True(t, IsAWSAlias("alias/aws/dynamodb"))
	assert.False(t, IsAWSAlias("alias/my-key"))
	assert.False(t, IsAWSAlias("my-key"))
	assert.False(t, IsAWSAlias(""))
}

func TestIsValidConditionOperator(t *testing.T) {
	valid := []string{
		"StringEquals", "StringNotEquals", "StringLike", "StringNotLike",
		"NumericEquals", "Bool", "IpAddress", "NotIpAddress",
		"ArnEquals", "ArnLike", "Null", "BinaryEquals",
		"ForAllValues:StringEquals", "ForAnyValue:StringEquals",
	}
	for _, op := range valid {
		assert.True(t, IsValidConditionOperator(op), op)
	}
	assert.False(t, IsValidConditionOperator("InvalidOp"))
	assert.False(t, IsValidConditionOperator(""))
}

func TestIsValidScheduleExpression(t *testing.T) {
	assert.True(t, IsValidScheduleExpression("rate(5 minutes)"))
	assert.True(t, IsValidScheduleExpression("rate(1 hour)"))
	assert.True(t, IsValidScheduleExpression("cron(0 10 * * ? *)"))
	assert.False(t, IsValidScheduleExpression("invalid"))
	assert.False(t, IsValidScheduleExpression(""))
}

func TestValidateAliasName(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, ValidateAliasName("alias/my-key"))
	})

	t.Run("aws alias", func(t *testing.T) {
		assert.NoError(t, ValidateAliasName("alias/aws/s3"))
	})

	t.Run("missing prefix", func(t *testing.T) {
		assert.Error(t, ValidateAliasName("my-key"))
	})

	t.Run("too short", func(t *testing.T) {
		assert.Error(t, ValidateAliasName("alias/"))
	})

	t.Run("invalid characters", func(t *testing.T) {
		assert.Error(t, ValidateAliasName("alias/my key"))
	})
}

func TestValidateBucketName(t *testing.T) {
	valid := []string{
		"my-bucket", "my.bucket", "my-bucket-123",
		"abc", "a1b", "log-bucket-1",
	}
	for _, name := range valid {
		t.Run("valid_"+name, func(t *testing.T) {
			assert.NoError(t, ValidateBucketName(name))
		})
	}

	t.Run("too short", func(t *testing.T) {
		assert.Error(t, ValidateBucketName("ab"))
	})
	t.Run("too long", func(t *testing.T) {
		assert.Error(t, ValidateBucketName("a"+string(make([]byte, 64))))
	})
	t.Run("uppercase", func(t *testing.T) {
		assert.Error(t, ValidateBucketName("MyBucket"))
	})
	t.Run("IP format", func(t *testing.T) {
		assert.Error(t, ValidateBucketName("192.168.1.1"))
	})
	t.Run("consecutive periods", func(t *testing.T) {
		assert.Error(t, ValidateBucketName("my..bucket"))
	})
	t.Run("period adjacent to hyphen", func(t *testing.T) {
		assert.Error(t, ValidateBucketName("my.-bucket"))
	})
	t.Run("starts with hyphen", func(t *testing.T) {
		assert.Error(t, ValidateBucketName("-mybucket"))
	})
	t.Run("ends with hyphen", func(t *testing.T) {
		assert.Error(t, ValidateBucketName("mybucket-"))
	})
}

func TestValidateHealthCheckConfig(t *testing.T) {
	t.Run("valid HTTP", func(t *testing.T) {
		cfg := &HealthCheckConfig{
			Type:             "HTTP",
			Port:             80,
			RequestInterval:  30,
			FailureThreshold: 3,
		}
		assert.NoError(t, ValidateHealthCheckConfig(cfg))
	})

	t.Run("nil config", func(t *testing.T) {
		assert.Error(t, ValidateHealthCheckConfig(nil))
	})

	t.Run("wrong type", func(t *testing.T) {
		assert.Error(t, ValidateHealthCheckConfig("string"))
	})

	t.Run("invalid IP", func(t *testing.T) {
		cfg := &HealthCheckConfig{
			IPAddress:        "not-an-ip",
			Type:             "TCP",
			Port:             80,
			RequestInterval:  30,
			FailureThreshold: 3,
		}
		assert.Error(t, ValidateHealthCheckConfig(cfg))
	})

	t.Run("port out of range", func(t *testing.T) {
		cfg := &HealthCheckConfig{
			Type:             "TCP",
			Port:             99999,
			RequestInterval:  30,
			FailureThreshold: 3,
		}
		assert.Error(t, ValidateHealthCheckConfig(cfg))
	})

	t.Run("invalid type", func(t *testing.T) {
		cfg := &HealthCheckConfig{
			Type:             "INVALID",
			Port:             80,
			RequestInterval:  30,
			FailureThreshold: 3,
		}
		assert.Error(t, ValidateHealthCheckConfig(cfg))
	})

	t.Run("interval out of range", func(t *testing.T) {
		cfg := &HealthCheckConfig{
			Type:             "TCP",
			Port:             80,
			RequestInterval:  60,
			FailureThreshold: 3,
		}
		assert.Error(t, ValidateHealthCheckConfig(cfg))
	})

	t.Run("threshold out of range", func(t *testing.T) {
		cfg := &HealthCheckConfig{
			Type:             "TCP",
			Port:             80,
			RequestInterval:  30,
			FailureThreshold: 20,
		}
		assert.Error(t, ValidateHealthCheckConfig(cfg))
	})
}

func TestParseHealthCheckConfig(t *testing.T) {
	m := map[string]interface{}{
		"ipAddress":                "10.0.0.1",
		"port":                     float64(443),
		"type":                     "HTTPS",
		"resourcePath":             "/health",
		"fullyQualifiedDomainName": "example.com",
		"requestInterval":          float64(10),
		"failureThreshold":         float64(2),
	}
	cfg, err := ParseHealthCheckConfig(m)
	assert.NoError(t, err)
	assert.Equal(t, "10.0.0.1", cfg.IPAddress)
	assert.Equal(t, 443, cfg.Port)
	assert.Equal(t, "HTTPS", cfg.Type)
	assert.Equal(t, "/health", cfg.ResourcePath)
	assert.Equal(t, "example.com", cfg.FullyQualifiedDomainName)
	assert.Equal(t, 10, cfg.RequestInterval)
	assert.Equal(t, 2, cfg.FailureThreshold)
}

func TestValidateHostedZoneName(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, ValidateHostedZoneName("example.com"))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Error(t, ValidateHostedZoneName(""))
	})

	t.Run("too long", func(t *testing.T) {
		assert.Error(t, ValidateHostedZoneName(string(make([]byte, 256))))
	})

	t.Run("invalid character", func(t *testing.T) {
		assert.Error(t, ValidateHostedZoneName("exam ple.com"))
	})

	t.Run("starts with hyphen", func(t *testing.T) {
		assert.Error(t, ValidateHostedZoneName("-example.com"))
	})

	t.Run("ends with dot", func(t *testing.T) {
		assert.Error(t, ValidateHostedZoneName("example.com."))
	})

	t.Run("consecutive hyphens", func(t *testing.T) {
		assert.Error(t, ValidateHostedZoneName("exa--mple.com"))
	})
}

func TestValidateKeyID(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		assert.NoError(t, ValidateKeyID("12345678-1234-1234-1234-123456789012"))
	})

	t.Run("valid alias", func(t *testing.T) {
		assert.NoError(t, ValidateKeyID("alias/my-key"))
	})

	t.Run("valid key ARN", func(t *testing.T) {
		assert.NoError(t, ValidateKeyID("arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"))
	})

	t.Run("valid alias ARN", func(t *testing.T) {
		assert.NoError(t, ValidateKeyID("arn:aws:kms:us-east-1:123456789012:alias/my-key"))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Error(t, ValidateKeyID(""))
	})

	t.Run("invalid format", func(t *testing.T) {
		assert.Error(t, ValidateKeyID("not-a-valid-key"))
	})

	t.Run("invalid ARN", func(t *testing.T) {
		assert.Error(t, ValidateKeyID("arn:aws:kms:us-east-1:123456789012:key/invalid"))
	})
}

func TestValidateKMSPolicyName(t *testing.T) {
	assert.NoError(t, ValidateKMSPolicyName("default"))
	assert.Error(t, ValidateKMSPolicyName("custom"))
	assert.Error(t, ValidateKMSPolicyName(""))
}

func TestValidatePasswordPolicy(t *testing.T) {
	t.Run("nil policy", func(t *testing.T) {
		assert.NoError(t, ValidatePasswordPolicy("anything", nil))
	})

	t.Run("meets all requirements", func(t *testing.T) {
		policy := &PasswordPolicy{
			MinimumLength:    8,
			RequireUppercase: true,
			RequireLowercase: true,
			RequireNumbers:   true,
			RequireSymbols:   true,
		}
		assert.NoError(t, ValidatePasswordPolicy("Abc123!@", policy))
	})

	t.Run("too short", func(t *testing.T) {
		policy := &PasswordPolicy{MinimumLength: 10}
		assert.Error(t, ValidatePasswordPolicy("short", policy))
	})

	t.Run("too long", func(t *testing.T) {
		policy := &PasswordPolicy{MaximumLength: 5}
		assert.Error(t, ValidatePasswordPolicy("longpassword", policy))
	})

	t.Run("missing uppercase", func(t *testing.T) {
		policy := &PasswordPolicy{RequireUppercase: true}
		assert.Error(t, ValidatePasswordPolicy("lowercase", policy))
	})

	t.Run("missing lowercase", func(t *testing.T) {
		policy := &PasswordPolicy{RequireLowercase: true}
		assert.Error(t, ValidatePasswordPolicy("UPPERCASE", policy))
	})

	t.Run("missing number", func(t *testing.T) {
		policy := &PasswordPolicy{RequireNumbers: true}
		assert.Error(t, ValidatePasswordPolicy("nonumber", policy))
	})

	t.Run("missing symbol", func(t *testing.T) {
		policy := &PasswordPolicy{RequireSymbols: true}
		assert.Error(t, ValidatePasswordPolicy("nosymbol", policy))
	})
}

func TestValidatePolicyName(t *testing.T) {
	assert.NoError(t, ValidatePolicyName("my-policy"))
	assert.Error(t, ValidatePolicyName(""))
	assert.Error(t, ValidatePolicyName(string(make([]byte, 129))))
	assert.Error(t, ValidatePolicyName("policy with spaces"))
}

func TestValidateProtocol(t *testing.T) {
	valid := []string{"http", "https", "email", "email-json", "sqs", "sms", "lambda", "application", "firehose"}
	for _, p := range valid {
		assert.NoError(t, ValidateProtocol(p), p)
	}
	assert.Error(t, ValidateProtocol("invalid"))
	assert.Error(t, ValidateProtocol(""))
}

func TestValidateRoleName(t *testing.T) {
	assert.NoError(t, ValidateRoleName("my-role"))
	assert.Error(t, ValidateRoleName(""))
	assert.Error(t, ValidateRoleName(string(make([]byte, 65))))
}

func TestValidateTopicName(t *testing.T) {
	assert.NoError(t, ValidateTopicName("my-topic"))
	assert.NoError(t, ValidateTopicName("MyTopic_123"))
	assert.Error(t, ValidateTopicName(""))
	assert.Error(t, ValidateTopicName(string(make([]byte, 257))))
	assert.Error(t, ValidateTopicName("topic with spaces"))
}

func TestValidateUserName(t *testing.T) {
	assert.NoError(t, ValidateUserName("my-user"))
	assert.NoError(t, ValidateUserName("user+name"))
	assert.Error(t, ValidateUserName(""))
	assert.Error(t, ValidateUserName(string(make([]byte, 65))))
	assert.Error(t, ValidateUserName("user name"))
}
