package timestreamquery

import (
	"strings"
	"testing"
)

// TestValidateQueryStringUnicodeLengths pins that QueryString @length(1,
// 262144) is counted in Unicode characters; the shape carries no pattern,
// so multibyte SQL literals are valid input and must not be rejected on
// byte length.
func TestValidateQueryStringUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateQueryString(strings.Repeat(cjk, 262144)); err != nil {
		t.Errorf("262144-character CJK query rejected: %v", err)
	}
	if err := validateQueryString(strings.Repeat(cjk, 262145)); err == nil {
		t.Error("262145-character CJK query accepted")
	}
	if err := validateQueryString(""); err == nil {
		t.Error("empty query accepted")
	}
}

// TestValidateTagKeyUnicodeLengths pins that TagKey @length(1,128) is
// counted in Unicode characters (the Timestream Query TagKey shape carries
// no pattern, unlike the general AWS tag guidance).
func TestValidateTagKeyUnicodeLengths(t *testing.T) {
	cjk := "\u65e5"

	if err := validateTagKey(strings.Repeat(cjk, 128)); err != nil {
		t.Errorf("128-character CJK tag key rejected: %v", err)
	}
	if err := validateTagKey(strings.Repeat(cjk, 129)); err == nil {
		t.Error("129-character CJK tag key accepted")
	}
	if err := validateTagValue(strings.Repeat(cjk, 256)); err != nil {
		t.Errorf("256-character CJK tag value rejected: %v", err)
	}
	if err := validateTagValue(strings.Repeat(cjk, 257)); err == nil {
		t.Error("257-character CJK tag value accepted")
	}
}

// TestValidateScheduleExpressionCronGrammar pins the full creation-time
// grammar: a cron() ScheduleExpression carries the shared AWS cron field
// grammar (an out-of-range field is rejected instead of being stored
// accepted and never triggering a run), the schedule is a cron or rate
// expression alone (the at() one-shot is an EventBridge Scheduler form),
// and rate values agree with their unit with no upper bound (the Timestream
// guide documents none).
func TestValidateScheduleExpressionCronGrammar(t *testing.T) {
	valid := []string{
		"cron(0 12 * * ? *)",
		"cron(59 23 31 DEC ? 2199)",
		"cron(0 0 L * ? *)",
		"rate(5 minutes)",
		"rate(1 minute)",
	}
	for _, expr := range valid {
		if err := validateScheduleExpression(expr); err != nil {
			t.Errorf("valid expression %q rejected: %v", expr, err)
		}
	}
	invalid := []string{
		"cron(60 12 * * ? *)",
		"cron(0 24 * * ? *)",
		"cron(0 12 32 * ? *)",
		"cron(0 12 1 13 ? *)",
		"cron(0 12 ? * 8 *)",
		"cron(0 12 * * ? 2200)",
		// Both day fields specified violates the ?-exclusivity rule.
		"cron(0 12 15 * MON 2027)",
		// The at() one-shot is not a Timestream schedule form.
		"at(2026-01-01T00:00:00)",
		// A malformed rate body or a value/unit disagreement is rejected.
		"rate(bogus)",
		"rate(5 minute)",
		"rate(0 minutes)",
	}
	for _, expr := range invalid {
		if err := validateScheduleExpression(expr); err == nil {
			t.Errorf("invalid expression %q accepted", expr)
		}
	}
}
