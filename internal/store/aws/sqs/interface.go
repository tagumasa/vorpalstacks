package sqs

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

// SQSStoreInterface defines operations for managing SQS queues and messages.
type SQSStoreInterface interface {
	CreateQueue(queue *Queue) (*Queue, error)
	GetQueue(queueURL string) (*Queue, error)
	GetQueueByName(queueName string) (*Queue, error)
	UpdateQueue(queue *Queue) error
	DeleteQueue(queueURL string) error
	ListQueues(opts common.ListOptions) (*common.ListResult[Queue], error)
	SetQueueAttributes(queueURL string, attributes map[string]string) error
	AddPermission(queueURL, label string, awsAccountIDs []string, actions []string) error
	RemovePermission(queueURL, label string) error
	ListQueueTags(queueURL string) (map[string]string, error)
	TagQueue(queueURL string, tags map[string]string) error
	UntagQueue(queueURL string, tagKeys []string) error
	ListDeadLetterSourceQueues(dlqARN string) ([]*Queue, error)
	GetMessageCounts(queueURL string) (visible, notVisible, delayed int32)

	SendMessage(queueURL string, message *Message) (*Message, error)
	ReceiveMessage(queueURL string, maxNumberOfMessages int32, visibilityTimeoutPtr *int32, waitTimeSeconds int32) ([]*Message, error)
	DeleteMessage(queueURL, receiptHandle string) error
	ChangeMessageVisibility(queueURL, receiptHandle string, visibilityTimeout int32) error
	PurgeQueue(queueURL string) error

	StartMessageMoveTask(sourceARN, destARN string, maxMessages int32) (*MessageMoveTask, error)
	CancelMessageMoveTask(taskId string) (*MessageMoveTask, error)
	ListMessageMoveTasks(sourceARN string) ([]*MessageMoveTask, error)
	GetMessageMoveTask(taskId string) (*MessageMoveTask, error)

	Storage() storage.BasicStorage
	GetAccountID() string
	GetRegion() string
}

var _ SQSStoreInterface = (*SQSStore)(nil)
