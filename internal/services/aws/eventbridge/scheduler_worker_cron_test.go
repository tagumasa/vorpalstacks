package eventbridge

import (
	"fmt"
	"testing"
	"time"
)

// fireCheck describes one shouldFireCron observation.
type fireCheck struct {
	expr string
	at   time.Time
	want bool
}

// runFireChecks evaluates each observation with a fresh rule ARN so the
// package-level last-fire bookkeeping cannot leak between cases.
func runFireChecks(t *testing.T, name string, checks []fireCheck) {
	t.Helper()
	for i, c := range checks {
		ruleARN := fmt.Sprintf("arn:aws:events:us-east-1:123456789012:rule/%s-%d", name, i)
		if got := shouldFireCron(ruleARN, c.expr, c.at); got != c.want {
			t.Errorf("shouldFireCron(%q, %s) = %v, want %v", c.expr, c.at.Format(time.RFC3339), got, c.want)
		}
	}
}

// TestShouldFireCronDayOfWeek pins the AWS day-of-week numbering
// (1=Sunday..7=Saturday) on the delegated evaluator.
func TestShouldFireCronDayOfWeek(t *testing.T) {
	runFireChecks(t, "dow", []fireCheck{
		{"cron(0 12 ? * 1 2027)", time.Date(2027, 1, 3, 12, 0, 7, 0, time.UTC), true},  // Sunday
		{"cron(0 12 ? * 1 2027)", time.Date(2027, 1, 4, 12, 0, 7, 0, time.UTC), false}, // Monday
		{"cron(0 12 ? * 7 2027)", time.Date(2027, 1, 2, 12, 0, 7, 0, time.UTC), true},  // Saturday
		{"cron(0 12 ? * MON-FRI 2027)", time.Date(2027, 1, 4, 12, 0, 7, 0, time.UTC), true},
		{"cron(0 12 ? * MON-FRI 2027)", time.Date(2027, 1, 9, 12, 0, 7, 0, time.UTC), false}, // Saturday
		{"cron(0 12 ? * SUN 2027)", time.Date(2027, 1, 3, 12, 0, 7, 0, time.UTC), true},
	})
}

// TestShouldFireCronLWHash pins the L, W and # day wildcards that the
// previous service-local matcher could not evaluate.
func TestShouldFireCronLWHash(t *testing.T) {
	runFireChecks(t, "lwh", []fireCheck{
		// Last Friday of January 2027 (Fridays: 1, 8, 15, 22, 29).
		{"cron(15 10 ? * 6L 2019-2022)", time.Date(2019, 1, 25, 10, 15, 0, 0, time.UTC), true},
		{"cron(15 10 ? * 6L 2019-2022)", time.Date(2019, 1, 18, 10, 15, 0, 0, time.UTC), false},
		{"cron(0 12 ? * 6L 2027)", time.Date(2027, 1, 29, 12, 0, 7, 0, time.UTC), true},
		{"cron(0 12 ? * 6L 2027)", time.Date(2027, 1, 22, 12, 0, 7, 0, time.UTC), false},
		// Nearest weekday: 2027-08-01 is a Sunday, so 1W fires on Monday.
		{"cron(0 12 1W * ? 2027)", time.Date(2027, 8, 2, 12, 0, 7, 0, time.UTC), true},
		{"cron(0 12 1W * ? 2027)", time.Date(2027, 8, 1, 12, 0, 7, 0, time.UTC), false},
		// Third Friday of January 2027.
		{"cron(0 12 ? * FRI#3 2027)", time.Date(2027, 1, 15, 12, 0, 7, 0, time.UTC), true},
		{"cron(0 12 ? * FRI#3 2027)", time.Date(2027, 1, 8, 12, 0, 7, 0, time.UTC), false},
		// Last day of the month.
		{"cron(0 0 L * ? *)", time.Date(2027, 1, 31, 0, 0, 7, 0, time.UTC), true},
		{"cron(0 0 L * ? *)", time.Date(2027, 1, 30, 0, 0, 7, 0, time.UTC), false},
		// Bare L in day-of-week is Saturday, every week.
		{"cron(0 12 ? * L 2027)", time.Date(2027, 1, 2, 12, 0, 7, 0, time.UTC), true},
	})
}

// TestShouldFireCronEquivalence pins the firing behaviour of the plain
// expressions that were valid before the delegation so the switch to the
// shared engine cannot regress them.
func TestShouldFireCronEquivalence(t *testing.T) {
	runFireChecks(t, "eq", []fireCheck{
		{"cron(0 12 * * ? *)", time.Date(2027, 1, 2, 12, 0, 30, 0, time.UTC), true},
		{"cron(0 12 * * ? *)", time.Date(2027, 1, 2, 12, 1, 0, 0, time.UTC), false},
		{"cron(0 12 * * ? *)", time.Date(2027, 1, 2, 13, 0, 0, 0, time.UTC), false},
		{"cron(0,30 8 * * ? *)", time.Date(2027, 1, 2, 8, 30, 1, 0, time.UTC), true},
		{"cron(0,30 8 * * ? *)", time.Date(2027, 1, 2, 8, 15, 1, 0, time.UTC), false},
		{"cron(*/10 * * * ? *)", time.Date(2027, 1, 2, 12, 10, 1, 0, time.UTC), true},
		{"cron(*/10 * * * ? *)", time.Date(2027, 1, 2, 12, 15, 1, 0, time.UTC), false},
		{"cron(0 12 15 JAN ? 2027)", time.Date(2027, 1, 15, 12, 0, 1, 0, time.UTC), true},
		{"cron(0 12 15 JAN ? 2027)", time.Date(2027, 2, 15, 12, 0, 1, 0, time.UTC), false},
		{"cron(0 12 ? * 1 2027)", time.Date(2028, 1, 2, 12, 0, 1, 0, time.UTC), false}, // year gate
		// Five-field cron is not a valid AWS rule expression.
		{"cron(0 12 * * ?)", time.Date(2027, 1, 2, 12, 0, 1, 0, time.UTC), false},
	})
}

// TestShouldFireCronOncePerMinute pins the once-per-minute cap: a second
// evaluation inside the same minute must not fire again, and the next
// minute fires again for an every-minute schedule.
func TestShouldFireCronOncePerMinute(t *testing.T) {
	ruleARN := "arn:aws:events:us-east-1:123456789012:rule/once-per-minute"
	now := time.Date(2027, 1, 2, 12, 0, 7, 0, time.UTC)
	if !shouldFireCron(ruleARN, "cron(* * * * ? *)", now) {
		t.Fatalf("first evaluation in the minute should fire")
	}
	if shouldFireCron(ruleARN, "cron(* * * * ? *)", now.Add(10*time.Second)) {
		t.Errorf("second evaluation in the same minute should not fire")
	}
	if !shouldFireCron(ruleARN, "cron(* * * * ? *)", now.Add(time.Minute)) {
		t.Errorf("evaluation in the next minute should fire")
	}
}

// TestShouldFireRateUnchanged pins the elapsed-time rate() semantics,
// which the delegation intentionally leaves in place.
func TestShouldFireRateUnchanged(t *testing.T) {
	ruleARN := "arn:aws:events:us-east-1:123456789012:rule/rate-unchanged"
	now := time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC)
	if !shouldFireRate(ruleARN, "rate(1 minute)", now) {
		t.Fatalf("first rate evaluation should fire")
	}
	if shouldFireRate(ruleARN, "rate(1 minute)", now.Add(30*time.Second)) {
		t.Errorf("rate(1 minute) should not fire after 30 seconds")
	}
	if !shouldFireRate(ruleARN, "rate(1 minute)", now.Add(time.Minute)) {
		t.Errorf("rate(1 minute) should fire after one minute")
	}
}
