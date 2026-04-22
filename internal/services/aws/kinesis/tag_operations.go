package kinesis

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
	"vorpalstacks/internal/utils/aws/types"
)

func (s *KinesisService) kinesisTagConfig(store *kinesisstore.KinesisStore, req *request.ParsedRequest) tags.TagHandlerConfig {
	return tags.TagHandlerConfig{
		Param: tags.KinesisStreamConfig,
		ResourceKey: func(rawKey string) string {
			name := request.GetParamLowerFirst(req.Parameters, "StreamName")
			if name != "" {
				return name
			}
			arn := request.GetParamLowerFirst(req.Parameters, "StreamARN")
			if arn == "" {
				arn = request.GetParamLowerFirst(req.Parameters, "ResourceARN")
			}
			if arn == "" {
				return rawKey
			}
			stream, err := store.GetStreamByARN(arn)
			if err != nil {
				return ""
			}
			return stream.StreamName
		},
		ParseTags: func(params map[string]interface{}) []types.Tag {
			return tags.ParseTags(params, "Tags")
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagList []types.Tag) error {
			return store.Tag(resourceKey, tags.ToMap(tagList))
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return store.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			return store.ListAsSlice(resourceKey)
		},
		FormatResponse: func(tagList []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"Tags":        tags.ToResponse(tagList),
				"HasMoreTags": false,
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: func(err error) error {
			switch err.(type) {
			case *tags.MissingResourceError:
				return ErrInvalidArgument
			case *tags.MissingTagsError:
				return ErrInvalidArgument
			case *tags.MissingTagKeysError:
				return ErrInvalidArgument
			}
			return s.mapStoreError(err)
		},
	}
}

// AddTagsToStream adds or overwrites tags on a Kinesis data stream.
func (s *KinesisService) AddTagsToStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, s.kinesisTagConfig(store, req))
}

// RemoveTagsFromStream removes the specified tags from a Kinesis data stream.
func (s *KinesisService) RemoveTagsFromStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleUntag(ctx, req, s.kinesisTagConfig(store, req))
}

// ListTagsForStream lists all tags assigned to a Kinesis data stream.
func (s *KinesisService) ListTagsForStream(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleList(ctx, req, s.kinesisTagConfig(store, req))
}

// TagResource adds or overwrites tags on a Kinesis resource identified by ARN.
func (s *KinesisService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, s.kinesisTagConfig(store, req))
}

// UntagResource removes the specified tags from a Kinesis resource identified by ARN.
func (s *KinesisService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleUntag(ctx, req, s.kinesisTagConfig(store, req))
}

// ListTagsForResource lists all tags assigned to a Kinesis resource identified by ARN.
func (s *KinesisService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleList(ctx, req, s.kinesisTagConfig(store, req))
}
