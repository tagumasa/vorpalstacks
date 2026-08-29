package dynamodb

// This file carries the PartiQL execution engines: the per-statement-type
// Cores backing ExecuteStatement and ExecuteTransaction, plus the SET clause
// appliers they share. The handlers in partiql_operations.go and
// partiql_transaction.go parse and dispatch; everything here owns validation
// and persistence.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func (s *DynamoDBService) executePartiQLSelectEnhanced(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams, consistentRead bool, limit int, nextToken string) (interface{}, error) {
	tableName, whereExpr, orderBy, selectCols := parseSelectStatementWithOrderBy(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if !store.Tables().Exists(tableName) {
		return nil, ErrTableNotFound
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}

	pkName := getHashKeyName(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)

	// Decode pagination offset before scanning so we can calculate
	// how many items to collect before early termination.
	startOffset := 0
	if nextToken != "" {
		if decoded, decErr := base64.StdEncoding.DecodeString(nextToken); decErr == nil {
			var offset int
			if json.Unmarshal(decoded, &offset) == nil && offset > 0 {
				startOffset = offset
			}
		}
	}

	// When no ORDER BY is present, filter inside the scan callback and
	// break early after collecting enough items. This reduces memory
	// from O(N) to O(limit) for filtered queries on large tables.
	needed := 0
	if orderBy == nil && limit > 0 {
		needed = startOffset + limit
	}

	var items []*dbstore.Item
	scannedCount := 0

	scanCallback := func(item *dbstore.Item) error {
		scannedCount++
		if whereExpr != nil && !evaluateExpr(item.Attributes, whereExpr, params) {
			return nil
		}
		if needed > 0 && len(items) >= needed {
			return errScanSufficient
		}
		items = append(items, item)
		return nil
	}

	if pkValue != "" {
		err = store.Items().ScanByPartitionKey(tableName, pkValue, scanCallback)
	} else {
		err = store.Items().Scan(tableName, scanCallback)
	}
	if err != nil && !errors.Is(err, errScanSufficient) {
		return nil, err
	}

	if orderBy != nil {
		items = sortItemsByOrderBy(items, orderBy)
	}
	items = items[startOffset:]

	var newNextToken string
	if limit > 0 && len(items) > limit {
		items = items[:limit]
		offsetBytes, _ := json.Marshal(startOffset + limit)
		newNextToken = base64.StdEncoding.EncodeToString(offsetBytes)
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		attrs := item.Attributes
		if selectCols != nil {
			filtered := make(map[string]*dbstore.AttributeValue)
			for _, col := range selectCols {
				if v, ok := attrs[col]; ok {
					filtered[col] = v
				}
			}
			attrs = filtered
		}
		result = append(result, buildItemResponse(attrs))
	}

	_ = consistentRead

	resp := map[string]interface{}{
		"Items":        result,
		"Count":        len(result),
		"ScannedCount": scannedCount,
	}
	if newNextToken != "" {
		resp["NextToken"] = newNextToken
	}
	return resp, nil
}

func (s *DynamoDBService) executePartiQLInsert(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams) (interface{}, error) {
	tableName, itemData := parseInsertStatementWithParams(statement, params)
	if tableName == "" || itemData == nil {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if !store.Tables().Exists(tableName) {
		return nil, ErrTableNotFound
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}

	keyAttrs := buildKeyFromSchema(table.KeySchema, itemData)
	if keyAttrs == nil {
		return nil, ErrInvalidParameter
	}

	var isNewItem bool

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		_, err := txn.GetItem(tableName, keyAttrs)
		if err != nil {
			if !dbstore.IsItemNotFound(err) {
				return err
			}
			isNewItem = true
		} else {
			return ErrConditionalCheckFailed
		}

		if err := txn.PutItem(tableName, keyAttrs, itemData); err != nil {
			return err
		}

		newItem := &dbstore.Item{
			TableName:  tableName,
			Key:        keyAttrs,
			Attributes: itemData,
		}
		if err := txn.PutIndexEntries(tableName, newItem); err != nil {
			return err
		}

		if isNewItem {
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			newItemSize := calculateItemSize(itemData)
			if err := txn.UpdateTableSize(tableName, newItemSize); err != nil {
				return err
			}
		}

		s.captureStreamChangeTxn(txn, store, table, dbstore.StreamEventInsert, keyAttrs, itemData, nil)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *DynamoDBService) executePartiQLUpdate(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams, returnValuesOnConditionCheckFailure string) (ret interface{}, err error) {
	defer func() {
		if r := resilience.RecoverPanic("dynamodb partiql update"); r != nil {
			err = fmt.Errorf("panic in executePartiQLUpdate: %v", r)
		}
	}()

	tableName, clauses, whereExpr := parseUpdateStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	// At least one clause must be present.
	if len(clauses.setAssignments) == 0 && len(clauses.removeAttrs) == 0 &&
		len(clauses.addAssignments) == 0 && len(clauses.deleteAssignments) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if !store.Tables().Exists(tableName) {
		return nil, ErrTableNotFound
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}

	// AWS PartiQL requires UPDATE to target a single item via partition-key
	// equality in the WHERE clause. Reject statements that lack this.
	if whereExpr == nil {
		return nil, ErrInvalidParameter
	}
	pkName := getHashKeyName(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)
	if pkValue == "" {
		return nil, ErrInvalidParameter
	}

	// Scan only items matching the partition key, then filter by the
	// full WHERE clause in the callback. Collect at most 2 items to
	// detect multi-item matches without loading the entire partition.
	var items []*dbstore.Item
	var preFilterItems []*dbstore.Item
	scannedCount := 0
	scanCallback := func(item *dbstore.Item) error {
		scannedCount++
		if !evaluateExpr(item.Attributes, whereExpr, params) {
			if returnValuesOnConditionCheckFailure == "ALL_OLD" {
				preFilterItems = append(preFilterItems, item)
			}
			return nil
		}
		items = append(items, item)
		if len(items) >= 2 {
			return errScanSufficient
		}
		return nil
	}
	err = store.Items().ScanByPartitionKey(tableName, pkValue, scanCallback)
	if err != nil && !errors.Is(err, errScanSufficient) {
		return nil, err
	}

	// AWS rejects UPDATE statements that match more than one item.
	if len(items) > 1 {
		return nil, NewAPIError("com.amazonaws.dynamodb.v20120810#ValidationException",
			"UPDATE statement must match exactly one item", http.StatusBadRequest)
	}

	if len(items) == 0 && len(preFilterItems) > 0 {
		oldItems := make([]map[string]interface{}, 0, len(preFilterItems))
		for _, pi := range preFilterItems {
			oldItems = append(oldItems, buildItemResponse(pi.Attributes))
		}
		return map[string]interface{}{
			"Items":        oldItems,
			"Count":        0,
			"ScannedCount": scannedCount,
		}, nil
	}

	updatedCount := 0
	for _, item := range items {
		oldSize := calculateItemSize(item.Attributes)

		oldItem := &dbstore.Item{
			TableName:  tableName,
			Key:        copyAttributes(item.Key),
			Attributes: copyAttributes(item.Attributes),
		}

		applySetAssignments(item.Attributes, clauses.setAssignments, params)
		applyRemoveAttrs(item.Attributes, clauses.removeAttrs)
		if err := applyAddAssignments(item.Attributes, clauses.addAssignments, params); err != nil {
			return nil, err
		}
		if err := applyDeleteAssignments(item.Attributes, clauses.deleteAssignments, params); err != nil {
			return nil, err
		}

		newSize := calculateItemSize(item.Attributes)
		sizeDelta := newSize - oldSize

		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			if err := txn.DeleteIndexEntries(tableName, oldItem); err != nil {
				return err
			}

			if err := txn.PutItem(tableName, item.Key, item.Attributes); err != nil {
				return err
			}

			updatedItem := &dbstore.Item{
				TableName:  tableName,
				Key:        item.Key,
				Attributes: item.Attributes,
			}
			if err := txn.PutIndexEntries(tableName, updatedItem); err != nil {
				return err
			}

			if sizeDelta != 0 {
				if err := txn.UpdateTableSize(tableName, sizeDelta); err != nil {
					return err
				}
			}

			s.captureStreamChangeTxn(txn, store, table, dbstore.StreamEventModify, item.Key, item.Attributes, oldItem.Attributes)

			return nil
		})
		if err != nil {
			return nil, err
		}
		updatedCount++
	}

	return map[string]interface{}{
		"Items":        []map[string]interface{}{},
		"Count":        updatedCount,
		"ScannedCount": scannedCount,
	}, nil
}

func (s *DynamoDBService) executePartiQLDelete(ctx context.Context, reqCtx *request.RequestContext, statement string, params *partiQLParams, returnValuesOnConditionCheckFailure string) (ret interface{}, err error) {
	defer func() {
		if r := resilience.RecoverPanic("dynamodb partiql delete"); r != nil {
			err = fmt.Errorf("panic in executePartiQLDelete: %v", r)
		}
	}()

	tableName, whereExpr := parseDeleteStatement(statement)
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if !store.Tables().Exists(tableName) {
		return nil, ErrTableNotFound
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, err
	}

	// AWS PartiQL requires DELETE to target a single item via partition-key
	// equality in the WHERE clause. Reject statements that lack this.
	if whereExpr == nil {
		return nil, ErrInvalidParameter
	}
	pkName := getHashKeyName(table)
	pkValue := extractPartitionKeyFromWhere(whereExpr, pkName, params)
	if pkValue == "" {
		return nil, ErrInvalidParameter
	}

	// Scan only items matching the partition key, then filter by the
	// full WHERE clause in the callback. Collect at most 2 items to
	// detect multi-item matches without loading the entire partition.
	var items []*dbstore.Item
	var preFilterItems []*dbstore.Item
	scannedCount := 0
	scanCallback := func(item *dbstore.Item) error {
		scannedCount++
		if !evaluateExpr(item.Attributes, whereExpr, params) {
			if returnValuesOnConditionCheckFailure == "ALL_OLD" {
				preFilterItems = append(preFilterItems, item)
			}
			return nil
		}
		items = append(items, item)
		if len(items) >= 2 {
			return errScanSufficient
		}
		return nil
	}
	err = store.Items().ScanByPartitionKey(tableName, pkValue, scanCallback)
	if err != nil && !errors.Is(err, errScanSufficient) {
		return nil, err
	}

	// AWS rejects DELETE statements that match more than one item.
	if len(items) > 1 {
		return nil, NewAPIError("com.amazonaws.dynamodb.v20120810#ValidationException",
			"DELETE statement must match exactly one item", http.StatusBadRequest)
	}

	if len(items) == 0 && len(preFilterItems) > 0 {
		oldItems := make([]map[string]interface{}, 0, len(preFilterItems))
		for _, pi := range preFilterItems {
			oldItems = append(oldItems, buildItemResponse(pi.Attributes))
		}
		return map[string]interface{}{
			"Items":        oldItems,
			"Count":        0,
			"ScannedCount": scannedCount,
		}, nil
	}

	deletedCount := 0
	for _, item := range items {
		itemSize := calculateItemSize(item.Attributes)
		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			if err := txn.DeleteIndexEntries(tableName, item); err != nil {
				return err
			}

			if err := txn.DeleteItem(tableName, item.Key); err != nil {
				return err
			}

			if err := txn.UpdateItemCount(tableName, -1); err != nil {
				return err
			}

			if err := txn.UpdateTableSize(tableName, -itemSize); err != nil {
				return err
			}

			s.captureStreamChangeTxn(txn, store, table, dbstore.StreamEventRemove, item.Key, nil, item.Attributes)

			return nil
		})
		if err != nil {
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

// applyAddAssignments performs DynamoDB ADD operations: if the attribute
// is a number, add the value numerically; if it is a set, add elements.
// Returns ErrInvalidParameter when the existing attribute and the ADD
// value have incompatible types.
func applyAddAssignments(attrs map[string]*dbstore.AttributeValue, assignments []setAssignment, params *partiQLParams) error {
	for _, asgn := range assignments {
		existing := attrs[asgn.attrName]
		addValue := exprToAttributeValueWithParams(asgn.value, params)
		if addValue == nil {
			continue
		}

		if existing == nil {
			// Attribute does not exist — ADD creates it.
			attrs[asgn.attrName] = addValue
			continue
		}

		// Numeric ADD.
		if existing.N != nil && addValue.N != nil {
			result, addErr := addNumbers(*existing.N, *addValue.N)
			if addErr != nil {
				return addErr
			}
			r := result
			attrs[asgn.attrName] = &dbstore.AttributeValue{N: &r}
			continue
		}

		// String set ADD.
		if existing.SS != nil && addValue.SS != nil {
			setMap := make(map[string]bool)
			for _, s := range existing.SS {
				setMap[s] = true
			}
			for _, s := range addValue.SS {
				setMap[s] = true
			}
			merged := make([]string, 0, len(setMap))
			for s := range setMap {
				merged = append(merged, s)
			}
			attrs[asgn.attrName] = dbstore.StringSet(merged)
			continue
		}

		// Number set ADD.
		if existing.NS != nil && addValue.NS != nil {
			setMap := make(map[string]bool)
			for _, n := range existing.NS {
				setMap[n] = true
			}
			for _, n := range addValue.NS {
				setMap[n] = true
			}
			merged := make([]string, 0, len(setMap))
			for n := range setMap {
				merged = append(merged, n)
			}
			attrs[asgn.attrName] = dbstore.NumberSet(merged)
			continue
		}

		// Binary set ADD.
		if existing.BS != nil && addValue.BS != nil {
			existingSet := make(map[string]bool)
			for _, b := range existing.BS {
				existingSet[string(b)] = true
			}
			for _, b := range addValue.BS {
				if !existingSet[string(b)] {
					existing.BS = append(existing.BS, b)
				}
			}
			attrs[asgn.attrName] = dbstore.BinarySet(existing.BS)
			continue
		}

		// No compatible type pair matched — type mismatch.
		return ErrInvalidParameter
	}
	return nil
}

// applyDeleteAssignments performs DynamoDB DELETE operations: remove
// elements from a set attribute (SS, NS, or BS). Returns
// ErrInvalidParameter when the existing attribute and the DELETE value
// have incompatible types.
func applyDeleteAssignments(attrs map[string]*dbstore.AttributeValue, assignments []setAssignment, params *partiQLParams) error {
	for _, asgn := range assignments {
		existing := attrs[asgn.attrName]
		delValue := exprToAttributeValueWithParams(asgn.value, params)
		if delValue == nil || existing == nil {
			continue
		}

		// String set DELETE.
		if existing.SS != nil && delValue.SS != nil {
			delSet := make(map[string]bool)
			for _, s := range delValue.SS {
				delSet[s] = true
			}
			var remaining []string
			for _, s := range existing.SS {
				if !delSet[s] {
					remaining = append(remaining, s)
				}
			}
			if len(remaining) == 0 {
				delete(attrs, asgn.attrName)
			} else {
				attrs[asgn.attrName] = dbstore.StringSet(remaining)
			}
			continue
		}

		// Number set DELETE.
		if existing.NS != nil && delValue.NS != nil {
			delSet := make(map[string]bool)
			for _, n := range delValue.NS {
				delSet[n] = true
			}
			var remaining []string
			for _, n := range existing.NS {
				if !delSet[n] {
					remaining = append(remaining, n)
				}
			}
			if len(remaining) == 0 {
				delete(attrs, asgn.attrName)
			} else {
				attrs[asgn.attrName] = dbstore.NumberSet(remaining)
			}
			continue
		}

		// Binary set DELETE.
		if existing.BS != nil && delValue.BS != nil {
			delSet := make(map[string]bool)
			for _, b := range delValue.BS {
				delSet[string(b)] = true
			}
			var remaining [][]byte
			for _, b := range existing.BS {
				if !delSet[string(b)] {
					remaining = append(remaining, b)
				}
			}
			if len(remaining) == 0 {
				delete(attrs, asgn.attrName)
			} else {
				attrs[asgn.attrName] = dbstore.BinarySet(remaining)
			}
			continue
		}

		// No compatible type pair matched — type mismatch.
		return ErrInvalidParameter
	}
	return nil
}
