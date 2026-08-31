package timestreamwrite

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	tsstore "vorpalstacks/internal/store/aws/timestream"
)

// WriteRecords writes time-series records to a Timestream table.
func (s *TimestreamWriteService) WriteRecords(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.writeRecordsCore(st, WriteRecordsInput{
		DatabaseName:     request.GetParamCaseInsensitive(req.Parameters, "DatabaseName"),
		TableName:        request.GetParamCaseInsensitive(req.Parameters, "TableName"),
		Records:          req.Parameters["Records"],
		CommonAttributes: req.Parameters["CommonAttributes"],
	})
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"RecordsIngested": map[string]interface{}{
			"Total":         result.IngestedCount,
			"MemoryStore":   result.IngestedCount,
			"MagneticStore": result.MagneticStoreCount,
		},
	}

	if len(result.RejectedRecords) > 0 {
		resp["RejectedRecords"] = s.formatRejectedRecords(result.RejectedRecords)
	}

	return resp, nil
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
