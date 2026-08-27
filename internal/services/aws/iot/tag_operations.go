package iot

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// iotTagConfig builds a TagHandlerConfig backed by the IoT store Cores.
func (s *IoTService) iotTagConfig(store iotstore.TagOps) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		TagFunc: func(_ context.Context, resourceKey string, tags []tagutil.Tag) error {
			return s.tagResourceCore(store, resourceKey, tagutil.ToMap(tags))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return s.untagResourceCore(store, resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]tagutil.Tag, error) {
			tagsMap, err := s.listTagsCore(store, resourceKey)
			if err != nil {
				return nil, err
			}
			return tagutil.MapToTags(tagsMap), nil
		},
		FormatResponse: func(tags []tagutil.Tag, _ string) (interface{}, error) {
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
				return awserrors.NewAWSError("ResourceNotFoundException", "The specified resource does not exist.", 404)
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
	return tagutil.HandleTag(ctx, req, s.iotTagConfig(store))
}

// UntagResource removes tags from an IoT resource.
func (s *IoTService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.iotTagConfig(store))
}

// ListTagsForResource lists tags attached to an IoT resource.
func (s *IoTService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.iotTagConfig(store))
}

func (s *IoTService) ListActiveViolations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	violations, err := s.listActiveViolationsCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "thingName"),
		request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"))
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(violations))
	for _, v := range violations {
		items = append(items, activeViolationResponse(v))
	}
	return paginatedMaps("activeViolations", items, req.Parameters)
}

func (s *IoTService) ListViolationEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	events, err := s.listViolationEventsCore(store, parseListOptions(req.Parameters),
		request.GetParamCaseInsensitive(req.Parameters, "securityProfileName"),
		request.GetParamCaseInsensitive(req.Parameters, "thingName"))
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(events))
	for _, e := range events {
		items = append(items, violationEventResponse(e))
	}
	return map[string]interface{}{
		"violationEvents": items,
	}, nil
}

// No ML behaviour-model training pipeline exists in this platform; an empty
// list matches the AWS wire shape.
func (s *IoTService) GetBehaviorModelTrainingSummaries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"summaries": []map[string]interface{}{},
	}, nil
}
