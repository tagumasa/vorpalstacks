package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunCloudTrailTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cloudtrail",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cloudtrail.NewFromConfig(cfg)
	ctx := context.Background()

	trailName := fmt.Sprintf("test-trail-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("cloudtrail", "ListTrails", func() error {
		_, err := client.ListTrails(ctx, &cloudtrail.ListTrailsInput{})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "CreateTrail", func() error {
		_, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:                       aws.String(trailName),
			S3BucketName:               aws.String("test-bucket"),
			IncludeGlobalServiceEvents: aws.Bool(true),
			IsMultiRegionTrail:         aws.Bool(true),
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrail", func() error {
		_, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
			Name: aws.String(trailName),
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails", func() error {
		_, err := client.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{
			TrailNameList: []string{trailName},
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "StartLogging", func() error {
		_, err := client.StartLogging(ctx, &cloudtrail.StartLoggingInput{
			Name: aws.String(trailName),
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "StopLogging", func() error {
		_, err := client.StopLogging(ctx, &cloudtrail.StopLoggingInput{
			Name: aws.String(trailName),
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus", func() error {
		_, err := client.GetTrailStatus(ctx, &cloudtrail.GetTrailStatusInput{
			Name: aws.String(trailName),
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail", func() error {
		_, err := client.UpdateTrail(ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String(trailName),
			S3BucketName: aws.String("updated-bucket"),
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "GetEventSelectors", func() error {
		_, err := client.GetEventSelectors(ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String(trailName),
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors", func() error {
		_, err := client.PutEventSelectors(ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName: aws.String(trailName),
			EventSelectors: []types.EventSelector{
				{
					ReadWriteType:           types.ReadWriteTypeAll,
					IncludeManagementEvents: aws.Bool(true),
				},
			},
		})
		return err
	}))

	var trailARN string
	getTrailResp, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
		Name: aws.String(trailName),
	})
	if err == nil && getTrailResp.Trail != nil && getTrailResp.Trail.TrailARN != nil {
		trailARN = *getTrailResp.Trail.TrailARN
	}

	results = append(results, r.RunTest("cloudtrail", "AddTags", func() error {
		resourceID := trailARN
		if resourceID == "" {
			resourceID = trailName
		}
		_, err := client.AddTags(ctx, &cloudtrail.AddTagsInput{
			ResourceId: aws.String(resourceID),
			TagsList: []types.Tag{
				{
					Key:   aws.String("Environment"),
					Value: aws.String("test"),
				},
				{
					Key:   aws.String("Owner"),
					Value: aws.String("test-user"),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "ListTags", func() error {
		resourceID := trailARN
		if resourceID == "" {
			resourceID = trailName
		}
		_, err := client.ListTags(ctx, &cloudtrail.ListTagsInput{
			ResourceIdList: []string{resourceID},
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "RemoveTags", func() error {
		resourceID := trailARN
		if resourceID == "" {
			resourceID = trailName
		}
		_, err := client.RemoveTags(ctx, &cloudtrail.RemoveTagsInput{
			ResourceId: aws.String(resourceID),
			TagsList: []types.Tag{
				{
					Key: aws.String("Environment"),
				},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "LookupEvents", func() error {
		_, err := client.LookupEvents(ctx, &cloudtrail.LookupEventsInput{
			MaxResults: aws.Int32(10),
		})
		return err
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteTrail", func() error {
		_, err := client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{
			Name: aws.String(trailName),
		})
		return err
	}))

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("cloudtrail", "GetTrail_NonExistent", func() error {
		_, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent trail")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DeleteTrail_NonExistent", func() error {
		_, err := client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent trail")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StartLogging_NonExistent", func() error {
		_, err := client.StartLogging(ctx, &cloudtrail.StartLoggingInput{
			Name: aws.String("nonexistent-trail-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent trail")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "DescribeTrails_NonExistent", func() error {
		resp, err := client.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{
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

	var verifyTrailName string
	results = append(results, r.RunTest("cloudtrail", "CreateTrail_ContentVerify", func() error {
		verifyTrailName = fmt.Sprintf("verify-trail-%d", time.Now().UnixNano())
		resp, err := client.CreateTrail(ctx, &cloudtrail.CreateTrailInput{
			Name:                       aws.String(verifyTrailName),
			S3BucketName:               aws.String("verify-bucket"),
			IncludeGlobalServiceEvents: aws.Bool(true),
			IsMultiRegionTrail:         aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		if resp.Name == nil || *resp.Name != verifyTrailName {
			return fmt.Errorf("trail name mismatch")
		}
		if resp.S3BucketName == nil || *resp.S3BucketName != "verify-bucket" {
			return fmt.Errorf("S3 bucket name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "UpdateTrail_VerifyChange", func() error {
		if verifyTrailName == "" {
			return fmt.Errorf("trail name not available")
		}
		_, err := client.UpdateTrail(ctx, &cloudtrail.UpdateTrailInput{
			Name:         aws.String(verifyTrailName),
			S3BucketName: aws.String("updated-verify-bucket"),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}
		resp, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{
			Name: aws.String(verifyTrailName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Trail == nil || resp.Trail.S3BucketName == nil || *resp.Trail.S3BucketName != "updated-verify-bucket" {
			return fmt.Errorf("S3 bucket name not updated, got %v", resp.Trail.S3BucketName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "PutEventSelectors_VerifyContent", func() error {
		if verifyTrailName == "" {
			return fmt.Errorf("trail name not available")
		}
		selectors := []types.EventSelector{
			{
				ReadWriteType:           types.ReadWriteTypeReadOnly,
				IncludeManagementEvents: aws.Bool(false),
				DataResources: []types.DataResource{
					{
						Type:   aws.String("AWS::S3::Object"),
						Values: []string{"arn:aws:s3:::"},
					},
				},
			},
		}
		_, err := client.PutEventSelectors(ctx, &cloudtrail.PutEventSelectorsInput{
			TrailName:      aws.String(verifyTrailName),
			EventSelectors: selectors,
		})
		if err != nil {
			return fmt.Errorf("put event selectors: %v", err)
		}
		resp, err := client.GetEventSelectors(ctx, &cloudtrail.GetEventSelectorsInput{
			TrailName: aws.String(verifyTrailName),
		})
		if err != nil {
			return fmt.Errorf("get event selectors: %v", err)
		}
		if len(resp.EventSelectors) != 1 {
			return fmt.Errorf("expected 1 event selector, got %d", len(resp.EventSelectors))
		}
		if resp.EventSelectors[0].ReadWriteType != types.ReadWriteTypeReadOnly {
			return fmt.Errorf("ReadWriteType mismatch, got %v", resp.EventSelectors[0].ReadWriteType)
		}
		return nil
	}))

	if verifyTrailName != "" {
		client.DeleteTrail(ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(verifyTrailName)})
	}

	return results
}
