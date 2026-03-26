// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// DeleteResourcePolicy deletes a resource policy from a DynamoDB table.
func (s *DynamoDBService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	_, err = store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.Tables().DeleteResourcePolicy(tableName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetResourcePolicy returns the resource policy for a DynamoDB table.
func (s *DynamoDBService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	_, err = store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	policy, err := store.Tables().GetResourcePolicy(tableName)
	if err != nil {
		return nil, err
	}

	if policy == "" {
		return nil, ErrPolicyNotFound
	}

	return map[string]interface{}{
		"Policy": policy,
	}, nil
}

// PutResourcePolicy creates or updates a resource policy for a DynamoDB table.
func (s *DynamoDBService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	policy, ok := req.Parameters["Policy"].(string)
	if !ok || policy == "" {
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
	_, err = store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.Tables().SetResourcePolicy(tableName, policy); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"RevisionId": "v1",
	}, nil
}
