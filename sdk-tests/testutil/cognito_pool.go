package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoPoolCoreTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "DescribeUserPool", func() error {
		resp, err := tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.UserPool == nil {
			return fmt.Errorf("UserPool is nil")
		}
		if resp.UserPool.Id == nil || *resp.UserPool.Id != tc.userPoolID {
			return fmt.Errorf("UserPool.Id mismatch: got %v, want %s", resp.UserPool.Id, tc.userPoolID)
		}
		if resp.UserPool.Name == nil || *resp.UserPool.Name == "" {
			return fmt.Errorf("UserPool.Name is nil or empty")
		}
		if resp.UserPool.Arn == nil || *resp.UserPool.Arn == "" {
			return fmt.Errorf("UserPool.Arn is nil or empty")
		}
		if resp.UserPool.Policies == nil || resp.UserPool.Policies.PasswordPolicy == nil {
			return fmt.Errorf("PasswordPolicy is nil")
		}
		if resp.UserPool.Policies.PasswordPolicy.MinimumLength == nil || *resp.UserPool.Policies.PasswordPolicy.MinimumLength != 8 {
			return fmt.Errorf("MinimumLength mismatch: got %v, want 8", resp.UserPool.Policies.PasswordPolicy.MinimumLength)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListUserPools", func() error {
		pools, err := paginate(func(next *string) ([]types.UserPoolDescriptionType, *string, error) {
			resp, err := tc.client.ListUserPools(tc.ctx, &cognitoidentityprovider.ListUserPoolsInput{
				MaxResults: aws.Int32(10),
				NextToken:  next,
			})
			if err != nil {
				return nil, nil, err
			}
			return resp.UserPools, resp.NextToken, nil
		})
		if err != nil {
			return err
		}
		found := false
		for _, pool := range pools {
			if pool.Id != nil && *pool.Id == tc.userPoolID {
				found = true
				if pool.Name == nil || *pool.Name == "" {
					return fmt.Errorf("pool Name is nil or empty in listing")
				}
			}
		}
		if !found {
			return fmt.Errorf("created pool %s not found in ListUserPools", tc.userPoolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateUserPool", func() error {
		_, err := tc.client.UpdateUserPool(tc.ctx, &cognitoidentityprovider.UpdateUserPoolInput{
			UserPoolId: aws.String(tc.userPoolID),
			Policies: &types.UserPoolPolicyType{
				PasswordPolicy: &types.PasswordPolicyType{
					MinimumLength:    aws.Int32(10),
					RequireUppercase: true,
					RequireLowercase: true,
					RequireNumbers:   true,
					RequireSymbols:   true,
				},
			},
		})
		if err != nil {
			return err
		}
		descResp, err := tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return fmt.Errorf("DescribeUserPool after update: %v", err)
		}
		pp := descResp.UserPool.Policies.PasswordPolicy
		if pp.MinimumLength == nil || *pp.MinimumLength != 10 {
			return fmt.Errorf("MinimumLength not updated: got %v, want 10", pp.MinimumLength)
		}
		if !pp.RequireSymbols {
			return fmt.Errorf("RequireSymbols not updated to true")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "SetUserPoolMfaConfig", func() error {
		_, err := tc.client.SetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
			UserPoolId:       aws.String(tc.userPoolID),
			MfaConfiguration: types.UserPoolMfaTypeOn,
		})
		if err != nil {
			return fmt.Errorf("SetUserPoolMfaConfig failed: %v", err)
		}
		mfaResp, err := tc.client.GetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return fmt.Errorf("GetUserPoolMfaConfig failed: %v", err)
		}
		if mfaResp.MfaConfiguration != types.UserPoolMfaTypeOn {
			return fmt.Errorf("expected MfaConfiguration ON, got %v", mfaResp.MfaConfiguration)
		}
		_, err = tc.client.SetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
			UserPoolId:       aws.String(tc.userPoolID),
			MfaConfiguration: types.UserPoolMfaTypeOff,
		})
		if err != nil {
			return fmt.Errorf("SetUserPoolMfaConfig (OFF) failed: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GetUserPoolMfaConfig", func() error {
		resp, err := tc.client.GetUserPoolMfaConfig(tc.ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.MfaConfiguration != types.UserPoolMfaTypeOff {
			return fmt.Errorf("expected MfaConfiguration OFF after reset, got %v", resp.MfaConfiguration)
		}
		return nil
	}))

	domainName := tc.unique("test-domain")
	results = append(results, r.RunTest("cognito", "CreateUserPoolDomain", func() error {
		resp, err := tc.client.CreateUserPoolDomain(tc.ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.CloudFrontDomain == nil || *resp.CloudFrontDomain == "" {
			return fmt.Errorf("CloudFrontDomain is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeUserPoolDomain", func() error {
		resp, err := tc.client.DescribeUserPoolDomain(tc.ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		if resp.DomainDescription == nil {
			return fmt.Errorf("DomainDescription is nil")
		}
		if resp.DomainDescription.UserPoolId == nil || *resp.DomainDescription.UserPoolId != tc.userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.DomainDescription.UserPoolId, tc.userPoolID)
		}
		if resp.DomainDescription.Domain == nil || *resp.DomainDescription.Domain != domainName {
			return fmt.Errorf("domain mismatch: got %v, want %s", resp.DomainDescription.Domain, domainName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateUserPoolDomain", func() error {
		udDomain := tc.unique("ud-domain")
		_, err := tc.client.CreateUserPoolDomain(tc.ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(udDomain),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return fmt.Errorf("create domain: %v", err)
		}
		defer tc.client.DeleteUserPoolDomain(tc.ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     aws.String(udDomain),
			UserPoolId: aws.String(tc.userPoolID),
		})
		resp, err := tc.client.UpdateUserPoolDomain(tc.ctx, &cognitoidentityprovider.UpdateUserPoolDomainInput{
			Domain:     aws.String(udDomain),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return fmt.Errorf("UpdateUserPoolDomain failed: %v", err)
		}
		if resp.CloudFrontDomain == nil || *resp.CloudFrontDomain == "" {
			return fmt.Errorf("CloudFrontDomain is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteUserPoolDomain", func() error {
		_, err := tc.client.DeleteUserPoolDomain(tc.ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.DescribeUserPoolDomain(tc.ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if err == nil {
			return fmt.Errorf("expected error describing deleted domain")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GetCSVHeader", func() error {
		resp, err := tc.client.GetCSVHeader(tc.ctx, &cognitoidentityprovider.GetCSVHeaderInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if len(resp.CSVHeader) == 0 {
			return fmt.Errorf("expected non-empty CSV header")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeRiskConfiguration", func() error {
		resp, err := tc.client.DescribeRiskConfiguration(tc.ctx, &cognitoidentityprovider.DescribeRiskConfigurationInput{
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.RiskConfiguration == nil {
			return fmt.Errorf("RiskConfiguration is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "TagResource", func() error {
		newPoolName := tc.unique("test-pool-tags")
		newPool, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(newPoolName),
		})
		if err != nil {
			return err
		}
		if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
			return fmt.Errorf("new pool Arn is nil")
		}
		_, err = tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			Tags: map[string]string{
				"Environment": "test",
				"Owner":       "test-user",
			},
		})
		if err != nil {
			return err
		}
		listResp, listErr := tc.client.ListTagsForResource(tc.ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: newPool.UserPool.Arn,
		})
		if listErr != nil {
			tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return listErr
		}
		if listResp.Tags == nil {
			tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("tags is nil after tagging")
		}
		if listResp.Tags["Environment"] != "test" {
			tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("tag Environment not found after TagResource")
		}
		tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: newPool.UserPool.Id,
		})
		return nil
	}))

	results = append(results, r.RunTest("cognito", "TagResource_ReservedPrefixRejected", func() error {
		newPoolName := tc.unique("test-pool-tagres")
		newPool, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(newPoolName),
		})
		if err != nil {
			return err
		}
		if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
			return fmt.Errorf("new pool Arn is nil")
		}
		defer tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: newPool.UserPool.Id,
		})
		_, err = tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			Tags: map[string]string{
				"aws:reserved": "v",
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for aws:-prefixed tag key")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListTagsForResource", func() error {
		newPoolName := tc.unique("test-pool-listtags")
		newPool, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(newPoolName),
		})
		if err != nil {
			return err
		}
		if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
			return fmt.Errorf("new pool Arn is nil")
		}
		_, err = tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			Tags: map[string]string{
				"Test": "value",
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListTagsForResource(tc.ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: newPool.UserPool.Arn,
		})
		if err != nil {
			tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return err
		}
		if resp.Tags == nil {
			tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("tags is nil")
		}
		if resp.Tags["Test"] != "value" {
			tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("expected tag Test=value, got %v", resp.Tags["Test"])
		}
		tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: newPool.UserPool.Id,
		})
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UntagResource", func() error {
		newPoolName := tc.unique("test-pool-untag")
		newPool, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(newPoolName),
		})
		if err != nil {
			return err
		}
		if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
			return fmt.Errorf("new pool Arn is nil")
		}
		_, err = tc.client.TagResource(tc.ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			Tags: map[string]string{
				"Test": "value",
			},
		})
		if err != nil {
			return err
		}
		_, err = tc.client.UntagResource(tc.ctx, &cognitoidentityprovider.UntagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			TagKeys:     []string{"Test"},
		})
		if err != nil {
			return err
		}
		listResp, listErr := tc.client.ListTagsForResource(tc.ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: newPool.UserPool.Arn,
		})
		if listErr != nil {
			tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return listErr
		}
		if _, exists := listResp.Tags["Test"]; exists {
			tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("tag Test should have been removed after UntagResource")
		}
		tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: newPool.UserPool.Id,
		})
		return nil
	}))

	return results
}
