package dynamodb

import (
	"context"

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

// backfillVectorIndex scans all existing items in a table and creates vector
// index entries for the specified index only, mirroring backfillGSI for
// vector indexes added via UpdateTable.
func (s *DynamoDBService) backfillVectorIndex(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName, indexName string) {
	defer func() { resilience.RecoverPanic("dynamodb vector index backfill") }()
	err := store.Items().Scan(tableName, func(item *dbstore.Item) error {
		return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			return txn.PutVectorEntriesForIndex(tableName, indexName, item)
		})
	})
	if err != nil {
		logs.Warn("failed to backfill vector index entries",
			logs.String("table", tableName),
			logs.String("index", indexName),
			logs.Err(err))
	}
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

// backfillGSI scans all existing items in a table and creates index entries
// for the specified GSI only. This is called when a new GSI is added via
// UpdateTable, matching AWS behaviour where newly created GSIs are
// automatically populated with existing items.
func (s *DynamoDBService) backfillGSI(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName, gsiName string) {
	defer func() { resilience.RecoverPanic("dynamodb GSI backfill") }()
	err := store.Items().Scan(tableName, func(item *dbstore.Item) error {
		return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			return txn.PutGSIEntriesForIndex(tableName, gsiName, item)
		})
	})
	if err != nil {
		logs.Warn("failed to backfill GSI index entries",
			logs.String("table", tableName),
			logs.String("index", gsiName),
			logs.Err(err))
	}
}
