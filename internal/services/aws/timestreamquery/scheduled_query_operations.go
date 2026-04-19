package timestreamquery

import (
	"context"
	"strconv"
	"time"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	tsstore "vorpalstacks/internal/store/aws/timestream"
	"vorpalstacks/internal/utils/aws/types"

	"github.com/google/uuid"
)

// CreateScheduledQuery creates a new scheduled query.
func (s *TimestreamQueryService) CreateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")
	if name == "" {
		return nil, ErrValidationException
	}

	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")
	if queryString == "" {
		return nil, ErrValidationException
	}

	scheduleConfig := s.parseScheduleConfiguration(req.Parameters)
	if scheduleConfig == nil {
		return nil, ErrValidationException
	}

	notificationConfig := s.parseNotificationConfiguration(req.Parameters)
	roleARN := request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryExecutionRoleArn")
	kmsKeyID := request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId")
	errorReportConfig := s.parseErrorReportConfiguration(req.Parameters)
	targetConfig := s.parseTargetConfiguration(req.Parameters)
	clientToken := request.GetParamCaseInsensitive(req.Parameters, "ClientToken")

	if roleARN != "" {
		validator := reqCtx.GetIAMValidator()
		if err := validator.ValidateRoleForServiceWithErrors(ctx, roleARN, iam.ServicePrincipalTimestream, &iam.RoleErrorFactories{
			RoleNotFoundError:        iam.NewTimestreamRoleError,
			RoleCannotBeAssumedError: iam.NewTimestreamRoleError,
			InvalidArnError:          iam.NewTimestreamRoleError,
		}); err != nil {
			return nil, err
		}
	}

	if clientToken == "" {
		clientToken = uuid.New().String()
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sq, err := st.scheduledQueryStore.CreateScheduledQuery(
		name,
		queryString,
		scheduleConfig,
		notificationConfig,
		roleARN,
		kmsKeyID,
		errorReportConfig,
		targetConfig,
		clientToken,
	)
	if err != nil {
		if err == tsstore.ErrScheduledQueryAlreadyExists {
			return nil, ErrResourceAlreadyExists
		}
		return nil, ErrInternalServer
	}

	if tagMap := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")); len(tagMap) > 0 {
		if tagErr := st.scheduledQueryStore.Tag(sq.ARN, tagMap); tagErr != nil {
			logs.Warn("failed to tag scheduled query", logs.Err(tagErr), logs.String("arn", sq.ARN))
		}
	}

	tags, _ := st.scheduledQueryStore.List(sq.ARN)

	return s.formatScheduledQueryResponse(sq, tags), nil
}

func (s *TimestreamQueryService) tagHandlerConfig(st *tsQueryStores) tagutil.TagHandlerConfig {
	dispatch := func(ctx context.Context, resourceARN string, fn func(ctx context.Context, resourceARN string) error) error {
		name := st.arnBuilder.Timestream().ParseScheduledQueryName(resourceARN)
		if name != "" {
			if _, err := st.scheduledQueryStore.GetScheduledQuery(name); err != nil {
				if err == tsstore.ErrScheduledQueryNotFound {
					return ErrResourceNotFound
				}
				return ErrInternalServer
			}
		}
		return fn(ctx, resourceARN)
	}

	listTags := func(ctx context.Context, resourceARN string) ([]types.Tag, error) {
		name := st.arnBuilder.Timestream().ParseScheduledQueryName(resourceARN)
		if name != "" {
			if _, err := st.scheduledQueryStore.GetScheduledQuery(name); err != nil {
				if err == tsstore.ErrScheduledQueryNotFound {
					return nil, ErrResourceNotFound
				}
				return nil, ErrInternalServer
			}
			return st.scheduledQueryStore.ListAsSlice(resourceARN)
		}
		return st.dbStore.ListAsSlice(resourceARN)
	}

	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:      "ResourceARN",
			TagsParam:          "Tags",
			TagKeysParam:       "TagKeys",
			TagKeyName:         "Key",
			TagValueName:       "Value",
			RequireTags:        true,
			RequireTagKeys:     true,
			RequireResource:    true,
			CaseInsensitiveRes: true,
		},
		ResourceKey: func(rawKey string) string { return rawKey },
		TagFunc: func(ctx context.Context, resourceKey string, tagSlice []types.Tag) error {
			return dispatch(ctx, resourceKey, func(ctx context.Context, resourceARN string) error {
				name := st.arnBuilder.Timestream().ParseScheduledQueryName(resourceARN)
				if name != "" {
					return st.scheduledQueryStore.TagFromSlice(resourceARN, tagSlice)
				}
				return st.dbStore.TagFromSlice(resourceARN, tagSlice)
			})
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return dispatch(ctx, resourceKey, func(ctx context.Context, resourceARN string) error {
				name := st.arnBuilder.Timestream().ParseScheduledQueryName(resourceARN)
				if name != "" {
					return st.scheduledQueryStore.Untag(resourceARN, tagKeys)
				}
				return st.dbStore.Untag(resourceARN, tagKeys)
			})
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			return listTags(ctx, resourceKey)
		},
		FormatResponse: func(tagSlice []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.MapToResponse(tagutil.ToMap(tagSlice)),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: func(err error) error {
			return err
		},
	}
}

// TagResource adds tags to a Timestream resource.
func (s *TimestreamQueryService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.tagHandlerConfig(st))
}

// UntagResource removes tags from a Timestream resource.
func (s *TimestreamQueryService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.tagHandlerConfig(st))
}

// DeleteScheduledQuery deletes a scheduled query.
func (s *TimestreamQueryService) DeleteScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn")
	if arn == "" {
		return nil, ErrValidationException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := st.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return nil, ErrValidationException
	}

	err = st.scheduledQueryStore.DeleteScheduledQuery(name)
	if err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	return response.EmptyResponse(), nil
}

// DescribeScheduledQuery returns the details of a scheduled query.
func (s *TimestreamQueryService) DescribeScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn")
	if arn == "" {
		return nil, ErrValidationException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := st.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return nil, ErrValidationException
	}

	sq, err := st.scheduledQueryStore.GetScheduledQuery(name)
	if err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	tags, _ := st.scheduledQueryStore.List(sq.ARN)

	return map[string]interface{}{
		"ScheduledQuery": s.formatScheduledQueryDescriptionResponse(sq, tags),
	}, nil
}

// ListScheduledQueries returns a list of scheduled queries.
func (s *TimestreamQueryService) ListScheduledQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	queries, err := st.scheduledQueryStore.ListScheduledQueries()
	if err != nil {
		return nil, ErrInternalServer
	}

	maxResults := 50
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil && val > 0 {
			maxResults = val
		}
	}

	offset := 0
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		if val, err := strconv.Atoi(nextToken); err == nil && val >= 0 {
			offset = val
		}
	}

	var scheduledQueries []map[string]interface{}
	for i, sq := range queries {
		if i < offset {
			continue
		}
		if len(scheduledQueries) >= maxResults {
			break
		}
		tags, _ := st.scheduledQueryStore.List(sq.ARN)
		scheduledQueries = append(scheduledQueries, s.formatScheduledQueryResponse(sq, tags))
	}

	response := map[string]interface{}{
		"ScheduledQueries": scheduledQueries,
	}

	if offset+maxResults < len(queries) {
		response["NextToken"] = strconv.Itoa(offset + maxResults)
	}

	return response, nil
}

// UpdateScheduledQuery updates an existing scheduled query.
func (s *TimestreamQueryService) UpdateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn")
	if arn == "" {
		return nil, ErrValidationException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := st.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return nil, ErrValidationException
	}

	state := tsstore.ScheduledQueryState(request.GetParamCaseInsensitive(req.Parameters, "State"))
	scheduleConfig := s.parseScheduleConfiguration(req.Parameters)
	notificationConfig := s.parseNotificationConfiguration(req.Parameters)
	kmsKeyID := request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId")
	errorReportConfig := s.parseErrorReportConfiguration(req.Parameters)
	targetConfig := s.parseTargetConfiguration(req.Parameters)

	_, err = st.scheduledQueryStore.UpdateScheduledQuery(
		name,
		state,
		scheduleConfig,
		notificationConfig,
		kmsKeyID,
		errorReportConfig,
		targetConfig,
	)
	if err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	return response.EmptyResponse(), nil
}

// ExecuteScheduledQuery executes a scheduled query immediately.
func (s *TimestreamQueryService) ExecuteScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn")
	if arn == "" {
		return nil, ErrValidationException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := st.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return nil, ErrValidationException
	}

	sq, err := st.scheduledQueryStore.GetScheduledQuery(name)
	if err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	now := time.Now().UTC()
	run, err := st.scheduledQueryRunStore.CreateRun(sq.ARN, now, now)
	if err != nil {
		return nil, ErrInternalServer
	}

	if err := st.scheduledQueryRunStore.UpdateRunStatus(run.ARN, tsstore.ScheduleRunStatusRunning, "", nil); err != nil {
		logs.Error("Failed to update scheduled query run to RUNNING", logs.String("arn", run.ARN), logs.Err(err))
	}

	result, execErr := s.executeSQLQuery(ctx, reqCtx, sq.QueryString)

	if execErr != nil {
		if err := st.scheduledQueryRunStore.UpdateRunStatus(run.ARN, tsstore.ScheduleRunStatusFailed, execErr.Error(), nil); err != nil {
			logs.Error("Failed to update scheduled query run to FAILED", logs.String("arn", run.ARN), logs.Err(err))
		}
		if err := st.scheduledQueryStore.UpdateLastRun(name, "FAILED", now); err != nil {
			logs.Error("Failed to update last run status for scheduled query", logs.String("name", name), logs.Err(err))
		}
		return nil, ErrQueryExecutionError
	}

	stats := &tsstore.ExecutionStats{
		QueryResultRows: int64(len(result.Rows)),
	}
	if err := st.scheduledQueryRunStore.UpdateRunStatus(run.ARN, tsstore.ScheduleRunStatusSucceeded, "", stats); err != nil {
		logs.Error("Failed to update scheduled query run to SUCCEEDED", logs.String("arn", run.ARN), logs.Err(err))
	}
	if err := st.scheduledQueryStore.UpdateLastRun(name, "SUCCESS", now); err != nil {
		logs.Error("Failed to update last run status for scheduled query", logs.String("name", name), logs.Err(err))
	}

	return map[string]interface{}{
		"ScheduledQueryRunArn": run.ARN,
	}, nil
}

func (s *TimestreamQueryService) parseScheduleConfiguration(params map[string]interface{}) *tsstore.ScheduleConfiguration {
	scheduleConfigRaw := request.GetMapParamCaseInsensitive(params, "ScheduleConfiguration")
	if scheduleConfigRaw == nil {
		return nil
	}

	expr, _ := scheduleConfigRaw["ScheduleExpression"].(string)
	if expr == "" {
		return nil
	}
	return &tsstore.ScheduleConfiguration{
		ScheduleExpression: expr,
	}
}

func (s *TimestreamQueryService) parseNotificationConfiguration(params map[string]interface{}) *tsstore.NotificationConfiguration {
	notifConfigRaw := request.GetMapParamCaseInsensitive(params, "NotificationConfiguration")
	if notifConfigRaw == nil {
		return nil
	}

	snsConfigRaw, ok := notifConfigRaw["SnsConfiguration"].(map[string]interface{})
	if !ok {
		return nil
	}
	topicARN, _ := snsConfigRaw["TopicArn"].(string)
	if topicARN == "" {
		return nil
	}
	return &tsstore.NotificationConfiguration{
		SNSConfiguration: &tsstore.SNSConfiguration{
			TopicARN: topicARN,
		},
	}
}

func (s *TimestreamQueryService) parseErrorReportConfiguration(params map[string]interface{}) *tsstore.ErrorReportConfiguration {
	errorReportRaw := request.GetMapParamCaseInsensitive(params, "ErrorReportConfiguration")
	if errorReportRaw == nil {
		return nil
	}

	s3ConfigRaw, ok := errorReportRaw["S3Configuration"].(map[string]interface{})
	if !ok {
		return nil
	}
	bucketName, _ := s3ConfigRaw["BucketName"].(string)
	if bucketName == "" {
		return nil
	}
	objectKeyPrefix, _ := s3ConfigRaw["ObjectKeyPrefix"].(string)
	return &tsstore.ErrorReportConfiguration{
		S3Configuration: &tsstore.S3ErrorReportConfiguration{
			BucketName:      bucketName,
			ObjectKeyPrefix: objectKeyPrefix,
		},
	}
}

func (s *TimestreamQueryService) parseTargetConfiguration(params map[string]interface{}) *tsstore.TargetConfiguration {
	targetConfigRaw := request.GetMapParamCaseInsensitive(params, "TargetConfiguration")
	if targetConfigRaw == nil {
		return nil
	}

	tsConfigRaw, ok := targetConfigRaw["TimestreamConfiguration"].(map[string]interface{})
	if !ok {
		return nil
	}
	databaseName, _ := tsConfigRaw["DatabaseName"].(string)
	tableName, _ := tsConfigRaw["TableName"].(string)
	if databaseName == "" || tableName == "" {
		return nil
	}
	timeColumn, _ := tsConfigRaw["TimeColumn"].(string)
	return &tsstore.TargetConfiguration{
		TimestreamConfiguration: &tsstore.TimestreamConfiguration{
			DatabaseName: databaseName,
			TableName:    tableName,
			TimeColumn:   timeColumn,
		},
	}
}

// ListTagsForResource returns the tags for a scheduled query.
func (s *TimestreamQueryService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.tagHandlerConfig(st))
}

func epochFloat(t time.Time) float64 {
	return float64(t.UnixNano()) / 1e9
}

func (s *TimestreamQueryService) formatScheduledQueryResponse(sq *tsstore.ScheduledQuery, tags map[string]string) map[string]interface{} {
	response := map[string]interface{}{
		"Arn":          sq.ARN,
		"Name":         sq.Name,
		"State":        sq.ScheduledQueryStatus,
		"CreationTime": epochFloat(sq.CreationTime),
	}

	if sq.ErrorReportConfiguration != nil && sq.ErrorReportConfiguration.S3Configuration != nil {
		response["ErrorReportConfiguration"] = map[string]interface{}{
			"S3Configuration": map[string]interface{}{
				"BucketName":      sq.ErrorReportConfiguration.S3Configuration.BucketName,
				"ObjectKeyPrefix": sq.ErrorReportConfiguration.S3Configuration.ObjectKeyPrefix,
			},
		}
	}

	if sq.LastRunStatus != "" {
		response["LastRunStatus"] = sq.LastRunStatus
	}

	if !sq.NextRunTime.IsZero() {
		response["NextInvocationTime"] = epochFloat(sq.NextRunTime)
	}

	if !sq.PreviousRunTime.IsZero() {
		response["PreviousInvocationTime"] = epochFloat(sq.PreviousRunTime)
	}

	if sq.TargetConfiguration != nil && sq.TargetConfiguration.TimestreamConfiguration != nil {
		tsConfig := sq.TargetConfiguration.TimestreamConfiguration
		response["TargetDestination"] = map[string]interface{}{
			"TimestreamDestination": map[string]interface{}{
				"DatabaseName": tsConfig.DatabaseName,
				"TableName":    tsConfig.TableName,
			},
		}
	}

	return response
}

func (s *TimestreamQueryService) formatScheduledQueryDescriptionResponse(sq *tsstore.ScheduledQuery, tags map[string]string) map[string]interface{} {
	response := map[string]interface{}{
		"Arn":          sq.ARN,
		"Name":         sq.Name,
		"State":        sq.ScheduledQueryStatus,
		"QueryString":  sq.QueryString,
		"CreationTime": epochFloat(sq.CreationTime),
	}

	if sq.ScheduleConfiguration != nil {
		response["ScheduleConfiguration"] = map[string]interface{}{
			"ScheduleExpression": sq.ScheduleConfiguration.ScheduleExpression,
		}
	}

	if sq.NotificationConfiguration != nil && sq.NotificationConfiguration.SNSConfiguration != nil {
		response["NotificationConfiguration"] = map[string]interface{}{
			"SnsConfiguration": map[string]interface{}{
				"TopicArn": sq.NotificationConfiguration.SNSConfiguration.TopicARN,
			},
		}
	}

	if sq.ScheduledQueryExecutionRoleARN != "" {
		response["ScheduledQueryExecutionRoleArn"] = sq.ScheduledQueryExecutionRoleARN
	}

	if sq.KmsKeyID != "" {
		response["KmsKeyId"] = sq.KmsKeyID
	}

	if sq.ErrorReportConfiguration != nil && sq.ErrorReportConfiguration.S3Configuration != nil {
		response["ErrorReportConfiguration"] = map[string]interface{}{
			"S3Configuration": map[string]interface{}{
				"BucketName":      sq.ErrorReportConfiguration.S3Configuration.BucketName,
				"ObjectKeyPrefix": sq.ErrorReportConfiguration.S3Configuration.ObjectKeyPrefix,
			},
		}
	}

	if sq.LastRunStatus != "" {
		response["LastRunSummary"] = map[string]interface{}{
			"RunStatus": sq.LastRunStatus,
		}
	}

	if !sq.NextRunTime.IsZero() {
		response["NextInvocationTime"] = epochFloat(sq.NextRunTime)
	}

	if !sq.PreviousRunTime.IsZero() {
		response["PreviousInvocationTime"] = epochFloat(sq.PreviousRunTime)
	}

	if sq.TargetConfiguration != nil && sq.TargetConfiguration.TimestreamConfiguration != nil {
		tsConfig := sq.TargetConfiguration.TimestreamConfiguration
		response["TargetDestination"] = map[string]interface{}{
			"TimestreamDestination": map[string]interface{}{
				"DatabaseName": tsConfig.DatabaseName,
				"TableName":    tsConfig.TableName,
			},
		}
	}

	return response
}
