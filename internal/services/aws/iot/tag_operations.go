package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
	"vorpalstacks/internal/utils/aws/types"
)

// iotTagConfig builds a TagHandlerConfig backed by the IoT store.
func iotTagConfig(store iotstore.TagOps) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		TagFunc: func(_ context.Context, resourceKey string, tags []types.Tag) error {
			return store.TagResource(resourceKey, tagutil.ToMap(tags))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.UntagResource(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			tagsMap, err := store.ListTags(resourceKey)
			if err != nil {
				return nil, err
			}
			return tagutil.MapToTags(tagsMap), nil
		},
		FormatResponse: func(tags []types.Tag, _ string) (interface{}, error) {
			tagList := make([]map[string]interface{}, 0, len(tags))
			for _, t := range tags {
				tagList = append(tagList, map[string]interface{}{
					"Key":   t.Key,
					"Value": t.Value,
				})
			}
			return map[string]interface{}{"tags": tagList}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		MapError: func(err error) error {
			if _, ok := err.(*tagutil.MissingResourceError); ok {
				return iotstore.ErrMissingParam
			}
			return err
		},
	}
}

// TagResource adds or overwrites tags on an IoT resource.
func (s *IoTService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, iotTagConfig(store))
}

// UntagResource removes tags from an IoT resource.
func (s *IoTService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, iotTagConfig(store))
}

// ListTagsForResource lists tags attached to an IoT resource.
func (s *IoTService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, iotTagConfig(store))
}

func (s *IoTService) ListActiveViolations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"activeViolations": []map[string]interface{}{},
	}, nil
}

func (s *IoTService) ListViolationEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"violationEvents": []map[string]interface{}{},
	}, nil
}

func (s *IoTService) GetBehaviorModelTrainingSummaries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"summaries": []map[string]interface{}{},
	}, nil
}
