package dynamodb

import (
	"context"
	"net/http"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Scan page Core — shared paginator for the Scan data plane
// ---------------------------------------------------------------------------

// scanPageOptions configures collectScanPage, the shared paginator for the
// Scan data plane and the index fallback reads of Query.
type scanPageOptions struct {
	table         *dbstore.Table
	indexName     string // secondary index whose membership filters the page; "" scans the base table
	limit         int    // page size; iteration stops after limit+1 qualifying items
	segment       int    // parallel-scan segment; -1 disables segment filtering
	totalSegments int
	// marker is the storage key of the item the previous page stopped at;
	// iteration resumes strictly after it, so a page never re-reads the
	// head of the table. Derived from the request's exclusive start key
	// with the store's key encoder.
	marker string
}

// scanPageResult carries one page of qualifying items plus the continuation
// flag derived from the limit+1 lookahead item.
type scanPageResult struct {
	items   []*dbstore.Item
	hasMore bool
}

// collectScanPage walks the table in storage order and applies, per item, the
// parallel-scan segment filter and the secondary-index membership filter,
// collecting at most limit+1 qualifying items. The lookahead item only sets
// hasMore; callers truncate to limit. Applying the filters during iteration
// — before the limit is counted — is what keeps index scans paginable: a
// page bounded before filtering would strand later index members behind a
// LastEvaluatedKey that is never emitted.
func (s *DynamoDBService) collectScanPage(store dbstore.DynamoDBStoreInterface, tableName string, opts scanPageOptions) (scanPageResult, error) {
	var hashName, sortName string
	isGSIIndex := false
	if opts.indexName != "" {
		hashName, sortName, isGSIIndex = indexKeyAttributeNames(opts.table, opts.indexName)
	}

	pkName := ""
	if opts.segment >= 0 {
		for _, ks := range opts.table.KeySchema {
			if ks.KeyType == dbstore.KeyTypeHash {
				pkName = ks.AttributeName
				break
			}
		}
	}

	result := scanPageResult{}

	_, err := store.Items().ScanWithOptions(tableName, dbstore.ScanOptions{Marker: opts.marker}, func(item *dbstore.Item) error {
		if opts.segment >= 0 && opts.totalSegments > 0 {
			pkAttr := item.Key[pkName]
			if pkAttr == nil {
				pkAttr = item.Attributes[pkName]
			}
			if pkAttr == nil {
				return nil
			}
			if int(md5SegmentHash(pkAttr)%uint32(opts.totalSegments)) != opts.segment {
				return nil
			}
		}

		if opts.indexName != "" && !isIndexMember(item, hashName, sortName, isGSIIndex) {
			return nil
		}

		result.items = append(result.items, item)
		if len(result.items) > opts.limit {
			return errScanSufficient
		}
		return nil
	})
	if err != nil && err != errScanSufficient {
		return scanPageResult{}, err
	}

	result.hasMore = len(result.items) > opts.limit
	if len(result.items) > opts.limit {
		result.items = result.items[:opts.limit]
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Query Core — key-condition read plane
// ---------------------------------------------------------------------------

// queryInput carries the raw wire parameters of a Query request; the Core
// applies every validation in its documented order.
type queryInput struct {
	Parameters map[string]interface{}
}

// queryCore is the single validation and persistence path of the Query data
// plane: table resolution, index and key-condition validation, the snapshot
// read, filtering, ordering, projection, and the response assembly.
func (s *DynamoDBService) queryCore(ctx context.Context, reqCtx *request.RequestContext, in queryInput) (map[string]interface{}, error) {
	params := in.Parameters
	table, err := s.validateAndGetTable(reqCtx, params)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	preamble, preErr := resolveReadIndexPreamble(table, params)
	if preErr != nil {
		return nil, preErr
	}
	indexName := preamble.indexName
	consistentRead := preamble.consistentRead
	limit := preamble.limit
	exclusiveStartKey := preamble.exclusiveStartKey
	scanIndexForward, err := validateBoolParam(params, "ScanIndexForward", true)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var allItems []*dbstore.Item
	// The typed partition-key equality value feeds the contributor insights
	// accounting of the query after the page is served.
	var queryPKValue *dbstore.AttributeValue
	keyCondExpr := request.GetStringParam(params, "KeyConditionExpression")
	exprAttrNames, err := parseExpressionAttributeNames(params)
	if err != nil {
		return nil, err
	}
	exprAttrValues, eavErr := parseExpressionAttributeValues(params)
	if eavErr != nil {
		return nil, eavErr
	}

	if keyCondExpr == "" {
		return nil, ErrInvalidParameter
	}

	projection, countOnly, allProjected, projErr := resolveProjectionSelection(table, indexName, params)
	if projErr != nil {
		return nil, projErr
	}

	// A Query must perform an equality test on the partition key of the
	// table or index being queried; DynamoDB answers any other key
	// condition shape with a ValidationException rather than widening the
	// read to a scan.
	//
	// Both read paths walk storage in sort-key order (the key encoder is
	// sort-correct), so the sort-key condition and the start key are applied
	// during iteration and the limit+1 lookahead bounds the read in both
	// directions — a page never materialises more of the partition than its
	// own size requires.

	if indexName != "" {
		hashKeyValue, hashKeyAttr, sortKeyCondition := extractIndexKeyCondition(table, indexName, keyCondExpr, exprAttrNames, exprAttrValues)
		queryPKValue = hashKeyAttr
		if hashKeyValue == "" {
			idxHashName, _, _ := indexKeyAttributeNames(table, indexName)
			return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
				"Query condition missed key schema element: "+idxHashName, http.StatusBadRequest)
		}
		indexMarker := ""
		if exclusiveStartKey != nil {
			indexMarker = indexMarkerFromStartKey(table, indexName, hashKeyValue, exclusiveStartKey)
			if indexMarker == "" {
				return nil, ErrInvalidParameter
			}
		}
		err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
			opts := dbstore.IndexQueryOptions{
				Reverse: !scanIndexForward,
				Marker:  indexMarker,
				// Fetch one item beyond the page so the response tail can
				// tell a full page apart from the end of the index range.
				Limit: limit + 1,
			}
			if sortKeyCondition != nil {
				cond := sortKeyCondition
				opts.Filter = func(item *dbstore.Item) bool {
					return sortKeyConditionMatches(item, cond)
				}
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
	} else {
		hashKeyValue, hashKeyAttr, sortKeyCondition := extractPrimaryKeyCondition(table, keyCondExpr, exprAttrNames, exprAttrValues)
		queryPKValue = hashKeyAttr
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
		baseMarker := ""
		if exclusiveStartKey != nil {
			startKey := primaryKeyFromStartKey(table, exclusiveStartKey)
			if startKey == nil {
				return nil, ErrInvalidParameter
			}
			baseMarker = dbstore.EncodeItemKey(table.Name, startKey, table)
		}
		pkOpts := dbstore.ScanOptions{
			Reverse: !scanIndexForward,
			Marker:  baseMarker,
			// Fetch one item beyond the page so the response tail can tell
			// a full page apart from the end of the partition.
			Limit: limit + 1,
		}
		if sortKeyCondition != nil {
			cond := sortKeyCondition
			pkOpts.Filter = func(item *dbstore.Item) bool {
				return sortKeyConditionMatches(item, cond)
			}
		}
		_, err = store.Items().ScanByPartitionKeyWithTable(tableName, table, hashKeyValue, pkOpts, func(item *dbstore.Item) error {
			allItems = append(allItems, item)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// Storage walks arrive sorted (ascending forward, descending reverse);
	// the sorts stay as a bounded safety net over the ≤ limit+1 page.
	if scanIndexForward {
		sortItemsBySortKeyWithIndex(table, allItems, indexName)
	} else {
		sortItemsReverseBySortKeyWithIndex(table, allItems, indexName)
	}

	var scannedItems []*dbstore.Item
	if len(allItems) > limit {
		scannedItems = allItems[:limit]
	} else {
		scannedItems = allItems
	}

	scannedCount := len(scannedItems)

	filterExpr := request.GetStringParam(params, "FilterExpression")
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

	resp := map[string]interface{}{
		"Count":        len(items),
		"ScannedCount": scannedCount,
	}
	if !countOnly {
		resp["Items"] = buildItemsResponse(items)
	}
	if hasMoreItems && len(scannedItems) > 0 {
		resp["LastEvaluatedKey"] = buildLastEvaluatedKeyWithIndex(scannedItems[len(scannedItems)-1], table, indexName)
	}

	returnConsumedCapacity := getReturnConsumedCapacity(params)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		capacityUnits := float64(scannedCount) * rcuPerItem(consistentRead, indexName, table)
		isLSI := indexName != "" && !isGSI(table, indexName)
		resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithIndex(tableName, indexName, capacityUnits, isLSI)
	}

	// A Query is one read event on the queried partition regardless of how
	// many items it returns. Index-scoped key series are not modelled, so a
	// query served by a global secondary index is not attributed.
	if indexName == "" || !isGSI(table, indexName) {
		s.recordQueryContributorEvent(ctx, store, table, queryPKValue)
	}

	return resp, nil
}
