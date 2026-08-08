package cloudwatchlogs

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/types"
)

// CreateLogGroup creates a new CloudWatch Logs log group.
func (s *LogsService) CreateLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	input := CreateLogGroupInput{
		LogGroupName:              logGroupName,
		KmsKeyId:                  request.GetParamLowerFirst(req.Parameters, "KmsKeyId"),
		LogGroupClass:             request.GetParamLowerFirst(req.Parameters, "LogGroupClass"),
		Tags:                      tags,
		DeletionProtectionEnabled: request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled"),
		Region:                    reqCtx.GetRegion(),
	}

	if _, err := s.createLogGroupCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteLogGroup deletes a CloudWatch Logs log group.
func (s *LogsService) DeleteLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := DeleteLogGroupInput{
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		Region:       reqCtx.GetRegion(),
	}

	if err := s.deleteLogGroupCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeLogGroups returns a list of CloudWatch Logs log groups.
func (s *LogsService) DescribeLogGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	prefix := request.GetParamLowerFirst(req.Parameters, "LogGroupNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	logGroupClass := request.GetParamLowerFirst(req.Parameters, "LogGroupClass")

	limit, err := validateListLimit(limit, 50, 50)
	if err != nil {
		return nil, err
	}

	input := ListLogGroupsInput{
		LogGroupNamePrefix: prefix,
		LogGroupClass:      logGroupClass,
		NextToken:          nextToken,
		Limit:              limit,
		Region:             reqCtx.GetRegion(),
	}

	result, err := s.listLogGroupsCore(input)
	if err != nil {
		return nil, err
	}

	logGroups := make([]map[string]interface{}, 0)
	for _, lg := range result.LogGroups {
		entry := map[string]interface{}{
			"logGroupName":      lg.Name,
			"arn":               lg.ARN + ":*",
			"creationTime":      lg.CreatedAt.UnixMilli(),
			"metricFilterCount": lg.MetricFilterCount,
			"storedBytes":       lg.StoredBytes,
			"logGroupArn":       lg.ARN,
			"logGroupClass":     lg.LogGroupClass,
		}
		if lg.RetentionInDays > 0 {
			entry["retentionInDays"] = lg.RetentionInDays
		}
		if lg.KmsKeyId != "" {
			entry["kmsKeyId"] = lg.KmsKeyId
		}
		if lg.DeletionProtectionEnabled {
			entry["deletionProtectionEnabled"] = lg.DeletionProtectionEnabled
		}
		logGroups = append(logGroups, entry)
	}

	resp := map[string]interface{}{
		"logGroups": logGroups,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

// ListLogGroups returns a list of CloudWatch Logs log groups.
func (s *LogsService) ListLogGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	prefix := request.GetParamLowerFirst(req.Parameters, "logGroupNamePattern")
	if prefix == "" {
		prefix = request.GetParamLowerFirst(req.Parameters, "LogGroupNamePattern")
	}
	if prefix != "" && len(prefix) > 1 && prefix[0] == '^' {
		prefix = prefix[1:]
		if idx := strings.Index(prefix, "$"); idx >= 0 {
			prefix = prefix[:idx]
		}
	}
	if prefix == "" {
		prefix = request.GetParamLowerFirst(req.Parameters, "LogGroupNamePrefix")
	}

	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))

	limit, err := validateListLimit(limit, 50, 1000)
	if err != nil {
		return nil, err
	}

	input := ListLogGroupsInput{
		LogGroupNamePrefix: prefix,
		NextToken:          nextToken,
		Limit:              limit,
		Region:             reqCtx.GetRegion(),
	}

	result, err := s.listLogGroupsCore(input)
	if err != nil {
		return nil, err
	}

	logGroups := make([]map[string]interface{}, 0)
	for _, lg := range result.LogGroups {
		entry := map[string]interface{}{
			"logGroupName": lg.Name,
			"logGroupArn":  lg.ARN,
		}
		logGroups = append(logGroups, entry)
	}

	resp := map[string]interface{}{
		"logGroups": logGroups,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

// PutRetentionPolicy sets the retention policy for a CloudWatch Logs log group.
func (s *LogsService) PutRetentionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if err := validateLogGroupName(logGroupName); err != nil {
		return nil, err
	}

	retentionInDays := int32(request.GetIntParam(req.Parameters, "retentionInDays"))
	if retentionInDays == 0 {
		retentionInDays = int32(request.GetIntParam(req.Parameters, "RetentionInDays"))
	}

	if !logsstore.IsValidRetentionDays(retentionInDays) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("%d is not a valid retention value. Allowed values: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653", retentionInDays),
			400)
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
	if err := validateLogGroupName(logGroupName); err != nil {
		return nil, err
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

// UntagLogGroup removes tags from the specified CloudWatch Logs log group.
// Deprecated: use UntagResource instead.
func (s *LogsService) UntagLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	tagKeys := request.GetStringList(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if err := store.Tags().Untag(lg.ARN, tagKeys); err != nil {
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
