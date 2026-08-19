package dynamodb

import (
	"net/http"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

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

// indexKeyAttributeNames resolves the hash and sort key attribute names of a
// secondary index, reporting whether the index is global.
func indexKeyAttributeNames(table *dbstore.Table, indexName string) (hashName, sortName string, gsi bool) {
	for _, idx := range table.GlobalSecondaryIndexes {
		if idx.IndexName == indexName {
			gsi = true
			for _, key := range idx.KeySchema {
				if key.KeyType == dbstore.KeyTypeHash {
					hashName = key.AttributeName
				} else if key.KeyType == dbstore.KeyTypeRange {
					sortName = key.AttributeName
				}
			}
			return hashName, sortName, gsi
		}
	}
	for _, idx := range table.LocalSecondaryIndexes {
		if idx.IndexName == indexName {
			for _, key := range idx.KeySchema {
				if key.KeyType == dbstore.KeyTypeHash {
					hashName = key.AttributeName
				} else if key.KeyType == dbstore.KeyTypeRange {
					sortName = key.AttributeName
				}
			}
			return hashName, sortName, gsi
		}
	}
	return "", "", false
}

// isIndexMember reports whether the item is projected into the index being
// scanned. Global indexes require the index hash key (and sort key when
// declared); local indexes share the base hash key and require only the index
// sort key.
func isIndexMember(item *dbstore.Item, hashName, sortName string, gsi bool) bool {
	if item == nil {
		return false
	}
	if gsi {
		if _, hasHash := item.Attributes[hashName]; !hasHash {
			return false
		}
		if sortName != "" {
			if _, hasSort := item.Attributes[sortName]; !hasSort {
				return false
			}
		}
		return true
	}
	if sortName != "" {
		_, exists := item.Attributes[sortName]
		return exists
	}
	if hashName != "" {
		_, exists := item.Attributes[hashName]
		return exists
	}
	return true
}

// rcuPerItem returns the read capacity units charged per evaluated item.
// Strongly consistent base-table and local-index reads consume 1.0 unit per
// item, eventually consistent reads 0.5. Global secondary indexes are
// eventually consistent only, so they always consume 0.5.
func rcuPerItem(consistentRead bool, indexName string, table *dbstore.Table) float64 {
	if indexName != "" && isGSI(table, indexName) {
		return 0.5
	}
	if consistentRead {
		return 1.0
	}
	return 0.5
}

// indexProjectedAttributes returns the attribute names available from an
// index read: the primary key, the index keys, and for INCLUDE projections
// the listed non-key attributes. A nil return means the projection is ALL
// and the whole item is available.
func indexProjectedAttributes(table *dbstore.Table, indexName string) map[string]bool {
	var proj *dbstore.Projection
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			proj = gsi.Projection
			break
		}
	}
	if proj == nil {
		for _, lsi := range table.LocalSecondaryIndexes {
			if lsi.IndexName == indexName {
				proj = lsi.Projection
				break
			}
		}
	}
	if proj == nil || proj.ProjectionType == ProjectionTypeAll {
		return nil
	}
	names := map[string]bool{}
	for _, ks := range table.KeySchema {
		names[ks.AttributeName] = true
	}
	hashName, sortName, _ := indexKeyAttributeNames(table, indexName)
	if hashName != "" {
		names[hashName] = true
	}
	if sortName != "" {
		names[sortName] = true
	}
	for _, name := range proj.NonKeyAttributes {
		names[name] = true
	}
	return names
}

// applyIndexProjection trims item attributes to the set projected into the
// index. ALL projections keep every attribute.
func applyIndexProjection(attrs map[string]*dbstore.AttributeValue, table *dbstore.Table, indexName string) map[string]*dbstore.AttributeValue {
	if indexName == "" {
		return attrs
	}
	projected := indexProjectedAttributes(table, indexName)
	if projected == nil {
		return attrs
	}
	result := make(map[string]*dbstore.AttributeValue, len(projected))
	for name, value := range attrs {
		if projected[name] {
			result[name] = value
		}
	}
	return result
}

// validateGSIProjectionRequest enforces that a global secondary index read
// only names attributes projected into the index: ALL_ATTRIBUTES is only
// valid when the index projects everything, and a ProjectionExpression may
// only reference projected attribute names. Local secondary indexes are not
// restricted because they can fetch unprojected attributes from the parent
// table.
func validateGSIProjectionRequest(table *dbstore.Table, indexName string, allProjected, countOnly bool, projection []string) error {
	projected := indexProjectedAttributes(table, indexName)
	if projected == nil {
		return nil
	}
	if !allProjected && !countOnly && projection == nil {
		return NewAPIError("com.amazon.coral.validate#ValidationException",
			"One or more parameter values were invalid: Select value ALL_ATTRIBUTES is not supported for global secondary index because not all attributes are projected into the index", http.StatusBadRequest)
	}
	for _, name := range projection {
		if !projected[name] {
			return NewAPIError("com.amazon.coral.validate#ValidationException",
				"One or more parameter values were invalid: attribute "+name+" is not projected into the global secondary index", http.StatusBadRequest)
		}
	}
	return nil
}

// mergeIndexKey returns the item's primary key extended with the secondary
// index key attributes. LastEvaluatedKey for an index read is composed the
// same way, so ExclusiveStartKey comparisons for index scans and queries must
// run against this merged form — comparing the bare primary key would never
// match because the index key attributes live in the item's attribute map.
func mergeIndexKey(item *dbstore.Item, table *dbstore.Table, indexName string) map[string]*dbstore.AttributeValue {
	merged := make(map[string]*dbstore.AttributeValue, len(item.Key)+2)
	for k, v := range item.Key {
		merged[k] = v
	}
	if indexName == "" {
		return merged
	}
	hashName, sortName, _ := indexKeyAttributeNames(table, indexName)
	for _, name := range []string{hashName, sortName} {
		if name == "" {
			continue
		}
		if _, exists := merged[name]; !exists {
			if attr, ok := item.Attributes[name]; ok {
				merged[name] = attr
			}
		}
	}
	return merged
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
