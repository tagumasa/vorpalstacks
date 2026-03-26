package waf

import (
	"context"
	"fmt"
	"log"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateRuleGroup creates a new rule group.
func (s *WAFService) CreateRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidParameterValue", "Name is required", 400)
	}

	description := request.GetStringParam(req.Parameters, "Description")
	capacity := int64(request.GetIntParam(req.Parameters, "Capacity"))
	if capacity == 0 {
		capacity = 10
	}

	visibilityConfig := convertPBToStoreVisibilityConfig(request.GetMapParam(req.Parameters, "VisibilityConfig"))
	rules := convertPBToStoreRules(req.Parameters["Rules"])

	id, err := generateID()
	if err != nil {
		return nil, err
	}
	ruleGroup, err := stores.ruleGroups.Create(id, name, description, capacity, rules, visibilityConfig)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, NewAPIError("AlreadyExistsException", fmt.Sprintf("RuleGroup already exists: %s", name), 400)
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

// GetRuleGroup returns the details of a rule group.
func (s *WAFService) GetRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "Id is required", 400)
	}

	ruleGroup, err := stores.ruleGroups.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "RuleGroup not found", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"RuleGroup": map[string]interface{}{
			"Id":               ruleGroup.ID,
			"Name":             ruleGroup.Name,
			"ARN":              ruleGroup.ARN,
			"Capacity":         ruleGroup.Capacity,
			"Description":      ruleGroup.Description,
			"Rules":            convertStoreRulesToPB(ruleGroup.Rules),
			"VisibilityConfig": convertStoreVisibilityConfigToPB(ruleGroup.VisibilityConfig),
		},
		"LockToken": ruleGroup.LockToken,
	}, nil
}

// ListRuleGroups returns a list of rule groups.
func (s *WAFService) ListRuleGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
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

	return map[string]interface{}{
		"RuleGroups": ruleGroups,
		"NextMarker": result.NextMarker,
	}, nil
}

// UpdateRuleGroup updates an existing rule group.
func (s *WAFService) UpdateRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "Id is required", 400)
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, NewAPIError("InvalidParameterValue", "LockToken is required", 400)
	}

	capacity := int64(request.GetIntParam(req.Parameters, "Capacity"))
	visibilityConfig := convertPBToStoreVisibilityConfig(request.GetMapParam(req.Parameters, "VisibilityConfig"))
	rules := convertPBToStoreRules(req.Parameters["Rules"])

	ruleGroup, err := stores.ruleGroups.Update(id, lockToken, capacity, rules, visibilityConfig)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, NewAPIError("WAFInvalidLockTokenException", "Lock token mismatch", 400)
		}
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "RuleGroup not found", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": ruleGroup.LockToken,
	}, nil
}

// DeleteRuleGroup deletes a rule group.
func (s *WAFService) DeleteRuleGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "Id is required", 400)
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, NewAPIError("InvalidParameterValue", "LockToken is required", 400)
	}

	err = stores.ruleGroups.Delete(id, lockToken)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, NewAPIError("WAFInvalidLockTokenException", "Lock token mismatch", 400)
		}
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "RuleGroup not found", 400)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}
