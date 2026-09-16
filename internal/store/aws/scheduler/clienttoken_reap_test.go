package scheduler

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// seedExpiredEntry writes a token entry whose CreatedAt lies beyond
// clientTokenTTL, straight through the backing bucket: the reaper's delete
// path only runs when an expired record exists, and no test can wait out
// the 24-hour TTL to produce one.
func seedExpiredEntry(t *testing.T, tokens *ClientTokenStore, token, resourceArn string) {
	t.Helper()
	key := clientTokenKey("schedule", resourceArn, token)
	entry := ClientTokenEntry{
		ResourceArn: resourceArn,
		CreatedAt:   time.Now().Add(-clientTokenTTL - time.Hour),
	}
	if err := tokens.store.PutRaw(key, mustMarshal(&entry)); err != nil {
		t.Fatalf("seed expired entry: %v", err)
	}
}

// TestReapExpiredKeepsFreshClaim pins the reaper's selectivity: an entry
// within its TTL survives reapExpired while an expired one is deleted, so a
// client retry still replays the first outcome instead of re-executing the
// operation.
func TestReapExpiredKeepsFreshClaim(t *testing.T) {
	store := newCompletionTestStore(t)
	tokens := store.ClientTokens()

	expiredArn := "arn:aws:scheduler:us-east-1:000000000000:schedule/default/reap-expired"
	seedExpiredEntry(t, tokens, "reap-expired", expiredArn)

	if _, claimed := tokens.LookupOrClaim("reap-pin", "arn:aws:scheduler:us-east-1:000000000000:schedule/default/reap-pin", "schedule"); !claimed {
		t.Fatal("the first LookupOrClaim must claim")
	}
	tokens.reapExpired()
	if _, claimed := tokens.LookupOrClaim("reap-pin", "arn:aws:scheduler:us-east-1:000000000000:schedule/default/reap-pin", "schedule"); claimed {
		t.Fatal("a fresh claim did not survive reapExpired — the reaper deleted it and the retry re-executed")
	}
	if data, err := tokens.store.GetRaw(clientTokenKey("schedule", expiredArn, "reap-expired")); err != nil || data != nil {
		t.Fatalf("the expired entry survived reapExpired (data=%v, err=%v)", data, err)
	}
}

// TestStopIsIdempotent pins Stop's contract: the store cache is closed
// without eviction, so a repeated Stop (a retried shutdown hook) must not
// close an already-closed channel and panic.
func TestStopIsIdempotent(t *testing.T) {
	store := newCompletionTestStore(t)
	store.Close()
	store.Close()
}

// TestReapExpiredConcurrentWithClaims runs the reaper, primed with expired
// entries, against a concurrent claim storm on one key. It pins the
// observable reaper contracts under that load: the post-storm entry is a
// live claim (a retried LookupOrClaim answers "already claimed"), and the
// seeded expired entries are consumed by the delete path during the storm,
// not only after it. The read-to-delete window inside a single scan
// iteration — the interleaving the mutex exists to close — is narrower than
// any black-box interleave a test can force (a mid-storm loss is re-claimed
// by the next iteration before it could be observed); the mutex remains the
// mechanism guaranteeing it, as its own comment records.
func TestReapExpiredConcurrentWithClaims(t *testing.T) {
	store := newCompletionTestStore(t)
	tokens := store.ClientTokens()
	const resourceArn = "arn:aws:scheduler:us-east-1:000000000000:schedule/default/reap-race"
	seedExpiredEntry(t, tokens, "reap-race", resourceArn)
	fillerArns := make([]string, 8)
	for i := range fillerArns {
		fillerArns[i] = fmt.Sprintf("arn:aws:scheduler:us-east-1:000000000000:schedule/default/reap-race-fill-%d", i)
		seedExpiredEntry(t, tokens, fmt.Sprintf("reap-race-fill-%d", i), fillerArns[i])
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				tokens.LookupOrClaim("reap-race", resourceArn, "schedule")
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 200; j++ {
			tokens.reapExpired()
		}
	}()
	wg.Wait()

	if _, claimed := tokens.LookupOrClaim("reap-race", resourceArn, "schedule"); claimed {
		t.Fatal("after the race the entry must be a live claim; the reaper deleted a claim a retry had just written")
	}
	// Every seeded expired entry is gone: the delete path ran through the
	// storm, not only after it.
	for i, arn := range fillerArns {
		if data, err := tokens.store.GetRaw(clientTokenKey("schedule", arn, fmt.Sprintf("reap-race-fill-%d", i))); err != nil || data != nil {
			t.Fatalf("seeded expired entry %d survived the reaper storm (data=%v, err=%v)", i, data, err)
		}
	}
}
