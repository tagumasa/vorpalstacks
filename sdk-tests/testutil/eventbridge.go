package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunEventBridgeTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "events",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := eventbridge.NewFromConfig(cfg)
	ctx := context.Background()

	busName := fmt.Sprintf("TestBus-%d", time.Now().UnixNano())
	ruleName := fmt.Sprintf("TestRule-%d", time.Now().UnixNano())
	targetID := fmt.Sprintf("TestTarget-%d", time.Now().UnixNano())

	results = append(results, r.RunTest("events", "CreateEventBus", func() error {
		resp, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.EventBusArn == nil || *resp.EventBusArn == "" {
			return fmt.Errorf("event bus ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeEventBus", func() error {
		resp, err := client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp.Name == nil || *resp.Name != busName {
			return fmt.Errorf("event bus name mismatch, got %v", resp.Name)
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("event bus ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.runEventBridgeBusTests(ctx, client, busName)...)
	results = append(results, r.runEventBridgeRuleTests(ctx, client, busName, ruleName)...)
	results = append(results, r.runEventBridgeTargetTests(ctx, client, busName, ruleName, targetID)...)
	results = append(results, r.runEventBridgeEventTests(ctx, client, busName)...)

	results = append(results, r.RunTest("events", "RemoveTargets", func() error {
		resp, err := client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
			Rule:         aws.String(ruleName),
			EventBusName: aws.String(busName),
			Ids:          []string{targetID},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.FailedEntryCount > 0 {
			return fmt.Errorf("failed entries: %d", resp.FailedEntryCount)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DisableRule", func() error {
		resp, err := client.DisableRule(ctx, &eventbridge.DisableRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		desc, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return fmt.Errorf("describe after disable: %v", err)
		}
		if desc.State != types.RuleStateDisabled {
			return fmt.Errorf("expected DISABLED state, got %v", desc.State)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "EnableRule", func() error {
		resp, err := client.EnableRule(ctx, &eventbridge.EnableRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		desc, err := client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return fmt.Errorf("describe after enable: %v", err)
		}
		if desc.State != types.RuleStateEnabled {
			return fmt.Errorf("expected ENABLED state, got %v", desc.State)
		}
		return nil
	}))

	results = append(results, r.runEventBridgeTagTests(ctx, client, busName, ruleName)...)

	results = append(results, r.RunTest("events", "DeleteRule", func() error {
		resp, err := client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		_, err = client.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
		})
		if err == nil {
			return fmt.Errorf("expected error describing deleted rule")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteEventBus", func() error {
		_, err := client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{
			Name: aws.String(busName),
		})
		if err != nil {
			return err
		}
		_, err = client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String(busName),
		})
		if err == nil {
			return fmt.Errorf("expected error describing deleted event bus")
		}
		return nil
	}))

	results = append(results, r.runEventBridgeArchiveTests(ctx, client)...)
	results = append(results, r.runEventBridgeConnectionTests(ctx, client)...)
	results = append(results, r.runEventBridgeEdgeTests(ctx, client)...)

	return results
}
