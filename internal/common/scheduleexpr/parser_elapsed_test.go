package scheduleexpr

import (
	"testing"
	"time"
)

// TestElapsedExecutionTimeRate pins the first-interval contract: the
// anchor itself never fires, and the latest fully elapsed period
// boundary is returned once now has reached it.
func TestElapsedExecutionTimeRate(t *testing.T) {
	creation := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		expr      string
		now       time.Time
		want      bool
		wantBound time.Time
	}{
		{"first interval not elapsed", "rate(5 minutes)", creation.Add(5 * time.Second), false, time.Time{}},
		{"boundary exactly reached", "rate(5 minutes)", creation.Add(5 * time.Minute), true, creation.Add(5 * time.Minute)},
		{"mid second period", "rate(5 minutes)", creation.Add(7 * time.Minute), true, creation.Add(5 * time.Minute)},
		{"late in third period", "rate(5 minutes)", creation.Add(14*time.Minute + 59*time.Second), true, creation.Add(10 * time.Minute)},
		{"hourly first interval", "rate(1 hour)", creation.Add(30 * time.Minute), false, time.Time{}},
		{"hourly second boundary", "rate(1 hour)", creation.Add(90 * time.Minute), true, creation.Add(time.Hour)},
		{"now before the anchor", "rate(5 minutes)", creation.Add(-time.Hour), false, time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ElapsedExecutionTime(tt.expr, tt.now, creation, nil)
			if ok != tt.want {
				t.Fatalf("ElapsedExecutionTime(%q) elapsed = %v, want %v", tt.expr, ok, tt.want)
			}
			if ok && !got.Equal(tt.wantBound) {
				t.Errorf("ElapsedExecutionTime(%q) = %v, want %v", tt.expr, got, tt.wantBound)
			}
		})
	}

	// A startDate overrides the anchor: with the anchor at 14:00 and now
	// at 12:00 no interval has elapsed, and once the interval after the
	// override has passed its boundary fires.
	future := creation.Add(2 * time.Hour)
	if _, ok := ElapsedExecutionTime("rate(5 minutes)", creation, creation, &future); ok {
		t.Error("rate(5 minutes) with a future startDate reported elapsed")
	}
	bound, ok := ElapsedExecutionTime("rate(5 minutes)", future.Add(5*time.Minute+30*time.Second), creation, &future)
	if !ok || !bound.Equal(future.Add(5*time.Minute)) {
		t.Errorf("rate(5 minutes) anchored at startDate = %v (%v), want %v (true)", bound, ok, future.Add(5*time.Minute))
	}
}

// A rate value of zero or a value too large for an integer parses no
// duration at all; the public resolvers must report "no boundary"
// instead of dividing by a zero duration.
func TestElapsedExecutionTimeRateRejectsDegenerateValues(t *testing.T) {
	creation := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	now := creation.Add(3 * time.Hour)

	for _, expr := range []string{"rate(0 minutes)", "rate(0 hours)", "rate(0 days)", "rate(99999999999999999999 minutes)"} {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%q panicked: %v", expr, r)
				}
			}()
			if _, ok := ElapsedExecutionTime(expr, now, creation, nil); ok {
				t.Errorf("%q reported an elapsed boundary", expr)
			}
			if next, err := NextExecutionTime(expr, now, creation, nil); err == nil && !next.IsZero() {
				t.Errorf("%q reported a next execution time %v", expr, next)
			}
		}()
	}
}

// TestElapsedExecutionTimeAt pins the fixed-time behaviour and the
// interpretation of the offset-less timestamp in now's location.
func TestElapsedExecutionTimeAt(t *testing.T) {
	creation := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	if _, ok := ElapsedExecutionTime("at(2026-01-01T12:00:00)", creation.Add(11*time.Hour+59*time.Minute), creation, nil); ok {
		t.Error("at() before the fixed time reported elapsed")
	}
	bound, ok := ElapsedExecutionTime("at(2026-01-01T12:00:00)", time.Date(2026, 1, 1, 13, 0, 0, 0, time.UTC), creation, nil)
	if !ok || !bound.Equal(time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)) {
		t.Errorf("at() past = %v (%v), want 12:00 UTC (true)", bound, ok)
	}

	// The same instant 00:30 UTC interpreted in a UTC+9 zone: the at()
	// timestamp 09:00 in that zone corresponds to 00:00 UTC and has
	// been reached, while the UTC reading of 09:00 has not.
	tokyo := time.FixedZone("UTC+9", 9*60*60)
	nowInZone := time.Date(2026, 1, 1, 9, 30, 0, 0, tokyo)
	if _, ok := ElapsedExecutionTime("at(2026-01-01T09:00:00)", nowInZone, creation, nil); !ok {
		t.Error("at(09:00) not elapsed at 09:30 in the UTC+9 zone; want elapsed")
	}
	if _, ok := ElapsedExecutionTime("at(2026-01-01T09:00:00)", nowInZone.In(time.UTC), creation, nil); ok {
		t.Error("at(09:00) elapsed at 00:30 UTC; want not elapsed")
	}
}

// TestElapsedExecutionTimeCron pins the latest-elapsed-minute floor: the
// current minute counts when it matches, a late evaluation recovers the
// missed boundary, and boundaries older than the recovery horizon are
// not recovered.
func TestElapsedExecutionTimeCron(t *testing.T) {
	creation := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		expr      string
		now       time.Time
		want      bool
		wantBound time.Time
	}{
		{
			"current matching minute",
			"cron(0/5 * * * ? *)",
			time.Date(2026, 1, 1, 12, 5, 30, 0, time.UTC),
			true,
			time.Date(2026, 1, 1, 12, 5, 0, 0, time.UTC),
		},
		{
			"late evaluation recovers the missed boundary",
			"cron(0/5 * * * ? *)",
			time.Date(2026, 1, 1, 12, 7, 0, 0, time.UTC),
			true,
			time.Date(2026, 1, 1, 12, 5, 0, 0, time.UTC),
		},
		{
			"previous boundary within the horizon",
			"cron(0/5 * * * ? *)",
			time.Date(2026, 1, 1, 12, 4, 0, 0, time.UTC),
			true,
			time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			"manual wildcard path recovers a late boundary",
			"cron(0 9 ? * 6L 2026)",
			time.Date(2026, 1, 31, 12, 0, 0, 0, time.UTC),
			true,
			time.Date(2026, 1, 30, 9, 0, 0, 0, time.UTC),
		},
		{
			"boundary beyond the horizon is not recovered",
			"cron(0 0 1 * ? *)",
			time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
			false,
			time.Time{},
		},
		{
			"boundary within the horizon fires",
			"cron(0 0 1 * ? *)",
			time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC),
			true,
			time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			"no matching minute in the horizon",
			"cron(0 12 1 1 ? 2030)",
			time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			false,
			time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ElapsedExecutionTime(tt.expr, tt.now, creation, nil)
			if ok != tt.want {
				t.Fatalf("ElapsedExecutionTime(%q) elapsed = %v, want %v", tt.expr, ok, tt.want)
			}
			if ok && !got.Equal(tt.wantBound) {
				t.Errorf("ElapsedExecutionTime(%q) = %v, want %v", tt.expr, got, tt.wantBound)
			}
		})
	}

	// A boundary that predates the schedule's creation was never this
	// schedule's to fire: a rule created at 12:01 does not fire for the
	// 12:00 boundary.
	createdLate := time.Date(2026, 1, 1, 12, 1, 0, 0, time.UTC)
	if _, ok := ElapsedExecutionTime("cron(0 12 * * ? *)", createdLate.Add(time.Minute), createdLate, nil); ok {
		t.Error("cron boundary before the creation fired; want not elapsed")
	}

	// The cron floor is evaluated in now's location: 09:30 in a UTC+9
	// zone has the 09:00 zone-local boundary elapsed.
	tokyo := time.FixedZone("UTC+9", 9*60*60)
	nowInZone := time.Date(2026, 1, 1, 9, 30, 0, 0, tokyo)
	bound, ok := ElapsedExecutionTime("cron(0 9 * * ? *)", nowInZone, creation, nil)
	if !ok || !bound.Equal(time.Date(2026, 1, 1, 9, 0, 0, 0, tokyo)) {
		t.Errorf("cron in UTC+9 = %v (%v), want 09:00 zone-local (true)", bound, ok)
	}
}
