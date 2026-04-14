// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
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
		return nil, fmt.Errorf("transact get items: get store: %w", err)
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
				if !dbstore.IsItemNotFound(err) {
					return fmt.Errorf("transact get item on %s: %w", gi.tableName, err)
				}
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
		return nil, fmt.Errorf("transact get items: snapshot view: %w", err)
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

// writeOperation represents a single write operation within a transaction.
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
		return nil, fmt.Errorf("transact write items: get store: %w", err)
	}

	cancellationReasons := make([]CancellationReason, len(transactItems))
	for i := range cancellationReasons {
		cancellationReasons[i] = CancellationReason{Code: "None"}
	}

	operations, err := parseTransactWriteItems(s, store, transactItems, cancellationReasons)
	if err != nil {
		return nil, err
	}

	if err := executeTransactWriteItems(ctx, store, operations, cancellationReasons); err != nil {
		return nil, err
	}

	return buildTransactWriteResponse(req.Parameters, operations), nil
}

func parseTransactWriteItems(s *DynamoDBService, store dbstore.DynamoDBStoreInterface, transactItems []interface{}, cancellationReasons []CancellationReason) ([]writeOperation, error) {
	usedKeys := make(map[string]bool)
	var operations []writeOperation

	for idx, item := range transactItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
			return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
		}

		op, err := parseWriteOperation(s, store, idx, itemMap, usedKeys, cancellationReasons)
		if err != nil {
			return nil, err
		}
		if op != nil {
			operations = append(operations, *op)
		}
	}

	return operations, nil
}

func parseWriteOperation(s *DynamoDBService, store dbstore.DynamoDBStoreInterface, idx int, itemMap map[string]interface{}, usedKeys map[string]bool, cancellationReasons []CancellationReason) (*writeOperation, error) {
	for _, opType := range []string{"Put", "Update", "Delete", "ConditionCheck"} {
		opMap, ok := itemMap[opType].(map[string]interface{})
		if !ok {
			continue
		}

		tableName := request.GetStringParam(opMap, "TableName")
		if tableName == "" {
			cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
			return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
		}

		if _, err := store.Tables().Get(tableName); err != nil {
			cancellationReasons[idx] = CancellationReason{Code: "ResourceNotFound"}
			return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
		}

		key, err := extractOperationKey(s, store, opType, opMap, tableName)
		if err != nil {
			cancellationReasons[idx] = CancellationReason{Code: err.(*opParseError).code}
			return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
		}

		keyStr := buildKeyString(tableName, key)
		if opType != "ConditionCheck" {
			if usedKeys[keyStr] {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			usedKeys[keyStr] = true
		}

		op := &writeOperation{
			idx:                          idx,
			opType:                       opType,
			tableName:                    tableName,
			key:                          key,
			conditionExpr:                request.GetStringParam(opMap, "ConditionExpression"),
			exprAttrNames:                parseExpressionAttributeNames(opMap),
			exprAttrValues:               parseExpressionAttributeValues(opMap),
			returnValuesOnConditionCheck: request.GetStringParam(opMap, "ReturnValuesOnConditionCheckFailure"),
		}

		if opType == "Put" {
			itemData := parseItem(opMap["Item"])
			if itemData == nil {
				cancellationReasons[idx] = CancellationReason{Code: "ValidationError"}
				return nil, NewTransactionCanceledError("Transaction canceled", cancellationReasons)
			}
			op.itemData = itemData
		}

		if opType == "Update" {
			op.updateReq = opMap
		}

		return op, nil
	}

	return nil, nil
}

// opParseError represents a parse error encountered during transaction operation parsing.
type opParseError struct {
	code string
	err  error
}

// Error returns the underlying error message for a transaction operation parse failure.
func (e *opParseError) Error() string {
	return e.err.Error()
}

func extractOperationKey(s *DynamoDBService, store dbstore.DynamoDBStoreInterface, opType string, opMap map[string]interface{}, tableName string) (map[string]*dbstore.AttributeValue, error) {
	if opType == "Put" {
		itemData := parseItem(opMap["Item"])
		if itemData == nil {
			return nil, &opParseError{code: "ValidationError", err: fmt.Errorf("invalid item data")}
		}
		table, err := store.Tables().Get(tableName)
		if err != nil {
			return nil, &opParseError{code: "ResourceNotFound", err: err}
		}
		key := s.extractKeyFromItem(table, itemData)
		if key == nil {
			return nil, &opParseError{code: "ValidationError", err: fmt.Errorf("failed to extract key")}
		}
		return key, nil
	}

	key := parseKey(opMap["Key"])
	if key == nil {
		return nil, &opParseError{code: "ValidationError", err: fmt.Errorf("invalid key")}
	}
	return key, nil
}

func executeTransactWriteItems(ctx context.Context, store dbstore.DynamoDBStoreInterface, operations []writeOperation, cancellationReasons []CancellationReason) error {
	twoPhase := store.Storage().TwoPhaseTransaction()

	tableStore, ok := store.Tables().(*dbstore.TableStore)
	if !ok {
		return fmt.Errorf("transaction failed")
	}

	for _, op := range operations {
		op := op
		twoPhase.AddValidator(storage.ValidatorFunc(func(ctx context.Context, txn storage.Transaction) error {
			return validateWriteOperation(ctx, txn, tableStore, op, cancellationReasons)
		}))
	}

	for _, op := range operations {
		op := op
		twoPhase.AddExecutor(storage.ExecutorFunc(func(ctx context.Context, txn storage.Transaction) error {
			dbTxn := dbstore.NewDynamoDBTxn(txn, tableStore)
			return executeWriteOperation(dbTxn, op)
		}))
	}

	if err := twoPhase.Commit(ctx); err != nil {
		if _, ok := err.(*TransactionCanceledError); ok {
			return err
		}
		return ErrTransactionCanceled
	}

	return nil
}

func validateWriteOperation(_ context.Context, txn storage.Transaction, tableStore *dbstore.TableStore, op writeOperation, cancellationReasons []CancellationReason) error {
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
			return fmt.Errorf("validator get item %s: %w", op.tableName, err)
		}
	} else {
		item = existingItem
	}

	if op.conditionExpr != "" {
		conditionMet, err := evaluateConditionExpression(item, op.conditionExpr, op.exprAttrNames, op.exprAttrValues)
		if err != nil {
			return fmt.Errorf("validator condition check %s: %w", op.tableName, err)
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
}

func executeWriteOperation(dbTxn *dbstore.DynamoDBTxn, op writeOperation) error {
	exists, err := dbTxn.ItemExists(op.tableName, op.key)
	if err != nil {
		return fmt.Errorf("check item exists %s: %w", op.tableName, err)
	}

	switch op.opType {
	case "Put":
		return executePutOp(dbTxn, op, exists)
	case "Update":
		return executeUpdateOp(dbTxn, op, exists)
	case "Delete":
		return executeDeleteOp(dbTxn, op, exists)
	case "ConditionCheck":
		return nil
	}

	return nil
}

func executePutOp(dbTxn *dbstore.DynamoDBTxn, op writeOperation, exists bool) error {
	var oldItem *dbstore.Item
	var oldItemSize int64
	if exists {
		var err error
		oldItem, err = dbTxn.GetItem(op.tableName, op.key)
		if err != nil && !dbstore.IsItemNotFound(err) {
			return fmt.Errorf("put get old item %s: %w", op.tableName, err)
		}
		if oldItem != nil {
			oldItemSize = calculateItemSize(oldItem.Attributes)
			if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
				return fmt.Errorf("put delete index entries %s: %w", op.tableName, err)
			}
		}
	}
	if err := dbTxn.PutItem(op.tableName, op.key, op.itemData); err != nil {
		return fmt.Errorf("put item %s: %w", op.tableName, err)
	}
	newItem := &dbstore.Item{
		TableName:  op.tableName,
		Key:        op.key,
		Attributes: op.itemData,
	}
	if err := dbTxn.PutIndexEntries(op.tableName, newItem); err != nil {
		return fmt.Errorf("put index entries %s: %w", op.tableName, err)
	}
	return updateTableMetrics(dbTxn, op.tableName, exists, oldItemSize, calculateItemSize(op.itemData))
}

func executeUpdateOp(dbTxn *dbstore.DynamoDBTxn, op writeOperation, exists bool) error {
	var oldItem *dbstore.Item
	var oldItemSize int64
	if exists {
		var err error
		oldItem, err = dbTxn.GetItem(op.tableName, op.key)
		if err != nil && !dbstore.IsItemNotFound(err) {
			return fmt.Errorf("update get old item %s: %w", op.tableName, err)
		}
		if oldItem != nil {
			oldItemSize = calculateItemSize(oldItem.Attributes)
			if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
				return fmt.Errorf("update delete index entries %s: %w", op.tableName, err)
			}
		}
	}

	attrs := make(map[string]*dbstore.AttributeValue)
	if oldItem != nil {
		attrs = copyAttributes(oldItem.Attributes)
	} else {
		for k, v := range op.key {
			attrs[k] = v
		}
	}

	updateExpr := request.GetStringParam(op.updateReq, "UpdateExpression")
	if updateExpr != "" {
		if err := applyUpdateExpression(attrs, updateExpr, op.exprAttrNames, op.exprAttrValues); err != nil {
			return fmt.Errorf("apply update expression %s: %w", op.tableName, err)
		}
	}

	if err := dbTxn.PutItem(op.tableName, op.key, attrs); err != nil {
		return fmt.Errorf("update put item %s: %w", op.tableName, err)
	}
	newItem := &dbstore.Item{
		TableName:  op.tableName,
		Key:        op.key,
		Attributes: attrs,
	}
	if err := dbTxn.PutIndexEntries(op.tableName, newItem); err != nil {
		return fmt.Errorf("update put index entries %s: %w", op.tableName, err)
	}
	return updateTableMetrics(dbTxn, op.tableName, exists, oldItemSize, calculateItemSize(attrs))
}

func executeDeleteOp(dbTxn *dbstore.DynamoDBTxn, op writeOperation, exists bool) error {
	var oldItem *dbstore.Item
	var oldItemSize int64
	if exists {
		var err error
		oldItem, err = dbTxn.GetItem(op.tableName, op.key)
		if err != nil && !dbstore.IsItemNotFound(err) {
			return fmt.Errorf("delete get old item %s: %w", op.tableName, err)
		}
		if oldItem != nil {
			oldItemSize = calculateItemSize(oldItem.Attributes)
			if err := dbTxn.DeleteIndexEntries(op.tableName, oldItem); err != nil {
				return fmt.Errorf("delete index entries %s: %w", op.tableName, err)
			}
		}
	}
	if err := dbTxn.DeleteItem(op.tableName, op.key); err != nil {
		return fmt.Errorf("delete item %s: %w", op.tableName, err)
	}
	if exists {
		if err := dbTxn.UpdateItemCount(op.tableName, -1); err != nil {
			return fmt.Errorf("delete decrement item count %s: %w", op.tableName, err)
		}
		if oldItemSize > 0 {
			if err := dbTxn.UpdateTableSize(op.tableName, -oldItemSize); err != nil {
				return fmt.Errorf("delete adjust table size %s: %w", op.tableName, err)
			}
		}
	}
	return nil
}

func updateTableMetrics(dbTxn *dbstore.DynamoDBTxn, tableName string, existed bool, oldSize, newSize int64) error {
	if !existed {
		if err := dbTxn.UpdateItemCount(tableName, 1); err != nil {
			return fmt.Errorf("increment item count %s: %w", tableName, err)
		}
		if err := dbTxn.UpdateTableSize(tableName, newSize); err != nil {
			return fmt.Errorf("update table size %s: %w", tableName, err)
		}
	} else if newSize != oldSize {
		if err := dbTxn.UpdateTableSize(tableName, newSize-oldSize); err != nil {
			return fmt.Errorf("adjust table size %s: %w", tableName, err)
		}
	}
	return nil
}

func buildTransactWriteResponse(params map[string]interface{}, operations []writeOperation) map[string]interface{} {
	resp := map[string]interface{}{}

	returnConsumedCapacity := getReturnConsumedCapacity(params)
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

	return resp
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
