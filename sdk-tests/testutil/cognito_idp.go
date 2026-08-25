package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoIDPTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "CreateIdentityProvider", func() error {
		resp, err := tc.client.CreateIdentityProvider(tc.ctx, &cognitoidentityprovider.CreateIdentityProviderInput{
			UserPoolId:   aws.String(tc.userPoolID),
			ProviderName: aws.String("TestProvider"),
			ProviderType: types.IdentityProviderTypeType("Facebook"),
			ProviderDetails: map[string]string{
				"client_id":        "test-client-id",
				"client_secret":    "test-client-secret",
				"authorize_scopes": "public_profile,email",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityProvider == nil {
			return fmt.Errorf("IdentityProvider is nil")
		}
		if resp.IdentityProvider.ProviderName == nil || *resp.IdentityProvider.ProviderName != "TestProvider" {
			return fmt.Errorf("ProviderName mismatch: got %v, want TestProvider", resp.IdentityProvider.ProviderName)
		}
		if resp.IdentityProvider.ProviderType != types.IdentityProviderTypeTypeFacebook {
			return fmt.Errorf("ProviderType mismatch: got %v, want Facebook", resp.IdentityProvider.ProviderType)
		}
		if resp.IdentityProvider.UserPoolId == nil || *resp.IdentityProvider.UserPoolId != tc.userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.IdentityProvider.UserPoolId, tc.userPoolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListIdentityProviders", func() error {
		resp, err := tc.client.ListIdentityProviders(tc.ctx, &cognitoidentityprovider.ListIdentityProvidersInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if len(resp.Providers) == 0 {
			return fmt.Errorf("expected at least one identity provider")
		}
		found := false
		for _, p := range resp.Providers {
			if p.ProviderName != nil && *p.ProviderName == "TestProvider" {
				found = true
				if p.ProviderType == "" {
					return fmt.Errorf("ProviderType is empty in listing")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("TestProvider not found in ListIdentityProviders")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeIdentityProvider", func() error {
		resp, err := tc.client.DescribeIdentityProvider(tc.ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   aws.String(tc.userPoolID),
			ProviderName: aws.String("TestProvider"),
		})
		if err != nil {
			return fmt.Errorf("DescribeIdentityProvider failed: %v", err)
		}
		if resp.IdentityProvider == nil {
			return fmt.Errorf("IdentityProvider is nil")
		}
		if resp.IdentityProvider.ProviderName == nil || *resp.IdentityProvider.ProviderName != "TestProvider" {
			return fmt.Errorf("ProviderName mismatch: got %v", resp.IdentityProvider.ProviderName)
		}
		if resp.IdentityProvider.ProviderType != types.IdentityProviderTypeTypeFacebook {
			return fmt.Errorf("ProviderType mismatch: got %v", resp.IdentityProvider.ProviderType)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateIdentityProvider", func() error {
		_, err := tc.client.UpdateIdentityProvider(tc.ctx, &cognitoidentityprovider.UpdateIdentityProviderInput{
			UserPoolId:      aws.String(tc.userPoolID),
			ProviderName:    aws.String("TestProvider"),
			ProviderDetails: map[string]string{"updated_key": "updated_value"},
		})
		if err != nil {
			return fmt.Errorf("UpdateIdentityProvider failed: %v", err)
		}
		descResp, err := tc.client.DescribeIdentityProvider(tc.ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   aws.String(tc.userPoolID),
			ProviderName: aws.String("TestProvider"),
		})
		if err != nil {
			return fmt.Errorf("DescribeIdentityProvider after update failed: %v", err)
		}
		if descResp.IdentityProvider.ProviderDetails == nil {
			return fmt.Errorf("ProviderDetails is nil after update")
		}
		if descResp.IdentityProvider.ProviderDetails["updated_key"] != "updated_value" {
			return fmt.Errorf("ProviderDetails not updated")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteIdentityProvider", func() error {
		delProvider := tc.unique("del-provider")
		_, err := tc.client.CreateIdentityProvider(tc.ctx, &cognitoidentityprovider.CreateIdentityProviderInput{
			UserPoolId:      aws.String(tc.userPoolID),
			ProviderName:    aws.String(delProvider),
			ProviderType:    types.IdentityProviderTypeTypeGoogle,
			ProviderDetails: map[string]string{"client_id": "test"},
		})
		if err != nil {
			return fmt.Errorf("create provider: %v", err)
		}
		_, err = tc.client.DeleteIdentityProvider(tc.ctx, &cognitoidentityprovider.DeleteIdentityProviderInput{
			UserPoolId:   aws.String(tc.userPoolID),
			ProviderName: aws.String(delProvider),
		})
		if err != nil {
			return fmt.Errorf("DeleteIdentityProvider failed: %v", err)
		}
		_, err = tc.client.DescribeIdentityProvider(tc.ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   aws.String(tc.userPoolID),
			ProviderName: aws.String(delProvider),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
