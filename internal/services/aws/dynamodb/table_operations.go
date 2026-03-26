// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/pagination"
	"vorpalstacks/internal/services/aws/common/request"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// CreateTable creates a new DynamoDB table.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_CreateTable.html
func (s *DynamoDBService) CreateTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableName := request.GetStringParam(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	if len(tableName) < 3 || len(tableName) > 255 {
		return nil, ErrInvalidParameter
	}

	if !tableNameRegex.MatchString(tableName) {
		return nil, ErrInvalidParameter
	}

	keySchema := parseKeySchema(req.Parameters)
	if len(keySchema) == 0 {
		return nil, ErrInvalidParameter
	}

	if err := validateKeySchema(keySchema); err != nil {
		return nil, err
	}

	attrDefs := parseAttributeDefinitions(req.Parameters)
	if err := validateAttributeDefinitions(keySchema, attrDefs); err != nil {
		return nil, err
	}

	billingMode := dbstore.BillingMode(request.GetStringParam(req.Parameters, "BillingMode"))
	if billingMode == "" {
		billingMode = dbstore.BillingModeProvisioned
	}

	var provThroughput *dbstore.ProvisionedThroughput
	if billingMode == dbstore.BillingModeProvisioned {
		provThroughput = parseProvisionedThroughput(req.Parameters)
		if provThroughput == nil {
			return nil, ErrInvalidParameter
		}
	}

	gsi := parseGlobalSecondaryIndexes(req.Parameters)
	lsi := parseLocalSecondaryIndexes(req.Parameters)

	if err := validateAllKeyAttributesInDefs(keySchema, gsi, lsi, attrDefs); err != nil {
		return nil, err
	}

	if err := validateLSIPartitionKey(keySchema, lsi); err != nil {
		return nil, err
	}

	if err := validateIndexNameUniqueness(gsi, lsi); err != nil {
		return nil, err
	}

	streamSpec := parseStreamSpecification(req.Parameters)
	tagList := tagutil.ParseTags(req.Parameters, "Tags")
	deletionProtectionEnabled := request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := store.Tables().Create(
		tableName,
		keySchema,
		attrDefs,
		billingMode,
		provThroughput,
		gsi,
		lsi,
		streamSpec,
		tagList,
		deletionProtectionEnabled,
	)
	if err != nil {
		if dbstore.IsTableAlreadyExists(err) {
			return nil, ErrTableAlreadyExists
		}
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}

// DeleteTable deletes a DynamoDB table.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_DeleteTable.html
func (s *DynamoDBService) DeleteTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	if table.DeletionProtectionEnabled {
		return nil, ErrTableDeletionProtected
	}

	tableName := table.Name

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Items().DeleteAllForTable(tableName); err != nil {
		return nil, err
	}

	if err := store.Tables().Delete(tableName); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}

// DescribeTable returns information about a DynamoDB table.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_DescribeTable.html
func (s *DynamoDBService) DescribeTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
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
	limit := pagination.GetMaxItems(req.Parameters, 100)
	marker := pagination.GetMarker(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tables, nextToken, err := store.Tables().List(marker, limit)
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
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	// Check if table is in ACTIVE state
	if table.Status != dbstore.TableStatusActive {
		return nil, ErrTableNotActive
	}

	if billingMode := request.GetStringParam(req.Parameters, "BillingMode"); billingMode != "" {
		table.BillingMode = dbstore.BillingMode(billingMode)
	}

	if provThroughput := parseProvisionedThroughput(req.Parameters); provThroughput != nil {
		table.ProvisionedThroughput = provThroughput
	}

	if attrDefs := parseAttributeDefinitions(req.Parameters); len(attrDefs) > 0 {
		table.AttributeDefinitions = mergeAttributeDefinitions(table.AttributeDefinitions, attrDefs)
	}

	if gsiUpdates, ok := req.Parameters["GlobalSecondaryIndexUpdates"].([]interface{}); ok {
		table.GlobalSecondaryIndexes = applyGSIUpdates(table.ARN, table.GlobalSecondaryIndexes, gsiUpdates)
	}

	if streamSpec := parseStreamSpecification(req.Parameters); streamSpec != nil {
		table.StreamSpecification = streamSpec
		if streamSpec.StreamEnabled {
			now := time.Now().UTC()
			table.StreamArn = table.ARN + "/stream/" + now.Format("20060102150405")
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

	table.LastUpdatedDateTime = time.Now().UTC()
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Tables().Put(table); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table),
	}, nil
}
