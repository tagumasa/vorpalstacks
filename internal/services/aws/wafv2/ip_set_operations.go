package wafv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateIPSet creates a new IP set containing the specified IP addresses.
func (s *WAFv2Service) CreateIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, validationError("Name is required")
	}

	_ = request.GetStringParam(req.Parameters, "Scope")

	ipAddressVersion := request.GetStringParam(req.Parameters, "IPAddressVersion")
	if ipAddressVersion == "" {
		ipAddressVersion = "IPV4"
	}

	var addresses []string
	if addrRaw := req.Parameters["Addresses"]; addrRaw != nil {
		if arr, ok := addrRaw.([]interface{}); ok {
			for _, a := range arr {
				if s, ok := a.(string); ok {
					addresses = append(addresses, s)
				}
			}
		}
	}

	description := request.GetStringParam(req.Parameters, "Description")

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	ipSet, err := stores.ipSets.Create(id, name, description, ipAddressVersion, addresses)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WafV2AlreadyExistsException", "IPSet already exists", 400)
		}
		return nil, err
	}

	if tags := tagutil.ParseTags(req.Parameters, "Tags"); len(tags) > 0 {
		if err := stores.tags.TagFromSlice(ipSet.ARN, tags); err != nil {
			logs.Warn("failed to persist tags for IPSet", logs.String("id", ipSet.ID), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"Summary": map[string]interface{}{
			"Id":          ipSet.ID,
			"Name":        ipSet.Name,
			"ARN":         ipSet.ARN,
			"Description": ipSet.Description,
			"LockToken":   ipSet.LockToken,
		},
	}, nil
}

// GetIPSet retrieves the details of the specified IP set.
func (s *WAFv2Service) GetIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, validationError("Id is required")
	}

	ipSet, err := stores.ipSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("IPSet")
		}
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
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	maxItems := request.GetIntParam(req.Parameters, "Limit")
	if maxItems == 0 {
		maxItems = 100
	}

	nextMarker := request.GetStringParam(req.Parameters, "NextMarker")

	result, err := stores.ipSets.List(nextMarker, maxItems)
	if err != nil {
		return nil, err
	}

	ipSets := make([]interface{}, 0, len(result.IPSets))
	for _, ips := range result.IPSets {
		ipSets = append(ipSets, map[string]interface{}{
			"Id":          ips.ID,
			"Name":        ips.Name,
			"ARN":         ips.ARN,
			"Description": ips.Description,
			"LockToken":   ips.LockToken,
		})
	}

	resp := map[string]interface{}{
		"IPSets": ipSets,
	}
	if result.NextMarker != "" {
		resp["NextMarker"] = result.NextMarker
	}
	return resp, nil
}

// UpdateIPSet updates the specified IP set with new addresses, returning a new lock token.
func (s *WAFv2Service) UpdateIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, validationError("Id is required")
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, validationError("LockToken is required")
	}

	var addresses []string
	if addrRaw := req.Parameters["Addresses"]; addrRaw != nil {
		if arr, ok := addrRaw.([]interface{}); ok {
			for _, a := range arr {
				if s, ok := a.(string); ok {
					addresses = append(addresses, s)
				}
			}
		}
	}

	ipSet, err := stores.ipSets.Update(id, lockToken, addresses)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("IPSet")
		}
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": ipSet.LockToken,
	}, nil
}

// DeleteIPSet permanently deletes the specified IP set.
func (s *WAFv2Service) DeleteIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, validationError("Id is required")
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, validationError("LockToken is required")
	}

	ipSet, err := stores.ipSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("IPSet")
		}
		return nil, err
	}

	if err := stores.ipSets.Delete(id, lockToken); err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		return nil, err
	}

	stores.tags.Delete(ipSet.ARN)

	return map[string]interface{}{}, nil
}
