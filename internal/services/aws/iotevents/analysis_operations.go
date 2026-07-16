package iotevents

import (
	"context"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTEventsService) StartDetectorModelAnalysis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	analysisID := uuid.New().String()
	return map[string]interface{}{
		"analysisId": analysisID,
	}, nil
}

func (s *IoTEventsService) DescribeDetectorModelAnalysis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	analysisID := request.GetParamCaseInsensitive(req.Parameters, "analysisId")
	if analysisID == "" {
		return nil, iotstore.ErrMissingParam
	}
	return map[string]interface{}{
		"analysisId": analysisID,
		"status":     "COMPLETE",
	}, nil
}

func (s *IoTEventsService) GetDetectorModelAnalysisResults(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	analysisID := request.GetParamCaseInsensitive(req.Parameters, "analysisId")
	if analysisID == "" {
		return nil, iotstore.ErrMissingParam
	}
	return map[string]interface{}{
		"analysisId": analysisID,
		"results":    []map[string]interface{}{},
	}, nil
}

func (s *IoTEventsService) PutLoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	opts := map[string]interface{}{}
	_ = store.GetGeneric("config/ioteventsLogging", &opts)
	if raw, ok := req.Parameters["detectorDebugLogging"]; ok {
		opts["detectorDebugLogging"] = raw
	}
	if raw, ok := req.Parameters["roleArn"]; ok {
		opts["roleArn"] = raw
	}
	if raw, ok := req.Parameters["level"]; ok {
		opts["level"] = raw
	}
	if err := store.PutGeneric("config/ioteventsLogging", opts); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTEventsService) DescribeLoggingOptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	opts := map[string]interface{}{}
	_ = store.GetGeneric("config/ioteventsLogging", &opts)
	if len(opts) == 0 {
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{
		"loggingOptions": opts,
	}, nil
}

func (s *IoTEventsService) ListInputRoutings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"inputRoutings": []map[string]interface{}{},
		"nextToken":     "",
	}, nil
}

func (s *IoTEventsService) ListDetectorModelVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	dm, err := store.GetDetectorModel(name)
	if err != nil {
		return nil, err
	}
	version := dm.DetectorModelVersion
	if version == "" {
		version = "1"
	}
	versions := []map[string]interface{}{
		{
			"detectorModelName":    dm.DetectorModelName,
			"detectorModelVersion": version,
			"creationTime":         dm.CreationDate.Unix(),
		},
	}
	return paginatedMaps("detectorModelVersionSummaries", versions, req.Parameters), nil
}
