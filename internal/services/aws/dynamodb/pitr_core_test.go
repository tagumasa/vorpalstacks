package dynamodb

import (
	"testing"
	"time"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestPitrEarliestRestorableTrailingWindow(t *testing.T) {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

	// Recovery enabled an hour ago: the journal only reaches back to the
	// enable moment, so that remains the earliest restorable time.
	fresh := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.Add(-time.Hour),
	}
	if got := pitrEarliestRestorable(fresh, now); !got.Equal(fresh.EarliestRestorableDateTime) {
		t.Fatalf("fresh enable: earliest = %v, want the enable time %v", got, fresh.EarliestRestorableDateTime)
	}

	// Recovery enabled 40 days ago with the default period: the trailing
	// 35-day edge is later than the enable moment and wins.
	aged := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.AddDate(0, 0, -40),
	}
	want := now.AddDate(0, 0, -pitrDefaultRecoveryPeriodDays)
	if got := pitrEarliestRestorable(aged, now); !got.Equal(want) {
		t.Fatalf("aged enable: earliest = %v, want the trailing edge %v", got, want)
	}

	// A configured 1-day period on a 3-day-old enable moves the edge to
	// now-1d.
	shortPeriod := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.AddDate(0, 0, -3),
		RecoveryPeriodInDays:       1,
	}
	want = now.AddDate(0, 0, -1)
	if got := pitrEarliestRestorable(shortPeriod, now); !got.Equal(want) {
		t.Fatalf("1-day period: earliest = %v, want %v", got, want)
	}

	// The trailing edge never precedes the enable moment: the journal holds
	// nothing older than it.
	shortEnable := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.Add(-12 * time.Hour),
		RecoveryPeriodInDays:       1,
	}
	if got := pitrEarliestRestorable(shortEnable, now); !got.Equal(shortEnable.EarliestRestorableDateTime) {
		t.Fatalf("recent enable with 1-day period: earliest = %v, want the enable time %v", got, shortEnable.EarliestRestorableDateTime)
	}
}
