package lambda

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// TagResource adds tags to a Lambda function.
func (s *LambdaService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := svcarn.ExtractFunctionNameFromARN(request.GetStringParam(req.Parameters, "Resource"))
	if functionName == "" {
		return nil, NewInvalidParameter("Resource", "Resource ARN is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.Functions.Get(functionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	tags := tagutil.ParseTags(req.Parameters, "Tags")
	if len(tags) == 0 {
		return nil, NewInvalidParameter("Tags", "Tags must be a map")
	}

	if err := store.Functions.AddTags(functionName, tags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a Lambda function.
func (s *LambdaService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := svcarn.ExtractFunctionNameFromARN(request.GetStringParam(req.Parameters, "Resource"))
	if functionName == "" {
		return nil, NewInvalidParameter("Resource", "Resource ARN is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	function, err := store.Functions.Get(functionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	tagKeysToRemove := tagutil.ParseTagKeys(req.Parameters, "TagKeys")
	newTags := tagutil.Remove(function.Tags, tagKeysToRemove)

	if err := store.Functions.SetTags(functionName, newTags); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTags lists the tags for a Lambda function.
func (s *LambdaService) ListTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	functionName := svcarn.ExtractFunctionNameFromARN(request.GetStringParam(req.Parameters, "Resource"))
	if functionName == "" {
		return nil, NewInvalidParameter("Resource", "Resource ARN is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tags, err := store.Functions.GetTags(functionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"Tags": tagutil.ToMap(tags),
	}, nil
}
