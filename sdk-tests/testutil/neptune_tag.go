package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
)

func (r *TestRunner) runNeptuneTagTests(tc *neptuneContext) []TestResult {
	var results []TestResult
	clusterARN := fmt.Sprintf("arn:aws:rds:%s:%s:cluster:%s", tc.region, tc.accountID, tc.clusterID)

	results = append(results, r.RunTest("neptune", "AddTagsToResource", func() error {
		_, err := tc.client.AddTagsToResource(tc.ctx, &neptune.AddTagsToResourceInput{
			ResourceName: aws.String(clusterARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Owner"), Value: aws.String("sdk-test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ListTagsForResource", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &neptune.ListTagsForResourceInput{
			ResourceName: aws.String(clusterARN),
		})
		if err != nil {
			return err
		}
		if resp.TagList == nil {
			return fmt.Errorf("expected TagList to be non-nil")
		}
		if len(resp.TagList) < 2 {
			return fmt.Errorf("expected at least 2 tags, got %d", len(resp.TagList))
		}
		foundEnv := false
		foundOwner := false
		for _, t := range resp.TagList {
			if t.Key != nil && *t.Key == "Environment" && t.Value != nil && *t.Value == "test" {
				foundEnv = true
			}
			if t.Key != nil && *t.Key == "Owner" && t.Value != nil && *t.Value == "sdk-test" {
				foundOwner = true
			}
		}
		if !foundEnv {
			return fmt.Errorf("expected tag Environment=test not found in TagList")
		}
		if !foundOwner {
			return fmt.Errorf("expected tag Owner=sdk-test not found in TagList")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptune", "RemoveTagsFromResource", func() error {
		_, err := tc.client.RemoveTagsFromResource(tc.ctx, &neptune.RemoveTagsFromResourceInput{
			ResourceName: aws.String(clusterARN),
			TagKeys:      []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("neptune", "ListTagsForResource_AfterRemove", func() error {
		resp, err := tc.client.ListTagsForResource(tc.ctx, &neptune.ListTagsForResourceInput{
			ResourceName: aws.String(clusterARN),
		})
		if err != nil {
			return err
		}
		for _, t := range resp.TagList {
			if t.Key != nil && *t.Key == "Environment" {
				return fmt.Errorf("expected tag Environment to be removed but still present")
			}
		}
		foundOwner := false
		for _, t := range resp.TagList {
			if t.Key != nil && *t.Key == "Owner" {
				foundOwner = true
			}
		}
		if !foundOwner {
			return fmt.Errorf("expected tag Owner to still be present after removing Environment")
		}
		return nil
	}))

	// Tag operations against a cluster that does not exist fail with the
	// resource-specific DBClusterNotFoundFault, as the service model
	// specifies for these operations.
	results = append(results, r.RunTest("neptune", "AddTagsToResource_NonExistentCluster", func() error {
		arn := fmt.Sprintf("arn:aws:rds:%s:%s:cluster:no-such-cluster-%d",
			tc.region, tc.accountID, time.Now().UnixNano())
		_, err := tc.client.AddTagsToResource(tc.ctx, &neptune.AddTagsToResourceInput{
			ResourceName: aws.String(arn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err := AssertErrorContains(err, "DBClusterNotFoundFault"); err != nil {
			return fmt.Errorf("AddTagsToResource: %v", err)
		}
		_, err = tc.client.RemoveTagsFromResource(tc.ctx, &neptune.RemoveTagsFromResourceInput{
			ResourceName: aws.String(arn),
			TagKeys:      []string{"Environment"},
		})
		if err := AssertErrorContains(err, "DBClusterNotFoundFault"); err != nil {
			return fmt.Errorf("RemoveTagsFromResource: %v", err)
		}
		_, err = tc.client.ListTagsForResource(tc.ctx, &neptune.ListTagsForResourceInput{
			ResourceName: aws.String(arn),
		})
		if err := AssertErrorContains(err, "DBClusterNotFoundFault"); err != nil {
			return fmt.Errorf("ListTagsForResource: %v", err)
		}
		return nil
	}))

	return results
}
