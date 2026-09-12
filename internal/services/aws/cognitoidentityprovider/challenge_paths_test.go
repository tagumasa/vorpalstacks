package cognitoidentityprovider

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"

	"github.com/fxamacker/cbor/v2"
)

// This file pins the challenge-answer paths that were previously advertised
// but unanswerable: the PASSWORD_SRP SRP leg, the WEB_AUTHN assertion
// ceremony, and the device SRP challenges that replace the MFA challenge.

// testSrpClient computes the client-side SRP values with the package's own
// primitives, mirroring the reference client maths: A from a random scalar,
// the shared secret from the server's B, and the HMAC claim over
// idA || idB || secretBlock || timestamp. The identity pair is pool name and
// username for password proofs, device group key and device key for device
// proofs.
type testSrpClient struct {
	a      *big.Int
	bigA   *big.Int
	salt   string
	idA    string
	idB    string
	secret string
}

func newTestSrpClient(t *testing.T, salt, idA, idB, secret string) *testSrpClient {
	t.Helper()
	rb := make([]byte, 32)
	if _, err := rand.Read(rb); err != nil {
		t.Fatal(err)
	}
	c := &testSrpClient{
		a:      new(big.Int).Mod(new(big.Int).SetBytes(rb), cognitoSrpN),
		salt:   salt,
		idA:    idA,
		idB:    idB,
		secret: secret,
	}
	c.bigA = new(big.Int).Exp(cognitoSrpG, c.a, cognitoSrpN)
	return c
}

func (c *testSrpClient) srpAHex() string { return c.bigA.Text(16) }

// claim computes the PASSWORD_CLAIM_SIGNATURE for the server's challenge
// parameters.
func (c *testSrpClient) claim(t *testing.T, srpBHex, secretBlockB64 string, now time.Time) (sig, ts string) {
	t.Helper()
	B := mustHexToBig(srpBHex)
	saltInt, _ := new(big.Int).SetString(c.salt, 16)
	userPassHash := hashSHA256Hex([]byte(c.idA + c.idB + ":" + c.secret))
	u := mustHexToBig(hexHash(padHex(c.bigA.Text(16)) + padHex(B.Text(16))))
	x := mustHexToBig(hexHash(padHex(saltInt.Text(16)) + userPassHash))
	gx := new(big.Int).Exp(cognitoSrpG, x, cognitoSrpN)
	base := new(big.Int).Sub(B, new(big.Int).Mul(cognitoSrpK, gx))
	base.Mod(base, cognitoSrpN)
	exp := new(big.Int).Add(c.a, new(big.Int).Mul(u, x))
	S := new(big.Int).Exp(base, exp, cognitoSrpN)
	K := computeHKDF(padHex(S.Text(16)), padHex(u.Text(16)))
	// The SDKs sign the TIMESTAMP claim with the 24-hour layout the server
	// parses for freshness; a 12-hour rendering would sit twelve hours off
	// the server clock for afternoon claims.
	ts = now.In(time.UTC).Format("Mon Jan 2 15:04:05 MST 2006")
	secretBlock, err := base64.StdEncoding.DecodeString(secretBlockB64)
	if err != nil {
		t.Fatal(err)
	}
	mac := hmac.New(sha256.New, K)
	mac.Write([]byte(c.idA + c.idB + string(secretBlock) + ts))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), ts
}

// testWebAuthnKey generates a P-256 key pair and the COSE EC2 public key
// encoding the server stores at registration.
func testWebAuthnKey(t *testing.T) (*ecdsa.PrivateKey, string) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	x := make([]byte, 32)
	y := make([]byte, 32)
	priv.PublicKey.X.FillBytes(x)
	priv.PublicKey.Y.FillBytes(y)
	cose, err := cbor.Marshal(map[int64]interface{}{
		1:  int64(2),  // kty: EC2
		3:  int64(-7), // alg: ES256
		-1: int64(1),  // crv: P-256
		-2: x,         // x coordinate
		-3: y,         // y coordinate
	})
	if err != nil {
		t.Fatal(err)
	}
	return priv, base64.RawURLEncoding.EncodeToString(cose)
}

// testWebAuthnAssertion builds the AuthenticationResponseJSON the client
// submits as CREDENTIAL: an assertion over the request options' challenge
// with the given signature counter, signed in the COSE raw fixed-width
// r||s encoding (RFC 9053 §2.1) that is the WebAuthn wire format.
func testWebAuthnAssertion(t *testing.T, priv *ecdsa.PrivateKey, credentialID, rpID, challenge, origin string, signCount uint32) string {
	t.Helper()
	return testWebAuthnAssertionSign(t, priv, credentialID, rpID, challenge, origin, signCount, false)
}

// testWebAuthnAssertionSign is testWebAuthnAssertion with the signature
// encoding selectable: the COSE raw r||s wire format, or ASN.1 DER for the
// rejection pin.
func testWebAuthnAssertionSign(t *testing.T, priv *ecdsa.PrivateKey, credentialID, rpID, challenge, origin string, signCount uint32, der bool) string {
	t.Helper()
	rpIDHash := sha256.Sum256([]byte(rpID))
	authData := make([]byte, 37)
	copy(authData, rpIDHash[:])
	authData[32] = webauthnAuthDataFlagUP
	binary.BigEndian.PutUint32(authData[33:37], signCount)

	clientData, err := json.Marshal(map[string]string{
		"type":      "webauthn.get",
		"challenge": challenge,
		"origin":    origin,
	})
	if err != nil {
		t.Fatal(err)
	}
	clientDataHash := sha256.Sum256(clientData)
	message := append(append([]byte{}, authData...), clientDataHash[:]...)
	digest := sha256.Sum256(message)
	var sig []byte
	if der {
		if sig, err = ecdsa.SignASN1(rand.Reader, priv, digest[:]); err != nil {
			t.Fatal(err)
		}
	} else {
		r, s, serr := ecdsa.Sign(rand.Reader, priv, digest[:])
		if serr != nil {
			t.Fatal(serr)
		}
		sig = make([]byte, 64)
		r.FillBytes(sig[:32])
		s.FillBytes(sig[32:])
	}

	assertion := map[string]interface{}{
		"id":    credentialID,
		"rawId": credentialID,
		"type":  "public-key",
		"response": map[string]string{
			"authenticatorData": base64.RawURLEncoding.EncodeToString(authData),
			"clientDataJSON":    base64.RawURLEncoding.EncodeToString(clientData),
			"signature":         base64.RawURLEncoding.EncodeToString(sig),
		},
	}
	buf, err := json.Marshal(assertion)
	if err != nil {
		t.Fatal(err)
	}
	return string(buf)
}

// testWebAuthnRequestOptions parses the CREDENTIAL_REQUEST_OPTIONS
// parameter of a WEB_AUTHN challenge.
func testWebAuthnRequestOptions(t *testing.T, resp map[string]interface{}) (challenge, rpID string, allowIDs []string) {
	t.Helper()
	params, _ := resp["ChallengeParameters"].(map[string]string)
	optionsJSON, ok := params["CREDENTIAL_REQUEST_OPTIONS"]
	if !ok {
		t.Fatalf("challenge carries no CREDENTIAL_REQUEST_OPTIONS: %#v", resp)
	}
	var options struct {
		Challenge        string           `json:"challenge"`
		RpID             string           `json:"rpId"`
		AllowCredentials []map[string]any `json:"allowCredentials"`
	}
	if err := json.Unmarshal([]byte(optionsJSON), &options); err != nil {
		t.Fatalf("parse CREDENTIAL_REQUEST_OPTIONS: %v", err)
	}
	for _, c := range options.AllowCredentials {
		if id, ok := c["id"].(string); ok {
			allowIDs = append(allowIDs, id)
		}
	}
	return options.Challenge, options.RpID, allowIDs
}

// seedSrpCredentials gives the environment user a native SRP verifier for
// the given password.
func seedSrpCredentials(t *testing.T, env *challengeTestEnv, password string) {
	t.Helper()
	saltHex, verifierHex, err := computeSrpVerifier(env.pool.ID, env.user.Username, password)
	if err != nil {
		t.Fatal(err)
	}
	env.user.SrpSalt = saltHex
	env.user.SrpVerifier = verifierHex
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}
}

// initiatePasswordSrp drives the PASSWORD_SRP answer against a PASSWORD_SRP
// session.
func respondToChallenge(env *challengeTestEnv, session, name string, responses map[string]interface{}) (interface{}, error) {
	return env.svc.RespondToAuthChallenge(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"ClientId":           challengeTestClientID,
		"ChallengeName":      name,
		"Session":            session,
		"ChallengeResponses": responses,
	}))
}

// A PREFERRED_CHALLENGE of PASSWORD_SRP issues the PASSWORD_SRP challenge
// shell; its SRP_A answer returns the PASSWORD_VERIFIER challenge whose
// claim completes the sign-in.
func TestPasswordSrpPreferredChallengeRoundTrip(t *testing.T) {
	env := newChallengeTestEnv(t)
	const password = "OldPass123!"
	seedSrpCredentials(t, env, password)

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":            "victim",
			"PREFERRED_CHALLENGE": "PASSWORD_SRP",
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	if m["ChallengeName"] != "PASSWORD_SRP" {
		t.Fatalf("expected PASSWORD_SRP challenge, got %#v", resp)
	}
	srpSession, _ := m["Session"].(string)

	client := newTestSrpClient(t, env.user.SrpSalt, poolIDName(t, env.pool.ID), "victim", password)
	verifierResp, err := respondToChallenge(env, srpSession, "PASSWORD_SRP", map[string]interface{}{
		"USERNAME": "victim",
		"SRP_A":    client.srpAHex(),
	})
	if err != nil {
		t.Fatal(err)
	}
	vm, _ := verifierResp.(map[string]interface{})
	if vm["ChallengeName"] != "PASSWORD_VERIFIER" {
		t.Fatalf("expected PASSWORD_VERIFIER after SRP_A, got %#v", verifierResp)
	}
	cp, _ := vm["ChallengeParameters"].(map[string]interface{})
	salt, _ := cp["SALT"].(string)
	srpB, _ := cp["SRP_B"].(string)
	secretBlock, _ := cp["SECRET_BLOCK"].(string)
	if salt == "" || srpB == "" || secretBlock == "" {
		t.Fatalf("missing verifier challenge parameters: %#v", cp)
	}

	// The PASSWORD_SRP session was burned; the claim answers the fresh
	// PASSWORD_VERIFIER session.
	sig, ts := client.claim(t, srpB, secretBlock, time.Now())
	authResp, err := respondToChallenge(env, vm["Session"].(string), "PASSWORD_VERIFIER", map[string]interface{}{
		"USERNAME":                    "victim",
		"PASSWORD_CLAIM_SIGNATURE":    sig,
		"PASSWORD_CLAIM_SECRET_BLOCK": secretBlock,
		"TIMESTAMP":                   ts,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := authResp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens after PASSWORD_VERIFIER, got %#v", authResp)
	}
}

// PREFERRED_CHALLENGE for a challenge the user cannot answer falls back to
// the selector response instead of issuing an unanswerable challenge.
func TestUserAuthPreferredChallengeFallsBackToSelector(t *testing.T) {
	env := newChallengeTestEnv(t)

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":            "victim",
			"PREFERRED_CHALLENGE": "WEB_AUTHN",
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	if _, ok := m["AvailableChallenges"]; !ok {
		t.Fatalf("expected selector response for unavailable preference, got %#v", resp)
	}
	for _, c := range m["AvailableChallenges"].([]string) {
		if c == "WEB_AUTHN" {
			t.Fatal("WEB_AUTHN advertised for a user without passkey credentials")
		}
	}
}

// The SELECT_CHALLENGE answer completes PASSWORD_SRP inline: SRP_A travels
// with the selection and the response is the PASSWORD_VERIFIER challenge.
func TestSelectChallengeInlinePasswordSrp(t *testing.T) {
	env := newChallengeTestEnv(t)
	const password = "OldPass123!"
	seedSrpCredentials(t, env, password)

	resp, err := initiateUserAuth(env, "victim")
	if err != nil {
		t.Fatal(err)
	}
	selector, _ := resp.(map[string]interface{})["Session"].(string)

	client := newTestSrpClient(t, env.user.SrpSalt, poolIDName(t, env.pool.ID), "victim", password)
	verifierResp, err := respondToChallenge(env, selector, "SELECT_CHALLENGE", map[string]interface{}{
		"USERNAME": "victim",
		"ANSWER":   "PASSWORD_SRP",
		"SRP_A":    client.srpAHex(),
	})
	if err != nil {
		t.Fatal(err)
	}
	vm, _ := verifierResp.(map[string]interface{})
	if vm["ChallengeName"] != "PASSWORD_VERIFIER" {
		t.Fatalf("expected PASSWORD_VERIFIER from inline SRP_A, got %#v", verifierResp)
	}
	cp, _ := vm["ChallengeParameters"].(map[string]interface{})
	sig, ts := client.claim(t, cp["SRP_B"].(string), cp["SECRET_BLOCK"].(string), time.Now())
	authResp, err := respondToChallenge(env, vm["Session"].(string), "PASSWORD_VERIFIER", map[string]interface{}{
		"USERNAME":                    "victim",
		"PASSWORD_CLAIM_SIGNATURE":    sig,
		"PASSWORD_CLAIM_SECRET_BLOCK": cp["SECRET_BLOCK"].(string),
		"TIMESTAMP":                   ts,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := authResp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens, got %#v", authResp)
	}
}

// The SELECT_CHALLENGE answer completes PASSWORD inline: the password
// travels with the selection and tokens issue in the selection response.
func TestSelectChallengeInlinePassword(t *testing.T) {
	env := newChallengeTestEnv(t)

	resp, err := initiateUserAuth(env, "victim")
	if err != nil {
		t.Fatal(err)
	}
	selector, _ := resp.(map[string]interface{})["Session"].(string)

	authResp, err := respondToChallenge(env, selector, "SELECT_CHALLENGE", map[string]interface{}{
		"USERNAME": "victim",
		"ANSWER":   "PASSWORD",
		"PASSWORD": "OldPass123!",
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := authResp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens from inline PASSWORD selection, got %#v", authResp)
	}

	// A wrong inline password is rejected and burns nothing it should not:
	// the selector session was consumed by the attempt.
	resp, err = initiateUserAuth(env, "victim")
	if err != nil {
		t.Fatal(err)
	}
	selector, _ = resp.(map[string]interface{})["Session"].(string)
	if _, err := respondToChallenge(env, selector, "SELECT_CHALLENGE", map[string]interface{}{
		"USERNAME": "victim",
		"ANSWER":   "PASSWORD",
		"PASSWORD": "WrongPass999!",
	}); err == nil {
		t.Fatal("inline PASSWORD selection accepted a wrong password")
	}
}

// poolIDName extracts the pool-name suffix used in the SRP inner hash.
func poolIDName(t *testing.T, poolID string) string {
	t.Helper()
	name, ok := poolNameFromID(poolID)
	if !ok {
		t.Fatalf("pool id %q has no name suffix", poolID)
	}
	return name
}

// The WEB_AUTHN challenge carries the assertion request options bound to
// the session; the CREDENTIAL answer verifies against them and completes
// the sign-in. A cloned-authenticator signature counter (non-advancing) is
// rejected on the next sign-in.
func TestWebAuthnPreferredChallengeAssertionRoundTrip(t *testing.T) {
	env := newChallengeTestEnv(t)
	const rpID = "auth.example.com"
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: rpID}
	})
	priv, coseB64 := testWebAuthnKey(t)
	if err := env.store.CreateWebAuthnCredential(&cognitostore.WebAuthnCredential{
		CredentialID: "unit-passkey-1",
		UserPoolID:   env.pool.ID,
		UserID:       env.user.ID,
		PublicKey:    coseB64,
	}); err != nil {
		t.Fatal(err)
	}

	initiate := func() map[string]interface{} {
		resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
			"AuthFlow": "USER_AUTH",
			"ClientId": challengeTestClientID,
			"AuthParameters": map[string]interface{}{
				"USERNAME":            "victim",
				"PREFERRED_CHALLENGE": "WEB_AUTHN",
			},
		}))
		if err != nil {
			t.Fatal(err)
		}
		m, _ := resp.(map[string]interface{})
		if m["ChallengeName"] != "WEB_AUTHN" {
			t.Fatalf("expected WEB_AUTHN challenge, got %#v", resp)
		}
		return m
	}

	m := initiate()
	challenge, gotRpID, allowIDs := testWebAuthnRequestOptions(t, m)
	if gotRpID != rpID {
		t.Fatalf("request options rpId %q, want %q", gotRpID, rpID)
	}
	if len(allowIDs) != 1 || allowIDs[0] != "unit-passkey-1" {
		t.Fatalf("allowCredentials %v, want [unit-passkey-1]", allowIDs)
	}

	assertion := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, challenge, "https://"+rpID, 1)
	authResp, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": assertion,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := authResp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens from WEB_AUTHN assertion, got %#v", authResp)
	}

	// The signature counter advanced to 1 with the verified assertion; a
	// second sign-in replaying counter 1 is a cloned authenticator and is
	// rejected, while counter 2 succeeds.
	m = initiate()
	challenge, _, _ = testWebAuthnRequestOptions(t, m)
	replayed := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, challenge, "https://"+rpID, 1)
	if _, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": replayed,
	}); err == nil {
		t.Fatal("WEB_AUTHN assertion with non-advancing signature counter accepted")
	}

	m = initiate()
	challenge, _, _ = testWebAuthnRequestOptions(t, m)
	fresh := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, challenge, "https://"+rpID, 2)
	if _, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": fresh,
	}); err != nil {
		t.Fatalf("WEB_AUTHN assertion with advancing counter rejected: %v", err)
	}
}

// A WEB_AUTHN assertion that answers a different challenge, comes from a
// foreign origin, or signs with the wrong key is rejected.
func TestWebAuthnAssertionRejectsForgedAnswers(t *testing.T) {
	env := newChallengeTestEnv(t)
	const rpID = "auth.example.com"
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: rpID}
	})
	priv, coseB64 := testWebAuthnKey(t)
	otherPriv, _ := testWebAuthnKey(t)
	if err := env.store.CreateWebAuthnCredential(&cognitostore.WebAuthnCredential{
		CredentialID: "unit-passkey-1",
		UserPoolID:   env.pool.ID,
		UserID:       env.user.ID,
		PublicKey:    coseB64,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":            "victim",
			"PREFERRED_CHALLENGE": "WEB_AUTHN",
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	challenge, _, _ := testWebAuthnRequestOptions(t, m)

	wrongChallenge := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, "some-other-challenge", "https://"+rpID, 1)
	if _, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": wrongChallenge,
	}); err == nil {
		t.Fatal("assertion for a foreign challenge accepted")
	}

	foreignOrigin := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, challenge, "https://evil.example.com", 1)
	if _, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": foreignOrigin,
	}); err == nil {
		t.Fatal("assertion from a foreign origin accepted")
	}

	wrongKey := testWebAuthnAssertion(t, otherPriv, "unit-passkey-1", rpID, challenge, "https://"+rpID, 1)
	if _, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": wrongKey,
	}); err == nil {
		t.Fatal("assertion signed with a foreign key accepted")
	}

	// A plain-http origin on a non-loopback host is not a secure context,
	// however well the rest of the assertion is formed.
	httpOrigin := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, challenge, "http://"+rpID, 1)
	if _, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": httpOrigin,
	}); err == nil {
		t.Fatal("assertion from a plain-http origin accepted")
	}
}

// The ES256 assertion signature travels in the COSE raw fixed-width r||s
// encoding; the ASN.1 DER form is not the WebAuthn wire format and the
// assertion is rejected.
func TestWebAuthnAssertionRejectsDERSignature(t *testing.T) {
	env := newChallengeTestEnv(t)
	const rpID = "auth.example.com"
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: rpID}
	})
	priv, coseB64 := testWebAuthnKey(t)
	if err := env.store.CreateWebAuthnCredential(&cognitostore.WebAuthnCredential{
		CredentialID: "unit-passkey-1",
		UserPoolID:   env.pool.ID,
		UserID:       env.user.ID,
		PublicKey:    coseB64,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":            "victim",
			"PREFERRED_CHALLENGE": "WEB_AUTHN",
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	challenge, _, _ := testWebAuthnRequestOptions(t, m)

	der := testWebAuthnAssertionSign(t, priv, "unit-passkey-1", rpID, challenge, "https://"+rpID, 1, true)
	if _, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": der,
	}); err == nil {
		t.Fatal("WEB_AUTHN assertion accepted an ASN.1 DER signature")
	}
}

// The selector response carries the assertion request options bound to the
// SELECT_CHALLENGE session, so the WEB_AUTHN selection completes inline.
func TestWebAuthnSelectionInlineCompletion(t *testing.T) {
	env := newChallengeTestEnv(t)
	const rpID = "auth.example.com"
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: rpID}
	})
	priv, coseB64 := testWebAuthnKey(t)
	if err := env.store.CreateWebAuthnCredential(&cognitostore.WebAuthnCredential{
		CredentialID: "unit-passkey-1",
		UserPoolID:   env.pool.ID,
		UserID:       env.user.ID,
		PublicKey:    coseB64,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := initiateUserAuth(env, "victim")
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	advertised := false
	for _, c := range m["AvailableChallenges"].([]string) {
		if c == "WEB_AUTHN" {
			advertised = true
		}
	}
	if !advertised {
		t.Fatalf("WEB_AUTHN not advertised: %#v", m["AvailableChallenges"])
	}
	challenge, _, _ := testWebAuthnRequestOptions(t, m)
	selector, _ := m["Session"].(string)

	assertion := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, challenge, "https://"+rpID, 1)
	authResp, err := respondToChallenge(env, selector, "SELECT_CHALLENGE", map[string]interface{}{
		"USERNAME":   "victim",
		"ANSWER":     "WEB_AUTHN",
		"CREDENTIAL": assertion,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := authResp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens from inline WEB_AUTHN selection, got %#v", authResp)
	}
}

// A passkey enrolled as the user's MFA factor challenges with WEB_AUTHN
// after the password verifies, and the assertion completes the sign-in
// without a further challenge.
func TestWebAuthnAsMfaChallenge(t *testing.T) {
	env := newChallengeTestEnv(t)
	const rpID = "auth.example.com"
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: rpID}
		p.MfaConfiguration = "ON"
	})
	priv, coseB64 := testWebAuthnKey(t)
	if err := env.store.CreateWebAuthnCredential(&cognitostore.WebAuthnCredential{
		CredentialID: "unit-passkey-1",
		UserPoolID:   env.pool.ID,
		UserID:       env.user.ID,
		PublicKey:    coseB64,
	}); err != nil {
		t.Fatal(err)
	}
	env.user.WebAuthnMfaEnabled = true
	if err := env.store.UpdateUser(env.user); err != nil {
		t.Fatal(err)
	}

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME": "victim",
			"PASSWORD": "OldPass123!",
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	if m["ChallengeName"] != "WEB_AUTHN" {
		t.Fatalf("expected WEB_AUTHN as the MFA challenge, got %#v", resp)
	}
	challenge, _, _ := testWebAuthnRequestOptions(t, m)

	assertion := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, challenge, "https://"+rpID, 1)
	authResp, err := respondToChallenge(env, m["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": assertion,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := authResp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens after the MFA assertion, got %#v", authResp)
	}
}

// testWebAuthnRegistrationCredential builds the registration Credential JSON
// the client submits: a "none"-format attestation whose authenticator data
// hashes the relying party id, carries the given flags byte, and embeds
// attested credential data with the attested credential ID and COSE key. The
// JSON-level id and publicKey copies carry clientID and coseB64, so a
// cross-check rejection is built by passing values that differ from the
// attested ones.
func testWebAuthnRegistrationCredential(t *testing.T, coseB64, attestedID, clientID, rpID, challenge, origin, clientType string, flags byte) string {
	t.Helper()
	cose, err := base64.RawURLEncoding.DecodeString(coseB64)
	if err != nil {
		t.Fatal(err)
	}
	rpIDHash := sha256.Sum256([]byte(rpID))
	idRaw := []byte(attestedID)
	var idLen [2]byte
	binary.BigEndian.PutUint16(idLen[:], uint16(len(idRaw)))
	authData := make([]byte, 0, 55+len(idRaw)+len(cose))
	authData = append(authData, rpIDHash[:]...)
	authData = append(authData, flags)
	authData = append(authData, 0, 0, 0, 0)
	authData = append(authData, make([]byte, 16)...)
	authData = append(authData, idLen[:]...)
	authData = append(authData, idRaw...)
	authData = append(authData, cose...)

	clientData, err := json.Marshal(map[string]string{
		"type":      clientType,
		"challenge": challenge,
		"origin":    origin,
	})
	if err != nil {
		t.Fatal(err)
	}
	attObj, err := cbor.Marshal(map[string]interface{}{
		"fmt":      "none",
		"attStmt":  map[string]interface{}{},
		"authData": authData,
	})
	if err != nil {
		t.Fatal(err)
	}
	buf, err := json.Marshal(map[string]interface{}{
		"id":        clientID,
		"publicKey": coseB64,
		"type":      "public-key",
		"response": map[string]string{
			"clientDataJSON":    base64.RawURLEncoding.EncodeToString(clientData),
			"attestationObject": base64.RawURLEncoding.EncodeToString(attObj),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return string(buf)
}

// The registration ceremony verifies the attested credential data: the
// stored credential's id and key come from the attestation object (flags UP
// and AT set), cross-checked against the client-asserted copies; the
// ceremony type, the origin scheme and the challenge are enforced, the
// creation options carry the pool's user-verification policy, and completion
// guessing is bounded by the pending registration's attempt budget.
func TestWebAuthnRegistrationCeremony(t *testing.T) {
	env := newChallengeTestEnv(t)
	const rpID = "auth.example.com"
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: rpID}
	})
	priv, coseB64 := testWebAuthnKey(t)

	accessToken, _, _, _, err := env.svc.CreateTokens(env.reqCtx, env.pool.ID, env.user.ID, challengeTestClientID, TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}
	start := func(wantUserVerification string) (challenge string) {
		t.Helper()
		resp, serr := env.svc.startWebAuthnRegistrationCore(env.reqCtx, StartWebAuthnRegistrationInput{AccessToken: accessToken})
		if serr != nil {
			t.Fatal(serr)
		}
		opts := resp.(map[string]interface{})["CredentialCreationOptions"].(map[string]interface{})
		if id, _ := opts["rp"].(map[string]interface{})["id"].(string); id != rpID {
			t.Fatalf("creation options rp.id %v, want %q", opts["rp"], rpID)
		}
		sel, _ := opts["authenticatorSelection"].(map[string]interface{})
		if sel == nil || sel["userVerification"] != wantUserVerification {
			t.Fatalf("authenticatorSelection %v, want userVerification %q", opts["authenticatorSelection"], wantUserVerification)
		}
		return opts["challenge"].(string)
	}
	complete := func(credentialJSON string) error {
		t.Helper()
		var obj interface{}
		if err := json.Unmarshal([]byte(credentialJSON), &obj); err != nil {
			t.Fatal(err)
		}
		_, cerr := env.svc.completeWebAuthnRegistrationCore(env.reqCtx, CompleteWebAuthnRegistrationInput{
			AccessToken: accessToken,
			Params:      map[string]interface{}{"Credential": obj},
		})
		return cerr
	}
	good := func(challenge string) string {
		return testWebAuthnRegistrationCredential(t, coseB64, "reg-passkey-1", "reg-passkey-1", rpID, challenge, "https://"+rpID, "webauthn.create", webauthnAuthDataFlagUP|webauthnAuthDataFlagAT)
	}

	// The pool's unconfigured user-verification policy defaults to
	// "preferred" in the creation options; a configured value is carried
	// through.
	start("preferred")
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration.UserVerification = "required"
	})
	start("required")
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration.UserVerification = ""
	})

	// A well-formed attestation registers: the stored id, friendly name and
	// key are the attested material, and the stored key verifies a later
	// assertion.
	challenge := start("preferred")
	if err := complete(good(challenge)); err != nil {
		t.Fatalf("well-formed registration rejected: %v", err)
	}
	if err := complete(good(challenge)); err == nil {
		t.Fatal("second completion of a consumed registration accepted")
	}
	stored, err := env.store.GetWebAuthnCredential(env.pool.ID, env.user.ID, "reg-passkey-1")
	if err != nil {
		t.Fatal(err)
	}
	if stored.FriendlyName != "reg-passkey-1" || stored.PublicKey != coseB64 {
		t.Fatalf("stored credential %q/%q, want the attested material", stored.FriendlyName, stored.PublicKey)
	}
	initResp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":            "victim",
			"PREFERRED_CHALLENGE": "WEB_AUTHN",
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	im, _ := initResp.(map[string]interface{})
	assertionChallenge, _, _ := testWebAuthnRequestOptions(t, im)
	assertion := testWebAuthnAssertion(t, priv, "reg-passkey-1", rpID, assertionChallenge, "https://"+rpID, 1)
	authResp, err := respondToChallenge(env, im["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": assertion,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := authResp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens from the registered credential's assertion, got %#v", authResp)
	}

	// A client-asserted id that differs from the attested credential ID is
	// an assembled response, not an authenticator product.
	if err := complete(testWebAuthnRegistrationCredential(t, coseB64, "attested-id", "forged-id", rpID, start("preferred"), "https://"+rpID, "webauthn.create", webauthnAuthDataFlagUP|webauthnAuthDataFlagAT)); err == nil {
		t.Fatal("client-asserted id diverging from the attested id accepted")
	}
	// The ceremony type must be webauthn.create.
	if err := complete(testWebAuthnRegistrationCredential(t, coseB64, "reg-passkey-2", "reg-passkey-2", rpID, start("preferred"), "https://"+rpID, "webauthn.get", webauthnAuthDataFlagUP|webauthnAuthDataFlagAT)); err == nil {
		t.Fatal("webauthn.get-typed registration accepted")
	}
	// User presence and attested credential data are both required.
	if err := complete(testWebAuthnRegistrationCredential(t, coseB64, "reg-passkey-3", "reg-passkey-3", rpID, start("preferred"), "https://"+rpID, "webauthn.create", webauthnAuthDataFlagAT)); err == nil {
		t.Fatal("registration without user presence accepted")
	}
	if err := complete(testWebAuthnRegistrationCredential(t, coseB64, "reg-passkey-4", "reg-passkey-4", rpID, start("preferred"), "https://"+rpID, "webauthn.create", webauthnAuthDataFlagUP)); err == nil {
		t.Fatal("registration without attested credential data accepted")
	}
	// A plain-http non-loopback origin is not a secure context.
	if err := complete(testWebAuthnRegistrationCredential(t, coseB64, "reg-passkey-5", "reg-passkey-5", rpID, start("preferred"), "http://"+rpID, "webauthn.create", webauthnAuthDataFlagUP|webauthnAuthDataFlagAT)); err == nil {
		t.Fatal("registration from a plain-http origin accepted")
	}
	// The ceremony must answer the pending challenge.
	if err := complete(testWebAuthnRegistrationCredential(t, coseB64, "reg-passkey-6", "reg-passkey-6", rpID, "foreign-challenge", "https://"+rpID, "webauthn.create", webauthnAuthDataFlagUP|webauthnAuthDataFlagAT)); err == nil {
		t.Fatal("registration answering a foreign challenge accepted")
	}

	// Completion guessing is bounded: the attempt budget exhausts the
	// pending registration and a subsequent valid completion finds nothing.
	challenge = start("preferred")
	for i := 0; i < maxChallengeAttempts; i++ {
		if err := complete(testWebAuthnRegistrationCredential(t, coseB64, "reg-passkey-7", "reg-passkey-7", rpID, "wrong-challenge", "https://"+rpID, "webauthn.create", webauthnAuthDataFlagUP|webauthnAuthDataFlagAT)); err == nil {
			t.Fatalf("wrong-challenge completion %d accepted", i)
		}
	}
	if err := complete(good(challenge)); err == nil {
		t.Fatal("valid completion accepted after the attempt budget was exhausted")
	}
}

// Concurrent assertions on one credential cannot lose a counter advance: the
// stored counter ends at the highest presented value, and a replayed value is
// rejected afterwards.
func TestWebAuthnSignatureCounterConcurrentAdvance(t *testing.T) {
	env := newChallengeTestEnv(t)
	const rpID = "auth.example.com"
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: rpID}
	})
	priv, coseB64 := testWebAuthnKey(t)
	if err := env.store.CreateWebAuthnCredential(&cognitostore.WebAuthnCredential{
		CredentialID: "unit-passkey-1",
		UserPoolID:   env.pool.ID,
		UserID:       env.user.ID,
		PublicKey:    coseB64,
	}); err != nil {
		t.Fatal(err)
	}

	session := &cognitostore.ChallengeSession{
		ChallengeData:  "concurrent-challenge",
		RelyingPartyID: rpID,
	}
	const rounds = 16
	assertions := make([]string, rounds)
	for i := range assertions {
		assertions[i] = testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, "concurrent-challenge", "https://"+rpID, uint32(i+1))
	}
	var wg sync.WaitGroup
	for _, a := range assertions {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			// A goroutine whose counter was overtaken by a larger value is
			// legitimately rejected; what must never happen is a lost
			// advance — the stored counter regressing below the maximum.
			_ = verifyWebAuthnAssertion(env.store, session, env.pool, env.user, a)
		}(a)
	}
	wg.Wait()

	stored, err := env.store.GetWebAuthnCredential(env.pool.ID, env.user.ID, "unit-passkey-1")
	if err != nil {
		t.Fatal(err)
	}
	if stored.SignCount != rounds {
		t.Fatalf("stored signature counter %d after %d concurrent assertions, want %d — an advance was lost", stored.SignCount, rounds, rounds)
	}

	replay := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, "concurrent-challenge", "https://"+rpID, rounds)
	if err := verifyWebAuthnAssertion(env.store, session, env.pool, env.user, replay); err == nil {
		t.Fatal("non-advancing signature counter accepted after concurrent assertions")
	}
}

// Selecting WEB_AUTHN in a SELECT_CHALLENGE response without an inline
// CREDENTIAL answer issues the WEB_AUTHN challenge with fresh assertion
// options; the follow-up assertion answers the new session.
func TestSelectChallengeWebAuthnIssuesAssertionOptions(t *testing.T) {
	env := newChallengeTestEnv(t)
	const rpID = "auth.example.com"
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{RelyingPartyId: rpID}
	})
	priv, coseB64 := testWebAuthnKey(t)
	if err := env.store.CreateWebAuthnCredential(&cognitostore.WebAuthnCredential{
		CredentialID: "unit-passkey-1",
		UserPoolID:   env.pool.ID,
		UserID:       env.user.ID,
		PublicKey:    coseB64,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := initiateUserAuth(env, "victim")
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	selector, _ := m["Session"].(string)

	authResp, err := respondToChallenge(env, selector, "SELECT_CHALLENGE", map[string]interface{}{
		"USERNAME": "victim",
		"ANSWER":   "WEB_AUTHN",
	})
	if err != nil {
		t.Fatal(err)
	}
	am, _ := authResp.(map[string]interface{})
	if am["ChallengeName"] != "WEB_AUTHN" {
		t.Fatalf("expected WEB_AUTHN after the selection, got %#v", authResp)
	}
	challenge, gotRpID, _ := testWebAuthnRequestOptions(t, am)
	if gotRpID != rpID {
		t.Fatalf("request options rpId %q, want %q", gotRpID, rpID)
	}

	assertion := testWebAuthnAssertion(t, priv, "unit-passkey-1", rpID, challenge, "https://"+rpID, 1)
	final, err := respondToChallenge(env, am["Session"].(string), "WEB_AUTHN", map[string]interface{}{
		"USERNAME":   "victim",
		"CREDENTIAL": assertion,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := final.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens from the follow-up WEB_AUTHN assertion, got %#v", final)
	}
}

// seedRememberedDevice registers a remembered device whose verifier derives
// from the documented device secret construction:
// x = H(salt || H(groupKey || deviceKey || ":" || deviceSecret)).
func seedRememberedDevice(t *testing.T, env *challengeTestEnv, deviceKey, deviceSecret string) {
	t.Helper()
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		t.Fatal(err)
	}
	saltHex := hex.EncodeToString(saltBytes)
	groupKey := poolIDName(t, env.pool.ID)
	verifier := ComputeVerifier(saltHex, groupKey, deviceKey, deviceSecret)
	if err := env.store.CreateDevice(&cognitostore.Device{
		DeviceKey:              deviceKey,
		UserPoolID:             env.pool.ID,
		UserID:                 env.user.ID,
		DeviceSecretVerifierB:  verifier.Text(16),
		DeviceSaltVerifier:     saltHex,
		DeviceRememberedStatus: "remembered",
	}); err != nil {
		t.Fatal(err)
	}
}

// A remembered device replaces the pool's MFA challenge with device SRP:
// the password sign-in continues with DEVICE_SRP_AUTH, the SRP_A answer
// returns DEVICE_PASSWORD_VERIFIER, and the claim over the device identity
// completes the sign-in and records the device authentication.
func TestDeviceSrpAuthReplacesMfaChallenge(t *testing.T) {
	env := newChallengeTestEnv(t)
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.MfaConfiguration = "ON"
		p.DeviceConfiguration = &cognitostore.DeviceConfiguration{ChallengeRequiredOnNewDevice: true}
	})
	const deviceKey = "us-east-1_device-1"
	const deviceSecret = "device-secret-value"
	seedRememberedDevice(t, env, deviceKey, deviceSecret)
	groupKey := poolIDName(t, env.pool.ID)

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":   "victim",
			"PASSWORD":   "OldPass123!",
			"DEVICE_KEY": deviceKey,
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	if m["ChallengeName"] != "DEVICE_SRP_AUTH" {
		t.Fatalf("expected DEVICE_SRP_AUTH to replace the MFA challenge, got %#v", resp)
	}
	deviceSession, _ := m["Session"].(string)

	// A DEVICE_SRP session bound to the device cannot be answered for a
	// different device key.
	client := newTestSrpClient(t, deviceSaltFor(t, env, deviceKey), groupKey, deviceKey, deviceSecret)
	if _, err := respondToChallenge(env, deviceSession, "DEVICE_SRP_AUTH", map[string]interface{}{
		"USERNAME":   "victim",
		"DEVICE_KEY": "us-east-1_other-device",
		"SRP_A":      client.srpAHex(),
	}); err == nil {
		t.Fatal("DEVICE_SRP_AUTH answered for a foreign device key")
	}

	verifierResp, err := respondToChallenge(env, deviceSession, "DEVICE_SRP_AUTH", map[string]interface{}{
		"USERNAME":   "victim",
		"DEVICE_KEY": deviceKey,
		"SRP_A":      client.srpAHex(),
	})
	if err != nil {
		t.Fatal(err)
	}
	vm, _ := verifierResp.(map[string]interface{})
	if vm["ChallengeName"] != "DEVICE_PASSWORD_VERIFIER" {
		t.Fatalf("expected DEVICE_PASSWORD_VERIFIER, got %#v", verifierResp)
	}
	cp, _ := vm["ChallengeParameters"].(map[string]interface{})
	if cp["SECRET_BLOCK"] == "" || cp["SRP_B"] == "" {
		t.Fatalf("missing device verifier challenge parameters: %#v", cp)
	}

	sig, ts := client.claim(t, cp["SRP_B"].(string), cp["SECRET_BLOCK"].(string), time.Now())
	authResp, err := respondToChallenge(env, vm["Session"].(string), "DEVICE_PASSWORD_VERIFIER", map[string]interface{}{
		"USERNAME":                    "victim",
		"DEVICE_KEY":                  deviceKey,
		"PASSWORD_CLAIM_SIGNATURE":    sig,
		"PASSWORD_CLAIM_SECRET_BLOCK": cp["SECRET_BLOCK"].(string),
		"TIMESTAMP":                   ts,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ar, ok := authResp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{}); !ok || ar["AccessToken"] == "" {
		t.Fatalf("expected tokens after the device proof, got %#v", authResp)
	}
	device, err := env.store.GetDevice(env.pool.ID, env.user.ID, deviceKey)
	if err != nil {
		t.Fatal(err)
	}
	if device.DeviceLastAuthenticatedDate.IsZero() {
		t.Fatal("device authentication not recorded")
	}
}

// A device proof computed from the wrong device secret fails with the
// uniform NotAuthorized error.
func TestDeviceSrpAuthRejectsWrongSecret(t *testing.T) {
	env := newChallengeTestEnv(t)
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.MfaConfiguration = "ON"
		p.DeviceConfiguration = &cognitostore.DeviceConfiguration{ChallengeRequiredOnNewDevice: true}
	})
	const deviceKey = "us-east-1_device-2"
	seedRememberedDevice(t, env, deviceKey, "real-secret")

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":   "victim",
			"PASSWORD":   "OldPass123!",
			"DEVICE_KEY": deviceKey,
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	if m["ChallengeName"] != "DEVICE_SRP_AUTH" {
		t.Fatalf("expected DEVICE_SRP_AUTH, got %#v", resp)
	}
	groupKey := poolIDName(t, env.pool.ID)
	client := newTestSrpClient(t, deviceSaltFor(t, env, deviceKey), groupKey, deviceKey, "wrong-secret")

	verifierResp, err := respondToChallenge(env, m["Session"].(string), "DEVICE_SRP_AUTH", map[string]interface{}{
		"USERNAME":   "victim",
		"DEVICE_KEY": deviceKey,
		"SRP_A":      client.srpAHex(),
	})
	if err != nil {
		t.Fatal(err)
	}
	vm, _ := verifierResp.(map[string]interface{})
	cp, _ := vm["ChallengeParameters"].(map[string]interface{})
	sig, ts := client.claim(t, cp["SRP_B"].(string), cp["SECRET_BLOCK"].(string), time.Now())
	if _, err := respondToChallenge(env, vm["Session"].(string), "DEVICE_PASSWORD_VERIFIER", map[string]interface{}{
		"USERNAME":                    "victim",
		"DEVICE_KEY":                  deviceKey,
		"PASSWORD_CLAIM_SIGNATURE":    sig,
		"PASSWORD_CLAIM_SECRET_BLOCK": cp["SECRET_BLOCK"].(string),
		"TIMESTAMP":                   ts,
	}); err == nil || !strings.Contains(err.Error(), "NotAuthorized") {
		t.Fatalf("expected NotAuthorized for a wrong device secret, got %v", err)
	}
}

// The freshness window admits an honestly clocked client and refuses stale,
// future-dated and malformed claims; the unpadded day-of-month the SDKs emit
// must parse.
func TestFreshSrpTimestampWindow(t *testing.T) {
	now := time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)
	fresh := func(d time.Duration) string {
		return now.Add(d).Format(srpTimestampLayout)
	}
	if !freshSrpTimestamp(fresh(0), now) {
		t.Fatal("current timestamp refused")
	}
	if !freshSrpTimestamp(fresh(-25*time.Second), now) {
		t.Fatal("timestamp 25s in the past refused")
	}
	if !freshSrpTimestamp(fresh(25*time.Second), now) {
		t.Fatal("timestamp 25s in the future refused")
	}
	singleDigitDay := time.Date(2026, 9, 3, 9, 5, 7, 0, time.UTC)
	if !freshSrpTimestamp(singleDigitDay.Format(srpTimestampLayout), singleDigitDay) {
		t.Fatal("unpadded day-of-month refused")
	}
	for name, ts := range map[string]string{
		"stale":     fresh(-srpTimestampWindow - time.Second),
		"future":    fresh(srpTimestampWindow + time.Second),
		"twelve-ho": now.Add(-12 * time.Hour).Format(srpTimestampLayout),
		"malformed": "not a timestamp",
		"rfc3339":   now.Format(time.RFC3339),
	} {
		if freshSrpTimestamp(ts, now) {
			t.Fatalf("%s timestamp accepted", name)
		}
	}
}

// A correctly signed but stale PASSWORD_VERIFIER claim is refused with the
// NotAuthorizedException the model documents for the exceeded response
// period, closing replay of a captured claim signature.
func TestPasswordVerifierStaleTimestampRefused(t *testing.T) {
	env := newChallengeTestEnv(t)
	const password = "OldPass123!"
	seedSrpCredentials(t, env, password)

	resp, err := initiateUserAuth(env, "victim")
	if err != nil {
		t.Fatal(err)
	}
	selector, _ := resp.(map[string]interface{})["Session"].(string)

	client := newTestSrpClient(t, env.user.SrpSalt, poolIDName(t, env.pool.ID), "victim", password)
	verifierResp, err := respondToChallenge(env, selector, "SELECT_CHALLENGE", map[string]interface{}{
		"USERNAME": "victim",
		"ANSWER":   "PASSWORD_SRP",
		"SRP_A":    client.srpAHex(),
	})
	if err != nil {
		t.Fatal(err)
	}
	vm, _ := verifierResp.(map[string]interface{})
	cp, _ := vm["ChallengeParameters"].(map[string]interface{})

	sig, ts := client.claim(t, cp["SRP_B"].(string), cp["SECRET_BLOCK"].(string), time.Now().Add(-5*time.Minute))
	if _, err := respondToChallenge(env, vm["Session"].(string), "PASSWORD_VERIFIER", map[string]interface{}{
		"USERNAME":                    "victim",
		"PASSWORD_CLAIM_SIGNATURE":    sig,
		"PASSWORD_CLAIM_SECRET_BLOCK": cp["SECRET_BLOCK"].(string),
		"TIMESTAMP":                   ts,
	}); err == nil || !strings.Contains(err.Error(), "NotAuthorized") {
		t.Fatalf("stale TIMESTAMP claim returned %v, want NotAuthorizedException", err)
	}
}

// The device SRP claim carries the same freshness contract: a stale claim is
// refused even though its signature verifies against the device secret.
func TestDeviceSrpStaleTimestampRefused(t *testing.T) {
	env := newChallengeTestEnv(t)
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.MfaConfiguration = "ON"
		p.DeviceConfiguration = &cognitostore.DeviceConfiguration{ChallengeRequiredOnNewDevice: true}
	})
	const deviceKey = "us-east-1_device-stale"
	const deviceSecret = "device-secret-value"
	seedRememberedDevice(t, env, deviceKey, deviceSecret)

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":   "victim",
			"PASSWORD":   "OldPass123!",
			"DEVICE_KEY": deviceKey,
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	deviceSession, _ := m["Session"].(string)

	client := newTestSrpClient(t, deviceSaltFor(t, env, deviceKey), poolIDName(t, env.pool.ID), deviceKey, deviceSecret)
	verifierResp, err := respondToChallenge(env, deviceSession, "DEVICE_SRP_AUTH", map[string]interface{}{
		"USERNAME":   "victim",
		"DEVICE_KEY": deviceKey,
		"SRP_A":      client.srpAHex(),
	})
	if err != nil {
		t.Fatal(err)
	}
	vm, _ := verifierResp.(map[string]interface{})
	cp, _ := vm["ChallengeParameters"].(map[string]interface{})

	sig, ts := client.claim(t, cp["SRP_B"].(string), cp["SECRET_BLOCK"].(string), time.Now().Add(-5*time.Minute))
	if _, err := respondToChallenge(env, vm["Session"].(string), "DEVICE_PASSWORD_VERIFIER", map[string]interface{}{
		"USERNAME":                    "victim",
		"DEVICE_KEY":                  deviceKey,
		"PASSWORD_CLAIM_SIGNATURE":    sig,
		"PASSWORD_CLAIM_SECRET_BLOCK": cp["SECRET_BLOCK"].(string),
		"TIMESTAMP":                   ts,
	}); err == nil || !strings.Contains(err.Error(), "NotAuthorized") {
		t.Fatalf("stale device TIMESTAMP claim returned %v, want NotAuthorizedException", err)
	}
}

func deviceSaltFor(t *testing.T, env *challengeTestEnv, deviceKey string) string {
	t.Helper()
	device, err := env.store.GetDevice(env.pool.ID, env.user.ID, deviceKey)
	if err != nil {
		t.Fatal(err)
	}
	return device.DeviceSaltVerifier
}

// A sign-in in a device-tracking pool mints the NewDeviceMetadata the
// client confirms with ConfirmDevice; a sign-in naming a remembered device
// gets no new device key.
func TestNewDeviceMetadataOnDeviceTrackingSignIn(t *testing.T) {
	env := newChallengeTestEnv(t)
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.DeviceConfiguration = &cognitostore.DeviceConfiguration{}
	})

	signIn := func(deviceKey string) map[string]interface{} {
		params := map[string]interface{}{
			"USERNAME": "victim",
			"PASSWORD": "OldPass123!",
		}
		if deviceKey != "" {
			params["DEVICE_KEY"] = deviceKey
		}
		resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
			"AuthFlow":       "USER_PASSWORD_AUTH",
			"ClientId":       challengeTestClientID,
			"AuthParameters": params,
		}))
		if err != nil {
			t.Fatal(err)
		}
		ar, ok := resp.(map[string]interface{})["AuthenticationResult"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected tokens, got %#v", resp)
		}
		return ar
	}

	ar := signIn("")
	meta, ok := ar["NewDeviceMetadata"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected NewDeviceMetadata on a device-tracking sign-in, got %#v", ar)
	}
	if meta["DeviceGroupKey"] != poolIDName(t, env.pool.ID) {
		t.Fatalf("DeviceGroupKey %v, want the pool name suffix", meta["DeviceGroupKey"])
	}
	if deviceKey, _ := meta["DeviceKey"].(string); !strings.HasPrefix(deviceKey, "us-east-1_") {
		t.Fatalf("DeviceKey %v not in the documented region_UUID format", meta["DeviceKey"])
	}

	const deviceKey = "us-east-1_device-3"
	seedRememberedDevice(t, env, deviceKey, "secret-3")
	ar = signIn(deviceKey)
	if _, ok := ar["NewDeviceMetadata"]; ok {
		t.Fatalf("identified remembered device got a new device key: %#v", ar["NewDeviceMetadata"])
	}
}

// Without ChallengeRequiredOnNewDevice the regular MFA challenge stays in
// place for a remembered device: substituting the device SRP proof for the
// second factor is that flag's documented effect, not a default of device
// tracking.
func TestDeviceSecondFactorRequiresChallengeFlag(t *testing.T) {
	env := newChallengeTestEnv(t)
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.MfaConfiguration = "ON"
		p.DeviceConfiguration = &cognitostore.DeviceConfiguration{}
	})
	const deviceKey = "us-east-1_device-noflag"
	seedRememberedDevice(t, env, deviceKey, "secret-noflag")

	resp, err := env.svc.InitiateAuth(context.Background(), env.reqCtx, challengeReq(map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": challengeTestClientID,
		"AuthParameters": map[string]interface{}{
			"USERNAME":   "victim",
			"PASSWORD":   "OldPass123!",
			"DEVICE_KEY": deviceKey,
		},
	}))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := resp.(map[string]interface{})
	if m["ChallengeName"] == "DEVICE_SRP_AUTH" {
		t.Fatal("remembered device replaced the MFA challenge without ChallengeRequiredOnNewDevice")
	}
	if m["ChallengeName"] != "MFA_SETUP" {
		t.Fatalf("expected the regular MFA challenge to stay in place, got %#v", resp)
	}
}

// ConfirmDevice follows the remembering axis the pool configures: under
// DeviceOnlyRememberedOnUserPrompt the device is created not_remembered and
// the response asks for the user's confirmation (UpdateDeviceStatus flips
// it); without the flag the device is immediately remembered. A duplicate
// device key is rejected, and device attributes pass the pool schema
// validation every attribute write passes.
func TestConfirmDeviceRememberingAxis(t *testing.T) {
	env := newChallengeTestEnv(t)
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.DeviceConfiguration = &cognitostore.DeviceConfiguration{DeviceOnlyRememberedOnUserPrompt: true}
	})
	accessToken, _, _, _, err := env.svc.CreateTokens(env.reqCtx, env.pool.ID, env.user.ID, challengeTestClientID, TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}
	confirm := func(deviceKey string, attrName string) (interface{}, error) {
		params := map[string]interface{}{}
		if attrName != "" {
			params["DeviceAttributes"] = []interface{}{
				map[string]interface{}{"Name": attrName, "Value": "UnitTest"},
			}
		}
		return env.svc.confirmDeviceCore(env.reqCtx, ConfirmDeviceInput{
			AccessToken: accessToken,
			DeviceKey:   deviceKey,
			Params:      params,
		})
	}

	resp, cerr := confirm("us-east-1_device-prompt", "given_name")
	if cerr != nil {
		t.Fatal(cerr)
	}
	if resp.(map[string]interface{})["UserConfirmationNecessary"] != true {
		t.Fatalf("UserConfirmationNecessary = %v, want true under DeviceOnlyRememberedOnUserPrompt", resp)
	}
	device, gerr := env.store.GetDevice(env.pool.ID, env.user.ID, "us-east-1_device-prompt")
	if gerr != nil {
		t.Fatal(gerr)
	}
	if device.DeviceRememberedStatus != "not_remembered" {
		t.Fatalf("device created %q under DeviceOnlyRememberedOnUserPrompt, want not_remembered", device.DeviceRememberedStatus)
	}
	if device.DeviceAttributes["given_name"] != "UnitTest" {
		t.Fatalf("device attributes %v, want the confirmed given_name", device.DeviceAttributes)
	}

	// The user confirms remembering through UpdateDeviceStatus.
	if _, uerr := env.svc.updateDeviceStatusCore(env.reqCtx, UpdateDeviceStatusInput{
		AccessToken:            accessToken,
		DeviceKey:              "us-east-1_device-prompt",
		DeviceRememberedStatus: "remembered",
	}); uerr != nil {
		t.Fatal(uerr)
	}
	device, gerr = env.store.GetDevice(env.pool.ID, env.user.ID, "us-east-1_device-prompt")
	if gerr != nil {
		t.Fatal(gerr)
	}
	if device.DeviceRememberedStatus != "remembered" {
		t.Fatalf("device status %q after UpdateDeviceStatus, want remembered", device.DeviceRememberedStatus)
	}

	// A duplicate device key is rejected, not silently overwritten.
	if _, cerr = confirm("us-east-1_device-prompt", ""); cerr == nil {
		t.Fatal("ConfirmDevice accepted an already-registered device key")
	}

	// An attribute outside the pool's schema is rejected.
	if _, cerr = confirm("us-east-1_device-schema", "custom:os"); cerr == nil {
		t.Fatal("ConfirmDevice accepted an attribute outside the pool schema")
	}

	// Without the prompt flag the device is immediately remembered and no
	// confirmation is asked for.
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.DeviceConfiguration.DeviceOnlyRememberedOnUserPrompt = false
	})
	resp, cerr = confirm("us-east-1_device-immediate", "")
	if cerr != nil {
		t.Fatal(cerr)
	}
	if resp.(map[string]interface{})["UserConfirmationNecessary"] != false {
		t.Fatalf("UserConfirmationNecessary = %v, want false without DeviceOnlyRememberedOnUserPrompt", resp)
	}
	device, gerr = env.store.GetDevice(env.pool.ID, env.user.ID, "us-east-1_device-immediate")
	if gerr != nil {
		t.Fatal(gerr)
	}
	if device.DeviceRememberedStatus != "remembered" {
		t.Fatalf("device created %q without the prompt flag, want remembered", device.DeviceRememberedStatus)
	}
}

// The registration ceremonies demand the pool's passkey enablement: a pool
// without a WebAuthnConfiguration answers WebAuthnNotEnabledException on
// both Start and Complete, a configuration without a relying party id on a
// pool without a hosted domain answers
// WebAuthnConfigurationMissingException, and a hosted domain alone satisfies
// the relying party id.
func TestWebAuthnRegistrationRequiresPoolEnablement(t *testing.T) {
	env := newChallengeTestEnv(t)
	accessToken, _, _, _, err := env.svc.CreateTokens(env.reqCtx, env.pool.ID, env.user.ID, challengeTestClientID, TokenGenerationAuthentication, nil)
	if err != nil {
		t.Fatal(err)
	}

	// An unconfigured pool: the passkey feature is not enabled.
	if _, serr := env.svc.startWebAuthnRegistrationCore(env.reqCtx, StartWebAuthnRegistrationInput{AccessToken: accessToken}); !strings.Contains(serr.Error(), "WebAuthnNotEnabledException") {
		t.Fatalf("start on an unconfigured pool returned %v, want WebAuthnNotEnabledException", serr)
	}

	// A configuration without a relying party id and no hosted domain.
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = &cognitostore.WebAuthnConfiguration{}
	})
	if _, serr := env.svc.startWebAuthnRegistrationCore(env.reqCtx, StartWebAuthnRegistrationInput{AccessToken: accessToken}); !strings.Contains(serr.Error(), "WebAuthnConfigurationMissingException") {
		t.Fatalf("start without a relying party id returned %v, want WebAuthnConfigurationMissingException", serr)
	}

	// A hosted domain satisfies the relying party id.
	const hostedDomain = "auth.example.com"
	if derr := env.store.SetUserPoolDomain(hostedDomain, &cognitostore.UserPoolDomain{
		Domain:     hostedDomain,
		UserPoolID: env.pool.ID,
	}); derr != nil {
		t.Fatal(derr)
	}
	started, serr := env.svc.startWebAuthnRegistrationCore(env.reqCtx, StartWebAuthnRegistrationInput{AccessToken: accessToken})
	if serr != nil {
		t.Fatalf("start with a hosted domain: %v", serr)
	}
	opts := started.(map[string]interface{})["CredentialCreationOptions"].(map[string]interface{})
	if id, _ := opts["rp"].(map[string]interface{})["id"].(string); id != hostedDomain {
		t.Fatalf("creation options rp.id %v, want the hosted domain %q", opts["rp"], hostedDomain)
	}

	// Removing the configuration again refuses the completion of the
	// pending registration: nothing registers on a pool without the
	// passkey feature.
	env.updatePool(t, func(p *cognitostore.UserPool) {
		p.WebAuthnConfiguration = nil
	})
	if _, cerr := env.svc.completeWebAuthnRegistrationCore(env.reqCtx, CompleteWebAuthnRegistrationInput{
		AccessToken: accessToken,
		Params:      map[string]interface{}{"Credential": "not-a-credential"},
	}); !strings.Contains(cerr.Error(), "WebAuthnNotEnabledException") {
		t.Fatalf("complete on an unconfigured pool returned %v, want WebAuthnNotEnabledException", cerr)
	}
}
