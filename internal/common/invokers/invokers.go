// Package invokers holds the cross-service invocation contracts: the
// interfaces and data types through which one AWS service consumes
// another service's operations without importing that service
// directly. Implementations are wired by the server application at
// startup and reached through the event bus's invoker registry.
package invokers

import (
	"context"
	"time"
)

// LambdaInvoker invokes a Lambda function and returns the status code and
// response payload. It extends the minimal lambdautil.Invoker gateway
// contract (InvokeForGateway alone) with the trigger invocation and ARN
// resolution that event-delivery consumers need; the two contracts are
// kept separate because their consumers do not overlap.
type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (statusCode int64, responsePayload []byte, err error)
	InvokeForTrigger(ctx context.Context, functionName string, payload []byte) (LambdaInvocation, error)
	GetFunctionARN(ctx context.Context, functionName string) (string, error)
}

// LambdaInvocation carries the outcome of a direct Lambda invocation for
// synchronous trigger consumers. A failed function execution is reported
// through FunctionError (the invoke transport succeeds even when the
// function itself raises an error), while an invocation-transport failure
// is reported through err.
type LambdaInvocation struct {
	StatusCode    int64
	Payload       []byte
	FunctionError string
}

// AppSyncInvoker executes a GraphQL mutation on an AppSync API for
// cross-service event delivery (EventBridge targets). The operation is the
// mutation document (AppSyncParameters.GraphQLOperation); variablesJSON is
// the transformed event payload as a JSON object — the shapes AWS
// documents for AppSync targets (a valid GraphQL mutation invoked for
// matched events with AWS_IAM authorisation, which the internal bus call
// plays on this platform).
type AppSyncInvoker interface {
	ExecuteGraphQLMutation(ctx context.Context, region, apiID, operation string, variablesJSON []byte) error
}

// SQSInvoker provides SQS operations for cross-service consumers.
// All methods require a region parameter to target the correct regional
// SQS store, enabling cross-region delivery (e.g. alarm actions targeting
// queues in a different region than the source service).
type SQSInvoker interface {
	GetQueueByName(ctx context.Context, region, queueName string) (queueURL string, err error)
	GetQueueARN(ctx context.Context, region, queueURL string) (queueARN string, err error)
	SendMessage(ctx context.Context, region, queueURL, body string, opts SQSSendOptions) (messageID string, md5OfBody string, err error)
	ReceiveMessage(ctx context.Context, region, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]ReceivedSQSMessage, error)
	DeleteMessage(ctx context.Context, region, queueURL, receiptHandle string) error
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
	// GetTopicPolicy returns the topic's access-policy document JSON, or an
	// empty string when the topic carries no policy. Consumers that must
	// verify the topic grants a service principal publish access before
	// accepting it as a notification destination (CloudTrail's
	// InsufficientSnsTopicPolicyException) evaluate the document this
	// returns.
	GetTopicPolicy(ctx context.Context, topicARN string) (string, error)
	ListSubscriptionsByTopic(ctx context.Context, topicARN string) ([]SubscriptionInfo, error)
	// PublishToTopic publishes through the service's publish path, so
	// cross-service publishes carry the same validation and delivery as the
	// Publish API. messageAttributes carries typed attributes: the
	// SQSMessageAttribute type is the invokers' protocol-neutral
	// typed-attribute carrier (DataType plus StringValue or the decoded
	// BinaryValue), named for its first consumer.
	PublishToTopic(ctx context.Context, topicARN string, message string, subject string, messageAttributes map[string]SQSMessageAttribute) (messageID string, err error)
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
	// ListShards lists the shards of the stream in the given region —
	// callers address a stream by its ARN and must pass the ARN's region so
	// cross-region reads reach the correct regional store. An empty region
	// addresses the server default region.
	ListShards(ctx context.Context, region, streamName string) ([]ShardInfo, error)
	// PutRecord writes to the stream in the given region — callers address a
	// stream by its ARN and must pass the ARN's region so cross-region
	// delivery reaches the correct regional store. An empty region addresses
	// the server default region (callers with no region information at
	// all). The return is the record's receipt: the sequence number and the
	// shard the record landed on — the model's PutRecordOutput members,
	// which consumers surface in their own responses.
	PutRecord(ctx context.Context, region, streamName, partitionKey string, data []byte) (sequenceNumber, shardID string, err error)
	// CreateShardIterator creates a shard iterator for the stream in the
	// given region (empty = the server default region). The timestamp is used
	// when iteratorType is AT_TIMESTAMP and may be nil otherwise.
	CreateShardIterator(ctx context.Context, region, streamName string, shardID string, iteratorType string, startingSequenceNumber string, timestamp *time.Time) (iteratorSequenceNumber string, err error)
	// GetRecords reads records from the stream in the given region (empty =
	// the server default region), from the position onwards. includeStart
	// also returns the starting record itself, mirroring the store's
	// AT_SEQUENCE_NUMBER semantics; checkpoint resumes stay strictly
	// after their position so consumed records are never redelivered.
	GetRecords(ctx context.Context, region, streamName string, shardID string, startingSequenceNumber string, limit int32, includeStart bool) (records []KinesisRecord, nextSequenceNumber string, err error)
	// StreamExists reports whether the stream addressed by the ARN exists in
	// the given region. DynamoDB streaming destinations may only target a
	// stream in the table's own region, so the ARN region must match and the
	// stream must exist there.
	StreamExists(ctx context.Context, region, streamARN string) (bool, error)
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

// EC2Invoker provides EC2 subnet, security group, and VPC lookup operations
// for cross-service consumers (e.g. Lambda VpcConfig validation).
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
	ScanWithPagination(ctx context.Context, region, tableName string, limit int, exclusiveStartKey string) ([]map[string]interface{}, string, error)
	QueryWithPagination(ctx context.Context, region, tableName, partitionKeyValue string, limit int, exclusiveStartKey string) ([]map[string]interface{}, string, error)
	// ContributorRules lists the DynamoDB contributor insights rule names
	// derived from the insights-enabled tables of one region. CloudWatch
	// merges them into DescribeInsightRules.
	ContributorRules(ctx context.Context, region string) ([]ContributorInsightRule, error)
	// ContributorStats returns the most accessed tracked keys of one table
	// inside the half-open time window, for the GetInsightRuleReport path.
	ContributorStats(ctx context.Context, region, tableName, layout string, start, end time.Time, max int) ([]ContributorKeyStat, error)
}

// ContributorInsightRule is one DynamoDB-derived contributor insights
// rule exposed through the eventbus.
type ContributorInsightRule struct {
	Name string
}

// ContributorKeyStat is one aggregated key of a DynamoDB contributor
// insights rule. Keys holds the decoded key attribute values.
type ContributorKeyStat struct {
	Keys  []string
	Count int64
	Units float64
}

// DynamoDBStreamsInvoker provides DynamoDB Streams operations for
// cross-service consumers (e.g. Lambda ESM poller). Consumers call these
// methods instead of holding a direct reference to the DynamoDB stream store.
type DynamoDBStreamsInvoker interface {
	GetRecords(ctx context.Context, region, tableName string, fromSeq int64, limit int) ([]DynamoDBStreamRecord, int64, error)
	GetLatestSequence(ctx context.Context, region, tableName string) (int64, error)
	// ShardIDForStream returns the deterministic shard identifier of a
	// DynamoDB stream ARN — the same value the DynamoDB Streams API reports
	// for the stream. Consumers that need to identify the shard of a stream
	// (for example Lambda's tumbling-window event envelope) use it.
	ShardIDForStream(streamARN string) string
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
	// SymmetricEncryptionKeyExists reports whether the key exists and is
	// a symmetric encryption key (KeySpec SYMMETRIC_DEFAULT, KeyUsage
	// ENCRYPT_DECRYPT) — the key class consumers such as EventBridge
	// Scheduler validate at configuration time via kms:DescribeKey.
	SymmetricEncryptionKeyExists(ctx context.Context, keyID string) bool
}

// KMSDataKeyResult carries the plaintext and encrypted data key returned by
// GenerateDataKey.
type KMSDataKeyResult struct {
	Plaintext      []byte
	CiphertextBlob []byte
}

// SecretsManagerInvoker provides secret lifecycle operations for
// cross-service consumers (e.g. EventBridge connection credentials). AWS
// stores a service-owned secret per connection ("When you create a
// connection and add authorization parameters, EventBridge creates a secret
// in AWS Secrets Manager", user guide) and surfaces it as the connection's
// SecretArn; this contract lets the events service follow the same design
// without importing the Secrets Manager store.
type SecretsManagerInvoker interface {
	// CreateServiceSecret creates a secret holding secretString and returns
	// its ARN. The name is the caller's service-qualified name.
	CreateServiceSecret(ctx context.Context, region, name, secretString, description string) (arn string, err error)
	// GetServiceSecretString returns the current secret string of the secret
	// addressed by ARN or name.
	GetServiceSecretString(ctx context.Context, region, secretId string) (string, error)
	// UpdateServiceSecretString replaces the secret string of the secret
	// addressed by ARN or name.
	UpdateServiceSecretString(ctx context.Context, region, secretId, secretString string) error
	// DeleteServiceSecret immediately deletes the secret addressed by ARN or
	// name (no recovery window — the caller owns the resource lifecycle).
	DeleteServiceSecret(ctx context.Context, region, secretId string) error
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
	// GetObjectVersion retrieves the content of a specific object
	// version; an empty versionID reads the latest version, matching the
	// S3 store's version-aware read.
	GetObjectVersion(ctx context.Context, region, bucket, key, versionID string, maxBytes int64) ([]byte, error)
	PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error
	// PutObjectWithMetadata stores the object carrying S3 object metadata
	// (x-amz-meta-*). Consumers delivering signed artifacts (CloudTrail
	// digest files carry their signature as object metadata) use this
	// instead of PutObject.
	PutObjectWithMetadata(ctx context.Context, region, bucket, key string, data []byte, contentType string, metadata map[string]string) error
	ListObjects(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]string, error)
	// ListObjectEntries lists object metadata under a prefix, mirroring the
	// ListObjectsV2 item shape Step Functions exposes to Distributed Map
	// ItemReaders (Etag, Key, LastModified, Size, StorageClass). Consumers
	// that need the full metadata records call this instead of ListObjects,
	// which returns keys only.
	ListObjectEntries(ctx context.Context, region, bucket, prefix string, maxKeys int) ([]S3ObjectEntry, error)
	// BucketExists reports whether the named bucket exists in the region.
	// Consumers use it to distinguish a missing source bucket from an
	// empty one, which listing alone cannot do.
	BucketExists(ctx context.Context, region, bucket string) (bool, error)
	// GetBucketPolicy returns the bucket's policy document JSON, or an
	// empty string when the bucket carries no policy. Consumers that must
	// verify the bucket grants a service principal write access before
	// accepting it as a destination (CloudTrail's InsufficientS3Bucket
	// PolicyException) evaluate the document this returns.
	GetBucketPolicy(ctx context.Context, region, bucket string) (string, error)
	// EnsureBucket creates the named bucket when it does not exist yet.
	// Consumers that own an internal service bucket (e.g. the Cognito
	// user-import upload bucket) use it instead of requiring operators to
	// provision the bucket manually.
	EnsureBucket(ctx context.Context, region, bucket string) error
	// DeleteObject removes an object. Consumers use it to purge transient
	// payloads after use, mirroring AWS services that delete uploaded
	// import files once the job completes.
	DeleteObject(ctx context.Context, region, bucket, key string) error
}

// S3ObjectEntry is one object's ListObjectsV2 metadata as surfaced to
// cross-service consumers (the Step Functions ItemReader listObjectsV2
// dataset item shape).
type S3ObjectEntry struct {
	Key          string
	ETag         string
	LastModified int64
	Size         int64
	StorageClass string
}

// WAFInvoker provides WAF WebACL association operations for cross-service
// consumers (e.g. CloudFront distribution management). Consumers call these
// methods instead of holding a direct reference to the WAF store.
type WAFInvoker interface {
	AssociateWebACL(webACLArn, resourceArn string) error
	DisassociateWebACL(webACLArn, resourceArn string) error
	// WebACLExists reports whether a Web ACL with the given ARN or ID
	// exists, so cross-service consumers can reject references to
	// non-existent Web ACLs before persisting them.
	WebACLExists(ctx context.Context, webACLIdOrArn string) bool
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

// PutLogEvents batch bounds from the operation's reference: "The maximum
// batch size is 1,048,576 bytes. This size is calculated as the sum of all
// event messages in UTF-8, plus 26 bytes for each log event", "The maximum
// number of log events in a batch is 10,000", and "Each log event can be
// no larger than 1 MB". A caller with more events than one batch holds
// splits the delivery into multiple calls; the events within one batch are
// in chronological order by timestamp.
const (
	PutLogEventsBatchMaxEvents = 10000
	PutLogEventsBatchMaxBytes  = 1048576
	PutLogEventsEventOverhead  = 26
)

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

// TLSCertificateMaterial is the PEM certificate material a certificate
// provider returns so another service's listener can terminate TLS: the
// leaf certificate, its private key, and the optional chain.
type TLSCertificateMaterial struct {
	Certificate      string
	PrivateKey       string
	CertificateChain string
}

// ACMCertificateProvider resolves the PEM material of an issued ACM
// certificate for TLS termination on a cross-service listener (e.g. the
// CloudFront distribution plane). Consumers call this instead of holding a
// direct reference to the ACM store.
type ACMCertificateProvider interface {
	CertificateMaterial(ctx context.Context, region, certArn string) (TLSCertificateMaterial, error)
}

// IAMServerCertificateProvider resolves the PEM material of an IAM server
// certificate for TLS termination on a cross-service listener. Consumers
// call this instead of holding a direct reference to the IAM store.
type IAMServerCertificateProvider interface {
	ServerCertificateMaterial(ctx context.Context, serverCertificateId string) (TLSCertificateMaterial, error)
}

// CognitoTokenValidator validates Cognito JWT access tokens for cross-service
// consumers (e.g. API Gateway COGNITO_USER_POOLS authorizer). The validator
// resolves the user pool by ID, fetches its JWKS public key, and verifies the
// token signature, expiration, issuer, and token_use claim.
type CognitoTokenValidator interface {
	ValidateTokenForPool(ctx context.Context, region, userPoolID, accessToken string) (subject string, err error)
}

// CognitoIDTokenClaimResolver validates a Cognito user-pool ID token for a
// specific pool and returns the token's string claims — the claim set
// identity-pool role mappings resolve against (cognito:roles,
// cognito:preferred_role and the custom: attributes Rules match). Consumers
// call this instead of parsing the token themselves: signature, issuer,
// expiry and token_use validation stay with the token-issuing service.
type CognitoIDTokenClaimResolver interface {
	IDTokenClaimsForPool(ctx context.Context, region, userPoolID, idToken string) (claims map[string]string, err error)
}
