package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/types"
)

// CreateLogGroup creates a new CloudWatch Logs log group.
func (s *LogsService) CreateLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	lg := logsstore.NewLogGroup(logGroupName, reqCtx.GetRegion(), s.accountID)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.CreateLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	if len(tags) > 0 {
		if err := store.Tags().Tag(lg.ARN, tags); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return response.EmptyResponse(), nil
}

// DeleteLogGroup deletes a CloudWatch Logs log group.
func (s *LogsService) DeleteLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteLogGroup(logGroupName); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeLogGroups returns a list of CloudWatch Logs log groups.
func (s *LogsService) DescribeLogGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	prefix := request.GetParamLowerFirst(req.Parameters, "LogGroupNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 50
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	groups, nextMarker, err := store.ListLogGroups(prefix, nextToken, int(limit))
	if err != nil {
		return nil, mapStoreError(err)
	}

	logGroups := make([]map[string]interface{}, 0)
	for _, lg := range groups {
		logGroups = append(logGroups, map[string]interface{}{
			"logGroupName":         lg.Name,
			"arn":                  lg.ARN + ":*",
			"creationTime":         lg.CreatedAt.UnixMilli(),
			"retentionInDays":      lg.RetentionInDays,
			"metricFilterCount":    lg.MetricFilterCount,
			"storedBytes":          lg.StoredBytes,
			"logGroupArn":          lg.ARN,
			"kmsKeyId":             "",
			"dataProtectionStatus": "",
		})
	}

	response := map[string]interface{}{
		"logGroups": logGroups,
	}
	if nextMarker != "" {
		response["nextToken"] = nextMarker
	}

	return response, nil
}

// ListLogGroups returns a list of CloudWatch Logs log groups.
func (s *LogsService) ListLogGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.DescribeLogGroups(ctx, reqCtx, req)
}

// PutRetentionPolicy sets the retention policy for a CloudWatch Logs log group.
func (s *LogsService) PutRetentionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	retentionInDays := int32(request.GetIntParam(req.Parameters, "retentionInDays"))
	if retentionInDays == 0 {
		retentionInDays = int32(request.GetIntParam(req.Parameters, "RetentionInDays"))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	lg.SetRetention(retentionInDays)
	if err := store.PutLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DeleteRetentionPolicy deletes the retention policy for a CloudWatch Logs log group.
func (s *LogsService) DeleteRetentionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	lg.SetRetention(0)
	if err := store.PutLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// TagResource adds tags to a CloudWatch Logs resource.
func (s *LogsService) tagHandlerConfig(store *logsstore.Store) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:      "ResourceArn",
			TagsParam:          "Tags",
			TagKeysParam:       "TagKeys",
			TagKeyName:         "Key",
			TagValueName:       "Value",
			RequireTags:        false,
			RequireTagKeys:     false,
			RequireResource:    true,
			CaseInsensitiveRes: true,
		},
		ParseTags: func(params map[string]interface{}) []types.Tag {
			return tagutil.MapToTags(tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(params, "Tags")))
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return request.GetStringList(params, "TagKeys")
		},
		TagFunc: func(_ context.Context, resourceKey string, tagSlice []types.Tag) error {
			return store.Tags().TagFromSlice(resourceKey, tagSlice)
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.Tags().Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			return store.Tags().ListAsSlice(resourceKey)
		},
		FormatResponse: func(tagSlice []types.Tag, _ string) (interface{}, error) {
			m := tagutil.ToMap(tagSlice)
			if m == nil {
				m = make(map[string]string)
			}
			return map[string]interface{}{
				"tags": m,
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: mapStoreError,
	}
}

// TagResource adds tags to a CloudWatch Logs log group.
func (s *LogsService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.tagHandlerConfig(store))
}

// UntagResource removes tags from a CloudWatch Logs resource.
func (s *LogsService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.tagHandlerConfig(store))
}

// ListTagsForResource lists the tags for a CloudWatch Logs resource.
func (s *LogsService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.tagHandlerConfig(store))
}

// TagLogGroup adds tags to the specified CloudWatch Logs log group.
func (s *LogsService) TagLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if err := store.Tags().Tag(lg.ARN, tags); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTagsLogGroup retrieves the tags for the specified CloudWatch Logs log group.
func (s *LogsService) ListTagsLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	tags, err := store.Tags().List(lg.ARN)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return map[string]interface{}{
		"tags": tags,
	}, nil
}
