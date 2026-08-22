package scheduleexpr

import (
	"fmt"
	"strings"
	"testing"
)

// TestValidateRuleExpression pins the EventBridge PutRule validation
// profile: at() is not a rule expression, rate() follows the AWS
// agreement contract without upper bounds (the PutRule model specifies
// only the 256-character length trait), and cron() carries full
// field-level validation.
func TestValidateRuleExpression(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want bool
	}{
		{"empty is valid", "", true},
		{"rate singular", "rate(1 minute)", true},
		{"rate plural agreement", "rate(5 minutes)", true},
		{"rate disagreement singular value", "rate(1 minutes)", false},
		{"rate disagreement plural value", "rate(5 hour)", false},
		{"rate zero", "rate(0 minutes)", false},
		{"rate weeks rejected", "rate(2 weeks)", false},
		{"rate no documented bound", "rate(1000000 minutes)", true},
		{"at rejected", "at(2026-01-01T00:00:00)", false},
		{"cron six fields", "cron(0 12 * * ? *)", true},
		{"cron five fields", "cron(0 12 * * ?)", false},
		{"cron names", "cron(0 12 1 JAN ? 2027)", true},
		{"cron minute out of range", "cron(99 12 * * ? *)", false},
		{"cron list and step", "cron(0,30/10 12 1-15 * ? *)", true},
		{"cron malformed", "cron(not a cron)", false},
		{"cron last friday of month", "cron(15 10 ? * 6L 2019-2022)", true},
		{"cron weekday nearest day", "cron(0 9 1W * ? *)", true},
		{"cron nth weekday of month", "cron(0 9 ? * FRI#3 2027)", true},
		{"cron last day of month", "cron(0 0 L * ? *)", true},
		{"cron last day of week", "cron(59 23 ? * L 2030)", true},
		{"cron both day fields specified", "cron(0 12 15 * FRI 2027)", false},
		{"cron names both day fields", "cron(0 12 1 JAN MON 2027)", false},
		{"cron question mark in minutes", "cron(? 12 * * ? *)", false},
		{"cron question mark in year", "cron(0 12 ? * * ?)", false},
		{"cron W in day of week", "cron(0 12 ? * 3W 2027)", false},
		{"cron hash in day of month", "cron(0 12 3#2 * ? 2027)", false},
		{"cron hash list in day of week", "cron(0 12 ? * 1#1,6#3 2027)", false},
		{"cron bare W in day of month", "cron(0 12 W * ? *)", false},
		{"length beyond 256 rejected", "cron(" + strings.Repeat("0 ", 200) + "*)", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRuleExpression(tt.expr); got != tt.want {
				t.Errorf("ValidateRuleExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

// TestValidateExpressionSchedulerProfile pins the differences of the
// Scheduler/Timestream profile: at() is accepted and rate() values carry
// the documented upper bounds.
func TestValidateExpressionSchedulerProfile(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want bool
	}{
		{"at accepted", "at(2026-01-01T00:00:00)", true},
		{"at malformed", "at(2026-13-01T00:00:00)", false},
		{"rate at minute bound", fmt.Sprintf("rate(%d minutes)", MaxRateMinutes), true},
		{"rate beyond minute bound", fmt.Sprintf("rate(%d minutes)", MaxRateMinutes+1), false},
		{"rate beyond hour bound", fmt.Sprintf("rate(%d hours)", MaxRateHours+1), false},
		{"rate beyond day bound", fmt.Sprintf("rate(%d days)", MaxRateDays+1), false},
		{"rate singular minute", "rate(1 minute)", true},
		{"cron day fields both asterisk", "cron(0 9 * * * *)", false},
		{"cron day fields both values", "cron(0 9 15 * MON *)", false},
		{"cron dom asterisk dow question", "cron(0 9 * * ? *)", true},
		{"cron dom question dow asterisk", "cron(0 9 ? * * *)", true},
		{"cron dom value dow question", "cron(0 9 15 * ? *)", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateExpression(tt.expr); got != tt.want {
				t.Errorf("ValidateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}
