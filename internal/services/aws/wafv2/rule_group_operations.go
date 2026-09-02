package wafv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateRuleGroup creates a new rule group with the specified rules and visibility configuration.
func (s *WAFv2Service) CreateRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	ruleGroup, err := s.createRuleGroupCore(stores, RuleGroupCreateInput{
		Name:                 request.GetStringParam(req.Parameters, "Name"),
		Scope:                request.GetStringParam(req.Parameters, "Scope"),
		Description:          request.GetStringParam(req.Parameters, "Description"),
		Capacity:             int64(request.GetIntParam(req.Parameters, "Capacity")),
		VisibilityConfig:     convertVisibilityConfig(request.GetMapParam(req.Parameters, "VisibilityConfig")),
		RulesRaw:             req.Parameters["Rules"],
		CustomResponseBodies: req.Parameters["CustomResponseBodies"],
		MonetizationConfig:   req.Parameters["MonetizationConfig"],
		Tags:                 tagutil.ParseTags(req.Parameters, "Tags"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Summary": buildRuleGroupSummary(ruleGroup),
	}, nil
}

// GetRuleGroup retrieves the details of the specified rule group.
func (s *WAFv2Service) GetRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	ruleGroup, err := s.getRuleGroupCore(stores, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
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
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	result, err := s.listRuleGroupsCore(stores, RuleGroupListInput{
		Scope:      request.GetStringParam(req.Parameters, "Scope"),
		Limit:      pagination.GetMaxItems(req.Parameters, 100, "Limit"),
		NextMarker: pagination.GetMarker(req.Parameters, "NextMarker"),
	})
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
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	ruleGroup, err := s.updateRuleGroupCore(stores, RuleGroupUpdateInput{
		Id:                  request.GetStringParam(req.Parameters, "Id"),
		LockToken:           request.GetStringParam(req.Parameters, "LockToken"),
		VisibilityConfigRaw: req.Parameters["VisibilityConfig"],
		RulesRaw:            req.Parameters["Rules"],
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": ruleGroup.LockToken,
	}, nil
}

// DeleteRuleGroup permanently deletes the specified rule group.
func (s *WAFv2Service) DeleteRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	if err := s.deleteRuleGroupCore(stores, request.GetStringParam(req.Parameters, "Id"), request.GetStringParam(req.Parameters, "LockToken")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
