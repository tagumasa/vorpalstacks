// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// CreateTable creates a new DynamoDB table.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CreateTable.html
func (s *DynamoDBService) CreateTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableName := request.GetStringParam(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	keySchema := parseKeySchema(req.Parameters)
	attrDefs := parseAttributeDefinitions(req.Parameters)

	billingMode := dbstore.BillingMode(request.GetStringParam(req.Parameters, "BillingMode"))

	var provThroughput *dbstore.ProvisionedThroughput
	if billingMode == dbstore.BillingModeProvisioned {
		provThroughput = parseProvisionedThroughput(req.Parameters)
	}

	gsi, err := parseGlobalSecondaryIndexes(req.Parameters)
	if err != nil {
		return nil, err
	}
	lsi, err := parseLocalSecondaryIndexes(req.Parameters)
	if err != nil {
		return nil, err
	}

	streamSpec, err := parseStreamSpecification(req.Parameters)
	if err != nil {
		return nil, err
	}
	tagList := tagutil.ParseTags(req.Parameters, "Tags")
	deletionProtectionEnabled := request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled")

	var warmThroughput *dbstore.WarmThroughput
	if wtMap, ok := req.Parameters["WarmThroughput"].(map[string]interface{}); ok {
		warmThroughput = &dbstore.WarmThroughput{
			ReadUnitsPerSecond:  request.GetInt64Param(wtMap, "ReadUnitsPerSecond"),
			WriteUnitsPerSecond: request.GetInt64Param(wtMap, "WriteUnitsPerSecond"),
		}
	}

	var onDemandThroughput *dbstore.OnDemandThroughput
	if odtMap, ok := req.Parameters["OnDemandThroughput"].(map[string]interface{}); ok {
		onDemandThroughput = &dbstore.OnDemandThroughput{
			MaxReadRequestUnits:  request.GetInt64Param(odtMap, "MaxReadRequestUnits"),
			MaxWriteRequestUnits: request.GetInt64Param(odtMap, "MaxWriteRequestUnits"),
		}
	}

	globalTableSourceArn := request.GetStringParam(req.Parameters, "GlobalTableSourceArn")

	var sseDesc *dbstore.SSEDescription
	if sseSpec, ok := req.Parameters["SSESpecification"].(map[string]interface{}); ok {
		var err error
		sseDesc, err = parseSSESpecification(sseSpec)
		if err != nil {
			return nil, err
		}
	}

	tableClass := request.GetStringParam(req.Parameters, "TableClass")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	table, err := s.createTableCore(store, CreateTableInput{
		TableName:                 tableName,
		KeySchema:                 keySchema,
		AttributeDefinitions:      attrDefs,
		BillingMode:               billingMode,
		ProvisionedThroughput:     provThroughput,
		GlobalSecondaryIndexes:    gsi,
		LocalSecondaryIndexes:     lsi,
		StreamSpecification:       streamSpec,
		Tags:                      tagList,
		DeletionProtectionEnabled: deletionProtectionEnabled,
		WarmThroughput:            warmThroughput,
		OnDemandThroughput:        onDemandThroughput,
		GlobalTableSourceArn:      globalTableSourceArn,
		SSEDescription:            sseDesc,
		TableClass:                tableClass,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}

// DeleteTable removes a DynamoDB table and all its data.
func (s *DynamoDBService) DeleteTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableName := request.GetStringParam(req.Parameters, "TableName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	deletedTable, err := s.deleteTableCore(ctx, store, tableName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(deletedTable),
	}, nil
}

// DescribeTable returns detailed metadata for a DynamoDB table.
func (s *DynamoDBService) DescribeTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableName := request.GetStringParam(req.Parameters, "TableName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	table, err := s.describeTableCore(store, tableName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Table": s.buildTableDescription(table),
	}, nil
}

// ListTables returns a list of DynamoDB tables.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_ListTables.html
func (s *DynamoDBService) ListTables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit := listTablesMaxLimit
	if _, ok := req.Parameters["Limit"]; ok {
		v := request.GetIntParam(req.Parameters, "Limit")
		if err := validateListTablesLimit(v); err != nil {
			return nil, err
		}
		limit = v
	}
	marker := pagination.GetMarker(req.Parameters, "ExclusiveStartTableName")
	if estn := request.GetStringParam(req.Parameters, "ExclusiveStartTableName"); estn != "" {
		if err := validateResourceName(estn); err != nil {
			return nil, err
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tables, nextToken, err := s.listTablesCore(store, marker, limit)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(tables))
	for i, t := range tables {
		names[i] = t.Name
	}

	resp := map[string]interface{}{
		"TableNames": names,
	}
	if nextToken != "" {
		resp["LastEvaluatedTableName"] = nextToken
	}

	return resp, nil
}

// UpdateTable updates a DynamoDB table.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_UpdateTable.html
func (s *DynamoDBService) UpdateTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// TableName must pass the same length and character checks as CreateTable.
	if err := validateTableName(request.GetStringParam(req.Parameters, "TableName")); err != nil {
		return nil, err
	}

	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}

	table = deepCopyTable(table)

	if billingMode := request.GetStringParam(req.Parameters, "BillingMode"); billingMode != "" {
		table.BillingMode = dbstore.BillingMode(billingMode)
	}

	if provThroughput := parseProvisionedThroughput(req.Parameters); provThroughput != nil {
		table.ProvisionedThroughput = provThroughput
	}

	// PROVISIONED billing mode requires ProvisionedThroughput.
	if err := validateBillingModeConsistency(table.BillingMode, table.ProvisionedThroughput); err != nil {
		return nil, err
	}

	if attrDefs := parseAttributeDefinitions(req.Parameters); len(attrDefs) > 0 {
		// New attribute definitions must cover all key attributes.
		if err := validateAttributeDefinitions(table.KeySchema, attrDefs); err != nil {
			return nil, err
		}
		table.AttributeDefinitions = mergeAttributeDefinitions(table.AttributeDefinitions, attrDefs)
	}

	// Record existing GSI names before applying updates so we can
	// backfill newly created GSIs with existing items.
	existingGSINames := make(map[string]bool)
	for _, g := range table.GlobalSecondaryIndexes {
		existingGSINames[g.IndexName] = true
	}

	if gsiUpdates, ok := req.Parameters["GlobalSecondaryIndexUpdates"].([]interface{}); ok {
		updatedGSIs, err := applyGSIUpdates(table.ARN, table.GlobalSecondaryIndexes, gsiUpdates)
		if err != nil {
			return nil, err
		}
		table.GlobalSecondaryIndexes = updatedGSIs
	}

	// Validate that all key attributes across table + GSIs + LSIs are
	// present in AttributeDefinitions, and that index names are unique.
	// These checks run on the merged post-update table state.
	if err := validateAllKeyAttributesInDefs(table.KeySchema, table.GlobalSecondaryIndexes, table.LocalSecondaryIndexes, table.AttributeDefinitions); err != nil {
		return nil, err
	}
	if err := validateIndexNameUniqueness(table.GlobalSecondaryIndexes, table.LocalSecondaryIndexes); err != nil {
		return nil, err
	}

	streamSpec, err := parseStreamSpecification(req.Parameters)
	if err != nil {
		return nil, err
	}
	if streamSpec != nil {
		table.StreamSpecification = streamSpec
		if streamSpec.StreamEnabled {
			now := time.Now().UTC()
			table.StreamArn = table.ARN + "/stream/" + now.Format("2006-01-02T15:04:05.000")
			table.LatestStreamLabel = now.Format("2006-01-02T15:04:05.000")
		} else {
			table.StreamArn = ""
			table.LatestStreamLabel = ""
		}
	}

	if sseSpec, ok := req.Parameters["SSESpecification"].(map[string]interface{}); ok {
		sseDesc, err := parseSSESpecification(sseSpec)
		if err != nil {
			return nil, err
		}
		table.SSEDescription = sseDesc
	}

	if _, ok := req.Parameters["DeletionProtectionEnabled"]; ok {
		table.DeletionProtectionEnabled = request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled")
	}

	if tableClass := request.GetStringParam(req.Parameters, "TableClass"); tableClass != "" {
		table.TableClass = tableClass
	}

	table.LastUpdatedDateTime = time.Now().UTC()
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Tables().Put(table); err != nil {
		return nil, err
	}

	// Backfill: when new GSIs are created via UpdateTable, AWS automatically
	// populates index entries for all existing items. We must do the same
	// so that Query on the new GSI returns pre-existing items.
	for _, g := range table.GlobalSecondaryIndexes {
		if existingGSINames[g.IndexName] {
			continue
		}
		s.backfillGSI(ctx, store, table.Name, g.IndexName)
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}
