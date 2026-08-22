package timestreamquery

import (
	"testing"
	"time"
)

// TestShouldFireQueryBoundaryOnce pins that a scheduled query fires the
// latest reached boundary on the first evaluation at or after it — a late
// evaluation (ticker gap longer than the boundary interval) still fires
// the pending boundary — and never fires the same boundary twice.
func TestShouldFireQueryBoundaryOnce(t *testing.T) {
	e := &ScheduledQueryEngine{}
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)
	key := "us-east-1/query-late"

	// The 12:05 boundary's window passed without an evaluation; the 12:07
	// evaluation must still fire it and report the boundary itself.
	next, ok := e.shouldFireQuery(key, "rate(5 minutes)", creation.Add(7*time.Minute), creation, time.Time{})
	if !ok || !next.Equal(creation.Add(5*time.Minute)) {
		t.Fatalf("late evaluation: ok=%v next=%s, want the 12:05 boundary", ok, next)
	}
	// The same boundary re-evaluated after the fire stays silent.
	e.lastFired.Store(key, creation.Add(7*time.Minute))
	if _, ok := e.shouldFireQuery(key, "rate(5 minutes)", creation.Add(7*time.Minute+30*time.Second), creation, time.Time{}); ok {
		t.Error("boundary fired twice")
	}
}

// TestShouldFireQueryFutureStaysSilent pins that cron() and at() forms
// stay silent until their match or timestamp is reached, and that a past
// at() time fires exactly once.
func TestShouldFireQueryFutureStaysSilent(t *testing.T) {
	e := &ScheduledQueryEngine{}
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)

	if _, ok := e.shouldFireQuery("us-east-1/query-cron", "cron(0 6 * * ? *)", creation, creation, time.Time{}); ok {
		t.Error("cron query fired on a non-matching minute")
	}
	if _, ok := e.shouldFireQuery("us-east-1/query-at", "at(2027-06-01T08:30:00)", creation, creation, time.Time{}); ok {
		t.Error("at() query fired before its time")
	}
	// A past at() time fires exactly once, on the first evaluation after it.
	if _, ok := e.shouldFireQuery("us-east-1/query-at-past", "at(2027-01-01T11:00:00)", creation.Add(2*time.Hour), creation, time.Time{}); !ok {
		t.Error("past at() query never fired")
	}
}

// TestShouldFireQueryRateWaitsFirstInterval pins the AWS rate() contract:
// the scheduled query does not run on the creation boundary — the first
// run happens one full interval after creation.
func TestShouldFireQueryRateWaitsFirstInterval(t *testing.T) {
	e := &ScheduledQueryEngine{}
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)

	if _, ok := e.shouldFireQuery("us-east-1/query-first", "rate(5 minutes)", creation.Add(5*time.Second), creation, time.Time{}); ok {
		t.Error("scheduled query ran on the creation boundary")
	}
	next, ok := e.shouldFireQuery("us-east-1/query-first", "rate(5 minutes)", creation.Add(5*time.Minute+1*time.Second), creation, time.Time{})
	if !ok || !next.Equal(creation.Add(5*time.Minute)) {
		t.Errorf("first run = %v (%v), want the 12:05 boundary (true)", next, ok)
	}
}

// TestShouldFireQueryCronLateBoundary pins that a cron matching minute
// missed by a slow ticker fires on the next evaluation instead of being
// lost, reporting the missed boundary as the scheduled time.
func TestShouldFireQueryCronLateBoundary(t *testing.T) {
	e := &ScheduledQueryEngine{}
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)

	// The 12:05 matching minute passed unobserved; the 12:07 evaluation
	// must still fire it.
	next, ok := e.shouldFireQuery("us-east-1/query-cron-late", "cron(0/5 * * * ? *)", creation.Add(7*time.Minute), creation, time.Time{})
	if !ok || !next.Equal(creation.Add(5*time.Minute)) {
		t.Errorf("late cron evaluation = %v (%v), want the 12:05 boundary (true)", next, ok)
	}
	// A boundary that predates the creation never fires.
	if _, ok := e.shouldFireQuery("us-east-1/query-cron-clamp", "cron(0 12 * * ? *)", creation.Add(26*time.Hour), creation.Add(25*time.Hour), time.Time{}); ok {
		t.Error("cron boundary before the query creation fired")
	}
}

// TestShouldFireQueryPersistedLastRunSuppressesRefire pins the restart
// contract: the boundary consumed by an earlier run — persisted on the
// record as PreviousRunTime after every run, successful or failed — keeps
// a fresh engine (empty in-memory state) from re-running it, while a
// genuinely new boundary still fires.
func TestShouldFireQueryPersistedLastRunSuppressesRefire(t *testing.T) {
	e := &ScheduledQueryEngine{}
	creation := time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)
	key := "us-east-1/query-restart"

	// The 12:05 boundary was delivered by a previous process; the run
	// stamped its invocation time (12:05:02) on the record.
	lastRun := creation.Add(5*time.Minute + 2*time.Second)
	if _, ok := e.shouldFireQuery(key, "rate(5 minutes)", creation.Add(7*time.Minute), creation, lastRun); ok {
		t.Error("delivered boundary re-ran after restart")
	}
	// A boundary newer than the persisted last run still fires.
	if _, ok := e.shouldFireQuery(key, "rate(5 minutes)", creation.Add(12*time.Minute), creation, lastRun); !ok {
		t.Error("new boundary did not fire after restart")
	}
	// No persisted run (zero value): the elapsed boundary fires.
	if _, ok := e.shouldFireQuery(key, "rate(5 minutes)", creation.Add(7*time.Minute), creation, time.Time{}); !ok {
		t.Error("boundary did not fire without a persisted last run")
	}
}
