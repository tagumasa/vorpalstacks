package cognitoidentity

import (
	"errors"
	"testing"
)

// Deleting a pool cascades across identities (and their ID index entries),
// developer identities, principal tag attribute maps and resource tags, and
// removes the pool record itself.
func TestDeleteIdentityPoolCascadesAllRecords(t *testing.T) {
	s := newTestStore(t)

	pool := NewIdentityPool("cascade-pool", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	identity := NewIdentity(pool.ID)
	identity.Logins = map[string]string{"graph.facebook.com": "token"}
	if err := s.CreateIdentity(identity); err != nil {
		t.Fatalf("create identity: %v", err)
	}
	if err := s.LinkDeveloperIdentity(&DeveloperIdentity{
		DeveloperUserIdentifier: "user-1",
		DeveloperProviderName:   "login.example.com",
		IdentityPoolID:          pool.ID,
		IdentityID:              identity.ID,
	}); err != nil {
		t.Fatalf("link developer identity: %v", err)
	}
	if err := s.SetPrincipalTagAttributeMap(pool.ID, "login.example.com", map[string]string{"email": "email"}, false); err != nil {
		t.Fatalf("set principal tag attribute map: %v", err)
	}
	if err := s.Tag(pool.Arn, map[string]string{"Team": "sdk"}); err != nil {
		t.Fatalf("tag pool: %v", err)
	}

	if _, err := s.GetIdentityByID(identity.ID); err != nil {
		t.Fatalf("identity should resolve before the cascade: %v", err)
	}

	if err := s.DeleteIdentityPool(pool.ID); err != nil {
		t.Fatalf("DeleteIdentityPool: %v", err)
	}

	if _, err := s.GetIdentityPool(pool.ID); !errors.Is(err, ErrIdentityPoolNotFound) {
		t.Fatalf("pool record survived the deletion: %v", err)
	}
	if _, err := s.GetIdentityByID(identity.ID); !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("identity survived the cascade: %v", err)
	}
	if _, err := s.GetIdentity(pool.ID, identity.ID); !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("identity record survived the cascade: %v", err)
	}
	if _, err := s.GetDeveloperIdentity(pool.ID, "login.example.com", "user-1"); !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("developer identity survived the cascade: %v", err)
	}
	if _, err := s.GetPrincipalTagAttributeMap(pool.ID, "login.example.com"); !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("principal tag attribute map survived the cascade: %v", err)
	}
	tags, err := s.List(pool.Arn)
	if err != nil {
		t.Fatalf("list tags after cascade: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("resource tags survived the cascade: %v", tags)
	}
}

// The prefix batch removes only the keys under its prefix; records of other
// pools in the same buckets are untouched.
func TestDeleteIdentityPoolLeavesOtherPoolsIntact(t *testing.T) {
	s := newTestStore(t)

	poolA := NewIdentityPool("cascade-a", false, "us-east-1")
	poolB := NewIdentityPool("cascade-b", false, "us-east-1")
	for _, pool := range []*IdentityPool{poolA, poolB} {
		if _, err := s.CreateIdentityPool(pool); err != nil {
			t.Fatalf("create pool %s: %v", pool.Name, err)
		}
	}
	identityB := NewIdentity(poolB.ID)
	if err := s.CreateIdentity(identityB); err != nil {
		t.Fatalf("create identity B: %v", err)
	}

	if err := s.DeleteIdentityPool(poolA.ID); err != nil {
		t.Fatalf("DeleteIdentityPool: %v", err)
	}

	if _, err := s.GetIdentityByID(identityB.ID); err != nil {
		t.Fatalf("identity of another pool was removed by the cascade: %v", err)
	}
	if _, err := s.GetIdentityPool(poolB.ID); err != nil {
		t.Fatalf("another pool was removed by the cascade: %v", err)
	}
}
