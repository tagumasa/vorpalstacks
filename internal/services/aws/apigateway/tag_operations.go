package apigateway

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/utils/aws/types"
)

// taggableResource holds the kind and identifiers extracted from a taggable
// resource ARN.
type taggableResource struct {
	kind string // "stage", "usageplan", "apikey", "domainname", "restapi"
	id1  string
	id2  string // stageName for stages only
}

// parseTaggableArn resolves a resource ARN into its kind and identifiers.
// This is the single source of truth for ARN-to-resource mapping in tag
// operations.
func parseTaggableArn(arnStr string) (taggableResource, error) {
	switch {
	case strings.Contains(arnStr, "/stages/"):
		apiId := extractResourceFromArn(arnStr, "/restapis/")
		stageName := extractResourceFromArn(arnStr, "/stages/")
		if apiId == "" || stageName == "" {
			return taggableResource{}, NewBadRequestException("invalid stage ARN")
		}
		return taggableResource{kind: "stage", id1: apiId, id2: stageName}, nil

	case strings.Contains(arnStr, "/usageplans/"):
		usagePlanId := extractResourceFromArn(arnStr, "/usageplans/")
		if usagePlanId == "" {
			return taggableResource{}, NewBadRequestException("invalid usage plan ARN")
		}
		return taggableResource{kind: "usageplan", id1: usagePlanId}, nil

	case strings.Contains(arnStr, "/apikeys/"):
		apiKeyId := extractResourceFromArn(arnStr, "/apikeys/")
		if apiKeyId == "" {
			return taggableResource{}, NewBadRequestException("invalid API key ARN")
		}
		return taggableResource{kind: "apikey", id1: apiKeyId}, nil

	case strings.Contains(arnStr, "/domainnames/"):
		domainName := extractResourceFromArn(arnStr, "/domainnames/")
		if domainName == "" {
			return taggableResource{}, NewBadRequestException("invalid domain name ARN")
		}
		return taggableResource{kind: "domainname", id1: domainName}, nil

	case strings.Contains(arnStr, "/restapis/"):
		apiId := extractResourceFromArn(arnStr, "/restapis/")
		if apiId == "" {
			return taggableResource{}, NewBadRequestException("invalid REST API ARN")
		}
		return taggableResource{kind: "restapi", id1: apiId}, nil

	default:
		return taggableResource{}, NewBadRequestException("resourceArn is required")
	}
}

// tagResource dispatches tag operations based on the resource ARN pattern.
func (s *APIGatewayService) tagResource(stores *apiGatewayStores, arnStr string, tagsMap map[string]string) error {
	r, err := parseTaggableArn(arnStr)
	if err != nil {
		return err
	}
	switch r.kind {
	case "stage":
		return stores.restApis.TagStage(r.id1, r.id2, tagsMap)
	case "usageplan":
		return stores.usage.TagUsagePlan(r.id1, tagsMap)
	case "apikey":
		return stores.usage.TagApiKey(r.id1, tagsMap)
	case "domainname":
		return stores.domains.TagDomainName(r.id1, tagsMap)
	case "restapi":
		return stores.restApis.Tag(r.id1, tagsMap)
	default:
		return NewBadRequestException("resourceArn is required")
	}
}

// untagResource dispatches untag operations based on the resource ARN pattern.
func (s *APIGatewayService) untagResource(stores *apiGatewayStores, arnStr string, tagKeys []string) error {
	r, err := parseTaggableArn(arnStr)
	if err != nil {
		return err
	}
	switch r.kind {
	case "stage":
		return stores.restApis.UntagStage(r.id1, r.id2, tagKeys)
	case "usageplan":
		return stores.usage.UntagUsagePlan(r.id1, tagKeys)
	case "apikey":
		return stores.usage.UntagApiKey(r.id1, tagKeys)
	case "domainname":
		return stores.domains.UntagDomainName(r.id1, tagKeys)
	case "restapi":
		return stores.restApis.Untag(r.id1, tagKeys)
	default:
		return NewBadRequestException("resourceArn is required")
	}
}

// getResourceTags dispatches get-tags operations based on the resource ARN
// pattern.
func (s *APIGatewayService) getResourceTags(stores *apiGatewayStores, arnStr string) ([]types.Tag, error) {
	r, err := parseTaggableArn(arnStr)
	if err != nil {
		return nil, err
	}
	switch r.kind {
	case "stage":
		return stores.restApis.GetStageTags(r.id1, r.id2)
	case "usageplan":
		return stores.usage.GetUsagePlanTags(r.id1)
	case "apikey":
		return stores.usage.GetApiKeyTags(r.id1)
	case "domainname":
		return stores.domains.GetDomainNameTags(r.id1)
	case "restapi":
		return stores.restApis.GetResourceTags(r.id1)
	default:
		return nil, NewBadRequestException("resourceArn is required")
	}
}

func extractResourceFromArn(arnStr, suffix string) string {
	idx := strings.Index(arnStr, suffix)
	if idx < 0 {
		return ""
	}
	rest := arnStr[idx+len(suffix):]
	if slashIdx := strings.Index(rest, "/"); slashIdx >= 0 {
		return rest[:slashIdx]
	}
	return rest
}

// apiGatewayMapTagError maps framework-level tag errors to API Gateway errors.
func apiGatewayMapTagError(err error) error {
	switch err.(type) {
	case *tagutil.MissingResourceError:
		return NewBadRequestException("resourceArn is required")
	case *tagutil.MissingTagsError:
		return NewBadRequestException("tags is required")
	case *tagutil.MissingTagKeysError:
		return NewBadRequestException("tagKeys is required")
	}
	return err
}

// apiGatewayTagConfig builds a TagHandlerConfig for the common tags framework.
// The HTTP-path TagResource/UntagResource/ListTagsForResource handlers
// delegate to tags.HandleTag/HandleUntag/HandleList using this config.
func (s *APIGatewayService) apiGatewayTagConfig(stores *apiGatewayStores, req *request.ParsedRequest) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.TagOperationConfig{
			ResourceParam:    "resourceArn",
			TagsParam:        "tags",
			TagKeysParam:     "tagKeys",
			TagKeyName:       "key",
			TagValueName:     "value",
			RequireTags:      true,
			RequireTagKeys:   true,
			RequireResource:  true,
			UseQueryFallback: true,
		},
		ParseTags: func(_ map[string]interface{}) []types.Tag {
			return tagutil.ParseTagsWithQueryFallback(req.Parameters, "tags")
		},
		ParseTagKeys: func(_ map[string]interface{}) []string {
			return tagutil.ParseTagKeysWithQueryFallback(req.Parameters, "tagKeys")
		},
		TagFunc: func(_ context.Context, resourceKey string, t []types.Tag) error {
			return s.tagResource(stores, resourceKey, tagutil.ToMap(t))
		},
		UntagFunc: func(_ context.Context, resourceKey string, tagKeys []string) error {
			return s.untagResource(stores, resourceKey, tagKeys)
		},
		ListFunc: func(_ context.Context, resourceKey string) ([]types.Tag, error) {
			return s.getResourceTags(stores, resourceKey)
		},
		FormatResponse: func(tags []types.Tag, _ string) (interface{}, error) {
			return map[string]interface{}{
				"tags": tagutil.ToMap(tags),
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return response.EmptyResponse(), nil
		},
		MapError: apiGatewayMapTagError,
	}
}

// TagResource adds tags to a resource in API Gateway.
func (s *APIGatewayService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.apiGatewayTagConfig(stores, req))
}

// UntagResource removes tags from a resource in API Gateway.
func (s *APIGatewayService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.apiGatewayTagConfig(stores, req))
}

// ListTagsForResource lists tags for a resource in API Gateway.
func (s *APIGatewayService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.apiGatewayTagConfig(stores, req))
}
