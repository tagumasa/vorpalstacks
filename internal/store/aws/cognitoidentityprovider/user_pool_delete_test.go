package cognitoidentityprovider

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// countPrefix walks a bucket prefix and returns the number of live keys.
func countPrefix(t *testing.T, b *common.BaseStore, prefix string) int {
	t.Helper()
	n := 0
	if err := b.ScanPrefix(prefix, func(string, []byte) error {
		n++
		return nil
	}); err != nil {
		t.Fatalf("scan prefix %q: %v", prefix, err)
	}
	return n
}

// countChallengeSessionsForPool counts the challenge sessions whose record
// belongs to the pool; the pool lives in the record, not in the key.
func countChallengeSessionsForPool(t *testing.T, s *CognitoStore, poolID string) int {
	t.Helper()
	n := 0
	if err := s.challengeSessionsStore.ForEach(func(_ string, value []byte) error {
		var session ChallengeSession
		if err := json.Unmarshal(value, &session); err == nil && session.UserPoolID == poolID {
			n++
		}
		return nil
	}); err != nil {
		t.Fatalf("walk challenge sessions: %v", err)
	}
	return n
}

// TestDeleteUserPoolCascadesEveryFamily seeds one record of every entity
// family a pool owns, deletes the pool, and asserts every bucket and key
// prefix the family lives in is empty — while a second pool's records
// survive untouched, guarding the prefixes against over-broad sweeps.
func TestDeleteUserPoolCascadesEveryFamily(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	s := NewCognitoStore(st, "000000000000", "us-east-1")

	pool, err := s.CreateUserPool(NewUserPool("cascade-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	user := NewUser(pool.ID, "alice")
	if err := s.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	expires := time.Now().Add(time.Hour).UTC()
	if err := s.CreateRefreshToken(&RefreshToken{Token: "rt-1", UserPoolID: pool.ID, UserID: user.ID, Expires: expires}); err != nil {
		t.Fatalf("create refresh token: %v", err)
	}
	if err := s.CreateIDToken(&IDToken{Token: "idt-1", UserPoolID: pool.ID, UserID: user.ID, Expires: expires}); err != nil {
		t.Fatalf("create id token: %v", err)
	}
	if err := s.CreateAccessToken(&AccessToken{Token: "at-1", UserPoolID: pool.ID, UserID: user.ID, Expires: expires}); err != nil {
		t.Fatalf("create access token: %v", err)
	}
	if err := s.CreateGroup(NewGroup(pool.ID, "graders")); err != nil {
		t.Fatalf("create group: %v", err)
	}
	if err := s.CreateUserPoolClient(&UserPoolClient{UserPoolID: pool.ID, ClientID: "client-1", ClientName: "web"}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	// Both challenge-session key shapes: a bare auth-flow session ID and a
	// composite WebAuthn registration key.
	now := time.Now().UTC()
	for _, session := range []*ChallengeSession{
		{SessionID: "SESSION_cascade", UserPoolID: pool.ID, Username: "alice", ChallengeName: "SMS_MFA", CreatedAt: now, ExpiresAt: expires},
		{SessionID: "webauthn-reg#" + pool.ID + "#" + user.ID, UserPoolID: pool.ID, ChallengeName: "MFA", CreatedAt: now, ExpiresAt: expires},
	} {
		if err := s.SaveChallengeSession(session); err != nil {
			t.Fatalf("save challenge session %s: %v", session.SessionID, err)
		}
	}
	if err := s.CreateResourceServer(&ResourceServer{UserPoolID: pool.ID, Identifier: "api.example", Name: "api"}); err != nil {
		t.Fatalf("create resource server: %v", err)
	}
	if err := s.CreateIdentityProvider(&IdentityProvider{UserPoolID: pool.ID, ProviderName: "Google", ProviderType: "Google"}); err != nil {
		t.Fatalf("create identity provider: %v", err)
	}
	if err := s.SetUserPoolDomain("cascade-domain", &UserPoolDomain{Domain: "cascade-domain", UserPoolID: pool.ID}); err != nil {
		t.Fatalf("set domain: %v", err)
	}
	if err := s.TagStore.Tag(pool.Arn, map[string]string{"Team": "core"}); err != nil {
		t.Fatalf("tag pool: %v", err)
	}
	if err := s.CreateDevice(&Device{UserPoolID: pool.ID, UserID: user.ID, DeviceKey: "device-1"}); err != nil {
		t.Fatalf("create device: %v", err)
	}
	if err := s.CreateAuthEvent(&AuthEvent{UserPoolID: pool.ID, UserID: user.ID, EventID: "evt-1", EventType: "SignIn", EventResponse: "Pass", CreationDate: now}); err != nil {
		t.Fatalf("create auth event: %v", err)
	}
	if err := s.CreateWebAuthnCredential(&WebAuthnCredential{UserPoolID: pool.ID, UserID: user.ID, CredentialID: "cred-1", PublicKey: "pub", CreatedAt: now}); err != nil {
		t.Fatalf("create webauthn credential: %v", err)
	}
	if err := s.CreateUserImportJob(&UserImportJob{JobID: "job-1", JobName: "import", UserPoolID: pool.ID, Status: "InProgress", CreationDate: now}); err != nil {
		t.Fatalf("create import job: %v", err)
	}
	if err := s.SaveLogDeliveryConfiguration(&LogDeliveryConfiguration{UserPoolID: pool.ID}); err != nil {
		t.Fatalf("save log delivery: %v", err)
	}
	for _, cfg := range []*RiskConfiguration{{UserPoolID: pool.ID}, {UserPoolID: pool.ID, ClientID: "client-1"}} {
		if err := s.SaveRiskConfiguration(cfg); err != nil {
			t.Fatalf("save risk configuration: %v", err)
		}
	}
	for _, ui := range []*UICustomization{{UserPoolID: pool.ID}, {UserPoolID: pool.ID, ClientID: "client-1"}} {
		if err := s.SaveUICustomization(ui); err != nil {
			t.Fatalf("save ui customisation: %v", err)
		}
	}
	if err := s.SaveManagedLoginBranding(&ManagedLoginBranding{UserPoolID: pool.ID, ManagedLoginBrandingId: "branding-1"}); err != nil {
		t.Fatalf("save branding: %v", err)
	}
	if err := s.SaveTerms(&Terms{UserPoolID: pool.ID, TermsID: "terms-1", TermsName: "tos"}); err != nil {
		t.Fatalf("save terms: %v", err)
	}
	if err := s.SaveUserPoolReplica(&UserPoolReplica{UserPoolID: pool.ID, RegionName: "us-west-2", Status: "ENABLED"}); err != nil {
		t.Fatalf("save replica: %v", err)
	}

	// A second pool seeds the same families; every record must survive the
	// first pool's deletion.
	pool2, err := s.CreateUserPool(NewUserPool("cascade-pool-2", "us-east-1"))
	if err != nil {
		t.Fatalf("create second pool: %v", err)
	}
	user2 := NewUser(pool2.ID, "bob")
	if err := s.CreateUser(user2); err != nil {
		t.Fatalf("create second user: %v", err)
	}
	if err := s.CreateDevice(&Device{UserPoolID: pool2.ID, UserID: user2.ID, DeviceKey: "device-2"}); err != nil {
		t.Fatalf("create second device: %v", err)
	}
	if err := s.CreateUserImportJob(&UserImportJob{JobID: "job-2", JobName: "import2", UserPoolID: pool2.ID, Status: "Created", CreationDate: now}); err != nil {
		t.Fatalf("create second import job: %v", err)
	}
	if err := s.SaveRiskConfiguration(&RiskConfiguration{UserPoolID: pool2.ID}); err != nil {
		t.Fatalf("save second risk configuration: %v", err)
	}
	if err := s.SetUserPoolDomain("cascade-domain-2", &UserPoolDomain{Domain: "cascade-domain-2", UserPoolID: pool2.ID}); err != nil {
		t.Fatalf("set second domain: %v", err)
	}
	if err := s.SaveChallengeSession(&ChallengeSession{SessionID: "SESSION_survivor", UserPoolID: pool2.ID, ChallengeName: "MFA", CreatedAt: now, ExpiresAt: expires}); err != nil {
		t.Fatalf("save survivor session: %v", err)
	}

	if err := s.DeleteUserPool(pool.ID); err != nil {
		t.Fatalf("delete pool: %v", err)
	}

	if _, err := s.GetUserPool(pool.ID); !errors.Is(err, ErrUserPoolNotFound) {
		t.Fatalf("pool record: got %v, want ErrUserPoolNotFound", err)
	}
	if _, err := s.GetUser(pool.ID, "alice"); !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("user: got %v, want ErrUserNotFound", err)
	}
	if _, err := s.GetGroup(pool.ID, "graders"); !errors.Is(err, ErrGroupNotFound) {
		t.Fatalf("group: got %v, want ErrGroupNotFound", err)
	}
	if _, err := s.GetUserPoolClient(pool.ID, "client-1"); !errors.Is(err, ErrClientNotFound) {
		t.Fatalf("client: got %v, want ErrClientNotFound", err)
	}
	if _, err := s.GetUserPoolDomain("cascade-domain"); !errors.Is(err, ErrUserPoolDomainNotFound) {
		t.Fatalf("domain: got %v, want ErrUserPoolDomainNotFound", err)
	}
	tags, err := s.TagStore.List(pool.Arn)
	if err != nil {
		t.Fatalf("list tags of deleted pool: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("tags of deleted pool: got %v, want none", tags)
	}

	for _, tc := range []struct {
		bucket *common.BaseStore
		name   string
	}{
		{s.refreshTokensStore, "refresh tokens"},
		{s.idTokensStore, "id tokens"},
		{s.accessTokensStore, "access tokens"},
		{s.devicesStore, "devices"},
		{s.authEventsStore, "auth events"},
		{s.webauthnStore, "webauthn credentials"},
		{s.userImportJobsStore, "import jobs"},
	} {
		if n := countPrefix(t, tc.bucket, pool.ID+"#"); n != 0 {
			t.Fatalf("%s of the deleted pool: %d keys remain", tc.name, n)
		}
	}
	if n := countChallengeSessionsForPool(t, s, pool.ID); n != 0 {
		t.Fatalf("challenge sessions of the deleted pool: %d remain", n)
	}
	for _, prefix := range []string{
		logDeliveryKey(pool.ID),
		riskConfigKey(pool.ID, ""),
		uiCustomizationKey(pool.ID, ""),
		managedLoginBrandingPrefix(pool.ID),
		termsPrefix(pool.ID),
		userPoolReplicaPrefix(pool.ID),
		resourceServerPrefix(pool.ID),
		identityProviderPrefix(pool.ID),
		"aliasidx:" + pool.ID + "#",
	} {
		if n := countPrefix(t, s.BaseStore, prefix); n != 0 {
			t.Fatalf("prefix %q of the deleted pool: %d keys remain", prefix, n)
		}
	}

	// The second pool's families are untouched by the cascade.
	if _, err := s.GetUser(pool2.ID, "bob"); err != nil {
		t.Fatalf("second pool user disturbed: %v", err)
	}
	if _, err := s.GetDevice(pool2.ID, user2.ID, "device-2"); err != nil {
		t.Fatalf("second pool device disturbed: %v", err)
	}
	job2, err := s.GetUserImportJob(pool2.ID, "job-2")
	if err != nil {
		t.Fatalf("second pool import job disturbed: %v", err)
	}
	if job2.Status != "Created" {
		t.Fatalf("second pool import job status: got %q, want Created", job2.Status)
	}
	if _, err := s.GetRiskConfiguration(pool2.ID, ""); err != nil {
		t.Fatalf("second pool risk configuration disturbed: %v", err)
	}
	if _, err := s.GetUserPoolDomain("cascade-domain-2"); err != nil {
		t.Fatalf("second pool domain disturbed: %v", err)
	}
	if n := countChallengeSessionsForPool(t, s, pool2.ID); n != 1 {
		t.Fatalf("second pool challenge sessions: got %d, want 1", n)
	}
	if n := countPrefix(t, s.userImportJobsStore, pool2.ID+"#"); n != 1 {
		t.Fatalf("second pool import job keys: got %d, want 1", n)
	}
}
