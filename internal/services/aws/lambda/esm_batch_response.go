package lambda

import (
	"encoding/json"
	"fmt"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// reportBatchItemFailures is the FunctionResponseTypes enum value that turns
// on partial batch responses.
const reportBatchItemFailures = "ReportBatchItemFailures"

// batchItemFailuresReport captures the documented StreamsEventResponse /
// SQSBatchResponse wire shape:
//
//	{"batchItemFailures": [{"itemIdentifier": "<id>"}]}
//
// The item identifier is the message id for SQS sources and the sequence
// number for Kinesis and DynamoDB stream sources. Tumbling-window
// invocations are out of scope: their response contract is the window
// state object, not a batch report.
type batchItemFailuresReport struct {
	BatchItemFailures *[]batchItemFailure `json:"batchItemFailures"`
}

type batchItemFailure struct {
	ItemIdentifier *string `json:"itemIdentifier"`
}

// batchResponseReport is the parsed outcome of one partial batch response:
// either a complete failure, or the set of item identifiers the function
// reported as failed (empty when the batch succeeded completely).
type batchResponseReport struct {
	completeFailure bool
	failedIDs       map[string]struct{}
}

// batchReportError marks a partial batch response Lambda must treat as a
// complete failure. "Lambda treats a batch as a complete failure if you
// return any of the following: an empty string itemIdentifier, a null
// itemIdentifier, an itemIdentifier with a bad key name. Lambda retries
// failures based on your retry strategy."
type batchReportError struct{}

func (e *batchReportError) Error() string {
	return "esm: batch item failures response was not a valid report"
}

// reportsBatchItemFailures answers whether the mapping asked for partial
// batch responses. "Even when your function code returns partial batch
// failure responses, these responses are not processed by Lambda unless the
// ReportBatchItemFailures feature is explicitly turned on."
func reportsBatchItemFailures(mapping *lambdastore.EventSourceMapping) bool {
	for _, t := range mapping.FunctionResponseTypes {
		if t == reportBatchItemFailures {
			return true
		}
	}
	return false
}

// parseBatchItemFailuresReport interprets a successful invocation's payload
// as a partial batch response. "Lambda treats a batch as a complete success
// if you return any of the following: an empty batchItemFailure list, a
// null batchItemFailure list, an empty EventResponse, a null EventResponse."
// A payload that is not recognisable as the report shape at all is a
// complete failure: with reporting enabled the response contract is the
// report, and an undeserialisable response is a function-level failure.
func parseBatchItemFailuresReport(payload []byte) batchResponseReport {
	completeFailure := batchResponseReport{completeFailure: true}
	if len(payload) == 0 {
		// An empty EventResponse: no return value at all.
		return batchResponseReport{}
	}
	doc, ok := finalJSONDocument(payload)
	if !ok {
		return completeFailure
	}
	var report batchItemFailuresReport
	if err := json.Unmarshal(doc, &report); err != nil {
		// A JSON null EventResponse deserialises without members.
		if string(doc) == "null" {
			return batchResponseReport{}
		}
		return completeFailure
	}
	if report.BatchItemFailures == nil || len(*report.BatchItemFailures) == 0 {
		return batchResponseReport{}
	}
	failed := make(map[string]struct{}, len(*report.BatchItemFailures))
	for _, item := range *report.BatchItemFailures {
		if item.ItemIdentifier == nil || *item.ItemIdentifier == "" {
			return completeFailure
		}
		failed[*item.ItemIdentifier] = struct{}{}
	}
	return batchResponseReport{failedIDs: failed}
}

// batchResponseSink adapts a successful invocation to the poller flow: a
// valid partial batch response is stored for the caller, an invalid one is
// returned as an error so the mapping's retry policy applies to the whole
// batch. Mappings without reporting ignore every response.
func batchResponseSink(mapping *lambdastore.EventSourceMapping, report *batchResponseReport) func(*lambdastore.InvocationResult) error {
	return func(result *lambdastore.InvocationResult) error {
		if !reportsBatchItemFailures(mapping) {
			return nil
		}
		parsed := parseBatchItemFailuresReport(result.Payload)
		if parsed.completeFailure {
			return &batchReportError{}
		}
		*report = parsed
		return nil
	}
}

// firstFailedItem returns the index of the first batch item whose
// identifier the function reported as failed, or -1 when nothing failed.
// "If the batchItemFailures array contains multiple items, Lambda uses the
// record with the lowest sequence number as the checkpoint. Lambda then
// retries all records starting from that checkpoint." Items are held in
// stream order, so the first matching item carries the lowest sequence
// number; identifiers that belong to no item of the batch are ignored —
// the checkpoint rule keys on the batch's own records.
func firstFailedItem(items []streamBatchItem, failed map[string]struct{}) int {
	if len(failed) == 0 {
		return -1
	}
	for i, item := range items {
		if _, ok := failed[item.seq]; ok {
			return i
		}
	}
	return -1
}

// streamPartialOutcome maps a parsed partial batch response onto the
// stream consumption cursor: a complete success consumes the whole batch,
// a partial response consumes only the prefix before the first failed
// record — reads pick up strictly after the cursor, so the failed record
// and everything behind it are re-read on the next cycle.
func streamPartialOutcome(items []streamBatchItem, report batchResponseReport) (lastConsumed string) {
	idx := firstFailedItem(items, report.failedIDs)
	switch {
	case idx < 0:
		return items[len(items)-1].seq
	case idx == 0:
		return ""
	default:
		return items[idx-1].seq
	}
}

// countReportedItems counts the batch's items whose identifiers the
// function reported as failed, mirroring the in-batch scoping of
// firstFailedItem.
func countReportedItems(items []streamBatchItem, failed map[string]struct{}) int {
	n := 0
	for _, item := range items {
		if _, ok := failed[item.seq]; ok {
			n++
		}
	}
	return n
}

// reportedSQSFailureCount counts the messages of one poll whose ids the
// function reported as failed; the number feeds the mapping's processing
// result so a partial batch stays visible.
func reportedSQSFailureCount(msgIDs []string, report batchResponseReport) int {
	n := 0
	for _, id := range msgIDs {
		if _, ok := report.failedIDs[id]; ok {
			n++
		}
	}
	return n
}

// sqsPartialResult renders the LastProcessingResult text after a partial
// batch response, mirroring the shape of the delete-failure report.
func sqsPartialResult(reported int) string {
	return fmt.Sprintf("%d message(s) reported in batchItemFailures", reported)
}

// streamPartialResult renders the LastProcessingResult text after a
// partial batch response on a stream event source, mirroring the SQS
// wording so partial responses stay visible on every event source.
func streamPartialResult(reported int) string {
	return fmt.Sprintf("%d record(s) reported in batchItemFailures", reported)
}
