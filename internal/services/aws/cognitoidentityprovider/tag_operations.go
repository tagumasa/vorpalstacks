package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	"vorpalstacks/internal/utils/aws/types"
)

func cognitoIdpMapError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError, *tagutil.MissingTagsError, *tagutil.MissingTagKeysError:
		return ErrInvalidParameter
	}
	return err
}

func cognitoIdpTagConfig(store cognitostore.CognitoStoreInterface) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		TagFunc: func(ctx context.Context, resourceKey string, tags []types.Tag) error {
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
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			tags, err := store.ListAsSlice(resourceKey)
			if err != nil {
				return nil, ErrInternalError
			}
			return tags, nil
		},
		FormatResponse: func(tags []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{"Tags": tagutil.ToMap(tags)}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: cognitoIdpMapError,
	}
}

// TagResource adds or overwrites tags on a Cognito IdP resource.
func (s *CognitoService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, cognitoIdpTagConfig(store))
}

// UntagResource removes the specified tags from a Cognito IdP resource.
func (s *CognitoService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, cognitoIdpTagConfig(store))
}

// ListTagsForResource lists all tags assigned to a Cognito IdP resource.
func (s *CognitoService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, cognitoIdpTagConfig(store))
}
