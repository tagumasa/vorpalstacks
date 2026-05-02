package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfResponseHeadersPolicyTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("cloudfront", "ListResponseHeadersPolicies_VerifyFields", func() error {
		resp, err := client.ListResponseHeadersPolicies(ctx, &cloudfront.ListResponseHeadersPoliciesInput{
			MaxItems: aws.Int32(100),
		})
		if err != nil {
			return err
		}
		if resp.ResponseHeadersPolicyList == nil {
			return fmt.Errorf("response headers policy list is nil")
		}
		if resp.ResponseHeadersPolicyList.MaxItems == nil {
			return fmt.Errorf("maxItems is nil")
		}
		return nil
	}))

	rhpName := tc.uniquePrefix("test-rhp")
	var rhpID string
	var rhpETag string
	results = append(results, tc.runner.RunTest("cloudfront", "CreateResponseHeadersPolicy_VerifyFields", func() error {
		resp, err := client.CreateResponseHeadersPolicy(ctx, &cloudfront.CreateResponseHeadersPolicyInput{
			ResponseHeadersPolicyConfig: &types.ResponseHeadersPolicyConfig{
				Name:    aws.String(rhpName),
				Comment: aws.String("Test RHP"),
				ServerTimingHeadersConfig: &types.ResponseHeadersPolicyServerTimingHeadersConfig{
					Enabled:      aws.Bool(true),
					SamplingRate: aws.Float64(0.5),
				},
				SecurityHeadersConfig: &types.ResponseHeadersPolicySecurityHeadersConfig{
					XSSProtection: &types.ResponseHeadersPolicyXSSProtection{
						Override:   aws.Bool(true),
						Protection: aws.Bool(true),
						ModeBlock:  aws.Bool(true),
					},
					ContentTypeOptions: &types.ResponseHeadersPolicyContentTypeOptions{
						Override: aws.Bool(true),
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.ResponseHeadersPolicy == nil {
			return fmt.Errorf("response headers policy is nil")
		}
		if resp.ResponseHeadersPolicy.Id == nil || *resp.ResponseHeadersPolicy.Id == "" {
			return fmt.Errorf("response headers policy ID is nil or empty")
		}
		cfg := resp.ResponseHeadersPolicy.ResponseHeadersPolicyConfig
		if cfg == nil {
			return fmt.Errorf("RHP config is nil")
		}
		if cfg.Name == nil || *cfg.Name != rhpName {
			return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(cfg.Name), rhpName)
		}
		if cfg.ServerTimingHeadersConfig == nil || cfg.ServerTimingHeadersConfig.Enabled == nil || !*cfg.ServerTimingHeadersConfig.Enabled {
			return fmt.Errorf("server timing headers config mismatch")
		}
		if resp.ETag == nil || *resp.ETag == "" {
			return fmt.Errorf("ETag is nil or empty")
		}
		rhpID = *resp.ResponseHeadersPolicy.Id
		rhpETag = *resp.ETag
		return nil
	}))

	if rhpID != "" {
		results = append(results, tc.runner.RunTest("cloudfront", "GetResponseHeadersPolicy_VerifyFields", func() error {
			resp, err := client.GetResponseHeadersPolicy(ctx, &cloudfront.GetResponseHeadersPolicyInput{
				Id: aws.String(rhpID),
			})
			if err != nil {
				return err
			}
			if resp.ResponseHeadersPolicy == nil {
				return fmt.Errorf("response headers policy is nil")
			}
			if resp.ResponseHeadersPolicy.Id == nil || *resp.ResponseHeadersPolicy.Id != rhpID {
				return fmt.Errorf("ID mismatch: got %q", aws.ToString(resp.ResponseHeadersPolicy.Id))
			}
			if resp.ResponseHeadersPolicy.LastModifiedTime == nil {
				return fmt.Errorf("lastModifiedTime is nil")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "GetResponseHeadersPolicyConfig_VerifyFields", func() error {
			resp, err := client.GetResponseHeadersPolicyConfig(ctx, &cloudfront.GetResponseHeadersPolicyConfigInput{
				Id: aws.String(rhpID),
			})
			if err != nil {
				return err
			}
			if resp.ResponseHeadersPolicyConfig == nil {
				return fmt.Errorf("response headers policy config is nil")
			}
			if resp.ResponseHeadersPolicyConfig.Name == nil || *resp.ResponseHeadersPolicyConfig.Name != rhpName {
				return fmt.Errorf("name mismatch: got %q, want %q", aws.ToString(resp.ResponseHeadersPolicyConfig.Name), rhpName)
			}
			if resp.ETag == nil || *resp.ETag == "" {
				return fmt.Errorf("ETag is nil or empty")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "UpdateResponseHeadersPolicy_VerifyFields", func() error {
			resp, err := client.UpdateResponseHeadersPolicy(ctx, &cloudfront.UpdateResponseHeadersPolicyInput{
				Id:      aws.String(rhpID),
				IfMatch: aws.String(rhpETag),
				ResponseHeadersPolicyConfig: &types.ResponseHeadersPolicyConfig{
					Name:    aws.String(rhpName + "-updated"),
					Comment: aws.String("Updated RHP"),
					CustomHeadersConfig: &types.ResponseHeadersPolicyCustomHeadersConfig{
						Quantity: aws.Int32(1),
						Items: []types.ResponseHeadersPolicyCustomHeader{
							{
								Header:   aws.String("X-Custom-Header"),
								Value:    aws.String("custom-value"),
								Override: aws.Bool(true),
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}
			if resp.ResponseHeadersPolicy == nil {
				return fmt.Errorf("updated RHP is nil")
			}
			cfg := resp.ResponseHeadersPolicy.ResponseHeadersPolicyConfig
			if cfg == nil {
				return fmt.Errorf("RHP config is nil")
			}
			if cfg.Name == nil || *cfg.Name != rhpName+"-updated" {
				return fmt.Errorf("name not updated: got %q", aws.ToString(cfg.Name))
			}
			if cfg.CustomHeadersConfig == nil || cfg.CustomHeadersConfig.Quantity == nil || *cfg.CustomHeadersConfig.Quantity != 1 {
				return fmt.Errorf("custom headers config not set")
			}
			return nil
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "DeleteResponseHeadersPolicy", func() error {
			_, err := client.DeleteResponseHeadersPolicy(ctx, &cloudfront.DeleteResponseHeadersPolicyInput{
				Id: aws.String(rhpID),
			})
			return err
		}))

		results = append(results, tc.runner.RunTest("cloudfront", "GetResponseHeadersPolicy_NonExistent", func() error {
			_, err := client.GetResponseHeadersPolicy(ctx, &cloudfront.GetResponseHeadersPolicyInput{
				Id: aws.String(rhpID),
			})
			return AssertErrorContains(err, "NoSuchResponseHeadersPolicy")
		}))
	}

	return results
}
