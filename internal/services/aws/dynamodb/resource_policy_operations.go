package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// DeleteResourcePolicy deletes a resource policy from a DynamoDB table.
func (s *DynamoDBService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	expectedRevisionId := request.GetStringParam(req.Parameters, "ExpectedRevisionId")
	if !validatePolicyRevisionId(expectedRevisionId) {
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
	if _, err := store.Tables().Get(tableName); err != nil {
		return nil, ErrResourceNotFound
	}

	if err := s.deleteResourcePolicyCore(store, tableName, expectedRevisionId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetResourcePolicy returns the resource policy for a DynamoDB table.
func (s *DynamoDBService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
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
	if _, err := store.Tables().Get(tableName); err != nil {
		return nil, ErrResourceNotFound
	}

	result, err := s.getResourcePolicyCore(store, tableName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Policy":     result.Policy,
		"RevisionId": result.RevisionId,
	}, nil
}

// PutResourcePolicy creates or updates a resource policy for a DynamoDB table.
func (s *DynamoDBService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if !validateResourceArnString(resourceArn) {
		return nil, ErrInvalidParameter
	}

	policy, ok := req.Parameters["Policy"].(string)
	if !ok || policy == "" {
		return nil, ErrInvalidParameter
	}

	expectedRevisionId := request.GetStringParam(req.Parameters, "ExpectedRevisionId")
	if !validatePolicyRevisionId(expectedRevisionId) {
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
	if _, err := store.Tables().Get(tableName); err != nil {
		return nil, ErrResourceNotFound
	}

	result, err := s.putResourcePolicyCore(store, PutResourcePolicyInput{
		TableName:          tableName,
		Policy:             policy,
		ExpectedRevisionId: expectedRevisionId,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"RevisionId": result.RevisionId,
	}, nil
}
