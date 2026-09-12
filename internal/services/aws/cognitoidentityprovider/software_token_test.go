package cognitoidentityprovider

import (
	"context"
	"encoding/base32"
	"fmt"
	"strings"
	"testing"
	"time"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// codeForSecret derives the TOTP code for the given secret at the current
// step plus offset, mirroring how an authenticator app computes codes.
func codeForSecret(t *testing.T, secret string, stepOffset int64) string {
	t.Helper()
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		t.Fatalf("decode secret: %v", err)
	}
	return totpCodeAt(key, time.Now().Unix()/30+stepOffset)
}

// wrongCodeForSecret returns a six-digit code that matches none of the three
// steps accepted by validateTOTPCode, so rejection cannot be accidental.
func wrongCodeForSecret(t *testing.T, secret string) string {
	t.Helper()
	valid := map[string]bool{}
	for _, off := range []int64{-1, 0, 1} {
		valid[codeForSecret(t, secret, off)] = true
	}
	for i := 0; i < 1000000; i++ {
		candidate := fmt.Sprintf("%06d", i)
		if !valid[candidate] {
			return candidate
		}
	}
	t.Fatal("no wrong code found")
	return ""
}

// The code generated for the current time step must be accepted.
func TestValidateTOTPCodeAcceptsCurrentStepCode(t *testing.T) {
	secret, err := generateTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	if !validateTOTPCode(secret, codeForSecret(t, secret, 0)) {
		t.Fatal("current-step code rejected")
	}
}

// The truncation modulus must be exactly 10^totpCodeDigits so the numeric
// range of a generated code matches its formatted width.
func TestTOTPCodeModulusMatchesDigitCount(t *testing.T) {
	want := uint32(1)
	for i := 0; i < totpCodeDigits; i++ {
		want *= 10
	}
	if totpCodeModulus != want {
		t.Fatalf("totpCodeModulus = %d, want 10^%d = %d", totpCodeModulus, totpCodeDigits, want)
	}
}

// Codes from one step either side (clock drift tolerance) must be accepted.
func TestValidateTOTPCodeAcceptsAdjacentSteps(t *testing.T) {
	secret, err := generateTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	for _, off := range []int64{-1, 1} {
		if !validateTOTPCode(secret, codeForSecret(t, secret, off)) {
			t.Fatalf("code at step offset %d rejected", off)
		}
	}
}

// A code matching no accepted step must be rejected — the inverted
// comparison accepted any mismatching code, bypassing TOTP entirely.
func TestValidateTOTPCodeRejectsWrongCode(t *testing.T) {
	secret, err := generateTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	if validateTOTPCode(secret, wrongCodeForSecret(t, secret)) {
		t.Fatal("wrong code accepted")
	}
	if validateTOTPCode(secret, "") {
		t.Fatal("empty code accepted")
	}
}

// A secret that is not valid base32 must never validate any code.
func TestValidateTOTPCodeRejectsInvalidSecret(t *testing.T) {
	if validateTOTPCode("not-base32!!", "123456") {
		t.Fatal("code accepted for an undecodable secret")
	}
}

// SOFTWARE_TOKEN_MFA sign-in must verify the TOTP against the enrolled
// secret: a wrong code fails with CodeMismatch and consumes the attempt
// budget, the correct code completes authentication.
func TestSoftwareTokenMfaChallengeVerifiesTOTP(t *testing.T) {
	env := newChallengeTestEnv(t)

	secret, err := generateTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	env.user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{
		Enabled:   true,
		SecretKey: secret,
		Verified:  true,
	}
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}

	sess, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "SOFTWARE_TOKEN_MFA", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	respond := func(code string) (interface{}, error) {
		return env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
			"ClientId":      challengeTestClientID,
			"ChallengeName": "SOFTWARE_TOKEN_MFA",
			"Session":       sess,
			"ChallengeResponses": map[string]interface{}{
				"USERNAME":                "victim",
				"SOFTWARE_TOKEN_MFA_CODE": code,
			},
		}))
	}

	if _, err := respond(wrongCodeForSecret(t, secret)); err == nil {
		t.Fatal("wrong TOTP accepted during SOFTWARE_TOKEN_MFA sign-in")
	}

	resp, err := respond(codeForSecret(t, secret, 0))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := resp.(map[string]interface{})["AuthenticationResult"]; !ok {
		t.Fatalf("expected AuthenticationResult, got %#v", resp)
	}
}

// The pool's per-factor MFA configuration gates both the challenge choice
// and the enable: a factor the pool disabled is neither challenged when
// enrolled nor enablable through the preference path.
func TestPerFactorMfaConfiguration(t *testing.T) {
	pool := &cognitostore.UserPool{
		MfaConfiguration:              "ON",
		MfaConfigurationSoftwareToken: &cognitostore.MfaConfigurationType{Enabled: false},
	}
	user := &cognitostore.User{
		SoftwareTokenMfa: &cognitostore.SoftwareTokenMfaSettings{Enabled: true, Verified: true},
	}
	if got := mfaChallengeFor(pool, user, false); got != "MFA_SETUP" {
		t.Fatalf("challenge %q with the software factor disabled in the pool, want MFA_SETUP", got)
	}
	pool.MfaConfigurationSoftwareToken = nil
	if got := mfaChallengeFor(pool, user, false); got != "SOFTWARE_TOKEN_MFA" {
		t.Fatalf("challenge %q with the per-factor configuration absent, want SOFTWARE_TOKEN_MFA", got)
	}
	pool.MfaConfigurationSms = &cognitostore.SmsMfaConfig{}
	user.SmsMfa = &cognitostore.SmsMfaSettings{Enabled: true}
	if got := mfaChallengeFor(pool, user, false); got != "SOFTWARE_TOKEN_MFA" {
		t.Fatalf("challenge %q with SMS configured without its SmsConfiguration, want the available software factor", got)
	}

	env := newChallengeTestEnv(t)
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.MfaConfiguration = "ON"
		p.MfaConfigurationSoftwareToken = &cognitostore.MfaConfigurationType{Enabled: false}
	})
	env.user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{Enabled: true, Verified: true, SecretKey: "JBSWY3DPEHPK3PXP"}
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}
	if err := validateMFAPrerequisites(map[string]interface{}{
		"SoftwareTokenMfaSettings": map[string]interface{}{"Enabled": true},
	}, env.store, env.user, env.pool); err == nil {
		t.Fatal("enabling software token MFA accepted on a pool that disabled the factor")
	}
}

// VerifySoftwareToken attempts are budgeted on both entry paths: the
// session path shares the challenge-session budget, and the access-token
// path — which carries no session — counts on the registration itself;
// once exhausted, even the correct code is refused.
func TestVerifySoftwareTokenAttemptBudget(t *testing.T) {
	env := newChallengeTestEnv(t)
	secret, err := generateTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	env.user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{SecretKey: secret}
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}
	accessToken, _, _, _, err := env.svc.CreateTokens(env.reqCtx, env.pool.ID, env.user.ID, challengeTestClientID, TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}
	verifyToken := func(code string) error {
		_, verr := env.svc.verifySoftwareTokenCore(context.Background(), env.reqCtx, VerifySoftwareTokenInput{AccessToken: accessToken, UserCode: code})
		return verr
	}

	wrong := wrongCodeForSecret(t, secret)
	for i := 0; i < maxChallengeAttempts; i++ {
		if err := verifyToken(wrong); err == nil {
			t.Fatalf("wrong enrolment code accepted on attempt %d", i)
		}
	}
	if err := verifyToken(codeForSecret(t, secret, 0)); err == nil {
		t.Fatal("correct code accepted after the registration's attempt budget was exhausted")
	}

	sess, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "MFA_SETUP", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	verifySession := func(code string) error {
		_, verr := env.svc.verifySoftwareTokenCore(context.Background(), env.reqCtx, VerifySoftwareTokenInput{Session: sess, UserCode: code})
		return verr
	}
	for i := 0; i < maxChallengeAttempts; i++ {
		if err := verifySession(wrong); err == nil {
			t.Fatalf("wrong enrolment code accepted on session attempt %d", i)
		}
	}
	if err := verifySession(codeForSecret(t, secret, 0)); err == nil {
		t.Fatal("correct code accepted after the session budget was exhausted")
	}
}

// An already-associated registration keeps its secret: a stray
// AssociateSoftwareToken call returns the existing SecretCode instead of
// destroying a working TOTP registration.
func TestAssociateSoftwareTokenKeepsExistingSecret(t *testing.T) {
	env := newChallengeTestEnv(t)
	const existing = "JBSWY3DPEHPK3PXP"
	env.user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{Enabled: true, Verified: true, SecretKey: existing}
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}
	accessToken, _, _, _, err := env.svc.CreateTokens(env.reqCtx, env.pool.ID, env.user.ID, challengeTestClientID, TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := env.svc.associateSoftwareTokenCore(context.Background(), env.reqCtx, AssociateSoftwareTokenInput{AccessToken: accessToken})
	if err != nil {
		t.Fatal(err)
	}
	if got := resp.(map[string]interface{})["SecretCode"]; got != existing {
		t.Fatalf("SecretCode %v, want the existing secret", got)
	}
	user, err := env.store.GetUser(env.pool.ID, "victim")
	if err != nil {
		t.Fatal(err)
	}
	if user.SoftwareTokenMfa == nil || !user.SoftwareTokenMfa.Verified || user.SoftwareTokenMfa.SecretKey != existing {
		t.Fatalf("working registration destroyed: %+v", user.SoftwareTokenMfa)
	}

	// Without an existing registration a fresh secret is issued.
	env.user.SoftwareTokenMfa = nil
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}
	resp, err = env.svc.associateSoftwareTokenCore(context.Background(), env.reqCtx, AssociateSoftwareTokenInput{AccessToken: accessToken})
	if err != nil {
		t.Fatal(err)
	}
	if got := resp.(map[string]interface{})["SecretCode"]; got == "" {
		t.Fatal("fresh association returned no secret")
	}
}

// AdminDeleteSoftwareToken carries no confirmed-status precondition: the
// operation's contract stops at the token's existence.
func TestAdminDeleteSoftwareTokenWithoutConfirmedStatus(t *testing.T) {
	env := newChallengeTestEnv(t)
	env.user.UserStatus = "FORCE_CHANGE_PASSWORD"
	env.user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{Enabled: true, SecretKey: "JBSWY3DPEHPK3PXP"}
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}
	if _, err := env.svc.adminDeleteSoftwareTokenCore(env.reqCtx, AdminDeleteSoftwareTokenInput{UserPoolID: env.pool.ID, Username: "victim"}); err != nil {
		t.Fatalf("AdminDeleteSoftwareToken rejected an unconfirmed user: %v", err)
	}
	user, err := env.store.GetUser(env.pool.ID, "victim")
	if err != nil {
		t.Fatal(err)
	}
	if user.SoftwareTokenMfa != nil {
		t.Fatal("software token registration survived AdminDeleteSoftwareToken")
	}
}
