package iot

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
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

	if _, err := store.GetThing(thingName); err == nil {
		return nil, iotstore.ErrThingAlreadyExists
	}

	attributes, _, _ := parseAttributePayload(req.Parameters)
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")

	if thingTypeName != "" {
		if _, err := store.GetThingType(thingTypeName); err != nil {
			return nil, iotstore.ErrThingTypeNotFound
		}
	}

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
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"thingName":        created.ThingName,
		"thingArn":         created.ThingARN,
		"thingId":          created.ThingID,
		"thingTypeName":    created.ThingTypeName,
		"attributes":       created.Attributes,
		"attributeNames":   created.AttributeNames,
		"version":          created.Version,
		"creationDate":     created.CreationDate.Unix(),
		"lastModifiedDate": created.LastModifiedDate.Unix(),
	}, nil
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
		return nil, iotstore.ErrThingNotFound
	}

	return map[string]interface{}{
		"thingName":        thing.ThingName,
		"thingArn":         thing.ThingARN,
		"thingId":          thing.ThingID,
		"thingTypeName":    thing.ThingTypeName,
		"attributes":       ensureMap(thing.Attributes),
		"attributeNames":   thing.AttributeNames,
		"version":          thing.Version,
		"creationDate":     thing.CreationDate.Unix(),
		"lastModifiedDate": thing.LastModifiedDate.Unix(),
	}, nil
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

	existing, err := store.GetThing(thingName)
	if err != nil {
		return nil, iotstore.ErrThingNotFound
	}

	attributes, merge, payloadProvided := parseAttributePayload(req.Parameters)
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")

	isRemove := request.GetBoolParam(req.Parameters, "removeThingType")
	if isRemove {
		existing.ThingTypeName = ""
	} else if thingTypeName != "" {
		if _, err := store.GetThingType(thingTypeName); err != nil {
			return nil, iotstore.ErrThingTypeNotFound
		}
		existing.ThingTypeName = thingTypeName
	}

	if payloadProvided {
		if merge {
			if existing.Attributes == nil {
				existing.Attributes = make(map[string]string)
			}
			for k, v := range attributes {
				if v == "" {
					delete(existing.Attributes, k)
				} else {
					existing.Attributes[k] = v
				}
			}
		} else {
			existing.Attributes = make(map[string]string)
			for k, v := range attributes {
				if v != "" {
					existing.Attributes[k] = v
				}
			}
		}
		existing.AttributeNames = mapKeys(existing.Attributes)
	}

	existing.Version++
	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateThing(existing); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"thingName":        existing.ThingName,
		"thingArn":         existing.ThingARN,
		"thingId":          existing.ThingID,
		"thingTypeName":    existing.ThingTypeName,
		"attributes":       existing.Attributes,
		"attributeNames":   existing.AttributeNames,
		"version":          existing.Version,
		"creationDate":     existing.CreationDate.Unix(),
		"lastModifiedDate": existing.LastModifiedDate.Unix(),
	}, nil
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

	if err := store.DeleteThing(thingName); err != nil {
		return nil, iotstore.ErrThingNotFound
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
	result, err := store.ListThings(opts)
	if err != nil {
		return nil, err
	}

	attributeName := request.GetParamCaseInsensitive(req.Parameters, "attributeName")
	attributeValue := request.GetParamCaseInsensitive(req.Parameters, "attributeValue")

	things := make([]map[string]interface{}, 0, len(result.Items))
	for _, t := range result.Items {
		if attributeName != "" {
			val, exists := t.Attributes[attributeName]
			if !exists || val != attributeValue {
				continue
			}
		}
		things = append(things, map[string]interface{}{
			"thingName":        t.ThingName,
			"thingArn":         t.ThingARN,
			"thingId":          t.ThingID,
			"thingTypeName":    t.ThingTypeName,
			"attributes":       t.Attributes,
			"attributeNames":   t.AttributeNames,
			"version":          t.Version,
			"creationDate":     t.CreationDate.Unix(),
			"lastModifiedDate": t.LastModifiedDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"things": things,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
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
		things = append(things, map[string]interface{}{
			"thingName":      t.ThingName,
			"thingArn":       t.ThingARN,
			"thingId":        t.ThingID,
			"thingTypeName":  t.ThingTypeName,
			"attributes":     t.Attributes,
			"attributeNames": t.AttributeNames,
			"version":        t.Version,
		})
	}

	resp := map[string]interface{}{
		"things": things,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
}

// store returns the IotStore for the request region, creating one lazily if needed.
func (s *IoTService) store(reqCtx *request.RequestContext) (iotstore.IotStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (iotstore.IotStoreInterface, error) {
		st, err := reqCtx.GetStorage()
		if err != nil {
			return nil, err
		}
		return iotstore.NewIotStore(st, s.accountID, reqCtx.GetRegion()), nil
	})
}

// parseAttributePayload extracts the attributes map from an attributePayload parameter.
func parseAttributePayload(params map[string]interface{}) (map[string]string, bool, bool) {
	attrs := make(map[string]string)

	// Try map form first (AWS SDK JSON protocol passes nested structs as maps).
	if m := request.GetMapParamCaseInsensitive(params, "attributePayload"); m != nil {
		if a, ok := m["attributes"].(map[string]interface{}); ok {
			for k, v := range a {
				if s, ok := v.(string); ok {
					attrs[k] = s
				}
			}
			merge := true
			if mb, ok := m["merge"].(bool); ok {
				merge = mb
			} else if ms, ok := m["merge"].(string); ok {
				merge = strings.ToLower(ms) == "true"
			}
			return attrs, merge, true
		}
	}

	attrStr := request.GetParamCaseInsensitive(params, "attributePayload")
	if attrStr == "" {
		return attrs, true, false
	}

	var payload struct {
		Attributes map[string]string `json:"attributes"`
		Merge      *bool             `json:"merge"`
	}
	if json.Unmarshal([]byte(attrStr), &payload) == nil && payload.Attributes != nil {
		merge := true
		if payload.Merge != nil {
			merge = *payload.Merge
		}
		return payload.Attributes, merge, true
	}

	var direct map[string]string
	if json.Unmarshal([]byte(attrStr), &direct) == nil {
		return direct, true, true
	}

	return attrs, true, false
}

// mapKeys returns the keys of a string map as a sorted slice.
func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// parseListOptions extracts pagination parameters (nextToken, maxResults) from
// the request parameters into ListOptions.
func parseListOptions(params map[string]interface{}) storecommon.ListOptions {
	opts := storecommon.ListOptions{}
	if token := request.GetParamCaseInsensitive(params, "nextToken"); token != "" {
		opts.Marker = token
	}
	if max := request.GetIntParam(params, "maxResults"); max > 0 {
		opts.MaxItems = max
	}
	return opts
}

func ensureMap(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}

// unwrapProps extracts a nested properties wrapper from request parameters.
// The AWS IoT awsJson1_1 protocol wraps operation parameters in structures
// (e.g., topicRulePayload, thingGroupProperties). The framework's JSON parser
// stores these as map[string]interface{} values at the top level of
// req.Parameters, so handlers must explicitly unwrap them.
// If wrapperKey is not found in params, the original params are returned
// unchanged.
func unwrapProps(params map[string]interface{}, wrapperKey string) map[string]interface{} {
	if m := request.GetMapParamCaseInsensitive(params, wrapperKey); m != nil {
		return m
	}
	return params
}
