package appsync

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/tags"
	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

func appsyncMapError(err error) error {
	switch err.(type) {
	case *tags.MissingResourceError:
		return NewBadRequestException("resourceArn is required")
	case *tags.MissingTagsError:
		return NewBadRequestException("tags are required")
	case *tags.MissingTagKeysError:
		return NewBadRequestException("tagKeys are required")
	}
	return err
}

func appsyncTagConfig(store *appsyncstore.AppSyncStore, req *request.ParsedRequest) tags.TagHandlerConfig {
	return tags.TagHandlerConfig{
		Param: tags.TagOperationConfig{
			ResourceParam:    "resourceArn",
			TagsParam:        "tags",
			TagKeysParam:     "tagKeys",
			TagKeyName:       "Key",
			TagValueName:     "Value",
			RequireTags:      true,
			RequireTagKeys:   true,
			RequireResource:  true,
			UseQueryFallback: false,
		},
		ParseTags: func(_ map[string]interface{}) []types.Tag {
			m := parseTags(req.Parameters)
			return tags.MapToTags(m)
		},
		ParseTagKeys: func(_ map[string]interface{}) []string {
			return parseTagKeysFromQuery(req)
		},
		TagFunc: func(_ context.Context, resourceKey string, tag []types.Tag) error {
			return store.TagStore.Tag(resourceKey, tags.ToMap(tag))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return store.TagStore.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			m, err := store.TagStore.List(resourceKey)
			if err != nil {
				return nil, err
			}
			return tags.MapToTags(m), nil
		},
		FormatResponse: func(tag []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"tags": tags.ToMap(tag),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		MapError: appsyncMapError,
		ValidateResource: func(_ context.Context, resourceArn string) error {
			return validateAppSyncResource(store, resourceArn)
		},
	}
}

func validateAppSyncResource(store *appsyncstore.AppSyncStore, resourceArn string) error {
	_, _, _, _, resource := arn.SplitARN(resourceArn)
	if resource == "" {
		return NewBadRequestException(fmt.Sprintf("Invalid resource ARN: %s", resourceArn))
	}

	switch {
	case strings.Contains(resource, "/channelNamespaces/"):
		prefix := "apis/"
		if !strings.HasPrefix(resource, prefix) {
			return NewBadRequestException(fmt.Sprintf("Invalid channel namespace ARN: %s", resourceArn))
		}
		withoutType := resource[len(prefix):]
		idx := strings.Index(withoutType, "/channelNamespaces/")
		if idx <= 0 {
			return NewBadRequestException(fmt.Sprintf("Invalid channel namespace ARN: %s", resourceArn))
		}
		apiId := withoutType[:idx]
		after := withoutType[idx+len("/channelNamespaces/"):]
		name := strings.SplitN(after, "/", 2)[0]
		_, err := store.GetChannelNamespace(apiId, name)
		if err != nil {
			return NewNotFoundException("Channel namespace")
		}

	case strings.Contains(resource, "/datasources/"):
		prefix := "apis/"
		if !strings.HasPrefix(resource, prefix) {
			return NewBadRequestException(fmt.Sprintf("Invalid data source ARN: %s", resourceArn))
		}
		withoutType := resource[len(prefix):]
		idx := strings.Index(withoutType, "/datasources/")
		if idx <= 0 {
			return NewBadRequestException(fmt.Sprintf("Invalid data source ARN: %s", resourceArn))
		}
		apiId := withoutType[:idx]
		after := withoutType[idx+len("/datasources/"):]
		name := strings.SplitN(after, "/", 2)[0]
		_, err := store.GetDataSource(apiId, name)
		if err != nil {
			return NewNotFoundException("Data source")
		}

	default:
		parts := strings.SplitN(resource, "/", 2)
		if len(parts) < 2 {
			return NewBadRequestException(fmt.Sprintf("Invalid resource ARN: %s", resourceArn))
		}
		id := parts[1]
		_, errApi := store.GetApiById(id)
		if errApi == nil {
			return nil
		}
		_, errGql := store.GetGraphqlApiById(id)
		if errGql == nil {
			return nil
		}
		return NewNotFoundException("API")
	}

	return nil
}

// TagResource adds or overwrites tags on an AppSync API.
func (s *AppSyncService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, appsyncTagConfig(store, req))
}

// UntagResource removes the specified tags from an AppSync API.
func (s *AppSyncService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleUntag(ctx, req, appsyncTagConfig(store, req))
}

// ListTagsForResource lists all tags assigned to an AppSync API.
func (s *AppSyncService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tags.HandleList(ctx, req, appsyncTagConfig(store, req))
}

func parseTagKeysFromQuery(req *request.ParsedRequest) []string {
	keys := req.QueryParams["tagKeys"]
	result := make([]string, 0, len(keys))
	for _, k := range keys {
		if k != "" {
			result = append(result, k)
		}
	}
	return result
}
