package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
	"vorpalstacks/internal/utils/aws/types"
)

// iotTagConfig builds a TagHandlerConfig backed by the IoT store.
func iotTagConfig(store iotstore.IotStoreInterface) tagutil.TagHandlerConfig {
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

// ListActiveViolations lists active Device Defender violations.
func (s *IoTService) ListActiveViolations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"activeViolations": []map[string]interface{}{},
	}, nil
}

// ListViolationEvents lists Device Defender violation events.
func (s *IoTService) ListViolationEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"violationEvents": []map[string]interface{}{},
	}, nil
}

// GetBehaviorModelTrainingSummaries retrieves ML model training summaries.
func (s *IoTService) GetBehaviorModelTrainingSummaries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"summaries": []map[string]interface{}{},
	}, nil
}

// ValidateSecurityProfileBehaviors validates security profile behaviours.
func (s *IoTService) ValidateSecurityProfileBehaviors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"valid": true,
	}, nil
}

// DescribeSecurityProfile retrieves a Device Defender security profile.
func (s *IoTService) DescribeSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{}, nil
}

// UpdateSecurityProfile updates a Device Defender security profile.
func (s *IoTService) UpdateSecurityProfile(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{}, nil
}
