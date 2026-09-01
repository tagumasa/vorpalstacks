package cloudwatch

import (
	"context"
	"errors"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// cwAlarmTagConfig builds the tag-engine configuration for the CloudWatch
// alarm tag operations. The closures are the single tag persistence path:
// each resolves the regional alarm store and delegates to it. Every tag
// operation first resolves the alarm behind the ARN, so a tag against a
// nonexistent alarm fails with the modelled ResourceNotFoundException
// instead of silently persisting tags under an unowned key.
func cwAlarmTagConfig(s *CloudWatchService, reqCtx *request.RequestContext) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardARNConfig,
		ValidateResource: func(_ context.Context, resourceKey string) error {
			// CloudWatch tags alarms; the resource field is "alarm:<name>".
			name := extractAlarmNameFromARN(resourceKey)
			if name == "" {
				return ErrAlarmNotFound
			}
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			if _, err := store.alarms.GetAlarm(name); err != nil {
				if errors.Is(err, cwstore.ErrAlarmNotFound) {
					return ErrAlarmNotFound
				}
				return err
			}
			return nil
		},
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
