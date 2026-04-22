package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
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

	action := req.QueryParams["Action"]
	if action == "" {
		action = req.Headers["Action"]
	}
	if action == "" && len(req.Body) > 0 {
		action = "PutItem"
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
			_ = json.Unmarshal(req.Body, &item)
		}
		if item == nil {
			item = make(map[string]interface{})
		}
		result, err = e.bus.DynamoDBInvoker().PutItem(ctx, e.region, tableName, nil, item)
	case "GetItem":
		var key map[string]interface{}
		if len(req.Body) > 0 {
			_ = json.Unmarshal(req.Body, &key)
		}
		if key == nil {
			key = make(map[string]interface{})
		}
		result, err = e.bus.DynamoDBInvoker().GetItem(ctx, e.region, tableName, key)
	case "DeleteItem":
		var key map[string]interface{}
		if len(req.Body) > 0 {
			_ = json.Unmarshal(req.Body, &key)
		}
		if key == nil {
			key = make(map[string]interface{})
		}
		err = e.bus.DynamoDBInvoker().DeleteItem(ctx, e.region, tableName, key)
		result = map[string]interface{}{}
	case "Scan":
		limit := 100
		if val := req.Headers["Limit"]; val != "" {
			fmt.Sscanf(val, "%d", &limit)
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

	responseJSON, _ := json.Marshal(map[string]interface{}{"Response": result})
	return &IntegrationResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            responseJSON,
		IsBase64Encoded: false,
	}, nil
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

		smName := ""
		if idx := strings.Index(stateMachineArn, ":stateMachine:"); idx >= 0 {
			smName = stateMachineArn[idx+len(":stateMachine:"):]
		}
		executionArn := fmt.Sprintf("%s:execution:%s", strings.TrimSuffix(stateMachineArn, ":stateMachine:"+smName), executionName)

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
