package waf

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateIPSet creates a new IP set containing IP address descriptors.
func (s *WAFService) CreateIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidParameterValue", "Name is required", 400)
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	var descriptors []interface{}
	if descRaw := req.Parameters["IPSetDescriptors"]; descRaw != nil {
		if descArr, ok := descRaw.([]interface{}); ok {
			descriptors = descArr
		}
	}

	var tags []wafstore.Tag
	if tagsRaw := req.Parameters["Tags"]; tagsRaw != nil {
		tags = parseTags(tagsRaw)
	}

	ipSet, err := stores.ipSets.Create(id, name, "", "IPV4", nil)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, NewAPIError("AlreadyExistsException", fmt.Sprintf("IPSet already exists: %s", name), 400)
		}
		return nil, err
	}

	if len(descriptors) > 0 {
		ipSet.IPSetDescriptors = descriptors
		if err := stores.ipSets.Put(ipSet.ID, ipSet); err != nil {
			logs.Warn("failed to persist descriptors for IPSet", logs.String("id", ipSet.ID), logs.Err(err))
		}
	}

	if len(tags) > 0 && ipSet.ARN != "" {
		ipSet.Tags = tags
		if err := stores.ipSets.Put(ipSet.ID, ipSet); err != nil {
			logs.Warn("failed to persist tags for IPSet", logs.String("id", ipSet.ID), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
		"IPSet": map[string]interface{}{
			"IPSetId":          ipSet.ID,
			"Name":             ipSet.Name,
			"IPSetDescriptors": ipSet.IPSetDescriptors,
		},
	}, nil
}

// GetIPSet returns the details of the specified IP set.
func (s *WAFService) GetIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "IPSetId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "IPSetId is required", 400)
	}

	ipSet, err := stores.ipSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "IPSet not found", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"IPSet": map[string]interface{}{
			"IPSetId":          ipSet.ID,
			"Name":             ipSet.Name,
			"IPSetDescriptors": ipSet.IPSetDescriptors,
		},
	}, nil
}

// ListIPSets returns a paginated list of all IP sets.
func (s *WAFService) ListIPSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	maxItems := request.GetIntParam(req.Parameters, "Limit")
	if maxItems == 0 {
		maxItems = 100
	}

	result, err := stores.ipSets.List("", maxItems)
	if err != nil {
		return nil, err
	}

	ipSets := make([]interface{}, 0, len(result.IPSets))
	for _, ips := range result.IPSets {
		ipSets = append(ipSets, map[string]interface{}{
			"IPSetId": ips.ID,
			"Name":    ips.Name,
		})
	}

	return map[string]interface{}{
		"IPSets":     ipSets,
		"NextMarker": result.NextMarker,
	}, nil
}

// UpdateIPSet inserts or deletes IP address descriptors within the specified IP set.
func (s *WAFService) UpdateIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "IPSetId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "IPSetId is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	ipSet, err := stores.ipSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "IPSet not found", 400)
		}
		return nil, err
	}

	if updatesRaw := req.Parameters["Updates"]; updatesRaw != nil {
		if updates, ok := updatesRaw.([]interface{}); ok {
			descriptors := make([]interface{}, len(ipSet.IPSetDescriptors))
			copy(descriptors, ipSet.IPSetDescriptors)
			for _, u := range updates {
				if update, ok := u.(map[string]interface{}); ok {
					action := request.GetStringParam(update, "Action")
					desc := update["IPSetDescriptor"]
					if descMap, ok := desc.(map[string]interface{}); ok {
						if action == "INSERT" {
							descriptors = append(descriptors, descMap)
						} else if action == "DELETE" {
							descriptors = removeDescriptor(descriptors, descMap)
						}
					}
				}
			}
			ipSet.IPSetDescriptors = descriptors
		}
	}

	if err := stores.ipSets.Put(id, ipSet); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
	}, nil
}

func removeDescriptor(descriptors []interface{}, target map[string]interface{}) []interface{} {
	tType, _ := target["Type"].(string)
	tValue, _ := target["Value"].(string)
	result := make([]interface{}, 0, len(descriptors))
	for _, d := range descriptors {
		if dm, ok := d.(map[string]interface{}); ok {
			if dm["Type"] == tType && dm["Value"] == tValue {
				continue
			}
		}
		result = append(result, d)
	}
	return result
}

// DeleteIPSet permanently deletes the specified IP set.
func (s *WAFService) DeleteIPSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "IPSetId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "IPSetId is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	if err := stores.ipSets.Delete(id, ""); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "IPSet not found", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
	}, nil
}
