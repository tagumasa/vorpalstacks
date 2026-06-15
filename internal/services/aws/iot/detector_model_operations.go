package iot

import (
	"context"
	"encoding/json"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateDetectorModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetDetectorModel(name); err == nil {
		return nil, iotstore.ErrDetectorModelAlreadyExists
	}

	defRaw := request.GetParamCaseInsensitive(req.Parameters, "detectorModelDefinition")
	if defRaw == "" {
		defRaw = "{}"
	}

	var def map[string]interface{}
	json.Unmarshal([]byte(defRaw), &def)

	dm := &iotstore.DetectorModel{
		DetectorModelName:        name,
		DetectorModelDescription: request.GetParamCaseInsensitive(req.Parameters, "detectorModelDescription"),
		RoleARN:                  request.GetParamCaseInsensitive(req.Parameters, "roleArn"),
		DetectorModelDefinition:  def,
		Status:                   "ACTIVE",
		CreationDate:             time.Now().UTC(),
		LastModifiedDate:         time.Now().UTC(),
	}

	created, err := store.CreateDetectorModel(dm)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"detectorModelArn":  created.DetectorModelARN,
		"detectorModelName": created.DetectorModelName,
	}, nil
}

func (s *IoTService) DescribeDetectorModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		return nil, iotstore.ErrDetectorModelNotFound
	}

	return detectorModelDetailResponse(dm), nil
}

func (s *IoTService) UpdateDetectorModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		return nil, iotstore.ErrDetectorModelNotFound
	}

	if desc := request.GetParamCaseInsensitive(req.Parameters, "detectorModelDescription"); desc != "" {
		dm.DetectorModelDescription = desc
	}
	if role := request.GetParamCaseInsensitive(req.Parameters, "roleArn"); role != "" {
		dm.RoleARN = role
	}
	if defRaw := request.GetParamCaseInsensitive(req.Parameters, "detectorModelDefinition"); defRaw != "" {
		var def map[string]interface{}
		if json.Unmarshal([]byte(defRaw), &def) == nil {
			dm.DetectorModelDefinition = def
		}
	}
	if status := request.GetParamCaseInsensitive(req.Parameters, "status"); status != "" {
		dm.Status = status
	}
	dm.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateDetectorModel(dm); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteDetectorModel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "detectorModelName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteDetectorModel(name); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListDetectorModels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	models, err := store.ListDetectorModels(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(models.Items))
	for _, dm := range models.Items {
		items = append(items, detectorModelResponse(dm))
	}

	return listResponse("detectorModels", items, models.NextMarker), nil
}
