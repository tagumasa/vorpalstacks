package timestreamwrite

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// CreateTable creates a new Timestream table.
func (s *TimestreamWriteService) CreateTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidParameter
	}

	if !isValidTimestreamName(databaseName) {
		return nil, ErrValidationException
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	if !isValidTimestreamName(tableName) {
		return nil, ErrValidationException
	}

	retentionProperties := s.parseRetentionProperties(req.Parameters["RetentionProperties"])
	schema := s.parseSchema(req.Parameters["Schema"])

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := st.tableStore.CreateTable(databaseName, tableName, retentionProperties, schema)
	if err != nil {
		if err == tsstore.ErrTableAlreadyExists {
			return nil, ErrResourceAlreadyExists
		}
		if err == tsstore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	if tags := tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"); len(tags) > 0 {
		_ = st.store.TagResource(table.ARN, tagutil.ToMap(tags))
	}

	tags, _ := st.store.ListTags(table.ARN)

	return map[string]interface{}{
		"Table": s.formatTableResponse(table, tags),
	}, nil
}

// DescribeTable returns information about a Timestream table.
func (s *TimestreamWriteService) DescribeTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidParameter
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := st.tableStore.GetTable(databaseName, tableName)
	if err != nil {
		if err == tsstore.ErrTableNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	tags, _ := st.store.ListTags(table.ARN)

	return map[string]interface{}{
		"Table": s.formatTableResponse(table, tags),
	}, nil
}

// ListTables returns a list of Timestream tables in a database.
func (s *TimestreamWriteService) ListTables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	maxResults := request.GetIntParam(req.Parameters, "MaxResults")
	if maxResults <= 0 {
		maxResults = 20
	}

	opts := common.ListOptions{MaxItems: maxResults}
	if nextToken != "" {
		opts.Marker = nextToken
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := st.tableStore.ListTables(databaseName, opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	tableList := make([]map[string]interface{}, 0)
	for _, table := range result.Items {
		tags, _ := st.store.ListTags(table.ARN)
		tableList = append(tableList, s.formatTableResponse(table, tags))
	}

	response := map[string]interface{}{
		"Tables": tableList,
	}
	if result.NextMarker != "" {
		response["NextToken"] = result.NextMarker
	}
	return response, nil
}

// UpdateTable modifies the retention properties or schema of a Timestream table.
func (s *TimestreamWriteService) UpdateTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidParameter
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	retentionProperties := s.parseRetentionProperties(req.Parameters["RetentionProperties"])
	schema := s.parseSchema(req.Parameters["Schema"])

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := st.tableStore.UpdateTable(databaseName, tableName, retentionProperties, schema)
	if err != nil {
		if err == tsstore.ErrTableNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	tags, _ := st.store.ListTags(table.ARN)

	return map[string]interface{}{
		"Table": s.formatTableResponse(table, tags),
	}, nil
}

// DeleteTable deletes a Timestream table.
func (s *TimestreamWriteService) DeleteTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidParameter
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	err = st.tableStore.DeleteTable(databaseName, tableName)
	if err != nil {
		if err == tsstore.ErrTableNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

func (s *TimestreamWriteService) parseRetentionProperties(data interface{}) *tsstore.RetentionProperties {
	if data == nil {
		return nil
	}

	propsMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	props := &tsstore.RetentionProperties{
		MemoryStoreRetentionPeriodInHours:  6,
		MagneticStoreRetentionPeriodInDays: 7300,
	}

	if v, ok := propsMap["MemoryStoreRetentionPeriodInHours"]; ok {
		if f, ok := v.(float64); ok {
			props.MemoryStoreRetentionPeriodInHours = int64(f)
		}
	} else if v, ok := propsMap["memoryStoreRetentionPeriodInHours"]; ok {
		if f, ok := v.(float64); ok {
			props.MemoryStoreRetentionPeriodInHours = int64(f)
		}
	}

	if v, ok := propsMap["MagneticStoreRetentionPeriodInDays"]; ok {
		if f, ok := v.(float64); ok {
			props.MagneticStoreRetentionPeriodInDays = int64(f)
		}
	} else if v, ok := propsMap["magneticStoreRetentionPeriodInDays"]; ok {
		if f, ok := v.(float64); ok {
			props.MagneticStoreRetentionPeriodInDays = int64(f)
		}
	}

	return props
}

func (s *TimestreamWriteService) parseSchema(data interface{}) *tsstore.Schema {
	if data == nil {
		return nil
	}

	schemaMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}

	schema := &tsstore.Schema{}

	var cpk interface{}
	if cpk, ok = schemaMap["CompositePartitionKey"]; !ok {
		cpk, ok = schemaMap["compositePartitionKey"]
	}
	if ok {
		if cpkList, ok := cpk.([]interface{}); ok {
			for _, item := range cpkList {
				if itemMap, ok := item.(map[string]interface{}); ok {
					pk := tsstore.PartitionKey{}
					if t, ok := itemMap["Type"].(string); ok {
						pk.Type = tsstore.PartitionKeyType(t)
					} else if t, ok := itemMap["type"].(string); ok {
						pk.Type = tsstore.PartitionKeyType(t)
					}
					if n, ok := itemMap["Name"].(string); ok {
						pk.Name = n
					} else if n, ok := itemMap["name"].(string); ok {
						pk.Name = n
					}
					if e, ok := itemMap["EnforcementInRecord"].(string); ok {
						pk.EnforcementInRecord = tsstore.EnforcementInRecord(e)
					} else if e, ok := itemMap["enforcementInRecord"].(string); ok {
						pk.EnforcementInRecord = tsstore.EnforcementInRecord(e)
					}
					schema.CompositePartitionKey = append(schema.CompositePartitionKey, pk)
				}
			}
		}
	}

	return schema
}

func (s *TimestreamWriteService) formatTableResponse(table *tsstore.Table, tags map[string]string) map[string]interface{} {
	response := map[string]interface{}{
		"Arn":             table.ARN,
		"TableName":       table.TableName,
		"DatabaseName":    table.DatabaseName,
		"TableStatus":     table.TableStatus,
		"CreationTime":    float64(table.CreationTime.Unix()) + float64(table.CreationTime.Nanosecond())/1e9,
		"LastUpdatedTime": float64(table.LastUpdatedTime.Unix()) + float64(table.LastUpdatedTime.Nanosecond())/1e9,
	}

	if table.RetentionProperties != nil {
		response["RetentionProperties"] = map[string]interface{}{
			"MemoryStoreRetentionPeriodInHours":  table.RetentionProperties.MemoryStoreRetentionPeriodInHours,
			"MagneticStoreRetentionPeriodInDays": table.RetentionProperties.MagneticStoreRetentionPeriodInDays,
		}
	}

	if table.Schema != nil && len(table.Schema.CompositePartitionKey) > 0 {
		var cpk []map[string]interface{}
		for _, pk := range table.Schema.CompositePartitionKey {
			pkMap := map[string]interface{}{
				"Type": string(pk.Type),
			}
			if pk.Name != "" {
				pkMap["Name"] = pk.Name
			}
			if pk.EnforcementInRecord != "" {
				pkMap["EnforcementInRecord"] = string(pk.EnforcementInRecord)
			}
			cpk = append(cpk, pkMap)
		}
		response["Schema"] = map[string]interface{}{
			"CompositePartitionKey": cpk,
		}
	}

	if len(tags) > 0 {
		response["Tags"] = tagutil.MapToResponse(tags)
	}

	return response
}
