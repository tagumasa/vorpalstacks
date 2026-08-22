package eventbridge

import (
	"testing"
	"time"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// TestSeedLastFireSuppressesRefireAfterRestart pins the restart re-seed
// of the dedup cache: a rule whose persisted fire marker survived a
// restart (empty in-memory cache) does not re-fire the boundary recorded
// in the marker, while later boundaries still fire.
func TestSeedLastFireSuppressesRefireAfterRestart(t *testing.T) {
	arn := "arn:aws:events:us-east-1:000000000000:rule/seed-probe"
	lastFireTimes.Delete(arn)
	defer lastFireTimes.Delete(arn)

	creation := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	// The 12:05 boundary fired before the restart and was persisted.
	rule := &eventsstore.Rule{
		ARN:                arn,
		ScheduleExpression: "rate(5 minutes)",
		CreatedAt:          creation,
		LastFiredAt:        creation.Add(5 * time.Minute),
	}

	seedLastFire(rule)
	if shouldFireSchedule(arn, rule.ScheduleExpression, creation.Add(7*time.Minute), creation) {
		t.Error("restart re-fired the boundary recorded in the persisted marker")
	}
	if !shouldFireSchedule(arn, rule.ScheduleExpression, creation.Add(10*time.Minute+time.Second), creation) {
		t.Error("boundary after the seeded marker did not fire")
	}

	// The seed never regresses an already-newer cached value.
	newer := creation.Add(10 * time.Minute)
	setLastFire(arn, newer)
	seedLastFire(rule)
	if last, ok := getLastFire(arn); !ok || !last.Equal(newer) {
		t.Errorf("seed regressed the cached fire time: %v (%v), want %v", last, ok, newer)
	}
}

// TestShouldFireScheduleRateAnchored pins the anchored rate() firing
// semantics: a rate rule does not fire on creation, the first fire
// happens one full interval after creation, boundaries stay pinned to
// the creation time, and each boundary fires exactly once.
func TestShouldFireScheduleRateAnchored(t *testing.T) {
	clearLastFire := func() { lastFireTimes.Delete("arn:aws:events:us-east-1:000000000000:rule/rate-probe") }
	clearLastFire()

	arn := "arn:aws:events:us-east-1:000000000000:rule/rate-probe"
	creation := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)

	// Mid-first-period: no fire yet.
	if shouldFireSchedule(arn, "rate(1 minute)", creation.Add(30*time.Second), creation) {
		t.Error("rate rule fired before the first interval elapsed")
	}
	// Evaluation before creation completes even later stays silent.
	if shouldFireSchedule(arn, "rate(1 minute)", creation.Add(59*time.Second), creation) {
		t.Error("rate rule fired before the first interval elapsed")
	}

	// One interval (plus a tick) after creation: fires once.
	if !shouldFireSchedule(arn, "rate(1 minute)", creation.Add(61*time.Second), creation) {
		t.Error("rate rule did not fire one interval after creation")
	}
	// Same period, later evaluation: no second fire.
	if shouldFireSchedule(arn, "rate(1 minute)", creation.Add(90*time.Second), creation) {
		t.Error("rate rule fired twice within one period")
	}

	// Second period boundary: fires again, anchored to creation — the
	// 12:02:00 and 12:02:30 evaluations both map to the 12:02:00
	// boundary, which fires only once.
	if !shouldFireSchedule(arn, "rate(1 minute)", creation.Add(2*time.Minute), creation) {
		t.Error("rate rule did not fire at the second boundary")
	}
	if shouldFireSchedule(arn, "rate(1 minute)", creation.Add(2*time.Minute+30*time.Second), creation) {
		t.Error("rate rule fired twice for the second boundary")
	}
}

// TestShouldFireScheduleRateNoDrift pins that boundaries derive from the
// creation time, not from the last fire: after a late first evaluation
// the next boundary is still a whole multiple of the interval from
// creation.
func TestShouldFireScheduleRateNoDrift(t *testing.T) {
	arn := "arn:aws:events:us-east-1:000000000000:rule/rate-drift-probe"
	lastFireTimes.Delete(arn)
	creation := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)

	// First evaluation happens late, at 12:04:30 (4.5 periods in): fires
	// the latest passed boundary (12:04:00) once.
	if !shouldFireSchedule(arn, "rate(1 minute)", creation.Add(4*time.Minute+30*time.Second), creation) {
		t.Fatal("late first evaluation did not fire")
	}
	// At 12:05:00 the next anchored boundary arrives — under elapsed-since-
	// -fire semantics the rule would stay silent until 12:05:30.
	if !shouldFireSchedule(arn, "rate(1 minute)", creation.Add(5*time.Minute), creation) {
		t.Error("boundary drifted forward with the late first fire")
	}
}

// TestShouldFireScheduleRateLongInterval pins the no-immediate-fire rule
// for long intervals: a rate(1 hour) rule created 30 minutes ago must
// stay silent.
func TestShouldFireScheduleRateLongInterval(t *testing.T) {
	arn := "arn:aws:events:us-east-1:000000000000:rule/rate-hour-probe"
	lastFireTimes.Delete(arn)
	creation := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)

	if shouldFireSchedule(arn, "rate(1 hour)", creation.Add(30*time.Minute), creation) {
		t.Error("hourly rate rule fired half an interval after creation")
	}
	if !shouldFireSchedule(arn, "rate(1 hour)", creation.Add(61*time.Minute), creation) {
		t.Error("hourly rate rule did not fire one interval after creation")
	}
}
