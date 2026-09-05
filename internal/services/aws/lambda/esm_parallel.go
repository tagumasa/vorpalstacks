package lambda

import (
	"context"
	"fmt"
	"sync"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/resilience"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// Parallel batch processing.
//
// ParallelizationFactor is "The number of batches to process from each
// shard concurrently": the poller fetches that many batches per shard per
// cycle and processes them concurrently, while "Lambda still ensures
// in-order processing at the partition-key level" — a batch only starts
// once every earlier batch sharing one of its record keys has completed.

// batchOutcome reports the result of one concurrently processed batch.
type batchOutcome struct {
	lastConsumed string // sequence of the last record consumed by the batch
	discarded    bool   // records were dropped after exhausting retries
	delivered    bool   // at least one record was delivered to the function
	// short marks a batch whose consumption ended on a partial batch
	// response before its final record: the contiguous consumed prefix
	// ends with this batch, so no later batch may move the checkpoint
	// past its cursor.
	short bool
	// reported counts this batch's items the function reported in its
	// batchItemFailures response.
	reported int
	err      error // the batch still fails with an infinite retry budget
}

// parallelizationFactorOf clamps the mapping's ParallelizationFactor into
// the documented range; 0 (unset) behaves as the default single batch.
func parallelizationFactorOf(mapping *lambdastore.EventSourceMapping) int {
	pf := int(mapping.ParallelizationFactor)
	if pf < int(lambdastore.MinParallelizationFactor) {
		return int(lambdastore.MinParallelizationFactor)
	}
	if pf > int(lambdastore.MaxParallelizationFactor) {
		return int(lambdastore.MaxParallelizationFactor)
	}
	return pf
}

// kinesisRecordKeys returns the partition keys of a Kinesis batch.
func kinesisRecordKeys(records []invokers.KinesisRecord) map[string]struct{} {
	keys := make(map[string]struct{}, len(records))
	for _, rec := range records {
		keys[rec.PartitionKey] = struct{}{}
	}
	return keys
}

// dynamoDBRecordKeys returns the canonical item keys of a DynamoDB Streams
// batch; Go renders map keys in sorted order, so fmt.Sprint is a stable
// canonical form.
func dynamoDBRecordKeys(records []invokers.DynamoDBStreamRecord) map[string]struct{} {
	keys := make(map[string]struct{}, len(records))
	for i := range records {
		if keyMap, ok := records[i].Dynamodb["Keys"].(map[string]interface{}); ok {
			keys[fmt.Sprint(keyMap)] = struct{}{}
		}
	}
	return keys
}

// runOrderedBatches launches one goroutine per batch. keys[i] holds the
// record keys of batch i; a batch waits for every earlier batch sharing a
// key before invoking, which preserves per-key ordering while disjoint
// batches run concurrently. Waiting only ever reaches backwards, so the
// schedule is deadlock-free.
func runOrderedBatches(ctx context.Context, count int, keys []map[string]struct{}, run func(ctx context.Context, idx int) batchOutcome) []batchOutcome {
	outcomes := make([]batchOutcome, count)
	dones := make([]chan struct{}, count)
	for i := range dones {
		dones[i] = make(chan struct{})
	}
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer resilience.RecoverPanic("lambda esm parallel batch")
			defer close(dones[i])
			for j := 0; j < i; j++ {
				if sharesAnyKey(keys[i], keys[j]) {
					select {
					case <-dones[j]:
					case <-ctx.Done():
						// A batch that never started must not fold into the
						// consumed prefix as a silent success: failing it
						// stops the prefix at this position so the checkpoint
						// cannot advance past unprocessed records.
						outcomes[i] = batchOutcome{err: ctx.Err()}
						return
					}
				}
			}
			outcomes[i] = run(ctx, i)
		}(i)
	}
	wg.Wait()
	return outcomes
}

// sharesAnyKey reports whether two record key sets intersect.
func sharesAnyKey(a, b map[string]struct{}) bool {
	for k := range a {
		if _, ok := b[k]; ok {
			return true
		}
	}
	return false
}

// prefixOutcome folds ordered batch outcomes into the checkpoint position:
// only the contiguous prefix of fully consumed batches advances it. The
// first failure suspends everything after it — those records are re-read
// on the next poll, matching the at-least-once contract. A batch that
// ended on a partial batch response terminates the prefix the same way:
// later batches already ran concurrently, but "Lambda uses the record
// with the lowest sequence number as the checkpoint. Lambda then retries
// all records starting from that checkpoint", so their cursors must not
// move the checkpoint past the reported failure.
func prefixOutcome(outcomes []batchOutcome) (lastConsumed string, delivered, discarded bool, reported int, failure error) {
	clamped := false
	for _, oc := range outcomes {
		delivered = delivered || oc.delivered
		discarded = discarded || oc.discarded
		reported += oc.reported
		if oc.err != nil {
			return lastConsumed, delivered, discarded, reported, oc.err
		}
		if clamped {
			continue
		}
		if oc.lastConsumed != "" {
			lastConsumed = oc.lastConsumed
		}
		if oc.short {
			clamped = true
		}
	}
	return lastConsumed, delivered, discarded, reported, nil
}
