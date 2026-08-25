package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

func (r *TestRunner) cognitoIdentityCredentialsTests(tc *cognitoIdentityContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity", func() error {
		resp, err := tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: aws.String(tc.identityID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityId != tc.identityID {
			return fmt.Errorf("expected identity ID %s, got %s", tc.identityID, *resp.IdentityId)
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		if resp.Credentials.AccessKeyId == nil || *resp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("AccessKeyId is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken", func() error {
		resp, err := tc.client.GetOpenIdToken(tc.ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String(tc.identityID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityId != tc.identityID {
			return fmt.Errorf("expected identity ID %s, got %s", tc.identityID, *resp.IdentityId)
		}
		if resp.Token == nil || *resp.Token == "" {
			return fmt.Errorf("token is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken_WithLogins", func() error {
		resp, err := tc.client.GetOpenIdToken(tc.ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String(tc.identityID),
			Logins: map[string]string{
				"graph.facebook.com": "new-token",
			},
		})
		if err != nil {
			return err
		}
		if resp.Token == nil || *resp.Token == "" {
			return fmt.Errorf("token is nil or empty")
		}
		return nil
	}))

	return results
}
