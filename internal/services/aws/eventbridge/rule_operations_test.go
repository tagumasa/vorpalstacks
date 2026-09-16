package eventbridge

import (
	"testing"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// TestRuleToMapShapeMembers pins the serialised Rule shapes against the
// model: the ListRules item shape (Rule) carries exactly {Arn,
// EventBusName, Name, State} plus the non-empty optionals — no timestamps,
// no CreatedBy — while DescribeRuleResponse additionally models CreatedBy.
func TestRuleToMapShapeMembers(t *testing.T) {
	r := &eventsstore.Rule{
		ARN:          "arn:aws:events:us-east-1:000000000000:rule/default/r1",
		Name:         "r1",
		EventBusName: "default",
		State:        eventsstore.RuleStateEnabled,
		Description:  "d",
		EventPattern: `{"source":["s"]}`,
		RoleARN:      "arn:aws:iam::000000000000:role/r",
		ManagedBy:    ".amazonaws",
		CreatedBy:    "arn:aws:iam::000000000000:root",
	}
	m := ruleToMap(r)
	for _, key := range []string{"CreationTime", "LastModifiedTime", "CreatedBy"} {
		if _, ok := m[key]; ok {
			t.Errorf("ListRules item carries unmodelled member %s", key)
		}
	}
	for _, key := range []string{"Arn", "Name", "EventBusName", "State", "Description", "EventPattern", "RoleArn", "ManagedBy"} {
		if _, ok := m[key]; !ok {
			t.Errorf("ListRules item missing modelled member %s", key)
		}
	}
}

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
