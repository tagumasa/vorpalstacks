package testutil

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// cognitoMFATests pins second-factor enforcement at sign-in: a pool with
// MFA ON never mints tokens on the primary credentials alone.
func (r *TestRunner) cognitoMFATests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "InitiateAuth_MFA_TOTPEnforced", func() error {
		mfaPoolID, cleanupPool, err := tc.createUserPool(tc.unique("mfa-pool"), func(input *cognitoidentityprovider.CreateUserPoolInput) {
			input.MfaConfiguration = types.UserPoolMfaTypeOn
		})
		if err != nil {
			return fmt.Errorf("create MFA pool: %v", err)
		}
		defer cleanupPool()
		mfaClientID, cleanupClient, err := tc.createPoolClient(mfaPoolID, tc.unique("mfa-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupClient()
		mfaUser := tc.unique("mfa-user")
		_, err = tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(mfaPoolID),
			Username:          aws.String(mfaUser),
			TemporaryPassword: aws.String("MfaTemp123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer tc.client.AdminDeleteUser(tc.ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(mfaPoolID),
			Username:   aws.String(mfaUser),
		})
		if _, err := tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(mfaPoolID),
			Username:   aws.String(mfaUser),
			Password:   aws.String("MfaPass123!"),
			Permanent:  true,
		}); err != nil {
			return fmt.Errorf("set permanent password: %v", err)
		}

		passwordAuth := func() (*cognitoidentityprovider.InitiateAuthOutput, error) {
			return tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
				ClientId: aws.String(mfaClientID),
				AuthFlow: types.AuthFlowTypeUserPasswordAuth,
				AuthParameters: map[string]string{
					"USERNAME": mfaUser,
					"PASSWORD": "MfaPass123!",
				},
			})
		}

		// MFA ON with no enrolled factor routes to MFA_SETUP.
		authResp, err := passwordAuth()
		if err != nil {
			return err
		}
		if authResp.ChallengeName != types.ChallengeNameTypeMfaSetup {
			return fmt.Errorf("expected MFA_SETUP for a factor-less user under MFA ON, got %v (AuthenticationResult %v)", authResp.ChallengeName, authResp.AuthenticationResult != nil)
		}
		setupSession := aws.ToString(authResp.Session)

		// Mid-sign-in enrolment: the MFA_SETUP session drives both
		// AssociateSoftwareToken and VerifySoftwareToken.
		assocResp, err := tc.client.AssociateSoftwareToken(tc.ctx, &cognitoidentityprovider.AssociateSoftwareTokenInput{
			Session: aws.String(setupSession),
		})
		if err != nil {
			return fmt.Errorf("AssociateSoftwareToken with the MFA_SETUP session: %v", err)
		}
		secretCode := aws.ToString(assocResp.SecretCode)
		if secretCode == "" {
			return fmt.Errorf("SecretCode is empty")
		}
		verifyResp, err := tc.client.VerifySoftwareToken(tc.ctx, &cognitoidentityprovider.VerifySoftwareTokenInput{
			Session:  aws.String(setupSession),
			UserCode: aws.String(totpCodeAt([]byte(secretCode), 0)),
		})
		if err != nil {
			return fmt.Errorf("VerifySoftwareToken with the MFA_SETUP session: %v", err)
		}
		if verifyResp.Status != types.VerifySoftwareTokenResponseTypeSuccess {
			return fmt.Errorf("verification status %v, want SUCCESS", verifyResp.Status)
		}

		// The enrolled factor turns the sign-in into SOFTWARE_TOKEN_MFA.
		authResp, err = passwordAuth()
		if err != nil {
			return err
		}
		if authResp.ChallengeName != types.ChallengeNameTypeSoftwareTokenMfa {
			return fmt.Errorf("expected SOFTWARE_TOKEN_MFA for the enrolled user, got %v (AuthenticationResult %v)", authResp.ChallengeName, authResp.AuthenticationResult != nil)
		}
		totpSession := aws.ToString(authResp.Session)

		// A wrong TOTP code is rejected; the code sits outside the ±1-step
		// drift window the validator allows, and the session stays
		// answerable within the failed-attempt budget.
		_, err = tc.client.RespondToAuthChallenge(tc.ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
			ClientId:      aws.String(mfaClientID),
			ChallengeName: types.ChallengeNameTypeSoftwareTokenMfa,
			Session:       aws.String(totpSession),
			ChallengeResponses: map[string]string{
				"USERNAME":                mfaUser,
				"SOFTWARE_TOKEN_MFA_CODE": totpCodeAt([]byte(secretCode), 2),
			},
		})
		if err == nil {
			return fmt.Errorf("wrong TOTP code accepted")
		}
		var codeMismatch *types.CodeMismatchException
		if !errors.As(err, &codeMismatch) {
			return fmt.Errorf("expected CodeMismatchException for the wrong TOTP code, got %v", err)
		}

		// The correct TOTP code completes the sign-in.
		challengeResp, err := tc.client.RespondToAuthChallenge(tc.ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
			ClientId:      aws.String(mfaClientID),
			ChallengeName: types.ChallengeNameTypeSoftwareTokenMfa,
			Session:       aws.String(totpSession),
			ChallengeResponses: map[string]string{
				"USERNAME":                mfaUser,
				"SOFTWARE_TOKEN_MFA_CODE": totpCodeAt([]byte(secretCode), 0),
			},
		})
		if err != nil {
			return fmt.Errorf("RespondToAuthChallenge with the correct TOTP: %v", err)
		}
		if challengeResp.AuthenticationResult == nil || challengeResp.AuthenticationResult.AccessToken == nil {
			return fmt.Errorf("AccessToken missing after the TOTP challenge")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminDeleteSoftwareToken_RemovesTOTPFactor", func() error {
		return r.runAdminDeleteSoftwareTokenTest(tc)
	}))

	return results
}

// runAdminDeleteSoftwareTokenTest pins the admin-side TOTP deletion: the
// first call removes the registration, a second call reports
// ResourceNotFoundException because no registration remains, and an unknown
// username keeps its UserNotFoundException.
func (r *TestRunner) runAdminDeleteSoftwareTokenTest(tc *cognitoIDPContext) error {
	clientID, cleanupClient, err := tc.createPoolClient(tc.userPoolID, tc.unique("deltoken-client"))
	if err != nil {
		return err
	}
	defer cleanupClient()

	username := tc.unique("deltoken-user")
	cleanupUser, err := tc.createConfirmedUser(username, "DelPass123!")
	if err != nil {
		return err
	}
	defer cleanupUser()

	authResp, err := tc.client.InitiateAuth(tc.ctx, &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(clientID),
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": "DelPass123!",
		},
	})
	if err != nil {
		return err
	}
	if authResp.AuthenticationResult == nil || authResp.AuthenticationResult.AccessToken == nil {
		return fmt.Errorf("AccessToken missing after password sign-in")
	}
	accessToken := *authResp.AuthenticationResult.AccessToken

	assoc, err := tc.client.AssociateSoftwareToken(tc.ctx, &cognitoidentityprovider.AssociateSoftwareTokenInput{
		AccessToken: aws.String(accessToken),
	})
	if err != nil {
		return fmt.Errorf("AssociateSoftwareToken: %v", err)
	}
	verify, err := tc.client.VerifySoftwareToken(tc.ctx, &cognitoidentityprovider.VerifySoftwareTokenInput{
		AccessToken: aws.String(accessToken),
		UserCode:    aws.String(totpCodeAt([]byte(*assoc.SecretCode), 0)),
	})
	if err != nil {
		return fmt.Errorf("VerifySoftwareToken: %v", err)
	}
	if verify.Status != types.VerifySoftwareTokenResponseTypeSuccess {
		return fmt.Errorf("VerifySoftwareToken status: got %q, want SUCCESS", verify.Status)
	}

	if _, err := tc.client.AdminDeleteSoftwareToken(tc.ctx, &cognitoidentityprovider.AdminDeleteSoftwareTokenInput{
		UserPoolId: aws.String(tc.userPoolID),
		Username:   aws.String(username),
	}); err != nil {
		return fmt.Errorf("AdminDeleteSoftwareToken: %v", err)
	}

	_, err = tc.client.AdminDeleteSoftwareToken(tc.ctx, &cognitoidentityprovider.AdminDeleteSoftwareTokenInput{
		UserPoolId: aws.String(tc.userPoolID),
		Username:   aws.String(username),
	})
	if err := expectAWSErrorCode(err, "ResourceNotFoundException"); err != nil {
		return fmt.Errorf("second AdminDeleteSoftwareToken: %w", err)
	}

	_, err = tc.client.AdminDeleteSoftwareToken(tc.ctx, &cognitoidentityprovider.AdminDeleteSoftwareTokenInput{
		UserPoolId: aws.String(tc.userPoolID),
		Username:   aws.String(tc.unique("no-such-user")),
	})
	return expectAWSErrorCode(err, "UserNotFoundException")
}
