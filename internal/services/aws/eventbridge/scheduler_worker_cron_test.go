package eventbridge

import (
	"fmt"
	"testing"
	"time"
)

// fireCheck describes one shouldFireSchedule observation. The creation
// time defaults to the evaluation instant so the table pins pure
// expression matching; window semantics have dedicated tests below.
type fireCheck struct {
	expr    string
	at      time.Time
	want    bool
	created time.Time
}

// runFireChecks evaluates each observation with a fresh rule ARN so the
// package-level last-fire bookkeeping cannot leak between cases.
func runFireChecks(t *testing.T, name string, checks []fireCheck) {
	t.Helper()
	for i, c := range checks {
		ruleARN := fmt.Sprintf("arn:aws:events:us-east-1:123456789012:rule/%s-%d", name, i)
		created := c.created
		if created.IsZero() {
			created = c.at
		}
		if got := shouldFireSchedule(ruleARN, c.expr, c.at, created); got != c.want {
			t.Errorf("shouldFireSchedule(%q, %s) = %v, want %v", c.expr, c.at.Format(time.RFC3339), got, c.want)
		}
	}
}

// TestShouldFireScheduleCronDayOfWeek pins the AWS day-of-week numbering
// (1=Sunday..7=Saturday) on the delegated evaluator.
func TestShouldFireScheduleCronDayOfWeek(t *testing.T) {
	runFireChecks(t, "dow", []fireCheck{
		{"cron(0 12 ? * 1 2027)", time.Date(2027, 1, 3, 12, 0, 7, 0, time.UTC), true, time.Time{}},  // Sunday
		{"cron(0 12 ? * 1 2027)", time.Date(2027, 1, 4, 12, 0, 7, 0, time.UTC), false, time.Time{}}, // Monday
		{"cron(0 12 ? * 7 2027)", time.Date(2027, 1, 2, 12, 0, 7, 0, time.UTC), true, time.Time{}},  // Saturday
		{"cron(0 12 ? * MON-FRI 2027)", time.Date(2027, 1, 4, 12, 0, 7, 0, time.UTC), true, time.Time{}},
		{"cron(0 12 ? * MON-FRI 2027)", time.Date(2027, 1, 9, 12, 0, 7, 0, time.UTC), false, time.Time{}}, // Saturday
		{"cron(0 12 ? * SUN 2027)", time.Date(2027, 1, 3, 12, 0, 7, 0, time.UTC), true, time.Time{}},
	})
}

// TestShouldFireScheduleCronLWHash pins the L, W and # day wildcards
// that the previous service-local matcher could not evaluate.
func TestShouldFireScheduleCronLWHash(t *testing.T) {
	runFireChecks(t, "lwh", []fireCheck{
		// Last Friday of January 2027 (Fridays: 1, 8, 15, 22, 29).
		{"cron(15 10 ? * 6L 2019-2022)", time.Date(2019, 1, 25, 10, 15, 0, 0, time.UTC), true, time.Time{}},
		{"cron(15 10 ? * 6L 2019-2022)", time.Date(2019, 1, 18, 10, 15, 0, 0, time.UTC), false, time.Time{}},
		{"cron(0 12 ? * 6L 2027)", time.Date(2027, 1, 29, 12, 0, 7, 0, time.UTC), true, time.Time{}},
		{"cron(0 12 ? * 6L 2027)", time.Date(2027, 1, 22, 12, 0, 7, 0, time.UTC), false, time.Time{}},
		// Nearest weekday: 2027-08-01 is a Sunday, so 1W fires on Monday.
		{"cron(0 12 1W * ? 2027)", time.Date(2027, 8, 2, 12, 0, 7, 0, time.UTC), true, time.Time{}},
		{"cron(0 12 1W * ? 2027)", time.Date(2027, 8, 1, 12, 0, 7, 0, time.UTC), false, time.Time{}},
		// Third Friday of January 2027.
		{"cron(0 12 ? * FRI#3 2027)", time.Date(2027, 1, 15, 12, 0, 7, 0, time.UTC), true, time.Time{}},
		{"cron(0 12 ? * FRI#3 2027)", time.Date(2027, 1, 8, 12, 0, 7, 0, time.UTC), false, time.Time{}},
		// Last day of the month.
		{"cron(0 0 L * ? *)", time.Date(2027, 1, 31, 0, 0, 7, 0, time.UTC), true, time.Time{}},
		{"cron(0 0 L * ? *)", time.Date(2027, 1, 30, 0, 0, 7, 0, time.UTC), false, time.Time{}},
		// Bare L in day-of-week is Saturday, every week.
		{"cron(0 12 ? * L 2027)", time.Date(2027, 1, 2, 12, 0, 7, 0, time.UTC), true, time.Time{}},
	})
}

// TestShouldFireScheduleCronEquivalence pins the firing behaviour of the
// plain expressions that were valid before the delegation so the switch
// to the shared engine cannot regress them.
func TestShouldFireScheduleCronEquivalence(t *testing.T) {
	runFireChecks(t, "eq", []fireCheck{
		{"cron(0 12 * * ? *)", time.Date(2027, 1, 2, 12, 0, 30, 0, time.UTC), true, time.Time{}},
		{"cron(0 12 * * ? *)", time.Date(2027, 1, 2, 12, 1, 0, 0, time.UTC), false, time.Time{}},
		{"cron(0 12 * * ? *)", time.Date(2027, 1, 2, 13, 0, 0, 0, time.UTC), false, time.Time{}},
		{"cron(0,30 8 * * ? *)", time.Date(2027, 1, 2, 8, 30, 1, 0, time.UTC), true, time.Time{}},
		{"cron(0,30 8 * * ? *)", time.Date(2027, 1, 2, 8, 15, 1, 0, time.UTC), false, time.Time{}},
		{"cron(*/10 * * * ? *)", time.Date(2027, 1, 2, 12, 10, 1, 0, time.UTC), true, time.Time{}},
		{"cron(*/10 * * * ? *)", time.Date(2027, 1, 2, 12, 15, 1, 0, time.UTC), false, time.Time{}},
		{"cron(0 12 15 JAN ? 2027)", time.Date(2027, 1, 15, 12, 0, 1, 0, time.UTC), true, time.Time{}},
		{"cron(0 12 15 JAN ? 2027)", time.Date(2027, 2, 15, 12, 0, 1, 0, time.UTC), false, time.Time{}},
		{"cron(0 12 ? * 1 2027)", time.Date(2028, 1, 2, 12, 0, 1, 0, time.UTC), false, time.Time{}}, // year gate
		// Five-field cron is not a valid AWS rule expression.
		{"cron(0 12 * * ?)", time.Date(2027, 1, 2, 12, 0, 1, 0, time.UTC), false, time.Time{}},
	})
}

// TestShouldFireScheduleCronLateBoundary pins the late-recovery contract
// shared with rate rules: a ticker gap that skips a matching minute
// still fires the missed boundary on the next evaluation instead of
// losing it silently.
func TestShouldFireScheduleCronLateBoundary(t *testing.T) {
	arn := "arn:aws:events:us-east-1:123456789012:rule/cron-late-probe"
	lastFireTimes.Delete(arn)
	creation := time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC)

	// Evaluations ran up to 12:04, then the sweep paused; the 12:05
	// boundary is missed and the next evaluation lands at 12:07.
	if !shouldFireSchedule(arn, "cron(0/5 * * * ? *)", creation.Add(7*time.Minute), creation) {
		t.Error("late evaluation did not fire the missed 12:05 boundary")
	}
	// The boundary fired once: a re-evaluation stays silent.
	if shouldFireSchedule(arn, "cron(0/5 * * * ? *)", creation.Add(7*time.Minute+30*time.Second), creation) {
		t.Error("missed boundary fired twice")
	}
	// The next boundary fires normally.
	if !shouldFireSchedule(arn, "cron(0/5 * * * ? *)", creation.Add(10*time.Minute+10*time.Second), creation) {
		t.Error("next boundary after recovery did not fire")
	}
}

// TestShouldFireScheduleCronCreationClamp pins that a boundary which
// predates the rule's creation never fires: a rule created at 12:01
// waits for the next 12:05-grid boundary, not the elapsed 12:00 one.
func TestShouldFireScheduleCronCreationClamp(t *testing.T) {
	arn := "arn:aws:events:us-east-1:123456789012:rule/cron-clamp-probe"
	lastFireTimes.Delete(arn)
	creation := time.Date(2027, 1, 2, 12, 1, 0, 0, time.UTC)

	if shouldFireSchedule(arn, "cron(0/5 * * * ? *)", creation.Add(time.Minute), creation) {
		t.Error("cron boundary before the rule creation fired")
	}
	if !shouldFireSchedule(arn, "cron(0/5 * * * ? *)", creation.Add(4*time.Minute+30*time.Second), creation) {
		t.Error("first boundary after the rule creation did not fire")
	}
}

// TestShouldFireScheduleCronOncePerMinute pins the once-per-minute cap:
// a second evaluation inside the same minute must not fire again, and
// the next minute fires again for an every-minute schedule.
func TestShouldFireScheduleCronOncePerMinute(t *testing.T) {
	ruleARN := "arn:aws:events:us-east-1:123456789012:rule/once-per-minute"
	now := time.Date(2027, 1, 2, 12, 0, 7, 0, time.UTC)
	if !shouldFireSchedule(ruleARN, "cron(* * * * ? *)", now, now) {
		t.Fatalf("first evaluation in the minute should fire")
	}
	if shouldFireSchedule(ruleARN, "cron(* * * * ? *)", now.Add(10*time.Second), now) {
		t.Errorf("second evaluation in the same minute should not fire")
	}
	if !shouldFireSchedule(ruleARN, "cron(* * * * ? *)", now.Add(time.Minute), now) {
		t.Errorf("evaluation in the next minute should fire")
	}
}
