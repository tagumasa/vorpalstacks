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
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(ueBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(ueBus)})

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
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name:        aws.String(uvBus),
			Description: aws.String("original"),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(uvBus)})

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
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(dupBus),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(dupBus)})

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
		_, err := client.CreateEventBus(ctx, &eventbridge.CreateEventBusInput{
			Name: aws.String(lnpBus),
		})
		if err != nil {
			return fmt.Errorf("create bus: %v", err)
		}
		defer client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(lnpBus)})

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

	return results
}
