package dynamodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Backup Core — single validation + persistence path for backup operations.
//
// These methods encapsulate backup lifecycle logic. Both the HTTP API
// handlers (backup_operations.go) and any future admin handler delegate to
// these methods to ensure identical behaviour.
// ---------------------------------------------------------------------------

// createBackupInput carries the raw wire parameters for CreateBackup.
type createBackupInput struct {
	Parameters map[string]interface{}
}

// createBackupCore validates the request, persists a new backup record in
// CREATING state and launches the background snapshot goroutine. The table
// must be ACTIVE. Returns the persisted backup in its initial state.
func (s *DynamoDBService) createBackupCore(ctx context.Context, reqCtx *request.RequestContext, in createBackupInput) (*dbstore.Backup, error) {
	// Table must be ACTIVE to create a backup. CreateBackup's Smithy model
	// declares TableNotFoundException (not ResourceNotFoundException), so
	// surface the individual error sentinel here.
	table, err := s.validateAndGetActiveTableWithErr(reqCtx, in.Parameters, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}

	backupName := request.GetStringParam(in.Parameters, "BackupName")
	if !validateResourceName(backupName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if store.Backups().Exists(backupName) {
		return nil, ErrBackupAlreadyExists
	}

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

// deleteBackupCore validates the request, then deletes a backup by ARN and
// returns the deleted backup record for response formatting.
func (s *DynamoDBService) deleteBackupCore(ctx context.Context, reqCtx *request.RequestContext, backupArn string) (*dbstore.Backup, error) {
	if !validateBackupArn(backupArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	backup, err := store.Backups().Get(backupArn)
	if err != nil {
		return nil, ErrBackupNotFound
	}
	if err := store.Backups().Delete(backup.BackupName); err != nil {
		return nil, err
	}
	return backup, nil
}

// describeBackupCore validates the request, then returns a backup by ARN.
func (s *DynamoDBService) describeBackupCore(ctx context.Context, reqCtx *request.RequestContext, backupArn string) (*dbstore.Backup, error) {
	if !validateBackupArn(backupArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

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

// listBackupsCore validates the request, then returns a filtered, paginated
// list of backups.
func (s *DynamoDBService) listBackupsCore(ctx context.Context, reqCtx *request.RequestContext, in ListBackupsCoreInput) (*ListBackupsCoreResult, error) {
	if in.Limit == 0 {
		in.Limit = listBackupsMaxLimit
	} else {
		if !validateListBackupsLimit(in.Limit) {
			return nil, ErrInvalidParameter
		}
	}
	if in.ExclusiveStartBackupArn != "" {
		if !validateBackupArn(in.ExclusiveStartBackupArn) {
			return nil, ErrInvalidParameter
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

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

// restoreTableFromBackupCore validates the request, then creates a new
// table from a backup snapshot. The target table must not already exist.
func (s *DynamoDBService) restoreTableFromBackupCore(ctx context.Context, reqCtx *request.RequestContext, in RestoreTableFromBackupCoreInput) (*dbstore.Table, error) {
	if !validateBackupArn(in.BackupArn) {
		return nil, ErrInvalidParameter
	}
	if !validateResourceName(in.TargetTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

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
	table.RestoreSummary = &dbstore.RestoreSummary{
		SourceBackupArn:   backup.BackupArn,
		SourceTableArn:    backup.SourceTableArn,
		RestoreDateTime:   backup.BackupCreationDateTime,
		RestoreInProgress: false,
	}
	if err := store.Tables().Put(table); err != nil {
		return nil, err
	}

	return table, nil
}

// restoreTableToPointInTimeInput carries the raw wire parameters for
// RestoreTableToPointInTime.
type restoreTableToPointInTimeInput struct {
	Parameters map[string]interface{}
}

// restoreTableToPointInTimeCore validates the request, then creates a new
// table by copying the source table's schema and items, applying any
// overrides. The target table must not already exist.
// restoreChunkSize is the number of items written per store.Update()
// transaction during restore. Each chunk is atomic; the table remains
// in CREATING status (invisible to clients) until all chunks complete.
const restoreChunkSize = 500

func (s *DynamoDBService) restoreTableToPointInTimeCore(ctx context.Context, reqCtx *request.RequestContext, in restoreTableToPointInTimeInput) (*dbstore.Table, error) {
	sourceTableName := request.GetStringParam(in.Parameters, "SourceTableName")
	sourceTableArn := request.GetStringParam(in.Parameters, "SourceTableArn")

	if sourceTableName == "" && sourceTableArn == "" {
		return nil, ErrInvalidParameter
	}
	if sourceTableName != "" && sourceTableArn != "" {
		return nil, ErrInvalidParameter
	}

	if sourceTableName != "" {
		if !validateResourceName(sourceTableName) {
			return nil, ErrInvalidParameter
		}
	}

	if sourceTableArn != "" {
		sourceTableName = svcarn.ParseTableARN(sourceTableArn)
		if sourceTableName == "" {
			return nil, ErrResourceNotFound
		}
	}

	targetTableName := request.GetStringParam(in.Parameters, "TargetTableName")
	if !validateResourceName(targetTableName) {
		return nil, ErrInvalidParameter
	}

	// RestoreTableToPointInTime declares TableNotFoundException (rather
	// than the general ResourceNotFoundException) in the Smithy model, so
	// use the individual error sentinel here.
	sourceTable, err := s.validateAndGetTableWithErr(reqCtx, map[string]interface{}{"TableName": sourceTableName}, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Point-in-time restore requires recovery to be enabled on the source
	// table; the request must then name a point inside the restorable
	// window.
	pitr, err := store.Tables().GetPointInTimeRecovery(sourceTableName)
	if err != nil {
		return nil, err
	}
	if pitr == nil || pitr.Status != dbstore.PITRStatusEnabled {
		return nil, ErrPITRNotEnabled
	}

	now := time.Now()
	restoreDateTime, hasRestoreDateTime := parseTimestampParam(in.Parameters, "RestoreDateTime")
	useLatest := false
	if raw, present := in.Parameters["UseLatestRestorableTime"]; present {
		if b, isBool := raw.(bool); isBool {
			useLatest = b
		}
	}
	switch {
	case useLatest:
		restoreDateTime = now
	case !hasRestoreDateTime:
		return nil, ErrInvalidParameter
	}
	if restoreDateTime.Before(pitr.EarliestRestorableDateTime) || restoreDateTime.After(now) {
		return nil, ErrInvalidRestoreTime
	}

	// Parse optional overrides.
	billingMode := sourceTable.BillingMode
	provThroughput := sourceTable.ProvisionedThroughput
	gsi := sourceTable.GlobalSecondaryIndexes
	lsi := sourceTable.LocalSecondaryIndexes
	var sseDesc *dbstore.SSEDescription

	if billingModeOverride := request.GetStringParam(in.Parameters, "BillingModeOverride"); billingModeOverride != "" {
		billingMode = dbstore.BillingMode(billingModeOverride)
	}
	if provOverride := parseProvisionedThroughputOverride(in.Parameters); provOverride != nil {
		provThroughput = provOverride
	}
	if sseSpec, ok := in.Parameters["SSESpecificationOverride"].(map[string]interface{}); ok {
		sseDesc, err = parseSSESpecification(sseSpec)
		if err != nil {
			return nil, err
		}
	}
	if gsiOverrides := parseGSIOverrideList(in.Parameters); len(gsiOverrides) > 0 {
		gsi, err = selectGSIOverrides(gsi, gsiOverrides)
		if err != nil {
			return nil, err
		}
	}
	if lsiOverrideList := parseLSIOverrideList(in.Parameters); len(lsiOverrideList) > 0 {
		lsi, err = selectLSIOverrides(lsi, lsiOverrideList)
		if err != nil {
			return nil, err
		}
	}

	if store.Tables().Exists(targetTableName) {
		return nil, ErrTableAlreadyExistsException
	}

	if !validateBillingModeConsistency(billingMode, provThroughput) {
		return nil, ErrInvalidParameter
	}

	table, err := store.Tables().Create(
		targetTableName, sourceTable.KeySchema, sourceTable.AttributeDefinitions,
		billingMode, provThroughput, gsi, lsi, nil, nil, false,
	)
	if err != nil {
		return nil, err
	}

	// Set table to CREATING status so clients cannot see partial data
	// during the restore. AWS uses the same status transition.
	table.Status = dbstore.TableStatusCreating
	if sseDesc != nil {
		table.SSEDescription = sseDesc
	}
	if err := store.Tables().Put(table); err != nil {
		return nil, err
	}

	sourceTableName = sourceTable.Name

	// Snapshot the source as of the restore point (the current state with
	// every journaled mutation newer than the restore time undone) and
	// write it to the target in buffered chunks. Each chunk is a single
	// atomic transaction; if any chunk fails, the partially-populated
	// CREATING table is left for startup recovery.
	snapshot, err := snapshotItemsAsOf(store, sourceTableName, restoreDateTime)
	if err != nil {
		return nil, err
	}
	buffer := make([]*dbstore.Item, 0, restoreChunkSize)
	for _, item := range snapshot {
		buffer = append(buffer, &dbstore.Item{
			TableName:  targetTableName,
			Key:        copyAttributes(item.Key),
			Attributes: copyAttributes(item.Attributes),
		})
		if len(buffer) >= restoreChunkSize {
			if err := s.flushRestoreChunk(ctx, store, targetTableName, buffer); err != nil {
				return nil, err
			}
			buffer = buffer[:0]
		}
	}
	// Flush remaining items.
	if len(buffer) > 0 {
		if err := s.flushRestoreChunk(ctx, store, targetTableName, buffer); err != nil {
			return nil, err
		}
	}

	// All items copied — re-fetch to preserve ItemCount/TableSizeBytes
	// updated during chunked restore, then transition to ACTIVE.
	table, err = store.Tables().Get(targetTableName)
	if err != nil {
		return nil, err
	}
	table.Status = dbstore.TableStatusActive
	table.RestoreSummary = &dbstore.RestoreSummary{
		SourceTableArn:    sourceTable.ARN,
		RestoreDateTime:   restoreDateTime,
		RestoreInProgress: false,
	}
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
