package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

func (r *TestRunner) runCloudTrailEdgeTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "GetTrail_NonExistent", func() error {
		_, err := tc.client.GetTrail(tc.ctx, &cloudtrail.GetTrailInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteTrail_NonExistent", func() error {
		_, err := tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StartLogging_NonExistent", func() error {
		_, err := tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails_NonExistent", func() error {
		resp, err := tc.client.DescribeTrails(tc.ctx, &cloudtrail.DescribeTrailsInput{
			TrailNameList: []string{"nonexistent-trail-xyz"},
		})
		if err != nil {
			return fmt.Errorf("DescribeTrails should not error for non-existent trail, got: %v", err)
		}
		if len(resp.TrailList) != 0 {
			return fmt.Errorf("expected empty trail list for non-existent trail, got %d", len(resp.TrailList))
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StopLogging_NonExistent", func() error {
		_, err := tc.client.StopLogging(tc.ctx, &cloudtrail.StopLoggingInput{
			Name: aws.String("nonexistent-stop-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus_NonExistent", func() error {
		_, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{
			Name: aws.String("nonexistent-status-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_NonExistent", func() error {
		_, err := tc.client.UpdateTrail(tc.ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String("nonexistent-update-xyz"),
			S3BucketName: aws.String("some-bucket"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetEventSelectors_NonExistent", func() error {
		_, err := tc.client.GetEventSelectors(tc.ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String("nonexistent-es-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors_NonExistent", func() error {
		_, err := tc.client.PutEventSelectors(tc.ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName: aws.String("nonexistent-es-xyz"),
			EventSelectors: []types.EventSelector{
				{ReadWriteType: types.ReadWriteTypeAll},
			},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutInsightSelectors_NonExistent", func() error {
		_, err := tc.client.PutInsightSelectors(tc.ctx, &cloudtrail.PutInsightSelectorsInput{
			TrailName: aws.String("nonexistent-is-xyz"),
			InsightSelectors: []types.InsightSelector{
				{InsightType: types.InsightTypeApiCallRateInsight},
			},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetInsightSelectors_NonExistent", func() error {
		_, err := tc.client.GetInsightSelectors(tc.ctx, &cloudtrail.GetInsightSelectorsInput{
			TrailName: aws.String("nonexistent-is-xyz"),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "AddTags_NonExistent", func() error {
		_, err := tc.client.AddTags(tc.ctx, &cloudtrail.AddTagsInput{
			ResourceId: aws.String("arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-tag-xyz"),
			TagsList:   []types.Tag{{Key: aws.String("K"), Value: aws.String("V")}},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "RemoveTags_NonExistent", func() error {
		_, err := tc.client.RemoveTags(tc.ctx, &cloudtrail.RemoveTagsInput{
			ResourceId: aws.String("arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-rm-xyz"),
			TagsList:   []types.Tag{{Key: aws.String("K")}},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "ListTags_NonExistent", func() error {
		_, err := tc.client.ListTags(tc.ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{"arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-lt-xyz"},
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetResourcePolicy_NonExistentTrail", func() error {
		fakeARN := "arn:aws:cloudtrail:us-east-1:123456789012:trail/nonexistent-grp-xyz"
		_, err := tc.client.GetResourcePolicy(tc.ctx, &cloudtrail.GetResourcePolicyInput{
			ResourceArn: aws.String(fakeARN),
		})
		if err := AssertErrorContains(err, "TrailNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
