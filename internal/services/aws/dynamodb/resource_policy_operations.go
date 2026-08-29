package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// DeleteResourcePolicy deletes a resource policy from a DynamoDB table.
func (s *DynamoDBService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteResourcePolicyCore(ctx, reqCtx, DeleteResourcePolicyInput{
		ResourceArn:        request.GetStringParam(req.Parameters, "ResourceArn"),
		ExpectedRevisionId: request.GetStringParam(req.Parameters, "ExpectedRevisionId"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetResourcePolicy returns the resource policy for a DynamoDB table.
func (s *DynamoDBService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getResourcePolicyCore(ctx, reqCtx, GetResourcePolicyInput{
		ResourceArn: request.GetStringParam(req.Parameters, "ResourceArn"),
	})
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
	policy, _ := req.Parameters["Policy"].(string)
	result, err := s.putResourcePolicyCore(ctx, reqCtx, PutResourcePolicyInput{
		ResourceArn:        request.GetStringParam(req.Parameters, "ResourceArn"),
		Policy:             policy,
		ExpectedRevisionId: request.GetStringParam(req.Parameters, "ExpectedRevisionId"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"RevisionId": result.RevisionId,
	}, nil
}
