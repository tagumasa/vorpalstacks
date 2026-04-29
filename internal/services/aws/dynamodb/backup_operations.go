// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// CreateBackup creates a backup of a DynamoDB table.
func (s *DynamoDBService) CreateBackup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	backupName := request.GetStringParam(req.Parameters, "BackupName")
	if backupName == "" {
		return nil, ErrInvalidParameter
	}

	if len(backupName) < 3 || len(backupName) > 255 {
		return nil, ErrInvalidParameter
	}

	if !tableNameRegex.MatchString(backupName) {
		return nil, ErrInvalidParameter
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
	if err := store.Backups().Put(backup); err != nil {
		return nil, err
	}

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
	if backupArn == "" {
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
	if backupArn == "" {
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
		limit = 100
	}
	exclusiveStartBackupArn := request.GetStringParam(req.Parameters, "ExclusiveStartBackupArn")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	backups, _, err := store.Backups().List("", 0, tableName)
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

	startIdx := 0
	if exclusiveStartBackupArn != "" {
		for i, b := range filteredBackups {
			if b.BackupArn == exclusiveStartBackupArn {
				startIdx = i + 1
				break
			}
		}
	}

	filteredBackups = filteredBackups[startIdx:]

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
	if backupArn == "" {
		return nil, ErrInvalidParameter
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
		if err == nil {
			keySchema = sourceTable.KeySchema
			attrDefs = sourceTable.AttributeDefinitions
			billingMode = sourceTable.BillingMode
			provThroughput = sourceTable.ProvisionedThroughput
			gsi = sourceTable.GlobalSecondaryIndexes
			lsi = sourceTable.LocalSecondaryIndexes
		} else {
			keySchema = []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}}
			attrDefs = []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}}
			billingMode = dbstore.BillingModePayPerRequest
		}
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

	sourceTable, err := store.Tables().Get(backup.SourceTableName)
	if err != nil && !dbstore.IsTableNotFound(err) {
		return nil, fmt.Errorf("failed to get source table %s: %w", backup.SourceTableName, err)
	}
	if sourceTable != nil {
		var itemsToCopy []*dbstore.Item
		err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
			return store.Items().Scan(backup.SourceTableName, func(item *dbstore.Item) error {
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
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}

// RestoreTableToPointInTime restores a table to a point in time.
func (s *DynamoDBService) RestoreTableToPointInTime(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	sourceTableName := request.GetStringParam(req.Parameters, "SourceTableName")
	if sourceTableName == "" {
		return nil, ErrInvalidParameter
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

	table, err := store.Tables().Create(
		targetTableName,
		sourceTable.KeySchema,
		sourceTable.AttributeDefinitions,
		sourceTable.BillingMode,
		sourceTable.ProvisionedThroughput,
		sourceTable.GlobalSecondaryIndexes,
		sourceTable.LocalSecondaryIndexes,
		nil,
		nil,
		false,
	)
	if err != nil {
		return nil, err
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
