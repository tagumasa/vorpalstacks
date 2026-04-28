package tags

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/utils/aws/types"
)

// TagHandlerConfig configures the handler framework for a single service's
// tag operations. Services provide this once; the three Handle* functions
// use it to eliminate boilerplate parameter parsing, store delegation, and
// response formatting.
type TagHandlerConfig struct {
	// Param is the TagOperationConfig for parameter extraction (resource key,
	// tags, tag keys). Reuses the existing GetResourceKey/GetTags/GetTagKeys.
	Param TagOperationConfig

	// ResourceKey transforms the raw resource key before passing it to the
	// store. For example, Lambda extracts the function name from an ARN.
	// If nil, the raw resource key is used as-is.
	ResourceKey func(rawKey string) string

	// ParseTags overrides tag parsing from request parameters.
	// If nil, GetTags(req.Parameters, Param) is used.
	ParseTags func(params map[string]interface{}) []types.Tag

	// ParseTagKeys overrides tag-key parsing from request parameters.
	// If nil, GetTagKeys(req.Parameters, Param) is used.
	ParseTagKeys func(params map[string]interface{}) []string

	// TagFunc tags a resource. Receives the (possibly transformed) resource
	// key and the parsed tags. Must call the underlying TagStore.
	TagFunc func(ctx context.Context, resourceKey string, tags []types.Tag) error

	// UntagFunc untags a resource. Receives the resource key and tag keys.
	UntagFunc func(ctx context.Context, resourceKey string, tagKeys []string) error

	// ListFunc lists tags for a resource. Returns []types.Tag.
	ListFunc func(ctx context.Context, resourceKey string) ([]types.Tag, error)

	// FormatResponse builds the response envelope for ListTags.
	// Receives []types.Tag and the raw resource key (before ResourceKey transform).
	// If nil, returns {"<TagsParam>": ToResponse(tags)}.
	FormatResponse func(tags []types.Tag, rawResourceKey string) (interface{}, error)

	// ValidateResource is called before tag/untag/list operations.
	// Receives the raw resource key. Return an error to reject the request.
	// If nil, no pre-validation is performed.
	ValidateResource func(ctx context.Context, rawResourceKey string) error

	// TagResponse builds the response for a successful TagResource operation.
	// If nil, EmptyResponse is used (or nil if EmptyResponse is also nil).
	// SNS uses this to return the tag list after tagging.
	TagResponse func(ctx context.Context, resourceKey string) (interface{}, error)

	// EmptyResponse returns the response for successful Untag operations.
	// If nil, returns nil (no body).
	EmptyResponse func() (interface{}, error)

	// MapError maps handler errors to service-specific errors.
	// Called on any non-nil error returned by the handler.
	MapError func(error) error
}

func applyMapError(cfg TagHandlerConfig, err error) error {
	if err != nil && cfg.MapError != nil {
		return cfg.MapError(err)
	}
	return err
}

// HandleTag processes a tag-resource request using the handler config.
// Flow: extract resource key → validate → parse tags → TagFunc → response.
func HandleTag(ctx context.Context, req *request.ParsedRequest, cfg TagHandlerConfig) (interface{}, error) {
	rawKey := GetResourceKey(req.Parameters, cfg.Param)
	if rawKey == "" {
		return nil, applyMapError(cfg, &MissingResourceError{Param: cfg.Param.ResourceParam})
	}

	resourceKey := rawKey
	if cfg.ResourceKey != nil {
		resourceKey = cfg.ResourceKey(rawKey)
	}

	if cfg.ValidateResource != nil {
		if err := cfg.ValidateResource(ctx, rawKey); err != nil {
			return nil, applyMapError(cfg, err)
		}
	}

	var tags []types.Tag
	if cfg.ParseTags != nil {
		tags = cfg.ParseTags(req.Parameters)
	} else {
		tags = GetTags(req.Parameters, cfg.Param)
	}

	if cfg.Param.RequireTags && len(tags) == 0 {
		return nil, applyMapError(cfg, &MissingTagsError{Param: cfg.Param.TagsParam})
	}

	if len(tags) > 0 && cfg.TagFunc != nil {
		if err := cfg.TagFunc(ctx, resourceKey, tags); err != nil {
			return nil, applyMapError(cfg, err)
		}
	}

	if cfg.TagResponse != nil {
		resp, err := cfg.TagResponse(ctx, resourceKey)
		return resp, applyMapError(cfg, err)
	}
	if cfg.EmptyResponse != nil {
		resp, err := cfg.EmptyResponse()
		return resp, applyMapError(cfg, err)
	}
	return nil, nil
}

// HandleUntag processes an untag-resource request using the handler config.
// Flow: extract resource key → validate → parse tag keys → UntagFunc → response.
func HandleUntag(ctx context.Context, req *request.ParsedRequest, cfg TagHandlerConfig) (interface{}, error) {
	rawKey := GetResourceKey(req.Parameters, cfg.Param)
	if rawKey == "" {
		return nil, applyMapError(cfg, &MissingResourceError{Param: cfg.Param.ResourceParam})
	}

	resourceKey := rawKey
	if cfg.ResourceKey != nil {
		resourceKey = cfg.ResourceKey(rawKey)
	}

	if cfg.ValidateResource != nil {
		if err := cfg.ValidateResource(ctx, rawKey); err != nil {
			return nil, applyMapError(cfg, err)
		}
	}

	var tagKeys []string
	if cfg.ParseTagKeys != nil {
		tagKeys = cfg.ParseTagKeys(req.Parameters)
	} else {
		tagKeys = GetTagKeys(req.Parameters, cfg.Param)
	}

	if cfg.Param.RequireTagKeys && len(tagKeys) == 0 {
		return nil, applyMapError(cfg, &MissingTagKeysError{Param: cfg.Param.TagKeysParam})
	}

	if len(tagKeys) > 0 && cfg.UntagFunc != nil {
		if err := cfg.UntagFunc(ctx, resourceKey, tagKeys); err != nil {
			return nil, applyMapError(cfg, err)
		}
	}

	if cfg.EmptyResponse != nil {
		resp, err := cfg.EmptyResponse()
		return resp, applyMapError(cfg, err)
	}
	return nil, nil
}

// HandleList processes a list-tags request using the handler config.
// Flow: extract resource key → validate → ListFunc → FormatResponse.
func HandleList(ctx context.Context, req *request.ParsedRequest, cfg TagHandlerConfig) (interface{}, error) {
	rawKey := GetResourceKey(req.Parameters, cfg.Param)
	if rawKey == "" {
		return nil, applyMapError(cfg, &MissingResourceError{Param: cfg.Param.ResourceParam})
	}

	resourceKey := rawKey
	if cfg.ResourceKey != nil {
		resourceKey = cfg.ResourceKey(rawKey)
	}

	if cfg.ValidateResource != nil {
		if err := cfg.ValidateResource(ctx, rawKey); err != nil {
			return nil, applyMapError(cfg, err)
		}
	}

	if cfg.ListFunc == nil {
		return nil, applyMapError(cfg, &MissingResourceError{Param: "ListFunc"})
	}

	tags, err := cfg.ListFunc(ctx, resourceKey)
	if err != nil {
		return nil, applyMapError(cfg, err)
	}

	if cfg.FormatResponse != nil {
		resp, err := cfg.FormatResponse(tags, rawKey)
		return resp, applyMapError(cfg, err)
	}

	var tagResponse []map[string]interface{}
	if cfg.Param.TagKeyName == "Key" && cfg.Param.TagValueName == "Value" {
		tagResponse = ToResponse(tags)
	} else {
		tagResponse = ToResponseWithKeyNames(tags, cfg.Param.TagKeyName, cfg.Param.TagValueName)
	}

	return map[string]interface{}{
		cfg.Param.TagsParam: tagResponse,
	}, nil
}

// MissingResourceError is returned when the resource key parameter is missing.
type MissingResourceError struct {
	Param string
}

// Error returns a human-readable error message.
func (e *MissingResourceError) Error() string {
	return e.Param + " is required"
}

// MissingTagsError is returned when the tags parameter is missing.
type MissingTagsError struct {
	Param string
}

// Error returns a human-readable error message.
func (e *MissingTagsError) Error() string {
	return e.Param + " is required"
}

// MissingTagKeysError is returned when the tag keys parameter is missing.
type MissingTagKeysError struct {
	Param string
}

// Error returns a human-readable error message.
func (e *MissingTagKeysError) Error() string {
	return e.Param + " is required"
}

// SliceTagsToKeys extracts keys from a []types.Tag as []string.
func SliceTagsToKeys(tags []types.Tag) []string {
	keys := make([]string, len(tags))
	for i, t := range tags {
		keys[i] = t.Key
	}
	return keys
}
