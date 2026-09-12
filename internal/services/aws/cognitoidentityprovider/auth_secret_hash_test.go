package cognitoidentityprovider

import (
	"context"
	"errors"
	"testing"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// A secret-bearing app client proves itself with SECRET_HASH on the
// authentication planes: the hash is required in the auth parameters of
// every InitiateAuth flow (admin plane included), in every challenge
// response, and on the refresh grant; a wrong or missing proof is
// NotAuthorizedException, while a no-secret client stays exempt.
func TestAuthPlaneSecretHashEnforcement(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("secret-hash-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	const secretClientID = "secret-hash-client"
	const secret = "unit-test-client-secret-0123456789abcdef"
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID:   pool.ID,
		ClientID:     secretClientID,
		ClientName:   "secret-hash",
		ClientSecret: secret,
	}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	user := cognitostore.NewUser(pool.ID, "sh-user")
	user.UserStatus = "CONFIRMED"
	user.Attributes = map[string]string{"sub": user.ID, "email": "sh@example.com"}
	hash, err := bcrypt.GenerateFromPassword([]byte("Sup3rSecret!"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	user.PasswordHash = string(hash)
	if err := store.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	passwordAuth := func(secretHash string) (interface{}, error) {
		params := map[string]interface{}{
			"USERNAME": "sh-user",
			"PASSWORD": "Sup3rSecret!",
		}
		if secretHash != "" {
			params["SECRET_HASH"] = secretHash
		}
		return svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
			"AuthFlow":       "USER_PASSWORD_AUTH",
			"ClientId":       secretClientID,
			"AuthParameters": params,
		}))
	}

	// The correct proof completes the authentication.
	resp, err := passwordAuth(secretHashFor(secret, "sh-user", secretClientID))
	if err != nil {
		t.Fatalf("correct SECRET_HASH rejected: %v", err)
	}
	result, ok := resp.(map[string]interface{})
	if !ok || result["AuthenticationResult"] == nil {
		t.Fatalf("USER_PASSWORD_AUTH returned no AuthenticationResult: %v", resp)
	}

	// A wrong proof and a missing proof are both NotAuthorizedException.
	if _, err := passwordAuth(secretHashFor(secret, "sh-user", "other-client")); !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("wrong SECRET_HASH returned %v, want NotAuthorized", err)
	}
	if _, err := passwordAuth(""); !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("missing SECRET_HASH returned %v, want NotAuthorized", err)
	}

	// The admin plane enforces the same proof on its AuthParameters.
	if _, err := svc.adminInitiateAuthCore(context.Background(), reqCtx, AdminInitiateAuthInput{
		UserPoolID: pool.ID,
		ClientID:   secretClientID,
		AuthFlow:   "ADMIN_USER_PASSWORD_AUTH",
		Params: map[string]interface{}{
			"AuthParameters": map[string]interface{}{
				"USERNAME": "sh-user",
				"PASSWORD": "Sup3rSecret!",
			},
		},
	}); !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("admin flow without SECRET_HASH returned %v, want NotAuthorized", err)
	}

	// Every challenge response carries the proof: the gate fires before the
	// session is even resolved.
	if _, err := svc.RespondToAuthChallenge(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"ClientId":           secretClientID,
		"ChallengeName":      "SMS_MFA",
		"Session":            "any-session",
		"ChallengeResponses": map[string]interface{}{"USERNAME": "sh-user", "SMS_MFA_CODE": "123456"},
	})); !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("challenge response without SECRET_HASH returned %v, want NotAuthorized", err)
	}

	// The refresh grant takes the proof against the token-resolved user:
	// the username-claim derivation verifies, a wrong hash does not.
	authResult, ok := result["AuthenticationResult"].(map[string]interface{})
	if !ok {
		t.Fatalf("AuthenticationResult shape: %v", result["AuthenticationResult"])
	}
	refreshToken, _ := authResult["RefreshToken"].(string)
	if refreshToken == "" {
		t.Fatal("USER_PASSWORD_AUTH issued no refresh token")
	}
	refresh := func(secretHash string) error {
		params := map[string]interface{}{"REFRESH_TOKEN": refreshToken}
		if secretHash != "" {
			params["SECRET_HASH"] = secretHash
		}
		_, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
			"AuthFlow":       "REFRESH_TOKEN_AUTH",
			"ClientId":       secretClientID,
			"AuthParameters": params,
		}))
		return err
	}
	if err := refresh(secretHashFor(secret, "sh-user", secretClientID)); err != nil {
		t.Fatalf("refresh with correct SECRET_HASH failed: %v", err)
	}
	if err := refresh("not-the-hash"); !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("refresh with wrong SECRET_HASH returned %v, want NotAuthorized", err)
	}
	if err := refresh(""); !errors.Is(err, ErrNotAuthorized) {
		t.Fatalf("refresh without SECRET_HASH returned %v, want NotAuthorized", err)
	}
}
