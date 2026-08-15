package cognitoidentityprovider

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

const challengeTestClientID = "test-client-id"

type challengeTestEnv struct {
	svc    *CognitoService
	reqCtx *request.RequestContext
	store  cognitostore.CognitoStoreInterface
	pool   *cognitostore.UserPool
	user   *cognitostore.User
}

func newChallengeTestEnv(t *testing.T) *challengeTestEnv {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
	svc := NewCognitoService("000000000000", "us-east-1")
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := store.CreateUserPool(cognitostore.NewUserPool("testpool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   challengeTestClientID,
		ClientName: "test-client",
	}); err != nil {
		t.Fatal(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("OldPass123!"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	user := cognitostore.NewUser(pool.ID, "victim")
	user.UserStatus = "CONFIRMED"
	user.PasswordHash = string(hash)
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}

	return &challengeTestEnv{svc: svc, reqCtx: reqCtx, store: store, pool: pool, user: user}
}

func challengeReq(params map[string]interface{}) *request.ParsedRequest {
	return &request.ParsedRequest{Parameters: params}
}

func respondToNewPassword(env *challengeTestEnv, session, newPassword string) (interface{}, error) {
	return env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"ClientId":      challengeTestClientID,
		"ChallengeName": "NEW_PASSWORD_REQUIRED",
		"Session":       session,
		"ChallengeResponses": map[string]interface{}{
			"USERNAME":     "victim",
			"NEW_PASSWORD": newPassword,
		},
	}))
}

// A PASSWORD_VERIFIER session issued by InitiateAuth(USER_SRP_AUTH) must not
// satisfy a NEW_PASSWORD_REQUIRED response; otherwise anyone could reset a
// victim's password without knowing the old one.
func TestNewPasswordChallengeRejectsForeignSessionType(t *testing.T) {
	env := newChallengeTestEnv(t)

	sess, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "PASSWORD_VERIFIER", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := respondToNewPassword(env, sess, "AttackPass123!"); err == nil {
		t.Fatal("NEW_PASSWORD_REQUIRED response accepted with a PASSWORD_VERIFIER session")
	}

	user, err := env.store.GetUser(env.pool.ID, "victim")
	if err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("AttackPass123!")) == nil {
		t.Fatal("password was overwritten via a foreign challenge session")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("OldPass123!")) != nil {
		t.Fatal("original password was disturbed")
	}
}

// A session bound to another app client must not answer a challenge for
// this client.
func TestNewPasswordChallengeRejectsForeignClient(t *testing.T) {
	env := newChallengeTestEnv(t)

	sess, err := mintChallengeSession(env.store, env.pool.ID, "other-client", "victim", "NEW_PASSWORD_REQUIRED", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := respondToNewPassword(env, sess, "NewPass456!"); err == nil {
		t.Fatal("NEW_PASSWORD_REQUIRED response accepted with a foreign client session")
	}
}

// The happy path rotates the password and issues tokens exactly once; the
// session is burned afterwards.
func TestNewPasswordChallengeHappyPathBurnsSession(t *testing.T) {
	env := newChallengeTestEnv(t)

	sess, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "NEW_PASSWORD_REQUIRED", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := respondToNewPassword(env, sess, "NewPass456!")
	if err != nil {
		t.Fatal(err)
	}
	authResult, ok := resp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{})
	if !ok || authResult["AccessToken"] == "" || authResult["RefreshToken"] == "" {
		t.Fatalf("unexpected AuthenticationResult: %#v", resp)
	}

	user, err := env.store.GetUser(env.pool.ID, "victim")
	if err != nil {
		t.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("NewPass456!")) != nil {
		t.Fatal("password was not rotated")
	}
	if user.UserStatus != "CONFIRMED" {
		t.Fatalf("unexpected status %q", user.UserStatus)
	}

	if _, err := respondToNewPassword(env, sess, "AgainPass789!"); err == nil {
		t.Fatal("burned session accepted a second response")
	}
}

// A disabled user must not be revived through the password challenge.
func TestNewPasswordChallengeRejectsDisabledUser(t *testing.T) {
	env := newChallengeTestEnv(t)

	env.user.Enabled = false
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}

	sess, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "NEW_PASSWORD_REQUIRED", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := respondToNewPassword(env, sess, "NewPass456!"); err == nil {
		t.Fatal("disabled user passed the password challenge")
	}
}

// MFA challenges must be answered with a session of the matching type, not
// with a selector session or a bare USERNAME.
func TestMfaChallengeRequiresTypedSession(t *testing.T) {
	env := newChallengeTestEnv(t)

	sel, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "SELECT_CHALLENGE", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	_, err = env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"ClientId":      challengeTestClientID,
		"ChallengeName": "SMS_OTP",
		"Session":       sel,
		"ChallengeResponses": map[string]interface{}{
			"USERNAME":     "victim",
			"SMS_MFA_CODE": "123456",
		},
	}))
	if err == nil {
		t.Fatal("SMS_OTP response accepted with a SELECT_CHALLENGE session")
	}

	// Sessionless responses must fail as well.
	_, err = env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"ClientId":      challengeTestClientID,
		"ChallengeName": "SMS_OTP",
		"ChallengeResponses": map[string]interface{}{
			"USERNAME":     "victim",
			"SMS_MFA_CODE": "123456",
		},
	}))
	if err == nil {
		t.Fatal("SMS_OTP response accepted without a session")
	}
}

// SELECT_CHALLENGE mints a fresh session of the selected type which then
// answers the selected challenge; wrong answers stay bounded by the
// failed-attempt budget.
func TestSelectChallengeMintsTypedSessionAndBoundedOTP(t *testing.T) {
	env := newChallengeTestEnv(t)

	env.user.ConfirmationCode = "123456"
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}

	sel, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "SELECT_CHALLENGE", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"ClientId":           challengeTestClientID,
		"ChallengeName":      "SELECT_CHALLENGE",
		"Session":            sel,
		"USERNAME":           "victim",
		"SELECTED_CHALLENGE": "SMS_OTP",
	}))
	if err != nil {
		t.Fatal(err)
	}
	m := resp.(map[string]interface{})
	if m["ChallengeName"] != "SMS_OTP" {
		t.Fatalf("unexpected challenge %v", m["ChallengeName"])
	}
	typedSession, _ := m["Session"].(string)
	if typedSession == "" || typedSession == sel {
		t.Fatalf("selector response must mint a fresh session, got %q", typedSession)
	}
	cs, err := env.store.GetChallengeSession(typedSession)
	if err != nil || cs.ChallengeName != "SMS_OTP" {
		t.Fatalf("minted session not typed SMS_OTP: %v %+v", err, cs)
	}

	// Correct code completes authentication.
	resp, err = env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"ClientId":      challengeTestClientID,
		"ChallengeName": "SMS_OTP",
		"Session":       typedSession,
		"USERNAME":      "victim",
		"SMS_MFA_CODE":  "123456",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := resp.(map[string]interface{})["AuthenticationResult"]; !ok {
		t.Fatalf("expected AuthenticationResult, got %#v", resp)
	}

	// Wrong answers are bounded: the budget exhausts the session.
	budgetSess, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "SMS_OTP", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < maxChallengeAttempts; i++ {
		_, err = env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
			"ClientId":      challengeTestClientID,
			"ChallengeName": "SMS_OTP",
			"Session":       budgetSess,
			"USERNAME":      "victim",
			"SMS_MFA_CODE":  "000000",
		}))
		if err == nil {
			t.Fatalf("wrong code accepted on attempt %d", i+1)
		}
	}
	if _, err := validateChallengeSession(env.store, budgetSess, "SMS_OTP", env.pool.ID, challengeTestClientID, "victim"); err == nil {
		t.Fatal("session still valid after exhausting the attempt budget")
	}
}

// A session minted for another challenge type must not overwrite a verified
// MFA configuration via AssociateSoftwareToken.
func TestAssociateSoftwareTokenRequiresMfaSetupSession(t *testing.T) {
	env := newChallengeTestEnv(t)

	env.user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{
		Enabled:   true,
		SecretKey: "SEEDSECRET",
		Verified:  true,
	}
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}

	verifierSess, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "PASSWORD_VERIFIER", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := env.svc.AssociateSoftwareToken(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"Session": verifierSess,
	})); err == nil {
		t.Fatal("AssociateSoftwareToken accepted a PASSWORD_VERIFIER session")
	}
	user, err := env.store.GetUser(env.pool.ID, "victim")
	if err != nil {
		t.Fatal(err)
	}
	if user.SoftwareTokenMfa == nil || user.SoftwareTokenMfa.SecretKey != "SEEDSECRET" || !user.SoftwareTokenMfa.Verified {
		t.Fatal("verified MFA configuration was disturbed via a foreign session")
	}

	mfaSess, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "MFA_SETUP", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := env.svc.AssociateSoftwareToken(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"Session": mfaSess,
	}))
	if err != nil {
		t.Fatal(err)
	}
	m := resp.(map[string]interface{})
	if m["SecretCode"] == "" {
		t.Fatalf("expected SecretCode, got %#v", resp)
	}
	user, err = env.store.GetUser(env.pool.ID, "victim")
	if err != nil {
		t.Fatal(err)
	}
	if user.SoftwareTokenMfa.Verified {
		t.Fatal("newly associated secret must not be verified")
	}
}

// A refresh token issued by one pool must not mint tokens for another pool.
func TestRefreshTokenRejectsForeignPool(t *testing.T) {
	env := newChallengeTestEnv(t)

	otherPool, err := env.store.CreateUserPool(cognitostore.NewUserPool("otherpool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}
	if err := env.store.CreateRefreshToken(&cognitostore.RefreshToken{
		Token:      "rt-test-value",
		Expires:    time.Now().Add(time.Hour),
		ClientID:   challengeTestClientID,
		UserPoolID: env.pool.ID,
		UserID:     env.user.ID,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := env.svc.refreshAuthToken(env.reqCtx, otherPool.ID, "rt-test-value"); err == nil {
		t.Fatal("refresh token accepted for a foreign user pool")
	}
}
