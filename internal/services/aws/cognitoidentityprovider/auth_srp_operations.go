package cognitoidentityprovider

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// handleUserSrpAuth handles the InitiateAuth USER_SRP_AUTH flow. It receives
// the client's SRP_A value, looks up the user, generates the server's
// ephemeral B and a fresh SECRET_BLOCK, persists the SRP session state, and
// returns a PASSWORD_VERIFIER challenge. The actual proof verification happens
// in respondToPasswordVerifier when the client responds.
//
// Per AWS spec, USER_ID_FOR_SRP returned to the client is the user's username
// (the same value the client must use in the inner hash and claim message).
func (s *CognitoService) handleUserSrpAuth(reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	username, srpAHex, err := parseSrpAuthParams(req)
	if err != nil {
		return nil, err
	}

	clientID := getClientId(req)
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
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}
	if user.UserStatus != "CONFIRMED" {
		// Users in FORCE_CHANGE_PASSWORD or other transitional states cannot
		// complete SRP; they must first complete the NEW_PASSWORD_REQUIRED
		// challenge via the non-SRP admin flow.
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, "SignIn", "Fail")
		return nil, ErrUserNotConfirmed
	}
	if user.SrpVerifier == "" || user.SrpSalt == "" {
		// The user was created before SRP support was added. They must reset
		// their password via ForgotPassword/AdminSetUserPassword to obtain a
		// verifier before USER_SRP_AUTH will succeed.
		s.recordAuthEvent(reqCtx, userPool.ID, user.ID, "SignIn", "Fail")
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

	sessionID := generateSessionID()
	challengeSession := &cognitostore.ChallengeSession{
		SessionID:     sessionID,
		UserPoolID:    userPool.ID,
		ClientID:      clientID,
		Username:      user.Username,
		ChallengeName: "PASSWORD_VERIFIER",
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(5 * time.Minute),
		SrpA:          srpAHex,
		SrpB:          B.Text(16),
		SrpPrivateB:   b.Text(16),
		SecretBlock:   base64.StdEncoding.EncodeToString(secretBlock),
	}
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	s.recordAuthEvent(reqCtx, userPool.ID, user.ID, "SignIn", "InProgress")

	return map[string]interface{}{
		"ChallengeName": "PASSWORD_VERIFIER",
		"Session":       sessionID,
		"ChallengeParameters": map[string]interface{}{
			"USERNAME":        user.Username,
			"USER_ID_FOR_SRP": user.Username,
			"SALT":            user.SrpSalt,
			"SECRET_BLOCK":    challengeSession.SecretBlock,
			"SRP_B":           B.Text(16),
		},
	}, nil
}

// parseSrpAuthParams extracts USERNAME and SRP_A from the AuthParameters
// block of an InitiateAuth USER_SRP_AUTH request. SRP_A is a lowercase hex
// string supplied by the client.
func parseSrpAuthParams(req *request.ParsedRequest) (username, srpAHex string, err error) {
	authParams := req.Parameters["AuthParameters"]
	if authParams == nil {
		return "", "", ErrInvalidParameter
	}
	params, ok := authParams.(map[string]interface{})
	if !ok {
		return "", "", ErrInvalidParameter
	}
	username, _ = params["USERNAME"].(string)
	srpAHex, _ = params["SRP_A"].(string)
	if username == "" || srpAHex == "" {
		return "", "", ErrInvalidParameter
	}
	if _, ok := new(big.Int).SetString(srpAHex, 16); !ok {
		return "", "", ErrInvalidParameter
	}
	return username, srpAHex, nil
}

// respondToPasswordVerifier handles the PASSWORD_VERIFIER challenge for both
// RespondToAuthChallenge and AdminRespondToAuthChallenge. It retrieves the
// SRP session established by handleUserSrpAuth, recomputes the shared secret
// from the stored private scalar b and the user's verifier, derives the
// HMAC key, and constant-time compares the expected claim against the client's
// PASSWORD_CLAIM_SIGNATURE.
//
// AWS Cognito returns a generic NotAuthorizedException on any mismatch
// (wrong password, tampered SRP_A, expired session, etc.) to avoid leaking
// which leg of the verification failed.
func (s *CognitoService) respondToPasswordVerifier(reqCtx *request.RequestContext, req *request.ParsedRequest, userPoolID, clientID, session string) (interface{}, error) {
	challengeResponses := req.Parameters["ChallengeResponses"]
	if challengeResponses == nil {
		return nil, ErrInvalidParameter
	}
	params, ok := challengeResponses.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	username, _ := params["USERNAME"].(string)
	claimSigB64, _ := params["PASSWORD_CLAIM_SIGNATURE"].(string)
	claimBlockB64, _ := params["PASSWORD_CLAIM_SECRET_BLOCK"].(string)
	timestamp, _ := params["TIMESTAMP"].(string)
	if username == "" || claimSigB64 == "" || claimBlockB64 == "" || timestamp == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	challengeSession, err := validateChallengeSession(store, session, "PASSWORD_VERIFIER", userPoolID, clientID, username)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	// Burn the session regardless of outcome so a leaked Session identifier
	// cannot be replayed.
	defer func() { _ = store.DeleteChallengeSession(session) }()

	if challengeSession.SecretBlock != claimBlockB64 {
		// The client must echo back the exact SECRET_BLOCK we issued. A
		// mismatch indicates the client did not use our challenge parameters.
		return nil, ErrNotAuthorized
	}

	A, ok := new(big.Int).SetString(challengeSession.SrpA, 16)
	if !ok {
		return nil, ErrInternalError
	}
	B, ok := new(big.Int).SetString(challengeSession.SrpB, 16)
	if !ok {
		return nil, ErrInternalError
	}
	b, ok := new(big.Int).SetString(challengeSession.SrpPrivateB, 16)
	if !ok {
		return nil, ErrInternalError
	}
	secretBlock, err := base64.StdEncoding.DecodeString(challengeSession.SecretBlock)
	if err != nil {
		return nil, ErrInternalError
	}
	clientSig, err := base64.StdEncoding.DecodeString(claimSigB64)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	user, err := store.GetUser(userPoolID, username)
	if err != nil {
		return nil, ErrNotAuthorized
	}
	verifier, ok := new(big.Int).SetString(user.SrpVerifier, 16)
	if !ok || verifier.Sign() == 0 {
		return nil, ErrNotAuthorized
	}

	K, err := DeriveServerKey(A, B, b, verifier)
	if err != nil {
		// ErrInvalidSrpA indicates a malicious or malformed client value.
		return nil, ErrNotAuthorized
	}

	poolName, ok := poolNameFromID(userPoolID)
	if !ok {
		return nil, ErrInternalError
	}
	// USER_ID_FOR_SRP for Cognito is the username (not the sub).
	expectedSig := VerifyClaim(K, poolName, user.Username, secretBlock, timestamp)

	if subtle.ConstantTimeCompare(clientSig, expectedSig) != 1 {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	if !user.Enabled {
		s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Fail")
		return nil, ErrNotAuthorized
	}

	accessToken, idToken, refreshToken, expiresIn, err := s.CreateTokens(reqCtx, userPoolID, user.ID, clientID, TokenGenerationAuthentication, parseClientMetadata(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}
	s.recordAuthEvent(reqCtx, userPoolID, user.ID, "SignIn", "Pass")

	return map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken":  accessToken,
			"IdToken":      idToken,
			"RefreshToken": refreshToken,
			"TokenType":    "Bearer",
			"ExpiresIn":    expiresIn,
		},
	}, nil
}

// computeSrpVerifier derives a fresh random salt and the matching SRP verifier
// for the supplied password. It must be invoked at every site that stores a
// password hash so that the user can later authenticate via USER_SRP_AUTH.
//
// userPoolID is the full Cognito pool ID (e.g. "us-east-1_abc123"); the part
// after the underscore (poolName) is required by Cognito's SRP variant inner
// hash. The returned saltHex and verifierHex are lowercase hex strings suitable
// for direct assignment to User.SrpSalt and User.SrpVerifier.
func computeSrpVerifier(userPoolID, username, password string) (saltHex, verifierHex string, err error) {
	idx := strings.Index(userPoolID, "_")
	if idx < 0 || idx == len(userPoolID)-1 {
		return "", "", fmt.Errorf("invalid user pool ID %q: missing region prefix", userPoolID)
	}
	poolName := userPoolID[idx+1:]
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}
	saltHex = hex.EncodeToString(salt)
	v := ComputeVerifier(saltHex, poolName, username, password)
	return saltHex, hex.EncodeToString(v.Bytes()), nil
}

// poolNameFromID extracts the portion of a Cognito user pool ID after the
// underscore (e.g. "us-east-1_abc123" => "abc123"). The pool name is used as
// part of the Cognito SRP inner hash and the claim message. The boolean is
// false when the ID does not contain a valid region/name separator.
func poolNameFromID(userPoolID string) (string, bool) {
	idx := strings.Index(userPoolID, "_")
	if idx < 0 || idx == len(userPoolID)-1 {
		return "", false
	}
	return userPoolID[idx+1:], true
}
