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
			m, err := parseTags(req.Parameters)
			if err != nil {
				return nil
			}
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
		ValidateTagsFunc: func(t []types.Tag) error {
			if err := tags.ValidateTags(t); err != nil {
				return NewBadRequestException(err.Error())
			}
			return nil
		},
		ValidateResource: func(_ context.Context, resourceArn string) error {
			return validateAppSyncResource(store, resourceArn)
		},
	}
}

// extractApiScopedResource splits an "apis/{apiId}/{rest}" resource path
// into the apiId and the remainder. Returns ok=false if the path does not
// match the expected prefix.
func extractApiScopedResource(resource string) (apiId, rest string, ok bool) {
	const prefix = "apis/"
	if !strings.HasPrefix(resource, prefix) {
		return "", "", false
	}
	after := resource[len(prefix):]
	idx := strings.Index(after, "/")
	if idx <= 0 {
		return "", "", false
	}
	return after[:idx], after[idx+1:], true
}

// firstSegment returns the first path segment of s (up to the next "/").
func firstSegment(s string) string {
	if idx := strings.Index(s, "/"); idx >= 0 {
		return s[:idx]
	}
	return s
}

func validateAppSyncResource(store *appsyncstore.AppSyncStore, resourceArn string) error {
	_, _, _, _, resource := arn.SplitARN(resourceArn)
	if resource == "" {
		return NewBadRequestException(fmt.Sprintf("Invalid resource ARN: %s", resourceArn))
	}

	switch {
	// --- Api-scoped sub-resources (apis/{apiId}/...) ---

	case strings.Contains(resource, "/channelNamespaces/"):
		apiId, rest, ok := extractApiScopedResource(resource)
		if !ok || !strings.HasPrefix(rest, "channelNamespaces/") {
			return NewBadRequestException(fmt.Sprintf("Invalid channel namespace ARN: %s", resourceArn))
		}
		name := firstSegment(strings.TrimPrefix(rest, "channelNamespaces/"))
		if _, err := store.GetChannelNamespace(apiId, name); err != nil {
			return NewNotFoundException("Channel namespace")
		}

	case strings.Contains(resource, "/datasources/"):
		apiId, rest, ok := extractApiScopedResource(resource)
		if !ok || !strings.HasPrefix(rest, "datasources/") {
			return NewBadRequestException(fmt.Sprintf("Invalid data source ARN: %s", resourceArn))
		}
		name := firstSegment(strings.TrimPrefix(rest, "datasources/"))
		if _, err := store.GetDataSource(apiId, name); err != nil {
			return NewNotFoundException("Data source")
		}

	case strings.Contains(resource, "/types/") && strings.Contains(resource, "/resolvers/"):
		// Must precede the /types/ case.
		// Path: apis/{apiId}/types/{typeName}/resolvers/{fieldName}
		apiId, rest, ok := extractApiScopedResource(resource)
		if !ok || !strings.HasPrefix(rest, "types/") {
			return NewBadRequestException(fmt.Sprintf("Invalid resolver ARN: %s", resourceArn))
		}
		afterType := strings.TrimPrefix(rest, "types/")
		typeName := firstSegment(afterType)
		afterTypeName := strings.TrimPrefix(afterType, typeName+"/")
		if !strings.HasPrefix(afterTypeName, "resolvers/") {
			return NewBadRequestException(fmt.Sprintf("Invalid resolver ARN: %s", resourceArn))
		}
		fieldName := firstSegment(strings.TrimPrefix(afterTypeName, "resolvers/"))
		if _, err := store.GetResolver(apiId, typeName, fieldName); err != nil {
			return NewNotFoundException("Resolver")
		}

	case strings.Contains(resource, "/functions/"):
		apiId, rest, ok := extractApiScopedResource(resource)
		if !ok || !strings.HasPrefix(rest, "functions/") {
			return NewBadRequestException(fmt.Sprintf("Invalid function ARN: %s", resourceArn))
		}
		functionId := firstSegment(strings.TrimPrefix(rest, "functions/"))
		if _, err := store.GetFunction(apiId, functionId); err != nil {
			return NewNotFoundException("Function")
		}

	case strings.Contains(resource, "/types/"):
		// Path: apis/{apiId}/types/{typeName}
		apiId, rest, ok := extractApiScopedResource(resource)
		if !ok || !strings.HasPrefix(rest, "types/") {
			return NewBadRequestException(fmt.Sprintf("Invalid type ARN: %s", resourceArn))
		}
		typeName := firstSegment(strings.TrimPrefix(rest, "types/"))
		if _, err := store.GetType(apiId, typeName); err != nil {
			return NewNotFoundException("Type")
		}

	case strings.Contains(resource, "/apikeys/"):
		apiId, rest, ok := extractApiScopedResource(resource)
		if !ok || !strings.HasPrefix(rest, "apikeys/") {
			return NewBadRequestException(fmt.Sprintf("Invalid API key ARN: %s", resourceArn))
		}
		keyId := firstSegment(strings.TrimPrefix(rest, "apikeys/"))
		if _, err := store.GetApiKey(apiId, keyId); err != nil {
			return NewNotFoundException("ApiKey")
		}

	case strings.HasSuffix(resource, "/ApiCaches"):
		apiId, _, ok := extractApiScopedResource(resource)
		if !ok {
			return NewBadRequestException(fmt.Sprintf("Invalid API cache ARN: %s", resourceArn))
		}
		if _, err := store.GetApiCache(apiId); err != nil {
			return NewNotFoundException("ApiCache")
		}

	// --- Top-level resources (not api-scoped) ---

	case strings.HasPrefix(resource, "domainnames/"):
		domainName := strings.TrimPrefix(resource, "domainnames/")
		if domainName == "" {
			return NewBadRequestException(fmt.Sprintf("Invalid domain name ARN: %s", resourceArn))
		}
		if _, err := store.GetDomainName(domainName); err != nil {
			return NewNotFoundException("DomainName")
		}

	default:
		// GraphQL API or Event API by id: apis/{apiId}
		if !strings.HasPrefix(resource, "apis/") {
			return NewBadRequestException(fmt.Sprintf("Invalid resource ARN: %s", resourceArn))
		}
		parts := strings.SplitN(resource, "/", 2)
		if len(parts) < 2 {
			return NewBadRequestException(fmt.Sprintf("Invalid resource ARN: %s", resourceArn))
		}
		id := parts[1]
		if _, errApi := store.GetApiById(id); errApi == nil {
			return nil
		}
		if _, errGql := store.GetGraphqlApiById(id); errGql == nil {
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
		return mapStoreError(err)
	}
	// Pre-validate tags so that specific validation errors (e.g. aws:
	// prefix reservation, 50-tag limit, key/value length) surface to the
	// caller. Without this, the ParseTags callback inside
	// appsyncTagConfig swallows the error and the handler converts the
	// resulting empty tag slice into a generic "tags are required".
	if _, err := parseTags(req.Parameters); err != nil {
		return nil, err
	}
	return tags.HandleTag(ctx, req, appsyncTagConfig(store, req))
}

// UntagResource removes the specified tags from an AppSync API.
func (s *AppSyncService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}
	return tags.HandleUntag(ctx, req, appsyncTagConfig(store, req))
}

// ListTagsForResource lists all tags assigned to an AppSync API.
func (s *AppSyncService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
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
