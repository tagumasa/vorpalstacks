package wafv2

import (
	"context"
	"fmt"
	"log"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateRuleGroup creates a new rule group with the specified rules and visibility configuration.
func (s *WAFv2Service) CreateRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, validationError("Name is required")
	}

	scope := request.GetStringParam(req.Parameters, "Scope")
	_ = scope

	description := request.GetStringParam(req.Parameters, "Description")
	capacity := int64(request.GetIntParam(req.Parameters, "Capacity"))
	if capacity == 0 {
		capacity = 10
	}

	visibilityConfig := convertVisibilityConfig(request.GetMapParam(req.Parameters, "VisibilityConfig"))
	rules := convertRules(req.Parameters["Rules"])

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	ruleGroup, err := stores.ruleGroups.Create(id, name, description, capacity, rules, visibilityConfig)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WafV2AlreadyExistsException", fmt.Sprintf("RuleGroup already exists: %s", name), 400)
		}
		return nil, err
	}

	if tags := parseTags(req.Parameters["Tags"]); len(tags) > 0 && ruleGroup.ARN != "" {
		ruleGroup.Tags = tags
		if err := stores.ruleGroups.Put(ruleGroup.ID, ruleGroup); err != nil {
			log.Printf("WARNING: failed to persist tags for RuleGroup %s: %v", ruleGroup.ID, err)
		}
	}

	return map[string]interface{}{
		"Summary": map[string]interface{}{
			"Id":          ruleGroup.ID,
			"Name":        ruleGroup.Name,
			"ARN":         ruleGroup.ARN,
			"Description": ruleGroup.Description,
			"LockToken":   ruleGroup.LockToken,
		},
	}, nil
}

// GetRuleGroup retrieves the details of the specified rule group.
func (s *WAFv2Service) GetRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, validationError("Id is required")
	}

	ruleGroup, err := stores.ruleGroups.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RuleGroup")
		}
		return nil, err
	}

	rulesResp := convertRulesToResponse(ruleGroup.Rules)

	return map[string]interface{}{
		"RuleGroup": map[string]interface{}{
			"Id":               ruleGroup.ID,
			"Name":             ruleGroup.Name,
			"ARN":              ruleGroup.ARN,
			"Capacity":         ruleGroup.Capacity,
			"Description":      ruleGroup.Description,
			"Rules":            rulesResp,
			"VisibilityConfig": convertVisibilityConfigToResponse(ruleGroup.VisibilityConfig),
		},
		"LockToken": ruleGroup.LockToken,
	}, nil
}

// ListRuleGroups returns a paginated list of all rule groups.
func (s *WAFv2Service) ListRuleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	maxItems := request.GetIntParam(req.Parameters, "Limit")
	if maxItems == 0 {
		maxItems = 100
	}

	result, err := stores.ruleGroups.List("", maxItems)
	if err != nil {
		return nil, err
	}

	ruleGroups := make([]interface{}, 0, len(result.RuleGroups))
	for _, rg := range result.RuleGroups {
		ruleGroups = append(ruleGroups, map[string]interface{}{
			"Id":          rg.ID,
			"Name":        rg.Name,
			"ARN":         rg.ARN,
			"Description": rg.Description,
			"LockToken":   rg.LockToken,
		})
	}

	resp := map[string]interface{}{
		"RuleGroups": ruleGroups,
	}
	if result.NextMarker != "" {
		resp["NextMarker"] = result.NextMarker
	}
	return resp, nil
}

// UpdateRuleGroup updates the specified rule group with new rules and configuration, returning a new lock token.
func (s *WAFv2Service) UpdateRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	capacity := int64(request.GetIntParam(req.Parameters, "Capacity"))
	visibilityConfig := convertVisibilityConfig(request.GetMapParam(req.Parameters, "VisibilityConfig"))
	rules := convertRules(req.Parameters["Rules"])

	ruleGroup, err := stores.ruleGroups.Update(id, lockToken, capacity, rules, visibilityConfig)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RuleGroup")
		}
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": ruleGroup.LockToken,
	}, nil
}

// DeleteRuleGroup permanently deletes the specified rule group.
func (s *WAFv2Service) DeleteRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	err = stores.ruleGroups.Delete(id, lockToken)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RuleGroup")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}
