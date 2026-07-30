package sesv2

import (
	"testing"
)

func TestIsValidIdentityFormat(t *testing.T) {
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
			if got := isValidIdentityFormat(tt.input); got != tt.want {
				t.Errorf("isValidIdentityFormat(%q) = %v, want %v", tt.input, got, tt.want)
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

func TestIsValidEmailAddress(t *testing.T) {
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
			if got := isValidEmailAddress(tt.input); got != tt.want {
				t.Errorf("isValidEmailAddress(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidContactListName(t *testing.T) {
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
			if got := isValidContactListName(tt.input); got != tt.want {
				t.Errorf("isValidContactListName(%q) = %v, want %v", tt.input, got, tt.want)
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
