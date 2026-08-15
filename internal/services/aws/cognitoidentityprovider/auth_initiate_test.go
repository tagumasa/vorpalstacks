package cognitoidentityprovider

import (
	"context"
	"testing"
	"time"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

func initiateUserAuth(env *challengeTestEnv, username string) (interface{}, error) {
	return env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_AUTH",
		"ClientId": challengeTestClientID,
		"Username": username,
	}))
}

// USER_AUTH issues the SELECT_CHALLENGE session before any credential
// verification, so the selector must only offer sign-in challenges. MFA
// challenges and MFA_SETUP are server-issued; accepting them here would let
// an unauthenticated caller mint an MFA_SETUP session and overwrite a
// victim's TOTP configuration through AssociateSoftwareToken.
func TestUserAuthSelectChallengeRejectsMfaChallenges(t *testing.T) {
	env := newChallengeTestEnv(t)

	env.user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{
		Enabled:   true,
		SecretKey: "SEEDSECRET",
		Verified:  true,
	}
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}

	resp, err := initiateUserAuth(env, "victim")
	if err != nil {
		t.Fatal(err)
	}
	selector, _ := resp.(map[string]interface{})["Session"].(string)
	if selector == "" {
		t.Fatalf("expected selector session, got %#v", resp)
	}

	// Rejected selections must not consume the selector session.
	for _, selected := range []string{"MFA_SETUP", "SOFTWARE_TOKEN_MFA", "SMS_MFA", "CUSTOM_CHALLENGE"} {
		_, err := env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
			"ClientId":           challengeTestClientID,
			"ChallengeName":      "SELECT_CHALLENGE",
			"Session":            selector,
			"USERNAME":           "victim",
			"SELECTED_CHALLENGE": selected,
		}))
		if err == nil {
			t.Fatalf("SELECT_CHALLENGE accepted non-selectable challenge %q", selected)
		}
	}

	resp, err = env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"ClientId":           challengeTestClientID,
		"ChallengeName":      "SELECT_CHALLENGE",
		"Session":            selector,
		"USERNAME":           "victim",
		"SELECTED_CHALLENGE": "PASSWORD",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if resp.(map[string]interface{})["ChallengeName"] != "PASSWORD" {
		t.Fatalf("expected PASSWORD session mint, got %#v", resp)
	}

	user, err := env.store.GetUser(env.pool.ID, "victim")
	if err != nil {
		t.Fatal(err)
	}
	if user.SoftwareTokenMfa == nil || user.SoftwareTokenMfa.SecretKey != "SEEDSECRET" || !user.SoftwareTokenMfa.Verified {
		t.Fatal("victim TOTP configuration disturbed via SELECT_CHALLENGE")
	}
}

// A FORCE_CHANGE_PASSWORD user must authenticate with the temporary
// password before NEW_PASSWORD_REQUIRED is issued; issuing the challenge on
// user status alone would let anyone take over the account by username only.
func TestPasswordAuthForceChangePasswordRequiresPassword(t *testing.T) {
	env := newChallengeTestEnv(t)

	hash, err := bcrypt.GenerateFromPassword([]byte("TempPass123!"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	env.user.UserStatus = "FORCE_CHANGE_PASSWORD"
	env.user.PasswordHash = string(hash)
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}

	initiate := func(password string) (interface{}, error) {
		return env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
			"AuthFlow": "USER_PASSWORD_AUTH",
			"ClientId": challengeTestClientID,
			"AuthParameters": map[string]interface{}{
				"USERNAME": "victim",
				"PASSWORD": password,
			},
		}))
	}

	if _, err := initiate("WrongPass999!"); err == nil {
		t.Fatal("NEW_PASSWORD_REQUIRED issued without verifying the temporary password")
	}

	resp, err := initiate("TempPass123!")
	if err != nil {
		t.Fatal(err)
	}
	if m, ok := resp.(map[string]interface{}); !ok || m["ChallengeName"] != "NEW_PASSWORD_REQUIRED" {
		t.Fatalf("expected NEW_PASSWORD_REQUIRED after correct temporary password, got %#v", resp)
	}
}

// SELECT_MFA_TYPE keeps its own answer set: SMS_MFA and SOFTWARE_TOKEN_MFA.
func TestSelectMfaTypeAcceptsMfaChoices(t *testing.T) {
	env := newChallengeTestEnv(t)

	env.user.SoftwareTokenMfa = &cognitostore.SoftwareTokenMfaSettings{
		Enabled:   true,
		SecretKey: "SEEDSECRET",
		Verified:  true,
	}
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}

	sel, err := mintChallengeSession(env.store, env.pool.ID, challengeTestClientID, "victim", "SELECT_MFA_TYPE", 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"ClientId":      challengeTestClientID,
		"ChallengeName": "SELECT_MFA_TYPE",
		"Session":       sel,
		"USERNAME":      "victim",
		"MFA_TYPE":      "SOFTWARE_TOKEN_MFA",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if resp.(map[string]interface{})["ChallengeName"] != "SOFTWARE_TOKEN_MFA" {
		t.Fatalf("expected SOFTWARE_TOKEN_MFA session mint, got %#v", resp)
	}
}
