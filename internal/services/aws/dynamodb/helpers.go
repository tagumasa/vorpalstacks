package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func (s *DynamoDBService) validateAndGetTable(reqCtx *request.RequestContext, params map[string]interface{}) (*dbstore.Table, error) {
	return s.validateAndGetTableWithErr(reqCtx, params, ErrTableNotFound)
}

// validateAndGetTableWithErr behaves like validateAndGetTable but lets the
// caller pick the not-found error sentinel. Operations whose Smithy model
// declares TableNotFoundException (rather than the general
// ResourceNotFoundException) must pass ErrTableNotFoundException here so the
// client receives the individual error code.
func (s *DynamoDBService) validateAndGetTableWithErr(reqCtx *request.RequestContext, params map[string]interface{}, notFoundErr *APIError) (*dbstore.Table, error) {
	tableName := request.GetStringParam(params, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := store.Tables().Get(tableName)
	if err != nil {
		if dbstore.IsTableNotFound(err) || commonstore.IsNotFound(err) {
			return nil, notFoundErr
		}
		return nil, err
	}
	return table, nil
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
// for the specified GSI. This is called when a new GSI is added via
// UpdateTable, matching AWS behaviour where newly created GSIs are
// automatically populated with existing items.
func (s *DynamoDBService) backfillGSI(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName, gsiName string) {
	defer func() { resilience.RecoverPanic("dynamodb GSI backfill") }()
	err := store.Items().Scan(tableName, func(item *dbstore.Item) error {
		return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			return txn.PutIndexEntries(tableName, item)
		})
	})
	if err != nil {
		logs.Warn("failed to backfill GSI index entries",
			logs.String("table", tableName),
			logs.String("index", gsiName),
			logs.Err(err))
	}
}
