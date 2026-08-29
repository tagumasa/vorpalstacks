package dynamodb

import (
	"context"
	"math"
	"time"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// CreateBackup creates a backup of a DynamoDB table.
func (s *DynamoDBService) CreateBackup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	backup, err := s.createBackupCore(ctx, reqCtx, createBackupInput{
		Parameters: req.Parameters,
	})
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
	backup, err := s.deleteBackupCore(ctx, reqCtx, request.GetStringParam(req.Parameters, "BackupArn"))
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
	backup, err := s.describeBackupCore(ctx, reqCtx, request.GetStringParam(req.Parameters, "BackupArn"))
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
	coreResult, err := s.listBackupsCore(ctx, reqCtx, ListBackupsCoreInput{
		TableName:               request.GetStringParam(req.Parameters, "TableName"),
		BackupTypeFilter:        request.GetStringParam(req.Parameters, "BackupTypeFilter"),
		TimeRangeLowerBound:     request.GetInt64Param(req.Parameters, "TimeRangeLowerBoundDateTime"),
		TimeRangeUpperBound:     request.GetInt64Param(req.Parameters, "TimeRangeUpperBoundDateTime"),
		Limit:                   request.GetIntParam(req.Parameters, "Limit"),
		ExclusiveStartBackupArn: request.GetStringParam(req.Parameters, "ExclusiveStartBackupArn"),
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
	table, err := s.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
		BackupArn:       request.GetStringParam(req.Parameters, "BackupArn"),
		TargetTableName: request.GetStringParam(req.Parameters, "TargetTableName"),
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
	table, err := s.restoreTableToPointInTimeCore(ctx, reqCtx, restoreTableToPointInTimeInput{
		Parameters: req.Parameters,
	})
	if err != nil {
		return nil, err
	}

	// The restore summary travels with the table description (both the
	// restore response and later DescribeTable reads) from the persisted
	// table record.
	description := s.buildTableDescription(table)

	return map[string]interface{}{
		"TableDescription": description,
	}, nil
}

// parseTimestampParam reads an epoch-seconds timestamp parameter. AWS JSON
// protocols serialise timestamps as numbers, which arrive as float64 after
// body decoding.
func parseTimestampParam(params map[string]interface{}, key string) (time.Time, bool) {
	switch v := params[key].(type) {
	case float64:
		seconds, fraction := math.Modf(v)
		return time.Unix(int64(seconds), int64(fraction*1e9)), true
	case int64:
		return time.Unix(v, 0), true
	}
	return time.Time{}, false
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
			if nkaList, ok := proj["NonKeyAttributes"].([]interface{}); ok {
				for _, nkaRaw := range nkaList {
					if nka, ok := nkaRaw.(string); ok {
						idx.Projection.NonKeyAttributes = append(idx.Projection.NonKeyAttributes, nka)
					}
				}
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
			if nkaList, ok := proj["NonKeyAttributes"].([]interface{}); ok {
				for _, nkaRaw := range nkaList {
					if nka, ok := nkaRaw.(string); ok {
						idx.Projection.NonKeyAttributes = append(idx.Projection.NonKeyAttributes, nka)
					}
				}
			}
		}
		result = append(result, idx)
	}
	return result
}

// selectGSIOverrides narrows the restored table's global secondary indexes
// to those named by the override list, applying the provided projection and
// throughput settings. Restore overrides select from the existing indexes:
// an override naming an unknown index, or one replacing the key schema, is
// a validation error.
func selectGSIOverrides(existing []*dbstore.GlobalSecondaryIndex, overrides []*dbstore.GlobalSecondaryIndex) ([]*dbstore.GlobalSecondaryIndex, error) {
	byName := make(map[string]*dbstore.GlobalSecondaryIndex, len(existing))
	for _, g := range existing {
		byName[g.IndexName] = g
	}
	selected := make([]*dbstore.GlobalSecondaryIndex, 0, len(overrides))
	for _, ov := range overrides {
		base, ok := byName[ov.IndexName]
		if !ok {
			return nil, ErrInvalidParameter
		}
		if len(ov.KeySchema) > 0 && !keySchemasEqual(ov.KeySchema, base.KeySchema) {
			return nil, ErrInvalidParameter
		}
		restored := *base
		if ov.Projection != nil {
			restored.Projection = ov.Projection
		}
		if ov.ProvisionedThroughput != nil {
			restored.ProvisionedThroughput = ov.ProvisionedThroughput
		}
		selected = append(selected, &restored)
	}
	return selected, nil
}

// selectLSIOverrides narrows the restored table's local secondary indexes
// to those named by the override list, applying the provided projection.
// Local secondary index key schemas cannot change, so an override naming an
// unknown index or a different key schema is a validation error.
func selectLSIOverrides(existing []*dbstore.LocalSecondaryIndex, overrides []*dbstore.LocalSecondaryIndex) ([]*dbstore.LocalSecondaryIndex, error) {
	byName := make(map[string]*dbstore.LocalSecondaryIndex, len(existing))
	for _, l := range existing {
		byName[l.IndexName] = l
	}
	selected := make([]*dbstore.LocalSecondaryIndex, 0, len(overrides))
	for _, ov := range overrides {
		base, ok := byName[ov.IndexName]
		if !ok {
			return nil, ErrInvalidParameter
		}
		if len(ov.KeySchema) > 0 && !keySchemasEqual(ov.KeySchema, base.KeySchema) {
			return nil, ErrInvalidParameter
		}
		restored := *base
		if ov.Projection != nil {
			restored.Projection = ov.Projection
		}
		selected = append(selected, &restored)
	}
	return selected, nil
}

// keySchemasEqual compares two key schemas element by element.
func keySchemasEqual(a, b []*dbstore.KeySchemaElement) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].AttributeName != b[i].AttributeName || a[i].KeyType != b[i].KeyType {
			return false
		}
	}
	return true
}
