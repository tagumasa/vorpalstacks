package apps

import (
	"context"
	"fmt"
	"strings"

	storecommon "vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/eventbus"
	storekinesis "vorpalstacks/internal/store/aws/kinesis"
	storesns "vorpalstacks/internal/store/aws/sns"
	storesqs "vorpalstacks/internal/store/aws/sqs"
)

// sqsInvokerAdapter adapts the SQS concrete store to the eventbus.SQSInvoker
// interface, so that cross-service consumers send messages through the bus
// instead of holding a direct store reference.
type sqsInvokerAdapter struct {
	store storesqs.SQSStoreInterface
}

// GetQueueByName looks up a queue by name and returns its URL.
func (a *sqsInvokerAdapter) GetQueueByName(_ context.Context, queueName string) (string, error) {
	q, err := a.store.GetQueueByName(queueName)
	if err != nil {
		return "", err
	}
	return q.URL, nil
}

// GetQueueARN looks up a queue by URL and returns its ARN.
func (a *sqsInvokerAdapter) GetQueueARN(_ context.Context, queueURL string) (string, error) {
	q, err := a.store.GetQueue(queueURL)
	if err != nil {
		return "", err
	}
	return q.ARN, nil
}

// SendMessage sends a message to the specified queue, returning the message ID and MD5.
func (a *sqsInvokerAdapter) SendMessage(_ context.Context, queueURL string, body string, delaySeconds int64, messageAttributes map[string]string) (string, string, error) {
	msg := &storesqs.Message{
		Body:              body,
		DelaySeconds:      int32(delaySeconds),
		MessageAttributes: convertToSQSMessageAttributes(messageAttributes),
	}
	result, err := a.store.SendMessage(queueURL, msg)
	if err != nil {
		return "", "", err
	}
	return result.ID, result.MD5OfBody, nil
}

// ReceiveMessage retrieves messages from the specified queue.
func (a *sqsInvokerAdapter) ReceiveMessage(_ context.Context, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]eventbus.ReceivedSQSMessage, error) {
	msgs, err := a.store.ReceiveMessage(queueURL, maxMessages, visibilityTimeout, waitTimeSeconds)
	if err != nil {
		return nil, err
	}
	out := make([]eventbus.ReceivedSQSMessage, len(msgs))
	for i, m := range msgs {
		out[i] = eventbus.ReceivedSQSMessage{
			MessageID:                       m.ID,
			ReceiptHandle:                   m.ReceiptHandle,
			Body:                            m.Body,
			MD5OfBody:                       m.MD5OfBody,
			MD5OfMessageAttributes:          m.MD5OfMessageAttributes,
			SentTimestamp:                   m.SentTimestamp,
			ApproximateReceiveCount:         m.ApproximateReceiveCount,
			ApproximateFirstReceiveTimestamp: m.ApproximateFirstReceiveTimestamp,
			SequenceNumber:                  m.SequenceNumber,
			MessageDeduplicationID:          m.MessageDeduplicationID,
			MessageGroupID:                  m.MessageGroupID,
		}
		if m.MessageAttributes != nil {
			out[i].MessageAttributes = convertFromSQSMessageAttributes(m.MessageAttributes)
		}
	}
	return out, nil
}

// DeleteMessage deletes a message from the specified queue.
func (a *sqsInvokerAdapter) DeleteMessage(_ context.Context, queueURL string, receiptHandle string) error {
	return a.store.DeleteMessage(queueURL, receiptHandle)
}

// snsInvokerAdapter adapts the SNS concrete store and publisher to the
// eventbus.SNSInvoker interface.
type snsInvokerAdapter struct {
	store     storesns.SNSStoreInterface
	publisher snsPublisher
}

// snsPublisher publishes a message to an SNS topic by ARN.
type snsPublisher interface {
	PublishToTopic(ctx context.Context, accountID, region, topicArn, message string) error
}

// GetTopic retrieves the topic ARN for the given topic ARN.
func (a *snsInvokerAdapter) GetTopic(_ context.Context, topicARN string) (string, error) {
	topic, err := a.store.GetTopic(topicARN)
	if err != nil {
		return "", err
	}
	return topic.Arn, nil
}

// ListSubscriptionsByTopic returns subscriptions for the given topic ARN.
func (a *snsInvokerAdapter) ListSubscriptionsByTopic(_ context.Context, topicARN string) ([]eventbus.SubscriptionInfo, error) {
	result, err := a.store.ListSubscriptionsByTopic(topicARN, storecommon.ListOptions{})
	if err != nil {
		return nil, err
	}
	out := make([]eventbus.SubscriptionInfo, len(result.Items))
	for i, sub := range result.Items {
		out[i] = eventbus.SubscriptionInfo{
			SubscriptionARN:     sub.SubscriptionArn,
			Protocol:            sub.Protocol,
			Endpoint:            sub.Endpoint,
			TopicARN:            sub.TopicArn,
			PendingConfirmation: sub.PendingConfirmation,
		}
	}
	return out, nil
}

// PublishToTopic publishes a message to the given SNS topic ARN.
func (a *snsInvokerAdapter) PublishToTopic(ctx context.Context, topicARN string, message string, subject string, messageAttributes map[string]string) (string, error) {
	if a.publisher == nil {
		return "", fmt.Errorf("sns: publisher not configured")
	}
	parts := strings.Split(topicARN, ":")
	if len(parts) < 5 {
		return "", fmt.Errorf("sns: invalid topic ARN: %s", topicARN)
	}
	accountID := parts[4]
	region := parts[3]
	if err := a.publisher.PublishToTopic(ctx, accountID, region, topicARN, message); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", topicARN), nil
}

// StoreMessage persists arbitrary data keyed by the given key.
func (a *snsInvokerAdapter) StoreMessage(_ context.Context, key string, data any) error {
	return a.store.Put(key, data)
}

// kinesisInvokerAdapter adapts the Kinesis concrete store to the
// eventbus.KinesisInvoker interface.
type kinesisInvokerAdapter struct {
	store storekinesis.KinesisStoreInterface
}

// ListShards lists the shards in the given Kinesis stream.
func (a *kinesisInvokerAdapter) ListShards(_ context.Context, streamName string) ([]eventbus.ShardInfo, error) {
	shards, err := a.store.ListShards(streamName, nil, "", 0)
	if err != nil {
		return nil, err
	}
	out := make([]eventbus.ShardInfo, len(shards))
	for i, s := range shards {
		endSeq := ""
		if s.SequenceNumberRange != nil {
			endSeq = s.SequenceNumberRange.EndingSequenceNumber
		}
		out[i] = eventbus.ShardInfo{
			ShardID:                s.ShardID,
			SequenceNumberRangeEnd: endSeq,
		}
	}
	return out, nil
}

// PutRecord puts a record into the given Kinesis stream on an open shard.
func (a *kinesisInvokerAdapter) PutRecord(_ context.Context, streamName string, partitionKey string, data []byte) (string, error) {
	shards, err := a.store.ListShards(streamName, nil, "", 0)
	if err != nil {
		return "", err
	}
	openShard := selectOpenShard(shards)
	if openShard == nil {
		return "", fmt.Errorf("kinesis: no open shard in stream %s", streamName)
	}
	record, err := a.store.PutRecord(streamName, openShard.ShardID, partitionKey, string(data))
	if err != nil {
		return "", err
	}
	return record.SequenceNumber, nil
}

// CreateShardIterator creates a shard iterator for the given stream and shard.
func (a *kinesisInvokerAdapter) CreateShardIterator(_ context.Context, streamName string, shardID string, iteratorType string, startingSequenceNumber string) (string, error) {
	iterator, err := a.store.CreateShardIterator(streamName, shardID, iteratorType, startingSequenceNumber, nil)
	if err != nil {
		return "", err
	}
	return iterator.SequenceNumber, nil
}

// GetRecords retrieves records from a Kinesis shard starting at the given sequence number.
func (a *kinesisInvokerAdapter) GetRecords(_ context.Context, streamName string, shardID string, startingSequenceNumber string, limit int32) ([]eventbus.KinesisRecord, string, error) {
	records, nextSeq, err := a.store.GetRecords(streamName, shardID, startingSequenceNumber, limit, true)
	if err != nil {
		return nil, "", err
	}
	out := make([]eventbus.KinesisRecord, len(records))
	for i, r := range records {
		out[i] = eventbus.KinesisRecord{
			SequenceNumber:              r.SequenceNumber,
			PartitionKey:                r.PartitionKey,
			Data:                        []byte(r.Data),
			ApproximateArrivalTimestamp: r.ApproximateArrivalTimestamp,
		}
	}
	return out, nextSeq, nil
}

// eventsInvokerAdapter adapts the EventBridge store Put function to the
// eventbus.EventsInvoker interface.
type eventsInvokerAdapter struct {
	putFn func(key string, data any) error
}

// PutEvent persists an event in the EventBridge store.
func (a *eventsInvokerAdapter) PutEvent(_ context.Context, key string, event any) error {
	return a.putFn(key, event)
}

// convertToSQSMessageAttributes converts a simple string map to SQS
// MessageAttributeValue map for SendMessage.
func convertToSQSMessageAttributes(attrs map[string]string) map[string]*storesqs.MessageAttributeValue {
	if attrs == nil {
		return nil
	}
	out := make(map[string]*storesqs.MessageAttributeValue, len(attrs))
	for k, v := range attrs {
		out[k] = &storesqs.MessageAttributeValue{
			StringValue: &v,
			DataType:    "String",
		}
	}
	return out
}

// convertFromSQSMessageAttributes converts SQS MessageAttributeValue map
// back to a simple string map for cross-service consumers.
func convertFromSQSMessageAttributes(attrs map[string]*storesqs.MessageAttributeValue) map[string]string {
	if attrs == nil {
		return nil
	}
	out := make(map[string]string, len(attrs))
	for k, v := range attrs {
		if v.StringValue != nil {
			out[k] = *v.StringValue
		}
	}
	return out
}

// selectOpenShard returns the first open shard (no ending sequence number)
// from the list, falling back to the last shard if all are closed.
func selectOpenShard(shards []*storekinesis.Shard) *storekinesis.Shard {
	for _, s := range shards {
		if s.SequenceNumberRange != nil && s.SequenceNumberRange.EndingSequenceNumber == "" {
			return s
		}
	}
	if len(shards) > 0 {
		return shards[len(shards)-1]
	}
	return nil
}
