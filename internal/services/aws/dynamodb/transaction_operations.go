// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/services/aws/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TransactGetItems performs multiple GetItem operations in a single transaction with snapshot isolation.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactGetItems.html
func (s *DynamoDBService) TransactGetItems(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	transactItems, ok := req.Parameters["TransactItems"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	if len(transactItems) > 100 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	type getItem struct {
		tableName  string
		key        map[string]*dbstore.AttributeValue
		projection []string
	}

	var getItems []getItem
	for _, item := range transactItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		getMap, ok := itemMap["Get"].(map[string]interface{})
		if !ok {
			continue
		}

		tableName := request.GetStringParam(getMap, "TableName")
		if tableName == "" {
			continue
		}

		key := parseKey(getMap["Key"])
		if key == nil {
			continue
		}

		projection := parseProjectionExpression(getMap)

		getItems = append(getItems, getItem{
			tableName:  tableName,
			key:        key,
			projection: projection,
		})
	}

	var responses []map[string]interface{}

	err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
		for _, gi := range getItems {
			dbItem, err := txn.GetItem(gi.tableName, gi.key)
			if err != nil {
				responses = append(responses, map[string]interface{}{"Item": nil})
				continue
			}

			attrs := dbItem.Attributes
			if gi.projection != nil {
				attrs = applyProjection(attrs, gi.projection)
			}

			responses = append(responses, map[string]interface{}{
				"Item": buildItemResponse(attrs),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"Responses": responses,
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		tableNames := make(map[string]bool)
		for _, gi := range getItems {
			tableNames[gi.tableName] = true
		}
		var consumedCapacities []map[string]interface{}
		for tableName := range tableNames {
			consumedCapacities = append(consumedCapacities, buildConsumedCapacityResponse(tableName, float64(len(responses))*0.5))
		}
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp, nil
}

// TransactWriteItems performs multiple write operations in a single transaction with ACID semantics.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactWriteItems.html
func (s *DynamoDBService) TransactWriteItems(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	transactItems, ok := req.Parameters["TransactItems"].([]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	if len(transactItems) > 100 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	type writeOperation struct {
		idx                          int
		opType                       string
		tableName                    string
		key                          map[string]*dbstore.AttributeValue
		itemData                     map[string]*dbstore.AttributeValue
		updateReq                    map[string]interface{}
		conditionExpr                string
		exprAttrNames                map[string]string
		exprAttrValues               map[string]*dbstore.AttributeValue
		returnValuesOnConditionCheck string
	}

	var operations []writeOperation
	cancellationReasons := make([]CancellationReason, len(transactItems))
	for i := range cancellationReasons {
		cancellationReasons[i] = CancellationReason{Code: "None"}
	}
	usedKeys := make(map[string]bool)

	for idx, item := range transactItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
			return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
		}

		if putReq, ok := itemMap["Put"].(map[string]interface{}); ok {
			tableName := request.GetStringParam(putReq, "TableName")
			if tableName == "" {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			table, err := store.Tables().Get(tableName)
			if err != nil {
				cancellationReasons[idx] = CancellationReason{Code: "ResourceNotFound"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			itemData := parseItem(putReq["Item"])
			if itemData == nil {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			key := s.extractKeyFromItem(table, itemData)
			if key == nil {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			keyStr := buildKeyString(tableName, key)
			if usedKeys[keyStr] {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			usedKeys[keyStr] = true

			operations = append(operations, writeOperation{
				idx:                          idx,
				opType:                       "Put",
				tableName:                    tableName,
				key:                          key,
				itemData:                     itemData,
				conditionExpr:                request.GetStringParam(putReq, "ConditionExpression"),
				exprAttrNames:                parseExpressionAttributeNames(putReq),
				exprAttrValues:               parseExpressionAttributeValues(putReq),
				returnValuesOnConditionCheck: request.GetStringParam(putReq, "ReturnValuesOnConditionCheckFailure"),
			})
		}

		if updateReq, ok := itemMap["Update"].(map[string]interface{}); ok {
			tableName := request.GetStringParam(updateReq, "TableName")
			if tableName == "" {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			if _, err := store.Tables().Get(tableName); err != nil {
				cancellationReasons[idx] = CancellationReason{Code: "ResourceNotFound"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			key := parseKey(updateReq["Key"])
			if key == nil {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			keyStr := buildKeyString(tableName, key)
			if usedKeys[keyStr] {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			usedKeys[keyStr] = true

			operations = append(operations, writeOperation{
				idx:                          idx,
				opType:                       "Update",
				tableName:                    tableName,
				key:                          key,
				updateReq:                    updateReq,
				conditionExpr:                request.GetStringParam(updateReq, "ConditionExpression"),
				exprAttrNames:                parseExpressionAttributeNames(updateReq),
				exprAttrValues:               parseExpressionAttributeValues(updateReq),
				returnValuesOnConditionCheck: request.GetStringParam(updateReq, "ReturnValuesOnConditionCheckFailure"),
			})
		}

		if deleteReq, ok := itemMap["Delete"].(map[string]interface{}); ok {
			tableName := request.GetStringParam(deleteReq, "TableName")
			if tableName == "" {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			if _, err := store.Tables().Get(tableName); err != nil {
				cancellationReasons[idx] = CancellationReason{Code: "ResourceNotFound"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			key := parseKey(deleteReq["Key"])
			if key == nil {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			keyStr := buildKeyString(tableName, key)
			if usedKeys[keyStr] {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			usedKeys[keyStr] = true

			operations = append(operations, writeOperation{
				idx:                          idx,
				opType:                       "Delete",
				tableName:                    tableName,
				key:                          key,
				conditionExpr:                request.GetStringParam(deleteReq, "ConditionExpression"),
				exprAttrNames:                parseExpressionAttributeNames(deleteReq),
				exprAttrValues:               parseExpressionAttributeValues(deleteReq),
				returnValuesOnConditionCheck: request.GetStringParam(deleteReq, "ReturnValuesOnConditionCheckFailure"),
			})
		}

		if checkReq, ok := itemMap["ConditionCheck"].(map[string]interface{}); ok {
			tableName := request.GetStringParam(checkReq, "TableName")
			if tableName == "" {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			if _, err := store.Tables().Get(tableName); err != nil {
				cancellationReasons[idx] = CancellationReason{Code: "ResourceNotFound"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			key := parseKey(checkReq["Key"])
			if key == nil {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}

			keyStr := buildKeyString(tableName, key)
			if usedKeys[keyStr] {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			usedKeys[keyStr] = true

			operations = append(operations, writeOperation{
				idx:                          idx,
				opType:                       "ConditionCheck",
				tableName:                    tableName,
				key:                          key,
				conditionExpr:                request.GetStringParam(checkReq, "ConditionExpression"),
				exprAttrNames:                parseExpressionAttributeNames(checkReq),
				exprAttrValues:               parseExpressionAttributeValues(checkReq),
				returnValuesOnConditionCheck: request.GetStringParam(checkReq, "ReturnValuesOnConditionCheckFailure"),
			})
		}
	}

	twoPhase := store.Storage().TwoPhaseTransaction()

	tableStore, ok := store.Tables().(*dbstore.TableStore)
	if !ok {
		return nil, fmt.Errorf("transaction failed")
	}

	for _, op := range operations {
		op := op
		twoPhase.AddValidator(storage.ValidatorFunc(func(ctx context.Context, txn storage.Transaction) error {
			dbTxn := dbstore.NewDynamoDBTxn(txn, tableStore)

			var item *dbstore.Item
			itemExists := true

			existingItem, err := dbTxn.GetItem(op.tableName, op.key)
			if err != nil {
				if dbstore.IsItemNotFound(err) {
					itemExists = false
					item = &dbstore.Item{
						TableName:  op.tableName,
						Key:        op.key,
						Attributes: make(map[string]*dbstore.AttributeValue),
					}
				} else {
					return err
				}
			} else {
				item = existingItem
			}

			if op.conditionExpr != "" {
				conditionMet, err := evaluateConditionExpression(item, op.conditionExpr, op.exprAttrNames, op.exprAttrValues)
				if err != nil {
					return err
				}
				if !conditionMet {
					reason := CancellationReason{Code: "ConditionalCheckFailed"}
					if op.returnValuesOnConditionCheck == "ALL_OLD" && item != nil && itemExists {
						reason.Item = buildItemResponse(item.Attributes)
					}
					cancellationReasons[op.idx] = reason
					return NewTransactionCanceledError("Transaction canceled", cancellationReasons)
				}
			}

			return nil
		}))
	}

	for _, op := range operations {
		op := op
		twoPhase.AddExecutor(storage.ExecutorFunc(func(ctx context.Context, txn storage.Transaction) error {
			dbTxn := dbstore.NewDynamoDBTxn(txn, tableStore)

			exists, err := dbTxn.ItemExists(op.tableName, op.key)
			if err != nil {
				return err
			}

			switch op.opType {
			case "Put":
				var oldItem *dbstore.Item
				var oldItemSize int64
				if exists {
					oldItem, err = dbTxn.GetItem(op.tableName, op.key)
					if err != nil && !dbstore.IsItemNotFound(err) {
						return err
					}
					if oldItem != nil {
						oldItemSize = calculateItemSize(oldItem.Attributes)
						if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
							return err
						}
					}
				}
				if err := dbTxn.PutItem(op.tableName, op.key, op.itemData); err != nil {
					return err
				}
				newItem := &dbstore.Item{
					TableName:  op.tableName,
					Key:        op.key,
					Attributes: op.itemData,
				}
				if err := dbTxn.PutIndexEntries(op.tableName, newItem); err != nil {
					return err
				}
				newItemSize := calculateItemSize(op.itemData)
				if !exists {
					if err := dbTxn.UpdateItemCount(op.tableName, 1); err != nil {
						return err
					}
					if err := dbTxn.UpdateTableSize(op.tableName, newItemSize); err != nil {
						return err
					}
				} else {
					if newItemSize != oldItemSize {
						if err := dbTxn.UpdateTableSize(op.tableName, newItemSize-oldItemSize); err != nil {
							return err
						}
					}
				}

			case "Update":
				var oldItem *dbstore.Item
				var oldItemSize int64
				if exists {
					oldItem, err = dbTxn.GetItem(op.tableName, op.key)
					if err != nil && !dbstore.IsItemNotFound(err) {
						return err
					}
					if oldItem != nil {
						oldItemSize = calculateItemSize(oldItem.Attributes)
						if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
							return err
						}
					}
				}
				var attrs map[string]*dbstore.AttributeValue
				if oldItem != nil {
					attrs = copyAttributes(oldItem.Attributes)
				} else {
					attrs = make(map[string]*dbstore.AttributeValue)
					for k, v := range op.key {
						attrs[k] = v
					}
				}

				updateExpr := request.GetStringParam(op.updateReq, "UpdateExpression")
				if updateExpr != "" {
					if err := applyUpdateExpression(attrs, updateExpr, op.exprAttrNames, op.exprAttrValues); err != nil {
						return err
					}
				}

				if err := dbTxn.PutItem(op.tableName, op.key, attrs); err != nil {
					return err
				}
				newItem := &dbstore.Item{
					TableName:  op.tableName,
					Key:        op.key,
					Attributes: attrs,
				}
				if err := dbTxn.PutIndexEntries(op.tableName, newItem); err != nil {
					return err
				}
				newItemSize := calculateItemSize(attrs)
				if !exists {
					if err := dbTxn.UpdateItemCount(op.tableName, 1); err != nil {
						return err
					}
					if err := dbTxn.UpdateTableSize(op.tableName, newItemSize); err != nil {
						return err
					}
				} else {
					if newItemSize != oldItemSize {
						if err := dbTxn.UpdateTableSize(op.tableName, newItemSize-oldItemSize); err != nil {
							return err
						}
					}
				}

			case "Delete":
				var oldItem *dbstore.Item
				var oldItemSize int64
				if exists {
					oldItem, err = dbTxn.GetItem(op.tableName, op.key)
					if err != nil && !dbstore.IsItemNotFound(err) {
						return err
					}
					if oldItem != nil {
						oldItemSize = calculateItemSize(oldItem.Attributes)
						if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
							return err
						}
					}
				}
				if err := dbTxn.DeleteItem(op.tableName, op.key); err != nil {
					return err
				}
				if exists {
					if err := dbTxn.UpdateItemCount(op.tableName, -1); err != nil {
						return err
					}
					if oldItemSize > 0 {
						if err := dbTxn.UpdateTableSize(op.tableName, -oldItemSize); err != nil {
							return err
						}
					}
				}
			}

			return nil
		}))
	}

	if err := twoPhase.Commit(ctx); err != nil {
		if _, ok := err.(*TransactionCanceledError); ok {
			return nil, err
		}
		return nil, ErrTransactionCanceled
	}

	resp := map[string]interface{}{}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		tableNames := make(map[string]bool)
		for _, op := range operations {
			tableNames[op.tableName] = true
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

func copyAttributes(attrs map[string]*dbstore.AttributeValue) map[string]*dbstore.AttributeValue {
	if attrs == nil {
		return nil
	}
	cpy := make(map[string]*dbstore.AttributeValue)
	for k, v := range attrs {
		cpy[k] = deepCopyAttributeValue(v)
	}
	return cpy
}

func deepCopyAttributeValue(v *dbstore.AttributeValue) *dbstore.AttributeValue {
	if v == nil {
		return nil
	}
	cpy := &dbstore.AttributeValue{}
	if v.S != nil {
		s := *v.S
		cpy.S = &s
	}
	if v.N != nil {
		n := *v.N
		cpy.N = &n
	}
	if v.B != nil {
		cpy.B = make([]byte, len(v.B))
		copy(cpy.B, v.B)
	}
	if v.SS != nil {
		cpy.SS = make([]string, len(v.SS))
		copy(cpy.SS, v.SS)
	}
	if v.NS != nil {
		cpy.NS = make([]string, len(v.NS))
		copy(cpy.NS, v.NS)
	}
	if v.BS != nil {
		cpy.BS = make([][]byte, len(v.BS))
		for i, b := range v.BS {
			cpy.BS[i] = make([]byte, len(b))
			copy(cpy.BS[i], b)
		}
	}
	if v.M != nil {
		cpy.M = make(map[string]*dbstore.AttributeValue)
		for k, val := range v.M {
			cpy.M[k] = deepCopyAttributeValue(val)
		}
	}
	if v.L != nil {
		cpy.L = make([]*dbstore.AttributeValue, len(v.L))
		for i, val := range v.L {
			cpy.L[i] = deepCopyAttributeValue(val)
		}
	}
	if v.NULL != nil {
		null := *v.NULL
		cpy.NULL = &null
	}
	if v.BOOL != nil {
		b := *v.BOOL
		cpy.BOOL = &b
	}
	return cpy
}

func buildKeyString(tableName string, key map[string]*dbstore.AttributeValue) string {
	names := make([]string, 0, len(key))
	for k := range key {
		names = append(names, k)
	}
	sort.Strings(names)

	result := tableName + "#"
	for _, k := range names {
		v := key[k]
		result += k + "="
		if v.S != nil {
			result += "S:" + escapeKeyPart(*v.S)
		} else if v.N != nil {
			result += "N:" + *v.N
		} else if v.B != nil {
			result += "B:" + base64.StdEncoding.EncodeToString(v.B)
		}
		result += ";"
	}
	return result
}

func escapeKeyPart(s string) string {
	result := ""
	for _, c := range s {
		switch c {
		case '\\', ';', '=', '#':
			result += "\\" + string(c)
		default:
			result += string(c)
		}
	}
	return result
}
