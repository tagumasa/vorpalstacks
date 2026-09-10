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
	"vorpalstacks/internal/eventbus"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Dispatch limits shared by the HTTP-plane data source dispatchers. Raw
// values appear here and nowhere else.
const (
	httpDispatchTimeout    = 30 * time.Second
	httpMaxResponseBytes   = 6 * 1024 * 1024
	dynamoDefaultPageLimit = 1000
	maxSQLStatements       = 2
)

var sharedHTTPClient = &http.Client{Timeout: httpDispatchTimeout}

// dispatchDataSource routes a resolver payload to the appropriate data source
// based on the data source type (DynamoDB, Lambda, HTTP, EventBridge, Neptune, None, etc.).
func (e *graphQLEngine) dispatchDataSource(
	ctx context.Context,
	reqCtx *request.RequestContext,
	apiId string,
	dataSourceName string,
	payload interface{},
	maxBatchSize int32,
) (interface{}, error) {
	if dataSourceName == "" {
		return nil, nil
	}

	ds, err := e.store.GetDataSource(apiId, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("data source %q not found: %w", dataSourceName, err)
	}

	switch ds.Type {
	case "AWS_LAMBDA":
		return e.dispatchLambda(ctx, reqCtx, ds, payload)
	case "AMAZON_DYNAMODB":
		return e.dispatchDynamoDB(ctx, reqCtx, ds, payload, maxBatchSize)
	case "HTTP":
		return e.dispatchHTTP(ctx, reqCtx, ds, payload)
	case "AMAZON_EVENTBRIDGE":
		return e.dispatchEventBridge(ctx, reqCtx, ds, payload)
	case "AMAZON_NEPTUNE":
		return e.dispatchNeptune(ctx, reqCtx, ds, payload)
	case "NONE":
		return e.dispatchNone()
	case "AMAZON_ELASTICSEARCH", "AMAZON_OPENSEARCH_SERVICE":
		return e.dispatchOpenSearch(ctx, reqCtx, ds, payload)
	case "RELATIONAL_DATABASE":
		return e.dispatchRDS(ctx, reqCtx, ds, payload)
	default:
		return nil, fmt.Errorf("unsupported data source type: %s", ds.Type)
	}
}

// dispatchLambda invokes a Lambda function via the event bus LambdaInvoker.
// The payload is passed as the Lambda event payload.
func (e *graphQLEngine) dispatchLambda(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	if e.bus == nil {
		return nil, fmt.Errorf("event bus not configured for Lambda invocation")
	}

	if ds.LambdaConfig == nil || ds.LambdaConfig.LambdaFunctionArn == "" {
		return nil, fmt.Errorf("lambda data source has no function ARN configured")
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Lambda payload: %w", err)
	}

	// Use the shared ARN utility instead of a local duplicate
	functionName := svcarn.ExtractFunctionNameFromARN(ds.LambdaConfig.LambdaFunctionArn)
	if functionName == "" {
		functionName = ds.LambdaConfig.LambdaFunctionArn
	}

	invoker := e.bus.LambdaInvoker()
	if invoker == nil {
		return nil, fmt.Errorf("Lambda invoker not configured on event bus")
	}

	_, responseBytes, err := invoker.InvokeForGateway(ctx, functionName, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("lambda invocation failed: %w", err)
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
	maxBatchSize int32,
) (interface{}, error) {
	if e.bus == nil {
		return nil, fmt.Errorf("event bus not configured for DynamoDB invocation")
	}

	invoker := e.bus.DynamoDBInvoker()
	if invoker == nil {
		return nil, fmt.Errorf("DynamoDB invoker not configured on event bus")
	}

	region := reqCtx.GetRegion()

	payloadMap, ok := toMap(payload)
	if !ok {
		return nil, fmt.Errorf("DynamoDB payload must be a JSON object")
	}

	operation := request.GetStringParam(payloadMap, "operation")
	tableName := request.GetStringParam(payloadMap, "table")
	if tableName == "" && ds.DynamodbConfig != nil {
		tableName = ds.DynamodbConfig.TableName
	}
	if tableName == "" {
		return nil, fmt.Errorf("DynamoDB table name is required")
	}

	keyMap := getMapFromMap(payloadMap, "key")
	itemMap := getMapFromMap(payloadMap, "item")
	if itemMap == nil {
		itemMap = getMapFromMap(payloadMap, "attributeValues")
	}

	switch operation {
	case "GetItem":
		if keyMap == nil {
			return nil, fmt.Errorf("DynamoDB GetItem requires 'key'")
		}
		return invoker.GetItem(ctx, region, tableName, keyMap)

	case "PutItem":
		if keyMap == nil && itemMap == nil {
			return nil, fmt.Errorf("DynamoDB PutItem requires 'key' or 'item'")
		}
		if keyMap == nil {
			keyMap = make(map[string]interface{})
		}
		if itemMap == nil {
			itemMap = make(map[string]interface{})
		}
		if idVal, ok := keyMap["id"]; ok && idVal == "" {
			newID := uuid.New().String()
			keyMap["id"] = newID
			itemMap["id"] = newID
		}
		return invoker.PutItem(ctx, region, tableName, keyMap, itemMap)

	case "DeleteItem":
		if keyMap == nil {
			return nil, fmt.Errorf("DynamoDB DeleteItem requires 'key'")
		}
		err := invoker.DeleteItem(ctx, region, tableName, keyMap)
		if err != nil {
			return nil, fmt.Errorf("DynamoDB DeleteItem failed: %w", err)
		}
		return map[string]interface{}{}, nil

	case "Scan":
		scanLimit := dynamoDefaultPageLimit
		if maxBatchSize > 0 && int(maxBatchSize) < scanLimit {
			scanLimit = int(maxBatchSize)
		}
		exclusiveStartKey := request.GetStringParam(payloadMap, "nextToken")
		items, nextToken, err := invoker.ScanWithPagination(ctx, region, tableName, scanLimit, exclusiveStartKey)
		if err != nil {
			return nil, fmt.Errorf("DynamoDB Scan failed: %w", err)
		}
		result := map[string]interface{}{"items": items}
		if nextToken != "" {
			result["scannedCount"] = len(items)
			result["nextToken"] = nextToken
		}
		return result, nil

	case "Query":
		if keyMap == nil {
			return nil, fmt.Errorf("DynamoDB Query requires 'key'")
		}
		pkVal := request.GetStringParam(keyMap, "id")
		if pkVal == "" {
			for _, v := range keyMap {
				pkVal = fmt.Sprintf("%v", v)
				break
			}
		}
		queryLimit := dynamoDefaultPageLimit
		if maxBatchSize > 0 && int(maxBatchSize) < queryLimit {
			queryLimit = int(maxBatchSize)
		}
		exclusiveStartKey := request.GetStringParam(payloadMap, "nextToken")
		items, nextToken, err := invoker.QueryWithPagination(ctx, region, tableName, pkVal, queryLimit, exclusiveStartKey)
		if err != nil {
			return nil, fmt.Errorf("DynamoDB Query failed: %w", err)
		}
		result := map[string]interface{}{"items": items}
		if nextToken != "" {
			result["scannedCount"] = len(items)
			result["nextToken"] = nextToken
		}
		return result, nil

	case "UpdateItem":
		if keyMap == nil {
			return nil, fmt.Errorf("DynamoDB UpdateItem requires 'key'")
		}
		err := invoker.UpdateItem(ctx, region, tableName, keyMap, itemMap)
		if err != nil {
			return nil, fmt.Errorf("DynamoDB UpdateItem failed: %w", err)
		}
		return map[string]interface{}{}, nil

	default:
		return nil, fmt.Errorf("unsupported DynamoDB operation: %s", operation)
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

	return postJSON(ctx, endpoint, payload, "HTTP")
}

// postJSON dispatches a JSON POST to an HTTP endpoint and decodes the JSON
// response, falling back to the raw body string when the response is not
// valid JSON. Shared by the HTTP and OpenSearch data source dispatchers.
func postJSON(ctx context.Context, endpoint string, payload interface{}, source string) (interface{}, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s payload: %w", source, err)
	}

	httpCtx, cancel := context.WithTimeout(ctx, httpDispatchTimeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(httpCtx, http.MethodPost, endpoint, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create %s request: %w", source, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := sharedHTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", source, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, httpMaxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s response: %w", source, err)
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
		return nil, fmt.Errorf("event bus not configured")
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
				ID:          uuid.New().String(),
				Timestamp:   time.Now().UTC(),
				Source:      "appsync.amazonaws.com",
				Region:      e.store.GetRegion(),
				AccountID:   e.store.GetAccountID(),
				EventBusARN: eventBusArn,
				Payload:     evMap,
			}
			if err := e.bus.Publish(ctx, publishEvent); err != nil {
				// The first delivery failure is the operation's failure:
				// reporting success for a publish that did not land
				// falsifies the delivery state.
				return nil, fmt.Errorf("EventBridge publish failed: %w", err)
			}
		}
	}

	return map[string]interface{}{"success": true}, nil
}

// dispatchNeptune forwards queries to the NeptuneGraph service via the event bus.
func (e *graphQLEngine) dispatchNeptune(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	if e.bus == nil {
		return nil, fmt.Errorf("event bus not configured for Neptune invocation")
	}

	invoker := e.bus.NeptuneGraphInvoker()
	if invoker == nil {
		return nil, fmt.Errorf("NeptuneGraph invoker not configured on event bus")
	}

	graphID := ""
	if ds.NeptuneConfig != nil {
		graphID = ds.NeptuneConfig.GraphID
	}
	if graphID == "" {
		return nil, fmt.Errorf("neptune data source has no graphId configured")
	}

	payloadMap, ok := toMap(payload)
	if !ok {
		return nil, fmt.Errorf("neptune payload must be a JSON object")
	}

	query := request.GetStringParam(payloadMap, "query")
	if query == "" {
		return nil, fmt.Errorf("neptune query is required")
	}

	language := request.GetStringParam(payloadMap, "language")
	if language == "" {
		language = "CYPHER"
	}

	parameters := getMapFromMap(payloadMap, "parameters")

	return invoker.ExecuteQueryOnGraph(ctx, graphID, query, language, parameters)
}

// dispatchNone returns null for NONE data sources. These are used for
// resolvers that don't interact with any external data source.
func (e *graphQLEngine) dispatchNone() (interface{}, error) {
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

	return postJSON(ctx, endpoint, payload, "OpenSearch")
}

// dispatchRDS forwards SQL queries to the RDS Data API via the EventBus
// RDSDataInvoker. It extracts SQL statements from the VTL or JS resolver
// payload and executes them against the configured RDS instance.
func (e *graphQLEngine) dispatchRDS(
	ctx context.Context,
	reqCtx *request.RequestContext,
	ds *appsyncstore.DataSource,
	payload interface{},
) (interface{}, error) {
	if e.bus == nil {
		return nil, fmt.Errorf("event bus not configured for RDS invocation")
	}

	invoker := e.bus.RDSDataInvoker()
	if invoker == nil {
		return nil, fmt.Errorf("RDSData invoker not configured on event bus")
	}

	if ds.RelationalDatabaseConfig == nil || ds.RelationalDatabaseConfig.RdsHttpEndpointConfig == nil {
		return nil, fmt.Errorf("RDS data source has no endpoint config")
	}

	rdsCfg := ds.RelationalDatabaseConfig.RdsHttpEndpointConfig
	resourceArn := rdsCfg.DbClusterIdentifier
	secretArn := rdsCfg.AwsSecretStoreArn
	database := rdsCfg.DatabaseName
	schema := rdsCfg.Schema

	// Extract SQL from payload (VTL format: {version, statements[]})
	payloadMap, ok := toMap(payload)
	if !ok {
		return nil, fmt.Errorf("RDS payload must be a JSON object")
	}

	statements := extractStatements(payloadMap)
	if len(statements) == 0 {
		return nil, fmt.Errorf("no SQL statements in payload")
	}
	if len(statements) > maxSQLStatements {
		return nil, fmt.Errorf("maximum %d SQL statements allowed per request", maxSQLStatements)
	}

	// The RDS request mapping template admits up to two statements per
	// request; they execute in order — the documented pattern is a mutation
	// followed by a read-back select — and the resolver returns the final
	// statement's result.
	var result interface{}
	for _, stmt := range statements {
		r, execErr := invoker.ExecuteStatement(ctx, resourceArn, secretArn, database, schema, stmt, false, "")
		if execErr != nil {
			return nil, fmt.Errorf("RDS execution failed: %w", execErr)
		}
		result = r
	}

	return result, nil
}

// extractStatements pulls SQL strings from VTL or JS resolver payloads.
func extractStatements(payload map[string]interface{}) []string {
	// VTL format: {"version":"2018-05-29","statements":["SELECT ..."]}
	if stmts, ok := payload["statements"]; ok {
		if stmtList, ok := stmts.([]interface{}); ok {
			var result []string
			for _, s := range stmtList {
				if str, ok := s.(string); ok && str != "" {
					result = append(result, str)
				}
			}
			return result
		}
	}

	// JS resolver payloads carry the SQL string in the "sql" member; the
	// typed check and the generic lookup read the same key, so one
	// extraction covers both.
	if sql := request.GetStringParam(payload, "sql"); sql != "" {
		return []string{sql}
	}

	return nil
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

// --- Generic map access helpers ---

func toMap(v interface{}) (map[string]interface{}, bool) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, true
	}
	return nil, false
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
