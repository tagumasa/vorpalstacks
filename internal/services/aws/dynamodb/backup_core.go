package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Backup Core — single validation + persistence path for backup operations.
//
// These methods encapsulate backup lifecycle logic. Both the HTTP API
// handlers (backup_operations.go) and any future admin handler delegate to
// these methods to ensure identical behaviour.
// ---------------------------------------------------------------------------

// createBackupCore persists a new backup record in CREATING state and
// launches the background snapshot goroutine. The table must already be
// validated and ACTIVE. Returns the persisted backup in its initial state.
func (s *DynamoDBService) createBackupCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, table *dbstore.Table, backupName string) (*dbstore.Backup, error) {
	backup, err := store.Backups().Create(backupName, table.Name, table.ARN, table.TableSizeBytes)
	if err != nil {
		return nil, err
	}

	backup.KeySchema = table.KeySchema
	backup.AttributeDefinitions = table.AttributeDefinitions
	backup.BillingMode = table.BillingMode
	backup.ProvisionedThroughput = table.ProvisionedThroughput
	backup.GlobalSecondaryIndexes = table.GlobalSecondaryIndexes
	backup.LocalSecondaryIndexes = table.LocalSecondaryIndexes
	backup.SourceTableCreationTime = table.CreationDateTime
	backup.SourceTableSizeBytes = table.TableSizeBytes
	backup.SourceTableItemCount = table.ItemCount
	backup.BackupSizeBytes = table.TableSizeBytes
	backup.BackupStatus = dbstore.BackupStatusCreating

	if err := store.Backups().Put(backup); err != nil {
		return nil, err
	}

	tableName := table.Name
	go func() {
		defer func() { resilience.RecoverPanic("dynamodb backup snapshot transition") }()
		var snapshotItems []*dbstore.Item
		if err := store.Items().Scan(tableName, func(item *dbstore.Item) error {
			snapshotItems = append(snapshotItems, &dbstore.Item{
				TableName:  tableName,
				Key:        copyAttributes(item.Key),
				Attributes: copyAttributes(item.Attributes),
			})
			return nil
		}); err != nil {
			logs.Error("Failed to scan items for backup",
				logs.Err(err),
				logs.String("tableName", tableName),
			)
			return
		}
		if err := store.Backups().SaveSnapshot(backupName, snapshotItems); err != nil {
			logs.Error("Failed to save backup snapshot",
				logs.Err(err),
				logs.String("backupName", backupName),
			)
			return
		}
		backup.BackupStatus = dbstore.BackupStatusAvailable
		if err := store.Backups().Put(backup); err != nil {
			logs.Error("Failed to update backup status to AVAILABLE",
				logs.Err(err),
				logs.String("backupName", backupName),
			)
		}
	}()

	return backup, nil
}

// deleteBackupCore deletes a backup by ARN and returns the deleted backup
// record for response formatting.
func (s *DynamoDBService) deleteBackupCore(store dbstore.DynamoDBStoreInterface, backupArn string) (*dbstore.Backup, error) {
	backup, err := store.Backups().Get(backupArn)
	if err != nil {
		return nil, ErrBackupNotFound
	}
	if err := store.Backups().Delete(backup.BackupName); err != nil {
		return nil, err
	}
	return backup, nil
}

// describeBackupCore returns a backup by ARN.
func (s *DynamoDBService) describeBackupCore(store dbstore.DynamoDBStoreInterface, backupArn string) (*dbstore.Backup, error) {
	backup, err := store.Backups().Get(backupArn)
	if err != nil {
		return nil, ErrBackupNotFound
	}
	return backup, nil
}

// ListBackupsCoreInput is the service-layer DTO for ListBackups.
type ListBackupsCoreInput struct {
	TableName               string
	BackupTypeFilter        string
	TimeRangeLowerBound     int64
	TimeRangeUpperBound     int64
	Limit                   int
	ExclusiveStartBackupArn string
}

// ListBackupsCoreResult is the service-layer result of ListBackups.
type ListBackupsCoreResult struct {
	Backups                []*dbstore.Backup
	LastEvaluatedBackupArn string
}

// listBackupsCore returns a filtered, paginated list of backups.
func (s *DynamoDBService) listBackupsCore(store dbstore.DynamoDBStoreInterface, in ListBackupsCoreInput) (*ListBackupsCoreResult, error) {
	marker := ""
	if in.ExclusiveStartBackupArn != "" {
		parts := strings.Split(in.ExclusiveStartBackupArn, "/")
		if len(parts) > 0 {
			marker = parts[len(parts)-1]
		}
	}

	fetchLimit := in.Limit
	if fetchLimit < listBackupsMinFetchSize {
		fetchLimit = listBackupsMinFetchSize
	}
	backups, _, err := store.Backups().List(marker, fetchLimit, in.TableName)
	if err != nil {
		return nil, err
	}

	var filtered []*dbstore.Backup
	for _, b := range backups {
		if in.BackupTypeFilter != "" && string(b.BackupType) != in.BackupTypeFilter {
			continue
		}
		backupTime := b.BackupCreationDateTime.Unix()
		if in.TimeRangeLowerBound > 0 && backupTime < in.TimeRangeLowerBound {
			continue
		}
		if in.TimeRangeUpperBound > 0 && backupTime > in.TimeRangeUpperBound {
			continue
		}
		filtered = append(filtered, b)
	}

	result := &ListBackupsCoreResult{}
	if len(filtered) > in.Limit {
		result.Backups = filtered[:in.Limit]
		result.LastEvaluatedBackupArn = filtered[in.Limit-1].BackupArn
	} else {
		result.Backups = filtered
	}
	return result, nil
}

// RestoreTableFromBackupCoreInput is the service-layer DTO for
// RestoreTableFromBackup.
type RestoreTableFromBackupCoreInput struct {
	BackupArn       string
	TargetTableName string
}

// restoreTableFromBackupCore creates a new table from a backup snapshot.
// The target table must not already exist.
func (s *DynamoDBService) restoreTableFromBackupCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, in RestoreTableFromBackupCoreInput) (*dbstore.Table, error) {
	backup, err := store.Backups().Get(in.BackupArn)
	if err != nil {
		return nil, ErrBackupNotFound
	}

	if store.Tables().Exists(in.TargetTableName) {
		return nil, ErrTableAlreadyExistsException
	}

	var keySchema []*dbstore.KeySchemaElement
	var attrDefs []*dbstore.AttributeDefinition
	var billingMode dbstore.BillingMode
	var provThroughput *dbstore.ProvisionedThroughput
	var gsi []*dbstore.GlobalSecondaryIndex
	var lsi []*dbstore.LocalSecondaryIndex

	if len(backup.KeySchema) > 0 {
		keySchema = backup.KeySchema
		attrDefs = backup.AttributeDefinitions
		billingMode = backup.BillingMode
		provThroughput = backup.ProvisionedThroughput
		gsi = backup.GlobalSecondaryIndexes
		lsi = backup.LocalSecondaryIndexes
	} else {
		sourceTable, err := store.Tables().Get(backup.SourceTableName)
		if err != nil {
			return nil, fmt.Errorf("backup %s has no key schema and source table %q not found: %w", in.BackupArn, backup.SourceTableName, err)
		}
		keySchema = sourceTable.KeySchema
		attrDefs = sourceTable.AttributeDefinitions
		billingMode = sourceTable.BillingMode
		provThroughput = sourceTable.ProvisionedThroughput
		gsi = sourceTable.GlobalSecondaryIndexes
		lsi = sourceTable.LocalSecondaryIndexes
	}

	if !validateBillingModeConsistency(billingMode, provThroughput) {
		return nil, ErrInvalidParameter
	}

	table, err := store.Tables().Create(
		in.TargetTableName, keySchema, attrDefs, billingMode, provThroughput,
		gsi, lsi, nil, nil, false,
	)
	if err != nil {
		return nil, err
	}

	// Set table to CREATING status so clients cannot see partial data.
	table.Status = dbstore.TableStatusCreating
	if err := store.Tables().Put(table); err != nil {
		return nil, err
	}

	snapshotItems, snapErr := store.Backups().GetSnapshot(backup.BackupName)
	if snapErr == nil && len(snapshotItems) > 0 {
		buffer := make([]*dbstore.Item, 0, restoreChunkSize)
		for _, item := range snapshotItems {
			buffer = append(buffer, item)
			if len(buffer) >= restoreChunkSize {
				if err := s.flushRestoreChunk(ctx, store, in.TargetTableName, buffer); err != nil {
					return nil, err
				}
				buffer = buffer[:0]
			}
		}
		if len(buffer) > 0 {
			if err := s.flushRestoreChunk(ctx, store, in.TargetTableName, buffer); err != nil {
				return nil, err
			}
		}
	}

	// All items copied — re-fetch to preserve ItemCount/TableSizeBytes
	// updated during chunked restore, then transition to ACTIVE.
	table, err = store.Tables().Get(in.TargetTableName)
	if err != nil {
		return nil, err
	}
	table.Status = dbstore.TableStatusActive
	if err := store.Tables().Put(table); err != nil {
		return nil, err
	}

	return table, nil
}

// RestoreTableToPointInTimeCoreInput is the service-layer DTO for
// RestoreTableToPointInTime.
type RestoreTableToPointInTimeCoreInput struct {
	SourceTable     *dbstore.Table
	TargetTableName string
	BillingMode     dbstore.BillingMode
	ProvThroughput  *dbstore.ProvisionedThroughput
	SSEDesc         *dbstore.SSEDescription
	GSI             []*dbstore.GlobalSecondaryIndex
	LSI             []*dbstore.LocalSecondaryIndex
}

// restoreTableToPointInTimeCore creates a new table by copying the source
// table's schema and items, applying any overrides. The target table must
// not already exist.
// restoreChunkSize is the number of items written per store.Update()
// transaction during restore. Each chunk is atomic; the table remains
// in CREATING status (invisible to clients) until all chunks complete.
const restoreChunkSize = 500

func (s *DynamoDBService) restoreTableToPointInTimeCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, in RestoreTableToPointInTimeCoreInput) (*dbstore.Table, error) {
	if store.Tables().Exists(in.TargetTableName) {
		return nil, ErrTableAlreadyExistsException
	}

	if !validateBillingModeConsistency(in.BillingMode, in.ProvThroughput) {
		return nil, ErrInvalidParameter
	}

	table, err := store.Tables().Create(
		in.TargetTableName, in.SourceTable.KeySchema, in.SourceTable.AttributeDefinitions,
		in.BillingMode, in.ProvThroughput, in.GSI, in.LSI, nil, nil, false,
	)
	if err != nil {
		return nil, err
	}

	// Set table to CREATING status so clients cannot see partial data
	// during the restore. AWS uses the same status transition.
	table.Status = dbstore.TableStatusCreating
	if in.SSEDesc != nil {
		table.SSEDescription = in.SSEDesc
	}
	if err := store.Tables().Put(table); err != nil {
		return nil, err
	}

	sourceTableName := in.SourceTable.Name

	// Scan source items and copy to target in buffered chunks. Each
	// chunk is a single atomic transaction; if any chunk fails, the
	// partially-populated CREATING table is left for startup recovery.
	buffer := make([]*dbstore.Item, 0, restoreChunkSize)

	err = store.Items().Scan(sourceTableName, func(item *dbstore.Item) error {
		buffer = append(buffer, &dbstore.Item{
			TableName:  in.TargetTableName,
			Key:        copyAttributes(item.Key),
			Attributes: copyAttributes(item.Attributes),
		})
		if len(buffer) >= restoreChunkSize {
			if err := s.flushRestoreChunk(ctx, store, in.TargetTableName, buffer); err != nil {
				return err
			}
			buffer = buffer[:0]
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Flush remaining items.
	if len(buffer) > 0 {
		if err := s.flushRestoreChunk(ctx, store, in.TargetTableName, buffer); err != nil {
			return nil, err
		}
	}

	// All items copied — re-fetch to preserve ItemCount/TableSizeBytes
	// updated during chunked restore, then transition to ACTIVE.
	table, err = store.Tables().Get(in.TargetTableName)
	if err != nil {
		return nil, err
	}
	table.Status = dbstore.TableStatusActive
	if err := store.Tables().Put(table); err != nil {
		return nil, err
	}

	return table, nil
}

// flushRestoreChunk writes a batch of items to the target table in a
// single atomic transaction, updating indexes, item count, and table
// size for each item.
func (s *DynamoDBService) flushRestoreChunk(ctx context.Context, store dbstore.DynamoDBStoreInterface, tableName string, items []*dbstore.Item) error {
	return store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
		for _, item := range items {
			itemSize := calculateItemSize(item.Attributes)
			if err := txn.PutItem(tableName, item.Key, item.Attributes); err != nil {
				return err
			}
			if err := txn.PutIndexEntries(tableName, item); err != nil {
				return err
			}
			if err := txn.UpdateItemCount(tableName, 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(tableName, itemSize); err != nil {
				return err
			}
		}
		return nil
	})
}
