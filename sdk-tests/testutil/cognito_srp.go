package testutil

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// This file exercises the USER_SRP_AUTH / PASSWORD_VERIFIER flow end-to-end
// against the VorpalStacks Cognito IDP implementation. The client-side SRP
// maths are implemented inline (rather than depending on alexrudd/cognito-srp,
// whose published pseudo-version does not compile against the current
// aws-sdk-go-v2 layout). The algorithm is a direct port of
// amazon-cognito-identity-js/src/AuthenticationHelper.js and has been
// verified to interoperate with the server-side implementation in
// internal/services/aws/cognitoidentityprovider/srp.go.

const (
	cognitoSrpNHex = "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
		"15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
		"ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
		"ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
		"F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
		"BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
		"43DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF"
	cognitoSrpGHex     = "2"
	cognitoSrpInfoBits = "Caldera Derived Key"
)

var (
	cognitoSrpN = cognitoSrpMustHex(cognitoSrpNHex)
	cognitoSrpG = cognitoSrpMustHex(cognitoSrpGHex)
	cognitoSrpK = cognitoSrpMustHex(cognitoSrpHexHash("00" + cognitoSrpNHex + "0" + cognitoSrpGHex))
)

func cognitoSrpMustHex(s string) *big.Int { n, _ := new(big.Int).SetString(s, 16); return n }
func cognitoSrpHashSHA256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
func cognitoSrpHexHash(s string) string {
	b, _ := hex.DecodeString(s)
	return cognitoSrpHashSHA256Hex(b)
}

func cognitoSrpPadHex(s string) string {
	if len(s)%2 == 1 {
		return "0" + s
	}
	if len(s) > 0 && strings.IndexByte("89abcdef", s[0]) >= 0 {
		return "00" + s
	}
	return s
}

// cognitoSrpClient holds the ephemeral client-side SRP state for a single
// authentication attempt. It must not be reused across challenges.
type cognitoSrpClient struct {
	poolName string
	username string
	password string
	a        *big.Int
	bigA     *big.Int
}

func newCognitoSrpClient(poolID, username, password string) (*cognitoSrpClient, error) {
	idx := strings.Index(poolID, "_")
	if idx < 0 || idx == len(poolID)-1 {
		return nil, fmt.Errorf("invalid pool ID %q", poolID)
	}
	c := &cognitoSrpClient{
		poolName: poolID[idx+1:],
		username: username,
		password: password,
	}
	rb := make([]byte, 128)
	if _, err := rand.Read(rb); err != nil {
		return nil, err
	}
	c.a = new(big.Int).Mod(new(big.Int).SetBytes(rb), cognitoSrpN)
	c.bigA = new(big.Int).Exp(cognitoSrpG, c.a, cognitoSrpN)
	if new(big.Int).Mod(c.bigA, cognitoSrpN).Sign() == 0 {
		return nil, fmt.Errorf("client safety check failed: A mod N == 0")
	}
	return c, nil
}

func (c *cognitoSrpClient) srpAHex() string { return c.bigA.Text(16) }

// passwordVerifier computes the client's PASSWORD_CLAIM_SIGNATURE for the
// challenge parameters returned by InitiateAuth. It mirrors the reference
// amazon-cognito-identity-js implementation.
func (c *cognitoSrpClient) passwordVerifier(saltHex, srpBHex, secretBlockB64 string, now time.Time) (signature, timestamp string, err error) {
	bigB := cognitoSrpMustHex(srpBHex)
	saltInt, ok := new(big.Int).SetString(saltHex, 16)
	if !ok {
		return "", "", fmt.Errorf("invalid SALT hex")
	}
	secretBlock, err := base64.StdEncoding.DecodeString(secretBlockB64)
	if err != nil {
		return "", "", fmt.Errorf("invalid SECRET_BLOCK: %w", err)
	}

	userPass := c.poolName + c.username + ":" + c.password
	userPassHash := cognitoSrpHashSHA256Hex([]byte(userPass))
	uVal := cognitoSrpMustHex(cognitoSrpHexHash(cognitoSrpPadHex(c.bigA.Text(16)) + cognitoSrpPadHex(bigB.Text(16))))
	xVal := cognitoSrpMustHex(cognitoSrpHexHash(cognitoSrpPadHex(saltInt.Text(16)) + userPassHash))
	gModPowXN := new(big.Int).Exp(cognitoSrpG, xVal, cognitoSrpN)
	intVal1 := new(big.Int).Sub(bigB, new(big.Int).Mul(cognitoSrpK, gModPowXN))
	intVal2 := new(big.Int).Add(c.a, new(big.Int).Mul(uVal, xVal))
	sVal := new(big.Int).Exp(intVal1, intVal2, cognitoSrpN)

	hkdfKey := cognitoSrpHKDF(cognitoSrpPadHex(sVal.Text(16)), cognitoSrpPadHex(uVal.Text(16)))

	timestamp = now.In(time.UTC).Format("Mon Jan 2 03:04:05 MST 2006")
	msg := c.poolName + c.username + string(secretBlock) + timestamp
	mac := hmac.New(sha256.New, hkdfKey)
	mac.Write([]byte(msg))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature, timestamp, nil
}

func cognitoSrpHKDF(ikmHex, saltHex string) []byte {
	ikm, _ := hex.DecodeString(ikmHex)
	salt, _ := hex.DecodeString(saltHex)
	mac := hmac.New(sha256.New, salt)
	mac.Write(ikm)
	prk := mac.Sum(nil)
	mac2 := hmac.New(sha256.New, prk)
	mac2.Write(append([]byte(cognitoSrpInfoBits), 1))
	return mac2.Sum(nil)[:16]
}

// cognitoSRPTests exercises the USER_SRP_AUTH / PASSWORD_VERIFIER flow.
func (r *TestRunner) cognitoSRPTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	// Happy path: create a confirmed user with a known permanent password,
	// then complete the full SRP exchange.
	results = append(results, r.RunTest("cognito", "InitiateAuth_USER_SRP_AUTH", func() error {
		srpClientID, cleanupSrpClientID, err := tc.createPoolClient(tc.userPoolID, tc.unique("srp-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupSrpClientID()

		srpUser := tc.unique("srp-user")
		srpPassword := "SrpPassword123!"
		cleanupSrpUser, err := tc.adminCreateUser(srpUser,
			func(input *cognitoidentityprovider.AdminCreateUserInput) {
				input.MessageAction = ""
			})
		if err != nil {
			return fmt.Errorf("admin create user: %v", err)
		}
		defer cleanupSrpUser()
		_, err = tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(srpUser),
			Password:   aws.String(srpPassword),
			Permanent:  true,
		})
		if err != nil {
			return fmt.Errorf("admin set user password: %v", err)
		}

		c, err := newCognitoSrpClient(tc.userPoolID, srpUser, srpPassword)
		if err != nil {
			return err
		}
		initResp, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "USER_SRP_AUTH",
			ClientId: aws.String(srpClientID),
			AuthParameters: map[string]string{
				"USERNAME": srpUser,
				"SRP_A":    c.srpAHex(),
			},
		})
		if err != nil {
			return fmt.Errorf("InitiateAuth USER_SRP_AUTH: %v", err)
		}
		if initResp.ChallengeName != "PASSWORD_VERIFIER" {
			return fmt.Errorf("expected ChallengeName=PASSWORD_VERIFIER, got %v", initResp.ChallengeName)
		}
		cp := initResp.ChallengeParameters
		if cp["SALT"] == "" || cp["SECRET_BLOCK"] == "" || cp["SRP_B"] == "" || cp["USER_ID_FOR_SRP"] == "" {
			return fmt.Errorf("missing challenge parameters: %+v", cp)
		}

		sig, ts, err := c.passwordVerifier(cp["SALT"], cp["SRP_B"], cp["SECRET_BLOCK"], time.Now())
		if err != nil {
			return fmt.Errorf("compute password verifier: %v", err)
		}
		respResp, err := tc.client.RespondToAuthChallenge(tc.ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName: "PASSWORD_VERIFIER",
			ClientId:      aws.String(srpClientID),
			Session:       initResp.Session,
			ChallengeResponses: map[string]string{
				"USERNAME":                    srpUser,
				"PASSWORD_CLAIM_SIGNATURE":    sig,
				"PASSWORD_CLAIM_SECRET_BLOCK": cp["SECRET_BLOCK"],
				"TIMESTAMP":                   ts,
			},
		})
		if err != nil {
			return fmt.Errorf("RespondToAuthChallenge PASSWORD_VERIFIER: %v", err)
		}
		if respResp.AuthenticationResult == nil || respResp.AuthenticationResult.AccessToken == nil {
			return fmt.Errorf("no access token in SRP auth result: %+v", respResp.AuthenticationResult)
		}
		return nil
	}))

	// Negative test: a client using the wrong password must be rejected.
	results = append(results, r.RunTest("cognito", "USER_SRP_AUTH_wrong_password_rejected", func() error {
		srpClientID, cleanupSrpClientID, err := tc.createPoolClient(tc.userPoolID, tc.unique("srp-neg"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupSrpClientID()

		srpUser := tc.unique("srp-wrong")
		cleanupSrpUser, err := tc.adminCreateUser(srpUser,
			func(input *cognitoidentityprovider.AdminCreateUserInput) {
				input.MessageAction = ""
			})
		if err != nil {
			return fmt.Errorf("admin create user: %v", err)
		}
		defer cleanupSrpUser()
		_, err = tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(srpUser),
			Password:   aws.String("CorrectPassword123!"),
			Permanent:  true,
		})
		if err != nil {
			return fmt.Errorf("admin set user password: %v", err)
		}

		// Client uses the WRONG password.
		c, err := newCognitoSrpClient(tc.userPoolID, srpUser, "WrongPassword456!")
		if err != nil {
			return err
		}
		initResp, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "USER_SRP_AUTH",
			ClientId: aws.String(srpClientID),
			AuthParameters: map[string]string{
				"USERNAME": srpUser,
				"SRP_A":    c.srpAHex(),
			},
		})
		if err != nil {
			return fmt.Errorf("InitiateAuth: %v", err)
		}
		cp := initResp.ChallengeParameters
		sig, ts, err := c.passwordVerifier(cp["SALT"], cp["SRP_B"], cp["SECRET_BLOCK"], time.Now())
		if err != nil {
			return err
		}
		_, err = tc.client.RespondToAuthChallenge(tc.ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName: "PASSWORD_VERIFIER",
			ClientId:      aws.String(srpClientID),
			Session:       initResp.Session,
			ChallengeResponses: map[string]string{
				"USERNAME":                    srpUser,
				"PASSWORD_CLAIM_SIGNATURE":    sig,
				"PASSWORD_CLAIM_SECRET_BLOCK": cp["SECRET_BLOCK"],
				"TIMESTAMP":                   ts,
			},
		})
		if err == nil {
			return fmt.Errorf("expected NotAuthorizedException for wrong password, got success")
		}
		if !strings.Contains(err.Error(), "NotAuthorized") {
			return fmt.Errorf("expected NotAuthorizedException, got: %v", err)
		}
		return nil
	}))

	// Choice-based sign-in: USER_AUTH issues a SELECT_CHALLENGE session, the
	// client selects PASSWORD via ANSWER inside ChallengeResponses (the
	// documented carrier), and completing the PASSWORD challenge must return
	// the full AuthenticationResult including the refresh token.
	results = append(results, r.RunTest("cognito", "RespondToAuthChallenge_SelectChallengePasswordFlow", func() error {
		selectClientID, cleanupSelectClient, err := tc.createPoolClient(tc.userPoolID, tc.unique("select-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupSelectClient()

		selectUser := tc.unique("select-user")
		selectPassword := "SelectPassword123!"
		cleanupSelectUser, err := tc.adminCreateUser(selectUser,
			func(input *cognitoidentityprovider.AdminCreateUserInput) {
				input.MessageAction = ""
			})
		if err != nil {
			return fmt.Errorf("admin create user: %v", err)
		}
		defer cleanupSelectUser()
		_, err = tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(selectUser),
			Password:   aws.String(selectPassword),
			Permanent:  true,
		})
		if err != nil {
			return fmt.Errorf("admin set user password: %v", err)
		}

		// USER_AUTH takes USERNAME inside AuthParameters per the AWS API
		// reference.
		initResp, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "USER_AUTH",
			ClientId: aws.String(selectClientID),
			AuthParameters: map[string]string{
				"USERNAME": selectUser,
			},
		})
		if err != nil {
			return fmt.Errorf("InitiateAuth USER_AUTH: %v", err)
		}
		// Without PREFERRED_CHALLENGE the documented USER_AUTH response
		// carries AvailableChallenges and a session (no ChallengeName).
		if initResp.Session == nil || *initResp.Session == "" {
			return fmt.Errorf("expected a session from InitiateAuth USER_AUTH")
		}
		hasPassword := false
		for _, c := range initResp.AvailableChallenges {
			if c == "PASSWORD" {
				hasPassword = true
			}
		}
		if !hasPassword {
			return fmt.Errorf("expected PASSWORD among available challenges, got %v", initResp.AvailableChallenges)
		}

		selectResp, err := tc.client.RespondToAuthChallenge(tc.ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName: "SELECT_CHALLENGE",
			ClientId:      aws.String(selectClientID),
			Session:       initResp.Session,
			ChallengeResponses: map[string]string{
				"USERNAME": selectUser,
				"ANSWER":   "PASSWORD",
			},
		})
		if err != nil {
			return fmt.Errorf("RespondToAuthChallenge SELECT_CHALLENGE: %v", err)
		}
		if selectResp.ChallengeName != "PASSWORD" {
			return fmt.Errorf("expected PASSWORD challenge after selection, got %q", selectResp.ChallengeName)
		}

		authResp, err := tc.client.RespondToAuthChallenge(tc.ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName: "PASSWORD",
			ClientId:      aws.String(selectClientID),
			Session:       selectResp.Session,
			ChallengeResponses: map[string]string{
				"USERNAME": selectUser,
				"PASSWORD": selectPassword,
			},
		})
		if err != nil {
			return fmt.Errorf("RespondToAuthChallenge PASSWORD: %v", err)
		}
		if authResp.AuthenticationResult == nil {
			return fmt.Errorf("expected AuthenticationResult after PASSWORD challenge")
		}
		if aws.ToString(authResp.AuthenticationResult.AccessToken) == "" || aws.ToString(authResp.AuthenticationResult.IdToken) == "" {
			return fmt.Errorf("expected access and ID tokens in the challenge sign-in result")
		}
		if aws.ToString(authResp.AuthenticationResult.RefreshToken) == "" {
			return fmt.Errorf("expected refresh token in the challenge sign-in result")
		}
		return nil
	}))

	return results
}
