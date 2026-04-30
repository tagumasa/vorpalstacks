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
		return nil, ErrNotFoundException
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
		return nil, ErrNotFoundException
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
	api, err := stores.restApis.Get(apiId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	patchOps, ok := req.Parameters["patchOperations"].([]interface{})
	if ok {
		for _, op := range patchOps {
			if opMap, ok := op.(map[string]interface{}); ok {
				path := ""
				value := ""
				if p, ok := opMap["path"].(string); ok {
					path = p
				}
				if v, ok := opMap["value"].(string); ok {
					value = v
				}

				switch path {
				case "/name":
					api.Name = value
				case "/description":
					api.Description = value
				case "/version":
					api.Version = value
				case "/apiKeySource":
					api.ApiKeySource = value
				case "/policy":
					api.Policy = value
				case "/minimumCompressionSize":
					api.MinimumCompressionSize = parseInt32(value)
				}
				if strings.HasPrefix(path, "/binaryMediaTypes") {
					opName := ""
					if o, ok := opMap["op"].(string); ok {
						opName = o
					}
					if opName == "add" {
						if path == "/binaryMediaTypes" || path == "/binaryMediaTypes/-" {
							if !containsMediaType(api.BinaryMediaTypes, value) {
								api.BinaryMediaTypes = append(api.BinaryMediaTypes, value)
							}
						} else if strings.HasPrefix(path, "/binaryMediaTypes/") {
							if !containsMediaType(api.BinaryMediaTypes, value) {
								api.BinaryMediaTypes = append(api.BinaryMediaTypes, value)
							}
						}
					} else if opName == "remove" {
						mediaType := strings.TrimPrefix(path, "/binaryMediaTypes/")
						mediaType = strings.ReplaceAll(mediaType, "~1", "/")
						mediaType = strings.ReplaceAll(mediaType, "~0", "~")
						target := mediaType
						idx, err := strconv.Atoi(mediaType)
						if err == nil && idx < len(api.BinaryMediaTypes) {
							target = api.BinaryMediaTypes[idx]
						} else if mediaType == "" {
							target = value
						}
						newTypes := []string{}
						for _, t := range api.BinaryMediaTypes {
							if t != target {
								newTypes = append(newTypes, t)
							}
						}
						api.BinaryMediaTypes = newTypes
					}
				}
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
	response := map[string]interface{}{
		"id":          api.Id,
		"name":        api.Name,
		"createdDate": timeutils.FormatEpochSeconds(api.CreatedDate),
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

func containsMediaType(types []string, t string) bool {
	for _, v := range types {
		if v == t {
			return true
		}
	}
	return false
}
