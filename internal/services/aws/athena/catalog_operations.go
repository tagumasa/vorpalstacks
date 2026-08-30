package athena

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// ListEngineVersions retrieves the list of available Athena engine versions.
// The list is service-defined (not user-configurable) and matches the AWS
// ListEngineVersions API output as of Athena engine version 3.
// Verified against https://docs.aws.amazon.com/athena/latest/ug/engine-versions-changing.html
func (s *AthenaService) ListEngineVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	}, nil
}

// ListDataCatalogs retrieves a list of all data catalogs in the Athena workgroup.
func (s *AthenaService) ListDataCatalogs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := ListDataCatalogsInput{
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	catalogs, nextMarker, err := listDataCatalogsCore(stores, input)
	if err != nil {
		return nil, err
	}

	summaries := make([]map[string]interface{}, 0, len(catalogs))
	for _, c := range catalogs {
		summaries = append(summaries, map[string]interface{}{
			"CatalogName": c.Name,
			"Type":        c.Type,
		})
	}

	return pagination.BuildListResponse("DataCatalogsSummary", summaries, nextMarker), nil
}

// GetDataCatalog retrieves metadata for the specified data catalog.
func (s *AthenaService) GetDataCatalog(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")

	catalog, err := s.getDataCatalogCore(reqCtx, name)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DataCatalog": s.dataCatalogToResponse(catalog),
	}, nil
}

// CreateDataCatalog creates a new data catalog in the Athena workgroup.
func (s *AthenaService) CreateDataCatalog(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	parametersRaw := request.GetMapParamCaseInsensitive(req.Parameters, "Parameters")
	input := CreateDataCatalogInput{
		Name:        request.GetParamCaseInsensitive(req.Parameters, "Name"),
		Description: request.GetParamCaseInsensitive(req.Parameters, "Description"),
		Type:        request.GetParamCaseInsensitive(req.Parameters, "Type"),
		Parameters:  convertMapToStringMap(parametersRaw),
		Tags:        tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")),
	}

	catalog, err := s.createDataCatalogCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DataCatalog": s.dataCatalogToResponse(catalog),
	}, nil
}

// DeleteDataCatalog deletes the specified data catalog.
func (s *AthenaService) DeleteDataCatalog(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")

	if err := s.deleteDataCatalogCore(reqCtx, name); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateDataCatalog updates the specified data catalog with new metadata.
// Per the Smithy model, Name and Type are both REQUIRED on this operation.
func (s *AthenaService) UpdateDataCatalog(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	parametersRaw := request.GetMapParamCaseInsensitive(req.Parameters, "Parameters")
	input := UpdateDataCatalogInput{
		Name:          request.GetParamCaseInsensitive(req.Parameters, "Name"),
		Description:   request.GetParamCaseInsensitive(req.Parameters, "Description"),
		Type:          request.GetParamCaseInsensitive(req.Parameters, "Type"),
		Parameters:    convertMapToStringMap(parametersRaw),
		HasParameters: parametersRaw != nil,
	}

	catalog, err := s.updateDataCatalogCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DataCatalog": s.dataCatalogToResponse(catalog),
	}, nil
}

// ListDatabases retrieves a list of databases in the specified data catalog.
func (s *AthenaService) ListDatabases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := ListDatabasesInput{
		CatalogName:   request.GetParamCaseInsensitive(req.Parameters, "CatalogName"),
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	databases, nextMarker, err := s.listDatabasesCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	dbList := make([]map[string]interface{}, 0, len(databases))
	for _, db := range databases {
		dbList = append(dbList, map[string]interface{}{
			"Name":        db.Name,
			"Description": db.Description,
		})
	}

	return pagination.BuildListResponse("DatabaseList", dbList, nextMarker), nil
}

// GetDatabase retrieves metadata for the specified database.
func (s *AthenaService) GetDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	catalogName := request.GetParamCaseInsensitive(req.Parameters, "CatalogName")
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")

	db, err := s.getDatabaseCore(reqCtx, catalogName, databaseName)
	if err != nil {
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
func (s *AthenaService) ListTableMetadata(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := ListTableMetadataInput{
		CatalogName:   request.GetParamCaseInsensitive(req.Parameters, "CatalogName"),
		DatabaseName:  request.GetParamCaseInsensitive(req.Parameters, "DatabaseName"),
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	tables, nextMarker, err := s.listTableMetadataCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	tableList := make([]map[string]interface{}, 0, len(tables))
	for _, t := range tables {
		tableList = append(tableList, map[string]interface{}{
			"Name":       t.Name,
			"TableType":  t.TableType,
			"CreateTime": nil,
		})
	}

	return pagination.BuildListResponse("TableMetadataList", tableList, nextMarker), nil
}

// GetTableMetadata retrieves metadata for the specified table.
func (s *AthenaService) GetTableMetadata(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	catalogName := request.GetParamCaseInsensitive(req.Parameters, "CatalogName")
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")

	table, err := s.getTableMetadataCore(reqCtx, catalogName, databaseName, tableName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableMetadata": s.tableMetadataToResponse(table),
	}, nil
}

func (s *AthenaService) dataCatalogToResponse(dc *athenastore.DataCatalog) map[string]interface{} {
	return map[string]interface{}{
		"Name":        dc.Name,
		"Description": dc.Description,
		"Type":        dc.Type,
		"Parameters":  dc.Parameters,
	}
}

func (s *AthenaService) tableMetadataToResponse(t *athenastore.TableMetadata) map[string]interface{} {
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
