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
		return mapStoreError(err)
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
		NeptuneConfig:            parseNeptuneConfig(req.Parameters),
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
		return mapStoreError(err)
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
		return mapStoreError(err)
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
		NeptuneConfig:            parseNeptuneConfig(req.Parameters),
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
		return mapStoreError(err)
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
		return mapStoreError(err)
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
