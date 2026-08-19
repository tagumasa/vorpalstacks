package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// oldItemAttributes safely extracts the attributes map from an item that
// may be nil. Used when passing old image data to stream capture.
func oldItemAttributes(item *dbstore.Item) map[string]*dbstore.AttributeValue {
	if item == nil {
		return nil
	}
	return item.Attributes
}

// PutItem creates a new item or replaces an existing item in the specified table.
func (s *DynamoDBService) PutItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetActiveTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	item, itemErr := parseItem(req.Parameters["Item"])
	if itemErr != nil || item == nil {
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

	if err := validateItemKeyTypes(table, item); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	exprAttrNames, exprAttrValues, err := getExpressionAttributes(req.Parameters)
	if err != nil {
		return nil, err
	}
	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")
	// PutItem and DeleteItem recognise only NONE and ALL_OLD (model
	// ReturnValue enum); any other value is rejected rather than ignored.
	if returnValues != "" && returnValues != "NONE" && returnValues != "ALL_OLD" {
		return nil, ErrInvalidParameter
	}

	// Build condition checker callback from ConditionExpression.
	var condChecker ConditionChecker
	if conditionExpr != "" {
		condChecker = func(existing *dbstore.Item, isNotFound bool) error {
			var evalItem *dbstore.Item
			if isNotFound {
				evalItem = &dbstore.Item{
					TableName:  tableName,
					Key:        key,
					Attributes: make(map[string]*dbstore.AttributeValue),
				}
			} else {
				evalItem = existing
			}
			met, evalErr := evaluateConditionExpression(evalItem, conditionExpr, exprAttrNames, exprAttrValues)
			if evalErr != nil {
				return evalErr
			}
			if !met {
				return ErrConditionalCheckFailed
			}
			return nil
		}
	}

	_, oldItem, err := s.putItemCore(ctx, store, reqCtx.GetRegion(), table, key, item, condChecker)
	if err != nil {
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

// GetItem retrieves an item from the specified table using the provided key.
func (s *DynamoDBService) GetItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	key, keyErr := parseKey(req.Parameters["Key"])
	if keyErr != nil || key == nil {
		return nil, ErrInvalidParameter
	}

	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	if err := validateKeyTypes(table, key); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	item, err := s.getItemCore(ctx, store, tableName, key)
	if err != nil {
		if dbstore.IsItemNotFound(err) {
			return response.EmptyResponse(), nil
		}
		return nil, err
	}

	projection, err := parseProjectionExpression(req.Parameters)
	if err != nil {
		return nil, err
	}
	if projection != nil {
		item.Attributes = applyProjection(item.Attributes, projection)
	}

	resp := map[string]interface{}{
		"Item": buildItemResponse(item.Attributes),
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		// Strongly consistent reads consume twice the capacity of eventually
		// consistent ones; the underlying store is always strongly
		// consistent, so the flag only affects the reported charge.
		capacityUnits := rcuPerItem(request.GetBoolParam(req.Parameters, "ConsistentRead"), "", table)
		resp["ConsumedCapacity"] = buildConsumedCapacityResponse(tableName, capacityUnits)
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

	key, keyErr := parseKey(req.Parameters["Key"])
	if keyErr != nil || key == nil {
		return nil, ErrInvalidParameter
	}

	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	if err := validateKeyTypes(table, key); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	exprAttrNames, exprAttrValues, err := getExpressionAttributes(req.Parameters)
	if err != nil {
		return nil, err
	}
	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")
	// PutItem and DeleteItem recognise only NONE and ALL_OLD (model
	// ReturnValue enum); any other value is rejected rather than ignored.
	if returnValues != "" && returnValues != "NONE" && returnValues != "ALL_OLD" {
		return nil, ErrInvalidParameter
	}

	// Build condition checker callback from ConditionExpression.
	var condChecker ConditionChecker
	if conditionExpr != "" {
		condChecker = func(existing *dbstore.Item, isNotFound bool) error {
			var evalItem *dbstore.Item
			if isNotFound {
				evalItem = &dbstore.Item{
					TableName:  tableName,
					Key:        key,
					Attributes: make(map[string]*dbstore.AttributeValue),
				}
			} else {
				evalItem = existing
			}
			met, evalErr := evaluateConditionExpression(evalItem, conditionExpr, exprAttrNames, exprAttrValues)
			if evalErr != nil {
				return evalErr
			}
			if !met {
				return ErrConditionalCheckFailed
			}
			return nil
		}
	}

	oldItem, err := s.deleteItemCore(ctx, store, reqCtx.GetRegion(), table, key, condChecker)
	if err != nil {
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
	table, err := s.validateAndGetActiveTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	key, keyErr := parseKey(req.Parameters["Key"])
	if keyErr != nil || key == nil {
		return nil, ErrInvalidParameter
	}
	if !validateKeyValueNotEmpty(key) {
		return nil, ErrInvalidParameter
	}

	if err := validateKeyTypes(table, key); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")
	updateExpr := request.GetStringParam(req.Parameters, "UpdateExpression")
	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	attrs := req.Parameters["AttributeUpdates"]

	if updateExpr != "" && attrs != nil {
		return nil, ErrInvalidParameter
	}

	exprAttrNames, exprAttrValues, err := getExpressionAttributes(req.Parameters)
	if err != nil {
		return nil, err
	}

	result, err := s.updateItemCore(ctx, store, reqCtx.GetRegion(), table, UpdateItemInput{
		Key:            key,
		UpdateExpr:     updateExpr,
		AttrUpdates:    attrs,
		ConditionExpr:  conditionExpr,
		ExprAttrNames:  exprAttrNames,
		ExprAttrValues: exprAttrValues,
		ReturnValues:   returnValues,
	})
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{}
	switch returnValues {
	case "ALL_OLD":
		if result.OldItem != nil {
			resp["Attributes"] = buildItemResponse(result.OldItem.Attributes)
		}
	case "ALL_NEW":
		if result.StoredItem != nil {
			resp["Attributes"] = buildItemResponse(result.StoredItem.Attributes)
		}
	case "UPDATED_OLD":
		if result.OldItem != nil && len(result.UpdatedAttrNames) > 0 {
			resp["Attributes"] = buildUpdatedAttributesResponse(result.OldItem.Attributes, result.UpdatedAttrNames)
		}
	case "UPDATED_NEW":
		if result.StoredItem != nil && len(result.UpdatedAttrNames) > 0 {
			resp["Attributes"] = buildUpdatedAttributesResponse(result.StoredItem.Attributes, result.UpdatedAttrNames)
		}
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		resp["ConsumedCapacity"] = buildConsumedCapacityResponse(table.Name, 1.0)
	}

	return resp, nil
}
