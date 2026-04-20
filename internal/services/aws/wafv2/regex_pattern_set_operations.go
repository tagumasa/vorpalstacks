package wafv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateRegexPatternSet creates a new regex pattern set containing the specified regular expressions.
func (s *WAFv2Service) CreateRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	var regularPatterns []string
	if rpRaw := req.Parameters["RegularExpressionList"]; rpRaw != nil {
		if arr, ok := rpRaw.([]interface{}); ok {
			for _, r := range arr {
				if m, ok := r.(map[string]interface{}); ok {
					if rs, ok := m["RegexString"].(string); ok {
						regularPatterns = append(regularPatterns, rs)
					}
				}
			}
		}
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	rps, err := stores.regexPatternSets.Create(id, name, description, regularPatterns)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WafV2AlreadyExistsException", "RegexPatternSet already exists", 400)
		}
		return nil, err
	}

	if tags := tagutil.ParseTags(req.Parameters, "Tags"); len(tags) > 0 {
		if err := stores.tags.TagFromSlice(rps.ARN, tags); err != nil {
			logs.Warn("failed to persist tags for RegexPatternSet", logs.String("id", rps.ID), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"Summary": map[string]interface{}{
			"Id":          rps.ID,
			"Name":        rps.Name,
			"ARN":         rps.ARN,
			"Description": rps.Description,
			"LockToken":   rps.LockToken,
		},
	}, nil
}

// GetRegexPatternSet retrieves the details of the specified regex pattern set.
func (s *WAFv2Service) GetRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "Id")
	if id == "" {
		return nil, validationError("Id is required")
	}

	rps, err := stores.regexPatternSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RegexPatternSet")
		}
		return nil, err
	}

	regexList := make([]interface{}, 0, len(rps.RegularPatterns))
	for _, p := range rps.RegularPatterns {
		regexList = append(regexList, map[string]interface{}{
			"RegexString": p,
		})
	}

	return map[string]interface{}{
		"RegexPatternSet": map[string]interface{}{
			"Id":                    rps.ID,
			"Name":                  rps.Name,
			"ARN":                   rps.ARN,
			"RegularExpressionList": regexList,
			"Description":           rps.Description,
		},
		"LockToken": rps.LockToken,
	}, nil
}

// ListRegexPatternSets returns a paginated list of all regex pattern sets.
func (s *WAFv2Service) ListRegexPatternSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	maxItems := pagination.GetMaxItems(req.Parameters, 100, "Limit")
	nextMarker := pagination.GetMarker(req.Parameters, "NextMarker")

	result, err := stores.regexPatternSets.List(nextMarker, maxItems)
	if err != nil {
		return nil, err
	}

	sets := make([]interface{}, 0, len(result.RegexPatternSets))
	for _, rps := range result.RegexPatternSets {
		sets = append(sets, map[string]interface{}{
			"Id":          rps.ID,
			"Name":        rps.Name,
			"ARN":         rps.ARN,
			"Description": rps.Description,
			"LockToken":   rps.LockToken,
		})
	}

	resp := map[string]interface{}{
		"RegexPatternSets": sets,
	}
	pagination.SetNextToken(resp, "NextMarker", result.NextMarker)
	return resp, nil
}

// UpdateRegexPatternSet updates the specified regex pattern set with new regular expressions, returning a new lock token.
func (s *WAFv2Service) UpdateRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	var regularPatterns []string
	if rpRaw := req.Parameters["RegularExpressionList"]; rpRaw != nil {
		if arr, ok := rpRaw.([]interface{}); ok {
			for _, r := range arr {
				if m, ok := r.(map[string]interface{}); ok {
					if rs, ok := m["RegexString"].(string); ok {
						regularPatterns = append(regularPatterns, rs)
					}
				}
			}
		}
	}

	rps, err := stores.regexPatternSets.Update(id, lockToken, regularPatterns)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RegexPatternSet")
		}
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": rps.LockToken,
	}, nil
}

// DeleteRegexPatternSet permanently deletes the specified regex pattern set.
func (s *WAFv2Service) DeleteRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	rps, err := stores.regexPatternSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RegexPatternSet")
		}
		return nil, err
	}

	if err := stores.regexPatternSets.Delete(id, lockToken); err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		return nil, err
	}

	stores.tags.Delete(rps.ARN)

	return map[string]interface{}{}, nil
}
