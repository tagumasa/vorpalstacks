package cognitoidentity

import (
	"errors"
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
