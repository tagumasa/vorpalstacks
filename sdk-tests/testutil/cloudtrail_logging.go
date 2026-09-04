package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
)

func (r *TestRunner) runCloudTrailLoggingTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cloudtrail", "StartLogging", func() error {
		name := tc.uniqueName("startlog")
		defer tc.deleteTrail(name)

		_, err := tc.createTrail(name, "startlog-bucket")
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
		if status.LatestDeliveryTime == nil {
			return fmt.Errorf("expected LatestDeliveryTime when logging is active")
		}
		return nil
	}))

	results = append(results, r.RunTest("cloudtrail", "StopLogging", func() error {
		name := tc.uniqueName("stoplog")
		defer tc.deleteTrail(name)

		_, err := tc.createTrail(name, "stoplog-bucket")
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
		name := tc.uniqueName("status")
		defer tc.deleteTrail(name)

		_, err := tc.createTrail(name, "status-bucket")
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

	return results
}
