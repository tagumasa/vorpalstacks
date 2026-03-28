package wafv2

import (
	"context"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// TagResource adds or overwrites tags on the specified WAFv2 resource.
func (s *WAFv2Service) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, validationError("ResourceARN is required")
	}

	newTags := tagutil.ParseTags(req.Parameters, "Tags")
	if len(newTags) == 0 {
		return nil, validationError("Tags are required")
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
	}

	return nil, notFoundError(resourceArn)
}

// UntagResource removes the specified tags from the WAFv2 resource.
func (s *WAFv2Service) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, validationError("ResourceARN is required")
	}

	tagKeysMap := tagutil.ParseTagKeys(req.Parameters, "TagKeys")
	if len(tagKeysMap) == 0 {
		return nil, validationError("TagKeys are required")
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
	}

	return nil, notFoundError(resourceArn)
}

// ListTagsForResource retrieves all tags associated with the specified WAFv2 resource.
func (s *WAFv2Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceARN")
	if resourceArn == "" {
		return nil, validationError("ResourceARN is required")
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
	}

	return nil, notFoundError(resourceArn)
}

func extractResourceTypeFromARN(arn string) string {
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return ""
	}
	resource := parts[5]
	subParts := strings.Split(resource, "/")
	if len(subParts) == 0 {
		return ""
	}
	first := strings.ToLower(subParts[0])
	switch first {
	case "ipset", "webacl", "rulegroup", "regexpatternset":
		return first
	case "regional", "cloudfront":
		if len(subParts) > 1 {
			return strings.ToLower(subParts[1])
		}
	}
	return ""
}

func tagWebACL(stores *wafv2Stores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	webACL, err := stores.webACLs.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("WebACL")
	}
	webACL.Tags = tagutil.Apply(webACL.Tags, newTags)
	if err := stores.webACLs.Put(webACL.ID, webACL); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagWebACL(stores *wafv2Stores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	webACL, err := stores.webACLs.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("WebACL")
	}
	webACL.Tags = tagutil.Remove(webACL.Tags, tagKeysMap)
	if err := stores.webACLs.Put(webACL.ID, webACL); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForWebACL(stores *wafv2Stores, arn string) (interface{}, error) {
	webACL, err := stores.webACLs.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("WebACL")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(webACL.Tags),
		},
	}, nil
}

func tagRuleGroup(stores *wafv2Stores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	ruleGroup, err := stores.ruleGroups.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("RuleGroup")
	}
	ruleGroup.Tags = tagutil.Apply(ruleGroup.Tags, newTags)
	if err := stores.ruleGroups.Put(ruleGroup.ID, ruleGroup); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagRuleGroup(stores *wafv2Stores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	ruleGroup, err := stores.ruleGroups.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("RuleGroup")
	}
	ruleGroup.Tags = tagutil.Remove(ruleGroup.Tags, tagKeysMap)
	if err := stores.ruleGroups.Put(ruleGroup.ID, ruleGroup); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForRuleGroup(stores *wafv2Stores, arn string) (interface{}, error) {
	ruleGroup, err := stores.ruleGroups.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("RuleGroup")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(ruleGroup.Tags),
		},
	}, nil
}

func tagIPSet(stores *wafv2Stores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	ipSet, err := stores.ipSets.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("IPSet")
	}
	ipSet.Tags = tagutil.Apply(ipSet.Tags, newTags)
	if err := stores.ipSets.Put(ipSet.ID, ipSet); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagIPSet(stores *wafv2Stores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	ipSet, err := stores.ipSets.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("IPSet")
	}
	ipSet.Tags = tagutil.Remove(ipSet.Tags, tagKeysMap)
	if err := stores.ipSets.Put(ipSet.ID, ipSet); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForIPSet(stores *wafv2Stores, arn string) (interface{}, error) {
	ipSet, err := stores.ipSets.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("IPSet")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(ipSet.Tags),
		},
	}, nil
}

func tagRegexPatternSet(stores *wafv2Stores, arn string, newTags []wafstore.Tag) (interface{}, error) {
	regexSet, err := stores.regexPatternSets.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("RegexPatternSet")
	}
	regexSet.Tags = tagutil.Apply(regexSet.Tags, newTags)
	if err := stores.regexPatternSets.Put(regexSet.ID, regexSet); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagRegexPatternSet(stores *wafv2Stores, arn string, tagKeysMap map[string]bool) (interface{}, error) {
	regexSet, err := stores.regexPatternSets.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("RegexPatternSet")
	}
	regexSet.Tags = tagutil.Remove(regexSet.Tags, tagKeysMap)
	if err := stores.regexPatternSets.Put(regexSet.ID, regexSet); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listTagsForRegexPatternSet(stores *wafv2Stores, arn string) (interface{}, error) {
	regexSet, err := stores.regexPatternSets.GetByARN(arn)
	if err != nil {
		return nil, notFoundError("RegexPatternSet")
	}
	return map[string]interface{}{
		"TagInfoForResource": map[string]interface{}{
			"ResourceARN": arn,
			"TagList":     tagutil.ToResponse(regexSet.Tags),
		},
	}, nil
}
