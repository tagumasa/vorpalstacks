package apps

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	storecommon "vorpalstacks/internal/store/aws/common"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	storekinesis "vorpalstacks/internal/store/aws/kinesis"
	storesns "vorpalstacks/internal/store/aws/sns"
	storesqs "vorpalstacks/internal/store/aws/sqs"
	wafstore "vorpalstacks/internal/store/aws/waf"
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
func (a *sqsInvokerAdapter) SendMessage(_ context.Context, queueURL string, body string, opts eventbus.SQSSendOptions) (string, string, error) {
	msg := &storesqs.Message{
		Body:                   body,
		DelaySeconds:           int32(opts.DelaySeconds),
		MessageAttributes:      buildSQSMessageAttributes(opts),
		MessageGroupID:         opts.MessageGroupID,
		MessageDeduplicationID: opts.MessageDeduplicationID,
	}
	result, err := a.store.SendMessage(queueURL, msg)
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

// ReceiveMessage retrieves messages from the specified queue.
func (a *sqsInvokerAdapter) ReceiveMessage(_ context.Context, queueURL string, maxMessages int32, visibilityTimeout *int32, waitTimeSeconds int32) ([]eventbus.ReceivedSQSMessage, error) {
	msgs, err := a.store.ReceiveMessage(queueURL, maxMessages, visibilityTimeout, waitTimeSeconds)
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
type dynamoDBInvokerAdapter struct {
	provider dynamoDBStoreProvider
}

func (a *dynamoDBInvokerAdapter) store(ctx context.Context, region string) (dynamodbstore.DynamoDBStoreInterface, error) {
	return a.provider.GetStoreForRegion(region)
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
	item, err := s.Items().Put(tableName, dynamoKey, dynamoAttrs)
	if err != nil {
		return nil, err
	}
	result := dynamoItemToPlainMap(item)
	if result == nil {
		result = map[string]interface{}{}
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
	return s.Items().Delete(tableName, dynamoKey)
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
	count := 0
	scanErr := s.Items().Scan(tableName, func(item *dynamodbstore.Item) error {
		if count >= limit {
			return errScanLimitReached
		}
		results = append(results, dynamoItemToPlainMap(item))
		count++
		return nil
	})
	if scanErr != nil && scanErr != errScanLimitReached {
		return nil, scanErr
	}
	return results, nil
}

// Query retrieves items from a DynamoDB table by partition key value.
func (a *dynamoDBInvokerAdapter) Query(ctx context.Context, region, tableName, partitionKeyValue string, limit int) ([]map[string]interface{}, error) {
	s, err := a.store(ctx, region)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 1000
	}
	var results []map[string]interface{}
	count := 0
	queryErr := s.Items().ScanByPartitionKey(tableName, partitionKeyValue, func(item *dynamodbstore.Item) error {
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
	_, err = s.Items().Put(tableName, dynamoKey, dynamoAttrs)
	return err
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

type wafInvokerAdapter struct {
	store *wafstore.WebACLAssociationStore
}

func (a *wafInvokerAdapter) AssociateWebACL(webACLArn, resourceArn string) error {
	return a.store.Associate(webACLArn, resourceArn)
}

func (a *wafInvokerAdapter) DisassociateWebACL(webACLArn, resourceArn string) error {
	return a.store.Disassociate(webACLArn, resourceArn)
}

type cloudWatchMetricInvokerAdapter struct {
	storageMgr *storage.RegionStorageManager
	dataPath   string
	stores     sync.Map
}

func (a *cloudWatchMetricInvokerAdapter) getOrCreateStore(region string) (*cwstore.MetricChunkStore, error) {
	if cached, ok := a.stores.Load(region); ok {
		if typed, ok := cached.(*cwstore.MetricChunkStore); ok {
			return typed, nil
		}
	}
	regionStorage, err := a.storageMgr.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store, err := cwstore.NewMetricChunkStoreWithIndex(regionStorage, region, a.dataPath)
	if err != nil {
		return nil, err
	}
	if actual, loaded := a.stores.LoadOrStore(region, store); loaded {
		store.Close()
		typed, ok := actual.(*cwstore.MetricChunkStore)
		if !ok {
			return nil, fmt.Errorf("cloudwatch metric: unexpected store type for region %s", region)
		}
		return typed, nil
	}
	return store, nil
}

func (a *cloudWatchMetricInvokerAdapter) PutMetricData(region, namespace string, metricName string, value float64, timestamp time.Time) error {
	store, err := a.getOrCreateStore(region)
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

type cloudTrailInvokerAdapter struct {
	storageMgr *storage.RegionStorageManager
	accountID  string
	stores     sync.Map
}

func (a *cloudTrailInvokerAdapter) getOrCreateStore(region string) (*cloudtrailstore.CloudTrailStore, error) {
	if cached, ok := a.stores.Load(region); ok {
		if typed, ok := cached.(*cloudtrailstore.CloudTrailStore); ok {
			return typed, nil
		}
	}
	regionStorage, err := a.storageMgr.GetStorage(region)
	if err != nil {
		return nil, err
	}
	ctStore := cloudtrailstore.NewCloudTrailStore(regionStorage, a.accountID, region)
	if actual, loaded := a.stores.LoadOrStore(region, ctStore); loaded {
		typed, ok := actual.(*cloudtrailstore.CloudTrailStore)
		if !ok {
			return nil, fmt.Errorf("cloudtrail: unexpected store type for region %s", region)
		}
		return typed, nil
	}
	return ctStore, nil
}

func (a *cloudTrailInvokerAdapter) LookupEvents(_ context.Context, region, accountID, username string, startTime, endTime time.Time, maxResults int32) ([]eventbus.CloudTrailEventInfo, string, error) {
	ctStore, err := a.getOrCreateStore(region)
	if err != nil {
		return nil, "", err
	}
	query := cloudtrailstore.EventQuery{
		MaxResults: maxResults,
		StartTime:  &startTime,
		EndTime:    &endTime,
		Username:   username,
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

type logsInvokerAdapter struct {
	storageMgr *storage.RegionStorageManager
	accountID  string
	dataPath   string
	stores     sync.Map
}

func (a *logsInvokerAdapter) getOrCreateStore(region string) (*logsstore.Store, error) {
	if cached, ok := a.stores.Load(region); ok {
		if typed, ok := cached.(*logsstore.Store); ok {
			return typed, nil
		}
	}
	regionStorage, err := a.storageMgr.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store, err := logsstore.NewStore(regionStorage, regionStorage.Bucket("logs-"+region), a.accountID, region, a.dataPath)
	if err != nil {
		return nil, err
	}
	if actual, loaded := a.stores.LoadOrStore(region, store); loaded {
		typed, ok := actual.(*logsstore.Store)
		if !ok {
			return nil, fmt.Errorf("logs: unexpected store type for region %s", region)
		}
		return typed, nil
	}
	return store, nil
}

func (a *logsInvokerAdapter) EnsureLogGroup(_ context.Context, region, logGroupName, accountID string) error {
	store, err := a.getOrCreateStore(region)
	if err != nil {
		return err
	}
	_, err = store.GetLogGroup(logGroupName)
	if err != nil {
		lg := logsstore.NewLogGroup(logGroupName, region, accountID)
		return store.CreateLogGroup(lg)
	}
	return nil
}

func (a *logsInvokerAdapter) EnsureLogStream(_ context.Context, region, logGroupName, logStreamName string) error {
	store, err := a.getOrCreateStore(region)
	if err != nil {
		return err
	}
	if _, err := store.GetLogStream(logGroupName, logStreamName); err == nil {
		return nil
	}
	ls := logsstore.NewLogStream(logStreamName, logGroupName)
	return store.CreateLogStream(ls)
}

func (a *logsInvokerAdapter) PutLogEvents(_ context.Context, region, logGroupName, logStreamName string, entries []eventbus.LogsLogEntry) error {
	store, err := a.getOrCreateStore(region)
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

func (a *neptuneGraphInvokerAdapter) ExecuteQueryOnGraph(ctx context.Context, graphID string, query string, language string, parameters map[string]interface{}) (interface{}, error) {
	return a.service.ExecuteQueryOnGraph(ctx, graphID, query, language, parameters)
}
