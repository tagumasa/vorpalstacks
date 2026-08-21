package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// Stream batch delivery: the retry policy, bisection on function error,
// and the wire-record helpers shared by the Kinesis and DynamoDB stream
// polling paths.

// esmRetryInterval is the pause between retries of a failed event source
// mapping batch.
const esmRetryInterval = time.Second

// esmFunctionError reports that the mapped function itself failed (the
// Handled or Unhandled classification) as opposed to an invoke-path
// infrastructure error. AWS retry, bisection and failure-destination
// semantics key off the function failing.
type esmFunctionError struct {
	classification  string
	statusCode      int64
	executedVersion string
}

func (e *esmFunctionError) Error() string {
	return fmt.Sprintf("function error: %s", e.classification)
}

// invokeWithRetry drives one batch invocation with the mapping's retry
// policy. onSuccess inspects a successful invocation result; a non-nil
// return treats the invocation as failed (function errors and responses
// without required members both count) and consumes a retry attempt.
// retryAttemptsOf returns how many times a batch is invoked before its
// budget is spent: MaximumRetryAttempts + 1, or a single attempt when
// the budget is unlimited (-1) and the poll loop does the retrying.
func retryAttemptsOf(mapping *lambdastore.EventSourceMapping) int {
	if mapping.MaximumRetryAttempts >= 0 {
		return int(mapping.MaximumRetryAttempts) + 1
	}
	return 1
}

func (p *esmPoller) invokeWithRetry(ctx context.Context, mapping *lambdastore.EventSourceMapping, payload []byte, onSuccess func(*lambdastore.InvocationResult) error) error {
	attempts := retryAttemptsOf(mapping)
	var lastErr error
	for i := 0; i < attempts; i++ {
		result, err := p.invokeLambda(ctx, mapping.FunctionArn, payload)
		switch {
		case err != nil:
			lastErr = err
		case result == nil:
			lastErr = fmt.Errorf("esm: invocation returned nil result")
		case result.FunctionError != "":
			lastErr = &esmFunctionError{
				classification:  result.FunctionError,
				statusCode:      result.StatusCode,
				executedVersion: result.ExecutedVersion,
			}
		case onSuccess != nil:
			if cerr := onSuccess(result); cerr != nil {
				lastErr = cerr
			} else {
				return nil
			}
		default:
			return nil
		}
		if i < attempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(esmRetryInterval):
			}
		}
	}
	return lastErr
}

// invokeBatchWithRetry invokes the mapped function honouring
// MaximumRetryAttempts: values >= 0 bound the retries (0 = a single
// attempt, the batch is discarded after it fails); -1 (the default) keeps
// the historical single-attempt poll. A response carrying a function error
// fails the batch just like an infrastructure error, because the invoke
// transport completes (HTTP 200) even when the function raises an error.
func (p *esmPoller) invokeBatchWithRetry(ctx context.Context, mapping *lambdastore.EventSourceMapping, payload []byte) error {
	return p.invokeWithRetry(ctx, mapping, payload, nil)
}

// streamBatchItem pairs one wire-format stream record with the sequence
// identifier used to advance the checkpoint. record holds either the
// Kinesis wire map or a DynamoDB Streams record; both marshal to the
// documented event forms.
type streamBatchItem struct {
	record interface{}
	seq    string
}

// kinesisRecordSeq reads the sequence number out of the Kinesis wire
// record form.
func kinesisRecordSeq(record map[string]interface{}) string {
	k, ok := record["kinesis"].(map[string]interface{})
	if !ok {
		return ""
	}
	seq, _ := k["sequenceNumber"].(string)
	return seq
}

// dynamoDBRecordSeq reads the sequence number out of the DynamoDB Streams
// wire record form.
func dynamoDBRecordSeq(record *eventbus.DynamoDBStreamRecord) string {
	if record == nil {
		return ""
	}
	seq, _ := record.Dynamodb["SequenceNumber"].(string)
	return seq
}

// marshalStreamBatch renders the batch event payload carrying the given
// stream records.
func marshalStreamBatch(items []streamBatchItem) []byte {
	records := make([]interface{}, len(items))
	for i, item := range items {
		records[i] = item.record
	}
	payload, err := json.Marshal(map[string]interface{}{"Records": records})
	if err != nil {
		return nil
	}
	return payload
}

// processStreamBatch delivers one stream batch with the documented
// error-handling semantics: the batch is retried honouring
// MaximumRetryAttempts; when it still fails and
// BisectBatchOnFunctionError is set, it is split in two and the halves are
// retried recursively until the failure is isolated to a single record. A
// batch that exhausted a finite retry budget is discarded so one poisoned
// record cannot block the shard; the default infinite budget (-1) reports
// the failure and consumes nothing, so the next poll re-reads the batch.
// The returned outcome's lastConsumed is the sequence identifier of the
// last consumed record; short marks a partial batch response whose cursor
// ends the consumed prefix, and discarded reports whether any record was
// dropped after exhausting retries.
func (p *esmPoller) processStreamBatch(ctx context.Context, mapping *lambdastore.EventSourceMapping, src streamSource, items []streamBatchItem) batchOutcome {
	var report batchResponseReport
	invokeErr := p.invokeWithRetry(ctx, mapping, marshalStreamBatch(items), batchResponseSink(mapping, &report))
	if invokeErr == nil {
		// A complete success consumes the whole batch; a partial batch
		// response consumes only the prefix before the first failed
		// record, which the next cycle re-reads.
		outcome := batchOutcome{lastConsumed: streamPartialOutcome(items, report), delivered: true}
		if outcome.lastConsumed != items[len(items)-1].seq {
			outcome.short = true
			outcome.reported = countReportedItems(items, report.failedIDs)
		}
		return outcome
	}
	if mapping.BisectBatchOnFunctionError && len(items) > 1 {
		mid := len(items) / 2
		first := p.processStreamBatch(ctx, mapping, src, items[:mid])
		if first.err != nil {
			// The first half is still failing with an infinite retry
			// budget; nothing after it may be consumed.
			return first
		}
		if first.short {
			// The first half ended on a partial batch response: its
			// cursor terminates the consumed prefix, so the second half
			// is not invoked — its records re-read from the partial
			// checkpoint on the next cycle.
			return first
		}
		second := p.processStreamBatch(ctx, mapping, src, items[mid:])
		if second.err != nil {
			// The second half is still failing with an infinite retry
			// budget; the first half's consumption survives it.
			return batchOutcome{
				lastConsumed: first.lastConsumed,
				discarded:    first.discarded || second.discarded,
				delivered:    first.delivered || second.delivered,
				err:          second.err,
			}
		}
		if second.lastConsumed == "" {
			// The second half's partial response consumed nothing; the
			// batch's cursor stays at the first half's end.
			second.lastConsumed = first.lastConsumed
		}
		second.discarded = second.discarded || first.discarded
		second.delivered = second.delivered || first.delivered
		return second
	}
	if mapping.MaximumRetryAttempts >= 0 {
		p.deliverDiscardedBatch(ctx, mapping, src, streamFailureBatchInfoOf(src, items),
			marshalStreamBatch(items), retryAttemptsOf(mapping), discardedBatchResponse(invokeErr))
		return batchOutcome{lastConsumed: items[len(items)-1].seq, discarded: true}
	}
	return batchOutcome{err: invokeErr}
}
