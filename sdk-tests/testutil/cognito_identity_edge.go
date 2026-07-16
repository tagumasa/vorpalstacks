package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

func (r *TestRunner) cognitoIdentityEdgeTests(ctx context.Context, client *cognitoidentity.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentityPool_NonExistent", func() error {
		_, err := client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentityPool_NonExistent", func() error {
		_, err := client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentity_NonExistent", func() error {
		_, err := client.DescribeIdentity(ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetId_NonExistentPool", func() error {
		_, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "CreateIdentityPool_WithTags", func() error {
		name := fmt.Sprintf("test-idpool-tags-%d", time.Now().UnixNano())
		resp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
			IdentityPoolName:               aws.String(name),
			AllowUnauthenticatedIdentities: true,
			IdentityPoolTags: map[string]string{
				"Env":  "production",
				"Cost": "high",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityPoolTags == nil {
			return fmt.Errorf("IdentityPoolTags is nil in CreateIdentityPool response")
		}
		if resp.IdentityPoolTags["Env"] != "production" {
			return fmt.Errorf("expected Env=production")
		}
		client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: resp.IdentityPoolId,
		})
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity_NonExistent", func() error {
		_, err := client.GetCredentialsForIdentity(ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken_NonExistent", func() error {
		_, err := client.GetOpenIdToken(ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentities_NonExistent", func() error {
		resp, err := client.DeleteIdentities(ctx, &cognitoidentity.DeleteIdentitiesInput{
			IdentityIdsToDelete: []string{"00000000-0000-0000-0000-000000000000"},
		})
		if err != nil {
			return err
		}
		if len(resp.UnprocessedIdentityIds) != 1 {
			return fmt.Errorf("expected 1 unprocessed identity")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetIdentityPoolRoles_NonExistent", func() error {
		_, err := client.GetIdentityPoolRoles(ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetPrincipalTagAttributeMap_NonExistentPool", func() error {
		_, err := client.GetPrincipalTagAttributeMap(ctx, &cognitoidentity.GetPrincipalTagAttributeMapInput{
			IdentityPoolId:       aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
			IdentityProviderName: aws.String("provider"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "ListIdentityPools_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgPools []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagIdPool-%s-%d", pgTs, i)
			resp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
				IdentityPoolName:               aws.String(name),
				AllowUnauthenticatedIdentities: true,
			})
			if err != nil {
				for _, pid := range pgPools {
					client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{IdentityPoolId: aws.String(pid)})
				}
				return fmt.Errorf("create identity pool %s: %v", name, err)
			}
			pgPools = append(pgPools, *resp.IdentityPoolId)
		}

		var allPools []string
		var nextToken *string
		for {
			resp, err := client.ListIdentityPools(ctx, &cognitoidentity.ListIdentityPoolsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, pid := range pgPools {
					client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{IdentityPoolId: aws.String(pid)})
				}
				return fmt.Errorf("list identity pools page: %v", err)
			}
			for _, p := range resp.IdentityPools {
				if p.IdentityPoolName != nil && strings.Contains(*p.IdentityPoolName, "PagIdPool-"+pgTs) {
					allPools = append(allPools, *p.IdentityPoolName)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, pid := range pgPools {
			client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{IdentityPoolId: aws.String(pid)})
		}
		if len(allPools) != 5 {
			return fmt.Errorf("expected 5 paginated identity pools, got %d", len(allPools))
		}
		return nil
	}))

	return results
}
