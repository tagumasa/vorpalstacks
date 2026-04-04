package wafv2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateWebACL creates a new web ACL with the specified default action, rules, and visibility configuration.
func (s *WAFv2Service) CreateWebACL(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, validationError("Name is required")
	}

	scope := request.GetStringParam(req.Parameters, "Scope")
	if scope == "" {
		scope = "REGIONAL"
	}

	description := request.GetStringParam(req.Parameters, "Description")
	capacity := int64(request.GetIntParam(req.Parameters, "Capacity"))
	if capacity == 0 {
		capacity = 1500
	}

	defaultAction := convertAction(request.GetMapParam(req.Parameters, "DefaultAction"))
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

	_, err = stores.webACLs.Get(id)
	if err == nil {
		return nil, newAPIError("WafV2AlreadyExistsException", fmt.Sprintf("WebACL already exists: %s", name), 400)
	}

	webACL, err := stores.webACLs.Create(id, name, description, scope, capacity, rules, defaultAction, visibilityConfig)
	if err != nil {
		return nil, err
	}

	if tags := parseTags(req.Parameters["Tags"]); len(tags) > 0 {
		webACL.Tags = tags
		if err := stores.webACLs.Put(webACL.ID, webACL); err != nil {
			logs.Warn("failed to persist tags for WebACL", logs.String("id", webACL.ID), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"Summary": map[string]interface{}{
			"Id":          webACL.ID,
			"Name":        webACL.Name,
			"ARN":         webACL.ARN,
			"Description": webACL.Description,
			"LockToken":   webACL.LockToken,
		},
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

// ListWebACLs returns a paginated list of all web ACLs.
func (s *WAFv2Service) ListWebACLs(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems := request.GetIntParam(req.Parameters, "Limit")
	if maxItems == 0 {
		maxItems = 100
	}

	nextMarker := request.GetStringParam(req.Parameters, "NextMarker")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.webACLs.List(nextMarker, maxItems)
	if err != nil {
		return nil, err
	}

	webACLs := make([]interface{}, 0, len(result.WebACLs))
	for _, wa := range result.WebACLs {
		webACLs = append(webACLs, map[string]interface{}{
			"Id":          wa.ID,
			"Name":        wa.Name,
			"ARN":         wa.ARN,
			"Description": wa.Description,
			"LockToken":   wa.LockToken,
		})
	}

	resp := map[string]interface{}{
		"WebACLs": webACLs,
	}
	if result.NextMarker != "" {
		resp["NextMarker"] = result.NextMarker
	}
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

	var daAction *wafstore.Action
	if defaultAction != nil {
		if a, ok := defaultAction.(*wafstore.Action); ok {
			daAction = a
		} else if m, ok := defaultAction.(map[string]interface{}); ok {
			daAction = mapToAction(m)
		}
	}

	updated, err := stores.webACLs.Update(id, lockToken, capacity, rules, daAction, visibilityConfig)
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

	if err := stores.webACLs.Delete(id, lockToken); err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("WebACL")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}
