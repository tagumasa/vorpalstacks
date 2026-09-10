package apigateway

import (
	"testing"

	"vorpalstacks/internal/core/storage"
)

// TestRecordUsagePropagatesReadFailure pins the error discipline of
// RecordUsage: only a missing record resets the counter to 1 — a record
// that cannot be read (here: corrupt payload) must surface the error
// instead of silently destroying the accumulated count.
func TestRecordUsagePropagatesReadFailure(t *testing.T) {
	ps, err := storage.NewPebbleStorage(&storage.Config{
		Path:       t.TempDir(),
		TTLEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewPebbleStorage: %v", err)
	}
	t.Cleanup(func() { ps.Close() })

	store := NewUsageStore(ps, "123456789012", "us-east-1")

	// Repeated recording accumulates from the stored count.
	first := &UsageRecord{UsagePlanID: "plan", APIKeyID: "key", Date: "2026-09-09"}
	if err := store.RecordUsage(first); err != nil {
		t.Fatalf("RecordUsage: %v", err)
	}
	if first.RequestCount != 1 {
		t.Fatalf("expected first count 1, got %d", first.RequestCount)
	}
	second := &UsageRecord{UsagePlanID: "plan", APIKeyID: "key", Date: "2026-09-09"}
	if err := store.RecordUsage(second); err != nil {
		t.Fatalf("RecordUsage again: %v", err)
	}
	if second.RequestCount != 2 {
		t.Fatalf("expected accumulated count 2, got %d", second.RequestCount)
	}

	// A fresh key starts at 1.
	fresh := &UsageRecord{UsagePlanID: "plan", APIKeyID: "other", Date: "2026-09-09"}
	if err := store.RecordUsage(fresh); err != nil {
		t.Fatalf("RecordUsage fresh: %v", err)
	}
	if fresh.RequestCount != 1 {
		t.Fatalf("expected fresh count 1, got %d", fresh.RequestCount)
	}

	// A record that cannot be decoded propagates the failure.
	if err := store.Put("usage#plan#corrupt#2026-09-09", "not-a-usage-record"); err != nil {
		t.Fatalf("seed corrupt record: %v", err)
	}
	if err := store.RecordUsage(&UsageRecord{UsagePlanID: "plan", APIKeyID: "corrupt", Date: "2026-09-09"}); err == nil {
		t.Fatal("RecordUsage swallowed a read failure and reset the counter")
	}
}

// TestActiveQuotaExtensionPeriods pins the lapse windows of a remaining-
// quota grant: DAY covers the grant date only, WEEK the seven days starting
// at it, MONTH the grant date's calendar month — before the grant and after
// the window the override contributes nothing.
func TestActiveQuotaExtensionPeriods(t *testing.T) {
	ps, err := storage.NewPebbleStorage(&storage.Config{
		Path:       t.TempDir(),
		TTLEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewPebbleStorage: %v", err)
	}
	t.Cleanup(func() { ps.Close() })

	store := NewUsageStore(ps, "123456789012", "us-east-1")

	cases := []struct {
		period string
		date   string
		want   int64
	}{
		{"DAY", "2026-08-31", 0},
		{"DAY", "2026-09-01", 50},
		{"DAY", "2026-09-02", 0},
		{"WEEK", "2026-09-01", 50},
		{"WEEK", "2026-09-07", 50},
		{"WEEK", "2026-09-08", 0},
		{"MONTH", "2026-09-30", 50},
		{"MONTH", "2026-10-01", 0},
	}
	for _, c := range cases {
		if err := store.SetUsageQuotaExtension(&UsageQuotaExtension{
			UsagePlanID: "plan", APIKeyID: "key", Period: c.period, GrantedDate: "2026-09-01", Net: 50,
		}); err != nil {
			t.Fatalf("set %s: %v", c.period, err)
		}
		if got := store.ActiveQuotaExtension("plan", "key", c.date); got != c.want {
			t.Errorf("ActiveQuotaExtension(%s, %s) = %d, want %d", c.period, c.date, got, c.want)
		}
	}

	// No grant on another key reads as zero.
	if got := store.ActiveQuotaExtension("plan", "other", "2026-09-01"); got != 0 {
		t.Errorf("ActiveQuotaExtension(unknown key) = %d, want 0", got)
	}
}
