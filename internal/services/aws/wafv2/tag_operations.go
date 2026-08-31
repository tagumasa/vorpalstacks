package wafv2

import (
	"context"
	"errors"
	"strings"

	"vorpalstacks/internal/common/request"
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
