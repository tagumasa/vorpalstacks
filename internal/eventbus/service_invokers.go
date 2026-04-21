package eventbus

import (
	"context"
	"time"
)

// LambdaInvoker invokes a Lambda function and returns the status code and
// response payload. This is the same contract as common.LambdaInvoker but
// defined here so that the Bus interface does not depend on the common
// package.
type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (statusCode int64, responsePayload []byte, err error)
}

// SQSInvoker provides SQS operations for cross-service consumers.
// Consumers call these methods instead of holding a direct reference to the
// SQS store.
type SQSInvoker interface {
	GetQueueByName(ctx context.Context, queueName string) (queueURL string, err error)
	GetQueueARN(ctx context.Context, queueURL string) (queueARN string, err error)
	SendMessage(ctx context.Context, queueURL string, body string, delaySeconds int64, messageAttributes map[string]string) (messageID string, md5OfBody string, err error)
	ReceiveMessage(ctx context.Context, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]ReceivedSQSMessage, error)
	DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error
}

// ReceivedSQSMessage carries the fields of an SQS message returned by
// ReceiveMessage that cross-service consumers need.
type ReceivedSQSMessage struct {
	MessageID                     string            `json:"messageId"`
	ReceiptHandle                 string            `json:"receiptHandle"`
	Body                          string            `json:"body"`
	MD5OfBody                     string            `json:"md5OfBody"`
	MessageAttributes             map[string]string `json:"messageAttributes,omitempty"`
	MD5OfMessageAttributes        string            `json:"md5OfMessageAttributes,omitempty"`
	SentTimestamp                 time.Time         `json:"sentTimestamp,omitempty"`
	ApproximateReceiveCount       int32             `json:"approximateReceiveCount,omitempty"`
	ApproximateFirstReceiveTimestamp time.Time       `json:"approximateFirstReceiveTimestamp,omitempty"`
	SequenceNumber                string            `json:"sequenceNumber,omitempty"`
	MessageDeduplicationID        string            `json:"messageDeduplicationId,omitempty"`
	MessageGroupID                string            `json:"messageGroupId,omitempty"`
}

// SNSInvoker provides SNS operations for cross-service consumers.
// Consumers call these methods instead of holding a direct reference to the
// SNS store.
type SNSInvoker interface {
	GetTopic(ctx context.Context, topicARN string) (string, error)
	ListSubscriptionsByTopic(ctx context.Context, topicARN string) ([]SubscriptionInfo, error)
	PublishToTopic(ctx context.Context, topicARN string, message string, subject string, messageAttributes map[string]string) (messageID string, err error)
	StoreMessage(ctx context.Context, key string, data any) error
}

// SubscriptionInfo carries the fields of an SNS subscription that
// cross-service consumers need.
type SubscriptionInfo struct {
	SubscriptionARN     string `json:"subscriptionArn"`
	Protocol            string `json:"protocol"`
	Endpoint            string `json:"endpoint"`
	TopicARN            string `json:"topicArn"`
	PendingConfirmation bool   `json:"pendingConfirmation"`
}

// KinesisInvoker provides Kinesis operations for cross-service consumers.
// Consumers call these methods instead of holding a direct reference to the
// Kinesis store.
type KinesisInvoker interface {
	ListShards(ctx context.Context, streamName string) ([]ShardInfo, error)
	PutRecord(ctx context.Context, streamName string, partitionKey string, data []byte) (sequenceNumber string, err error)
	CreateShardIterator(ctx context.Context, streamName string, shardID string, iteratorType string, startingSequenceNumber string) (iteratorSequenceNumber string, err error)
	GetRecords(ctx context.Context, streamName string, shardID string, startingSequenceNumber string, limit int32) (records []KinesisRecord, nextSequenceNumber string, err error)
}

// ShardInfo carries the fields of a Kinesis shard that cross-service
// consumers need.
type ShardInfo struct {
	ShardID                       string
	SequenceNumberRangeEnd        string
}

// KinesisRecord carries the fields of a Kinesis record that cross-service
// consumers need.
type KinesisRecord struct {
	SequenceNumber              string
	PartitionKey                string
	Data                        []byte
	ApproximateArrivalTimestamp time.Time
}

// EventsInvoker provides EventBridge store operations for cross-service
// consumers. Consumers call these methods instead of holding a direct
// reference to the EventBridge store.
type EventsInvoker interface {
	PutEvent(ctx context.Context, key string, event any) error
}

// EC2Invoker provides EC2 subnet and security group lookup operations for
// cross-service consumers. Consumers call these methods instead of using
// the generic GetInvoker/Invoke pattern.
type EC2Invoker interface {
	LookupSubnet(ctx context.Context, region string, subnetId string) (vpcId string, availabilityZone string, err error)
	LookupSecurityGroup(ctx context.Context, region string, groupId string) (vpcId string, err error)
}

// DynamoDBInvoker provides DynamoDB item operations for cross-service
// consumers (e.g. AppSync GraphQL resolvers). Consumers call these methods
// instead of holding a direct reference to the DynamoDB store.
//
// Keys and attribute values use map[string]interface{} where each value is
// one of: string, float64, bool, nil ([]byte and number-as-string are
// represented as strings). This avoids a dependency on the store package.
type DynamoDBInvoker interface {
	GetItem(ctx context.Context, region, tableName string, key map[string]interface{}) (map[string]interface{}, error)
	PutItem(ctx context.Context, region, tableName string, key map[string]interface{}, attributes map[string]interface{}) (map[string]interface{}, error)
	DeleteItem(ctx context.Context, region, tableName string, key map[string]interface{}) error
	Scan(ctx context.Context, region, tableName string, limit int) ([]map[string]interface{}, error)
	Query(ctx context.Context, region, tableName, partitionKeyValue string, limit int) ([]map[string]interface{}, error)
	UpdateItem(ctx context.Context, region, tableName string, key map[string]interface{}, attributes map[string]interface{}) error
}
