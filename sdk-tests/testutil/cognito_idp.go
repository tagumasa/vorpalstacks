package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoIDPTests(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "CreateIdentityProvider", func() error {
		resp, err := client.CreateIdentityProvider(ctx, &cognitoidentityprovider.CreateIdentityProviderInput{
			UserPoolId:   aws.String(userPoolID),
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
		if resp.IdentityProvider.UserPoolId == nil || *resp.IdentityProvider.UserPoolId != userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.IdentityProvider.UserPoolId, userPoolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListIdentityProviders", func() error {
		resp, err := client.ListIdentityProviders(ctx, &cognitoidentityprovider.ListIdentityProvidersInput{
			UserPoolId: aws.String(userPoolID),
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
		resp, err := client.DescribeIdentityProvider(ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   aws.String(userPoolID),
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
		_, err := client.UpdateIdentityProvider(ctx, &cognitoidentityprovider.UpdateIdentityProviderInput{
			UserPoolId:      aws.String(userPoolID),
			ProviderName:    aws.String("TestProvider"),
			ProviderDetails: map[string]string{"updated_key": "updated_value"},
		})
		if err != nil {
			return fmt.Errorf("UpdateIdentityProvider failed: %v", err)
		}
		descResp, err := client.DescribeIdentityProvider(ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   aws.String(userPoolID),
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
		delProvider := fmt.Sprintf("del-provider-%d", time.Now().UnixNano())
		_, err := client.CreateIdentityProvider(ctx, &cognitoidentityprovider.CreateIdentityProviderInput{
			UserPoolId:      aws.String(userPoolID),
			ProviderName:    aws.String(delProvider),
			ProviderType:    types.IdentityProviderTypeTypeGoogle,
			ProviderDetails: map[string]string{"client_id": "test"},
		})
		if err != nil {
			return fmt.Errorf("create provider: %v", err)
		}
		_, err = client.DeleteIdentityProvider(ctx, &cognitoidentityprovider.DeleteIdentityProviderInput{
			UserPoolId:   aws.String(userPoolID),
			ProviderName: aws.String(delProvider),
		})
		if err != nil {
			return fmt.Errorf("DeleteIdentityProvider failed: %v", err)
		}
		_, err = client.DescribeIdentityProvider(ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   aws.String(userPoolID),
			ProviderName: aws.String(delProvider),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
