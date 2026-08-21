package cognitoidentity

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

func cognitoIdentityMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError, *tagutil.MissingTagsError, *tagutil.MissingTagKeysError:
		return ErrInvalidParameter
	}
	if errors.Is(err, cognitoidentitystore.ErrIdentityPoolNotFound) {
		return ErrResourceNotFound
	}
	return err
}

func cognitoIdentityTagConfig(store cognitoidentitystore.CognitoIdentityStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		TagFunc: func(ctx context.Context, resourceKey string, tags []tagutil.Tag) error {
			if err := store.Tag(resourceKey, tagutil.ToMap(tags)); err != nil {
				return ErrInternalError
			}
			return nil
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			if err := store.Untag(resourceKey, tagKeys); err != nil {
				return ErrInternalError
			}
			return nil
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			m, err := store.List(resourceKey)
			if err != nil {
				return nil, ErrInternalError
			}
			return tagutil.MapToTags(m), nil
		},
		FormatResponse: func(tags []tagutil.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{"Tags": tagutil.ToMap(tags)}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: cognitoIdentityMapError,
	}
}

// TagResource adds or overwrites tags on a Cognito Identity pool.
func (s *CognitoIdentityService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, cognitoIdentityTagConfig(store))
}

// UntagResource removes the specified tags from a Cognito Identity pool.
func (s *CognitoIdentityService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, cognitoIdentityTagConfig(store))
}

// ListTagsForResource lists all tags assigned to a Cognito Identity pool.
func (s *CognitoIdentityService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, cognitoIdentityTagConfig(store))
}
