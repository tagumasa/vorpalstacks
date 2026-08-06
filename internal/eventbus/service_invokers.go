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
	GetFunctionARN(ctx context.Context, functionName string) (string, error)
}

// SQSInvoker provides SQS operations for cross-service consumers.
// Consumers call these methods instead of holding a direct reference to the
// SQS store.
type SQSInvoker interface {
	GetQueueByName(ctx context.Context, queueName string) (queueURL string, err error)
	GetQueueARN(ctx context.Context, queueURL string) (queueARN string, err error)
	SendMessage(ctx context.Context, queueURL string, body string, opts SQSSendOptions) (messageID string, md5OfBody string, err error)
	ReceiveMessage(ctx context.Context, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]ReceivedSQSMessage, error)
	DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error
}

// SQSSendOptions carries optional parameters for cross-service SQS SendMessage
// calls. Callers that do not need FIFO semantics leave MessageGroupID and
// MessageDeduplicationID empty.
type SQSSendOptions struct {
	DelaySeconds           int64
	MessageAttributes      map[string]string
	TypedMessageAttributes map[string]SQSMessageAttribute
	MessageGroupID         string
	MessageDeduplicationID string
}

// SQSMessageAttribute carries a typed SQS message attribute value.
// DataType is one of "String", "Number", or "Binary". For Binary type the
// caller stores the raw bytes in BinaryValue; for String/Number the value
// goes in StringValue.
type SQSMessageAttribute struct {
	DataType    string
	StringValue string
	BinaryValue []byte
}

// ReceivedSQSMessage carries the fields of an SQS message returned by
// ReceiveMessage that cross-service consumers need.
type ReceivedSQSMessage struct {
	MessageID                        string            `json:"messageId"`
	ReceiptHandle                    string            `json:"receiptHandle"`
	Body                             string            `json:"body"`
	MD5OfBody                        string            `json:"md5OfBody"`
	MessageAttributes                map[string]string `json:"messageAttributes,omitempty"`
	MD5OfMessageAttributes           string            `json:"md5OfMessageAttributes,omitempty"`
	SentTimestamp                    time.Time         `json:"sentTimestamp,omitempty"`
	ApproximateReceiveCount          int32             `json:"approximateReceiveCount,omitempty"`
	ApproximateFirstReceiveTimestamp time.Time         `json:"approximateFirstReceiveTimestamp,omitempty"`
	SequenceNumber                   string            `json:"sequenceNumber,omitempty"`
	MessageDeduplicationID           string            `json:"messageDeduplicationId,omitempty"`
	MessageGroupID                   string            `json:"messageGroupId,omitempty"`
}

// SNSInvoker provides SNS operations for cross-service consumers.
// Consumers call these methods instead of holding a direct reference to the
// SNS store.
type SNSInvoker interface {
	GetTopic(ctx context.Context, topicARN string) (string, error)
	ListSubscriptionsByTopic(ctx context.Context, topicARN string) ([]SubscriptionInfo, error)
	PublishToTopic(ctx context.Context, topicARN string, message string, subject string, messageAttributes map[string]string) (messageID string, err error)
	StoreMessage(ctx context.Context, key string, data any) error
	DeleteStoredMessage(ctx context.Context, key string) error
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
	ShardID                string
	SequenceNumberRangeEnd string
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

// EC2Invoker provides EC2 subnet, security group, and VPC lookup operations
// for cross-service consumers. Consumers call these methods instead of using
// the generic GetInvoker/Invoke pattern.
type EC2Invoker interface {
	LookupSubnet(ctx context.Context, region string, subnetId string) (vpcId string, availabilityZone string, err error)
	LookupSecurityGroup(ctx context.Context, region string, groupId string) (vpcId string, err error)
	LookupVPC(ctx context.Context, region string, vpcId string) error
}

// SubnetUsageChecker checks whether a subnet is referenced by a service's
// resources. EC2 calls all registered checkers before deleting a subnet to
// prevent orphaned VPC configuration references (e.g. Lambda VpcConfig,
// Neptune DB subnet groups).
type SubnetUsageChecker interface {
	IsSubnetInUse(ctx context.Context, region, subnetId string) bool
}

// SecurityGroupUsageChecker checks whether a security group is referenced
// by a service's resources. EC2 calls all registered checkers before
// deleting a security group to prevent orphaned VPC configuration
// references (e.g. Lambda VpcConfig.SecurityGroupIds).
type SecurityGroupUsageChecker interface {
	IsSecurityGroupInUse(ctx context.Context, region, sgId string) bool
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

// DynamoDBStreamsInvoker provides DynamoDB Streams operations for
// cross-service consumers (e.g. Lambda ESM poller). Consumers call these
// methods instead of holding a direct reference to the DynamoDB stream store.
type DynamoDBStreamsInvoker interface {
	GetRecords(ctx context.Context, region, tableName string, fromSeq int64, limit int) ([]DynamoDBStreamRecord, int64, error)
	GetLatestSequence(ctx context.Context, region, tableName string) (int64, error)
}

// DynamoDBStreamRecord carries a DynamoDB Streams record for cross-service
// consumption. The fields match the AWS DynamoDB Streams event format.
type DynamoDBStreamRecord struct {
	EventID        string                 `json:"eventID"`
	EventName      string                 `json:"eventName"`
	EventVersion   string                 `json:"eventVersion"`
	EventSource    string                 `json:"eventSource"`
	AWSRegion      string                 `json:"awsRegion"`
	Dynamodb       map[string]interface{} `json:"dynamodb"`
	EventSourceARN string                 `json:"eventSourceARN"`
}

// NeptuneGraphInvoker provides NeptuneGraph query execution for cross-service
// consumers (e.g. AppSync GraphQL resolvers). Consumers call these methods
// instead of holding a direct reference to the NeptuneGraph service.
type NeptuneGraphInvoker interface {
	ExecuteQueryOnGraph(ctx context.Context, graphID string, query string, language string, parameters map[string]interface{}) (interface{}, error)
}

// KMSInvoker provides KMS encryption operations for cross-service consumers
// (e.g. S3 SSE-KMS envelope encryption). Consumers call these methods instead
// of holding a direct reference to the KMS store or HSM backend. The
// sourceArn parameter carries the calling service's resource ARN (e.g. an S3
// bucket ARN) for grant constraint evaluation (aws:SourceArn).
type KMSInvoker interface {
	GenerateDataKey(ctx context.Context, keyID string, keySpec string, encryptionContext map[string]string, sourceArn string) (*KMSDataKeyResult, error)
	Decrypt(ctx context.Context, keyID string, ciphertext []byte, encryptionContext map[string]string, sourceArn string) ([]byte, error)
	KeyExists(ctx context.Context, keyID string) bool
}

// KMSDataKeyResult carries the plaintext and encrypted data key returned by
// GenerateDataKey.
type KMSDataKeyResult struct {
	Plaintext      []byte
	CiphertextBlob []byte
}

// IAMPrincipalResolver resolves an access key ID to a username for audit
// logging. Consumers call this instead of holding a direct reference to the
// IAM store.
type IAMPrincipalResolver interface {
	ResolvePrincipal(ctx context.Context, accessKeyID string) (username string, err error)
}

// S3Invoker provides S3 object read and write operations for cross-service
// consumers (e.g. NeptuneData bulk loader, DynamoDB import/export). Consumers
// call these methods instead of holding a direct reference to the S3 store.
type S3Invoker interface {
	GetObject(ctx context.Context, region, bucket, key string, maxBytes int64) ([]byte, error)
	PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error
	ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error)
}

// WAFInvoker provides WAF WebACL association operations for cross-service
// consumers (e.g. CloudFront distribution management). Consumers call these
// methods instead of holding a direct reference to the WAF store.
type WAFInvoker interface {
	AssociateWebACL(webACLArn, resourceArn string) error
	DisassociateWebACL(webACLArn, resourceArn string) error
}

// CloudWatchMetricInvoker provides CloudWatch metric data operations for
// cross-service consumers (e.g. CloudWatch Logs metric filter evaluation).
// Consumers call these methods instead of holding a direct reference to the
// CloudWatch metric store.
type CloudWatchMetricInvoker interface {
	PutMetricData(region, namespace string, metricName string, value float64, timestamp time.Time) error
}

// CloudWatchAlarmInvoker provides CloudWatch alarm state operations for
// cross-service consumers (e.g. IoT rule cloudwatchAlarm action).
type CloudWatchAlarmInvoker interface {
	SetAlarmState(region, alarmName, stateValue, stateReason string) error
}

// TimestreamInvoker provides Timestream WriteRecords operations for
// cross-service consumers (e.g. IoT rule timestream action).
type TimestreamInvoker interface {
	WriteRecords(region, databaseName, tableName string, dimensions map[string]string, measureName string, measureValue string, measureType string, timestamp time.Time) error
}

// CloudTrailEventInfo carries the fields of a CloudTrail event that
// cross-service consumers need for last-accessed analysis.
type CloudTrailEventInfo struct {
	EventID     string
	EventName   string
	EventSource string
	EventTime   time.Time
	Username    string
}

// CloudTrailInvoker provides CloudTrail event lookup for cross-service
// consumers (e.g. IAM GenerateServiceLastAccessedDetails). Consumers call
// these methods instead of holding a direct reference to the CloudTrail store.
//
// The CloudTrail store is scoped to a single account at construction time
// (NewCloudTrailStore receives the account ID), so LookupEvents does not
// take an accountID parameter: the caller's account is implicit in the
// store resolved by region.
//
// nextToken supports pagination: pass the token returned by the previous
// call to fetch the next page. Per AWS LookupEvents spec, maxResults must
// be 1-50.
type CloudTrailInvoker interface {
	LookupEvents(ctx context.Context, region, username, nextToken string, startTime, endTime time.Time, maxResults int32) ([]CloudTrailEventInfo, string, error)
}

// LogsLogEntry carries a single log entry for cross-service delivery.
type LogsLogEntry struct {
	Timestamp int64
	Message   string
}

// LogsInvoker provides CloudWatch Logs write operations for cross-service
// consumers (e.g. Lambda function log delivery). Consumers call these methods
// instead of holding a direct reference to the Logs store.
type LogsInvoker interface {
	EnsureLogGroup(ctx context.Context, region, logGroupName, accountID string) error
	EnsureLogStream(ctx context.Context, region, logGroupName, logStreamName string) error
	PutLogEvents(ctx context.Context, region, logGroupName, logStreamName string, entries []LogsLogEntry) error
}

// RDSDataInvoker provides RDS Data API SQL execution for cross-service consumers
// (e.g. AppSync RELATIONAL_DATABASE resolvers). Consumers call these methods
// instead of holding a direct reference to the RDS Data service.
type RDSDataInvoker interface {
	ExecuteStatement(ctx context.Context, resourceArn, secretArn, database, schema, sql string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error)
	BatchExecuteStatement(ctx context.Context, resourceArn, secretArn, database, schema, sql string, parameterSets [][]interface{}) (interface{}, error)
	BeginTransaction(ctx context.Context, resourceArn, secretArn, database, schema string) (string, error)
	CommitTransaction(ctx context.Context, resourceArn, secretArn, transactionId string) error
	RollbackTransaction(ctx context.Context, resourceArn, secretArn, transactionId string) error
}

// ACMInvoker provides ACM certificate usage tracking for cross-service
// consumers (e.g. CloudFront distribution management, API Gateway domain
// management). When a service associates an ACM certificate with a resource,
// it must register the resource ARN so that DeleteCertificate can enforce
// the ResourceInUseError guard. Consumers call these methods instead of
// holding a direct reference to the ACM store.
type ACMInvoker interface {
	RegisterCertificateUsage(ctx context.Context, region, certArn, resourceArn string) error
	UnregisterCertificateUsage(ctx context.Context, region, certArn, resourceArn string) error
	CertificateExists(ctx context.Context, region, certArn string) bool
}

// CognitoTokenValidator validates Cognito JWT access tokens for cross-service
// consumers (e.g. API Gateway COGNITO_USER_POOLS authorizer). The validator
// resolves the user pool by ID, fetches its JWKS public key, and verifies the
// token signature, expiration, issuer, and token_use claim.
type CognitoTokenValidator interface {
	ValidateTokenForPool(ctx context.Context, region, userPoolID, accessToken string) (subject string, err error)
}
