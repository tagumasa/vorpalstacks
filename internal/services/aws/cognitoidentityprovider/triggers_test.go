package cognitoidentityprovider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// stubTriggerInvoker is a LambdaInvoker stand-in that replays a canned
// invocation outcome so the trigger error-classification paths can be
// exercised without a real Lambda runtime.
type stubTriggerInvoker struct {
	functionError string
	payload       []byte
	invokeErr     error
}

func (s *stubTriggerInvoker) InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error) {
	return 200, s.payload, s.invokeErr
}

func (s *stubTriggerInvoker) InvokeForTrigger(ctx context.Context, functionName string, payload []byte) (invokers.LambdaInvocation, error) {
	return invokers.LambdaInvocation{
		StatusCode:    200,
		Payload:       s.payload,
		FunctionError: s.functionError,
	}, s.invokeErr
}

func (s *stubTriggerInvoker) GetFunctionARN(ctx context.Context, functionName string) (string, error) {
	return functionName, nil
}

func newTriggerTestService(t *testing.T, invoker invokers.LambdaInvoker) *CognitoService {
	t.Helper()
	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("failed to start event bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	bus.SetLambdaInvoker(invoker)
	svc := NewCognitoService("123456789012", "us-east-1")
	svc.bus = bus
	if _, err := eventbus.SubscribeTyped[*eventbus.CognitoTriggerEvent](bus, svc.handleCognitoTrigger); err != nil {
		t.Fatalf("failed to subscribe trigger handler: %v", err)
	}
	return svc
}

const triggerTestARN = "arn:aws:lambda:us-east-1:123456789012:function:migration"

// A Lambda function that raises an error must fail the migration with
// NotAuthorizedException semantics (an incorrect-password outcome), not an
// internal error.
func TestInvokeTriggerFunctionErrorClassifiedAsIncorrectPassword(t *testing.T) {
	svc := newTriggerTestService(t, &stubTriggerInvoker{functionError: "Unhandled"})

	_, err := svc.invokeTrigger(
		context.Background(),
		UserMigrationAuthentication, "us-east-1_POOL", "user", "client",
		triggerTestARN,
		map[string]interface{}{"password": "wrong"},
		map[string]interface{}{"finalUserStatus": "CONFIRMED"},
		true,
	)
	if err == nil {
		t.Fatal("expected an error from a failed trigger function")
	}
	var fnErr *lambdaFunctionError
	if !errors.As(err, &fnErr) {
		t.Fatalf("expected a lambdaFunctionError, got %T: %v", err, err)
	}
	if got := classifyMigrationFailure(err); got != ErrIncorrectPassword {
		t.Fatalf("expected ErrIncorrectPassword, got %v", got)
	}
}

// An invocation-transport failure is an infrastructure error, not a
// wrong-password outcome.
func TestInvokeTriggerTransportFailureClassifiedAsInternalError(t *testing.T) {
	svc := newTriggerTestService(t, &stubTriggerInvoker{invokeErr: errors.New("function not found")})

	_, err := svc.invokeTrigger(
		context.Background(),
		UserMigrationAuthentication, "us-east-1_POOL", "user", "client",
		triggerTestARN,
		map[string]interface{}{"password": "secret"},
		map[string]interface{}{"finalUserStatus": "CONFIRMED"},
		true,
	)
	if err == nil {
		t.Fatal("expected an error from a failed invocation")
	}
	if got := classifyMigrationFailure(err); got != ErrInternalError {
		t.Fatalf("expected ErrInternalError, got %v", got)
	}
}

// A successful invocation returns the unmarshalled Lambda response.
func TestInvokeTriggerSuccessReturnsResponse(t *testing.T) {
	svc := newTriggerTestService(t, &stubTriggerInvoker{payload: []byte(`{"finalUserStatus":"RESET_REQUIRED"}`)})

	resp, err := svc.invokeTrigger(
		context.Background(),
		UserMigrationAuthentication, "us-east-1_POOL", "user", "client",
		triggerTestARN,
		map[string]interface{}{"password": "secret"},
		map[string]interface{}{"finalUserStatus": "CONFIRMED"},
		true,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp["finalUserStatus"] != "RESET_REQUIRED" {
		t.Fatalf("expected the Lambda-provided finalUserStatus, got %v", resp["finalUserStatus"])
	}
}

// triggerInvocation records one Lambda invocation the dispatch stub served.
type triggerInvocation struct {
	Source string
	Code   string
}

// dispatchTriggerInvoker routes each trigger invocation by its payload's
// triggerSource: a configured outcome is replayed, anything else succeeds
// with an empty payload (the caller's responseDefaults), and every
// invocation is recorded so tests can assert which sources fired with which
// code parameter.
type dispatchTriggerInvoker struct {
	outcomes    map[string]invokers.LambdaInvocation
	invocations []triggerInvocation
}

func (d *dispatchTriggerInvoker) InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error) {
	return 200, nil, nil
}

func (d *dispatchTriggerInvoker) InvokeForTrigger(ctx context.Context, functionName string, payload []byte) (invokers.LambdaInvocation, error) {
	var body struct {
		TriggerSource string `json:"triggerSource"`
		Request       struct {
			CodeParameter string `json:"codeParameter"`
		} `json:"request"`
	}
	_ = json.Unmarshal(payload, &body)
	d.invocations = append(d.invocations, triggerInvocation{Source: body.TriggerSource, Code: body.Request.CodeParameter})
	if outcome, ok := d.outcomes[body.TriggerSource]; ok {
		return outcome, nil
	}
	return invokers.LambdaInvocation{StatusCode: 200}, nil
}

func (d *dispatchTriggerInvoker) GetFunctionARN(ctx context.Context, functionName string) (string, error) {
	return functionName, nil
}

// fired reports whether a trigger source was invoked, returning its recorded
// code parameter.
func (d *dispatchTriggerInvoker) fired(source string) (string, bool) {
	for _, inv := range d.invocations {
		if inv.Source == source {
			return inv.Code, true
		}
	}
	return "", false
}

// newTriggerServiceEnv builds a service wired to both the dispatch stub and
// a throwaway storage manager, plus a pool/client pair carrying the given
// Lambda configuration.
func newTriggerServiceEnv(t *testing.T, invoker *dispatchTriggerInvoker, config *cognitostore.LambdaConfig) (*CognitoService, *request.RequestContext, string, *dispatchTriggerInvoker) {
	t.Helper()
	svc := newTriggerTestService(t, invoker)
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc.SetStorageManager(mgr)
	reqCtx := request.NewRequestContext(context.Background(), mgr, "123456789012", "us-east-1")
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	pool, err := store.CreateUserPool(cognitostore.NewUserPool("triggerpool", "us-east-1"))
	if err != nil {
		t.Fatal(err)
	}
	if config != nil {
		if err := store.UpdateUserPoolFunc(pool.ID, func(p *cognitostore.UserPool) error {
			p.LambdaConfig = config
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.CreateUserPoolClient(&cognitostore.UserPoolClient{
		UserPoolID: pool.ID,
		ClientID:   "trigger-client",
		ClientName: "trigger",
	}); err != nil {
		t.Fatal(err)
	}
	return svc, reqCtx, pool.ID, invoker
}

const preAuthTestARN = "arn:aws:lambda:us-east-1:123456789012:function:pre-auth"

// A PreAuthentication Lambda that raises an error denies the sign-in with
// UserLambdaValidationException; a transport failure is an internal error.
func TestPreAuthenticationFailuresFailSignIn(t *testing.T) {
	config := &cognitostore.LambdaConfig{PreAuthentication: preAuthTestARN}

	svc := newTriggerTestService(t, &stubTriggerInvoker{functionError: "Unhandled"})
	if err := invokePreAuthentication(context.Background(), svc, "us-east-1_POOL", "user", "client", config, nil, nil); err != ErrUserLambdaValidation {
		t.Fatalf("raised function error returned %v, want UserLambdaValidationException", err)
	}

	svc = newTriggerTestService(t, &stubTriggerInvoker{invokeErr: errors.New("function not found")})
	if err := invokePreAuthentication(context.Background(), svc, "us-east-1_POOL", "user", "client", config, nil, nil); err != ErrInternalError {
		t.Fatalf("transport failure returned %v, want InternalErrorException", err)
	}
}

// A PostAuthentication Lambda failure fails the completed authentication —
// the documented contract its own comment states.
func TestPostAuthenticationFailureFailsSignIn(t *testing.T) {
	config := &cognitostore.LambdaConfig{PostAuthentication: preAuthTestARN}

	svc := newTriggerTestService(t, &stubTriggerInvoker{functionError: "Unhandled"})
	if err := invokePostAuthentication(context.Background(), svc, "us-east-1_POOL", "user", "client", config, nil, nil); err != ErrUserLambdaValidation {
		t.Fatalf("raised function error returned %v, want UserLambdaValidationException", err)
	}

	svc = newTriggerTestService(t, &stubTriggerInvoker{invokeErr: errors.New("function not found")})
	if err := invokePostAuthentication(context.Background(), svc, "us-east-1_POOL", "user", "client", config, nil, nil); err != ErrInternalError {
		t.Fatalf("transport failure returned %v, want InternalErrorException", err)
	}
}

// The CUSTOM_AUTH flow fails closed when DefineAuthChallenge raises: the
// challenge selection is the Lambda's decision, so no challenge session is
// minted behind its back.
func TestCustomAuthFlowFailsClosedOnTriggerError(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		DefineAuthChallenge: {StatusCode: 200, FunctionError: "Unhandled"},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{DefineAuthChallenge: preAuthTestARN})
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUser(cognitostore.NewUser(poolID, "custom-user")); err != nil {
		t.Fatal(err)
	}

	_, err = svc.customAuthFlow(context.Background(), reqCtx, "trigger-client", "custom-user", "")
	if err != ErrUserLambdaValidation {
		t.Fatalf("custom auth with a raising DefineAuthChallenge returned %v, want UserLambdaValidationException", err)
	}
}

// The PreSignUp response's auto-verify flags mark the matching attributes
// verified at creation.
func TestPreSignUpAutoVerifyMarksAttributes(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		PreSignUpSignUp: {StatusCode: 200, Payload: []byte(`{"autoConfirmUser":false,"autoVerifyEmail":true,"autoVerifyPhone":true}`)},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{PreSignUp: preAuthTestARN})

	if _, err := svc.signUpCore(context.Background(), reqCtx, SignUpInput{
		ClientID: "trigger-client",
		Username: "signup-user",
		Password: "NewPass123!",
		UserAttributes: map[string]string{
			"email":        "signup@example.com",
			"phone_number": "+12065550100",
		},
	}); err != nil {
		t.Fatal(err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.GetUser(poolID, "signup-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.Attributes["email_verified"] != "true" {
		t.Fatalf("autoVerifyEmail did not mark email verified: %v", user.Attributes)
	}
	if user.Attributes["phone_number_verified"] != "true" {
		t.Fatalf("autoVerifyPhone did not mark phone verified: %v", user.Attributes)
	}
}

// The PreTokenGeneration response's role overrides carry into the token's
// cognito:roles and cognito:preferred_role claims.
func TestPreTokenGenerationRoleClaimsApplied(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		TokenGenerationAuthentication: {StatusCode: 200, Payload: []byte(`{"claimsOverrideDetails":{"groupOverrideDetails":{"iamRolesToOverride":["arn:aws:iam::123456789012:role/first","arn:aws:iam::123456789012:role/second"],"preferredRole":"arn:aws:iam::123456789012:role/preferred"}}}`)},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{PreTokenGeneration: preAuthTestARN})
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUser(cognitostore.NewUser(poolID, "claims-user")); err != nil {
		t.Fatal(err)
	}

	accessToken, _, _, _, err := svc.CreateTokens(reqCtx, poolID, storeUserID(t, store, poolID, "claims-user"), "trigger-client", TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}
	claims := jwtPayloadClaims(t, accessToken)
	if claims["cognito:roles"] != "arn:aws:iam::123456789012:role/first,arn:aws:iam::123456789012:role/second" {
		t.Fatalf("cognito:roles claim = %v", claims["cognito:roles"])
	}
	if claims["cognito:preferred_role"] != "arn:aws:iam::123456789012:role/preferred" {
		t.Fatalf("cognito:preferred_role claim = %v", claims["cognito:preferred_role"])
	}
}

// ResendConfirmationCode fires the CustomMessage ResendCode source with the
// reissued code as the code parameter.
func TestResendConfirmationCodeFiresCustomMessage(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		CustomMessageSignUp: {StatusCode: 200},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{CustomMessage: preAuthTestARN})
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateUser(cognitostore.NewUser(poolID, "resend-user")); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.resendConfirmationCodeCore(context.Background(), reqCtx, ResendConfirmationCodeInput{
		ClientID: "trigger-client", Username: "resend-user",
	}); err != nil {
		t.Fatal(err)
	}
	user, err := store.GetUser(poolID, "resend-user")
	if err != nil {
		t.Fatal(err)
	}
	code, fired := invoker.fired(CustomMessageResendCode)
	if !fired {
		t.Fatal("CustomMessage ResendCode never fired")
	}
	if code != user.SignUpCode {
		t.Fatalf("ResendCode code parameter %q does not match the stored sign-up code %q", code, user.SignUpCode)
	}
}

// GetUserAttributeVerificationCode fires the CustomMessage
// VerifyUserAttribute source with the issued code.
func TestGetUserAttributeVerificationCodeFiresCustomMessage(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		TokenGenerationAuthentication: {StatusCode: 200},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{CustomMessage: preAuthTestARN})
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	user := cognitostore.NewUser(poolID, "verify-user")
	user.UserStatus = "CONFIRMED"
	user.Enabled = true
	user.Attributes = map[string]string{"email": "verify@example.com"}
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}
	accessToken, _, _, _, err := svc.CreateTokens(reqCtx, poolID, user.ID, "trigger-client", TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := svc.getUserAttributeVerificationCodeCore(context.Background(), reqCtx, GetUserAttributeVerificationCodeInput{
		AccessToken: accessToken, AttributeName: "email",
	}); err != nil {
		t.Fatal(err)
	}
	code, fired := invoker.fired(CustomMessageVerifyUserAttribute)
	if !fired {
		t.Fatal("CustomMessage VerifyUserAttribute never fired")
	}
	stored, err := store.GetUser(poolID, "verify-user")
	if err != nil {
		t.Fatal(err)
	}
	av := stored.AttributeVerificationCodes["email"]
	if av == nil || av.Code == "" || av.Code != code {
		t.Fatalf("VerifyUserAttribute code parameter %q does not match the stored code %v", code, av)
	}
}

// Updating an email attribute enters the verification round: the verified
// flag drops, a VerifyUserAttribute code is minted and the CustomMessage
// UpdateUserAttribute source fires with it.
func TestAttributeUpdateEntersVerificationRound(t *testing.T) {
	invoker := &dispatchTriggerInvoker{}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{CustomMessage: preAuthTestARN})
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	user := cognitostore.NewUser(poolID, "update-user")
	user.UserStatus = "CONFIRMED"
	user.Attributes = map[string]string{
		"email":          "old@example.com",
		"email_verified": "true",
	}
	if err := store.CreateUser(user); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.adminUpdateUserAttributesCore(context.Background(), reqCtx, AdminUpdateUserAttributesInput{
		UserPoolID:     poolID,
		Username:       "update-user",
		UserAttributes: map[string]string{"email": "new@example.com"},
	}); err != nil {
		t.Fatal(err)
	}
	updated, err := store.GetUser(poolID, "update-user")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Attributes["email"] != "new@example.com" {
		t.Fatalf("email not updated: %v", updated.Attributes)
	}
	if updated.Attributes["email_verified"] != "" {
		t.Fatal("updated email kept its verified flag")
	}
	av := updated.AttributeVerificationCodes["email"]
	if av == nil || av.Code == "" {
		t.Fatal("attribute update minted no verification code")
	}
	code, fired := invoker.fired(CustomMessageUpdateUserAttribute)
	if !fired {
		t.Fatal("CustomMessage UpdateUserAttribute never fired")
	}
	if code != av.Code {
		t.Fatalf("UpdateUserAttribute code parameter %q does not match the minted code %q", code, av.Code)
	}
}

// storeUserID resolves a user's internal ID for token generation.
func storeUserID(t *testing.T, store cognitostore.CognitoStoreInterface, poolID, username string) string {
	t.Helper()
	user, err := store.GetUser(poolID, username)
	if err != nil {
		t.Fatal(err)
	}
	return user.ID
}

// A PreSignUp Lambda that raises an error rejects the sign-up with
// UserLambdaValidationException — the documented rejection path — not an
// infrastructure-class internal error.
func TestPreSignUpFailureClassifiedAsUserLambdaValidation(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		PreSignUpSignUp: {StatusCode: 200, FunctionError: "Unhandled"},
	}}
	svc, reqCtx, _, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{PreSignUp: preAuthTestARN})

	_, err := svc.signUpCore(context.Background(), reqCtx, SignUpInput{
		ClientID:       "trigger-client",
		Username:       "reject-user",
		Password:       "NewPass123!",
		UserAttributes: map[string]string{"email": "reject@example.com"},
	})
	if err != ErrUserLambdaValidation {
		t.Fatalf("raising PreSignUp returned %v, want UserLambdaValidationException", err)
	}
}

// A PreSignUp response that injects attributes outside the pool schema is
// a malformed trigger response, and one that sets the verified-claim keys
// directly has them stripped: the auto-verify flags are the only sanctioned
// path.
func TestPreSignUpAttributeInjectionValidated(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		PreSignUpSignUp: {StatusCode: 200, Payload: []byte(`{"autoConfirmUser":true,"userAttributes":{"custom:rank":"gold"}}`)},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{PreSignUp: preAuthTestARN})

	if _, err := svc.signUpCore(context.Background(), reqCtx, SignUpInput{
		ClientID:       "trigger-client",
		Username:       "schema-user",
		Password:       "NewPass123!",
		UserAttributes: map[string]string{"email": "schema@example.com"},
	}); err != ErrInvalidLambdaResponse {
		t.Fatalf("schema-violating PreSignUp attributes returned %v, want InvalidLambdaResponseException", err)
	}

	invoker = &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		PreSignUpSignUp: {StatusCode: 200, Payload: []byte(`{"autoConfirmUser":true,"userAttributes":{"email":"inject@example.com","email_verified":"true"}}`)},
	}}
	svc, reqCtx, poolID, _ = newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{PreSignUp: preAuthTestARN})
	if _, err := svc.signUpCore(context.Background(), reqCtx, SignUpInput{
		ClientID:       "trigger-client",
		Username:       "verify-user",
		Password:       "NewPass123!",
		UserAttributes: map[string]string{"email": "verify@example.com"},
	}); err != nil {
		t.Fatal(err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.GetUser(poolID, "verify-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.Attributes["email_verified"] != "" {
		t.Fatalf("PreSignUp response self-certified email: %v", user.Attributes)
	}
}

// The UserMigration response's finalUserStatus and messageAction carry
// fixed vocabularies; a value outside them is a malformed trigger response.
func TestUserMigrationEnumValidation(t *testing.T) {
	config := &cognitostore.LambdaConfig{UserMigration: triggerTestARN}

	svc := newTriggerTestService(t, &stubTriggerInvoker{payload: []byte(`{"userAttributes":{"email":"m@example.com"},"finalUserStatus":"ARCHIVED"}`)})
	if _, err := invokeUserMigration(context.Background(), svc, "us-east-1_POOL", "user", "client", "secret", config, nil, nil); err != ErrInvalidLambdaResponse {
		t.Fatalf("finalUserStatus=ARCHIVED returned %v, want InvalidLambdaResponseException", err)
	}

	svc = newTriggerTestService(t, &stubTriggerInvoker{payload: []byte(`{"userAttributes":{"email":"m@example.com"},"messageAction":"resend-please"}`)})
	if _, err := invokeUserMigration(context.Background(), svc, "us-east-1_POOL", "user", "client", "secret", config, nil, nil); err != ErrInvalidLambdaResponse {
		t.Fatalf("messageAction=resend-please returned %v, want InvalidLambdaResponseException", err)
	}

	svc = newTriggerTestService(t, &stubTriggerInvoker{payload: []byte(`{"userAttributes":{"email":"m@example.com"},"finalUserStatus":"RESET_REQUIRED","messageAction":"RESEND"}`)})
	result, err := invokeUserMigration(context.Background(), svc, "us-east-1_POOL", "user", "client", "secret", config, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.FinalUserStatus != "RESET_REQUIRED" || result.MessageAction != "RESEND" {
		t.Fatalf("documented enum values rejected: %+v", result)
	}
}

// A UserMigration response that injects attributes outside the pool schema
// is rejected; the verified flags it self-certifies are the documented
// trust of the administrative migration and reach the stored record.
func TestUserMigrationAttributesValidated(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		UserMigrationAuthentication: {StatusCode: 200, Payload: []byte(`{"userAttributes":{"custom:bogus":"1"},"finalUserStatus":"CONFIRMED","messageAction":"SUPPRESS"}`)},
	}}
	svc, reqCtx, _, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{UserMigration: triggerTestARN})

	if _, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "bogus-migration-user",
			"PASSWORD": "SomePass123!",
		},
	})); err != ErrInvalidLambdaResponse {
		t.Fatalf("schema-violating migration attributes returned %v, want InvalidLambdaResponseException", err)
	}

	invoker = &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		UserMigrationAuthentication: {StatusCode: 200, Payload: []byte(`{"userAttributes":{"email":"migrate@example.com","email_verified":"true"},"finalUserStatus":"CONFIRMED","messageAction":"RESEND"}`)},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{UserMigration: triggerTestARN})
	if _, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "verified-migration-user",
			"PASSWORD": "SomePass123!",
		},
	})); err != nil {
		t.Fatal(err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.GetUser(poolID, "verified-migration-user")
	if err != nil {
		t.Fatal(err)
	}
	// The migration trigger is an administrative trust boundary: its
	// documented example self-certifies the verified flag ("email_verified":
	// "true"), which the RESET_REQUIRED resolution's forgot-password flow
	// requires.
	if user.Attributes["email_verified"] != "true" {
		t.Fatalf("migration response's verified flag was dropped: %v", user.Attributes)
	}
}

// PostConfirmation is non-blocking by contract: a failing trigger is logged
// inside the invoker and never fails the confirmed sign-up.
func TestPostConfirmationNonBlocking(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		PreSignUpSignUp:               {StatusCode: 200, Payload: []byte(`{"autoConfirmUser":true}`)},
		PostConfirmationConfirmSignUp: {StatusCode: 200, FunctionError: "Unhandled"},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{PreSignUp: preAuthTestARN, PostConfirmation: preAuthTestARN})

	if _, err := svc.signUpCore(context.Background(), reqCtx, SignUpInput{
		ClientID:       "trigger-client",
		Username:       "post-confirm-user",
		Password:       "NewPass123!",
		UserAttributes: map[string]string{"email": "post@example.com"},
	}); err != nil {
		t.Fatalf("failing PostConfirmation failed the sign-up: %v", err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.GetUser(poolID, "post-confirm-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.UserStatus != "CONFIRMED" {
		t.Fatalf("user status %q, want CONFIRMED", user.UserStatus)
	}
	if !invoker.firedOnce(PostConfirmationConfirmSignUp) {
		t.Fatal("PostConfirmation never fired")
	}
}

// firedOnce reports whether exactly the given source fired at least once.
func (d *dispatchTriggerInvoker) firedOnce(source string) bool {
	_, ok := d.fired(source)
	return ok
}

// jwtPayloadClaims decodes a JWT's payload segment into a claim map.
func jwtPayloadClaims(t *testing.T, token string) map[string]interface{} {
	t.Helper()
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("token is not a three-part JWT: %q", token)
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode token payload: %v", err)
	}
	claims := make(map[string]interface{})
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("unmarshal token claims: %v", err)
	}
	return claims
}

// A SUPPRESS migration still migrates the sign-in password: messageAction
// governs the welcome message alone, so the migrated account completes the
// sign-in and verifies against the same password on the next one.
func TestUserMigrationSuppressMigratesCredentials(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		UserMigrationAuthentication: {StatusCode: 200, Payload: []byte(`{"userAttributes":{"email":"suppress@example.com"},"finalUserStatus":"CONFIRMED","messageAction":"SUPPRESS"}`)},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{UserMigration: triggerTestARN})

	authReq := challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "suppress-migration-user",
			"PASSWORD": "SomePass123!",
		},
	})
	if _, err := svc.InitiateAuth(context.Background(), reqCtx, authReq); err != nil {
		t.Fatalf("SUPPRESS migration sign-in failed: %v", err)
	}

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	migrated, err := store.GetUser(poolID, "suppress-migration-user")
	if err != nil {
		t.Fatal(err)
	}
	if migrated.PasswordHash == "" {
		t.Fatal("SUPPRESS migration stored no password credentials")
	}

	// The second sign-in finds the native credentials: the migration
	// trigger does not fire for a known user and the same password
	// verifies.
	if _, err := svc.InitiateAuth(context.Background(), reqCtx, authReq); err != nil {
		t.Fatalf("second sign-in against the migrated credentials failed: %v", err)
	}
}

// The migration response's omitted members resolve to the documented
// defaults: finalUserStatus RESET_REQUIRED and a welcome message on SMS.
// The authentication ends with PasswordResetRequiredException, the stored
// user keeps no native credentials, and later sign-ins keep surfacing the
// reset error.
func TestUserMigrationOmittedResponseDefaults(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		UserMigrationAuthentication: {StatusCode: 200, Payload: []byte(`{"userAttributes":{"email":"m@example.com"}}`)},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{UserMigration: triggerTestARN})

	_, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "defaulted-migration-user",
			"PASSWORD": "SomePass123!",
		},
	}))
	if err != ErrPasswordResetRequired {
		t.Fatalf("omitted-response migration returned %v, want PasswordResetRequiredException", err)
	}

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.GetUser(poolID, "defaulted-migration-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.UserStatus != "RESET_REQUIRED" || !user.MigratedAwaitingReset {
		t.Fatalf("omitted response resolved to %s/%v, want RESET_REQUIRED awaiting reset", user.UserStatus, user.MigratedAwaitingReset)
	}
	if user.PasswordHash != "" || user.SrpVerifier != "" {
		t.Fatal("RESET_REQUIRED migration stored native credentials; the forgot-password flow must set the password")
	}

	// The next sign-in surfaces the same reset error.
	_, err = svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "defaulted-migration-user",
			"PASSWORD": "SomePass123!",
		},
	}))
	if err != ErrPasswordResetRequired {
		t.Fatalf("post-migration sign-in returned %v, want PasswordResetRequiredException", err)
	}
}

// An explicit CONFIRMED/SUPPRESS migration retains the sign-in password as
// native credentials and completes the authentication; the retained
// password works at the next sign-in.
func TestUserMigrationConfirmedRetainsPassword(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		UserMigrationAuthentication: {StatusCode: 200, Payload: []byte(`{"userAttributes":{"email":"c@example.com"},"finalUserStatus":"CONFIRMED","messageAction":"SUPPRESS"}`)},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{UserMigration: triggerTestARN})

	resp, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "confirmed-migration-user",
			"PASSWORD": "SomePass123!",
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := resp.(map[string]interface{})["AuthenticationResult"]; !ok {
		t.Fatalf("CONFIRMED migration returned %#v, want AuthenticationResult", resp)
	}

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.GetUser(poolID, "confirmed-migration-user")
	if err != nil {
		t.Fatal(err)
	}
	if user.UserStatus != "CONFIRMED" || user.PasswordHash == "" || user.MigratedAwaitingReset {
		t.Fatalf("CONFIRMED migration stored %s (hash %q, awaiting %v)", user.UserStatus, user.PasswordHash, user.MigratedAwaitingReset)
	}

	resp, err = svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "confirmed-migration-user",
			"PASSWORD": "SomePass123!",
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := resp.(map[string]interface{})["AuthenticationResult"]; !ok {
		t.Fatalf("retained password failed the next sign-in: %#v", resp)
	}
}

// forceAliasCreation TRUE migrates a conflicting alias from the previous
// holder to the migrated user; the omitted default (false) leaves the
// conflict failing the sign-in.
func TestUserMigrationForceAliasCreation(t *testing.T) {
	holdout := func(t *testing.T) (*CognitoService, *request.RequestContext, string) {
		invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
			UserMigrationAuthentication: {StatusCode: 200, Payload: []byte(`{"userAttributes":{"email":"shared@example.com","email_verified":"true"},"finalUserStatus":"CONFIRMED","messageAction":"SUPPRESS","forceAliasCreation":true}`)},
		}}
		svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{UserMigration: triggerTestARN})
		store, err := svc.store(reqCtx)
		if err != nil {
			t.Fatal(err)
		}
		// The pool must treat email as an alias for claims to form.
		if err := store.UpdateUserPoolFunc(poolID, func(p *cognitostore.UserPool) error {
			p.AliasAttributes = []string{"email"}
			return nil
		}); err != nil {
			t.Fatal(err)
		}
		holder := cognitostore.NewUser(poolID, "holder")
		holder.UserStatus = "CONFIRMED"
		holder.Attributes = map[string]string{"email": "shared@example.com", "email_verified": "true"}
		if err := store.CreateUser(holder); err != nil {
			t.Fatal(err)
		}
		return svc, reqCtx, poolID
	}

	svc, reqCtx, poolID := holdout(t)
	if _, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "alias-migrant",
			"PASSWORD": "SomePass123!",
		},
	})); err != nil {
		t.Fatalf("forceAliasCreation migration failed: %v", err)
	}
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	holder, err := store.GetUser(poolID, "holder")
	if err != nil {
		t.Fatal(err)
	}
	if holder.Attributes["email_verified"] != "false" {
		t.Fatalf("previous holder kept the verified alias: %v", holder.Attributes)
	}
	migrant, err := store.GetUser(poolID, "alias-migrant")
	if err != nil {
		t.Fatal(err)
	}
	if migrant.Attributes["email_verified"] != "true" {
		t.Fatalf("migrated user did not receive the verified alias: %v", migrant.Attributes)
	}

	// The omitted default: the conflict fails the migration.
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		UserMigrationAuthentication: {StatusCode: 200, Payload: []byte(`{"userAttributes":{"email":"shared@example.com","email_verified":"true"},"finalUserStatus":"CONFIRMED","messageAction":"SUPPRESS"}`)},
	}}
	svc, reqCtx, poolID, _ = newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{UserMigration: triggerTestARN})
	store, err = svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.UpdateUserPoolFunc(poolID, func(p *cognitostore.UserPool) error {
		p.AliasAttributes = []string{"email"}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	holder2 := cognitostore.NewUser(poolID, "holder2")
	holder2.UserStatus = "CONFIRMED"
	holder2.Attributes = map[string]string{"email": "shared@example.com", "email_verified": "true"}
	if err := store.CreateUser(holder2); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.InitiateAuth(context.Background(), reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": "trigger-client",
		"AuthParameters": map[string]interface{}{
			"USERNAME": "alias-migrant-2",
			"PASSWORD": "SomePass123!",
		},
	})); err != ErrIncorrectPassword {
		t.Fatalf("alias conflict without forceAliasCreation returned %v, want the NotAuthorized family", err)
	}
	if _, err := store.GetUser(poolID, "alias-migrant-2"); err == nil {
		t.Fatal("conflicted migration left a user behind")
	}
}

// SignUp's error list carries UsernameExistsException alone: an alias
// duplicate a PreSignUp-activated verified alias produces surfaces under
// that error, not AliasExistsException.
func TestSignUpAliasCollisionIsUsernameExists(t *testing.T) {
	invoker := &dispatchTriggerInvoker{outcomes: map[string]invokers.LambdaInvocation{
		PreSignUpSignUp: {StatusCode: 200, Payload: []byte(`{"autoConfirmUser":true,"autoVerifyEmail":true}`)},
	}}
	svc, reqCtx, poolID, _ := newTriggerServiceEnv(t, invoker, &cognitostore.LambdaConfig{PreSignUp: triggerTestARN})
	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.UpdateUserPoolFunc(poolID, func(p *cognitostore.UserPool) error {
		p.AliasAttributes = []string{"email"}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	holder := cognitostore.NewUser(poolID, "signup-holder")
	holder.UserStatus = "CONFIRMED"
	holder.Attributes = map[string]string{"email": "dupe@example.com", "email_verified": "true"}
	if err := store.CreateUser(holder); err != nil {
		t.Fatal(err)
	}

	_, err = svc.signUpCore(context.Background(), reqCtx, SignUpInput{
		ClientID:       "trigger-client",
		Username:       "signup-duplicate",
		Password:       "NewPass123!",
		UserAttributes: map[string]string{"email": "dupe@example.com"},
	})
	if err != ErrUserAlreadyExists {
		t.Fatalf("alias-colliding SignUp returned %v (%T), want UsernameExistsException", err, err)
	}
}
