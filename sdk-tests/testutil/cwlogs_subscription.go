package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func (tc *cwlogsTestCtx) subscriptionTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "PutSubscriptionFilter_VerifyFields", func() error {
		sfName := tc.uniquePrefix("SFGroup")
		if err := tc.createLogGroup(sfName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(sfName)

		roleARN, cleanup, err := tc.createSubscriptionRole()
		if err != nil {
			return err
		}
		defer cleanup()

		destARN := "arn:aws:lambda:us-east-1:000000000000:function:test-func"
		if err := tc.putSubscriptionFilter(sfName, "TestSub", "ERROR", destARN, roleARN); err != nil {
			return fmt.Errorf("put subscription filter: %v", err)
		}

		resp, err := tc.client.DescribeSubscriptionFilters(tc.ctx, &cloudwatchlogs.DescribeSubscriptionFiltersInput{
			LogGroupName: aws.String(sfName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.SubscriptionFilters) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp.SubscriptionFilters))
		}
		f := resp.SubscriptionFilters[0]
		if f.FilterName == nil || *f.FilterName != "TestSub" {
			return fmt.Errorf("filterName mismatch: got %q", aws.ToString(f.FilterName))
		}
		if f.FilterPattern == nil || *f.FilterPattern != "ERROR" {
			return fmt.Errorf("filterPattern mismatch: got %q", aws.ToString(f.FilterPattern))
		}
		if f.DestinationArn == nil || *f.DestinationArn != destARN {
			return fmt.Errorf("destinationArn mismatch: got %q", aws.ToString(f.DestinationArn))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeSubscriptionFilters_Basic", func() error {
		dsfName := tc.uniquePrefix("DSFGroup")
		if err := tc.createLogGroup(dsfName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(dsfName)

		roleARN, cleanup, err := tc.createSubscriptionRole()
		if err != nil {
			return err
		}
		defer cleanup()

		tc.client.PutSubscriptionFilter(tc.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(dsfName),
			FilterName:     aws.String("DescSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String("arn:aws:lambda:us-east-1:000000000000:function:test"),
			RoleArn:        aws.String(roleARN),
		})

		resp, err := tc.client.DescribeSubscriptionFilters(tc.ctx, &cloudwatchlogs.DescribeSubscriptionFiltersInput{
			LogGroupName: aws.String(dsfName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.SubscriptionFilters) != 1 {
			return fmt.Errorf("expected 1 filter, got %d", len(resp.SubscriptionFilters))
		}
		if *resp.SubscriptionFilters[0].FilterName != "DescSub" {
			return fmt.Errorf("filter name mismatch: got %q", *resp.SubscriptionFilters[0].FilterName)
		}
		if *resp.SubscriptionFilters[0].DestinationArn != "arn:aws:lambda:us-east-1:000000000000:function:test" {
			return fmt.Errorf("destination arn mismatch")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DeleteSubscriptionFilter_Basic", func() error {
		delSFName := tc.uniquePrefix("DelSFGroup")
		if err := tc.createLogGroup(delSFName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(delSFName)

		roleARN, cleanup, err := tc.createSubscriptionRole()
		if err != nil {
			return err
		}
		defer cleanup()

		tc.client.PutSubscriptionFilter(tc.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
			LogGroupName:   aws.String(delSFName),
			FilterName:     aws.String("DelSub"),
			FilterPattern:  aws.String("ERROR"),
			DestinationArn: aws.String("arn:aws:lambda:us-east-1:000000000000:function:test"),
			RoleArn:        aws.String(roleARN),
		})

		_, err = tc.client.DeleteSubscriptionFilter(tc.ctx, &cloudwatchlogs.DeleteSubscriptionFilterInput{
			LogGroupName: aws.String(delSFName),
			FilterName:   aws.String("DelSub"),
		})
		if err != nil {
			return fmt.Errorf("delete subscription filter: %v", err)
		}

		resp, err := tc.client.DescribeSubscriptionFilters(tc.ctx, &cloudwatchlogs.DescribeSubscriptionFiltersInput{
			LogGroupName: aws.String(delSFName),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.SubscriptionFilters) != 0 {
			return fmt.Errorf("expected 0 filters after delete, got %d", len(resp.SubscriptionFilters))
		}
		return nil
	}))

	return results
}
