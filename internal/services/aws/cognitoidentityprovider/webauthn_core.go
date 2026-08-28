package cognitoidentityprovider

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// StartWebAuthnRegistrationInput carries the wire parameters of
// StartWebAuthnRegistration.
type StartWebAuthnRegistrationInput struct {
	AccessToken string
}

// CompleteWebAuthnRegistrationInput carries the wire parameters of
// CompleteWebAuthnRegistration. Params holds the raw request parameter map;
// the nested Credential structure is read from it inside the Core.
type CompleteWebAuthnRegistrationInput struct {
	AccessToken string
	Params      map[string]interface{}
}

// webauthnRegistrationSessionKey builds the deterministic challenge-session
// key for a user's pending WebAuthn registration. AWS binds the pending
// registration to the signed-in user server-side — one outstanding
// registration per user, the most recent Start wins.
func webauthnRegistrationSessionKey(userPoolID, userID string) string {
	return "webauthn-reg#" + userPoolID + "#" + userID
}

// webauthnCredParam is one entry of the credential creation options'
// pubKeyCredParams: the key type and COSE algorithm identifier the user
// pool accepts for passkey credentials.
type webauthnCredParam struct {
	Type string
	Alg  int64
}

// webauthnCredParams is the single source of the key algorithms the user
// pool offers at StartWebAuthnRegistration and accepts at
// CompleteWebAuthnRegistration (ES256 and RS256).
var webauthnCredParams = []webauthnCredParam{
	{Type: "public-key", Alg: -7},
	{Type: "public-key", Alg: -257},
}

// webauthnAlgSupported reports whether alg is one of the offered key
// algorithms.
func webauthnAlgSupported(alg int64) bool {
	for _, p := range webauthnCredParams {
		if p.Alg == alg {
			return true
		}
	}
	return false
}

// webauthnOriginAllowed reports whether an origin URL aligns with the user
// pool relying party id: per WebAuthn relying-party scoping, the origin's
// effective domain must be the RP ID itself or a subdomain of it.
func webauthnOriginAllowed(origin, rpID string) bool {
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	host := u.Hostname()
	return host == rpID || strings.HasSuffix(host, "."+rpID)
}

// ListWebAuthnCredentialsInput carries the wire parameters of
// ListWebAuthnCredentials. Params holds the raw request parameter map for
// the MaxResults member.
type ListWebAuthnCredentialsInput struct {
	AccessToken string
	NextToken   string
	Params      map[string]interface{}
}

// DeleteWebAuthnCredentialInput carries the wire parameters of
// DeleteWebAuthnCredential.
type DeleteWebAuthnCredentialInput struct {
	AccessToken  string
	CredentialID string
}

// startWebAuthnRegistrationCore starts a WebAuthn credential registration
// flow.
func (s *CognitoService) startWebAuthnRegistrationCore(reqCtx *request.RequestContext, in StartWebAuthnRegistrationInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	tokenRecord, err := s.validateAccessTokenRecord(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		return nil, ErrInternalError
	}
	challengeB64 := base64.RawURLEncoding.EncodeToString(challenge)

	// Store the challenge in a session for CompleteWebAuthnRegistration binding
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.GetUserByID(tokenRecord.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	sessionID := webauthnRegistrationSessionKey(user.UserPoolID, user.ID)
	challengeSession := &cognitostore.ChallengeSession{
		SessionID:  sessionID,
		UserPoolID: user.UserPoolID,
		// The pending registration is bound to the app client whose token
		// started it; Complete must arrive with the same client's token.
		ClientID:      tokenRecord.ClientID,
		Username:      user.Username,
		ChallengeName: "WEB_AUTHN_REGISTRATION",
		CreatedAt:     time.Now().UTC(),
		ExpiresAt:     time.Now().UTC().Add(srpChallengeSessionTTL),
		ChallengeData: challengeB64,
	}

	// Use custom domain as RP ID if configured, otherwise default to host.
	rpID := cognitoIdpHost(s.region)
	if domain, err := store.GetUserPoolDomainByPool(user.UserPoolID); err == nil && domain.Domain != "" {
		rpID = domain.Domain
	}
	// The completion verifies origin and rpIdHash against the exact RP ID
	// offered here, so persist it alongside the challenge.
	challengeSession.RelyingPartyID = rpID
	if err := store.SaveChallengeSession(challengeSession); err != nil {
		return nil, ErrInternalError
	}

	pubKeyCredParams := make([]map[string]interface{}, 0, len(webauthnCredParams))
	for _, p := range webauthnCredParams {
		pubKeyCredParams = append(pubKeyCredParams, map[string]interface{}{"type": p.Type, "alg": p.Alg})
	}

	options := map[string]interface{}{
		"challenge": challengeB64,
		"rp": map[string]interface{}{
			"name": "Cognito",
			"id":   rpID,
		},
		"user": map[string]interface{}{
			"id":          base64.RawURLEncoding.EncodeToString([]byte(user.ID)),
			"name":        user.ID,
			"displayName": user.ID,
		},
		"pubKeyCredParams": pubKeyCredParams,
		"timeout":          60000,
		"attestation":      "none",
		"authenticatorSelection": map[string]interface{}{
			"authenticatorAttachment": "platform",
			"userVerification":        "preferred",
		},
	}

	return map[string]interface{}{
		"CredentialCreationOptions": options,
	}, nil
}

// completeWebAuthnRegistrationCore completes a WebAuthn credential
// registration.
func (s *CognitoService) completeWebAuthnRegistrationCore(reqCtx *request.RequestContext, in CompleteWebAuthnRegistrationInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	tokenRecord, err := s.validateAccessTokenRecord(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(tokenRecord.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// The pending registration is bound to the signed-in user server-side;
	// no session identifier travels on the wire.
	sessionKey := webauthnRegistrationSessionKey(user.UserPoolID, user.ID)
	challengeSession, err := store.GetChallengeSession(sessionKey)
	if err != nil {
		return nil, ErrWebAuthnChallengeNotFound
	}

	// The registration must be completed through the same app client that
	// started it; a token issued to any other client is rejected.
	if challengeSession.ClientID != "" && challengeSession.ClientID != tokenRecord.ClientID {
		return nil, ErrWebAuthnClientMismatch
	}

	if challengeSession.ChallengeName != "WEB_AUTHN_REGISTRATION" {
		return nil, ErrInvalidParameter
	}
	if time.Now().UTC().After(challengeSession.ExpiresAt) {
		return nil, ErrWebAuthnChallengeNotFound
	}

	credentialRaw, ok := in.Params["Credential"]
	if !ok {
		return nil, ErrInvalidParameter
	}

	credentialBytes, _ := json.Marshal(credentialRaw)
	var credential struct {
		ID        string `json:"id"`
		PublicKey string `json:"publicKey"`
		Type      string `json:"type"`
		Response  struct {
			ClientDataJSON    string `json:"clientDataJSON"`
			AttestationObject string `json:"attestationObject"`
		} `json:"response"`
	}
	json.Unmarshal(credentialBytes, &credential)
	if credential.ID == "" {
		return nil, ErrInvalidParameter
	}

	// Verify the challenge and origin in clientDataJSON against the pending
	// registration: the challenge must match, and the origin's domain must
	// align with the user pool relying party id the options were issued under.
	rpID := challengeSession.RelyingPartyID
	if rpID == "" {
		// Sessions issued before the RP ID was persisted carry none; resolve
		// it exactly as Start does so verification stays well-defined.
		rpID = cognitoIdpHost(s.region)
		if domain, derr := store.GetUserPoolDomainByPool(user.UserPoolID); derr == nil && domain.Domain != "" {
			rpID = domain.Domain
		}
	}
	if credential.Response.ClientDataJSON == "" {
		return nil, ErrInvalidParameter
	}
	clientDataBytes, decErr := base64.RawURLEncoding.DecodeString(credential.Response.ClientDataJSON)
	if decErr != nil {
		clientDataBytes, decErr = base64.StdEncoding.DecodeString(credential.Response.ClientDataJSON)
		if decErr != nil {
			return nil, ErrInvalidParameter
		}
	}
	var clientData struct {
		Type      string `json:"type"`
		Challenge string `json:"challenge"`
		Origin    string `json:"origin"`
	}
	if err := json.Unmarshal(clientDataBytes, &clientData); err != nil {
		return nil, ErrInvalidParameter
	}
	if clientData.Challenge != challengeSession.ChallengeData {
		return nil, ErrNotAuthorized
	}
	if !webauthnOriginAllowed(clientData.Origin, rpID) {
		return nil, ErrWebAuthnOriginNotAllowed
	}

	// The credential must be scoped to the user pool's relying party: the
	// leading 32 bytes of the attestation's authenticator data are the
	// SHA-256 hash of the relying party id the authenticator signed for.
	attBytes, decErr := base64.RawURLEncoding.DecodeString(credential.Response.AttestationObject)
	if decErr != nil {
		if attBytes, decErr = base64.StdEncoding.DecodeString(credential.Response.AttestationObject); decErr != nil {
			return nil, ErrWebAuthnCredentialNotSupported
		}
	}
	var attestation struct {
		AuthData []byte `json:"authData"`
	}
	if err := cbor.Unmarshal(attBytes, &attestation); err != nil || len(attestation.AuthData) < 32 {
		return nil, ErrWebAuthnCredentialNotSupported
	}
	rpIDHash := sha256.Sum256([]byte(rpID))
	if !bytes.Equal(attestation.AuthData[:32], rpIDHash[:]) {
		return nil, ErrWebAuthnRelyingPartyMismatch
	}

	// The credential's COSE key must carry one of the algorithms the
	// credential creation options offered; map key 3 is the COSE algorithm
	// identifier.
	pubBytes, decErr := base64.RawURLEncoding.DecodeString(credential.PublicKey)
	if decErr != nil {
		if pubBytes, decErr = base64.StdEncoding.DecodeString(credential.PublicKey); decErr != nil {
			return nil, ErrWebAuthnCredentialNotSupported
		}
	}
	var coseKey map[int64]interface{}
	if err := cbor.Unmarshal(pubBytes, &coseKey); err != nil {
		return nil, ErrWebAuthnCredentialNotSupported
	}
	alg, ok := coseKey[3].(int64)
	if !ok || !webauthnAlgSupported(alg) {
		return nil, ErrWebAuthnCredentialNotSupported
	}

	// FriendlyCredentialName is "an automatically-generated friendly name
	// for the passkey credential" backed by StringType @length(0,131072)
	// with no pattern: the credential ID itself is the generated name, in
	// full — no truncation.
	cred := &cognitostore.WebAuthnCredential{
		CredentialID: credential.ID,
		FriendlyName: credential.ID,
		UserPoolID:   user.UserPoolID,
		UserID:       user.ID,
		PublicKey:    credential.PublicKey,
		CreatedAt:    time.Now().UTC(),
	}

	if err := store.CreateWebAuthnCredential(cred); err != nil {
		return nil, ErrInternalError
	}

	// The pending challenge is single-use; a failed delete would leave it
	// replayable until its TTL, so surface it as a warning.
	if err := store.DeleteChallengeSession(sessionKey); err != nil {
		log.Printf("warning: failed to consume the WebAuthn registration challenge session for user %s in pool %s: %v", user.ID, user.UserPoolID, err)
	}

	return response.EmptyResponse(), nil
}

// listWebAuthnCredentialsCore lists registered WebAuthn credentials.
func (s *CognitoService) listWebAuthnCredentialsCore(reqCtx *request.RequestContext, in ListWebAuthnCredentialsInput) (interface{}, error) {
	if in.AccessToken == "" {
		return nil, ErrInvalidParameter
	}

	userID, err := s.ValidateAccessToken(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	maxResults := maxWebAuthnCredentialListLimit
	if mr := request.GetIntParam(in.Params, "MaxResults"); mr > 0 {
		maxResults = mr
	}
	// Smithy WebAuthnCredentialsQueryLimitType: range {min: 0, max: 20}
	if maxResults > maxWebAuthnCredentialListLimit {
		return nil, ErrInvalidParameter
	}

	result, err := store.ListWebAuthnCredentialsPaginated(user.UserPoolID, user.ID, storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, c := range result.Items {
		formatted = append(formatted, map[string]interface{}{
			"CredentialId":           c.CredentialID,
			"FriendlyCredentialName": c.FriendlyName,
			"CreatedAt":              c.CreatedAt.Unix(),
		})
	}

	resp := map[string]interface{}{"Credentials": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// deleteWebAuthnCredentialCore deletes a WebAuthn credential.
func (s *CognitoService) deleteWebAuthnCredentialCore(reqCtx *request.RequestContext, in DeleteWebAuthnCredentialInput) (interface{}, error) {
	if in.AccessToken == "" || in.CredentialID == "" {
		return nil, ErrInvalidParameter
	}

	tokenRecord, err := s.validateAccessTokenRecord(reqCtx, in.AccessToken)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(tokenRecord.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Deleting an unregistered credential is a resource-not-found error,
	// not a silent success.
	if _, err := store.GetWebAuthnCredential(user.UserPoolID, user.ID, in.CredentialID); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.DeleteWebAuthnCredential(user.UserPoolID, user.ID, in.CredentialID); err != nil {
		return nil, ErrInternalError
	}

	return response.EmptyResponse(), nil
}
