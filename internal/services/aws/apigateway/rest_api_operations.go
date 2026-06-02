package apigateway

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	store "vorpalstacks/internal/store/aws/apigateway"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/timeutils"
)

var apiGatewayArnRegex = regexp.MustCompile(`restapis/([^/]+)`)

func getPathParam(req *request.ParsedRequest, key string) string {
	if req.PathParams != nil {
		return req.PathParams[key]
	}
	return ""
}

// resolveBinaryMediaTypeToRemove determines which media type to remove based on
// the JSON Patch path. The path may be an index ("/binaryMediaTypes/2"), a
// value ("/binaryMediaTypes/image~1png"), or empty ("/binaryMediaTypes/-").
func resolveBinaryMediaTypeToRemove(path, value string, current []string) string {
	segment := strings.TrimPrefix(path, "/binaryMediaTypes/")
	segment = strings.ReplaceAll(segment, "~1", "/")
	segment = strings.ReplaceAll(segment, "~0", "~")
	if idx, err := strconv.Atoi(segment); err == nil && idx < len(current) {
		return current[idx]
	}
	if segment == "" || segment == "-" {
		return value
	}
	return segment
}

// removeString returns a new slice with all occurrences of target removed.
func removeString(slice []string, target string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != target {
			result = append(result, s)
		}
	}
	return result
}

func getRestApiId(req *request.ParsedRequest) string {
	apiId := request.GetStringParam(req.Parameters, "restApiId")
	if apiId == "" {
		apiId = getPathParam(req, "restApiId")
	}
	if apiId == "" {
		if arnStr := getPathParam(req, "resourceArn"); arnStr != "" {
			if matches := apiGatewayArnRegex.FindStringSubmatch(arnStr); len(matches) > 1 {
				apiId = matches[1]
			}
		}
	}
	if apiId == "" {
		if matches := apiGatewayArnRegex.FindStringSubmatch(request.GetStringParam(req.Parameters, "resourceArn")); len(matches) > 1 {
			apiId = matches[1]
		}
	}
	return apiId
}

// CreateRestApi creates a new REST API in API Gateway.
func (s *APIGatewayService) CreateRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "name")
	if name == "" {
		return nil, NewBadRequestException("name is required")
	}

	api := &store.RestApi{
		Name:                   name,
		Description:            request.GetStringParam(req.Parameters, "description"),
		Version:                request.GetStringParam(req.Parameters, "version"),
		ApiKeySource:           request.GetStringParam(req.Parameters, "apiKeySource"),
		Policy:                 request.GetStringParam(req.Parameters, "policy"),
		MinimumCompressionSize: int32(request.GetIntParam(req.Parameters, "minimumCompressionSize")),
	}

	if binaryMediaTypes, ok := req.Parameters["binaryMediaTypes"].([]interface{}); ok {
		for _, t := range binaryMediaTypes {
			if ts, ok := t.(string); ok {
				api.BinaryMediaTypes = append(api.BinaryMediaTypes, ts)
			}
		}
	}

	if endpointConfig, ok := req.Parameters["endpointConfiguration"].(map[string]interface{}); ok {
		api.EndpointConfiguration = &store.EndpointConfiguration{}
		if types, ok := endpointConfig["types"].([]interface{}); ok {
			for _, t := range types {
				if ts, ok := t.(string); ok {
					api.EndpointConfiguration.Types = append(api.EndpointConfiguration.Types, ts)
				}
			}
		}
	}
	if api.EndpointConfiguration == nil || len(api.EndpointConfiguration.Types) == 0 {
		api.EndpointConfiguration = &store.EndpointConfiguration{
			Types: []string{"REGIONAL"},
		}
	}

	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		api.Tags = tagutil.MapInterfaceToTags(tags)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := stores.restApis.Create(api)
	if err != nil {
		return nil, err
	}

	return s.toRestApiResponse(created), nil
}

// GetRestApi retrieves a REST API from API Gateway.
func (s *APIGatewayService) GetRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	api, err := stores.restApis.Get(apiId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toRestApiResponse(api), nil
}

// DeleteRestApi deletes a REST API from API Gateway.
func (s *APIGatewayService) DeleteRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.restApis.Delete(apiId); err != nil {
		return nil, toApiGatewayError(err)
	}

	return response.EmptyResponse(), nil
}

// UpdateRestApi updates an existing REST API in API Gateway.
func (s *APIGatewayService) UpdateRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	stores.keyLocker.Lock(apiId)
	defer stores.keyLocker.Unlock(apiId)

	api, err := stores.restApis.Get(apiId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	for _, po := range parsePatchOperations(req.Parameters) {
		switch {
		case po.Path == "/name":
			api.Name = po.Value
		case po.Path == "/description":
			api.Description = po.Value
		case po.Path == "/version":
			api.Version = po.Value
		case po.Path == "/apiKeySource":
			api.ApiKeySource = po.Value
		case po.Path == "/policy":
			api.Policy = po.Value
		case po.Path == "/minimumCompressionSize":
			v, err := parseInt32(po.Value)
			if err != nil {
				return nil, NewBadRequestException("invalid minimumCompressionSize: not a number")
			}
			api.MinimumCompressionSize = v
		case strings.HasPrefix(po.Path, "/binaryMediaTypes"):
			if po.Op == "add" {
				if !containsAny(api.BinaryMediaTypes, po.Value) {
					api.BinaryMediaTypes = append(api.BinaryMediaTypes, po.Value)
				}
			} else if po.Op == "remove" {
				target := resolveBinaryMediaTypeToRemove(po.Path, po.Value, api.BinaryMediaTypes)
				api.BinaryMediaTypes = removeString(api.BinaryMediaTypes, target)
			}
		}
	}

	if err := stores.restApis.Update(api); err != nil {
		return nil, err
	}

	return s.toRestApiResponse(api), nil
}

// GetRestApis lists all REST APIs in API Gateway.
func (s *APIGatewayService) GetRestApis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxItems := request.GetIntParam(req.Parameters, "limit")
	if maxItems <= 0 {
		maxItems = 25
	}
	marker := request.GetStringParam(req.Parameters, "position")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := stores.restApis.List(common.ListOptions{
		Marker:   marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, api := range result.Items {
		items = append(items, s.toRestApiResponse(api))
	}

	response := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		response["position"] = result.NextMarker
	}

	return response, nil
}

func (s *APIGatewayService) toRestApiResponse(api *store.RestApi) map[string]interface{} {
	var rootResourceId string
	for _, res := range api.Resources {
		if res.Path == "/" {
			rootResourceId = res.Id
			break
		}
	}

	response := map[string]interface{}{
		"id":              api.Id,
		"name":            api.Name,
		"createdDate":     timeutils.FormatEpochSeconds(api.CreatedDate),
		"rootResourceId":  rootResourceId,
	}

	if api.Description != "" {
		response["description"] = api.Description
	}
	if api.Version != "" {
		response["version"] = api.Version
	}
	if api.ApiKeySource != "" {
		response["apiKeySource"] = api.ApiKeySource
	}
	if api.MinimumCompressionSize > 0 {
		response["minimumCompressionSize"] = api.MinimumCompressionSize
	}
	if api.Policy != "" {
		response["policy"] = api.Policy
	}
	if len(api.BinaryMediaTypes) > 0 {
		response["binaryMediaTypes"] = api.BinaryMediaTypes
	}
	if api.EndpointConfiguration != nil {
		response["endpointConfiguration"] = map[string]interface{}{
			"types": api.EndpointConfiguration.Types,
		}
	} else {
		response["endpointConfiguration"] = map[string]interface{}{
			"types": []string{"EDGE"},
		}
	}
	if len(api.Tags) > 0 {
		response["tags"] = tagutil.ToMap(api.Tags)
	}

	return response
}
