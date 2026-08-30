package athena

import (
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// --- DTOs ---

// ListDataCatalogsInput carries the raw MaxResults window
// (presence-flagged) and pagination marker.
type ListDataCatalogsInput struct {
	MaxResults    int
	HasMaxResults bool
	NextToken     string
}

// CreateDataCatalogInput carries the parsed wire members of a
// CreateDataCatalog request (Parameters already converted to string map).
type CreateDataCatalogInput struct {
	Name        string
	Description string
	Type        string
	Parameters  map[string]string
	Tags        map[string]string
}

// UpdateDataCatalogInput carries the parsed wire members of an
// UpdateDataCatalog request; HasParameters distinguishes an explicitly
// provided empty Parameters map from an omitted member.
type UpdateDataCatalogInput struct {
	Name          string
	Description   string
	Type          string
	Parameters    map[string]string
	HasParameters bool
}

// ListDatabasesInput carries the catalog name plus the raw MaxResults window
// (presence-flagged) and pagination marker.
type ListDatabasesInput struct {
	CatalogName   string
	MaxResults    int
	HasMaxResults bool
	NextToken     string
}

// ListTableMetadataInput carries the catalog/database names plus the raw
// MaxResults window (presence-flagged) and pagination marker.
type ListTableMetadataInput struct {
	CatalogName   string
	DatabaseName  string
	MaxResults    int
	HasMaxResults bool
	NextToken     string
}

// --- Core functions ---

// listDataCatalogsCore lists data catalogs with the built-in AwsDataCatalog
// entry prepended and pages them by name with the documented window
// semantics (default 50, range 2-50) applied after the walk, matching the
// original validation position.
func listDataCatalogsCore(stores *athenaStores, input ListDataCatalogsInput) ([]*athenastore.DataCatalog, string, error) {
	catalogs, err := stores.dataCatalogStore.ListDataCatalogs()
	if err != nil {
		return nil, "", err
	}

	all := append([]*athenastore.DataCatalog{
		{Name: "AwsDataCatalog", Type: "GLUE"},
	}, catalogs...)

	maxResults, err := resolveMaxResults(input.MaxResults, input.HasMaxResults, 50, 2, 50)
	if err != nil {
		return nil, "", err
	}

	pageResult := pagination.PaginateSlice(all, input.NextToken, maxResults, func(item *athenastore.DataCatalog) string {
		return item.Name
	})

	return pageResult.Items, pageResult.NextMarker, nil
}

// getDataCatalogCore fetches a data catalog; the built-in AwsDataCatalog is
// answered with the documented default catalog definition without touching
// storage, and the store is acquired only for user-defined catalogs, the
// order the original handler applied.
func (s *AthenaService) getDataCatalogCore(reqCtx *request.RequestContext, name string) (*athenastore.DataCatalog, error) {
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	if name == "AwsDataCatalog" {
		return &athenastore.DataCatalog{
			Name:        "AwsDataCatalog",
			Type:        "GLUE",
			Description: "The default AWS data catalog",
			Parameters:  map[string]string{},
		}, nil
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	catalog, err := stores.dataCatalogStore.GetDataCatalog(name)
	if err != nil {
		if err == athenastore.ErrDataCatalogNotFound {
			return nil, dataCatalogNotFound(name)
		}
		return nil, err
	}
	return catalog, nil
}

// createDataCatalogCore validates the create request, persists the catalog
// and applies its create-time tags. The store is acquired only after the
// request has been validated, the order the original handler applied.
func (s *AthenaService) createDataCatalogCore(reqCtx *request.RequestContext, input CreateDataCatalogInput) (*athenastore.DataCatalog, error) {
	if input.Name == "" {
		return nil, ErrInvalidRequestException
	}
	if input.Name == "AwsDataCatalog" {
		return nil, ErrInvalidRequestException
	}
	if err := validateCatalogNameString(input.Name); err != nil {
		return nil, err
	}

	if input.Description != "" {
		if err := validateDescriptionString(input.Description); err != nil {
			return nil, err
		}
	}
	catalogType := input.Type
	if catalogType == "" {
		catalogType = "GLUE"
	}
	if err := validateDataCatalogType(catalogType); err != nil {
		return nil, err
	}

	if err := validateTags(input.Tags); err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	catalog := &athenastore.DataCatalog{
		Name:        input.Name,
		Description: input.Description,
		Type:        catalogType,
		Parameters:  input.Parameters,
	}

	if err := stores.dataCatalogStore.CreateDataCatalog(catalog); err != nil {
		if err == athenastore.ErrDataCatalogAlreadyExists {
			return nil, alreadyExistsInvalidRequest("DataCatalog", input.Name)
		}
		return nil, err
	}

	if len(input.Tags) > 0 {
		if err := stores.dataCatalogStore.Tag(input.Name, input.Tags); err != nil {
			logs.Warn("failed to tag data catalog", logs.String("catalog", input.Name), logs.Err(err))
		}
	}

	return catalog, nil
}

// deleteDataCatalogCore deletes a data catalog, refusing the built-in
// AwsDataCatalog and mapping the store not-found sentinel. The store is
// acquired only after the name validation, the order the original handler
// applied.
func (s *AthenaService) deleteDataCatalogCore(reqCtx *request.RequestContext, name string) error {
	if name == "" {
		return ErrInvalidRequestException
	}

	if name == "AwsDataCatalog" {
		return ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := stores.dataCatalogStore.DeleteDataCatalog(name); err != nil {
		if err == athenastore.ErrDataCatalogNotFound {
			return dataCatalogNotFound(name)
		}
		return err
	}
	return nil
}

// updateDataCatalogCore validates the update request (per the Smithy model,
// Name and Type are both REQUIRED), applies it to the stored record and
// persists it.
func (s *AthenaService) updateDataCatalogCore(reqCtx *request.RequestContext, input UpdateDataCatalogInput) (*athenastore.DataCatalog, error) {
	if input.Name == "" {
		return nil, ErrInvalidRequestException
	}

	if input.Name == "AwsDataCatalog" {
		return nil, ErrInvalidRequestException
	}

	if input.Type == "" {
		return nil, invalidRequestParameter("Type is required for UpdateDataCatalog")
	}
	if err := validateDataCatalogType(input.Type); err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	catalog, err := stores.dataCatalogStore.GetDataCatalog(input.Name)
	if err != nil {
		if err == athenastore.ErrDataCatalogNotFound {
			return nil, dataCatalogNotFound(input.Name)
		}
		return nil, err
	}

	catalog.Type = input.Type

	if input.Description != "" {
		if err := validateDescriptionString(input.Description); err != nil {
			return nil, err
		}
		catalog.Description = input.Description
	}

	if input.HasParameters {
		catalog.Parameters = input.Parameters
	}

	if err := stores.dataCatalogStore.UpdateDataCatalog(catalog); err != nil {
		return nil, err
	}

	return catalog, nil
}

// listDatabasesCore lists databases in a catalog (injecting the default
// database into AwsDataCatalog when absent) and pages them by name with the
// documented window semantics (default 50, range 1-50) applied before the
// walk, matching the original validation position. The store is acquired
// only after the window has been resolved.
func (s *AthenaService) listDatabasesCore(reqCtx *request.RequestContext, input ListDatabasesInput) ([]*athenastore.Database, string, error) {
	catalogName := input.CatalogName
	if catalogName == "" {
		catalogName = "AwsDataCatalog"
	}

	maxResults, err := resolveMaxResults(input.MaxResults, input.HasMaxResults, 50, 1, 50)
	if err != nil {
		return nil, "", err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}

	databases, err := stores.databaseStore.ListDatabases(catalogName)
	if err != nil {
		return nil, "", err
	}

	all := make([]*athenastore.Database, 0, len(databases)+1)
	for _, db := range databases {
		all = append(all, db)
	}

	if catalogName == "AwsDataCatalog" {
		hasDefault := false
		for _, db := range all {
			if db.Name == "default" {
				hasDefault = true
				break
			}
		}
		if !hasDefault {
			all = append(all, &athenastore.Database{
				Name:        "default",
				Description: "Default database",
			})
		}
	}

	pageResult := pagination.PaginateSlice(all, input.NextToken, maxResults, func(item *athenastore.Database) string {
		return item.Name
	})

	return pageResult.Items, pageResult.NextMarker, nil
}

// getDatabaseCore fetches database metadata; the AwsDataCatalog "default"
// database is answered with the documented definition without touching
// storage, and the store is acquired only for other databases, the order
// the original handler applied.
func (s *AthenaService) getDatabaseCore(reqCtx *request.RequestContext, catalogName, databaseName string) (*athenastore.Database, error) {
	if catalogName == "" {
		catalogName = "AwsDataCatalog"
	}

	if databaseName == "" {
		return nil, ErrInvalidRequestException
	}

	if catalogName == "AwsDataCatalog" && databaseName == "default" {
		return &athenastore.Database{
			Name:        "default",
			Description: "Default database",
			Parameters:  map[string]string{"EXTERNAL": "TRUE"},
		}, nil
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	db, err := stores.databaseStore.GetDatabase(catalogName, databaseName)
	if err != nil {
		if err == athenastore.ErrDatabaseNotFound {
			return nil, ErrMetadataException
		}
		return nil, err
	}
	return db, nil
}

// listTableMetadataCore lists table metadata in a database and pages it by
// name with the documented window semantics (default 50, range 1-50)
// applied after the walk, matching the original validation position.
func (s *AthenaService) listTableMetadataCore(reqCtx *request.RequestContext, input ListTableMetadataInput) ([]*athenastore.TableMetadata, string, error) {
	catalogName := input.CatalogName
	if catalogName == "" {
		catalogName = "AwsDataCatalog"
	}

	if input.DatabaseName == "" {
		return nil, "", ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}

	tables, err := stores.tableStore.ListTables(catalogName, input.DatabaseName)
	if err != nil {
		return nil, "", err
	}

	maxResults, err := resolveMaxResults(input.MaxResults, input.HasMaxResults, 50, 1, 50)
	if err != nil {
		return nil, "", err
	}

	pageResult := pagination.PaginateSlice(tables, input.NextToken, maxResults, func(item *athenastore.TableMetadata) string {
		return item.Name
	})

	return pageResult.Items, pageResult.NextMarker, nil
}

// getTableMetadataCore fetches table metadata, mapping the store not-found
// sentinel onto MetadataException.
func (s *AthenaService) getTableMetadataCore(reqCtx *request.RequestContext, catalogName, databaseName, tableName string) (*athenastore.TableMetadata, error) {
	if catalogName == "" {
		catalogName = "AwsDataCatalog"
	}

	if databaseName == "" {
		return nil, ErrInvalidRequestException
	}

	if tableName == "" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	table, err := stores.tableStore.GetTable(catalogName, databaseName, tableName)
	if err != nil {
		if err == athenastore.ErrTableNotFound {
			return nil, ErrMetadataException
		}
		return nil, err
	}
	return table, nil
}
