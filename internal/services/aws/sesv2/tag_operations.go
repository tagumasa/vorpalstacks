package sesv2

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/aws/types"
)

func sesv2TagConfig(s *SESv2Service, reqCtx *request.RequestContext) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		TagFunc: func(ctx context.Context, resourceKey string, tags []types.Tag) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.TagFromSlice(resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			store, err := s.store(reqCtx)
			if err != nil {
				return nil, err
			}
			return store.ListAsSlice(resourceKey)
		},
		EmptyResponse: func() (interface{}, error) { return response.EmptyResponse(), nil },
	}
}

// TagResource adds tags to an SESv2 resource.
func (s *SESv2Service) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleTag(ctx, req, sesv2TagConfig(s, reqCtx))
}

// UntagResource removes tags from an SESv2 resource.
func (s *SESv2Service) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleUntag(ctx, req, sesv2TagConfig(s, reqCtx))
}

// ListTagsForResource lists tags for an SESv2 resource.
func (s *SESv2Service) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleList(ctx, req, sesv2TagConfig(s, reqCtx))
}
