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

	keySchema := parseKeySchema(req.Parameters)
	attrDefs := parseAttributeDefinitions(req.Parameters)

	billingMode := dbstore.BillingMode(request.GetStringParam(req.Parameters, "BillingMode"))
	provThroughput := parseProvisionedThroughput(req.Parameters)

	gsi, err := parseGlobalSecondaryIndexes(req.Parameters)
	if err != nil {
		return nil, err
	}
	lsi, err := parseLocalSecondaryIndexes(req.Parameters)
	if err != nil {
		return nil, err
	}
	vectorIdx, err := parseVectorIndexes(req.Parameters)
	if err != nil {
		return nil, err
	}

	streamSpec, err := parseStreamSpecification(req.Parameters)
	if err != nil {
		return nil, err
	}
	tagList := tagutil.ParseTags(req.Parameters, "Tags")

	warmThroughput := parseWarmThroughput(req.Parameters)
	onDemandThroughput, odtErr := parseOnDemandThroughput(req.Parameters)
	if odtErr != nil {
		return nil, odtErr
	}

	globalTableSourceArn := request.GetStringParam(req.Parameters, "GlobalTableSourceArn")

	var sseDesc *dbstore.SSEDescription
	var sseDisable bool
	desc, disable, err := parseSSESpecification(req.Parameters["SSESpecification"])
	if err != nil {
		return nil, err
	}
	sseDesc, sseDisable = desc, disable

	tableClass := request.GetStringParam(req.Parameters, "TableClass")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourcePolicy := request.GetStringParam(req.Parameters, "ResourcePolicy")
	_, resourcePolicySet := req.Parameters["ResourcePolicy"]

	table, err := s.createTableCore(ctx, reqCtx, store, CreateTableInput{
		TableName:              tableName,
		KeySchema:              keySchema,
		AttributeDefinitions:   attrDefs,
		BillingMode:            billingMode,
		ProvisionedThroughput:  provThroughput,
		GlobalSecondaryIndexes: gsi,
		LocalSecondaryIndexes:  lsi,
		VectorIndexes:          vectorIdx,
		StreamSpecification:    streamSpec,
		Tags:                   tagList,
		DeletionProtectionRaw:  req.Parameters["DeletionProtectionEnabled"],
		WarmThroughput:         warmThroughput,
		OnDemandThroughput:     onDemandThroughput,
		GlobalTableSourceArn:   globalTableSourceArn,
		SSEDescription:         sseDesc,
		SSEDisable:             sseDisable,
		TableClass:             tableClass,
		ResourcePolicySet:      resourcePolicySet,
		ResourcePolicy:         resourcePolicy,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table, nil),
	}, nil
}

// DeleteTable removes a DynamoDB table and all its data.
func (s *DynamoDBService) DeleteTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tableName := request.GetStringParam(req.Parameters, "TableName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	deletedTable, err := s.deleteTableCore(ctx, store, reqCtx.GetRegion(), tableName)
	if err != nil {
		return nil, err
	}

	replicas := s.replicasForTable(store, deletedTable.Name)
	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(deletedTable, replicas),
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
		"Table": s.buildTableDescription(table, s.replicasForTable(store, table.Name)),
	}, nil
}

// ListTables returns a list of DynamoDB tables.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_ListTables.html
func (s *DynamoDBService) ListTables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit := listTablesMaxLimit
	if _, ok := req.Parameters["Limit"]; ok {
		limit = request.GetIntParam(req.Parameters, "Limit")
	}
	marker := pagination.GetMarker(req.Parameters, "ExclusiveStartTableName")

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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	onDemandThroughput, odtErr := parseOnDemandThroughput(req.Parameters)
	if odtErr != nil {
		return nil, odtErr
	}
	in := UpdateTableInput{
		TableName:             request.GetStringParam(req.Parameters, "TableName"),
		BillingMode:           request.GetStringParam(req.Parameters, "BillingMode"),
		ProvisionedThroughput: parseProvisionedThroughput(req.Parameters),
		AttributeDefinitions:  parseAttributeDefinitions(req.Parameters),
		TableClass:            request.GetStringParam(req.Parameters, "TableClass"),
		OnDemandThroughput:    onDemandThroughput,
		WarmThroughput:        parseWarmThroughput(req.Parameters),
	}
	if gsiUpdates, ok := req.Parameters["GlobalSecondaryIndexUpdates"].([]interface{}); ok {
		in.GSIUpdates = gsiUpdates
	}
	if viUpdates, ok := req.Parameters["VectorIndexUpdates"].([]interface{}); ok {
		in.VectorIndexUpdates = viUpdates
	}
	if replicaUpdates, ok := req.Parameters["ReplicaUpdates"].([]interface{}); ok {
		in.ReplicaUpdates = replicaUpdates
	}
	in.RequestRegion = reqCtx.GetRegion()
	in.MultiRegionConsistency = request.GetStringParam(req.Parameters, "MultiRegionConsistency")
	if witnessUpdates, ok := req.Parameters["GlobalTableWitnessUpdates"].([]interface{}); ok {
		in.GlobalTableWitnessUpdates = witnessUpdates
	}
	in.GlobalTableSettingsReplicationMode = request.GetStringParam(req.Parameters, "GlobalTableSettingsReplicationMode")
	streamSpec, streamErr := parseStreamSpecification(req.Parameters)
	if streamErr != nil {
		return nil, streamErr
	}
	in.StreamSpecification = streamSpec
	sseDesc, disable, sseErr := parseSSESpecification(req.Parameters["SSESpecification"])
	if sseErr != nil {
		return nil, sseErr
	}
	in.SSESpecification = sseDesc
	in.SSEDisable = disable
	if raw, ok := req.Parameters["DeletionProtectionEnabled"]; ok {
		in.DeletionProtectionSet = true
		in.DeletionProtectionRaw = raw
	}

	table, err := s.updateTableCore(ctx, reqCtx, store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableDescription": s.buildTableDescription(table, s.replicasForTable(store, table.Name)),
	}, nil
}
