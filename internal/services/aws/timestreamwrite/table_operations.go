package timestreamwrite

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// CreateTable creates a new Timestream table.
func (s *TimestreamWriteService) CreateTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrValidationException
	}

	if !isValidTimestreamName(databaseName) {
		return nil, ErrValidationException
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrValidationException
	}

	if !isValidTimestreamName(tableName) {
		return nil, ErrValidationException
	}

	retentionProperties, err := s.parseRetentionProperties(req.Parameters["RetentionProperties"])
	if err != nil {
		return nil, err
	}
	schema := s.parseSchema(req.Parameters["Schema"])
	magneticStoreWriteProperties, err := s.parseMagneticStoreWriteProperties(req.Parameters["MagneticStoreWriteProperties"])
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := st.tableStore.CreateTable(databaseName, tableName, retentionProperties, schema, magneticStoreWriteProperties)
	if err != nil {
		if err == tsstore.ErrTableAlreadyExists {
			return nil, ErrConflictException
		}
		if err == tsstore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	if tags := tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"); len(tags) > 0 {
		if err := st.store.Tag(table.ARN, tagutil.ToMap(tags)); err != nil {
			logs.Warn("Failed to tag newly created table", logs.String("tableArn", table.ARN), logs.Err(err))
		}
	}

	tags, _ := st.store.List(table.ARN)

	return map[string]interface{}{
		"Table": s.formatTableResponse(table, tags),
	}, nil
}

// DescribeTable returns information about a Timestream table.
func (s *TimestreamWriteService) DescribeTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrValidationException
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrValidationException
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

	tags, _ := st.store.List(table.ARN)

	return map[string]interface{}{
		"Table": s.formatTableResponse(table, tags),
	}, nil
}

// ListTables returns a list of Timestream tables in a database.
func (s *TimestreamWriteService) ListTables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrValidationException
	}
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 20, "MaxResults")
	if maxResults > maxListTablesResults {
		maxResults = maxListTablesResults
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
		tags, _ := st.store.List(table.ARN)
		tableList = append(tableList, s.formatTableResponse(table, tags))
	}

	response := map[string]interface{}{
		"Tables": tableList,
	}
	pagination.SetNextToken(response, "NextToken", result.NextMarker)
	return response, nil
}

// UpdateTable modifies the retention properties or schema of a Timestream table.
func (s *TimestreamWriteService) UpdateTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrValidationException
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrValidationException
	}

	retentionProperties, err := s.parseRetentionProperties(req.Parameters["RetentionProperties"])
	if err != nil {
		return nil, err
	}
	schema := s.parseSchema(req.Parameters["Schema"])
	magneticStoreWriteProperties, err := s.parseMagneticStoreWriteProperties(req.Parameters["MagneticStoreWriteProperties"])
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := st.tableStore.UpdateTable(databaseName, tableName, retentionProperties, schema, magneticStoreWriteProperties)
	if err != nil {
		if err == tsstore.ErrTableNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	tags, _ := st.store.List(table.ARN)

	return map[string]interface{}{
		"Table": s.formatTableResponse(table, tags),
	}, nil
}

// DeleteTable deletes a Timestream table.
func (s *TimestreamWriteService) DeleteTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrValidationException
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrValidationException
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
	st.recordStore.DeleteTableChunks(databaseName, tableName)

	return response.EmptyResponse(), nil
}

func (s *TimestreamWriteService) parseRetentionProperties(data interface{}) (*tsstore.RetentionProperties, error) {
	if data == nil {
		return nil, nil
	}

	propsMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	props := &tsstore.RetentionProperties{}

	memHours, memPresent := getRetentionPolicyValue(propsMap, "MemoryStoreRetentionPeriodInHours")
	magDays, magPresent := getRetentionPolicyValue(propsMap, "MagneticStoreRetentionPeriodInDays")

	if !memPresent || !magPresent {
		return nil, ErrValidationException
	}

	if memHours < 1 || memHours > 8766 {
		return nil, ErrValidationException
	}

	if magDays < 1 || magDays > 73000 {
		return nil, ErrValidationException
	}

	props.MemoryStoreRetentionPeriodInHours = memHours
	props.MagneticStoreRetentionPeriodInDays = magDays

	return props, nil
}

func getRetentionPolicyValue(propsMap map[string]interface{}, key string) (int64, bool) {
	v, ok := propsMap[key]
	if !ok {
		v, ok = propsMap[strings.ToLower(key[:1])+key[1:]]
		if !ok {
			return 0, false
		}
	}
	if f, ok := v.(float64); ok {
		return int64(f), true
	}
	return 0, false
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

// parseMagneticStoreWriteProperties parses the MagneticStoreWriteProperties
// parameter from a CreateTable/UpdateTable request. EnableMagneticStoreWrites
// is REQUIRED by Smithy when MagneticStoreWriteProperties is provided.
func (s *TimestreamWriteService) parseMagneticStoreWriteProperties(data interface{}) (*tsstore.MagneticStoreWriteProperties, error) {
	if data == nil {
		return nil, nil
	}

	mMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	result := &tsstore.MagneticStoreWriteProperties{}

	enableRaw, hasEnable := mMap["EnableMagneticStoreWrites"]
	if !hasEnable {
		enableRaw, hasEnable = mMap["enableMagneticStoreWrites"]
	}
	if !hasEnable {
		return nil, ErrValidationException
	}
	if enable, ok := enableRaw.(bool); ok {
		result.EnableMagneticStoreWrites = enable
	} else {
		return nil, ErrValidationException
	}

	if rdlRaw, ok := mMap["MagneticStoreRejectedDataLocation"]; ok {
		if rdlMap, ok := rdlRaw.(map[string]interface{}); ok {
			result.MagneticStoreRejectedDataLocation = &tsstore.MagneticStoreRejectedDataLocation{}
			if s3Raw, ok := rdlMap["S3Configuration"]; ok {
				if s3Map, ok := s3Raw.(map[string]interface{}); ok {
					s3Config := &tsstore.MagneticStoreWriteS3Configuration{}
					if bucket, ok := s3Map["BucketName"].(string); ok {
						s3Config.BucketName = bucket
					}
					if prefix, ok := s3Map["ObjectKeyPrefix"].(string); ok {
						s3Config.ObjectKeyPrefix = prefix
					}
					if enc, ok := s3Map["EncryptionOption"].(string); ok {
						s3Config.EncryptionOption = enc
					}
					if kms, ok := s3Map["KmsKeyId"].(string); ok {
						s3Config.KmsKeyId = kms
					}
					result.MagneticStoreRejectedDataLocation.S3Configuration = s3Config
				}
			}
		}
	}

	return result, nil
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

	if table.MagneticStoreWriteProperties != nil {
		mswp := map[string]interface{}{
			"EnableMagneticStoreWrites": table.MagneticStoreWriteProperties.EnableMagneticStoreWrites,
		}
		if rdl := table.MagneticStoreWriteProperties.MagneticStoreRejectedDataLocation; rdl != nil {
			rdlMap := map[string]interface{}{}
			if s3 := rdl.S3Configuration; s3 != nil {
				s3Map := map[string]interface{}{}
				if s3.BucketName != "" {
					s3Map["BucketName"] = s3.BucketName
				}
				if s3.ObjectKeyPrefix != "" {
					s3Map["ObjectKeyPrefix"] = s3.ObjectKeyPrefix
				}
				if s3.EncryptionOption != "" {
					s3Map["EncryptionOption"] = s3.EncryptionOption
				}
				if s3.KmsKeyId != "" {
					s3Map["KmsKeyId"] = s3.KmsKeyId
				}
				rdlMap["S3Configuration"] = s3Map
			}
			mswp["MagneticStoreRejectedDataLocation"] = rdlMap
		}
		response["MagneticStoreWriteProperties"] = mswp
	}

	if len(tags) > 0 {
		response["Tags"] = tagutil.MapToResponse(tags)
	}

	return response
}
