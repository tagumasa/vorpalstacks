package testutil

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func (r *TestRunner) runEventBridgeEventTests(ctx context.Context, client *eventbridge.Client, busName string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("events", "PutEvents", func() error {
		event, _ := json.Marshal(map[string]interface{}{
			"source":      "com.example.test",
			"detail-type": "TestEvent",
			"detail":      map[string]string{"message": "test"},
		})
		resp, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("TestEvent"),
					Detail:       aws.String(string(event)),
					EventBusName: aws.String(busName),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.FailedEntryCount > 0 {
			return fmt.Errorf("expected 0 failed entries, got %d", resp.FailedEntryCount)
		}
		if len(resp.Entries) != 1 {
			return fmt.Errorf("expected 1 entry, got %d", len(resp.Entries))
		}
		if resp.Entries[0].EventId == nil || *resp.Entries[0].EventId == "" {
			return fmt.Errorf("expected non-empty event ID")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutEvents_DefaultBus", func() error {
		event, _ := json.Marshal(map[string]interface{}{
			"source":      "com.test.default",
			"detail-type": "DefaultBusEvent",
			"detail":      map[string]string{"key": "value"},
		})
		resp, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					Source:     aws.String("com.test.default"),
					DetailType: aws.String("DefaultBusEvent"),
					Detail:     aws.String(string(event)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put events: %v", err)
		}
		if resp.FailedEntryCount != 0 {
			return fmt.Errorf("expected 0 failed entries, got %d", resp.FailedEntryCount)
		}
		if len(resp.Entries) != 1 {
			return fmt.Errorf("expected 1 entry result, got %d", len(resp.Entries))
		}
		if resp.Entries[0].EventId == nil || *resp.Entries[0].EventId == "" {
			return fmt.Errorf("expected non-empty event ID")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutEvents_MultipleEntries", func() error {
		event1, _ := json.Marshal(map[string]interface{}{
			"source": "com.test.multi", "detail-type": "Event1",
			"detail": map[string]string{"id": "1"},
		})
		event2, _ := json.Marshal(map[string]interface{}{
			"source": "com.test.multi", "detail-type": "Event2",
			"detail": map[string]string{"id": "2"},
		})
		resp, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					Source:     aws.String("com.test.multi"),
					DetailType: aws.String("Event1"),
					Detail:     aws.String(string(event1)),
				},
				{
					Source:     aws.String("com.test.multi"),
					DetailType: aws.String("Event2"),
					Detail:     aws.String(string(event2)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put events: %v", err)
		}
		if resp.FailedEntryCount != 0 {
			return fmt.Errorf("expected 0 failed entries, got %d", resp.FailedEntryCount)
		}
		if len(resp.Entries) != 2 {
			return fmt.Errorf("expected 2 entry results, got %d", len(resp.Entries))
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "TestEventPattern_Match", func() error {
		pattern, _ := json.Marshal(map[string]interface{}{
			"source": []string{"com.example.custom"},
		})
		event, _ := json.Marshal(map[string]interface{}{
			"source":      "com.example.custom",
			"detail-type": "TestEvent",
		})
		resp, err := client.TestEventPattern(ctx, &eventbridge.TestEventPatternInput{
			EventPattern: aws.String(string(pattern)),
			Event:        aws.String(string(event)),
		})
		if err != nil {
			return fmt.Errorf("test event pattern: %v", err)
		}
		if !resp.Result {
			return fmt.Errorf("expected pattern to match, got false")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "TestEventPattern_NoMatch", func() error {
		pattern, _ := json.Marshal(map[string]interface{}{
			"source": []string{"com.example.other"},
		})
		event, _ := json.Marshal(map[string]interface{}{
			"source":      "com.example.custom",
			"detail-type": "TestEvent",
		})
		resp, err := client.TestEventPattern(ctx, &eventbridge.TestEventPatternInput{
			EventPattern: aws.String(string(pattern)),
			Event:        aws.String(string(event)),
		})
		if err != nil {
			return fmt.Errorf("test event pattern: %v", err)
		}
		if resp.Result {
			return fmt.Errorf("expected pattern not to match, got true")
		}
		return nil
	}))

	return results
}
