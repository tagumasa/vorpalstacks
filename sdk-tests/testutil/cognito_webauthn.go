package testutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitodoc "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/document"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// webauthnCOSEKey builds a minimal CBOR COSE EC2 public key whose COSE
// algorithm identifier (map key 3) is the given negative COSE encoding
// byte (-7 is 0x26, -8 is 0x27).
func webauthnCOSEKey(algByte byte) string {
	key := []byte{0xA5, 0x01, 0x02, 0x03, algByte, 0x20, 0x01, 0x21, 0x58, 0x20}
	key = append(key, make([]byte, 32)...)
	key = append(key, 0x22, 0x58, 0x20)
	key = append(key, make([]byte, 32)...)
	return base64.RawURLEncoding.EncodeToString(key)
}

// webauthnAttestation builds a minimal "none"-format attestation object (a
// three-entry CBOR map) whose authenticator data leads with the SHA-256 of
// the given relying party id, asserts user presence with attested credential
// data present, and embeds the credential ID and COSE public key that the
// completion cross-checks its client-asserted copies against.
func webauthnAttestation(rpID, credentialID, publicKeyB64 string) (string, error) {
	cose, err := base64.RawURLEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(rpID))
	idRaw := []byte(credentialID)
	var idLen [2]byte
	binary.BigEndian.PutUint16(idLen[:], uint16(len(idRaw)))
	authData := make([]byte, 0, 55+len(idRaw)+len(cose))
	authData = append(authData, sum[:]...)
	authData = append(authData, 0x41) // flags: user present, attested credential data present
	authData = append(authData, 0, 0, 0, 0)
	authData = append(authData, make([]byte, 16)...) // AAGUID
	authData = append(authData, idLen[:]...)
	authData = append(authData, idRaw...)
	authData = append(authData, cose...)

	att := []byte{
		0xA3,
		0x63, 'f', 'm', 't', 0x64, 'n', 'o', 'n', 'e',
		0x67, 'a', 't', 't', 'S', 't', 'm', 't', 0xA0,
		0x68, 'a', 'u', 't', 'h', 'D', 'a', 't', 'a',
	}
	if l := len(authData); l < 256 {
		att = append(att, 0x58, byte(l))
	} else {
		att = append(att, 0x59, byte(l>>8), byte(l))
	}
	att = append(att, authData...)
	return base64.RawURLEncoding.EncodeToString(att), nil
}

func webauthnClientData(challenge, origin string) (string, error) {
	clientData, err := json.Marshal(map[string]string{
		"type":      "webauthn.create",
		"challenge": challenge,
		"origin":    origin,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(clientData), nil
}

func webauthnCredential(attestationRpID, id, publicKey, clientDataB64 string) (cognitodoc.Interface, error) {
	attestation, err := webauthnAttestation(attestationRpID, id, publicKey)
	if err != nil {
		return nil, err
	}
	return cognitodoc.NewLazyDocument(map[string]interface{}{
		"id":        id,
		"publicKey": publicKey,
		"type":      "public-key",
		"response": map[string]interface{}{
			"clientDataJSON":    clientDataB64,
			"attestationObject": attestation,
		},
	}), nil
}

// cognitoWebAuthnTests covers the passkey registration flow: Start returns
// only CredentialCreationOptions, and Complete binds the pending challenge
// to the signed-in user and the starting app client server-side — no session
// member travels on the wire, and the completion is checked against the
// challenge, the registration origin, the attestation's relying-party hash
// and the credential's key algorithm.
func (r *TestRunner) cognitoWebAuthnTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "WebAuthn_RegistrationRoundTrip", func() error {
		// The ceremony runs on its own pool with the passkey feature
		// enabled — Start/CompleteWebAuthnRegistration refuse pools without
		// a WebAuthnConfiguration with WebAuthnNotEnabledException.
		waPoolID, cleanupWaPool, err := tc.createUserPool(tc.unique("webauthn-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupWaPool()
		if _, err := tc.client.SetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
			UserPoolId:            aws.String(waPoolID),
			WebAuthnConfiguration: &types.WebAuthnConfigurationType{RelyingPartyId: aws.String("auth.webauthn.example.com")},
		}); err != nil {
			return fmt.Errorf("enable passkeys: %v", err)
		}

		waUser := tc.unique("webauthn-user")
		if _, err := tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:    aws.String(waPoolID),
			Username:      aws.String(waUser),
			MessageAction: types.MessageActionTypeSuppress,
		}); err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer func() {
			_, _ = tc.client.AdminDeleteUser(tc.ctx, &cognitoidentityprovider.AdminDeleteUserInput{
				UserPoolId: aws.String(waPoolID),
				Username:   aws.String(waUser),
			})
		}()
		if _, err := tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(waPoolID),
			Username:   aws.String(waUser),
			Password:   aws.String("WebauthnPass123!"),
			Permanent:  true,
		}); err != nil {
			return fmt.Errorf("set permanent password: %v", err)
		}
		auth := func(clientID string) (string, error) {
			resp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
				UserPoolId: aws.String(waPoolID),
				ClientId:   aws.String(clientID),
				AuthFlow:   "ADMIN_NO_SRP_AUTH",
				AuthParameters: map[string]string{
					"USERNAME": waUser,
					"PASSWORD": "WebauthnPass123!",
				},
			})
			if err != nil {
				return "", err
			}
			return aws.ToString(resp.AuthenticationResult.AccessToken), nil
		}

		clientA, cleanupA, err := tc.createPoolClient(waPoolID, tc.unique("webauthn-client-a"))
		if err != nil {
			return fmt.Errorf("create client A: %v", err)
		}
		defer cleanupA()
		tokenA, err := auth(clientA)
		if err != nil {
			return fmt.Errorf("auth client A: %v", err)
		}

		// Completing without a pending Start registration is rejected with
		// the documented challenge-not-found exception.
		cred, err := webauthnCredential("", "no-start-cred", webauthnCOSEKey(0x26), "")
		if err != nil {
			return err
		}
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  cred,
		})
		if err := expectAWSErrorCode(err, "WebAuthnChallengeNotFoundException"); err != nil {
			return fmt.Errorf("complete without start: %w", err)
		}

		startOut, err := tc.client.StartWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.StartWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
		})
		if err != nil {
			return fmt.Errorf("StartWebAuthnRegistration: %v", err)
		}
		// The SDK output shape has no Session member — the model response
		// carries only CredentialCreationOptions.
		var options struct {
			Challenge string `json:"challenge"`
			Rp        struct {
				ID string `json:"id"`
			} `json:"rp"`
		}
		if unmarshaller, ok := startOut.CredentialCreationOptions.(interface {
			UnmarshalSmithyDocument(interface{}) error
		}); !ok {
			return fmt.Errorf("CredentialCreationOptions is not unmarshallable")
		} else if err := unmarshaller.UnmarshalSmithyDocument(&options); err != nil {
			return fmt.Errorf("unmarshal CredentialCreationOptions: %v", err)
		}
		if options.Challenge == "" || options.Rp.ID == "" {
			return fmt.Errorf("CredentialCreationOptions challenge or rp.id is empty")
		}
		validClientData, err := webauthnClientData(options.Challenge, "https://"+options.Rp.ID)
		if err != nil {
			return err
		}

		// A token issued to a different app client cannot complete the
		// registration started through client A.
		clientB, cleanupB, err := tc.createPoolClient(waPoolID, tc.unique("webauthn-client-b"))
		if err != nil {
			return fmt.Errorf("create client B: %v", err)
		}
		defer cleanupB()
		tokenB, err := auth(clientB)
		if err != nil {
			return fmt.Errorf("auth client B: %v", err)
		}
		cred, err = webauthnCredential(options.Rp.ID, "test-credential-1", webauthnCOSEKey(0x26), validClientData)
		if err != nil {
			return err
		}
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenB),
			Credential:  cred,
		})
		if err := expectAWSErrorCode(err, "WebAuthnClientMismatchException"); err != nil {
			return fmt.Errorf("complete with client B token: %w", err)
		}

		// An origin that does not align with the relying party id is
		// rejected.
		evilOriginData, err := webauthnClientData(options.Challenge, "https://evil.example.com")
		if err != nil {
			return err
		}
		cred, err = webauthnCredential(options.Rp.ID, "test-credential-1", webauthnCOSEKey(0x26), evilOriginData)
		if err != nil {
			return err
		}
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  cred,
		})
		if err := expectAWSErrorCode(err, "WebAuthnOriginNotAllowedException"); err != nil {
			return fmt.Errorf("complete with foreign origin: %w", err)
		}

		// An attestation whose authenticator data hashes a different
		// relying party id is rejected.
		cred, err = webauthnCredential("other-rp.example.com", "test-credential-1", webauthnCOSEKey(0x26), validClientData)
		if err != nil {
			return err
		}
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  cred,
		})
		if err := expectAWSErrorCode(err, "WebAuthnRelyingPartyMismatchException"); err != nil {
			return fmt.Errorf("complete with foreign rpIdHash: %w", err)
		}

		// A credential whose COSE key algorithm was not offered in the
		// creation options (-8 EdDSA; the pool offers ES256 and RS256) is
		// rejected.
		cred, err = webauthnCredential(options.Rp.ID, "test-credential-1", webauthnCOSEKey(0x27), validClientData)
		if err != nil {
			return err
		}
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  cred,
		})
		if err := expectAWSErrorCode(err, "WebAuthnCredentialNotSupportedException"); err != nil {
			return fmt.Errorf("complete with unsupported alg: %w", err)
		}

		cred, err = webauthnCredential(options.Rp.ID, "test-credential-1", webauthnCOSEKey(0x26), validClientData)
		if err != nil {
			return err
		}
		if _, err := tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  cred,
		}); err != nil {
			return fmt.Errorf("CompleteWebAuthnRegistration: %v", err)
		}

		// The completed challenge is single-use: a second completion with
		// the same data no longer finds a pending registration.
		cred, err = webauthnCredential(options.Rp.ID, "test-credential-1", webauthnCOSEKey(0x26), validClientData)
		if err != nil {
			return err
		}
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  cred,
		})
		if err := expectAWSErrorCode(err, "WebAuthnChallengeNotFoundException"); err != nil {
			return fmt.Errorf("second completion: %w", err)
		}

		// Deleting the registered credential succeeds; deleting it again,
		// or any unregistered id, is a resource-not-found error.
		if _, err := tc.client.DeleteWebAuthnCredential(tc.ctx, &cognitoidentityprovider.DeleteWebAuthnCredentialInput{
			AccessToken:  aws.String(tokenA),
			CredentialId: aws.String("test-credential-1"),
		}); err != nil {
			return fmt.Errorf("delete registered credential: %v", err)
		}
		_, err = tc.client.DeleteWebAuthnCredential(tc.ctx, &cognitoidentityprovider.DeleteWebAuthnCredentialInput{
			AccessToken:  aws.String(tokenA),
			CredentialId: aws.String("test-credential-1"),
		})
		return expectAWSErrorCode(err, "ResourceNotFoundException")
	}))

	// The WEB_AUTHN sign-in round trip on a pool with passkeys configured:
	// InitiateAuth(USER_AUTH, PREFERRED_CHALLENGE=WEB_AUTHN) issues the
	// challenge with CREDENTIAL_REQUEST_OPTIONS, and the CREDENTIAL answer —
	// a W3C AuthenticationResponseJSON signed by the registered key —
	// completes the sign-in with tokens.
	results = append(results, r.RunTest("cognito", "WebAuthn_AssertionSignInRoundTrip", func() error {
		const rpID = "auth.example.com"
		poolResp, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(tc.unique("webauthn-signin-pool")),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		poolID := *poolResp.UserPool.Id
		defer func() {
			_, _ = tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: aws.String(poolID),
			})
		}()

		if _, err := tc.client.SetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
			UserPoolId:            aws.String(poolID),
			MfaConfiguration:      types.UserPoolMfaTypeOff,
			WebAuthnConfiguration: &types.WebAuthnConfigurationType{RelyingPartyId: aws.String(rpID)},
		}); err != nil {
			return fmt.Errorf("configure passkeys: %v", err)
		}

		clientID, cleanupClient, err := tc.createPoolClient(poolID, tc.unique("webauthn-signin-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupClient()

		waUser := tc.unique("webauthn-signin-user")
		const waPassword = "WebAuthnPass123!"
		if _, err := tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(poolID),
			Username:          aws.String(waUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		}); err != nil {
			return fmt.Errorf("admin create user: %v", err)
		}
		defer func() {
			_, _ = tc.client.AdminDeleteUser(tc.ctx, &cognitoidentityprovider.AdminDeleteUserInput{
				UserPoolId: aws.String(poolID),
				Username:   aws.String(waUser),
			})
		}()
		if _, err := tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(poolID),
			Username:   aws.String(waUser),
			Password:   aws.String(waPassword),
			Permanent:  true,
		}); err != nil {
			return fmt.Errorf("set permanent password: %v", err)
		}

		authResp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(poolID),
			ClientId:   aws.String(clientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": waUser,
				"PASSWORD": waPassword,
			},
		})
		if err != nil {
			return fmt.Errorf("admin initiate auth: %v", err)
		}
		accessToken := aws.ToString(authResp.AuthenticationResult.AccessToken)

		startOut, err := tc.client.StartWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.StartWebAuthnRegistrationInput{
			AccessToken: aws.String(accessToken),
		})
		if err != nil {
			return fmt.Errorf("StartWebAuthnRegistration: %v", err)
		}
		var options struct {
			Challenge string `json:"challenge"`
			Rp        struct {
				ID string `json:"id"`
			} `json:"rp"`
		}
		if unmarshaller, ok := startOut.CredentialCreationOptions.(interface {
			UnmarshalSmithyDocument(interface{}) error
		}); !ok {
			return fmt.Errorf("CredentialCreationOptions is not unmarshallable")
		} else if err := unmarshaller.UnmarshalSmithyDocument(&options); err != nil {
			return fmt.Errorf("unmarshal CredentialCreationOptions: %v", err)
		}

		priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		const credentialID = "signin-passkey-1"
		clientData, err := webauthnClientData(options.Challenge, "https://"+rpID)
		if err != nil {
			return err
		}
		cred, err := webauthnCredential(rpID, credentialID, webauthnCOSEKeyFromECDSA(priv), clientData)
		if err != nil {
			return err
		}
		if _, err := tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(accessToken),
			Credential:  cred,
		}); err != nil {
			return fmt.Errorf("CompleteWebAuthnRegistration: %v", err)
		}

		initResp, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeUserAuth,
			ClientId: aws.String(clientID),
			AuthParameters: map[string]string{
				"USERNAME":            waUser,
				"PREFERRED_CHALLENGE": "WEB_AUTHN",
			},
		})
		if err != nil {
			return fmt.Errorf("InitiateAuth USER_AUTH WEB_AUTHN: %v", err)
		}
		if string(initResp.ChallengeName) != "WEB_AUTHN" {
			return fmt.Errorf("expected ChallengeName=WEB_AUTHN, got %v", initResp.ChallengeName)
		}
		var requestOptions struct {
			Challenge        string           `json:"challenge"`
			RpID             string           `json:"rpId"`
			AllowCredentials []map[string]any `json:"allowCredentials"`
		}
		if err := json.Unmarshal([]byte(initResp.ChallengeParameters["CREDENTIAL_REQUEST_OPTIONS"]), &requestOptions); err != nil {
			return fmt.Errorf("parse CREDENTIAL_REQUEST_OPTIONS: %v", err)
		}
		if requestOptions.RpID != rpID || len(requestOptions.AllowCredentials) != 1 {
			return fmt.Errorf("unexpected request options: %+v", requestOptions)
		}
		if id, _ := requestOptions.AllowCredentials[0]["id"].(string); id != credentialID {
			return fmt.Errorf("allowCredentials id %q, want %q", id, credentialID)
		}

		assertion := webauthnAssertion(priv, credentialID, rpID, requestOptions.Challenge, "https://"+rpID, 1)
		challengeResp, err := tc.client.RespondToAuthChallenge(tc.ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName: "WEB_AUTHN",
			ClientId:      aws.String(clientID),
			Session:       initResp.Session,
			ChallengeResponses: map[string]string{
				"USERNAME":   waUser,
				"CREDENTIAL": assertion,
			},
		})
		if err != nil {
			return fmt.Errorf("RespondToAuthChallenge WEB_AUTHN: %v", err)
		}
		if challengeResp.AuthenticationResult == nil || aws.ToString(challengeResp.AuthenticationResult.AccessToken) == "" {
			return fmt.Errorf("no access token after the passkey assertion: %+v", challengeResp.AuthenticationResult)
		}
		return nil
	}))

	return results
}

// webauthnCOSEKeyFromECDSA builds the COSE EC2 public key of a P-256 key
// pair: the same five-entry CBOR map as webauthnCOSEKey, with the real
// curve coordinates.
func webauthnCOSEKeyFromECDSA(priv *ecdsa.PrivateKey) string {
	x := make([]byte, 32)
	y := make([]byte, 32)
	priv.PublicKey.X.FillBytes(x)
	priv.PublicKey.Y.FillBytes(y)
	key := []byte{0xA5, 0x01, 0x02, 0x03, 0x26, 0x20, 0x01, 0x21, 0x58, 0x20}
	key = append(key, x...)
	key = append(key, 0x22, 0x58, 0x20)
	key = append(key, y...)
	return base64.RawURLEncoding.EncodeToString(key)
}

// webauthnAssertion builds the AuthenticationResponseJSON the client
// submits as the CREDENTIAL answer of a WEB_AUTHN challenge: the assertion
// signs authenticatorData || SHA-256(clientDataJSON) with the credential's
// private key, in the COSE raw fixed-width r||s signature encoding that is
// the WebAuthn wire format.
func webauthnAssertion(priv *ecdsa.PrivateKey, credentialID, rpID, challenge, origin string, signCount uint32) string {
	rpIDHash := sha256.Sum256([]byte(rpID))
	authData := make([]byte, 37)
	copy(authData, rpIDHash[:])
	authData[32] = 0x01 // user present
	binary.BigEndian.PutUint32(authData[33:37], signCount)

	clientData, _ := json.Marshal(map[string]string{
		"type":      "webauthn.get",
		"challenge": challenge,
		"origin":    origin,
	})
	clientDataHash := sha256.Sum256(clientData)
	message := append(append([]byte{}, authData...), clientDataHash[:]...)
	digest := sha256.Sum256(message)
	r, s, _ := ecdsa.Sign(rand.Reader, priv, digest[:])
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])

	assertion, _ := json.Marshal(map[string]interface{}{
		"id":    credentialID,
		"rawId": credentialID,
		"type":  "public-key",
		"response": map[string]string{
			"authenticatorData": base64.RawURLEncoding.EncodeToString(authData),
			"clientDataJSON":    base64.RawURLEncoding.EncodeToString(clientData),
			"signature":         base64.RawURLEncoding.EncodeToString(sig),
		},
	})
	return string(assertion)
}
