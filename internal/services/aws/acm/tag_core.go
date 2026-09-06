package acm

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

func (s *ACMService) acmTagConfig(stores *acmStores) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:    "CertificateArn",
			TagsParam:        "Tags",
			TagKeysParam:     "Tags",
			TagKeyName:       "Key",
			TagValueName:     "Value",
			RequireTags:      true,
			RequireTagKeys:   true,
			RequireResource:  true,
			UseQueryFallback: true,
		},
		ResourceKey: func(rawKey string) string {
			return rawKey
		},
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			_, err := s.fetchCertificate(stores, resourceKey)
			return err
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagList []tagutil.Tag) error {
			cert, err := s.fetchCertificate(stores, resourceKey)
			if err != nil {
				return err
			}
			cert.Tags = tagutil.Apply(cert.Tags, tagList)
			return stores.certificates.Update(cert)
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			tagsToRemove := tagutil.ParseTagsWithQueryFallback(params, "Tags")
			keys := make([]string, 0, len(tagsToRemove))
			for _, t := range tagsToRemove {
				keys = append(keys, t.Key)
			}
			return keys
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			cert, err := s.fetchCertificate(stores, resourceKey)
			if err != nil {
				return err
			}
			tagKeySet := make(map[string]bool, len(tagKeys))
			for _, k := range tagKeys {
				tagKeySet[k] = true
			}
			cert.Tags = tagutil.Remove(cert.Tags, tagKeySet)
			return stores.certificates.Update(cert)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			cert, err := s.fetchCertificate(stores, resourceKey)
			if err != nil {
				return nil, err
			}
			return cert.Tags, nil
		},
		FormatResponse: func(tagList []tagutil.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"Tags": tagutil.ToResponse(tagList),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: func(err error) error {
			switch err.(type) {
			case *tagutil.MissingResourceError:
				return awserrors.NewResourceNotFoundException("certificate", "")
			case *tagutil.MissingTagsError:
				return awserrors.NewValidationException("Tags are required")
			case *tagutil.MissingTagKeysError:
				return awserrors.NewValidationException("Tag keys are required")
			}
			return err
		},
	}
}

// acmGenericTagConfig returns a TagHandlerConfig for the generic TagResource /
// UntagResource / ListTagsForResource operations. These use "ResourceArn" as
// the resource parameter and "TagKeys" (a plain list of key strings) for
// untagging, instead of the certificate-specific parameter names.
func (s *ACMService) acmGenericTagConfig(stores *acmStores) tagutil.TagHandlerConfig {
	config := s.acmTagConfig(stores)
	config.Param.ResourceParam = "ResourceArn"
	config.Param.TagKeysParam = "TagKeys"
	config.ParseTagKeys = func(params map[string]interface{}) []string {
		return tagutil.ParseTagKeysWithQueryFallback(params, "TagKeys")
	}
	return config
}
