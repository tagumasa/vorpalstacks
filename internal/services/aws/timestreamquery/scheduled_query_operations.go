package timestreamquery

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// CreateScheduledQuery creates a new scheduled query.
func (s *TimestreamQueryService) CreateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sq, err := s.createScheduledQueryCore(ctx, st, CreateScheduledQueryInput{
		Name:                  request.GetParamCaseInsensitive(req.Parameters, "Name"),
		QueryString:           request.GetParamCaseInsensitive(req.Parameters, "QueryString"),
		ScheduleConfigRaw:     request.GetMapParamCaseInsensitive(req.Parameters, "ScheduleConfiguration"),
		NotificationConfigRaw: request.GetMapParamCaseInsensitive(req.Parameters, "NotificationConfiguration"),
		ErrorReportConfigRaw:  request.GetMapParamCaseInsensitive(req.Parameters, "ErrorReportConfiguration"),
		TargetConfigRaw:       request.GetMapParamCaseInsensitive(req.Parameters, "TargetConfiguration"),
		RoleARN:               request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryExecutionRoleArn"),
		KmsKeyID:              request.GetParamCaseInsensitive(req.Parameters, "KmsKeyId"),
		ClientToken:           request.GetParamCaseInsensitive(req.Parameters, "ClientToken"),
		Tags:                  tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")),
		IAMValidator:          reqCtx.GetIAMValidator(),
	})
	if err != nil {
		return nil, err
	}

	return s.formatScheduledQueryResponse(sq), nil
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
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteScheduledQueryCore(st, request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeScheduledQuery returns the details of a scheduled query.
func (s *TimestreamQueryService) DescribeScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeScheduledQueryCore(st, request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ScheduledQuery": s.formatScheduledQueryDescriptionResponse(result),
	}, nil
}

// ListScheduledQueries returns a list of scheduled queries.
func (s *TimestreamQueryService) ListScheduledQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.ListScheduledQueriesCore(st, ListScheduledQueriesInput{
		MaxResultsRaw: request.GetParamCaseInsensitive(req.Parameters, "MaxResults"),
		NextTokenRaw:  request.GetParamCaseInsensitive(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	var scheduledQueries []map[string]interface{}
	for _, summary := range result.Summaries {
		scheduledQueries = append(scheduledQueries, formatScheduledQuerySummaryResponse(summary))
	}

	resp := map[string]interface{}{
		"ScheduledQueries": scheduledQueries,
	}

	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}

	return resp, nil
}

// UpdateScheduledQuery updates an existing scheduled query.
func (s *TimestreamQueryService) UpdateScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.updateScheduledQueryCore(st,
		request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn"),
		request.GetParamCaseInsensitive(req.Parameters, "State")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ExecuteScheduledQuery executes a scheduled query immediately.
func (s *TimestreamQueryService) ExecuteScheduledQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	run, err := s.executeScheduledQueryCore(ctx, st, ExecuteScheduledQueryInput{
		ScheduledQueryARN: request.GetParamCaseInsensitive(req.Parameters, "ScheduledQueryArn"),
		ClientToken:       request.GetParamCaseInsensitive(req.Parameters, "ClientToken"),
		InvocationTimeStr: request.GetParamCaseInsensitive(req.Parameters, "InvocationTime"),
		InvocationTimeRaw: req.Parameters["InvocationTime"],
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ScheduledQueryRunArn": run.ARN,
	}, nil
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

	result, err := s.listTagsForResourceCore(ctx, st,
		tagutil.GetResourceKey(req.Parameters, cfg.Param),
		request.GetParamCaseInsensitive(req.Parameters, "MaxResults"),
		request.GetParamCaseInsensitive(req.Parameters, "NextToken"))
	if err != nil {
		return nil, err
	}

	tagResponse := tagutil.ToResponseWithKeyNames(result.Tags, cfg.Param.TagKeyName, cfg.Param.TagValueName)

	resp := map[string]interface{}{
		cfg.Param.TagsParam: tagResponse,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
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

// formatScheduledQueryDescriptionResponse converts a
// DescribeScheduledQueryResult DTO into the JSON map representation of the
// ScheduledQueryDescription shape. The trigger-aware run statuses arrive
// precomputed in the DTO, so this formatter makes no store calls.
func (s *TimestreamQueryService) formatScheduledQueryDescriptionResponse(result *DescribeScheduledQueryResult) map[string]interface{} {
	sq := result.SQ
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
		if result.LastRun != nil {
			lastRun := result.LastRun
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
	if len(result.FailedRuns) > 0 {
		recentlyFailed := make([]map[string]interface{}, 0, len(result.FailedRuns))
		for i, fr := range result.FailedRuns {
			summary := map[string]interface{}{
				"RunStatus": result.FailedRunStatuses[i],
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
