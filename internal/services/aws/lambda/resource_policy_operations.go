package lambda

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// PutResourcePolicy replaces the resource-based policy of a Lambda
// resource with the submitted policy document.
func (s *LambdaService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.putResourcePolicyCore(store, &PutResourcePolicyInput{
		ResourceArn: request.GetStringParam(req.Parameters, "ResourceArn"),
		Policy:      request.GetStringParam(req.Parameters, "Policy"),
		RevisionId:  request.GetStringParam(req.Parameters, "RevisionId"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"Policy":     result.Policy,
		"RevisionId": result.RevisionId,
	}, nil
}

// GetResourcePolicy retrieves the resource-based policy attached to a
// Lambda resource together with the policy revision.
func (s *LambdaService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.getResourcePolicyCore(store, request.GetStringParam(req.Parameters, "ResourceArn"))
	if err != nil {
		return nil, err
	}
	resp := map[string]interface{}{
		"Policy": result.Policy,
	}
	if result.RevisionId != "" {
		resp["RevisionId"] = result.RevisionId
	}
	return resp, nil
}

// DeleteResourcePolicy removes the resource-based policy from a Lambda
// resource.
func (s *LambdaService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteResourcePolicyCore(store, &DeleteResourcePolicyInput{
		ResourceArn: request.GetStringParam(req.Parameters, "ResourceArn"),
		RevisionId:  request.GetStringParam(req.Parameters, "RevisionId"),
	}); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}
