package dynamodb

import (
	"context"
	"crypto/md5"
	"encoding/binary"

	"vorpalstacks/internal/common/request"
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
		if !validateIndexName(indexName) {
			return nil, ErrInvalidParameter
		}
		if !validateIndexExists(table, indexName) {
			return nil, ErrIndexNotFound
		}
	}

	// ConsistentRead is accepted for API compatibility. Single-instance
	// Pebble provides strong consistency for all reads; the flag cannot
	// relax consistency because there is no replica lag.
	_ = request.GetBoolParam(req.Parameters, "ConsistentRead")
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit <= 0 {
		limit = dataPlaneQueryDefaultLimit
	}
	if limit > dataPlaneQueryMaxLimit {
		limit = dataPlaneQueryMaxLimit
	}
	exclusiveStartKey, eskErr := parseExclusiveStartKey(req.Parameters)
	if eskErr != nil {
		return nil, eskErr
	}
	scanIndexForward, err := validateBoolParam(req.Parameters, "ScanIndexForward", true)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var allItems []*dbstore.Item
	keyCondExpr := request.GetStringParam(req.Parameters, "KeyConditionExpression")
	exprAttrNames, err := parseExpressionAttributeNames(req.Parameters)
	if err != nil {
		return nil, err
	}
	exprAttrValues, eavErr := parseExpressionAttributeValues(req.Parameters)
	if eavErr != nil {
		return nil, eavErr
	}

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
			_, err = store.Items().ScanWithOptions(tableName, dbstore.ScanOptions{Limit: limit + 1}, func(item *dbstore.Item) error {
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
			pkOpts := dbstore.ScanOptions{}
			if sortKeyCondition == nil && exclusiveStartKey == nil {
				pkOpts.Limit = limit + 1
			}
			_, err = store.Items().ScanByPartitionKeyWithTable(tableName, table, hashKeyValue, pkOpts, func(item *dbstore.Item) error {
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
			scanOpts := dbstore.ScanOptions{}
			if exclusiveStartKey == nil {
				scanOpts.Limit = limit
			}
			_, err = store.Items().ScanWithOptions(tableName, scanOpts, func(item *dbstore.Item) error {
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
		if alreadySortedAscending && skName != "" {
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

	projection, err := parseProjectionExpression(req.Parameters)
	if err != nil {
		return nil, err
	}
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
		if !validateIndexName(indexName) {
			return nil, ErrInvalidParameter
		}
		if !validateIndexExists(table, indexName) {
			return nil, ErrIndexNotFound
		}
	}

	// ConsistentRead is accepted for API compatibility. Single-instance
	// Pebble provides strong consistency for all reads; the flag cannot
	// relax consistency because there is no replica lag.
	_ = request.GetBoolParam(req.Parameters, "ConsistentRead")
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit <= 0 {
		limit = dataPlaneQueryDefaultLimit
	}
	if limit > dataPlaneQueryMaxLimit {
		limit = dataPlaneQueryMaxLimit
	}
	exclusiveStartKey, eskErr := parseExclusiveStartKey(req.Parameters)
	if eskErr != nil {
		return nil, eskErr
	}

	// Validate parallel Scan parameters (Smithy ScanSegment: range 0-999999,
	// ScanTotalSegments: range 1-1000000).
	segment := -1
	totalSegments := 0
	if _, ok := req.Parameters["Segment"]; ok {
		segment = request.GetIntParam(req.Parameters, "Segment")
		if !validateScanSegment(segment) {
			return nil, ErrInvalidParameter
		}
	}
	if _, ok := req.Parameters["TotalSegments"]; ok {
		totalSegments = request.GetIntParam(req.Parameters, "TotalSegments")
		if !validateScanTotalSegments(totalSegments) {
			return nil, ErrInvalidParameter
		}
	}
	parallelScan := segment >= 0 && totalSegments > 0
	pkName := ""
	if parallelScan {
		for _, ks := range table.KeySchema {
			if ks.KeyType == dbstore.KeyTypeHash {
				pkName = ks.AttributeName
				break
			}
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	var allItems []*dbstore.Item
	scanLimit := limit + 1
	if exclusiveStartKey != nil || parallelScan {
		scanLimit = 0
	}
	_, err = store.Items().ScanWithOptions(tableName, dbstore.ScanOptions{Limit: scanLimit}, func(item *dbstore.Item) error {
		if parallelScan {
			pkAttr := item.Key[pkName]
			if pkAttr == nil {
				pkAttr = item.Attributes[pkName]
			}
			if pkAttr == nil {
				return nil
			}
			h := md5SegmentHash(pkAttr)
			if int(h%uint32(totalSegments)) != segment {
				return nil
			}
		}
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
		scanNames, namesErr := parseExpressionAttributeNames(req.Parameters)
		if namesErr != nil {
			return nil, namesErr
		}
		scanValues, scanValsErr := parseExpressionAttributeValues(req.Parameters)
		if scanValsErr != nil {
			return nil, scanValsErr
		}
		items = filterByExpression(scannedItems, filterExpr, scanNames, scanValues)
	} else {
		items = scannedItems
	}

	hasMoreItems := len(allItems) > limit

	projection, err := parseProjectionExpression(req.Parameters)
	if err != nil {
		return nil, err
	}
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

// md5SegmentHash computes the MD5 hash of an AttributeValue for parallel
// scan segment assignment, matching DynamoDB's internal partition key
// hash function. The hash is computed over the raw value bytes (S string,
// N number string, or B binary) and the first 4 bytes are interpreted as
// a big-endian uint32.
func md5SegmentHash(av *dbstore.AttributeValue) uint32 {
	h := md5.New()
	if av.S != nil {
		h.Write([]byte(*av.S))
	} else if av.N != nil {
		h.Write([]byte(*av.N))
	} else if av.B != nil {
		h.Write(av.B)
	}
	sum := h.Sum(nil)
	return binary.BigEndian.Uint32(sum[:4])
}
