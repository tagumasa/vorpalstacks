package apps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	dynamodbsvc "vorpalstacks/internal/services/aws/dynamodb"
	"vorpalstacks/internal/services/aws/rds/rdsdata"
	svcwafv2 "vorpalstacks/internal/services/aws/wafv2"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	storecommon "vorpalstacks/internal/store/aws/common"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	storekinesis "vorpalstacks/internal/store/aws/kinesis"
	storesns "vorpalstacks/internal/store/aws/sns"
	storesqs "vorpalstacks/internal/store/aws/sqs"
	timestreamstore "vorpalstacks/internal/store/aws/timestream"
	wafstore "vorpalstacks/internal/store/aws/waf"
	"vorpalstacks/internal/utils/aws/arn"
)

// sqsStoreProvider resolves the per-region SQS store owned by the SQS
// service. The adapter must not construct its own store: the store keeps
// per-instance deduplication, sequence, and receive-attempt caches, so a
// second instance would diverge from the API plane.
type sqsStoreProvider interface {
	GetStoreForRegion(region string) (storesqs.SQSStoreInterface, error)
}

// sqsInvokerAdapter adapts the SQS service store to the eventbus.SQSInvoker
// interface. This enables cross-region SQS delivery (e.g. alarm actions
// targeting queues in a different region than the source service).
type sqsInvokerAdapter struct {
	provider sqsStoreProvider
}

func (a *sqsInvokerAdapter) getStore(region string) (storesqs.SQSStoreInterface, error) {
	if region == "" {
		region = defaults.DefaultRegion
	}
	return a.provider.GetStoreForRegion(region)
}

// GetQueueByName looks up a queue by name in the specified region and returns its URL.
func (a *sqsInvokerAdapter) GetQueueByName(_ context.Context, region, queueName string) (string, error) {
	store, err := a.getStore(region)
	if err != nil {
		return "", err
	}
	q, err := store.GetQueueByName(queueName)
	if err != nil {
		return "", err
	}
	return q.URL, nil
}

// GetQueueARN looks up a queue by URL in the specified region and returns its ARN.
func (a *sqsInvokerAdapter) GetQueueARN(_ context.Context, region, queueURL string) (string, error) {
	store, err := a.getStore(region)
	if err != nil {
		return "", err
	}
	q, err := store.GetQueue(queueURL)
	if err != nil {
		return "", err
	}
	return q.ARN, nil
}

// SendMessage sends a message to the specified queue in the given region.
func (a *sqsInvokerAdapter) SendMessage(_ context.Context, region, queueURL, body string, opts eventbus.SQSSendOptions) (string, string, error) {
	store, err := a.getStore(region)
	if err != nil {
		return "", "", err
	}
	msg := &storesqs.Message{
		Body:                   body,
		DelaySeconds:           int32(opts.DelaySeconds),
		MessageAttributes:      buildSQSMessageAttributes(opts),
		MessageGroupID:         opts.MessageGroupID,
		MessageDeduplicationID: opts.MessageDeduplicationID,
	}
	result, err := store.SendMessage(queueURL, msg)
	if err != nil {
		return "", "", err
	}
	return result.ID, result.MD5OfBody, nil
}

// buildSQSMessageAttributes selects the richer of the two attribute maps
// from SQSSendOptions. TypedMessageAttributes (with DataType) takes
// precedence; the legacy string map is used as a fallback.
func buildSQSMessageAttributes(opts eventbus.SQSSendOptions) map[string]*storesqs.MessageAttributeValue {
	if len(opts.TypedMessageAttributes) > 0 {
		out := make(map[string]*storesqs.MessageAttributeValue, len(opts.TypedMessageAttributes))
		for k, v := range opts.TypedMessageAttributes {
			attr := &storesqs.MessageAttributeValue{DataType: v.DataType}
			if len(v.BinaryValue) > 0 {
				attr.BinaryValue = v.BinaryValue
			} else {
				attr.StringValue = &v.StringValue
			}
			out[k] = attr
		}
		return out
	}
	return convertToSQSMessageAttributes(opts.MessageAttributes)
}

// ReceiveMessage retrieves messages from the specified queue in the given region.
func (a *sqsInvokerAdapter) ReceiveMessage(_ context.Context, region, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]eventbus.ReceivedSQSMessage, error) {
	store, err := a.getStore(region)
	if err != nil {
		return nil, err
	}
	msgs, err := store.ReceiveMessage(queueURL, maxMessages, visibilityTimeout, waitTimeSeconds, "")
	if err != nil {
		return nil, err
	}
	out := make([]eventbus.ReceivedSQSMessage, len(msgs))
	for i, m := range msgs {
		out[i] = eventbus.ReceivedSQSMessage{
			MessageID:                        m.ID,
			ReceiptHandle:                    m.ReceiptHandle,
			Body:                             m.Body,
			MD5OfBody:                        m.MD5OfBody,
			MD5OfMessageAttributes:           m.MD5OfMessageAttributes,
			SentTimestamp:                    m.SentTimestamp,
			ApproximateReceiveCount:          m.ApproximateReceiveCount,
			ApproximateFirstReceiveTimestamp: m.ApproximateFirstReceiveTimestamp,
			SequenceNumber:                   m.SequenceNumber,
			MessageDeduplicationID:           m.MessageDeduplicationID,
			MessageGroupID:                   m.MessageGroupID,
		}
		if m.MessageAttributes != nil {
			out[i].MessageAttributes = convertFromSQSMessageAttributes(m.MessageAttributes)
		}
	}
	return out, nil
}

// DeleteMessage deletes a message from the specified queue in the given region.
func (a *sqsInvokerAdapter) DeleteMessage(_ context.Context, region, queueURL, receiptHandle string) error {
	store, err := a.getStore(region)
	if err != nil {
		return err
	}
	return store.DeleteMessage(queueURL, receiptHandle)
}

// snsStoreForInvoker is the minimal store interface needed by the SNS invoker
// adapter. It intentionally includes Put (from BaseStore) because the event
// bus needs to persist delivery metadata — this does not belong on the public
// SNSStoreInterface.
type snsStoreForInvoker interface {
	GetTopic(topicArn string) (*storesns.Topic, error)
	ListSubscriptionsByTopic(topicArn string, opts storecommon.ListOptions) (*storecommon.ListResult[storesns.Subscription], error)
	Put(key string, data interface{}) error
}

// snsInvokerAdapter adapts the SNS concrete store and publisher to the
// eventbus.SNSInvoker interface.
type snsInvokerAdapter struct {
	store     snsStoreForInvoker
	kvStore   kvDeleter
	publisher snsPublisher
}

// kvDeleter provides raw key-value deletion for message cleanup.
type kvDeleter interface {
	Delete(key string) error
}

// snsPublisher publishes a message to an SNS topic by ARN and returns the
// generated message ID.
type snsPublisher interface {
	PublishToTopic(ctx context.Context, accountID, region, topicArn, message, subject string, messageAttributes map[string]string) (string, error)
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
	msgID, err := a.publisher.PublishToTopic(ctx, accountID, region, topicARN, message, subject, messageAttributes)
	if err != nil {
		return "", err
	}
	return msgID, nil
}

// StoreMessage persists arbitrary data keyed by the given key.
func (a *snsInvokerAdapter) StoreMessage(_ context.Context, key string, data any) error {
	return a.store.Put(key, data)
}

// DeleteStoredMessage removes a previously stored message by key, used for
// cleanup when delivery fails after persistence.
func (a *snsInvokerAdapter) DeleteStoredMessage(_ context.Context, key string) error {
	return a.kvStore.Delete(key)
}

// kinesisStoreProvider resolves the per-region Kinesis store owned by the
// Kinesis service. The adapter must not construct its own store: the store
// keeps in-memory sequence and shard-id counters, so a second instance
// would diverge from the API plane.
type kinesisStoreProvider interface {
	GetStoreForRegion(region string) (*storekinesis.KinesisStore, error)
}

// kinesisInvokerAdapter adapts the Kinesis service store to the
// eventbus.KinesisInvoker interface.
type kinesisInvokerAdapter struct {
	provider kinesisStoreProvider
	// defaultRegion is the region for interface methods that do not carry
	// a region parameter.
	defaultRegion string
}

// getStore returns the KinesisStore for the region, resolved through the
// owning service so the API plane and cross-service writers share one
// instance per region.
func (a *kinesisInvokerAdapter) getStore(region string) (*storekinesis.KinesisStore, error) {
	if region == "" {
		region = a.defaultRegion
	}
	return a.provider.GetStoreForRegion(region)
}

// defaultStore returns the KinesisStore of the adapter's default region.
func (a *kinesisInvokerAdapter) defaultStore() (*storekinesis.KinesisStore, error) {
	return a.provider.GetStoreForRegion(a.defaultRegion)
}

// StreamExists reports whether the stream addressed by the ARN exists in the
// given region. A DynamoDB streaming destination may only target a stream in
// the table's own region, so the ARN region must match the requested region.
func (a *kinesisInvokerAdapter) StreamExists(_ context.Context, region, streamARN string) (bool, error) {
	parsed, err := arn.ParseARN(streamARN)
	if err != nil {
		return false, nil
	}
	if parsed.Region != region {
		return false, nil
	}
	streamName := strings.TrimPrefix(parsed.Resource, "stream/")
	if streamName == "" || streamName == parsed.Resource {
		return false, nil
	}
	store, err := a.getStore(region)
	if err != nil {
		return false, err
	}
	if _, err := store.GetStream(streamName); err != nil {
		if errors.Is(err, storekinesis.ErrStreamNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListShards lists the shards in the given Kinesis stream.
func (a *kinesisInvokerAdapter) ListShards(_ context.Context, streamName string) ([]eventbus.ShardInfo, error) {
	store, err := a.defaultStore()
	if err != nil {
		return nil, err
	}
	shards, err := store.ListShards(streamName, nil, "", 0)
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
	store, err := a.defaultStore()
	if err != nil {
		return "", err
	}
	shards, err := store.ListShards(streamName, nil, "", 0)
	if err != nil {
		return "", err
	}
	openShard := selectOpenShard(shards)
	if openShard == nil {
		return "", fmt.Errorf("kinesis: no open shard in stream %s", streamName)
	}
	record, err := store.PutRecord(streamName, openShard.ShardID, partitionKey, string(data))
	if err != nil {
		return "", err
	}
	return record.SequenceNumber, nil
}

// CreateShardIterator creates a shard iterator for the given stream and shard.
// The timestamp parameter is honoured for AT_TIMESTAMP iterators.
func (a *kinesisInvokerAdapter) CreateShardIterator(_ context.Context, streamName string, shardID string, iteratorType string, startingSequenceNumber string, timestamp *time.Time) (string, error) {
	store, err := a.defaultStore()
	if err != nil {
		return "", err
	}
	iterator, err := store.CreateShardIterator(streamName, shardID, iteratorType, startingSequenceNumber, timestamp)
	if err != nil {
		return "", err
	}
	return iterator.SequenceNumber, nil
}

// GetRecords retrieves records from a Kinesis shard strictly after the
// given sequence number. The poller resumes from checkpoints and chains
// batch reads, so re-including the boundary record would redeliver it on
// every cycle; the public API keeps the same exclusive semantics for
// every iterator type except AT_SEQUENCE_NUMBER. includeStart re-enables
// the inclusive read for the poller's initial LATEST anchor.
func (a *kinesisInvokerAdapter) GetRecords(_ context.Context, streamName string, shardID string, startingSequenceNumber string, limit int32, includeStart bool) ([]eventbus.KinesisRecord, string, error) {
	store, err := a.defaultStore()
	if err != nil {
		return nil, "", err
	}
	records, nextSeq, err := store.GetRecords(streamName, shardID, startingSequenceNumber, limit, includeStart)
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
// from the list, or nil if all shards are closed.
func selectOpenShard(shards []*storekinesis.Shard) *storekinesis.Shard {
	for _, s := range shards {
		if s.SequenceNumberRange != nil && s.SequenceNumberRange.EndingSequenceNumber == "" {
			return s
		}
	}
	return nil
}

// dynamoDBStoreProvider is a minimal interface for obtaining a DynamoDB store
// by region, satisfied by DynamoDBService.
type dynamoDBStoreProvider interface {
	GetStoreForRegion(region string) (dynamodbstore.DynamoDBStoreInterface, error)
}

// dynamoDBInvokerAdapter adapts the DynamoDB store to the eventbus.DynamoDBInvoker
// interface, so that cross-service consumers (e.g. AppSync GraphQL resolvers)
// perform item operations through the bus instead of holding a direct store reference.
// Writes go through the store transaction so PITR journaling and
// contributor-insights accounting commit together with the item, exactly like
// direct data-plane writes.
type dynamoDBInvokerAdapter struct {
	provider dynamoDBStoreProvider
}

func (a *dynamoDBInvokerAdapter) store(ctx context.Context, region string) (dynamodbstore.DynamoDBStoreInterface, error) {
	return a.provider.GetStoreForRegion(region)
}

// recordContributorReads credits one read event per returned key, mirroring
// how the data plane counts every read item. Monitoring must never fail the
// read it observes, so failures are logged and dropped.
func (a *dynamoDBInvokerAdapter) recordContributorReads(ctx context.Context, s dynamodbstore.DynamoDBStoreInterface, tableName string, keys []map[string]*dynamodbstore.AttributeValue) {
	if len(keys) == 0 {
		return
	}
	if err := s.RecordContributorReads(ctx, tableName, keys); err != nil {
		logs.Warn("failed to record contributor reads",
			logs.String("table", tableName), logs.Err(err))
	}
}

// recordContributorQuery credits a query as a single read event on the
// partition-key series, regardless of how many items the query returned.
// The partition value is typed from the table's attribute definitions so
// the event lands on the same key series as item writes. Failures are
// logged and dropped like read events.
func (a *dynamoDBInvokerAdapter) recordContributorQuery(ctx context.Context, s dynamodbstore.DynamoDBStoreInterface, table *dynamodbstore.Table, partitionKeyValue string) {
	if !table.ContributorInsightsEnabled {
		return
	}
	pkName := ""
	for _, ks := range table.KeySchema {
		if ks.KeyType == dynamodbstore.KeyTypeHash {
			pkName = ks.AttributeName
			break
		}
	}
	if pkName == "" {
		return
	}
	pkType := dynamodbstore.ScalarAttributeTypeS
	for _, def := range table.AttributeDefinitions {
		if def.AttributeName == pkName {
			pkType = def.AttributeType
			break
		}
	}
	value := &dynamodbstore.AttributeValue{}
	switch pkType {
	case dynamodbstore.ScalarAttributeTypeN:
		value.N = &partitionKeyValue
	case dynamodbstore.ScalarAttributeTypeB:
		value.B = []byte(partitionKeyValue)
	default:
		value.S = &partitionKeyValue
	}
	if err := s.RecordContributorQuery(ctx, table.Name, map[string]*dynamodbstore.AttributeValue{pkName: value}); err != nil {
		logs.Warn("failed to record contributor query event",
			logs.String("table", table.Name), logs.Err(err))
	}
}

// GetItem retrieves a single item from DynamoDB by key.
func (a *dynamoDBInvokerAdapter) GetItem(ctx context.Context, region, tableName string, key map[string]interface{}) (map[string]interface{}, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, err
	}
	dynamoKey := dynamoMapToKey(key)
	item, err := s.Items().Get(tableName, dynamoKey)
	if err != nil {
		return nil, err
	}
	a.recordContributorReads(ctx, s, tableName, []map[string]*dynamodbstore.AttributeValue{dynamoKey})
	return dynamoItemToPlainMap(item), nil
}

// PutItem creates or replaces an item in DynamoDB.
func (a *dynamoDBInvokerAdapter) PutItem(ctx context.Context, region, tableName string, key, attributes map[string]interface{}) (map[string]interface{}, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, err
	}
	dynamoKey := dynamoMapToKey(key)
	dynamoAttrs := dynamoMapToAttrs(attributes)
	if err := s.Update(ctx, func(txn *dynamodbstore.DynamoDBTxn) error {
		return txn.PutItem(tableName, dynamoKey, dynamoAttrs)
	}); err != nil {
		return nil, err
	}
	result := make(map[string]interface{}, len(key)+len(attributes))
	for k, v := range attributes {
		result[k] = v
	}
	for k, v := range key {
		result[k] = v
	}
	return result, nil
}

// DeleteItem removes an item from DynamoDB by key.
func (a *dynamoDBInvokerAdapter) DeleteItem(ctx context.Context, region, tableName string, key map[string]interface{}) error {
	s, err := a.store(ctx, region)
	if err != nil {
		return err
	}
	dynamoKey := dynamoMapToKey(key)
	return s.Update(ctx, func(txn *dynamodbstore.DynamoDBTxn) error {
		return txn.DeleteItem(tableName, dynamoKey)
	})
}

var errScanLimitReached = fmt.Errorf("scan limit reached")

// Scan scans all items in a DynamoDB table up to the given limit.
func (a *dynamoDBInvokerAdapter) Scan(ctx context.Context, region, tableName string, limit int) ([]map[string]interface{}, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 1000
	}
	var results []map[string]interface{}
	readKeys := make([]map[string]*dynamodbstore.AttributeValue, 0, limit)
	count := 0
	scanErr := s.Items().Scan(tableName, func(item *dynamodbstore.Item) error {
		if count >= limit {
			return errScanLimitReached
		}
		results = append(results, dynamoItemToPlainMap(item))
		readKeys = append(readKeys, item.Key)
		count++
		return nil
	})
	if scanErr != nil && scanErr != errScanLimitReached {
		return nil, scanErr
	}
	a.recordContributorReads(ctx, s, tableName, readKeys)
	return results, nil
}

// Query retrieves items from a DynamoDB table by partition key value.
func (a *dynamoDBInvokerAdapter) Query(ctx context.Context, region, tableName, partitionKeyValue string, limit int) ([]map[string]interface{}, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, err
	}
	table, err := s.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 1000
	}
	var results []map[string]interface{}
	count := 0
	_, queryErr := s.Items().ScanByPartitionKeyWithTable(tableName, table, partitionKeyValue, dynamodbstore.ScanOptions{}, func(item *dynamodbstore.Item) error {
		if count >= limit {
			return errScanLimitReached
		}
		results = append(results, dynamoItemToPlainMap(item))
		count++
		return nil
	})
	if queryErr != nil && queryErr != errScanLimitReached {
		return nil, queryErr
	}
	a.recordContributorQuery(ctx, s, table, partitionKeyValue)
	return results, nil
}

// UpdateItem replaces the attributes of an existing item in DynamoDB.
func (a *dynamoDBInvokerAdapter) UpdateItem(ctx context.Context, region, tableName string, key, attributes map[string]interface{}) error {
	s, err := a.store(ctx, region)
	if err != nil {
		return err
	}
	dynamoKey := dynamoMapToKey(key)
	dynamoAttrs := dynamoMapToAttrs(attributes)
	return s.Update(ctx, func(txn *dynamodbstore.DynamoDBTxn) error {
		return txn.PutItem(tableName, dynamoKey, dynamoAttrs)
	})
}

// ScanWithPagination performs a paginated scan of a DynamoDB table.
// Returns items up to the given limit and a next-marker for subsequent calls.
// An empty marker means no more items remain.
func (a *dynamoDBInvokerAdapter) ScanWithPagination(ctx context.Context, region, tableName string, limit int, exclusiveStartKey string) ([]map[string]interface{}, string, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, "", err
	}
	if limit <= 0 {
		limit = 1000
	}
	opts := dynamodbstore.ScanOptions{Limit: limit, Marker: exclusiveStartKey}
	var results []map[string]interface{}
	readKeys := make([]map[string]*dynamodbstore.AttributeValue, 0, limit)
	nextMarker, scanErr := s.Items().ScanWithOptions(tableName, opts, func(item *dynamodbstore.Item) error {
		results = append(results, dynamoItemToPlainMap(item))
		readKeys = append(readKeys, item.Key)
		return nil
	})
	if scanErr != nil && !errors.Is(scanErr, errScanLimitReached) {
		return nil, "", scanErr
	}
	a.recordContributorReads(ctx, s, tableName, readKeys)
	return results, nextMarker, nil
}

// QueryWithPagination performs a paginated query of a DynamoDB table by partition key.
// Returns items up to the given limit and a next-marker for subsequent calls.
func (a *dynamoDBInvokerAdapter) QueryWithPagination(ctx context.Context, region, tableName, partitionKeyValue string, limit int, exclusiveStartKey string) ([]map[string]interface{}, string, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, "", err
	}
	table, err := s.Tables().Get(tableName)
	if err != nil {
		return nil, "", err
	}
	if limit <= 0 {
		limit = 1000
	}
	opts := dynamodbstore.ScanOptions{Limit: limit, Marker: exclusiveStartKey}
	var results []map[string]interface{}
	nextMarker, queryErr := s.Items().ScanByPartitionKeyWithTable(tableName, table, partitionKeyValue, opts, func(item *dynamodbstore.Item) error {
		results = append(results, dynamoItemToPlainMap(item))
		return nil
	})
	if queryErr != nil && !errors.Is(queryErr, errScanLimitReached) {
		return nil, "", queryErr
	}
	a.recordContributorQuery(ctx, s, table, partitionKeyValue)
	return results, nextMarker, nil
}

// ContributorRules lists the DynamoDB contributor insights rule names
// derived from the insights-enabled tables of one region.
func (a *dynamoDBInvokerAdapter) ContributorRules(ctx context.Context, region string) ([]eventbus.ContributorInsightRule, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, err
	}
	var rules []eventbus.ContributorInsightRule
	marker := ""
	for {
		tables, next, err := s.Tables().List(marker, 0)
		if err != nil {
			return nil, err
		}
		for _, t := range tables {
			for _, name := range dynamodbsvc.ContributorInsightsRuleNames(t) {
				rules = append(rules, eventbus.ContributorInsightRule{Name: name})
			}
		}
		if next == "" {
			break
		}
		marker = next
	}
	return rules, nil
}

// ContributorStats returns the most accessed tracked keys of one table
// inside the half-open time window.
func (a *dynamoDBInvokerAdapter) ContributorStats(ctx context.Context, region, tableName, layout string, start, end time.Time, max int) ([]eventbus.ContributorKeyStat, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, err
	}
	stats, err := s.Contributors().TopKeys(tableName, layout, start, end, max)
	if err != nil {
		return nil, err
	}
	out := make([]eventbus.ContributorKeyStat, 0, len(stats))
	for _, stat := range stats {
		var keys []string
		_ = json.Unmarshal([]byte(stat.Key), &keys)
		// The aggregation encodes each key value with a type prefix to
		// keep same-text values of different types distinct; the report
		// exposes the bare values.
		for i, k := range keys {
			if _, rest, found := strings.Cut(k, ":"); found {
				keys[i] = rest
			}
		}
		out = append(out, eventbus.ContributorKeyStat{
			Keys:  keys,
			Count: stat.Count,
			Units: stat.Units,
		})
	}
	return out, nil
}

func dynamoMapToKey(m map[string]interface{}) map[string]*dynamodbstore.AttributeValue {
	if m == nil {
		return nil
	}
	out := make(map[string]*dynamodbstore.AttributeValue, len(m))
	for k, v := range m {
		out[k] = dynamoInterfaceToAV(v)
	}
	return out
}

func dynamoMapToAttrs(m map[string]interface{}) map[string]*dynamodbstore.AttributeValue {
	return dynamoMapToKey(m)
}

func dynamoInterfaceToAV(v interface{}) *dynamodbstore.AttributeValue {
	if v == nil {
		null := true
		return &dynamodbstore.AttributeValue{NULL: &null}
	}
	switch val := v.(type) {
	case string:
		return &dynamodbstore.AttributeValue{S: &val}
	case float64:
		s := fmt.Sprintf("%g", val)
		return &dynamodbstore.AttributeValue{N: &s}
	case int:
		s := fmt.Sprintf("%d", val)
		return &dynamodbstore.AttributeValue{N: &s}
	case int64:
		s := fmt.Sprintf("%d", val)
		return &dynamodbstore.AttributeValue{N: &s}
	case bool:
		return &dynamodbstore.AttributeValue{BOOL: &val}
	case map[string]interface{}:
		m := make(map[string]*dynamodbstore.AttributeValue)
		for mk, mv := range val {
			m[mk] = dynamoInterfaceToAV(mv)
		}
		return &dynamodbstore.AttributeValue{M: m}
	case []interface{}:
		l := make([]*dynamodbstore.AttributeValue, len(val))
		for i, item := range val {
			l[i] = dynamoInterfaceToAV(item)
		}
		return &dynamodbstore.AttributeValue{L: l}
	case []byte:
		return &dynamodbstore.AttributeValue{B: val}
	default:
		s := fmt.Sprintf("%v", val)
		return &dynamodbstore.AttributeValue{S: &s}
	}
}

func dynamoItemToPlainMap(item *dynamodbstore.Item) map[string]interface{} {
	if item == nil {
		return nil
	}
	result := make(map[string]interface{})
	if item.Key != nil {
		for k, v := range item.Key {
			result[k] = dynamoAVToInterface(v)
		}
	}
	if item.Attributes != nil {
		for k, v := range item.Attributes {
			result[k] = dynamoAVToInterface(v)
		}
	}
	return result
}

func dynamoAVToInterface(av *dynamodbstore.AttributeValue) interface{} {
	if av == nil {
		return nil
	}
	if av.NULL != nil && *av.NULL {
		return nil
	}
	if av.S != nil {
		return *av.S
	}
	if av.N != nil {
		if f, err := strconv.ParseFloat(*av.N, 64); err == nil {
			return f
		}
		return *av.N
	}
	if av.BOOL != nil {
		return *av.BOOL
	}
	if av.M != nil {
		m := make(map[string]interface{})
		for k, v := range av.M {
			m[k] = dynamoAVToInterface(v)
		}
		return m
	}
	if av.L != nil {
		l := make([]interface{}, len(av.L))
		for i, v := range av.L {
			l[i] = dynamoAVToInterface(v)
		}
		return l
	}
	if av.B != nil {
		return string(av.B)
	}
	return nil
}

type neptuneGraphInvokerAdapter struct {
	service interface {
		ExecuteQueryOnGraph(ctx context.Context, graphID string, query string, language string, parameters map[string]interface{}) (interface{}, error)
	}
}

// rdsDataInvokerAdapter adapts the RDS Data API service to the eventbus.RDSDataInvoker
// interface, so that cross-service consumers (e.g. AppSync GraphQL resolvers)
// execute SQL through the bus instead of holding a direct service reference.
type rdsDataInvokerAdapter struct {
	service interface {
		ExecuteStatementForInvoker(ctx context.Context, resourceArn, secretArn, database, schema, sql string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error)
		BatchExecuteStatementForInvoker(ctx context.Context, resourceArn, secretArn, database, schema, sql string, parameterSets [][]rdsdata.SqlParameter) (interface{}, error)
		BeginTransactionForInvoker(ctx context.Context, resourceArn, secretArn, database, schema string) (string, error)
		CommitTransactionForInvoker(ctx context.Context, resourceArn, secretArn, transactionId string) error
		RollbackTransactionForInvoker(ctx context.Context, resourceArn, secretArn, transactionId string) error
	}
}

func (a *rdsDataInvokerAdapter) ExecuteStatement(ctx context.Context, resourceArn, secretArn, database, schema, sql string, includeResultMetadata bool, formatRecordsAs string) (interface{}, error) {
	return a.service.ExecuteStatementForInvoker(ctx, resourceArn, secretArn, database, schema, sql, includeResultMetadata, formatRecordsAs)
}

func (a *rdsDataInvokerAdapter) BatchExecuteStatement(ctx context.Context, resourceArn, secretArn, database, schema, sql string, parameterSets [][]interface{}) (interface{}, error) {
	paramSets := make([][]rdsdata.SqlParameter, len(parameterSets))
	for i, rawSet := range parameterSets {
		params := make([]rdsdata.SqlParameter, len(rawSet))
		for j, raw := range rawSet {
			p, err := rdsdata.SqlParameterFromInterface(raw)
			if err != nil {
				return nil, err
			}
			params[j] = p
		}
		paramSets[i] = params
	}
	return a.service.BatchExecuteStatementForInvoker(ctx, resourceArn, secretArn, database, schema, sql, paramSets)
}

func (a *rdsDataInvokerAdapter) BeginTransaction(ctx context.Context, resourceArn, secretArn, database, schema string) (string, error) {
	return a.service.BeginTransactionForInvoker(ctx, resourceArn, secretArn, database, schema)
}

func (a *rdsDataInvokerAdapter) CommitTransaction(ctx context.Context, resourceArn, secretArn, transactionId string) error {
	return a.service.CommitTransactionForInvoker(ctx, resourceArn, secretArn, transactionId)
}

func (a *rdsDataInvokerAdapter) RollbackTransaction(ctx context.Context, resourceArn, secretArn, transactionId string) error {
	return a.service.RollbackTransactionForInvoker(ctx, resourceArn, secretArn, transactionId)
}

type wafInvokerAdapter struct {
	store *wafstore.WebACLAssociationStore
	// provider resolves Web ACL existence via the WAFv2 service. It is
	// attached after the optional services initialise because cross-service
	// wiring runs before the WAFv2 service exists. When nil (WAFv2 not
	// initialised) existence cannot be verified and callers keep the
	// historical accept-and-associate behaviour.
	provider *svcwafv2.WAFv2Service
}

// AssociateWebACL links a WAF WebACL to a resource.
func (a *wafInvokerAdapter) AssociateWebACL(webACLArn, resourceArn string) error {
	return a.store.Associate(webACLArn, resourceArn)
}

// DisassociateWebACL removes the WAF WebACL association from a resource.
func (a *wafInvokerAdapter) DisassociateWebACL(webACLArn, resourceArn string) error {
	return a.store.Disassociate(resourceArn)
}

// WebACLExists reports whether the referenced Web ACL exists.
func (a *wafInvokerAdapter) WebACLExists(ctx context.Context, webACLIdOrArn string) bool {
	if a.provider == nil {
		return true
	}
	return a.provider.WebACLExistsForInvoker(ctx, webACLIdOrArn)
}

// SetWebACLProvider wires the WAFv2 service used for Web ACL existence
// checks. Cross-service wiring runs before the optional WAFv2 service is
// created, so the provider is attached when initWAFv2 completes.
func (a *wafInvokerAdapter) SetWebACLProvider(p *svcwafv2.WAFv2Service) {
	a.provider = p
}

// cwAlarmStoreProvider resolves the per-region alarm store owned by the
// CloudWatch service. The adapter must not construct its own store: the
// service owns the store-group lifecycle, and a second instance would
// silently diverge the moment the store gains in-memory state.
type cwAlarmStoreProvider interface {
	AlarmStoreForRegion(region string) (*cwstore.AlarmStore, error)
}

type cloudWatchAlarmInvokerAdapter struct {
	provider cwAlarmStoreProvider
}

func (a *cloudWatchAlarmInvokerAdapter) getStore(region string) (*cwstore.AlarmStore, error) {
	if region == "" {
		region = defaults.DefaultRegion
	}
	return a.provider.AlarmStoreForRegion(region)
}

func (a *cloudWatchAlarmInvokerAdapter) SetAlarmState(region, alarmName, stateValue, stateReason string) error {
	store, err := a.getStore(region)
	if err != nil {
		return err
	}
	return store.SetAlarmState(alarmName, stateValue, stateReason, "")
}

// timestreamRecordStoreProvider resolves the per-region record store owned
// by the Timestream Write service. The adapter must not construct its own
// store: the store buffers writes in per-instance memory, so a second
// instance would diverge from the API read plane.
type timestreamRecordStoreProvider interface {
	RecordStoreForRegion(region string) (*timestreamstore.RecordStore, error)
}

type timestreamInvokerAdapter struct {
	provider timestreamRecordStoreProvider
}

func (a *timestreamInvokerAdapter) getStore(region string) (*timestreamstore.RecordStore, error) {
	if region == "" {
		region = defaults.DefaultRegion
	}
	return a.provider.RecordStoreForRegion(region)
}

func (a *timestreamInvokerAdapter) WriteRecords(region, databaseName, tableName string, dimensions map[string]string, measureName string, measureValue string, measureType string, timestamp time.Time) error {
	store, err := a.getStore(region)
	if err != nil {
		return err
	}
	dims := make([]timestreamstore.Dimension, 0, len(dimensions))
	for k, v := range dimensions {
		dims = append(dims, timestreamstore.Dimension{Name: k, Value: v})
	}
	mvt := timestreamstore.MeasureValueTypeVarchar
	if measureType != "" {
		mvt = timestreamstore.MeasureValueType(measureType)
	}
	records := []timestreamstore.Record{{
		Dimensions:       dims,
		MeasureName:      measureName,
		MeasureValue:     measureValue,
		MeasureValueType: mvt,
		Time:             fmt.Sprintf("%d", timestamp.UnixMilli()),
		TimeUnit:         timestreamstore.TimeUnitMilliseconds,
	}}
	_, err = store.WriteRecords(databaseName, tableName, records)
	return err
}

// cwMetricStoreProvider resolves the per-region metric store owned by the
// CloudWatch service. The adapter must not construct its own store: the
// MetricChunkStore tracks chunk files in per-instance memory and runs
// background goroutines, so a second instance would orphan chunks and leak
// goroutines that nothing ever closes.
type cwMetricStoreProvider interface {
	MetricStoreForRegion(region string) (*cwstore.MetricChunkStore, error)
}

type cloudWatchMetricInvokerAdapter struct {
	provider cwMetricStoreProvider
}

func (a *cloudWatchMetricInvokerAdapter) getStore(region string) (*cwstore.MetricChunkStore, error) {
	if region == "" {
		region = defaults.DefaultRegion
	}
	return a.provider.MetricStoreForRegion(region)
}

// PutMetricData writes a single metric datum to CloudWatch in the given region.
func (a *cloudWatchMetricInvokerAdapter) PutMetricData(region, namespace string, metricName string, value float64, timestamp time.Time) error {
	store, err := a.getStore(region)
	if err != nil {
		return err
	}
	datum := cwstore.MetricDatum{
		Namespace:  namespace,
		MetricName: metricName,
		Value:      value,
		Timestamp:  timestamp,
	}
	return store.PutMetricData(namespace, []cwstore.MetricDatum{datum})
}

// cloudTrailStoreProvider resolves the per-region CloudTrail store owned by
// the CloudTrail service. The adapter must not construct its own store: the
// service owns the store lifecycle, and a second instance would silently
// diverge the moment the store gains in-memory state.
type cloudTrailStoreProvider interface {
	GetStoreForRegion(region string) (cloudtrailstore.CloudTrailStoreInterface, error)
}

type cloudTrailInvokerAdapter struct {
	provider cloudTrailStoreProvider
}

func (a *cloudTrailInvokerAdapter) getStore(region string) (cloudtrailstore.CloudTrailStoreInterface, error) {
	if region == "" {
		region = defaults.DefaultRegion
	}
	return a.provider.GetStoreForRegion(region)
}

// LookupEvents queries CloudTrail for events matching the given criteria.
// nextToken supports pagination by forwarding it to EventQuery.NextToken.
func (a *cloudTrailInvokerAdapter) LookupEvents(_ context.Context, region, username, nextToken string, startTime, endTime time.Time, maxResults int32) ([]eventbus.CloudTrailEventInfo, string, error) {
	ctStore, err := a.getStore(region)
	if err != nil {
		return nil, "", err
	}
	query := cloudtrailstore.EventQuery{
		MaxResults: maxResults,
		StartTime:  &startTime,
		EndTime:    &endTime,
		Username:   username,
		NextToken:  nextToken,
	}
	events, nextToken, err := ctStore.LookupEvents(query)
	if err != nil {
		return nil, "", err
	}
	out := make([]eventbus.CloudTrailEventInfo, 0, len(events))
	for _, e := range events {
		if e == nil {
			continue
		}
		username := ""
		if e.UserIdentity != nil {
			username = e.UserIdentity.UserName
		}
		out = append(out, eventbus.CloudTrailEventInfo{
			EventID:     e.EventID,
			EventName:   e.EventName,
			EventSource: e.EventSource,
			EventTime:   e.EventTime,
			Username:    username,
		})
	}
	return out, nextToken, nil
}

// logsStoreProvider resolves the per-region CloudWatch Logs store owned by
// the CloudWatch Logs service. The adapter must not construct its own store:
// every writer and the API read plane must share the one store instance per
// region, or the per-store chunk sequence and in-memory registry state would
// silently diverge between instances.
type logsStoreProvider interface {
	GetStoreForRegion(region string) (*logsstore.Store, error)
}

type logsInvokerAdapter struct {
	provider logsStoreProvider
}

// EnsureLogGroup creates the log group if it does not already exist.
func (a *logsInvokerAdapter) EnsureLogGroup(_ context.Context, region, logGroupName, accountID string) error {
	store, err := a.provider.GetStoreForRegion(region)
	if err != nil {
		return err
	}
	_, err = store.GetLogGroup(logGroupName)
	if err != nil {
		lg := logsstore.NewLogGroup(logGroupName, region, accountID)
		// A concurrent creator winning the race is fine: the group exists,
		// which is all this method has to guarantee.
		if createErr := store.CreateLogGroup(lg); createErr != nil && !errors.Is(createErr, logsstore.ErrLogGroupAlreadyExists) {
			return createErr
		}
	}
	return nil
}

// EnsureLogStream creates the log stream if it does not already exist.
func (a *logsInvokerAdapter) EnsureLogStream(_ context.Context, region, logGroupName, logStreamName string) error {
	store, err := a.provider.GetStoreForRegion(region)
	if err != nil {
		return err
	}
	if _, err := store.GetLogStream(logGroupName, logStreamName); err == nil {
		return nil
	}
	ls := logsstore.NewLogStream(logStreamName, logGroupName)
	if createErr := store.CreateLogStream(ls); createErr != nil && !errors.Is(createErr, logsstore.ErrLogStreamAlreadyExists) {
		return createErr
	}
	return nil
}

// PutLogEvents writes log entries to the specified log stream.
func (a *logsInvokerAdapter) PutLogEvents(_ context.Context, region, logGroupName, logStreamName string, entries []eventbus.LogsLogEntry) error {
	store, err := a.provider.GetStoreForRegion(region)
	if err != nil {
		return err
	}
	storeEvents := make([]logsstore.LogEntry, len(entries))
	for i, e := range entries {
		storeEvents[i] = logsstore.LogEntry{Timestamp: e.Timestamp, Message: e.Message}
	}
	_, err = store.PutLogEvents(logGroupName, logStreamName, storeEvents)
	return err
}

// ExecuteQueryOnGraph runs a graph query (Cypher/Gremlin/openCypher) against the identified graph.
func (a *neptuneGraphInvokerAdapter) ExecuteQueryOnGraph(ctx context.Context, graphID string, query string, language string, parameters map[string]interface{}) (interface{}, error) {
	return a.service.ExecuteQueryOnGraph(ctx, graphID, query, language, parameters)
}
