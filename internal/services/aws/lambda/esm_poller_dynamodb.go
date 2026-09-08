package lambda

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/invokers"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// processDynamoDBStreamsMapping polls a DynamoDB stream for new records,
// batches them, and invokes the mapped Lambda function with the records
// in DynamoDB Streams event format. Checkpoints are persisted per mapping
// to survive restarts.
func (p *esmPoller) processDynamoDBStreamsMapping(ctx context.Context, mapping *lambdastore.EventSourceMapping) {
	if p.bus == nil {
		return
	}

	invoker := p.bus.DynamoDBStreamsInvoker()
	if invoker == nil {
		return
	}

	// Extract table name from the stream ARN:
	// arn:aws:dynamodb:region:account:table/tableName/stream/label
	_, _, streamRegion, _, resource := arnutil.SplitARN(mapping.EventSourceArn)
	if streamRegion == "" {
		p.log("failed to parse DynamoDB stream ARN", "arn", mapping.EventSourceArn)
		return
	}

	tableName := resource
	if idx := strings.Index(resource, "table/"); idx != -1 {
		rest := resource[idx+len("table/"):]
		if slashIdx := strings.Index(rest, "/"); slashIdx != -1 {
			tableName = rest[:slashIdx]
		} else {
			tableName = rest
		}
	}

	cycle := p.newStreamCycle(mapping, "DynamoDB ESM")
	windowSeconds := int64(mapping.TumblingWindowInSeconds)

	// Load the checkpoint for this mapping and resolve the read position.
	checkpointKey := fmt.Sprintf("ddb:%s", mapping.UUID)
	readFrom := p.streamReadPosition(checkpointKey, cycle.windowed, cycle.buffered)

	shardID := invoker.ShardIDForStream(mapping.EventSourceArn)
	src := streamSource{kind: streamSourceDynamoDB, streamArn: mapping.EventSourceArn, shardID: shardID}

	// Reads are exclusive of the read position, so the numeric sequence
	// of a record is a valid "read everything after this" cursor.
	var fromSeq int64
	anchorLATEST := int64(0)
	anchored := false
	if readFrom != "" {
		if v, err := strconv.ParseInt(readFrom, 10, 64); err == nil {
			fromSeq = v
		}
	} else {
		// No durable read position yet: honour the user-configured
		// StartingPosition for the initial read, mirroring the Kinesis
		// path. AT_TIMESTAMP never reaches the poller for DynamoDB —
		// mapping creation rejects it.
		switch mapping.StartingPosition {
		case "LATEST":
			latest, aerr := invoker.GetLatestSequence(ctx, streamRegion, tableName)
			if aerr != nil {
				p.log("failed to read latest DynamoDB stream sequence",
					"table", tableName, "error", aerr.Error())
				return
			}
			anchorLATEST = latest
			anchored = true
			// The anchor is read inclusively: the mapping activates
			// immediately on this platform, without the activation lag
			// that on AWS lets records published between
			// CreateEventSourceMapping and the iterator anchoring still
			// be delivered, so the anchor record — the stream's latest
			// at first poll — must be delivered too.
			fromSeq = latest - 1
			if fromSeq < 0 {
				fromSeq = 0
			}
		default:
			// TRIM_HORIZON (or unset): read from the beginning of the
			// stream. Sequence numbers are assigned from one upwards, so
			// zero reads every record.
			fromSeq = 0
		}
	}

	pf := parallelizationFactorOf(mapping)

	// Fetch up to ParallelizationFactor batches; a chained read picks up
	// strictly after the previous batch's last record. The invoker answers
	// with the last record's sequence, which is a valid read position
	// because reads are exclusive of it.
	var batches []streamFetchedBatch
	from := fromSeq
	for i := 0; i < pf; i++ {
		records, nextSeq, gerr := invoker.GetRecords(ctx, streamRegion, tableName, from, int(cycle.batchSize))
		if gerr != nil {
			p.log("failed to get DynamoDB stream records",
				"table", tableName, "error", gerr.Error())
			if len(batches) == 0 {
				return
			}
			break
		}
		if len(records) == 0 {
			break
		}
		batches = append(batches, p.prepareDynamoDBBatch(ctx, mapping, src, records, nextSeq, windowSeconds))
		if int32(len(records)) < cycle.batchSize {
			break
		}
		from = nextSeq
	}

	anchor := ""
	if anchored {
		anchor = strconv.FormatInt(anchorLATEST, 10)
	}
	cycle.dispatchBatches(ctx, checkpointKey, src, batches, anchor)
	cycle.finish()
}

// prepareDynamoDBBatch applies age expiry and event filtering and renders
// one fetched DynamoDB Streams batch into delivery items with their
// tumbling-window pairings.
func (p *esmPoller) prepareDynamoDBBatch(ctx context.Context, mapping *lambdastore.EventSourceMapping, src streamSource, records []invokers.DynamoDBStreamRecord, nextSeq int64, windowSeconds int64) streamFetchedBatch {
	batch := streamFetchedBatch{keys: dynamoDBRecordKeys(records), readSeq: strconv.FormatInt(nextSeq, 10)}
	records = p.discardExpiredDynamoDBRecords(ctx, mapping, src, records)
	records = filterDynamoDBRecords(records, mapping.FilterCriteria)
	batch.items = make([]streamBatchItem, len(records))
	batch.witems = make([]windowedStreamItem, len(records))
	for i := range records {
		batch.items[i] = streamBatchItem{record: &records[i], seq: dynamoDBRecordSeq(&records[i])}
		batch.witems[i] = windowedStreamItem{
			item:        batch.items[i],
			windowStart: windowStartOf(ddbArrivalUnix(&records[i]), windowSeconds),
		}
	}
	return batch
}
