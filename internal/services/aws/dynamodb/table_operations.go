// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

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
		if !validateListTablesLimit(v) {
			return nil, ErrInvalidParameter
		}
		limit = v
	}
	marker := pagination.GetMarker(req.Parameters, "ExclusiveStartTableName")
	if estn := request.GetStringParam(req.Parameters, "ExclusiveStartTableName"); estn != "" {
		if !validateResourceName(estn) {
			return nil, ErrInvalidParameter
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
// UpdateTable updates a DynamoDB table.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_UpdateTable.html
func (s *DynamoDBService) UpdateTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := UpdateTableInput{
		TableName:             request.GetStringParam(req.Parameters, "TableName"),
		BillingMode:           request.GetStringParam(req.Parameters, "BillingMode"),
		ProvisionedThroughput: parseProvisionedThroughput(req.Parameters),
		AttributeDefinitions:  parseAttributeDefinitions(req.Parameters),
		TableClass:            request.GetStringParam(req.Parameters, "TableClass"),
	}
	if gsiUpdates, ok := req.Parameters["GlobalSecondaryIndexUpdates"].([]interface{}); ok {
		in.GSIUpdates = gsiUpdates
	}
	streamSpec, streamErr := parseStreamSpecification(req.Parameters)
	if streamErr != nil {
		return nil, streamErr
	}
	in.StreamSpecification = streamSpec
	if sseSpec, ok := req.Parameters["SSESpecification"].(map[string]interface{}); ok {
		sseDesc, sseErr := parseSSESpecification(sseSpec)
		if sseErr != nil {
			return nil, sseErr
		}
		in.SSESpecification = sseDesc
	}
	if _, ok := req.Parameters["DeletionProtectionEnabled"]; ok {
		in.DeletionProtectionSet = true
		in.DeletionProtection = request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled")
	}

	table, err := s.updateTableCore(ctx, store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}
