package wafv2

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateRegexPatternSet creates a new regex pattern set containing the specified regular expressions.
func (s *WAFv2Service) CreateRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	rps, err := s.createRegexPatternSetCore(stores, RegexPatternSetCreateInput{
		Name:                  request.GetStringParam(req.Parameters, "Name"),
		Scope:                 request.GetStringParam(req.Parameters, "Scope"),
		Description:           request.GetStringParam(req.Parameters, "Description"),
		RegularExpressionList: req.Parameters["RegularExpressionList"],
		Tags:                  tagutil.ParseTags(req.Parameters, "Tags"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Summary": buildRegexPatternSetSummary(rps),
	}, nil
}

// GetRegexPatternSet retrieves the details of the specified regex pattern set.
func (s *WAFv2Service) GetRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	rps, err := s.getRegexPatternSetCore(stores, request.GetStringParam(req.Parameters, "Id"))
	if err != nil {
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
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	result, err := s.listRegexPatternSetsCore(stores, RegexPatternSetListInput{
		Scope:      request.GetStringParam(req.Parameters, "Scope"),
		Limit:      pagination.GetMaxItems(req.Parameters, 100, "Limit"),
		NextMarker: pagination.GetMarker(req.Parameters, "NextMarker"),
	})
	if err != nil {
		return nil, err
	}

	sets := make([]interface{}, 0, len(result.RegexPatternSets))
	for _, rps := range result.RegexPatternSets {
		sets = append(sets, buildRegexPatternSetSummary(rps))
	}

	resp := map[string]interface{}{
		"RegexPatternSets": sets,
	}
	pagination.SetNextToken(resp, "NextMarker", result.NextMarker)
	return resp, nil
}

// UpdateRegexPatternSet updates the specified regex pattern set with new regular expressions, returning a new lock token.
func (s *WAFv2Service) UpdateRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	rps, err := s.updateRegexPatternSetCore(stores, RegexPatternSetUpdateInput{
		Id:                    request.GetStringParam(req.Parameters, "Id"),
		LockToken:             request.GetStringParam(req.Parameters, "LockToken"),
		RegularExpressionList: req.Parameters["RegularExpressionList"],
		Description:           request.GetStringParam(req.Parameters, "Description"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"NextLockToken": rps.LockToken,
	}, nil
}

// DeleteRegexPatternSet permanently deletes the specified regex pattern set.
func (s *WAFv2Service) DeleteRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.storeForScope(reqCtx, request.GetStringParam(req.Parameters, "Scope"))
	if err != nil {
		return nil, err
	}

	if err := s.deleteRegexPatternSetCore(stores, request.GetStringParam(req.Parameters, "Id"), request.GetStringParam(req.Parameters, "LockToken")); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
