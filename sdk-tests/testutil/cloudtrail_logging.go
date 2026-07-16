package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailLoggingTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "StartLogging", func() error {
		name := fmt.Sprintf("startlog-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("startlog-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{
			Name: aws.String(name),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get status after start: %v", err)
		}
		if status.IsLogging == nil || !*status.IsLogging {
			return fmt.Errorf("expected IsLogging=true after StartLogging")
		}
		if status.StartLoggingTime == nil {
			return fmt.Errorf("expected StartLoggingTime to be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StopLogging", func() error {
		name := fmt.Sprintf("stoplog-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("stoplog-bucket"),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("start logging: %v", err)
		}

		resp, err := tc.client.StopLogging(tc.ctx, &cloudtrail.StopLoggingInput{
			Name: aws.String(name),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get status after stop: %v", err)
		}
		if status.IsLogging != nil && *status.IsLogging {
			return fmt.Errorf("expected IsLogging=false after StopLogging")
		}
		if status.StopLoggingTime == nil {
			return fmt.Errorf("expected StopLoggingTime to be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus", func() error {
		name := fmt.Sprintf("status-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("status-bucket"),
		})
		if err != nil {
			return err
		}

		resp, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{
			Name: aws.String(name),
		})
		if err != nil {
			return err
		}
		if resp.IsLogging != nil && *resp.IsLogging {
			return fmt.Errorf("expected IsLogging=false for new trail")
		}
		if resp.LatestDeliveryTime != nil {
			return fmt.Errorf("expected nil LatestDeliveryTime when not logging")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus_AfterStart", func() error {
		name := fmt.Sprintf("status-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("status-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("start logging: %v", err)
		}

		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get status: %v", err)
		}
		if status.IsLogging == nil || !*status.IsLogging {
			return fmt.Errorf("expected IsLogging=true after StartLogging")
		}
		if status.StartLoggingTime == nil {
			return fmt.Errorf("expected StartLoggingTime to be set")
		}
		if status.LatestDeliveryTime == nil {
			return fmt.Errorf("expected LatestDeliveryTime when logging is active")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus_AfterStop", func() error {
		name := fmt.Sprintf("stopstat-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("stopstat-bucket"),
		})
		if err != nil {
			return err
		}

		_, err = tc.client.StartLogging(tc.ctx, &cloudtrail.StartLoggingInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("start logging: %v", err)
		}
		_, err = tc.client.StopLogging(tc.ctx, &cloudtrail.StopLoggingInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("stop logging: %v", err)
		}

		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get status: %v", err)
		}
		if status.IsLogging != nil && *status.IsLogging {
			return fmt.Errorf("expected IsLogging=false after StopLogging")
		}
		if status.StopLoggingTime == nil {
			return fmt.Errorf("expected StopLoggingTime to be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "GetTrailStatus_DefaultIsNotLogging", func() error {
		name := fmt.Sprintf("islog-%d", time.Now().UnixNano())
		defer tc.client.DeleteTrail(tc.ctx, &cloudtrail.DeleteTrailInput{Name: aws.String(name)})

		_, err := tc.client.CreateTrail(tc.ctx, &cloudtrail.CreateTrailInput{
			Name:         aws.String(name),
			S3BucketName: aws.String("islog-bucket"),
		})
		if err != nil {
			return err
		}

		status, err := tc.client.GetTrailStatus(tc.ctx, &cloudtrail.GetTrailStatusInput{Name: aws.String(name)})
		if err != nil {
			return fmt.Errorf("get status: %v", err)
		}
		if status.IsLogging != nil && *status.IsLogging {
			return fmt.Errorf("expected IsLogging=false for newly created trail")
		}
		if status.LatestDeliveryTime != nil {
			return fmt.Errorf("expected nil LatestDeliveryTime when not logging")
		}
		return nil
	}))

	return results
}
