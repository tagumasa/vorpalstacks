package wafv2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
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
		return nil, invalidParamError("Name is required")
	}

	scope := request.GetStringParam(req.Parameters, "Scope")
	if err := validateScope(scope); err != nil {
		return nil, err
	}

	description := request.GetStringParam(req.Parameters, "Description")
	capacity := int64(request.GetIntParam(req.Parameters, "Capacity"))
	if capacity <= 0 {
		return nil, invalidParamError("Capacity is required and must be greater than 0")
	}

	visibilityConfig := convertVisibilityConfig(request.GetMapParam(req.Parameters, "VisibilityConfig"))
	rules, err := parseRules(req.Parameters["Rules"])
	if err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	ruleGroup := &wafstore.RuleGroup{
		ID:                   id,
		Name:                 name,
		Description:          description,
		Capacity:             capacity,
		Rules:                rules,
		VisibilityConfig:     visibilityConfig,
		Scope:                scope,
		CustomResponseBodies: req.Parameters["CustomResponseBodies"],
		MonetizationConfig:   req.Parameters["MonetizationConfig"],
	}

	ruleGroup, err = stores.ruleGroups.Create(ruleGroup)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WAFDuplicateItemException", fmt.Sprintf("AWS WAF couldn't perform the operation because some resource in your request is a duplicate of an existing one: %s", name), 400)
		}
		return nil, err
	}

	if tags := tagutil.ParseTags(req.Parameters, "Tags"); len(tags) > 0 && ruleGroup.ARN != "" {
		if err := stores.tags.TagFromSlice(ruleGroup.ARN, tags); err != nil {
			logs.Warn("failed to persist tags for RuleGroup", logs.String("id", ruleGroup.ID), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"Summary": buildRuleGroupSummary(ruleGroup),
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
		return nil, invalidParamError("Id is required")
	}

	ruleGroup, err := stores.ruleGroups.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RuleGroup")
		}
		return nil, err
	}

	rulesResp := convertRulesToResponse(ruleGroup.Rules)

	rgMap := map[string]interface{}{
		"Id":               ruleGroup.ID,
		"Name":             ruleGroup.Name,
		"ARN":              ruleGroup.ARN,
		"Capacity":         ruleGroup.Capacity,
		"Description":      ruleGroup.Description,
		"Rules":            rulesResp,
		"VisibilityConfig": convertVisibilityConfigToResponse(ruleGroup.VisibilityConfig),
	}
	if ruleGroup.LabelNamespace != "" {
		rgMap["LabelNamespace"] = ruleGroup.LabelNamespace
	}
	if ruleGroup.CustomResponseBodies != nil {
		rgMap["CustomResponseBodies"] = ruleGroup.CustomResponseBodies
	}
	if ruleGroup.AvailableLabels != nil {
		rgMap["AvailableLabels"] = ruleGroup.AvailableLabels
	}
	if ruleGroup.ConsumedLabels != nil {
		rgMap["ConsumedLabels"] = ruleGroup.ConsumedLabels
	}
	if ruleGroup.MonetizationConfig != nil {
		rgMap["MonetizationConfig"] = ruleGroup.MonetizationConfig
	}

	return map[string]interface{}{
		"RuleGroup": rgMap,
		"LockToken": ruleGroup.LockToken,
	}, nil
}

// ListRuleGroups returns a paginated list of all rule groups.
func (s *WAFv2Service) ListRuleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	scope := request.GetStringParam(req.Parameters, "Scope")
	if err := validateScope(scope); err != nil {
		return nil, err
	}
	maxItems := pagination.GetMaxItems(req.Parameters, 100, "Limit")
	nextMarker := pagination.GetMarker(req.Parameters, "NextMarker")

	result, err := stores.ruleGroups.List(nextMarker, maxItems, scope)
	if err != nil {
		return nil, err
	}

	ruleGroups := make([]interface{}, 0, len(result.RuleGroups))
	for _, rg := range result.RuleGroups {
		ruleGroups = append(ruleGroups, buildRuleGroupSummary(rg))
	}

	resp := map[string]interface{}{
		"RuleGroups": ruleGroups,
	}
	pagination.SetNextToken(resp, "NextMarker", result.NextMarker)
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
		return nil, invalidParamError("Id is required")
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, invalidParamError("LockToken is required")
	}

	ruleGroup, err := stores.ruleGroups.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RuleGroup")
		}
		return nil, err
	}

	capacity := ruleGroup.Capacity
	if c := int64(request.GetIntParam(req.Parameters, "Capacity")); c > 0 {
		capacity = c
	}
	visibilityConfig := ruleGroup.VisibilityConfig
	if vcRaw := req.Parameters["VisibilityConfig"]; vcRaw != nil {
		if vc, ok := vcRaw.(map[string]interface{}); ok {
			visibilityConfig = convertVisibilityConfig(vc)
		}
	}
	var rules []*wafstore.Rule
	if rulesRaw := req.Parameters["Rules"]; rulesRaw != nil {
		parsed, pErr := parseRules(rulesRaw)
		if pErr != nil {
			return nil, pErr
		}
		rules = parsed
	}

	ruleGroup, err = stores.ruleGroups.Update(id, lockToken, capacity, rules, visibilityConfig)
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
		return nil, invalidParamError("Id is required")
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, invalidParamError("LockToken is required")
	}

	deleted, err := stores.ruleGroups.Delete(id, lockToken)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RuleGroup")
		}
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		return nil, err
	}

	_ = stores.tags.Delete(deleted.ARN)

	return response.EmptyResponse(), nil
}
