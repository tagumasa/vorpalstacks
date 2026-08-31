package wafv2

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// wafv2TagConfig builds the shared tag-handler configuration bound to the
// given store group (relocated verbatim from tag_operations.go so the
// store-touching closures live on the Core layer).
func wafv2TagConfig(stores *wafv2Stores) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardARNConfig,
		ParseTagKeys: func(params map[string]interface{}) []string {
			tagKeysMap := tagutil.ParseTagKeys(params, "TagKeys")
			tagKeys := make([]string, 0, len(tagKeysMap))
			for k := range tagKeysMap {
				tagKeys = append(tagKeys, k)
			}
			return tagKeys
		},
		TagFunc: func(_ context.Context, resourceKey string, tag []tagutil.Tag) error {
			return stores.tags.Tag(resourceKey, tagutil.ToMap(tag))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return stores.tags.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return stores.tags.ListAsSlice(resourceKey)
		},
		FormatResponse: func(tag []tagutil.Tag, rawResourceKey string) (interface{}, error) {
			return map[string]interface{}{
				"TagInfoForResource": map[string]interface{}{
					"ResourceARN": rawResourceKey,
					"TagList":     tagutil.ToResponse(tag),
				},
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: wafv2MapError,
		ValidateResource: func(_ context.Context, resourceArn string) error {
			return validateWAFv2Resource(stores, resourceArn)
		},
		ValidateTagsFunc: tagutil.ValidateTags,
	}
}

// validateWAFv2Resource verifies that the tagged ARN resolves to an
// existing WAFv2 resource of a supported type (relocated verbatim from
// tag_operations.go).
func validateWAFv2Resource(stores *wafv2Stores, resourceArn string) error {
	resourceType := extractResourceTypeFromARN(resourceArn)
	if resourceType == "" {
		return invalidParamError(fmt.Sprintf("Unable to detect resource type from ARN: %s", resourceArn))
	}

	switch resourceType {
	case "webacl":
		_, err := stores.webACLs.GetByARN(resourceArn)
		if err != nil {
			return notFoundError("WebACL")
		}
	case "rulegroup":
		_, err := stores.ruleGroups.GetByARN(resourceArn)
		if err != nil {
			return notFoundError("RuleGroup")
		}
	case "ipset":
		_, err := stores.ipSets.GetByARN(resourceArn)
		if err != nil {
			return notFoundError("IPSet")
		}
	case "regexpatternset":
		_, err := stores.regexPatternSets.GetByARN(resourceArn)
		if err != nil {
			return notFoundError("RegexPatternSet")
		}
	default:
		return invalidParamError(fmt.Sprintf("Unsupported WAFv2 resource type: %s", resourceType))
	}

	return nil
}
