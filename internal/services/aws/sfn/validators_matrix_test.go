package sfn

import (
	"strings"
	"testing"
	"time"
)

// TestLambdaFamilyExactActionHole pins the matrix rule: only the carried
// lambda action validates; other lambda-family actions are refused at
// creation like every other non-carried family member.
func TestLambdaFamilyExactActionHole(t *testing.T) {
	def := `{"StartAt": "T", "States": {"T": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invokeBatch", "End": true}}}`
	err := validateDefinitionStructure(def, "STANDARD")
	if err == nil || !strings.Contains(err.Error(), "not a carried integration API") {
		t.Fatalf("lambda:invokeBatch error = %v, want the unsupported-integration rejection", err)
	}
	ok := `{"StartAt": "T", "States": {"T": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke", "End": true}}}`
	if err := validateDefinitionStructure(ok, "STANDARD"); err != nil {
		t.Fatalf("lambda:invoke was rejected: %v", err)
	}
}

// TestResourceNameCountsCharacters pins that the 1-80 name bound counts
// Unicode characters, not UTF-8 bytes.
func TestResourceNameCountsCharacters(t *testing.T) {
	long := strings.Repeat("昭", 81) // 81 characters, 243 bytes
	if err := validateResourceName(long); err == nil {
		t.Fatal("an 81-character name was accepted past the 80-character bound")
	}
	fit := strings.Repeat("昭", 80) // 80 characters, 240 bytes
	if err := validateResourceName(fit); err != nil {
		t.Fatalf("an 80-character name was rejected: %v", err)
	}
	if err := validateExecutionName(fit); err != nil {
		t.Fatalf("an 80-character execution name was rejected: %v", err)
	}
}

// TestWaitTimestampFractionDigits pins the RFC3339 fraction: one or more
// digits — not only three, six, or nine.
func TestWaitTimestampFractionDigits(t *testing.T) {
	for _, value := range []string{
		"2026-03-14T01:59:00Z",
		"2026-03-14T01:59:00.5Z",
		"2026-03-14T01:59:00.25Z",
		"2026-03-14T01:59:00.123456Z",
		"2026-03-14T01:59:00.123456789Z",
	} {
		if _, ok := parseASLTimestamp(value); !ok {
			t.Fatalf("RFC3339 timestamp %q was rejected", value)
		}
	}
	for _, value := range []string{
		"2026-03-14t01:59:00Z",  // lowercase T
		"2026-03-14T01:59:00z",  // lowercase Z
		"2026-03-14 01:59:00Z",  // missing T
		"2026-03-14T01:59:00.Z", // empty fraction
	} {
		if _, ok := parseASLTimestamp(value); ok {
			t.Fatalf("non-conforming timestamp %q was accepted", value)
		}
	}
	if _, ok := parseASLTimestamp("2026-03-14T01:59:00.5Z"); !ok {
		t.Fatal("unreachable")
	}
	_ = time.Time{}
}
