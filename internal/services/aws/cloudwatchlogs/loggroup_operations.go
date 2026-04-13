package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/internal/utils/aws/arn"
)

// CreateLogGroup creates a new CloudWatch Logs log group.
func (s *LogsService) CreateLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	if logGroupName == "" {
		return nil, ErrMissingParameter
	}

	lg := logsstore.NewLogGroup(logGroupName, reqCtx.GetRegion(), s.accountID)

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))
	if len(tags) > 0 {
		lg.Tags = tags
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.CreateLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
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
func (s *LogsService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	if resourceARN == "" {
		return nil, ErrMissingParameter
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	logGroupName := arn.ExtractLogGroupNameFromARN(resourceARN)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if lg.Tags == nil {
		lg.Tags = make(map[string]string)
	}
	for k, v := range tags {
		lg.Tags[k] = v
	}

	if err := store.PutLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a CloudWatch Logs resource.
func (s *LogsService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	if resourceARN == "" {
		return nil, ErrMissingParameter
	}

	tagKeys := request.GetStringList(req.Parameters, "TagKeys")

	logGroupName := arn.ExtractLogGroupNameFromARN(resourceARN)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if lg.Tags == nil {
		lg.Tags = make(map[string]string)
	}
	for _, k := range tagKeys {
		delete(lg.Tags, k)
	}

	if err := store.PutLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists the tags for a CloudWatch Logs resource.
func (s *LogsService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	if resourceARN == "" {
		return nil, ErrMissingParameter
	}

	logGroupName := arn.ExtractLogGroupNameFromARN(resourceARN)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	tags := lg.Tags
	if tags == nil {
		tags = make(map[string]string)
	}
	return map[string]interface{}{
		"tags": tags,
	}, nil
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

	if lg.Tags == nil {
		lg.Tags = make(map[string]string)
	}
	for k, v := range tags {
		lg.Tags[k] = v
	}

	if err := store.PutLogGroup(lg); err != nil {
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

	return map[string]interface{}{
		"tags": lg.Tags,
	}, nil
}
