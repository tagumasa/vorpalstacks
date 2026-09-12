package cognitoidentityprovider

import (
	"errors"
	"testing"
)

// The secondary-index getters resolve identifiers through their index only.
// The ghost records below are record-shaped keys with values that cannot be
// parsed as user or client records: a bucket scan would fail on them, so the
// sentinel answers prove the unknown identifiers were resolved without a
// scan.
func TestIndexGettersUnknownIDsAreNotFoundWithoutScan(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("index-lookups", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, "alice")
	user.ProviderName = "Google"
	user.ProviderAttributeValue = "google-123"
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := s.CreateUserPoolClient(&UserPoolClient{ClientID: "client-7", ClientName: "web", UserPoolID: pool.ID}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	if err := s.usersStore.Put(userPoolUserKey(pool.ID, "ghost"), "not a user record"); err != nil {
		t.Fatalf("seed ghost user record: %v", err)
	}
	if err := s.clientsStore.Put(userPoolClientKey(pool.ID, "ghost-client"), "not a client record"); err != nil {
		t.Fatalf("seed ghost client record: %v", err)
	}

	if _, err := s.GetUserByID("no-such-id"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetUserByID unknown id: err = %v, want ErrUserNotFound", err)
	}
	if _, err := s.GetUserByProvider(pool.ID, "Google", "no-such-value"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetUserByProvider unknown value: err = %v, want ErrUserNotFound", err)
	}
	if _, err := s.GetUserPoolByClientID("no-such-client"); !errors.Is(err, ErrClientNotFound) {
		t.Fatalf("GetUserPoolByClientID unknown client: err = %v, want ErrClientNotFound", err)
	}
}

// The index is written together with the record it resolves, so a record
// whose index entry is missing is not resurrected by scanning the bucket and
// rebuilding the entry — the identifier is simply unknown. The round trips
// with intact index entries guard against the lookup being broken outright.
func TestIndexGettersDoNotResurrectUnindexedRecords(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("unindexed-records", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	unindexed := NewUser(pool.ID, "alice")
	unindexed.ProviderName = "Google"
	unindexed.ProviderAttributeValue = "google-123"
	if err := s.CreateUser(unindexed); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := s.CreateUserPoolClient(&UserPoolClient{ClientID: "client-7", ClientName: "web", UserPoolID: pool.ID}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	if err := s.usersStore.Delete(userIndexKey(unindexed.ID)); err != nil {
		t.Fatalf("drop user index: %v", err)
	}
	if err := s.usersStore.Delete(providerIndexKey(pool.ID, "Google", "google-123")); err != nil {
		t.Fatalf("drop provider index: %v", err)
	}
	if err := s.clientsStore.Delete(clientIndexKey("client-7")); err != nil {
		t.Fatalf("drop client index: %v", err)
	}

	if _, err := s.GetUserByID(unindexed.ID); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetUserByID without index: err = %v, want ErrUserNotFound", err)
	}
	if _, err := s.GetUserByProvider(pool.ID, "Google", "google-123"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetUserByProvider without index: err = %v, want ErrUserNotFound", err)
	}
	if _, err := s.GetUserPoolByClientID("client-7"); !errors.Is(err, ErrClientNotFound) {
		t.Fatalf("GetUserPoolByClientID without index: err = %v, want ErrClientNotFound", err)
	}
	if s.usersStore.Exists(userIndexKey(unindexed.ID)) {
		t.Fatal("lookup must not rebuild the user index entry")
	}
	if s.clientsStore.Exists(clientIndexKey("client-7")) {
		t.Fatal("lookup must not rebuild the client index entry")
	}

	indexed := NewUser(pool.ID, "bob")
	indexed.ProviderName = "Google"
	indexed.ProviderAttributeValue = "google-456"
	if err := s.CreateUser(indexed); err != nil {
		t.Fatalf("create indexed user: %v", err)
	}
	if err := s.CreateUserPoolClient(&UserPoolClient{ClientID: "client-8", ClientName: "mobile", UserPoolID: pool.ID}); err != nil {
		t.Fatalf("create indexed client: %v", err)
	}
	got, err := s.GetUserByID(indexed.ID)
	if err != nil {
		t.Fatalf("GetUserByID round trip: %v", err)
	}
	if got.Username != "bob" {
		t.Fatalf("GetUserByID round trip: username = %q, want bob", got.Username)
	}
	if _, err := s.GetUserByProvider(pool.ID, "Google", "google-456"); err != nil {
		t.Fatalf("GetUserByProvider round trip: %v", err)
	}
	if _, err := s.GetUserPoolByClientID("client-8"); err != nil {
		t.Fatalf("GetUserPoolByClientID round trip: %v", err)
	}
}

// The provider index entry is invalidated when the federated identity changes
// or is unlinked: the stale entry must stop resolving the user, so a later
// lookup of the released identity reports it as unknown instead of
// resurrecting the unlinked user.
func TestUpdateUserInvalidatesProviderIndex(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("provider-index-invalidation", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, "alice")
	user.ProviderName = "Google"
	user.ProviderAttributeValue = "google-123"
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// Relinking under a different identity removes the previous entry while
	// the new one resolves.
	changed, err := s.GetUser(pool.ID, "alice")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	changed.ProviderName = "Facebook"
	changed.ProviderAttributeValue = "fb-456"
	if err := s.UpdateUser(changed); err != nil {
		t.Fatalf("update user: %v", err)
	}
	if _, err := s.GetUserByProvider(pool.ID, "Google", "google-123"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetUserByProvider previous identity after change: err = %v, want ErrUserNotFound", err)
	}
	got, err := s.GetUserByProvider(pool.ID, "Facebook", "fb-456")
	if err != nil {
		t.Fatalf("GetUserByProvider new identity: %v", err)
	}
	if got.Username != "alice" {
		t.Fatalf("GetUserByProvider new identity: username = %q, want alice", got.Username)
	}

	// Unlinking clears the entry outright.
	unlinked, err := s.GetUser(pool.ID, "alice")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	unlinked.ProviderName = ""
	unlinked.ProviderAttributeName = ""
	unlinked.ProviderAttributeValue = ""
	if err := s.UpdateUser(unlinked); err != nil {
		t.Fatalf("unlink user: %v", err)
	}
	if _, err := s.GetUserByProvider(pool.ID, "Facebook", "fb-456"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("GetUserByProvider after unlink: err = %v, want ErrUserNotFound", err)
	}
	if s.usersStore.Exists(providerIndexKey(pool.ID, "Facebook", "fb-456")) {
		t.Fatal("unlink must remove the provider index entry")
	}
}
