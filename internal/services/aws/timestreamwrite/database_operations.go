package timestreamwrite

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/aws/types"
)

// DescribeEndpoints returns information about the Timestream Write endpoints.
func (s *TimestreamWriteService) DescribeEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"Endpoints": []map[string]interface{}{
			{
				"Address":              s.serverHost,
				"CachePeriodInMinutes": 1440,
			},
		},
	}, nil
}

// CreateDatabase creates a new Timestream database.
func (s *TimestreamWriteService) CreateDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	kmsKeyID := request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId")

	var tags []types.Tag
	tagsProvided := false
	if parsedTags := tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"); len(parsedTags) > 0 {
		tags = parsedTags
		tagsProvided = true
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createDatabaseCore(ctx, stores, CreateDatabaseInput{
		DatabaseName: databaseName,
		KmsKeyId:     kmsKeyID,
		Tags:         tags,
		TagsProvided: tagsProvided,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Database": formatDatabaseResponse(result),
	}, nil
}

// DescribeDatabase returns information about a Timestream database.
func (s *TimestreamWriteService) DescribeDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeDatabaseCore(ctx, stores, DescribeDatabaseInput{
		DatabaseName: databaseName,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Database": formatDatabaseResponse(result),
	}, nil
}

// ListDatabases returns a list of Timestream databases.
func (s *TimestreamWriteService) ListDatabases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, maxListDatabasesResults, "MaxResults")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listDatabasesCore(ctx, stores, ListDatabasesInput{
		NextToken: nextToken,
		MaxItems:  maxResults,
	})
	if err != nil {
		return nil, err
	}

	dbList := make([]map[string]interface{}, 0, len(result.Databases))
	for i := range result.Databases {
		dbList = append(dbList, formatDatabaseResponse(&result.Databases[i]))
	}

	resp := map[string]interface{}{
		"Databases": dbList,
	}
	pagination.SetNextToken(resp, "NextToken", result.NextToken)
	return resp, nil
}

// UpdateDatabase modifies the KMS key for a Timestream database.
func (s *TimestreamWriteService) UpdateDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	kmsKeyID := request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.updateDatabaseCore(ctx, stores, UpdateDatabaseInput{
		DatabaseName: databaseName,
		KmsKeyId:     kmsKeyID,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Database": formatDatabaseResponse(result),
	}, nil
}

// DeleteDatabase deletes a Timestream database.
func (s *TimestreamWriteService) DeleteDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteDatabaseCore(ctx, stores, DeleteDatabaseInput{
		DatabaseName: databaseName,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// formatDatabaseResponse converts a DatabaseResult to the HTTP API JSON
// response format (epoch float64 timestamps).
func formatDatabaseResponse(db *DatabaseResult) map[string]interface{} {
	resp := map[string]interface{}{
		"Arn":             db.ARN,
		"DatabaseName":    db.DatabaseName,
		"TableCount":      db.TableCount,
		"CreationTime":    float64(db.CreationTime.Unix()) + float64(db.CreationTime.Nanosecond())/1e9,
		"LastUpdatedTime": float64(db.LastUpdatedTime.Unix()) + float64(db.LastUpdatedTime.Nanosecond())/1e9,
	}

	if db.KmsKeyId != "" {
		resp["KmsKeyId"] = db.KmsKeyId
	}

	if len(db.Tags) > 0 {
		resp["Tags"] = tagutil.MapToResponse(tagutil.ToMap(db.Tags))
	}

	return resp
}
