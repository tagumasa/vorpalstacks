package timestreamwrite

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	tsstore "vorpalstacks/internal/store/aws/timestream"
	"vorpalstacks/internal/utils/aws/types"
)

// WriteRecords writes time-series records to a Timestream table.
func (s *TimestreamWriteService) WriteRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	databaseName := request.GetParamCaseInsensitive(req.Parameters, "DatabaseName")
	if databaseName == "" {
		return nil, ErrValidationException
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrValidationException
	}

	records := s.parseRecords(req.Parameters["Records"])

	if len(records) < 1 || len(records) > 100 {
		return nil, ErrValidationException
	}

	commonAttributes := s.parseCommonAttributes(req.Parameters["CommonAttributes"])
	if commonAttributes != nil {
		merged, err := s.mergeCommonAttributes(records, commonAttributes)
		if err != nil {
			return nil, err
		}
		records = merged
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rejectedRecords, err := st.recordStore.WriteRecords(databaseName, tableName, records)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	ingestedCount := int64(len(records) - len(rejectedRecords))

	response := map[string]interface{}{
		"RecordsIngested": map[string]interface{}{
			"Total":         ingestedCount,
			"MemoryStore":   ingestedCount,
			"MagneticStore": int64(0),
		},
	}

	if len(rejectedRecords) > 0 {
		response["RejectedRecords"] = s.formatRejectedRecords(rejectedRecords)
	}

	return response, nil
}

func (s *TimestreamWriteService) parseRecords(data interface{}) []tsstore.Record {
	var records []tsstore.Record

	recordsList, ok := data.([]interface{})
	if !ok {
		return records
	}

	for _, item := range recordsList {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		records = append(records, s.parseSingleRecord(itemMap))
	}

	return records
}

// parseSingleRecord parses a single record from a map[string]interface{}.
func (s *TimestreamWriteService) parseSingleRecord(itemMap map[string]interface{}) tsstore.Record {
	record := tsstore.Record{
		Dimensions:       s.parseDimensions(itemMap["Dimensions"]),
		MeasureName:      request.GetStringParam(itemMap, "MeasureName"),
		MeasureValue:     request.GetStringParam(itemMap, "MeasureValue"),
		MeasureValueType: tsstore.MeasureValueType(request.GetStringParam(itemMap, "MeasureValueType")),
		Time:             request.GetStringParam(itemMap, "Time"),
		TimeUnit:         tsstore.TimeUnit(request.GetStringParam(itemMap, "TimeUnit")),
		Version:          getIntFromMap(itemMap, "Version"),
	}

	record.MeasureValues = s.parseMeasureValues(itemMap["MeasureValues"])

	return record
}

func (s *TimestreamWriteService) parseDimensions(data interface{}) []tsstore.Dimension {
	var dimensions []tsstore.Dimension

	dimsList, ok := data.([]interface{})
	if !ok {
		return dimensions
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
		if dim.DimensionValueType == "" {
			dim.DimensionValueType = tsstore.DimensionValueTypeVarchar
		}

		dimensions = append(dimensions, dim)
	}

	return dimensions
}

// parseCommonAttributes parses the CommonAttributes record from the request.
// CommonAttributes contains measure, dimension, time, and version attributes
// shared across all records in the WriteRecords request.
func (s *TimestreamWriteService) parseCommonAttributes(data interface{}) *tsstore.Record {
	if data == nil {
		return nil
	}
	itemMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	rec := s.parseSingleRecord(itemMap)
	return &rec
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

func (s *TimestreamWriteService) formatRejectedRecords(records []tsstore.RejectedRecord) []map[string]interface{} {
	var result []map[string]interface{}
	for _, r := range records {
		result = append(result, map[string]interface{}{
			"RecordIndex":     r.RecordIndex,
			"Reason":          r.Reason,
			"ExistingVersion": r.ExistingVersion,
		})
	}
	return result
}

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
		TagFunc: func(_ context.Context, resourceKey string, tagSlice []types.Tag) error {
			return st.store.TagFromSlice(resourceKey, tagSlice)
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return st.store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			return st.store.ListAsSlice(resourceKey)
		},
		FormatResponse: func(tagSlice []types.Tag, _ string) (interface{}, error) {
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

// TagResource adds tags to a Timestream resource.
func (s *TimestreamWriteService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.tagHandlerConfig(st))
}

// UntagResource removes tags from a Timestream resource.
func (s *TimestreamWriteService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.tagHandlerConfig(st))
}

// ListTagsForResource returns the tags for a Timestream resource.
func (s *TimestreamWriteService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.tagHandlerConfig(st))
}

func getIntFromMap(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case int64:
			return n
		case float64:
			return int64(n)
		}
	}
	return 0
}
