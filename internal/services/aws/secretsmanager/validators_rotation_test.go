package secretsmanager

import (
	"testing"
	"time"

	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

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

// TestValidateScheduleExpressionForms pins that the rotation schedule
// accepts the rate()/cron() forms AWS documents for RotationRules and
// rejects structural errors and the EventBridge-only at() form.
func TestValidateScheduleExpressionForms(t *testing.T) {
	valid := []string{
		"rate(1 minute)",
		"rate(30 days)",
		"cron(0 16 1,15 * ? *)",
	}
	for _, expr := range valid {
		if err := validateScheduleExpression(expr); err != nil {
			t.Errorf("validateScheduleExpression(%q) unexpected error: %v", expr, err)
		}
	}
	invalid := []string{
		"at(2026-01-01T00:00:00)", // one-shot form, not part of this contract
		"rate(1 days)",            // unit must agree with the value
		"cron(0 16 1,15 * *)",     // five fields
		"weekly",                  // not a schedule form
	}
	for _, expr := range invalid {
		if err := validateScheduleExpression(expr); err == nil {
			t.Errorf("validateScheduleExpression(%q) accepted invalid expression", expr)
		}
	}
}

// TestComputeNextRotationDate pins the schedule semantics: a fresh rate()
// schedule first fires one full period after configuration, a mid-period
// schedule fires at the next boundary, an overdue boundary is returned as a
// past due time so the checker fires immediately, cron() resolves to its
// next matching minute, and AutomaticallyAfterDays keeps its
// offset-from-last-rotation meaning.
func TestComputeNextRotationDate(t *testing.T) {
	now := time.Date(2026, 8, 26, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name  string
		rules *secretsmanagerstore.RotationRules
		last  time.Time
		want  time.Time
	}{
		{
			name:  "fresh rate schedule counts from configuration",
			rules: &secretsmanagerstore.RotationRules{ScheduleExpression: "rate(10 days)"},
			last:  time.Time{},
			want:  now.AddDate(0, 0, 10),
		},
		{
			name:  "mid-period rate fires at next boundary",
			rules: &secretsmanagerstore.RotationRules{ScheduleExpression: "rate(10 days)"},
			last:  now.AddDate(0, 0, -4),
			want:  now.AddDate(0, 0, 6),
		},
		{
			name:  "overdue rate boundary is owed now",
			rules: &secretsmanagerstore.RotationRules{ScheduleExpression: "rate(10 days)"},
			last:  now.AddDate(0, 0, -12),
			want:  now.AddDate(0, 0, -2),
		},
		{
			name:  "cron matching the current minute is due immediately",
			rules: &secretsmanagerstore.RotationRules{ScheduleExpression: "cron(0 12 * * ? *)"},
			last:  now.AddDate(0, 0, -1),
			want:  now, // daily noon, called at noon -> this minute
		},
		{
			name:  "cron resolves to next matching minute",
			rules: &secretsmanagerstore.RotationRules{ScheduleExpression: "cron(30 9 * * ? *)"},
			last:  now.AddDate(0, 0, -1),
			want:  time.Date(2026, 8, 27, 9, 30, 0, 0, time.UTC),
		},
		{
			name:  "days keep offset from last rotation",
			rules: &secretsmanagerstore.RotationRules{AutomaticallyAfterDays: 30},
			last:  now.AddDate(0, 0, -2),
			want:  now.AddDate(0, 0, 28),
		},
		{
			name:  "no rules means no schedule",
			rules: nil,
			last:  now,
			want:  time.Time{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := computeNextRotationDate(tc.rules, tc.last, now)
			if !got.Equal(tc.want) {
				t.Errorf("computeNextRotationDate(%+v) = %v, want %v", tc.rules, got, tc.want)
			}
		})
	}
}
