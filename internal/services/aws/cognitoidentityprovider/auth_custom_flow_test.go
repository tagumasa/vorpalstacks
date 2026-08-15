package cognitoidentityprovider

import (
	"context"
	"testing"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// A DefineAuthChallenge response may only nominate the flow-level challenges
// of the custom authentication flow. MFA challenges — above all MFA_SETUP —
// are service-issued after credential verification: accepting MFA_SETUP here
// would let an unauthenticated caller mint an MFA_SETUP-typed session and
// overwrite a victim's TOTP configuration through AssociateSoftwareToken.
func TestResolveCustomFlowChallengeRejectsMfaSetup(t *testing.T) {
	rejected := []string{
		"MFA_SETUP",
		"SOFTWARE_TOKEN_MFA",
		"SMS_MFA",
		"EMAIL_OTP",
		"NEW_PASSWORD_REQUIRED",
		"SELECT_CHALLENGE",
	}
	for _, name := range rejected {
		if _, err := resolveCustomFlowChallenge(map[string]interface{}{"challengeName": name}); err == nil {
			t.Fatalf("DefineAuthChallenge nomination %q accepted", name)
		}
	}
	// A non-string challengeName is a malformed payload; it is treated as
	// absent and falls back to the CUSTOM_CHALLENGE default rather than
	// minting a typed session from an unusable value.
	if got, err := resolveCustomFlowChallenge(map[string]interface{}{"challengeName": 42}); err != nil || got != "CUSTOM_CHALLENGE" {
		t.Fatalf("non-string challengeName must default to CUSTOM_CHALLENGE, got %q err %v", got, err)
	}

	accepted := map[string]string{
		"CUSTOM_CHALLENGE":  "CUSTOM_CHALLENGE",
		"SRP_A":             "SRP_A",
		"PASSWORD_VERIFIER": "PASSWORD_VERIFIER",
	}
	for name, want := range accepted {
		got, err := resolveCustomFlowChallenge(map[string]interface{}{"challengeName": name})
		if err != nil {
			t.Fatalf("legitimate nomination %q rejected: %v", name, err)
		}
		if got != want {
			t.Fatalf("nomination %q resolved to %q", name, got)
		}
	}

	if got, err := resolveCustomFlowChallenge(nil); err != nil || got != "CUSTOM_CHALLENGE" {
		t.Fatalf("absent Lambda response must default to CUSTOM_CHALLENGE, got %q err %v", got, err)
	}
	if got, err := resolveCustomFlowChallenge(map[string]interface{}{}); err != nil || got != "CUSTOM_CHALLENGE" {
		t.Fatalf("empty Lambda response must default to CUSTOM_CHALLENGE, got %q err %v", got, err)
	}
}

// The un-Lambda'd CUSTOM_AUTH path issues CUSTOM_CHALLENGE without any
// credential exchange; with the MFA_SETUP nomination closed there must be no
// way for that response to carry any other challenge type.
func TestCustomAuthWithoutLambdaIssuesCustomChallenge(t *testing.T) {
	env := newChallengeTestEnv(t)

	env.pool.LambdaConfig = &cognitostore.LambdaConfig{
		DefineAuthChallenge: "arn:aws:lambda:us-east-1:000000000000:function:define",
	}
	if err := env.store.UpdateUserPool(env.pool); err != nil {
		t.Fatal(err)
	}

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "CUSTOM_AUTH",
		"ClientId": challengeTestClientID,
		"Username": "victim",
		// Deliberately no PASSWORD: the flow must still not hand out an
		// MFA-capable challenge session.
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, ok := resp.(map[string]interface{})
	if !ok || m["ChallengeName"] != "CUSTOM_CHALLENGE" {
		t.Fatalf("expected CUSTOM_CHALLENGE, got %#v", resp)
	}
}
