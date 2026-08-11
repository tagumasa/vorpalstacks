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
	if err := validateScheduledQueryName(name); err != nil {
		return nil, err
	}

	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")
	if queryString == "" {
		return nil, ErrValidationException
	}
	if err := validateQueryString(queryString); err != nil {
		return nil, err
	}

	scheduleConfig := s.parseScheduleConfiguration(req.Parameters)
	if scheduleConfig == nil {
		return nil, ErrValidationException
	}
	if err := validateScheduleExpression(scheduleConfig.ScheduleExpression); err != nil {
		return nil, err
	}

	notificationConfig := s.parseNotificationConfiguration(req.Parameters)
	roleARN := request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryExecutionRoleArn")
	if roleARN == "" {
		return nil, ErrValidationException
	}
	kmsKeyID := request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId")
	if kmsKeyID != "" && len(kmsKeyID) > maxAmazonResourceName {
		return nil, ErrValidationException
	}
	errorReportConfig, err := s.parseErrorReportConfiguration(req.Parameters)
	if err != nil {
		return nil, err
	}
	targetConfig, err := s.parseTargetConfiguration(req.Parameters)
	if err != nil {
		return nil, err
	}
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
	} else {
		if err := validateClientToken(clientToken); err != nil {
			return nil, err
		}
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
			return nil, ErrConflictException
		}
		return nil, ErrInternalServer
	}

	if tagMap := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")); len(tagMap) > 0 {
		if tagErr := st.scheduledQueryStore.Tag(sq.ARN, tagMap); tagErr != nil {
			logs.Warn("failed to tag scheduled query", logs.Err(tagErr), logs.String("arn", sq.ARN))
		}
	}

	return s.formatScheduledQueryResponse(sq), nil
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
	if err := validateAmazonResourceName(arn); err != nil {
		return nil, err
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

	runs, runErr := st.scheduledQueryRunStore.ListRuns(sq.ARN)
	if runErr == nil {
		for _, run := range runs {
			_ = st.scheduledQueryRunStore.DeleteRun(run.ARN)
		}
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
	if err := validateAmazonResourceName(arn); err != nil {
		return nil, err
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

	// Query the most recent run to populate LastRunSummary with full
	// details (InvocationTime, TriggerTime, ExecutionStats,
	// FailureReason) per the Smithy ScheduledQueryRunSummary shape.
	// Also collect failed runs for RecentlyFailedRuns.
	var lastRun *tsstore.ScheduledQueryRun
	var failedRuns []*tsstore.ScheduledQueryRun
	runs, runErr := st.scheduledQueryRunStore.ListRuns(sq.ARN)
	if runErr == nil && len(runs) > 0 {
		lastRun = runs[len(runs)-1]
		for _, r := range runs {
			if r.RunStatus == tsstore.ScheduleRunStatusFailed {
				failedRuns = append(failedRuns, r)
			}
		}
	}

	return map[string]interface{}{
		"ScheduledQuery": s.formatScheduledQueryDescriptionResponse(sq, lastRun, failedRuns),
	}, nil
}

// ListScheduledQueries returns a list of scheduled queries.
func (s *TimestreamQueryService) ListScheduledQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	maxResults := 0
	hasMaxResults := false
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxStr != "" {
		if val, atoiErr := strconv.Atoi(maxStr); atoiErr == nil {
			maxResults = val
			hasMaxResults = true
		}
	}

	offset := 0
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		val, atoiErr := strconv.Atoi(nextToken)
		if atoiErr != nil || val < 0 {
			return nil, ErrValidationException
		}
		offset = val
	}

	result, err := s.ListScheduledQueriesCore(st, ListScheduledQueriesInput{
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     offset,
	})
	if err != nil {
		return nil, err
	}

	var scheduledQueries []map[string]interface{}
	for _, summary := range result.Summaries {
		scheduledQueries = append(scheduledQueries, formatScheduledQuerySummaryResponse(summary))
	}

	response := map[string]interface{}{
		"ScheduledQueries": scheduledQueries,
	}

	if result.NextToken != "" {
		response["NextToken"] = result.NextToken
	}

	return response, nil
}

// UpdateScheduledQuery updates an existing scheduled query.
func (s *TimestreamQueryService) UpdateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn")
	if arn == "" {
		return nil, ErrValidationException
	}
	if err := validateAmazonResourceName(arn); err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := st.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return nil, ErrValidationException
	}

	state := tsstore.ScheduledQueryStatus(request.GetParamCaseInsensitive(req.Parameters, "State"))
	if state != tsstore.ScheduledQueryStatusEnabled && state != tsstore.ScheduledQueryStatusDisabled {
		return nil, ErrValidationException
	}

	_, err = st.scheduledQueryStore.UpdateScheduledQuery(name, state)
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
	if err := validateAmazonResourceName(arn); err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := st.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return nil, ErrValidationException
	}

	if ct := request.GetParamCaseInsensitive(req.Parameters, "ClientToken"); ct != "" {
		if err := validateClientToken(ct); err != nil {
			return nil, err
		}
	}

	sq, err := st.scheduledQueryStore.GetScheduledQuery(name)
	if err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	// Parse InvocationTime (defaults to now). QueryInsights is accepted
	// but currently not applied to the query execution (would require
	// extending the query executor to collect per-query insights).
	now := time.Now().UTC()
	invocationTimeRaw := request.GetParamCaseInsensitive(req.Parameters, "InvocationTime")
	if invocationTimeRaw == "" {
		if raw, exists := req.Parameters["InvocationTime"]; exists && raw != nil {
			switch v := raw.(type) {
			case float64:
				if v > 1e12 {
					now = time.UnixMilli(int64(v)).UTC()
				} else {
					now = time.Unix(int64(v), 0).UTC()
				}
			case int64:
				now = time.Unix(v, 0).UTC()
			case int:
				now = time.Unix(int64(v), 0).UTC()
			}
		} else {
			return nil, ErrValidationException
		}
	} else {
		if parsed, err := time.Parse(time.RFC3339Nano, invocationTimeRaw); err == nil {
			now = parsed.UTC()
		} else if parsed, err := time.Parse(time.RFC3339, invocationTimeRaw); err == nil {
			now = parsed.UTC()
		}
	}

	run, err := s.executeScheduledQueryInternal(ctx, st, sq, now, tsstore.TriggerTypeManual)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ScheduledQueryRunArn": run.ARN,
	}, nil
}

// executeScheduledQueryInternal executes a scheduled query and records the
// run lifecycle (PENDING → RUNNING → SUCCEEDED/FAILED). Shared by the manual
// ExecuteScheduledQuery handler and the background auto-trigger engine.
func (s *TimestreamQueryService) executeScheduledQueryInternal(ctx context.Context, st *tsQueryStores, sq *tsstore.ScheduledQuery, invocationTime time.Time, triggerType string) (*tsstore.ScheduledQueryRun, error) {
	name := sq.Name

	run, err := st.scheduledQueryRunStore.CreateRun(sq.ARN, invocationTime, invocationTime, triggerType)
	if err != nil {
		return nil, ErrInternalServer
	}

	if err := st.scheduledQueryRunStore.UpdateRunStatus(run.ARN, tsstore.ScheduleRunStatusRunning, "", nil); err != nil {
		logs.Error("Failed to update scheduled query run to RUNNING", logs.String("arn", run.ARN), logs.Err(err))
		return nil, ErrInternalServer
	}

	result, execErr := s.executeSQLQuery(ctx, st, sq.QueryString)

	if execErr != nil {
		if err := st.scheduledQueryRunStore.UpdateRunStatus(run.ARN, tsstore.ScheduleRunStatusFailed, execErr.Error(), nil); err != nil {
			logs.Error("Failed to update scheduled query run to FAILED", logs.String("arn", run.ARN), logs.Err(err))
		}
		if err := st.scheduledQueryStore.UpdateLastRun(name, tsstore.ScheduledQueryRunStatusFromTrigger(triggerType, tsstore.ScheduleRunStatusFailed), invocationTime); err != nil {
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
	if err := st.scheduledQueryStore.UpdateLastRun(name, tsstore.ScheduledQueryRunStatusFromTrigger(triggerType, tsstore.ScheduleRunStatusSucceeded), invocationTime); err != nil {
		logs.Error("Failed to update last run status for scheduled query", logs.String("name", name), logs.Err(err))
	}

	return run, nil
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

func (s *TimestreamQueryService) parseErrorReportConfiguration(params map[string]interface{}) (*tsstore.ErrorReportConfiguration, error) {
	errorReportRaw := request.GetMapParamCaseInsensitive(params, "ErrorReportConfiguration")
	if errorReportRaw == nil {
		return nil, nil
	}

	s3ConfigRaw, ok := errorReportRaw["S3Configuration"].(map[string]interface{})
	if !ok {
		return nil, ErrValidationException
	}
	bucketName, _ := s3ConfigRaw["BucketName"].(string)
	if bucketName == "" {
		return nil, ErrValidationException
	}
	if err := validateS3BucketName(bucketName); err != nil {
		return nil, err
	}
	objectKeyPrefix, _ := s3ConfigRaw["ObjectKeyPrefix"].(string)
	if objectKeyPrefix != "" {
		if err := validateS3ObjectKeyPrefix(objectKeyPrefix); err != nil {
			return nil, err
		}
	}
	encryptionOption, _ := s3ConfigRaw["EncryptionOption"].(string)
	if encryptionOption != "" && !validateS3EncryptionOption(encryptionOption) {
		return nil, ErrValidationException
	}
	return &tsstore.ErrorReportConfiguration{
		S3Configuration: &tsstore.S3ErrorReportConfiguration{
			BucketName:       bucketName,
			ObjectKeyPrefix:  objectKeyPrefix,
			EncryptionOption: encryptionOption,
		},
	}, nil
}

func (s *TimestreamQueryService) parseTargetConfiguration(params map[string]interface{}) (*tsstore.TargetConfiguration, error) {
	targetConfigRaw := request.GetMapParamCaseInsensitive(params, "TargetConfiguration")
	if targetConfigRaw == nil {
		return nil, nil
	}

	tsConfigRaw, ok := targetConfigRaw["TimestreamConfiguration"].(map[string]interface{})
	if !ok {
		return nil, ErrValidationException
	}
	databaseName, _ := tsConfigRaw["DatabaseName"].(string)
	tableName, _ := tsConfigRaw["TableName"].(string)
	if databaseName == "" || tableName == "" {
		return nil, ErrValidationException
	}
	timeColumn, _ := tsConfigRaw["TimeColumn"].(string)
	measureNameColumn, _ := tsConfigRaw["MeasureNameColumn"].(string)

	var dimensionMappings []tsstore.QueryDimensionMapping
	if dmList, ok := tsConfigRaw["DimensionMappings"].([]interface{}); ok {
		for _, dm := range dmList {
			dmMap, ok := dm.(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := dmMap["Name"].(string)
			dvt, _ := dmMap["DimensionValueType"].(string)
			if name != "" {
				if dvt != "" && !validateDimensionValueType(dvt) {
					return nil, ErrValidationException
				}
				dimensionMappings = append(dimensionMappings, tsstore.QueryDimensionMapping{
					Name:               name,
					DimensionValueType: dvt,
				})
			}
		}
	}

	var multiMeasureMappings *tsstore.MultiMeasureMappings
	if mmm, ok := tsConfigRaw["MultiMeasureMappings"].(map[string]interface{}); ok {
		mmmResult, err := parseMultiMeasureMappings(mmm)
		if err != nil {
			return nil, err
		}
		multiMeasureMappings = mmmResult
	}

	var mixedMeasureMappings []tsstore.MixedMeasureMapping
	if mmmList, ok := tsConfigRaw["MixedMeasureMappings"].([]interface{}); ok {
		for _, m := range mmmList {
			if mmmMap, ok := m.(map[string]interface{}); ok {
				mapping, err := parseMixedMeasureMapping(mmmMap)
				if err != nil {
					return nil, err
				}
				if mapping != nil {
					mixedMeasureMappings = append(mixedMeasureMappings, *mapping)
				}
			}
		}
	}

	return &tsstore.TargetConfiguration{
		TimestreamConfiguration: &tsstore.TimestreamConfiguration{
			DatabaseName:         databaseName,
			TableName:            tableName,
			TimeColumn:           timeColumn,
			DimensionMappings:    dimensionMappings,
			MultiMeasureMappings: multiMeasureMappings,
			MixedMeasureMappings: mixedMeasureMappings,
			MeasureNameColumn:    measureNameColumn,
		},
	}, nil
}

func parseMultiMeasureMappings(raw map[string]interface{}) (*tsstore.MultiMeasureMappings, error) {
	mmm := &tsstore.MultiMeasureMappings{}
	if v, ok := raw["TargetMultiMeasureName"].(string); ok {
		mmm.TargetMultiMeasureName = v
	}
	if list, ok := raw["MultiMeasureAttributeMappings"].([]interface{}); ok {
		for _, item := range list {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			srcCol, _ := itemMap["SourceColumn"].(string)
			targetName, _ := itemMap["TargetMultiMeasureAttributeName"].(string)
			measureType, _ := itemMap["MeasureValueType"].(string)
			if measureType != "" && !validateMeasureValueType(measureType) {
				return nil, ErrValidationException
			}
			if srcCol != "" {
				mmm.MultiMeasureAttributeMappings = append(mmm.MultiMeasureAttributeMappings, tsstore.MultiMeasureAttributeMapping{
					SourceColumn:                    &tsstore.SourceColumn{Name: srcCol},
					TargetMultiMeasureAttributeName: targetName,
					MeasureValueMeasureValueType:    tsstore.MeasureValueType(measureType),
				})
			}
		}
	}
	if len(mmm.MultiMeasureAttributeMappings) == 0 {
		return nil, nil
	}
	return mmm, nil
}

// parseMixedMeasureMapping parses a single MixedMeasureMapping from a raw map.
func parseMixedMeasureMapping(raw map[string]interface{}) (*tsstore.MixedMeasureMapping, error) {
	mapping := &tsstore.MixedMeasureMapping{}
	if v, ok := raw["MeasureName"].(string); ok {
		mapping.MeasureName = v
	}
	if v, ok := raw["SourceColumn"].(string); ok {
		mapping.SourceColumn = v
	}
	if v, ok := raw["TargetMeasureName"].(string); ok {
		mapping.TargetMeasureName = v
	}
	if v, ok := raw["MeasureValueType"].(string); ok {
		if v != "" && !validateMeasureValueType(v) {
			return nil, ErrValidationException
		}
		mapping.MeasureValueMeasureValueType = tsstore.MeasureValueType(v)
	}
	if list, ok := raw["MultiMeasureAttributeMappings"].([]interface{}); ok {
		for _, item := range list {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			srcCol, _ := itemMap["SourceColumn"].(string)
			targetName, _ := itemMap["TargetMultiMeasureAttributeName"].(string)
			measureType, _ := itemMap["MeasureValueType"].(string)
			if measureType != "" && !validateMeasureValueType(measureType) {
				return nil, ErrValidationException
			}
			if srcCol != "" {
				mapping.MultiMeasureAttributeMappings = append(mapping.MultiMeasureAttributeMappings, tsstore.MultiMeasureAttributeMapping{
					SourceColumn:                    &tsstore.SourceColumn{Name: srcCol},
					TargetMultiMeasureAttributeName: targetName,
					MeasureValueMeasureValueType:    tsstore.MeasureValueType(measureType),
				})
			}
		}
	}
	if mapping.SourceColumn == "" && mapping.MeasureName == "" && len(mapping.MultiMeasureAttributeMappings) == 0 {
		return nil, nil
	}
	return mapping, nil
}

// ListTagsForResource returns the tags for a scheduled query.
// Implements MaxResults validation (range {1, 200}) and NextToken
// pagination per the Smithy model.
func (s *TimestreamQueryService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cfg := s.tagHandlerConfig(st)

	rawKey := tagutil.GetResourceKey(req.Parameters, cfg.Param)
	if rawKey == "" {
		return nil, ErrValidationException
	}
	if err := validateAmazonResourceName(rawKey); err != nil {
		return nil, err
	}
	resourceKey := rawKey
	if cfg.ResourceKey != nil {
		resourceKey = cfg.ResourceKey(rawKey)
	}

	if cfg.ListFunc == nil {
		return nil, ErrInternalServer
	}

	tags, err := cfg.ListFunc(ctx, resourceKey)
	if err != nil {
		return nil, err
	}

	maxResults := 0
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxStr != "" {
		val, atoiErr := strconv.Atoi(maxStr)
		if atoiErr != nil {
			return nil, ErrValidationException
		}
		if err := validateMaxResultsTags(val); err != nil {
			return nil, err
		}
		maxResults = val
	}

	offset := 0
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		val, atoiErr := strconv.Atoi(nextToken)
		if atoiErr != nil || val < 0 {
			return nil, ErrValidationException
		}
		offset = val
	}

	totalTags := len(tags)
	if offset > totalTags {
		offset = totalTags
	}
	end := totalTags
	if maxResults > 0 && offset+maxResults < end {
		end = offset + maxResults
	}
	pagedTags := tags[offset:end]

	tagResponse := tagutil.ToResponseWithKeyNames(pagedTags, cfg.Param.TagKeyName, cfg.Param.TagValueName)

	resp := map[string]interface{}{
		cfg.Param.TagsParam: tagResponse,
	}
	if end < totalTags {
		resp["NextToken"] = strconv.Itoa(end)
	}

	return resp, nil
}

func epochFloat(t time.Time) float64 {
	return float64(t.UnixNano()) / 1e9
}

func (s *TimestreamQueryService) formatScheduledQueryBaseResponse(sq *tsstore.ScheduledQuery) map[string]interface{} {
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

	if !sq.NextRunTime.IsZero() {
		response["NextInvocationTime"] = epochFloat(sq.NextRunTime)
	}

	if !sq.PreviousRunTime.IsZero() {
		response["PreviousInvocationTime"] = epochFloat(sq.PreviousRunTime)
	}

	return response
}

func (s *TimestreamQueryService) formatScheduledQueryResponse(sq *tsstore.ScheduledQuery) map[string]interface{} {
	response := s.formatScheduledQueryBaseResponse(sq)

	// TargetDestination is the ListScheduledQueries representation of the
	// target — a read-only projection containing only DatabaseName and
	// TableName (Smithy: ScheduledQuery.TargetDestination).
	if sq.TargetConfiguration != nil && sq.TargetConfiguration.TimestreamConfiguration != nil {
		tsConfig := sq.TargetConfiguration.TimestreamConfiguration
		response["TargetDestination"] = map[string]interface{}{
			"TimestreamDestination": map[string]interface{}{
				"DatabaseName": tsConfig.DatabaseName,
				"TableName":    tsConfig.TableName,
			},
		}
	}

	if sq.LastRunStatus != "" {
		response["LastRunStatus"] = sq.LastRunStatus
	}

	return response
}

// formatScheduledQuerySummaryResponse converts a ScheduledQuerySummary DTO
// (returned by ListScheduledQueriesCore) into the JSON map representation
// expected by the HTTP API response.
func formatScheduledQuerySummaryResponse(summary *ScheduledQuerySummary) map[string]interface{} {
	response := map[string]interface{}{
		"Arn":          summary.ARN,
		"Name":         summary.Name,
		"State":        summary.State,
		"CreationTime": epochFloat(summary.CreationTime),
	}

	if summary.ErrorReportConfiguration != nil && summary.ErrorReportConfiguration.S3Configuration != nil {
		response["ErrorReportConfiguration"] = map[string]interface{}{
			"S3Configuration": map[string]interface{}{
				"BucketName":      summary.ErrorReportConfiguration.S3Configuration.BucketName,
				"ObjectKeyPrefix": summary.ErrorReportConfiguration.S3Configuration.ObjectKeyPrefix,
			},
		}
	}

	if !summary.NextRunTime.IsZero() {
		response["NextInvocationTime"] = epochFloat(summary.NextRunTime)
	}

	if !summary.PreviousRunTime.IsZero() {
		response["PreviousInvocationTime"] = epochFloat(summary.PreviousRunTime)
	}

	if summary.TargetConfiguration != nil && summary.TargetConfiguration.TimestreamConfiguration != nil {
		tsConfig := summary.TargetConfiguration.TimestreamConfiguration
		response["TargetDestination"] = map[string]interface{}{
			"TimestreamDestination": map[string]interface{}{
				"DatabaseName": tsConfig.DatabaseName,
				"TableName":    tsConfig.TableName,
			},
		}
	}

	if summary.LastRunStatus != "" {
		response["LastRunStatus"] = summary.LastRunStatus
	}

	return response
}

func (s *TimestreamQueryService) formatScheduledQueryDescriptionResponse(sq *tsstore.ScheduledQuery, lastRun *tsstore.ScheduledQueryRun, failedRuns []*tsstore.ScheduledQueryRun) map[string]interface{} {
	response := s.formatScheduledQueryBaseResponse(sq)

	response["QueryString"] = sq.QueryString

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

	// TargetConfiguration is the DescribeScheduledQuery representation of
	// the target — the full write-side configuration (Smithy:
	// ScheduledQueryDescription.TargetConfiguration). This is distinct
	// from TargetDestination which is a read-only projection used in
	// ListScheduledQueries.
	if sq.TargetConfiguration != nil && sq.TargetConfiguration.TimestreamConfiguration != nil {
		tsConfig := sq.TargetConfiguration.TimestreamConfiguration
		tsMap := map[string]interface{}{
			"DatabaseName": tsConfig.DatabaseName,
			"TableName":    tsConfig.TableName,
		}
		if tsConfig.TimeColumn != "" {
			tsMap["TimeColumn"] = tsConfig.TimeColumn
		}
		if tsConfig.MeasureNameColumn != "" {
			tsMap["MeasureNameColumn"] = tsConfig.MeasureNameColumn
		}
		if len(tsConfig.DimensionMappings) > 0 {
			dmList := make([]map[string]interface{}, len(tsConfig.DimensionMappings))
			for i, dm := range tsConfig.DimensionMappings {
				dmList[i] = map[string]interface{}{
					"Name":               dm.Name,
					"DimensionValueType": dm.DimensionValueType,
				}
			}
			tsMap["DimensionMappings"] = dmList
		}
		if tsConfig.MultiMeasureMappings != nil && len(tsConfig.MultiMeasureMappings.MultiMeasureAttributeMappings) > 0 {
			mmmMap := map[string]interface{}{}
			if tsConfig.MultiMeasureMappings.TargetMultiMeasureName != "" {
				mmmMap["TargetMultiMeasureName"] = tsConfig.MultiMeasureMappings.TargetMultiMeasureName
			}
			attrs := make([]map[string]interface{}, len(tsConfig.MultiMeasureMappings.MultiMeasureAttributeMappings))
			for i, a := range tsConfig.MultiMeasureMappings.MultiMeasureAttributeMappings {
				attrMap := map[string]interface{}{}
				if a.SourceColumn != nil {
					attrMap["SourceColumn"] = a.SourceColumn.Name
				}
				attrMap["TargetMultiMeasureAttributeName"] = a.TargetMultiMeasureAttributeName
				attrMap["MeasureValueType"] = a.MeasureValueMeasureValueType
				attrs[i] = attrMap
			}
			mmmMap["MultiMeasureAttributeMappings"] = attrs
			tsMap["MultiMeasureMappings"] = mmmMap
		}
		if len(tsConfig.MixedMeasureMappings) > 0 {
			mmmList := make([]map[string]interface{}, len(tsConfig.MixedMeasureMappings))
			for i, m := range tsConfig.MixedMeasureMappings {
				mMap := map[string]interface{}{}
				if m.MeasureName != "" {
					mMap["MeasureName"] = m.MeasureName
				}
				if m.SourceColumn != "" {
					mMap["SourceColumn"] = m.SourceColumn
				}
				if m.TargetMeasureName != "" {
					mMap["TargetMeasureName"] = m.TargetMeasureName
				}
				mMap["MeasureValueType"] = m.MeasureValueMeasureValueType
				if len(m.MultiMeasureAttributeMappings) > 0 {
					attrs := make([]map[string]interface{}, len(m.MultiMeasureAttributeMappings))
					for j, a := range m.MultiMeasureAttributeMappings {
						attrMap := map[string]interface{}{}
						if a.SourceColumn != nil {
							attrMap["SourceColumn"] = a.SourceColumn.Name
						}
						attrMap["TargetMultiMeasureAttributeName"] = a.TargetMultiMeasureAttributeName
						attrMap["MeasureValueType"] = a.MeasureValueMeasureValueType
						attrs[j] = attrMap
					}
					mMap["MultiMeasureAttributeMappings"] = attrs
				}
				mmmList[i] = mMap
			}
			tsMap["MixedMeasureMappings"] = mmmList
		}
		response["TargetConfiguration"] = map[string]interface{}{
			"TimestreamConfiguration": tsMap,
		}
	}

	if sq.LastRunStatus != "" {
		summary := map[string]interface{}{
			"RunStatus": sq.LastRunStatus,
		}
		if lastRun != nil {
			if !lastRun.InvocationTime.IsZero() {
				summary["InvocationTime"] = epochFloat(lastRun.InvocationTime)
			}
			if !lastRun.TriggerTime.IsZero() {
				summary["TriggerTime"] = epochFloat(lastRun.TriggerTime)
			}
			if lastRun.Error != "" {
				summary["FailureReason"] = lastRun.Error
			}
			if lastRun.ExecutionStats != nil {
				stats := map[string]interface{}{}
				if lastRun.ExecutionStats.DataWrites > 0 {
					stats["DataWrites"] = lastRun.ExecutionStats.DataWrites
				}
				if lastRun.ExecutionStats.BytesMetered > 0 {
					stats["BytesMetered"] = lastRun.ExecutionStats.BytesMetered
				}
				if lastRun.ExecutionStats.QueryResultRows > 0 {
					stats["QueryResultRows"] = lastRun.ExecutionStats.QueryResultRows
				}
				if lastRun.ExecutionStats.CumulativeBytesScanned > 0 {
					stats["CumulativeBytesScanned"] = lastRun.ExecutionStats.CumulativeBytesScanned
				}
				if lastRun.ExecutionStats.ExecutionTimeInMillis > 0 {
					stats["ExecutionTimeInMillis"] = lastRun.ExecutionStats.ExecutionTimeInMillis
				}
				if lastRun.ExecutionStats.RecordsIngested > 0 {
					stats["RecordsIngested"] = lastRun.ExecutionStats.RecordsIngested
				}
				if len(stats) > 0 {
					summary["ExecutionStats"] = stats
				}
			}
		}
		response["LastRunSummary"] = summary
	}

	// Populate RecentlyFailedRuns (Smithy: ScheduledQueryRunSummaryList).
	if len(failedRuns) > 0 {
		recentlyFailed := make([]map[string]interface{}, 0, len(failedRuns))
		for _, fr := range failedRuns {
			summary := map[string]interface{}{
				"RunStatus": tsstore.ScheduledQueryRunStatusFromTrigger(fr.TriggerType, fr.RunStatus),
			}
			if !fr.InvocationTime.IsZero() {
				summary["InvocationTime"] = epochFloat(fr.InvocationTime)
			}
			if !fr.TriggerTime.IsZero() {
				summary["TriggerTime"] = epochFloat(fr.TriggerTime)
			}
			if fr.Error != "" {
				summary["FailureReason"] = fr.Error
			}
			if fr.ExecutionStats != nil {
				stats := map[string]interface{}{}
				if fr.ExecutionStats.DataWrites > 0 {
					stats["DataWrites"] = fr.ExecutionStats.DataWrites
				}
				if fr.ExecutionStats.BytesMetered > 0 {
					stats["BytesMetered"] = fr.ExecutionStats.BytesMetered
				}
				if fr.ExecutionStats.QueryResultRows > 0 {
					stats["QueryResultRows"] = fr.ExecutionStats.QueryResultRows
				}
				if fr.ExecutionStats.CumulativeBytesScanned > 0 {
					stats["CumulativeBytesScanned"] = fr.ExecutionStats.CumulativeBytesScanned
				}
				if fr.ExecutionStats.ExecutionTimeInMillis > 0 {
					stats["ExecutionTimeInMillis"] = fr.ExecutionStats.ExecutionTimeInMillis
				}
				if fr.ExecutionStats.RecordsIngested > 0 {
					stats["RecordsIngested"] = fr.ExecutionStats.RecordsIngested
				}
				if len(stats) > 0 {
					summary["ExecutionStats"] = stats
				}
			}
			recentlyFailed = append(recentlyFailed, summary)
		}
		response["RecentlyFailedRuns"] = recentlyFailed
	}

	return response
}
