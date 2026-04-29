package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// PutItem creates a new item or replaces an existing item in the specified table.
func (s *DynamoDBService) PutItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetActiveTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	item := parseItem(req.Parameters["Item"])
	if item == nil {
		return nil, ErrInvalidParameter
	}

	if itemSize := calculateItemSize(item); itemSize > maxItemSizeBytes {
		return nil, ErrInvalidParameter
	}

	key := s.extractKeyFromItem(table, item)
	if key == nil {
		return nil, ErrMissingKey
	}

	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	exprAttrNames, exprAttrValues := getExpressionAttributes(req.Parameters)
	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")

	var oldItem *dbstore.Item
	var isNewItem bool
	var storedItem *dbstore.Item

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, err := txn.GetItem(tableName, key)
		if err != nil {
			if dbstore.IsItemNotFound(err) {
				isNewItem = true
				if conditionExpr != "" {
					syntheticItem := &dbstore.Item{
						TableName:  tableName,
						Key:        key,
						Attributes: make(map[string]*dbstore.AttributeValue),
					}
					conditionMet, evalErr := evaluateConditionExpression(syntheticItem, conditionExpr, exprAttrNames, exprAttrValues)
					if evalErr != nil {
						return evalErr
					}
					if !conditionMet {
						return ErrConditionalCheckFailed
					}
				}
			} else {
				return err
			}
		} else {
			oldItem = existingItem
			if oldItem != nil {
				if err := txn.DeleteIndexEntries(tableName, oldItem); err != nil {
					return err
				}
			}
			if conditionExpr != "" {
				conditionMet, err := evaluateConditionExpression(existingItem, conditionExpr, exprAttrNames, exprAttrValues)
				if err != nil {
					return err
				}
				if !conditionMet {
					return ErrConditionalCheckFailed
				}
			}
		}

		if err := txn.PutItem(tableName, key, item); err != nil {
			return err
		}

		storedItem = &dbstore.Item{
			TableName:  tableName,
			Key:        key,
			Attributes: item,
		}
		if err := txn.PutIndexEntries(tableName, storedItem); err != nil {
			return err
		}

		newItemSize := calculateItemSize(item)
		if isNewItem {
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(tableName, newItemSize); err != nil {
				return err
			}
		} else if oldItem != nil {
			oldItemSize := calculateItemSize(oldItem.Attributes)
			if newItemSize != oldItemSize {
				if err := txn.UpdateTableSize(tableName, newItemSize-oldItemSize); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		if err == ErrConditionalCheckFailed {
			return nil, err
		}
		return nil, err
	}

	resp := map[string]interface{}{}
	if returnValues == "ALL_OLD" && oldItem != nil {
		resp["Attributes"] = buildItemResponse(oldItem.Attributes)
	} else if returnValues == "ALL_NEW" && storedItem != nil {
		resp["Attributes"] = buildItemResponse(storedItem.Attributes)
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		resp["ConsumedCapacity"] = buildConsumedCapacityResponse(tableName, 1.0)
	}

	return resp, nil
}

// GetItem retrieves an item from the specified table using the provided key.
func (s *DynamoDBService) GetItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	key := parseKey(req.Parameters["Key"])
	if key == nil {
		return nil, ErrMissingKey
	}

	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	item, err := store.Items().Get(tableName, key)
	if err != nil {
		if dbstore.IsItemNotFound(err) {
			return response.EmptyResponse(), nil
		}
		return nil, err
	}

	projection := parseProjectionExpression(req.Parameters)
	if projection != nil {
		item.Attributes = applyProjection(item.Attributes, projection)
	}

	resp := map[string]interface{}{
		"Item": buildItemResponse(item.Attributes),
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		resp["ConsumedCapacity"] = buildConsumedCapacityResponse(tableName, 0.5)
	}

	return resp, nil
}

// DeleteItem removes an item from the specified table using the provided key.
func (s *DynamoDBService) DeleteItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	key := parseKey(req.Parameters["Key"])
	if key == nil {
		return nil, ErrMissingKey
	}

	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	exprAttrNames, exprAttrValues := getExpressionAttributes(req.Parameters)
	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")

	var oldItem *dbstore.Item

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, err := txn.GetItem(tableName, key)
		if err != nil {
			if dbstore.IsItemNotFound(err) {
				if conditionExpr != "" {
					syntheticItem := &dbstore.Item{
						TableName:  tableName,
						Key:        key,
						Attributes: make(map[string]*dbstore.AttributeValue),
					}
					conditionMet, evalErr := evaluateConditionExpression(syntheticItem, conditionExpr, exprAttrNames, exprAttrValues)
					if evalErr != nil {
						return evalErr
					}
					if !conditionMet {
						return ErrConditionalCheckFailed
					}
				}
				return nil
			}
			return err
		}

		oldItem = existingItem

		if conditionExpr != "" {
			conditionMet, err := evaluateConditionExpression(existingItem, conditionExpr, exprAttrNames, exprAttrValues)
			if err != nil {
				return err
			}
			if !conditionMet {
				return ErrConditionalCheckFailed
			}
		}

		if oldItem != nil {
			if err := txn.DeleteIndexEntries(tableName, oldItem); err != nil {
				return err
			}
		}

		if err := txn.DeleteItem(tableName, key); err != nil {
			return err
		}

		if oldItem != nil {
			if err := txn.UpdateItemCount(tableName, -1); err != nil {
				return err
			}
			oldItemSize := calculateItemSize(oldItem.Attributes)
			if err := txn.UpdateTableSize(tableName, -oldItemSize); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if err == ErrConditionalCheckFailed {
			return nil, err
		}
		return nil, err
	}

	resp := map[string]interface{}{}
	if returnValues == "ALL_OLD" && oldItem != nil {
		resp["Attributes"] = buildItemResponse(oldItem.Attributes)
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		resp["ConsumedCapacity"] = buildConsumedCapacityResponse(tableName, 1.0)
	}

	return resp, nil
}

// UpdateItem edits an existing item's attributes or creates a new item if it does not exist.
func (s *DynamoDBService) UpdateItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	key := parseKey(req.Parameters["Key"])
	if key == nil {
		return nil, ErrMissingKey
	}

	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")
	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	exprAttrNames, exprAttrValues := getExpressionAttributes(req.Parameters)

	updateExpr := request.GetStringParam(req.Parameters, "UpdateExpression")
	attrs := req.Parameters["AttributeUpdates"]

	if updateExpr != "" && attrs != nil {
		return nil, ErrInvalidParameter
	}

	var oldItem *dbstore.Item
	var storedItem *dbstore.Item
	var updatedAttrNames []string
	var oldItemSize int64

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, err := txn.GetItem(tableName, key)
		isNewItem := false
		var item *dbstore.Item

		if err != nil {
			if dbstore.IsItemNotFound(err) {
				isNewItem = true
				if conditionExpr != "" {
					syntheticItem := &dbstore.Item{
						TableName:  tableName,
						Key:        key,
						Attributes: make(map[string]*dbstore.AttributeValue),
					}
					conditionMet, evalErr := evaluateConditionExpression(syntheticItem, conditionExpr, exprAttrNames, exprAttrValues)
					if evalErr != nil {
						return evalErr
					}
					if !conditionMet {
						return ErrConditionalCheckFailed
					}
				}
				item = &dbstore.Item{
					TableName:  tableName,
					Key:        key,
					Attributes: make(map[string]*dbstore.AttributeValue),
				}
				for k, v := range key {
					item.Attributes[k] = v
				}
			} else {
				return err
			}
		} else {
			item = existingItem
			oldItemSize = calculateItemSize(item.Attributes)
			if conditionExpr != "" {
				conditionMet, err := evaluateConditionExpression(item, conditionExpr, exprAttrNames, exprAttrValues)
				if err != nil {
					return err
				}
				if !conditionMet {
					return ErrConditionalCheckFailed
				}
			}

			if returnValues == "ALL_OLD" || returnValues == "UPDATED_OLD" {
				oldItem = &dbstore.Item{
					Attributes: make(map[string]*dbstore.AttributeValue),
				}
				for k, v := range item.Attributes {
					oldItem.Attributes[k] = deepCopyAttributeValue(v)
				}
			}

			if err := txn.DeleteIndexEntries(tableName, item); err != nil {
				return err
			}
		}

		if updateExpr != "" {
			var err error
			updatedAttrNames, err = applyUpdateExpressionWithTracking(item.Attributes, updateExpr, parseExpressionAttributeNames(req.Parameters), parseExpressionAttributeValues(req.Parameters))
			if err != nil {
				if err.Error() == "TYPE_MISMATCH: Type mismatch for attribute to update" {
					return ErrInvalidParameter
				}
				return err
			}
		} else if attrs != nil {
			updatedAttrNames, err = applyAttributeUpdatesWithTracking(item.Attributes, attrs)
			if err != nil {
				return err
			}
		}

		if itemSize := calculateItemSize(item.Attributes); itemSize > maxItemSizeBytes {
			return ErrInvalidParameter
		}

		if err := txn.PutItem(tableName, key, item.Attributes); err != nil {
			return err
		}

		storedItem = &dbstore.Item{
			TableName:  tableName,
			Key:        key,
			Attributes: item.Attributes,
		}
		if err := txn.PutIndexEntries(tableName, storedItem); err != nil {
			return err
		}

		newItemSize := calculateItemSize(item.Attributes)
		if isNewItem {
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(tableName, newItemSize); err != nil {
				return err
			}
		} else {
			if newItemSize != oldItemSize {
				if err := txn.UpdateTableSize(tableName, newItemSize-oldItemSize); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		if err == ErrConditionalCheckFailed {
			return nil, err
		}
		return nil, err
	}

	resp := map[string]interface{}{}
	switch returnValues {
	case "ALL_OLD":
		if oldItem != nil {
			resp["Attributes"] = buildItemResponse(oldItem.Attributes)
		}
	case "ALL_NEW":
		if storedItem != nil {
			resp["Attributes"] = buildItemResponse(storedItem.Attributes)
		}
	case "UPDATED_OLD":
		if oldItem != nil && len(updatedAttrNames) > 0 {
			resp["Attributes"] = buildUpdatedAttributesResponse(oldItem.Attributes, updatedAttrNames)
		}
	case "UPDATED_NEW":
		if storedItem != nil && len(updatedAttrNames) > 0 {
			resp["Attributes"] = buildUpdatedAttributesResponse(storedItem.Attributes, updatedAttrNames)
		}
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		resp["ConsumedCapacity"] = buildConsumedCapacityResponse(tableName, 1.0)
	}

	return resp, nil
}
