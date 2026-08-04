package timestreamquery

import (
	"context"
	"regexp"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	tsstore "vorpalstacks/internal/store/aws/timestream"
	"vorpalstacks/internal/utils/aws/types"

	"github.com/google/uuid"
)

// scheduledQueryNamePattern implements the Smithy ScheduledQueryName pattern:
// ^[a-zA-Z0-9|!\-_*'()]([a-zA-Z0-9]|[!\-_*'()/.])+$
var scheduledQueryNamePattern = regexp.MustCompile(`^[a-zA-Z0-9|!\-_*'()]([a-zA-Z0-9]|[!\-_*'()/.])+$`)

// validateScheduledQueryName validates the Name parameter against the Smithy
// ScheduledQueryName shape: length {min:1, max:64} and the pattern above.
func validateScheduledQueryName(name string) error {
	if len(name) < 1 || len(name) > 64 {
		return awserrors.NewAWSError("ValidationException",
			"ScheduledQueryName must be between 1 and 64 characters.", 400)
	}
	if !scheduledQueryNamePattern.MatchString(name) {
		return awserrors.NewAWSError("ValidationException",
			"ScheduledQueryName does not match the required pattern.", 400)
	}
	return nil
}

// validateClientToken validates a ClientToken against the Smithy
// ClientToken shape: length {min:32, max:128}. Only called when the
// client explicitly provides a token (empty tokens are replaced with a
// generated UUID by the handler).
func validateClientToken(token string) error {
	if len(token) < 32 || len(token) > 128 {
		return awserrors.NewAWSError("ValidationException",
			"ClientToken must be between 32 and 128 characters.", 400)
	}
	return nil
}

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

	scheduleConfig := s.parseScheduleConfiguration(req.Parameters)
	if scheduleConfig == nil {
		return nil, ErrValidationException
	}

	notificationConfig := s.parseNotificationConfiguration(req.Parameters)
	roleARN := request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryExecutionRoleArn")
	kmsKeyID := request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId")
	// KmsKeyId targets Smithy StringValue2048 — length 1-2048.
	if kmsKeyID != "" && len(kmsKeyID) > 2048 {
		return nil, ErrValidationException
	}
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
	if maxResults > maxListScheduledQueries {
		maxResults = maxListScheduledQueries
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
		scheduledQueries = append(scheduledQueries, s.formatScheduledQueryResponse(sq))
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
	if invocationTimeStr := request.GetParamCaseInsensitive(req.Parameters, "InvocationTime"); invocationTimeStr != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, invocationTimeStr); err == nil {
			now = parsed.UTC()
		} else if parsed, err := time.Parse(time.RFC3339, invocationTimeStr); err == nil {
			now = parsed.UTC()
		}
	}
	run, err := st.scheduledQueryRunStore.CreateRun(sq.ARN, now, now, tsstore.TriggerTypeManual)
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
		if err := st.scheduledQueryStore.UpdateLastRun(name, tsstore.ScheduledQueryRunStatusManualTriggerFailure, now); err != nil {
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
	if err := st.scheduledQueryStore.UpdateLastRun(name, tsstore.ScheduledQueryRunStatusManualTriggerSuccess, now); err != nil {
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
	encryptionOption, _ := s3ConfigRaw["EncryptionOption"].(string)
	return &tsstore.ErrorReportConfiguration{
		S3Configuration: &tsstore.S3ErrorReportConfiguration{
			BucketName:       bucketName,
			ObjectKeyPrefix:  objectKeyPrefix,
			EncryptionOption: encryptionOption,
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
	measureNameColumn, _ := tsConfigRaw["MeasureNameColumn"].(string)

	// DimensionMappings (Smithy: [{Name, DimensionValueType}]).
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
				dimensionMappings = append(dimensionMappings, tsstore.QueryDimensionMapping{
					Name:               name,
					DimensionValueType: dvt,
				})
			}
		}
	}

	// MultiMeasureMappings (Smithy: {TargetMultiMeasureName, MultiMeasureAttributeMappings}).
	var multiMeasureMappings *tsstore.MultiMeasureMappings
	if mmm, ok := tsConfigRaw["MultiMeasureMappings"].(map[string]interface{}); ok {
		multiMeasureMappings = parseMultiMeasureMappings(mmm)
	}

	// MixedMeasureMappings (Smithy: list of MixedMeasureMapping).
	var mixedMeasureMappings []tsstore.MixedMeasureMapping
	if mmmList, ok := tsConfigRaw["MixedMeasureMappings"].([]interface{}); ok {
		for _, m := range mmmList {
			if mmmMap, ok := m.(map[string]interface{}); ok {
				if mapping := parseMixedMeasureMapping(mmmMap); mapping != nil {
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
	}
}

// parseMultiMeasureMappings parses MultiMeasureMappings from a raw map
// (Smithy: {TargetMultiMeasureName, MultiMeasureAttributeMappings}).
func parseMultiMeasureMappings(raw map[string]interface{}) *tsstore.MultiMeasureMappings {
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
		return nil
	}
	return mmm
}

// parseMixedMeasureMapping parses a single MixedMeasureMapping from a raw map.
func parseMixedMeasureMapping(raw map[string]interface{}) *tsstore.MixedMeasureMapping {
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
		return nil
	}
	return mapping
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
		if atoiErr != nil || val < 1 || val > maxListTagsForResource {
			return nil, awserrors.NewAWSError("ValidationException",
				"MaxResults must be between 1 and 200.", 400)
		}
		maxResults = val
	}

	offset := 0
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		if val, atoiErr := strconv.Atoi(nextToken); atoiErr == nil && val >= 0 {
			offset = val
		}
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
