package dynamodb

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

type statementType int

const (
	statementTypeUnknown statementType = iota
	statementTypeRead
	statementTypeWrite
)

// ExecuteTransaction executes a TransactWriteItems operation for PartiQL statements.
func (s *DynamoDBService) ExecuteTransaction(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statements, ok := req.Parameters["TransactStatements"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	if len(statements) == 0 {
		return nil, ErrInvalidParameter
	}

	if len(statements) > 100 {
		return nil, ErrInvalidParameter
	}

	parsedStatements := make([]struct {
		statement string
		params    *partiQLParams
		stmtType  statementType
	}, 0, len(statements))

	var hasRead, hasWrite bool

	for _, stmt := range statements {
		stmtMap, ok := stmt.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		statement, _ := stmtMap["Statement"].(string)
		if statement == "" {
			return nil, ErrInvalidParameter
		}

		params := parsePartiQLParams(stmtMap)

		upperStmt := strings.ToUpper(strings.TrimSpace(statement))
		var stmtType statementType
		switch {
		case strings.HasPrefix(upperStmt, "SELECT"):
			stmtType = statementTypeRead
			hasRead = true
		case strings.HasPrefix(upperStmt, "INSERT"),
			strings.HasPrefix(upperStmt, "UPDATE"),
			strings.HasPrefix(upperStmt, "DELETE"):
			stmtType = statementTypeWrite
			hasWrite = true
		default:
			return nil, ErrInvalidParameter
		}

		parsedStatements = append(parsedStatements, struct {
			statement string
			params    *partiQLParams
			stmtType  statementType
		}{statement, params, stmtType})
	}

	if hasRead && hasWrite {
		return nil, ErrTransactionConflict
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	responses := make([]map[string]interface{}, len(parsedStatements))

	if hasRead {
		err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for i, ps := range parsedStatements {
				result, err := s.executePartiQLSelectInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				if err != nil {
					return err
				}

				var item interface{} = nil
				if resultMap, ok := result.(map[string]interface{}); ok {
					if items, ok := resultMap["Items"].([]map[string]interface{}); ok && len(items) > 0 {
						item = items[0]
					}
				}
				responses[i] = map[string]interface{}{"Item": item}
			}
			return nil
		})
	} else {
		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for i, ps := range parsedStatements {
				upperStmt := strings.ToUpper(strings.TrimSpace(ps.statement))
				var result interface{}
				var execErr error

				switch {
				case strings.HasPrefix(upperStmt, "INSERT"):
					result, execErr = s.executePartiQLInsertInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				case strings.HasPrefix(upperStmt, "UPDATE"):
					result, execErr = s.executePartiQLUpdateInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				case strings.HasPrefix(upperStmt, "DELETE"):
					result, execErr = s.executePartiQLDeleteInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				}

				if execErr != nil {
					return execErr
				}

				var item interface{} = nil
				if resultMap, ok := result.(map[string]interface{}); ok {
					if items, ok := resultMap["Items"].([]map[string]interface{}); ok && len(items) > 0 {
						item = items[0]
					}
				}
				responses[i] = map[string]interface{}{"Item": item}
			}
			return nil
		})
	}

	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"Responses": responses,
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		tableNames := make(map[string]bool)
		for _, ps := range parsedStatements {
			if tn := extractTableNameFromStatement(ps.statement); tn != "" {
				tableNames[tn] = true
			}
		}
		var consumedCapacities []map[string]interface{}
		for tableName := range tableNames {
			consumedCapacities = append(consumedCapacities, buildConsumedCapacityResponse(tableName, 2.0))
		}
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp, nil
}

func (s *DynamoDBService) executePartiQLSelectInTxn(ctx context.Context, reqCtx *request.RequestContext, txn *dbstore.DynamoDBTxn, statement string, params *partiQLParams) (interface{}, error) {
	tableName, whereExpr := parseSelectStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	table, err := txn.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	var items []*dbstore.Item

	pkName := getPartitionKeyNameFromTable(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)

	if pkValue != "" {
		err = txn.ScanByPartitionKey(tableName, pkValue, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		err = txn.Scan(tableName, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	scannedCount := len(items)
	if whereExpr != nil {
		items = filterItemsByExpr(items, whereExpr, params)
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, buildItemResponse(item.Attributes))
	}

	return map[string]interface{}{
		"Items":        result,
		"Count":        len(result),
		"ScannedCount": scannedCount,
	}, nil
}

func (s *DynamoDBService) executePartiQLInsertInTxn(ctx context.Context, reqCtx *request.RequestContext, txn *dbstore.DynamoDBTxn, statement string, params *partiQLParams) (interface{}, error) {
	tableName, itemData := parseInsertStatement(statement)
	if tableName == "" || itemData == nil {
		return nil, ErrInvalidParameter
	}

	table, err := txn.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	keyAttrs := buildKeyFromSchema(table.KeySchema, itemData)
	if keyAttrs == nil {
		return nil, ErrInvalidParameter
	}

	_, err = txn.GetItem(tableName, keyAttrs)
	if err != nil {
		if !dbstore.IsItemNotFound(err) {
			return nil, err
		}
	} else {
		return nil, ErrConditionalCheckFailed
	}

	if err := txn.PutItem(tableName, keyAttrs, itemData); err != nil {
		return nil, err
	}

	newItem := &dbstore.Item{
		TableName:  tableName,
		Key:        keyAttrs,
		Attributes: itemData,
	}
	if err := txn.PutIndexEntries(tableName, newItem); err != nil {
		return nil, err
	}

	if err := txn.UpdateItemCount(tableName, 1); err != nil {
		return nil, err
	}

	newItemSize := calculateItemSize(itemData)
	if err := txn.UpdateTableSize(tableName, newItemSize); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        0,
		"ScannedCount": 0,
	}, nil
}

func (s *DynamoDBService) executePartiQLUpdateInTxn(ctx context.Context, reqCtx *request.RequestContext, txn *dbstore.DynamoDBTxn, statement string, params *partiQLParams) (interface{}, error) {
	tableName, assignments, whereExpr := parseUpdateStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	table, err := txn.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	var items []*dbstore.Item

	pkName := getPartitionKeyNameFromTable(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)

	if pkValue != "" {
		err = txn.ScanByPartitionKey(tableName, pkValue, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
	} else {
		err = txn.Scan(tableName, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
	}
	if err != nil {
		return nil, err
	}

	scannedCount := len(items)
	if whereExpr != nil {
		items = filterItemsByExpr(items, whereExpr, params)
	}

	updatedCount := 0
	for _, item := range items {
		oldSize := calculateItemSize(item.Attributes)

		oldItem := &dbstore.Item{
			TableName:  tableName,
			Key:        copyAttributes(item.Key),
			Attributes: copyAttributes(item.Attributes),
		}

		applySetAssignments(item.Attributes, assignments, params)

		newSize := calculateItemSize(item.Attributes)
		sizeDelta := newSize - oldSize

		if err := txn.DeleteIndexEntries(tableName, oldItem); err != nil {
			return nil, err
		}

		if err := txn.PutItem(tableName, item.Key, item.Attributes); err != nil {
			return nil, err
		}

		updatedItem := &dbstore.Item{
			TableName:  tableName,
			Key:        item.Key,
			Attributes: item.Attributes,
		}
		if err := txn.PutIndexEntries(tableName, updatedItem); err != nil {
			return nil, err
		}

		if sizeDelta != 0 {
			if err := txn.UpdateTableSize(tableName, sizeDelta); err != nil {
				return nil, err
			}
		}

		updatedCount++
	}

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        updatedCount,
		"ScannedCount": scannedCount,
	}, nil
}

func (s *DynamoDBService) executePartiQLDeleteInTxn(ctx context.Context, reqCtx *request.RequestContext, txn *dbstore.DynamoDBTxn, statement string, params *partiQLParams) (interface{}, error) {
	tableName, whereExpr := parseDeleteStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	table, err := txn.GetTable(tableName)
	if err != nil {
		return nil, err
	}

	var items []*dbstore.Item

	pkName := getPartitionKeyNameFromTable(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)

	if pkValue != "" {
		err = txn.ScanByPartitionKey(tableName, pkValue, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
	} else {
		err = txn.Scan(tableName, func(item *dbstore.Item) error {
			items = append(items, item)
			return nil
		})
	}
	if err != nil {
		return nil, err
	}

	scannedCount := len(items)
	if whereExpr != nil {
		items = filterItemsByExpr(items, whereExpr, params)
	}

	deletedCount := 0
	for _, item := range items {
		itemSize := calculateItemSize(item.Attributes)
		if err := txn.DeleteIndexEntries(tableName, item); err != nil {
			return nil, err
		}

		if err := txn.DeleteItem(tableName, item.Key); err != nil {
			return nil, err
		}

		if err := txn.UpdateItemCount(tableName, -1); err != nil {
			return nil, err
		}

		if err := txn.UpdateTableSize(tableName, -itemSize); err != nil {
			return nil, err
		}

		deletedCount++
	}

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        deletedCount,
		"ScannedCount": scannedCount,
	}, nil
}

func getPartitionKeyNameFromTable(table *dbstore.Table) string {
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			return ks.AttributeName
		}
	}
	return ""
}

// BatchExecuteStatement executes multiple PartiQL statements in a batch.
func (s *DynamoDBService) BatchExecuteStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statements, ok := req.Parameters["Statements"].([]interface{})
	if !ok {
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
