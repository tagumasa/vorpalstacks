package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---- Dimensions ----------------------------------------------------

func (s *IoTService) CreateDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, err := s.bulkCreate(reqCtx, "dimension", req, "name", map[string]interface{}{
		"type":         request.GetParamCaseInsensitive(req.Parameters, "type"),
		"stringValues": request.GetStringList(req.Parameters, "stringValues"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name": rec["name"],
		"arn":  iotstore.BuildDimensionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
	}, nil
}
func (s *IoTService) DeleteDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.bulkDelete(reqCtx, "dimension", req, "name"); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, _, exists, err := s.bulkGet(reqCtx, "dimension", req, "name")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDimensionNotFound
	}
	return map[string]interface{}{
		"name":             rec["name"],
		"arn":              iotstore.BuildDimensionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"type":             rec["type"],
		"stringValues":     rec["stringValues"],
		"creationDate":     rec["creationDate"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}
func (s *IoTService) ListDimensions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy: ListDimensionsResponse.dimensionNames is list<DimensionName> (string).
	items, err := s.bulkList(reqCtx, "dimension")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(items))
	for _, item := range items {
		if name, ok := item["name"].(string); ok {
			names = append(names, name)
		}
	}
	return paginatedStrings("dimensionNames", names, req.Parameters), nil
}
func (s *IoTService) UpdateDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	rec, exists, err := s.bulkUpdate(reqCtx, "dimension", req, "name", map[string]interface{}{
		"stringValues": request.GetStringList(req.Parameters, "stringValues"),
	})
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrDimensionNotFound
	}
	return map[string]interface{}{
		"name":             rec["name"],
		"arn":              iotstore.BuildDimensionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), bulkName(rec)),
		"stringValues":     rec["stringValues"],
		"lastModifiedDate": rec["lastModifiedDate"],
	}, nil
}
