package eventbridge

import "testing"

// TestIsValidScheduleExpression pins the AWS contract for scheduled rule
// expressions: rate() accepts only minute/hour/day units, the value is a
// positive number, and the unit agrees in number with the value. cron()
// keeps its six-field validation, and at() is EventBridge Scheduler syntax
// that scheduled rules must reject.
func TestIsValidScheduleExpression(t *testing.T) {
	valid := []string{
		"",
		"rate(1 minute)",
		"rate(5 minutes)",
		"rate(1 hour)",
		"rate(3 hours)",
		"rate(1 day)",
		"rate(7 days)",
		"cron(0 12 * * ? *)",
	}
	for _, expr := range valid {
		if !isValidScheduleExpression(expr) {
			t.Errorf("isValidScheduleExpression(%q) = false, want true", expr)
		}
	}

	invalid := []string{
		// week is not a valid unit for scheduled rules
		"rate(1 week)",
		"rate(2 weeks)",
		// a value of 1 requires a singular unit, values above 1 a plural one
		"rate(1 minutes)",
		"rate(1 hours)",
		"rate(1 days)",
		"rate(5 minute)",
		"rate(3 hour)",
		"rate(7 day)",
		// the value must be a positive number
		"rate(0 minutes)",
		"rate(0 days)",
		// scheduler-only syntax and malformed expressions
		"at(2026-01-01T12:00:00)",
		"rate(5minutes)",
		"rate(-1 minutes)",
	}
	for _, expr := range invalid {
		if isValidScheduleExpression(expr) {
			t.Errorf("isValidScheduleExpression(%q) = true, want false", expr)
		}
	}
}
