package wafv2

import (
	"fmt"
	"regexp"

	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// parseRegularExpressionList parses the raw RegularExpressionList member
// and compiles every RegexString it carries. RegularExpressionList is a
// required member of CreateRegexPatternSet and UpdateRegexPatternSet: an
// omitted member and entries that are not pattern objects are rejected.
// An object without a RegexString member keeps its AWS semantics — the
// member is not required on the Regex shape, so the entry is accepted as
// empty and contributes no pattern.
func parseRegularExpressionList(raw interface{}) ([]string, error) {
	if raw == nil {
		return nil, invalidParamError("RegularExpressionList is a required member")
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return nil, invalidParamError("RegularExpressionList must be a list of pattern objects")
	}
	var patterns []string
	for _, r := range arr {
		m, ok := r.(map[string]interface{})
		if !ok {
			return nil, invalidParamError("RegularExpressionList entries must be objects with a RegexString member")
		}
		if rs, ok := m["RegexString"].(string); ok {
			if _, err := regexp.Compile(rs); err != nil {
				return nil, invalidParamError(fmt.Sprintf("Invalid regex pattern: %s", rs))
			}
			patterns = append(patterns, rs)
		}
	}
	return patterns, nil
}

// RegexPatternSetCreateInput is the transport-agnostic input for creating a
// regex pattern set.
type RegexPatternSetCreateInput struct {
	Name                  string
	Scope                 string
	Description           string
	RegularExpressionList interface{}
	Tags                  []types.Tag
}

// createRegexPatternSetCore is the single entry point for creating a regex
// pattern set. The validation ladder order matches the original handler:
// name, scope, description, patterns.
func (s *WAFv2Service) createRegexPatternSetCore(stores *wafv2Stores, in RegexPatternSetCreateInput) (*wafstore.RegexPatternSet, error) {
	if err := validateEntityName(in.Name); err != nil {
		return nil, err
	}
	if err := validateScope(in.Scope); err != nil {
		return nil, err
	}
	if err := validateEntityDescription(in.Description); err != nil {
		return nil, err
	}
	regularPatterns, err := parseRegularExpressionList(in.RegularExpressionList)
	if err != nil {
		return nil, err
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	rps, err := stores.regexPatternSets.Create(id, in.Name, in.Description, regularPatterns, in.Scope)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, newAPIError("WAFDuplicateItemException", "AWS WAF couldn't perform the operation because some resource in your request is a duplicate of an existing one", 400)
		}
		return nil, err
	}

	if len(in.Tags) > 0 {
		if err := stores.tags.TagFromSlice(rps.ARN, in.Tags); err != nil {
			logs.Warn("failed to persist tags for RegexPatternSet", logs.String("id", rps.ID), logs.Err(err))
		}
	}

	return rps, nil
}

// getRegexPatternSetCore is the single entry point for retrieving a regex
// pattern set.
func (s *WAFv2Service) getRegexPatternSetCore(stores *wafv2Stores, id string) (*wafstore.RegexPatternSet, error) {
	if id == "" {
		return nil, invalidParamError("Id is required")
	}

	rps, err := stores.regexPatternSets.Get(id)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RegexPatternSet")
		}
		return nil, err
	}

	return rps, nil
}

// RegexPatternSetListInput is the transport-agnostic input for listing
// regex pattern sets.
type RegexPatternSetListInput struct {
	Scope      string
	Limit      int
	NextMarker string
}

// listRegexPatternSetsCore is the single entry point for listing regex
// pattern sets.
func (s *WAFv2Service) listRegexPatternSetsCore(stores *wafv2Stores, in RegexPatternSetListInput) (*wafstore.RegexPatternSetListResult, error) {
	if err := validateScope(in.Scope); err != nil {
		return nil, err
	}

	return stores.regexPatternSets.List(in.NextMarker, in.Limit, in.Scope)
}

// RegexPatternSetUpdateInput is the transport-agnostic input for updating
// a regex pattern set.
type RegexPatternSetUpdateInput struct {
	Id                    string
	LockToken             string
	RegularExpressionList interface{}
	Description           string
}

// updateRegexPatternSetCore is the single entry point for updating a regex
// pattern set.
func (s *WAFv2Service) updateRegexPatternSetCore(stores *wafv2Stores, in RegexPatternSetUpdateInput) (*wafstore.RegexPatternSet, error) {
	if in.Id == "" {
		return nil, invalidParamError("Id is required")
	}

	if in.LockToken == "" {
		return nil, invalidParamError("LockToken is required")
	}

	regularPatterns, err := parseRegularExpressionList(in.RegularExpressionList)
	if err != nil {
		return nil, err
	}

	rps, err := stores.regexPatternSets.Update(in.Id, in.LockToken, regularPatterns, in.Description)
	if err != nil {
		if wafstore.IsLockTokenMismatch(err) {
			return nil, lockTokenError()
		}
		if wafstore.IsNotFound(err) {
			return nil, notFoundError("RegexPatternSet")
		}
		return nil, err
	}

	return rps, nil
}

// deleteRegexPatternSetCore is the single entry point for deleting a regex
// pattern set, including the tag cleanup on the deleted ARN.
func (s *WAFv2Service) deleteRegexPatternSetCore(stores *wafv2Stores, id, lockToken string) error {
	if id == "" {
		return invalidParamError("Id is required")
	}

	if lockToken == "" {
		return invalidParamError("LockToken is required")
	}

	deleted, err := stores.regexPatternSets.Delete(id, lockToken)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return notFoundError("RegexPatternSet")
		}
		if wafstore.IsLockTokenMismatch(err) {
			return lockTokenError()
		}
		return err
	}

	if deleted.ARN != "" {
		if err := stores.tags.Delete(deleted.ARN); err != nil {
			logs.Warn("failed to clean up tags for deleted RegexPatternSet", logs.String("id", id), logs.Err(err))
		}
	}

	return nil
}
