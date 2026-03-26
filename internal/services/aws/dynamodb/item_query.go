package dynamodb

import (
	"context"

	"vorpalstacks/internal/services/aws/common/pagination"
	"vorpalstacks/internal/services/aws/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// Query retrieves items based on their key condition expression.
func (s *DynamoDBService) Query(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	indexName := request.GetStringParam(req.Parameters, "IndexName")
	if indexName != "" {
		if !validateIndexExists(table, indexName) {
			return nil, ErrIndexNotFound
		}
	}

	_ = request.GetBoolParam(req.Parameters, "ConsistentRead")
	limit := pagination.GetMaxItems(req.Parameters, 100)
	exclusiveStartKey := parseExclusiveStartKey(req.Parameters)
	scanIndexForward := getBoolParamWithDefault(req.Parameters, "ScanIndexForward", true)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var allItems []*dbstore.Item
	keyCondExpr := request.GetStringParam(req.Parameters, "KeyConditionExpression")
	exprAttrNames := parseExpressionAttributeNames(req.Parameters)
	exprAttrValues := parseExpressionAttributeValues(req.Parameters)

	if keyCondExpr == "" {
		return nil, ErrInvalidParameter
	}

	// alreadySortedAscending is set when items are returned in ascending sort-key
	// order directly from storage (i.e. via ScanByPartitionKeyWithTable), so that
	// ScanIndexForward=false can be satisfied with a simple O(n) reverse instead of
	// an O(n log n) sort.
	alreadySortedAscending := false

	if indexName != "" {
		hashKeyValue, sortKeyCondition := extractIndexKeyCondition(table, indexName, keyCondExpr, exprAttrNames, exprAttrValues)
		if hashKeyValue != "" {
			err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
				opts := dbstore.IndexQueryOptions{
					Reverse: !scanIndexForward,
				}
				if sortKeyCondition == nil {
					opts.Limit = limit
				}
				var items []*dbstore.Item
				var queryErr error
				if isGSI(table, indexName) {
					items, queryErr = txn.QueryByGSI(tableName, indexName, hashKeyValue, opts)
				} else {
					items, queryErr = txn.QueryByLSI(tableName, indexName, hashKeyValue, opts)
				}
				if queryErr != nil {
					return queryErr
				}
				allItems = items
				return nil
			})
			if err != nil {
				return nil, err
			}
			if sortKeyCondition != nil {
				allItems = filterBySortKeyCondition(allItems, sortKeyCondition)
			}
		} else {
			err = store.Items().Scan(tableName, func(item *dbstore.Item) error {
				allItems = append(allItems, item)
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	} else {
		hashKeyValue, sortKeyCondition := extractPrimaryKeyCondition(table, keyCondExpr, exprAttrNames, exprAttrValues)
		if hashKeyValue != "" {
			err = store.Items().ScanByPartitionKeyWithTable(tableName, table, hashKeyValue, func(item *dbstore.Item) error {
				allItems = append(allItems, item)
				return nil
			})
			if err != nil {
				return nil, err
			}
			if sortKeyCondition != nil {
				allItems = filterBySortKeyCondition(allItems, sortKeyCondition)
			}
			alreadySortedAscending = true
		} else {
			err = store.Items().Scan(tableName, func(item *dbstore.Item) error {
				allItems = append(allItems, item)
				return nil
			})
			if err != nil {
				return nil, err
			}
			allItems = filterByKeyCondition(allItems, keyCondExpr, exprAttrNames, exprAttrValues)
		}
	}

	if !scanIndexForward {
		skName := getSortKeyName(table, indexName)
		if alreadySortedAscending && skName != "" && getSortKeyType(table, skName) == "S" {
			for i, j := 0, len(allItems)-1; i < j; i, j = i+1, j-1 {
				allItems[i], allItems[j] = allItems[j], allItems[i]
			}
		} else {
			sortItemsReverseBySortKeyWithIndex(table, allItems, indexName)
		}
	} else {
		sortItemsBySortKeyWithIndex(table, allItems, indexName)
	}

	if exclusiveStartKey != nil {
		allItems = skipToKeyMap(allItems, exclusiveStartKey, table, indexName)
	}

	var scannedItems []*dbstore.Item
	if len(allItems) > limit {
		scannedItems = allItems[:limit]
	} else {
		scannedItems = allItems
	}

	scannedCount := len(scannedItems)

	filterExpr := request.GetStringParam(req.Parameters, "FilterExpression")
	var items []*dbstore.Item
	if filterExpr != "" {
		items = filterByExpression(scannedItems, filterExpr, exprAttrNames, exprAttrValues)
	} else {
		items = scannedItems
	}

	hasMoreItems := len(allItems) > limit

	projection := parseProjectionExpression(req.Parameters)
	if projection != nil {
		for _, item := range items {
			item.Attributes = applyProjection(item.Attributes, projection)
		}
	}

	collectionKey := tableName
	if indexName != "" {
		collectionKey = indexName
	}
	sizeGB := float64(table.TableSizeBytes) / (1024.0 * 1024.0 * 1024.0)
	resp := map[string]interface{}{
		"Items":        buildItemsResponse(items),
		"Count":        len(items),
		"ScannedCount": scannedCount,
		"ItemCollectionMetrics": map[string]interface{}{
			collectionKey: map[string]interface{}{
				"SizeEstimateRangeGB": []float64{sizeGB, sizeGB},
			},
		},
	}
	if indexName != "" {
		resp["IndexName"] = indexName
	}
	if hasMoreItems && len(scannedItems) > 0 {
		resp["LastEvaluatedKey"] = buildLastEvaluatedKeyWithIndex(scannedItems[len(scannedItems)-1], table, indexName)
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		capacityUnits := float64(scannedCount) * 0.5
		isLSI := indexName != "" && !isGSI(table, indexName)
		resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithIndex(tableName, indexName, capacityUnits, isLSI)
	}

	return resp, nil
}

// Scan retrieves all items in the specified table or index.
func (s *DynamoDBService) Scan(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	indexName := request.GetStringParam(req.Parameters, "IndexName")
	if indexName != "" {
		if !validateIndexExists(table, indexName) {
			return nil, ErrIndexNotFound
		}
	}

	_ = request.GetBoolParam(req.Parameters, "ConsistentRead")
	limit := pagination.GetMaxItems(req.Parameters, 100)
	exclusiveStartKey := parseExclusiveStartKey(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	var allItems []*dbstore.Item
	err = store.Items().Scan(tableName, func(item *dbstore.Item) error {
		allItems = append(allItems, item)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if indexName != "" {
		allItems = filterItemsByIndexMembership(allItems, table, indexName)
	}

	if exclusiveStartKey != nil {
		allItems = skipToKeyMap(allItems, exclusiveStartKey, table, indexName)
	}

	var scannedItems []*dbstore.Item
	if len(allItems) > limit {
		scannedItems = allItems[:limit]
	} else {
		scannedItems = allItems
	}

	scannedCount := len(scannedItems)

	filterExpr := request.GetStringParam(req.Parameters, "FilterExpression")
	var items []*dbstore.Item
	if filterExpr != "" {
		items = filterByExpression(scannedItems, filterExpr, parseExpressionAttributeNames(req.Parameters), parseExpressionAttributeValues(req.Parameters))
	} else {
		items = scannedItems
	}

	hasMoreItems := len(allItems) > limit

	projection := parseProjectionExpression(req.Parameters)
	if projection != nil {
		for _, item := range items {
			item.Attributes = applyProjection(item.Attributes, projection)
		}
	}

	collectionKey := tableName
	if indexName != "" {
		collectionKey = indexName
	}
	sizeGB := float64(table.TableSizeBytes) / (1024.0 * 1024.0 * 1024.0)
	resp := map[string]interface{}{
		"Items":        buildItemsResponse(items),
		"Count":        len(items),
		"ScannedCount": scannedCount,
		"ItemCollectionMetrics": map[string]interface{}{
			collectionKey: map[string]interface{}{
				"SizeEstimateRangeGB": []float64{sizeGB, sizeGB},
			},
		},
	}
	if indexName != "" {
		resp["IndexName"] = indexName
	}
	if hasMoreItems && len(scannedItems) > 0 {
		resp["LastEvaluatedKey"] = buildLastEvaluatedKeyWithIndex(scannedItems[len(scannedItems)-1], table, indexName)
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		capacityUnits := float64(scannedCount) * 0.5
		isLSI := indexName != "" && !isGSI(table, indexName)
		resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithIndex(tableName, indexName, capacityUnits, isLSI)
	}

	return resp, nil
}

func filterItemsByIndexMembership(items []*dbstore.Item, table *dbstore.Table, indexName string) []*dbstore.Item {
	var hashKeyName, sortKeyName string
	var isGSI bool
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			isGSI = true
			for _, key := range gsi.KeySchema {
				if key.KeyType == "HASH" {
					hashKeyName = key.AttributeName
				} else if key.KeyType == "RANGE" {
					sortKeyName = key.AttributeName
				}
			}
			break
		}
	}
	if !isGSI {
		for _, lsi := range table.LocalSecondaryIndexes {
			if lsi.IndexName == indexName {
				for _, key := range lsi.KeySchema {
					if key.KeyType == "HASH" {
						hashKeyName = key.AttributeName
					} else if key.KeyType == "RANGE" {
						sortKeyName = key.AttributeName
					}
				}
				break
			}
		}
	}
	if hashKeyName == "" && sortKeyName == "" {
		return items
	}
	var result []*dbstore.Item
	for _, item := range items {
		if isGSI {
			_, hasHash := item.Attributes[hashKeyName]
			_, hasSort := item.Attributes[sortKeyName]
			if hasHash && (sortKeyName == "" || hasSort) {
				result = append(result, item)
			}
		} else {
			if sortKeyName != "" {
				if _, exists := item.Attributes[sortKeyName]; exists {
					result = append(result, item)
				}
			} else if hashKeyName != "" {
				if _, exists := item.Attributes[hashKeyName]; exists {
					result = append(result, item)
				}
			}
		}
	}
	return result
}
