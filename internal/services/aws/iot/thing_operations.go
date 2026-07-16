package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// CreateThing creates a new IoT thing with the given name, optional attributes
// and thing type. Returns ResourceAlreadyExistsException if the thing name is
// already taken.
func (s *IoTService) CreateThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	attributes, _, _ := parseAttributePayload(req.Parameters)
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")

	thing := &iotstore.Thing{
		ThingName:        thingName,
		ThingTypeName:    thingTypeName,
		Attributes:       attributes,
		AttributeNames:   mapKeys(attributes),
		Version:          1,
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateThing(thing)
	if err != nil {
		return nil, err
	}

	return thingResponse(created), nil
}

// DescribeThing retrieves the details of an existing IoT thing including
// attributes, version, and metadata.
func (s *IoTService) DescribeThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	thing, err := store.GetThing(thingName)
	if err != nil {
		return nil, err
	}

	resp := thingDescribeResponse(thing)

	// AWS: DescribeThing returns billingGroupName when the thing belongs to
	// a billing group (at most one per AWS constraints).
	if groups, _ := store.ListBillingGroupsForThing(thingName); len(groups) > 0 {
		resp["billingGroupName"] = groups[0]
	}

	return resp, nil
}

// UpdateThing modifies the attributes and thing type of an existing thing.
// Supports attribute merging and removal. Increments the version counter.
func (s *IoTService) UpdateThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	attributes, merge, payloadProvided := parseAttributePayload(req.Parameters)
	opts := iotstore.ThingUpdateOpts{
		Attributes:      attributes,
		MergeAttributes: merge,
		PayloadProvided: payloadProvided,
		ThingTypeName:   request.GetParamCaseInsensitive(req.Parameters, "thingTypeName"),
		RemoveThingType: request.GetBoolParam(req.Parameters, "removeThingType"),
	}

	updated, err := store.UpdateThing(thingName, opts)
	if err != nil {
		return nil, err
	}

	return thingResponse(updated), nil
}

// DeleteThing removes a thing from the registry. Returns ResourceNotFoundException
// if the thing does not exist.
func (s *IoTService) DeleteThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	arn := iotstore.BuildThingARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), thingName)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteThing(thingName); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// ListThings returns a paginated list of things, optionally filtered by
// attribute name and value.
func (s *IoTService) ListThings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	attributeName := request.GetParamCaseInsensitive(req.Parameters, "attributeName")
	attributeValue := request.GetParamCaseInsensitive(req.Parameters, "attributeValue")

	result, err := store.ListThings(opts, attributeName, attributeValue)
	if err != nil {
		return nil, err
	}

	things := make([]map[string]interface{}, 0, len(result.Items))
	for _, t := range result.Items {
		things = append(things, thingResponse(t))
	}

	return listResponse("things", things, result.NextMarker), nil
}

// ListThingsForThingType returns things that belong to the specified thing type.
func (s *IoTService) ListThingsForThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")
	if thingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}

	opts := parseListOptions(req.Parameters)
	result, err := store.ListThingsForThingType(thingTypeName, opts)
	if err != nil {
		return nil, err
	}

	things := make([]map[string]interface{}, 0, len(result.Items))
	for _, t := range result.Items {
		things = append(things, thingResponse(t))
	}

	return listResponse("things", things, result.NextMarker), nil
}

// store returns the singleton IotStore for the request region, shared with
// iotevents and the MQTT auth provider via iotstore.GetOrCreateStore.
func (s *IoTService) store(reqCtx *request.RequestContext) (iotstore.IotStoreInterface, error) {
	st, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	return iotstore.GetOrCreateStore(st, s.accountID, reqCtx.GetRegion()), nil
}
