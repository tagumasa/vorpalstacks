package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

func (r *TestRunner) cognitoIdentityDeveloperTests(ctx context.Context, client *cognitoidentity.Client, poolID, identityID string) []TestResult {
	var results []TestResult

	var devIdentityID string
	results = append(results, r.RunTest("cognito-identity", "GetOpenIdTokenForDeveloperIdentity", func() error {
		resp, err := client.GetOpenIdTokenForDeveloperIdentity(ctx, &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"my-dev-provider": "dev-user-1",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		if resp.Token == nil || *resp.Token == "" {
			return fmt.Errorf("token is nil or empty")
		}
		devIdentityID = *resp.IdentityId
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "LookupDeveloperIdentity", func() error {
		resp, err := client.LookupDeveloperIdentity(ctx, &cognitoidentity.LookupDeveloperIdentityInput{
			IdentityPoolId:          aws.String(poolID),
			DeveloperUserIdentifier: aws.String("dev-user-1"),
		})
		if err != nil {
			return err
		}
		if len(resp.DeveloperUserIdentifierList) < 1 {
			return fmt.Errorf("expected at least 1 developer user identifier")
		}
		found := false
		for _, id := range resp.DeveloperUserIdentifierList {
			if id == "dev-user-1" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected dev-user-1 in DeveloperUserIdentifierList")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdTokenForDeveloperIdentity_Reuse", func() error {
		resp, err := client.GetOpenIdTokenForDeveloperIdentity(ctx, &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"my-dev-provider": "dev-user-1",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		if *resp.IdentityId != devIdentityID {
			return fmt.Errorf("expected same IdentityId for same developer user, got %s want %s", *resp.IdentityId, devIdentityID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "MergeDeveloperIdentities", func() error {
		_, err := client.GetOpenIdTokenForDeveloperIdentity(ctx, &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"my-dev-provider": "dev-user-2",
			},
		})
		if err != nil {
			return err
		}

		resp, err := client.MergeDeveloperIdentities(ctx, &cognitoidentity.MergeDeveloperIdentitiesInput{
			SourceUserIdentifier:      aws.String("dev-user-1"),
			DestinationUserIdentifier: aws.String("dev-user-2"),
			DeveloperProviderName:     aws.String("my-dev-provider"),
			IdentityPoolId:            aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty after merge")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "UnlinkDeveloperIdentity", func() error {
		_, err := client.UnlinkDeveloperIdentity(ctx, &cognitoidentity.UnlinkDeveloperIdentityInput{
			IdentityPoolId:          aws.String(poolID),
			IdentityId:              aws.String(identityID),
			DeveloperProviderName:   aws.String("my-dev-provider"),
			DeveloperUserIdentifier: aws.String("dev-user-1"),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "UnlinkIdentity", func() error {
		_, err := client.UnlinkIdentity(ctx, &cognitoidentity.UnlinkIdentityInput{
			IdentityId: aws.String(identityID),
			Logins: map[string]string{
				"graph.facebook.com": "token",
			},
			LoginsToRemove: []string{"graph.facebook.com"},
		})
		if err != nil {
			return err
		}
		resp, err := client.DescribeIdentity(ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String(identityID),
		})
		if err != nil {
			return err
		}
		for _, l := range resp.Logins {
			if l == "graph.facebook.com" {
				return fmt.Errorf("graph.facebook.com should have been unlinked")
			}
		}
		return nil
	}))

	return results
}
