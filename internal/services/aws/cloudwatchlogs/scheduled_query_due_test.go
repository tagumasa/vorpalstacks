package cloudwatchlogs

import (
	"testing"
	"time"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

func dueQuery(expr string, creation, lastBoundary time.Time) *logsstore.ScheduledQuery {
	sq := &logsstore.ScheduledQuery{
		ScheduleExpression: expr,
		CreationTime:       creation.UnixMilli(),
	}
	if !lastBoundary.IsZero() {
		sq.LastExecutedBoundary = lastBoundary.UnixMilli()
	}
	return sq
}

// TestScheduledQueryDueOncePerBoundary pins that each schedule boundary
// runs a scheduled query exactly once: after a run the query is not due
// again until the next boundary is reached.
func TestScheduledQueryDueOncePerBoundary(t *testing.T) {
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)
	sq := dueQuery("rate(1 hour)", creation, time.Time{})

	// Never run: the first run waits one full interval — the creation
	// boundary itself never runs the query.
	if _, due := scheduledQueryDue(sq, creation.Add(30*time.Second)); due {
		t.Error("query ran on the creation boundary")
	}
	// After the first boundary has been consumed, an evaluation before
	// the next boundary must not run it again (the next boundary is a
	// full hour after the first).
	sq.LastExecutedBoundary = creation.Add(time.Hour).UnixMilli()
	if _, due := scheduledQueryDue(sq, creation.Add(time.Hour+90*time.Second)); due {
		t.Error("query re-ran within one boundary period")
	}
	// The next boundary is due.
	if _, due := scheduledQueryDue(sq, creation.Add(2*time.Hour)); !due {
		t.Error("query did not run at the next boundary")
	}
}

// TestScheduledQueryDueWaitsForFirstBoundary pins that a query that has
// never run waits for its first reached boundary instead of running
// immediately: cron() and at() forms run at their first matching time,
// and rate() runs one full interval after creation.
func TestScheduledQueryDueWaitsForFirstBoundary(t *testing.T) {
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)

	cron := dueQuery("cron(0 6 * * ? *)", creation, time.Time{})
	if _, due := scheduledQueryDue(cron, creation.Add(time.Minute)); due {
		t.Error("cron query ran before its first matching minute")
	}

	at := dueQuery("at(2030-01-01T00:00:00)", creation, time.Time{})
	if _, due := scheduledQueryDue(at, creation.Add(time.Minute)); due {
		t.Error("future at() query ran immediately")
	}

	rate := dueQuery("rate(1 hour)", creation, time.Time{})
	if _, due := scheduledQueryDue(rate, creation.Add(time.Minute)); due {
		t.Error("rate query ran on the creation boundary")
	}
	if _, due := scheduledQueryDue(rate, creation.Add(time.Hour)); !due {
		t.Error("rate query did not run one interval after creation")
	}
}

// TestScheduledQueryDueTimezone pins that the schedule expression is
// evaluated in the query's configured timezone: a cron(0 9 * * ? *) query
// with an Asia/Tokyo timezone is due from 00:00 UTC (09:00 Tokyo)
// onwards, while the same expression without a timezone stays on UTC.
func TestScheduledQueryDueTimezone(t *testing.T) {
	creation := time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC)

	tokyo := dueQuery("cron(0 9 * * ? *)", creation, time.Time{})
	tokyo.Timezone = "Asia/Tokyo"
	// 23:59:30 UTC is 08:59:30 Tokyo: the boundary has not been reached.
	if _, due := scheduledQueryDue(tokyo, time.Date(2026, 12, 31, 23, 59, 30, 0, time.UTC)); due {
		t.Error("tokyo query due before the 09:00 Tokyo boundary")
	}
	tokyoRun := dueQuery("cron(0 9 * * ? *)", creation, time.Time{})
	tokyoRun.Timezone = "Asia/Tokyo"
	if _, due := scheduledQueryDue(tokyoRun, time.Date(2027, 1, 1, 0, 0, 30, 0, time.UTC)); !due {
		t.Error("tokyo query not due at 09:00:30 Tokyo time")
	}

	utc := dueQuery("cron(0 9 * * ? *)", creation, time.Time{})
	if _, due := scheduledQueryDue(utc, time.Date(2027, 1, 1, 0, 0, 30, 0, time.UTC)); due {
		t.Error("utc query due before the 09:00 UTC boundary")
	}
	if _, due := scheduledQueryDue(utc, time.Date(2027, 1, 1, 9, 0, 30, 0, time.UTC)); !due {
		t.Error("utc query not due at 09:00:30 UTC")
	}

	// An invalid timezone falls back to UTC instead of blocking runs.
	invalid := dueQuery("cron(0 9 * * ? *)", creation, time.Time{})
	invalid.Timezone = "Not/AZone"
	if _, due := scheduledQueryDue(invalid, time.Date(2027, 1, 1, 9, 0, 30, 0, time.UTC)); !due {
		t.Error("invalid timezone did not fall back to UTC")
	}
}

// TestScheduledQueryDueScheduleWindow pins the optional execution window:
// the query never runs before scheduleStartTime or after scheduleEndTime.
func TestScheduledQueryDueScheduleWindow(t *testing.T) {
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)

	start := dueQuery("rate(1 hour)", creation, time.Time{})
	start.ScheduleStartTime = creation.Add(time.Hour).UnixMilli()
	if _, due := scheduledQueryDue(start, creation.Add(30*time.Minute)); due {
		t.Error("query ran before scheduleStartTime")
	}
	if _, due := scheduledQueryDue(start, creation.Add(90*time.Minute)); !due {
		t.Error("query did not run after scheduleStartTime reached the first boundary")
	}

	end := dueQuery("rate(1 hour)", creation, time.Time{})
	end.ScheduleEndTime = creation.Add(2 * time.Hour).UnixMilli()
	if _, due := scheduledQueryDue(end, creation.Add(3*time.Hour)); due {
		t.Error("query ran after scheduleEndTime")
	}
}

// TestScheduledQueryDueLateExecutionClockDoesNotSkipBoundary pins that a
// late execution suppresses only the boundary it consumed: when the
// execution wall clock ran past the next boundary, the next boundary
// must still run (late) instead of being silently skipped. The
// deduplication marker records the consumed boundary value, never the
// execution clock; lastTriggeredTime stays the execution clock surfaced
// on the wire.
func TestScheduledQueryDueLateExecutionClockDoesNotSkipBoundary(t *testing.T) {
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)
	sq := dueQuery("rate(1 minute)", creation, time.Time{})

	firstBoundary := creation.Add(time.Minute)
	// The first boundary ran; its execution clock advanced 70 seconds,
	// past the second boundary.
	sq.LastExecutedBoundary = firstBoundary.UnixMilli()
	sq.LastTriggeredTime = firstBoundary.Add(70 * time.Second).UnixMilli()

	// The second boundary (firstBoundary + one minute) was never
	// executed; an evaluation 65 seconds after the first boundary must
	// still run it even though the recorded execution clock is ahead of
	// it, and the evaluated boundary is exactly the second one.
	boundary, due := scheduledQueryDue(sq, firstBoundary.Add(65*time.Second))
	if !due {
		t.Fatal("the unexecuted second boundary was skipped because the execution clock ran past it")
	}
	if want := firstBoundary.Add(time.Minute); !boundary.Equal(want) {
		t.Fatalf("evaluated boundary = %v, want %v", boundary, want)
	}
}
