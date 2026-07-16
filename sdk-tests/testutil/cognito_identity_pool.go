package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
)

func (r *TestRunner) cognitoIdentityPoolTests(ctx context.Context, client *cognitoidentity.Client, poolID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentityPool", func() error {
		resp, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if resp.IdentityPoolName == nil || *resp.IdentityPoolName == "" {
			return fmt.Errorf("IdentityPoolName is nil or empty")
		}
		if resp.IdentityPoolId == nil || *resp.IdentityPoolId != poolID {
			return fmt.Errorf("IdentityPoolId mismatch: got %v, want %s", resp.IdentityPoolId, poolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "CreateIdentityPool_WithOptions", func() error {
		name := fmt.Sprintf("test-idpool-opts-%d", time.Now().UnixNano())
		resp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
			IdentityPoolName:               aws.String(name),
			AllowUnauthenticatedIdentities: false,
			AllowClassicFlow:               aws.Bool(true),
			DeveloperProviderName:          aws.String("my-dev-provider"),
			OpenIdConnectProviderARNs:      []string{"arn:aws:iam::123456789012:oidc-provider/example.com"},
			SamlProviderARNs:               []string{"arn:aws:iam::123456789012:saml-provider/example.com"},
			CognitoIdentityProviders: []types.CognitoIdentityProvider{
				{
					ProviderName:         aws.String("cognito-idp.us-east-1.amazonaws.com/us-east-1_xxxxx"),
					ClientId:             aws.String("abc123"),
					ServerSideTokenCheck: aws.Bool(true),
				},
			},
			SupportedLoginProviders: map[string]string{
				"graph.facebook.com": "1234567890",
			},
		})
		if err != nil {
			return err
		}
		if *resp.IdentityPoolName != name {
			return fmt.Errorf("expected name %s, got %s", name, *resp.IdentityPoolName)
		}
		if resp.AllowUnauthenticatedIdentities {
			return fmt.Errorf("expected AllowUnauthenticatedIdentities false")
		}
		if resp.AllowClassicFlow == nil || !*resp.AllowClassicFlow {
			return fmt.Errorf("expected AllowClassicFlow true")
		}
		if resp.DeveloperProviderName == nil || *resp.DeveloperProviderName != "my-dev-provider" {
			return fmt.Errorf("expected DeveloperProviderName my-dev-provider")
		}
		if len(resp.OpenIdConnectProviderARNs) != 1 {
			return fmt.Errorf("expected 1 OpenIdConnectProviderARN")
		}
		if len(resp.SamlProviderARNs) != 1 {
			return fmt.Errorf("expected 1 SamlProviderARN")
		}
		if len(resp.CognitoIdentityProviders) != 1 {
			return fmt.Errorf("expected 1 CognitoIdentityProvider")
		}
		if len(resp.SupportedLoginProviders) != 1 {
			return fmt.Errorf("expected 1 SupportedLoginProvider")
		}
		client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: resp.IdentityPoolId,
		})
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "ListIdentityPools", func() error {
		var found bool
		var nextToken *string
		for {
			resp, err := client.ListIdentityPools(ctx, &cognitoidentity.ListIdentityPoolsInput{
				MaxResults: aws.Int32(10),
				NextToken:  nextToken,
			})
			if err != nil {
				return err
			}
			for _, p := range resp.IdentityPools {
				if p.IdentityPoolId != nil && *p.IdentityPoolId == poolID {
					found = true
					if p.IdentityPoolName == nil || *p.IdentityPoolName == "" {
						return fmt.Errorf("IdentityPoolName is nil or empty in listing")
					}
				}
			}
			if found || resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		if !found {
			return fmt.Errorf("created pool %s not found in ListIdentityPools", poolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "UpdateIdentityPool", func() error {
		descResp, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return fmt.Errorf("describe before update: %v", err)
		}
		newName := *descResp.IdentityPoolName + "-updated"
		_, err = client.UpdateIdentityPool(ctx, &cognitoidentity.UpdateIdentityPoolInput{
			IdentityPoolId:                 aws.String(poolID),
			IdentityPoolName:               aws.String(newName),
			AllowUnauthenticatedIdentities: false,
		})
		if err != nil {
			return err
		}
		resp, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityPoolName != newName {
			return fmt.Errorf("expected updated name %s, got %s", newName, *resp.IdentityPoolName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "UpdateIdentityPool_ReplacesTags", func() error {
		name := fmt.Sprintf("test-tagreplace-%d", time.Now().UnixNano())
		createResp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
			IdentityPoolName:               aws.String(name),
			AllowUnauthenticatedIdentities: true,
			IdentityPoolTags: map[string]string{
				"Old1": "v1",
				"Old2": "v2",
			},
		})
		if err != nil {
			return err
		}
		tagPoolID := *createResp.IdentityPoolId
		defer client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{IdentityPoolId: aws.String(tagPoolID)})

		_, err = client.UpdateIdentityPool(ctx, &cognitoidentity.UpdateIdentityPoolInput{
			IdentityPoolId:                 aws.String(tagPoolID),
			IdentityPoolName:               aws.String(name + "-v2"),
			AllowUnauthenticatedIdentities: true,
			IdentityPoolTags: map[string]string{
				"Old2": "v2-updated",
				"New":  "tag",
			},
		})
		if err != nil {
			return err
		}

		descResp, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(tagPoolID),
		})
		if err != nil {
			return err
		}
		if descResp.IdentityPoolTags == nil {
			return fmt.Errorf("IdentityPoolTags is nil after update")
		}
		if _, exists := descResp.IdentityPoolTags["Old1"]; exists {
			return fmt.Errorf("Old1 should have been removed by tag replacement")
		}
		if descResp.IdentityPoolTags["Old2"] != "v2-updated" {
			return fmt.Errorf("expected Old2=v2-updated, got %s", descResp.IdentityPoolTags["Old2"])
		}
		if descResp.IdentityPoolTags["New"] != "tag" {
			return fmt.Errorf("expected New=tag, got %s", descResp.IdentityPoolTags["New"])
		}
		if len(descResp.IdentityPoolTags) != 2 {
			return fmt.Errorf("expected 2 tags after replacement, got %d", len(descResp.IdentityPoolTags))
		}
		return nil
	}))

	return results
}
