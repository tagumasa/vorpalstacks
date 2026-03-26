package waf

import (
	"context"
	"fmt"
	"log"

	"vorpalstacks/internal/services/aws/common/request"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// CreateRegexPatternSet creates a new regex pattern set containing regular expression patterns.
func (s *WAFService) CreateRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidParameterValue", "Name is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	var regexPatternStrings []string
	if rsRaw := req.Parameters["RegexPatternStrings"]; rsRaw != nil {
		if arr, ok := rsRaw.([]interface{}); ok {
			for _, r := range arr {
				if s, ok := r.(string); ok {
					regexPatternStrings = append(regexPatternStrings, s)
				}
			}
		}
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}
	rps, err := stores.regexPatternSets.Create(id, name, "", regexPatternStrings)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, NewAPIError("AlreadyExistsException", fmt.Sprintf("RegexPatternSet already exists: %s", name), 400)
		}
		return nil, err
	}

	rps.RegexPatternStrings = regexPatternStrings
	if err := stores.regexPatternSets.Put(rps.ID, rps); err != nil {
		log.Printf("WARNING: failed to persist patterns for RegexPatternSet %s: %v", rps.ID, err)
	}

	if tags := parseTags(req.Parameters["Tags"]); len(tags) > 0 && rps.ARN != "" {
		rps.Tags = tags
		if err := stores.regexPatternSets.Put(rps.ID, rps); err != nil {
			log.Printf("WARNING: failed to persist tags for RegexPatternSet %s: %v", rps.ID, err)
		}
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
		"RegexPatternSet": map[string]interface{}{
			"RegexPatternSetId":   rps.ID,
			"Name":                rps.Name,
			"RegexPatternStrings": rps.RegexPatternStrings,
		},
	}, nil
}

// GetRegexPatternSet returns the details of the specified regex pattern set.
func (s *WAFService) GetRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "RegexPatternSetId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "RegexPatternSetId is required", 400)
	}

	rps, err := stores.regexPatternSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "RegexPatternSet not found", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"RegexPatternSet": map[string]interface{}{
			"RegexPatternSetId":   rps.ID,
			"Name":                rps.Name,
			"RegexPatternStrings": rps.RegexPatternStrings,
		},
	}, nil
}

// ListRegexPatternSets returns a paginated list of all regex pattern sets.
func (s *WAFService) ListRegexPatternSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	maxItems := request.GetIntParam(req.Parameters, "Limit")
	if maxItems == 0 {
		maxItems = 100
	}

	result, err := stores.regexPatternSets.List("", maxItems)
	if err != nil {
		return nil, err
	}

	sets := make([]interface{}, 0, len(result.RegexPatternSets))
	for _, rps := range result.RegexPatternSets {
		sets = append(sets, map[string]interface{}{
			"RegexPatternSetId": rps.ID,
			"Name":              rps.Name,
		})
	}

	return map[string]interface{}{
		"RegexPatternSets": sets,
		"NextMarker":       result.NextMarker,
	}, nil
}

// UpdateRegexPatternSet inserts or deletes regex pattern strings within the specified regex pattern set.
func (s *WAFService) UpdateRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "RegexPatternSetId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "RegexPatternSetId is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	rps, err := stores.regexPatternSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "RegexPatternSet not found", 400)
		}
		return nil, err
	}

	if updatesRaw := req.Parameters["Updates"]; updatesRaw != nil {
		if updates, ok := updatesRaw.([]interface{}); ok {
			patterns := make([]string, len(rps.RegexPatternStrings))
			copy(patterns, rps.RegexPatternStrings)
			for _, u := range updates {
				if update, ok := u.(map[string]interface{}); ok {
					action := request.GetStringParam(update, "Action")
					if action == "INSERT" {
						if rs, ok := update["RegexPatternString"].(string); ok {
							patterns = append(patterns, rs)
						}
					} else if action == "DELETE" {
						if rs, ok := update["RegexPatternString"].(string); ok {
							patterns = removeString(patterns, rs)
						}
					}
				}
			}
			rps.RegexPatternStrings = patterns
		}
	}

	if err := stores.regexPatternSets.Put(id, rps); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
	}, nil
}

func removeString(slice []string, target string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != target {
			result = append(result, s)
		}
	}
	return result
}

// DeleteRegexPatternSet permanently deletes the specified regex pattern set.
func (s *WAFService) DeleteRegexPatternSet(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	id := request.GetStringParam(req.Parameters, "RegexPatternSetId")
	if id == "" {
		id = request.GetStringParam(req.Parameters, "Id")
	}
	if id == "" {
		return nil, NewAPIError("InvalidParameterValue", "RegexPatternSetId is required", 400)
	}

	changeToken := request.GetStringParam(req.Parameters, "ChangeToken")

	if err := stores.regexPatternSets.Delete(id, ""); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, NewAPIError("WAFNonexistentItemException", "RegexPatternSet not found", 400)
		}
		return nil, err
	}

	return map[string]interface{}{
		"ChangeToken": changeToken,
	}, nil
}
