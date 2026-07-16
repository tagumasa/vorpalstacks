package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	commonstore "vorpalstacks/internal/store/aws/common"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func (s *DynamoDBService) validateAndGetTable(reqCtx *request.RequestContext, params map[string]interface{}) (*dbstore.Table, error) {
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
			return nil, ErrTableNotFound
		}
		return nil, err
	}
	return table, nil
}

func (s *DynamoDBService) validateAndGetActiveTable(reqCtx *request.RequestContext, params map[string]interface{}) (*dbstore.Table, error) {
	table, err := s.validateAndGetTable(reqCtx, params)
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
