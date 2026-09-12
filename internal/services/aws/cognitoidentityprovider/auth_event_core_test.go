package cognitoidentityprovider

import (
	"context"
	"testing"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// authEventTypes collects the EventType values recorded for a user, walking
// every page of the event history.
func authEventTypes(t *testing.T, svc *CognitoService, reqCtx *request.RequestContext, poolID, username string) map[string]string {
	t.Helper()
	events := map[string]string{}
	next := ""
	for {
		in := AdminListUserAuthEventsInput{UserPoolID: poolID, Username: username}
		if next != "" {
			in.NextToken = next
		}
		resp, err := svc.adminListUserAuthEventsCore(reqCtx, in)
		if err != nil {
			t.Fatalf("list auth events: %v", err)
		}
		for _, raw := range resp.(map[string]interface{})["AuthEvents"].([]map[string]interface{}) {
			events[raw["EventType"].(string)] = raw["EventResponse"].(string)
		}
		nt, _ := resp.(map[string]interface{})["NextToken"].(string)
		if nt == "" {
			return events
		}
		next = nt
	}
}

// The account-lifecycle operations record their completed outcome in the
// user's event history: self-service sign-up (SignUp), password-reset code
// issuance (ForgotPassword), an access-token password change
// (PasswordChange) and a confirmation-code resend (ResendCode) each leave a
// Pass event, so AdminListUserAuthEvents reports the full Smithy EventType
// inventory and not the sign-in family alone.
func TestLifecycleOperationsRecordAuthEvents(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("events-pool", "us-east-1"))
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID, ClientID: "events-client", ClientName: "events",
	}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	ctx := context.Background()

	// Self-service sign-up: an UNCONFIRMED user with an outstanding code.
	if _, err := svc.signUpCore(ctx, reqCtx, SignUpInput{
		ClientID: "events-client", Username: "events-user", Password: "EventsPass1!",
	}); err != nil {
		t.Fatalf("sign up: %v", err)
	}

	// Password-reset code issuance.
	if _, err := svc.forgotPasswordCore(ctx, reqCtx, ForgotPasswordInput{
		ClientID: "events-client", Username: "events-user",
	}); err != nil {
		t.Fatalf("forgot password: %v", err)
	}

	// Confirmation-code resend (the sign-up user is still unconfirmed).
	if _, err := svc.resendConfirmationCodeCore(ctx, reqCtx, ResendConfirmationCodeInput{
		ClientID: "events-client", Username: "events-user",
	}); err != nil {
		t.Fatalf("resend confirmation code: %v", err)
	}

	// Access-token password change on a separate confirmed user.
	confirmed := cognitostore.NewUser(pool.ID, "events-confirmed")
	confirmed.UserStatus = "CONFIRMED"
	confirmed.Attributes = map[string]string{"sub": confirmed.ID}
	if err := setNativePasswordCredentials(confirmed, nil, "OldEvents1!"); err != nil {
		t.Fatalf("seed credentials: %v", err)
	}
	if err := store.CreateUser(confirmed); err != nil {
		t.Fatalf("create confirmed user: %v", err)
	}
	accessToken, _, _, _, err := svc.CreateTokens(reqCtx, pool.ID, confirmed.ID, "events-client", TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatalf("create tokens: %v", err)
	}
	if _, err := svc.changePasswordCore(reqCtx, ChangePasswordInput{
		AccessToken: accessToken, PreviousPassword: "OldEvents1!", NewPassword: "NewEvents1!",
	}); err != nil {
		t.Fatalf("change password: %v", err)
	}

	events := authEventTypes(t, svc, reqCtx, pool.ID, "events-user")
	for _, eventType := range []string{"SignUp", "ForgotPassword", "ResendCode"} {
		if response, ok := events[eventType]; !ok || response != "Pass" {
			t.Fatalf("event %s = (%q, %v), want a recorded Pass", eventType, response, ok)
		}
	}

	confirmedEvents := authEventTypes(t, svc, reqCtx, pool.ID, "events-confirmed")
	if response, ok := confirmedEvents["PasswordChange"]; !ok || response != "Pass" {
		t.Fatalf("PasswordChange event = (%q, %v), want a recorded Pass", response, ok)
	}
}
