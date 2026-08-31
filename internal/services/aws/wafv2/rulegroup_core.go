package wafv2

import (
	"fmt"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	storecommon "vorpalstacks/internal/store/aws/common"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// RuleGroupCreateInput is the transport-agnostic input for creating a rule
// group. VisibilityConfig arrives already converted because its wire
// conversion is pure, while Rules travel raw because the rule parse (which
// carries the statement validation) must run after the capacity and
// visibility-config validations in the original failure precedence.
type RuleGroupCreateInput struct {
	Name                 string
	Scope                string
	Description          string
	Capacity             int64
	VisibilityConfig     *wafstore.VisibilityConfig
	RulesRaw             interface{}
	CustomResponseBodies interface{}
	MonetizationConfig   interface{}
	Tags                 []types.Tag
}

// createRuleGroupCore is the single entry point for creating a rule group.
func (s *WAFv2Service) createRuleGroupCore(stores *wafv2Stores, in RuleGroupCreateInput) (*wafstore.RuleGroup, error) {
	if err := validateEntityName(in.Name); err != nil {
		return nil, err
	}
	if err := validateScope(in.Scope); err != nil {
		return nil, err
	}
	if err := validateEntityDescription(in.Description); err != nil {
		return nil, err
	}
	if in.Capacity < wafstore.MinRuleGroupCapacity {
		return nil, invalidParamError("Capacity is required and must be at least 1")
	}
	if in.Capacity > wafstore.MaxWebACLCapacity {
		return nil, limitsExceededError(in.Capacity)
	}
	if err := validateVisibilityConfig(in.VisibilityConfig); err != nil {
		return nil, err
	}
	rules, err := parseRules(in.RulesRaw)
	if err != nil {
		return nil, err
	}
	// Per the Smithy documentation for Capacity, WAF enforces the
	// declared limit whenever rules are added or modified — including
	// at creation. The capacity bound subsumes the global WCU quota
	// check because capacity itself was already validated against
	// MaxWebACLCapacity above.
	if consumed := calculateRulesCapacity(rules); consumed > in.Capacity {
		return nil, limitsExceededError(consumed)
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	ruleGroup := &wafstore.RuleGroup{
		ID:                   id,
		Name:                 in.Name,
		Description:          in.Description,
		Capacity:             in.Capacity,
		Rules:                rules,
		VisibilityConfig:     in.VisibilityConfig,
		Scope:                in.Scope,
		CustomResponseBodies: in.CustomResponseBodies,
		MonetizationConfig:   in.MonetizationConfig,
	}

	ruleGroup, err = stores.ruleGroups.Create(ruleGroup)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WAFDuplicateItemException", fmt.Sprintf("AWS WAF couldn't perform the operation because some resource in your request is a duplicate of an existing one: %s", in.Name), 400)
		}
		return nil, err
	}

	if len(in.Tags) > 0 && ruleGroup.ARN != "" {
		if err := stores.tags.TagFromSlice(ruleGroup.ARN, in.Tags); err != nil {
			logs.Warn("failed to persist tags for RuleGroup", logs.String("id", ruleGroup.ID), logs.Err(err))
		}
	}

	return ruleGroup, nil
}

// getRuleGroupCore is the single entry point for retrieving a rule group.
func (s *WAFv2Service) getRuleGroupCore(stores *wafv2Stores, id string) (*wafstore.RuleGroup, error) {
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

	return ruleGroup, nil
}

// RuleGroupListInput is the transport-agnostic input for listing rule
// groups.
type RuleGroupListInput struct {
	Scope      string
	Limit      int
	NextMarker string
}

// listRuleGroupsCore is the single entry point for listing rule groups.
func (s *WAFv2Service) listRuleGroupsCore(stores *wafv2Stores, in RuleGroupListInput) (*wafstore.RuleGroupListResult, error) {
	if err := validateScope(in.Scope); err != nil {
		return nil, err
	}

	return stores.ruleGroups.List(in.NextMarker, in.Limit, in.Scope)
}

// RuleGroupUpdateInput is the transport-agnostic input for updating a rule
// group. VisibilityConfig and Rules travel raw because their wire
// conversion must run after the pre-update fetch of the declared capacity,
// matching the original failure precedence.
type RuleGroupUpdateInput struct {
	Id                  string
	LockToken           string
	VisibilityConfigRaw interface{}
	RulesRaw            interface{}
}

// updateRuleGroupCore is the single entry point for updating a rule group.
func (s *WAFv2Service) updateRuleGroupCore(stores *wafv2Stores, in RuleGroupUpdateInput) (*wafstore.RuleGroup, error) {
	if in.Id == "" {
		return nil, invalidParamError("Id is required")
	}

	if in.LockToken == "" {
		return nil, invalidParamError("LockToken is required")
	}

	// UpdateRuleGroupRequest has no Capacity member in the Smithy
	// model: a rule group's capacity is fixed at creation, and WAF
	// enforces that declared limit whenever rules are added or
	// modified, so the existing rule group is fetched for its declared
	// capacity. VisibilityConfig is required on every call and Rules
	// omitted means an empty rule list.
	existing, err := stores.ruleGroups.Get(in.Id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RuleGroup")
		}
		return nil, err
	}

	vcRaw := in.VisibilityConfigRaw
	if vcRaw == nil {
		return nil, invalidParamError("VisibilityConfig is required")
	}
	vcMap, ok := vcRaw.(map[string]interface{})
	if !ok {
		return nil, invalidParamError("VisibilityConfig must be an object")
	}
	visibilityConfig := convertVisibilityConfig(vcMap)
	if err := validateVisibilityConfig(visibilityConfig); err != nil {
		return nil, err
	}

	var rules []*wafstore.Rule
	if in.RulesRaw != nil {
		parsed, pErr := parseRules(in.RulesRaw)
		if pErr != nil {
			return nil, pErr
		}
		rules = parsed
	}

	if consumed := calculateRulesCapacity(rules); consumed > existing.Capacity {
		return nil, limitsExceededError(consumed)
	}

	ruleGroup, err := stores.ruleGroups.Update(in.Id, in.LockToken, rules, visibilityConfig)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RuleGroup")
		}
		return nil, err
	}

	return ruleGroup, nil
}

// deleteRuleGroupCore is the single entry point for deleting a rule group.
func (s *WAFv2Service) deleteRuleGroupCore(stores *wafv2Stores, id, lockToken string) error {
	if id == "" {
		return invalidParamError("Id is required")
	}

	if lockToken == "" {
		return invalidParamError("LockToken is required")
	}

	// AWS rejects deletion of a rule group that is still referenced by
	// any web ACL rule (WAFAssociatedItemException). The check spans
	// the request region; cross-region references are impossible
	// because ARNs embed the creating region and resources are scoped
	// per region.
	existing, err := stores.ruleGroups.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return notFoundError("RuleGroup")
		}
		return err
	}
	if err := ensureRuleGroupNotReferenced(stores, existing.ARN); err != nil {
		return err
	}

	deleted, err := stores.ruleGroups.Delete(id, lockToken)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return notFoundError("RuleGroup")
		}
		if wafstore.IsLockTokenMismatch(err) {
			return lockTokenError()
		}
		return err
	}

	if deleted.ARN != "" {
		if err := stores.tags.Delete(deleted.ARN); err != nil {
			logs.Warn("failed to clean up tags for deleted RuleGroup", logs.String("id", id), logs.Err(err))
		}
	}

	return nil
}

// ensureRuleGroupNotReferenced scans every WebACL in the same region for
// a RuleGroupReferenceStatement whose ARN matches the given rule group.
// AWS rejects deletion of a rule group that any web ACL still uses,
// returning WAFAssociatedItemException.
func ensureRuleGroupNotReferenced(stores *wafv2Stores, ruleGroupArn string) error {
	if ruleGroupArn == "" {
		return nil
	}
	webACLs, err := storecommon.ListAll[wafstore.WebACL](stores.webACLs.BaseStore)
	if err != nil {
		return err
	}
	for _, acl := range webACLs {
		if acl == nil {
			continue
		}
		for _, rule := range acl.Rules {
			if rule == nil || rule.Statement == nil {
				continue
			}
			if ref := rule.Statement.RuleGroupReferenceStatement; ref != nil && ref.ARN == ruleGroupArn {
				return associatedItemError(fmt.Sprintf("RuleGroup %s is still referenced by WebACL %s.", ruleGroupArn, acl.ARN))
			}
		}
	}
	return nil
}
