package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/eventbus"
)

var dynamodbRegex = regexp.MustCompile(`dynamodb:action/[^/]+/([^/]+)`)

func isDynamoDBURI(uri string) bool {
	return strings.Contains(uri, ":dynamodb:")
}

func (e *AWSExecutor) executeDynamoDB(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.bus == nil || e.bus.DynamoDBInvoker() == nil {
		return nil, &IntegrationError{
			Message:  "DynamoDB invoker not configured",
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	req = applyRequestParameterMapping(req)

	processedBody, pErr := processRequestBody(req)
	if pErr != nil {
		return nil, pErr
	}
	req.Body = processedBody

	action := req.QueryParams["Action"]
	if action == "" {
		action = req.Headers["Action"]
	}
	if action == "" {
		return nil, &IntegrationError{
			Message:  "DynamoDB Action not specified in integration request",
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}

	var tableName string
	if matches := dynamodbRegex.FindStringSubmatch(req.URI); len(matches) >= 2 {
		tableName = matches[1]
	}
	if tableName == "" {
		tableName = req.Headers["TableName"]
	}
	if tableName == "" {
		return nil, &IntegrationError{
			Message:  "Table name not specified in DynamoDB integration URI",
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}

	var result interface{}
	var err error

	switch action {
	case "PutItem":
		var item map[string]interface{}
		if len(req.Body) > 0 {
			if jsonErr := json.Unmarshal(req.Body, &item); jsonErr != nil {
				return nil, &IntegrationError{
					Message:  fmt.Sprintf("DynamoDB PutItem: failed to parse body: %v", jsonErr),
					Type:     "BadRequestException",
					HTTPCode: 400,
				}
			}
		}
		if item == nil {
			item = make(map[string]interface{})
		}
		result, err = e.bus.DynamoDBInvoker().PutItem(ctx, e.region, tableName, nil, item)
	case "GetItem":
		var key map[string]interface{}
		if len(req.Body) > 0 {
			if jsonErr := json.Unmarshal(req.Body, &key); jsonErr != nil {
				return nil, &IntegrationError{
					Message:  fmt.Sprintf("DynamoDB GetItem: failed to parse body: %v", jsonErr),
					Type:     "BadRequestException",
					HTTPCode: 400,
				}
			}
		}
		if key == nil {
			key = make(map[string]interface{})
		}
		result, err = e.bus.DynamoDBInvoker().GetItem(ctx, e.region, tableName, key)
	case "DeleteItem":
		var key map[string]interface{}
		if len(req.Body) > 0 {
			if jsonErr := json.Unmarshal(req.Body, &key); jsonErr != nil {
				return nil, &IntegrationError{
					Message:  fmt.Sprintf("DynamoDB DeleteItem: failed to parse body: %v", jsonErr),
					Type:     "BadRequestException",
					HTTPCode: 400,
				}
			}
		}
		if key == nil {
			key = make(map[string]interface{})
		}
		err = e.bus.DynamoDBInvoker().DeleteItem(ctx, e.region, tableName, key)
		result = map[string]interface{}{}
	case "Scan":
		limit := 100
		if val := req.Headers["Limit"]; val != "" {
			if n, err := fmt.Sscanf(val, "%d", &limit); err != nil || n != 1 {
				limit = 100
			}
		}
		if limit < 1 {
			limit = 100
		}
		if limit > 1000 {
			limit = 1000
		}
		result, err = e.bus.DynamoDBInvoker().Scan(ctx, e.region, tableName, limit)
	case "Query":
		partitionKey := ""
		if len(req.Body) > 0 {
			var body map[string]interface{}
			if jsonErr := json.Unmarshal(req.Body, &body); jsonErr == nil {
				if pk, ok := body["partitionKey"].(string); ok {
					partitionKey = pk
				}
			}
		}
		result, err = e.bus.DynamoDBInvoker().Query(ctx, e.region, tableName, partitionKey, 100)
	case "UpdateItem":
		var key, attrs map[string]interface{}
		if len(req.Body) > 0 {
			var body map[string]interface{}
			if jsonErr := json.Unmarshal(req.Body, &body); jsonErr == nil {
				if k, ok := body["key"].(map[string]interface{}); ok {
					key = k
				}
				if a, ok := body["attributes"].(map[string]interface{}); ok {
					attrs = a
				}
			}
		}
		if key == nil {
			key = make(map[string]interface{})
		}
		if attrs == nil {
			attrs = make(map[string]interface{})
		}
		err = e.bus.DynamoDBInvoker().UpdateItem(ctx, e.region, tableName, key, attrs)
		result = map[string]interface{}{}
	default:
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Unsupported DynamoDB action: %s", action),
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}

	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("DynamoDB operation failed: %v", err),
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	// Apply FilterExpression and ProjectionExpression as post-query filters
	// on Scan/Query result sets. The invoker interface is intentionally
	// minimal, so expressions are evaluated in the integration layer.
	if slice, ok := result.([]map[string]interface{}); ok {
		filterExpr := req.Headers["FilterExpression"]
		projectionExpr := req.Headers["ProjectionExpression"]
		if filterExpr != "" || projectionExpr != "" {
			exprValues := parseDynamoExpressionAttributeValues(req.Headers["ExpressionAttributeValues"])
			result = applyDynamoExpressions(slice, filterExpr, projectionExpr, exprValues)
		}
	}

	responseJSON, _ := json.Marshal(map[string]interface{}{"Response": result})

	statusCode := 200
	headers := map[string]string{"Content-Type": "application/json"}
	body := responseJSON

	// Apply integration response pipeline: selection, templates, parameter
	// mapping, and content handling. For AWS backends, the selection
	// pattern is matched against the HTTP status code.
	if req.IntegrationResponses != nil {
		respConfig := matchIntegrationResponse(req.IntegrationResponses, fmt.Sprintf("%d", statusCode), statusCode)
		if respConfig != nil {
			if respConfig.ResponseTemplates != nil {
				contentType := selectResponseContentType(respConfig.ResponseTemplates, headers, req.Headers)
				if tmpl, ok := respConfig.ResponseTemplates[contentType]; ok && tmpl != "" {
					transformed, tErr := applyResponseTemplate(tmpl, body, req)
					if tErr != nil {
						return nil, &IntegrationError{
							Message:  fmt.Sprintf("Failed to apply response template: %v", tErr),
							Type:     "InternalServerError",
							HTTPCode: 500,
						}
					}
					body = transformed
				}
			}

			if respConfig.StatusCode != "" {
				if parsed, ok := parseStatusCode(respConfig.StatusCode); ok {
					statusCode = parsed
				}
			}

			headers = applyResponseParameterMapping(respConfig.ResponseHeaders, headers, string(body))

			if respConfig.ContentHandling != "" {
				body = applyContentHandlingResponse(body, respConfig.ContentHandling)
			}
		}
	}

	return &IntegrationResponse{
		StatusCode:      statusCode,
		Headers:         headers,
		Body:            body,
		IsBase64Encoded: false,
	}, nil
}

// parseDynamoExpressionAttributeValues parses the JSON-encoded
// ExpressionAttributeValues map that API Gateway passes as an integration
// header. The values arrive in DynamoDB typed format (e.g. {":v": {"S": "x"}}),
// so each value is unwrapped to its native Go type to match the flat item
// format returned by the DynamoDB invoker. Returns nil if the string is
// empty or unparseable.
func parseDynamoExpressionAttributeValues(raw string) map[string]interface{} {
	if raw == "" {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil
	}
	for k, v := range m {
		m[k] = unwrapDynamoDBValue(v)
	}
	return m
}

// unwrapDynamoDBValue extracts the inner value from a DynamoDB typed
// attribute value (e.g. {"S": "hello"} → "hello", {"N": "42"} → 42.0).
// Map (M) and List (L) types are unwrapped recursively so that deeply
// nested typed values are fully flattened to match the item format
// returned by the DynamoDB invoker.
func unwrapDynamoDBValue(v interface{}) interface{} {
	m, ok := v.(map[string]interface{})
	if !ok || len(m) != 1 {
		return v
	}
	for typeKey, inner := range m {
		switch typeKey {
		case "N":
			if s, ok := inner.(string); ok {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					return f
				}
			}
			return inner
		case "S", "B", "BOOL", "NULL":
			return inner
		case "SS", "NS", "BS":
			return inner
		case "L":
			if items, ok := inner.([]interface{}); ok {
				result := make([]interface{}, len(items))
				for i, item := range items {
					result[i] = unwrapDynamoDBValue(item)
				}
				return result
			}
			return inner
		case "M":
			if entries, ok := inner.(map[string]interface{}); ok {
				result := make(map[string]interface{}, len(entries))
				for k, val := range entries {
					result[k] = unwrapDynamoDBValue(val)
				}
				return result
			}
			return inner
		default:
			return inner
		}
	}
	return v
}

// applyDynamoExpressions applies FilterExpression (item filtering) and
// ProjectionExpression (attribute projection) to the result set returned by
// the DynamoDB invoker.
func applyDynamoExpressions(items []map[string]interface{}, filterExpr, projectionExpr string, exprValues map[string]interface{}) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, item := range items {
		if filterExpr == "" || evaluateFilterExpression(item, filterExpr, exprValues) {
			filtered = append(filtered, item)
		}
	}
	if projectionExpr != "" {
		attrs := parseProjectionExpression(projectionExpr)
		for i, item := range filtered {
			projected := make(map[string]interface{}, len(attrs))
			for _, attr := range attrs {
				if v, ok := item[attr]; ok {
					projected[attr] = v
				}
			}
			filtered[i] = projected
		}
	}
	return filtered
}

// parseProjectionExpression splits a comma-separated ProjectionExpression
// into individual attribute names, trimming whitespace.
func parseProjectionExpression(expr string) []string {
	parts := strings.Split(expr, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// evaluateFilterExpression evaluates a simplified DynamoDB FilterExpression
// against a single item. Supports comparison operators (=, <>, <, >, <=, >=),
// attribute_exists / attribute_not_exists, BETWEEN ... AND ..., and IN (...).
// Value placeholders (:name) are resolved from exprValues.
func evaluateFilterExpression(item map[string]interface{}, expr string, exprValues map[string]interface{}) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return true
	}

	if strings.HasPrefix(expr, "attribute_exists(") {
		path := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(expr, "attribute_exists("), ")"))
		_, ok := item[path]
		return ok
	}
	if strings.HasPrefix(expr, "attribute_not_exists(") {
		path := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(expr, "attribute_not_exists("), ")"))
		_, ok := item[path]
		return !ok
	}

	upper := strings.ToUpper(expr)
	if strings.Contains(upper, " BETWEEN ") {
		return evaluateBetween(item, expr, exprValues)
	}
	if strings.Contains(upper, " IN (") {
		return evaluateIn(item, expr, exprValues)
	}

	return evaluateComparison(item, expr, exprValues)
}

// evaluateBetween handles "attr BETWEEN :low AND :high" expressions.
func evaluateBetween(item map[string]interface{}, expr string, exprValues map[string]interface{}) bool {
	upper := strings.ToUpper(expr)
	idx := strings.Index(upper, " BETWEEN ")
	if idx < 0 {
		return false
	}
	lhs := strings.TrimSpace(expr[:idx])
	rest := strings.TrimSpace(expr[idx+len(" BETWEEN "):])
	andIdx := strings.Index(strings.ToUpper(rest), " AND ")
	if andIdx < 0 {
		return false
	}
	lowStr := strings.TrimSpace(rest[:andIdx])
	highStr := strings.TrimSpace(rest[andIdx+len(" AND "):])

	itemVal, exists := item[lhs]
	if !exists {
		return false
	}

	lowVal := resolveOperand(lowStr, exprValues)
	highVal := resolveOperand(highStr, exprValues)

	return compareValues(itemVal, ">=", lowVal) && compareValues(itemVal, "<=", highVal)
}

// evaluateIn handles "attr IN (:v1, :v2, ...)" expressions.
func evaluateIn(item map[string]interface{}, expr string, exprValues map[string]interface{}) bool {
	upper := strings.ToUpper(expr)
	idx := strings.Index(upper, " IN (")
	if idx < 0 {
		return false
	}
	lhs := strings.TrimSpace(expr[:idx])
	listPart := strings.TrimSuffix(strings.TrimSpace(expr[idx+len(" IN ("):]), ")")

	itemVal, exists := item[lhs]
	if !exists {
		return false
	}

	for _, valStr := range strings.Split(listPart, ",") {
		cmpVal := resolveOperand(strings.TrimSpace(valStr), exprValues)
		if compareValues(itemVal, "=", cmpVal) {
			return true
		}
	}
	return false
}

// resolveOperand resolves a value placeholder (prefixed with ":") from
// exprValues, or returns the literal value as-is.
func resolveOperand(operand string, exprValues map[string]interface{}) interface{} {
	if strings.HasPrefix(operand, ":") {
		return exprValues[operand]
	}
	return operand
}

// evaluateComparison handles a single binary comparison: attr op operand.
func evaluateComparison(item map[string]interface{}, expr string, exprValues map[string]interface{}) bool {
	tokens := tokeniseComparison(expr)
	if len(tokens) != 3 {
		return true
	}
	attr, op, operand := tokens[0], tokens[1], tokens[2]

	itemVal, exists := item[attr]
	if !exists {
		return false
	}

	return compareValues(itemVal, op, resolveOperand(operand, exprValues))
}

// tokeniseComparison splits a comparison expression into [attr, op, operand].
func tokeniseComparison(expr string) []string {
	for _, op := range []string{"<=", ">=", "<>", "=", "<", ">"} {
		if idx := strings.Index(expr, op); idx >= 0 {
			return []string{
				strings.TrimSpace(expr[:idx]),
				op,
				strings.TrimSpace(expr[idx+len(op):]),
			}
		}
	}
	return nil
}

// compareValues compares two values using the given operator.
func compareValues(left interface{}, op string, right interface{}) bool {
	lf, lok := toFloat64(left)
	rf, rok := toFloat64(right)
	if lok && rok {
		switch op {
		case "=":
			return lf == rf
		case "<>":
			return lf != rf
		case "<":
			return lf < rf
		case ">":
			return lf > rf
		case "<=":
			return lf <= rf
		case ">=":
			return lf >= rf
		}
	}

	ls := fmt.Sprintf("%v", left)
	rs := fmt.Sprintf("%v", right)
	switch op {
	case "=":
		return ls == rs
	case "<>":
		return ls != rs
	case "<":
		return ls < rs
	case ">":
		return ls > rs
	case "<=":
		return ls <= rs
	case ">=":
		return ls >= rs
	}
	return false
}

// toFloat64 attempts to convert a value to float64 for numeric comparison.
func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}

var kinesisRegex = regexp.MustCompile(`kinesis:action/[^/]+/([^/]+)`)

func isKinesisURI(uri string) bool {
	return strings.Contains(uri, ":kinesis:")
}

func (e *AWSExecutor) executeKinesis(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.bus == nil || e.bus.KinesisInvoker() == nil {
		return nil, &IntegrationError{
			Message:  "Kinesis invoker not configured",
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	req = applyRequestParameterMapping(req)

	action := req.QueryParams["Action"]
	if action == "" {
		action = req.Headers["Action"]
	}
	if action == "" {
		action = "PutRecord"
	}

	var streamName string
	if matches := kinesisRegex.FindStringSubmatch(req.URI); len(matches) >= 2 {
		streamName = matches[1]
	}
	if streamName == "" {
		streamName = req.Headers["StreamName"]
	}
	if streamName == "" {
		return nil, &IntegrationError{
			Message:  "Stream name not specified in Kinesis integration URI",
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}

	switch action {
	case "PutRecord":
		partitionKey := req.Headers["PartitionKey"]
		if partitionKey == "" {
			partitionKey = "default"
		}

		data := req.Body
		if len(data) == 0 {
			data = []byte(req.Headers["Data"])
		}

		sequenceNumber, err := e.bus.KinesisInvoker().PutRecord(ctx, streamName, partitionKey, data)
		if err != nil {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Kinesis PutRecord failed: %v", err),
				Type:     "InternalServerError",
				HTTPCode: 500,
			}
		}

		responseJSON, _ := json.Marshal(map[string]interface{}{
			"PutRecordResponse": map[string]interface{}{
				"PutRecordResult": map[string]string{
					"SequenceNumber": sequenceNumber,
					"ShardId":        "shardId-000000000000",
				},
				"ResponseMetadata": map[string]string{
					"RequestId": fmt.Sprintf("%x", time.Now().UnixNano()),
				},
			},
		})
		return &IntegrationResponse{
			StatusCode:      200,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            responseJSON,
			IsBase64Encoded: false,
		}, nil

	case "ListShards":
		shards, err := e.bus.KinesisInvoker().ListShards(ctx, streamName)
		if err != nil {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Kinesis ListShards failed: %v", err),
				Type:     "InternalServerError",
				HTTPCode: 500,
			}
		}

		type shardInfo struct {
			ShardId string `json:"ShardId"`
		}
		shardItems := make([]shardInfo, 0, len(shards))
		for _, s := range shards {
			shardItems = append(shardItems, shardInfo{ShardId: s.ShardID})
		}

		responseJSON, _ := json.Marshal(map[string]interface{}{"Shards": shardItems})
		return &IntegrationResponse{
			StatusCode:      200,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            responseJSON,
			IsBase64Encoded: false,
		}, nil

	default:
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Unsupported Kinesis action: %s", action),
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}
}

// buildExecutionArn constructs a Step Functions execution ARN from a state
// machine ARN by replacing ":stateMachine:" with ":execution:" and appending
// the execution name. Returns empty string if the input ARN is malformed.
func buildExecutionArn(stateMachineArn, executionName string) string {
	idx := strings.Index(stateMachineArn, ":stateMachine:")
	if idx < 0 {
		return ""
	}
	prefix := stateMachineArn[:idx]
	smName := stateMachineArn[idx+len(":stateMachine:"):]
	if smName == "" {
		return ""
	}
	return fmt.Sprintf("%s:execution:%s:%s", prefix, smName, executionName)
}

var sfnRegex = regexp.MustCompile(`states:action/[^/]+/([^/]+)`)

func isSFNURI(uri string) bool {
	return strings.Contains(uri, ":states:")
}

func (e *AWSExecutor) executeStepFunctions(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.bus == nil {
		return nil, &IntegrationError{
			Message:  "Event bus not configured",
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	req = applyRequestParameterMapping(req)

	action := req.QueryParams["Action"]
	if action == "" {
		action = req.Headers["Action"]
	}
	if action == "" {
		action = "StartExecution"
	}

	switch action {
	case "StartExecution":
		var stateMachineArn string
		if matches := sfnRegex.FindStringSubmatch(req.URI); len(matches) >= 2 {
			stateMachineArn = matches[1]
		}
		if stateMachineArn == "" {
			stateMachineArn = req.Headers["StateMachineArn"]
		}
		if stateMachineArn == "" {
			return nil, &IntegrationError{
				Message:  "State machine ARN not specified in Step Functions integration URI",
				Type:     "BadRequestException",
				HTTPCode: 400,
			}
		}

		input := string(req.Body)
		if input == "" {
			input = "{}"
		}

		executionName := req.Headers["ExecutionName"]
		if executionName == "" {
			executionName = fmt.Sprintf("apigw-%x", time.Now().UnixNano())
		}

		// Build execution ARN from state machine ARN:
		// arn:aws:states:{region}:{account}:stateMachine:{name}
		// → arn:aws:states:{region}:{account}:execution:{name}:{executionName}
		executionArn := buildExecutionArn(stateMachineArn, executionName)
		if executionArn == "" {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Failed to construct execution ARN from state machine ARN: %s", stateMachineArn),
				Type:     "BadRequestException",
				HTTPCode: 400,
			}
		}

		evt := &eventbus.StepFunctionsStartExecutionEvent{
			StateMachineArn: stateMachineArn,
			Input:           input,
		}
		evt.ExecutionArn = executionArn
		evt.Region = e.region
		evt.AccountID = e.accountID

		if err := e.bus.Publish(ctx, evt); err != nil {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Step Functions execution failed: %v", err),
				Type:     "InternalServerError",
				HTTPCode: 500,
			}
		}

		responseJSON, _ := json.Marshal(map[string]interface{}{
			"StartExecutionResponse": map[string]interface{}{
				"StartExecutionResult": map[string]string{
					"ExecutionArn": executionArn,
				},
				"ResponseMetadata": map[string]string{
					"RequestId": fmt.Sprintf("%x", time.Now().UnixNano()),
				},
			},
		})
		return &IntegrationResponse{
			StatusCode:      200,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            responseJSON,
			IsBase64Encoded: false,
		}, nil

	default:
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Unsupported Step Functions action: %s", action),
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}
}
