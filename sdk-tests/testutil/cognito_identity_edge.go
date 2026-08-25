package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

func (r *TestRunner) cognitoIdentityEdgeTests(tc *cognitoIdentityContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentityPool_NonExistent", func() error {
		_, err := tc.client.DescribeIdentityPool(tc.ctx, &cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(fmt.Sprintf("%s:00000000-0000-0000-0000-000000000000", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentityPool_NonExistent", func() error {
		_, err := tc.client.DeleteIdentityPool(tc.ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: aws.String(fmt.Sprintf("%s:00000000-0000-0000-0000-000000000000", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentity_NonExistent", func() error {
		_, err := tc.client.DescribeIdentity(tc.ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String(fmt.Sprintf("%s:00000000-0000-0000-0000-000000000000", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetId_NonExistentPool", func() error {
		_, err := tc.client.GetId(tc.ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(fmt.Sprintf("%s:00000000-0000-0000-0000-000000000000", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "CreateIdentityPool_WithTags", func() error {
		name := tc.unique("test-idpool-tags")
		resp, err := tc.client.CreateIdentityPool(tc.ctx, &cognitoidentity.CreateIdentityPoolInput{
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
		tc.client.DeleteIdentityPool(tc.ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: resp.IdentityPoolId,
		})
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity_NonExistent", func() error {
		_, err := tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: aws.String(fmt.Sprintf("%s:00000000-0000-0000-0000-000000000000", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken_NonExistent", func() error {
		_, err := tc.client.GetOpenIdToken(tc.ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String(fmt.Sprintf("%s:00000000-0000-0000-0000-000000000000", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentities_NonExistent", func() error {
		resp, err := tc.client.DeleteIdentities(tc.ctx, &cognitoidentity.DeleteIdentitiesInput{
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
		_, err := tc.client.GetIdentityPoolRoles(tc.ctx, &cognitoidentity.GetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(fmt.Sprintf("%s:00000000-0000-0000-0000-000000000000", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetPrincipalTagAttributeMap_NonExistentPool", func() error {
		_, err := tc.client.GetPrincipalTagAttributeMap(tc.ctx, &cognitoidentity.GetPrincipalTagAttributeMapInput{
			IdentityPoolId:       aws.String("us-east-1:00000000-0000-0000-0000-000000000000"),
			IdentityProviderName: aws.String("provider"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "ListIdentityPools_Pagination", func() error {
		var cleanupPools []func()
		defer func() {
			for _, cleanup := range cleanupPools {
				cleanup()
			}
		}()
		for i := 0; i < 5; i++ {
			_, cleanupPool, err := tc.createIdPool(fmt.Sprintf("PagIdPool-%s-%d", tc.ts, i))
			if err != nil {
				return err
			}
			cleanupPools = append(cleanupPools, cleanupPool)
		}

		allPools, err := paginate(func(next *string) ([]string, *string, error) {
			resp, err := tc.client.ListIdentityPools(tc.ctx, &cognitoidentity.ListIdentityPoolsInput{
				MaxResults: aws.Int32(2),
				NextToken:  next,
			})
			if err != nil {
				return nil, nil, err
			}
			var names []string
			for _, p := range resp.IdentityPools {
				if p.IdentityPoolName != nil {
					names = append(names, *p.IdentityPoolName)
				}
			}
			return names, resp.NextToken, nil
		})
		if err != nil {
			return fmt.Errorf("list identity pools page: %v", err)
		}

		sawOwn := 0
		for _, name := range allPools {
			if strings.Contains(name, "PagIdPool-"+tc.ts) {
				sawOwn++
			}
		}
		if sawOwn != 5 {
			return fmt.Errorf("expected 5 paginated identity pools, got %d", sawOwn)
		}
		return nil
	}))

	return results
}
