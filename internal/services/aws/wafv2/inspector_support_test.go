package wafv2

import (
	"testing"
	"time"

	"vorpalstacks/internal/services/aws/wafv2/inspection"
)

func TestRateTrackerWindowSlides(t *testing.T) {
	tracker := newRateTracker()
	key := inspection.RateKey{WebACLARN: "arn", RuleName: "r", Value: "203.0.113.9"}
	now := time.Now()

	if got := tracker.Hit(key, time.Minute, now); got != 1 {
		t.Fatalf("first hit = %d, want 1", got)
	}
	if got := tracker.Hit(key, time.Minute, now.Add(30*time.Second)); got != 2 {
		t.Fatalf("second hit = %d, want 2", got)
	}
	// A hit after the first two left the one-minute window counts fresh.
	if got := tracker.Hit(key, time.Minute, now.Add(2*time.Minute)); got != 1 {
		t.Fatalf("post-window hit = %d, want 1", got)
	}

	// Keys aggregate independently.
	other := inspection.RateKey{WebACLARN: "arn", RuleName: "r", Value: "198.51.100.7"}
	if got := tracker.Hit(other, time.Minute, now.Add(2*time.Minute)); got != 1 {
		t.Fatalf("other-key hit = %d, want 1", got)
	}
}

func TestRateTrackerActiveIPKeys(t *testing.T) {
	tracker := newRateTracker()
	now := time.Now()
	tracker.Hit(inspection.RateKey{WebACLARN: "arn", RuleName: "rate", Value: "203.0.113.9"}, time.Minute, now)
	tracker.Hit(inspection.RateKey{WebACLARN: "arn", RuleName: "rate", Value: "2001:db8::1"}, time.Minute, now)
	tracker.Hit(inspection.RateKey{WebACLARN: "arn", RuleName: "rate", Value: "not-an-ip"}, time.Minute, now)
	tracker.Hit(inspection.RateKey{WebACLARN: "arn", RuleName: "other", Value: "192.0.2.1"}, time.Minute, now)

	ipv4, ipv6 := tracker.ActiveIPKeys("arn", "rate")
	if len(ipv4) != 1 || ipv4[0] != "203.0.113.9" {
		t.Fatalf("ipv4 keys = %v", ipv4)
	}
	if len(ipv6) != 1 || ipv6[0] != "2001:db8::1" {
		t.Fatalf("ipv6 keys = %v", ipv6)
	}
}

func TestRateTrackerSweepsIdleKeys(t *testing.T) {
	tracker := newRateTracker()
	now := time.Now()
	idle := inspection.RateKey{WebACLARN: "arn", RuleName: "rate", Value: "203.0.113.9"}
	active := inspection.RateKey{WebACLARN: "arn", RuleName: "rate", Value: "198.51.100.7"}
	tracker.Hit(idle, time.Minute, now)
	tracker.Hit(active, time.Minute, now)

	// A hit far in the future triggers the retention sweep: the key
	// idle past the longest evaluation window disappears, the recently
	// active key stays.
	later := now.Add(rateTrackerRetention + 2*time.Minute)
	tracker.Hit(active, time.Minute, later)
	tracker.mu.Lock()
	if _, ok := tracker.keys[idle]; ok {
		t.Fatal("idle key survived the retention sweep")
	}
	if _, ok := tracker.keys[active]; !ok {
		t.Fatal("active key was swept")
	}
	tracker.mu.Unlock()
}
