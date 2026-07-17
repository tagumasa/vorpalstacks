package sts

import (
	"encoding/base32"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDurationSecondsExtended(t *testing.T) {
	tests := []struct {
		name        string
		duration    int
		expected    int
		expectError bool
	}{
		{name: "default when zero", duration: 0, expected: 3600},
		{name: "minimum", duration: 900, expected: 900},
		{name: "maximum extended", duration: 129600, expected: 129600},
		{name: "valid middle", duration: 43200, expected: 43200},
		{name: "below minimum", duration: 899, expectError: true},
		{name: "above maximum extended", duration: 129601, expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateDurationSecondsExtended(tt.duration)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidDurationExtended, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateRootDurationSeconds(t *testing.T) {
	tests := []struct {
		name        string
		duration    int
		expected    int
		expectError bool
	}{
		{name: "default when zero", duration: 0, expected: 900},
		{name: "one second", duration: 1, expected: 1},
		{name: "maximum 900", duration: 900, expected: 900},
		{name: "valid middle", duration: 300, expected: 300},
		{name: "above maximum", duration: 901, expectError: true},
		{name: "negative", duration: -1, expectError: true},
		{name: "old role max rejected", duration: 43200, expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateRootDurationSeconds(tt.duration)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidRootDuration, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateFederationDurationSeconds(t *testing.T) {
	tests := []struct {
		name        string
		duration    int
		isRoot      bool
		expected    int
		expectError bool
	}{
		{name: "non-root default", duration: 0, isRoot: false, expected: 43200},
		{name: "root default", duration: 0, isRoot: true, expected: 3600},
		{name: "non-root minimum", duration: 900, isRoot: false, expected: 900},
		{name: "non-root maximum", duration: 129600, isRoot: false, expected: 129600},
		{name: "root at cap", duration: 3600, isRoot: true, expected: 3600},
		{name: "root above cap", duration: 3601, isRoot: true, expectError: true},
		{name: "root way above cap", duration: 129600, isRoot: true, expectError: true},
		{name: "below minimum", duration: 899, isRoot: false, expectError: true},
		{name: "above maximum", duration: 129601, isRoot: false, expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateFederationDurationSeconds(tt.duration, tt.isRoot)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateWebIdentityDurationSeconds(t *testing.T) {
	tests := []struct {
		name        string
		duration    int
		expected    int
		expectError bool
	}{
		{name: "default when zero", duration: 0, expected: 300},
		{name: "minimum 60", duration: 60, expected: 60},
		{name: "maximum 3600", duration: 3600, expected: 3600},
		{name: "valid middle", duration: 300, expected: 300},
		{name: "below minimum", duration: 59, expectError: true},
		{name: "above maximum", duration: 3601, expectError: true},
		{name: "old role min rejected", duration: 900, expected: 900},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateWebIdentityDurationSeconds(tt.duration)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidWebIdentityDuration, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateRoleSessionName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{name: "empty", input: "", expectError: true},
		{name: "one char too short", input: "a", expectError: true},
		{name: "two chars valid", input: "ab"},
		{name: "64 chars valid", input: "abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz01"},
		{name: "65 chars too long", input: "abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz012", expectError: true},
		{name: "valid alphanumeric", input: "MySession123"},
		{name: "valid special chars", input: "user.name@domain-test+value=value,com"},
		{name: "invalid space", input: "my session", expectError: true},
		{name: "invalid slash", input: "my/session", expectError: true},
		{name: "invalid colon", input: "my:session", expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRoleSessionName(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSourceIdentity(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{name: "empty allowed", input: ""},
		{name: "valid name", input: "Alice"},
		{name: "aws prefix rejected", input: "aws:internal", expectError: true},
		{name: "one char too short", input: "a", expectError: true},
		{name: "two chars valid", input: "ab"},
		{name: "65 chars too long", input: "abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz012", expectError: true},
		{name: "invalid space", input: "my identity", expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSourceIdentity(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidSourceIdentity, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateMFACredentials(t *testing.T) {
	tests := []struct {
		name         string
		serialNumber string
		tokenCode    string
		expectError  bool
	}{
		{name: "both empty ok", serialNumber: "", tokenCode: ""},
		{name: "serial without code", serialNumber: "GAHT12345678", tokenCode: "", expectError: true},
		{name: "code without serial", serialNumber: "", tokenCode: "123456", expectError: true},
		{name: "valid pair", serialNumber: "GAHT12345678", tokenCode: "123456"},
		{name: "valid ARN serial", serialNumber: "arn:aws:iam::123456789012:mfa/user", tokenCode: "000000"},
		{name: "serial too short", serialNumber: "short", tokenCode: "123456", expectError: true},
		{name: "serial invalid chars", serialNumber: "GAHT 12345678", tokenCode: "123456", expectError: true},
		{name: "code not six digits", serialNumber: "GAHT12345678", tokenCode: "12345", expectError: true},
		{name: "code with letters", serialNumber: "GAHT12345678", tokenCode: "12345a", expectError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMFACredentials(tt.serialNumber, tt.tokenCode)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyTOTP(t *testing.T) {
	// Generate a known TOTP seed and compute the expected code for the
	// current time window, then verify round-trip.
	seed := "JBSWY3DPEHPK3PXP"
	now := time.Now().UTC()
	expectedCode := computeHOTP(decodeBase32(t, seed), now.Unix()/30)

	t.Run("valid current code", func(t *testing.T) {
		assert.True(t, verifyTOTP(seed, expectedCode, now))
	})

	t.Run("valid code in previous window", func(t *testing.T) {
		prevCode := computeHOTP(decodeBase32(t, seed), now.Unix()/30-1)
		assert.True(t, verifyTOTP(seed, prevCode, now))
	})

	t.Run("valid code in next window", func(t *testing.T) {
		nextCode := computeHOTP(decodeBase32(t, seed), now.Unix()/30+1)
		assert.True(t, verifyTOTP(seed, nextCode, now))
	})

	t.Run("invalid code", func(t *testing.T) {
		assert.False(t, verifyTOTP(seed, "000000", now.Add(2*time.Hour)))
	})

	t.Run("invalid seed", func(t *testing.T) {
		assert.False(t, verifyTOTP("!!!invalid!!!", expectedCode, now))
	})

	t.Run("unpadded base32 seed", func(t *testing.T) {
		// Many authenticator apps distribute seeds without base32 padding.
		// JBSWY3DPEHPK3PXP is 16 chars (no padding needed), but
		// GEZDGNBVGY3TQOJQ is also unpadded and valid.
		unpaddedSeed := "GEZDGNBVGY3TQOJQ"
		code := computeHOTP(decodeBase32(t, unpaddedSeed), now.Unix()/30)
		assert.True(t, verifyTOTP(unpaddedSeed, code, now))
	})
}

func decodeBase32(t *testing.T, seed string) []byte {
	t.Helper()
	dec, err := base32.StdEncoding.DecodeString(seed)
	require.NoError(t, err)
	return dec
}

func TestComputePackedPolicySize(t *testing.T) {
	params := map[string]interface{}{
		"PolicyArns.member.1.arn": "arn:aws:iam::123456789012:policy/test1",
		"PolicyArns.member.2.arn": "arn:aws:iam::123456789012:policy/test2",
	}

	t.Run("empty policy and params", func(t *testing.T) {
		size := computePackedPolicySize("", map[string]interface{}{}, nil)
		assert.Equal(t, int32(0), size)
	})

	t.Run("policy only", func(t *testing.T) {
		policy := `{"Version":"2012-10-17"}`
		size := computePackedPolicySize(policy, map[string]interface{}{}, nil)
		expected := int32((len(policy) * 100) / 2048)
		assert.Equal(t, expected, size)
	})

	t.Run("policy with arns", func(t *testing.T) {
		policy := `{"Version":"2012-10-17"}`
		size := computePackedPolicySize(policy, params, nil)
		totalLen := len(policy) + len("arn:aws:iam::123456789012:policy/test1") + len("arn:aws:iam::123456789012:policy/test2")
		expected := int32((totalLen * 100) / 2048)
		assert.Equal(t, expected, size)
	})

	t.Run("policy with tags", func(t *testing.T) {
		policy := `{"Version":"2012-10-17"}`
		tags := map[string]string{"Department": "Engineering", "Project": "Alpha"}
		size := computePackedPolicySize(policy, map[string]interface{}{}, tags)
		tagLen := len("Department") + len("Engineering") + len("Project") + len("Alpha")
		expected := int32(((len(policy) + tagLen) * 100) / 2048)
		assert.Equal(t, expected, size)
	})

	t.Run("tags only", func(t *testing.T) {
		tags := map[string]string{"Env": "prod"}
		size := computePackedPolicySize("", map[string]interface{}{}, tags)
		expected := int32(((len("Env") + len("prod")) * 100) / 2048)
		assert.Equal(t, expected, size)
	})
}

func TestExtractSessionTags(t *testing.T) {
	t.Run("empty params returns nil", func(t *testing.T) {
		tags, err := extractSessionTags(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Nil(t, tags)
	})

	t.Run("single tag", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags.member.1.Key":   "Department",
			"Tags.member.1.Value": "Engineering",
		}
		tags, err := extractSessionTags(params)
		assert.NoError(t, err)
		assert.Equal(t, "Engineering", tags["Department"])
	})

	t.Run("multiple tags", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags.member.1.Key":   "Department",
			"Tags.member.1.Value": "Engineering",
			"Tags.member.2.Key":   "Project",
			"Tags.member.2.Value": "Alpha",
		}
		tags, err := extractSessionTags(params)
		assert.NoError(t, err)
		assert.Len(t, tags, 2)
	})

	t.Run("empty value allowed", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags.member.1.Key":   "Env",
			"Tags.member.1.Value": "",
		}
		tags, err := extractSessionTags(params)
		assert.NoError(t, err)
		assert.Equal(t, "", tags["Env"])
	})

	t.Run("invalid key characters", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags.member.1.Key":   "Bad;Key",
			"Tags.member.1.Value": "val",
		}
		_, err := extractSessionTags(params)
		assert.Error(t, err)
	})

	t.Run("key too long", func(t *testing.T) {
		longKey := strings.Repeat("a", 129)
		params := map[string]interface{}{
			"Tags.member.1.Key":   longKey,
			"Tags.member.1.Value": "val",
		}
		_, err := extractSessionTags(params)
		assert.Error(t, err)
	})

	t.Run("duplicate keys rejected", func(t *testing.T) {
		params := map[string]interface{}{
			"Tags.member.1.Key":   "Dept",
			"Tags.member.1.Value": "Eng",
			"Tags.member.2.Key":   "Dept",
			"Tags.member.2.Value": "Sales",
		}
		_, err := extractSessionTags(params)
		assert.Error(t, err)
	})
}

func TestExtractTransitiveTagKeys(t *testing.T) {
	t.Run("empty params", func(t *testing.T) {
		keys, err := extractTransitiveTagKeys(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Nil(t, keys)
	})

	t.Run("single key", func(t *testing.T) {
		params := map[string]interface{}{
			"TransitiveTagKeys.member.1": "Dept",
		}
		keys, err := extractTransitiveTagKeys(params)
		assert.NoError(t, err)
		assert.Len(t, keys, 1)
		assert.Contains(t, keys, "Dept")
	})

	t.Run("multiple keys", func(t *testing.T) {
		params := map[string]interface{}{
			"TransitiveTagKeys.member.1": "Dept",
			"TransitiveTagKeys.member.2": "Project",
		}
		keys, err := extractTransitiveTagKeys(params)
		assert.NoError(t, err)
		assert.Len(t, keys, 2)
	})

	t.Run("duplicate keys rejected", func(t *testing.T) {
		params := map[string]interface{}{
			"TransitiveTagKeys.member.1": "Dept",
			"TransitiveTagKeys.member.2": "Dept",
		}
		_, err := extractTransitiveTagKeys(params)
		assert.Error(t, err)
	})

	t.Run("invalid key characters", func(t *testing.T) {
		params := map[string]interface{}{
			"TransitiveTagKeys.member.1": "Bad;Key",
		}
		_, err := extractTransitiveTagKeys(params)
		assert.Error(t, err)
	})
}

func TestExtractPolicyNameFromArn(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "root-task arn", input: "arn:aws:iam::aws:policy/root-task/IAMAuditRootUserCredentials", expected: "IAMAuditRootUserCredentials"},
		{name: "simple policy arn", input: "arn:aws:iam::aws:policy/MyPolicy", expected: "MyPolicy"},
		{name: "bare name", input: "IAMAuditRootUserCredentials", expected: "IAMAuditRootUserCredentials"},
		{name: "empty", input: "", expected: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPolicyNameFromArn(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAllowedRootTaskPolicyNames(t *testing.T) {
	valid := []string{
		"IAMAuditRootUserCredentials",
		"IAMCreateRootUserPassword",
		"IAMDeleteRootUserCredentials",
		"S3UnlockBucketPolicy",
		"SQSUnlockQueuePolicy",
	}
	for _, name := range valid {
		t.Run("valid_"+name, func(t *testing.T) {
			assert.True(t, allowedRootTaskPolicyNames[name])
		})
	}

	t.Run("invalid policy name", func(t *testing.T) {
		assert.False(t, allowedRootTaskPolicyNames["AdministratorAccess"])
	})
}
