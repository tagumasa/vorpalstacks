package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
		// The model's Event member documentation makes id, account,
		// source, time, region, resources and detail-type mandatory.
		event, _ := json.Marshal(map[string]interface{}{
			"id":          "test-event-pattern",
			"account":     "000000000000",
			"source":      "com.example.custom",
			"time":        "2015-11-11T21:29:54Z",
			"region":      "us-east-1",
			"resources":   []string{},
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
		// The model's Event member documentation makes id, account,
		// source, time, region, resources and detail-type mandatory.
		event, _ := json.Marshal(map[string]interface{}{
			"id":          "test-event-pattern",
			"account":     "000000000000",
			"source":      "com.example.custom",
			"time":        "2015-11-11T21:29:54Z",
			"region":      "us-east-1",
			"resources":   []string{},
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

	results = append(results, r.RunTest("events", "TestEventPattern_OrMatch", func() error {
		pattern, _ := json.Marshal(map[string]interface{}{
			"$or": []map[string]interface{}{
				{"source": []string{"com.example.other"}},
				{"detail-type": []string{"TestEvent"}},
			},
		})
		event, _ := json.Marshal(map[string]interface{}{
			"id":          "test-event-pattern-or",
			"account":     "000000000000",
			"source":      "com.example.custom",
			"time":        "2015-11-11T21:29:54Z",
			"region":      "us-east-1",
			"resources":   []string{},
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
			return fmt.Errorf("expected the $or second disjunct to match, got false")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "TestEventPattern_InvalidPattern", func() error {
		_, err := client.TestEventPattern(ctx, &eventbridge.TestEventPatternInput{
			EventPattern: aws.String(`{"detail": {"state": [{"regexp": "x"}]}}`),
			Event:        aws.String(`{"id": "i", "account": "a", "source": "s", "time": "t", "region": "r", "resources": [], "detail-type": "d"}`),
		})
		if err == nil {
			return fmt.Errorf("expected InvalidEventPatternException for an unknown operator")
		}
		var iep *types.InvalidEventPatternException
		if !errors.As(err, &iep) {
			return fmt.Errorf("expected InvalidEventPatternException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutEvents_EntryContract", func() error {
		resp, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					EventBusName: aws.String(busName),
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("EntryContract"),
					Detail:       aws.String(`{"k":"v"}`),
				},
				{
					EventBusName: aws.String(busName),
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("EntryContract"),
					Detail:       aws.String(`not-json{`),
				},
				{
					EventBusName: aws.String(busName),
					Source:       aws.String("com.example.test"),
					Detail:       aws.String(`{"k":"v"}`),
				},
				{
					EventBusName: aws.String(busName),
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("EntryContract"),
					Detail:       aws.String(`null`),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put events: %v", err)
		}
		if resp.FailedEntryCount != 3 {
			return fmt.Errorf("expected 3 failed entries, got %d", resp.FailedEntryCount)
		}
		if resp.Entries[0].EventId == nil || *resp.Entries[0].EventId == "" {
			return fmt.Errorf("the complete entry must succeed")
		}
		// The per-entry codes are the documented ErrorCode vocabulary:
		// MalformedDetail for invalid Detail JSON (the null literal is a
		// valid JSON value but not an object), InvalidArgument for an
		// incomplete entry.
		if code := aws.ToString(resp.Entries[1].ErrorCode); code != "MalformedDetail" {
			return fmt.Errorf("malformed Detail must fail as MalformedDetail, got %q", code)
		}
		if code := aws.ToString(resp.Entries[2].ErrorCode); code != "InvalidArgument" {
			return fmt.Errorf("an entry missing DetailType must fail as InvalidArgument, got %q", code)
		}
		if code := aws.ToString(resp.Entries[3].ErrorCode); code != "MalformedDetail" {
			return fmt.Errorf("a null Detail must fail as MalformedDetail, got %q", code)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutEvents_WholeRequestFailures", func() error {
		// A request in which no entry includes Source, DetailType and
		// Detail fails entirely rather than per entry.
		_, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{{EventBusName: aws.String(busName)}},
		})
		if err == nil {
			return fmt.Errorf("a request with no complete entry must fail entirely")
		}

		// The summed entry size must stay under 1 MB (1,048,576 bytes).
		bigDetail := `{"pad":"` + strings.Repeat("a", 1<<20) + `"}`
		_, err = client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					EventBusName: aws.String(busName),
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("EntryContract"),
					Detail:       aws.String(bigDetail),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("a request at the 1 MB entry-size ceiling must fail entirely")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "PutEvents_IngressContract", func() error {
		// EventBusName accepts the ARN form ("The name or ARN of the event
		// bus to receive the event").
		busARN := fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", r.region, r.accountID, busName)
		resp, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					EventBusName: aws.String(busARN),
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("IngressContract"),
					Detail:       aws.String(`{"k":"v"}`),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put events with an ARN-form bus name: %v", err)
		}
		if resp.FailedEntryCount != 0 || resp.Entries[0].EventId == nil {
			return fmt.Errorf("the ARN-form bus name must resolve and succeed: %+v", resp)
		}

		// A non-existent bus is a documented silent drop: a 200 response,
		// no failed entry, and an event ID.
		drop, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					EventBusName: aws.String("no-such-bus-ingress-contract"),
					Source:       aws.String("com.example.test"),
					DetailType:   aws.String("IngressContract"),
					Detail:       aws.String(`{"k":"v"}`),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put events to a non-existent bus: %v", err)
		}
		if drop.FailedEntryCount != 0 {
			return fmt.Errorf("a non-existent bus must not fail the entry (documented drop), got %d failed", drop.FailedEntryCount)
		}
		if drop.Entries[0].EventId == nil {
			return fmt.Errorf("a dropped-bus entry still reports an event ID: %+v", drop)
		}
		return nil
	}))

	return results
}
