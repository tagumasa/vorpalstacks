// Package timestreamwrite provides Timestream Write service operations for vorpalstacks.
package timestreamwrite

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// CreateBatchLoadTask creates a new batch load task.
func (s *TimestreamWriteService) CreateBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var recordVersion int64
	if rv := request.GetParamCaseInsensitive(req.Parameters, "RecordVersion"); rv != "" {
		if parsed, ok := parseint64(rv); ok {
			recordVersion = parsed
		}
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createBatchLoadTaskCore(st, CreateBatchLoadTaskInput{
		TargetDatabaseName:      request.GetParamCaseInsensitive(req.Parameters, "TargetDatabaseName"),
		TargetTableName:         request.GetParamCaseInsensitive(req.Parameters, "TargetTableName"),
		DataSourceConfiguration: req.Parameters["DataSourceConfiguration"],
		DataModelConfiguration:  req.Parameters["DataModelConfiguration"],
		ReportConfiguration:     req.Parameters["ReportConfiguration"],
		RecordVersion:           recordVersion,
		ClientToken:             request.GetParamCaseInsensitive(req.Parameters, "ClientToken"),
		Region:                  reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TaskId": result.TaskId,
	}, nil
}

// DescribeBatchLoadTask returns information about a batch load task.
func (s *TimestreamWriteService) DescribeBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := s.describeBatchLoadTaskCore(st, DescribeBatchLoadTaskInput{
		TaskId: request.GetParamCaseInsensitive(req.Parameters, "TaskId"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"BatchLoadTaskDescription": s.formatBatchLoadTaskDescription(task),
	}, nil
}

// ListBatchLoadTasks returns a list of batch load tasks.
func (s *TimestreamWriteService) ListBatchLoadTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var taskStatus tsstore.BatchLoadStatus
	if status := request.GetParamCaseInsensitive(req.Parameters, "TaskStatus"); status != "" {
		taskStatus = tsstore.BatchLoadStatus(status)
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listBatchLoadTasksCore(st, ListBatchLoadTasksInput{
		TaskStatus: taskStatus,
		NextToken:  pagination.GetMarker(req.Parameters, "NextToken"),
		MaxResults: pagination.GetMaxItems(req.Parameters, 20, "MaxResults"),
	})
	if err != nil {
		return nil, err
	}

	taskList := make([]map[string]interface{}, 0)
	for _, task := range result.Tasks {
		taskList = append(taskList, s.formatBatchLoadTask(task))
	}

	response := map[string]interface{}{
		"BatchLoadTasks": taskList,
	}
	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}

// ResumeBatchLoadTask resumes a stopped batch load task.
func (s *TimestreamWriteService) ResumeBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.resumeBatchLoadTaskCore(st, ResumeBatchLoadTaskInput{
		TaskId: request.GetParamCaseInsensitive(req.Parameters, "TaskId"),
		Region: reqCtx.GetRegion(),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *TimestreamWriteService) formatBatchLoadTask(task *tsstore.BatchLoadTask) map[string]interface{} {
	response := map[string]interface{}{
		"TaskId":       task.TaskId,
		"TaskStatus":   string(task.TaskStatus),
		"DatabaseName": task.DatabaseName,
		"TableName":    task.TableName,
	}

	if !task.CreationTime.IsZero() {
		response["CreationTime"] = float64(task.CreationTime.UnixNano()) / 1e9
	}
	if !task.LastUpdatedTime.IsZero() {
		response["LastUpdatedTime"] = float64(task.LastUpdatedTime.UnixNano()) / 1e9
	}
	if !task.ResumableUntil.IsZero() {
		response["ResumableUntil"] = float64(task.ResumableUntil.UnixNano()) / 1e9
	}

	return response
}

func (s *TimestreamWriteService) formatBatchLoadTaskDescription(task *tsstore.BatchLoadTaskDescription) map[string]interface{} {
	response := map[string]interface{}{
		"TaskId":             task.TaskId,
		"TargetDatabaseName": task.TargetDatabaseName,
		"TargetTableName":    task.TargetTableName,
		"TaskStatus":         string(task.TaskStatus),
	}

	if !task.CreationTime.IsZero() {
		response["CreationTime"] = float64(task.CreationTime.UnixNano()) / 1e9
	}
	if !task.LastUpdatedTime.IsZero() {
		response["LastUpdatedTime"] = float64(task.LastUpdatedTime.UnixNano()) / 1e9
	}
	if !task.ResumableUntil.IsZero() {
		response["ResumableUntil"] = float64(task.ResumableUntil.UnixNano()) / 1e9
	}
	if task.ErrorMessage != "" {
		response["ErrorMessage"] = task.ErrorMessage
	}
	if task.RecordVersion > 0 {
		response["RecordVersion"] = task.RecordVersion
	}
	if task.DataSourceConfiguration != nil {
		response["DataSourceConfiguration"] = formatDataSourceConfiguration(task.DataSourceConfiguration)
	}
	if task.DataModelConfiguration != nil {
		response["DataModelConfiguration"] = formatDataModelConfiguration(task.DataModelConfiguration)
	}
	if task.ReportConfiguration != nil {
		response["ReportConfiguration"] = formatReportConfiguration(task.ReportConfiguration)
	}
	if task.ProgressReport != nil {
		response["ProgressReport"] = formatProgressReport(task.ProgressReport)
	}

	return response
}

func formatDataSourceConfiguration(config *tsstore.DataSourceConfiguration) map[string]interface{} {
	result := map[string]interface{}{
		"DataFormat": string(config.DataFormat),
	}

	if config.DataSourceS3Configuration != nil {
		result["DataSourceS3Configuration"] = map[string]interface{}{
			"BucketName":      config.DataSourceS3Configuration.BucketName,
			"ObjectKeyPrefix": config.DataSourceS3Configuration.ObjectKeyPrefix,
		}
	}

	if config.CsvConfiguration != nil {
		result["CsvConfiguration"] = map[string]interface{}{
			"ColumnSeparator": config.CsvConfiguration.ColumnSeparator,
			"EscapeChar":      config.CsvConfiguration.EscapeChar,
			"NullValue":       config.CsvConfiguration.NullValue,
			"QuoteChar":       config.CsvConfiguration.QuoteChar,
			"TrimWhiteSpace":  config.CsvConfiguration.TrimWhiteSpace,
		}
	}

	return result
}

func formatDataModelConfiguration(config *tsstore.DataModelConfiguration) map[string]interface{} {
	result := map[string]interface{}{}

	if config.DataModel != nil {
		dm := map[string]interface{}{}
		if len(config.DataModel.DimensionMappings) > 0 {
			var mappings []map[string]interface{}
			for _, m := range config.DataModel.DimensionMappings {
				mapping := map[string]interface{}{}
				if m.SourceColumn != nil {
					mapping["SourceColumn"] = m.SourceColumn.Name
				}
				if m.DestinationColumn != nil {
					mapping["DestinationColumn"] = m.DestinationColumn.Name
				}
				mappings = append(mappings, mapping)
			}
			dm["DimensionMappings"] = mappings
		}
		if config.DataModel.MeasureNameColumn != "" {
			dm["MeasureNameColumn"] = config.DataModel.MeasureNameColumn
		}
		if config.DataModel.TimeColumn != "" {
			dm["TimeColumn"] = config.DataModel.TimeColumn
		}
		if config.DataModel.TimeUnit != "" {
			dm["TimeUnit"] = string(config.DataModel.TimeUnit)
		}
		if len(config.DataModel.MixedMeasureMappings) > 0 {
			var mmmList []map[string]interface{}
			for _, mmm := range config.DataModel.MixedMeasureMappings {
				m := map[string]interface{}{}
				if mmm.MeasureName != "" {
					m["MeasureName"] = mmm.MeasureName
				}
				if mmm.SourceColumn != "" {
					m["SourceColumn"] = mmm.SourceColumn
				}
				if mmm.TargetMeasureName != "" {
					m["TargetMeasureName"] = mmm.TargetMeasureName
				}
				m["MeasureValueType"] = string(mmm.MeasureValueMeasureValueType)
				if len(mmm.MultiMeasureAttributeMappings) > 0 {
					m["MultiMeasureAttributeMappings"] = formatMultiMeasureAttributeMappings(mmm.MultiMeasureAttributeMappings)
				}
				mmmList = append(mmmList, m)
			}
			dm["MixedMeasureMappings"] = mmmList
		}
		if config.DataModel.MultiMeasureMappings != nil {
			mmm := map[string]interface{}{}
			if config.DataModel.MultiMeasureMappings.TargetMultiMeasureName != "" {
				mmm["TargetMultiMeasureName"] = config.DataModel.MultiMeasureMappings.TargetMultiMeasureName
			}
			if len(config.DataModel.MultiMeasureMappings.MultiMeasureAttributeMappings) > 0 {
				mmm["MultiMeasureAttributeMappings"] = formatMultiMeasureAttributeMappings(config.DataModel.MultiMeasureMappings.MultiMeasureAttributeMappings)
			}
			dm["MultiMeasureMappings"] = mmm
		}
		result["DataModel"] = dm
	}

	if config.DataModelS3Configuration != nil {
		result["DataModelS3Configuration"] = map[string]interface{}{
			"BucketName": config.DataModelS3Configuration.BucketName,
			"ObjectKey":  config.DataModelS3Configuration.ObjectKey,
		}
	}

	return result
}

// formatMultiMeasureAttributeMappings formats attribute mappings for output.
func formatMultiMeasureAttributeMappings(attrs []tsstore.MultiMeasureAttributeMapping) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(attrs))
	for _, am := range attrs {
		m := map[string]interface{}{}
		if am.SourceColumn != nil {
			m["SourceColumn"] = am.SourceColumn.Name
		}
		if am.TargetMultiMeasureAttributeName != "" {
			m["TargetMultiMeasureAttributeName"] = am.TargetMultiMeasureAttributeName
		}
		m["MeasureValueType"] = string(am.MeasureValueMeasureValueType)
		result = append(result, m)
	}
	return result
}

func formatReportConfiguration(config *tsstore.ReportConfiguration) map[string]interface{} {
	result := map[string]interface{}{}

	if config.ReportS3Configuration != nil {
		result["ReportS3Configuration"] = map[string]interface{}{
			"BucketName":       config.ReportS3Configuration.BucketName,
			"EncryptionOption": string(config.ReportS3Configuration.EncryptionOption),
			"KmsKeyId":         config.ReportS3Configuration.KmsKeyId,
			"ObjectKeyPrefix":  config.ReportS3Configuration.ObjectKeyPrefix,
		}
	}

	return result
}

func formatProgressReport(report *tsstore.BatchLoadProgressReport) map[string]interface{} {
	return map[string]interface{}{
		"BytesMetered":            report.BytesMetered,
		"FileFailures":            report.FileFailures,
		"ParseFailures":           report.ParseFailures,
		"RecordIngestionFailures": report.RecordIngestionFailures,
		"RecordsIngested":         report.RecordsIngested,
		"RecordsProcessed":        report.RecordsProcessed,
	}
}
