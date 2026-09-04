package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func (r *TestRunner) runEventBridgeEdgeTests(ctx context.Context, client *eventbridge.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("events", "DescribeEventBus_NonExistent", func() error {
		_, err := client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String("nonexistent-bus-xyz-12345"),
		})
		return expectEventBridgeResourceNotFound(err)
	}))

	results = append(results, r.RunTest("events", "DeleteEventBus_NonExistent", func() error {
		_, err := client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{
			Name: aws.String("nonexistent-bus-xyz-12345"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent event bus")
		}
		return expectAWSErrorCode(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("events", "DescribeRule_NonExistent", func() error {
		_, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name: aws.String("nonexistent-rule-xyz-12345"),
		})
		return expectEventBridgeResourceNotFound(err)
	}))

	results = append(results, r.RunTest("events", "DeleteRule_NonExistent", func() error {
		_, err := client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{
			Name: aws.String("nonexistent-rule-xyz-12345"),
		})
		return expectEventBridgeResourceNotFound(err)
	}))

	results = append(results, r.RunTest("events", "DescribeConnection_NonExistent", func() error {
		_, err := client.DescribeConnection(ctx, &eventbridge.DescribeConnectionInput{
			Name: aws.String("nonexistent-conn-xyz-12345"),
		})
		return expectEventBridgeResourceNotFound(err)
	}))

	results = append(results, r.RunTest("events", "DeleteConnection_NonExistent", func() error {
		_, err := client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{
			Name: aws.String("nonexistent-conn-xyz-12345"),
		})
		return expectEventBridgeResourceNotFound(err)
	}))

	results = append(results, r.RunTest("events", "DescribeApiDestination_NonExistent", func() error {
		_, err := client.DescribeApiDestination(ctx, &eventbridge.DescribeApiDestinationInput{
			Name: aws.String("nonexistent-apidest-xyz-12345"),
		})
		return expectEventBridgeResourceNotFound(err)
	}))

	results = append(results, r.RunTest("events", "DeleteApiDestination_NonExistent", func() error {
		_, err := client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{
			Name: aws.String("nonexistent-apidest-xyz-12345"),
		})
		return expectEventBridgeResourceNotFound(err)
	}))

	return results
}

func (r *TestRunner) runEventBridgeTagTests(ctx context.Context, client *eventbridge.Client, busName, ruleName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("events", "TagResource", func() error {
		ruleARN := fmt.Sprintf("arn:aws:events:%s:%s:rule/%s/%s", r.region, r.accountID, busName, ruleName)
		resp, err := client.TagResource(ctx, &eventbridge.TagResourceInput{
			ResourceARN: aws.String(ruleARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		listResp, err := client.ListTagsForResource(ctx, &eventbridge.ListTagsForResourceInput{
			ResourceARN: aws.String(ruleARN),
		})
		if err != nil {
			return fmt.Errorf("list tags after tagging: %v", err)
		}
		found := false
		for _, t := range listResp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				found = true
				if t.Value == nil || *t.Value != "test" {
					return fmt.Errorf("tag value mismatch, got %v", t.Value)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("tag Environment not found after TagResource")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListTagsForResource", func() error {
		ruleARN := fmt.Sprintf("arn:aws:events:%s:%s:rule/%s/%s", r.region, r.accountID, busName, ruleName)
		resp, err := client.ListTagsForResource(ctx, &eventbridge.ListTagsForResourceInput{
			ResourceARN: aws.String(ruleARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil || len(resp.Tags) == 0 {
			return fmt.Errorf("expected at least 1 tag")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UntagResource", func() error {
		ruleARN := fmt.Sprintf("arn:aws:events:%s:%s:rule/%s/%s", r.region, r.accountID, busName, ruleName)
		resp, err := client.UntagResource(ctx, &eventbridge.UntagResourceInput{
			ResourceARN: aws.String(ruleARN),
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		listResp, err := client.ListTagsForResource(ctx, &eventbridge.ListTagsForResourceInput{
			ResourceARN: aws.String(ruleARN),
		})
		if err != nil {
			return fmt.Errorf("list tags after untagging: %v", err)
		}
		for _, t := range listResp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				return fmt.Errorf("tag Environment still present after UntagResource")
			}
		}
		return nil
	}))

	return results
}
