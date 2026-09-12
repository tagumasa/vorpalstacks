package cognitoidentityprovider

import (
	"errors"
	"testing"
	"time"
)

// newDomainTestPools creates two distinct pools for domain binding tests.
func newDomainTestPools(t *testing.T, s *CognitoStore) (*UserPool, *UserPool) {
	t.Helper()
	poolA, err := s.CreateUserPool(NewUserPool("domain-pool-a", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool A: %v", err)
	}
	poolB, err := s.CreateUserPool(NewUserPool("domain-pool-b", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool B: %v", err)
	}
	return poolA, poolB
}

func newDomainEntry(domain string, poolID string) *UserPoolDomain {
	return &UserPoolDomain{
		Domain:      domain,
		UserPoolID:  poolID,
		Status:      "ACTIVE",
		CreatedDate: time.Now().UTC(),
	}
}

// A domain string binds to at most one pool, a pool owns at most one domain,
// and a delete requires the owning pool's ID. Every violation has its own
// sentinel, and a missing domain reports the domain sentinel — not the pool
// one.
func TestUserPoolDomainBindingRules(t *testing.T) {
	s := newUserPoolTestStore(t)
	poolA, poolB := newDomainTestPools(t, s)

	if err := s.SetUserPoolDomain("alpha", newDomainEntry("alpha", poolA.ID)); err != nil {
		t.Fatalf("initial bind: %v", err)
	}

	// Claiming pool A's domain for pool B is rejected.
	if err := s.SetUserPoolDomain("alpha", newDomainEntry("alpha", poolB.ID)); !errors.Is(err, ErrUserPoolDomainInUse) {
		t.Fatalf("cross-pool claim returned %v, want ErrUserPoolDomainInUse", err)
	}

	// A second domain string for pool A is rejected too.
	if err := s.SetUserPoolDomain("beta", newDomainEntry("beta", poolA.ID)); !errors.Is(err, ErrUserPoolAlreadyHasDomain) {
		t.Fatalf("second domain for one pool returned %v, want ErrUserPoolAlreadyHasDomain", err)
	}

	// Re-binding a pool's own domain is the update path and succeeds.
	updated := newDomainEntry("alpha", poolA.ID)
	updated.Status = "UPDATING"
	if err := s.SetUserPoolDomain("alpha", updated); err != nil {
		t.Fatalf("same-pool re-bind: %v", err)
	}
	stored, err := s.GetUserPoolDomain("alpha")
	if err != nil || stored.Status != "UPDATING" {
		t.Fatalf("re-bind did not persist: entry=%v err=%v", stored, err)
	}

	// A missing domain reports the domain sentinel, not the pool sentinel.
	if _, err := s.GetUserPoolDomain("no-such-domain"); !errors.Is(err, ErrUserPoolDomainNotFound) {
		t.Fatalf("missing domain returned %v, want ErrUserPoolDomainNotFound", err)
	}
	if _, err := s.GetUserPoolDomainByPool(poolB.ID); !errors.Is(err, ErrUserPoolDomainNotFound) {
		t.Fatalf("domainless pool returned %v, want ErrUserPoolDomainNotFound", err)
	}

	// Deleting with a foreign pool ID fails and leaves the binding intact.
	if err := s.DeleteUserPoolDomain(poolB.ID, "alpha"); !errors.Is(err, ErrUserPoolDomainNotFound) {
		t.Fatalf("foreign delete returned %v, want ErrUserPoolDomainNotFound", err)
	}
	if _, err := s.GetUserPoolDomain("alpha"); err != nil {
		t.Fatalf("foreign delete removed the binding: %v", err)
	}

	// Deleting a missing domain reports not found rather than succeeding.
	if err := s.DeleteUserPoolDomain(poolA.ID, "no-such-domain"); !errors.Is(err, ErrUserPoolDomainNotFound) {
		t.Fatalf("missing delete returned %v, want ErrUserPoolDomainNotFound", err)
	}

	// The owning pool deletes its own domain.
	if err := s.DeleteUserPoolDomain(poolA.ID, "alpha"); err != nil {
		t.Fatalf("owning delete: %v", err)
	}
	if _, err := s.GetUserPoolDomain("alpha"); !errors.Is(err, ErrUserPoolDomainNotFound) {
		t.Fatalf("deleted domain returned %v, want ErrUserPoolDomainNotFound", err)
	}
}

// Domain storage folds case: the production router lowercases the Host
// header before domain extraction, so a domain created in mixed case
// resolves through any spelling of the same name.
func TestUserPoolDomainCaseFolding(t *testing.T) {
	s := newUserPoolTestStore(t)
	poolA, _ := newDomainTestPools(t, s)

	if err := s.SetUserPoolDomain("AuthPool", newDomainEntry("AuthPool", poolA.ID)); err != nil {
		t.Fatalf("set mixed-case domain: %v", err)
	}
	entry, err := s.GetUserPoolDomain("authpool")
	if err != nil {
		t.Fatalf("lowercase lookup of a mixed-case domain failed: %v", err)
	}
	if entry.UserPoolID != poolA.ID {
		t.Fatalf("case-folded lookup returned pool %s", entry.UserPoolID)
	}
	if _, err := s.GetUserPoolDomain("AUTHPOOL"); err != nil {
		t.Fatalf("uppercase lookup of a mixed-case domain failed: %v", err)
	}
}
