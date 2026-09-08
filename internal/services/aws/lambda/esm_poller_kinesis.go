package lambda

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/invokers"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

func (p *esmPoller) processKinesisMapping(ctx context.Context, mapping *lambdastore.EventSourceMapping) {
	if p.bus == nil {
		return
	}

	_, _, streamRegion, _, resource := arnutil.SplitARN(mapping.EventSourceArn)
	if streamRegion == "" {
		p.log("failed to parse Kinesis event source ARN", "arn", mapping.EventSourceArn)
		return
	}

	streamName := resource
	if idx := strings.Index(resource, "stream/"); idx != -1 {
		streamName = resource[idx+len("stream/"):]
	}

	shards, err := p.bus.KinesisInvoker().ListShards(ctx, streamName)
	if err != nil {
		p.log("failed to list shards for Kinesis ESM", "stream", streamName, "error", err)
		return
	}

	cycle := p.newStreamCycle(mapping, "Kinesis ESM")
	windowSeconds := int64(mapping.TumblingWindowInSeconds)
	pf := parallelizationFactorOf(mapping)

	fnName := arnutil.ExtractFunctionNameFromARN(mapping.FunctionArn)
	if fnName == "" {
		p.log("failed to extract function name from ARN", "arn", mapping.FunctionArn)
		return
	}

	for _, shard := range shards {
		cpKey := fmt.Sprintf("%s:%s:%s", mapping.UUID, streamName, shard.ShardID)
		src := streamSource{kind: streamSourceKinesis, streamArn: mapping.EventSourceArn, shardID: shard.ShardID}

		if shard.SequenceNumberRangeEnd != "" {
			// A closed shard never yields new records; deliver the final
			// invocation for any window still open on it and flush any
			// records its batching window was still gathering.
			if cycle.windowed {
				p.closeEndedShardWindow(ctx, mapping, cpKey, shard.ShardID)
			}
			if cycle.buffered {
				if items := p.dropStreamBuffer(cpKey); len(items) > 0 {
					if out := p.flushStreamBuffer(ctx, mapping, src, cpKey, "", items); out.err != nil {
						p.log("failed to flush batching window of ended shard",
							"mapping", mapping.UUID, "shard", shard.ShardID, "error", out.err)
					}
				}
			}
			continue
		}

		readFrom := p.streamReadPosition(cpKey, cycle.windowed, cycle.buffered)

		var iteratorType, iteratorSeqNum string
		var iteratorTimestamp *time.Time
		if readFrom != "" {
			iteratorType = "AFTER_SEQUENCE_NUMBER"
			iteratorSeqNum = readFrom
		} else {
			// Honour the user-configured StartingPosition for the initial read.
			// Default to TRIM_HORIZON for backward compatibility.
			switch mapping.StartingPosition {
			case "LATEST":
				iteratorType = "LATEST"
			case "AT_TIMESTAMP":
				iteratorType = "AT_TIMESTAMP"
				if !mapping.StartingPositionTimestamp.IsZero() {
					ts := mapping.StartingPositionTimestamp.UTC()
					iteratorTimestamp = &ts
				}
			default:
				iteratorType = "TRIM_HORIZON"
			}
		}

		iteratorSeq, err := p.bus.KinesisInvoker().CreateShardIterator(ctx, streamName, shard.ShardID, iteratorType, iteratorSeqNum, iteratorTimestamp)
		if err != nil {
			p.log("failed to create shard iterator", "stream", streamName, "shard", shard.ShardID, "error", err)
			continue
		}

		// Fetch up to ParallelizationFactor batches; a chained read picks
		// up strictly after the previous batch's last record. The first
		// read of a LATEST anchor includes the anchor record itself: the
		// mapping activates immediately on this platform, without the
		// activation lag that on AWS lets records published between
		// CreateEventSourceMapping and the iterator anchoring still be
		// delivered, so the anchor record — the stream's latest at first
		// poll — must be read inclusively.
		anchorInitialLATEST := readFrom == "" && iteratorType == "LATEST"
		var batches []streamFetchedBatch
		pos := iteratorSeq
		for i := 0; i < pf; i++ {
			records, next, gerr := p.bus.KinesisInvoker().GetRecords(ctx, streamName, shard.ShardID, pos, cycle.batchSize, i == 0 && anchorInitialLATEST)
			if gerr != nil {
				p.log("failed to get records from Kinesis", "stream", streamName, "shard", shard.ShardID, "error", gerr)
				break
			}
			if len(records) == 0 {
				break
			}
			batches = append(batches, p.prepareKinesisBatch(ctx, mapping, src, shard.ShardID, streamRegion, cycle.windowed, windowSeconds, records))
			if int32(len(records)) < cycle.batchSize {
				break
			}
			pos = next
		}

		anchor := ""
		if readFrom == "" {
			anchor = iteratorSeq
		}
		cycle.dispatchBatches(ctx, cpKey, src, batches, anchor)
	}

	cycle.finish()
}

// renderKinesisRecords builds the wire-format maps for one batch, recording
// each record's arrival time for tumbling window boundaries.
func (p *esmPoller) renderKinesisRecords(shardID, streamRegion string, mapping *lambdastore.EventSourceMapping, windowed bool, records []invokers.KinesisRecord) ([]map[string]interface{}, map[string]int64) {
	arrivals := make(map[string]int64, len(records))
	rendered := make([]map[string]interface{}, 0, len(records))
	for _, rec := range records {
		arrivals[rec.SequenceNumber] = rec.ApproximateArrivalTimestamp.Unix()
		record := map[string]interface{}{
			"kinesis": map[string]interface{}{
				"kinesisSchemaVersion":        "1.0",
				"partitionKey":                rec.PartitionKey,
				"sequenceNumber":              rec.SequenceNumber,
				"data":                        string(rec.Data),
				"approximateArrivalTimestamp": rec.ApproximateArrivalTimestamp.Format(timeutils.ISO8601UTCFormat),
			},
			"eventSource":       "aws:kinesis",
			"eventVersion":      "1.0",
			"eventID":           fmt.Sprintf("%s:%s:%s", shardID, rec.SequenceNumber, mapping.UUID),
			"awsRegion":         streamRegion,
			"eventName":         "aws:kinesis:record",
			"invokeIdentityArn": arnutil.NewARNBuilder(p.lambdaSvc.accountID, "").IAM().Role("vorpalstacks-lambda"),
		}
		if windowed {
			// The documented KinesisTimeWindowEvent carries the source
			// ARN on every record as well as at the envelope level.
			record["eventSourceARN"] = mapping.EventSourceArn
		}
		rendered = append(rendered, record)
	}
	return rendered, arrivals
}

// dropExpiredKinesisRecords removes records older than
// MaximumRecordAgeInSeconds (-1, the default, keeps every record) and
// reports the expired batch to the on-failure destination: "Lambda
// retries until the records expire, exceed the maximum age ... If the
// error handling measures fail, Lambda discards the records".
func (p *esmPoller) dropExpiredKinesisRecords(ctx context.Context, mapping *lambdastore.EventSourceMapping, src streamSource, shardID, streamRegion string, windowed bool, records []invokers.KinesisRecord) []invokers.KinesisRecord {
	if mapping.MaximumRecordAgeInSeconds <= 0 {
		return records
	}
	cutoff := time.Now().Add(-time.Duration(mapping.MaximumRecordAgeInSeconds) * time.Second)
	fresh := records[:0]
	var expired []invokers.KinesisRecord
	for _, rec := range records {
		if rec.ApproximateArrivalTimestamp.After(cutoff) {
			fresh = append(fresh, rec)
			continue
		}
		expired = append(expired, rec)
	}
	if len(expired) > 0 {
		rendered, _ := p.renderKinesisRecords(shardID, streamRegion, mapping, windowed, expired)
		items := make([]streamBatchItem, len(rendered))
		for i, kr := range rendered {
			items[i] = streamBatchItem{record: kr, seq: kinesisRecordSeq(kr)}
		}
		p.deliverDiscardedBatch(ctx, mapping, src, streamFailureBatchInfoOf(src, items),
			marshalStreamBatch(items), 0, uninvokedBatchResponse())
	}
	return fresh
}

// prepareKinesisBatch applies age expiry and event filtering and renders
// one fetched Kinesis batch into delivery items with their tumbling-window
// pairings.
func (p *esmPoller) prepareKinesisBatch(ctx context.Context, mapping *lambdastore.EventSourceMapping, src streamSource, shardID, streamRegion string, windowed bool, windowSeconds int64, records []invokers.KinesisRecord) streamFetchedBatch {
	batch := streamFetchedBatch{keys: kinesisRecordKeys(records)}
	if len(records) > 0 {
		batch.readSeq = records[len(records)-1].SequenceNumber
	}
	records = p.dropExpiredKinesisRecords(ctx, mapping, src, shardID, streamRegion, windowed, records)
	rendered, arrivals := p.renderKinesisRecords(shardID, streamRegion, mapping, windowed, records)
	rendered = filterKinesisRecords(rendered, mapping.FilterCriteria)
	batch.items = make([]streamBatchItem, len(rendered))
	batch.witems = make([]windowedStreamItem, len(rendered))
	for i, kr := range rendered {
		batch.items[i] = streamBatchItem{record: kr, seq: kinesisRecordSeq(kr)}
		batch.witems[i] = windowedStreamItem{
			item:        batch.items[i],
			windowStart: windowStartOf(arrivals[batch.items[i].seq], windowSeconds),
		}
	}
	return batch
}
