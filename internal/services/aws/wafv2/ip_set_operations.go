package wafv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateIPSet creates a new IP set containing the specified IP addresses.
func (s *WAFv2Service) CreateIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	ipSet, err := s.createIPSetCore(stores, IPSetCreateInput{
		Name:             request.GetStringParam(req.Parameters, "Name"),
		Scope:            request.GetStringParam(req.Parameters, "Scope"),
		IPAddressVersion: request.GetStringParam(req.Parameters, "IPAddressVersion"),
		AddressesRaw:     req.Parameters["Addresses"],
		Description:      request.GetStringParam(req.Parameters, "Description"),
		Tags:             tagutil.ParseTags(req.Parameters, "Tags"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Summary": buildIPSetSummary(ipSet),
	}, nil
}

// GetIPSet retrieves the details of the specified IP set.
func (s *WAFv2Service) GetIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	ipSet, err := s.getIPSetCore(stores, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"IPSet": map[string]interface{}{
			"Id":               ipSet.ID,
			"Name":             ipSet.Name,
			"ARN":              ipSet.ARN,
			"IPAddressVersion": ipSet.IPAddressVersion,
			"Addresses":        ipSet.Addresses,
			"Description":      ipSet.Description,
		},
		"LockToken": ipSet.LockToken,
	}, nil
}

// ListIPSets returns a paginated list of all IP sets.
func (s *WAFv2Service) ListIPSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	result, err := s.listIPSetsCore(stores, IPSetListInput{
		Scope:      request.GetStringParam(req.Parameters, "Scope"),
		Limit:      pagination.GetMaxItems(req.Parameters, 100, "Limit"),
		NextMarker: pagination.GetMarker(req.Parameters, "NextMarker"),
	})
	if err != nil {
		return nil, err
	}

	ipSets := make([]interface{}, 0, len(result.IPSets))
	for _, ips := range result.IPSets {
		ipSets = append(ipSets, buildIPSetSummary(ips))
	}

	resp := map[string]interface{}{
		"IPSets": ipSets,
	}
	pagination.SetNextToken(resp, "NextMarker", result.NextMarker)
	return resp, nil
}

// UpdateIPSet updates the specified IP set with new addresses, returning a new lock token.
func (s *WAFv2Service) UpdateIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	ipSet, err := s.updateIPSetCore(stores, IPSetUpdateInput{
		Id:           request.GetStringParam(req.Parameters, "Id"),
		LockToken:    request.GetStringParam(req.Parameters, "LockToken"),
		AddressesRaw: req.Parameters["Addresses"],
		Description:  request.GetStringParam(req.Parameters, "Description"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": ipSet.LockToken,
	}, nil
}

// DeleteIPSet permanently deletes the specified IP set.
func (s *WAFv2Service) DeleteIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	if err := s.deleteIPSetCore(stores, request.GetStringParam(req.Parameters, "Id"), request.GetStringParam(req.Parameters, "LockToken")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
