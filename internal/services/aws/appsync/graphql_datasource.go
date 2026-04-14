package appsync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	dynamodbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// dispatchDataSource routes a resolver payload to the appropriate data source
// based on the data source type (DynamoDB, Lambda, HTTP, EventBridge, Neptune, None, etc.).
func (e *graphQLEngine) dispatchDataSource(
	ctx context.Context,
	reqCtx *request.RequestContext,
	apiId string,
	dataSourceName string,
	payload interface{},
) (interface{}, error) {
	if dataSourceName == "" {
		return nil, nil
	}

	ds, err := e.store.GetDataSource(apiId, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Data source %q not found: %w", dataSourceName, err)
	}

	switch ds.Type {
	case "AWS_LAMBDA":
		return e.dispatchLambda(ctx, reqCtx, ds, payload)
	case "AMAZON_DYNAMODB":
		return e.dispatchDynamoDB(ctx, reqCtx, ds, payload)
	case "HTTP":
		return e.dispatchHTTP(ctx, reqCtx, ds, payload)
	case "AMAZON_EVENTBRIDGE":
		return e.dispatchEventBridge(ctx, reqCtx, ds, payload)
	case "AMAZON_NEPTUNE":
		return e.dispatchNeptune(ctx, reqCtx, ds, payload)
	case "NONE":
		return e.dispatchNone(payload)
	case "AMAZON_ELASTICSEARCH", "AMAZON_OPENSEARCH_SERVICE":
		return e.dispatchOpenSearch(ctx, reqCtx, ds, payload)
	case "RELATIONAL_DATABASE":
		return e.dispatchRDS(ctx, reqCtx, ds, payload)
	default:
		return nil, fmt.Errorf("Unsupported data source type: %s", ds.Type)
	}
}

// dispatchLambda invokes a Lambda function via the injected LambdaInvoker.
// The payload is passed as the Lambda event payload.
func (e *graphQLEngine) dispatchLambda(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	if e.lambdaInvoker == nil {
		return nil, fmt.Errorf("Lambda invoker not configured")
	}

	if ds.LambdaConfig == nil || ds.LambdaConfig.LambdaFunctionArn == "" {
		return nil, fmt.Errorf("Lambda data source has no function ARN configured")
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal Lambda payload: %w", err)
	}

	// Use the shared ARN utility instead of a local duplicate
	functionName := svcarn.ExtractFunctionNameFromARN(ds.LambdaConfig.LambdaFunctionArn)
	if functionName == "" {
		functionName = ds.LambdaConfig.LambdaFunctionArn
	}

	_, responseBytes, err := e.lambdaInvoker.InvokeForGateway(ctx, functionName, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("Lambda invocation failed: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		return string(responseBytes), nil
	}
	return result, nil
}

// dispatchDynamoDB handles DynamoDB data source operations.
// Supports GetItem, PutItem, DeleteItem, UpdateItem, and Scan operations
// identified by the "operation" key in the payload.
func (e *graphQLEngine) dispatchDynamoDB(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	basicStorage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, fmt.Errorf("DynamoDB storage not available: %w", err)
	}
	txnStorage, ok := basicStorage.(storage.TransactionalStorageWith2PC)
	if !ok {
		return nil, fmt.Errorf("DynamoDB storage does not support transactions")
	}
	dynamoDBStore := dynamodbstore.NewDynamoDBStore(txnStorage, reqCtx.GetAccountID(), reqCtx.GetRegion())

	payloadMap, ok := toMap(payload)
	if !ok {
		return nil, fmt.Errorf("DynamoDB payload must be a JSON object")
	}

	operation := getStringFromMap(payloadMap, "operation")
	tableName := getStringFromMap(payloadMap, "table")
	if tableName == "" && ds.DynamodbConfig != nil {
		tableName = ds.DynamodbConfig.TableName
	}
	if tableName == "" {
		return nil, fmt.Errorf("DynamoDB table name is required")
	}

	keyMap := getMapFromMap(payloadMap, "key")
	itemMap := getMapFromMap(payloadMap, "item")

	switch operation {
	case "GetItem":
		if keyMap == nil {
			return nil, fmt.Errorf("DynamoDB GetItem requires 'key'")
		}
		dynamoKey := convertToDynamoKey(keyMap)
		item, err := dynamoDBStore.Items().Get(tableName, dynamoKey)
		if err != nil {
			return map[string]interface{}{}, nil
		}
		return dynamoItemToMap(item), nil

	case "PutItem":
		if keyMap == nil && itemMap == nil {
			return nil, fmt.Errorf("DynamoDB PutItem requires 'key' or 'item'")
		}
		dynamoKey := convertToDynamoKey(keyMap)
		attrs := convertToDynamoAttrs(itemMap)
		_, err := dynamoDBStore.Items().Put(tableName, dynamoKey, attrs)
		if err != nil {
			return nil, fmt.Errorf("DynamoDB PutItem failed: %w", err)
		}
		return map[string]interface{}{}, nil

	case "DeleteItem":
		if keyMap == nil {
			return nil, fmt.Errorf("DynamoDB DeleteItem requires 'key'")
		}
		dynamoKey := convertToDynamoKey(keyMap)
		err := dynamoDBStore.Items().Delete(tableName, dynamoKey)
		if err != nil {
			return nil, fmt.Errorf("DynamoDB DeleteItem failed: %w", err)
		}
		return map[string]interface{}{}, nil

	case "Scan":
		var items []interface{}
		err := dynamoDBStore.Items().Scan(tableName, func(item *dynamodbstore.Item) error {
			items = append(items, dynamoItemToMap(item))
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("DynamoDB Scan failed: %w", err)
		}
		return map[string]interface{}{"items": items}, nil

	case "Query":
		if keyMap == nil {
			return nil, fmt.Errorf("DynamoDB Query requires 'key'")
		}
		pkVal := getStringFromMap(keyMap, "id")
		if pkVal == "" {
			for _, v := range keyMap {
				pkVal = fmt.Sprintf("%v", v)
				break
			}
		}
		var items []interface{}
		err := dynamoDBStore.Items().ScanByPartitionKey(tableName, pkVal, func(item *dynamodbstore.Item) error {
			items = append(items, dynamoItemToMap(item))
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("DynamoDB Query failed: %w", err)
		}
		return map[string]interface{}{"items": items}, nil

	case "UpdateItem":
		if keyMap == nil {
			return nil, fmt.Errorf("DynamoDB UpdateItem requires 'key'")
		}
		dynamoKey := convertToDynamoKey(keyMap)
		attrs := convertToDynamoAttrs(itemMap)
		_, err := dynamoDBStore.Items().Put(tableName, dynamoKey, attrs)
		if err != nil {
			return nil, fmt.Errorf("DynamoDB UpdateItem failed: %w", err)
		}
		return map[string]interface{}{}, nil

	default:
		return nil, fmt.Errorf("Unsupported DynamoDB operation: %s", operation)
	}
}

// dispatchHTTP proxies requests to an HTTP data source endpoint.
func (e *graphQLEngine) dispatchHTTP(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	if ds.HttpConfig == nil || ds.HttpConfig.Endpoint == "" {
		return nil, fmt.Errorf("HTTP data source has no endpoint configured")
	}

	endpoint := ds.HttpConfig.Endpoint
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal HTTP payload: %w", err)
	}

	httpCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(httpCtx, http.MethodPost, endpoint, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("Failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 6*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("Failed to read HTTP response: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return string(body), nil
	}
	return result, nil
}

// dispatchEventBridge publishes an event to EventBridge via the event bus.
// AppSync EventBridge data sources send events through the EventBus pipeline
// rather than directly to the EventBridge store, matching the pattern used
// by other services (EventBridge itself, Scheduler, SNS, etc.).
func (e *graphQLEngine) dispatchEventBridge(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	payloadMap, ok := toMap(payload)
	if !ok {
		return nil, fmt.Errorf("EventBridge payload must be a JSON object")
	}

	if e.bus == nil {
		return nil, fmt.Errorf("Event bus not configured")
	}

	eventBusArn := ""
	if ds.EventBridgeConfig != nil {
		eventBusArn = ds.EventBridgeConfig.EventBusArn
	}

	events := getArrayFromMap(payloadMap, "events")
	if len(events) == 0 {
		events = []interface{}{payloadMap}
	}

	for _, ev := range events {
		if evMap, ok := ev.(map[string]interface{}); ok {
			publishEvent := &appsyncPublishEvent{
				ID:          generateUUID(),
				Timestamp:   time.Now().UTC(),
				Source:      "aws.appsync",
				Region:      e.store.GetRegion(),
				AccountID:   e.store.GetAccountID(),
				EventBusARN: eventBusArn,
				Payload:     evMap,
			}
			if err := e.bus.Publish(ctx, publishEvent); err != nil {
				logs.Warn("Failed to publish AppSync event to bus",
					logs.String("eventBus", eventBusArn),
					logs.Err(err))
			}
		}
	}

	return map[string]interface{}{"success": true}, nil
}

// dispatchNeptune forwards queries to the Neptune data service via the graph DB manager.
func (e *graphQLEngine) dispatchNeptune(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	return nil, fmt.Errorf("Neptune data source not yet implemented")
}

// dispatchNone returns null for NONE data sources. These are used for
// resolvers that don't interact with any external data source.
func (e *graphQLEngine) dispatchNone(payload interface{}) (interface{}, error) {
	return nil, nil
}

// dispatchOpenSearch forwards queries to an OpenSearch/Elasticsearch cluster.
func (e *graphQLEngine) dispatchOpenSearch(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	endpoint := ""
	if ds.OpenSearchServiceConfig != nil {
		endpoint = ds.OpenSearchServiceConfig.Endpoint
	} else if ds.ElasticsearchConfig != nil {
		endpoint = ds.ElasticsearchConfig.Endpoint
	}
	if endpoint == "" {
		return nil, fmt.Errorf("OpenSearch data source has no endpoint configured")
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal OpenSearch payload: %w", err)
	}

	httpCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(httpCtx, http.MethodPost, endpoint, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("Failed to create OpenSearch request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("OpenSearch request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 6*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("Failed to read OpenSearch response: %w", err)
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return string(body), nil
	}
	return result, nil
}

// dispatchRDS forwards queries to an RDS Data API endpoint.
func (e *graphQLEngine) dispatchRDS(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	return nil, fmt.Errorf("RDS data source not yet implemented")
}

// --- Event types ---

// appsyncPublishEvent is published when an AppSync resolver dispatches
// to an EventBridge data source.
type appsyncPublishEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Source      string                 `json:"source"`
	Region      string                 `json:"region"`
	AccountID   string                 `json:"account_id"`
	EventBusARN string                 `json:"event_bus_arn"`
	Payload     map[string]interface{} `json:"payload"`
}

// EventType returns the event type identifier for AppSync EventBridge publish events.
func (e *appsyncPublishEvent) EventType() string { return "appsync:eventbridge-publish" }

// EventID returns the unique identifier for this event.
func (e *appsyncPublishEvent) EventID() string { return e.ID }

// EventTimestamp returns the time at which this event was created.
func (e *appsyncPublishEvent) EventTimestamp() time.Time { return e.Timestamp }

// EventSource returns the source of this event (e.g. "aws.appsync").
func (e *appsyncPublishEvent) EventSource() string { return e.Source }

// EventRegion returns the AWS region associated with this event.
func (e *appsyncPublishEvent) EventRegion() string { return e.Region }

// EventAccountID returns the AWS account ID associated with this event.
func (e *appsyncPublishEvent) EventAccountID() string { return e.AccountID }

// EventDepth returns the propagation depth of this event within the bus pipeline.
func (e *appsyncPublishEvent) EventDepth() int { return 0 }

// SetEventDepth sets the propagation depth of this event within the bus pipeline.
func (e *appsyncPublishEvent) SetEventDepth(d int) {}

// EventCaller returns the caller context that produced this event.
func (e *appsyncPublishEvent) EventCaller() eventbus.CallerContext { return eventbus.CallerContext{} }

// --- DynamoDB conversion helpers ---

// convertToDynamoKey converts a map[string]interface{} to a DynamoDB key
// format using the store's AttributeValue type.
func convertToDynamoKey(keyMap map[string]interface{}) map[string]*dynamodbstore.AttributeValue {
	if keyMap == nil {
		return nil
	}
	result := make(map[string]*dynamodbstore.AttributeValue)
	for k, v := range keyMap {
		result[k] = interfaceToAttributeValue(v)
	}
	return result
}

// convertToDynamoAttrs converts a flat map to DynamoDB attribute format
// using the store's AttributeValue type.
func convertToDynamoAttrs(itemMap map[string]interface{}) map[string]*dynamodbstore.AttributeValue {
	if itemMap == nil {
		return nil
	}
	result := make(map[string]*dynamodbstore.AttributeValue)
	for k, v := range itemMap {
		result[k] = interfaceToAttributeValue(v)
	}
	return result
}

// interfaceToAttributeValue converts a Go interface{} value to a DynamoDB AttributeValue.
func interfaceToAttributeValue(v interface{}) *dynamodbstore.AttributeValue {
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
			m[mk] = interfaceToAttributeValue(mv)
		}
		return &dynamodbstore.AttributeValue{M: m}
	case []interface{}:
		l := make([]*dynamodbstore.AttributeValue, len(val))
		for i, item := range val {
			l[i] = interfaceToAttributeValue(item)
		}
		return &dynamodbstore.AttributeValue{L: l}
	case *dynamodbstore.AttributeValue:
		return val
	default:
		s := fmt.Sprintf("%v", val)
		return &dynamodbstore.AttributeValue{S: &s}
	}
}

// dynamoItemToMap converts a DynamoDB Item to a plain map[string]interface{}
// suitable for JSON serialisation in the GraphQL response.
func dynamoItemToMap(item *dynamodbstore.Item) map[string]interface{} {
	if item == nil {
		return nil
	}
	result := make(map[string]interface{})
	if item.Key != nil {
		for k, v := range item.Key {
			result[k] = attributeValueToInterface(v)
		}
	}
	if item.Attributes != nil {
		for k, v := range item.Attributes {
			result[k] = attributeValueToInterface(v)
		}
	}
	return result
}

// attributeValueToInterface converts a DynamoDB AttributeValue back to a
// plain Go interface{} for JSON serialisation.
func attributeValueToInterface(av *dynamodbstore.AttributeValue) interface{} {
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
		return *av.N
	}
	if av.BOOL != nil {
		return *av.BOOL
	}
	if av.M != nil {
		m := make(map[string]interface{})
		for k, v := range av.M {
			m[k] = attributeValueToInterface(v)
		}
		return m
	}
	if av.L != nil {
		l := make([]interface{}, len(av.L))
		for i, v := range av.L {
			l[i] = attributeValueToInterface(v)
		}
		return l
	}
	if av.SS != nil {
		return av.SS
	}
	if av.NS != nil {
		return av.NS
	}
	if av.B != nil {
		return string(av.B)
	}
	return nil
}

// generateUUID generates a new UUID string.
func generateUUID() string {
	return uuid.New().String()
}

// --- Generic map access helpers ---

func toMap(v interface{}) (map[string]interface{}, bool) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, true
	}
	return nil, false
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getMapFromMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key]; ok {
		if sub, ok := v.(map[string]interface{}); ok {
			return sub
		}
	}
	return nil
}

func getArrayFromMap(m map[string]interface{}, key string) []interface{} {
	if v, ok := m[key]; ok {
		if arr, ok := v.([]interface{}); ok {
			return arr
		}
	}
	return nil
}
