package cognitoidentityprovider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"golang.org/x/crypto/bcrypt"
)

// InitiateAuthInput carries the wire parameters of InitiateAuth. Params holds
// the raw request parameter map; the nested AuthParameters block is read from
// it inside the Core so the extraction semantics stay with the flow logic.
type InitiateAuthInput struct {
	AuthFlow       string
	ClientID       string
	Params         map[string]interface{}
	ValidationData map[string]string
	ClientMetadata map[string]string
}

// AdminInitiateAuthInput carries the wire parameters of AdminInitiateAuth.
type AdminInitiateAuthInput struct {
	UserPoolID     string
	AuthFlow       string
	ClientID       string
	Params         map[string]interface{}
	ValidationData map[string]string
	ClientMetadata map[string]string
}

// initiateAuthCore validates the requested flow and dispatches to the
// flow-specific implementation.
func (s *CognitoService) initiateAuthCore(ctx context.Context, reqCtx *request.RequestContext, in InitiateAuthInput) (interface{}, error) {
	if in.AuthFlow == "" || in.ClientID == "" {
		return nil, ErrInvalidParameter
	}
	if !validateInitiateAuthFlow(in.AuthFlow) {
		return nil, ErrInvalidParameter
	}

	// A secret-bearing client proves itself on every authentication request
	// — the model's AuthParameters contract: "Add a SECRET_HASH parameter if
	// your app client has a client secret". The refresh flows verify against
	// the token-resolved user inside the shared refresh path.
	if in.AuthFlow != "REFRESH_TOKEN_AUTH" && in.AuthFlow != "REFRESH_TOKEN" {
		if err := s.verifyAuthPlaneSecretHash(reqCtx, in.ClientID, authMember(in.Params, "USERNAME"), authMember(in.Params, "SECRET_HASH")); err != nil {
			return nil, err
		}
	}

	switch in.AuthFlow {
	case "USER_PASSWORD_AUTH":
		return s.userPasswordAuthFlow(ctx, reqCtx, in)
	case "USER_SRP_AUTH":
		return s.userSrpAuthFlow(reqCtx, in.ClientID, in.Params)
	case "REFRESH_TOKEN_AUTH", "REFRESH_TOKEN":
		return s.refreshTokenAuthFlow(reqCtx, in)
	case "CUSTOM_AUTH":
		return s.customAuthFlow(ctx, reqCtx, in.ClientID,
			authMember(in.Params, "USERNAME"),
			authMember(in.Params, "PASSWORD"))
	case "USER_AUTH":
		return s.userAuthFlow(reqCtx, in.ClientID,
			authMember(in.Params, "USERNAME"),
			authMember(in.Params, "PREFERRED_CHALLENGE"))
	default:
		return nil, ErrInvalidParameter
	}
}

// adminInitiateAuthCore validates the requested admin flow, resolves the
// named user pool and dispatches to the flow-specific implementation.
func (s *CognitoService) adminInitiateAuthCore(ctx context.Context, reqCtx *request.RequestContext, in AdminInitiateAuthInput) (interface{}, error) {
	if in.UserPoolID == "" || in.ClientID == "" || in.AuthFlow == "" {
		return nil, ErrInvalidParameter
	}
	if !validateAdminInitiateAuthFlow(in.AuthFlow) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPool(in.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	// The admin plane carries the same AuthParameters contract: "Add a
	// SECRET_HASH parameter if your app client has a client secret". The
	// refresh flow verifies against the token-resolved user instead.
	if in.AuthFlow != "REFRESH_TOKEN_AUTH" && in.AuthFlow != "REFRESH_TOKEN" {
		if err := s.verifyAuthPlaneSecretHash(reqCtx, in.ClientID, authMember(in.Params, "USERNAME"), authMember(in.Params, "SECRET_HASH")); err != nil {
			return nil, err
		}
	}

	switch in.AuthFlow {
	case "ADMIN_NO_SRP_AUTH", "ADMIN_USER_PASSWORD_AUTH":
		return s.adminNoSrpAuthFlow(ctx, reqCtx, in.UserPoolID, in.ClientID, userPool.LambdaConfig, in.Params, in.ValidationData, in.ClientMetadata)
	case "CUSTOM_AUTH":
		return s.customAuthFlow(ctx, reqCtx, in.ClientID,
			authMember(in.Params, "USERNAME"),
			authMember(in.Params, "PASSWORD"))
	case "REFRESH_TOKEN_AUTH", "REFRESH_TOKEN":
		return s.adminRefreshTokenAuthFlow(reqCtx, in)
	default:
		return nil, ErrInvalidParameter
	}
}

// customFlowChallengeNames lists the challenge types a DefineAuthChallenge
// Lambda response may nominate as the next challenge in a custom
// authentication flow. Anything outside this set is a malformed response.
// Per the Amazon Cognito developer guide, the DefineAuthChallenge response
// designates only the flow-level challenges (CUSTOM_CHALLENGE, SRP_A,
// PASSWORD_VERIFIER); MFA challenges — including MFA_SETUP — are issued by
// the service itself after credential verification ("You don't need to
// invoke any MFA challenges in your define auth challenge function").
// Accepting MFA_SETUP here would mint an MFA_SETUP-typed session without
// any verified credential, and AssociateSoftwareToken accepts exactly that
// session type, so an unauthenticated caller could overwrite a victim's
// TOTP configuration.
var customFlowChallengeNames = map[string]bool{
	"CUSTOM_CHALLENGE":  true,
	"SRP_A":             true,
	"PASSWORD_VERIFIER": true,
}

// resolveCustomFlowChallenge applies the customFlowChallengeNames contract
// to a DefineAuthChallenge response and returns the challenge the flow
// continues with. A nil or empty response keeps the CUSTOM_CHALLENGE
// default; a name outside the set is a malformed Lambda response.
func resolveCustomFlowChallenge(lambdaResult map[string]interface{}) (string, error) {
	if lambdaResult == nil {
		return "CUSTOM_CHALLENGE", nil
	}
	cn, ok := lambdaResult["challengeName"].(string)
	if !ok || cn == "" {
		return "CUSTOM_CHALLENGE", nil
	}
	if !customFlowChallengeNames[cn] {
		return "", ErrInvalidLambdaResponse
	}
	return cn, nil
}

// customAuthFlow implements the CUSTOM_AUTH flow. The client provides a
// USERNAME (and optionally a PASSWORD for the initial verification). The
// server invokes the DefineAuthChallenge and CreateAuthChallenge Lambda
// triggers to produce a custom challenge that the client must answer via
// RespondToAuthChallenge.
func (s *CognitoService) customAuthFlow(ctx context.Context, reqCtx *request.RequestContext, clientID, username, password string) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, username)
	if err != nil {
		// Return the same NotAuthorized error used by other auth flows to
		// prevent user enumeration via distinguishable session or error.
		return nil, ErrNotAuthorized
	}

	// If a password is provided, verify it as part of the initial auth.
	if password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			return nil, ErrNotAuthorized
		}
	}

	// Invoke DefineAuthChallenge Lambda trigger to determine the challenge.
	// A trigger failure fails the flow closed: the challenge selection is
	// the Lambda's decision, so falling through to the default challenge on
	// an error would bypass the pool's authentication policy.
	var lambdaResult map[string]interface{}
	if userPool.LambdaConfig != nil && userPool.LambdaConfig.DefineAuthChallenge != "" {
		res, derr := s.invokeTrigger(ctx, DefineAuthChallenge, userPool.ID, username, clientID,
			userPool.LambdaConfig.DefineAuthChallenge,
			map[string]interface{}{
				"userAttributes":  userAttributesMap(user),
				"challengeResult": "FAILED",
			},
			map[string]interface{}{
				"challengeName":      "",
				"issueTokens":        false,
				"failAuthentication": false,
			},
			true,
		)
		if derr != nil {
			return nil, classifyTriggerFailure(derr)
		}
		lambdaResult = res
	}

	sessionID := generateSessionID()
	challengeParams := map[string]string{
		"USERNAME": username,
	}

	challengeName, err := resolveCustomFlowChallenge(lambdaResult)
	if err != nil {
		return nil, err
	}

	// Invoke CreateAuthChallenge Lambda trigger to produce challenge parameters.
	if userPool.LambdaConfig != nil && userPool.LambdaConfig.CreateAuthChallenge != "" {
		createResult, cerr := s.invokeTrigger(ctx, CreateAuthChallenge, userPool.ID, username, clientID,
			userPool.LambdaConfig.CreateAuthChallenge,
			map[string]interface{}{
				"userAttributes": userAttributesMap(user),
				"challengeName":  challengeName,
			},
			map[string]interface{}{
				"publicChallengeParameters": map[string]string{},
			},
			true,
		)
		if cerr != nil {
			return nil, classifyTriggerFailure(cerr)
		}
		if createResult != nil {
			if params, ok := createResult["publicChallengeParameters"].(map[string]interface{}); ok {
				for k, v := range params {
					if vs, ok := v.(string); ok {
						challengeParams[k] = vs
					}
				}
			}
		}
	}

	challengeSession := &cognitostore.ChallengeSession{
		SessionID:     sessionID,
		UserPoolID:    userPool.ID,
		ClientID:      clientID,
		Username:      username,
		ChallengeName: challengeName,
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(challengeSessionTTL),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"ChallengeName":       challengeName,
		"Session":             sessionID,
		"ChallengeParameters": challengeParams,
	}, nil
}

// userAuthFlow implements the USER_AUTH flow (choice-based authentication).
// The server inspects the user's configured auth factors and returns the set
// of available challenges for the client to choose from. A PREFERRED_CHALLENGE
// parameter names the challenge to issue directly; a preference that is not
// among the available challenges is a hint the flow cannot honour and falls
// back to the selector response, so an unknown or unenrolled preference
// neither fails the documented flow nor leaks user state.
func (s *CognitoService) userAuthFlow(reqCtx *request.RequestContext, clientID, username, preferredChallenge string) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, username)

	// The selector session is issued regardless of whether the user
	// exists, so InitiateAuth(USER_AUTH) does not leak account existence;
	// every subsequent challenge response resolves the user again and
	// fails uniformly with NotAuthorized for unknown users. The
	// password-family challenges are advertised unconditionally — their
	// presence carries no account information — while the credential-backed
	// factors (passkeys, software token, OTPs) appear only for a user that
	// has enrolled them.
	var webauthnCreds []*cognitostore.WebAuthnCredential
	webauthnAvailable := false
	if err == nil && user != nil && userPool.WebAuthnConfiguration != nil && userPool.WebAuthnConfiguration.RelyingPartyId != "" {
		creds, cerr := listUserWebAuthnCredentials(store, userPool.ID, user.ID)
		if cerr != nil {
			return nil, ErrInternalError
		}
		webauthnCreds = creds
		webauthnAvailable = len(webauthnCreds) > 0
	}

	// A configured sign-in policy restricts the first authentication
	// factors the pool permits. PASSWORD and PASSWORD_SRP are both
	// password first factors; SOFTWARE_TOKEN_MFA is a second factor and
	// stays outside the policy's scope.
	firstFactorAllowed := func(factor string) bool {
		if userPool.SignInPolicy == nil || len(userPool.SignInPolicy.AllowedFirstAuthFactors) == 0 {
			return true
		}
		for _, f := range userPool.SignInPolicy.AllowedFirstAuthFactors {
			if f == factor {
				return true
			}
		}
		return false
	}

	var available []string
	if firstFactorAllowed("PASSWORD") {
		available = append(available, "PASSWORD")
	}
	if err == nil && user != nil && user.SrpVerifier != "" && firstFactorAllowed("PASSWORD") {
		// Any user with a stored SRP verifier can take the PASSWORD_SRP
		// leg of choice-based authentication.
		available = append(available, "PASSWORD_SRP")
	}
	if webauthnAvailable && firstFactorAllowed("WEB_AUTHN") {
		available = append(available, "WEB_AUTHN")
	}

	if err == nil && user != nil {
		if user.SoftwareTokenMfa != nil && user.SoftwareTokenMfa.Verified {
			available = append(available, "SOFTWARE_TOKEN_MFA")
		}
		if isAttributeVerified(user.Attributes, "phone_number") && firstFactorAllowed("SMS_OTP") {
			available = append(available, "SMS_OTP")
		}
		if isAttributeVerified(user.Attributes, "email") && firstFactorAllowed("EMAIL_OTP") {
			available = append(available, "EMAIL_OTP")
		}
	}

	// A preferred challenge is issued directly. PASSWORD_SRP carries only
	// the challenge shell — the client answers it with SRP_A; the OTP
	// selections deliver their code with the challenge.
	if user != nil && selectChallengeChoices[preferredChallenge] {
		for _, c := range available {
			if c != preferredChallenge {
				continue
			}
			switch preferredChallenge {
			case "WEB_AUTHN":
				return s.issueWebAuthnChallenge(store, userPool, clientID, user, "")
			case "PASSWORD", "PASSWORD_SRP", "SMS_OTP", "EMAIL_OTP":
				session, merr := mintChallengeSession(store, userPool.ID, clientID, username, preferredChallenge, challengeSessionTTL)
				if merr != nil {
					return nil, ErrInternalError
				}
				s.deliverChallengeCode(reqCtx, userPool.ID, preferredChallenge)
				return map[string]interface{}{
					"ChallengeName": preferredChallenge,
					"Session":       session,
					"ChallengeParameters": map[string]string{
						"USERNAME": username,
					},
				}, nil
			}
		}
	}

	// No preference to honour: the selector session is the response. When
	// the user can sign in with a passkey the response also carries the
	// assertion request options, bound to this session, so the client can
	// complete WEB_AUTHN inside the selection response as documented.
	challengeSession := newChallengeSession(userPool.ID, clientID, username, "SELECT_CHALLENGE", challengeSessionTTL)
	challengeParameters := map[string]string{"USERNAME": username}
	if webauthnAvailable {
		challengeB64, rpID, optionsJSON, oerr := s.buildWebAuthnRequestOptions(store, userPool, webauthnCreds)
		if oerr != nil {
			return nil, oerr
		}
		challengeSession.ChallengeData = challengeB64
		challengeSession.RelyingPartyID = rpID
		challengeParameters["CREDENTIAL_REQUEST_OPTIONS"] = optionsJSON
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{
		"AvailableChallenges": available,
		"Session":             challengeSession.SessionID,
		"ChallengeParameters": challengeParameters,
	}, nil
}

// authenticateUser contains the shared authentication logic used by both
// USER_PASSWORD_AUTH (InitiateAuth) and ADMIN_NO_SRP_AUTH (AdminInitiateAuth).
// It handles user lookup, migration, challenge detection, credential
// verification, trigger invocation, and token generation. The deviceKey
// carries the request's DEVICE_KEY auth parameter: a remembered device
// replaces the pool's MFA challenge with device SRP authentication. The
// trailing nonce parameter is the OIDC nonce an authorisation request
// supplied: non-empty only on the hosted-UI plane, it is bound into the ID
// token per the authorize-endpoint contract ("The nonce value that you
// provide is included in the ID token that Amazon Cognito issues"); the API
// plane passes an empty string.
func (s *CognitoService) authenticateUser(
	ctx context.Context,
	reqCtx *request.RequestContext,
	userPoolID, clientID, username, password, deviceKey string,
	lambdaConfig *cognitostore.LambdaConfig,
	validationData, clientMetadata map[string]string,
	nonce string,
) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		migrationResult, migrationErr := invokeUserMigration(ctx, s, userPoolID, username, clientID, password, lambdaConfig, validationData, clientMetadata)
		// A Lambda function that raises an error rejects the migration and
		// fails the sign-in as a wrong password; an invocation-transport
		// failure is an infrastructure error. A missing trigger means no
		// migration path exists and authentication fails as an incorrect
		// password.
		if migrationErr != nil {
			return nil, classifyMigrationFailure(migrationErr)
		}
		if migrationResult == nil {
			return nil, ErrIncorrectPassword
		}

		migratedUser := cognitostore.NewUser(userPoolID, username)
		if migrationResult.UserAttributes != nil {
			migratedUser.Attributes = migrationResult.UserAttributes
		}
		if migratedUser.Attributes == nil {
			migratedUser.Attributes = make(map[string]string)
		}
		// Every user's record carries its immutable identifier as the
		// required sub attribute; the migration trigger's response can
		// never supply it, so any stale entry is replaced by the
		// service-assigned value.
		delete(migratedUser.Attributes, "sub")
		// The migration response replaced the whole attribute map, so the
		// schema gate applies to what actually gets stored. The verified
		// flags survive the transfer: the migration trigger is an
		// administrative trust boundary whose documented example sets them
		// ("email_verified": "true"), the forceAliasCreation semantics
		// operate on verified alias claims, and the RESET_REQUIRED
		// resolution's forgot-password flow requires a verified contact
		// attribute. A map outside the pool's schema is a malformed
		// trigger response.
		migrationPool, mperr := store.GetUserPool(userPoolID)
		if mperr != nil {
			return nil, ErrResourceNotFound
		}
		if err := validateUserAttributesAgainstSchema(migrationPool, migratedUser.Attributes, false); err != nil {
			return nil, ErrInvalidLambdaResponse
		}
		migratedUser.Attributes["sub"] = migratedUser.ID
		migratedUser.UserStatus = migrationResult.FinalUserStatus
		// The sign-in password migrates into native credentials only when
		// the resolved status lets the user keep signing in with it
		// ("auto-confirm your users so that they can sign in with their
		// previous passwords" — CONFIRMED and FORCE_CHANGE_PASSWORD); a
		// RESET_REQUIRED resolution keeps no native credentials, and the
		// user sets a password through the forgot-password flow instead.
		if password != "" && migratedUser.UserStatus != "RESET_REQUIRED" {
			// The pool policy governs the password history the migrated
			// credentials retain.
			if err := setNativePasswordCredentials(migratedUser, migrationPool.PasswordPolicy, password); err != nil {
				if errors.Is(err, ErrPasswordHistoryViolation) {
					return nil, ErrPasswordHistoryViolation
				}
				return nil, ErrInternalError
			}
		}
		if migratedUser.UserStatus == "RESET_REQUIRED" {
			migratedUser.MigratedAwaitingReset = true
		}
		// enableSMSMFA requires a phone number in the response attributes
		// ("Your user's attributes in the request parameters must include a
		// phone number, or else the migration of that user will fail").
		if migrationResult.EnableSMSMFA {
			if _, ok := migratedUser.Attributes["phone_number"]; !ok {
				return nil, ErrInvalidLambdaResponse
			}
			migratedUser.SmsMfa = &cognitostore.SmsMfaSettings{Enabled: true}
		}
		createErr := store.CreateUser(migratedUser)
		if errors.Is(createErr, cognitostore.ErrAliasExists) && migrationResult.ForceAliasCreation {
			// forceAliasCreation TRUE migrates a conflicting alias "from
			// the previous user to the newly created user" — the create
			// retries with alias-claim takeover.
			createErr = store.CreateUserMigrateAliasClaims(migratedUser)
		}
		if createErr != nil {
			if errors.Is(createErr, cognitostore.ErrUserAlreadyExists) || errors.Is(createErr, cognitostore.ErrAliasExists) {
				return nil, ErrIncorrectPassword
			}
			return nil, ErrInternalError
		}
		user = migratedUser
		// A RESET_REQUIRED migration ends the authentication here: the
		// model documents that the client "must handle the
		// PasswordResetRequiredException during the authentication flow".
		if migratedUser.UserStatus == "RESET_REQUIRED" {
			s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
			return nil, ErrPasswordResetRequired
		}
		// An omitted messageAction sends the welcome message ("If your
		// function doesn't return this parameter, Amazon Cognito sends the
		// welcome message"), on the resolved delivery mediums.
		if migrationResult.MessageAction != "SUPPRESS" {
			s.publishNotificationLogCore(reqCtx, userPoolID, fmt.Sprintf("UserMigration welcome message via %s", strings.Join(migrationResult.DesiredDeliveryMediums, "/")))
		}
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrNotAuthorized
	}

	// A password a password reset has deactivated cannot complete a sign-in:
	// the reset "responds with a PasswordResetRequiredException error" and
	// the user completes the forgot-password flow instead. This covers users
	// whose native credentials a reset deactivated (AdminResetUserPassword
	// keeps the stored hash) and users the UserMigration trigger created
	// with a RESET_REQUIRED resolution; the CSV-import exception below keeps
	// its own any-password NEW_PASSWORD_REQUIRED flow only until a reset
	// sets the marker.
	if user.UserStatus == "RESET_REQUIRED" &&
		((user.PasswordHash != "" && user.PasswordHashAlgo == "") || user.MigratedAwaitingReset) {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrPasswordResetRequired
	}

	// The password must verify before any challenge is issued. Amazon
	// Cognito returns NEW_PASSWORD_REQUIRED only for users who signed in
	// successfully with their temporary password; issuing it on user status
	// alone would let anyone reset a FORCE_CHANGE_PASSWORD or RESET_REQUIRED
	// account by username only. Users imported from a CSV are the documented
	// exception: "the first time they sign in, they can enter any password.
	// Amazon Cognito prompts them to enter a new password" — an imported
	// RESET_REQUIRED user with no stored hash goes straight to the
	// NEW_PASSWORD_REQUIRED challenge, and a user holding an imported hash
	// verifies against the import algorithm, migrating to the native
	// bcrypt+SRP credentials on success.
	switch {
	case user.PasswordHashAlgo != "":
		if !verifyImportedPasswordHash(user.PasswordHashAlgo, user.PasswordHash, password) {
			s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
			return nil, ErrIncorrectPassword
		}
		// The pool is only loaded on this rare migration path: the password
		// policy it carries governs the history recording the native
		// credentials below retain.
		migrationPool, mperr := store.GetUserPool(userPoolID)
		if mperr != nil {
			return nil, ErrResourceNotFound
		}
		if err := s.migrateImportedCredentials(store, migrationPool.PasswordPolicy, user, password); err != nil {
			return nil, err
		}
	case user.PasswordHash == "":
		if user.UserStatus != "RESET_REQUIRED" {
			s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
			return nil, ErrIncorrectPassword
		}
	default:
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
			return nil, ErrIncorrectPassword
		}
	}

	if user.UserStatus == "FORCE_CHANGE_PASSWORD" || user.UserStatus == "RESET_REQUIRED" {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseInProgress)
		return s.newPasswordChallenge(reqCtx, userPoolID, clientID, user)
	}

	if user.UserStatus != "CONFIRMED" {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrUserNotConfirmed
	}

	attrs := userAttributesMap(user)
	if err := invokePreAuthentication(ctx, s, userPoolID, username, clientID, lambdaConfig, attrs, clientMetadata); err != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
		return nil, err
	}

	// The pool's MFA configuration applies once the primary credentials
	// have verified: a second-factor challenge replaces the token issuance,
	// and PostAuthentication only fires when the authentication completes.
	// A remembered device named by DEVICE_KEY replaces the MFA challenge
	// with device SRP authentication.
	pool, err := store.GetUserPool(userPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	challenge, err := s.nextChallengeAfterPrimary(reqCtx, store, pool, clientID, user, deviceKey)
	if err != nil {
		return nil, ErrInternalError
	}
	if challenge != nil {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponseInProgress)
		return challenge, nil
	}

	if err := invokePostAuthentication(ctx, s, userPoolID, username, clientID, lambdaConfig, attrs, clientMetadata); err != nil {
		return nil, err
	}

	var idClaims map[string]interface{}
	if nonce != "" {
		idClaims = map[string]interface{}{"nonce": nonce}
	}
	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokensWithIDClaims(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, clientMetadata, idClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	s.recordAuthEvent(reqCtx, userPoolID, user.ID, username, clientID, authEventSignIn, authEventResponsePass)
	return s.withNewDeviceMetadata(store, pool, user, deviceKey, authResult(accessToken, idToken, refreshToken, expiresIn)), nil
}

// migrateImportedCredentials replaces an imported password hash with the
// native bcrypt+SRP pair after successful verification, mirroring AWS's
// transparent credential migration at first sign-in. The rewrite is a
// serialised read-modify-write so a concurrent password change's history
// cannot be overwritten by the stale snapshot the caller authenticated
// against.
func (s *CognitoService) migrateImportedCredentials(store cognitostore.CognitoStoreInterface, policy *cognitostore.PasswordPolicy, user *cognitostore.User, password string) error {
	err := store.UpdateUserFunc(user.UserPoolID, user.Username, func(u *cognitostore.User) error {
		return setNativePasswordCredentials(u, policy, password)
	})
	if err != nil {
		if errors.Is(err, ErrPasswordHistoryViolation) {
			return ErrPasswordHistoryViolation
		}
		var awsErr *awserrors.AWSError
		if errors.As(err, &awsErr) {
			return err
		}
		return ErrInternalError
	}
	return nil
}

// newPasswordChallenge creates a NEW_PASSWORD_REQUIRED challenge session and
// returns the challenge response.
func (s *CognitoService) newPasswordChallenge(reqCtx *request.RequestContext, userPoolID, clientID string, user *cognitostore.User) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session := generateSessionID()
	challengeSession := &cognitostore.ChallengeSession{
		SessionID:     session,
		UserPoolID:    userPoolID,
		ClientID:      clientID,
		Username:      user.Username,
		ChallengeName: "NEW_PASSWORD_REQUIRED",
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(srpChallengeSessionTTL),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}
	return map[string]interface{}{
		"ChallengeName": "NEW_PASSWORD_REQUIRED",
		"Session":       session,
		"ChallengeParameters": map[string]interface{}{
			"USER_ID_FOR_SRP":    user.Username,
			"requiredAttributes": "[]",
		},
	}, nil
}

// userPasswordAuthFlow implements the InitiateAuth USER_PASSWORD_AUTH flow.
func (s *CognitoService) userPasswordAuthFlow(ctx context.Context, reqCtx *request.RequestContext, in InitiateAuthInput) (interface{}, error) {
	username, password, err := authParamsFrom(in.Params)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	userPool, err := store.GetUserPoolByClientID(in.ClientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	deviceKey := authMember(in.Params, "DEVICE_KEY")
	return s.authenticateUser(ctx, reqCtx, userPool.ID, in.ClientID, username, password, deviceKey, userPool.LambdaConfig, in.ValidationData, in.ClientMetadata, "")
}

// adminNoSrpAuthFlow implements the AdminInitiateAuth
// ADMIN_NO_SRP_AUTH / ADMIN_USER_PASSWORD_AUTH flows.
func (s *CognitoService) adminNoSrpAuthFlow(
	ctx context.Context, reqCtx *request.RequestContext,
	userPoolID, clientID string,
	lambdaConfig *cognitostore.LambdaConfig,
	params map[string]interface{},
	validationData, clientMetadata map[string]string,
) (interface{}, error) {
	username, password, err := authParamsFrom(params)
	if err != nil {
		return nil, err
	}

	return s.authenticateUser(ctx, reqCtx, userPoolID, clientID, username, password, authMember(params, "DEVICE_KEY"), lambdaConfig, validationData, clientMetadata, "")
}

// userSrpAuthFlow handles the InitiateAuth USER_SRP_AUTH flow. It receives
// the client's SRP_A value, looks up the user and defers to the shared SRP
// server leg, which generates the server's ephemeral B and a fresh
// SECRET_BLOCK, persists the SRP session state, and returns a
// PASSWORD_VERIFIER challenge. The actual proof verification happens in
// respondToPasswordVerifierCore when the client responds.
//
// Per AWS spec, USER_ID_FOR_SRP returned to the client is the user's username
// (the same value the client must use in the inner hash and claim message).
func (s *CognitoService) userSrpAuthFlow(reqCtx *request.RequestContext, clientID string, params map[string]interface{}) (interface{}, error) {
	username, srpAHex, err := srpAuthParamsFrom(params)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	userPool, err := store.GetUserPoolByClientID(clientID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	user, err := store.GetUser(userPool.ID, username)
	if err != nil {
		// Avoid revealing whether the username exists; mirror the NotAuthorized
		// semantics used elsewhere in the auth flow.
		return nil, ErrNotAuthorized
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrNotAuthorized
	}

	return s.issueSrpPasswordVerifier(reqCtx, store, userPool, clientID, user, srpAHex)
}

// issueSrpPasswordVerifier performs the server leg of the password SRP
// exchange shared by InitiateAuth(USER_SRP_AUTH), the PASSWORD_SRP challenge
// answer and the inline PASSWORD_SRP selection: it generates the server
// ephemeral B and a fresh SECRET_BLOCK against the user's stored verifier,
// persists the PASSWORD_VERIFIER session state and renders the challenge
// response the client computes its claim against.
func (s *CognitoService) issueSrpPasswordVerifier(
	reqCtx *request.RequestContext,
	store cognitostore.CognitoStoreInterface,
	userPool *cognitostore.UserPool,
	clientID string,
	user *cognitostore.User,
	srpAHex string,
) (interface{}, error) {
	if user.UserStatus != "CONFIRMED" {
		// Users in FORCE_CHANGE_PASSWORD or other transitional states cannot
		// complete SRP; they must first complete the NEW_PASSWORD_REQUIRED
		// challenge via the non-SRP admin flow.
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, user.Username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrUserNotConfirmed
	}
	if user.SrpVerifier == "" || user.SrpSalt == "" {
		// A user whose credentials were imported from a CSV carries only a
		// password hash, no SRP material, so the client's proof cannot be
		// checked. Setting a password through ForgotPassword or
		// AdminSetUserPassword provisions the verifier; until then USER_SRP_AUTH
		// cannot succeed for them.
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, user.Username, clientID, authEventSignIn, authEventResponseFail)
		return nil, ErrNotAuthorized
	}

	verifier, ok := new(big.Int).SetString(user.SrpVerifier, 16)
	if !ok {
		return nil, ErrInternalError
	}
	B, b, err := GenerateB(verifier)
	if err != nil {
		return nil, ErrInternalError
	}

	secretBlock := make([]byte, 16)
	if _, err := rand.Read(secretBlock); err != nil {
		return nil, ErrInternalError
	}

	challengeSession := newChallengeSession(userPool.ID, clientID, user.Username, "PASSWORD_VERIFIER", srpChallengeSessionTTL)
	challengeSession.SrpA = srpAHex
	challengeSession.SrpB = B.Text(16)
	challengeSession.SrpPrivateB = b.Text(16)
	challengeSession.SecretBlock = base64.StdEncoding.EncodeToString(secretBlock)
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	s.recordAuthEvent(reqCtx, userPool.ID, user.ID, user.Username, clientID, authEventSignIn, authEventResponseInProgress)

	return map[string]interface{}{
		"ChallengeName": "PASSWORD_VERIFIER",
		"Session":       challengeSession.SessionID,
		"ChallengeParameters": map[string]interface{}{
			"USERNAME":        user.Username,
			"USER_ID_FOR_SRP": user.Username,
			"SALT":            user.SrpSalt,
			"SECRET_BLOCK":    challengeSession.SecretBlock,
			"SRP_B":           B.Text(16),
		},
	}, nil
}

// refreshTokenAuthFlow implements the InitiateAuth REFRESH_TOKEN_AUTH flow.
func (s *CognitoService) refreshTokenAuthFlow(reqCtx *request.RequestContext, in InitiateAuthInput) (interface{}, error) {
	refreshToken, err := refreshTokenFrom(in.Params)
	if err != nil {
		return nil, err
	}
	return s.refreshAuthToken(reqCtx, "", in.ClientID, refreshToken, authMember(in.Params, "SECRET_HASH"), true)
}

// adminRefreshTokenAuthFlow implements the AdminInitiateAuth
// REFRESH_TOKEN_AUTH flow with an explicit user pool binding.
func (s *CognitoService) adminRefreshTokenAuthFlow(reqCtx *request.RequestContext, in AdminInitiateAuthInput) (interface{}, error) {
	refreshToken, err := refreshTokenFrom(in.Params)
	if err != nil {
		return nil, err
	}
	return s.refreshAuthToken(reqCtx, in.UserPoolID, in.ClientID, refreshToken, authMember(in.Params, "SECRET_HASH"), true)
}

// refreshAuthToken contains the shared refresh-token flow for both InitiateAuth
// and AdminInitiateAuth. It validates the refresh token, looks up the user, and
// issues new access/ID tokens. The SDK auth planes pass the client's
// SECRET_HASH proof (secretHashEnforced with the presented secretHash); the
// hosted UI token endpoint authenticates its client by the OAuth client-secret
// contract before calling and passes no proof.
func (s *CognitoService) refreshAuthToken(reqCtx *request.RequestContext, userPoolID, clientID, refreshToken, secretHash string, secretHashEnforced bool) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rt, err := store.GetRefreshTokenByValue(refreshToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// A refresh token is scoped to the user pool and the app client that
	// issued it. When the caller names either (the Admin flows carry
	// UserPoolId, every InitiateAuth carries ClientId, the hosted UI
	// carries both) they must match the token's own bindings; otherwise
	// the request could mint tokens signed with another pool's keys or
	// attributed to a client the token was never issued to.
	if userPoolID != "" && userPoolID != rt.UserPoolID {
		return nil, ErrNotAuthorized
	}
	if clientID != "" && clientID != rt.ClientID {
		return nil, ErrNotAuthorized
	}
	poolID := rt.UserPoolID

	user, err := store.GetUserByID(rt.UserID)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	// A disabled user cannot mint new tokens through a refresh grant.
	if !user.Enabled {
		return nil, ErrNotAuthorized
	}

	if secretHashEnforced {
		client, err := store.GetUserPoolClient(poolID, rt.ClientID)
		if err != nil {
			return nil, ErrInternalError
		}
		// The guide's REFRESH_TOKEN_AUTH rule computes the hash's username
		// element from the username claim when the pool takes username
		// sign-in, and from the sub claim otherwise — both the username and
		// the user ID verify, so either derivation is accepted.
		if verifyClientSecretHash(client, user.Username, secretHash) != nil &&
			verifyClientSecretHash(client, user.ID, secretHash) != nil {
			return nil, ErrNotAuthorized
		}
	}

	attrs := userAttributesMap(user)
	if err := invokePostAuthentication(reqCtx, s, poolID, user.Username, rt.ClientID, nil, attrs, nil); err != nil {
		return nil, err
	}

	accessToken, idToken, _, expiresIn, err := s.CreateTokensWithSessionScope(reqCtx, poolID, user.ID, rt.ClientID, TokenGenerationRefreshTokens, rt.Scope, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	return authResultNoRefresh(accessToken, idToken, expiresIn), nil
}

// authMember reads an InitiateAuth AuthParameters member. The documented
// carrier is the AuthParameters map (per the AWS API reference: USER_AUTH
// takes USERNAME, CUSTOM_AUTH takes USERNAME, USER_PASSWORD_AUTH takes
// USERNAME and PASSWORD); every key passed here is one of the documented
// uppercase AuthParameters names, and nothing outside the map is read.
func authMember(params map[string]interface{}, keys ...string) string {
	authParams, ok := params["AuthParameters"].(map[string]interface{})
	if !ok {
		return ""
	}
	for _, k := range keys {
		if s, ok := authParams[k].(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// authParamsFrom extracts USERNAME and PASSWORD from the AuthParameters
// block of an InitiateAuth or AdminInitiateAuth request.
func authParamsFrom(params map[string]interface{}) (username, password string, err error) {
	authParams := params["AuthParameters"]
	if authParams == nil {
		return "", "", ErrInvalidParameter
	}
	m, ok := authParams.(map[string]interface{})
	if !ok {
		return "", "", ErrInvalidParameter
	}
	username, _ = m["USERNAME"].(string)
	password, _ = m["PASSWORD"].(string)
	if username == "" || password == "" {
		return "", "", ErrInvalidParameter
	}
	return username, password, nil
}

// srpAuthParamsFrom extracts USERNAME and SRP_A from the AuthParameters
// block of an InitiateAuth USER_SRP_AUTH request. SRP_A is a lowercase hex
// string supplied by the client.
func srpAuthParamsFrom(params map[string]interface{}) (username, srpAHex string, err error) {
	authParams := params["AuthParameters"]
	if authParams == nil {
		return "", "", ErrInvalidParameter
	}
	m, ok := authParams.(map[string]interface{})
	if !ok {
		return "", "", ErrInvalidParameter
	}
	username, _ = m["USERNAME"].(string)
	srpAHex, _ = m["SRP_A"].(string)
	if username == "" || srpAHex == "" {
		return "", "", ErrInvalidParameter
	}
	if _, ok := new(big.Int).SetString(srpAHex, 16); !ok {
		return "", "", ErrInvalidParameter
	}
	return username, srpAHex, nil
}

// refreshTokenFrom extracts the REFRESH_TOKEN from AuthParameters.
func refreshTokenFrom(params map[string]interface{}) (string, error) {
	authParams := params["AuthParameters"]
	if authParams == nil {
		return "", ErrInvalidParameter
	}
	m, ok := authParams.(map[string]interface{})
	if !ok {
		return "", ErrInvalidParameter
	}
	refreshToken, _ := m["REFRESH_TOKEN"].(string)
	if refreshToken == "" {
		return "", ErrInvalidParameter
	}
	return refreshToken, nil
}

func authResult(accessToken, idToken, refreshToken string, expiresIn int64) map[string]interface{} {
	return map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken":  accessToken,
			"IdToken":      idToken,
			"RefreshToken": refreshToken,
			"TokenType":    "Bearer",
			"ExpiresIn":    expiresIn,
		},
	}
}

func authResultNoRefresh(accessToken, idToken string, expiresIn int64) map[string]interface{} {
	return map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken": accessToken,
			"IdToken":     idToken,
			"TokenType":   "Bearer",
			"ExpiresIn":   expiresIn,
		},
	}
}
