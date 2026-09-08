package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/logs"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

const (
	// sqsReceiveMessageMax is the maximum number of messages that can be
	// returned by a single SQS ReceiveMessage call.
	sqsReceiveMessageMax = int32(10)
)

// esmSQSRecord represents a single SQS message formatted as an ESM event
// record in the Lambda event payload. The structure matches the AWS Lambda
// SQS event format documented at:
// https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html
type esmSQSRecord struct {
	MessageID                        string                 `json:"messageId"`
	ReceiptHandle                    string                 `json:"receiptHandle"`
	Body                             string                 `json:"body"`
	MD5OfBody                        string                 `json:"md5OfBody"`
	MD5OfMessageAttributes           string                 `json:"md5OfMessageAttributes"`
	MessageAttributes                map[string]interface{} `json:"messageAttributes,omitempty"`
	EventSourceARN                   string                 `json:"eventSourceArn"`
	EventSource                      string                 `json:"eventSource"`
	AWSRegion                        string                 `json:"awsRegion"`
	ApproximateReceiveCount          string                 `json:"approximateReceiveCount"`
	ApproximateFirstReceiveTimestamp string                 `json:"approximateFirstReceiveTimestamp"`
	SentTimestamp                    string                 `json:"sentTimestamp"`
	SenderID                         string                 `json:"senderId"`
}

// esmSQSEvent is the full Lambda event payload for an SQS event source
// mapping. It contains an array of SQS message records.
type esmSQSEvent struct {
	Records []esmSQSRecord `json:"Records"`
}

func (p *esmPoller) processSQSMapping(ctx context.Context, mapping *lambdastore.EventSourceMapping) {
	if p.bus == nil {
		return
	}

	_, _, region, accountID, _ := arnutil.SplitARN(mapping.EventSourceArn)
	if region == "" || accountID == "" {
		p.log("failed to parse event source ARN", "arn", mapping.EventSourceArn)
		return
	}

	queueName := arnutil.ExtractQueueNameFromARN(mapping.EventSourceArn)
	if queueName == "" {
		p.log("failed to extract queue name from ARN", "arn", mapping.EventSourceArn)
		return
	}

	queueURL, err := p.bus.SQSInvoker().GetQueueByName(ctx, region, queueName)
	if err != nil {
		p.log("sqs queue not found by name", "queue", queueName, "mapping", mapping.UUID, "error", err)
		return
	}

	batchSize := clampESMBatchSize(mapping.BatchSize, mapping.EventSourceArn)

	perCallMax := sqsReceiveMessageMax
	if perCallMax > batchSize {
		perCallMax = batchSize
	}

	// The receive wait is bounded to one poll interval: the poll cycle runs
	// synchronously (pollAll waits for every mapping to finish), so a
	// receive that long-polls longer than the interval would stall the
	// whole cycle and delay every other mapping's poll. Long polling lives
	// in the SQS store; the poller only needs each cycle to stay
	// responsive, and the batching window still governs record gathering.
	sqsReceiveWaitSeconds := int32(p.interval / time.Second)
	if sqsReceiveWaitSeconds < 1 {
		sqsReceiveWaitSeconds = 1
	}

	var allMessages []invokers.ReceivedSQSMessage
	remaining := batchSize
	for remaining > 0 {
		fetchCount := perCallMax
		if fetchCount > remaining {
			fetchCount = remaining
		}

		msgs, err := p.bus.SQSInvoker().ReceiveMessage(ctx, region, queueURL, fetchCount, nil, sqsReceiveWaitSeconds)
		if err != nil {
			p.log("sqs receive failed", "queue", queueName, "mapping", mapping.UUID, "error", err)
			break
		}

		if len(msgs) == 0 {
			break
		}

		allMessages = append(allMessages, msgs...)
		remaining -= int32(len(msgs))

		if int32(len(msgs)) < fetchCount {
			break
		}
	}

	if len(allMessages) == 0 {
		return
	}

	records := make([]esmSQSRecord, 0, len(allMessages))
	receiptHandles := make([]string, 0, len(allMessages))
	messageIDs := make([]string, 0, len(allMessages))
	for _, msg := range allMessages {
		records = append(records, receivedSQSMessageToRecord(msg, mapping.EventSourceArn, region))
		receiptHandles = append(receiptHandles, msg.ReceiptHandle)
		messageIDs = append(messageIDs, msg.MessageID)
	}

	// Apply event filtering before invocation.
	records = filterSQSRecords(records, mapping.FilterCriteria)
	if len(records) == 0 {
		// All messages were filtered out.  AWS treats filtered-out
		// messages as successfully processed — delete them from the
		// queue to prevent infinite re-polling after visibility
		// timeout expires.
		for _, handle := range receiptHandles {
			if err := p.bus.SQSInvoker().DeleteMessage(ctx, region, queueURL, handle); err != nil {
				p.log("failed to delete filtered-out message", "queue", queueName, "error", err)
			}
		}
		if err := p.esmStore.SetProcessingResult(mapping.UUID, "No errors."); err != nil {
			logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.String("error", err.Error()))
		}
		return
	}

	eventPayload := esmSQSEvent{Records: records}
	payload, err := json.Marshal(eventPayload)
	if err != nil {
		p.log("failed to marshal ESM event payload", "queue", queueName, "error", err)
		return
	}

	fnName := arnutil.ExtractFunctionNameFromARN(mapping.FunctionArn)
	if fnName == "" {
		p.log("failed to extract function name from ARN", "arn", mapping.FunctionArn)
		return
	}

	var report batchResponseReport
	invokeErr := p.invokeWithRetry(ctx, mapping, payload, batchResponseSink(mapping, &report))

	if invokeErr != nil {
		p.log("lambda invocation failed", "function", fnName, "queue", queueName, "error", invokeErr)
		if err := p.esmStore.SetProcessingResult(mapping.UUID, invokeErr.Error()); err != nil {
			logs.Warn("esm: failed to set state after SQS invocation error",
				logs.String("mapping", mapping.UUID), logs.Err(err))
		}
		return
	}

	// A partial batch response deletes only the messages the function did
	// not report: "To make messages id2 and id4 visible again in your
	// queue, your function should return" their identifiers — the reported
	// messages return with the queue's visibility timeout.
	deleteFailures := 0
	for i, handle := range receiptHandles {
		if _, failed := report.failedIDs[messageIDs[i]]; failed {
			continue
		}
		if err := p.bus.SQSInvoker().DeleteMessage(ctx, region, queueURL, handle); err != nil {
			p.log("failed to delete message", "queue", queueName, "error", err)
			deleteFailures++
		}
	}

	lastResult := "No errors."
	if deleteFailures > 0 {
		lastResult = fmt.Sprintf("%d message(s) failed to delete", deleteFailures)
	} else if reported := reportedSQSFailureCount(messageIDs, report); reported > 0 {
		lastResult = sqsPartialResult(reported)
	}

	if err := p.esmStore.SetProcessingResult(mapping.UUID, lastResult); err != nil {
		logs.Error("esm: failed to set state", logs.String("mapping", mapping.UUID), logs.String("error", err.Error()))
	}
}

// receivedSQSMessageToRecord converts an invokers.ReceivedSQSMessage into an
// ESM SQS record matching the Lambda event format.
func receivedSQSMessageToRecord(msg invokers.ReceivedSQSMessage, eventSourceArn, region string) esmSQSRecord {
	record := esmSQSRecord{
		MessageID:               msg.MessageID,
		ReceiptHandle:           msg.ReceiptHandle,
		Body:                    msg.Body,
		MD5OfBody:               msg.MD5OfBody,
		MD5OfMessageAttributes:  msg.MD5OfMessageAttributes,
		EventSourceARN:          eventSourceArn,
		EventSource:             "aws:sqs",
		AWSRegion:               region,
		ApproximateReceiveCount: fmt.Sprintf("%d", msg.ApproximateReceiveCount),
		SentTimestamp:           fmt.Sprintf("%d", msg.SentTimestamp.UnixMilli()),
	}

	if msg.ApproximateFirstReceiveTimestamp.IsZero() {
		record.ApproximateFirstReceiveTimestamp = record.SentTimestamp
	} else {
		record.ApproximateFirstReceiveTimestamp = fmt.Sprintf("%d", msg.ApproximateFirstReceiveTimestamp.UnixMilli())
	}

	if msg.SequenceNumber != "" {
		if record.MessageAttributes == nil {
			record.MessageAttributes = make(map[string]interface{})
		}
		record.MessageAttributes["SequenceNumber"] = map[string]string{
			"stringValue": msg.SequenceNumber,
			"dataType":    "String",
		}
	}

	if msg.MessageDeduplicationID != "" {
		if record.MessageAttributes == nil {
			record.MessageAttributes = make(map[string]interface{})
		}
		record.MessageAttributes["MessageDeduplicationId"] = map[string]string{
			"stringValue": msg.MessageDeduplicationID,
			"dataType":    "String",
		}
	}

	if msg.MessageGroupID != "" {
		if record.MessageAttributes == nil {
			record.MessageAttributes = make(map[string]interface{})
		}
		record.MessageAttributes["MessageGroupId"] = map[string]string{
			"stringValue": msg.MessageGroupID,
			"dataType":    "String",
		}
	}

	return record
}
