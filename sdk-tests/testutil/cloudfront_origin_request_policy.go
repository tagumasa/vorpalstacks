package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfOriginRequestPolicyTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("cloudfront", "ListOriginRequestPolicies_VerifyFields", func() error {
		resp, err := client.ListOriginRequestPolicies(ctx, &cloudfront.ListOriginRequestPoliciesInput{
			MaxItems: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.OriginRequestPolicyList == nil {
			return fmt.Errorf("origin request policy list is nil")
		}
		if resp.OriginRequestPolicyList.MaxItems == nil {
			return fmt.Errorf("maxItems is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetOriginRequestPolicy_Managed", func() error {
		resp, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
			Id: aws.String("88a5eaf4-2fd4-4709-b370-b4c650ea3fcf"),
		})
		if err != nil {
			return err
		}
		if resp.OriginRequestPolicy == nil {
			return fmt.Errorf("origin request policy is nil")
		}
		if resp.OriginRequestPolicy.Id == nil {
			return fmt.Errorf("origin request policy ID is nil")
		}
		if resp.OriginRequestPolicy.OriginRequestPolicyConfig == nil {
			return fmt.Errorf("origin request policy config is nil")
		}
		return nil
	}))

	orpName := tc.uniquePrefix("test-orp")
	var orpID string
	var orpETag string
	results = append(results, tc.runner.RunTest("cloudfront", "CreateOriginRequestPolicy_VerifyFields", func() error {
		resp, err := client.CreateOriginRequestPolicy(ctx, &cloudfront.CreateOriginRequestPolicyInput{
			OriginRequestPolicyConfig: &types.OriginRequestPolicyConfig{
				Name:    aws.String(orpName),
				Comment: aws.String("Test ORP"),
				CookiesConfig: &types.OriginRequestPolicyCookiesConfig{
					CookieBehavior: types.OriginRequestPolicyCookieBehaviorNone,
				},
				HeadersConfig: &types.OriginRequestPolicyHeadersConfig{
					HeaderBehavior: types.OriginRequestPolicyHeaderBehaviorNone,
				},
				QueryStringsConfig: &types.OriginRequestPolicyQueryStringsConfig{
					QueryStringBehavior: types.OriginRequestPolicyQueryStringBehaviorNone,
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.OriginRequestPolicy == nil {
			return fmt.Errorf("origin request policy is nil")
		}
		if resp.OriginRequestPolicy.Id == nil || *resp.OriginRequestPolicy.Id == "" {
			return fmt.Errorf("origin request policy ID is nil or empty")
		}
		cfg := resp.OriginRequestPolicy.OriginRequestPolicyConfig
		if cfg == nil {
			return fmt.Errorf("origin request policy config is nil")
		}
		if cfg.Name == nil || *cfg.Name != orpName {
			return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(cfg.Name), orpName)
		}
		if cfg.CookiesConfig == nil || cfg.CookiesConfig.CookieBehavior != types.OriginRequestPolicyCookieBehaviorNone {
			return fmt.Errorf("cookie behavior mismatch")
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		orpID = *resp.OriginRequestPolicy.Id
		orpETag = *resp.ETag
		return nil
	}))

	if orpID != "" {
		results = append(results, tc.runner.RunTest("cloudfront", "GetOriginRequestPolicy_Custom_VerifyFields", func() error {
			resp, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
				Id: aws.String(orpID),
			})
			if err != nil {
				return err
			}
			if resp.OriginRequestPolicy == nil {
				return fmt.Errorf("origin request policy is nil")
			}
			if resp.OriginRequestPolicy.Id == nil || *resp.OriginRequestPolicy.Id != orpID {
				return fmt.Errorf("ID mismatch: got %q, want %q", aws.ToString(resp.OriginRequestPolicy.Id), orpID)
			}
			if resp.OriginRequestPolicy.LastModifiedTime == nil {
				return fmt.Errorf("lastModifiedTime is nil")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "GetOriginRequestPolicyConfig_VerifyFields", func() error {
			resp, err := client.GetOriginRequestPolicyConfig(ctx, &cloudfront.GetOriginRequestPolicyConfigInput{
				Id: aws.String(orpID),
			})
			if err != nil {
				return err
			}
			if resp.OriginRequestPolicyConfig == nil {
				return fmt.Errorf("origin request policy config is nil")
			}
			if resp.OriginRequestPolicyConfig.Name == nil || *resp.OriginRequestPolicyConfig.Name != orpName {
				return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(resp.OriginRequestPolicyConfig.Name), orpName)
			}
			if resp.ETag == nil || *resp.ETag == "" {
				return fmt.Errorf("ETag is nil or empty")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "UpdateOriginRequestPolicy_VerifyFields", func() error {
			resp, err := client.UpdateOriginRequestPolicy(ctx, &cloudfront.UpdateOriginRequestPolicyInput{
				Id:      aws.String(orpID),
				IfMatch: aws.String(orpETag),
				OriginRequestPolicyConfig: &types.OriginRequestPolicyConfig{
					Name:    aws.String("updated-orp"),
					Comment: aws.String("Updated test ORP"),
					CookiesConfig: &types.OriginRequestPolicyCookiesConfig{
						CookieBehavior: types.OriginRequestPolicyCookieBehaviorAll,
					},
					HeadersConfig: &types.OriginRequestPolicyHeadersConfig{
						HeaderBehavior: types.OriginRequestPolicyHeaderBehaviorAllViewer,
					},
					QueryStringsConfig: &types.OriginRequestPolicyQueryStringsConfig{
						QueryStringBehavior: types.OriginRequestPolicyQueryStringBehaviorAll,
					},
				},
			})
			if err != nil {
				return err
			}
			if resp.OriginRequestPolicy == nil {
				return fmt.Errorf("updated ORP is nil")
			}
			cfg := resp.OriginRequestPolicy.OriginRequestPolicyConfig
			if cfg == nil {
				return fmt.Errorf("ORP config is nil")
			}
			if cfg.Name == nil || *cfg.Name != "updated-orp" {
				return fmt.Errorf("name not updated: got %q", aws.ToString(cfg.Name))
			}
			if cfg.CookiesConfig == nil || cfg.CookiesConfig.CookieBehavior != types.OriginRequestPolicyCookieBehaviorAll {
				return fmt.Errorf("cookie behavior not updated")
			}
			if cfg.HeadersConfig == nil || cfg.HeadersConfig.HeaderBehavior != types.OriginRequestPolicyHeaderBehaviorAllViewer {
				return fmt.Errorf("header behavior not updated")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "DeleteOriginRequestPolicy", func() error {
			_, err := client.DeleteOriginRequestPolicy(ctx, &cloudfront.DeleteOriginRequestPolicyInput{
				Id: aws.String(orpID),
			})
			return err
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "GetOriginRequestPolicy_NonExistent", func() error {
			_, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
				Id: aws.String(orpID),
			})
			return AssertErrorContains(err, "NoSuchOriginRequestPolicy")
		}))
	}

	return results
}
