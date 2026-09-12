package cognitoidentityprovider

import (
	"testing"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// An explicitly supplied MinimumLength of zero is an out-of-range value,
// not the "unset" marker the stored-policy check tolerates. The dotted
// query-form key, by contrast, is not a second accepted shape: the member
// reads as absent, so no policy is parsed at all.
func TestParsePasswordPolicyExplicitMinimumLengthZeroRejected(t *testing.T) {
	jsonPath := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies": map[string]interface{}{
			"PasswordPolicy": map[string]interface{}{
				"MinimumLength": 0,
			},
		},
	}}
	if _, err := parsePasswordPolicyWithBase(jsonPath, nil); err == nil {
		t.Fatal("expected an explicit MinimumLength of 0 to be rejected")
	}

	flat := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies.PasswordPolicy.MinimumLength": "0",
	}}
	policy, err := parsePasswordPolicyWithBase(flat, nil)
	if err != nil {
		t.Fatalf("query-form member must not surface a parse error: %v", err)
	}
	if policy != nil {
		t.Fatalf("query-form member parsed as a policy: %+v", policy)
	}
}

// The Smithy range boundaries {6, 99} apply to explicitly supplied values;
// an absent member keeps the stored or default value untouched.
func TestParsePasswordPolicyMinimumLengthBoundaries(t *testing.T) {
	for _, tc := range []struct {
		value   float64
		wantErr bool
	}{
		{6, false},
		{8, false},
		{99, false},
		{100, true},
	} {
		req := &request.ParsedRequest{Parameters: map[string]interface{}{
			"Policies": map[string]interface{}{
				"PasswordPolicy": map[string]interface{}{
					"MinimumLength": tc.value,
				},
			},
		}}
		policy, err := parsePasswordPolicyWithBase(req, nil)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("expected MinimumLength %v to be rejected", tc.value)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for MinimumLength %v: %v", tc.value, err)
		}
		if policy.MinimumLength != 8 && tc.value == 8 {
			t.Fatalf("expected the parsed MinimumLength to be 8, got %d", policy.MinimumLength)
		}
	}
}

// An absent MinimumLength member leaves the base policy value in place.
func TestParsePasswordPolicyAbsentMinimumLengthKeepsBase(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies": map[string]interface{}{
			"PasswordPolicy": map[string]interface{}{
				"RequireSymbols": false,
			},
		},
	}}
	policy, err := parsePasswordPolicyWithBase(req, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.MinimumLength != 8 {
		t.Fatalf("expected the default MinimumLength of 8, got %d", policy.MinimumLength)
	}
	if policy.RequireSymbols {
		t.Fatal("expected RequireSymbols to parse as false")
	}
}

// A present member of the wrong JSON type is a wire-shape violation: the
// parser rejects it instead of coercing or ignoring it.
func TestParsePasswordPolicyRejectsWrongTypeMembers(t *testing.T) {
	for name, value := range map[string]interface{}{
		"string": "eight",
		"bool":   true,
	} {
		req := &request.ParsedRequest{Parameters: map[string]interface{}{
			"Policies": map[string]interface{}{
				"PasswordPolicy": map[string]interface{}{"MinimumLength": value},
			},
		}}
		if _, err := parsePasswordPolicyWithBase(req, nil); err == nil {
			t.Fatalf("%s-typed MinimumLength accepted", name)
		}
	}

	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies": map[string]interface{}{
			"PasswordPolicy": map[string]interface{}{"RequireUppercase": "true"},
		},
	}}
	if _, err := parsePasswordPolicyWithBase(req, nil); err == nil {
		t.Fatal("string-typed RequireUppercase accepted")
	}
}

// The dotted LambdaConfig query-form keys read as absent: no trigger is
// configured through them.
func TestParseLambdaConfigIgnoresQueryFormKeys(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"LambdaConfig.PreSignUp": "arn:aws:lambda:us-east-1:000000000000:function:pre",
	}}
	if config := parseLambdaConfigWithBase(req, nil); config != nil {
		t.Fatalf("query-form LambdaConfig keys parsed: %+v", config)
	}

	jsonReq := &request.ParsedRequest{Parameters: map[string]interface{}{
		"LambdaConfig": map[string]interface{}{
			"PreSignUp": "arn:aws:lambda:us-east-1:000000000000:function:pre",
		},
	}}
	config := parseLambdaConfigWithBase(jsonReq, nil)
	if config == nil || config.PreSignUp == "" {
		t.Fatalf("JSON map member not parsed: %+v", config)
	}
}

// An AccountRecoverySetting member present but carrying no mechanisms parses
// as "not set", like every sibling config parser — it must not surface an
// empty setting that would overwrite the stored one.
func TestParseAccountRecoverySettingEmptyMechanismsIsUnset(t *testing.T) {
	for name, params := range map[string]map[string]interface{}{
		"empty object": {"AccountRecoverySetting": map[string]interface{}{}},
		"zero mechanisms": {"AccountRecoverySetting": map[string]interface{}{
			"RecoveryMechanisms": []interface{}{},
		}},
	} {
		req := &request.ParsedRequest{Parameters: params}
		if v := parseAccountRecoverySetting(req); v != nil {
			t.Fatalf("%s parsed as %+v, want nil", name, v)
		}
	}

	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"AccountRecoverySetting": map[string]interface{}{
			"RecoveryMechanisms": []interface{}{
				map[string]interface{}{"Name": "verified_email", "Priority": float64(1)},
			},
		},
	}}
	v := parseAccountRecoverySetting(req)
	if v == nil || len(v.RecoveryMechanisms) != 1 || v.RecoveryMechanisms[0].Name != "verified_email" {
		t.Fatalf("populated setting misparsed: %+v", v)
	}
}

// The sign-in policy names the permitted first authentication factors:
// valid factors are stored, off-enum values, SOFTWARE_TOKEN (documented as
// unsupported as a first factor) and length-boundary violations are
// rejected, and an absent member parses as no policy.
func TestParseSignInPolicy(t *testing.T) {
	valid := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies": map[string]interface{}{
			"SignInPolicy": map[string]interface{}{
				"AllowedFirstAuthFactors": []interface{}{"PASSWORD", "EMAIL_OTP"},
			},
		},
	}}
	policy, err := parseSignInPolicy(valid)
	if err != nil {
		t.Fatalf("valid sign-in policy rejected: %v", err)
	}
	if policy == nil || len(policy.AllowedFirstAuthFactors) != 2 ||
		policy.AllowedFirstAuthFactors[0] != "PASSWORD" || policy.AllowedFirstAuthFactors[1] != "EMAIL_OTP" {
		t.Fatalf("valid sign-in policy misparsed: %+v", policy)
	}

	absent := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies": map[string]interface{}{"PasswordPolicy": map[string]interface{}{"MinimumLength": 8}},
	}}
	if v, err := parseSignInPolicy(absent); err != nil || v != nil {
		t.Fatalf("absent member parsed as %+v (err %v), want nil/nil", v, err)
	}

	empty := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies": map[string]interface{}{"SignInPolicy": map[string]interface{}{}},
	}}
	if v, err := parseSignInPolicy(empty); err != nil || v == nil || v.AllowedFirstAuthFactors != nil {
		t.Fatalf("empty SignInPolicy parsed as %+v (err %v), want a set policy with no factors", v, err)
	}

	for name, factors := range map[string][]interface{}{
		"off-enum":       {"PASSWORD", "PASSKEY"},
		"software token": {"SOFTWARE_TOKEN"},
		"non-string":     {"PASSWORD", 7},
		"below minimum":  {},
		"above maximum":  {"PASSWORD", "EMAIL_OTP", "SMS_OTP", "WEB_AUTHN", "PASSWORD", "SMS_OTP"},
	} {
		req := &request.ParsedRequest{Parameters: map[string]interface{}{
			"Policies": map[string]interface{}{
				"SignInPolicy": map[string]interface{}{
					"AllowedFirstAuthFactors": factors,
				},
			},
		}}
		if v, err := parseSignInPolicy(req); err == nil {
			t.Fatalf("%s accepted: %+v", name, v)
		}
	}
}

// The End User Messaging SMS configuration replaces the SNS delivery
// configuration: both at once, a structure without the required CallerArn,
// and a Region other than the pool's own are all rejected, and a valid
// structure is stored in full.
func TestParseSmsConfigurationEumsSms(t *testing.T) {
	valid := &request.ParsedRequest{Parameters: map[string]interface{}{
		"SmsConfiguration": map[string]interface{}{
			"EumsSms": map[string]interface{}{
				"CallerArn":            "arn:aws:iam::123456789012:role/eums",
				"ExternalId":           "ext-id",
				"OriginationIdentity":  "+10000000000",
				"ConfigurationSetName": "cs",
				"InEntityId":           "in-entity",
				"InTemplateId":         "in-template",
			},
		},
	}}
	cfg, err := parseSmsConfiguration(valid, "us-east-1")
	if err != nil {
		t.Fatalf("valid EumsSms configuration rejected: %v", err)
	}
	if cfg == nil || cfg.EumsSms == nil || cfg.EumsSms.CallerArn != "arn:aws:iam::123456789012:role/eums" ||
		cfg.EumsSms.ExternalId != "ext-id" || cfg.EumsSms.OriginationIdentity != "+10000000000" ||
		cfg.EumsSms.ConfigurationSetName != "cs" || cfg.EumsSms.InEntityId != "in-entity" ||
		cfg.EumsSms.InTemplateId != "in-template" {
		t.Fatalf("EumsSms misparsed: %+v", cfg)
	}
	if cfg.SnsCallerArn != "" {
		t.Fatalf("SNS members must stay unset alongside EumsSms: %+v", cfg)
	}

	region := &request.ParsedRequest{Parameters: map[string]interface{}{
		"SmsConfiguration": map[string]interface{}{
			"EumsSms": map[string]interface{}{
				"CallerArn": "arn:aws:iam::123456789012:role/eums",
				"Region":    "us-east-1",
			},
		},
	}}
	if _, err := parseSmsConfiguration(region, "us-east-1"); err != nil {
		t.Fatalf("Region matching the pool Region rejected: %v", err)
	}

	for name, smsMap := range map[string]map[string]interface{}{
		"combined with SNS": {
			"SnsCallerArn": "arn:aws:iam::123456789012:role/sns",
			"EumsSms":      map[string]interface{}{"CallerArn": "arn:aws:iam::123456789012:role/eums"},
		},
		"missing CallerArn": {
			"EumsSms": map[string]interface{}{"ExternalId": "ext-id"},
		},
		"foreign Region": {
			"EumsSms": map[string]interface{}{
				"CallerArn": "arn:aws:iam::123456789012:role/eums",
				"Region":    "eu-west-1",
			},
		},
	} {
		req := &request.ParsedRequest{Parameters: map[string]interface{}{"SmsConfiguration": smsMap}}
		if v, err := parseSmsConfiguration(req, "us-east-1"); err == nil {
			t.Fatalf("%s accepted: %+v", name, v)
		}
	}

	sns := &request.ParsedRequest{Parameters: map[string]interface{}{
		"SmsConfiguration": map[string]interface{}{
			"SnsCallerArn": "arn:aws:iam::123456789012:role/sns",
		},
	}}
	cfg, err = parseSmsConfiguration(sns, "us-east-1")
	if err != nil {
		t.Fatalf("SNS-only configuration rejected: %v", err)
	}
	if cfg == nil || cfg.SnsCallerArn == "" || cfg.EumsSms != nil {
		t.Fatalf("SNS-only configuration misparsed: %+v", cfg)
	}
}

// The CustomAuthMode member of AdvancedSecurityAdditionalFlows accepts the
// enabled-mode enum only — the value has no OFF. The SDK validates the
// closed enum client-side, so the server-side rejection is pinned here.
func TestApplyUserPoolConfigMembersCustomAuthMode(t *testing.T) {
	valid := &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserPoolAddOns": map[string]interface{}{
			"AdvancedSecurityAdditionalFlows": map[string]interface{}{
				"CustomAuthMode": "AUDIT",
			},
		},
	}}
	pool := &cognitostore.UserPool{}
	if err := applyUserPoolConfigMembers(pool, valid, "us-east-1"); err != nil {
		t.Fatalf("valid CustomAuthMode rejected: %v", err)
	}
	if pool.UserPoolAddOns == nil || pool.UserPoolAddOns.AdvancedSecurityAdditionalFlows == nil ||
		pool.UserPoolAddOns.AdvancedSecurityAdditionalFlows.CustomAuthMode != "AUDIT" {
		t.Fatalf("CustomAuthMode misparsed: %+v", pool.UserPoolAddOns)
	}

	invalid := &request.ParsedRequest{Parameters: map[string]interface{}{
		"UserPoolAddOns": map[string]interface{}{
			"AdvancedSecurityAdditionalFlows": map[string]interface{}{
				"CustomAuthMode": "OFF",
			},
		},
	}}
	if err := applyUserPoolConfigMembers(&cognitostore.UserPool{}, invalid, "us-east-1"); err == nil {
		t.Fatal("CustomAuthMode outside the enabled-mode enum accepted")
	}
}

// A present AllowedFirstAuthFactors member of the wrong type is malformed
// and rejected — it must not fail open to an all-factors policy.
func TestParseSignInPolicyRejectsWrongTypeMember(t *testing.T) {
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"Policies": map[string]interface{}{
			"SignInPolicy": map[string]interface{}{
				"AllowedFirstAuthFactors": "PASSWORD",
			},
		},
	}}
	if v, err := parseSignInPolicy(req); err == nil {
		t.Fatalf("wrong-type AllowedFirstAuthFactors accepted: %+v", v)
	}
}
