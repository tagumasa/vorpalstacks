package cognitoidentityprovider

import (
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// mintTestTokens creates one live token of each family for the pool's user
// and returns their values in refresh, ID, access order.
func mintTestTokens(t *testing.T, s *CognitoStore, poolID, userID string) (string, string, string) {
	t.Helper()
	expires := time.Now().Add(time.Hour)
	rt := NewRefreshToken(poolID, userID, "client-1", "openid", expires)
	if err := s.CreateRefreshToken(rt); err != nil {
		t.Fatalf("create refresh token: %v", err)
	}
	it := NewIDToken(poolID, userID, "client-1", "openid", expires, nil)
	if err := s.CreateIDToken(it); err != nil {
		t.Fatalf("create ID token: %v", err)
	}
	at := NewAccessToken(poolID, userID, "client-1", "openid", expires)
	if err := s.CreateAccessToken(at); err != nil {
		t.Fatalf("create access token: %v", err)
	}
	return rt.Token, it.Token, at.Token
}

// keyExists reports whether a raw key resolves in the bucket; index
// maintenance is asserted against it directly.
func keyExists(t *testing.T, store *common.BaseStore, key string) bool {
	t.Helper()
	var probe string
	return store.Get(key, &probe) == nil
}

// The token-value getters resolve through the value index, not a bucket
// scan: the ghost keys below hold values that cannot parse as token
// records, so a scan would fail on them, yet each minted token resolves
// and an unknown value reports the not-found sentinel.
func TestTokenValueIndexResolvesWithoutScan(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("token-index", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, "alice")
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := s.refreshTokensStore.Put(tokenKey(pool.ID, user.ID, "ghost-rt"), "not a token record"); err != nil {
		t.Fatalf("seed ghost refresh record: %v", err)
	}
	if err := s.idTokensStore.Put(tokenKey(pool.ID, user.ID, "ghost-it"), "not a token record"); err != nil {
		t.Fatalf("seed ghost ID record: %v", err)
	}
	if err := s.accessTokensStore.Put(tokenKey(pool.ID, user.ID, "ghost-at"), "not a token record"); err != nil {
		t.Fatalf("seed ghost access record: %v", err)
	}
	refresh, id, access := mintTestTokens(t, s, pool.ID, user.ID)

	if _, err := s.GetRefreshTokenByValue("no-such-refresh"); !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("unknown refresh value: err = %v, want ErrTokenNotFound", err)
	}
	if _, err := s.GetIDTokenByValue("no-such-id"); !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("unknown ID value: err = %v, want ErrTokenNotFound", err)
	}
	if _, err := s.GetAccessTokenByValue("no-such-access"); !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("unknown access value: err = %v, want ErrTokenNotFound", err)
	}

	rt, err := s.GetRefreshTokenByValue(refresh)
	if err != nil || rt == nil || rt.Token != refresh {
		t.Fatalf("refresh by value: (%v, %+v)", err, rt)
	}
	it, err := s.GetIDTokenByValue(id)
	if err != nil || it == nil || it.Token != id {
		t.Fatalf("ID by value: (%v, %+v)", err, it)
	}
	at, err := s.GetAccessTokenByValue(access)
	if err != nil || at == nil || at.Token != access {
		t.Fatalf("access by value: (%v, %+v)", err, at)
	}
}

// Every deletion path takes the value-index entry with the record it
// removes: explicit per-token deletes, expiry discovered through either
// getter, the per-user refresh sweep, and the all-tokens sweep.
func TestTokenValueIndexMaintainedOnDelete(t *testing.T) {
	s := newUserPoolTestStore(t)
	pool, err := s.CreateUserPool(NewUserPool("token-index-del", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, "bob")
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// Explicit per-token delete.
	refresh, _, access := mintTestTokens(t, s, pool.ID, user.ID)
	if err := s.DeleteAccessToken(pool.ID, user.ID, access); err != nil {
		t.Fatalf("delete access token: %v", err)
	}
	if _, err := s.GetAccessTokenByValue(access); !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("access value after delete: err = %v, want ErrTokenNotFound", err)
	}
	if keyExists(t, s.accessTokensStore, tokenIndexKey(access)) {
		t.Fatal("access index entry survived deletion")
	}
	if err := s.DeleteRefreshToken(pool.ID, user.ID, refresh); err != nil {
		t.Fatalf("delete refresh token: %v", err)
	}
	if _, err := s.GetRefreshTokenByValue(refresh); !errors.Is(err, ErrTokenNotFound) {
		t.Fatalf("refresh value after delete: err = %v, want ErrTokenNotFound", err)
	}
	if keyExists(t, s.refreshTokensStore, tokenIndexKey(refresh)) {
		t.Fatal("refresh index entry survived deletion")
	}

	// Expiry discovered through the value getter removes record and index.
	expired := NewIDToken(pool.ID, user.ID, "client-1", "openid", time.Now().Add(-time.Minute), nil)
	if err := s.CreateIDToken(expired); err != nil {
		t.Fatalf("create expired ID token: %v", err)
	}
	if _, err := s.GetIDTokenByValue(expired.Token); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expired ID by value: err = %v, want ErrTokenExpired", err)
	}
	if keyExists(t, s.idTokensStore, tokenKey(pool.ID, user.ID, expired.Token)) {
		t.Fatal("expired ID record survived expiry")
	}
	if keyExists(t, s.idTokensStore, tokenIndexKey(expired.Token)) {
		t.Fatal("expired ID index entry survived expiry")
	}

	// Expiry discovered through the primary-key getter removes the index too.
	expiredRefresh := NewRefreshToken(pool.ID, user.ID, "client-1", "openid", time.Now().Add(-time.Minute))
	if err := s.CreateRefreshToken(expiredRefresh); err != nil {
		t.Fatalf("create expired refresh token: %v", err)
	}
	if _, err := s.GetRefreshToken(pool.ID, user.ID, expiredRefresh.Token); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expired refresh by key: err = %v, want ErrTokenExpired", err)
	}
	if keyExists(t, s.refreshTokensStore, tokenIndexKey(expiredRefresh.Token)) {
		t.Fatal("expired refresh index entry survived primary-key expiry")
	}

	// The per-user refresh sweep cleans its index entries.
	r2 := NewRefreshToken(pool.ID, user.ID, "client-1", "openid", time.Now().Add(time.Hour))
	if err := s.CreateRefreshToken(r2); err != nil {
		t.Fatalf("create second refresh token: %v", err)
	}
	if err := s.DeleteAllRefreshTokensForUser(pool.ID, user.ID); err != nil {
		t.Fatalf("refresh sweep: %v", err)
	}
	if keyExists(t, s.refreshTokensStore, tokenIndexKey(r2.Token)) {
		t.Fatal("swept refresh index entry survived")
	}

	// The all-tokens sweep cleans every family's index entries.
	r3, i3, a3 := mintTestTokens(t, s, pool.ID, user.ID)
	if err := s.DeleteUserTokens(pool.ID, user.ID); err != nil {
		t.Fatalf("token sweep: %v", err)
	}
	if keyExists(t, s.refreshTokensStore, tokenIndexKey(r3)) {
		t.Fatal("refresh index entry survived the token sweep")
	}
	if keyExists(t, s.idTokensStore, tokenIndexKey(i3)) {
		t.Fatal("ID index entry survived the token sweep")
	}
	if keyExists(t, s.accessTokensStore, tokenIndexKey(a3)) {
		t.Fatal("access index entry survived the token sweep")
	}
}
