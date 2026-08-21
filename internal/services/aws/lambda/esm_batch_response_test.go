package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

func TestParseBatchItemFailuresReport(t *testing.T) {
	cases := []struct {
		name         string
		payload      string
		completeFail bool
		failed       []string
	}{
		{name: "empty payload is a complete success", payload: ``},
		{name: "null event response is a complete success", payload: `null`},
		{name: "empty event response is a complete success", payload: `{}`},
		{name: "null batchItemFailures is a complete success", payload: `{"batchItemFailures":null}`},
		{name: "empty batchItemFailures is a complete success", payload: `{"batchItemFailures":[]}`},
		{
			name:    "valid identifiers are collected",
			payload: `{"batchItemFailures":[{"itemIdentifier":"a"},{"itemIdentifier":"b"}]}`,
			failed:  []string{"a", "b"},
		},
		{name: "empty string itemIdentifier is a complete failure", payload: `{"batchItemFailures":[{"itemIdentifier":""}]}`, completeFail: true},
		{name: "null itemIdentifier is a complete failure", payload: `{"batchItemFailures":[{"itemIdentifier":null}]}`, completeFail: true},
		{name: "bad key name is a complete failure", payload: `{"batchItemFailures":[{"identifier":"a"}]}`, completeFail: true},
		{name: "non-report payload is a complete failure", payload: `"ok"`, completeFail: true},
		{name: "invalid JSON payload is a complete failure", payload: `not json`, completeFail: true},
		{
			name:         "one invalid entry fails the whole report",
			payload:      `{"batchItemFailures":[{"itemIdentifier":"a"},{"itemIdentifier":null}]}`,
			completeFail: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			report := parseBatchItemFailuresReport([]byte(tc.payload))
			if report.completeFailure != tc.completeFail {
				t.Fatalf("completeFailure = %v, want %v", report.completeFailure, tc.completeFail)
			}
			if len(report.failedIDs) != len(tc.failed) {
				t.Fatalf("failed ids = %v, want %v", report.failedIDs, tc.failed)
			}
			for _, want := range tc.failed {
				if _, ok := report.failedIDs[want]; !ok {
					t.Fatalf("failed ids %v must contain %q", report.failedIDs, want)
				}
			}
		})
	}
}

func TestStreamPartialOutcome(t *testing.T) {
	items := []streamBatchItem{
		{seq: "10"}, {seq: "20"}, {seq: "30"},
	}
	failed := func(ids ...string) map[string]struct{} {
		m := make(map[string]struct{}, len(ids))
		for _, id := range ids {
			m[id] = struct{}{}
		}
		return m
	}
	if got := streamPartialOutcome(items, batchResponseReport{}); got != "30" {
		t.Fatalf("complete success must consume the batch, got cursor %q", got)
	}
	if got := streamPartialOutcome(items, batchResponseReport{failedIDs: failed("30")}); got != "20" {
		t.Fatalf("a failed tail must consume its prefix, got cursor %q", got)
	}
	if got := streamPartialOutcome(items, batchResponseReport{failedIDs: failed("20", "30")}); got != "10" {
		t.Fatalf("the lowest failed record rules the checkpoint, got cursor %q", got)
	}
	if got := streamPartialOutcome(items, batchResponseReport{failedIDs: failed("10")}); got != "" {
		t.Fatalf("a failed head must consume nothing, got cursor %q", got)
	}
	if got := streamPartialOutcome(items, batchResponseReport{failedIDs: failed("999")}); got != "30" {
		t.Fatalf("identifiers outside the batch are ignored, got cursor %q", got)
	}
}

// scriptedBatchReporter answers each invocation with a partial batch
// response that reports failSeq as failed for the first failCalls
// invocations containing it. needle spells the sequence member of the
// wire form the payloads carry; the DynamoDB and Kinesis event shapes
// capitalise it differently.
type scriptedBatchReporter struct {
	mu        sync.Mutex
	payloads  []string
	failSeq   string
	failCalls int
	needle    string
}

func (f *scriptedBatchReporter) seqNeedle() string {
	if f.needle == "" {
		return `"SequenceNumber":%q`
	}
	return f.needle + `:%q`
}

func (f *scriptedBatchReporter) invoke(_ context.Context, _ string, payload []byte) (*lambdastore.InvocationResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.payloads = append(f.payloads, string(payload))
	failures := []interface{}{}
	if f.failCalls > 0 && strings.Contains(string(payload), fmt.Sprintf(f.seqNeedle(), f.failSeq)) {
		failures = append(failures, map[string]interface{}{"itemIdentifier": f.failSeq})
		f.failCalls--
	}
	raw, err := json.Marshal(map[string]interface{}{"batchItemFailures": failures})
	if err != nil {
		return nil, err
	}
	return &lambdastore.InvocationResult{StatusCode: 200, Payload: raw}, nil
}

func (f *scriptedBatchReporter) snapshot() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.payloads...)
}

func deliveryCount(payloads []string, seq string) int {
	return deliveryCountForm(payloads, seq, `"SequenceNumber":%q`)
}

// deliveryCountForm counts the payloads carrying seq under the given
// wire-form needle.
func deliveryCountForm(payloads []string, seq, needleForm string) int {
	n := 0
	needle := fmt.Sprintf(needleForm, seq)
	for _, p := range payloads {
		if strings.Contains(p, needle) {
			n++
		}
	}
	return n
}

// TestDynamoDBPartialBatchResponseHoldsCheckpoint drives the plain DDB
// delivery path with a partial batch response: the reported record's
// predecessor becomes the checkpoint and the failed record (with everything
// behind it) is re-read on the next cycle.
func TestDynamoDBPartialBatchResponseHoldsCheckpoint(t *testing.T) {
	p, stream, mapping, _ := newDDBAnchorPoller(t, "TRIM_HORIZON")
	mapping.FunctionResponseTypes = []string{"ReportBatchItemFailures"}
	reporter := &scriptedBatchReporter{failSeq: "2", failCalls: 1}
	p.invoke = reporter.invoke

	stream.publish(1)
	stream.publish(2)
	stream.publish(3)

	p.processDynamoDBStreamsMapping(context.Background(), mapping)
	if got := p.kinesisCP["ddb:"+mapping.UUID]; got != "1" {
		t.Fatalf("a partial response must checkpoint before the failed record, got %q", got)
	}

	p.processDynamoDBStreamsMapping(context.Background(), mapping)
	if got := p.kinesisCP["ddb:"+mapping.UUID]; got != "3" {
		t.Fatalf("a later successful cycle must advance the checkpoint, got %q", got)
	}

	payloads := reporter.snapshot()
	if got := deliveryCount(payloads, "1"); got != 1 {
		t.Fatalf("the record before the failure must be delivered once, got %d", got)
	}
	if got := deliveryCount(payloads, "2"); got != 2 {
		t.Fatalf("the failed record must be re-delivered, got %d deliveries", got)
	}
	if got := deliveryCount(payloads, "3"); got != 2 {
		t.Fatalf("the record after the failure is retried with it (at-least-once), got %d deliveries", got)
	}
}

// TestDynamoDBInvalidBatchReportDiscardsBatch pins the complete-failure
// wiring: a response that is not a valid report consumes the retry budget
// and ends in the documented discard path.
func TestDynamoDBInvalidBatchReportDiscardsBatch(t *testing.T) {
	p, stream, mapping, _ := newDDBAnchorPoller(t, "TRIM_HORIZON")
	mapping.FunctionResponseTypes = []string{"ReportBatchItemFailures"}
	mapping.MaximumRetryAttempts = 0
	invocations := 0
	p.invoke = func(_ context.Context, _ string, payload []byte) (*lambdastore.InvocationResult, error) {
		invocations++
		// Not the report shape: a complete failure by contract.
		return &lambdastore.InvocationResult{StatusCode: 200, Payload: []byte(`"ok"`)}, nil
	}

	stream.publish(1)
	stream.publish(2)

	p.processDynamoDBStreamsMapping(context.Background(), mapping)

	if got := p.kinesisCP["ddb:"+mapping.UUID]; got != "2" {
		t.Fatalf("an exhausted invalid report must discard and advance past the batch, got %q", got)
	}
	result, err := p.esmStore.Get(mapping.UUID)
	if err != nil {
		t.Fatalf("get mapping: %v", err)
	}
	if !strings.Contains(result.LastProcessingResult, "discarded") {
		t.Fatalf("the discard must be reported, got %q", result.LastProcessingResult)
	}
}
