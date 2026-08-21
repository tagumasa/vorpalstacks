package scheduleexpr

import (
	"testing"
	"time"
)

// TestNextExecutionTimeCronDayOfWeek pins the AWS day-of-week numbering
// (1=Sunday..7=Saturday, or SUN-SAT) on every evaluation path: plain
// numerics go through the standard parser after an offset conversion,
// names map directly, and L/W/# expressions are evaluated manually.
func TestNextExecutionTimeCronDayOfWeek(t *testing.T) {
	// 2027-01-01 is a Friday.
	ref := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		expr string
		want time.Time
	}{
		{"numeric 1 is sunday", "cron(0 12 ? * 1 2027)", time.Date(2027, 1, 3, 12, 0, 0, 0, time.UTC)},
		{"numeric 7 is saturday", "cron(0 12 ? * 7 2027)", time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC)},
		{"name sunday", "cron(0 12 ? * SUN 2027)", time.Date(2027, 1, 3, 12, 0, 0, 0, time.UTC)},
		{"range 2-6 is monday to friday", "cron(0 12 ? * 2-6 2027)", time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)},
		{"step 1/2 is sunday tuesday thursday saturday", "cron(0 12 ? * 1/2 2027)", time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC)},
		{"year field gates earlier matches", "cron(0 12 ? * 1 2028)", time.Date(2028, 1, 2, 12, 0, 0, 0, time.UTC)},
		{"last friday of january", "cron(0 12 ? * 6L 2027)", time.Date(2027, 1, 29, 12, 0, 0, 0, time.UTC)},
		{"third friday of january", "cron(0 12 ? * FRI#3 2027)", time.Date(2027, 1, 15, 12, 0, 0, 0, time.UTC)},
		{"bare L is the last day of the week", "cron(0 12 ? * L 2027)", time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextExecutionTime(tt.expr, ref, ref, nil)
			if err != nil {
				t.Fatalf("NextExecutionTime(%q) error: %v", tt.expr, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("NextExecutionTime(%q) = %s, want %s", tt.expr, got, tt.want)
			}
		})
	}
}

// TestNextExecutionTimeCronDom pins the day-of-month wildcards: L (last
// day of the month) and nW (the weekday nearest day n).
func TestNextExecutionTimeCronDom(t *testing.T) {
	tests := []struct {
		name string
		expr string
		ref  time.Time
		want time.Time
	}{
		// 2027-01-01 is a Friday; January has 31 days.
		{"last day of month", "cron(0 0 L * ? *)", time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2027, 1, 31, 0, 0, 0, 0, time.UTC)},
		// 2027-01-01 itself is a Friday, so 1W fires on the same day.
		{"nearest weekday fires on the day itself", "cron(0 12 1W * ? 2027)", time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2027, 1, 1, 12, 0, 0, 0, time.UTC)},
		// 2027-08-01 is a Sunday, so 1W shifts to Monday 2027-08-02.
		{"nearest weekday shifts off sunday", "cron(0 12 1W * ? 2027)", time.Date(2027, 8, 1, 0, 0, 0, 0, time.UTC), time.Date(2027, 8, 2, 12, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextExecutionTime(tt.expr, tt.ref, tt.ref, nil)
			if err != nil {
				t.Fatalf("NextExecutionTime(%q) error: %v", tt.expr, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("NextExecutionTime(%q) = %s, want %s", tt.expr, got, tt.want)
			}
		})
	}
}

// TestNextExecutionTimeStandard pins the previously supported plain
// expressions so the day-of-week conversion cannot regress them.
func TestNextExecutionTimeStandard(t *testing.T) {
	tests := []struct {
		name string
		expr string
		ref  time.Time
		want time.Time
	}{
		{"daily noon", "cron(0 12 * * ? *)", time.Date(2027, 1, 1, 13, 0, 0, 0, time.UTC), time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC)},
		{"minute list", "cron(0,30 8 * * ? *)", time.Date(2027, 1, 1, 8, 15, 0, 0, time.UTC), time.Date(2027, 1, 1, 8, 30, 0, 0, time.UTC)},
		{"dom range", "cron(0 12 1-15 * ? *)", time.Date(2027, 1, 15, 13, 0, 0, 0, time.UTC), time.Date(2027, 2, 1, 12, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextExecutionTime(tt.expr, tt.ref, tt.ref, nil)
			if err != nil {
				t.Fatalf("NextExecutionTime(%q) error: %v", tt.expr, err)
			}
			if !got.Equal(tt.want) {
				t.Errorf("NextExecutionTime(%q) = %s, want %s", tt.expr, got, tt.want)
			}
		})
	}
}

// TestNextExecutionTimeRateAndAt pins the rate anchoring (creation-time
// based boundaries, StartDate override) and the at() one-shot form.
func TestNextExecutionTimeRateAndAt(t *testing.T) {
	creation := time.Date(2027, 1, 1, 10, 0, 0, 0, time.UTC)
	now := time.Date(2027, 1, 1, 10, 12, 0, 0, time.UTC)
	got, err := NextExecutionTime("rate(5 minutes)", now, creation, nil)
	if err != nil {
		t.Fatalf("rate error: %v", err)
	}
	if want := time.Date(2027, 1, 1, 10, 10, 0, 0, time.UTC); !got.Equal(want) {
		t.Errorf("rate(5 minutes) = %s, want %s", got, want)
	}

	start := time.Date(2027, 1, 1, 9, 1, 0, 0, time.UTC)
	got, err = NextExecutionTime("rate(5 minutes)", now, creation, &start)
	if err != nil {
		t.Fatalf("rate with start date error: %v", err)
	}
	if want := time.Date(2027, 1, 1, 10, 11, 0, 0, time.UTC); !got.Equal(want) {
		t.Errorf("rate(5 minutes) with StartDate = %s, want %s", got, want)
	}

	got, err = NextExecutionTime("at(2027-06-01T08:30:00)", now, creation, nil)
	if err != nil {
		t.Fatalf("at error: %v", err)
	}
	if want := time.Date(2027, 6, 1, 8, 30, 0, 0, time.UTC); !got.Equal(want) {
		t.Errorf("at() = %s, want %s", got, want)
	}
}
