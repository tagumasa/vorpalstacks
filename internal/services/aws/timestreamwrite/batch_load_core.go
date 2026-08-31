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

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	storecommon "vorpalstacks/internal/store/aws/common"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Output DTOs — batch load
// ---------------------------------------------------------------------------

// CreateBatchLoadTaskInput carries every field that CreateBatchLoadTask needs.
// The three configuration members hold the raw request values; their structure
// parsing and validation run inside the Core so both protocol planes share
// one validation path. Region is the request region used to drive the
// asynchronous batch execution.
type CreateBatchLoadTaskInput struct {
	TargetDatabaseName      string
	TargetTableName         string
	DataSourceConfiguration interface{}
	DataModelConfiguration  interface{}
	ReportConfiguration     interface{}
	RecordVersion           int64
	ClientToken             string
	Region                  string
}

// CreateBatchLoadTaskResult carries the created (or idempotently matched)
// task identifier.
type CreateBatchLoadTaskResult struct {
	TaskId string
}

// DescribeBatchLoadTaskInput carries the fields for DescribeBatchLoadTask.
type DescribeBatchLoadTaskInput struct {
	TaskId string
}

// ListBatchLoadTasksInput carries the fields for ListBatchLoadTasks.
type ListBatchLoadTasksInput struct {
	TaskStatus tsstore.BatchLoadStatus
	NextToken  string
	MaxResults int
}

// ListBatchLoadTasksResult is the paginated result of listing batch load tasks.
type ListBatchLoadTasksResult struct {
	Tasks     []*tsstore.BatchLoadTask
	NextToken string
}

// ResumeBatchLoadTaskInput carries the fields for ResumeBatchLoadTask.
// Region drives the asynchronous batch execution, as in CreateBatchLoadTask.
type ResumeBatchLoadTaskInput struct {
	TaskId string
	Region string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createBatchLoadTaskCore is the single entry point for batch load task
// creation shared by the HTTP API and any admin surface.
func (s *TimestreamWriteService) createBatchLoadTaskCore(st *tsWriteStores, in CreateBatchLoadTaskInput) (*CreateBatchLoadTaskResult, error) {
	targetDatabaseName := in.TargetDatabaseName
	if targetDatabaseName == "" {
		return nil, ErrValidationException
	}

	targetTableName := in.TargetTableName
	if targetTableName == "" {
		return nil, ErrValidationException
	}

	dataSourceConfig, err := parseDataSourceConfiguration(in.DataSourceConfiguration)
	if err != nil {
		return nil, err
	}
	if dataSourceConfig == nil {
		return nil, ErrValidationException
	}
	// Validate DataFormat enum (Smithy BatchLoadDataFormat: CSV only).
	if dataSourceConfig.DataFormat != "" && dataSourceConfig.DataFormat != tsstore.BatchLoadDataFormatCsv {
		return nil, ErrValidationException
	}

	reportConfig, err := parseReportConfiguration(in.ReportConfiguration)
	if err != nil {
		return nil, err
	}
	if reportConfig == nil {
		return nil, ErrValidationException
	}

	dataModelConfig, err := parseDataModelConfiguration(in.DataModelConfiguration)
	if err != nil {
		return nil, err
	}
	// DimensionMappings is REQUIRED when DataModel is provided (Smithy).
	if dataModelConfig != nil && dataModelConfig.DataModel != nil && len(dataModelConfig.DataModel.DimensionMappings) == 0 {
		return nil, ErrValidationException
	}

	clientToken := in.ClientToken
	// Validate ClientToken length (Smithy ClientRequestToken: 1-64).
	if clientToken != "" && !validateClientToken(clientToken) {
		return nil, ErrValidationException
	}

	// @idempotencyToken — if ClientToken matches an existing task,
	// return that task instead of creating a new one. The AWS SDK
	// auto-generates a UUID for the ClientToken (Smithy idempotencyToken
	// trait), but the TaskId is always server-generated (uppercase
	// alphanumeric, matching BatchLoadTaskId pattern ^[A-Z0-9]+$).
	// Properly propagate storage errors. Only ErrBatchLoadTaskNotFound
	// is an acceptable fallback (no match found); any other error must be
	// returned rather than silently creating a duplicate task.
	if clientToken != "" {
		existingTask, gerr := st.batchLoadStore.FindByClientToken(clientToken)
		if gerr == nil {
			return &CreateBatchLoadTaskResult{
				TaskId: existingTask.TaskId,
			}, nil
		}
		if gerr != tsstore.ErrBatchLoadTaskNotFound {
			return nil, ErrInternalServer
		}
	}

	taskId := generateTaskId()
	task, err := st.batchLoadStore.CreateBatchLoadTask(taskId, clientToken, targetDatabaseName, targetTableName, dataSourceConfig, dataModelConfig, reportConfig, in.RecordVersion)
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
		s.executeBatchLoad(ctx, st, taskId, in.Region)
	}()

	return &CreateBatchLoadTaskResult{
		TaskId: task.TaskId,
	}, nil
}

// describeBatchLoadTaskCore is the single entry point for batch load task
// description.
func (s *TimestreamWriteService) describeBatchLoadTaskCore(st *tsWriteStores, in DescribeBatchLoadTaskInput) (*tsstore.BatchLoadTaskDescription, error) {
	taskId := in.TaskId
	if taskId == "" {
		return nil, ErrValidationException
	}
	// Validate BatchLoadTaskId pattern (Smithy: ^[A-Z0-9]+$, 3-32 chars).
	if !batchLoadTaskIdRegex.MatchString(taskId) {
		return nil, ErrValidationException
	}

	task, err := st.batchLoadStore.GetBatchLoadTask(taskId)
	if err != nil {
		if err == tsstore.ErrBatchLoadTaskNotFound {
			return nil, ErrResourceNotFound
		}
		return nil, s.mapStoreError(err)
	}

	return task, nil
}

// listBatchLoadTasksCore is the single entry point for listing batch load
// tasks.
func (s *TimestreamWriteService) listBatchLoadTasksCore(st *tsWriteStores, in ListBatchLoadTasksInput) (*ListBatchLoadTasksResult, error) {
	maxResults := in.MaxResults
	if maxResults > maxListBatchLoadTasksResults {
		maxResults = maxListBatchLoadTasksResults
	}

	opts := storecommon.ListOptions{MaxItems: maxResults}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}

	result, err := st.batchLoadStore.ListBatchLoadTasks(in.TaskStatus, opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return &ListBatchLoadTasksResult{
		Tasks:     result.Items,
		NextToken: result.NextMarker,
	}, nil
}

// resumeBatchLoadTaskCore is the single entry point for resuming a stopped
// batch load task.
func (s *TimestreamWriteService) resumeBatchLoadTaskCore(st *tsWriteStores, in ResumeBatchLoadTaskInput) error {
	taskId := in.TaskId
	if taskId == "" {
		return ErrValidationException
	}
	// Validate BatchLoadTaskId pattern (Smithy: ^[A-Z0-9]+$, 3-32 chars).
	if !batchLoadTaskIdRegex.MatchString(taskId) {
		return ErrValidationException
	}

	err := st.batchLoadStore.TransitionBatchLoadTaskStatus(taskId, []tsstore.BatchLoadStatus{
		tsstore.BatchLoadStatusProgressStopped,
		tsstore.BatchLoadStatusPendingResume,
		tsstore.BatchLoadStatusFailed,
	}, tsstore.BatchLoadStatusPendingResume)
	if err != nil {
		if err == tsstore.ErrBatchLoadTaskNotFound {
			return ErrResourceNotFound
		}
		if err == tsstore.ErrInvalidTransition {
			return ErrValidationException
		}
		return s.mapStoreError(err)
	}

	s.batchWg.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(s.batchCtx, 5*time.Minute)
		defer cancel()
		defer s.batchWg.Done()
		defer func() { resilience.RecoverPanic("timestreamwrite batch load resume") }()
		s.executeBatchLoad(ctx, st, taskId, in.Region)
	}()

	return nil
}

// executeBatchLoad drives the asynchronous S3 → CSV → records ingestion for
// one batch load task. It runs on the service's batch context and reports
// all progress through the batch load store.
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

		data, err := s.s3Invoker.GetObject(ctx, region, s3Config.BucketName, key, maxBatchLoadObjectSize)
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

		records, rejectedFromParse := parseCSVToRecords(data, csvConfig, dataModel)
		progress.ParseFailures += int64(len(rejectedFromParse))
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

func parseDataSourceConfiguration(config interface{}) (*tsstore.DataSourceConfiguration, error) {
	if config == nil {
		return nil, nil
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, ErrValidationException
	}

	result := &tsstore.DataSourceConfiguration{
		DataFormat: tsstore.BatchLoadDataFormat("CSV"),
	}

	if format, ok := configMap["DataFormat"].(string); ok {
		result.DataFormat = tsstore.BatchLoadDataFormat(format)
	}

	// DataSourceS3Configuration is REQUIRED per Smithy.
	s3Config, hasS3 := configMap["DataSourceS3Configuration"].(map[string]interface{})
	if !hasS3 {
		return nil, ErrValidationException
	}
	// BucketName is REQUIRED per Smithy when DataSourceS3Configuration is provided.
	bucket, hasBucket := s3Config["BucketName"].(string)
	if !hasBucket || bucket == "" {
		return nil, ErrValidationException
	}
	if !validateS3BucketName(bucket) {
		return nil, ErrValidationException
	}
	result.DataSourceS3Configuration = &tsstore.DataSourceS3Configuration{
		BucketName: bucket,
	}
	if prefix, ok := s3Config["ObjectKeyPrefix"].(string); ok {
		if !validateS3ObjectKeyPrefix(prefix) {
			return nil, ErrValidationException
		}
		result.DataSourceS3Configuration.ObjectKeyPrefix = prefix
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

	return result, nil
}

func parseReportConfiguration(config interface{}) (*tsstore.ReportConfiguration, error) {
	if config == nil {
		return nil, nil
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, ErrValidationException
	}

	result := &tsstore.ReportConfiguration{}

	if s3Config, ok := configMap["ReportS3Configuration"].(map[string]interface{}); ok {
		// ReportS3Configuration is optional per Smithy, but when provided,
		// BucketName is REQUIRED.
		bucket, hasBucket := s3Config["BucketName"].(string)
		if !hasBucket || bucket == "" {
			return nil, ErrValidationException
		}
		if !validateS3BucketName(bucket) {
			return nil, ErrValidationException
		}
		result.ReportS3Configuration = &tsstore.ReportS3Configuration{
			BucketName: bucket,
		}
		if enc, ok := s3Config["EncryptionOption"].(string); ok {
			if !validateEncryptionOption(enc) {
				return nil, ErrValidationException
			}
			result.ReportS3Configuration.EncryptionOption = tsstore.S3EncryptionOption(enc)
		}
		if kms, ok := s3Config["KmsKeyId"].(string); ok {
			if kms != "" && !validateKmsKeyId(kms) {
				return nil, ErrValidationException
			}
			result.ReportS3Configuration.KmsKeyId = kms
		}
		if prefix, ok := s3Config["ObjectKeyPrefix"].(string); ok {
			if !validateS3ObjectKeyPrefix(prefix) {
				return nil, ErrValidationException
			}
			result.ReportS3Configuration.ObjectKeyPrefix = prefix
		}
	}

	return result, nil
}

func parseDataModelConfiguration(config interface{}) (*tsstore.DataModelConfiguration, error) {
	if config == nil {
		return nil, nil
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, ErrValidationException
	}

	result := &tsstore.DataModelConfiguration{}

	if dataModel, ok := configMap["DataModel"].(map[string]interface{}); ok {
		dm, err := parseDataModel(dataModel)
		if err != nil {
			return nil, err
		}
		result.DataModel = dm
	}

	if s3Config, ok := configMap["DataModelS3Configuration"].(map[string]interface{}); ok {
		result.DataModelS3Configuration = &tsstore.DataModelS3Configuration{}
		bucket, hasBucket := s3Config["BucketName"].(string)
		if !hasBucket || bucket == "" {
			return nil, ErrValidationException
		}
		if !validateS3BucketName(bucket) {
			return nil, ErrValidationException
		}
		result.DataModelS3Configuration.BucketName = bucket
		if key, ok := s3Config["ObjectKey"].(string); ok && key != "" {
			if !validateS3ObjectKey(key) {
				return nil, ErrValidationException
			}
			result.DataModelS3Configuration.ObjectKey = key
		}
	}

	return result, nil
}

func parseDataModel(dataModel map[string]interface{}) (*tsstore.DataModel, error) {
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
		// MultiMeasureAttributeMappings is REQUIRED per Smithy when
		// MultiMeasureMappings is provided. Reject empty mappings rather
		// than silently discarding the entire DataModel.
		if len(mmm.MultiMeasureAttributeMappings) == 0 {
			return nil, ErrValidationException
		}
		result.MultiMeasureMappings = mmm
	}

	return result, nil
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
// configuration. Returns the parsed records along with per-record rejection
// details so callers can report which rows failed and why.
func parseCSVToRecords(data []byte, csvConfig *tsstore.CsvConfiguration, dataModelConfig *tsstore.DataModelConfiguration) ([]tsstore.Record, []tsstore.RejectedRecord) {
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
		return nil, nil
	}
	if len(rows) == 0 {
		return nil, nil
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
	var rejectedRecords []tsstore.RejectedRecord

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
			rejectedRecords = append(rejectedRecords, tsstore.RejectedRecord{
				RecordIndex: int64(rowIdx),
				Reason:      "record has no Time and no Dimensions",
			})
			continue
		}

		records = append(records, record)
	}

	return records, rejectedRecords
}

func parseint64(s string) (int64, bool) {
	result, err := strconv.ParseInt(s, 10, 64)
	return result, err == nil
}
