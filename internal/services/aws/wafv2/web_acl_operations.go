package wafv2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateWebACL creates a new web ACL with the specified default action, rules, and visibility configuration.
func (s *WAFv2Service) CreateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
	description := request.GetStringParam(req.Parameters, "Description")
	defaultAction := convertAction(request.GetMapParam(req.Parameters, "DefaultAction"))
	rules, err := parseRules(req.Parameters["Rules"])
	if err != nil {
		return nil, err
	}
	visibilityConfig := convertVisibilityConfig(request.GetMapParam(req.Parameters, "VisibilityConfig"))

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	webACL, err := s.createWebACLCore(stores, CreateWebACLInput{
		Name:                   request.GetStringParam(req.Parameters, "Name"),
		Description:            description,
		Scope:                  scope,
		DefaultAction:          defaultAction,
		Rules:                  rules,
		VisibilityConfig:       visibilityConfig,
		CustomResponseBodies:   req.Parameters["CustomResponseBodies"],
		CaptchaConfig:          req.Parameters["CaptchaConfig"],
		ChallengeConfig:        req.Parameters["ChallengeConfig"],
		TokenDomains:           req.Parameters["TokenDomains"],
		AssociationConfig:      req.Parameters["AssociationConfig"],
		ApplicationConfig:      req.Parameters["ApplicationConfig"],
		MonetizationConfig:     req.Parameters["MonetizationConfig"],
		DataProtectionConfig:   req.Parameters["DataProtectionConfig"],
		OnSourceDDoSProtection: req.Parameters["OnSourceDDoSProtectionConfig"],
		Tags:                   tagutil.ParseTags(req.Parameters, "Tags"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Summary": buildWebACLSummary(webACL),
	}, nil
}

// GetWebACL retrieves the details of the specified web ACL.
func (s *WAFv2Service) GetWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	webACL, err := s.getWebACLCore(reqCtx, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
		return nil, err
	}

	webACLMap := map[string]interface{}{
		"Id":               webACL.ID,
		"Name":             webACL.Name,
		"ARN":              webACL.ARN,
		"Description":      webACL.Description,
		"DefaultAction":    convertActionToResponse(webACL.DefaultAction),
		"Rules":            convertRulesToResponse(webACL.Rules),
		"VisibilityConfig": convertVisibilityConfigToResponse(webACL.VisibilityConfig),
		"Capacity":         webACL.Capacity,
	}
	if webACL.LabelNamespace != "" {
		webACLMap["LabelNamespace"] = webACL.LabelNamespace
	}
	if webACL.CustomResponseBodies != nil {
		webACLMap["CustomResponseBodies"] = webACL.CustomResponseBodies
	}
	if webACL.CaptchaConfig != nil {
		webACLMap["CaptchaConfig"] = webACL.CaptchaConfig
	}
	if webACL.ChallengeConfig != nil {
		webACLMap["ChallengeConfig"] = webACL.ChallengeConfig
	}
	if webACL.TokenDomains != nil {
		webACLMap["TokenDomains"] = webACL.TokenDomains
	}
	if webACL.AssociationConfig != nil {
		webACLMap["AssociationConfig"] = webACL.AssociationConfig
	}
	if webACL.ApplicationConfig != nil {
		webACLMap["ApplicationConfig"] = webACL.ApplicationConfig
	}
	if webACL.MonetizationConfig != nil {
		webACLMap["MonetizationConfig"] = webACL.MonetizationConfig
	}
	if webACL.DataProtectionConfig != nil {
		webACLMap["DataProtectionConfig"] = webACL.DataProtectionConfig
	}
	if webACL.OnSourceDDoSProtection != nil {
		webACLMap["OnSourceDDoSProtectionConfig"] = webACL.OnSourceDDoSProtection
	}

	resp := map[string]interface{}{
		"WebACL":    webACLMap,
		"LockToken": webACL.LockToken,
	}

	region := reqCtx.GetRegion()
	if region != "" {
		resp["ApplicationIntegrationURL"] = fmt.Sprintf("https://console.aws.amazon.com/wafv2/%sv2/home#/webacl/%s/%s",
			region, webACL.ID, webACL.Scope)
	}

	return resp, nil
}

// ListWebACLs returns a paginated list of all web ACLs filtered by scope.
func (s *WAFv2Service) ListWebACLs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
	maxItems := pagination.GetMaxItems(req.Parameters, 100, "Limit")
	nextMarker := pagination.GetMarker(req.Parameters, "NextMarker")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listWebACLsCore(stores, ListWebACLsInput{
		Scope:      scope,
		Limit:      maxItems,
		NextMarker: nextMarker,
	})
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"WebACLs": buildWebACLSummaryList(result.WebACLs),
	}
	pagination.SetNextToken(resp, "NextMarker", result.NextMarker)
	return resp, nil
}

// UpdateWebACL updates the specified web ACL with new rules, default action, and visibility configuration, returning a new lock token.
func (s *WAFv2Service) UpdateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	nextLockToken, err := s.updateWebACLCore(reqCtx, WebACLUpdateInput{
		Id:                        request.GetStringParam(req.Parameters, "Id"),
		LockToken:                 request.GetStringParam(req.Parameters, "LockToken"),
		DefaultActionRaw:          req.Parameters["DefaultAction"],
		VisibilityConfigRaw:       req.Parameters["VisibilityConfig"],
		RulesRaw:                  req.Parameters["Rules"],
		Description:               request.GetStringParam(req.Parameters, "Description"),
		CustomResponseBodies:      req.Parameters["CustomResponseBodies"],
		CaptchaConfig:             req.Parameters["CaptchaConfig"],
		ChallengeConfig:           req.Parameters["ChallengeConfig"],
		TokenDomains:              req.Parameters["TokenDomains"],
		AssociationConfig:         req.Parameters["AssociationConfig"],
		ApplicationConfig:         req.Parameters["ApplicationConfig"],
		MonetizationConfig:        req.Parameters["MonetizationConfig"],
		DataProtectionConfig:      req.Parameters["DataProtectionConfig"],
		OnSourceDDoSProtectionCfg: req.Parameters["OnSourceDDoSProtectionConfig"],
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": nextLockToken,
	}, nil
}

// DeleteWebACL permanently deletes the specified web ACL.
func (s *WAFv2Service) DeleteWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	lockToken := request.GetStringParam(req.Parameters, "LockToken")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	assocStores, err := s.allAssociationStores(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := s.deleteWebACLCore(stores, assocStores, id, lockToken); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
