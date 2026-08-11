// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CreateBackup creates a backup of a DynamoDB table.
func (s *DynamoDBService) CreateBackup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Table must be ACTIVE to create a backup.
	table, err := s.validateAndGetActiveTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	backupName := request.GetStringParam(req.Parameters, "BackupName")
	if err := validateResourceName(backupName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if store.Backups().Exists(backupName) {
		return nil, ErrBackupAlreadyExists
	}

	backup, err := store.Backups().Create(backupName, tableName, table.ARN, table.TableSizeBytes)
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

	// Set CREATING status and persist immediately so DescribeBackup
	// shows the correct state while the snapshot is being taken.
	backup.BackupStatus = dbstore.BackupStatusCreating
	if err := store.Backups().Put(backup); err != nil {
		return nil, err
	}

	// Take the snapshot in the background and transition to AVAILABLE.
	go func() {
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

	return map[string]interface{}{
		"BackupDetails": map[string]interface{}{
			"BackupArn":              backup.BackupArn,
			"BackupName":             backup.BackupName,
			"BackupSizeBytes":        backup.BackupSizeBytes,
			"BackupStatus":           string(backup.BackupStatus),
			"BackupType":             string(backup.BackupType),
			"BackupCreationDateTime": backup.BackupCreationDateTime.Unix(),
			"SourceTableName":        backup.SourceTableName,
			"SourceTableArn":         backup.SourceTableArn,
		},
	}, nil
}

// DeleteBackup deletes a DynamoDB table backup.
func (s *DynamoDBService) DeleteBackup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	backupArn := request.GetStringParam(req.Parameters, "BackupArn")
	if err := validateBackupArn(backupArn); err != nil {
		return nil, err
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

	return map[string]interface{}{
		"BackupDescription": map[string]interface{}{
			"BackupDetails": map[string]interface{}{
				"BackupArn":              backup.BackupArn,
				"BackupName":             backup.BackupName,
				"BackupStatus":           "DELETED",
				"BackupType":             string(backup.BackupType),
				"BackupCreationDateTime": backup.BackupCreationDateTime.Unix(),
			},
			"SourceTableDetails": map[string]interface{}{
				"TableName": backup.SourceTableName,
			},
		},
	}, nil
}

// DescribeBackup returns information about a DynamoDB table backup.
func (s *DynamoDBService) DescribeBackup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	backupArn := request.GetStringParam(req.Parameters, "BackupArn")
	if err := validateBackupArn(backupArn); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	backup, err := store.Backups().Get(backupArn)
	if err != nil {
		return nil, ErrBackupNotFound
	}

	return map[string]interface{}{
		"BackupDescription": map[string]interface{}{
			"BackupDetails": map[string]interface{}{
				"BackupArn":              backup.BackupArn,
				"BackupName":             backup.BackupName,
				"BackupSizeBytes":        backup.BackupSizeBytes,
				"BackupStatus":           string(backup.BackupStatus),
				"BackupType":             string(backup.BackupType),
				"BackupCreationDateTime": backup.BackupCreationDateTime.Unix(),
			},
			"SourceTableDetails": map[string]interface{}{
				"TableName":             backup.SourceTableName,
				"TableArn":              backup.SourceTableArn,
				"TableSizeBytes":        backup.SourceTableSizeBytes,
				"TableCreationDateTime": backup.SourceTableCreationTime.Unix(),
				"ItemCount":             backup.SourceTableItemCount,
				"KeySchema":             buildKeySchemaResponse(backup.KeySchema),
				"ProvisionedThroughput": func() map[string]interface{} {
					if backup.ProvisionedThroughput != nil {
						return map[string]interface{}{
							"ReadCapacityUnits":      backup.ProvisionedThroughput.ReadCapacityUnits,
							"WriteCapacityUnits":     backup.ProvisionedThroughput.WriteCapacityUnits,
							"NumberOfDecreasesToday": backup.ProvisionedThroughput.NumberOfDecreasesToday,
						}
					}
					return map[string]interface{}{
						"ReadCapacityUnits":      int64(0),
						"WriteCapacityUnits":     int64(0),
						"NumberOfDecreasesToday": int64(0),
					}
				}(),
			},
		},
	}, nil
}

// ListBackups lists the backups of a DynamoDB table.
func (s *DynamoDBService) ListBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableName := request.GetStringParam(req.Parameters, "TableName")
	backupTypeFilter := request.GetStringParam(req.Parameters, "BackupTypeFilter")
	timeRangeLowerBound := request.GetInt64Param(req.Parameters, "TimeRangeLowerBoundDateTime")
	timeRangeUpperBound := request.GetInt64Param(req.Parameters, "TimeRangeUpperBoundDateTime")
	limit := request.GetIntParam(req.Parameters, "Limit")
	if limit == 0 {
		limit = listBackupsMaxLimit
	} else {
		if err := validateListBackupsLimit(limit); err != nil {
			return nil, err
		}
	}
	exclusiveStartBackupArn := request.GetStringParam(req.Parameters, "ExclusiveStartBackupArn")
	if exclusiveStartBackupArn != "" {
		if err := validateBackupArn(exclusiveStartBackupArn); err != nil {
			return nil, err
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Convert ExclusiveStartBackupArn to backup name for the store marker
	// (store keys are backup names, not ARNs).
	marker := ""
	if exclusiveStartBackupArn != "" {
		parts := strings.Split(exclusiveStartBackupArn, "/")
		if len(parts) > 0 {
			marker = parts[len(parts)-1]
		}
	}

	// Fetch with a larger window than limit so post-filter pagination works.
	fetchLimit := limit
	if fetchLimit < 100 {
		fetchLimit = 100
	}
	backups, _, err := store.Backups().List(marker, fetchLimit, tableName)
	if err != nil {
		return nil, err
	}

	var filteredBackups []*dbstore.Backup
	for _, b := range backups {
		if backupTypeFilter != "" && string(b.BackupType) != backupTypeFilter {
			continue
		}

		backupTime := b.BackupCreationDateTime.Unix()
		if timeRangeLowerBound > 0 && backupTime < timeRangeLowerBound {
			continue
		}
		if timeRangeUpperBound > 0 && backupTime > timeRangeUpperBound {
			continue
		}

		filteredBackups = append(filteredBackups, b)
	}

	backupSummaries := make([]map[string]interface{}, 0)
	hasMore := len(filteredBackups) > limit

	if len(filteredBackups) > limit {
		filteredBackups = filteredBackups[:limit]
	}

	for _, b := range filteredBackups {
		backupSummaries = append(backupSummaries, map[string]interface{}{
			"BackupArn":              b.BackupArn,
			"BackupName":             b.BackupName,
			"BackupSizeBytes":        b.BackupSizeBytes,
			"BackupStatus":           string(b.BackupStatus),
			"BackupType":             string(b.BackupType),
			"BackupCreationDateTime": b.BackupCreationDateTime.Unix(),
			"TableName":              b.SourceTableName,
			"TableArn":               b.SourceTableArn,
		})
	}

	resp := map[string]interface{}{
		"BackupSummaries": backupSummaries,
	}

	if hasMore && len(backupSummaries) > 0 {
		resp["LastEvaluatedBackupArn"] = backupSummaries[len(backupSummaries)-1]["BackupArn"]
	}

	return resp, nil
}

// RestoreTableFromBackup restores a table from a DynamoDB backup.
func (s *DynamoDBService) RestoreTableFromBackup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	backupArn := request.GetStringParam(req.Parameters, "BackupArn")
	if err := validateBackupArn(backupArn); err != nil {
		return nil, err
	}

	targetTableName := request.GetStringParam(req.Parameters, "TargetTableName")
	if targetTableName == "" {
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

	if store.Tables().Exists(targetTableName) {
		return nil, ErrTableAlreadyExists
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
			return nil, fmt.Errorf("backup %s has no key schema and source table %q not found: %w", backupArn, backup.SourceTableName, err)
		}
		keySchema = sourceTable.KeySchema
		attrDefs = sourceTable.AttributeDefinitions
		billingMode = sourceTable.BillingMode
		provThroughput = sourceTable.ProvisionedThroughput
		gsi = sourceTable.GlobalSecondaryIndexes
		lsi = sourceTable.LocalSecondaryIndexes
	}

	table, err := store.Tables().Create(
		targetTableName,
		keySchema,
		attrDefs,
		billingMode,
		provThroughput,
		gsi,
		lsi,
		nil,
		nil,
		false,
	)
	if err != nil {
		return nil, err
	}

	// Restore items from the backup snapshot. This works even if the
	// original source table has been deleted.
	snapshotItems, snapErr := store.Backups().GetSnapshot(backup.BackupName)
	if snapErr == nil && len(snapshotItems) > 0 {
		for _, item := range snapshotItems {
			itemSize := calculateItemSize(item.Attributes)
			err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
				if err := txn.PutItem(targetTableName, item.Key, item.Attributes); err != nil {
					return err
				}
				if err := txn.PutIndexEntries(targetTableName, item); err != nil {
					return err
				}
				if err := txn.UpdateItemCount(targetTableName, 1); err != nil {
					return err
				}
				if err := txn.UpdateTableSize(targetTableName, itemSize); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}

// RestoreTableToPointInTime restores a table to a point in time.
func (s *DynamoDBService) RestoreTableToPointInTime(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	sourceTableName := request.GetStringParam(req.Parameters, "SourceTableName")
	sourceTableArn := request.GetStringParam(req.Parameters, "SourceTableArn")

	if sourceTableName == "" && sourceTableArn == "" {
		return nil, ErrInvalidParameter
	}
	if sourceTableName != "" && sourceTableArn != "" {
		return nil, ErrInvalidParameter
	}

	if sourceTableArn != "" {
		sourceTableName = svcarn.ParseTableARN(sourceTableArn)
		if sourceTableName == "" {
			return nil, ErrResourceNotFound
		}
	}

	targetTableName := request.GetStringParam(req.Parameters, "TargetTableName")
	if targetTableName == "" {
		return nil, ErrInvalidParameter
	}

	sourceTable, err := s.validateAndGetTable(reqCtx, map[string]interface{}{"TableName": sourceTableName})
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if store.Tables().Exists(targetTableName) {
		return nil, ErrTableAlreadyExists
	}

	keySchema := sourceTable.KeySchema
	attrDefs := sourceTable.AttributeDefinitions
	billingMode := sourceTable.BillingMode
	provThroughput := sourceTable.ProvisionedThroughput
	gsi := sourceTable.GlobalSecondaryIndexes
	lsi := sourceTable.LocalSecondaryIndexes
	var sseDesc *dbstore.SSEDescription

	if billingModeOverride := request.GetStringParam(req.Parameters, "BillingModeOverride"); billingModeOverride != "" {
		billingMode = dbstore.BillingMode(billingModeOverride)
	}

	if provOverride := parseProvisionedThroughputOverride(req.Parameters); provOverride != nil {
		provThroughput = provOverride
	}

	if sseSpec, ok := req.Parameters["SSESpecificationOverride"].(map[string]interface{}); ok {
		sseDesc, err = parseSSESpecification(sseSpec)
		if err != nil {
			return nil, err
		}
	}

	if gsiOverrides := parseGSIOverrideList(req.Parameters); len(gsiOverrides) > 0 {
		gsi = applyGSIOverrides(gsi, gsiOverrides)
	}

	if lsiOverrideList := parseLSIOverrideList(req.Parameters); len(lsiOverrideList) > 0 {
		lsi = lsiOverrideList
	}

	if err := validateBillingModeConsistency(billingMode, provThroughput); err != nil {
		return nil, err
	}

	table, err := store.Tables().Create(
		targetTableName,
		keySchema,
		attrDefs,
		billingMode,
		provThroughput,
		gsi,
		lsi,
		nil,
		nil,
		false,
	)
	if err != nil {
		return nil, err
	}
	if sseDesc != nil {
		table.SSEDescription = sseDesc
	}

	var itemsToCopy []*dbstore.Item
	err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
		return store.Items().Scan(sourceTableName, func(item *dbstore.Item) error {
			copiedItem := &dbstore.Item{
				TableName:  targetTableName,
				Key:        copyAttributes(item.Key),
				Attributes: copyAttributes(item.Attributes),
			}
			itemsToCopy = append(itemsToCopy, copiedItem)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	for _, item := range itemsToCopy {
		itemSize := calculateItemSize(item.Attributes)
		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			if err := txn.PutItem(targetTableName, item.Key, item.Attributes); err != nil {
				return err
			}
			if err := txn.PutIndexEntries(targetTableName, item); err != nil {
				return err
			}
			if err := txn.UpdateItemCount(targetTableName, 1); err != nil {
				return err
			}
			if err := txn.UpdateTableSize(targetTableName, itemSize); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}

func parseProvisionedThroughputOverride(params map[string]interface{}) *dbstore.ProvisionedThroughput {
	ptMap, ok := params["ProvisionedThroughputOverride"].(map[string]interface{})
	if !ok {
		return nil
	}
	rcu := int64(0)
	if v, ok := ptMap["ReadCapacityUnits"]; ok {
		if f, ok := v.(float64); ok {
			rcu = int64(f)
		}
	}
	wcu := int64(0)
	if v, ok := ptMap["WriteCapacityUnits"]; ok {
		if f, ok := v.(float64); ok {
			wcu = int64(f)
		}
	}
	if rcu == 0 && wcu == 0 {
		return nil
	}
	return &dbstore.ProvisionedThroughput{
		ReadCapacityUnits:  rcu,
		WriteCapacityUnits: wcu,
	}
}

func parseGSIOverrideList(params map[string]interface{}) []*dbstore.GlobalSecondaryIndex {
	rawList, ok := params["GlobalSecondaryIndexOverride"].([]interface{})
	if !ok {
		return nil
	}
	var result []*dbstore.GlobalSecondaryIndex
	for _, raw := range rawList {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		idx := &dbstore.GlobalSecondaryIndex{}
		if name, ok := m["IndexName"].(string); ok {
			idx.IndexName = name
		}
		if ksList, ok := m["KeySchema"].([]interface{}); ok {
			for _, ksRaw := range ksList {
				ksMap, ok := ksRaw.(map[string]interface{})
				if !ok {
					continue
				}
				ks := &dbstore.KeySchemaElement{}
				if n, ok := ksMap["AttributeName"].(string); ok {
					ks.AttributeName = n
				}
				if t, ok := ksMap["KeyType"].(string); ok {
					ks.KeyType = dbstore.KeyType(t)
				}
				idx.KeySchema = append(idx.KeySchema, ks)
			}
		}
		if proj, ok := m["Projection"].(map[string]interface{}); ok {
			idx.Projection = &dbstore.Projection{}
			if pt, ok := proj["ProjectionType"].(string); ok {
				idx.Projection.ProjectionType = pt
			}
		}
		if pt, ok := m["ProvisionedThroughput"].(map[string]interface{}); ok {
			idx.ProvisionedThroughput = &dbstore.ProvisionedThroughput{}
			if v, ok := pt["ReadCapacityUnits"].(float64); ok {
				idx.ProvisionedThroughput.ReadCapacityUnits = int64(v)
			}
			if v, ok := pt["WriteCapacityUnits"].(float64); ok {
				idx.ProvisionedThroughput.WriteCapacityUnits = int64(v)
			}
		}
		result = append(result, idx)
	}
	return result
}

func parseLSIOverrideList(params map[string]interface{}) []*dbstore.LocalSecondaryIndex {
	rawList, ok := params["LocalSecondaryIndexOverride"].([]interface{})
	if !ok {
		return nil
	}
	var result []*dbstore.LocalSecondaryIndex
	for _, raw := range rawList {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		idx := &dbstore.LocalSecondaryIndex{}
		if name, ok := m["IndexName"].(string); ok {
			idx.IndexName = name
		}
		if ksList, ok := m["KeySchema"].([]interface{}); ok {
			for _, ksRaw := range ksList {
				ksMap, ok := ksRaw.(map[string]interface{})
				if !ok {
					continue
				}
				ks := &dbstore.KeySchemaElement{}
				if n, ok := ksMap["AttributeName"].(string); ok {
					ks.AttributeName = n
				}
				if t, ok := ksMap["KeyType"].(string); ok {
					ks.KeyType = dbstore.KeyType(t)
				}
				idx.KeySchema = append(idx.KeySchema, ks)
			}
		}
		if proj, ok := m["Projection"].(map[string]interface{}); ok {
			idx.Projection = &dbstore.Projection{}
			if pt, ok := proj["ProjectionType"].(string); ok {
				idx.Projection.ProjectionType = pt
			}
		}
		result = append(result, idx)
	}
	return result
}

func applyGSIOverrides(existing []*dbstore.GlobalSecondaryIndex, overrides []*dbstore.GlobalSecondaryIndex) []*dbstore.GlobalSecondaryIndex {
	byName := make(map[string]*dbstore.GlobalSecondaryIndex)
	for _, g := range existing {
		byName[g.IndexName] = g
	}
	for _, ov := range overrides {
		if existingGSI, ok := byName[ov.IndexName]; ok {
			if len(ov.KeySchema) > 0 {
				existingGSI.KeySchema = ov.KeySchema
			}
			if ov.Projection != nil {
				existingGSI.Projection = ov.Projection
			}
			if ov.ProvisionedThroughput != nil {
				existingGSI.ProvisionedThroughput = ov.ProvisionedThroughput
			}
		} else {
			byName[ov.IndexName] = ov
		}
	}
	result := make([]*dbstore.GlobalSecondaryIndex, 0, len(byName))
	for _, g := range byName {
		result = append(result, g)
	}
	return result
}
