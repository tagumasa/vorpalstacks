package appsync

import (
	"context"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/common/request"
)

// CreateDataSource creates a new data source for a GraphQL API.
// POST /v1/apis/{apiId}/datasources
func (s *AppSyncService) CreateDataSource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	dsType := request.GetStringParam(req.Parameters, "type")
	if dsType == "" {
		return nil, NewBadRequestException("type is required")
	}

	ds := &appsyncstore.DataSource{
		ApiId:                    apiId,
		Name:                     name,
		Type:                     dsType,
		Description:              request.GetStringParam(req.Parameters, "description"),
		ServiceRoleArn:           request.GetStringParam(req.Parameters, "serviceRoleArn"),
		DynamodbConfig:           parseDynamoDBConfig(req.Parameters),
		ElasticsearchConfig:      parseElasticsearchConfig(req.Parameters),
		EventBridgeConfig:        parseEventBridgeConfig(req.Parameters),
		HttpConfig:               parseHttpConfig(req.Parameters),
		LambdaConfig:             parseLambdaDataSourceConfig(req.Parameters),
		MetricsConfig:            request.GetStringParam(req.Parameters, "metricsConfig"),
		OpenSearchServiceConfig:  parseOpenSearchServiceConfig(req.Parameters),
		RelationalDatabaseConfig: parseRelationalDatabaseConfig(req.Parameters),
	}

	created, err := store.CreateDataSource(ds)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"dataSource": dataSourceToMap(created),
	}, nil
}

// GetDataSource retrieves a data source by API ID and name.
// GET /v1/apis/{apiId}/datasources/{name}
func (s *AppSyncService) GetDataSource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	name := request.GetStringParam(req.Parameters, "name")
	if apiId == "" || name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}

	ds, err := store.GetDataSource(apiId, name)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"dataSource": dataSourceToMap(ds),
	}, nil
}

// UpdateDataSource updates an existing data source.
// POST /v1/apis/{apiId}/datasources/{name}
func (s *AppSyncService) UpdateDataSource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	name := request.GetStringParam(req.Parameters, "name")
	if apiId == "" || name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}

	dsType := request.GetStringParam(req.Parameters, "type")
	if dsType == "" {
		return nil, NewBadRequestException("type is required")
	}

	ds := &appsyncstore.DataSource{
		ApiId:                    apiId,
		Name:                     name,
		Type:                     dsType,
		Description:              request.GetStringParam(req.Parameters, "description"),
		ServiceRoleArn:           request.GetStringParam(req.Parameters, "serviceRoleArn"),
		DynamodbConfig:           parseDynamoDBConfig(req.Parameters),
		ElasticsearchConfig:      parseElasticsearchConfig(req.Parameters),
		EventBridgeConfig:        parseEventBridgeConfig(req.Parameters),
		HttpConfig:               parseHttpConfig(req.Parameters),
		LambdaConfig:             parseLambdaDataSourceConfig(req.Parameters),
		MetricsConfig:            request.GetStringParam(req.Parameters, "metricsConfig"),
		OpenSearchServiceConfig:  parseOpenSearchServiceConfig(req.Parameters),
		RelationalDatabaseConfig: parseRelationalDatabaseConfig(req.Parameters),
	}

	updated, err := store.UpdateDataSource(ds)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{
		"dataSource": dataSourceToMap(updated),
	}, nil
}

// DeleteDataSource removes a data source.
// DELETE /v1/apis/{apiId}/datasources/{name}
func (s *AppSyncService) DeleteDataSource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	name := request.GetStringParam(req.Parameters, "name")
	if apiId == "" || name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}

	if err := store.DeleteDataSource(apiId, name); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// ListDataSources returns a paginated list of data sources for a GraphQL API.
// GET /v1/apis/{apiId}/datasources
func (s *AppSyncService) ListDataSources(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	opts := parsePaginationOptions(req)
	dataSources, nextToken, err := store.ListDataSources(apiId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]interface{}, 0, len(dataSources))
	for _, ds := range dataSources {
		items = append(items, dataSourceToMap(ds))
	}

	response := map[string]interface{}{
		"dataSources": items,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

// dataSourceToMap converts a DataSource struct to a response map with correct wire format.
func dataSourceToMap(ds *appsyncstore.DataSource) map[string]interface{} {
	m := map[string]interface{}{
		"name":          ds.Name,
		"type":          ds.Type,
		"dataSourceArn": ds.DataSourceArn,
	}

	if ds.Description != "" {
		m["description"] = ds.Description
	}
	if ds.ServiceRoleArn != "" {
		m["serviceRoleArn"] = ds.ServiceRoleArn
	}
	if ds.DynamodbConfig != nil {
		m["dynamodbConfig"] = dynamodbConfigToMap(ds.DynamodbConfig)
	}
	if ds.ElasticsearchConfig != nil {
		m["elasticsearchConfig"] = elasticsearchConfigToMap(ds.ElasticsearchConfig)
	}
	if ds.EventBridgeConfig != nil {
		m["eventBridgeConfig"] = eventBridgeConfigToMap(ds.EventBridgeConfig)
	}
	if ds.HttpConfig != nil {
		m["httpConfig"] = httpConfigToMap(ds.HttpConfig)
	}
	if ds.LambdaConfig != nil {
		m["lambdaConfig"] = lambdaDataSourceConfigToMap(ds.LambdaConfig)
	}
	if ds.MetricsConfig != "" {
		m["metricsConfig"] = ds.MetricsConfig
	}
	if ds.OpenSearchServiceConfig != nil {
		m["openSearchServiceConfig"] = openSearchServiceConfigToMap(ds.OpenSearchServiceConfig)
	}
	if ds.RelationalDatabaseConfig != nil {
		m["relationalDatabaseConfig"] = relationalDatabaseConfigToMap(ds.RelationalDatabaseConfig)
	}
	if len(ds.Tags) > 0 {
		m["tags"] = ds.Tags
	}

	return m
}

// parseDynamoDBConfig parses a DynamoDB data source config from request parameters.
func parseDynamoDBConfig(params map[string]interface{}) *appsyncstore.DynamodbDataSourceConfig {
	raw := request.GetMapParam(params, "dynamodbConfig")
	if raw == nil {
		return nil
	}
	cfg := &appsyncstore.DynamodbDataSourceConfig{
		AwsRegion:            request.GetStringParam(raw, "awsRegion"),
		TableName:            request.GetStringParam(raw, "tableName"),
		UseCallerCredentials: request.GetBoolParam(raw, "useCallerCredentials"),
		Versioned:            request.GetBoolParam(raw, "versioned"),
	}
	if deltaRaw := request.GetMapParam(raw, "deltaSyncConfig"); deltaRaw != nil {
		cfg.DeltaSyncConfig = &appsyncstore.DeltaSyncConfig{
			BaseTableTTL:       request.GetInt64Param(deltaRaw, "baseTableTTL"),
			DeltaSyncTableName: request.GetStringParam(deltaRaw, "deltaSyncTableName"),
			DeltaSyncTableTTL:  request.GetInt64Param(deltaRaw, "deltaSyncTableTTL"),
		}
	}
	return cfg
}

// parseElasticsearchConfig parses an Elasticsearch data source config from request parameters.
func parseElasticsearchConfig(params map[string]interface{}) *appsyncstore.ElasticsearchDataSourceConfig {
	raw := request.GetMapParam(params, "elasticsearchConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.ElasticsearchDataSourceConfig{
		AwsRegion: request.GetStringParam(raw, "awsRegion"),
		Endpoint:  request.GetStringParam(raw, "endpoint"),
	}
}

// parseEventBridgeConfig parses an EventBridge data source config from request parameters.
func parseEventBridgeConfig(params map[string]interface{}) *appsyncstore.EventBridgeDataSourceConfig {
	raw := request.GetMapParam(params, "eventBridgeConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.EventBridgeDataSourceConfig{
		EventBusArn: request.GetStringParam(raw, "eventBusArn"),
	}
}

// parseHttpConfig parses an HTTP data source config from request parameters.
func parseHttpConfig(params map[string]interface{}) *appsyncstore.HttpDataSourceConfig {
	raw := request.GetMapParam(params, "httpConfig")
	if raw == nil {
		return nil
	}
	cfg := &appsyncstore.HttpDataSourceConfig{
		Endpoint: request.GetStringParam(raw, "endpoint"),
	}
	if authRaw := request.GetMapParam(raw, "authorizationConfig"); authRaw != nil {
		cfg.AuthorizationConfig = &appsyncstore.AuthorizationConfig{
			AuthorizationType: request.GetStringParam(authRaw, "authorizationType"),
		}
		if iamRaw := request.GetMapParam(authRaw, "awsIamConfig"); iamRaw != nil {
			cfg.AuthorizationConfig.AwsIamConfig = &appsyncstore.AwsIamConfig{
				SigningRegion:      request.GetStringParam(iamRaw, "signingRegion"),
				SigningServiceName: request.GetStringParam(iamRaw, "signingServiceName"),
			}
		}
	}
	return cfg
}

// parseLambdaDataSourceConfig parses a Lambda data source config from request parameters.
func parseLambdaDataSourceConfig(params map[string]interface{}) *appsyncstore.LambdaDataSourceConfig {
	raw := request.GetMapParam(params, "lambdaConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.LambdaDataSourceConfig{
		LambdaFunctionArn: request.GetStringParam(raw, "lambdaFunctionArn"),
	}
}

// parseOpenSearchServiceConfig parses an OpenSearch Service data source config from request parameters.
func parseOpenSearchServiceConfig(params map[string]interface{}) *appsyncstore.OpenSearchServiceDataSourceConfig {
	raw := request.GetMapParam(params, "openSearchServiceConfig")
	if raw == nil {
		return nil
	}
	return &appsyncstore.OpenSearchServiceDataSourceConfig{
		AwsRegion: request.GetStringParam(raw, "awsRegion"),
		Endpoint:  request.GetStringParam(raw, "endpoint"),
	}
}

// parseRelationalDatabaseConfig parses a relational database data source config from request parameters.
func parseRelationalDatabaseConfig(params map[string]interface{}) *appsyncstore.RelationalDatabaseDataSourceConfig {
	raw := request.GetMapParam(params, "relationalDatabaseConfig")
	if raw == nil {
		return nil
	}
	cfg := &appsyncstore.RelationalDatabaseDataSourceConfig{
		RelationalDatabaseSourceType: request.GetStringParam(raw, "relationalDatabaseSourceType"),
	}
	if rdsRaw := request.GetMapParam(raw, "rdsHttpEndpointConfig"); rdsRaw != nil {
		cfg.RdsHttpEndpointConfig = &appsyncstore.RdsHttpEndpointConfig{
			AwsRegion:           request.GetStringParam(rdsRaw, "awsRegion"),
			AwsSecretStoreArn:   request.GetStringParam(rdsRaw, "awsSecretStoreArn"),
			DatabaseName:        request.GetStringParam(rdsRaw, "databaseName"),
			DbClusterIdentifier: request.GetStringParam(rdsRaw, "dbClusterIdentifier"),
			Schema:              request.GetStringParam(rdsRaw, "schema"),
		}
	}
	return cfg
}

// dynamodbConfigToMap converts a DynamodbDataSourceConfig to a wire-format map.
func dynamodbConfigToMap(c *appsyncstore.DynamodbDataSourceConfig) map[string]interface{} {
	m := map[string]interface{}{
		"awsRegion": c.AwsRegion,
		"tableName": c.TableName,
	}
	if c.UseCallerCredentials {
		m["useCallerCredentials"] = c.UseCallerCredentials
	}
	if c.Versioned {
		m["versioned"] = c.Versioned
	}
	if c.DeltaSyncConfig != nil {
		dsc := map[string]interface{}{}
		if c.DeltaSyncConfig.BaseTableTTL > 0 {
			dsc["baseTableTTL"] = c.DeltaSyncConfig.BaseTableTTL
		}
		if c.DeltaSyncConfig.DeltaSyncTableName != "" {
			dsc["deltaSyncTableName"] = c.DeltaSyncConfig.DeltaSyncTableName
		}
		if c.DeltaSyncConfig.DeltaSyncTableTTL > 0 {
			dsc["deltaSyncTableTTL"] = c.DeltaSyncConfig.DeltaSyncTableTTL
		}
		m["deltaSyncConfig"] = dsc
	}
	return m
}

// elasticsearchConfigToMap converts an ElasticsearchDataSourceConfig to a wire-format map.
func elasticsearchConfigToMap(c *appsyncstore.ElasticsearchDataSourceConfig) map[string]interface{} {
	return map[string]interface{}{
		"awsRegion": c.AwsRegion,
		"endpoint":  c.Endpoint,
	}
}

// eventBridgeConfigToMap converts an EventBridgeDataSourceConfig to a wire-format map.
func eventBridgeConfigToMap(c *appsyncstore.EventBridgeDataSourceConfig) map[string]interface{} {
	return map[string]interface{}{
		"eventBusArn": c.EventBusArn,
	}
}

// httpConfigToMap converts an HttpDataSourceConfig to a wire-format map.
func httpConfigToMap(c *appsyncstore.HttpDataSourceConfig) map[string]interface{} {
	m := map[string]interface{}{}
	if c.Endpoint != "" {
		m["endpoint"] = c.Endpoint
	}
	if c.AuthorizationConfig != nil {
		ac := map[string]interface{}{
			"authorizationType": c.AuthorizationConfig.AuthorizationType,
		}
		if c.AuthorizationConfig.AwsIamConfig != nil {
			iam := map[string]interface{}{}
			if c.AuthorizationConfig.AwsIamConfig.SigningRegion != "" {
				iam["signingRegion"] = c.AuthorizationConfig.AwsIamConfig.SigningRegion
			}
			if c.AuthorizationConfig.AwsIamConfig.SigningServiceName != "" {
				iam["signingServiceName"] = c.AuthorizationConfig.AwsIamConfig.SigningServiceName
			}
			ac["awsIamConfig"] = iam
		}
		m["authorizationConfig"] = ac
	}
	return m
}

// lambdaDataSourceConfigToMap converts a LambdaDataSourceConfig to a wire-format map.
func lambdaDataSourceConfigToMap(c *appsyncstore.LambdaDataSourceConfig) map[string]interface{} {
	return map[string]interface{}{
		"lambdaFunctionArn": c.LambdaFunctionArn,
	}
}

// openSearchServiceConfigToMap converts an OpenSearchServiceDataSourceConfig to a wire-format map.
func openSearchServiceConfigToMap(c *appsyncstore.OpenSearchServiceDataSourceConfig) map[string]interface{} {
	return map[string]interface{}{
		"awsRegion": c.AwsRegion,
		"endpoint":  c.Endpoint,
	}
}

// relationalDatabaseConfigToMap converts a RelationalDatabaseDataSourceConfig to a wire-format map.
func relationalDatabaseConfigToMap(c *appsyncstore.RelationalDatabaseDataSourceConfig) map[string]interface{} {
	m := map[string]interface{}{}
	if c.RelationalDatabaseSourceType != "" {
		m["relationalDatabaseSourceType"] = c.RelationalDatabaseSourceType
	}
	if c.RdsHttpEndpointConfig != nil {
		rds := map[string]interface{}{}
		if c.RdsHttpEndpointConfig.AwsRegion != "" {
			rds["awsRegion"] = c.RdsHttpEndpointConfig.AwsRegion
		}
		if c.RdsHttpEndpointConfig.AwsSecretStoreArn != "" {
			rds["awsSecretStoreArn"] = c.RdsHttpEndpointConfig.AwsSecretStoreArn
		}
		if c.RdsHttpEndpointConfig.DatabaseName != "" {
			rds["databaseName"] = c.RdsHttpEndpointConfig.DatabaseName
		}
		if c.RdsHttpEndpointConfig.DbClusterIdentifier != "" {
			rds["dbClusterIdentifier"] = c.RdsHttpEndpointConfig.DbClusterIdentifier
		}
		if c.RdsHttpEndpointConfig.Schema != "" {
			rds["schema"] = c.RdsHttpEndpointConfig.Schema
		}
		m["rdsHttpEndpointConfig"] = rds
	}
	return m
}
