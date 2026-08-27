package iot

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// CreateThing creates a new IoT thing with the given name, optional attributes
// and thing type. Returns ResourceAlreadyExistsException if the thing name is
// already taken.
func (s *IoTService) CreateThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	attributes, _, _ := parseAttributePayload(req.Parameters)
	result, err := s.createThingCore(store, CreateThingInput{
		ThingName:        request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		ThingTypeName:    request.GetParamCaseInsensitive(req.Parameters, "thingTypeName"),
		BillingGroupName: request.GetParamCaseInsensitive(req.Parameters, "billingGroupName"),
		Attributes:       attributes,
	})
	if err != nil {
		return nil, err
	}

	return thingResponse(result.Thing), nil
}

// DescribeThing retrieves the details of an existing IoT thing including
// attributes, version, and metadata.
func (s *IoTService) DescribeThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.describeThingCore(store, thingName)
	if err != nil {
		return nil, err
	}

	resp := thingDescribeResponse(result.Thing)
	if result.BillingGroupName != "" {
		resp["billingGroupName"] = result.BillingGroupName
	}

	return resp, nil
}

// UpdateThing modifies the attributes and thing type of an existing thing.
// Supports attribute merging and removal. Increments the version counter.
func (s *IoTService) UpdateThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	attributes, merge, payloadProvided := parseAttributePayload(req.Parameters)
	result, err := s.updateThingCore(store, UpdateThingInput{
		ThingName:       request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		Attributes:      attributes,
		MergeAttributes: merge,
		PayloadProvided: payloadProvided,
		ThingTypeName:   request.GetParamCaseInsensitive(req.Parameters, "thingTypeName"),
		RemoveThingType: request.GetBoolParam(req.Parameters, "removeThingType"),
	})
	if err != nil {
		return nil, err
	}

	return thingResponse(result.Thing), nil
}

// DeleteThing removes a thing from the registry. Returns ResourceNotFoundException
// if the thing does not exist.
func (s *IoTService) DeleteThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteThingCore(store, thingName, reqCtx.GetAccountID(), reqCtx.GetRegion()); err != nil {
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
	result, err := s.listThingsCore(store, ListThingsInput{
		AttributeName:  request.GetParamCaseInsensitive(req.Parameters, "attributeName"),
		AttributeValue: request.GetParamCaseInsensitive(req.Parameters, "attributeValue"),
		NextToken:      opts.Marker,
		MaxItems:       opts.MaxItems,
	})
	if err != nil {
		return nil, err
	}

	things := make([]map[string]interface{}, 0, len(result.Things))
	for _, t := range result.Things {
		things = append(things, thingResponse(t))
	}

	return listResponse("things", things, result.NextToken), nil
}

// ListThingsForThingType returns things that belong to the specified thing type.
func (s *IoTService) ListThingsForThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")
	opts := parseListOptions(req.Parameters)
	result, err := s.listThingsCore(store, ListThingsInput{
		ThingTypeName: thingTypeName,
		NextToken:     opts.Marker,
		MaxItems:      opts.MaxItems,
	})
	if err != nil {
		return nil, err
	}

	things := make([]map[string]interface{}, 0, len(result.Things))
	for _, t := range result.Things {
		things = append(things, thingResponse(t))
	}

	return listResponse("things", things, result.NextToken), nil
}

// store returns the singleton IotStore for the request region, shared with
// the MQTT auth provider via iotstore.GetOrCreateStore.
func (s *IoTService) store(reqCtx *request.RequestContext) (iotstore.IotStoreInterface, error) {
	st, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	return iotstore.GetOrCreateStore(st, s.accountID, reqCtx.GetRegion()), nil
}

// ---------------------------------------------------------------------------
// GetThingConnectivityData: returns MQTT connection status. Without a real
// MQTT broker feeding connection state, we return connected=false which is
// honest for a platform that has no active device connections.
// ---------------------------------------------------------------------------

func (s *IoTService) GetThingConnectivityData(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getThingConnectivityCore(store, thingName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"thingName":        result.ThingName,
		"connected":        result.Connected,
		"timestamp":        result.TimestampSeconds,
		"connectTime":      result.ConnectTime,
		"disconnectReason": "",
	}, nil
}

// ---------------------------------------------------------------------------
// V2 principal/thing listing with richer output.
// ---------------------------------------------------------------------------

func (s *IoTService) ListPrincipalThingsV2(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	principal := request.GetParamCaseInsensitive(req.Parameters, "principal")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	things, err := s.listPrincipalThingsV2Core(store, principal)
	if err != nil {
		return nil, err
	}
	objects := make([]map[string]interface{}, 0, len(things))
	for _, t := range things {
		objects = append(objects, map[string]interface{}{
			"thingName":          t,
			"thingPrincipalType": "EXCLUSIVE_THING",
		})
	}
	return paginatedMaps("principalThingObjects", objects, req.Parameters)
}

// extractCertIDFromPrincipal extracts the certificate ID from an IoT
// principal ARN (e.g. arn:aws:iot:us-east-1:123:cert/abcdef).
func extractCertIDFromPrincipal(principal string) string {
	idx := strings.LastIndex(principal, "cert/")
	if idx < 0 {
		return ""
	}
	return principal[idx+len("cert/"):]
}

func (s *IoTService) ListThingPrincipalsV2(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	principals, err := s.listThingPrincipalsV2Core(store, thingName)
	if err != nil {
		return nil, err
	}
	objects := make([]map[string]interface{}, 0, len(principals))
	for _, p := range principals {
		objects = append(objects, map[string]interface{}{
			"principal":          p,
			"thingPrincipalType": "EXCLUSIVE_THING",
		})
	}
	return paginatedMaps("thingPrincipalObjects", objects, req.Parameters)
}
