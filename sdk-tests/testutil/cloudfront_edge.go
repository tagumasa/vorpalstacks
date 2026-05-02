package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfEdgeTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("cloudfront", "ListKeyGroups_VerifyFields", func() error {
		resp, err := client.ListKeyGroups(ctx, &cloudfront.ListKeyGroupsInput{
			MaxItems: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.KeyGroupList == nil {
			return fmt.Errorf("key group list is nil")
		}
		if resp.KeyGroupList.MaxItems == nil {
			return fmt.Errorf("maxItems is nil")
		}
		if resp.KeyGroupList.Quantity == nil {
			return fmt.Errorf("quantity is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetDistribution_StaleETag", func() error {
		_, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
			Id: aws.String("nonexistent-dist-id"),
		})
		return AssertErrorContains(err, "NoSuchDistribution")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetInvalidation_NonExistentDist", func() error {
		_, err := client.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
			DistributionId: aws.String("nonexistent-dist-id"),
			Id:             aws.String("any-inv-id"),
		})
		return AssertErrorContains(err, "NoSuchDistribution")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "CreateInvalidation_NonExistentDist", func() error {
		_, err := client.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
			DistributionId: aws.String("nonexistent-dist-id"),
			InvalidationBatch: &types.InvalidationBatch{
				CallerReference: aws.String("ref"),
				Paths: &types.Paths{
					Quantity: aws.Int32(1),
					Items:    []string{"/"},
				},
			},
		})
		return AssertErrorContains(err, "NoSuchDistribution")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetCachePolicy_NonExistent", func() error {
		_, err := client.GetCachePolicy(ctx, &cloudfront.GetCachePolicyInput{
			Id: aws.String("nonexistent-policy-id"),
		})
		return AssertErrorContains(err, "NoSuchCachePolicy")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetOriginRequestPolicy_NonExistent", func() error {
		_, err := client.GetOriginRequestPolicy(ctx, &cloudfront.GetOriginRequestPolicyInput{
			Id: aws.String("nonexistent-policy-id"),
		})
		return AssertErrorContains(err, "NoSuchOriginRequestPolicy")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetResponseHeadersPolicy_NonExistent", func() error {
		_, err := client.GetResponseHeadersPolicy(ctx, &cloudfront.GetResponseHeadersPolicyInput{
			Id: aws.String("nonexistent-policy-id"),
		})
		return AssertErrorContains(err, "NoSuchResponseHeadersPolicy")
	}))

	return results
}
