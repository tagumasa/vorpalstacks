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

// CreateWebACL creates a new web ACL with the specified default action, rules, and visibility configuration.
func (s *WAFv2Service) CreateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, validationError("Name is required")
	}

	scope := request.GetStringParam(req.Parameters, "Scope")
	if err := validateScope(scope); err != nil {
		return nil, err
	}

	description := request.GetStringParam(req.Parameters, "Description")
	capacity := int64(request.GetIntParam(req.Parameters, "Capacity"))
	if capacity == 0 {
		capacity = 1500
	}

	defaultAction := convertAction(request.GetMapParam(req.Parameters, "DefaultAction"))
	if err := validateDefaultAction(defaultAction); err != nil {
		return nil, err
	}
	rules := convertRules(req.Parameters["Rules"])
	visibilityConfig := convertVisibilityConfig(request.GetMapParam(req.Parameters, "VisibilityConfig"))

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	webACL, err := stores.webACLs.Create(id, name, description, scope, capacity, rules, defaultAction, visibilityConfig)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WAFDuplicateItemException", fmt.Sprintf("AWS WAF couldn't perform the operation because some resource in your request is a duplicate of an existing one: %s", name), 400)
		}
		return nil, err
	}

	if tags := tagutil.ParseTags(req.Parameters, "Tags"); len(tags) > 0 {
		if err := stores.tags.TagFromSlice(webACL.ARN, tags); err != nil {
			logs.Warn("failed to persist tags for WebACL", logs.String("id", webACL.ID), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"Summary": buildWebACLSummary(webACL),
	}, nil
}

// GetWebACL retrieves the details of the specified web ACL.
func (s *WAFv2Service) GetWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, validationError("Id is required")
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

	resp := map[string]interface{}{
		"WebACL": map[string]interface{}{
			"Id":               webACL.ID,
			"Name":             webACL.Name,
			"ARN":              webACL.ARN,
			"Description":      webACL.Description,
			"DefaultAction":    convertActionToResponse(webACL.DefaultAction),
			"Rules":            convertRulesToResponse(webACL.Rules),
			"VisibilityConfig": convertVisibilityConfigToResponse(webACL.VisibilityConfig),
			"Capacity":         webACL.Capacity,
			"Scope":            webACL.Scope,
		},
		"LockToken": webACL.LockToken,
	}

	return resp, nil
}

// ListWebACLs returns a paginated list of all web ACLs filtered by scope.
func (s *WAFv2Service) ListWebACLs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
	if err := validateScope(scope); err != nil {
		return nil, err
	}
	maxItems := pagination.GetMaxItems(req.Parameters, 100, "Limit")
	nextMarker := pagination.GetMarker(req.Parameters, "NextMarker")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.webACLs.List(nextMarker, maxItems, scope)
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
		return nil, validationError("Id is required")
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, validationError("LockToken is required")
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
		rules = convertRules(rulesRaw)
	}

	daAction := convertAction(nil)
	if defaultAction != nil {
		if a, ok := defaultAction.(*wafstore.Action); ok {
			daAction = a
		} else if m, ok := defaultAction.(map[string]interface{}); ok {
			daAction = convertAction(m)
		}
	}

	updated, err := stores.webACLs.Update(id, lockToken, capacity, rules, daAction, visibilityConfig, request.GetStringParam(req.Parameters, "Description"))
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
	if id == "" {
		return nil, validationError("Id is required")
	}

	lockToken := request.GetStringParam(req.Parameters, "LockToken")
	if lockToken == "" {
		return nil, validationError("LockToken is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	deleted, err := stores.webACLs.Delete(id, lockToken)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		return nil, err
	}

	_ = stores.tags.Delete(deleted.ARN)

	return response.EmptyResponse(), nil
}
