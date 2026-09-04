package testutil

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitodoc "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/document"
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

// webauthnAttestation builds a minimal "none"-format attestation object
// (a three-entry CBOR map) whose authenticator data leads with the
// SHA-256 of the given relying party id.
func webauthnAttestation(rpID string) string {
	sum := sha256.Sum256([]byte(rpID))
	att := []byte{
		0xA3,
		0x63, 'f', 'm', 't', 0x64, 'n', 'o', 'n', 'e',
		0x67, 'a', 't', 't', 'S', 't', 'm', 't', 0xA0,
		0x68, 'a', 'u', 't', 'h', 'D', 'a', 't', 'a', 0x58, 0x25,
	}
	att = append(att, sum[:]...)
	// flags byte + 4-byte sign count complete the 37-byte auth data.
	att = append(att, 0x00, 0, 0, 0, 0)
	return base64.RawURLEncoding.EncodeToString(att)
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

func webauthnCredential(id, publicKey, clientDataB64, attestationB64 string) cognitodoc.Interface {
	return cognitodoc.NewLazyDocument(map[string]interface{}{
		"id":        id,
		"publicKey": publicKey,
		"type":      "public-key",
		"response": map[string]interface{}{
			"clientDataJSON":    clientDataB64,
			"attestationObject": attestationB64,
		},
	})
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
		waUser := tc.unique("webauthn-user")
		cleanupUser, err := tc.createConfirmedUser(waUser, "WebAuthnPass123!")
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupUser()
		auth := func(clientID string) (string, error) {
			resp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
				UserPoolId: aws.String(tc.userPoolID),
				ClientId:   aws.String(clientID),
				AuthFlow:   "ADMIN_NO_SRP_AUTH",
				AuthParameters: map[string]string{
					"USERNAME": waUser,
					"PASSWORD": "WebAuthnPass123!",
				},
			})
			if err != nil {
				return "", err
			}
			return aws.ToString(resp.AuthenticationResult.AccessToken), nil
		}

		clientA, cleanupA, err := tc.createPoolClient(tc.userPoolID, tc.unique("webauthn-client-a"))
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
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  webauthnCredential("no-start-cred", webauthnCOSEKey(0x26), "", ""),
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
		clientB, cleanupB, err := tc.createPoolClient(tc.userPoolID, tc.unique("webauthn-client-b"))
		if err != nil {
			return fmt.Errorf("create client B: %v", err)
		}
		defer cleanupB()
		tokenB, err := auth(clientB)
		if err != nil {
			return fmt.Errorf("auth client B: %v", err)
		}
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenB),
			Credential:  webauthnCredential("test-credential-1", webauthnCOSEKey(0x26), validClientData, webauthnAttestation(options.Rp.ID)),
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
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  webauthnCredential("test-credential-1", webauthnCOSEKey(0x26), evilOriginData, webauthnAttestation(options.Rp.ID)),
		})
		if err := expectAWSErrorCode(err, "WebAuthnOriginNotAllowedException"); err != nil {
			return fmt.Errorf("complete with foreign origin: %w", err)
		}

		// An attestation whose authenticator data hashes a different
		// relying party id is rejected.
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  webauthnCredential("test-credential-1", webauthnCOSEKey(0x26), validClientData, webauthnAttestation("other-rp.example.com")),
		})
		if err := expectAWSErrorCode(err, "WebAuthnRelyingPartyMismatchException"); err != nil {
			return fmt.Errorf("complete with foreign rpIdHash: %w", err)
		}

		// A credential whose COSE key algorithm was not offered in the
		// creation options (-8 EdDSA; the pool offers ES256 and RS256) is
		// rejected.
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  webauthnCredential("test-credential-1", webauthnCOSEKey(0x27), validClientData, webauthnAttestation(options.Rp.ID)),
		})
		if err := expectAWSErrorCode(err, "WebAuthnCredentialNotSupportedException"); err != nil {
			return fmt.Errorf("complete with unsupported alg: %w", err)
		}

		if _, err := tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  webauthnCredential("test-credential-1", webauthnCOSEKey(0x26), validClientData, webauthnAttestation(options.Rp.ID)),
		}); err != nil {
			return fmt.Errorf("CompleteWebAuthnRegistration: %v", err)
		}

		// The completed challenge is single-use: a second completion with
		// the same data no longer finds a pending registration.
		_, err = tc.client.CompleteWebAuthnRegistration(tc.ctx, &cognitoidentityprovider.CompleteWebAuthnRegistrationInput{
			AccessToken: aws.String(tokenA),
			Credential:  webauthnCredential("test-credential-1", webauthnCOSEKey(0x26), validClientData, webauthnAttestation(options.Rp.ID)),
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

	return results
}
