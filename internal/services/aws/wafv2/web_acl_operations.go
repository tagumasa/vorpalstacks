package wafv2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	wafstore "vorpalstacks/internal/store/aws/waf"
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
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, invalidParamError("Id is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	webACL, err := stores.webACLs.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
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
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, invalidParamError("Id is required")
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, invalidParamError("LockToken is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	webACL, err := stores.webACLs.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	capacity := webACL.Capacity
	defaultAction := webACL.DefaultAction
	visibilityConfig := webACL.VisibilityConfig

	if c := int64(request.GetIntParam(req.Parameters, "Capacity")); c > 0 {
		capacity = c
	}
	if daRaw := req.Parameters["DefaultAction"]; daRaw != nil {
		if da, ok := daRaw.(map[string]interface{}); ok {
			defaultAction = convertAction(da)
		}
	}
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

	daAction := convertAction(nil)
	if defaultAction != nil {
		if a, ok := defaultAction.(*wafstore.Action); ok {
			daAction = a
		} else if m, ok := defaultAction.(map[string]interface{}); ok {
			daAction = convertAction(m)
		}
	}

	if v := req.Parameters["TokenDomains"]; v != nil {
		if err := validateTokenDomains(v); err != nil {
			return nil, err
		}
	}
	if v := req.Parameters["CustomResponseBodies"]; v != nil {
		if err := validateCustomResponseBodies(v); err != nil {
			return nil, err
		}
	}

	updated, err := stores.webACLs.Update(id, lockToken, capacity, rules, daAction, visibilityConfig, request.GetStringParam(req.Parameters, "Description"), func(webACL *wafstore.WebACL) {
		if v := req.Parameters["CustomResponseBodies"]; v != nil {
			webACL.CustomResponseBodies = v
		}
		if v := req.Parameters["CaptchaConfig"]; v != nil {
			webACL.CaptchaConfig = v
		}
		if v := req.Parameters["ChallengeConfig"]; v != nil {
			webACL.ChallengeConfig = v
		}
		if v := req.Parameters["TokenDomains"]; v != nil {
			webACL.TokenDomains = v
		}
		if v := req.Parameters["AssociationConfig"]; v != nil {
			webACL.AssociationConfig = v
		}
		if v := req.Parameters["ApplicationConfig"]; v != nil {
			webACL.ApplicationConfig = v
		}
		if v := req.Parameters["MonetizationConfig"]; v != nil {
			webACL.MonetizationConfig = v
		}
		if v := req.Parameters["DataProtectionConfig"]; v != nil {
			webACL.DataProtectionConfig = v
		}
		if v := req.Parameters["OnSourceDDoSProtectionConfig"]; v != nil {
			webACL.OnSourceDDoSProtection = v
		}
	})
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": updated.LockToken,
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

	if _, err := s.deleteWebACLCore(stores, id, lockToken); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
