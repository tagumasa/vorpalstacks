package cognitoidentityprovider

import (
	"context"
	"errors"
	"testing"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// A pool whose email is both an alias attribute and auto-verified implements
// the staged activation model: sign-up accepts a duplicate value (Amazon
// Cognito documents this as preventing account enumeration through the
// public sign-up API), and the second confirmation is the documented
// AliasExistsException point.
func TestSignUpStagedAliasActivation(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool := cognitostore.NewUserPool("alias-stage-pool", "us-east-1")
	pool.AliasAttributes = []string{"email"}
	pool.AutoVerifiedAttributes = []string{"email"}
	if _, err := store.CreateUserPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID, ClientID: "alias-client", ClientName: "alias",
	}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	ctx := context.Background()
	shared := "staged@example.com"

	for _, name := range []string{"staged-first", "staged-second"} {
		if _, err := svc.signUpCore(ctx, reqCtx, SignUpInput{
			ClientID: "alias-client", Username: name, Password: "StagedPass123!",
			UserAttributes: map[string]string{"email": shared},
		}); err != nil {
			t.Fatalf("sign-up %s with a duplicate email must succeed: %v", name, err)
		}
	}

	first, err := store.GetUser(pool.ID, "staged-first")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.confirmSignUpCore(ctx, reqCtx, ConfirmSignUpInput{
		ClientID: "alias-client", Username: "staged-first", ConfirmationCode: first.SignUpCode,
	}); err != nil {
		t.Fatalf("confirm the first holder: %v", err)
	}

	second, err := store.GetUser(pool.ID, "staged-second")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.confirmSignUpCore(ctx, reqCtx, ConfirmSignUpInput{
		ClientID: "alias-client", Username: "staged-second", ConfirmationCode: second.SignUpCode,
	}); !errors.Is(err, ErrAliasExists) {
		t.Fatalf("confirming the duplicate returned %v, want AliasExistsException", err)
	}
}

// AdminCreateUser activates a verified alias at creation: the conflict is
// rejected without ForceAliasCreation, and the parameter migrates the claim
// to the new user — the previous holder keeps the value but can no longer
// sign in with it.
func TestAdminCreateUserForceAliasCreation(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool := cognitostore.NewUserPool("alias-admin-pool", "us-east-1")
	pool.AliasAttributes = []string{"email"}
	if _, err := store.CreateUserPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	ctx := context.Background()
	shared := "forced@example.com"
	verifiedAttrs := func() map[string]string {
		return map[string]string{"email": shared, "email_verified": "true"}
	}

	if _, err := svc.adminCreateUserCore(ctx, reqCtx, AdminCreateUserInput{
		UserPoolID: pool.ID, Username: "alias-original", MessageAction: "SUPPRESS",
		UserAttributes: verifiedAttrs(),
	}); err != nil {
		t.Fatalf("create the original holder: %v", err)
	}

	if _, err := svc.adminCreateUserCore(ctx, reqCtx, AdminCreateUserInput{
		UserPoolID: pool.ID, Username: "alias-challenger", MessageAction: "SUPPRESS",
		UserAttributes: verifiedAttrs(),
	}); !errors.Is(err, ErrAliasExists) {
		t.Fatalf("duplicate verified alias returned %v, want AliasExistsException", err)
	}

	if _, err := svc.adminCreateUserCore(ctx, reqCtx, AdminCreateUserInput{
		UserPoolID: pool.ID, Username: "alias-taker", MessageAction: "SUPPRESS",
		ForceAliasCreation: true, UserAttributes: verifiedAttrs(),
	}); err != nil {
		t.Fatalf("forced create: %v", err)
	}
	original, err := store.GetUser(pool.ID, "alias-original")
	if err != nil {
		t.Fatal(err)
	}
	if original.Attributes["email"] != shared || original.Attributes["email_verified"] != "false" {
		t.Fatalf("previous holder after migration = %v", original.Attributes)
	}
	if resolved, err := store.GetUser(pool.ID, shared); err != nil || resolved.Username != "alias-taker" {
		t.Fatalf("alias lookup after migration = (%v, %v), want alias-taker", resolved, err)
	}

	// Re-activating the previous holder's claim hits the migrated owner.
	if _, err := svc.adminUpdateUserAttributesCore(ctx, reqCtx, AdminUpdateUserAttributesInput{
		UserPoolID: pool.ID, Username: "alias-original",
		UserAttributes: map[string]string{"email_verified": "true"},
	}); !errors.Is(err, ErrAliasExists) {
		t.Fatalf("re-activating a migrated alias returned %v, want AliasExistsException", err)
	}
}

// Verifying a duplicate email through the attribute-verification round
// claims the value and is rejected — the model lists
// AliasExistsException on VerifyUserAttribute.
func TestVerifyUserAttributeAliasConflict(t *testing.T) {
	svc, reqCtx, store := newContractTestService(t)
	pool := cognitostore.NewUserPool("alias-verify-pool", "us-east-1")
	pool.AliasAttributes = []string{"email"}
	if _, err := store.CreateUserPool(pool); err != nil {
		t.Fatalf("create pool: %v", err)
	}
	// Token issuance reads the client record and fails closed when it is
	// missing, so the fixture registers the app client it issues against.
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID, ClientID: "verify-client", ClientName: "verify",
	}); err != nil {
		t.Fatalf("create client: %v", err)
	}
	shared := "verify@example.com"

	holder := cognitostore.NewUser(pool.ID, "verify-holder")
	holder.UserStatus = "CONFIRMED"
	holder.Attributes = map[string]string{"sub": holder.ID, "email": shared, "email_verified": "true"}
	if err := store.CreateUser(holder); err != nil {
		t.Fatal(err)
	}

	duplicate := cognitostore.NewUser(pool.ID, "verify-duplicate")
	duplicate.UserStatus = "CONFIRMED"
	duplicate.Attributes = map[string]string{"sub": duplicate.ID, "email": shared}
	if err := store.CreateUser(duplicate); err != nil {
		t.Fatal(err)
	}
	accessToken, _, _, _, err := svc.CreateTokens(reqCtx, pool.ID, duplicate.ID, "verify-client", TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatalf("create tokens: %v", err)
	}

	if _, err := svc.getUserAttributeVerificationCodeCore(context.Background(), reqCtx, GetUserAttributeVerificationCodeInput{
		AccessToken: accessToken, AttributeName: "email",
	}); err != nil {
		t.Fatalf("issue verification code: %v", err)
	}
	stored, err := store.GetUser(pool.ID, "verify-duplicate")
	if err != nil {
		t.Fatal(err)
	}
	codeEntry := stored.AttributeVerificationCodes["email"]
	if codeEntry == nil || codeEntry.Code == "" {
		t.Fatal("verification code was not minted")
	}

	if _, err := svc.verifyUserAttributeCore(reqCtx, VerifyUserAttributeInput{
		AccessToken: accessToken, AttributeName: "email", Code: codeEntry.Code,
	}); !errors.Is(err, ErrAliasExists) {
		t.Fatalf("verifying a duplicate alias returned %v, want AliasExistsException", err)
	}
}
