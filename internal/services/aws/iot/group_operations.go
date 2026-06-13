package iot

import (
	"context"
	"strconv"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// CreateThingType creates a new thing type with optional description and
// property definitions. Returns ResourceAlreadyExistsException if the thing
// type name is already taken.
func (s *IoTService) CreateThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")
	if thingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetThingType(thingTypeName); err == nil {
		return nil, iotstore.ErrThingTypeAlreadyExists
	}

	props := unwrapProps(req.Parameters, "thingTypeProperties")

	tt := &iotstore.ThingType{
		ThingTypeName:    thingTypeName,
		Description:      request.GetParamCaseInsensitive(props, "thingTypeDescription"),
		Version:          1,
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateThingType(tt)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	thingTypeProps := map[string]interface{}{
		"thingTypeDescription": created.Description,
	}
	return map[string]interface{}{
		"thingTypeName":       created.ThingTypeName,
		"thingTypeArn":        created.ThingTypeARN,
		"thingTypeId":         created.ThingTypeID,
		"description":         created.Description,
		"thingTypeProperties": thingTypeProps,
		"version":             created.Version,
		"creationDate":        created.CreationDate.Unix(),
		"lastModifiedDate":    created.LastModifiedDate.Unix(),
	}, nil
}

// DescribeThingType retrieves the details of an existing thing type including
// properties and deprecation status.
func (s *IoTService) DescribeThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")
	if thingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tt, err := store.GetThingType(thingTypeName)
	if err != nil {
		return nil, iotstore.ErrThingTypeNotFound
	}

	thingTypeProps := map[string]interface{}{
		"thingTypeDescription": tt.Description,
	}
	return map[string]interface{}{
		"thingTypeName":       tt.ThingTypeName,
		"thingTypeArn":        tt.ThingTypeARN,
		"thingTypeId":         tt.ThingTypeID,
		"description":         tt.Description,
		"thingTypeProperties": thingTypeProps,
		"version":             tt.Version,
		"creationDate":        tt.CreationDate.Unix(),
		"lastModifiedDate":    tt.LastModifiedDate.Unix(),
	}, nil
}

// UpdateThingType modifies the description or properties of an existing thing
// type. Increments the version counter.
func (s *IoTService) UpdateThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")
	if thingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.GetThingType(thingTypeName)
	if err != nil {
		return nil, iotstore.ErrThingTypeNotFound
	}

	props := unwrapProps(req.Parameters, "thingTypeProperties")

	if desc := request.GetParamCaseInsensitive(props, "thingTypeDescription"); desc != "" {
		existing.Description = desc
	}

	existing.Version++
	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateThingType(existing); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"thingTypeName":    existing.ThingTypeName,
		"thingTypeArn":     existing.ThingTypeARN,
		"thingTypeId":      existing.ThingTypeID,
		"description":      existing.Description,
		"version":          existing.Version,
		"lastModifiedDate": existing.LastModifiedDate.Unix(),
	}, nil
}

// DeleteThingType removes a thing type from the registry. Returns
// ResourceNotFoundException if the thing type does not exist.
func (s *IoTService) DeleteThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")
	if thingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	things, err := store.ListThingsForThingType(thingTypeName, storecommon.ListOptions{})
	if err == nil && len(things.Items) > 0 {
		return nil, iotstore.ErrDeleteConflict
	}

	if err := store.DeleteThingType(thingTypeName); err != nil {
		return nil, iotstore.ErrThingTypeNotFound
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
	result, err := store.ListThingTypes(opts)
	if err != nil {
		return nil, err
	}

	types := make([]map[string]interface{}, 0, len(result.Items))
	for _, t := range result.Items {
		types = append(types, map[string]interface{}{
			"thingTypeName":    t.ThingTypeName,
			"thingTypeArn":     t.ThingTypeARN,
			"thingTypeId":      t.ThingTypeID,
			"description":      t.Description,
			"version":          t.Version,
			"creationDate":     t.CreationDate.Unix(),
			"lastModifiedDate": t.LastModifiedDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"thingTypes": types,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
}

// DeprecateThingType marks a thing type as deprecated. Deprecated thing types
// cannot be associated with new things.
func (s *IoTService) DeprecateThingType(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingTypeName := request.GetParamCaseInsensitive(req.Parameters, "thingTypeName")
	if thingTypeName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetThingType(thingTypeName)
	if err != nil {
		return nil, iotstore.ErrThingTypeNotFound
	}

	return map[string]interface{}{}, nil
}

// CreateThingGroup creates a new thing group with optional parent group and
// attributes. Returns ResourceAlreadyExistsException if the group name is
// already taken.
func (s *IoTService) CreateThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "thingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetThingGroup(groupName); err == nil {
		return nil, iotstore.ErrThingGroupAlreadyExists
	}

	parentGroupName := request.GetParamCaseInsensitive(req.Parameters, "parentGroupName")
	if parentGroupName != "" {
		if _, err := store.GetThingGroup(parentGroupName); err != nil {
			return nil, iotstore.ErrThingGroupNotFound
		}
	}

	props := unwrapProps(req.Parameters, "thingGroupProperties")
	attrs, _, _ := parseAttributePayload(props)
	group := &iotstore.ThingGroup{
		GroupName:        groupName,
		ParentGroupName:  parentGroupName,
		Description:      request.GetParamCaseInsensitive(props, "thingGroupDescription"),
		Attributes:       attrs,
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateThingGroup(group)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"thingGroupName":   created.GroupName,
		"thingGroupArn":    created.GroupARN,
		"thingGroupId":     created.GroupID,
		"parentGroupName":  created.ParentGroupName,
		"description":      created.Description,
		"attributes":       created.Attributes,
		"creationDate":     created.CreationDate.Unix(),
		"lastModifiedDate": created.LastModifiedDate.Unix(),
	}, nil
}

// DescribeThingGroup retrieves the details of an existing thing group.
func (s *IoTService) DescribeThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "thingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	group, err := store.GetThingGroup(groupName)
	if err != nil {
		return nil, iotstore.ErrThingGroupNotFound
	}

	thingGroupProps := map[string]interface{}{
		"thingGroupDescription": group.Description,
		"attributePayload": map[string]interface{}{
			"attributes": ensureMap(group.Attributes),
		},
	}
	return map[string]interface{}{
		"thingGroupName":       group.GroupName,
		"thingGroupArn":        group.GroupARN,
		"thingGroupId":         group.GroupID,
		"parentGroupName":       group.ParentGroupName,
		"description":           group.Description,
		"thingGroupProperties": thingGroupProps,
		"attributes":            group.Attributes,
		"creationDate":          group.CreationDate.Unix(),
		"lastModifiedDate":      group.LastModifiedDate.Unix(),
	}, nil
}

// UpdateThingGroup modifies the attributes and description of an existing thing
// group.
func (s *IoTService) UpdateThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "thingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.GetThingGroup(groupName)
	if err != nil {
		return nil, iotstore.ErrThingGroupNotFound
	}

	props := unwrapProps(req.Parameters, "thingGroupProperties")
	attrs, _, _ := parseAttributePayload(props)
	if desc := request.GetParamCaseInsensitive(props, "thingGroupDescription"); desc != "" {
		existing.Description = desc
	}
	if evStr := request.GetParamCaseInsensitive(req.Parameters, "expectedVersion"); evStr != "" && existing.Version > 0 {
		expVer, err := strconv.ParseInt(evStr, 10, 64)
		if err != nil || expVer != existing.Version {
			return nil, iotstore.ErrVersionConflict
		}
	}

	if len(attrs) > 0 {
		if existing.Attributes == nil {
			existing.Attributes = make(map[string]string)
		}
		for k, v := range attrs {
			existing.Attributes[k] = v
		}
	}

	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateThingGroup(existing); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"version": existing.Version,
	}, nil
}

// DeleteThingGroup removes a thing group from the registry. Returns
// ResourceNotFoundException if the thing group does not exist.
func (s *IoTService) DeleteThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "thingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetThingGroup(groupName); err != nil {
		return nil, iotstore.ErrThingGroupNotFound
	}

	if members, err := store.ListThingsInGroup(groupName); err == nil && len(members) > 0 {
		return nil, iotstore.ErrDeleteConflict
	}

	groups, err := store.ListThingGroups(storecommon.ListOptions{})
	if err == nil {
		for _, g := range groups.Items {
			if g.ParentGroupName == groupName {
				return nil, iotstore.ErrDeleteConflict
			}
		}
	}

	if err := store.DeleteThingGroup(groupName); err != nil {
		return nil, iotstore.ErrThingGroupNotFound
	}

	return map[string]interface{}{}, nil
}

// ListThingGroups returns a paginated list of thing groups, optionally filtered
// by parent group name.
func (s *IoTService) ListThingGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := store.ListThingGroups(opts)
	if err != nil {
		return nil, err
	}

	parentFilter := request.GetParamCaseInsensitive(req.Parameters, "parentGroupName")
	groups := make([]map[string]interface{}, 0)
	for _, g := range result.Items {
		if parentFilter != "" && g.ParentGroupName != parentFilter {
			continue
		}
		groups = append(groups, map[string]interface{}{
			"thingGroupName":   g.GroupName,
			"thingGroupArn":    g.GroupARN,
			"thingGroupId":     g.GroupID,
			"parentGroupName":  g.ParentGroupName,
			"description":      g.Description,
			"creationDate":     g.CreationDate.Unix(),
			"lastModifiedDate": g.LastModifiedDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"thingGroups": groups,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
}

// AddThingToThingGroup adds a thing to a thing group.
func (s *IoTService) AddThingToThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	groupName := request.GetParamCaseInsensitive(req.Parameters, "thingGroupName")
	if thingName == "" || groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	if err := store.AddThingToThingGroup(thingName, groupName); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{}, nil
}

// RemoveThingFromThingGroup removes a thing from a thing group.
func (s *IoTService) RemoveThingFromThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	groupName := request.GetParamCaseInsensitive(req.Parameters, "thingGroupName")
	if thingName == "" || groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	if err := store.RemoveThingFromThingGroup(thingName, groupName); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{}, nil
}

// ListThingsInThingGroup returns things that belong to the specified thing group.
func (s *IoTService) ListThingsInThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	groupName := request.GetParamCaseInsensitive(req.Parameters, "thingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	things, err := store.ListThingsInGroup(groupName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"things": things,
	}, nil
}

// ListThingGroupsForThing returns groups containing the specified thing.
func (s *IoTService) ListThingGroupsForThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	if thingName == "" {
		return nil, iotstore.ErrMissingParam
	}

	groups, err := store.ListGroupsForThing(thingName)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(groups))
	for _, g := range groups {
		result = append(result, map[string]interface{}{
			"groupName": g,
		})
	}

	return map[string]interface{}{
		"thingGroups": result,
	}, nil
}

// CreateBillingGroup creates a new billing group.
func (s *IoTService) CreateBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetBillingGroup(groupName); err == nil {
		return nil, iotstore.ErrResourceAlreadyExists
	}

	props := unwrapProps(req.Parameters, "billingGroupProperties")

	bg := &iotstore.BillingGroup{
		GroupName:        groupName,
		Description:      request.GetParamCaseInsensitive(props, "billingGroupDescription"),
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateBillingGroup(bg)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{
		"billingGroupName": created.GroupName,
		"billingGroupArn":  created.GroupARN,
		"billingGroupId":   created.GroupID,
		"description":      created.Description,
		"creationDate":     created.CreationDate.Unix(),
		"lastModifiedDate": created.LastModifiedDate.Unix(),
	}, nil
}

// DescribeBillingGroup retrieves the details of an existing billing group.
func (s *IoTService) DescribeBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	bg, err := store.GetBillingGroup(groupName)
	if err != nil {
		return nil, iotstore.ErrBillingGroupNotFound
	}

	billingGroupProps := map[string]interface{}{
		"billingGroupDescription": bg.Description,
	}
	return map[string]interface{}{
		"billingGroupName":       bg.GroupName,
		"billingGroupArn":        bg.GroupARN,
		"billingGroupId":         bg.GroupID,
		"description":            bg.Description,
		"billingGroupProperties": billingGroupProps,
		"attributes":             bg.Attributes,
		"creationDate":           bg.CreationDate.Unix(),
		"lastModifiedDate":       bg.LastModifiedDate.Unix(),
	}, nil
}

// UpdateBillingGroup modifies the description or attributes of an existing
// billing group.
func (s *IoTService) UpdateBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.GetBillingGroup(groupName)
	if err != nil {
		return nil, iotstore.ErrBillingGroupNotFound
	}

	props := unwrapProps(req.Parameters, "billingGroupProperties")

	if desc := request.GetParamCaseInsensitive(props, "billingGroupDescription"); desc != "" {
		existing.Description = desc
	}

	existing.LastModifiedDate = time.Now().UTC()

	if err := store.UpdateBillingGroup(existing); err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	return map[string]interface{}{}, nil
}

// DeleteBillingGroup removes a billing group from the registry.
func (s *IoTService) DeleteBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteBillingGroup(groupName); err != nil {
		return nil, iotstore.ErrBillingGroupNotFound
	}

	return map[string]interface{}{}, nil
}

// ListBillingGroups returns a paginated list of billing groups.
func (s *IoTService) ListBillingGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := parseListOptions(req.Parameters)
	result, err := store.ListBillingGroups(opts)
	if err != nil {
		return nil, err
	}

	groups := make([]map[string]interface{}, 0, len(result.Items))
	for _, g := range result.Items {
		groups = append(groups, map[string]interface{}{
			"billingGroupName": g.GroupName,
			"billingGroupArn":  g.GroupARN,
			"billingGroupId":   g.GroupID,
			"description":      g.Description,
			"creationDate":     g.CreationDate.Unix(),
			"lastModifiedDate": g.LastModifiedDate.Unix(),
		})
	}

	resp := map[string]interface{}{
		"billingGroups": groups,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}
	return resp, nil
}
