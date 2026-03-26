package waf

import (
	"context"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// TagResource adds tags to a WAF resource.
func (s *WAFService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, NewValidationException("ResourceARN is required")
	}

	newTags := tagutil.ParseTags(req.Parameters, "Tags")
	if len(newTags) == 0 {
		return nil, NewValidationException("Tags are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceType := extractResourceTypeFromARN(resourceArn)
	switch resourceType {
	case "webacl":
		return tagWebACL(stores, resourceArn, newTags)
	case "rulegroup":
		return tagRuleGroup(stores, resourceArn, newTags)
	case "ipset":
		return tagIPSet(stores, resourceArn, newTags)
	case "regexpatternset":
		return tagRegexPatternSet(stores, resourceArn, newTags)
	case "rule":
		return tagRule(stores, resourceArn, newTags)
	}

	return nil, NewResourceNotFoundException("Resource not found: " + resourceArn)
}

// UntagResource removes tags from a WAF resource.
func (s *WAFService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, NewValidationException("ResourceARN is required")
	}

	tagKeysMap := tagutil.ParseTagKeys(req.Parameters, "TagKeys")
	if len(tagKeysMap) == 0 {
		return nil, NewValidationException("TagKeys are required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceType := extractResourceTypeFromARN(resourceArn)
	switch resourceType {
	case "webacl":
		return untagWebACL(stores, resourceArn, tagKeysMap)
	case "rulegroup":
		return untagRuleGroup(stores, resourceArn, tagKeysMap)
	case "ipset":
		return untagIPSet(stores, resourceArn, tagKeysMap)
	case "regexpatternset":
		return untagRegexPatternSet(stores, resourceArn, tagKeysMap)
	case "rule":
		return untagRule(stores, resourceArn, tagKeysMap)
	}

	return nil, NewResourceNotFoundException("Resource not found: " + resourceArn)
}

// ListTagsForResource returns the tags for a WAF resource.
func (s *WAFService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, NewValidationException("ResourceARN is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceType := extractResourceTypeFromARN(resourceArn)
	switch resourceType {
	case "webacl":
		return listTagsForWebACL(stores, resourceArn)
	case "rulegroup":
		return listTagsForRuleGroup(stores, resourceArn)
	case "ipset":
		return listTagsForIPSet(stores, resourceArn)
	case "regexpatternset":
		return listTagsForRegexPatternSet(stores, resourceArn)
	case "rule":
		return listTagsForRule(stores, resourceArn)
	}

	return nil, NewResourceNotFoundException("Resource not found: " + resourceArn)
}

func extractResourceTypeFromARN(arn string) string {
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return ""
	}
	resource := parts[5]
	subParts := strings.Split(resource, "/")
	if len(subParts) > 0 {
		return strings.ToLower(subParts[0])
	}
	return ""
}

func extractResourceIDFromARN(arn string) string {
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return ""
	}
	resource := parts[5]
	subParts := strings.Split(resource, "/")
	if len(subParts) > 1 {
		return subParts[1]
	}
	return ""
}

func tagWebACL(stores *wafStores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	webACL, err := stores.webACLs.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("WebACL not found")
	}
	webACL.Tags = tagutil.Apply(webACL.Tags, newTags)
	if err := stores.webACLs.Put(webACL.ID, webACL); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagWebACL(stores *wafStores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	webACL, err := stores.webACLs.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("WebACL not found")
	}
	webACL.Tags = tagutil.Remove(webACL.Tags, tagKeysMap)
	if err := stores.webACLs.Put(webACL.ID, webACL); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForWebACL(stores *wafStores, arn string) (interface{}, error) {
	webACL, err := stores.webACLs.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("WebACL not found")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(webACL.Tags),
		},
	}, nil
}

func tagRuleGroup(stores *wafStores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	ruleGroup, err := stores.ruleGroups.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("RuleGroup not found")
	}
	ruleGroup.Tags = tagutil.Apply(ruleGroup.Tags, newTags)
	if err := stores.ruleGroups.Put(ruleGroup.ID, ruleGroup); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagRuleGroup(stores *wafStores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	ruleGroup, err := stores.ruleGroups.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("RuleGroup not found")
	}
	ruleGroup.Tags = tagutil.Remove(ruleGroup.Tags, tagKeysMap)
	if err := stores.ruleGroups.Put(ruleGroup.ID, ruleGroup); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForRuleGroup(stores *wafStores, arn string) (interface{}, error) {
	ruleGroup, err := stores.ruleGroups.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("RuleGroup not found")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(ruleGroup.Tags),
		},
	}, nil
}

func tagIPSet(stores *wafStores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	ipSet, err := stores.ipSets.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("IPSet not found")
	}
	ipSet.Tags = tagutil.Apply(ipSet.Tags, newTags)
	if err := stores.ipSets.Put(ipSet.ID, ipSet); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagIPSet(stores *wafStores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	ipSet, err := stores.ipSets.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("IPSet not found")
	}
	ipSet.Tags = tagutil.Remove(ipSet.Tags, tagKeysMap)
	if err := stores.ipSets.Put(ipSet.ID, ipSet); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForIPSet(stores *wafStores, arn string) (interface{}, error) {
	ipSet, err := stores.ipSets.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("IPSet not found")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(ipSet.Tags),
		},
	}, nil
}

func tagRegexPatternSet(stores *wafStores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	regexSet, err := stores.regexPatternSets.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("RegexPatternSet not found")
	}
	regexSet.Tags = tagutil.Apply(regexSet.Tags, newTags)
	if err := stores.regexPatternSets.Put(regexSet.ID, regexSet); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagRegexPatternSet(stores *wafStores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	regexSet, err := stores.regexPatternSets.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("RegexPatternSet not found")
	}
	regexSet.Tags = tagutil.Remove(regexSet.Tags, tagKeysMap)
	if err := stores.regexPatternSets.Put(regexSet.ID, regexSet); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForRegexPatternSet(stores *wafStores, arn string) (interface{}, error) {
	regexSet, err := stores.regexPatternSets.GetByARN(arn)
	if err != nil {
		return nil, NewResourceNotFoundException("RegexPatternSet not found")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(regexSet.Tags),
		},
	}, nil
}

func tagRule(stores *wafStores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	id := extractResourceIDFromARN(arn)
	if id == "" {
		return nil, NewResourceNotFoundException("Rule not found")
	}
	rule, err := stores.ruleGroups.GetRule(id)
	if err != nil {
		return nil, NewResourceNotFoundException("Rule not found")
	}
	rule.Tags = tagutil.Apply(rule.Tags, newTags)
	if err := stores.ruleGroups.CreateRule(id, rule); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagRule(stores *wafStores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	id := extractResourceIDFromARN(arn)
	if id == "" {
		return nil, NewResourceNotFoundException("Rule not found")
	}
	rule, err := stores.ruleGroups.GetRule(id)
	if err != nil {
		return nil, NewResourceNotFoundException("Rule not found")
	}
	rule.Tags = tagutil.Remove(rule.Tags, tagKeysMap)
	if err := stores.ruleGroups.CreateRule(id, rule); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForRule(stores *wafStores, arn string) (interface{}, error) {
	id := extractResourceIDFromARN(arn)
	if id == "" {
		return nil, NewResourceNotFoundException("Rule not found")
	}
	rule, err := stores.ruleGroups.GetRule(id)
	if err != nil {
		return nil, NewResourceNotFoundException("Rule not found")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(rule.Tags),
		},
	}, nil
}
