package timestreamwrite

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// WriteRecords writes time-series records to a Timestream table.
func (s *Service) WriteRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

func (s *Service) parseRecords(data interface{}) []tsstore.Record {
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
			MeasureName:      getStringFromMap(itemMap, "MeasureName"),
			MeasureValue:     getStringFromMap(itemMap, "MeasureValue"),
			MeasureValueType: tsstore.MeasureValueType(getStringFromMap(itemMap, "MeasureValueType")),
			Time:             getStringFromMap(itemMap, "Time"),
			TimeUnit:         tsstore.TimeUnit(getStringFromMap(itemMap, "TimeUnit")),
			Version:          getIntFromMap(itemMap, "Version"),
		}

		record.MeasureValues = s.parseMeasureValues(itemMap["MeasureValues"])

		records = append(records, record)
	}

	return records
}

func (s *Service) parseDimensions(data interface{}) []tsstore.Dimension {
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
			Name:               getStringFromMap(itemMap, "Name"),
			Value:              getStringFromMap(itemMap, "Value"),
			DimensionValueType: tsstore.DimensionValueType(getStringFromMap(itemMap, "DimensionValueType")),
		}
		if dim.DimensionValueType == "" {
			dim.DimensionValueType = tsstore.DimensionValueTypeVarchar
		}

		dimensions = append(dimensions, dim)
	}

	return dimensions
}

func (s *Service) parseMeasureValues(data interface{}) []tsstore.MeasureValue {
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
			Name:  getStringFromMap(itemMap, "Name"),
			Value: getStringFromMap(itemMap, "Value"),
			Type:  tsstore.MeasureValueType(getStringFromMap(itemMap, "Type")),
		}

		measureValues = append(measureValues, mv)
	}

	return measureValues
}

func (s *Service) formatRejectedRecords(records []tsstore.RejectedRecord) []map[string]interface{} {
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

// TagResource adds tags to a Timestream resource.
func (s *Service) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamCaseInsensitive(req.Parameters, "ResourceARN")
	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	tagMap := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))
	if len(tagMap) == 0 {
		return nil, ErrInvalidParameter
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.store.TagResource(resourceARN, tagMap); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a Timestream resource.
func (s *Service) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamCaseInsensitive(req.Parameters, "ResourceARN")
	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	tagKeys := tagutil.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys")
	if len(tagKeys) == 0 {
		return nil, ErrInvalidParameter
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.store.UntagResource(resourceARN, tagKeys); err != nil {
		return nil, s.mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource returns the tags for a Timestream resource.
func (s *Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetParamCaseInsensitive(req.Parameters, "ResourceARN")
	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tagMap, err := st.store.ListTags(resourceARN)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	var tagList []map[string]string
	for k, v := range tagMap {
		tagList = append(tagList, map[string]string{
			"Key":   k,
			"Value": v,
		})
	}

	return map[string]interface{}{
		"Tags": tagList,
	}, nil
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
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
