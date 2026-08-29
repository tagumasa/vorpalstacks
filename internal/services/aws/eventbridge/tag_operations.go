package eventbridge

import (
	"context"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
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

// TagResource adds tags to an EventBridge resource.
func (s *EventsService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if err := validateResourceArnParam(resourceArn); err != nil {
		return nil, err
	}

	newTags := tagutil.ParseTags(req.Parameters, "Tags")
	if len(newTags) == 0 {
		return nil, awserrors.NewValidationException("Tags are required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.tagResourceCore(ctx, store, resourceArn, newTags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from an EventBridge resource.
func (s *EventsService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if err := validateResourceArnParam(resourceArn); err != nil {
		return nil, err
	}

	tagKeysMap := tagutil.ParseTagKeys(req.Parameters, "TagKeys")
	if len(tagKeysMap) == 0 {
		tagKeysMap = tagutil.ParseTagKeys(req.Parameters, "tagKeys")
	}
	if len(tagKeysMap) == 0 {
		return nil, awserrors.NewValidationException("TagKeys are required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.untagResourceCore(ctx, store, resourceArn, tagKeysMap); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists tags for an EventBridge resource.
func (s *EventsService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceARN")
	if err := validateResourceArnParam(resourceArn); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tagSlice, err := s.listTagsForResourceCore(ctx, store, resourceArn)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"Tags": tagListToMaps(tagSlice)}, nil
}
