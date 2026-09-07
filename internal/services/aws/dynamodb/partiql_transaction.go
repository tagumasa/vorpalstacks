package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ExecuteTransaction executes a TransactWriteItems operation for PartiQL statements.
func (s *DynamoDBService) ExecuteTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statements, ok := req.Parameters["TransactStatements"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	return s.executeTransactionCore(ctx, reqCtx, executeTransactionInput{
		TransactStatements: statements,
		Parameters:         req.Parameters,
	})
}

// BatchExecuteStatement executes multiple PartiQL statements in a batch.
func (s *DynamoDBService) BatchExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statements, ok := req.Parameters["Statements"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	if !validatePartiQLBatchCount(len(statements)) {
		return nil, ErrInvalidParameter
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	var responses []map[string]interface{}
	var consumedCapacities []map[string]interface{}

	for _, stmt := range statements {
		stmtMap, ok := stmt.(map[string]interface{})
		if !ok {
			responses = append(responses, map[string]interface{}{
				"Error": map[string]interface{}{
					"Code":    "ValidationError",
					"Message": "Invalid statement format",
				},
			})
			continue
		}

		statement, _ := stmtMap["Statement"].(string)
		if statement == "" {
			responses = append(responses, map[string]interface{}{
				"Error": map[string]interface{}{
					"Code":    "ValidationError",
					"Message": "Statement is required",
				},
			})
			continue
		}

		tableName := extractTableNameFromStatement(statement)
		consistentRead := request.GetBoolParam(stmtMap, "ConsistentRead")

		params := parsePartiQLParams(stmtMap)
		result, err := s.ExecuteStatement(ctx, reqCtx, &request.ParsedRequest{
			Parameters: map[string]interface{}{
				"Statement":      statement,
				"Parameters":     params.Parameters,
				"ConsistentRead": consistentRead,
			},
		})

		response := map[string]interface{}{
			"TableName": tableName,
		}

		if err != nil {
			errorCode, errorMessage := mapPartiQLError(err)
			response["Error"] = map[string]interface{}{
				"Code":    errorCode,
				"Message": errorMessage,
			}
			if apiErr, ok := err.(*APIError); ok && apiErr.Code == "com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException" {
				if errMap, ok := response["Error"].(map[string]interface{}); ok {
					errMap["Item"] = nil
				}
			}
		} else {
			var item interface{} = nil
			if resultMap, ok := result.(map[string]interface{}); ok {
				if items, ok := resultMap["Items"].([]map[string]interface{}); ok && len(items) > 0 {
					item = items[0]
				}
				if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
					if cc, ok := resultMap["ConsumedCapacity"].(map[string]interface{}); ok {
						consumedCapacities = append(consumedCapacities, cc)
					}
				}
			}
			response["Item"] = item
		}

		responses = append(responses, response)
	}

	resp := map[string]interface{}{
		"Responses": responses,
	}

	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp, nil
}

func mapPartiQLError(err error) (string, string) {
	if apiErr, ok := err.(*APIError); ok {
		switch apiErr.Code {
		case "com.amazonaws.dynamodb.v20120810#ConditionalCheckFailedException":
			return "ConditionalCheckFailed", "The conditional request failed"
		case "com.amazonaws.dynamodb.v20120810#ResourceNotFoundException":
			return "ResourceNotFound", "Requested resource not found"
		case "com.amazonaws.dynamodb.v20120810#ResourceInUseException":
			return "ResourceInUse", apiErr.Message
		case "com.amazon.coral.validate#ValidationException":
			return "ValidationError", apiErr.Message
		case "com.amazonaws.dynamodb.v20120810#TransactionConflictException":
			return "TransactionConflict", "Transaction is ongoing for the item"
		case "com.amazonaws.dynamodb.v20120810#DuplicateItemException":
			return "DuplicateItem", "Item already exists"
		default:
			return "InternalServerError", apiErr.Message
		}
	}
	return "InternalServerError", err.Error()
}
