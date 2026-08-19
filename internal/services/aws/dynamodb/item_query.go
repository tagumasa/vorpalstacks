package dynamodb

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"net/http"

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
	// relax consistency because there is no replica lag. Global secondary
	// indexes are eventually consistent in the AWS contract, so a strongly
	// consistent read against one is rejected.
	consistentRead := request.GetBoolParam(req.Parameters, "ConsistentRead")
	if indexName != "" && isGSI(table, indexName) && consistentRead {
		return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
			"Consistent reads are not supported on global secondary indexes", http.StatusBadRequest)
	}
	// An explicit Limit below 1 is invalid (valid range minimum 1); a plain
	// value check cannot tell "unset" apart from an explicit zero, so
	// presence is checked first.
	if _, ok := req.Parameters["Limit"]; ok && request.GetIntParam(req.Parameters, "Limit") <= 0 {
		return nil, ErrInvalidParameter
	}
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
	if exclusiveStartKey != nil {
		if err := validateKeyTypes(table, exclusiveStartKey); err != nil {
			return nil, err
		}
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

	projection, projErr := parseProjectionExpression(req.Parameters)
	if projErr != nil {
		return nil, projErr
	}
	countOnly, allProjected, selErr := parseSelectParam(req.Parameters, indexName, projection != nil)
	if selErr != nil {
		return nil, selErr
	}
	if indexName != "" && isGSI(table, indexName) {
		// A global secondary index can only serve attributes projected into
		// it; local secondary indexes may fall back to the parent table.
		if gsiProjErr := validateGSIProjectionRequest(table, indexName, allProjected, countOnly, projection); gsiProjErr != nil {
			return nil, gsiProjErr
		}
	}

	// A Query must perform an equality test on the partition key of the
	// table or index being queried; DynamoDB answers any other key
	// condition shape with a ValidationException rather than widening the
	// read to a scan.
	//
	// alreadySortedAscending is set when items are returned in ascending sort-key
	// order directly from storage (i.e. via ScanByPartitionKeyWithTable), so that
	// ScanIndexForward=false can be satisfied with a simple O(n) reverse instead of
	// an O(n log n) sort.
	alreadySortedAscending := false

	if indexName != "" {
		hashKeyValue, sortKeyCondition := extractIndexKeyCondition(table, indexName, keyCondExpr, exprAttrNames, exprAttrValues)
		if hashKeyValue == "" {
			idxHashName, _, _ := indexKeyAttributeNames(table, indexName)
			return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
				"Query condition missed key schema element: "+idxHashName, http.StatusBadRequest)
		}
		err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
			opts := dbstore.IndexQueryOptions{
				Reverse: !scanIndexForward,
			}
			if sortKeyCondition == nil && exclusiveStartKey == nil {
				// Fetch one item beyond the page so the response tail can
				// tell a full page apart from the end of the index range.
				opts.Limit = limit + 1
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
		hashKeyValue, sortKeyCondition := extractPrimaryKeyCondition(table, keyCondExpr, exprAttrNames, exprAttrValues)
		if hashKeyValue == "" {
			pkAttrName := ""
			for _, ks := range table.KeySchema {
				if ks.KeyType == dbstore.KeyTypeHash {
					pkAttrName = ks.AttributeName
					break
				}
			}
			return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
				"Query condition missed key schema element: "+pkAttrName, http.StatusBadRequest)
		}
		pkOpts := dbstore.ScanOptions{}
		// The lookahead cap bounds the ascending storage scan, so it can
		// only serve forward pages; a reverse page reads its items from the
		// far end of the partition and must not cap the scan.
		if sortKeyCondition == nil && exclusiveStartKey == nil && scanIndexForward {
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
		allItems = skipToKeyMap(allItems, exclusiveStartKey, table, indexName, scanIndexForward)
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

	if projection != nil {
		for _, item := range items {
			item.Attributes = applyProjection(item.Attributes, projection)
		}
	} else if allProjected {
		for _, item := range items {
			item.Attributes = applyIndexProjection(item.Attributes, table, indexName)
		}
	}

	collectionKey := tableName
	if indexName != "" {
		collectionKey = indexName
	}
	sizeGB := float64(table.TableSizeBytes) / (1024.0 * 1024.0 * 1024.0)
	resp := map[string]interface{}{
		"Count":        len(items),
		"ScannedCount": scannedCount,
		"ItemCollectionMetrics": map[string]interface{}{
			collectionKey: map[string]interface{}{
				"SizeEstimateRangeGB": []float64{sizeGB, sizeGB},
			},
		},
	}
	if !countOnly {
		resp["Items"] = buildItemsResponse(items)
	}
	if hasMoreItems && len(scannedItems) > 0 {
		resp["LastEvaluatedKey"] = buildLastEvaluatedKeyWithIndex(scannedItems[len(scannedItems)-1], table, indexName)
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		capacityUnits := float64(scannedCount) * rcuPerItem(consistentRead, indexName, table)
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
	// relax consistency because there is no replica lag. Global secondary
	// indexes are eventually consistent in the AWS contract, so a strongly
	// consistent read against one is rejected.
	consistentRead := request.GetBoolParam(req.Parameters, "ConsistentRead")
	if indexName != "" && isGSI(table, indexName) && consistentRead {
		return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
			"Consistent reads are not supported on global secondary indexes", http.StatusBadRequest)
	}
	// An explicit Limit below 1 is invalid (valid range minimum 1); a plain
	// value check cannot tell "unset" apart from an explicit zero, so
	// presence is checked first.
	if _, ok := req.Parameters["Limit"]; ok && request.GetIntParam(req.Parameters, "Limit") <= 0 {
		return nil, ErrInvalidParameter
	}
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
	if exclusiveStartKey != nil {
		if err := validateKeyTypes(table, exclusiveStartKey); err != nil {
			return nil, err
		}
	}

	// Validate parallel Scan parameters (Smithy ScanSegment: range 0-999999,
	// ScanTotalSegments: range 1-1000000). The two parameters form a pair:
	// specifying one without the other is a validation error.
	_, hasSegment := req.Parameters["Segment"]
	_, hasTotalSegments := req.Parameters["TotalSegments"]
	if hasSegment != hasTotalSegments {
		return nil, ErrInvalidParameter
	}
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	projection, projErr := parseProjectionExpression(req.Parameters)
	if projErr != nil {
		return nil, projErr
	}
	countOnly, allProjected, selErr := parseSelectParam(req.Parameters, indexName, projection != nil)
	if selErr != nil {
		return nil, selErr
	}
	if indexName != "" && isGSI(table, indexName) {
		// A global secondary index can only serve attributes projected into
		// it; local secondary indexes may fall back to the parent table.
		if gsiProjErr := validateGSIProjectionRequest(table, indexName, allProjected, countOnly, projection); gsiProjErr != nil {
			return nil, gsiProjErr
		}
	}

	// The page is built during iteration — segment filtering, secondary
	// index membership, and the exclusive-start-key skip all apply before
	// the limit is counted, so every qualifying item remains reachable
	// through pagination regardless of how the surrounding table entries
	// are interleaved.
	page, pageErr := s.collectScanPage(store, tableName, scanPageOptions{
		table:             table,
		indexName:         indexName,
		limit:             limit,
		segment:           segment,
		totalSegments:     totalSegments,
		exclusiveStartKey: exclusiveStartKey,
	})
	if pageErr != nil {
		return nil, pageErr
	}
	allItems := page.items
	hasMoreItems := page.hasMore

	scannedItems := allItems
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

	if projection != nil {
		for _, item := range items {
			item.Attributes = applyProjection(item.Attributes, projection)
		}
	} else if allProjected {
		for _, item := range items {
			item.Attributes = applyIndexProjection(item.Attributes, table, indexName)
		}
	}

	collectionKey := tableName
	if indexName != "" {
		collectionKey = indexName
	}
	sizeGB := float64(table.TableSizeBytes) / (1024.0 * 1024.0 * 1024.0)
	resp := map[string]interface{}{
		"Count":        len(items),
		"ScannedCount": scannedCount,
		"ItemCollectionMetrics": map[string]interface{}{
			collectionKey: map[string]interface{}{
				"SizeEstimateRangeGB": []float64{sizeGB, sizeGB},
			},
		},
	}
	if !countOnly {
		resp["Items"] = buildItemsResponse(items)
	}
	if hasMoreItems && len(scannedItems) > 0 {
		resp["LastEvaluatedKey"] = buildLastEvaluatedKeyWithIndex(scannedItems[len(scannedItems)-1], table, indexName)
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		capacityUnits := float64(scannedCount) * rcuPerItem(consistentRead, indexName, table)
		isLSI := indexName != "" && !isGSI(table, indexName)
		resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithIndex(tableName, indexName, capacityUnits, isLSI)
	}

	return resp, nil
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
