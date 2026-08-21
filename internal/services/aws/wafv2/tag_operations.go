package wafv2

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func wafv2MapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError:
		return invalidParamError("ResourceARN is required")
	case *tagutil.MissingTagsError:
		return invalidParamError("Tags are required")
	case *tagutil.MissingTagKeysError:
		return invalidParamError("TagKeys are required")
	}
	if errors.Is(err, tagutil.ErrReservedTagKey) || errors.Is(err, tagutil.ErrTooManyTags) || errors.Is(err, tagutil.ErrInvalidTagKey) || errors.Is(err, tagutil.ErrInvalidTagValue) {
		return invalidParamError(err.Error())
	}
	return err
}

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

// TagResource adds or overwrites tags on a WAFv2 resource.
func (s *WAFv2Service) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, wafv2TagConfig(stores))
}

// UntagResource removes the specified tags from a WAFv2 resource.
func (s *WAFv2Service) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, wafv2TagConfig(stores))
}

// ListTagsForResource lists all tags assigned to a WAFv2 resource.
func (s *WAFv2Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, wafv2TagConfig(stores))
}

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

// extractResourceTypeFromARN reads the WAFv2 resource type from the ARN
// resource field. WAFv2 resource paths are either scope-prefixed
// (regional|cloudfront/<type>/<name>/<id>) or type-direct
// (<type>/<name>/<id>).
func extractResourceTypeFromARN(arn string) string {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	subParts := strings.Split(resource, "/")
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
