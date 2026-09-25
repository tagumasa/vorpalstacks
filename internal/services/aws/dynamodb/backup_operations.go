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
		"BackupDetails": buildBackupDetailsResponse(backup, backup.BackupStatus),
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
			"BackupDetails":      buildBackupDetailsResponse(backup, dbstore.BackupStatusDeleted),
			"SourceTableDetails": buildSourceTableDetailsResponse(backup),
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
			"BackupDetails":      buildBackupDetailsResponse(backup, backup.BackupStatus),
			"SourceTableDetails": buildSourceTableDetailsResponse(backup),
		},
	}, nil
}

// ListBackups lists the backups of a DynamoDB table.
func (s *DynamoDBService) ListBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// The model members are TimeRangeLowerBound/TimeRangeUpperBound, and
	// presence decides: an explicit epoch bound filters, an omitted member
	// does not — the zero epoch is a real bound, not an omission.
	lowerBound, lowerSet, lowerErr := intParamWithPresence(req.Parameters, "TimeRangeLowerBound")
	if lowerErr != nil {
		return nil, lowerErr
	}
	upperBound, upperSet, upperErr := intParamWithPresence(req.Parameters, "TimeRangeUpperBound")
	if upperErr != nil {
		return nil, upperErr
	}
	in := ListBackupsCoreInput{
		TableName:               request.GetStringParam(req.Parameters, "TableName"),
		BackupTypeFilter:        request.GetStringParam(req.Parameters, "BackupTypeFilter"),
		Limit:                   request.GetIntParam(req.Parameters, "Limit"),
		ExclusiveStartBackupArn: request.GetStringParam(req.Parameters, "ExclusiveStartBackupArn"),
	}
	if lowerSet {
		lb := int64(lowerBound)
		in.TimeRangeLowerBound = &lb
	}
	if upperSet {
		ub := int64(upperBound)
		in.TimeRangeUpperBound = &ub
	}
	coreResult, err := s.listBackupsCore(ctx, reqCtx, in)
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
			"TableId":                b.SourceTableId,
			"BillingModeSummary": map[string]interface{}{
				"BillingMode": string(b.BillingMode),
			},
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
	overrides, err := parseRestoreOverridesWire(req.Parameters)
	if err != nil {
		return nil, err
	}
	table, err := s.restoreTableFromBackupCore(ctx, reqCtx, RestoreTableFromBackupCoreInput{
		BackupArn:       request.GetStringParam(req.Parameters, "BackupArn"),
		TargetTableName: request.GetStringParam(req.Parameters, "TargetTableName"),
		Overrides:       overrides,
	})
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table, s.replicasForTable(store, table.Name)),
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// The restore summary travels with the table description (both the
	// restore response and later DescribeTable reads) from the persisted
	// table record.
	description := s.buildTableDescription(table, s.replicasForTable(store, table.Name))

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

// indexOverrideFields holds the override members shared by the global and
// local secondary index families: both carry index name, key schema and
// projection on the wire, and every member is optional — an override that
// omits one leaves the restored index's value untouched.
type indexOverrideFields struct {
	IndexName  string
	KeySchema  []*dbstore.KeySchemaElement
	Projection *dbstore.Projection
}

// parseIndexOverrideEntry extracts the shared members of one override
// entry.
func parseIndexOverrideEntry(m map[string]interface{}) indexOverrideFields {
	var f indexOverrideFields
	if name, ok := m["IndexName"].(string); ok {
		f.IndexName = name
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
			f.KeySchema = append(f.KeySchema, ks)
		}
	}
	if proj, ok := m["Projection"].(map[string]interface{}); ok {
		f.Projection = &dbstore.Projection{}
		if pt, ok := proj["ProjectionType"].(string); ok {
			f.Projection.ProjectionType = dbstore.ProjectionType(pt)
		}
		if nkaList, ok := proj["NonKeyAttributes"].([]interface{}); ok {
			for _, nkaRaw := range nkaList {
				if nka, ok := nkaRaw.(string); ok {
					f.Projection.NonKeyAttributes = append(f.Projection.NonKeyAttributes, nka)
				}
			}
		}
	}
	return f
}

// parseGSIOverrideList extracts the GlobalSecondaryIndexOverride parameter
// from a RestoreTableToPointInTime request.
func parseGSIOverrideList(params map[string]interface{}) ([]*dbstore.GlobalSecondaryIndex, error) {
	rawList, ok := params["GlobalSecondaryIndexOverride"].([]interface{})
	if !ok {
		return nil, nil
	}
	result := make([]*dbstore.GlobalSecondaryIndex, 0, len(rawList))
	for _, raw := range rawList {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		f := parseIndexOverrideEntry(m)
		idx := &dbstore.GlobalSecondaryIndex{
			IndexName:  f.IndexName,
			KeySchema:  f.KeySchema,
			Projection: f.Projection,
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
		if odt, ok := m["OnDemandThroughput"].(map[string]interface{}); ok {
			parsed, odtErr := parseOnDemandThroughputMap(odt)
			if odtErr != nil {
				return nil, odtErr
			}
			idx.OnDemandThroughput = parsed
		}
		if wt, ok := m["WarmThroughput"].(map[string]interface{}); ok {
			idx.WarmThroughput = parseWarmThroughputMap(wt)
		}
		result = append(result, idx)
	}
	return result, nil
}

// parseLSIOverrideList extracts the LocalSecondaryIndexOverride parameter
// from a RestoreTableToPointInTime request.
func parseLSIOverrideList(params map[string]interface{}) []*dbstore.LocalSecondaryIndex {
	rawList, ok := params["LocalSecondaryIndexOverride"].([]interface{})
	if !ok {
		return nil
	}
	result := make([]*dbstore.LocalSecondaryIndex, 0, len(rawList))
	for _, raw := range rawList {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		f := parseIndexOverrideEntry(m)
		result = append(result, &dbstore.LocalSecondaryIndex{
			IndexName:  f.IndexName,
			KeySchema:  f.KeySchema,
			Projection: f.Projection,
		})
	}
	return result
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

// parseOnDemandThroughputOverride extracts the OnDemandThroughputOverride
// member shared by both restore inputs; an absent member returns nil.
func parseOnDemandThroughputOverride(params map[string]interface{}) (*dbstore.OnDemandThroughput, error) {
	odt, ok := params["OnDemandThroughputOverride"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	return parseOnDemandThroughputMap(odt)
}

// parseVectorIndexOverrideList extracts the VectorIndexOverride member
// shared by both restore inputs; an absent member returns nil.
func parseVectorIndexOverrideList(params map[string]interface{}) ([]*dbstore.VectorIndex, error) {
	rawList, ok := params["VectorIndexOverride"].([]interface{})
	if !ok {
		return nil, nil
	}
	return parseVectorIndexList(rawList)
}

// parseRestoreOverridesWire reads the override member family the model
// defines identically on RestoreTableFromBackupInput and
// RestoreTableToPointInTimeInput. Every member is optional; an absent
// member leaves the zero value, which the core reads as "keep the base
// setting".
func parseRestoreOverridesWire(params map[string]interface{}) (RestoreOverrides, error) {
	var overrides RestoreOverrides
	if bm := request.GetStringParam(params, "BillingModeOverride"); bm != "" {
		overrides.BillingMode = dbstore.BillingMode(bm)
	}
	overrides.ProvisionedThroughput = parseProvisionedThroughputOverride(params)
	odtOverride, odtErr := parseOnDemandThroughputOverride(params)
	if odtErr != nil {
		return overrides, odtErr
	}
	overrides.OnDemandThroughput = odtOverride
	desc, disable, err := parseSSESpecification(params["SSESpecificationOverride"])
	if err != nil {
		return overrides, err
	}
	overrides.SSEDescription = desc
	overrides.SSEDisable = disable
	gsiOverrides, gsiErr := parseGSIOverrideList(params)
	if gsiErr != nil {
		return overrides, gsiErr
	}
	overrides.GlobalSecondaryIndexes = gsiOverrides
	overrides.LocalSecondaryIndexes = parseLSIOverrideList(params)
	vectorIdx, err := parseVectorIndexOverrideList(params)
	if err != nil {
		return overrides, err
	}
	overrides.VectorIndexes = vectorIdx
	return overrides, nil
}
