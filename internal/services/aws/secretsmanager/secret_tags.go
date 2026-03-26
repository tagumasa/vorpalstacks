package secretsmanager

import (
	"context"

	"vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
)

// TagResource tags a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_TagResource.html
func (s *SecretsManagerService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))
	if len(tags) > 0 {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		if err := store.TagSecret(secret.Name, tags); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return response.EmptyResponse(), nil
}

// UntagResource removes tags from a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_UntagResource.html
func (s *SecretsManagerService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	tagKeys := request.GetStringList(req.Parameters, "TagKeys")
	if len(tagKeys) > 0 {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		if err := store.UntagSecret(secret.Name, tagKeys); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return response.EmptyResponse(), nil
}

// ListTagsForResource lists the tags for a secret.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_ListTagsForResource.html
func (s *SecretsManagerService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tags, err := store.ListSecretTags(secret.Name)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"Tags": tagutil.MapToResponse(tags),
	}, nil
}
