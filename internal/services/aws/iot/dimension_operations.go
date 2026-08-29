package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
)

// ---- Dimensions ----------------------------------------------------

func (s *IoTService) CreateDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tagList := tags.ParseTagsWithQueryFallback(req.Parameters, "tags")
	recTags := make(map[string]string, len(tagList))
	for _, t := range tagList {
		recTags[t.Key] = t.Value
	}
	result, err := s.createDimensionCore(store, CreateDimensionInput{
		Name:               request.GetParamCaseInsensitive(req.Parameters, "name"),
		Type:               request.GetParamCaseInsensitive(req.Parameters, "type"),
		StringValues:       request.GetStringList(req.Parameters, "stringValues"),
		Tags:               recTags,
		ClientRequestToken: request.GetParamCaseInsensitive(req.Parameters, "clientRequestToken"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name": result.Name,
		"arn":  result.Arn,
	}, nil
}
func (s *IoTService) DeleteDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteDimensionCore(store, request.GetParamCaseInsensitive(req.Parameters, "name")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) DescribeDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec, err := s.describeDimensionCore(store, request.GetParamCaseInsensitive(req.Parameters, "name"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name":             rec.Rec["name"],
		"arn":              rec.Arn,
		"type":             rec.Rec["type"],
		"stringValues":     rec.Rec["stringValues"],
		"creationDate":     rec.Rec["creationDate"],
		"lastModifiedDate": rec.Rec["lastModifiedDate"],
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
	return paginatedStrings("dimensionNames", names, req.Parameters)
}
func (s *IoTService) UpdateDimension(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.updateDimensionCore(store, UpdateDimensionInput{
		Name:         request.GetParamCaseInsensitive(req.Parameters, "name"),
		StringValues: request.GetStringList(req.Parameters, "stringValues"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name":             result.Rec["name"],
		"arn":              result.Arn,
		"type":             result.Rec["type"],
		"stringValues":     result.Rec["stringValues"],
		"creationDate":     result.Rec["creationDate"],
		"lastModifiedDate": result.Rec["lastModifiedDate"],
	}, nil
}
