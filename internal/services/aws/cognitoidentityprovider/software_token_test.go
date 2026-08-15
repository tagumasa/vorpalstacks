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
			"ClientId":                challengeTestClientID,
			"ChallengeName":           "SOFTWARE_TOKEN_MFA",
			"Session":                 sess,
			"USERNAME":                "victim",
			"SOFTWARE_TOKEN_MFA_CODE": code,
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
