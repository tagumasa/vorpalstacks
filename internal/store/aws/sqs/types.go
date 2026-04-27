package sqs

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"
)

// Queue represents an SQS queue.
type Queue struct {
	Name                          string                 `json:"name"`
	URL                           string                 `json:"url"`
	ARN                           string                 `json:"arn"`
	Region                        string                 `json:"region"`
	AccountID                     string                 `json:"accountId"`
	CreatedTimestamp              time.Time              `json:"createdTimestamp"`
	LastModifiedTimestamp         time.Time              `json:"lastModifiedTimestamp"`
	VisibilityTimeout             int32                  `json:"visibilityTimeout"`
	MaximumMessageSize            int32                  `json:"maximumMessageSize"`
	MessageRetentionPeriod        int32                  `json:"messageRetentionPeriod"`
	DelaySeconds                  int32                  `json:"delaySeconds"`
	ReceiveMessageWaitTimeSeconds int32                  `json:"receiveMessageWaitTimeSeconds"`
	Attributes                    map[string]string      `json:"attributes"`
	Tags                          map[string]string      `json:"tags"`
	Permissions                   map[string]*Permission `json:"permissions,omitempty"`
	Policy                        string                 `json:"policy,omitempty"`
	RedrivePolicy                 *RedrivePolicy         `json:"redrivePolicy,omitempty"`
	FifoQueue                     bool                   `json:"fifoQueue"`
	ContentBasedDeduplication     bool                   `json:"contentBasedDeduplication"`
}

// Permission represents an SQS queue permission.
type Permission struct {
	Label         string   `json:"label"`
	AWSAccountIDs []string `json:"awsAccountIds"`
	Actions       []string `json:"actions"`
}

// RedrivePolicy represents the redrive policy for an SQS queue.
type RedrivePolicy struct {
	DeadLetterTargetARN string `json:"deadLetterTargetArn"`
	MaxReceiveCount     int32  `json:"maxReceiveCount"`
}

// ParseRedrivePolicy parses a RedrivePolicy JSON string, accepting both
// string and integer formats for maxReceiveCount per AWS convention.
func ParseRedrivePolicy(data string) (*RedrivePolicy, error) {
	var raw struct {
		DeadLetterTargetARN string          `json:"deadLetterTargetArn"`
		MaxReceiveCount     json.RawMessage `json:"maxReceiveCount"`
	}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, err
	}
	var count int32
	if len(raw.MaxReceiveCount) > 0 {
		trimmed := string(raw.MaxReceiveCount)
		if trimmed[0] == '"' {
			unquoted := trimmed[1 : len(trimmed)-1]
			v, err := strconv.ParseInt(unquoted, 10, 32)
			if err != nil {
				return nil, err
			}
			count = int32(v)
		} else {
			v, err := strconv.ParseInt(trimmed, 10, 32)
			if err != nil {
				return nil, err
			}
			count = int32(v)
		}
	}
	return &RedrivePolicy{
		DeadLetterTargetARN: raw.DeadLetterTargetARN,
		MaxReceiveCount:     count,
	}, nil
}

// Message represents an SQS message.
type Message struct {
	ID                               string                            `json:"id"`
	Body                             string                            `json:"body"`
	MD5OfBody                        string                            `json:"md5OfBody"`
	MD5OfMessageAttributes           string                            `json:"md5OfMessageAttributes"`
	MessageAttributes                map[string]*MessageAttributeValue `json:"messageAttributes"`
	ReceiptHandle                    string                            `json:"receiptHandle"`
	QueueURL                         string                            `json:"queueUrl"`
	QueueARN                         string                            `json:"queueArn"`
	SentTimestamp                    time.Time                         `json:"sentTimestamp"`
	DelaySeconds                     int32                             `json:"delaySeconds"`
	VisibilityTimeout                int32                             `json:"visibilityTimeout"`
	ReceivedAt                       time.Time                         `json:"receivedAt,omitempty"`
	VisibleAfter                     time.Time                         `json:"visibleAfter,omitempty"`
	ApproximateReceiveCount          int32                             `json:"approximateReceiveCount"`
	ApproximateFirstReceiveTimestamp time.Time                         `json:"approximateFirstReceiveTimestamp"`
	SequenceNumber                   string                            `json:"sequenceNumber,omitempty"`
	MessageDeduplicationID           string                            `json:"messageDeduplicationId,omitempty"`
	MessageGroupID                   string                            `json:"messageGroupId,omitempty"`
	Attributes                       map[string]string                 `json:"attributes"`
}

// MessageAttributeValue represents an SQS message attribute.
type MessageAttributeValue struct {
	StringValue      *string  `json:"stringValue,omitempty"`
	BinaryValue      []byte   `json:"binaryValue,omitempty"`
	StringListValues []string `json:"stringListValues,omitempty"`
	BinaryListValues [][]byte `json:"binaryListValues,omitempty"`
	DataType         string   `json:"dataType"`
}

// NewQueue creates a new Queue with default values.
func NewQueue(name, region, accountID string) *Queue {
	now := time.Now().UTC()
	return &Queue{
		Name:                          name,
		Region:                        region,
		AccountID:                     accountID,
		CreatedTimestamp:              now,
		LastModifiedTimestamp:         now,
		VisibilityTimeout:             30,
		MaximumMessageSize:            262144,
		MessageRetentionPeriod:        345600,
		DelaySeconds:                  0,
		ReceiveMessageWaitTimeSeconds: 0,
		Attributes:                    make(map[string]string),
		Tags:                          make(map[string]string),
		Permissions:                   make(map[string]*Permission),
	}
}

// NewMessage creates a new Message with the given body.
func NewMessage(body string) *Message {
	return &Message{
		Body:              body,
		MD5OfBody:         calculateMD5(body),
		SentTimestamp:     time.Now().UTC(),
		MessageAttributes: make(map[string]*MessageAttributeValue),
		Attributes:        make(map[string]string),
	}
}

func calculateMD5(s string) string {
	if s == "" {
		return "d41d8cd98f00b204e9800998ecf8427e"
	}
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// MessageMoveTask represents a message move task between queues.
type MessageMoveTask struct {
	TaskId              string    `json:"taskId"`
	SourceQueueARN      string    `json:"sourceQueueArn"`
	DestinationQueueARN string    `json:"destinationQueueArn"`
	Status              string    `json:"status"`
	MaxNumberOfMessages int32     `json:"maxNumberOfMessages"`
	MovedMessages       int32     `json:"movedMessages"`
	FailureMessages     int32     `json:"failureMessages"`
	StartTime           time.Time `json:"startTime"`
	EndTime             time.Time `json:"endTime,omitempty"`
}
