package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateThingType creates a new thing type with optional description and
// property definitions. Returns ResourceAlreadyExistsException if the thing
// type name is already taken.
func (s *IoTService) CreateThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	props := unwrapProps(req.Parameters, "thingTypeProperties")
	in := CreateThingTypeInput{
		ThingTypeName:        request.GetParamCaseInsensitive(req.Parameters, "thingTypeName"),
		Description:          request.GetParamCaseInsensitive(props, "thingTypeDescription"),
		SearchableAttributes: parseSearchableAttributes(props),
	}

	created, err := s.createThingTypeCore(store, in)
	if err != nil {
		return nil, err
	}

	return thingTypeResponse(created), nil
}

// DescribeThingType retrieves the details of an existing thing type including
// properties and deprecation status.
func (s *IoTService) DescribeThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tt, err := s.describeThingTypeCore(store, thingTypeName)
	if err != nil {
		return nil, err
	}

	return thingTypeDescribeResponse(tt), nil
}

// UpdateThingType modifies the description or properties of an existing thing
// type. Increments the version counter.
func (s *IoTService) UpdateThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	props := unwrapProps(req.Parameters, "thingTypeProperties")
	in := UpdateThingTypeInput{
		ThingTypeName:        request.GetParamCaseInsensitive(req.Parameters, "thingTypeName"),
		Description:          request.GetParamCaseInsensitive(props, "thingTypeDescription"),
		SearchableAttributes: parseSearchableAttributes(props),
	}

	updated, err := s.updateThingTypeCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"thingTypeName":    updated.ThingTypeName,
		"thingTypeArn":     updated.ThingTypeARN,
		"thingTypeId":      updated.ThingTypeID,
		"description":      updated.Description,
		"version":          updated.Version,
		"lastModifiedDate": updated.LastModifiedDate.Unix(),
	}, nil
}

// DeleteThingType removes a thing type from the registry. Returns
// ResourceNotFoundException if the thing type does not exist.
func (s *IoTService) DeleteThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteThingTypeCore(store, thingTypeName); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListThingTypes returns a paginated list of thing types.
func (s *IoTService) ListThingTypes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := s.listThingTypesCore(store, opts)
	if err != nil {
		return nil, err
	}

	types := make([]map[string]interface{}, 0, len(result.Items))
	for _, t := range result.Items {
		types = append(types, thingTypeResponse(t))
	}

	return listResponse("thingTypes", types, result.NextMarker), nil
}

// DeprecateThingType marks a thing type as deprecated. Deprecated thing types
// cannot be associated with new things.
func (s *IoTService) DeprecateThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	undoDeprecate := request.GetBoolParam(req.Parameters, "undoDeprecate")
	if err := s.deprecateThingTypeCore(store, thingTypeName, undoDeprecate); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func parseSearchableAttributes(props map[string]interface{}) []string {
	raw, ok := props["searchableAttributes"]
	if !ok {
		return nil
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(list))
	for _, item := range list {
		if name, ok := item.(string); ok && name != "" {
			result = append(result, name)
		}
	}
	return result
}
