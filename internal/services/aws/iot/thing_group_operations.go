package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

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

	parentGroupName := request.GetParamCaseInsensitive(req.Parameters, "parentGroupName")

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
		return nil, err
	}

	return thingGroupResponse(created), nil
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
		return nil, err
	}

	return thingGroupDescribeResponse(group), nil
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

	props := unwrapProps(req.Parameters, "thingGroupProperties")
	attrs, _, _ := parseAttributePayload(props)

	expectedVersion := int64(request.GetIntParam(req.Parameters, "expectedVersion"))
	opts := iotstore.ThingGroupUpdateOpts{
		Description:     request.GetParamCaseInsensitive(props, "thingGroupDescription"),
		Attributes:      attrs,
		ExpectedVersion: expectedVersion,
	}

	updated, err := store.UpdateThingGroup(groupName, opts)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"version": updated.Version,
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

	arn := iotstore.BuildThingGroupARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), groupName)
	_ = store.DeleteAllTags(arn)

	if err := store.DeleteThingGroup(groupName); err != nil {
		return nil, err
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
	parentFilter := request.GetParamCaseInsensitive(req.Parameters, "parentGroupName")
	result, err := store.ListThingGroups(opts, parentFilter)
	if err != nil {
		return nil, err
	}

	groups := make([]map[string]interface{}, 0, len(result.Items))
	for _, g := range result.Items {
		groups = append(groups, groupNameAndArnResponse(g.GroupName, g.GroupARN))
	}

	return listResponse("thingGroups", groups, result.NextMarker), nil
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
		return nil, err
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
		return nil, err
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

	return paginatedStrings("things", things, req.Parameters), nil
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

	return paginatedMaps("thingGroups", result, req.Parameters), nil
}
