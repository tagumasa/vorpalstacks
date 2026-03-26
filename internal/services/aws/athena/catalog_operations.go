package athena

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// ListEngineVersions retrieves the list of available Athena engine versions.
func (s *Service) ListEngineVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"EngineVersions": []map[string]interface{}{
			{
				"SelectedEngineVersion":  "AUTO",
				"EffectiveEngineVersion": "Athena engine version 3",
			},
			{
				"SelectedEngineVersion":  "Athena engine version 3",
				"EffectiveEngineVersion": "Athena engine version 3",
			},
			{
				"SelectedEngineVersion":  "Athena engine version 2",
				"EffectiveEngineVersion": "Athena engine version 2",
			},
		},
		"NextToken": "",
	}, nil
}

// ListDataCatalogs retrieves a list of all data catalogs in the Athena workgroup.
func (s *Service) ListDataCatalogs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	catalogs, err := stores.dataCatalogStore.ListDataCatalogs()
	if err != nil {
		return nil, err
	}

	var summaries []map[string]interface{}
	for _, c := range catalogs {
		summaries = append(summaries, map[string]interface{}{
			"CatalogName": c.Name,
			"Type":        c.Type,
		})
	}

	summaries = append([]map[string]interface{}{
		{
			"CatalogName": "AwsDataCatalog",
			"Type":        "GLUE",
		},
	}, summaries...)

	return map[string]interface{}{
		"DataCatalogsSummary": summaries,
		"NextToken":           "",
	}, nil
}

// GetDataCatalog retrieves metadata for the specified data catalog.
func (s *Service) GetDataCatalog(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	if name == "AwsDataCatalog" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	catalog, err := stores.dataCatalogStore.GetDataCatalog(name)
	if err != nil {
		if err == athenastore.ErrDataCatalogNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return map[string]interface{}{
		"DataCatalog": s.dataCatalogToResponse(catalog),
	}, nil
}

// CreateDataCatalog creates a new data catalog in the Athena workgroup.
func (s *Service) CreateDataCatalog(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
	catalogType := request.GetParamCaseInsensitive(req.Parameters, "Type")
	if catalogType == "" {
		catalogType = "GLUE"
	}

	parametersRaw := request.GetMapParamCaseInsensitive(req.Parameters, "Parameters")
	parameters := convertMapToStringMap(parametersRaw)

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	catalog := &athenastore.DataCatalog{
		Name:        name,
		Description: description,
		Type:        catalogType,
		Parameters:  parameters,
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.dataCatalogStore.CreateDataCatalog(catalog); err != nil {
		if err == athenastore.ErrDataCatalogAlreadyExists {
			return nil, ErrResourceAlreadyExistsException
		}
		return nil, err
	}

	if len(tags) > 0 {
		_ = stores.dataCatalogStore.TagResource(name, tags)
	}

	return map[string]interface{}{
		"DataCatalog": s.dataCatalogToResponse(catalog),
	}, nil
}

// DeleteDataCatalog deletes the specified data catalog.
func (s *Service) DeleteDataCatalog(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	if name == "AwsDataCatalog" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.dataCatalogStore.DeleteDataCatalog(name); err != nil {
		if err == athenastore.ErrDataCatalogNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateDataCatalog updates the specified data catalog with new metadata.
func (s *Service) UpdateDataCatalog(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	if name == "AwsDataCatalog" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	catalog, err := stores.dataCatalogStore.GetDataCatalog(name)
	if err != nil {
		if err == athenastore.ErrDataCatalogNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
	if description != "" {
		catalog.Description = description
	}

	parametersRaw := request.GetMapParamCaseInsensitive(req.Parameters, "Parameters")
	if parametersRaw != nil {
		catalog.Parameters = convertMapToStringMap(parametersRaw)
	}

	if err := stores.dataCatalogStore.UpdateDataCatalog(catalog); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DataCatalog": s.dataCatalogToResponse(catalog),
	}, nil
}

// ListDatabases retrieves a list of databases in the specified data catalog.
func (s *Service) ListDatabases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	catalogName := request.GetParamCaseInsensitive(req.Parameters, "CatalogName")
	if catalogName == "" {
		catalogName = "AwsDataCatalog"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	databases, err := stores.databaseStore.ListDatabases(catalogName)
	if err != nil {
		return nil, err
	}

	var dbList []map[string]interface{}
	for _, db := range databases {
		dbList = append(dbList, map[string]interface{}{
			"Name":        db.Name,
			"Description": db.Description,
		})
	}

	if catalogName == "AwsDataCatalog" && len(dbList) == 0 {
		dbList = append(dbList, map[string]interface{}{
			"Name":        "default",
			"Description": "Default database",
		})
	}

	return map[string]interface{}{
		"DatabaseList": dbList,
		"NextToken":    "",
	}, nil
}

// GetDatabase retrieves metadata for the specified database.
func (s *Service) GetDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	catalogName := request.GetParamCaseInsensitive(req.Parameters, "CatalogName")
	if catalogName == "" {
		catalogName = "AwsDataCatalog"
	}

	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidRequestException
	}

	if catalogName == "AwsDataCatalog" && databaseName == "default" {
		return map[string]interface{}{
			"Database": map[string]interface{}{
				"Name":        "default",
				"Description": "Default database",
				"Parameters": map[string]string{
					"EXTERNAL": "TRUE",
				},
			},
		}, nil
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	db, err := stores.databaseStore.GetDatabase(catalogName, databaseName)
	if err != nil {
		if err == athenastore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return map[string]interface{}{
		"Database": map[string]interface{}{
			"Name":        db.Name,
			"Description": db.Description,
			"Parameters":  db.Parameters,
		},
	}, nil
}

// ListTableMetadata retrieves metadata for all tables in the specified database.
func (s *Service) ListTableMetadata(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	catalogName := request.GetParamCaseInsensitive(req.Parameters, "CatalogName")
	if catalogName == "" {
		catalogName = "AwsDataCatalog"
	}

	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tables, err := stores.tableStore.ListTables(catalogName, databaseName)
	if err != nil {
		return nil, err
	}

	var tableList []map[string]interface{}
	for _, t := range tables {
		tableList = append(tableList, map[string]interface{}{
			"Name":       t.Name,
			"TableType":  t.TableType,
			"CreateTime": nil,
		})
	}

	return map[string]interface{}{
		"TableMetadataList": tableList,
		"NextToken":         "",
	}, nil
}

// GetTableMetadata retrieves metadata for the specified table.
func (s *Service) GetTableMetadata(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	catalogName := request.GetParamCaseInsensitive(req.Parameters, "CatalogName")
	if catalogName == "" {
		catalogName = "AwsDataCatalog"
	}

	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidRequestException
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
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
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return map[string]interface{}{
		"TableMetadata": s.tableMetadataToResponse(table),
	}, nil
}

func (s *Service) dataCatalogToResponse(dc *athenastore.DataCatalog) map[string]interface{} {
	return map[string]interface{}{
		"Name":        dc.Name,
		"Description": dc.Description,
		"Type":        dc.Type,
		"Parameters":  dc.Parameters,
	}
}

func (s *Service) tableMetadataToResponse(t *athenastore.TableMetadata) map[string]interface{} {
	var columns []map[string]interface{}
	for _, c := range t.Columns {
		columns = append(columns, map[string]interface{}{
			"Name":    c.Name,
			"Type":    c.Type,
			"Comment": c.Comment,
		})
	}

	var partitionKeys []map[string]interface{}
	for _, p := range t.PartitionKeys {
		partitionKeys = append(partitionKeys, map[string]interface{}{
			"Name":    p.Name,
			"Type":    p.Type,
			"Comment": p.Comment,
		})
	}

	return map[string]interface{}{
		"Name":          t.Name,
		"DatabaseName":  t.DatabaseName,
		"Description":   t.Description,
		"TableType":     t.TableType,
		"Columns":       columns,
		"PartitionKeys": partitionKeys,
		"Parameters":    t.Parameters,
	}
}

func convertMapToStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if vs, ok := v.(string); ok {
			result[k] = vs
		} else {
			result[k] = ""
		}
	}
	return result
}
