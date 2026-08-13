package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CreateBackup creates a backup of a DynamoDB table.
func (s *DynamoDBService) CreateBackup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Table must be ACTIVE to create a backup. CreateBackup's Smithy model
	// declares TableNotFoundException (not ResourceNotFoundException), so
	// surface the individual error sentinel here.
	table, err := s.validateAndGetActiveTableWithErr(reqCtx, req.Parameters, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}

	backupName := request.GetStringParam(req.Parameters, "BackupName")
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

	backup, err := s.createBackupCore(ctx, store, table, backupName)
	if err != nil {
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
	if !validateBackupArn(backupArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	backup, err := s.deleteBackupCore(store, backupArn)
	if err != nil {
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
	if !validateBackupArn(backupArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	backup, err := s.describeBackupCore(store, backupArn)
	if err != nil {
		return nil, err
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
		if !validateListBackupsLimit(limit) {
			return nil, ErrInvalidParameter
		}
	}
	exclusiveStartBackupArn := request.GetStringParam(req.Parameters, "ExclusiveStartBackupArn")
	if exclusiveStartBackupArn != "" {
		if !validateBackupArn(exclusiveStartBackupArn) {
			return nil, ErrInvalidParameter
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	coreResult, err := s.listBackupsCore(store, ListBackupsCoreInput{
		TableName:               tableName,
		BackupTypeFilter:        backupTypeFilter,
		TimeRangeLowerBound:     timeRangeLowerBound,
		TimeRangeUpperBound:     timeRangeUpperBound,
		Limit:                   limit,
		ExclusiveStartBackupArn: exclusiveStartBackupArn,
	})
	if err != nil {
		return nil, err
	}

	backupSummaries := make([]map[string]interface{}, 0, len(coreResult.Backups))
	for _, b := range coreResult.Backups {
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
	if coreResult.LastEvaluatedBackupArn != "" {
		resp["LastEvaluatedBackupArn"] = coreResult.LastEvaluatedBackupArn
	}

	return resp, nil
}

// RestoreTableFromBackup restores a table from a DynamoDB backup.
func (s *DynamoDBService) RestoreTableFromBackup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	backupArn := request.GetStringParam(req.Parameters, "BackupArn")
	if !validateBackupArn(backupArn) {
		return nil, ErrInvalidParameter
	}

	targetTableName := request.GetStringParam(req.Parameters, "TargetTableName")
	if !validateResourceName(targetTableName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	table, err := s.restoreTableFromBackupCore(ctx, store, RestoreTableFromBackupCoreInput{
		BackupArn:       backupArn,
		TargetTableName: targetTableName,
	})
	if err != nil {
		return nil, err
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

	targetTableName := request.GetStringParam(req.Parameters, "TargetTableName")
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

	// Parse optional overrides.
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

	table, err := s.restoreTableToPointInTimeCore(ctx, store, RestoreTableToPointInTimeCoreInput{
		SourceTable:     sourceTable,
		TargetTableName: targetTableName,
		BillingMode:     billingMode,
		ProvThroughput:  provThroughput,
		SSEDesc:         sseDesc,
		GSI:             gsi,
		LSI:             lsi,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}

// parseProvisionedThroughputOverride extracts the ProvisionedThroughputOverride
// parameter from a RestoreTableToPointInTime request.
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

// parseGSIOverrideList extracts the GlobalSecondaryIndexOverride parameter
// from a RestoreTableToPointInTime request.
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

// parseLSIOverrideList extracts the LocalSecondaryIndexOverride parameter
// from a RestoreTableToPointInTime request.
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

// applyGSIOverrides merges override GSI definitions into the existing set,
// matching by IndexName.
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
