package cognitoidentityprovider

import (
	"context"
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// newCodesTestEnv builds a pool, a client and an unconfirmed user carrying
// both an outstanding sign-up code and an outstanding reset code.
func newCodesTestEnv(t *testing.T) (*CognitoService, *request.RequestContext, string) {
	t.Helper()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := request.NewRequestContext(context.Background(), mgr, "000000000000", "us-east-1")
	svc := NewCognitoService("000000000000", "us-east-1")
	svc.SetStorageManager(mgr)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := store.CreateUserPool(cognitostore.NewUserPool("codespool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   "codes-client",
		ClientName: "codes",
	}); err != nil {
		t.Fatal(err)
	}

	user := cognitostore.NewUser(pool.ID, "codes-user")
	user.UserStatus = "UNCONFIRMED"
	user.SignUpCode = "111222"
	user.SignUpCodeExpiry = time.Now().Add(time.Hour)
	user.PasswordResetCode = "333444"
	user.PasswordResetExpiry = time.Now().Add(time.Hour)
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}
	return svc, reqCtx, pool.ID
}

// Verification codes are purpose-bound: an outstanding reset code cannot
// confirm sign-up, a sign-up code cannot confirm a password reset, and the
// matching code of each flow succeeds on its own flow alone.
func TestVerificationCodesArePurposeBound(t *testing.T) {
	svc, reqCtx, poolID := newCodesTestEnv(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	_, err = svc.confirmSignUpCore(ctx, reqCtx, ConfirmSignUpInput{
		ClientID: "codes-client", Username: "codes-user", ConfirmationCode: "333444",
	})
	if !errors.Is(err, ErrCodeMismatch) {
		t.Fatalf("reset code confirmed sign-up: %v", err)
	}

	_, err = svc.confirmForgotPasswordCore(reqCtx, ConfirmForgotPasswordInput{
		ClientID: "codes-client", Username: "codes-user", Password: "NewPass123!", ConfirmationCode: "111222",
	})
	if !errors.Is(err, ErrCodeMismatch) {
		t.Fatalf("sign-up code confirmed a password reset: %v", err)
	}

	if _, err := svc.confirmSignUpCore(ctx, reqCtx, ConfirmSignUpInput{
		ClientID: "codes-client", Username: "codes-user", ConfirmationCode: "111222",
	}); err != nil {
		t.Fatalf("matching sign-up code rejected: %v", err)
	}
	user, err := store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.UserStatus != "CONFIRMED" {
		t.Fatalf("user status %s after sign-up confirmation", user.UserStatus)
	}
	if user.PasswordResetCode != "333444" {
		t.Fatalf("sign-up confirmation disturbed the reset code: %q", user.PasswordResetCode)
	}

	// The reset flow issues a fresh code and confirms with it alone.
	if _, err := svc.forgotPasswordCore(ctx, reqCtx, ForgotPasswordInput{
		ClientID: "codes-client", Username: "codes-user",
	}); err != nil {
		t.Fatal(err)
	}
	user, err = store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.PasswordResetCode == "" || user.PasswordResetCode == "333444" {
		t.Fatalf("forgot password did not issue a fresh reset code: %q", user.PasswordResetCode)
	}
	if _, err := svc.confirmForgotPasswordCore(reqCtx, ConfirmForgotPasswordInput{
		ClientID: "codes-client", Username: "codes-user", Password: "NewPass123!", ConfirmationCode: user.PasswordResetCode,
	}); err != nil {
		t.Fatalf("matching reset code rejected: %v", err)
	}
}

// ResendConfirmationCode reissues the code that confirms a NEW account; a
// CONFIRMED user has nothing left to confirm and the call is refused without
// minting a code.
func TestResendConfirmationCodeRejectsConfirmedUser(t *testing.T) {
	svc, reqCtx, poolID := newCodesTestEnv(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	user, err := store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	user.UserStatus = "CONFIRMED"
	user.SignUpCode = ""
	if err := store.UpdateUser(user); err != nil {
		t.Fatal(err)
	}

	_, err = svc.resendConfirmationCodeCore(context.Background(), reqCtx, ResendConfirmationCodeInput{
		ClientID: "codes-client", Username: "codes-user",
	})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("resend for a confirmed user returned %v, want InvalidParameterException", err)
	}
	user, err = store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.SignUpCode != "" {
		t.Fatal("refused resend minted a sign-up code")
	}
}

// AdminResetUserPassword issues the reset code its documented follow-up
// ConfirmForgotPassword verifies: a fresh purpose-bound code lands on the
// user, the outstanding sign-up code stays unusable for the reset, and
// confirming with the issued code completes the flow.
func TestAdminResetUserPasswordIssuesPurposeBoundCode(t *testing.T) {
	svc, reqCtx, poolID := newCodesTestEnv(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	user, err := store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	user.UserStatus = "CONFIRMED"
	user.PasswordResetCode = ""
	if err := store.UpdateUser(user); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.adminResetUserPasswordCore(context.Background(), reqCtx, AdminResetUserPasswordInput{
		UserPoolID: poolID, Username: "codes-user",
	}); err != nil {
		t.Fatal(err)
	}
	user, err = store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.UserStatus != "RESET_REQUIRED" {
		t.Fatalf("user status %s after admin reset, want RESET_REQUIRED", user.UserStatus)
	}
	if user.PasswordResetCode == "" {
		t.Fatal("admin reset issued no reset code")
	}
	if time.Now().After(user.PasswordResetExpiry) {
		t.Fatal("admin reset code already expired")
	}

	_, err = svc.confirmForgotPasswordCore(reqCtx, ConfirmForgotPasswordInput{
		ClientID: "codes-client", Username: "codes-user", Password: "NewPass123!", ConfirmationCode: user.SignUpCode,
	})
	if !errors.Is(err, ErrCodeMismatch) {
		t.Fatalf("sign-up code confirmed an admin reset: %v", err)
	}

	if _, err := svc.confirmForgotPasswordCore(reqCtx, ConfirmForgotPasswordInput{
		ClientID: "codes-client", Username: "codes-user", Password: "NewPass123!", ConfirmationCode: user.PasswordResetCode,
	}); err != nil {
		t.Fatalf("confirm with the admin-issued code failed: %v", err)
	}
	user, err = store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.UserStatus != "CONFIRMED" {
		t.Fatalf("user status %s after confirming the admin reset", user.UserStatus)
	}
	if user.PasswordResetCode != "" {
		t.Fatal("confirmation did not consume the reset code")
	}
}

// An administrator reset deactivates imported credentials with everything
// else: the next USER_PASSWORD_AUTH sign-in refuses even the correct
// imported password with PasswordResetRequiredException, and the
// forgot-password flow restores native sign-in.
func TestAdminResetDeactivatesImportedCredentials(t *testing.T) {
	svc, reqCtx, poolID := newCodesTestEnv(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	const importedPassword = "ImportedPass123!"
	importedHash, err := bcrypt.GenerateFromPassword([]byte(importedPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	user.UserStatus = "RESET_REQUIRED"
	user.SignUpCode = ""
	user.PasswordHash = string(importedHash)
	user.PasswordHashAlgo = "BCRYPT"
	if err := store.UpdateUser(user); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.adminResetUserPasswordCore(context.Background(), reqCtx, AdminResetUserPasswordInput{
		UserPoolID: poolID, Username: "codes-user",
	}); err != nil {
		t.Fatal(err)
	}

	// The imported password no longer verifies — the reset refused it before
	// the import-verification branch could run.
	passwordAuth := func(password string) error {
		_, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
			"AuthFlow":       "USER_PASSWORD_AUTH",
			"ClientId":       "codes-client",
			"AuthParameters": map[string]interface{}{"USERNAME": "codes-user", "PASSWORD": password},
		}))
		return err
	}
	if err := passwordAuth(importedPassword); !errors.Is(err, ErrPasswordResetRequired) {
		t.Fatalf("post-reset sign-in with the imported password returned %v, want PasswordResetRequiredException", err)
	}

	// The forgot-password flow completes the reset and native sign-in works.
	reset, err := store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	if reset.PasswordResetCode == "" {
		t.Fatal("admin reset issued no reset code")
	}
	if _, err := svc.confirmForgotPasswordCore(reqCtx, ConfirmForgotPasswordInput{
		ClientID: "codes-client", Username: "codes-user", Password: "NewPass123!", ConfirmationCode: reset.PasswordResetCode,
	}); err != nil {
		t.Fatalf("confirm forgot password: %v", err)
	}
	if err := passwordAuth("NewPass123!"); err != nil {
		t.Fatalf("sign-in after the completed reset failed: %v", err)
	}
}

// An administrator reset deactivates a hash-less imported user's sign-in
// too: a user with no stored credentials of any kind would otherwise stay
// on the CSV-import any-password NEW_PASSWORD_REQUIRED flow, which the
// reset supersedes — every post-reset sign-in answers
// PasswordResetRequiredException until the forgot-password flow completes
// and native sign-in works again.
func TestAdminResetGatesHashlessImportedSignIn(t *testing.T) {
	svc, reqCtx, poolID := newCodesTestEnv(t)
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}

	user, err := store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	user.UserStatus = "RESET_REQUIRED"
	user.SignUpCode = ""
	user.PasswordHash = ""
	user.PasswordHashAlgo = ""
	if err := store.UpdateUser(user); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.adminResetUserPasswordCore(context.Background(), reqCtx, AdminResetUserPasswordInput{
		UserPoolID: poolID, Username: "codes-user",
	}); err != nil {
		t.Fatal(err)
	}

	// Any password — the import flow would have accepted it and issued the
	// NEW_PASSWORD_REQUIRED challenge — must answer the reset refusal
	// instead.
	passwordAuth := func(password string) error {
		_, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
			"AuthFlow":       "USER_PASSWORD_AUTH",
			"ClientId":       "codes-client",
			"AuthParameters": map[string]interface{}{"USERNAME": "codes-user", "PASSWORD": password},
		}))
		return err
	}
	if err := passwordAuth("AnyPassword123!"); !errors.Is(err, ErrPasswordResetRequired) {
		t.Fatalf("post-reset sign-in with an arbitrary password returned %v, want PasswordResetRequiredException", err)
	}

	// The forgot-password flow completes the reset and native sign-in works.
	reset, err := store.GetUser(poolID, "codes-user")
	if err != nil {
		t.Fatal(err)
	}
	if reset.PasswordResetCode == "" {
		t.Fatal("admin reset issued no reset code")
	}
	if _, err := svc.confirmForgotPasswordCore(reqCtx, ConfirmForgotPasswordInput{
		ClientID: "codes-client", Username: "codes-user", Password: "NewPass123!", ConfirmationCode: reset.PasswordResetCode,
	}); err != nil {
		t.Fatalf("confirm forgot password: %v", err)
	}
	if err := passwordAuth("NewPass123!"); err != nil {
		t.Fatalf("sign-in after the completed reset failed: %v", err)
	}
}
