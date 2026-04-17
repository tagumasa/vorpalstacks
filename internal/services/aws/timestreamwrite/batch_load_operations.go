// Package timestreamwrite provides Timestream Write service operations for vorpalstacks.
package timestreamwrite

import (
	"context"
	cryptorand "crypto/rand"
	"math/big"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// CreateBatchLoadTask creates a new batch load task.
func (s *TimestreamWriteService) CreateBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	targetDatabaseName := request.GetParamCaseInsensitive(req.Parameters, "TargetDatabaseName")
	if targetDatabaseName == "" {
		return nil, ErrInvalidParameter
	}

	targetTableName := request.GetParamCaseInsensitive(req.Parameters, "TargetTableName")
	if targetTableName == "" {
		return nil, ErrInvalidParameter
	}

	dataSourceConfig := parseDataSourceConfiguration(req.Parameters["DataSourceConfiguration"])
	if dataSourceConfig == nil {
		return nil, ErrInvalidParameter
	}

	reportConfig := parseReportConfiguration(req.Parameters["ReportConfiguration"])
	if reportConfig == nil {
		return nil, ErrInvalidParameter
	}

	dataModelConfig := parseDataModelConfiguration(req.Parameters["DataModelConfiguration"])

	var recordVersion int64
	if rv := request.GetParamCaseInsensitive(req.Parameters, "RecordVersion"); rv != "" {
		if parsed, ok := parseint64(rv); ok {
			recordVersion = parsed
		}
	}

	clientToken := request.GetParamCaseInsensitive(req.Parameters, "ClientToken")
	taskId := clientToken
	if taskId == "" {
		taskId = generateTaskId()
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	task, err := st.batchLoadStore.CreateBatchLoadTask(taskId, targetDatabaseName, targetTableName, dataSourceConfig, dataModelConfig, reportConfig, recordVersion)
	if err != nil {
		if err == tsstore.ErrBatchLoadTaskAlreadyExists {
			return nil, ErrResourceAlreadyExists
		}
		if err == tsstore.ErrTableNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	s.batchWg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		defer s.batchWg.Done()
		s.simulateBatchLoad(ctx, st.batchLoadStore, taskId)
	}()

	return map[string]interface{}{
		"TaskId": task.TaskId,
	}, nil
}

// DescribeBatchLoadTask returns information about a batch load task.
func (s *TimestreamWriteService) DescribeBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "TaskId")
	if taskId == "" {
		return nil, ErrInvalidParameter
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

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tasks, err := st.batchLoadStore.ListBatchLoadTasks(taskStatus)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	taskList := make([]map[string]interface{}, 0)
	for _, task := range tasks {
		taskList = append(taskList, s.formatBatchLoadTask(task))
	}

	response := map[string]interface{}{
		"BatchLoadTasks": taskList,
	}

	return response, nil
}

// ResumeBatchLoadTask resumes a stopped batch load task.
func (s *TimestreamWriteService) ResumeBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "TaskId")
	if taskId == "" {
		return nil, ErrInvalidParameter
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

	if task.TaskStatus != tsstore.BatchLoadStatusProgressStopped && task.TaskStatus != tsstore.BatchLoadStatusPendingResume {
		return nil, ErrInvalidParameter
	}

	if err := st.batchLoadStore.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusPendingResume, ""); err != nil {
		return nil, s.mapStoreError(err)
	}

	s.batchWg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		defer s.batchWg.Done()
		s.simulateBatchLoad(ctx, st.batchLoadStore, taskId)
	}()

	return response.EmptyResponse(), nil
}

// DeleteBatchLoadTask deletes a batch load task.
func (s *TimestreamWriteService) DeleteBatchLoadTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "TaskId")
	if taskId == "" {
		return nil, ErrInvalidParameter
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.batchLoadStore.DeleteBatchLoadTask(taskId); err != nil {
		if err == tsstore.ErrBatchLoadTaskNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

func (s *TimestreamWriteService) simulateBatchLoad(ctx context.Context, store *tsstore.BatchLoadTaskStore, taskId string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	defer func() {
		if ctx.Err() != nil {
			_ = store.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusFailed, "Simulation timed out")
		}
	}()

	if err := store.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusInProgress, ""); err != nil {
		logs.Error("Failed to update batch load task status", logs.Err(err))
		return
	}

	timer := time.NewTimer(2 * time.Second)
	select {
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
		return
	}

	progress := &tsstore.BatchLoadProgressReport{
		RecordsProcessed:        1000,
		RecordsIngested:         950,
		FileFailures:            0,
		ParseFailures:           10,
		BytesMetered:            102400,
		RecordIngestionFailures: 40,
	}
	if err := store.UpdateBatchLoadTaskProgress(taskId, progress); err != nil {
		logs.Error("Failed to update batch load task progress", logs.Err(err))
	}

	timer = time.NewTimer(1 * time.Second)
	select {
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
		return
	}

	progress.RecordsProcessed = 2000
	progress.RecordsIngested = 1900
	progress.BytesMetered = 204800
	if err := store.UpdateBatchLoadTaskProgress(taskId, progress); err != nil {
		logs.Error("Failed to update batch load task progress", logs.Err(err))
	}

	timer = time.NewTimer(1 * time.Second)
	select {
	case <-timer.C:
	case <-ctx.Done():
		timer.Stop()
		return
	}

	progress.RecordsProcessed = 3000
	progress.RecordsIngested = 2850
	progress.BytesMetered = 307200
	if err := store.UpdateBatchLoadTaskProgress(taskId, progress); err != nil {
		logs.Error("Failed to update batch load task progress", logs.Err(err))
	}

	if err := store.UpdateBatchLoadTaskStatus(taskId, tsstore.BatchLoadStatusSucceeded, ""); err != nil {
		logs.Error("Failed to update batch load task status", logs.Err(err))
	}
}

func (s *TimestreamWriteService) formatBatchLoadTask(task *tsstore.BatchLoadTask) map[string]interface{} {
	response := map[string]interface{}{
		"TaskId":             task.TaskId,
		"TaskStatus":         string(task.TaskStatus),
		"TargetDatabaseName": task.DatabaseName,
		"TargetTableName":    task.TableName,
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
	return time.Now().UTC().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		num, _ := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[num.Int64()]
	}
	return string(b)
}

func parseint64(s string) (int64, bool) {
	result, err := strconv.ParseInt(s, 10, 64)
	return result, err == nil
}
