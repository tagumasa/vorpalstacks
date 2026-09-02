package wafv2

import (
	"net/netip"
	"sync"
	"time"

	"vorpalstacks/internal/services/aws/wafv2/inspection"
)

// rateTracker implements inspection.RateTracker with an in-memory
// sliding window. AWS WAF evaluates rates roughly every 10 seconds and
// does not expose counter state through the API beyond
// GetRateBasedStatementManagedKeys, so in-memory bookkeeping matches
// the observable behaviour; a platform restart resets the counters,
// which the managed-keys API reflects as an empty key set.
type rateTracker struct {
	mu        sync.Mutex
	bucket    time.Duration
	keys      map[inspection.RateKey][]time.Time
	lastSweep time.Time
}

// newRateTracker creates a tracker with a 1-second bucket granularity,
// which keeps the sliding-window error well under the ~10-second rate
// evaluation cadence while bounding per-key bookkeeping.
func newRateTracker() *rateTracker {
	return &rateTracker{
		bucket: time.Second,
		keys:   make(map[inspection.RateKey][]time.Time),
	}
}

// rateTrackerRetention bounds how long a key's entry survives without
// activity: the longest evaluation window a rate-based statement can
// configure is 600 seconds, so a key whose newest timestamp fell
// outside that horizon can never count again and its entry is swept.
const rateTrackerRetention = 10 * time.Minute

// rateTrackerSweepInterval spaces the amortised retention sweeps so
// the per-request cost of pruning stays negligible.
const rateTrackerSweepInterval = time.Minute

// Hit increments the counter for the key and returns the number of
// requests recorded for the key inside the window ending at now.
func (t *rateTracker) Hit(key inspection.RateKey, window time.Duration, now time.Time) int64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := now.Add(-window)
	times := t.keys[key]
	kept := times[:0]
	for _, ts := range times {
		if ts.After(cutoff) {
			kept = append(kept, ts)
		}
	}
	kept = append(kept, now)
	t.keys[key] = kept
	if now.Sub(t.lastSweep) >= rateTrackerSweepInterval {
		t.sweepLocked(now)
		t.lastSweep = now
	}
	return int64(len(kept))
}

// sweepLocked drops entries whose newest timestamp lies outside the
// retention horizon, bounding the tracker's memory to keys that were
// active within the longest possible evaluation window. Without the
// sweep, address rotation would grow the key set without limit.
func (t *rateTracker) sweepLocked(now time.Time) {
	cutoff := now.Add(-rateTrackerRetention)
	for key, times := range t.keys {
		newest := time.Time{}
		for _, ts := range times {
			if ts.After(newest) {
				newest = ts
			}
		}
		if newest.Before(cutoff) {
			delete(t.keys, key)
		}
	}
}

// ActiveIPKeys returns the aggregation keys currently tracked for the
// given web ACL rule that are valid IP addresses, classified by IP
// version. This backs GetRateBasedStatementManagedKeys, whose response
// shape (ManagedKeysIPV4 / ManagedKeysIPV6 with IPAddressVersion and
// Addresses) only carries address keys.
func (t *rateTracker) ActiveIPKeys(webACLARN, ruleName string) (ipv4, ipv6 []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for key := range t.keys {
		if key.WebACLARN != webACLARN || key.RuleName != ruleName {
			continue
		}
		addr, err := netip.ParseAddr(key.Value)
		if err != nil {
			continue
		}
		if addr.Is4() || addr.Is4In6() {
			ipv4 = append(ipv4, key.Value)
		} else {
			ipv6 = append(ipv6, key.Value)
		}
	}
	return ipv4, ipv6
}
