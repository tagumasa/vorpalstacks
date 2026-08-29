package dynamodb

import (
	"context"
	"sort"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
)

func dynamodbMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError:
		return ErrInvalidParameter
	case *tagutil.MissingTagsError:
		return ErrInvalidParameter
	case *tagutil.MissingTagKeysError:
		return ErrInvalidParameter
	}
	return err
}

// TagResource adds or overwrites tags on a DynamoDB table.
func (s *DynamoDBService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tags := tagutil.ParseTags(req.Parameters, "Tags")
	if tagutil.HasDuplicateKeys(tags) {
		return nil, ErrInvalidParameter
	}
	for _, tag := range tags {
		if !validateTagKey(tag.Key) {
			return nil, ErrInvalidParameter
		}
		if !validateTagValue(tag.Value) {
			return nil, ErrInvalidParameter
		}
	}
	return tagutil.HandleTag(ctx, req, dynamodbTagConfig(store))
}

// UntagResource removes the specified tags from a DynamoDB table.
func (s *DynamoDBService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, dynamodbTagConfig(store))
}

// ListTagsForResource lists tags assigned to a DynamoDB table with pagination.
// AWS paginates tags; the marker is the tag key of the last returned entry.
func (s *DynamoDBService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cfg := dynamodbTagConfig(store)
	rawKey := tagutil.GetResourceKey(req.Parameters, cfg.Param)
	if rawKey == "" {
		return nil, dynamodbMapError(&tagutil.MissingResourceError{Param: "ResourceArn"})
	}
	resourceKey := rawKey
	if cfg.ResourceKey != nil {
		resourceKey = cfg.ResourceKey(rawKey)
	}
	if cfg.ValidateResource != nil {
		if err := cfg.ValidateResource(ctx, resourceKey); err != nil {
			return nil, dynamodbMapError(err)
		}
	}

	allTags, err := cfg.ListFunc(ctx, resourceKey)
	if err != nil {
		return nil, err
	}

	// Sort tags for deterministic pagination.
	sort.Slice(allTags, func(i, j int) bool {
		return allTags[i].Key < allTags[j].Key
	})

	// Apply NextToken pagination (marker = tag key).
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	startIdx := 0
	if nextToken != "" {
		for i, t := range allTags {
			if t.Key == nextToken {
				startIdx = i + 1
				break
			}
		}
	}

	pageSize := listTagsForResourceDefaultPageSize
	remaining := allTags[startIdx:]
	hasMore := len(remaining) > pageSize
	if hasMore {
		remaining = remaining[:pageSize]
	}

	resp := map[string]interface{}{
		"Tags": tagutil.ConvertToMapSlice(remaining),
	}
	if hasMore && len(remaining) > 0 {
		resp["NextToken"] = remaining[len(remaining)-1].Key
	}

	return resp, nil
}
