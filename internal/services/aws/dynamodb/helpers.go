package dynamodb

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func (s *DynamoDBService) validateAndGetTable(reqCtx *request.RequestContext, params map[string]interface{}) (*dbstore.Table, error) {
	return s.validateAndGetTableWithErr(reqCtx, params, ErrTableNotFound)
}

func (s *DynamoDBService) validateAndGetActiveTable(reqCtx *request.RequestContext, params map[string]interface{}) (*dbstore.Table, error) {
	return s.validateAndGetActiveTableWithErr(reqCtx, params, ErrTableNotFound)
}

// validateAndGetActiveTableWithErr is the not-found-configurable variant of
// validateAndGetActiveTable. See validateAndGetTableWithErr for the rationale.
func (s *DynamoDBService) validateAndGetActiveTableWithErr(reqCtx *request.RequestContext, params map[string]interface{}, notFoundErr *APIError) (*dbstore.Table, error) {
	table, err := s.validateAndGetTableWithErr(reqCtx, params, notFoundErr)
	if err != nil {
		return nil, err
	}

	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}
	return table, nil
}

// backfillIndexEntries scans a table's existing items and stages the named
// index's entry for each through the per-item function the index family
// supplies. The helper runs on the request's goroutine: a panic converts
// to the returned error so the operation fails over the half-built index
// instead of succeeding silently, while a scan-level failure stays a
// logged warning — the table record carrying the index has already
// committed, and a re-run does not retry the population: UpdateTable
// skips an index name already present before the request, so recovering
// an unpopulated index means removing and re-adding it. Re-running the
// population itself would be safe — both families' entries are keyed
// upserts — the skip keeps a re-declared create from redoing the scan.
func (s *DynamoDBService) backfillIndexEntries(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName, indexName, family string, stageEntry func(txn *dbstore.DynamoDBTxn, item *dbstore.Item) error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			resilience.LogPanic(fmt.Sprintf("dynamodb %s backfill", family), r)
			err = fmt.Errorf("%s backfill panicked: %v", family, r)
		}
	}()
	if scanErr := store.Items().Scan(tableName, func(item *dbstore.Item) error {
		return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			return stageEntry(txn, item)
		})
	}); scanErr != nil {
		logs.Warn("failed to backfill index entries",
			logs.String("table", tableName),
			logs.String("index", indexName),
			logs.String("family", family),
			logs.Err(scanErr))
	}
	return nil
}

// backfillVectorIndex populates a vector index's entries for existing
// items, added via UpdateTable — the vector family's staging of the shared
// backfill lifecycle.
func (s *DynamoDBService) backfillVectorIndex(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName, indexName string) error {
	return s.backfillIndexEntries(ctx, store, tableName, indexName, "vector index", func(txn *dbstore.DynamoDBTxn, item *dbstore.Item) error {
		return txn.PutVectorEntriesForIndex(tableName, indexName, item)
	})
}

// applyRestoredVectorIndexes attaches vector index metadata to a restored
// table, cloning each definition so the backup or source table records are
// never mutated, and re-deriving index ARNs against the restored table's ARN.
func applyRestoredVectorIndexes(table *dbstore.Table, vectorIdx []*dbstore.VectorIndex) {
	if len(vectorIdx) == 0 {
		return
	}
	copied := make([]*dbstore.VectorIndex, len(vectorIdx))
	for i, vi := range vectorIdx {
		clone := *vi
		clone.IndexArn = table.ARN + "/index/" + vi.IndexName
		copied[i] = &clone
	}
	table.VectorIndexes = copied
}

func validateIndexExists(table *dbstore.Table, indexName string) bool {
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			return true
		}
	}
	for _, lsi := range table.LocalSecondaryIndexes {
		if lsi.IndexName == indexName {
			return true
		}
	}
	return false
}

// backfillGSI populates a GSI's entries for existing items, added via
// UpdateTable, matching AWS behaviour where newly created GSIs are
// automatically populated — the GSI family's staging of the shared
// backfill lifecycle.
func (s *DynamoDBService) backfillGSI(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName, gsiName string) error {
	return s.backfillIndexEntries(ctx, store, tableName, gsiName, "GSI", func(txn *dbstore.DynamoDBTxn, item *dbstore.Item) error {
		return txn.PutGSIEntriesForIndex(tableName, gsiName, item)
	})
}
