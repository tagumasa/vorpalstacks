package cognitoidentityprovider

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"math/big"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fxamacker/cbor/v2"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// webAuthnRelyingPartyID resolves the relying party id that every WebAuthn
// ceremony for the pool runs under: the configured RelyingPartyId when set,
// otherwise the pool's hosted domain, otherwise the regional Cognito IDP
// host. Registration and assertion share this resolution so the rpIdHash a
// credential was created under can never drift from the rpId the assertion
// options carry.
func (s *CognitoService) webAuthnRelyingPartyID(store cognitostore.CognitoStoreInterface, pool *cognitostore.UserPool) string {
	if pool.WebAuthnConfiguration != nil && pool.WebAuthnConfiguration.RelyingPartyId != "" {
		return pool.WebAuthnConfiguration.RelyingPartyId
	}
	if domain, err := store.GetUserPoolDomainByPool(pool.ID); err == nil && domain.Domain != "" {
		return domain.Domain
	}
	// The fallback host is regional, so the region comes from the pool's own
	// ARN — a pool outside the service's constructor region must still
	// resolve its own regional host.
	region := s.region
	if parsed, err := svcarn.ParseARN(pool.Arn); err == nil && parsed.Region != "" {
		region = parsed.Region
	}
	return cognitoIdpHost(region)
}

// listUserWebAuthnCredentials walks every WebAuthn credential a user has
// registered in the pool. The store list is capped at the wire page size, so
// pages are followed to the end; the per-user count stays small. A page
// error fails closed: a partial list would both starve the assertion
// allow-list and flip the WebAuthn-enrolled decision for a user whose
// credentials are registered.
func listUserWebAuthnCredentials(store cognitostore.CognitoStoreInterface, userPoolID, userID string) ([]*cognitostore.WebAuthnCredential, error) {
	var creds []*cognitostore.WebAuthnCredential
	marker := ""
	for {
		result, err := store.ListWebAuthnCredentialsPaginated(userPoolID, userID, storecommon.ListOptions{
			MaxItems: maxWebAuthnCredentialListLimit,
			Marker:   marker,
		})
		if err != nil {
			return nil, err
		}
		creds = append(creds, result.Items...)
		if !result.IsTruncated || result.NextMarker == "" {
			return creds, nil
		}
		marker = result.NextMarker
	}
}

// webauthnAssertionTimeoutMs is the request timeout offered in the
// CREDENTIAL_REQUEST_OPTIONS, matching the documented example responses.
const webauthnAssertionTimeoutMs = 180000

// webAuthnUserVerification resolves the user-verification policy for
// assertion ceremonies: the pool's configured value, defaulting to
// "preferred".
func webAuthnUserVerification(pool *cognitostore.UserPool) string {
	if pool.WebAuthnConfiguration != nil && pool.WebAuthnConfiguration.UserVerification != "" {
		return pool.WebAuthnConfiguration.UserVerification
	}
	return "preferred"
}

// buildWebAuthnRequestOptions generates a fresh WebAuthn assertion challenge
// and renders the PublicKeyCredentialRequestOptions JSON that travels as the
// challenge's CREDENTIAL_REQUEST_OPTIONS parameter: the challenge, the
// timeout, the relying party id, the user's registered credentials as the
// allow-list, and the user-verification policy.
func (s *CognitoService) buildWebAuthnRequestOptions(store cognitostore.CognitoStoreInterface, pool *cognitostore.UserPool, creds []*cognitostore.WebAuthnCredential) (challengeB64, rpID, optionsJSON string, err error) {
	challenge := make([]byte, 32)
	if _, rerr := rand.Read(challenge); rerr != nil {
		return "", "", "", ErrInternalError
	}
	challengeB64 = base64.RawURLEncoding.EncodeToString(challenge)
	rpID = s.webAuthnRelyingPartyID(store, pool)

	allow := make([]map[string]interface{}, 0, len(creds))
	for _, c := range creds {
		allow = append(allow, map[string]interface{}{
			"type":       "public-key",
			"id":         c.CredentialID,
			"transports": []string{},
		})
	}
	options := map[string]interface{}{
		"challenge":        challengeB64,
		"timeout":          webauthnAssertionTimeoutMs,
		"rpId":             rpID,
		"allowCredentials": allow,
		"userVerification": webAuthnUserVerification(pool),
	}
	buf, merr := json.Marshal(options)
	if merr != nil {
		return "", "", "", ErrInternalError
	}
	return challengeB64, rpID, string(buf), nil
}

// challengeRoleMFA marks sessions issued by the MFA machinery as the
// sign-in's second factor (the empty default is the primary role).
const challengeRoleMFA = "MFA"

// issueWebAuthnChallenge mints a WEB_AUTHN challenge session carrying the
// assertion challenge and the relying party id, and renders the challenge
// response with the CREDENTIAL_REQUEST_OPTIONS parameter. The role records
// whether the challenge is the primary sign-in method (client-selected) or
// the second factor (MFA machinery).
func (s *CognitoService) issueWebAuthnChallenge(store cognitostore.CognitoStoreInterface, pool *cognitostore.UserPool, clientID string, user *cognitostore.User, role string) (map[string]interface{}, error) {
	creds, lerr := listUserWebAuthnCredentials(store, pool.ID, user.ID)
	if lerr != nil {
		return nil, ErrInternalError
	}
	if len(creds) == 0 {
		// The assertion allow-list is the user's registered credentials; a
		// WebAuthn challenge for a user without any is unanswerable.
		return nil, ErrNotAuthorized
	}
	challengeB64, rpID, optionsJSON, err := s.buildWebAuthnRequestOptions(store, pool, creds)
	if err != nil {
		return nil, err
	}
	session := newChallengeSession(pool.ID, clientID, user.Username, "WEB_AUTHN", challengeSessionTTL)
	session.ChallengeData = challengeB64
	session.RelyingPartyID = rpID
	session.ChallengeRole = role
	if err := store.SaveChallengeSession(session); err != nil {
		return nil, ErrInternalError
	}
	return map[string]interface{}{
		"ChallengeName": "WEB_AUTHN",
		"Session":       session.SessionID,
		"ChallengeParameters": map[string]string{
			"USERNAME":                   user.Username,
			"CREDENTIAL_REQUEST_OPTIONS": optionsJSON,
		},
	}, nil
}

// coseKeyInt reads a COSE map value as a signed integer. The CBOR decoder
// materialises positive integers as uint64 and negative ones as int64 when
// the target is interface{}, so both widths are accepted.
func coseKeyInt(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case uint64:
		return int64(n), true
	}
	return 0, false
}

// verifyECDSACOSESignature verifies an ES256 assertion signature in the COSE
// encoding RFC 9053 §2.1 defines — "the two integers are then concatenated
// together to form a byte string that is the resulting signature", each
// integer zero-padded to the curve's byte size, which for P-256 is the
// fixed-width 64-byte r||s form. ASN.1 DER is not the WebAuthn wire format
// and is not accepted.
func verifyECDSACOSESignature(pub *ecdsa.PublicKey, digest, sig []byte) bool {
	size := (pub.Curve.Params().BitSize + 7) / 8
	if len(sig) != 2*size {
		return false
	}
	r := new(big.Int).SetBytes(sig[:size])
	s := new(big.Int).SetBytes(sig[size:])
	return ecdsa.Verify(pub, digest, r, s)
}

// webauthnCredLocks serialises the per-credential signature-counter advance:
// the read-modify-write of the stored counter must be owned by one goroutine
// at a time, or two concurrent assertions both read the old count and one
// write is lost — regressing the clone-detection baseline. A fixed stripe
// array bounds the memory: the credential count is unbounded, the lock
// count must not be, and collisions merely share a mutex.
const webauthnCredLockStripes = 64

var webauthnCredLocks [webauthnCredLockStripes]sync.Mutex

// lockWebAuthnCredential returns the mutex owning one credential's counter
// read-modify-write.
func lockWebAuthnCredential(userPoolID, userID, credentialID string) *sync.Mutex {
	h := fnv.New32a()
	h.Write([]byte(userPoolID + "/" + userID + "/" + credentialID))
	return &webauthnCredLocks[h.Sum32()%webauthnCredLockStripes]
}

// webauthnAuthenticationResponseJSON is the W3C AuthenticationResponseJSON
// dict (WebAuthn-3 §5.1.1) that the client submits as the CREDENTIAL answer,
// with the binary members base64url-encoded.
type webauthnAuthenticationResponseJSON struct {
	ID       string `json:"id"`
	RawID    string `json:"rawId"`
	Type     string `json:"type"`
	Response struct {
		AuthenticatorData string `json:"authenticatorData"`
		ClientDataJSON    string `json:"clientDataJSON"`
		Signature         string `json:"signature"`
		UserHandle        string `json:"userHandle"`
	} `json:"response"`
}

// webauthnAuthDataFlagUP is the authenticator-data user-present flag; UV is
// the user-verified flag and AT the attested-credential-data-present flag
// (WebAuthn-3 §6.1).
const (
	webauthnAuthDataFlagUP byte = 0x01
	webauthnAuthDataFlagUV byte = 0x04
	webauthnAuthDataFlagAT byte = 0x40
)

// webauthnDecodeB64 decodes a WebAuthn binary member, accepting either
// base64url (the wire format) or standard base64.
func webauthnDecodeB64(v string) ([]byte, error) {
	if b, err := base64.RawURLEncoding.DecodeString(v); err == nil {
		return b, nil
	}
	return base64.StdEncoding.DecodeString(v)
}

// verifyWebAuthnAssertion verifies a CREDENTIAL answer against the challenge
// session that issued it, per the WebAuthn-3 §7.2 assertion verification:
// the client data must be a webauthn.get ceremony for this challenge from an
// origin aligned with the session's relying party id; the authenticator data
// must hash the relying party id and assert user presence (and user
// verification when the pool requires it); and the signature over
// authenticatorData || SHA-256(clientDataJSON) must verify with the stored
// credential's COSE public key. The signature counter must advance when both
// counters are non-zero — a non-advancing counter indicates a cloned
// authenticator and the assertion is rejected.
func verifyWebAuthnAssertion(
	store cognitostore.CognitoStoreInterface,
	session *cognitostore.ChallengeSession,
	pool *cognitostore.UserPool,
	user *cognitostore.User,
	credentialJSON string,
) error {
	var cred webauthnAuthenticationResponseJSON
	if err := json.Unmarshal([]byte(credentialJSON), &cred); err != nil {
		return ErrInvalidParameter
	}

	clientDataBytes, err := webauthnDecodeB64(cred.Response.ClientDataJSON)
	if err != nil {
		return ErrInvalidParameter
	}
	var clientData struct {
		Type      string `json:"type"`
		Challenge string `json:"challenge"`
		Origin    string `json:"origin"`
	}
	if err := json.Unmarshal(clientDataBytes, &clientData); err != nil {
		return ErrInvalidParameter
	}
	if clientData.Type != "webauthn.get" {
		return ErrNotAuthorized
	}
	if session.ChallengeData == "" || clientData.Challenge != session.ChallengeData {
		return ErrNotAuthorized
	}
	if !webauthnOriginAllowed(clientData.Origin, session.RelyingPartyID) {
		return ErrNotAuthorized
	}

	authData, err := webauthnDecodeB64(cred.Response.AuthenticatorData)
	if err != nil || len(authData) < 37 {
		return ErrInvalidParameter
	}
	rpIDHash := sha256.Sum256([]byte(session.RelyingPartyID))
	if !bytes.Equal(authData[:32], rpIDHash[:]) {
		return ErrNotAuthorized
	}
	if authData[32]&webauthnAuthDataFlagUP == 0 {
		return ErrNotAuthorized
	}
	if webAuthnUserVerification(pool) == "required" && authData[32]&webauthnAuthDataFlagUV == 0 {
		return ErrNotAuthorized
	}

	// The response identifies a public-key credential, and its id and rawId
	// members carry the same credential identifier (the JSON string and the
	// base64url-encoded bytes); a divergence is a malformed or forged
	// response.
	if cred.Type != "" && cred.Type != "public-key" {
		return ErrNotAuthorized
	}
	credentialID := cred.ID
	if credentialID == "" {
		credentialID = cred.RawID
	}
	if credentialID == "" {
		return ErrInvalidParameter
	}
	if cred.ID != "" && cred.RawID != "" && cred.ID != cred.RawID {
		return ErrNotAuthorized
	}
	stored, err := store.GetWebAuthnCredential(pool.ID, user.ID, credentialID)
	if err != nil {
		return ErrNotAuthorized
	}
	// When the authenticator echoes a user handle it must identify the same
	// user the challenge was issued for.
	if cred.Response.UserHandle != "" {
		handle, derr := webauthnDecodeB64(cred.Response.UserHandle)
		if derr != nil || !bytes.Equal(handle, []byte(user.ID)) {
			return ErrNotAuthorized
		}
	}

	sig, err := webauthnDecodeB64(cred.Response.Signature)
	if err != nil {
		return ErrInvalidParameter
	}
	pubBytes, err := webauthnDecodeB64(stored.PublicKey)
	if err != nil {
		return ErrInternalError
	}
	var coseKey map[int64]interface{}
	if err := cbor.Unmarshal(pubBytes, &coseKey); err != nil {
		return ErrInternalError
	}
	alg, algOK := coseKeyInt(coseKey[3])
	if !algOK || !webauthnAlgSupported(alg) {
		return ErrNotAuthorized
	}
	// The assertion signature input is authenticatorData || SHA-256(
	// clientDataJSON); ES256 and RS256 both sign the SHA-256 digest of
	// that input.
	clientDataHash := sha256.Sum256(clientDataBytes)
	message := make([]byte, 0, len(authData)+len(clientDataHash))
	message = append(message, authData...)
	message = append(message, clientDataHash[:]...)
	digest := sha256.Sum256(message)
	kty, ktyOK := coseKeyInt(coseKey[1])
	if !ktyOK {
		return ErrNotAuthorized
	}
	switch kty {
	case 2: // EC2
		if alg != -7 {
			return ErrNotAuthorized
		}
		crv, crvOK := coseKeyInt(coseKey[-1])
		if !crvOK || crv != 1 { // P-256
			return ErrNotAuthorized
		}
		x, okX := coseKey[-2].([]byte)
		y, okY := coseKey[-3].([]byte)
		if !okX || !okY || len(x) != 32 || len(y) != 32 {
			return ErrNotAuthorized
		}
		pub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: new(big.Int).SetBytes(x), Y: new(big.Int).SetBytes(y)}
		if !verifyECDSACOSESignature(pub, digest[:], sig) {
			return ErrNotAuthorized
		}
	case 3: // RSA
		if alg != -257 {
			return ErrNotAuthorized
		}
		n, okN := coseKey[-1].([]byte)
		e, okE := coseKey[-2].([]byte)
		if !okN || !okE || len(n) == 0 || len(e) == 0 {
			return ErrNotAuthorized
		}
		exponent := 0
		for _, b := range e {
			exponent = exponent<<8 | int(b)
		}
		if exponent < 3 || exponent%2 == 0 {
			return ErrNotAuthorized
		}
		pub := &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: exponent}
		if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig); err != nil {
			return ErrNotAuthorized
		}
	default:
		return ErrNotAuthorized
	}

	// Signature counter: a non-advancing counter on a counter-equipped
	// authenticator signals a cloned credential. Counters advance (and are
	// persisted) whenever the new value is non-zero. The check-and-advance is
	// a read-modify-write over the stored credential, owned by a
	// per-credential mutex and re-read under the lock, so concurrent
	// assertions cannot interleave and lose an advance; a failed persist is
	// logged rather than swallowed — the clone-detection baseline must not
	// regress silently.
	assertionCount := binary.BigEndian.Uint32(authData[33:37])
	if assertionCount != 0 || stored.SignCount != 0 {
		mu := lockWebAuthnCredential(pool.ID, user.ID, credentialID)
		mu.Lock()
		current, gerr := store.GetWebAuthnCredential(pool.ID, user.ID, credentialID)
		if gerr != nil {
			mu.Unlock()
			return ErrNotAuthorized
		}
		nonAdvancing := assertionCount != 0 && current.SignCount != 0 && assertionCount <= current.SignCount
		if !nonAdvancing && assertionCount > current.SignCount {
			current.SignCount = assertionCount
			if uerr := store.UpdateWebAuthnCredential(current); uerr != nil {
				logs.Warn("cognito webauthn signature counter persist failed", logs.String("pool", pool.ID), logs.String("user", user.ID), logs.String("credential", credentialID), logs.Err(uerr))
			}
		}
		mu.Unlock()
		if nonAdvancing {
			return ErrNotAuthorized
		}
	}
	return nil
}

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

// webauthnCredentialIDMaxLength is the WebAuthn credential-ID bound: "At
// most 1023 bytes long" (WebAuthn-3 Terminology).
const webauthnCredentialIDMaxLength = 1023

// attestedCredentialData is the credential material the authenticator
// attested inside registration authenticator data: the credential ID in its
// raw byte form and its base64url encoding, and the COSE public key as the
// exact attested byte string.
type attestedCredentialData struct {
	CredentialIDRaw []byte
	CredentialID    string
	COSEKey         []byte
}

// parseAttestedCredentialData extracts the attested credential data that
// follows the 37-byte authenticator-data header: AAGUID (16 bytes), a
// big-endian 2-byte credential-ID length, the credential ID, then exactly
// one CBOR item — the COSE key — taken as raw bytes so the stored key is the
// attested byte string rather than a re-encoding of it.
func parseAttestedCredentialData(authData []byte) (attestedCredentialData, bool) {
	var out attestedCredentialData
	if len(authData) < 37+16+2 {
		return out, false
	}
	rest := authData[53:]
	idLen := int(binary.BigEndian.Uint16(rest[:2]))
	if idLen == 0 || idLen > webauthnCredentialIDMaxLength {
		return out, false
	}
	rest = rest[2:]
	if len(rest) < idLen {
		return out, false
	}
	out.CredentialIDRaw = rest[:idLen]
	out.CredentialID = base64.RawURLEncoding.EncodeToString(rest[:idLen])
	dec := cbor.NewDecoder(bytes.NewReader(rest[idLen:]))
	var raw cbor.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return out, false
	}
	out.COSEKey = raw
	return out, true
}

// webauthnOriginAllowed reports whether an origin URL aligns with the user
// pool relying party id: per WebAuthn relying-party scoping, the origin's
// effective domain must be the RP ID itself or a subdomain of it, and the
// origin must be a secure context — https, or http on a loopback host (the
// one insecure origin browsers treat as potentially trustworthy).
func webauthnOriginAllowed(origin, rpID string) bool {
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}
	host := u.Hostname()
	if host != rpID && !strings.HasSuffix(host, "."+rpID) {
		return false
	}
	if u.Scheme == "https" {
		return true
	}
	return u.Scheme == "http" && (host == "localhost" || host == "127.0.0.1" || host == "::1")
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

	// The registration and assertion ceremonies run under the same relying
	// party id; persist it alongside the challenge so the completion
	// verifies origin and rpIdHash against the exact value offered here.
	// Registration additionally demands the passkey feature: the pool must
	// carry a WebAuthnConfiguration, and the relying party id must come from
	// that configuration or the pool's hosted domain — the regional-host
	// fallback serves assertion resolution only, never a registration.
	pool, err := store.GetUserPool(user.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}
	if pool.WebAuthnConfiguration == nil {
		return nil, ErrWebAuthnNotEnabled
	}
	rpID := pool.WebAuthnConfiguration.RelyingPartyId
	if rpID == "" {
		if domain, derr := store.GetUserPoolDomainByPool(pool.ID); derr == nil && domain.Domain != "" {
			rpID = domain.Domain
		}
	}
	if rpID == "" {
		return nil, ErrWebAuthnConfigurationMissing
	}
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
			"id":   challengeSession.RelyingPartyID,
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
			// The pool's configured user-verification policy governs the
			// registration ceremony just as it governs assertions.
			"userVerification": webAuthnUserVerification(pool),
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

	// A completion on a pool whose passkey feature is not enabled — never
	// configured, or removed after the registration started — registers
	// nothing.
	if pool, perr := store.GetUserPool(user.UserPoolID); perr != nil {
		return nil, ErrResourceNotFound
	} else if pool.WebAuthnConfiguration == nil {
		return nil, ErrWebAuthnNotEnabled
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
	// The pending registration carries the same failed-attempt budget as
	// every other challenge answer, so completion guessing stays bounded.
	if challengeSession.FailedAttempts >= maxChallengeAttempts {
		return nil, ErrNotAuthorized
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

	// Ceremony failures consume the pending registration's attempt budget,
	// like every other challenge answer.
	fail := func(err error) (interface{}, error) {
		recordChallengeFailure(store, challengeSession)
		return nil, err
	}

	// Verify the challenge and origin in clientDataJSON against the pending
	// registration: the ceremony must be a webauthn.create for this
	// challenge, and the origin's domain must align with the user pool
	// relying party id the options were issued under. Start persists that id
	// with the session; a session without one is malformed and its ceremony
	// cannot be verified.
	rpID := challengeSession.RelyingPartyID
	if rpID == "" {
		return nil, ErrWebAuthnChallengeNotFound
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
	if clientData.Type != "webauthn.create" {
		return fail(ErrNotAuthorized)
	}
	if clientData.Challenge != challengeSession.ChallengeData {
		return fail(ErrNotAuthorized)
	}
	if !webauthnOriginAllowed(clientData.Origin, rpID) {
		return fail(ErrWebAuthnOriginNotAllowed)
	}

	// The credential must be scoped to the user pool's relying party: the
	// leading 32 bytes of the attestation's authenticator data are the
	// SHA-256 hash of the relying party id the authenticator signed for.
	attBytes, decErr := base64.RawURLEncoding.DecodeString(credential.Response.AttestationObject)
	if decErr != nil {
		if attBytes, decErr = base64.StdEncoding.DecodeString(credential.Response.AttestationObject); decErr != nil {
			return fail(ErrWebAuthnCredentialNotSupported)
		}
	}
	var attestation struct {
		AuthData []byte `json:"authData"`
	}
	if err := cbor.Unmarshal(attBytes, &attestation); err != nil || len(attestation.AuthData) < 37 {
		return fail(ErrWebAuthnCredentialNotSupported)
	}
	rpIDHash := sha256.Sum256([]byte(rpID))
	if !bytes.Equal(attestation.AuthData[:32], rpIDHash[:]) {
		return fail(ErrWebAuthnRelyingPartyMismatch)
	}
	// The registration ceremony asserts user presence, and the credential
	// material is authoritative only inside the attested credential data the
	// authenticator embedded: the UP and AT flags must both be set, and the
	// credential ID and COSE key below are taken from the attestation, not
	// from the client-asserted JSON copies.
	flags := attestation.AuthData[32]
	if flags&webauthnAuthDataFlagUP == 0 || flags&webauthnAuthDataFlagAT == 0 {
		return fail(ErrWebAuthnCredentialNotSupported)
	}
	attested, aok := parseAttestedCredentialData(attestation.AuthData)
	if !aok {
		return fail(ErrWebAuthnCredentialNotSupported)
	}

	// The attested credential's COSE key must carry one of the algorithms
	// the credential creation options offered; map key 3 is the COSE
	// algorithm identifier, read with the decoder-width-tolerant helper the
	// assertion path uses.
	var attestedKey map[int64]interface{}
	if err := cbor.Unmarshal(attested.COSEKey, &attestedKey); err != nil {
		return fail(ErrWebAuthnCredentialNotSupported)
	}
	alg, algOK := coseKeyInt(attestedKey[3])
	if !algOK || !webauthnAlgSupported(alg) {
		return fail(ErrWebAuthnCredentialNotSupported)
	}

	// The client-asserted id is a copy of the attested credential ID, in
	// either the raw-byte string form or the base64url encoding W3C JSON
	// responses use; anything else means the registration response was
	// assembled rather than produced by an authenticator.
	if credential.ID != string(attested.CredentialIDRaw) && credential.ID != attested.CredentialID {
		return fail(ErrWebAuthnCredentialNotSupported)
	}
	if credential.PublicKey != "" {
		clientKey, kerr := webauthnDecodeB64(credential.PublicKey)
		if kerr != nil || !bytes.Equal(clientKey, attested.COSEKey) {
			return fail(ErrWebAuthnCredentialNotSupported)
		}
	}

	// FriendlyCredentialName is "an automatically-generated friendly name
	// for the passkey credential" backed by StringType @length(0,131072)
	// with no pattern: the credential ID itself is the generated name, in
	// full — no truncation. The id is the client-asserted copy the
	// cross-check bound to the attested material, and the stored key is the
	// attested byte string.
	cred := &cognitostore.WebAuthnCredential{
		CredentialID: credential.ID,
		FriendlyName: credential.ID,
		UserPoolID:   user.UserPoolID,
		UserID:       user.ID,
		PublicKey:    base64.RawURLEncoding.EncodeToString(attested.COSEKey),
		CreatedAt:    time.Now().UTC(),
	}

	if err := store.CreateWebAuthnCredential(cred); err != nil {
		return nil, ErrInternalError
	}

	// The pending challenge is single-use; a failed delete would leave it
	// replayable until its TTL, so surface it as a warning.
	if err := store.DeleteChallengeSession(sessionKey); err != nil {
		logs.Warn("cognito webauthn registration challenge session consumption failed", logs.String("user", user.ID), logs.String("pool", user.UserPoolID), logs.Err(err))
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

	// Smithy WebAuthnCredentialsQueryLimitType: range {min: 0, max: 20}
	maxResults, err := parseListLimit(in.Params, "MaxResults", maxWebAuthnCredentialListLimit)
	if err != nil {
		return nil, err
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
