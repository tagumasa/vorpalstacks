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
)

// This file holds the tag-operation Core methods shared by the HTTP API
// surface. The handlers are thin adapters over the shared tag machinery;
// all validation, resolution and persistence live here so behaviour cannot
// drift between surfaces.

// tagResourceCore validates the secret and the tag payload, enforces the
// merged tag quota, and applies the tags. It is the single mutation path
// for TagResource.
func (s *SecretsManagerService) tagResourceCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, secretId string, tags []tagutil.Tag) error {
	if err := validateSecretId(secretId); err != nil {
		return err
	}
	// Enforce AWS Secrets Manager tag quotas.
	if err := validateSecretTags(tags); err != nil {
		return err
	}
	secret, err := resolveSecret(store, secretId)
	if err != nil {
		return err
	}
	// Check total tag count after merge.
	existingTags, _ := store.ListSecretTags(secret.Name)
	merged := make(map[string]string)
	for k, v := range existingTags {
		merged[k] = v
	}
	for _, t := range tags {
		merged[t.Key] = t.Value
	}
	if len(merged) > maxTagsPerSecret {
		return errors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("You can't have more than %d tags on a secret.", maxTagsPerSecret), http.StatusBadRequest)
	}
	return store.TagSecret(secret.Name, tagutil.ToMap(tags))
}

// untagResourceCore validates the secret and the tag keys, then removes the
// tags.
func (s *SecretsManagerService) untagResourceCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, secretId string, tagKeys []string) error {
	if err := validateSecretId(secretId); err != nil {
		return err
	}
	secret, err := resolveSecret(store, secretId)
	if err != nil {
		return err
	}
	if err := validateUntagKeys(tagKeys); err != nil {
		return err
	}
	return store.UntagSecret(secret.Name, tagKeys)
}

// listTagsForResourceCore validates the secret and returns its tags.
func (s *SecretsManagerService) listTagsForResourceCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, secretId string) ([]tagutil.Tag, error) {
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}
	secret, err := resolveSecret(store, secretId)
	if err != nil {
		return nil, err
	}
	m, err := store.ListSecretTags(secret.Name)
	if err != nil {
		return nil, err
	}
	return tagutil.MapToTags(m), nil
}

// secretsManagerTagConfig wires the shared tag machinery onto the Core
// methods above.
func secretsManagerTagConfig(store secretsmanagerstore.SecretStoreInterface, s *SecretsManagerService) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam: "SecretId",
			TagsParam:     "Tags",
			TagKeysParam:  "TagKeys",
			TagKeyName:    "Key",
			TagValueName:  "Value",
		},
		ResourceKey: func(rawKey string) string { return rawKey },
		ParseTags: func(params map[string]interface{}) []tagutil.Tag {
			return tagutil.ParseTagsWithQueryFallback(params, "Tags")
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return request.GetStringList(params, "TagKeys")
		},
		TagFunc: func(_ context.Context, secretId string, tags []tagutil.Tag) error {
			return s.tagResourceCore(context.Background(), store, secretId, tags)
		},
		UntagFunc: func(_ context.Context, secretId string, tagKeys []string) error {
			return s.untagResourceCore(context.Background(), store, secretId, tagKeys)
		},
		ListFunc: func(_ context.Context, secretId string) ([]tagutil.Tag, error) {
			return s.listTagsForResourceCore(context.Background(), store, secretId)
		},
		FormatResponse: func(tags []tagutil.Tag, _ string) (interface{}, error) {
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
