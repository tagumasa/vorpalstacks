package eventbridge

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	tagutil "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// putPatternRule is the shared fixture for the rule-plane contract tests:
// a pattern rule on the given bus.
func putPatternRule(t *testing.T, svc *EventsService, store *eventsstore.EventsStore, bus, name string) {
	t.Helper()
	if _, err := svc.putRuleCore(context.Background(), store, PutRuleInput{
		Name:                 name,
		EventBusName:         bus,
		EventBusNameProvided: true,
		EventPatternSet:      true,
		EventPattern:         `{"source":["contract.test"]}`,
	}); err != nil {
		t.Fatalf("put rule %s on %s: %v", name, bus, err)
	}
}

// TestPutRuleUpsertValidatesMembers pins that the upsert (already-exists)
// path rejects invalid member values exactly like the create path: the
// validation runs once before the create/upsert fork, and the atomic
// mutation only assigns pre-validated values.
func TestPutRuleUpsertValidatesMembers(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()

	// A pattern rule on the default bus; the upsert attempts below must
	// all be rejected before any assignment reaches the store.
	putPatternRule(t, svc, store, "default", "upsert-rule")

	cases := []struct {
		name string
		mod  func(*PutRuleInput)
	}{
		{"over-long description", func(in *PutRuleInput) {
			in.DescriptionSet = true
			in.Description = strings.Repeat("d", 513)
		}},
		{"over-long event pattern", func(in *PutRuleInput) {
			in.EventPatternSet = true
			in.EventPattern = `{"source":["` + strings.Repeat("x", 5000) + `"]}`
		}},
		{"structurally invalid event pattern", func(in *PutRuleInput) {
			in.EventPatternSet = true
			in.EventPattern = `{"source":`
		}},
		{"invalid schedule expression", func(in *PutRuleInput) {
			in.ScheduleExpressionSet = true
			in.ScheduleExpression = "rate(1 week)"
		}},
		{"invalid state", func(in *PutRuleInput) {
			in.StateSet = true
			in.State = "BOGUS"
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input := PutRuleInput{Name: "upsert-rule"}
			c.mod(&input)
			_, err := svc.putRuleCore(ctx, store, input)
			if err == nil {
				t.Fatalf("upsert with %s must be rejected", c.name)
			}
		})
	}
}

// TestPutRuleUpsertReplacesOmittedMembers pins the replacement contract:
// "If you are updating an existing rule, the rule is replaced with what
// you specify in this PutRule command. If you omit arguments in PutRule,
// the old values for those arguments are not kept. Instead, they are
// replaced with null values", "Rules are enabled by default", and "any
// tags you specify in the PutRule operation are ignored" (model
// operation documentation).
func TestPutRuleUpsertReplacesOmittedMembers(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()

	// A rule carrying every optional member, disabled and tagged.
	if _, err := svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                 "replace-rule",
		EventBusName:         "default",
		EventBusNameProvided: true,
		EventPatternSet:      true,
		EventPattern:         `{"source":["replace.test"]}`,
		DescriptionSet:       true,
		Description:          "original",
		RoleArnSet:           true,
		RoleArn:              "arn:aws:iam::000000000000:role/original",
		StateSet:             true,
		State:                "DISABLED",
		Tags:                 []tagutil.Tag{{Key: "stage", Value: "original"}},
	}); err != nil {
		t.Fatalf("create rule: %v", err)
	}

	// Upsert supplying the pattern alone: the omitted Description and
	// RoleArn clear, the omitted State resets to the default ENABLED,
	// and the update's tags are ignored (the creation tags remain).
	if _, err := svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                 "replace-rule",
		EventBusName:         "default",
		EventBusNameProvided: true,
		EventPatternSet:      true,
		EventPattern:         `{"source":["replace.test"]}`,
		Tags:                 []tagutil.Tag{{Key: "stage", Value: "ignored"}},
	}); err != nil {
		t.Fatalf("replace upsert: %v", err)
	}
	after, err := store.GetRule(ctx, "default", "replace-rule")
	if err != nil {
		t.Fatalf("get rule after replace: %v", err)
	}
	if after.Description != "" || after.RoleARN != "" {
		t.Fatalf("omitted members must clear: description=%q roleArn=%q", after.Description, after.RoleARN)
	}
	if after.State != eventsstore.RuleStateEnabled {
		t.Fatalf("omitted State must reset to ENABLED, got %q", after.State)
	}
	tags, err := store.TagStore.ListAsSlice(after.ARN)
	if err != nil {
		t.Fatalf("list tags: %v", err)
	}
	if len(tags) != 1 || tags[0].Key != "stage" || tags[0].Value != "original" {
		t.Fatalf("update tags must be ignored, tags=%+v", tags)
	}

	// An explicitly supplied State survives the replacement.
	if _, err := svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                 "replace-rule",
		EventBusName:         "default",
		EventBusNameProvided: true,
		EventPatternSet:      true,
		EventPattern:         `{"source":["replace.test"]}`,
		StateSet:             true,
		State:                "DISABLED",
	}); err != nil {
		t.Fatalf("explicit-state upsert: %v", err)
	}
	if after, err = store.GetRule(ctx, "default", "replace-rule"); err != nil {
		t.Fatalf("get rule after explicit state: %v", err)
	}
	if after.State != eventsstore.RuleStateDisabled {
		t.Fatalf("explicit State must survive the replacement, got %q", after.State)
	}
}

// TestPutRulePerBusQuota pins the documented per-bus rule-count quota at
// the Core layer: the 301st NEW rule on a bus answers
// LimitExceededException ("Maximum number of rules an account can have
// per event bus" — 300, EventBridge quotas page), while an upsert of an
// existing rule at the quota succeeds — the quota bounds the creation
// of new rules alone.
func TestPutRulePerBusQuota(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()

	for i := 0; i < eventsstore.MaxRulesPerEventBus; i++ {
		if err := store.CreateRule(ctx, &eventsstore.Rule{
			Name:         fmt.Sprintf("quota-rule-%03d", i),
			EventBusName: "default",
			State:        eventsstore.RuleStateEnabled,
		}); err != nil {
			t.Fatalf("seed rule %d: %v", i, err)
		}
	}

	_, err := svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                 "quota-overflow",
		EventBusName:         "default",
		EventBusNameProvided: true,
		EventPatternSet:      true,
		EventPattern:         `{"source":["quota.test"]}`,
	})
	if err == nil || !strings.Contains(err.Error(), "maximum of "+strconv.Itoa(eventsstore.MaxRulesPerEventBus)+" rules") {
		t.Fatalf("create at cap: got %v, want the LimitExceededException", err)
	}

	if _, err := svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                 "quota-rule-000",
		EventBusName:         "default",
		EventBusNameProvided: true,
		EventPatternSet:      true,
		EventPattern:         `{"source":["quota.test.updated"]}`,
	}); err != nil {
		t.Fatalf("upsert of an existing rule at the cap must succeed: %v", err)
	}
}

// TestDeleteRuleIdempotent pins the DeleteRule idempotency contract: "If
// you call delete rule multiple times for the same rule, all calls will
// succeed. When you call delete rule for a non-existent custom eventbus,
// ResourceNotFoundException is returned" (model operation documentation).
func TestDeleteRuleIdempotent(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()

	if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "idem-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	putPatternRule(t, svc, store, "idem-bus", "idem-rule")

	if err := svc.deleteRuleCore(ctx, store, DeleteRuleInput{EventBusName: "idem-bus", EventBusNameProvided: true, Name: "idem-rule"}); err != nil {
		t.Fatalf("first delete: %v", err)
	}
	if err := svc.deleteRuleCore(ctx, store, DeleteRuleInput{EventBusName: "idem-bus", EventBusNameProvided: true, Name: "idem-rule"}); err != nil {
		t.Fatalf("repeat delete must succeed, got: %v", err)
	}

	// A non-existent custom bus is the documented 404 case.
	err := svc.deleteRuleCore(ctx, store, DeleteRuleInput{EventBusName: "missing-bus", EventBusNameProvided: true, Name: "any"})
	if err == nil || !strings.Contains(err.Error(), "Event bus 'missing-bus' does not exist") {
		t.Fatalf("missing custom bus: got %v, want ResourceNotFoundException", err)
	}

	// The default bus always exists on AWS; here it may not have been
	// created yet, but no rule can exist without it, so the delete succeeds.
	if err := svc.deleteRuleCore(ctx, store, DeleteRuleInput{Name: "any-rule"}); err != nil {
		t.Fatalf("delete on uncreated default bus must succeed, got: %v", err)
	}
}

// TestEventBusNameArnFormResolved pins the EventBusNameOrArn contract: every
// rule-plane operation accepts "the name or ARN of the event bus" and the ARN
// form resolves to the same canonical name-addressed record.
func TestEventBusNameArnFormResolved(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()

	if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "arn-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	busARN := "arn:aws:events:us-east-1:000000000000:event-bus/arn-bus"

	// PutRule with the ARN form creates the rule on the named bus.
	if _, err := svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                 "arn-rule",
		EventBusName:         busARN,
		EventBusNameProvided: true,
		EventPatternSet:      true,
		EventPattern:         `{"source":["contract.test"]}`,
	}); err != nil {
		t.Fatalf("put rule via ARN form: %v", err)
	}
	if _, err := store.GetRule(ctx, "arn-bus", "arn-rule"); err != nil {
		t.Fatalf("rule must be stored under the canonical bus name: %v", err)
	}

	// DescribeRule resolves both forms to the same record.
	for _, form := range []string{"arn-bus", busARN} {
		res, err := svc.describeRuleCore(ctx, store, DescribeRuleInput{
			Name: "arn-rule", EventBusName: form, EventBusNameProvided: true,
		})
		if err != nil {
			t.Fatalf("describe via %q: %v", form, err)
		}
		if res.Rule.EventBusName != "arn-bus" {
			t.Fatalf("describe via %q: bus %q, want arn-bus", form, res.Rule.EventBusName)
		}
	}

	// DeleteRule via the ARN form removes it; the missing-bus 404 still
	// names the bus, not the rule.
	if err := svc.deleteRuleCore(ctx, store, DeleteRuleInput{
		Name: "arn-rule", EventBusName: busARN, EventBusNameProvided: true,
	}); err != nil {
		t.Fatalf("delete via ARN form: %v", err)
	}
}

// TestScheduledRuleDefaultBusOnly pins the user-guide restriction "You can
// only create scheduled rules using the default event bus" on both the
// create path and the upsert path.
func TestScheduledRuleDefaultBusOnly(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()

	if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "sched-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}

	// Create path: a schedule expression on a custom bus is rejected even
	// when the expression itself is valid.
	_, err := svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                  "sched-rule",
		EventBusName:          "sched-bus",
		EventBusNameProvided:  true,
		ScheduleExpressionSet: true,
		ScheduleExpression:    "rate(5 minutes)",
	})
	if err == nil || !strings.Contains(err.Error(), "Scheduled rules can only be created on the default event bus") {
		t.Fatalf("create on custom bus: got %v, want the default-bus ValidationException", err)
	}

	// The same expression on the default bus is accepted.
	if _, err := svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                  "sched-rule",
		ScheduleExpressionSet: true,
		ScheduleExpression:    "rate(5 minutes)",
	}); err != nil {
		t.Fatalf("create on default bus: %v", err)
	}

	// Upsert path: an existing pattern rule on a custom bus cannot be moved
	// onto a schedule either.
	putPatternRule(t, svc, store, "sched-bus", "pattern-rule")
	_, err = svc.putRuleCore(ctx, store, PutRuleInput{
		Name:                  "pattern-rule",
		EventBusName:          "sched-bus",
		EventBusNameProvided:  true,
		ScheduleExpressionSet: true,
		ScheduleExpression:    "cron(0 12 * * ? *)",
	})
	if err == nil || !strings.Contains(err.Error(), "Scheduled rules can only be created on the default event bus") {
		t.Fatalf("upsert onto schedule on custom bus: got %v, want the default-bus ValidationException", err)
	}
}

// TestListRuleNamesByTargetBeyondFirstTargetPage pins the inner per-rule
// target scan of listRuleNamesByTargetCore: a rule whose matching target
// sits past the first 100-entry page is still reported. The 5-targets-
// per-rule API quota keeps live rules single-page, so the fixture writes
// the targets at the store level directly.
func TestListRuleNamesByTargetBeyondFirstTargetPage(t *testing.T) {
	ctx := t.Context()
	svc := newEventBusCoreTestService()
	store := newEventBusCoreTestStore(t)
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}
	if err := store.CreateRule(ctx, &eventsstore.Rule{
		Name:         "many-targets",
		EventBusName: "default",
		State:        eventsstore.RuleStateEnabled,
	}); err != nil {
		t.Fatalf("create rule: %v", err)
	}
	for i := 0; i < 101; i++ {
		arn := fmt.Sprintf("arn:aws:sqs:us-east-1:000000000000:queue/filler-%03d", i)
		if i == 100 {
			// The matching reference sits on the second page.
			arn = "arn:aws:sqs:us-east-1:000000000000:queue/the-probe"
		}
		if err := store.PutTarget(ctx, &eventsstore.Target{
			ID:           fmt.Sprintf("t%d", i),
			RuleName:     "many-targets",
			EventBusName: "default",
			ARN:          arn,
		}); err != nil {
			t.Fatalf("put target %d: %v", i, err)
		}
	}

	result, err := svc.listRuleNamesByTargetCore(ctx, store, ListRuleNamesByTargetInput{
		TargetArn: "arn:aws:sqs:us-east-1:000000000000:queue/the-probe",
	})
	if err != nil {
		t.Fatalf("listRuleNamesByTargetCore: %v", err)
	}
	if len(result.RuleNames) != 1 || result.RuleNames[0] != "many-targets" {
		t.Fatalf("rule names = %v, want the rule whose target sits on the second page", result.RuleNames)
	}
}
