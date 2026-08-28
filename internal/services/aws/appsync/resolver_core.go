package appsync

import (
	"fmt"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

// resolverInput carries the parsed resolver request payload shared by the
// create and update operations. HasMaxBatchSize distinguishes an explicitly
// supplied maxBatchSize from an omitted one.
type resolverInput struct {
	ApiId                   string
	TypeName                string
	FieldName               string
	Kind                    string
	DataSourceName          string
	RequestMappingTemplate  string
	ResponseMappingTemplate string
	Code                    string
	MetricsConfig           string
	Runtime                 *appsyncstore.AppSyncRuntime
	CachingConfig           *appsyncstore.CachingConfig
	PipelineConfig          *appsyncstore.PipelineConfig
	SyncConfig              *appsyncstore.SyncConfig
	MaxBatchSize            int32
	HasMaxBatchSize         bool
}

// createResolverCore validates the request and persists a new resolver for a
// GraphQL API type and field.
func (s *AppSyncService) createResolverCore(store *appsyncstore.AppSyncStore, in resolverInput) (*appsyncstore.Resolver, error) {
	if in.ApiId == "" || in.TypeName == "" || in.FieldName == "" {
		return nil, NewBadRequestException("apiId, typeName, and fieldName are required")
	}
	if err := validateResourceName(in.TypeName); err != nil {
		return nil, err
	}
	if err := validateResourceName(in.FieldName); err != nil {
		return nil, err
	}
	if err := validateGraphqlApiExists(store, in.ApiId); err != nil {
		return nil, err
	}

	r := &appsyncstore.Resolver{
		ApiId:                   in.ApiId,
		TypeName:                in.TypeName,
		FieldName:               in.FieldName,
		Kind:                    in.Kind,
		DataSourceName:          in.DataSourceName,
		RequestMappingTemplate:  in.RequestMappingTemplate,
		ResponseMappingTemplate: in.ResponseMappingTemplate,
		Runtime:                 in.Runtime,
		Code:                    in.Code,
		CachingConfig:           in.CachingConfig,
		MaxBatchSize:            in.MaxBatchSize,
		MetricsConfig:           in.MetricsConfig,
		PipelineConfig:          in.PipelineConfig,
		SyncConfig:              in.SyncConfig,
	}

	if err := validateResolverPayload(r, in.HasMaxBatchSize); err != nil {
		return nil, err
	}

	created, err := store.CreateResolver(r)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	s.schemaCache.Delete(in.ApiId)

	return created, nil
}

// getResolverCore fetches a resolver by API ID, type name, and field name.
func (s *AppSyncService) getResolverCore(store *appsyncstore.AppSyncStore, apiId, typeName, fieldName string) (*appsyncstore.Resolver, error) {
	if apiId == "" || typeName == "" || fieldName == "" {
		return nil, NewBadRequestException("apiId, typeName, and fieldName are required")
	}

	r, err := store.GetResolver(apiId, typeName, fieldName)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return r, nil
}

// updateResolverCore validates the request and applies an update to an
// existing resolver.
func (s *AppSyncService) updateResolverCore(store *appsyncstore.AppSyncStore, in resolverInput) (*appsyncstore.Resolver, error) {
	if in.ApiId == "" || in.TypeName == "" || in.FieldName == "" {
		return nil, NewBadRequestException("apiId, typeName, and fieldName are required")
	}

	r := &appsyncstore.Resolver{
		ApiId:                   in.ApiId,
		TypeName:                in.TypeName,
		FieldName:               in.FieldName,
		Kind:                    in.Kind,
		DataSourceName:          in.DataSourceName,
		RequestMappingTemplate:  in.RequestMappingTemplate,
		ResponseMappingTemplate: in.ResponseMappingTemplate,
		Runtime:                 in.Runtime,
		Code:                    in.Code,
		CachingConfig:           in.CachingConfig,
		MaxBatchSize:            in.MaxBatchSize,
		MetricsConfig:           in.MetricsConfig,
		PipelineConfig:          in.PipelineConfig,
		SyncConfig:              in.SyncConfig,
	}

	if err := validateResolverPayload(r, in.HasMaxBatchSize); err != nil {
		return nil, err
	}

	updated, err := store.UpdateResolver(r)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	s.schemaCache.Delete(in.ApiId)

	return updated, nil
}

// deleteResolverCore removes a resolver and invalidates the engine's schema
// cache for the API.
func (s *AppSyncService) deleteResolverCore(store *appsyncstore.AppSyncStore, apiId, typeName, fieldName string) error {
	if apiId == "" || typeName == "" || fieldName == "" {
		return NewBadRequestException("apiId, typeName, and fieldName are required")
	}

	if err := store.DeleteResolver(apiId, typeName, fieldName); err != nil {
		return mapStoreErrorE(err)
	}

	s.schemaCache.Delete(apiId)

	return nil
}

// listResolversCore lists the resolvers attached to one type of a GraphQL API.
func (s *AppSyncService) listResolversCore(store *appsyncstore.AppSyncStore, apiId, typeName string, maxResults int, nextToken string) ([]*appsyncstore.Resolver, string, error) {
	if apiId == "" || typeName == "" {
		return nil, "", NewBadRequestException("apiId and typeName are required")
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	resolvers, nextToken, err := store.ListResolvers(apiId, typeName, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return resolvers, nextToken, nil
}

// listResolversByFunctionCore lists the resolvers that reference a given
// AppSync function.
func (s *AppSyncService) listResolversByFunctionCore(store *appsyncstore.AppSyncStore, apiId, functionId string, maxResults int, nextToken string) ([]*appsyncstore.Resolver, string, error) {
	if apiId == "" || functionId == "" {
		return nil, "", NewBadRequestException("apiId and functionId are required")
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	resolvers, nextToken, err := store.ListResolversByFunction(apiId, functionId, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return resolvers, nextToken, nil
}

// validateResolverPayload applies the shared kind/caching/runtime/metrics/
// batch-size/code/template checks to a built resolver in the order the
// operations historically applied them.
func validateResolverPayload(r *appsyncstore.Resolver, hasMaxBatchSize bool) error {
	if r.Kind != "" && !validateResolverKind(r.Kind) {
		return NewBadRequestException(fmt.Sprintf("Invalid resolver kind: %s. Valid values: UNIT, PIPELINE", r.Kind))
	}
	if err := validateCachingConfig(r.CachingConfig); err != nil {
		return err
	}
	if err := validateAppSyncRuntime(r.Runtime); err != nil {
		return err
	}
	if r.MetricsConfig != "" && !validateEnabledDisabled(r.MetricsConfig) {
		return NewBadRequestException(fmt.Sprintf("Invalid metricsConfig: %s", r.MetricsConfig))
	}
	if hasMaxBatchSize {
		if err := validateMaxBatchSize(r.MaxBatchSize); err != nil {
			return err
		}
	}
	if r.Code != "" {
		if err := validateCode(r.Code); err != nil {
			return err
		}
	}
	if r.RequestMappingTemplate != "" {
		if err := validateMappingTemplate(r.RequestMappingTemplate); err != nil {
			return err
		}
	}
	if r.ResponseMappingTemplate != "" {
		if err := validateMappingTemplate(r.ResponseMappingTemplate); err != nil {
			return err
		}
	}
	return nil
}
