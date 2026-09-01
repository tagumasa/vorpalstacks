package timestreamquery

import (
	"context"
	"strconv"
	"time"
	"unicode/utf8"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	tsstore "vorpalstacks/internal/store/aws/timestream"

	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// DTO types — service-layer representations that decouple the admin handler
// from store types. The admin handler only sees these types and never imports
// store packages directly. The HTTP (AWS API) handlers build the same DTOs
// from the wire request; nested configuration members stay in raw wire form
// and the Core owns their parsing and validation.
// ---------------------------------------------------------------------------

// ScheduledQuerySummary is the projection of a ScheduledQuery used by the
// ListScheduledQueries admin console operation.
type ScheduledQuerySummary struct {
	ARN                      string
	Name                     string
	CreationTime             time.Time
	State                    string
	LastRunStatus            string
	PreviousRunTime          time.Time
	NextRunTime              time.Time
	ErrorReportConfiguration *tsstore.ErrorReportConfiguration
	TargetConfiguration      *tsstore.TargetConfiguration
}

// ListScheduledQueriesInput is the DTO input for the ListScheduledQueries
// operation on both planes. The pagination members stay in their raw wire
// form; the admin console sends none (empty strings select the defaults).
type ListScheduledQueriesInput struct {
	Region        string
	MaxResultsRaw string
	NextTokenRaw  string
}

// ListScheduledQueriesResult is the DTO result for the ListScheduledQueries
// admin console operation.
type ListScheduledQueriesResult struct {
	Summaries []*ScheduledQuerySummary
	NextToken string
}

// CreateScheduledQueryInput is the DTO input for the CreateScheduledQuery
// operation.
type CreateScheduledQueryInput struct {
	Name                  string
	QueryString           string
	ScheduleConfigRaw     map[string]interface{}
	NotificationConfigRaw map[string]interface{}
	ErrorReportConfigRaw  map[string]interface{}
	TargetConfigRaw       map[string]interface{}
	RoleARN               string
	KmsKeyID              string
	ClientToken           string
	Tags                  map[string]string
	IAMValidator          *iam.IAMValidator
}

// DescribeScheduledQueryResult is the DTO result for the
// DescribeScheduledQuery operation. FailedRunStatuses carries the
// trigger-aware run status for each entry of FailedRuns, precomputed here
// so response formatting needs no store calls.
type DescribeScheduledQueryResult struct {
	SQ                *tsstore.ScheduledQuery
	LastRun           *tsstore.ScheduledQueryRun
	FailedRuns        []*tsstore.ScheduledQueryRun
	FailedRunStatuses []string
}

// ExecuteScheduledQueryInput is the DTO input for the ExecuteScheduledQuery
// operation. InvocationTime keeps its raw wire forms (the string form and
// the untyped raw value); the Core owns the RFC3339/epoch parsing and
// validation.
type ExecuteScheduledQueryInput struct {
	ScheduledQueryARN string
	ClientToken       string
	InvocationTimeStr string
	InvocationTimeRaw interface{}
}

// ListTagsForResourceResult is the DTO result for the ListTagsForResource
// operation: one page of tags with the next-page token when more remain.
type ListTagsForResourceResult struct {
	Tags      []tagutil.Tag
	NextToken string
}

// ---------------------------------------------------------------------------
// Core methods — the protocol handlers delegate to these, passing the store
// group obtained via s.store (HTTP) or admin_handler_convert.go (admin).
// ---------------------------------------------------------------------------

// ListScheduledQueriesCore returns a paginated list of scheduled query
// summaries for both protocol planes. It operates on the provided store
// group and returns service-layer DTOs.
func (s *TimestreamQueryService) ListScheduledQueriesCore(stores *tsQueryStores, input ListScheduledQueriesInput) (*ListScheduledQueriesResult, error) {
	queries, err := stores.scheduledQueryStore.ListScheduledQueries()
	if err != nil {
		return nil, ErrInternalServer
	}

	maxResults := 0
	hasMaxResults := false
	if input.MaxResultsRaw != "" {
		if val, atoiErr := strconv.Atoi(input.MaxResultsRaw); atoiErr == nil {
			maxResults = val
			hasMaxResults = true
		}
	}

	offset := 0
	if input.NextTokenRaw != "" {
		val, atoiErr := strconv.Atoi(input.NextTokenRaw)
		if atoiErr != nil || val < 0 {
			return nil, ErrValidationException
		}
		offset = val
	}

	if hasMaxResults {
		if err := validateMaxResultsInRange(maxResults, "MaxResults", rangeMaxScheduledQueriesResults); err != nil {
			return nil, err
		}
	} else {
		maxResults = maxListScheduledQueries
	}

	var summaries []*ScheduledQuerySummary
	for i, sq := range queries {
		if i < offset {
			continue
		}
		if len(summaries) >= maxResults {
			break
		}
		summary := &ScheduledQuerySummary{
			ARN:                      sq.ARN,
			Name:                     sq.Name,
			CreationTime:             sq.CreationTime,
			State:                    string(sq.ScheduledQueryStatus),
			LastRunStatus:            sq.LastRunStatus,
			PreviousRunTime:          sq.PreviousRunTime,
			NextRunTime:              sq.NextRunTime,
			ErrorReportConfiguration: sq.ErrorReportConfiguration,
			TargetConfiguration:      sq.TargetConfiguration,
		}
		summaries = append(summaries, summary)
	}

	result := &ListScheduledQueriesResult{
		Summaries: summaries,
	}

	if offset+maxResults < len(queries) {
		result.NextToken = strconv.Itoa(offset + maxResults)
	}

	return result, nil
}

// createScheduledQueryCore validates the request, creates the scheduled
// query and applies the initial tags.
func (s *TimestreamQueryService) createScheduledQueryCore(ctx context.Context, stores *tsQueryStores, input CreateScheduledQueryInput) (*tsstore.ScheduledQuery, error) {
	name := input.Name
	if name == "" {
		return nil, ErrValidationException
	}
	if err := validateScheduledQueryName(name); err != nil {
		return nil, err
	}

	queryString := input.QueryString
	if queryString == "" {
		return nil, ErrValidationException
	}
	if err := validateQueryString(queryString); err != nil {
		return nil, err
	}

	scheduleConfig := parseScheduleConfiguration(input.ScheduleConfigRaw)
	if scheduleConfig == nil {
		return nil, ErrValidationException
	}
	if err := validateScheduleExpression(scheduleConfig.ScheduleExpression); err != nil {
		return nil, err
	}

	notificationConfig := parseNotificationConfiguration(input.NotificationConfigRaw)
	roleARN := input.RoleARN
	if roleARN == "" {
		return nil, ErrValidationException
	}
	kmsKeyID := input.KmsKeyID
	// StringValue2048 @length(1,2048) counts Unicode characters (no
	// pattern).
	if kmsKeyID != "" && utf8.RuneCountInString(kmsKeyID) > maxAmazonResourceName {
		return nil, ErrValidationException
	}
	errorReportConfig, err := parseErrorReportConfiguration(input.ErrorReportConfigRaw)
	if err != nil {
		return nil, err
	}
	targetConfig, err := parseTargetConfiguration(input.TargetConfigRaw)
	if err != nil {
		return nil, err
	}
	clientToken := input.ClientToken

	if roleARN != "" {
		validator := input.IAMValidator
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

	sq, err := stores.scheduledQueryStore.CreateScheduledQuery(
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

	if len(input.Tags) > 0 {
		if tagErr := stores.scheduledQueryStore.Tag(sq.ARN, input.Tags); tagErr != nil {
			logs.Warn("failed to tag scheduled query", logs.Err(tagErr), logs.String("arn", sq.ARN))
		}
	}

	return sq, nil
}

// deleteScheduledQueryCore validates the ARN and deletes the scheduled
// query together with its run records.
func (s *TimestreamQueryService) deleteScheduledQueryCore(stores *tsQueryStores, arn string) error {
	if arn == "" {
		return ErrValidationException
	}
	if err := validateAmazonResourceName(arn); err != nil {
		return err
	}

	name := stores.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return ErrValidationException
	}

	sq, err := stores.scheduledQueryStore.GetScheduledQuery(name)
	if err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return ErrResourceNotFound
		}
		return ErrInternalServer
	}

	// Delete all runs belonging to this scheduled query before deleting
	// the query itself, so no orphaned run records can survive. A run
	// deletion failure aborts the whole operation (fail-closed) rather
	// than silently leaving orphans behind.
	runs, runErr := stores.scheduledQueryRunStore.ListRuns(sq.ARN)
	if runErr != nil {
		logs.Error("Failed to list scheduled query runs for deletion", logs.String("query", name), logs.Err(runErr))
		return ErrInternalServer
	}
	for _, run := range runs {
		if err := stores.scheduledQueryRunStore.DeleteRun(run.ARN); err != nil {
			logs.Error("Failed to delete scheduled query run", logs.String("run", run.ARN), logs.Err(err))
			return ErrInternalServer
		}
	}

	err = stores.scheduledQueryStore.DeleteScheduledQuery(name)
	if err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return ErrResourceNotFound
		}
		return ErrInternalServer
	}

	return nil
}

// describeScheduledQueryCore validates the ARN and loads the scheduled
// query with its run history.
func (s *TimestreamQueryService) describeScheduledQueryCore(stores *tsQueryStores, arn string) (*DescribeScheduledQueryResult, error) {
	if arn == "" {
		return nil, ErrValidationException
	}
	if err := validateAmazonResourceName(arn); err != nil {
		return nil, err
	}

	name := stores.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return nil, ErrValidationException
	}

	sq, err := stores.scheduledQueryStore.GetScheduledQuery(name)
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
	result := &DescribeScheduledQueryResult{SQ: sq}
	runs, runErr := stores.scheduledQueryRunStore.ListRuns(sq.ARN)
	if runErr == nil && len(runs) > 0 {
		result.LastRun = runs[len(runs)-1]
		for _, r := range runs {
			if r.RunStatus == tsstore.ScheduleRunStatusFailed {
				result.FailedRuns = append(result.FailedRuns, r)
				result.FailedRunStatuses = append(result.FailedRunStatuses,
					tsstore.ScheduledQueryRunStatusFromTrigger(r.TriggerType, r.RunStatus))
			}
		}
	}

	return result, nil
}

// updateScheduledQueryCore validates the ARN and state and updates the
// scheduled query state.
func (s *TimestreamQueryService) updateScheduledQueryCore(stores *tsQueryStores, arn, state string) error {
	if arn == "" {
		return ErrValidationException
	}
	if err := validateAmazonResourceName(arn); err != nil {
		return err
	}

	name := stores.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return ErrValidationException
	}

	sqState := tsstore.ScheduledQueryStatus(state)
	if sqState != tsstore.ScheduledQueryStatusEnabled && sqState != tsstore.ScheduledQueryStatusDisabled {
		return ErrValidationException
	}

	if _, err := stores.scheduledQueryStore.UpdateScheduledQuery(name, sqState); err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return ErrResourceNotFound
		}
		return ErrInternalServer
	}

	return nil
}

// executeScheduledQueryCore validates the request and executes the
// scheduled query manually.
func (s *TimestreamQueryService) executeScheduledQueryCore(ctx context.Context, stores *tsQueryStores, input ExecuteScheduledQueryInput) (*tsstore.ScheduledQueryRun, error) {
	arn := input.ScheduledQueryARN
	if arn == "" {
		return nil, ErrValidationException
	}
	if err := validateAmazonResourceName(arn); err != nil {
		return nil, err
	}

	name := stores.arnBuilder.Timestream().ParseScheduledQueryName(arn)
	if name == "" {
		return nil, ErrValidationException
	}

	if ct := input.ClientToken; ct != "" {
		if err := validateClientToken(ct); err != nil {
			return nil, err
		}
	}

	sq, err := stores.scheduledQueryStore.GetScheduledQuery(name)
	if err != nil {
		if err == tsstore.ErrScheduledQueryNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, ErrInternalServer
	}

	// InvocationTime is REQUIRED per Smithy and must be present as a
	// string in RFC3339 form or as an epoch number; a missing value or a
	// present-but-malformed string is rejected rather than silently
	// replaced with the current time. QueryInsights is accepted but
	// currently not applied to the query execution (would require
	// extending the query executor to collect per-query insights).
	now := time.Now().UTC()
	if input.InvocationTimeStr == "" {
		if input.InvocationTimeRaw != nil {
			switch v := input.InvocationTimeRaw.(type) {
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
		if parsed, err := time.Parse(time.RFC3339Nano, input.InvocationTimeStr); err == nil {
			now = parsed.UTC()
		} else if parsed, err := time.Parse(time.RFC3339, input.InvocationTimeStr); err == nil {
			now = parsed.UTC()
		} else {
			return nil, ErrValidationException
		}
	}

	return s.executeScheduledQueryInternal(ctx, stores, sq, now, tsstore.TriggerTypeManual)
}

// listTagsForResourceCore validates the resource ARN, loads the tags
// through the shared tag dispatch config and applies MaxResults/NextToken
// pagination.
func (s *TimestreamQueryService) listTagsForResourceCore(ctx context.Context, stores *tsQueryStores, rawKey, maxResultsRaw, nextTokenRaw string) (*ListTagsForResourceResult, error) {
	cfg := s.tagHandlerConfig(stores)

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

	// The Core mirrors the shared HandleList flow, so it carries the same
	// pre-list resource validation instead of bypassing it.
	if cfg.ValidateResource != nil {
		if err := cfg.ValidateResource(ctx, resourceKey); err != nil {
			return nil, err
		}
	}

	if cfg.ListFunc == nil {
		return nil, ErrInternalServer
	}

	tags, err := cfg.ListFunc(ctx, resourceKey)
	if err != nil {
		return nil, err
	}

	maxResults := 0
	if maxResultsRaw != "" {
		val, atoiErr := strconv.Atoi(maxResultsRaw)
		if atoiErr != nil {
			return nil, ErrValidationException
		}
		if err := validateMaxResultsInRange(val, "MaxResults", rangeMaxTagsForResourceResult); err != nil {
			return nil, err
		}
		maxResults = val
	}

	offset := 0
	if nextTokenRaw != "" {
		val, atoiErr := strconv.Atoi(nextTokenRaw)
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

	result := &ListTagsForResourceResult{
		Tags: tags[offset:end],
	}
	if end < totalTags {
		result.NextToken = strconv.Itoa(end)
	}

	return result, nil
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

	listTags := func(ctx context.Context, resourceARN string) ([]tagutil.Tag, error) {
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
		ValidateResource: func(_ context.Context, resourceKey string) error {
			// Timestream Query tags scheduled queries plus the shared
			// database/table namespace; every kind resolves its resource
			// first so a tag against a nonexistent resource fails with the
			// modelled ResourceNotFoundException.
			if name := st.arnBuilder.Timestream().ParseScheduledQueryName(resourceKey); name != "" {
				if _, err := st.scheduledQueryStore.GetScheduledQuery(name); err != nil {
					if err == tsstore.ErrScheduledQueryNotFound {
						return ErrResourceNotFound
					}
					return ErrInternalServer
				}
				return nil
			}
			database := st.arnBuilder.Timestream().ParseDatabaseName(resourceKey)
			if database == "" {
				return ErrResourceNotFound
			}
			if table := st.arnBuilder.Timestream().ParseTableName(resourceKey); table != "" {
				if _, err := st.tableStore.GetTable(database, table); err != nil {
					if err == tsstore.ErrTableNotFound {
						return ErrResourceNotFound
					}
					return ErrInternalServer
				}
				return nil
			}
			if _, err := st.dbStore.GetDatabase(database); err != nil {
				if err == tsstore.ErrDatabaseNotFound {
					return ErrResourceNotFound
				}
				return ErrInternalServer
			}
			return nil
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagSlice []tagutil.Tag) error {
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
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return listTags(ctx, resourceKey)
		},
		FormatResponse: func(tagSlice []tagutil.Tag, _ string) (interface{}, error) {
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

// parseScheduleConfiguration parses a ScheduleConfiguration from its raw
// wire map. A nil map or a missing ScheduleExpression yields nil, which the
// Core rejects as a validation error.
func parseScheduleConfiguration(scheduleConfigRaw map[string]interface{}) *tsstore.ScheduleConfiguration {
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

// parseNotificationConfiguration parses a NotificationConfiguration from
// its raw wire map. A nil map or a missing SNS topic yields nil (the
// notification configuration is optional).
func parseNotificationConfiguration(notifConfigRaw map[string]interface{}) *tsstore.NotificationConfiguration {
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

// parseErrorReportConfiguration parses an ErrorReportConfiguration from its
// raw wire map. A nil map yields (nil, nil) (the configuration is
// optional); a present map with an invalid S3 configuration is a
// validation error.
func parseErrorReportConfiguration(errorReportRaw map[string]interface{}) (*tsstore.ErrorReportConfiguration, error) {
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

// parseTargetConfiguration parses a TargetConfiguration from its raw wire
// map. A nil map yields (nil, nil) (the configuration is optional); a
// present map without a TimestreamConfiguration is a validation error.
func parseTargetConfiguration(targetConfigRaw map[string]interface{}) (*tsstore.TargetConfiguration, error) {
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
