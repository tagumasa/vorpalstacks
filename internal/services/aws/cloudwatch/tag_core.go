package cloudwatch

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// cwAlarmTagConfig builds the tag-engine configuration for the CloudWatch
// alarm tag operations. The closures are the single tag persistence path:
// each resolves the regional alarm store and delegates to it.
func cwAlarmTagConfig(s *CloudWatchService, reqCtx *request.RequestContext) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardARNConfig,
		TagFunc: func(_ context.Context, resourceKey string, tags []tagutil.Tag) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.alarms.TagFromSlice(resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.alarms.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			store, err := s.store(reqCtx)
			if err != nil {
				return nil, err
			}
			return store.alarms.ListAsSlice(resourceKey)
		},
		EmptyResponse: func() (interface{}, error) { return response.EmptyResponse(), nil },
	}
}
