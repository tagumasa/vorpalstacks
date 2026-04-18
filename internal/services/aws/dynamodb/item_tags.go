package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TagResource adds tags to the specified DynamoDB resource.
func (s *DynamoDBService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	tableName := svcarn.ParseTableARN(resourceArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	newTags := tagutil.ParseTags(req.Parameters, "Tags")
	mergedTags := tagutil.Apply(table.Tags, newTags)

	if err := store.Tables().SetTags(tableName, mergedTags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from the specified DynamoDB resource.
func (s *DynamoDBService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	tableName := svcarn.ParseTableARN(resourceArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	tagKeys := tagutil.ParseTagKeys(req.Parameters, "TagKeys")
	remainingTags := tagutil.Remove(table.Tags, tagKeys)

	if err := store.Tables().SetTags(tableName, remainingTags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists all tags associated with the specified DynamoDB resource.
func (s *DynamoDBService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	tableName := svcarn.ParseTableARN(resourceArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"Tags": tagutil.ConvertToMapSlice(table.Tags),
	}, nil
}
