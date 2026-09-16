package eventbridge

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// extractRuleInfoFromArn splits a rule ARN into its event bus and rule name
// parts.
func extractRuleInfoFromArn(arn string) (eventBusName, ruleName string) {
	_, _, _, _, resource := svcarn.SplitARN(arn)
	if resource == "" {
		return "", ""
	}
	parts := strings.Split(resource, "/")
	if len(parts) >= 3 && parts[0] == "rule" {
		eventBusName = parts[1]
		ruleName = parts[2]
	} else if len(parts) >= 2 {
		eventBusName = "default"
		ruleName = parts[len(parts)-1]
	}
	return eventBusName, ruleName
}

// TagResource adds tags to an EventBridge resource. The required-member
// rejections (ResourceARN, Tags) and the limit validation live in
// tagResourceCore.
func (s *EventsService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	newTags := tagutil.ParseTags(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.tagResourceCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "ResourceARN"), newTags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from an EventBridge resource. The
// required-member rejections (ResourceARN, TagKeys) live in
// untagResourceCore.
func (s *EventsService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tagKeysMap := tagutil.ParseTagKeys(req.Parameters, "TagKeys")
	if len(tagKeysMap) == 0 {
		tagKeysMap = tagutil.ParseTagKeys(req.Parameters, "tagKeys")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.untagResourceCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "ResourceARN"), tagKeysMap); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for an EventBridge resource. The
// ResourceARN requirement lives in listTagsForResourceCore's resolver.
func (s *EventsService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tagSlice, err := s.listTagsForResourceCore(ctx, store, request.GetParamLowerFirst(req.Parameters, "ResourceARN"))
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"Tags": tagListToMaps(tagSlice)}, nil
}
