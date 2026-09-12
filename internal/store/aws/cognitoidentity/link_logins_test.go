package cognitoidentity

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

// LinkLogins attaches public-provider entries to the named identity and is
// idempotent for pairs the identity already carries.
func TestLinkLoginsAttachesEntries(t *testing.T) {
	s := newTestStore(t)
	pool := NewIdentityPool("link-pool", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	identity := NewIdentity(pool.ID)
	if err := s.CreateIdentity(identity); err != nil {
		t.Fatalf("create identity: %v", err)
	}

	logins := map[string]string{"graph.facebook.com": "token-1"}
	if err := s.LinkLogins(pool.ID, identity.ID, logins); err != nil {
		t.Fatalf("link logins: %v", err)
	}
	got, err := s.GetIdentity(pool.ID, identity.ID)
	if err != nil {
		t.Fatalf("get identity: %v", err)
	}
	if got.Logins["graph.facebook.com"] != "token-1" {
		t.Fatalf("login was not attached: %v", got.Logins)
	}

	if err := s.LinkLogins(pool.ID, identity.ID, logins); err != nil {
		t.Fatalf("re-linking the same pair failed: %v", err)
	}
}

// A login already linked to another identity of the pool is rejected with
// ErrLoginConflict; the model documents ResourceConflictException for this
// case ("a login which is already linked to another account").
func TestLinkLoginsRejectsLoginLinkedToAnotherIdentity(t *testing.T) {
	s := newTestStore(t)
	pool := NewIdentityPool("link-conflict-pool", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	first := NewIdentity(pool.ID)
	if err := s.CreateIdentity(first); err != nil {
		t.Fatalf("create first identity: %v", err)
	}
	second := NewIdentity(pool.ID)
	if err := s.CreateIdentity(second); err != nil {
		t.Fatalf("create second identity: %v", err)
	}

	if err := s.LinkLogins(pool.ID, first.ID, map[string]string{"graph.facebook.com": "token-1"}); err != nil {
		t.Fatalf("link first: %v", err)
	}
	err := s.LinkLogins(pool.ID, second.ID, map[string]string{"graph.facebook.com": "token-1"})
	if !errors.Is(err, ErrLoginConflict) {
		t.Fatalf("linking a login owned by another identity returned %v, want ErrLoginConflict", err)
	}

	// A different token for the same provider is a distinct login and links.
	if err := s.LinkLogins(pool.ID, second.ID, map[string]string{"graph.facebook.com": "token-2"}); err != nil {
		t.Fatalf("link second with a distinct token: %v", err)
	}
}

// Concurrent merges for one identity cannot lose each other's links: the
// read-merge-write runs under the pool lock.
func TestMergeLoginsConcurrentMergesAllLand(t *testing.T) {
	s := newTestStore(t)
	pool := NewIdentityPool("merge-logins-pool", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	identity := NewIdentity(pool.ID)
	if err := s.CreateIdentity(identity); err != nil {
		t.Fatalf("create identity: %v", err)
	}

	const racers = 8
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < racers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, _ = s.MergeLogins(pool.ID, identity.ID, map[string]string{
				fmt.Sprintf("provider-%d.example.com", i): fmt.Sprintf("token-%d", i),
			})
		}(i)
	}
	close(start)
	wg.Wait()

	got, err := s.GetIdentity(pool.ID, identity.ID)
	if err != nil {
		t.Fatalf("get identity: %v", err)
	}
	for i := 0; i < racers; i++ {
		provider := fmt.Sprintf("provider-%d.example.com", i)
		if got.Logins[provider] != fmt.Sprintf("token-%d", i) {
			t.Fatalf("login %s missing after concurrent merges: %v", provider, got.Logins)
		}
	}
}

// MergeLogins returns the merged Logins map read inside the pool lock: both
// the pre-existing link and the merged link are present in the return value.
func TestMergeLoginsReturnsMergedMap(t *testing.T) {
	s := newTestStore(t)
	created, err := s.CreateIdentityPool(NewIdentityPool("merge-return-pool", false, "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	identity := NewIdentity(created.ID)
	if err := s.CreateIdentity(identity); err != nil {
		t.Fatalf("create identity: %v", err)
	}
	if err := s.LinkLogins(created.ID, identity.ID, map[string]string{"a.example.com": "token-a"}); err != nil {
		t.Fatalf("link first login: %v", err)
	}

	merged, err := s.MergeLogins(created.ID, identity.ID, map[string]string{"b.example.com": "token-b"})
	if err != nil {
		t.Fatalf("merge logins: %v", err)
	}
	if len(merged) != 2 || merged["a.example.com"] != "token-a" || merged["b.example.com"] != "token-b" {
		t.Fatalf("merged map = %#v, want both providers", merged)
	}
}
