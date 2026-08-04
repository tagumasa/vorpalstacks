// Package timestreamwrite provides Timestream Write service operations for vorpalstacks.
package timestreamwrite

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/csv"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	storecommon "vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// CreateBatchLoadTask creates a new batch load task.
func (s *TimestreamWriteService) CreateBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	targetDatabaseName := request.GetParamCaseInsensitive(req.Parameters, "TargetDatabaseName")
	if targetDatabaseName == "" {
		return nil, ErrValidationException
	}

	targetTableName := request.GetParamCaseInsensitive(req.Parameters, "TargetTableName")
	if targetTableName == "" {
		return nil, ErrValidationException
	}

	dataSourceConfig := parseDataSourceConfiguration(req.Parameters["DataSourceConfiguration"])
	if dataSourceConfig == nil {
		return nil, ErrValidationException
	}
	// Validate DataFormat enum (Smithy BatchLoadDataFormat: CSV only).
	if dataSourceConfig.DataFormat != "" && dataSourceConfig.DataFormat != tsstore.BatchLoadDataFormatCsv {
		return nil, ErrValidationException
	}

	reportConfig := parseReportConfiguration(req.Parameters["ReportConfiguration"])
	if reportConfig == nil {
		return nil, ErrValidationException
	}

	dataModelConfig := parseDataModelConfiguration(req.Parameters["DataModelConfiguration"])
	// DimensionMappings is REQUIRED when DataModel is provided (Smithy).
	if dataModelConfig != nil && dataModelConfig.DataModel != nil && len(dataModelConfig.DataModel.DimensionMappings) == 0 {
		return nil, ErrValidationException
	}

	var recordVersion int64
	if rv := request.GetParamCaseInsensitive(req.Parameters, "RecordVersion"); rv != "" {
		if parsed, ok := parseint64(rv); ok {
			recordVersion = parsed
		}
	}

	clientToken := request.GetParamCaseInsensitive(req.Parameters, "ClientToken")
	// Validate ClientToken length (Smithy ClientRequestToken: 1-64).
	if clientToken != "" && !validateClientToken(clientToken) {
		return nil, ErrValidationException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// @idempotencyToken — if ClientToken matches an existing task,
	// return that task instead of creating a new one. The AWS SDK
	// auto-generates a UUID for the ClientToken (Smithy idempotencyToken
	// trait), but the TaskId is always server-generated (uppercase
	// alphanumeric, matching BatchLoadTaskId pattern ^[A-Z0-9]+$).
	if clientToken != "" {
		if existingTask, gerr := st.batchLoadStore.FindByClientToken(clientToken); gerr == nil {
			return map[string]interface{}{
				"TaskId": existingTask.TaskId,
			}, nil
		}
	}

	taskId := generateTaskId()
	task, err := st.batchLoadStore.CreateBatchLoadTask(taskId, clientToken, targetDatabaseName, targetTableName, dataSourceConfig, dataModelConfig, reportConfig, recordVersion)
	if err != nil {
		if err == tsstore.ErrBatchLoadTaskAlreadyExists {
			return nil, ErrConflictException
		}
		if err == tsstore.ErrTableNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	s.batchWg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(s.batchCtx, 5*time.Minute)
		defer cancel()
		defer s.batchWg.Done()
		defer func() { resilience.RecoverPanic("timestreamwrite batch load") }()
		s.executeBatchLoad(ctx, st, taskId, reqCtx.GetRegion())
	}()

	return map[string]interface{}{
		"TaskId": task.TaskId,
	}, nil
}

// DescribeBatchLoadTask returns information about a batch load task.
func (s *TimestreamWriteService) DescribeBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "TaskId")
	if taskId == "" {
		return nil, ErrValidationException
	}
	// Validate BatchLoadTaskId pattern (Smithy: ^[A-Z0-9]+$, 3-32 chars).
	if !batchLoadTaskIdRegex.MatchString(taskId) {
		return nil, ErrValidationException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	task, err := st.batchLoadStore.GetBatchLoadTask(taskId)
	if err != nil {
		if err == tsstore.ErrBatchLoadTaskNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
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

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 20, "MaxResults")
	if maxResults > maxListBatchLoadTasksResults {
		maxResults = maxListBatchLoadTasksResults
	}

	opts := storecommon.ListOptions{MaxItems: maxResults}
	if nextToken != "" {
		opts.Marker = nextToken
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := st.batchLoadStore.ListBatchLoadTasks(taskStatus, opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	taskList := make([]map[string]interface{}, 0)
	for _, task := range result.Items {
		taskList = append(taskList, s.formatBatchLoadTask(task))
	}

	response := map[string]interface{}{
		"BatchLoadTasks": taskList,
	}
	pagination.SetNextToken(response, "NextToken", result.NextMarker)

	return response, nil
}

// ResumeBatchLoadTask resumes a stopped batch load task.
func (s *TimestreamWriteService) ResumeBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "TaskId")
	if taskId == "" {
		return nil, ErrValidationException
	}
	// Validate BatchLoadTaskId pattern (Smithy: ^[A-Z0-9]+$, 3-32 chars).
	if !batchLoadTaskIdRegex.MatchString(taskId) {
		return nil, ErrValidationException
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	err = st.batchLoadStore.TransitionBatchLoadTaskStatus(taskId, []tsstore.BatchLoadStatus{
		tsstore.BatchLoadStatusProgressStopped,
		tsstore.BatchLoadStatusPendingResume,
		tsstore.BatchLoadStatusFailed,
	}, tsstore.BatchLoadStatusPendingResume)
	if err != nil {
		if err == tsstore.ErrBatchLoadTaskNotFound {
			return nil, ErrResourceNotFound
		}
		if err == tsstore.ErrInvalidTransition {
			return nil, ErrValidationException
		}
		return nil, s.mapStoreError(err)
	}

	s.batchWg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(s.batchCtx, 5*time.Minute)
		defer cancel()
		defer s.batchWg.Done()
		defer func() { resilience.RecoverPanic("timestreamwrite batch load resume") }()
		s.executeBatchLoad(ctx, st, taskId, reqCtx.GetRegion())
	}()

	return response.EmptyResponse(), nil
}

func (s *TimestreamWriteService) executeBatchLoad(ctx context.Context, st *tsWriteStores, taskId, region string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	defer func() {
		if ctx.Err() != nil {
			if err := st.batchLoadStore.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusFailed, "Batch load cancelled"); err != nil {
				logs.Error("Failed to update batch load task status to FAILED", logs.String("taskId", taskId), logs.Err(err))
			}
		}
	}()

	if s.s3Invoker == nil {
		if err := st.batchLoadStore.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusFailed, "S3 service not available for batch load"); err != nil {
			logs.Error("Failed to update batch load task status", logs.Err(err))
		}
		return
	}

	task, err := st.batchLoadStore.GetBatchLoadTask(taskId)
	if err != nil {
		logs.Error("Failed to get batch load task", logs.String("taskId", taskId), logs.Err(err))
		return
	}

	if task.DataSourceConfiguration == nil || task.DataSourceConfiguration.DataSourceS3Configuration == nil {
		if err := st.batchLoadStore.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusFailed, "Missing S3 data source configuration"); err != nil {
			logs.Error("Failed to update batch load task status", logs.Err(err))
		}
		return
	}

	s3Config := task.DataSourceConfiguration.DataSourceS3Configuration
	csvConfig := task.DataSourceConfiguration.CsvConfiguration
	dataModel := task.DataModelConfiguration

	if err := st.batchLoadStore.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusInProgress, ""); err != nil {
		logs.Error("Failed to update batch load task status to IN_PROGRESS", logs.Err(err))
		return
	}

	keys, err := s.s3Invoker.ListObjects(ctx, region, s3Config.BucketName, s3Config.ObjectKeyPrefix, 0)
	if err != nil {
		if err := st.batchLoadStore.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusFailed, fmt.Sprintf("Failed to list S3 objects: %v", err)); err != nil {
			logs.Error("Failed to update batch load task status", logs.Err(err))
		}
		return
	}

	processedSet := make(map[string]bool, len(task.ProcessedS3Keys))
	for _, k := range task.ProcessedS3Keys {
		processedSet[k] = true
	}

	progress := task.ProgressReport
	if progress == nil {
		progress = &tsstore.BatchLoadProgressReport{}
	}

	for _, key := range keys {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if processedSet[key] {
			continue
		}

		data, err := s.s3Invoker.GetObject(ctx, region, s3Config.BucketName, key, 0)
		if err != nil {
			progress.FileFailures++
			logs.Warn("Failed to get S3 object, continuing",
				logs.String("taskId", taskId), logs.String("key", key), logs.Err(err))
			if err := st.batchLoadStore.UpdateBatchLoadTaskProgress(taskId, progress); err != nil {
				logs.Error("Failed to update progress", logs.Err(err))
			}
			continue
		}

		progress.BytesMetered += int64(len(data))

		records, parseFailures := parseCSVToRecords(data, csvConfig, dataModel)
		progress.ParseFailures += int64(parseFailures)
		progress.RecordsProcessed += int64(len(records))

		for i := 0; i < len(records); i += 100 {
			end := i + 100
			if end > len(records) {
				end = len(records)
			}
			batch := records[i:end]

			rejected, werr := st.recordStore.WriteRecords(task.TargetDatabaseName, task.TargetTableName, batch)
			if werr != nil {
				logs.Warn("Failed to write records batch",
					logs.String("taskId", taskId), logs.String("key", key), logs.Err(werr))
				progress.RecordIngestionFailures += int64(len(batch))
			} else {
				progress.RecordsIngested += int64(len(batch) - len(rejected))
				progress.RecordIngestionFailures += int64(len(rejected))
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}

		processedSet[key] = true
		task.ProcessedS3Keys = append(task.ProcessedS3Keys, key)
		if err := st.batchLoadStore.UpdateBatchLoadTaskProgress(taskId, progress); err != nil {
			logs.Error("Failed to update progress", logs.Err(err))
		}
		if err := st.batchLoadStore.SaveProcessedKeys(taskId, task.ProcessedS3Keys); err != nil {
			logs.Warn("Failed to save processed keys", logs.String("taskId", taskId), logs.Err(err))
		}
	}

	if err := st.batchLoadStore.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusSucceeded, ""); err != nil {
		logs.Error("Failed to update batch load task status to SUCCEEDED", logs.Err(err))
	}
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

func parseDataSourceConfiguration(config interface{}) *tsstore.DataSourceConfiguration {
	if config == nil {
		return nil
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil
	}

	result := &tsstore.DataSourceConfiguration{
		DataFormat: tsstore.BatchLoadDataFormat("CSV"),
	}

	if format, ok := configMap["DataFormat"].(string); ok {
		result.DataFormat = tsstore.BatchLoadDataFormat(format)
	}

	if s3Config, ok := configMap["DataSourceS3Configuration"].(map[string]interface{}); ok {
		result.DataSourceS3Configuration = &tsstore.DataSourceS3Configuration{}
		if bucket, ok := s3Config["BucketName"].(string); ok {
			result.DataSourceS3Configuration.BucketName = bucket
		}
		if prefix, ok := s3Config["ObjectKeyPrefix"].(string); ok {
			result.DataSourceS3Configuration.ObjectKeyPrefix = prefix
		}
	}

	if csvConfig, ok := configMap["CsvConfiguration"].(map[string]interface{}); ok {
		result.CsvConfiguration = &tsstore.CsvConfiguration{}
		if sep, ok := csvConfig["ColumnSeparator"].(string); ok {
			result.CsvConfiguration.ColumnSeparator = sep
		}
		if esc, ok := csvConfig["EscapeChar"].(string); ok {
			result.CsvConfiguration.EscapeChar = esc
		}
		if null, ok := csvConfig["NullValue"].(string); ok {
			result.CsvConfiguration.NullValue = null
		}
		if quote, ok := csvConfig["QuoteChar"].(string); ok {
			result.CsvConfiguration.QuoteChar = quote
		}
		if trim, ok := csvConfig["TrimWhiteSpace"].(bool); ok {
			result.CsvConfiguration.TrimWhiteSpace = trim
		}
	}

	return result
}

func parseReportConfiguration(config interface{}) *tsstore.ReportConfiguration {
	if config == nil {
		return nil
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil
	}

	result := &tsstore.ReportConfiguration{}

	if s3Config, ok := configMap["ReportS3Configuration"].(map[string]interface{}); ok {
		result.ReportS3Configuration = &tsstore.ReportS3Configuration{}
		if bucket, ok := s3Config["BucketName"].(string); ok {
			result.ReportS3Configuration.BucketName = bucket
		}
		if enc, ok := s3Config["EncryptionOption"].(string); ok {
			result.ReportS3Configuration.EncryptionOption = tsstore.S3EncryptionOption(enc)
		}
		if kms, ok := s3Config["KmsKeyId"].(string); ok {
			result.ReportS3Configuration.KmsKeyId = kms
		}
		if prefix, ok := s3Config["ObjectKeyPrefix"].(string); ok {
			result.ReportS3Configuration.ObjectKeyPrefix = prefix
		}
	}

	return result
}

func parseDataModelConfiguration(config interface{}) *tsstore.DataModelConfiguration {
	if config == nil {
		return nil
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil
	}

	result := &tsstore.DataModelConfiguration{}

	if dataModel, ok := configMap["DataModel"].(map[string]interface{}); ok {
		result.DataModel = parseDataModel(dataModel)
	}

	if s3Config, ok := configMap["DataModelS3Configuration"].(map[string]interface{}); ok {
		result.DataModelS3Configuration = &tsstore.DataModelS3Configuration{}
		if bucket, ok := s3Config["BucketName"].(string); ok {
			result.DataModelS3Configuration.BucketName = bucket
		}
		if key, ok := s3Config["ObjectKey"].(string); ok {
			result.DataModelS3Configuration.ObjectKey = key
		}
	}

	return result
}

func parseDataModel(dataModel map[string]interface{}) *tsstore.DataModel {
	result := &tsstore.DataModel{}

	if mappings, ok := dataModel["DimensionMappings"].([]interface{}); ok {
		for _, m := range mappings {
			if mapping, ok := m.(map[string]interface{}); ok {
				dm := tsstore.DimensionMapping{}
				if src, ok := mapping["SourceColumn"].(string); ok {
					dm.SourceColumn = &tsstore.SourceColumn{Name: src}
				} else if srcMap, ok := mapping["SourceColumn"].(map[string]interface{}); ok {
					if name, ok := srcMap["Name"].(string); ok {
						dm.SourceColumn = &tsstore.SourceColumn{Name: name}
					}
				}
				if dst, ok := mapping["DestinationColumn"].(string); ok {
					dm.DestinationColumn = &tsstore.DestinationColumn{Name: dst}
				} else if dstMap, ok := mapping["DestinationColumn"].(map[string]interface{}); ok {
					if name, ok := dstMap["Name"].(string); ok {
						dm.DestinationColumn = &tsstore.DestinationColumn{Name: name}
					}
				}
				result.DimensionMappings = append(result.DimensionMappings, dm)
			}
		}
	}

	if col, ok := dataModel["MeasureNameColumn"].(string); ok {
		result.MeasureNameColumn = col
	}
	if col, ok := dataModel["TimeColumn"].(string); ok {
		result.TimeColumn = col
	}
	if unit, ok := dataModel["TimeUnit"].(string); ok {
		result.TimeUnit = tsstore.TimeUnit(unit)
	}

	if mmmList, ok := dataModel["MixedMeasureMappings"].([]interface{}); ok {
		for _, m := range mmmList {
			mapping, ok := m.(map[string]interface{})
			if !ok {
				continue
			}
			mmm := tsstore.MixedMeasureMapping{}
			if name, ok := mapping["MeasureName"].(string); ok {
				mmm.MeasureName = name
			}
			if col, ok := mapping["SourceColumn"].(string); ok {
				mmm.SourceColumn = col
			}
			if name, ok := mapping["TargetMeasureName"].(string); ok {
				mmm.TargetMeasureName = name
			}
			if mvt, ok := mapping["MeasureValueType"].(string); ok {
				mmm.MeasureValueMeasureValueType = tsstore.MeasureValueType(mvt)
			}
			if attrs, ok := mapping["MultiMeasureAttributeMappings"].([]interface{}); ok {
				mmm.MultiMeasureAttributeMappings = parseMultiMeasureAttributeMappings(attrs)
			}
			result.MixedMeasureMappings = append(result.MixedMeasureMappings, mmm)
		}
	}

	if mmmMap, ok := dataModel["MultiMeasureMappings"].(map[string]interface{}); ok {
		mmm := &tsstore.MultiMeasureMappings{}
		if name, ok := mmmMap["TargetMultiMeasureName"].(string); ok {
			mmm.TargetMultiMeasureName = name
		}
		if attrs, ok := mmmMap["MultiMeasureAttributeMappings"].([]interface{}); ok {
			mmm.MultiMeasureAttributeMappings = parseMultiMeasureAttributeMappings(attrs)
		}
		if len(mmm.MultiMeasureAttributeMappings) > 0 || mmm.TargetMultiMeasureName != "" {
			result.MultiMeasureMappings = mmm
		}
	}

	return result
}

// parseMultiMeasureAttributeMappings parses a list of MultiMeasureAttributeMapping
// from the request. SourceColumn is REQUIRED per Smithy.
func parseMultiMeasureAttributeMappings(attrs []interface{}) []tsstore.MultiMeasureAttributeMapping {
	result := make([]tsstore.MultiMeasureAttributeMapping, 0, len(attrs))
	for _, a := range attrs {
		attrMap, ok := a.(map[string]interface{})
		if !ok {
			continue
		}
		am := tsstore.MultiMeasureAttributeMapping{}
		if src, ok := attrMap["SourceColumn"].(string); ok {
			am.SourceColumn = &tsstore.SourceColumn{Name: src}
		} else if srcMap, ok := attrMap["SourceColumn"].(map[string]interface{}); ok {
			if name, ok := srcMap["Name"].(string); ok {
				am.SourceColumn = &tsstore.SourceColumn{Name: name}
			}
		}
		if name, ok := attrMap["TargetMultiMeasureAttributeName"].(string); ok {
			am.TargetMultiMeasureAttributeName = name
		}
		if mvt, ok := attrMap["MeasureValueType"].(string); ok {
			am.MeasureValueMeasureValueType = tsstore.MeasureValueType(mvt)
		}
		result = append(result, am)
	}
	return result
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

func generateTaskId() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		num, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

// parseCSVToRecords parses CSV data into Timestream records using the DataModel
// configuration. Returns the parsed records and the number of parse failures.
func parseCSVToRecords(data []byte, csvConfig *tsstore.CsvConfiguration, dataModelConfig *tsstore.DataModelConfiguration) ([]tsstore.Record, int) {
	reader := csv.NewReader(bytes.NewReader(data))

	if csvConfig != nil {
		if csvConfig.ColumnSeparator != "" {
			runes := []rune(csvConfig.ColumnSeparator)
			if len(runes) > 0 {
				reader.Comma = runes[0]
			}
		}
		reader.TrimLeadingSpace = csvConfig.TrimWhiteSpace
	}

	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, 0
	}
	if len(rows) == 0 {
		return nil, 0
	}

	header := rows[0]
	colIndex := make(map[string]int, len(header))
	for i, col := range header {
		colName := strings.TrimSpace(col)
		if csvConfig != nil && csvConfig.TrimWhiteSpace {
			colName = strings.TrimSpace(colName)
		}
		colIndex[colName] = i
	}

	var dataModel *tsstore.DataModel
	if dataModelConfig != nil {
		dataModel = dataModelConfig.DataModel
	}

	var timeCol, measureNameCol string
	var timeUnit tsstore.TimeUnit
	var dimensionMappings []tsstore.DimensionMapping
	var multiMeasureMappings *tsstore.MultiMeasureMappings
	var mixedMeasureMappings []tsstore.MixedMeasureMapping
	if dataModel != nil {
		timeCol = dataModel.TimeColumn
		measureNameCol = dataModel.MeasureNameColumn
		timeUnit = dataModel.TimeUnit
		dimensionMappings = dataModel.DimensionMappings
		multiMeasureMappings = dataModel.MultiMeasureMappings
		mixedMeasureMappings = dataModel.MixedMeasureMappings
	}

	records := make([]tsstore.Record, 0, len(rows)-1)
	parseFailures := 0

	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]

		getCol := func(name string) string {
			if idx, ok := colIndex[name]; ok && idx < len(row) {
				val := row[idx]
				if csvConfig != nil && csvConfig.TrimWhiteSpace {
					val = strings.TrimSpace(val)
				}
				if csvConfig != nil && csvConfig.NullValue != "" && val == csvConfig.NullValue {
					return ""
				}
				return val
			}
			return ""
		}

		record := tsstore.Record{}

		for _, dm := range dimensionMappings {
			srcName := ""
			if dm.SourceColumn != nil {
				srcName = dm.SourceColumn.Name
			}
			dstName := ""
			if dm.DestinationColumn != nil {
				dstName = dm.DestinationColumn.Name
			}
			if srcName == "" {
				srcName = dstName
			}
			val := getCol(srcName)
			if dstName == "" {
				dstName = srcName
			}
			record.Dimensions = append(record.Dimensions, tsstore.Dimension{
				Name:  dstName,
				Value: val,
			})
		}

		if timeCol != "" {
			record.Time = getCol(timeCol)
		}
		if timeUnit != "" {
			record.TimeUnit = timeUnit
		}

		// Measure parsing: MultiMeasureMappings, MixedMeasureMappings, or
		// plain MeasureNameColumn. These three modes are mutually exclusive.
		if multiMeasureMappings != nil && len(multiMeasureMappings.MultiMeasureAttributeMappings) > 0 {
			// MultiMeasureMappings: all attribute columns combine into one
			// multi-measure record. The record MeasureName comes from
			// MeasureNameColumn (if set) or TargetMultiMeasureName.
			if measureNameCol != "" {
				record.MeasureName = getCol(measureNameCol)
			}
			if record.MeasureName == "" && multiMeasureMappings.TargetMultiMeasureName != "" {
				record.MeasureName = multiMeasureMappings.TargetMultiMeasureName
			}
			for _, am := range multiMeasureMappings.MultiMeasureAttributeMappings {
				srcName := ""
				if am.SourceColumn != nil {
					srcName = am.SourceColumn.Name
				}
				val := getCol(srcName)
				name := am.TargetMultiMeasureAttributeName
				if name == "" {
					name = srcName
				}
				record.MeasureValues = append(record.MeasureValues, tsstore.MeasureValue{
					Name:  name,
					Value: val,
					Type:  am.MeasureValueMeasureValueType,
				})
			}
			if record.MeasureValueType == "" {
				record.MeasureValueType = tsstore.MeasureValueTypeMulti
			}
		} else if len(mixedMeasureMappings) > 0 {
			// MixedMeasureMappings: heterogeneous measures where each row's
			// MeasureNameColumn value selects which mapping applies.
			measureNameVal := ""
			if measureNameCol != "" {
				measureNameVal = getCol(measureNameCol)
			}
			for _, mmm := range mixedMeasureMappings {
				if measureNameVal != "" && mmm.MeasureName != measureNameVal && mmm.TargetMeasureName != measureNameVal {
					continue
				}
				targetName := mmm.TargetMeasureName
				if targetName == "" {
					targetName = mmm.MeasureName
				}
				if targetName == "" {
					targetName = measureNameVal
				}
				record.MeasureName = targetName

				if mmm.SourceColumn != "" {
					record.MeasureValue = getCol(mmm.SourceColumn)
					record.MeasureValueType = mmm.MeasureValueMeasureValueType
				} else if len(mmm.MultiMeasureAttributeMappings) > 0 {
					for _, am := range mmm.MultiMeasureAttributeMappings {
						srcName := ""
						if am.SourceColumn != nil {
							srcName = am.SourceColumn.Name
						}
						val := getCol(srcName)
						name := am.TargetMultiMeasureAttributeName
						if name == "" {
							name = srcName
						}
						record.MeasureValues = append(record.MeasureValues, tsstore.MeasureValue{
							Name:  name,
							Value: val,
							Type:  am.MeasureValueMeasureValueType,
						})
					}
					if record.MeasureValueType == "" {
						record.MeasureValueType = tsstore.MeasureValueTypeMulti
					}
				}
				break
			}
		} else if measureNameCol != "" {
			record.MeasureName = getCol(measureNameCol)
		}

		if record.Time == "" && len(record.Dimensions) == 0 {
			parseFailures++
			continue
		}

		records = append(records, record)
	}

	return records, parseFailures
}

func parseint64(s string) (int64, bool) {
	result, err := strconv.ParseInt(s, 10, 64)
	return result, err == nil
}
