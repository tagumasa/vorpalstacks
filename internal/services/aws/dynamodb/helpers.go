package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func (s *DynamoDBService) validateAndGetTable(reqCtx *request.RequestContext, params map[string]interface{}) (*dbstore.Table, error) {
	return s.validateAndGetTableWithErr(reqCtx, params, ErrTableNotFound)
}

// storeRecordMissing reports whether an error is a store read failure of
// the not-found kind, in any sentinel family the store layers use: the
// common not-found class and the dynamodb store's own raw sentinels.
func storeRecordMissing(err error) bool {
	return commonstore.IsNotFound(err) ||
		errors.Is(err, dbstore.ErrTableNotFound) ||
		errors.Is(err, dbstore.ErrBackupNotFound)
}

// describeByArn runs the shared shape of every ARN-addressed describe
// core (backup, export, import): a length-validated ARN addresses one
// record of the regional store, and a failed read maps to the family's
// not-found sentinel rather than the store's raw error.
func describeByArn[T any](s *DynamoDBService, reqCtx *request.RequestContext, arn string, arnValid bool, notFound error, get func(dbstore.DynamoDBStoreInterface) (T, error)) (T, error) {
	var zero T
	if !arnValid {
		return zero, ErrInvalidParameter
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return zero, err
	}
	got, err := get(store)
	if err != nil {
		// Absence maps to the family's not-found sentinel; a storage fault
		// (bucket I/O, an unreadable record) is reported as the storage
		// error it is — never masquerading as absence.
		if storeRecordMissing(err) {
			return zero, notFound
		}
		return zero, err
	}
	return got, nil
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
// never mutated, and re-deriving index ARNs against the restored table's
// region and account through the ARN builder.
func applyRestoredVectorIndexes(arnBuilder *svcarn.DynamoDBBuilder, table *dbstore.Table, vectorIdx []*dbstore.VectorIndex) {
	if len(vectorIdx) == 0 {
		return
	}
	copied := make([]*dbstore.VectorIndex, len(vectorIdx))
	for i, vi := range vectorIdx {
		clone := *vi
		clone.IndexArn = arnBuilder.Index(table.Name, vi.IndexName)
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
