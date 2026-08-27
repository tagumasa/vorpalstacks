package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateThingGroup creates a new thing group with optional parent group and
// attributes. Returns ResourceAlreadyExistsException if the group name is
// already taken.
func (s *IoTService) CreateThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	props := unwrapProps(req.Parameters, "thingGroupProperties")
	attrs, _, _ := parseAttributePayload(props)
	in := CreateThingGroupInput{
		GroupName:       request.GetParamCaseInsensitive(req.Parameters, "thingGroupName"),
		ParentGroupName: request.GetParamCaseInsensitive(req.Parameters, "parentGroupName"),
		Description:     request.GetParamCaseInsensitive(props, "thingGroupDescription"),
		Attributes:      attrs,
	}

	created, err := s.createThingGroupCore(store, in)
	if err != nil {
		return nil, err
	}

	return thingGroupResponse(created), nil
}

// DescribeThingGroup retrieves the details of an existing thing group.
func (s *IoTService) DescribeThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "thingGroupName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	group, err := s.describeThingGroupCore(store, groupName)
	if err != nil {
		return nil, err
	}

	return thingGroupDescribeResponse(group), nil
}

// UpdateThingGroup modifies the attributes and description of an existing thing
// group.
func (s *IoTService) UpdateThingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	props := unwrapProps(req.Parameters, "thingGroupProperties")
	attrs, _, _ := parseAttributePayload(props)

	in := UpdateThingGroupInput{
		GroupName:          request.GetParamCaseInsensitive(req.Parameters, "thingGroupName"),
		Description:        request.GetParamCaseInsensitive(props, "thingGroupDescription"),
		Attributes:         attrs,
		ExpectedVersion:    int64(request.GetIntParam(req.Parameters, "expectedVersion")),
		PropertiesProvided: request.GetMapParamCaseInsensitive(req.Parameters, "thingGroupProperties") != nil,
	}

	updated, err := s.updateThingGroupCore(store, in)
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteThingGroupCore(store, groupName); err != nil {
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
	result, err := s.listThingGroupsCore(store, opts, parentFilter)
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

	if err := s.addThingToThingGroupCore(store, thingName, groupName); err != nil {
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

	if err := s.removeThingFromThingGroupCore(store, thingName, groupName); err != nil {
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

	things, err := s.listThingsInThingGroupCore(store, groupName)
	if err != nil {
		return nil, err
	}

	return paginatedStrings("things", things, req.Parameters)
}

// ListThingGroupsForThing returns groups containing the specified thing.
func (s *IoTService) ListThingGroupsForThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")

	groups, err := s.listThingGroupsForThingCore(store, thingName)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(groups))
	for _, g := range groups {
		result = append(result, map[string]interface{}{
			"groupName": g,
		})
	}

	return paginatedMaps("thingGroups", result, req.Parameters)
}

// ---------------------------------------------------------------------------
// UpdateThingGroupsForThing: atomically add/remove thing group memberships.
// ---------------------------------------------------------------------------

func (s *IoTService) UpdateThingGroupsForThing(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	thingName := request.GetParamCaseInsensitive(req.Parameters, "thingName")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	add := request.GetStringList(req.Parameters, "thingGroupsToAdd")
	remove := request.GetStringList(req.Parameters, "thingGroupsToRemove")
	if err := s.updateThingGroupsForThingCore(store, thingName, add, remove); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
