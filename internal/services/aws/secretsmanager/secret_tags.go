package secretsmanager

import (
	"context"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
	"vorpalstacks/internal/utils/aws/types"
)

func secretsManagerMapError(err error) error {
	switch e := err.(type) {
	case *tagutil.MissingResourceError:
		return errors.ErrMissingParameter
	case *tagutil.MissingTagsError:
		_ = e
		return errors.ErrMissingParameter
	case *tagutil.MissingTagKeysError:
		_ = e
		return errors.ErrMissingParameter
	}
	return err
}

func secretsManagerTagConfig(store secretsmanagerstore.SecretStoreInterface, secretName string) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam: "SecretId",
			TagsParam:     "Tags",
			TagKeysParam:  "TagKeys",
			TagKeyName:    "Key",
			TagValueName:  "Value",
		},
		ResourceKey: func(_ string) string { return secretName },
		ParseTags: func(params map[string]interface{}) []types.Tag {
			return tagutil.ParseTagsWithQueryFallback(params, "Tags")
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return request.GetStringList(params, "TagKeys")
		},
		TagFunc: func(_ context.Context, resourceKey string, tags []types.Tag) error {
			return store.TagSecret(resourceKey, tagutil.ToMap(tags))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.UntagSecret(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			m, err := store.ListSecretTags(resourceKey)
			if err != nil {
				return nil, err
			}
			return tagutil.MapToTags(m), nil
		},
		FormatResponse: func(tags []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.ToResponse(tags),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: secretsManagerMapError,
	}
}

// TagResource adds or overwrites tags on a Secrets Manager secret.
func (s *SecretsManagerService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	return tagutil.HandleTag(ctx, req, secretsManagerTagConfig(store, secret.Name))
}

// UntagResource removes the specified tags from a Secrets Manager secret.
func (s *SecretsManagerService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
	return tagutil.HandleUntag(ctx, req, secretsManagerTagConfig(store, secret.Name))
}

// ListTagsForResource lists all tags assigned to a Secrets Manager secret.
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
	return tagutil.HandleList(ctx, req, secretsManagerTagConfig(store, secret.Name))
}
