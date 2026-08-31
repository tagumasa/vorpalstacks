package timestreamwrite

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Output DTOs — record ingestion
// ---------------------------------------------------------------------------

// WriteRecordsInput carries every field that WriteRecords needs. Records and
// CommonAttributes hold the raw request values: per-record parsing, merging
// and enum validation run inside the Core so both protocol planes share one
// validation path.
type WriteRecordsInput struct {
	DatabaseName     string
	TableName        string
	Records          interface{}
	CommonAttributes interface{}
}

// WriteRecordsResult carries the ingestion counts and any per-record
// rejections.
type WriteRecordsResult struct {
	IngestedCount      int64
	MagneticStoreCount int64
	RejectedRecords    []tsstore.RejectedRecord
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// writeRecordsCore is the single entry point for record ingestion shared by
// the HTTP API and any admin surface.
func (s *TimestreamWriteService) writeRecordsCore(st *tsWriteStores, in WriteRecordsInput) (*WriteRecordsResult, error) {
	databaseName := in.DatabaseName
	if databaseName == "" {
		return nil, ErrValidationException
	}

	tableName := in.TableName
	if tableName == "" {
		return nil, ErrValidationException
	}

	records, err := s.parseRecords(in.Records)
	if err != nil {
		return nil, err
	}

	if len(records) < 1 || len(records) > 100 {
		return nil, ErrValidationException
	}

	commonAttributes, err := s.parseCommonAttributes(in.CommonAttributes)
	if err != nil {
		return nil, err
	}
	if commonAttributes != nil {
		merged, err := s.mergeCommonAttributes(records, commonAttributes)
		if err != nil {
			return nil, err
		}
		records = merged
	}

	// Validate the final (post-merge) enum values of every record: an
	// enum member outside the Smithy enum is a malformed request and
	// must fail the whole call, not be silently persisted.
	if err := validateRecordEnums(records); err != nil {
		return nil, err
	}

	rejectedRecords, err := st.recordStore.WriteRecords(databaseName, tableName, records)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// When all records are rejected, raise RejectedRecordsException per
	// Smithy (HTTP 419). The SDK expects errors.As(&types.RejectedRecordsException{}).
	if len(rejectedRecords) == len(records) && len(records) > 0 {
		exc := *ErrRejectedRecordsException
		exc.WithRawField("RejectedRecords", s.formatRejectedRecords(rejectedRecords))
		return nil, &exc
	}

	ingestedCount := int64(len(records) - len(rejectedRecords))

	// When the table has EnableMagneticStoreWrites=true, ingested records
	// are also written to the magnetic store. Reflect this in the count.
	magneticStoreCount := int64(0)
	if table, terr := st.tableStore.GetTable(databaseName, tableName); terr == nil {
		if table.MagneticStoreWriteProperties != nil && table.MagneticStoreWriteProperties.EnableMagneticStoreWrites {
			magneticStoreCount = ingestedCount
		}
	}

	return &WriteRecordsResult{
		IngestedCount:      ingestedCount,
		MagneticStoreCount: magneticStoreCount,
		RejectedRecords:    rejectedRecords,
	}, nil
}

func (s *TimestreamWriteService) parseRecords(data interface{}) ([]tsstore.Record, error) {
	var records []tsstore.Record

	recordsList, ok := data.([]interface{})
	if !ok {
		return records, nil
	}

	for _, item := range recordsList {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		rec, err := s.parseSingleRecord(itemMap)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

// parseSingleRecord parses a single record from a map[string]interface{}.
func (s *TimestreamWriteService) parseSingleRecord(itemMap map[string]interface{}) (tsstore.Record, error) {
	dimensions, err := s.parseDimensions(itemMap["Dimensions"])
	if err != nil {
		return tsstore.Record{}, err
	}

	record := tsstore.Record{
		Dimensions:       dimensions,
		MeasureName:      request.GetStringParam(itemMap, "MeasureName"),
		MeasureValue:     request.GetStringParam(itemMap, "MeasureValue"),
		MeasureValueType: tsstore.MeasureValueType(request.GetStringParam(itemMap, "MeasureValueType")),
		Time:             request.GetStringParam(itemMap, "Time"),
		TimeUnit:         tsstore.TimeUnit(request.GetStringParam(itemMap, "TimeUnit")),
		Version:          getIntFromMap(itemMap, "Version"),
	}

	record.MeasureValues = s.parseMeasureValues(itemMap["MeasureValues"])

	return record, nil
}

// parseDimensions parses dimension list from request data.
// Dimension.Name and Dimension.Value are both REQUIRED per Smithy. Reject
// empty values rather than silently persisting corrupted data.
func (s *TimestreamWriteService) parseDimensions(data interface{}) ([]tsstore.Dimension, error) {
	var dimensions []tsstore.Dimension

	dimsList, ok := data.([]interface{})
	if !ok {
		return dimensions, nil
	}

	for _, item := range dimsList {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		dim := tsstore.Dimension{
			Name:               request.GetStringParam(itemMap, "Name"),
			Value:              request.GetStringParam(itemMap, "Value"),
			DimensionValueType: tsstore.DimensionValueType(request.GetStringParam(itemMap, "DimensionValueType")),
		}

		if dim.Name == "" {
			return nil, awserrors.NewValidationException("Dimension.Name is required")
		}
		if dim.Value == "" {
			return nil, awserrors.NewValidationException("Dimension.Value is required")
		}

		if dim.DimensionValueType == "" {
			dim.DimensionValueType = tsstore.DimensionValueTypeVarchar
		}

		dimensions = append(dimensions, dim)
	}

	return dimensions, nil
}

// parseCommonAttributes parses the CommonAttributes record from the request.
func (s *TimestreamWriteService) parseCommonAttributes(data interface{}) (*tsstore.Record, error) {
	if data == nil {
		return nil, nil
	}
	itemMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	rec, err := s.parseSingleRecord(itemMap)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// mergeCommonAttributes merges CommonAttributes into each individual record.
// Per AWS spec: dimension names must not overlap between common and individual
// records (ValidationException). For scalar fields (MeasureName, MeasureValue,
// MeasureValueType, Time, TimeUnit, Version), the individual record's value
// takes precedence when set; otherwise the common value is used.
func (s *TimestreamWriteService) mergeCommonAttributes(records []tsstore.Record, common *tsstore.Record) ([]tsstore.Record, error) {
	if common == nil {
		return records, nil
	}

	commonDimNames := make(map[string]bool, len(common.Dimensions))
	for _, d := range common.Dimensions {
		commonDimNames[d.Name] = true
	}

	for i := range records {
		rec := &records[i]

		for _, d := range rec.Dimensions {
			if commonDimNames[d.Name] {
				return nil, ErrValidationException
			}
		}

		rec.Dimensions = append(rec.Dimensions, common.Dimensions...)

		if rec.MeasureName == "" {
			rec.MeasureName = common.MeasureName
		}
		if rec.MeasureValue == "" {
			rec.MeasureValue = common.MeasureValue
		}
		if rec.MeasureValueType == "" {
			rec.MeasureValueType = common.MeasureValueType
		}
		if len(rec.MeasureValues) == 0 {
			rec.MeasureValues = common.MeasureValues
		}
		if rec.Time == "" {
			rec.Time = common.Time
		}
		if rec.TimeUnit == "" {
			rec.TimeUnit = common.TimeUnit
		}
		if rec.Version == 0 {
			rec.Version = common.Version
		}
	}

	return records, nil
}

// validateRecordEnums validates the final enum values of every record
// after CommonAttributes merging has produced the definitive field
// values. MeasureValueType, TimeUnit, DimensionValueType and
// MeasureValue.Type are Smithy enums; invalid values are rejected with
// ValidationException rather than silently persisted or reinterpreted
// by the store layer.
func validateRecordEnums(records []tsstore.Record) error {
	for i := range records {
		rec := &records[i]
		if rec.MeasureValueType != "" && !validateMeasureValueType(string(rec.MeasureValueType)) {
			return awserrors.NewValidationException(fmt.Sprintf("Records[%d].MeasureValueType has invalid value %q", i, rec.MeasureValueType))
		}
		// Per the Smithy documentation for MeasureValues, the list is
		// "only allowed for type MULTI. For scalar values, use the
		// MeasureValue attribute of the record directly" — enforce the
		// pairing in both directions.
		if rec.MeasureValueType == tsstore.MeasureValueTypeMulti && len(rec.MeasureValues) == 0 {
			return awserrors.NewValidationException(fmt.Sprintf("Records[%d].MeasureValues is required when MeasureValueType is MULTI", i))
		}
		if rec.MeasureValueType != tsstore.MeasureValueTypeMulti && len(rec.MeasureValues) > 0 {
			return awserrors.NewValidationException(fmt.Sprintf("Records[%d].MeasureValues is only allowed when MeasureValueType is MULTI", i))
		}
		if rec.TimeUnit != "" && !validateTimeUnit(string(rec.TimeUnit)) {
			return awserrors.NewValidationException(fmt.Sprintf("Records[%d].TimeUnit has invalid value %q", i, rec.TimeUnit))
		}
		for j, d := range rec.Dimensions {
			if d.DimensionValueType != "" && !validateDimensionValueType(string(d.DimensionValueType)) {
				return awserrors.NewValidationException(fmt.Sprintf("Records[%d].Dimensions[%d].DimensionValueType has invalid value %q", i, j, d.DimensionValueType))
			}
		}
		for j, mv := range rec.MeasureValues {
			if mv.Type == "" {
				return awserrors.NewValidationException(fmt.Sprintf("Records[%d].MeasureValues[%d].Type is required", i, j))
			}
			if !validateMeasureValueType(string(mv.Type)) {
				return awserrors.NewValidationException(fmt.Sprintf("Records[%d].MeasureValues[%d].Type has invalid value %q", i, j, mv.Type))
			}
		}
	}
	return nil
}

func (s *TimestreamWriteService) parseMeasureValues(data interface{}) []tsstore.MeasureValue {
	var measureValues []tsstore.MeasureValue

	mvList, ok := data.([]interface{})
	if !ok {
		return measureValues
	}

	for _, item := range mvList {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		mv := tsstore.MeasureValue{
			Name:  request.GetStringParam(itemMap, "Name"),
			Value: request.GetStringParam(itemMap, "Value"),
			Type:  tsstore.MeasureValueType(request.GetStringParam(itemMap, "Type")),
		}

		measureValues = append(measureValues, mv)
	}

	return measureValues
}

// tagHandlerConfig builds the shared tag handler configuration bound to the
// given store bundle; the tag store closures are the single tag persistence
// path for both protocol planes.
func (s *TimestreamWriteService) tagHandlerConfig(st *tsWriteStores) tagutil.TagHandlerConfig {
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
		TagFunc: func(_ context.Context, resourceKey string, tagSlice []tagutil.Tag) error {
			return st.store.TagFromSlice(resourceKey, tagSlice)
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return st.store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return st.store.ListAsSlice(resourceKey)
		},
		FormatResponse: func(tagSlice []tagutil.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.MapToResponse(tagutil.ToMap(tagSlice)),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: s.mapStoreError,
	}
}
