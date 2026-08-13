package secretsmanager

import (
	"context"
	"fmt"
	"net/http"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
	"vorpalstacks/internal/utils/aws/types"
)

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
		MapError: func(err error) error {
			switch err.(type) {
			case *tagutil.MissingResourceError:
				return errors.ErrMissingParameter
			case *tagutil.MissingTagsError:
				return errors.ErrMissingParameter
			case *tagutil.MissingTagKeysError:
				return errors.ErrMissingParameter
			}
			return mapStoreError(err)
		},
	}
}

// TagResource adds or overwrites tags on a Secrets Manager secret.
func (s *SecretsManagerService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}
	// Enforce AWS Secrets Manager tag quotas.
	newTags := tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")
	if err := validateSecretTags(newTags); err != nil {
		return nil, err
	}
	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Check total tag count after merge.
	existingTags, _ := store.ListSecretTags(secret.Name)
	merged := make(map[string]string)
	for k, v := range existingTags {
		merged[k] = v
	}
	for _, t := range newTags {
		merged[t.Key] = t.Value
	}
	if len(merged) > maxTagsPerSecret {
		return nil, errors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("You can't have more than %d tags on a secret.", maxTagsPerSecret), http.StatusBadRequest)
	}
	return tagutil.HandleTag(ctx, req, secretsManagerTagConfig(store, secret.Name))
}

// UntagResource removes the specified tags from a Secrets Manager secret.
func (s *SecretsManagerService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}
	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}
	tagKeys := request.GetStringList(req.Parameters, "TagKeys")
	if err := validateUntagKeys(tagKeys); err != nil {
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
	if err := validateSecretId(secretId); err != nil {
		return nil, err
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
