package appsync

import (
	"fmt"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// dataSourceInput carries the parsed data-source request payload shared by
// the create and update operations.
type dataSourceInput struct {
	ApiId                    string
	Name                     string
	Type                     string
	Description              string
	ServiceRoleArn           string
	MetricsConfig            string
	DynamodbConfig           *appsyncstore.DynamodbDataSourceConfig
	ElasticsearchConfig      *appsyncstore.ElasticsearchDataSourceConfig
	EventBridgeConfig        *appsyncstore.EventBridgeDataSourceConfig
	HttpConfig               *appsyncstore.HttpDataSourceConfig
	LambdaConfig             *appsyncstore.LambdaDataSourceConfig
	NeptuneConfig            *appsyncstore.NeptuneDataSourceConfig
	OpenSearchServiceConfig  *appsyncstore.OpenSearchServiceDataSourceConfig
	RelationalDatabaseConfig *appsyncstore.RelationalDatabaseDataSourceConfig
}

// createDataSourceCore validates the request and persists a new data source.
func (s *AppSyncService) createDataSourceCore(store *appsyncstore.AppSyncStore, in dataSourceInput) (*appsyncstore.DataSource, error) {
	if in.ApiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}
	if err := validateGraphqlApiExists(store, in.ApiId); err != nil {
		return nil, err
	}

	if in.Name == "" {
		return nil, NewBadRequestException("name is required")
	}
	if err := validateResourceName(in.Name); err != nil {
		return nil, err
	}

	if in.Type == "" {
		return nil, NewBadRequestException("type is required")
	}
	if !validateDataSourceType(in.Type) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid data source type: %s", in.Type))
	}

	if in.MetricsConfig != "" && !validateEnabledDisabled(in.MetricsConfig) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid metricsConfig: %s", in.MetricsConfig))
	}

	if err := validateDescription(in.Description); err != nil {
		return nil, err
	}

	if err := checkRelationalDatabaseSourceType(in.RelationalDatabaseConfig); err != nil {
		return nil, err
	}

	ds := &appsyncstore.DataSource{
		ApiId:                    in.ApiId,
		Name:                     in.Name,
		Type:                     in.Type,
		Description:              in.Description,
		ServiceRoleArn:           in.ServiceRoleArn,
		DynamodbConfig:           in.DynamodbConfig,
		ElasticsearchConfig:      in.ElasticsearchConfig,
		EventBridgeConfig:        in.EventBridgeConfig,
		HttpConfig:               in.HttpConfig,
		LambdaConfig:             in.LambdaConfig,
		MetricsConfig:            in.MetricsConfig,
		NeptuneConfig:            in.NeptuneConfig,
		OpenSearchServiceConfig:  in.OpenSearchServiceConfig,
		RelationalDatabaseConfig: in.RelationalDatabaseConfig,
	}

	created, err := store.CreateDataSource(ds)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return created, nil
}

// getDataSourceCore fetches a data source by API ID and name.
func (s *AppSyncService) getDataSourceCore(store *appsyncstore.AppSyncStore, apiId, name string) (*appsyncstore.DataSource, error) {
	if apiId == "" || name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}

	ds, err := store.GetDataSource(apiId, name)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return ds, nil
}

// updateDataSourceCore validates the request and applies a full-record update
// to an existing data source.
func (s *AppSyncService) updateDataSourceCore(store *appsyncstore.AppSyncStore, in dataSourceInput) (*appsyncstore.DataSource, error) {
	if in.ApiId == "" || in.Name == "" {
		return nil, NewBadRequestException("apiId and name are required")
	}
	if err := validateResourceName(in.Name); err != nil {
		return nil, err
	}

	if in.Type == "" {
		return nil, NewBadRequestException("type is required")
	}
	if !validateDataSourceType(in.Type) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid data source type: %s", in.Type))
	}

	if in.MetricsConfig != "" && !validateEnabledDisabled(in.MetricsConfig) {
		return nil, NewBadRequestException(fmt.Sprintf("Invalid metricsConfig: %s", in.MetricsConfig))
	}

	if err := validateDescription(in.Description); err != nil {
		return nil, err
	}

	if err := checkRelationalDatabaseSourceType(in.RelationalDatabaseConfig); err != nil {
		return nil, err
	}

	ds := &appsyncstore.DataSource{
		ApiId:                    in.ApiId,
		Name:                     in.Name,
		Type:                     in.Type,
		Description:              in.Description,
		ServiceRoleArn:           in.ServiceRoleArn,
		DynamodbConfig:           in.DynamodbConfig,
		ElasticsearchConfig:      in.ElasticsearchConfig,
		EventBridgeConfig:        in.EventBridgeConfig,
		HttpConfig:               in.HttpConfig,
		LambdaConfig:             in.LambdaConfig,
		MetricsConfig:            in.MetricsConfig,
		NeptuneConfig:            in.NeptuneConfig,
		OpenSearchServiceConfig:  in.OpenSearchServiceConfig,
		RelationalDatabaseConfig: in.RelationalDatabaseConfig,
	}

	updated, err := store.UpdateDataSource(ds)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return updated, nil
}

// deleteDataSourceCore removes a data source.
func (s *AppSyncService) deleteDataSourceCore(store *appsyncstore.AppSyncStore, apiId, name string) error {
	if apiId == "" || name == "" {
		return NewBadRequestException("apiId and name are required")
	}

	if err := store.DeleteDataSource(apiId, name); err != nil {
		return mapStoreErrorE(err)
	}

	return nil
}

// listDataSourcesCore lists the data sources of a GraphQL API.
func (s *AppSyncService) listDataSourcesCore(store *appsyncstore.AppSyncStore, apiId string, maxResults int, nextToken string) ([]*appsyncstore.DataSource, string, error) {
	if apiId == "" {
		return nil, "", NewBadRequestException("apiId is required")
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	dataSources, nextToken, err := store.ListDataSources(apiId, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return dataSources, nextToken, nil
}

// checkRelationalDatabaseSourceType rejects an off-enum
// relationalDatabaseSourceType on a supplied relational database config.
func checkRelationalDatabaseSourceType(cfg *appsyncstore.RelationalDatabaseDataSourceConfig) error {
	if cfg != nil && cfg.RelationalDatabaseSourceType != "" && !validateRelationalDatabaseSourceType(cfg.RelationalDatabaseSourceType) {
		return NewBadRequestException(fmt.Sprintf("Invalid relationalDatabaseSourceType: %s", cfg.RelationalDatabaseSourceType))
	}
	return nil
}
