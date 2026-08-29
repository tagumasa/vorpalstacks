package dynamodb

import (
	"context"
	"errors"
	"net/http"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Condition checker callbacks
// ---------------------------------------------------------------------------

// ConditionChecker evaluates a condition expression against the existing
// item state within a transaction. It returns nil if the condition is
// satisfied, or ErrConditionalCheckFailed (or any other error) otherwise.
// When existing is nil and isNotFound is true, the caller should evaluate
// the condition against a synthetic empty item.
type ConditionChecker func(existing *dbstore.Item, isNotFound bool) error

// ---------------------------------------------------------------------------
// Core item functions — single validation + persistence path
//
// These functions encapsulate the full item lifecycle including DynamoDB
// Streams capture, Kinesis Data Stream destinations, and global table
// replication. Both the HTTP API handlers and the admin gRPC handler
// delegate to these functions to ensure identical side-effect behaviour.
// ---------------------------------------------------------------------------

// getItemCore retrieves a single item by primary key. It returns
// dbstore.ItemNotFound when the item does not exist.
func (s *DynamoDBService) getItemCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName string, key map[string]*dbstore.AttributeValue) (*dbstore.Item, error) {
	item, err := store.Items().Get(tableName, key)
	if err == nil {
		s.recordContributorReads(ctx, store, tableName, []map[string]*dbstore.AttributeValue{key})
	}
	return item, err
}

// recordContributorReads counts one read event per key in the contributor
// access aggregation of tables with contributor insights enabled. Reads
// count as one unit of ConsumedThroughputUnits. Failures are ignored:
// monitoring must never fail the read it observes.
func (s *DynamoDBService) recordContributorReads(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName string, keys []map[string]*dbstore.AttributeValue) {
	if len(keys) == 0 {
		return
	}
	if err := store.RecordContributorReads(ctx, tableName, keys); err != nil {
		logs.Warn("failed to record contributor reads",
			logs.String("table", tableName), logs.Err(err))
	}
}

// recordQueryContributorEvent counts the single read event a Query
// contributes to the partition-key series of tables with contributor
// insights enabled. Failures are ignored: monitoring must never fail the
// read it observes.
func (s *DynamoDBService) recordQueryContributorEvent(ctx context.Context, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, pkValue *dbstore.AttributeValue) {
	if !table.ContributorInsightsEnabled || pkValue == nil {
		return
	}
	pkName := ""
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			pkName = ks.AttributeName
			break
		}
	}
	if pkName == "" {
		return
	}
	if err := store.RecordContributorQuery(ctx, table.Name, map[string]*dbstore.AttributeValue{pkName: pkValue}); err != nil {
		logs.Warn("failed to record contributor query event",
			logs.String("table", table.Name), logs.Err(err))
	}
}

// ---------------------------------------------------------------------------
// Scan Core — paginated item listing
// ---------------------------------------------------------------------------

// ScanItemsInput is the service-layer DTO for a paginated item scan. It is
// used by the admin Scan handler and any other caller that needs a simple
// paginated listing without filter expressions, parallel segments, or
// projection — features that the data-plane Scan handler applies on top.
type ScanItemsInput struct {
	TableName string
	Limit     int
	Marker    string // empty string means first page
}

// ScanItemsResult is the service-layer result of a paginated item scan.
type ScanItemsResult struct {
	Items      []*dbstore.Item
	NextMarker string // empty string means no more pages
}

// scanItemsCore returns a single page of items for the specified table.
// It applies the standard default-and-cap limit logic and delegates to
// store.Items().List. The admin Scan handler calls this method so that
// no admin code path touches the store layer directly.
//
// The data-plane Scan handler (item_query.go) uses a richer code path with
// parallel-segment, filter-expression, and projection support; this Core
// method covers the simpler paginated-listing case needed by the admin
// console.
func (s *DynamoDBService) scanItemsCore(store dbstore.DynamoDBStoreInterface, in ScanItemsInput) (*ScanItemsResult, error) {
	lim := in.Limit
	if lim <= 0 {
		lim = dataPlaneQueryDefaultLimit
	}
	if lim > dataPlaneQueryMaxLimit {
		lim = dataPlaneQueryMaxLimit
	}
	items, nextMarker, err := store.Items().List(in.TableName, in.Marker, lim)
	if err != nil {
		return nil, err
	}
	return &ScanItemsResult{
		Items:      items,
		NextMarker: nextMarker,
	}, nil
}

// putItemCore creates or replaces an item, applying all side effects:
// DynamoDB Streams capture (within the transaction), Kinesis Data Stream
// destinations (after commit), and global table replication.
//
// The optional conditionChecker is invoked inside the transaction after
// loading the existing item. Pass nil to skip condition evaluation (admin
// console path).
//
// Returns the stored item and the previous item (nil if new).
func (s *DynamoDBService) putItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	table *dbstore.Table,
	key, item map[string]*dbstore.AttributeValue,
	conditionChecker ConditionChecker,
) (storedItem, oldItem *dbstore.Item, err error) {
	var isNew bool

	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, getErr := txn.GetItem(table.Name, key)
		if getErr != nil {
			if dbstore.IsItemNotFound(getErr) {
				isNew = true
			} else {
				return getErr
			}
		} else {
			oldItem = existingItem
			if oldItem != nil {
				if delErr := txn.DeleteIndexEntries(table.Name, oldItem); delErr != nil {
					return delErr
				}
			}
		}

		if conditionChecker != nil {
			if condErr := conditionChecker(existingItem, isNew); condErr != nil {
				return condErr
			}
		}

		if putErr := txn.PutItem(table.Name, key, item); putErr != nil {
			return putErr
		}

		storedItem = &dbstore.Item{
			TableName:  table.Name,
			Key:        key,
			Attributes: item,
		}
		if putIdxErr := txn.PutIndexEntries(table.Name, storedItem); putIdxErr != nil {
			return putIdxErr
		}

		newItemSize := calculateItemSize(item)
		if isNew {
			if upErr := txn.UpdateItemCount(table.Name, 1); upErr != nil {
				return upErr
			}
			if upErr := txn.UpdateTableSize(table.Name, newItemSize); upErr != nil {
				return upErr
			}
		} else if oldItem != nil {
			oldItemSize := calculateItemSize(oldItem.Attributes)
			if newItemSize != oldItemSize {
				if upErr := txn.UpdateTableSize(table.Name, newItemSize-oldItemSize); upErr != nil {
					return upErr
				}
			}
		}

		eventName := dbstore.StreamEventModify
		if isNew {
			eventName = dbstore.StreamEventInsert
		}
		s.captureStreamChangeTxn(txn, store, table, eventName, key, storedItem.Attributes, oldItemAttributes(oldItem))

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	eventName := dbstore.StreamEventModify
	if isNew {
		eventName = dbstore.StreamEventInsert
	}
	s.sendToKinesisDestinations(table, eventName, key, storedItem.Attributes, oldItemAttributes(oldItem))

	s.replicateToGlobalTableReplicas(store, region, table.Name, replicaPutOp(table, key, item))

	return storedItem, oldItem, nil
}

// deleteItemCore removes an item by primary key, applying all side effects:
// DynamoDB Streams capture (within the transaction), Kinesis Data Stream
// destinations (after commit), and global table replication.
//
// The optional conditionChecker is invoked inside the transaction after
// loading the existing item. Pass nil to skip condition evaluation (admin
// console path).
//
// Returns the previous item (nil if it did not exist).
func (s *DynamoDBService) deleteItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	table *dbstore.Table,
	key map[string]*dbstore.AttributeValue,
	conditionChecker ConditionChecker,
) (oldItem *dbstore.Item, err error) {
	err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, getErr := txn.GetItem(table.Name, key)
		if getErr != nil {
			if dbstore.IsItemNotFound(getErr) {
				if conditionChecker != nil {
					if condErr := conditionChecker(nil, true); condErr != nil {
						return condErr
					}
				}
				return nil
			}
			return getErr
		}

		if conditionChecker != nil {
			if condErr := conditionChecker(existingItem, false); condErr != nil {
				return condErr
			}
		}

		oldItem = existingItem

		if oldItem != nil {
			if delIdxErr := txn.DeleteIndexEntries(table.Name, oldItem); delIdxErr != nil {
				return delIdxErr
			}
			if delErr := txn.DeleteItem(table.Name, key); delErr != nil {
				return delErr
			}
			if upErr := txn.UpdateItemCount(table.Name, -1); upErr != nil {
				return upErr
			}
			oldItemSize := calculateItemSize(oldItem.Attributes)
			if upErr := txn.UpdateTableSize(table.Name, -oldItemSize); upErr != nil {
				return upErr
			}
			s.captureStreamChangeTxn(txn, store, table, dbstore.StreamEventRemove, key, nil, oldItem.Attributes)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if oldItem != nil {
		s.sendToKinesisDestinations(table, dbstore.StreamEventRemove, key, nil, oldItem.Attributes)
	}

	s.replicateToGlobalTableReplicas(store, region, table.Name, replicaDeleteOp(table, key))

	return oldItem, nil
}

// isItemNotFound reports whether the error is the store's item-not-found
// sentinel, so transport layers can branch on it without importing the store
// package.
func isItemNotFound(err error) bool {
	return dbstore.IsItemNotFound(err)
}

// itemCollectionWriteRef names one item collection a write touched: the
// table record plus the written item's primary key.
type itemCollectionWriteRef struct {
	tableName string
	table     *dbstore.Table
	key       map[string]*dbstore.AttributeValue
}

// itemCollectionKey resolves the partition key name and value of a written
// item from the table's key schema. An empty name or a nil value means no
// item collection key could be derived.
func itemCollectionKey(table *dbstore.Table, key map[string]*dbstore.AttributeValue) (string, *dbstore.AttributeValue) {
	if table == nil {
		return "", nil
	}
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			return ks.AttributeName, key[ks.AttributeName]
		}
	}
	return "", nil
}

// tableHasItemCollections reports whether the table tracks item
// collections: the Developer Guide scopes item collections to tables that
// have one or more local secondary indexes, and the write operations'
// documentation states ItemCollectionMetrics is not returned for a table
// without them.
func tableHasItemCollections(table *dbstore.Table) bool {
	return table != nil && len(table.LocalSecondaryIndexes) > 0
}

// itemCollectionMetricsEntry renders one ItemCollectionMetrics entry: the
// written item's partition key value as the ItemCollectionKey and the
// table's size estimate as the two-element SizeEstimateRangeGB.
func itemCollectionMetricsEntry(table *dbstore.Table, pkName string, pkValue *dbstore.AttributeValue) map[string]interface{} {
	sizeGB := float64(table.TableSizeBytes) / (1024.0 * 1024.0 * 1024.0)
	return map[string]interface{}{
		"ItemCollectionKey":   buildItemResponse(map[string]*dbstore.AttributeValue{pkName: pkValue}),
		"SizeEstimateRangeGB": []float64{sizeGB, sizeGB},
	}
}

// buildItemCollectionMetricsEntry renders the single-object
// ItemCollectionMetrics response member of the single-item writes
// (PutItem, UpdateItem, DeleteItem) for ReturnItemCollectionMetrics=SIZE.
// It returns nil when no entry applies: the table has no local secondary
// indexes or the written key carries no partition key value.
func buildItemCollectionMetricsEntry(table *dbstore.Table, key map[string]*dbstore.AttributeValue) map[string]interface{} {
	if !tableHasItemCollections(table) {
		return nil
	}
	pkName, pkValue := itemCollectionKey(table, key)
	if pkName == "" || pkValue == nil {
		return nil
	}
	return itemCollectionMetricsEntry(table, pkName, pkValue)
}

// buildItemCollectionMetricsPerTable renders the ItemCollectionMetrics
// response member for the SIZE setting of ReturnItemCollectionMetrics on
// the batched write planes (BatchWriteItem, TransactWriteItems): one entry
// per distinct item collection (table plus partition key value), carrying
// the same entry the single-item writes report. A delete of a non-existent
// key is a successful no-op here, and the documented contract — item
// collections "affected by individual DeleteItem or PutItem operations" —
// does not determine whether such a delete affects its collection, so it
// registers like any other successful write.
func buildItemCollectionMetricsPerTable(writes []itemCollectionWriteRef) map[string]interface{} {
	perTable := make(map[string]interface{})
	seen := make(map[string]bool)
	for _, w := range writes {
		if !tableHasItemCollections(w.table) {
			continue
		}
		pkName, pkValue := itemCollectionKey(w.table, w.key)
		if pkName == "" || pkValue == nil {
			continue
		}
		identity := w.tableName + "\x00" + avToString(pkValue)
		if seen[identity] {
			continue
		}
		seen[identity] = true
		entry := itemCollectionMetricsEntry(w.table, pkName, pkValue)
		existing, _ := perTable[w.tableName].([]map[string]interface{})
		perTable[w.tableName] = append(existing, entry)
	}
	if len(perTable) == 0 {
		return nil
	}
	return perTable
}

// ---------------------------------------------------------------------------
// Scan page Core — shared paginator for the Scan data plane
// ---------------------------------------------------------------------------

// scanPageOptions configures collectScanPage, the shared paginator for the
// Scan data plane and the index fallback reads of Query.
type scanPageOptions struct {
	table             *dbstore.Table
	indexName         string // secondary index whose membership filters the page; "" scans the base table
	limit             int    // page size; iteration stops after limit+1 qualifying items
	segment           int    // parallel-scan segment; -1 disables segment filtering
	totalSegments     int
	exclusiveStartKey map[string]*dbstore.AttributeValue
}

// scanPageResult carries one page of qualifying items plus the continuation
// flag derived from the limit+1 lookahead item.
type scanPageResult struct {
	items   []*dbstore.Item
	hasMore bool
}

// collectScanPage walks the table in storage order and applies, per item, the
// parallel-scan segment filter, the secondary-index membership filter, and the
// exclusive-start-key skip, collecting at most limit+1 qualifying items. The
// lookahead item only sets hasMore; callers truncate to limit. Applying the
// filters during iteration — before the limit is counted — is what keeps
// index scans paginable: a page bounded before filtering would strand later
// index members behind a LastEvaluatedKey that is never emitted.
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
	started := opts.exclusiveStartKey == nil

	_, err := store.Items().ScanWithOptions(tableName, dbstore.ScanOptions{}, func(item *dbstore.Item) error {
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

		if !started {
			itemStartKey := mergeIndexKey(item, opts.table, opts.indexName)
			if itemKeyMatches(itemStartKey, opts.exclusiveStartKey) {
				started = true
				return nil
			}
			if itemKeySortsAfter(itemStartKey, opts.exclusiveStartKey, opts.table, opts.indexName) {
				started = true
			} else {
				return nil
			}
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

	indexName := request.GetStringParam(params, "IndexName")
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
	consistentRead := request.GetBoolParam(params, "ConsistentRead")
	if indexName != "" && isGSI(table, indexName) && consistentRead {
		return nil, NewAPIError("com.amazon.coral.validate#ValidationException",
			"Consistent reads are not supported on global secondary indexes", http.StatusBadRequest)
	}
	// An explicit Limit below 1 is invalid (valid range minimum 1); a plain
	// value check cannot tell "unset" apart from an explicit zero, so
	// presence is checked first.
	if _, ok := params["Limit"]; ok && request.GetIntParam(params, "Limit") <= 0 {
		return nil, ErrInvalidParameter
	}
	limit := request.GetIntParam(params, "Limit")
	if limit <= 0 {
		limit = dataPlaneQueryDefaultLimit
	}
	if limit > dataPlaneQueryMaxLimit {
		limit = dataPlaneQueryMaxLimit
	}
	exclusiveStartKey, eskErr := parseExclusiveStartKey(params)
	if eskErr != nil {
		return nil, eskErr
	}
	if exclusiveStartKey != nil {
		if err := validateKeyTypes(table, exclusiveStartKey); err != nil {
			return nil, err
		}
	}
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

	projection, projErr := parseProjectionExpression(params)
	if projErr != nil {
		return nil, projErr
	}
	countOnly, allProjected, selErr := parseSelectParam(params, indexName, projection != nil)
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
		hashKeyValue, hashKeyAttr, sortKeyCondition := extractIndexKeyCondition(table, indexName, keyCondExpr, exprAttrNames, exprAttrValues)
		queryPKValue = hashKeyAttr
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

// ---------------------------------------------------------------------------
// Batch Core — BatchGetItem / BatchWriteItem
// ---------------------------------------------------------------------------

// batchGetItemInput carries the already-typed RequestItems member plus the
// raw wire parameters (consumed for the ReturnConsumedCapacity reporting).
type batchGetItemInput struct {
	RequestItems map[string]interface{}
	Parameters   map[string]interface{}
}

// batchGetItemCore is the single validation and persistence path of
// BatchGetItem: per-table key parsing, duplicate detection, reads, projection,
// and the UnprocessedKeys reporting.
func (s *DynamoDBService) batchGetItemCore(ctx context.Context, reqCtx *request.RequestContext, in batchGetItemInput) (map[string]interface{}, error) {
	requestItems := in.RequestItems

	totalKeys := 0
	for _, tableReq := range requestItems {
		if tr, ok := tableReq.(map[string]interface{}); ok {
			if keys, ok := tr["Keys"].([]interface{}); ok {
				totalKeys += len(keys)
			}
		}
	}
	if totalKeys > batchGetMaxTotalItems {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	responses := make(map[string]interface{})
	unprocessed := make(map[string]interface{})
	tableReadCounts := make(map[string]int)
	tableConsistentRead := make(map[string]bool)
	seenGetKeys := make(map[string]bool)

	for tableName, tableReq := range requestItems {
		tr, ok := tableReq.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		if !store.Tables().Exists(tableName) {
			return nil, ErrTableNotFound
		}

		batchTable, tblErr := store.Tables().Get(tableName)
		if tblErr != nil || batchTable == nil {
			return nil, ErrTableNotFound
		}

		keys, ok := tr["Keys"].([]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		// ConsistentRead is accepted for API compatibility. Single-instance
		// Pebble provides strong consistency for all reads; the flag is
		// honoured in the reported capacity charge.
		tableConsistentRead[tableName] = request.GetBoolParam(tr, "ConsistentRead")

		var tableItems []map[string]interface{}
		var unprocessedKeys []interface{}
		var foundKeys []map[string]*dbstore.AttributeValue

		for _, k := range keys {
			key, keyErr := parseKey(k)
			if keyErr != nil || key == nil {
				unprocessedKeys = append(unprocessedKeys, k)
				continue
			}

			// A wrong-typed key rejects the whole batch request, matching
			// the BatchGetItem contract.
			if err := validateKeyTypes(batchTable, key); err != nil {
				return nil, err
			}

			// Duplicate keys in one request are rejected rather than
			// deduplicated.
			keyStr := buildKeyString(tableName, key)
			if seenGetKeys[keyStr] {
				return nil, ErrDuplicateKeys
			}
			seenGetKeys[keyStr] = true

			item, err := store.Items().Get(tableName, key)
			if err != nil {
				if isItemNotFound(err) {
					// Requests for nonexistent items consume the minimum
					// read capacity units according to the read type.
					tableReadCounts[tableName]++
					continue
				}
				unprocessedKeys = append(unprocessedKeys, k)
				continue
			}
			tableReadCounts[tableName]++
			foundKeys = append(foundKeys, key)

			projection, projErr := parseProjectionExpression(tr)
			if projErr != nil {
				unprocessedKeys = append(unprocessedKeys, k)
				continue
			}
			if projection != nil {
				item.Attributes = applyProjection(item.Attributes, projection)
			}

			tableItems = append(tableItems, buildItemResponse(item.Attributes))
		}

		if len(tableItems) > 0 {
			responses[tableName] = tableItems
		}
		if len(unprocessedKeys) > 0 {
			unprocessed[tableName] = map[string]interface{}{"Keys": unprocessedKeys}
		}
		// Every item the batch actually read counts as one read event per
		// tracked key layout; requests for nonexistent items read nothing.
		s.recordContributorReads(ctx, store, tableName, foundKeys)
	}

	resp := map[string]interface{}{
		"Responses":       responses,
		"UnprocessedKeys": unprocessed,
	}

	returnConsumedCapacity := getReturnConsumedCapacity(in.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		var consumedCapacity []interface{}
		for tableName, readCount := range tableReadCounts {
			capacityUnits := float64(readCount) * rcuPerItem(tableConsistentRead[tableName], "", nil)
			consumedCapacity = append(consumedCapacity, buildConsumedCapacityResponse(tableName, capacityUnits))
		}
		if len(consumedCapacity) > 0 {
			resp["ConsumedCapacity"] = consumedCapacity
		}
	}

	return resp, nil
}

// batchWriteItemInput carries the already-typed RequestItems member plus the
// raw wire parameters (consumed for the ReturnConsumedCapacity reporting).
type batchWriteItemInput struct {
	RequestItems map[string]interface{}
	Parameters   map[string]interface{}
}

// batchWriteItemCore is the single validation and persistence path of
// BatchWriteItem: per-table write parsing, duplicate-key rejection,
// transactional writes with stream capture, and the asynchronous Kinesis and
// global-table replication side effects.
func (s *DynamoDBService) batchWriteItemCore(ctx context.Context, reqCtx *request.RequestContext, in batchWriteItemInput) (map[string]interface{}, error) {
	requestItems := in.RequestItems

	totalItems := 0
	for _, tableReq := range requestItems {
		if writes, ok := tableReq.([]interface{}); ok {
			totalItems += len(writes)
		}
	}
	if totalItems > batchWriteMaxItems {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	unprocessed := make(map[string]interface{})
	var metricsWrites []itemCollectionWriteRef

	type writeOp struct {
		tableName string
		opType    string
		key       map[string]*dbstore.AttributeValue
		item      map[string]*dbstore.AttributeValue
		rawReq    map[string]interface{}
	}

	var allWrites []writeOp
	tableCache := make(map[string]*dbstore.Table)
	// Primary keys already targeted in this request, per table: the whole
	// batch write is rejected when the same item appears twice, whether as
	// two puts or as a put plus a delete.
	seenWriteKeys := make(map[string]bool)

	for tableName, tableReq := range requestItems {
		writes, ok := tableReq.([]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		table, tableErr := store.Tables().Get(tableName)
		if tableErr != nil {
			if dbstore.IsTableNotFound(tableErr) {
				return nil, ErrTableNotFound
			}
			unprocessed[tableName] = writes
			continue
		}
		tableCache[tableName] = table

		for _, w := range writes {
			writeReq, ok := w.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidParameter
			}

			if putReq, ok := writeReq["PutRequest"].(map[string]interface{}); ok {
				item, itemErr := parseItem(putReq["Item"])
				if itemErr != nil || item == nil {
					return nil, ErrInvalidParameter
				}

				// Item size must not exceed 400 KB (same as PutItem).
				if calculateItemSize(item) > maxItemSizeBytes {
					return nil, ErrInvalidParameter
				}

				key := s.extractKeyFromItem(table, item)
				if key == nil {
					return nil, ErrInvalidParameter
				}

				// Key attribute values must not be empty.
				if !validateKeyValueNotEmpty(key) {
					return nil, ErrInvalidParameter
				}

				// Key attribute types must match the table schema and any
				// index key definitions.
				if err := validateItemKeyTypes(table, item); err != nil {
					return nil, err
				}

				keyStr := buildKeyString(tableName, key)
				if seenWriteKeys[keyStr] {
					return nil, ErrDuplicateKeys
				}
				seenWriteKeys[keyStr] = true

				allWrites = append(allWrites, writeOp{
					tableName: tableName,
					opType:    "Put",
					key:       key,
					item:      item,
					rawReq:    writeReq,
				})
			}

			if delReq, ok := writeReq["DeleteRequest"].(map[string]interface{}); ok {
				key, keyErr := parseKey(delReq["Key"])
				if keyErr != nil || key == nil {
					return nil, ErrInvalidParameter
				}

				// Key attribute values must not be empty.
				if !validateKeyValueNotEmpty(key) {
					return nil, ErrInvalidParameter
				}

				if err := validateKeyTypes(table, key); err != nil {
					return nil, err
				}

				keyStr := buildKeyString(tableName, key)
				if seenWriteKeys[keyStr] {
					return nil, ErrDuplicateKeys
				}
				seenWriteKeys[keyStr] = true

				allWrites = append(allWrites, writeOp{
					tableName: tableName,
					opType:    "Delete",
					key:       key,
					rawReq:    writeReq,
				})
			}
		}
	}

	for _, op := range allWrites {
		var batchIsNew bool
		var batchOldAttrs map[string]*dbstore.AttributeValue

		opErr := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			switch op.opType {
			case "Put":
				existingItem, err := txn.GetItem(op.tableName, op.key)
				isNewItem := dbstore.IsItemNotFound(err)
				batchIsNew = isNewItem
				if err != nil && !isNewItem {
					return err
				}
				if existingItem != nil {
					batchOldAttrs = existingItem.Attributes
					if err := txn.DeleteIndexEntries(op.tableName, existingItem); err != nil {
						return err
					}
				}
				if err := txn.PutItem(op.tableName, op.key, op.item); err != nil {
					return err
				}
				newItem := &dbstore.Item{
					TableName:  op.tableName,
					Key:        op.key,
					Attributes: op.item,
				}
				if err := txn.PutIndexEntries(op.tableName, newItem); err != nil {
					return err
				}
				newItemSize := calculateItemSize(op.item)
				if isNewItem {
					if err := txn.UpdateItemCount(op.tableName, 1); err != nil {
						return err
					}
					if err := txn.UpdateTableSize(op.tableName, newItemSize); err != nil {
						return err
					}
				} else if existingItem != nil {
					oldItemSize := calculateItemSize(existingItem.Attributes)
					if err := txn.UpdateTableSize(op.tableName, newItemSize-oldItemSize); err != nil {
						return err
					}
				}

			case "Delete":
				existingItem, err := txn.GetItem(op.tableName, op.key)
				if dbstore.IsItemNotFound(err) {
					return nil
				}
				if err != nil {
					return err
				}
				if existingItem != nil {
					batchOldAttrs = existingItem.Attributes
					if err := txn.DeleteIndexEntries(op.tableName, existingItem); err != nil {
						return err
					}
				}
				if err := txn.DeleteItem(op.tableName, op.key); err != nil {
					return err
				}
				if err := txn.UpdateItemCount(op.tableName, -1); err != nil {
					return err
				}
				if existingItem != nil {
					oldItemSize := calculateItemSize(existingItem.Attributes)
					if err := txn.UpdateTableSize(op.tableName, -oldItemSize); err != nil {
						return err
					}
				}
			}

			table := tableCache[op.tableName]
			if op.opType == "Put" {
				eventName := dbstore.StreamEventModify
				if batchIsNew {
					eventName = dbstore.StreamEventInsert
				}
				s.captureStreamChangeTxn(txn, store, table, eventName, op.key, op.item, batchOldAttrs)
			} else if op.opType == "Delete" && batchOldAttrs != nil {
				s.captureStreamChangeTxn(txn, store, table, dbstore.StreamEventRemove, op.key, nil, batchOldAttrs)
			}

			return nil
		})

		if opErr != nil {
			var unprocessedItems []interface{}
			if existing, ok := unprocessed[op.tableName].([]interface{}); ok {
				unprocessedItems = existing
			}
			unprocessed[op.tableName] = append(unprocessedItems, op.rawReq)
		} else {
			table := tableCache[op.tableName]
			metricsWrites = append(metricsWrites, itemCollectionWriteRef{tableName: op.tableName, table: table, key: op.key})
			if op.opType == "Put" {
				eventName := dbstore.StreamEventModify
				if batchIsNew {
					eventName = dbstore.StreamEventInsert
				}
				s.sendToKinesisDestinations(table, eventName, op.key, op.item, batchOldAttrs)
				if table != nil {
					s.replicateToGlobalTableReplicas(store, reqCtx.GetRegion(), op.tableName, replicaPutOp(table, op.key, op.item))
				}
			} else if op.opType == "Delete" {
				if batchOldAttrs != nil {
					s.sendToKinesisDestinations(table, dbstore.StreamEventRemove, op.key, nil, batchOldAttrs)
				}
				if table != nil {
					s.replicateToGlobalTableReplicas(store, reqCtx.GetRegion(), op.tableName, replicaDeleteOp(table, op.key))
				}
			}
		}
	}

	resp := map[string]interface{}{
		"UnprocessedItems": unprocessed,
	}

	// ReturnItemCollectionMetrics=SIZE asks for one entry per item
	// collection the batch actually wrote.
	if request.GetStringParam(in.Parameters, "ReturnItemCollectionMetrics") == "SIZE" {
		if metrics := buildItemCollectionMetricsPerTable(metricsWrites); metrics != nil {
			resp["ItemCollectionMetrics"] = metrics
		}
	}

	returnConsumedCapacity := getReturnConsumedCapacity(in.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		var consumedCapacity []interface{}
		for tableName := range requestItems {
			consumedCapacity = append(consumedCapacity, buildConsumedCapacityResponse(tableName, 1.0))
		}
		if len(consumedCapacity) > 0 {
			resp["ConsumedCapacity"] = consumedCapacity
		}
	}

	return resp, nil
}

// ---------------------------------------------------------------------------
// Transport-agnostic DTO for UpdateItem
// ---------------------------------------------------------------------------

// UpdateItemInput carries every field that UpdateItem needs. Both the HTTP
// API handler and admin gRPC handler build this struct and delegate to
// updateItemCore.
type UpdateItemInput struct {
	Key            map[string]*dbstore.AttributeValue
	UpdateExpr     string
	AttrUpdates    interface{} // legacy AttributeUpdates map (mutually exclusive with UpdateExpr)
	ConditionExpr  string
	ExprAttrNames  map[string]string
	ExprAttrValues map[string]*dbstore.AttributeValue
	ReturnValues   string
}

// UpdateItemResult holds the output of updateItemCore for response formatting.
type UpdateItemResult struct {
	StoredItem       *dbstore.Item
	OldItem          *dbstore.Item
	UpdatedAttrNames []string
	WasNewItem       bool
}

// updateItemCore is the single entry point for item updates shared by the
// HTTP API and admin gRPC handler. It evaluates conditions, applies update
// expressions, persists within a transaction, and fires all side effects
// (Streams, Kinesis, global table replication).
func (s *DynamoDBService) updateItemCore(
	ctx context.Context,
	store dbstore.DynamoDBStoreInterface,
	region string,
	table *dbstore.Table,
	in UpdateItemInput,
) (*UpdateItemResult, error) {
	tableName := table.Name

	var oldItem *dbstore.Item
	var storedItem *dbstore.Item
	var updatedAttrNames []string
	var oldItemSize int64
	var wasNewItem bool

	err := store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		existingItem, err := txn.GetItem(tableName, in.Key)
		isNewItem := false
		var item *dbstore.Item

		if err != nil {
			if dbstore.IsItemNotFound(err) {
				isNewItem = true
				if in.ConditionExpr != "" {
					syntheticItem := &dbstore.Item{
						TableName:  tableName,
						Key:        in.Key,
						Attributes: make(map[string]*dbstore.AttributeValue),
					}
					conditionMet, evalErr := evaluateConditionExpression(syntheticItem, in.ConditionExpr, in.ExprAttrNames, in.ExprAttrValues)
					if evalErr != nil {
						return evalErr
					}
					if !conditionMet {
						return ErrConditionalCheckFailed
					}
				}
				item = &dbstore.Item{
					TableName:  tableName,
					Key:        in.Key,
					Attributes: make(map[string]*dbstore.AttributeValue),
				}
				for k, v := range in.Key {
					item.Attributes[k] = v
				}
			} else {
				return err
			}
		} else {
			item = existingItem
			oldItemSize = calculateItemSize(item.Attributes)
			if in.ConditionExpr != "" {
				conditionMet, err := evaluateConditionExpression(item, in.ConditionExpr, in.ExprAttrNames, in.ExprAttrValues)
				if err != nil {
					return err
				}
				if !conditionMet {
					return ErrConditionalCheckFailed
				}
			}

			streamNeedsOld := table.StreamSpecification != nil &&
				table.StreamSpecification.StreamEnabled &&
				(table.StreamSpecification.StreamViewType == dbstore.StreamViewTypeOldImage ||
					table.StreamSpecification.StreamViewType == dbstore.StreamViewTypeNewAndOldImages)
			if in.ReturnValues == "ALL_OLD" || in.ReturnValues == "UPDATED_OLD" || streamNeedsOld {
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

		if in.UpdateExpr != "" {
			paths := extractUpdatedPaths(in.UpdateExpr, in.ExprAttrNames)
			if err := validateNotKeyAttributes(table, paths); err != nil {
				return err
			}
			var err error
			updatedAttrNames, err = applyUpdateExpressionWithTracking(item.Attributes, in.UpdateExpr, in.ExprAttrNames, in.ExprAttrValues)
			if err != nil {
				if errors.Is(err, ErrTypeMismatch) {
					return ErrInvalidParameter
				}
				return err
			}
		} else if in.AttrUpdates != nil {
			var attrNames []string
			if attrMap, ok := in.AttrUpdates.(map[string]interface{}); ok {
				for k := range attrMap {
					attrNames = append(attrNames, k)
				}
			}
			if err := validateNotKeyAttributes(table, attrNames); err != nil {
				return err
			}
			updatedAttrNames, err = applyAttributeUpdatesWithTracking(item.Attributes, in.AttrUpdates)
			if err != nil {
				return err
			}
		}

		if itemSize := calculateItemSize(item.Attributes); itemSize > maxItemSizeBytes {
			return ErrInvalidParameter
		}

		if err := txn.PutItem(tableName, in.Key, item.Attributes); err != nil {
			return err
		}

		storedItem = &dbstore.Item{
			TableName:  tableName,
			Key:        in.Key,
			Attributes: item.Attributes,
		}
		if err := txn.PutIndexEntries(tableName, storedItem); err != nil {
			return err
		}

		newItemSize := calculateItemSize(item.Attributes)
		wasNewItem = isNewItem
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

		eventName := dbstore.StreamEventModify
		if wasNewItem {
			eventName = dbstore.StreamEventInsert
		}
		s.captureStreamChangeTxn(txn, store, table, eventName, in.Key, storedItem.Attributes, oldItemAttributes(oldItem))

		return nil
	})
	if err != nil {
		return nil, err
	}

	{
		eventName := dbstore.StreamEventModify
		if wasNewItem {
			eventName = dbstore.StreamEventInsert
		}
		s.sendToKinesisDestinations(table, eventName, in.Key, storedItem.Attributes, oldItemAttributes(oldItem))
	}

	repKey := in.Key
	repAttrs := storedItem.Attributes
	repItemSize := calculateItemSize(repAttrs)
	s.replicateToGlobalTableReplicas(store, region, tableName, func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		return destStore.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			existing, getErr := txn.GetItem(tableName, repKey)
			isNew := false
			if getErr != nil {
				if dbstore.IsItemNotFound(getErr) {
					isNew = true
				} else {
					return getErr
				}
			}
			if !isNew && existing != nil {
				if err := txn.DeleteIndexEntries(tableName, existing); err != nil {
					return err
				}
			}
			if err := txn.PutItem(tableName, repKey, repAttrs); err != nil {
				return err
			}
			repItem := &dbstore.Item{
				TableName:  tableName,
				Key:        repKey,
				Attributes: repAttrs,
			}
			if err := txn.PutIndexEntries(tableName, repItem); err != nil {
				return err
			}
			if isNew {
				if err := txn.UpdateItemCount(tableName, 1); err != nil {
					return err
				}
				if err := txn.UpdateTableSize(tableName, repItemSize); err != nil {
					return err
				}
			} else {
				oldSize := calculateItemSize(existing.Attributes)
				if err := txn.UpdateTableSize(tableName, repItemSize-oldSize); err != nil {
					return err
				}
			}
			return nil
		})
	})

	return &UpdateItemResult{
		StoredItem:       storedItem,
		OldItem:          oldItem,
		UpdatedAttrNames: updatedAttrNames,
		WasNewItem:       wasNewItem,
	}, nil
}
