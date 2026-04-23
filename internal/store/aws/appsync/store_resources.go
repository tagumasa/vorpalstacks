package appsync

import (
	"encoding/json"

	"vorpalstacks/internal/store/aws/common"
)

// --- Function (AppSync) ---

// CreateFunction persists a new AppSync function (a reusable resolver unit).
func (s *AppSyncStore) CreateFunction(f *FunctionConfiguration) (*FunctionConfiguration, error) {
	if f.ApiId == "" || f.Name == "" || f.DataSourceName == "" {
		return nil, common.NewStoreError("appsync", "create_function", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	if f.FunctionId == "" {
		f.FunctionId = s.GenerateId()
	}

	key := f.ApiId + "/" + f.FunctionId
	if s.functionsStore.Exists(key) {
		return nil, ErrFunctionAlreadyExists
	}

	f.FunctionArn = s.BuildFunctionARN(f.ApiId, f.FunctionId)
	if f.FunctionVersion == "" {
		f.FunctionVersion = "2018-05-29"
	}

	if err := s.functionsStore.Put(key, f); err != nil {
		return nil, err
	}
	return f, nil
}

// GetFunction retrieves an AppSync function by API ID and function ID.
func (s *AppSyncStore) GetFunction(apiId, functionId string) (*FunctionConfiguration, error) {
	key := apiId + "/" + functionId
	var f FunctionConfiguration
	if err := s.functionsStore.Get(key, &f); err != nil {
		return nil, ErrFunctionNotFound
	}
	return &f, nil
}

// UpdateFunction merges non-zero fields from the input into the existing function.
func (s *AppSyncStore) UpdateFunction(f *FunctionConfiguration) (*FunctionConfiguration, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := f.ApiId + "/" + f.FunctionId
	existing, err := s.GetFunction(f.ApiId, f.FunctionId)
	if err != nil {
		return nil, err
	}

	if f.Name != "" {
		existing.Name = f.Name
	}
	if f.DataSourceName != "" {
		existing.DataSourceName = f.DataSourceName
	}
	if f.Description != "" {
		existing.Description = f.Description
	}
	if f.FunctionVersion != "" {
		existing.FunctionVersion = f.FunctionVersion
	}
	if f.RequestMappingTemplate != "" {
		existing.RequestMappingTemplate = f.RequestMappingTemplate
	}
	if f.ResponseMappingTemplate != "" {
		existing.ResponseMappingTemplate = f.ResponseMappingTemplate
	}
	if f.Runtime != nil {
		existing.Runtime = f.Runtime
	}
	if f.Code != "" {
		existing.Code = f.Code
	}
	if f.MaxBatchSize > 0 {
		existing.MaxBatchSize = f.MaxBatchSize
	}
	if f.SyncConfig != nil {
		existing.SyncConfig = f.SyncConfig
	}

	if err := s.functionsStore.Put(key, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteFunction removes an AppSync function by API ID and function ID.
func (s *AppSyncStore) DeleteFunction(apiId, functionId string) error {
	key := apiId + "/" + functionId
	if !s.functionsStore.Exists(key) {
		return ErrFunctionNotFound
	}
	return s.functionsStore.Delete(key)
}

// ListFunctions returns a paginated list of AppSync functions for a given GraphQL API.
func (s *AppSyncStore) ListFunctions(apiId string, opts common.ListOptions) ([]*FunctionConfiguration, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[FunctionConfiguration](s.functionsStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- Type ---

// CreateType persists a new type definition scoped to a GraphQL API.
// Extracts the type name from the SDL definition if the Name field is not set.
func (s *AppSyncStore) CreateType(t *Type) (*Type, error) {
	if t.ApiId == "" || t.Definition == "" || t.Format == "" {
		return nil, common.NewStoreError("appsync", "create_type", common.ErrInvalidInput)
	}

	if t.Name == "" {
		t.Name = extractTypeName(t.Definition)
	}
	if t.Name == "" {
		return nil, common.NewStoreError("appsync", "create_type", common.ErrInvalidInput)
	}

	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := t.ApiId + "/" + t.Name
	if s.typesStore.Exists(key) {
		return nil, ErrTypeAlreadyExists
	}

	t.Arn = s.BuildTypeARN(t.ApiId, t.Name)

	if err := s.typesStore.Put(key, t); err != nil {
		return nil, err
	}
	return t, nil
}

// GetType retrieves a type definition by API ID and type name.
func (s *AppSyncStore) GetType(apiId, typeName string) (*Type, error) {
	key := apiId + "/" + typeName
	var t Type
	if err := s.typesStore.Get(key, &t); err != nil {
		return nil, ErrTypeNotFound
	}
	return &t, nil
}

// UpdateType merges non-zero fields from the input into the existing type definition.
func (s *AppSyncStore) UpdateType(t *Type) (*Type, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := t.ApiId + "/" + t.Name
	existing, err := s.GetType(t.ApiId, t.Name)
	if err != nil {
		return nil, err
	}

	if t.Definition != "" {
		existing.Definition = t.Definition
	}
	if t.Format != "" {
		existing.Format = t.Format
	}
	if t.Description != "" {
		existing.Description = t.Description
	}

	if err := s.typesStore.Put(key, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteType removes a type definition by API ID and type name.
func (s *AppSyncStore) DeleteType(apiId, typeName string) error {
	key := apiId + "/" + typeName
	if !s.typesStore.Exists(key) {
		return ErrTypeNotFound
	}
	return s.typesStore.Delete(key)
}

// ListTypes returns a paginated list of type definitions for a given GraphQL API.
func (s *AppSyncStore) ListTypes(apiId string, opts common.ListOptions) ([]*Type, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[Type](s.typesStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- Schema Creation ---

// SaveSchemaCreationStatus persists the status of a schema creation operation.
func (s *AppSyncStore) SaveSchemaCreationStatus(apiId string, status *SchemaCreationStatus) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	return s.schemaStatusesStore.Put(apiId, status)
}

// GetSchemaCreationStatus retrieves the status of a schema creation operation.
func (s *AppSyncStore) GetSchemaCreationStatus(apiId string) (*SchemaCreationStatus, error) {
	var status SchemaCreationStatus
	if err := s.schemaStatusesStore.Get(apiId, &status); err != nil {
		return nil, ErrSchemaCreationNotFound
	}
	return &status, nil
}

// --- Environment Variables ---

// SaveEnvironmentVariables persists environment variables for a GraphQL API.
func (s *AppSyncStore) SaveEnvironmentVariables(apiId string, envVars *EnvironmentVariables) error {
	return s.envVariablesStore.Put(apiId, envVars)
}

// GetEnvironmentVariables retrieves environment variables for a GraphQL API.
func (s *AppSyncStore) GetEnvironmentVariables(apiId string) (*EnvironmentVariables, error) {
	var envVars EnvironmentVariables
	if err := s.envVariablesStore.Get(apiId, &envVars); err != nil {
		return &EnvironmentVariables{EnvironmentVariables: map[string]string{}}, nil
	}
	return &envVars, nil
}

// --- API Keys ---
// API keys are stored with composite keys "apiId/keyId" to ensure uniqueness
// across GraphQL APIs and to support scoped listing by apiId prefix.

// CreateApiKey persists a new API key.
func (s *AppSyncStore) CreateApiKey(apiId string, apiKey *ApiKey) error {
	return s.apiKeysStore.Put(apiId+"/"+apiKey.Id, apiKey)
}

// GetApiKey retrieves an API key by ID (requires apiId for composite key lookup).
func (s *AppSyncStore) GetApiKey(apiId, id string) (*ApiKey, error) {
	var apiKey ApiKey
	if err := s.apiKeysStore.Get(apiId+"/"+id, &apiKey); err != nil {
		return nil, ErrApiKeyNotFound
	}
	return &apiKey, nil
}

// UpdateApiKey updates an existing API key.
func (s *AppSyncStore) UpdateApiKey(apiId string, apiKey *ApiKey) error {
	return s.apiKeysStore.Put(apiId+"/"+apiKey.Id, apiKey)
}

// DeleteApiKey removes an API key by ID.
func (s *AppSyncStore) DeleteApiKey(apiId, id string) error {
	key := apiId + "/" + id
	if !s.apiKeysStore.Exists(key) {
		return ErrApiKeyNotFound
	}
	return s.apiKeysStore.Delete(key)
}

// ListApiKeys lists API keys for a given API, optionally filtered by expiry.
func (s *AppSyncStore) ListApiKeys(apiId string, opts common.ListOptions) ([]*ApiKey, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   apiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[ApiKey](s.apiKeysStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- API Cache ---

// CreateApiCache persists a new API cache configuration.
func (s *AppSyncStore) CreateApiCache(apiId string, cache *ApiCache) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	if s.apiCachesStore.Exists(apiId) {
		return ErrApiCacheAlreadyExists
	}

	return s.apiCachesStore.Put(apiId, cache)
}

// GetApiCache retrieves the API cache configuration for a given API.
func (s *AppSyncStore) GetApiCache(apiId string) (*ApiCache, error) {
	var cache ApiCache
	if err := s.apiCachesStore.Get(apiId, &cache); err != nil {
		return nil, ErrApiCacheNotFound
	}
	return &cache, nil
}

// UpdateApiCache updates the API cache configuration.
func (s *AppSyncStore) UpdateApiCache(apiId string, cache *ApiCache) error {
	return s.apiCachesStore.Put(apiId, cache)
}

// DeleteApiCache removes the API cache configuration.
// Validates existence first, as BaseStore.Delete silently succeeds on missing keys.
func (s *AppSyncStore) DeleteApiCache(apiId string) error {
	if !s.apiCachesStore.Exists(apiId) {
		return ErrApiCacheNotFound
	}
	return s.apiCachesStore.Delete(apiId)
}

// --- Domain Names ---

// CreateDomainName persists a new domain name configuration.
func (s *AppSyncStore) CreateDomainName(domainName *DomainNameConfig) error {
	return s.domainNamesStore.Put(domainName.DomainName, domainName)
}

// GetDomainName retrieves a domain name configuration.
func (s *AppSyncStore) GetDomainName(domainName string) (*DomainNameConfig, error) {
	var config DomainNameConfig
	if err := s.domainNamesStore.Get(domainName, &config); err != nil {
		return nil, ErrDomainNameNotFound
	}
	return &config, nil
}

// UpdateDomainName updates a domain name configuration.
func (s *AppSyncStore) UpdateDomainName(domainName *DomainNameConfig) error {
	return s.domainNamesStore.Put(domainName.DomainName, domainName)
}

// DeleteDomainName removes a domain name configuration.
// Validates existence first, as BaseStore.Delete silently succeeds on missing keys.
func (s *AppSyncStore) DeleteDomainName(domainName string) error {
	if !s.domainNamesStore.Exists(domainName) {
		return ErrDomainNameNotFound
	}
	return s.domainNamesStore.Delete(domainName)
}

// ListDomainNames lists all domain names.
func (s *AppSyncStore) ListDomainNames(opts common.ListOptions) ([]*DomainNameConfig, string, error) {
	result, err := common.List[DomainNameConfig](s.domainNamesStore, opts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}

// --- API Associations ---

// AssociateApi persists an API association with a domain name.
func (s *AppSyncStore) AssociateApi(association *ApiAssociation) error {
	return s.apiAssociationsStore.Put(association.DomainName, association)
}

// GetApiAssociation retrieves the API association for a domain name.
func (s *AppSyncStore) GetApiAssociation(domainName string) (*ApiAssociation, error) {
	var assoc ApiAssociation
	if err := s.apiAssociationsStore.Get(domainName, &assoc); err != nil {
		return nil, ErrApiAssociationNotFound
	}
	return &assoc, nil
}

// DisassociateApi removes the API association for a domain name.
// Validates existence first, as BaseStore.Delete silently succeeds on missing keys.
func (s *AppSyncStore) DisassociateApi(domainName string) error {
	if !s.apiAssociationsStore.Exists(domainName) {
		return ErrApiAssociationNotFound
	}
	return s.apiAssociationsStore.Delete(domainName)
}

// --- Merged API Associations ---
// Associations are stored with composite keys "mergedApiId/associationId" to
// support efficient prefix-based listing by merged API. For source-side lookups
// (where only associationId is known), GetMergedApiAssociationById performs a
// full scan with filter.

// CreateMergedApiAssociation persists a merged API association.
func (s *AppSyncStore) CreateMergedApiAssociation(assoc *SourceApiAssociation) error {
	key := assoc.MergedApiId + "/" + assoc.AssociationId
	return s.mergedApiAssociationsStore.Put(key, assoc)
}

// GetMergedApiAssociation retrieves a merged API association by merged API ID and association ID.
func (s *AppSyncStore) GetMergedApiAssociation(mergedApiId, associationId string) (*SourceApiAssociation, error) {
	key := mergedApiId + "/" + associationId
	var assoc SourceApiAssociation
	if err := s.mergedApiAssociationsStore.Get(key, &assoc); err != nil {
		return nil, ErrMergedApiAssociationNotFound
	}
	return &assoc, nil
}

// GetMergedApiAssociationById retrieves a merged API association by its UUID alone.
// This is required for source-side operations (e.g. DisassociateMergedGraphqlApi)
// where the mergedApiId is not available from the request path.
func (s *AppSyncStore) GetMergedApiAssociationById(associationId string) (*SourceApiAssociation, error) {
	var found *SourceApiAssociation
	err := s.mergedApiAssociationsStore.ForEach(func(key string, value []byte) error {
		var assoc SourceApiAssociation
		if err := json.Unmarshal(value, &assoc); err != nil {
			return err
		}
		if assoc.AssociationId == associationId {
			found = &assoc
		}
		return nil
	})
	if err != nil || found == nil {
		return nil, ErrMergedApiAssociationNotFound
	}
	return found, nil
}

// UpdateMergedApiAssociation updates an existing merged API association.
func (s *AppSyncStore) UpdateMergedApiAssociation(assoc *SourceApiAssociation) error {
	key := assoc.MergedApiId + "/" + assoc.AssociationId
	return s.mergedApiAssociationsStore.Put(key, assoc)
}

// DeleteMergedApiAssociation removes a merged API association by merged API ID and association ID.
func (s *AppSyncStore) DeleteMergedApiAssociation(mergedApiId, associationId string) error {
	key := mergedApiId + "/" + associationId
	return s.mergedApiAssociationsStore.Delete(key)
}

// ListMergedApiAssociations lists merged API associations for a given merged API.
func (s *AppSyncStore) ListMergedApiAssociations(mergedApiId string, opts common.ListOptions) ([]*SourceApiAssociation, string, error) {
	prefixOpts := common.ListOptions{
		Prefix:   mergedApiId + "/",
		Marker:   opts.Marker,
		MaxItems: opts.MaxItems,
	}
	result, err := common.List[SourceApiAssociation](s.mergedApiAssociationsStore, prefixOpts, nil)
	if err != nil {
		return nil, "", err
	}
	var nextToken string
	if result.IsTruncated {
		nextToken = result.NextMarker
	}
	return result.Items, nextToken, nil
}
