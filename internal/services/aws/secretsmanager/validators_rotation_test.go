package secretsmanager

import "testing"

func TestValidateRotationRules_ExclusiveScheduleFields(t *testing.T) {
	tests := []struct {
		name          string
		days          int
		scheduleExpr  string
		duration      string
		expectError   bool
		errorContains string
	}{
		{"days only", 30, "", "", false, ""},
		{"schedule only", 0, "rate(30 days)", "", false, ""},
		{"neither", 0, "", "", false, ""},
		{"days with duration", 30, "", "24h", false, ""},
		{"both days and schedule", 30, "rate(30 days)", "", true, "but not both"},
		{"invalid days", 1001, "", "", true, ""},
		{"invalid schedule charset", 0, "rate(30 days);", "", true, ""},
		{"invalid duration format", 0, "rate(30 days)", "24x", true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRotationRules(tt.days, tt.scheduleExpr, tt.duration)
			if tt.expectError && err == nil {
				t.Errorf("validateRotationRules(%d, %q, %q) expected error, got nil", tt.days, tt.scheduleExpr, tt.duration)
			}
			if !tt.expectError && err != nil {
				t.Errorf("validateRotationRules(%d, %q, %q) unexpected error: %v", tt.days, tt.scheduleExpr, tt.duration, err)
			}
			if tt.expectError && err != nil && tt.errorContains != "" {
				if msg := err.Error(); !contains(msg, tt.errorContains) {
					t.Errorf("error message %q does not contain %q", msg, tt.errorContains)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
