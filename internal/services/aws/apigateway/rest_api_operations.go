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
	input := CreateRestApiInput{
		Name:                      request.GetStringParam(req.Parameters, "name"),
		Description:               request.GetStringParam(req.Parameters, "description"),
		Version:                   request.GetStringParam(req.Parameters, "version"),
		ApiKeySource:              request.GetStringParam(req.Parameters, "apiKeySource"),
		Policy:                    request.GetStringParam(req.Parameters, "policy"),
		DisableExecuteApiEndpoint: request.GetBoolParam(req.Parameters, "disableExecuteApiEndpoint"),
		SecurityPolicy:            request.GetStringParam(req.Parameters, "securityPolicy"),
		EndpointAccessMode:        request.GetStringParam(req.Parameters, "endpointAccessMode"),
		CloneFrom:                 request.GetStringParam(req.Parameters, "cloneFrom"),
	}

	if _, ok := req.Parameters["minimumCompressionSize"]; ok {
		v := int32(request.GetIntParam(req.Parameters, "minimumCompressionSize"))
		input.MinimumCompressionSize = &v
	}

	if binaryMediaTypes, ok := req.Parameters["binaryMediaTypes"].([]interface{}); ok {
		for _, t := range binaryMediaTypes {
			if ts, ok := t.(string); ok {
				input.BinaryMediaTypes = append(input.BinaryMediaTypes, ts)
			}
		}
	}

	if endpointConfig, ok := req.Parameters["endpointConfiguration"].(map[string]interface{}); ok {
		if types, ok := endpointConfig["types"].([]interface{}); ok {
			for _, t := range types {
				if ts, ok := t.(string); ok {
					input.EndpointTypes = append(input.EndpointTypes, ts)
				}
			}
		}
	}

	if tags, ok := req.Parameters["tags"].(map[string]interface{}); ok {
		input.Tags = tagutil.MapInterfaceToTags(tags)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, err := s.createRestApiCore(stores, input)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	return s.toRestApiResponse(created), nil
}

// GetRestApi retrieves a REST API from API Gateway.
func (s *APIGatewayService) GetRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	api, err := s.getRestApiCore(stores, getRestApiId(req))
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toRestApiResponse(api), nil
}

// DeleteRestApi deletes a REST API from API Gateway. This is a cascading
// operation that deletes all associated resources (stages, deployments,
// models, authorizers, etc.) embedded in the RestApi document. It also
// removes any base path mappings that reference the deleted API.
func (s *APIGatewayService) DeleteRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteRestApiCore(stores, getRestApiId(req)); err != nil {
		return nil, toApiGatewayError(err)
	}
	return response.EmptyResponse(), nil
}

// UpdateRestApi updates an existing REST API in API Gateway.
func (s *APIGatewayService) UpdateRestApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ops, err := parsePatchOperations(req.Parameters)
	if err != nil {
		return nil, err
	}
	api, err := s.updateRestApiCore(stores, getRestApiId(req), ops)
	if err != nil {
		return nil, toApiGatewayError(err)
	}
	return s.toRestApiResponse(api), nil
}

// GetRestApis lists all REST APIs in API Gateway.
func (s *APIGatewayService) GetRestApis(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	limit, err := ResolvePaginationLimit(req.Parameters)
	if err != nil {
		return nil, err
	}
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listRestApisCore(stores, limit, request.GetStringParam(req.Parameters, "position"))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, api := range result.Items {
		items = append(items, s.toRestApiResponse(api))
	}

	resp := map[string]interface{}{
		"item": items,
	}
	if result.IsTruncated {
		resp["position"] = result.NextMarker
	}

	return resp, nil
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
		"id":                        api.Id,
		"name":                      api.Name,
		"createdDate":               timeutils.FormatEpochSeconds(api.CreatedDate),
		"rootResourceId":            rootResourceId,
		"disableExecuteApiEndpoint": api.DisableExecuteApiEndpoint,
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
	if api.MinimumCompressionSize != nil {
		response["minimumCompressionSize"] = *api.MinimumCompressionSize
	}
	if api.Policy != "" {
		response["policy"] = api.Policy
	}
	if api.SecurityPolicy != "" {
		response["securityPolicy"] = api.SecurityPolicy
	}
	if api.EndpointAccessMode != "" {
		response["endpointAccessMode"] = api.EndpointAccessMode
	}
	if len(api.Warnings) > 0 {
		response["warnings"] = api.Warnings
	}
	if api.ApiStatus != "" {
		response["apiStatus"] = api.ApiStatus
	}
	if api.ApiStatusMessage != "" {
		response["apiStatusMessage"] = api.ApiStatusMessage
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
			"types": []string{"REGIONAL"},
		}
	}
	if len(api.Tags) > 0 {
		response["tags"] = tagutil.ToMap(api.Tags)
	}

	return response
}
