package dynamodb

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"net/http"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// scanCoreResult is the typed page a Scan produced, before wire or proto
// serialisation.
type scanCoreResult struct {
	// Items are the items surviving the filter, with the projection already
	// applied. IncludeItems reports whether Select asked for them at all
	// (Select=COUNT reads the page but returns no Items member).
	Items        []*dbstore.Item
	IncludeItems bool
	Count        int // items surviving the filter
	ScannedCount int // items read before the filter
	// LastEvaluatedKey is the primary key of the last scanned item, extended
	// with the index key attributes for an index scan; nil means no further
	// pages.
	LastEvaluatedKey map[string]*dbstore.AttributeValue
	IndexName        string
	CapacityUnits    float64
}

// readIndexPreamble is the result of the validation prefix the Scan and
// Query Cores share.
type readIndexPreamble struct {
	indexName         string
	consistentRead    bool
	limit             int
	exclusiveStartKey map[string]*dbstore.AttributeValue
}

// resolveReadIndexPreamble applies the validation prefix the Scan and Query
// data planes share, in the order both planes apply it: index resolution,
// the ConsistentRead GSI rejection, the Limit presence check and clamp, and
// the exclusive-start-key parse with its shape validation.
func resolveReadIndexPreamble(table *dbstore.Table, params map[string]interface{}) (*readIndexPreamble, error) {
	indexName := request.GetStringParam(params, "IndexName")
	if indexName != "" {
		if !validateResourceName(indexName) {
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
	return &readIndexPreamble{
		indexName:         indexName,
		consistentRead:    consistentRead,
		limit:             limit,
		exclusiveStartKey: exclusiveStartKey,
	}, nil
}

// resolveProjectionSelection applies the projection/Select validation the
// Scan and Query data planes share: the projection expression parse, the
// Select resolution, and the GSI projection rules.
func resolveProjectionSelection(table *dbstore.Table, indexName string, params map[string]interface{}) (projection []string, countOnly, allProjected bool, err error) {
	projection, projErr := parseProjectionExpression(params)
	if projErr != nil {
		return nil, false, false, projErr
	}
	countOnly, allProjected, selErr := parseSelectParam(params, indexName, projection != nil)
	if selErr != nil {
		return nil, false, false, selErr
	}
	if indexName != "" && isGSI(table, indexName) {
		// A global secondary index can only serve attributes projected into
		// it; local secondary indexes may fall back to the parent table.
		if gsiProjErr := validateGSIProjectionRequest(table, indexName, allProjected, countOnly, projection); gsiProjErr != nil {
			return nil, false, false, gsiProjErr
		}
	}
	return projection, countOnly, allProjected, nil
}

// scanCore is the single validation and persistence path of the Scan data
// plane and the admin console listing: index resolution, ConsistentRead and
// parallel-segment validation, the paginated read, filter and projection
// application, and the contributor read accounting. The caller resolves the
// store and the table; every other parameter arrives in its wire shape so
// the admin plane can share the same validation by supplying a minimal
// parameter map (table, limit, exclusive start key only).
func (s *DynamoDBService) scanCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, params map[string]interface{}) (*scanCoreResult, error) {
	tableName := table.Name

	preamble, preErr := resolveReadIndexPreamble(table, params)
	if preErr != nil {
		return nil, preErr
	}
	indexName := preamble.indexName
	consistentRead := preamble.consistentRead
	limit := preamble.limit

	// The start key is "the primary key of the first item that this
	// operation will evaluate": it must carry the full primary key so the
	// page can resume at that item's storage position instead of re-reading
	// the head of the table.
	startMarker := ""
	if preamble.exclusiveStartKey != nil {
		startKey := primaryKeyFromStartKey(table, preamble.exclusiveStartKey)
		if startKey == nil {
			return nil, ErrInvalidParameter
		}
		startMarker = dbstore.EncodeItemKey(table.Name, startKey, table)
	}

	// Validate parallel Scan parameters (Smithy ScanSegment: range 0-999999,
	// ScanTotalSegments: range 1-1000000). The two parameters form a pair:
	// specifying one without the other is a validation error.
	_, hasSegment := params["Segment"]
	_, hasTotalSegments := params["TotalSegments"]
	if hasSegment != hasTotalSegments {
		return nil, ErrInvalidParameter
	}
	segment := -1
	totalSegments := 0
	if _, ok := params["Segment"]; ok {
		segment = request.GetIntParam(params, "Segment")
		if !validateScanSegment(segment) {
			return nil, ErrInvalidParameter
		}
	}
	if _, ok := params["TotalSegments"]; ok {
		totalSegments = request.GetIntParam(params, "TotalSegments")
		if !validateScanTotalSegments(totalSegments) {
			return nil, ErrInvalidParameter
		}
	}

	projection, countOnly, allProjected, projErr := resolveProjectionSelection(table, indexName, params)
	if projErr != nil {
		return nil, projErr
	}

	// The page is built during iteration — segment filtering and secondary
	// index membership apply before the limit is counted, so every
	// qualifying item remains reachable through pagination regardless of how
	// the surrounding table entries are interleaved; the start marker
	// resumes the walk strictly after the previous page's last item.
	page, pageErr := s.collectScanPage(store, tableName, scanPageOptions{
		table:         table,
		indexName:     indexName,
		limit:         limit,
		segment:       segment,
		totalSegments: totalSegments,
		marker:        startMarker,
	})
	if pageErr != nil {
		return nil, pageErr
	}
	scannedItems := page.items
	scannedCount := len(scannedItems)

	filterExpr := request.GetStringParam(params, "FilterExpression")
	var items []*dbstore.Item
	if filterExpr != "" {
		scanNames, namesErr := parseExpressionAttributeNames(params)
		if namesErr != nil {
			return nil, namesErr
		}
		scanValues, scanValsErr := parseExpressionAttributeValues(params)
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

	result := &scanCoreResult{
		Items:         items,
		IncludeItems:  !countOnly,
		Count:         len(items),
		ScannedCount:  scannedCount,
		IndexName:     indexName,
		CapacityUnits: float64(scannedCount) * rcuPerItem(consistentRead, indexName, table),
	}
	if page.hasMore && len(scannedItems) > 0 {
		result.LastEvaluatedKey = mergeIndexKey(scannedItems[len(scannedItems)-1], table, indexName)
	}

	// Every item the page read counts as one read event per tracked key
	// layout, regardless of the filter expression applied afterwards.
	scanReadKeys := make([]map[string]*dbstore.AttributeValue, 0, len(scannedItems))
	for _, item := range scannedItems {
		scanReadKeys = append(scanReadKeys, item.Key)
	}
	s.recordContributorReads(ctx, store, tableName, scanReadKeys)

	return result, nil
}

// md5SegmentHash computes the MD5 hash of an AttributeValue for parallel
// scan segment assignment. AWS does not document how segments are
// assigned, so any deterministic hash over the partition value is
// behaviour-compatible; the hash covers the raw value bytes (S string,
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
