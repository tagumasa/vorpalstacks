package cognitoidentityprovider

import (
	"errors"
	"testing"
)

func newAliasClaimUser(poolID, username, email string, verified bool) *User {
	u := NewUser(poolID, username)
	u.Attributes = map[string]string{"email": email}
	if verified {
		u.Attributes["email_verified"] = "true"
	}
	return u
}

// A verified email configured as an alias is exclusive to its holder:
// another user claiming it is rejected with ErrAliasExists, the migration
// variant takes the claim over (un-verifying the previous holder), and
// GetUser resolves the alias value to the claiming username.
func TestAliasClaimLifecycle(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool := NewUserPool("alias-pool", "us-east-1")
	pool.AliasAttributes = []string{"email"}
	if _, err := s.CreateUserPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}

	holder := newAliasClaimUser(pool.ID, "holder", "shared@example.com", true)
	if err := s.CreateUser(holder); err != nil {
		t.Fatalf("create holder: %v", err)
	}
	if resolved, err := s.GetUser(pool.ID, "shared@example.com"); err != nil || resolved.Username != "holder" {
		t.Fatalf("alias lookup = (%v, %v), want holder", resolved, err)
	}

	// An unverified duplicate coexists: the claim activates only on
	// verification.
	unverified := newAliasClaimUser(pool.ID, "unverified", "shared@example.com", false)
	if err := s.CreateUser(unverified); err != nil {
		t.Fatalf("create unverified duplicate: %v", err)
	}

	// Activating the duplicate's claim is the conflict.
	unverified.Attributes["email_verified"] = "true"
	if err := s.UpdateUser(unverified); !errors.Is(err, ErrAliasExists) {
		t.Fatalf("verifying a duplicate returned %v, want ErrAliasExists", err)
	}
	// The rejected write left no partial state: the index still resolves to
	// the holder.
	if resolved, err := s.GetUser(pool.ID, "shared@example.com"); err != nil || resolved.Username != "holder" {
		t.Fatalf("alias lookup after rejected claim = (%v, %v), want holder", resolved, err)
	}

	// The migration variant takes the claim over and un-verifies the
	// previous holder, who keeps the attribute value but can no longer sign
	// in with it.
	migrated := newAliasClaimUser(pool.ID, "migrated", "shared@example.com", true)
	if err := s.CreateUserMigrateAliasClaims(migrated); err != nil {
		t.Fatalf("migrating create: %v", err)
	}
	previous, err := s.GetUser(pool.ID, "holder")
	if err != nil {
		t.Fatal(err)
	}
	if previous.Attributes["email"] != "shared@example.com" || previous.Attributes["email_verified"] == "true" {
		t.Fatalf("previous holder attributes after migration = %v", previous.Attributes)
	}
	if resolved, err := s.GetUser(pool.ID, "shared@example.com"); err != nil || resolved.Username != "migrated" {
		t.Fatalf("alias lookup after migration = (%v, %v), want migrated", resolved, err)
	}

	// Changing the claim releases the old value.
	migrated.Attributes["email"] = "moved@example.com"
	if err := s.UpdateUser(migrated); err != nil {
		t.Fatalf("update claim: %v", err)
	}
	if _, err := s.GetUser(pool.ID, "shared@example.com"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("released alias still resolves: %v", err)
	}
	if resolved, err := s.GetUser(pool.ID, "moved@example.com"); err != nil || resolved.Username != "migrated" {
		t.Fatalf("new alias value lookup = (%v, %v), want migrated", resolved, err)
	}

	// Deleting the holder releases the claim for reassignment.
	if err := s.DeleteUser(pool.ID, "migrated"); err != nil {
		t.Fatalf("delete migrated: %v", err)
	}
	if _, err := s.GetUser(pool.ID, "moved@example.com"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("alias of a deleted user still resolves: %v", err)
	}
	reborn := newAliasClaimUser(pool.ID, "reborn", "moved@example.com", true)
	if err := s.CreateUser(reborn); err != nil {
		t.Fatalf("claim a released value: %v", err)
	}
}

// preferred_username has no verification round and claims on presence, so
// the second write of the same value is rejected outright.
func TestPreferredUsernameClaimOnPresence(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool := NewUserPool("preferred-pool", "us-east-1")
	pool.AliasAttributes = []string{"preferred_username"}
	if _, err := s.CreateUserPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}

	first := NewUser(pool.ID, "first")
	first.Attributes = map[string]string{"preferred_username": "nick"}
	if err := s.CreateUser(first); err != nil {
		t.Fatalf("create first: %v", err)
	}
	second := NewUser(pool.ID, "second")
	second.Attributes = map[string]string{"preferred_username": "nick"}
	if err := s.CreateUser(second); !errors.Is(err, ErrAliasExists) {
		t.Fatalf("duplicate preferred_username returned %v, want ErrAliasExists", err)
	}
}

// A pool without alias configuration neither claims values nor resolves
// them: duplicate emails coexist and the email lookup reports
// ErrUserNotFound.
func TestAliasRulesOnlyApplyToConfiguredPools(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("plain-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	if err := s.CreateUser(newAliasClaimUser(pool.ID, "one", "dup@example.com", true)); err != nil {
		t.Fatalf("create one: %v", err)
	}
	if err := s.CreateUser(newAliasClaimUser(pool.ID, "two", "dup@example.com", true)); err != nil {
		t.Fatalf("duplicate verified email in a pool without aliases must coexist: %v", err)
	}
	if _, err := s.GetUser(pool.ID, "dup@example.com"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("email lookup in a pool without aliases returned %v, want ErrUserNotFound", err)
	}
}

// Deleting a pool removes its alias index entries with the rest of its
// records.
func TestDeleteUserPoolSweepsAliasIndex(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool := NewUserPool("swept-pool", "us-east-1")
	pool.AliasAttributes = []string{"email"}
	if _, err := s.CreateUserPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := s.CreateUser(newAliasClaimUser(pool.ID, "claimer", "sweep@example.com", true)); err != nil {
		t.Fatalf("create claimer: %v", err)
	}

	if err := s.DeleteUserPool(pool.ID); err != nil {
		t.Fatalf("delete pool: %v", err)
	}
	remaining := 0
	_ = s.usersStore.ScanPrefix("aliasidx:"+pool.ID+"#", func(key string, _ []byte) error {
		remaining++
		return nil
	})
	if remaining != 0 {
		t.Fatalf("pool deletion left %d alias index entries behind", remaining)
	}
}
