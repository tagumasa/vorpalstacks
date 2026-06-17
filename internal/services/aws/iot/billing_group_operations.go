package iot

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

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

	props := unwrapProps(req.Parameters, "billingGroupProperties")

	bg := &iotstore.BillingGroup{
		GroupName:        groupName,
		Description:      request.GetParamCaseInsensitive(props, "billingGroupDescription"),
		CreationDate:     time.Now().UTC(),
		LastModifiedDate: time.Now().UTC(),
	}

	created, err := store.CreateBillingGroup(bg)
	if err != nil {
		return nil, err
	}

	return billingGroupResponse(created), nil
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

	return billingGroupDescribeResponse(bg), nil
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

	props := unwrapProps(req.Parameters, "billingGroupProperties")

	expectedVersion := int64(request.GetIntParam(req.Parameters, "expectedVersion"))
	opts := iotstore.BillingGroupUpdateOpts{
		Description:     request.GetParamCaseInsensitive(props, "billingGroupDescription"),
		ExpectedVersion: expectedVersion,
	}

	updated, err := store.UpdateBillingGroup(groupName, opts)
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
	if groupName == "" {
		return nil, iotstore.ErrMissingParam
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteBillingGroup(groupName); err != nil {
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
	result, err := store.ListBillingGroups(opts)
	if err != nil {
		return nil, err
	}

	groups := make([]map[string]interface{}, 0, len(result.Items))
	for _, g := range result.Items {
		groups = append(groups, billingGroupResponse(g))
	}

	return listResponse("billingGroups", groups, result.NextMarker), nil
}
