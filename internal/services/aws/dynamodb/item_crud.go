package dynamodb

import (
	"context"
	"net/http"

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

	// Tolerant extraction: a missing Item member parses to nil and
	// putItemCore rejects it (the derived key cannot be empty); a present
	// but malformed one reports its own cause.
	item, itemErr := parseItem(req.Parameters["Item"])
	if itemErr != nil {
		return nil, NewAPIError("com.amazon.coral.validate#ValidationException", itemErr.Error(), http.StatusBadRequest)
	}

	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	exprAttrNames, exprAttrValues, err := getExpressionAttributes(req.Parameters)
	if err != nil {
		return nil, err
	}
	condition := conditionSpec{Expr: conditionExpr, Names: exprAttrNames, Values: exprAttrValues}
	if legacyCond, lexErr := resolveLegacyWriteCondition(req.Parameters, conditionExpr, exprAttrNames, exprAttrValues); lexErr != nil {
		return nil, lexErr
	} else if legacyCond != nil {
		condition = *legacyCond
	}
	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")
	returnValuesOnConditionCheckFailure := request.GetStringParam(req.Parameters, "ReturnValuesOnConditionCheckFailure")
	// The response-shaping enums travel in the input; the core's validation
	// is the single enum check and rejects an unknown value before anything
	// executes. The handler reads the same values to shape its response.
	returnConsumedCapacity := request.GetStringParam(req.Parameters, "ReturnConsumedCapacity")
	collectionMetrics := request.GetStringParam(req.Parameters, "ReturnItemCollectionMetrics")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.putItemCore(ctx, store, reqCtx.GetRegion(), PutItemCoreInput{
		Table:                               table,
		Item:                                item,
		Condition:                           condition,
		ReturnValues:                        returnValues,
		ReturnValuesOnConditionCheckFailure: returnValuesOnConditionCheckFailure,
		ReturnConsumedCapacity:              returnConsumedCapacity,
		ReturnItemCollectionMetrics:         collectionMetrics,
	})
	if err != nil {
		return nil, err
	}
	key := result.StoredItem.Key

	resp := map[string]interface{}{}
	if returnValues == "ALL_OLD" && result.OldItem != nil {
		resp["Attributes"] = buildItemResponse(result.OldItem.Attributes)
	}

	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithVector(table.Name, result.WriteUnits, vectorWriteCapacityForItems(table, result.StoredItem))
	}

	if collectionMetrics == "SIZE" {
		if entry := buildItemCollectionMetricsEntry(table, key); entry != nil {
			resp["ItemCollectionMetrics"] = entry
		}
	}

	return resp, nil
}

// GetItem retrieves an item from the specified table using the provided key.
func (s *DynamoDBService) GetItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	// Tolerant extraction: a missing Key member parses to nil and
	// getItemCore rejects it via the empty-key rule; a present but
	// malformed one reports its own cause.
	key, keyErr := parseKey(req.Parameters["Key"])
	if keyErr != nil {
		return nil, NewAPIError("com.amazon.coral.validate#ValidationException", keyErr.Error(), http.StatusBadRequest)
	}

	// The consumed-capacity member travels in the core call; the core's
	// validation is the single enum check. The handler reads the same value
	// to shape its response.
	returnConsumedCapacity := request.GetStringParam(req.Parameters, "ReturnConsumedCapacity")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	item, err := s.getItemCore(ctx, store, table, key, req.Parameters["ConsistentRead"], returnConsumedCapacity)
	if err != nil {
		if isItemNotFound(err) {
			return response.EmptyResponse(), nil
		}
		return nil, err
	}

	projection, err := resolveProjectionMembers(req.Parameters)
	if err != nil {
		return nil, err
	}
	if projection != nil {
		item.Item.Attributes = applyProjection(item.Item.Attributes, projection)
	}

	resp := map[string]interface{}{
		"Item": buildItemResponse(item.Item.Attributes),
	}

	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		// The Core reports the charge over the item's full size as read,
		// before any projection narrows the returned attributes. GetItem
		// accesses no indexes, so INDEXES returns the table's detail alone
		// ("some operations, such as GetItem and BatchGetItem, do not
		// access any indexes at all"); TOTAL reports only the aggregate.
		if returnConsumedCapacity == "INDEXES" {
			resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithIndex(table.Name, "", item.ReadUnits, false)
		} else {
			resp["ConsumedCapacity"] = buildConsumedCapacityResponse(table.Name, item.ReadUnits)
		}
	}

	return resp, nil
}

// DeleteItem removes an item from the specified table using the provided key.
func (s *DynamoDBService) DeleteItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// A write: resolve through the ACTIVE-required path like PutItem and
	// UpdateItem — a table mid-restore (CREATING) must not be mutated.
	table, err := s.validateAndGetActiveTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	// Tolerant extraction: a missing Key member parses to nil and
	// deleteItemCore rejects it via the empty-key rule; a present but
	// malformed one reports its own cause.
	key, keyErr := parseKey(req.Parameters["Key"])
	if keyErr != nil {
		return nil, NewAPIError("com.amazon.coral.validate#ValidationException", keyErr.Error(), http.StatusBadRequest)
	}

	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	exprAttrNames, exprAttrValues, err := getExpressionAttributes(req.Parameters)
	if err != nil {
		return nil, err
	}
	condition := conditionSpec{Expr: conditionExpr, Names: exprAttrNames, Values: exprAttrValues}
	if legacyCond, lexErr := resolveLegacyWriteCondition(req.Parameters, conditionExpr, exprAttrNames, exprAttrValues); lexErr != nil {
		return nil, lexErr
	} else if legacyCond != nil {
		condition = *legacyCond
	}
	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")
	returnValuesOnConditionCheckFailure := request.GetStringParam(req.Parameters, "ReturnValuesOnConditionCheckFailure")
	returnConsumedCapacity := request.GetStringParam(req.Parameters, "ReturnConsumedCapacity")
	collectionMetrics := request.GetStringParam(req.Parameters, "ReturnItemCollectionMetrics")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.deleteItemCore(ctx, store, reqCtx.GetRegion(), DeleteItemCoreInput{
		Table:                               table,
		Key:                                 key,
		Condition:                           condition,
		ReturnValues:                        returnValues,
		ReturnValuesOnConditionCheckFailure: returnValuesOnConditionCheckFailure,
		ReturnConsumedCapacity:              returnConsumedCapacity,
		ReturnItemCollectionMetrics:         collectionMetrics,
	})
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{}
	if returnValues == "ALL_OLD" && result.OldItem != nil {
		resp["Attributes"] = buildItemResponse(result.OldItem.Attributes)
	}

	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithVector(table.Name, result.WriteUnits, vectorWriteCapacityForItems(table, result.OldItem))
	}

	if collectionMetrics == "SIZE" {
		if entry := buildItemCollectionMetricsEntry(table, key); entry != nil {
			resp["ItemCollectionMetrics"] = entry
		}
	}

	return resp, nil
}

// UpdateItem edits an existing item's attributes or creates a new item if it does not exist.
func (s *DynamoDBService) UpdateItem(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetActiveTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	// Tolerant extraction: a missing Key member parses to nil and
	// updateItemCore rejects it via the empty-key rule; a present but
	// malformed one reports its own cause.
	key, keyErr := parseKey(req.Parameters["Key"])
	if keyErr != nil {
		return nil, NewAPIError("com.amazon.coral.validate#ValidationException", keyErr.Error(), http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	returnValues := request.GetStringParam(req.Parameters, "ReturnValues")
	returnValuesOnConditionCheckFailure := request.GetStringParam(req.Parameters, "ReturnValuesOnConditionCheckFailure")
	updateExpr := request.GetStringParam(req.Parameters, "UpdateExpression")
	conditionExpr := request.GetStringParam(req.Parameters, "ConditionExpression")
	attrs := req.Parameters["AttributeUpdates"]

	returnConsumedCapacity := request.GetStringParam(req.Parameters, "ReturnConsumedCapacity")
	collectionMetrics := request.GetStringParam(req.Parameters, "ReturnItemCollectionMetrics")

	exprAttrNames, exprAttrValues, err := getExpressionAttributes(req.Parameters)
	if err != nil {
		return nil, err
	}
	condition := conditionSpec{Expr: conditionExpr, Names: exprAttrNames, Values: exprAttrValues}
	if legacyCond, lexErr := resolveLegacyWriteCondition(req.Parameters, conditionExpr, exprAttrNames, exprAttrValues); lexErr != nil {
		return nil, lexErr
	} else if legacyCond != nil {
		condition = *legacyCond
	}

	result, err := s.updateItemCore(ctx, store, reqCtx.GetRegion(), table, UpdateItemInput{
		Key:                                 key,
		UpdateExpr:                          updateExpr,
		AttrUpdates:                         attrs,
		ConditionExpr:                       condition.Expr,
		ExprAttrNames:                       condition.Names,
		ExprAttrValues:                      condition.Values,
		ReturnValues:                        returnValues,
		ReturnValuesOnConditionCheckFailure: returnValuesOnConditionCheckFailure,
		ReturnConsumedCapacity:              returnConsumedCapacity,
		ReturnItemCollectionMetrics:         collectionMetrics,
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

	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithVector(table.Name, result.WriteUnits, vectorWriteCapacityForItems(table, result.StoredItem))
	}

	if collectionMetrics == "SIZE" {
		if entry := buildItemCollectionMetricsEntry(table, key); entry != nil {
			resp["ItemCollectionMetrics"] = entry
		}
	}

	return resp, nil
}
