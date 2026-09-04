package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func (r *TestRunner) runEventBridgeBusTests(ctx context.Context, client *eventbridge.Client, busName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("events", "ListEventBuses", func() error {
		resp, err := client.ListEventBuses(ctx, &eventbridge.ListEventBusesInput{})
		if err != nil {
			return err
		}
		if resp.EventBuses == nil {
			return fmt.Errorf("event buses list is nil")
		}
		found := false
		for _, b := range resp.EventBuses {
			if b.Name != nil && *b.Name == busName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected bus %s in list", busName)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateEventBus", func() error {
		ueBus := fmt.Sprintf("UeBus-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, ueBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		resp, err := client.UpdateEventBus(ctx, &eventbridge.UpdateEventBusInput{
			Name:        aws.String(ueBus),
			Description: aws.String("updated description"),
		})
		if err != nil {
			return fmt.Errorf("update event bus: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Description == nil || *resp.Description != "updated description" {
			return fmt.Errorf("description mismatch, got %v", resp.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateEventBus_VerifyDescription", func() error {
		uvBus := fmt.Sprintf("UvBus-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, uvBus, func(input *eventbridge.CreateEventBusInput) {
			input.Description = aws.String("original")
		})
		if err != nil {
			return err
		}
		defer cleanupBus()

		_, err = client.UpdateEventBus(ctx, &eventbridge.UpdateEventBusInput{
			Name:        aws.String(uvBus),
			Description: aws.String("updated"),
		})
		if err != nil {
			return fmt.Errorf("update: %v", err)
		}

		desc, err := client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String(uvBus),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if desc.Description == nil || *desc.Description != "updated" {
			return fmt.Errorf("description not updated, got %v", desc.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "CreateEventBus_DuplicateName", func() error {
		dupBus := fmt.Sprintf("DupBus-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, dupBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		_, err = client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(dupBus),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate event bus name")
		}
		var riu *types.ResourceAlreadyExistsException
		if !errors.As(err, &riu) {
			return fmt.Errorf("expected ResourceAlreadyExistsException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListEventBuses_NamePrefix", func() error {
		lnpBus := fmt.Sprintf("LnpPrefixBus-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, lnpBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		resp, err := client.ListEventBuses(ctx, &eventbridge.ListEventBusesInput{
			NamePrefix: aws.String("LnpPrefixBus"),
		})
		if err != nil {
			return fmt.Errorf("list event buses with prefix: %v", err)
		}
		if resp == nil || resp.EventBuses == nil {
			return fmt.Errorf("response or event buses is nil")
		}
		found := false
		for _, bus := range resp.EventBuses {
			if bus.Name != nil && *bus.Name == lnpBus {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected bus %s in filtered list", lnpBus)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "EventBus_LogConfigRoundTrip", func() error {
		lcBus := fmt.Sprintf("LcBus-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, lcBus, func(input *eventbridge.CreateEventBusInput) {
			input.LogConfig = &types.LogConfig{
				IncludeDetail: types.IncludeDetailNone,
				Level:         types.LevelInfo,
			}
		})
		if err != nil {
			return err
		}
		defer cleanupBus()

		desc, err := client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String(lcBus),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if desc.LogConfig == nil {
			return fmt.Errorf("log config is nil after create")
		}
		if desc.LogConfig.IncludeDetail != types.IncludeDetailNone || desc.LogConfig.Level != types.LevelInfo {
			return fmt.Errorf("create echo mismatch: got %s/%s, want NONE/INFO",
				desc.LogConfig.IncludeDetail, desc.LogConfig.Level)
		}

		resp, err := client.UpdateEventBus(ctx, &eventbridge.UpdateEventBusInput{
			Name: aws.String(lcBus),
			LogConfig: &types.LogConfig{
				IncludeDetail: types.IncludeDetailFull,
				Level:         types.LevelTrace,
			},
		})
		if err != nil {
			return fmt.Errorf("update log config: %v", err)
		}
		if resp.LogConfig == nil || resp.LogConfig.IncludeDetail != types.IncludeDetailFull ||
			resp.LogConfig.Level != types.LevelTrace {
			return fmt.Errorf("update echo mismatch: got %+v, want FULL/TRACE", resp.LogConfig)
		}

		desc, err = client.DescribeEventBus(ctx, &eventbridge.DescribeEventBusInput{
			Name: aws.String(lcBus),
		})
		if err != nil {
			return fmt.Errorf("describe after update: %v", err)
		}
		if desc.LogConfig == nil || desc.LogConfig.IncludeDetail != types.IncludeDetailFull ||
			desc.LogConfig.Level != types.LevelTrace {
			return fmt.Errorf("persisted log config mismatch: got %+v, want FULL/TRACE", desc.LogConfig)
		}
		return nil
	}))

	// A non-empty out-of-enum value must be rejected on both the create and
	// update paths. The typed SDK omits empty enum members on the wire, so
	// the explicit-empty-string rejection is pinned by the Core-layer unit
	// tests of the service package instead.
	results = append(results, r.RunTest("events", "EventBus_LogConfigEnumRejected", func() error {
		leBus := fmt.Sprintf("LeBus-%d", time.Now().UnixNano())
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(leBus),
			LogConfig: &types.LogConfig{
				IncludeDetail: types.IncludeDetail("BOGUS"),
			},
		})
		if err == nil {
			defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(leBus)})
			return fmt.Errorf("expected error for out-of-enum IncludeDetail on create")
		}
		if codeErr := expectAWSErrorCode(err, "ValidationException"); codeErr != nil {
			return fmt.Errorf("create rejection: %v", codeErr)
		}

		cleanupBus, err := createEventBridgeTestBus(ctx, client, leBus)
		if err != nil {
			return err
		}
		defer cleanupBus()

		_, err = client.UpdateEventBus(ctx, &eventbridge.UpdateEventBusInput{
			Name: aws.String(leBus),
			LogConfig: &types.LogConfig{
				Level: types.Level("BOGUS"),
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for out-of-enum Level on update")
		}
		if codeErr := expectAWSErrorCode(err, "ValidationException"); codeErr != nil {
			return fmt.Errorf("update rejection: %v", codeErr)
		}
		return nil
	}))

	return results
}
