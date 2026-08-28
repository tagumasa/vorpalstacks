package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// parseDataSourceInput extracts the shared data-source payload members from
// the request parameters.
func parseDataSourceInput(req *request.ParsedRequest) (dataSourceInput, error) {
	in := dataSourceInput{
		ApiId:          request.GetStringParam(req.Parameters, "apiId"),
		Name:           request.GetStringParam(req.Parameters, "name"),
		Type:           request.GetStringParam(req.Parameters, "type"),
		Description:    request.GetStringParam(req.Parameters, "description"),
		ServiceRoleArn: request.GetStringParam(req.Parameters, "serviceRoleArn"),
		MetricsConfig:  request.GetStringParam(req.Parameters, "metricsConfig"),
	}

	var err error
	in.RelationalDatabaseConfig, err = parseRelationalDatabaseConfig(req.Parameters)
	if err != nil {
		return dataSourceInput{}, err
	}
	in.DynamodbConfig, err = parseDynamoDBConfig(req.Parameters)
	if err != nil {
		return dataSourceInput{}, err
	}
	in.HttpConfig, err = parseHttpConfig(req.Parameters)
	if err != nil {
		return dataSourceInput{}, err
	}
	in.ElasticsearchConfig = parseElasticsearchConfig(req.Parameters)
	in.EventBridgeConfig = parseEventBridgeConfig(req.Parameters)
	in.LambdaConfig = parseLambdaDataSourceConfig(req.Parameters)
	in.NeptuneConfig = parseNeptuneConfig(req.Parameters)
	in.OpenSearchServiceConfig = parseOpenSearchServiceConfig(req.Parameters)

	return in, nil
}

// CreateDataSource creates a new data source for a GraphQL API.
// POST /v1/apis/{apiId}/datasources
func (s *AppSyncService) CreateDataSource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	in, err := parseDataSourceInput(req)
	if err != nil {
		return nil, err
	}

	created, err := s.createDataSourceCore(store, in)
	if err != nil {
		return nil, err
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

	ds, err := s.getDataSourceCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "name"))
	if err != nil {
		return nil, err
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

	in, err := parseDataSourceInput(req)
	if err != nil {
		return nil, err
	}

	updated, err := s.updateDataSourceCore(store, in)
	if err != nil {
		return nil, err
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

	if err := s.deleteDataSourceCore(store, request.GetStringParam(req.Parameters, "apiId"), request.GetStringParam(req.Parameters, "name")); err != nil {
		return nil, err
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

	dataSources, nextToken, err := s.listDataSourcesCore(store, apiId, request.GetIntParam(req.Parameters, "maxResults"), request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
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
