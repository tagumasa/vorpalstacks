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
		return nil, ErrInvalidParameter
	}

	tableName := request.GetParamCaseInsensitive(req.Parameters, "TableName")
	if tableName == "" {
		return nil, ErrInvalidParameter
	}

	records := s.parseRecords(req.Parameters["Records"])

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rejectedRecords, err := st.recordStore.WriteRecords(databaseName, tableName, records)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	response := map[string]interface{}{
		"RecordsIngested": map[string]interface{}{
			"Total":         int64(len(records)),
			"MemoryStore":   int64(len(records) - len(rejectedRecords)),
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

		records = append(records, record)
	}

	return records
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
