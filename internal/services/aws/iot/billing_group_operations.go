package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// CreateBillingGroup creates a new billing group.
func (s *IoTService) CreateBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	props := unwrapProps(req.Parameters, "billingGroupProperties")
	in := CreateBillingGroupInput{
		GroupName:   request.GetParamCaseInsensitive(req.Parameters, "billingGroupName"),
		Description: request.GetParamCaseInsensitive(props, "billingGroupDescription"),
	}

	created, err := s.createBillingGroupCore(store, in)
	if err != nil {
		return nil, err
	}

	return billingGroupResponse(created), nil
}

// DescribeBillingGroup retrieves the details of an existing billing group.
func (s *IoTService) DescribeBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	bg, err := s.describeBillingGroupCore(store, groupName)
	if err != nil {
		return nil, err
	}

	return billingGroupDescribeResponse(bg), nil
}

// UpdateBillingGroup modifies the description or attributes of an existing
// billing group.
func (s *IoTService) UpdateBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	props := unwrapProps(req.Parameters, "billingGroupProperties")
	in := UpdateBillingGroupInput{
		GroupName:          request.GetParamCaseInsensitive(req.Parameters, "billingGroupName"),
		Description:        request.GetParamCaseInsensitive(props, "billingGroupDescription"),
		ExpectedVersion:    int64(request.GetIntParam(req.Parameters, "expectedVersion")),
		PropertiesProvided: request.GetMapParamCaseInsensitive(req.Parameters, "billingGroupProperties") != nil,
	}

	updated, err := s.updateBillingGroupCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"version": updated.Version,
	}, nil
}

// DeleteBillingGroup removes a billing group from the registry.
func (s *IoTService) DeleteBillingGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	groupName := request.GetParamCaseInsensitive(req.Parameters, "billingGroupName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteBillingGroupCore(store, groupName); err != nil {
		return nil, err
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
	result, err := s.listBillingGroupsCore(store, opts)
	if err != nil {
		return nil, err
	}

	groups := make([]map[string]interface{}, 0, len(result.Items))
	for _, g := range result.Items {
		groups = append(groups, groupNameAndArnResponse(g.GroupName, g.GroupARN))
	}

	return listResponse("billingGroups", groups, result.NextMarker), nil
}
