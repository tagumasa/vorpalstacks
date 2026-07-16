package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func cfInvalidationTests(tc *cfTestContext) []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	callerRef := tc.uniqueCallerRef("test-inv")
	distID, distETag, err := tc.createDistribution(callerRef, "For invalidation tests", "inv-test.example.com")
	if err != nil {
		results = append(results, TestResult{
			Service:  "cloudfront",
			TestName: "CreateDistribution_ForInvalidation",
			Status:   "FAIL",
			Error:    err.Error(),
		})
		return results
	}
	results = append(results, TestResult{
		Service:  "cloudfront",
		TestName: "CreateDistribution_ForInvalidation",
		Status:   "PASS",
	})

	invCallerRef := tc.uniqueCallerRef("inv-ref")
	var invID string
	results = append(results, tc.runner.RunTest("cloudfront", "CreateInvalidation_VerifyFields", func() error {
		resp, err := client.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
			DistributionId: aws.String(distID),
			InvalidationBatch: &types.InvalidationBatch{
				CallerReference: aws.String(invCallerRef),
				Paths: &types.Paths{
					Quantity: aws.Int32(2),
					Items: []string{
						"/index.html",
						"/images/*",
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Invalidation == nil {
			return fmt.Errorf("invalidation is nil")
		}
		if resp.Invalidation.Id == nil || *resp.Invalidation.Id == "" {
			return fmt.Errorf("invalidation ID is nil or empty")
		}
		if resp.Invalidation.Status == nil || *resp.Invalidation.Status == "" {
			return fmt.Errorf("invalidation status is nil or empty")
		}
		if resp.Invalidation.CreateTime == nil {
			return fmt.Errorf("invalidation createTime is nil")
		}
		if resp.Invalidation.InvalidationBatch == nil {
			return fmt.Errorf("invalidationBatch is nil")
		}
		if resp.Location == nil || *resp.Location == "" {
			return fmt.Errorf("location is nil or empty")
		}
		invID = *resp.Invalidation.Id
		return nil
	}))

	if invID != "" {
		results = append(results, tc.runner.RunTest("cloudfront", "GetInvalidation_VerifyFields", func() error {
			resp, err := client.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
				DistributionId: aws.String(distID),
				Id:             aws.String(invID),
			})
			if err != nil {
				return err
			}
			if resp.Invalidation == nil {
				return fmt.Errorf("invalidation is nil")
			}
			if resp.Invalidation.Id == nil || *resp.Invalidation.Id != invID {
				return fmt.Errorf("invalidation ID mismatch: got %q, want %q", aws.ToString(resp.Invalidation.Id), invID)
			}
			if resp.Invalidation.Status == nil {
				return fmt.Errorf("invalidation status is nil")
			}
			if resp.Invalidation.CreateTime == nil {
				return fmt.Errorf("invalidation createTime is nil")
			}
			if resp.Invalidation.InvalidationBatch == nil {
				return fmt.Errorf("invalidationBatch is nil")
			}
			return nil
		}))
	}

	results = append(results, tc.runner.RunTest("cloudfront", "ListInvalidations_ContainsCreated", func() error {
		resp, err := client.ListInvalidations(ctx, &cloudfront.ListInvalidationsInput{
			DistributionId: aws.String(distID),
		})
		if err != nil {
			return err
		}
		if resp.InvalidationList == nil {
			return fmt.Errorf("invalidation list is nil")
		}
		if resp.InvalidationList.Quantity == nil || *resp.InvalidationList.Quantity < 1 {
			return fmt.Errorf("expected at least 1 invalidation, got 0")
		}
		found := false
		for _, inv := range resp.InvalidationList.Items {
			if inv.Id != nil && *inv.Id == invID {
				found = true
				if inv.Status == nil || *inv.Status == "" {
					return fmt.Errorf("invalidation status is nil or empty in list")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created invalidation %q not found in list", invID)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "GetInvalidation_NonExistent", func() error {
		_, err := client.GetInvalidation(ctx, &cloudfront.GetInvalidationInput{
			DistributionId: aws.String(distID),
			Id:             aws.String("NONEXISTENT-ID"),
		})
		return AssertErrorContains(err, "NoSuchInvalidation")
	}))

	results = append(results, tc.runner.RunTest("cloudfront", "Cleanup_InvalidationDist", func() error {
		return tc.disableAndDeleteDistribution(distID, distETag)
	}))

	return results
}
