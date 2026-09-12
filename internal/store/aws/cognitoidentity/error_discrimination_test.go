package cognitoidentity

import (
	"errors"
	"testing"
)

// A corrupt or unreadable identity-index entry must surface as the storage
// error itself, not as the not-found sentinel: only a genuinely missing key
// maps to ErrIdentityNotFound.
func TestGetIdentityByIDCorruptIndexIsNotNotFound(t *testing.T) {
	s := newTestStore(t)
	pool := NewIdentityPool("corrupt-index-pool", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	identity := NewIdentity(pool.ID)
	if err := s.CreateIdentity(identity); err != nil {
		t.Fatalf("create identity: %v", err)
	}

	if err := s.identitiesStore.PutRaw(identityIndexKey(identity.ID), []byte("not-json")); err != nil {
		t.Fatalf("corrupt the index entry: %v", err)
	}
	_, err := s.GetIdentityByID(identity.ID)
	if err == nil {
		t.Fatal("lookup through a corrupt index entry unexpectedly succeeded")
	}
	if errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("corrupt index entry reported as not-found: %v", err)
	}
}

// A corrupt or unreadable identity record must surface as the storage error
// itself during UnlinkLogins, not as the not-found sentinel.
func TestUnlinkLoginsCorruptRecordIsNotNotFound(t *testing.T) {
	s := newTestStore(t)
	pool := NewIdentityPool("corrupt-record-pool", false, "us-east-1")
	if _, err := s.CreateIdentityPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	identity := NewIdentity(pool.ID)
	if err := s.CreateIdentity(identity); err != nil {
		t.Fatalf("create identity: %v", err)
	}

	if err := s.identitiesStore.PutRaw(IdentityPoolIdentityKey(pool.ID, identity.ID), []byte("not-json")); err != nil {
		t.Fatalf("corrupt the identity record: %v", err)
	}
	err := s.UnlinkLogins(pool.ID, identity.ID, []string{"graph.facebook.com"})
	if err == nil {
		t.Fatal("unlink through a corrupt record unexpectedly succeeded")
	}
	if errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("corrupt identity record reported as not-found: %v", err)
	}
}
