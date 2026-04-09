package timestreamwrite

import (
	"context"
	"regexp"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	"vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

var timestreamNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]{2,255}$`)

func isValidTimestreamName(name string) bool {
	return timestreamNameRegex.MatchString(name)
}

// DescribeEndpoints returns information about the Timestream Write endpoints.
func (s *Service) DescribeEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
func (s *Service) CreateDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidParameter
	}

	if !isValidTimestreamName(databaseName) {
		return nil, ErrValidationException
	}

	kmsKeyID := request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	db, err := st.store.CreateDatabase(databaseName, kmsKeyID)
	if err != nil {
		if err == tsstore.ErrDatabaseAlreadyExists {
			return nil, ErrResourceAlreadyExists
		}
		return nil, s.mapStoreError(err)
	}

	if tags := tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"); len(tags) > 0 {
		if err := st.store.TagResource(db.ARN, tagutil.ToMap(tags)); err != nil {
			return nil, s.mapStoreError(err)
		}
	}

	tags, _ := st.store.ListTags(db.ARN)

	return map[string]interface{}{
		"Database": s.formatDatabaseResponse(db, tags),
	}, nil
}

// DescribeDatabase returns information about a Timestream database.
func (s *Service) DescribeDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidParameter
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	db, err := st.store.GetDatabase(databaseName)
	if err != nil {
		if err == tsstore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	tags, _ := st.store.ListTags(db.ARN)

	return map[string]interface{}{
		"Database": s.formatDatabaseResponse(db, tags),
	}, nil
}

// ListDatabases returns a list of Timestream databases.
func (s *Service) ListDatabases(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	result, err := st.store.ListDatabases(opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	dbList := make([]map[string]interface{}, 0)
	for _, db := range result.Items {
		tags, _ := st.store.ListTags(db.ARN)
		dbList = append(dbList, s.formatDatabaseResponse(db, tags))
	}

	response := map[string]interface{}{
		"Databases": dbList,
	}
	if result.NextMarker != "" {
		response["NextToken"] = result.NextMarker
	}
	return response, nil
}

// UpdateDatabase modifies the KMS key for a Timestream database.
func (s *Service) UpdateDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidParameter
	}

	kmsKeyID := request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	db, err := st.store.UpdateDatabase(databaseName, kmsKeyID)
	if err != nil {
		if err == tsstore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	tags, _ := st.store.ListTags(db.ARN)

	return map[string]interface{}{
		"Database": s.formatDatabaseResponse(db, tags),
	}, nil
}

// DeleteDatabase deletes a Timestream database.
func (s *Service) DeleteDatabase(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrInvalidParameter
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	err = st.store.DeleteDatabase(databaseName)
	if err != nil {
		if err == tsstore.ErrDatabaseNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

func (s *Service) formatDatabaseResponse(db *tsstore.Database, tags map[string]string) map[string]interface{} {
	response := map[string]interface{}{
		"Arn":             db.ARN,
		"DatabaseName":    db.DatabaseName,
		"TableCount":      db.TableCount,
		"CreationTime":    float64(db.CreationTime.Unix()) + float64(db.CreationTime.Nanosecond())/1e9,
		"LastUpdatedTime": float64(db.LastUpdatedTime.Unix()) + float64(db.LastUpdatedTime.Nanosecond())/1e9,
	}

	if db.KmsKeyId != "" {
		response["KmsKeyId"] = db.KmsKeyId
	}

	if len(tags) > 0 {
		response["Tags"] = tagutil.MapToResponse(tags)
	}

	return response
}
