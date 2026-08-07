package sesv2

import (
	"fmt"
	"testing"
)

func TestValidateIdentityFormat(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"user@example.com", true},
		{"example.com", true},
		{"sub.example.com", true},
		{"", false},
		{"@@", false},
		{"@@@", false},
		{"a@b@c", false},
		{"nodomain", false},
		{"user@", false},
		{"@example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := validateIdentityFormat(tt.input); got != tt.want {
				t.Errorf("validateIdentityFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateFromAddress(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"user@example.com", true},
		{"example.com", true},
		{"", false},
		{"@@", false},
		{"a@b@c", false},
		{"user@", false},
		{"@example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := validateFromAddress(tt.input); got != tt.want {
				t.Errorf("validateFromAddress(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateEmailAddress(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"user@example.com", true},
		// Recipient addresses cannot be bare domains
		{"example.com", false},
		{"", false},
		{"@@", false},
		{"a@b@c", false},
		{"user@", false},
		{"@example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := validateEmailAddress(tt.input); got != tt.want {
				t.Errorf("validateEmailAddress(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateContactListName(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"MyList", true},
		{"my-list-123", true},
		{"list_name", true},
		{"", false},
		{"list with spaces", false},
		{"list!", false},
		{"list.name", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := validateContactListName(tt.input); got != tt.want {
				t.Errorf("validateContactListName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUnsupportedHandlebarsDetection(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"plain_text", "Hello world", false},
		{"simple_var", "{{name}}", false},
		{"if_block", "{{#if x}}yes{{/if}}", false},
		{"each_block", "{{#each items}}{{this}}{{/each}}", false},
		{"partial", "{{> header}}", true},
		{"atkey", "{{@key}}", true},
		{"atindex", "{{@index}}", true},
		{"block_param", "{{#each items as |item|}}{{item}}{{/each}}", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := false
			for _, p := range unsupportedHandlebarsPatterns {
				if p.MatchString(tt.input) {
					matched = true
					break
				}
			}
			if matched != tt.want {
				t.Errorf("unsupportedHandlebars check for %q: got %v, want %v", tt.input, matched, tt.want)
			}
		})
	}
}

func TestParseListManagementOptions(t *testing.T) {
	t.Run("nil_params", func(t *testing.T) {
		opts := parseListManagementOptions(nil)
		if opts != nil {
			t.Error("expected nil for nil params")
		}
	})

	t.Run("empty_map", func(t *testing.T) {
		opts := parseListManagementOptions(map[string]interface{}{})
		if opts != nil {
			t.Error("expected nil for empty map")
		}
	})

	t.Run("with_contact_list", func(t *testing.T) {
		params := map[string]interface{}{
			"ListManagementOptions": map[string]interface{}{
				"ContactListName": "my-list",
				"TopicName":       "news",
			},
		}
		opts := parseListManagementOptions(params)
		if opts == nil {
			t.Fatal("expected non-nil opts")
		}
		if opts.ContactListName != "my-list" {
			t.Errorf("ContactListName = %q, want %q", opts.ContactListName, "my-list")
		}
		if opts.TopicName != "news" {
			t.Errorf("TopicName = %q, want %q", opts.TopicName, "news")
		}
	})

	t.Run("without_contact_list", func(t *testing.T) {
		params := map[string]interface{}{
			"ListManagementOptions": map[string]interface{}{
				"TopicName": "news",
			},
		}
		opts := parseListManagementOptions(params)
		if opts != nil {
			t.Error("expected nil when ContactListName is absent")
		}
	})
}

func TestRejectTenantName(t *testing.T) {
	t.Run("no_tenant_name", func(t *testing.T) {
		err := rejectTenantName(map[string]interface{}{})
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("with_tenant_name", func(t *testing.T) {
		err := rejectTenantName(map[string]interface{}{
			"TenantName": "my-tenant",
		})
		if err == nil {
			t.Error("expected error for TenantName, got nil")
		}
	})
}

func TestValidateIdentityFormat_Whitespace(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{" user@example.com", false},
		{"user@example.com ", false},
		{"user\t@example.com", false},
		{"user\n@example.com", false},
		{"user@ example.com", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.input), func(t *testing.T) {
			if got := validateIdentityFormat(tt.input); got != tt.want {
				t.Errorf("validateIdentityFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateMailType(t *testing.T) {
	if !validateMailType("MARKETING") {
		t.Error("MARKETING should be valid")
	}
	if !validateMailType("TRANSACTIONAL") {
		t.Error("TRANSACTIONAL should be valid")
	}
	if validateMailType("INVALID") {
		t.Error("INVALID should be rejected")
	}
}

func TestValidateTlsPolicy(t *testing.T) {
	if !validateTlsPolicy("REQUIRE") {
		t.Error("REQUIRE should be valid")
	}
	if !validateTlsPolicy("OPTIONAL") {
		t.Error("OPTIONAL should be valid")
	}
	if validateTlsPolicy("NONE") {
		t.Error("NONE should be rejected")
	}
}

func TestValidateHttpsPolicy(t *testing.T) {
	if !validateHttpsPolicy("REQUIRE") {
		t.Error("REQUIRE should be valid")
	}
	if !validateHttpsPolicy("REQUIRE_OPEN_ONLY") {
		t.Error("REQUIRE_OPEN_ONLY should be valid")
	}
	if !validateHttpsPolicy("OPTIONAL") {
		t.Error("OPTIONAL should be valid")
	}
	if validateHttpsPolicy("NONE") {
		t.Error("NONE should be rejected")
	}
}

func TestValidateMaxDeliverySeconds(t *testing.T) {
	tests := []struct {
		input int32
		want  bool
	}{
		{300, true},
		{50400, true},
		{3600, true},
		{299, false},
		{50401, false},
		{0, false},
		{-1, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			if got := validateMaxDeliverySeconds(tt.input); got != tt.want {
				t.Errorf("validateMaxDeliverySeconds(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateConfigurationSetName(t *testing.T) {
	if !validateConfigurationSetName("my-config-set") {
		t.Error("my-config-set should be valid")
	}
	if validateConfigurationSetName("") {
		t.Error("empty name should be rejected")
	}
	if validateConfigurationSetName("name with spaces") {
		t.Error("spaces should be rejected")
	}
}

func TestValidateDkimSigningAttributesOrigin(t *testing.T) {
	if !validateDkimSigningAttributesOrigin("AWS_SES") {
		t.Error("AWS_SES should be valid")
	}
	if !validateDkimSigningAttributesOrigin("EXTERNAL") {
		t.Error("EXTERNAL should be valid")
	}
	if !validateDkimSigningAttributesOrigin("AWS_SES_US_EAST_1") {
		t.Error("AWS_SES_US_EAST_1 should be valid")
	}
	if validateDkimSigningAttributesOrigin("INVALID") {
		t.Error("INVALID should be rejected")
	}
}

func TestValidateEventTypes(t *testing.T) {
	if !validateEventTypes([]string{"SEND", "BOUNCE", "DELIVERY"}) {
		t.Error("valid event types should pass")
	}
	if validateEventTypes([]string{"SEND", "INVALID"}) {
		t.Error("invalid event type should be rejected")
	}
}

func TestCountEventDestinations(t *testing.T) {
	if countEventDestinations(true, false, false, false, false) != 1 {
		t.Error("one destination should return 1")
	}
	if countEventDestinations(true, true, false, false, false) != 2 {
		t.Error("two destinations should return 2")
	}
	if countEventDestinations(false, false, false, false, false) != 0 {
		t.Error("zero destinations should return 0")
	}
}

func TestValidatePolicyJSON(t *testing.T) {
	if err := validatePolicyJSON(`{"Version":"2012-10-17"}`); err != nil {
		t.Errorf("valid JSON should pass: %v", err)
	}
	if err := validatePolicyJSON(""); err == nil {
		t.Error("empty string should fail")
	}
	if err := validatePolicyJSON("not json"); err == nil {
		t.Error("invalid JSON should fail")
	}
}
