package iot

import (
	"context"
	"encoding/json"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

func (s *IoTService) CreateInput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "inputName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetInput(name); err == nil {
		return nil, iotstore.ErrInputAlreadyExists
	}

	defRaw := request.GetParamCaseInsensitive(req.Parameters, "inputDefinition")
	if defRaw == "" {
		defRaw = "{}"
	}
	var def map[string]interface{}
	json.Unmarshal([]byte(defRaw), &def)

	inp := &iotstore.Input{
		InputName:        name,
		InputDescription: request.GetParamCaseInsensitive(req.Parameters, "inputDescription"),
		InputDefinition:  def,
		Status:           "ACTIVE",
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateInput(inp)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"inputArn":  created.InputARN,
		"inputName": created.InputName,
	}, nil
}

func (s *IoTService) DescribeInput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "inputName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	inp, err := store.GetInput(name)
	if err != nil {
		return nil, iotstore.ErrInputNotFound
	}

	return inputDetailResponse(inp), nil
}

func (s *IoTService) UpdateInput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "inputName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	inp, err := store.GetInput(name)
	if err != nil {
		return nil, iotstore.ErrInputNotFound
	}

	if desc := request.GetParamCaseInsensitive(req.Parameters, "inputDescription"); desc != "" {
		inp.InputDescription = desc
	}
	if defRaw := request.GetParamCaseInsensitive(req.Parameters, "inputDefinition"); defRaw != "" {
		var def map[string]interface{}
		if json.Unmarshal([]byte(defRaw), &def) == nil {
			inp.InputDefinition = def
		}
	}
	if status := request.GetParamCaseInsensitive(req.Parameters, "status"); status != "" {
		inp.Status = status
	}
	inp.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateInput(inp); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) DeleteInput(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "inputName")
	if name == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteInput(name); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *IoTService) ListInputs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	inputs, err := store.ListInputs(parseListOptions(req.Parameters))
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(inputs.Items))
	for _, inp := range inputs.Items {
		items = append(items, inputResponse(inp))
	}

	return listResponse("inputs", items, inputs.NextMarker), nil
}
