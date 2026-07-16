package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoPoolCoreTests(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "DescribeUserPool", func() error {
		resp, err := client.DescribeUserPool(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.UserPool == nil {
			return fmt.Errorf("UserPool is nil")
		}
		if resp.UserPool.Id == nil || *resp.UserPool.Id != userPoolID {
			return fmt.Errorf("UserPool.Id mismatch: got %v, want %s", resp.UserPool.Id, userPoolID)
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
		var found bool
		var nextToken *string
		for {
			resp, err := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{
				MaxResults: aws.Int32(10),
				NextToken:  nextToken,
			})
			if err != nil {
				return err
			}
			for _, pool := range resp.UserPools {
				if pool.Id != nil && *pool.Id == userPoolID {
					found = true
					if pool.Name == nil || *pool.Name == "" {
						return fmt.Errorf("pool Name is nil or empty in listing")
					}
				}
			}
			if found || resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		if !found {
			return fmt.Errorf("created pool %s not found in ListUserPools", userPoolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateUserPool", func() error {
		_, err := client.UpdateUserPool(ctx, &cognitoidentityprovider.UpdateUserPoolInput{
			UserPoolId: aws.String(userPoolID),
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
		descResp, err := client.DescribeUserPool(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(userPoolID),
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
		_, err := client.SetUserPoolMfaConfig(ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
			UserPoolId:       aws.String(userPoolID),
			MfaConfiguration: types.UserPoolMfaTypeOn,
		})
		if err != nil {
			return fmt.Errorf("SetUserPoolMfaConfig failed: %v", err)
		}
		mfaResp, err := client.GetUserPoolMfaConfig(ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("GetUserPoolMfaConfig failed: %v", err)
		}
		if mfaResp.MfaConfiguration != types.UserPoolMfaTypeOn {
			return fmt.Errorf("expected MfaConfiguration ON, got %v", mfaResp.MfaConfiguration)
		}
		_, err = client.SetUserPoolMfaConfig(ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
			UserPoolId:       aws.String(userPoolID),
			MfaConfiguration: types.UserPoolMfaTypeOff,
		})
		if err != nil {
			return fmt.Errorf("SetUserPoolMfaConfig (OFF) failed: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GetUserPoolMfaConfig", func() error {
		resp, err := client.GetUserPoolMfaConfig(ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.MfaConfiguration != types.UserPoolMfaTypeOff {
			return fmt.Errorf("expected MfaConfiguration OFF after reset, got %v", resp.MfaConfiguration)
		}
		return nil
	}))

	domainName := fmt.Sprintf("test-domain-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("cognito", "CreateUserPoolDomain", func() error {
		resp, err := client.CreateUserPoolDomain(ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(userPoolID),
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
		resp, err := client.DescribeUserPoolDomain(ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		if resp.DomainDescription == nil {
			return fmt.Errorf("DomainDescription is nil")
		}
		if resp.DomainDescription.UserPoolId == nil || *resp.DomainDescription.UserPoolId != userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.DomainDescription.UserPoolId, userPoolID)
		}
		if resp.DomainDescription.Domain == nil || *resp.DomainDescription.Domain != domainName {
			return fmt.Errorf("domain mismatch: got %v, want %s", resp.DomainDescription.Domain, domainName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateUserPoolDomain", func() error {
		udDomain := fmt.Sprintf("ud-domain-%d", time.Now().UnixNano())
		_, err := client.CreateUserPoolDomain(ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			Domain:     aws.String(udDomain),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("create domain: %v", err)
		}
		defer client.DeleteUserPoolDomain(ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     aws.String(udDomain),
			UserPoolId: aws.String(userPoolID),
		})
		resp, err := client.UpdateUserPoolDomain(ctx, &cognitoidentityprovider.UpdateUserPoolDomainInput{
			Domain:     aws.String(udDomain),
			UserPoolId: aws.String(userPoolID),
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
		_, err := client.DeleteUserPoolDomain(ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
			Domain:     aws.String(domainName),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return err
		}
		_, err = client.DescribeUserPoolDomain(ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
			Domain: aws.String(domainName),
		})
		if err == nil {
			return fmt.Errorf("expected error describing deleted domain")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GetCSVHeader", func() error {
		resp, err := client.GetCSVHeader(ctx, &cognitoidentityprovider.GetCSVHeaderInput{
			UserPoolId: aws.String(userPoolID),
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
		resp, err := client.DescribeRiskConfiguration(ctx, &cognitoidentityprovider.DescribeRiskConfigurationInput{
			UserPoolId: aws.String(userPoolID),
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
		newPoolName := fmt.Sprintf("test-pool-tags-%d", time.Now().UnixNano())
		newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(newPoolName),
		})
		if err != nil {
			return err
		}
		if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
			return fmt.Errorf("new pool Arn is nil")
		}
		_, err = client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			Tags: map[string]string{
				"Environment": "test",
				"Owner":       "test-user",
			},
		})
		if err != nil {
			return err
		}
		listResp, listErr := client.ListTagsForResource(ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: newPool.UserPool.Arn,
		})
		if listErr != nil {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return listErr
		}
		if listResp.Tags == nil {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("tags is nil after tagging")
		}
		if listResp.Tags["Environment"] != "test" {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("tag Environment not found after TagResource")
		}
		client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: newPool.UserPool.Id,
		})
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListTagsForResource", func() error {
		newPoolName := fmt.Sprintf("test-pool-listtags-%d", time.Now().UnixNano())
		newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(newPoolName),
		})
		if err != nil {
			return err
		}
		if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
			return fmt.Errorf("new pool Arn is nil")
		}
		_, err = client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			Tags: map[string]string{
				"Test": "value",
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: newPool.UserPool.Arn,
		})
		if err != nil {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return err
		}
		if resp.Tags == nil {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("tags is nil")
		}
		if resp.Tags["Test"] != "value" {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("expected tag Test=value, got %v", resp.Tags["Test"])
		}
		client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: newPool.UserPool.Id,
		})
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UntagResource", func() error {
		newPoolName := fmt.Sprintf("test-pool-untag-%d", time.Now().UnixNano())
		newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(newPoolName),
		})
		if err != nil {
			return err
		}
		if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
			return fmt.Errorf("new pool Arn is nil")
		}
		_, err = client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			Tags: map[string]string{
				"Test": "value",
			},
		})
		if err != nil {
			return err
		}
		_, err = client.UntagResource(ctx, &cognitoidentityprovider.UntagResourceInput{
			ResourceArn: newPool.UserPool.Arn,
			TagKeys:     []string{"Test"},
		})
		if err != nil {
			return err
		}
		listResp, listErr := client.ListTagsForResource(ctx, &cognitoidentityprovider.ListTagsForResourceInput{
			ResourceArn: newPool.UserPool.Arn,
		})
		if listErr != nil {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return listErr
		}
		if _, exists := listResp.Tags["Test"]; exists {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return fmt.Errorf("tag Test should have been removed after UntagResource")
		}
		client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: newPool.UserPool.Id,
		})
		return nil
	}))

	return results
}
