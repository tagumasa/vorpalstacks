package cognitoidentity

import (
	"errors"
	"fmt"
	"testing"
)

// GetIdentityByID resolves identities of different pools through the
// identity-ID index, and the index tracks identity deletion.
func TestGetIdentityByIDResolvesAcrossPools(t *testing.T) {
	s := newTestStore(t)

	poolA := NewIdentityPool("pool-a", false, "us-east-1")
	poolB := NewIdentityPool("pool-b", false, "us-east-1")
	for _, pool := range []*IdentityPool{poolA, poolB} {
		if _, err := s.CreateIdentityPool(pool); err != nil {
			t.Fatalf("create pool %s: %v", pool.Name, err)
		}
	}

	identityA := NewIdentity(poolA.ID)
	if err := s.CreateIdentity(identityA); err != nil {
		t.Fatalf("create identity A: %v", err)
	}
	identityB := NewIdentity(poolB.ID)
	if err := s.CreateIdentity(identityB); err != nil {
		t.Fatalf("create identity B: %v", err)
	}

	for _, want := range []*Identity{identityA, identityB} {
		got, err := s.GetIdentityByID(want.ID)
		if err != nil {
			t.Fatalf("GetIdentityByID(%s): %v", want.ID, err)
		}
		if got.ID != want.ID || got.IdentityPoolID != want.IdentityPoolID {
			t.Fatalf("resolved identity {ID: %s, Pool: %s}, want {ID: %s, Pool: %s}",
				got.ID, got.IdentityPoolID, want.ID, want.IdentityPoolID)
		}
	}

	if err := s.DeleteIdentity(poolA.ID, identityA.ID); err != nil {
		t.Fatalf("DeleteIdentity: %v", err)
	}
	if _, err := s.GetIdentityByID(identityA.ID); !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("GetIdentityByID after delete returned %v, want ErrIdentityNotFound", err)
	}

	if _, err := s.GetIdentityByID("us-east-1:00000000-0000-0000-0000-000000000000"); !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("GetIdentityByID for an unknown ID returned %v, want ErrIdentityNotFound", err)
	}
}

// The login index answers the conflict check and the by-logins lookup across
// a pool holding many identities, tracks token refreshes, unlinks, identity
// deletion and the pool cascade, and leaves no residue behind.
func TestLoginIndexConflictAndLifecycle(t *testing.T) {
	s := newTestStore(t)
	created, err := s.CreateIdentityPool(NewIdentityPool("login-index-pool", false, "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	const racers = 200
	for i := 0; i < racers; i++ {
		identity := NewIdentity(created.ID)
		if err := s.CreateIdentity(identity); err != nil {
			t.Fatalf("create identity %d: %v", i, err)
		}
		if err := s.LinkLogins(created.ID, identity.ID, map[string]string{
			"graph.facebook.com": fmt.Sprintf("token-%d", i),
		}); err != nil {
			t.Fatalf("link identity %d: %v", i, err)
		}
	}

	// The conflicting link is detected from the index regardless of how many
	// identities the pool holds, and a distinct token links freely.
	claimant, err := s.GetOrCreateIdentityByLogins(created.ID, map[string]string{"graph.facebook.com": "token-7"})
	if err != nil {
		t.Fatalf("resolve claimant: %v", err)
	}
	other := NewIdentity(created.ID)
	if err := s.CreateIdentity(other); err != nil {
		t.Fatalf("create other: %v", err)
	}
	if err := s.LinkLogins(created.ID, other.ID, map[string]string{"graph.facebook.com": "token-7"}); !errors.Is(err, ErrLoginConflict) {
		t.Fatalf("conflicting link returned %v, want ErrLoginConflict", err)
	}
	if err := s.LinkLogins(created.ID, other.ID, map[string]string{"graph.facebook.com": "token-fresh"}); err != nil {
		t.Fatalf("link with a distinct token: %v", err)
	}

	// A token refresh re-points the index: the old token links elsewhere,
	// the new one still resolves to its owner.
	if err := s.LinkLogins(created.ID, other.ID, map[string]string{"graph.facebook.com": "token-refreshed"}); err != nil {
		t.Fatalf("refresh token: %v", err)
	}
	if err := s.LinkLogins(created.ID, claimant.ID, map[string]string{"graph.facebook.com": "token-fresh"}); err != nil {
		t.Fatalf("released token no longer links: %v", err)
	}
	resolved, err := s.GetOrCreateIdentityByLogins(created.ID, map[string]string{"graph.facebook.com": "token-refreshed"})
	if err != nil {
		t.Fatalf("resolve refreshed token: %v", err)
	}
	if resolved.ID != other.ID {
		t.Fatalf("refreshed token resolved %s, want %s", resolved.ID, other.ID)
	}

	// Unlink releases the claim; deletion releases the identity's remaining
	// entries; the pool cascade leaves no index residue at all.
	if err := s.UnlinkLogins(created.ID, other.ID, []string{"graph.facebook.com"}); err != nil {
		t.Fatalf("unlink: %v", err)
	}
	if err := s.LinkLogins(created.ID, claimant.ID, map[string]string{"graph.facebook.com": "token-refreshed"}); err != nil {
		t.Fatalf("link a released token: %v", err)
	}
	if err := s.DeleteIdentity(created.ID, claimant.ID); err != nil {
		t.Fatalf("delete identity: %v", err)
	}
	// The deleted identity's login resolves to a freshly created identity,
	// never to the removed record.
	reborn, err := s.GetOrCreateIdentityByLogins(created.ID, map[string]string{"graph.facebook.com": "token-7"})
	if err != nil {
		t.Fatalf("resolve after deletion: %v", err)
	}
	if reborn.ID == claimant.ID {
		t.Fatal("deleted identity's login still resolves to the removed record")
	}
	if err := s.DeleteIdentityPool(created.ID); err != nil {
		t.Fatalf("delete pool: %v", err)
	}
	residue := 0
	if err := s.identitiesStore.ScanPrefix(loginIndexKeyPrefix, func(string, []byte) error {
		residue++
		return nil
	}); err != nil {
		t.Fatalf("scan index residue: %v", err)
	}
	if residue != 0 {
		t.Fatalf("pool cascade left %d login index entries behind", residue)
	}
}

// A developer-identity merge that discards a token (the destination keeps
// its own for the provider) releases the token's index entry in the same
// commit — the entry must not outlive the source identity it names.
func TestMergeDeveloperIdentitiesReleasesDiscardedTokens(t *testing.T) {
	s := newTestStore(t)
	created, err := s.CreateIdentityPool(NewIdentityPool("merge-index-pool", false, "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	const provider = "graph.facebook.com"
	newLinkedIdentity := func(token string) *Identity {
		t.Helper()
		identity := NewIdentity(created.ID)
		identity.Logins = map[string]string{provider: token}
		if err := s.CreateIdentity(identity); err != nil {
			t.Fatalf("create identity %s: %v", token, err)
		}
		return identity
	}
	source := newLinkedIdentity("merge-src-token")
	dest := newLinkedIdentity("merge-dst-token")
	for _, link := range []struct{ user, identityID string }{
		{"user-a", source.ID}, {"user-b", dest.ID},
	} {
		if _, err := s.EnsureDeveloperIdentity(created.ID, "dev-provider", link.user, link.identityID); err != nil {
			t.Fatalf("link developer identity %s: %v", link.user, err)
		}
	}

	if _, err := s.MergeDeveloperIdentities(created.ID, "dev-provider", "user-a", "user-b"); err != nil {
		t.Fatalf("merge: %v", err)
	}

	// The discarded token is freed with its owner: a third identity links it
	// without the conflict the stale entry used to answer.
	third := NewIdentity(created.ID)
	if err := s.CreateIdentity(third); err != nil {
		t.Fatalf("create third identity: %v", err)
	}
	if err := s.LinkLogins(created.ID, third.ID, map[string]string{provider: "merge-src-token"}); err != nil {
		t.Fatalf("link the discarded token after merge: %v", err)
	}
	// The destination's own token still resolves to the destination.
	resolved, err := s.GetOrCreateIdentityByLogins(created.ID, map[string]string{provider: "merge-dst-token"})
	if err != nil {
		t.Fatalf("resolve destination token: %v", err)
	}
	if resolved.ID != dest.ID {
		t.Fatalf("destination token resolved %s, want %s", resolved.ID, dest.ID)
	}
}

// The by-logins lookup distinguishes a stale index entry (identity gone →
// no match, the create path self-heals) from a failing record read, which
// propagates instead of letting GetOrCreateIdentityByLogins duplicate the
// identity.
func TestFindIdentityByLoginsPropagatesReadFailures(t *testing.T) {
	s := newTestStore(t)
	created, err := s.CreateIdentityPool(NewIdentityPool("read-fail-pool", false, "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	const provider = "graph.facebook.com"
	identity := NewIdentity(created.ID)
	identity.Logins = map[string]string{provider: "corrupt-token"}
	if err := s.CreateIdentity(identity); err != nil {
		t.Fatalf("create identity: %v", err)
	}
	logins := map[string]string{provider: "corrupt-token"}

	// An unreadable record behind the index is a store failure, not a miss.
	bucket := s.identitiesStore.Bucket()
	if err := bucket.Put([]byte(IdentityPoolIdentityKey(created.ID, identity.ID)), []byte("not-json")); err != nil {
		t.Fatalf("corrupt record: %v", err)
	}
	if _, err := s.findIdentityByLogins(created.ID, logins); err == nil || errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("unreadable record returned %v, want a propagated store error", err)
	}
	if _, err := s.GetOrCreateIdentityByLogins(created.ID, logins); err == nil || errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("GetOrCreate with an unreadable record returned %v, want a propagated store error", err)
	}

	// A stale entry (record deleted) stays a miss and self-heals by creating.
	if err := bucket.Delete([]byte(IdentityPoolIdentityKey(created.ID, identity.ID))); err != nil {
		t.Fatalf("remove record: %v", err)
	}
	reborn, err := s.GetOrCreateIdentityByLogins(created.ID, logins)
	if err != nil {
		t.Fatalf("resolve after record removal: %v", err)
	}
	if reborn.ID == identity.ID {
		t.Fatal("stale index entry resolved to the removed record")
	}
}
