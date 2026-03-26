package waf

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/waf"
	wafconnect "vorpalstacks/internal/pb/aws/waf/wafconnect"
)

// AdminHandler provides WAF (classic) service administration functionality.
// It implements the WAFServiceHandler interface for gRPC-Web communication.
type AdminHandler struct{}

var _ wafconnect.WAFServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new WAF AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListWebACLs lists web access control lists in WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListWebACLs(ctx context.Context, req *connect.Request[pb.ListWebACLsRequest]) (*connect.Response[pb.ListWebACLsResponse], error) {
	return connect.NewResponse(&pb.ListWebACLsResponse{
		Webacls: []*pb.WebACLSummary{},
	}), nil
}

// CreateByteMatchSet creates a byte match set for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateByteMatchSet(ctx context.Context, req *connect.Request[pb.CreateByteMatchSetRequest]) (*connect.Response[pb.CreateByteMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateByteMatchSet is not implemented"))
}

// CreateGeoMatchSet creates a geographic match set for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateGeoMatchSet(ctx context.Context, req *connect.Request[pb.CreateGeoMatchSetRequest]) (*connect.Response[pb.CreateGeoMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateGeoMatchSet is not implemented"))
}

// CreateIPSet creates an IP set for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateIPSet(ctx context.Context, req *connect.Request[pb.CreateIPSetRequest]) (*connect.Response[pb.CreateIPSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateIPSet is not implemented"))
}

// CreateRateBasedRule creates a rate-based rule for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRateBasedRule(ctx context.Context, req *connect.Request[pb.CreateRateBasedRuleRequest]) (*connect.Response[pb.CreateRateBasedRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateRateBasedRule is not implemented"))
}

// CreateRegexMatchSet creates a regex match set for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRegexMatchSet(ctx context.Context, req *connect.Request[pb.CreateRegexMatchSetRequest]) (*connect.Response[pb.CreateRegexMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateRegexMatchSet is not implemented"))
}

// CreateRegexPatternSet creates a regex pattern set for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRegexPatternSet(ctx context.Context, req *connect.Request[pb.CreateRegexPatternSetRequest]) (*connect.Response[pb.CreateRegexPatternSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateRegexPatternSet is not implemented"))
}

// CreateRule creates a rule for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRule(ctx context.Context, req *connect.Request[pb.CreateRuleRequest]) (*connect.Response[pb.CreateRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateRule is not implemented"))
}

// CreateRuleGroup creates a rule group for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRuleGroup(ctx context.Context, req *connect.Request[pb.CreateRuleGroupRequest]) (*connect.Response[pb.CreateRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateRuleGroup is not implemented"))
}

// CreateSizeConstraintSet creates a size constraint set for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateSizeConstraintSet(ctx context.Context, req *connect.Request[pb.CreateSizeConstraintSetRequest]) (*connect.Response[pb.CreateSizeConstraintSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateSizeConstraintSet is not implemented"))
}

// CreateSqlInjectionMatchSet creates an SQL injection match set for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateSqlInjectionMatchSet(ctx context.Context, req *connect.Request[pb.CreateSqlInjectionMatchSetRequest]) (*connect.Response[pb.CreateSqlInjectionMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateSqlInjectionMatchSet is not implemented"))
}

// CreateWebACL creates a web access control list for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateWebACL(ctx context.Context, req *connect.Request[pb.CreateWebACLRequest]) (*connect.Response[pb.CreateWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateWebACL is not implemented"))
}

// CreateWebACLMigrationStack creates a migration stack for a web ACL.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateWebACLMigrationStack(ctx context.Context, req *connect.Request[pb.CreateWebACLMigrationStackRequest]) (*connect.Response[pb.CreateWebACLMigrationStackResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateWebACLMigrationStack is not implemented"))
}

// CreateXssMatchSet creates an XSS match set for WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateXssMatchSet(ctx context.Context, req *connect.Request[pb.CreateXssMatchSetRequest]) (*connect.Response[pb.CreateXssMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.CreateXssMatchSet is not implemented"))
}

// DeleteByteMatchSet deletes a byte match set from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteByteMatchSet(ctx context.Context, req *connect.Request[pb.DeleteByteMatchSetRequest]) (*connect.Response[pb.DeleteByteMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteByteMatchSet is not implemented"))
}

// DeleteGeoMatchSet deletes a geographic match set from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteGeoMatchSet(ctx context.Context, req *connect.Request[pb.DeleteGeoMatchSetRequest]) (*connect.Response[pb.DeleteGeoMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteGeoMatchSet is not implemented"))
}

// DeleteIPSet deletes an IP set from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteIPSet(ctx context.Context, req *connect.Request[pb.DeleteIPSetRequest]) (*connect.Response[pb.DeleteIPSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteIPSet is not implemented"))
}

// DeleteLoggingConfiguration deletes a logging configuration from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteLoggingConfiguration(ctx context.Context, req *connect.Request[pb.DeleteLoggingConfigurationRequest]) (*connect.Response[pb.DeleteLoggingConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteLoggingConfiguration is not implemented"))
}

// DeletePermissionPolicy deletes a permission policy from a WAF resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeletePermissionPolicy(ctx context.Context, req *connect.Request[pb.DeletePermissionPolicyRequest]) (*connect.Response[pb.DeletePermissionPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeletePermissionPolicy is not implemented"))
}

// DeleteRateBasedRule deletes a rate-based rule from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRateBasedRule(ctx context.Context, req *connect.Request[pb.DeleteRateBasedRuleRequest]) (*connect.Response[pb.DeleteRateBasedRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteRateBasedRule is not implemented"))
}

// DeleteRegexMatchSet deletes a regex match set from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRegexMatchSet(ctx context.Context, req *connect.Request[pb.DeleteRegexMatchSetRequest]) (*connect.Response[pb.DeleteRegexMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteRegexMatchSet is not implemented"))
}

// DeleteRegexPatternSet deletes a regex pattern set from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRegexPatternSet(ctx context.Context, req *connect.Request[pb.DeleteRegexPatternSetRequest]) (*connect.Response[pb.DeleteRegexPatternSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteRegexPatternSet is not implemented"))
}

// DeleteRule deletes a rule from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRule(ctx context.Context, req *connect.Request[pb.DeleteRuleRequest]) (*connect.Response[pb.DeleteRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteRule is not implemented"))
}

// DeleteRuleGroup deletes a rule group from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRuleGroup(ctx context.Context, req *connect.Request[pb.DeleteRuleGroupRequest]) (*connect.Response[pb.DeleteRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteRuleGroup is not implemented"))
}

// DeleteSizeConstraintSet deletes a size constraint set from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteSizeConstraintSet(ctx context.Context, req *connect.Request[pb.DeleteSizeConstraintSetRequest]) (*connect.Response[pb.DeleteSizeConstraintSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteSizeConstraintSet is not implemented"))
}

// DeleteSqlInjectionMatchSet deletes an SQL injection match set from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteSqlInjectionMatchSet(ctx context.Context, req *connect.Request[pb.DeleteSqlInjectionMatchSetRequest]) (*connect.Response[pb.DeleteSqlInjectionMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteSqlInjectionMatchSet is not implemented"))
}

// DeleteWebACL deletes a web access control list from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteWebACL(ctx context.Context, req *connect.Request[pb.DeleteWebACLRequest]) (*connect.Response[pb.DeleteWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteWebACL is not implemented"))
}

// DeleteXssMatchSet deletes an XSS match set from WAF.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteXssMatchSet(ctx context.Context, req *connect.Request[pb.DeleteXssMatchSetRequest]) (*connect.Response[pb.DeleteXssMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.DeleteXssMatchSet is not implemented"))
}

// GetByteMatchSet returns details about a byte match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetByteMatchSet(ctx context.Context, req *connect.Request[pb.GetByteMatchSetRequest]) (*connect.Response[pb.GetByteMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetByteMatchSet is not implemented"))
}

// GetChangeToken returns the change token for making updates to WAF resources.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetChangeToken(ctx context.Context, req *connect.Request[pb.GetChangeTokenRequest]) (*connect.Response[pb.GetChangeTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetChangeToken is not implemented"))
}

// GetChangeTokenStatus returns the status of a change token.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetChangeTokenStatus(ctx context.Context, req *connect.Request[pb.GetChangeTokenStatusRequest]) (*connect.Response[pb.GetChangeTokenStatusResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetChangeTokenStatus is not implemented"))
}

// GetGeoMatchSet returns details about a geographic match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetGeoMatchSet(ctx context.Context, req *connect.Request[pb.GetGeoMatchSetRequest]) (*connect.Response[pb.GetGeoMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetGeoMatchSet is not implemented"))
}

// GetIPSet returns details about an IP set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetIPSet(ctx context.Context, req *connect.Request[pb.GetIPSetRequest]) (*connect.Response[pb.GetIPSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetIPSet is not implemented"))
}

// GetLoggingConfiguration returns details about a logging configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetLoggingConfiguration(ctx context.Context, req *connect.Request[pb.GetLoggingConfigurationRequest]) (*connect.Response[pb.GetLoggingConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetLoggingConfiguration is not implemented"))
}

// GetPermissionPolicy returns the permission policy for a WAF resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetPermissionPolicy(ctx context.Context, req *connect.Request[pb.GetPermissionPolicyRequest]) (*connect.Response[pb.GetPermissionPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetPermissionPolicy is not implemented"))
}

// GetRateBasedRule returns details about a rate-based rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRateBasedRule(ctx context.Context, req *connect.Request[pb.GetRateBasedRuleRequest]) (*connect.Response[pb.GetRateBasedRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetRateBasedRule is not implemented"))
}

// GetRateBasedRuleManagedKeys returns the managed keys for a rate-based rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRateBasedRuleManagedKeys(ctx context.Context, req *connect.Request[pb.GetRateBasedRuleManagedKeysRequest]) (*connect.Response[pb.GetRateBasedRuleManagedKeysResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetRateBasedRuleManagedKeys is not implemented"))
}

// GetRegexMatchSet returns details about a regex match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRegexMatchSet(ctx context.Context, req *connect.Request[pb.GetRegexMatchSetRequest]) (*connect.Response[pb.GetRegexMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetRegexMatchSet is not implemented"))
}

// GetRegexPatternSet returns details about a regex pattern set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRegexPatternSet(ctx context.Context, req *connect.Request[pb.GetRegexPatternSetRequest]) (*connect.Response[pb.GetRegexPatternSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetRegexPatternSet is not implemented"))
}

// GetRule returns details about a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRule(ctx context.Context, req *connect.Request[pb.GetRuleRequest]) (*connect.Response[pb.GetRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetRule is not implemented"))
}

// GetRuleGroup returns details about a rule group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRuleGroup(ctx context.Context, req *connect.Request[pb.GetRuleGroupRequest]) (*connect.Response[pb.GetRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetRuleGroup is not implemented"))
}

// GetSampledRequests returns sampled requests for a web ACL.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSampledRequests(ctx context.Context, req *connect.Request[pb.GetSampledRequestsRequest]) (*connect.Response[pb.GetSampledRequestsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetSampledRequests is not implemented"))
}

// GetSizeConstraintSet returns details about a size constraint set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSizeConstraintSet(ctx context.Context, req *connect.Request[pb.GetSizeConstraintSetRequest]) (*connect.Response[pb.GetSizeConstraintSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetSizeConstraintSet is not implemented"))
}

// GetSqlInjectionMatchSet returns details about an SQL injection match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSqlInjectionMatchSet(ctx context.Context, req *connect.Request[pb.GetSqlInjectionMatchSetRequest]) (*connect.Response[pb.GetSqlInjectionMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetSqlInjectionMatchSet is not implemented"))
}

// GetWebACL returns details about a web access control list.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetWebACL(ctx context.Context, req *connect.Request[pb.GetWebACLRequest]) (*connect.Response[pb.GetWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetWebACL is not implemented"))
}

// GetXssMatchSet returns details about an XSS match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetXssMatchSet(ctx context.Context, req *connect.Request[pb.GetXssMatchSetRequest]) (*connect.Response[pb.GetXssMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.GetXssMatchSet is not implemented"))
}

// ListActivatedRulesInRuleGroup returns rules that are activated in a rule group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListActivatedRulesInRuleGroup(ctx context.Context, req *connect.Request[pb.ListActivatedRulesInRuleGroupRequest]) (*connect.Response[pb.ListActivatedRulesInRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListActivatedRulesInRuleGroup is not implemented"))
}

// ListByteMatchSets returns all byte match sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListByteMatchSets(ctx context.Context, req *connect.Request[pb.ListByteMatchSetsRequest]) (*connect.Response[pb.ListByteMatchSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListByteMatchSets is not implemented"))
}

// ListGeoMatchSets returns all geographic match sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListGeoMatchSets(ctx context.Context, req *connect.Request[pb.ListGeoMatchSetsRequest]) (*connect.Response[pb.ListGeoMatchSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListGeoMatchSets is not implemented"))
}

// ListIPSets returns all IP sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListIPSets(ctx context.Context, req *connect.Request[pb.ListIPSetsRequest]) (*connect.Response[pb.ListIPSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListIPSets is not implemented"))
}

// ListLoggingConfigurations returns all logging configurations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListLoggingConfigurations(ctx context.Context, req *connect.Request[pb.ListLoggingConfigurationsRequest]) (*connect.Response[pb.ListLoggingConfigurationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListLoggingConfigurations is not implemented"))
}

// ListRateBasedRules returns all rate-based rules.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRateBasedRules(ctx context.Context, req *connect.Request[pb.ListRateBasedRulesRequest]) (*connect.Response[pb.ListRateBasedRulesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListRateBasedRules is not implemented"))
}

// ListRegexMatchSets returns all regex match sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRegexMatchSets(ctx context.Context, req *connect.Request[pb.ListRegexMatchSetsRequest]) (*connect.Response[pb.ListRegexMatchSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListRegexMatchSets is not implemented"))
}

// ListRegexPatternSets returns all regex pattern sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRegexPatternSets(ctx context.Context, req *connect.Request[pb.ListRegexPatternSetsRequest]) (*connect.Response[pb.ListRegexPatternSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListRegexPatternSets is not implemented"))
}

// ListRuleGroups returns all rule groups.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRuleGroups(ctx context.Context, req *connect.Request[pb.ListRuleGroupsRequest]) (*connect.Response[pb.ListRuleGroupsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListRuleGroups is not implemented"))
}

// ListRules returns all rules.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRules(ctx context.Context, req *connect.Request[pb.ListRulesRequest]) (*connect.Response[pb.ListRulesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListRules is not implemented"))
}

// ListSizeConstraintSets returns all size constraint sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSizeConstraintSets(ctx context.Context, req *connect.Request[pb.ListSizeConstraintSetsRequest]) (*connect.Response[pb.ListSizeConstraintSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListSizeConstraintSets is not implemented"))
}

// ListSqlInjectionMatchSets returns all SQL injection match sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSqlInjectionMatchSets(ctx context.Context, req *connect.Request[pb.ListSqlInjectionMatchSetsRequest]) (*connect.Response[pb.ListSqlInjectionMatchSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListSqlInjectionMatchSets is not implemented"))
}

// ListSubscribedRuleGroups returns rule groups that are subscribed to.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSubscribedRuleGroups(ctx context.Context, req *connect.Request[pb.ListSubscribedRuleGroupsRequest]) (*connect.Response[pb.ListSubscribedRuleGroupsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListSubscribedRuleGroups is not implemented"))
}

// ListTagsForResource returns tags for a WAF resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceRequest]) (*connect.Response[pb.ListTagsForResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListTagsForResource is not implemented"))
}

// ListXssMatchSets returns all XSS match sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListXssMatchSets(ctx context.Context, req *connect.Request[pb.ListXssMatchSetsRequest]) (*connect.Response[pb.ListXssMatchSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.ListXssMatchSets is not implemented"))
}

// PutLoggingConfiguration creates or updates a logging configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutLoggingConfiguration(ctx context.Context, req *connect.Request[pb.PutLoggingConfigurationRequest]) (*connect.Response[pb.PutLoggingConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.PutLoggingConfiguration is not implemented"))
}

// PutPermissionPolicy attaches a permission policy to a WAF resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutPermissionPolicy(ctx context.Context, req *connect.Request[pb.PutPermissionPolicyRequest]) (*connect.Response[pb.PutPermissionPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.PutPermissionPolicy is not implemented"))
}

// TagResource adds tags to a WAF resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[pb.TagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.TagResource is not implemented"))
}

// UntagResource removes tags from a WAF resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[pb.UntagResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UntagResource is not implemented"))
}

// UpdateByteMatchSet updates a byte match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateByteMatchSet(ctx context.Context, req *connect.Request[pb.UpdateByteMatchSetRequest]) (*connect.Response[pb.UpdateByteMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateByteMatchSet is not implemented"))
}

// UpdateGeoMatchSet updates a geographic match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateGeoMatchSet(ctx context.Context, req *connect.Request[pb.UpdateGeoMatchSetRequest]) (*connect.Response[pb.UpdateGeoMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateGeoMatchSet is not implemented"))
}

// UpdateIPSet updates an IP set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateIPSet(ctx context.Context, req *connect.Request[pb.UpdateIPSetRequest]) (*connect.Response[pb.UpdateIPSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateIPSet is not implemented"))
}

// UpdateRateBasedRule updates a rate-based rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRateBasedRule(ctx context.Context, req *connect.Request[pb.UpdateRateBasedRuleRequest]) (*connect.Response[pb.UpdateRateBasedRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateRateBasedRule is not implemented"))
}

// UpdateRegexMatchSet updates a regex match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRegexMatchSet(ctx context.Context, req *connect.Request[pb.UpdateRegexMatchSetRequest]) (*connect.Response[pb.UpdateRegexMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateRegexMatchSet is not implemented"))
}

// UpdateRegexPatternSet updates a regex pattern set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRegexPatternSet(ctx context.Context, req *connect.Request[pb.UpdateRegexPatternSetRequest]) (*connect.Response[pb.UpdateRegexPatternSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateRegexPatternSet is not implemented"))
}

// UpdateRule updates a rule.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRule(ctx context.Context, req *connect.Request[pb.UpdateRuleRequest]) (*connect.Response[pb.UpdateRuleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateRule is not implemented"))
}

// UpdateRuleGroup updates a rule group.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRuleGroup(ctx context.Context, req *connect.Request[pb.UpdateRuleGroupRequest]) (*connect.Response[pb.UpdateRuleGroupResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateRuleGroup is not implemented"))
}

// UpdateSizeConstraintSet updates a size constraint set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateSizeConstraintSet(ctx context.Context, req *connect.Request[pb.UpdateSizeConstraintSetRequest]) (*connect.Response[pb.UpdateSizeConstraintSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateSizeConstraintSet is not implemented"))
}

// UpdateSqlInjectionMatchSet updates an SQL injection match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateSqlInjectionMatchSet(ctx context.Context, req *connect.Request[pb.UpdateSqlInjectionMatchSetRequest]) (*connect.Response[pb.UpdateSqlInjectionMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateSqlInjectionMatchSet is not implemented"))
}

// UpdateWebACL updates a web access control list.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateWebACL(ctx context.Context, req *connect.Request[pb.UpdateWebACLRequest]) (*connect.Response[pb.UpdateWebACLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateWebACL is not implemented"))
}

// UpdateXssMatchSet updates an XSS match set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateXssMatchSet(ctx context.Context, req *connect.Request[pb.UpdateXssMatchSetRequest]) (*connect.Response[pb.UpdateXssMatchSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("waf.WAFService.UpdateXssMatchSet is not implemented"))
}
