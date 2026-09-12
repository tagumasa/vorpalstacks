package cognitoidentityprovider

import (
	"sync"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// newUserIntegrityStore opens a throwaway store for the user-record
// integrity tests.
func newUserIntegrityStore(t *testing.T) *CognitoStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewCognitoStore(st, "000000000000", "us-east-1")
}

// An attribute update written from a membership-stale snapshot must not
// resurrect a group membership a concurrent operation removed, and a
// concurrent membership operation must not drop the attribute write: the
// single record guard plus the membership-ownership rule make both survive
// in either serial order.
func TestUserAttributeUpdateAndMembershipBothSurvive(t *testing.T) {
	s := newUserIntegrityStore(t)
	pool, err := s.CreateUserPool(NewUserPool("integrity-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, "bob")
	user.Attributes = map[string]string{"email": "bob@example.com"}
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := s.CreateGroup(NewGroup(pool.ID, "team")); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if err := s.AddUserToGroup(pool.ID, "team", "bob"); err != nil {
		t.Fatalf("add to group: %v", err)
	}

	// The stale-snapshot sequence: the attribute writer reads the user
	// (still a team member), the membership removal completes, then the
	// attribute writer persists from its snapshot.
	stale, err := s.GetUser(pool.ID, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.RemoveUserFromGroup(pool.ID, "team", "bob"); err != nil {
		t.Fatalf("remove from group: %v", err)
	}
	stale.Attributes["email"] = "moved@example.com"
	if err := s.UpdateUser(stale); err != nil {
		t.Fatalf("stale attribute update: %v", err)
	}
	after, err := s.GetUser(pool.ID, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if after.Attributes["email"] != "moved@example.com" {
		t.Fatalf("attribute write lost: %v", after.Attributes)
	}
	if len(after.Groups) != 0 {
		t.Fatalf("removed membership resurrected by the stale write: %v", after.Groups)
	}

	// The concurrent form: an attribute update and a group deletion race;
	// whatever the interleaving, the record keeps the attribute change and
	// loses the membership.
	if err := s.CreateGroup(NewGroup(pool.ID, "crew")); err != nil {
		t.Fatal(err)
	}
	if err := s.AddUserToGroup(pool.ID, "crew", "bob"); err != nil {
		t.Fatal(err)
	}
	snapshot, err := s.GetUser(pool.ID, "bob")
	if err != nil {
		t.Fatal(err)
	}
	snapshot.Attributes["email"] = "raced@example.com"
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_ = s.UpdateUser(snapshot)
	}()
	go func() {
		defer wg.Done()
		_ = s.DeleteGroup(pool.ID, "crew")
	}()
	wg.Wait()
	final, err := s.GetUser(pool.ID, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if final.Attributes["email"] != "raced@example.com" {
		t.Fatalf("attribute write lost under concurrency: %v", final.Attributes)
	}
	if len(final.Groups) != 0 {
		t.Fatalf("membership survived its group's deletion under concurrency: %v", final.Groups)
	}
}

// A DeleteUser that cannot decode the record must propagate the error
// instead of deleting it: the record's indexes, group references and alias
// claims would otherwise be orphaned by a partial deletion.
func TestDeleteUserUndecodableRecordPropagates(t *testing.T) {
	s := newUserIntegrityStore(t)
	pool, err := s.CreateUserPool(NewUserPool("integrity-pool-2", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, "carol")
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	key := userPoolUserKey(pool.ID, "carol")
	if err := s.usersStore.Put(key, "not-a-user-record"); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteUser(pool.ID, "carol"); err == nil {
		t.Fatal("DeleteUser succeeded on an undecodable record")
	}
	if !s.usersStore.Exists(key) {
		t.Fatal("DeleteUser removed an undecodable record instead of propagating the error")
	}
}

// The alias uniqueness gate fails closed: an alias index entry that cannot
// be decoded must not read as an unclaimed value a second user can take
// over.
func TestCreateUserAliasClaimFailsClosedOnUndecodableIndex(t *testing.T) {
	s := newUserIntegrityStore(t)
	pool := NewUserPool("integrity-pool-3", "us-east-1")
	pool.AliasAttributes = []string{"email"}
	pool, err := s.CreateUserPool(pool)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	first := NewUser(pool.ID, "dave")
	first.Attributes = map[string]string{"email": "shared@example.com", "email_verified": "true"}
	if err := s.CreateUser(first); err != nil {
		t.Fatalf("create first holder: %v", err)
	}
	if err := s.usersStore.Put(aliasIndexKey(pool.ID, "email", "shared@example.com"), "{not-json"); err != nil {
		t.Fatal(err)
	}

	second := NewUser(pool.ID, "erin")
	second.Attributes = map[string]string{"email": "shared@example.com", "email_verified": "true"}
	if err := s.CreateUser(second); err == nil {
		t.Fatal("alias claim accepted over an undecodable uniqueness entry")
	}
}

// User creation racing a pool deletion leaves no residue: the creator's
// pool-existence check runs under the pool lock, and the cascade re-sweeps
// after the tombstone, so every user lands either before a sweep or not at
// all.
func TestDeleteUserPoolSweepsRacingUserCreates(t *testing.T) {
	s := newUserIntegrityStore(t)
	pool, err := s.CreateUserPool(NewUserPool("integrity-pool-4", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	// Seed one user so the cascade has a non-empty family to sweep.
	if err := s.CreateUser(NewUser(pool.ID, "seed")); err != nil {
		t.Fatalf("create seed user: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 32; i++ {
			_ = s.CreateUser(NewUser(pool.ID, "racer"))
		}
	}()
	go func() {
		defer wg.Done()
		if err := s.DeleteUserPool(pool.ID); err != nil {
			t.Errorf("delete pool: %v", err)
		}
	}()
	wg.Wait()

	if n := countPrefix(t, s.usersStore, pool.ID+"#"); n != 0 {
		t.Fatalf("%d user records survived the pool deletion", n)
	}
	if n := countPrefix(t, s.usersStore, "userci:"+pool.ID+"#"); n != 0 {
		t.Fatalf("%d case-insensitive index entries survived the pool deletion", n)
	}
	if n := countPrefix(t, s.usersStore, "aliasidx:"+pool.ID+"#"); n != 0 {
		t.Fatalf("%d alias index entries survived the pool deletion", n)
	}
}
