package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// On-failure destination delivery for stream event source mappings.
//
// "To retain records of failed event source mapping invocations, add a
// destination to your function's event source mapping. Each record sent
// to the destination is a JSON document containing metadata about the
// failed invocation. For Amazon S3 destinations, Lambda also sends the
// entire invocation record along with the metadata. You can configure any
// Amazon SNS topic, Amazon SQS queue, Amazon S3 bucket, or Kafka as a
// destination." (Lambda Developer Guide, Kinesis and DynamoDB error
// handling pages). The DestinationConfig model member scopes the feature
// to stream sources (Kinesis, DynamoDB Streams and the Kafka family), so
// only the stream paths deliver; SQS event sources keep using the
// queue's own redrive semantics.

// failureConditionRetryExhausted is the only condition value AWS
// documents for stream on-failure records: every invocation-record
// example in the Kinesis and DynamoDB error handling pages carries
// "RetryAttemptsExhausted", and no distinct value is documented for
// records discarded because they reached MaximumRecordAgeInSeconds, so
// every stream discard reports this condition.
const failureConditionRetryExhausted = "RetryAttemptsExhausted"

// failureRecordVersion mirrors the documented "version": "1.0" member of
// the invocation record.
const failureRecordVersion = "1.0"

// streamSourceKind identifies which stream service a batch came from; it
// selects the batch-info member name of the destination record
// ("KinesisBatchInfo" versus "DDBStreamBatchInfo").
type streamSourceKind string

const (
	streamSourceKinesis  streamSourceKind = "kinesis"
	streamSourceDynamoDB streamSourceKind = "dynamodb"
)

// streamSource identifies the stream and shard a delivered batch came
// from, for the on-failure destination record.
type streamSource struct {
	kind      streamSourceKind
	streamArn string
	shardID   string
}

// failureRequestContext is the requestContext member of the invocation
// record; approximateInvokeCount reports how many times the batch was
// invoked before being discarded.
type failureRequestContext struct {
	RequestID              string `json:"requestId"`
	FunctionARN            string `json:"functionArn"`
	Condition              string `json:"condition"`
	ApproximateInvokeCount int    `json:"approximateInvokeCount"`
}

// failureResponseContext is the responseContext member of the invocation
// record, describing the last failed invocation.
type failureResponseContext struct {
	StatusCode      int    `json:"statusCode"`
	ExecutedVersion string `json:"executedVersion"`
	FunctionError   string `json:"functionError"`
}

// streamFailureBatchInfo carries the batch placement fields. AWS
// documents exactly this field set for both KinesisBatchInfo and
// DDBStreamBatchInfo, so "use the streamArn, shardId,
// startSequenceNumber, and endSequenceNumber fields to obtain the full
// original record" from the stream.
type streamFailureBatchInfo struct {
	ShardID                         string `json:"shardId"`
	StartSequenceNumber             string `json:"startSequenceNumber"`
	EndSequenceNumber               string `json:"endSequenceNumber"`
	ApproximateArrivalOfFirstRecord string `json:"approximateArrivalOfFirstRecord"`
	ApproximateArrivalOfLastRecord  string `json:"approximateArrivalOfLastRecord"`
	BatchSize                       int    `json:"batchSize"`
	StreamARN                       string `json:"streamArn"`
}

// streamFailureRecord is the JSON document sent to the destination.
type streamFailureRecord struct {
	RequestContext     failureRequestContext   `json:"requestContext"`
	ResponseContext    failureResponseContext  `json:"responseContext"`
	Version            string                  `json:"version"`
	Timestamp          string                  `json:"timestamp"`
	KinesisBatchInfo   *streamFailureBatchInfo `json:"KinesisBatchInfo,omitempty"`
	DDBStreamBatchInfo *streamFailureBatchInfo `json:"DDBStreamBatchInfo,omitempty"`
	// Payload carries the whole invocation record, "Only available in S3".
	Payload string `json:"payload,omitempty"`
}

// streamFailureBatchInfoOf derives the batch placement of a discarded
// batch from its first and last items. Arrivals keep each source's
// documented form: the Kinesis event's ISO-8601 millisecond timestamp
// and the DynamoDB record's second-precision RFC 3339 creation time.
func streamFailureBatchInfoOf(src streamSource, items []streamBatchItem) streamFailureBatchInfo {
	info := streamFailureBatchInfo{
		ShardID:   src.shardID,
		StreamARN: src.streamArn,
		BatchSize: len(items),
	}
	if len(items) == 0 {
		return info
	}
	info.StartSequenceNumber = items[0].seq
	info.EndSequenceNumber = items[len(items)-1].seq
	info.ApproximateArrivalOfFirstRecord = itemArrivalOf(items[0])
	info.ApproximateArrivalOfLastRecord = itemArrivalOf(items[len(items)-1])
	return info
}

// itemArrivalOf reads a batch item's record insertion time in the
// documented destination-record form.
func itemArrivalOf(item streamBatchItem) string {
	switch record := item.record.(type) {
	case map[string]interface{}:
		kinesis, ok := record["kinesis"].(map[string]interface{})
		if !ok {
			return ""
		}
		arrival, _ := kinesis["approximateArrivalTimestamp"].(string)
		return arrival
	case *eventbus.DynamoDBStreamRecord:
		if record == nil {
			return ""
		}
		return time.Unix(ddbArrivalUnix(record), 0).UTC().Format(time.RFC3339)
	}
	return ""
}

// discardedBatchResponse renders the responseContext of a batch dropped
// after exhausting its retry budget. The documented example pairs
// statusCode 200, executedVersion "$LATEST" and functionError
// "Unhandled"; a function error keeps its own classification, status
// code and executed version.
func discardedBatchResponse(err error) failureResponseContext {
	response := failureResponseContext{
		StatusCode:      200,
		ExecutedVersion: "$LATEST",
		FunctionError:   "Unhandled",
	}
	if fnErr, ok := err.(*esmFunctionError); ok {
		if fnErr.statusCode != 0 {
			response.StatusCode = int(fnErr.statusCode)
		}
		if fnErr.executedVersion != "" {
			response.ExecutedVersion = fnErr.executedVersion
		}
		if fnErr.classification != "" {
			response.FunctionError = fnErr.classification
		}
	}
	return response
}

// uninvokedBatchResponse renders the responseContext of records
// discarded by MaximumRecordAgeInSeconds before any invocation. The
// documented record shape requires the response members; the discarded
// records were never invoked, so the documented example values stand in.
func uninvokedBatchResponse() failureResponseContext {
	return failureResponseContext{
		StatusCode:      200,
		ExecutedVersion: "$LATEST",
		FunctionError:   "Unhandled",
	}
}

// failureDestinationObjectKey renders the documented S3 object naming
// convention: "aws/lambda/<ESM-UUID>/<shardID>/YYYY/MM/DD/
// YYYY-MM-DDTHH.MM.SS-<Random UUID>".
func failureDestinationObjectKey(mappingUUID, shardID string, now time.Time) string {
	return fmt.Sprintf("aws/lambda/%s/%s/%s/%s-%s",
		mappingUUID,
		shardID,
		now.UTC().Format("2006/01/02"),
		now.UTC().Format("2006-01-02T15.04.05"),
		uuid.NewString())
}

// deliverDiscardedBatch sends the invocation record for one discarded
// batch to the mapping's on-failure destination. A delivery failure is
// logged only: the records are already dropped and "Lambda discards the
// records and continues processing batches from the stream", so the
// poller must not stall on an unreachable destination.
func (p *esmPoller) deliverDiscardedBatch(ctx context.Context, mapping *lambdastore.EventSourceMapping, src streamSource, batch streamFailureBatchInfo, payload []byte, invokeCount int, response failureResponseContext) {
	if mapping.DestinationConfig == nil || mapping.DestinationConfig.OnFailure == nil {
		return
	}
	destination := mapping.DestinationConfig.OnFailure.Destination
	if destination == "" || p.bus == nil {
		return
	}

	record := streamFailureRecord{
		RequestContext: failureRequestContext{
			RequestID:              uuid.NewString(),
			FunctionARN:            mapping.FunctionArn,
			Condition:              failureConditionRetryExhausted,
			ApproximateInvokeCount: invokeCount,
		},
		ResponseContext: response,
		Version:         failureRecordVersion,
		Timestamp:       time.Now().UTC().Format(timeutils.ISO8601UTCFormat),
	}
	switch src.kind {
	case streamSourceKinesis:
		batch.StreamARN = src.streamArn
		record.KinesisBatchInfo = &batch
	case streamSourceDynamoDB:
		batch.StreamARN = src.streamArn
		record.DDBStreamBatchInfo = &batch
	default:
		return
	}

	_, service, _, _, resource := arnutil.SplitARN(destination)
	if service == "s3" {
		// "In addition to all of the fields from the previous example for
		// SQS and SNS destinations, the payload field contains the
		// original invocation record as an escaped JSON string."
		record.Payload = string(payload)
	}
	body, err := json.Marshal(record)
	if err != nil {
		p.log("failed to marshal on-failure destination record",
			"mapping", mapping.UUID, "error", err.Error())
		return
	}

	switch service {
	case "sqs":
		queueName := arnutil.ExtractQueueNameFromARN(destination)
		_, _, queueRegion, _, _ := arnutil.SplitARN(destination)
		queueURL, gerr := p.bus.SQSInvoker().GetQueueByName(ctx, queueRegion, queueName)
		if gerr != nil {
			p.log("failed to resolve on-failure SQS destination",
				"destination", destination, "error", gerr.Error())
			return
		}
		if _, _, serr := p.bus.SQSInvoker().SendMessage(ctx, queueRegion, queueURL, string(body), eventbus.SQSSendOptions{}); serr != nil {
			p.log("failed to deliver to on-failure SQS destination",
				"destination", destination, "error", serr.Error())
		}
	case "sns":
		if _, perr := p.bus.SNSInvoker().PublishToTopic(ctx, destination, string(body), "", nil); perr != nil {
			p.log("failed to deliver to on-failure SNS destination",
				"destination", destination, "error", perr.Error())
		}
	case "s3":
		key := failureDestinationObjectKey(mapping.UUID, src.shardID, time.Now())
		if perr := p.bus.S3Invoker().PutObject(ctx, "", resource, key, body, "application/json"); perr != nil {
			p.log("failed to deliver to on-failure S3 destination",
				"destination", destination, "error", perr.Error())
		}
	default:
		p.log("unsupported on-failure destination service",
			"destination", destination)
	}
}

// discardExpiredDynamoDBRecords delivers stream records already older
// than MaximumRecordAgeInSeconds ("Lambda retries until the records
// expire, exceed the maximum age ... If the error handling measures
// fail, Lambda discards the records") to the on-failure destination and
// returns the fresh remainder. -1, the default, keeps every record.
func (p *esmPoller) discardExpiredDynamoDBRecords(ctx context.Context, mapping *lambdastore.EventSourceMapping, src streamSource, records []eventbus.DynamoDBStreamRecord) []eventbus.DynamoDBStreamRecord {
	if mapping.MaximumRecordAgeInSeconds <= 0 {
		return records
	}
	cutoff := time.Now().Add(-time.Duration(mapping.MaximumRecordAgeInSeconds) * time.Second)
	fresh := records[:0]
	var expired []eventbus.DynamoDBStreamRecord
	for i := range records {
		if time.Unix(ddbArrivalUnix(&records[i]), 0).After(cutoff) {
			fresh = append(fresh, records[i])
			continue
		}
		expired = append(expired, records[i])
	}
	if len(expired) > 0 {
		items := make([]streamBatchItem, len(expired))
		for i := range expired {
			items[i] = streamBatchItem{record: &expired[i], seq: dynamoDBRecordSeq(&expired[i])}
		}
		p.deliverDiscardedBatch(ctx, mapping, src, streamFailureBatchInfoOf(src, items),
			marshalStreamBatch(items), 0, uninvokedBatchResponse())
	}
	return fresh
}
