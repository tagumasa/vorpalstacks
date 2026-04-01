package apigateway

import (
	"context"
	"encoding/json"
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
	store "vorpalstacks/internal/store/aws/apigateway"
)

func (s *APIGatewayService) TestInvokeMethod(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId, resourceId := getApiIdAndResourceId(req)
	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	httpMethod := request.GetStringParam(req.Parameters, "httpMethod")
	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	method, err := stores.restApis.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, ErrNotFoundException
	}

	body := request.GetStringParam(req.Parameters, "body")
	pathWithQueryString := request.GetStringParam(req.Parameters, "pathWithQueryString")

	responseBody := ""
	responseStatus := 200

	if method.MethodIntegration != nil {
		switch method.MethodIntegration.Type {
		case "MOCK":
			responseBody = `{"message": "Mock integration response"}`
		case "AWS_PROXY", "AWS":
			if method.MethodIntegration.Uri != "" {
				responseBody = `{"message": "Test invocation simulated successfully"}`
			} else {
				responseBody = ""
				responseStatus = 502
			}
		case "HTTP", "HTTP_PROXY":
			responseBody = `{"message": "Test invocation simulated successfully"}`
		default:
			responseBody = `{"message": "Test invocation simulated"}`
		}
	} else {
		responseStatus = 502
		responseBody = `{"message": "No integration configured"}`
	}

	headers := make(map[string]string)
	if h, ok := req.Parameters["headers"].(map[string]interface{}); ok {
		for k, v := range h {
			if vs, ok := v.(string); ok {
				headers[k] = vs
			}
		}
	}

	stageVariables := make(map[string]string)
	if sv, ok := req.Parameters["stageVariables"].(map[string]interface{}); ok {
		for k, v := range sv {
			if vs, ok := v.(string); ok {
				stageVariables[k] = vs
			}
		}
	}

	result := map[string]interface{}{
		"status": responseStatus,
		"log":    "TestInvokeMethod completed successfully",
	}

	if len(headers) > 0 {
		result["headers"] = headers
	}
	if body != "" {
		result["body"] = body
	}
	if pathWithQueryString != "" {
		result["pathWithQueryString"] = pathWithQueryString
	}
	if len(stageVariables) > 0 {
		result["stageVariables"] = stageVariables
	}
	if responseBody != "" {
		result["body"] = responseBody
	}

	return result, nil
}

func (s *APIGatewayService) TestInvokeAuthorizer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	apiId := getRestApiId(req)
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	authorizerId := request.GetStringParam(req.Parameters, "authorizerId")
	if authorizerId == "" {
		return nil, NewBadRequestException("authorizerId is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	authorizer, err := stores.restApis.GetAuthorizer(apiId, authorizerId)
	if err != nil {
		return nil, ErrNotFoundException
	}

	headers := make(map[string]string)
	if h, ok := req.Parameters["headers"].(map[string]interface{}); ok {
		for k, v := range h {
			if vs, ok := v.(string); ok {
				headers[k] = vs
			}
		}
	}

	stageVariables := make(map[string]string)
	if sv, ok := req.Parameters["stageVariables"].(map[string]interface{}); ok {
		for k, v := range sv {
			if vs, ok := v.(string); ok {
				stageVariables[k] = vs
			}
		}
	}

	multiValueHeaders := make(map[string][]string)
	if mvh, ok := req.Parameters["multiValueHeaders"].(map[string]interface{}); ok {
		for k, v := range mvh {
			if arr, ok := v.([]interface{}); ok {
				for _, item := range arr {
					if vs, ok := item.(string); ok {
						multiValueHeaders[k] = append(multiValueHeaders[k], vs)
					}
				}
			}
		}
	}

	additionalContext := make(map[string]string)
	if ac, ok := req.Parameters["additionalContext"].(map[string]interface{}); ok {
		for k, v := range ac {
			if vs, ok := v.(string); ok {
				additionalContext[k] = vs
			}
		}
	}

	result := map[string]interface{}{
		"clientStatus": 200,
		"log":          "TestInvokeAuthorizer completed successfully",
		"latency":      1,
	}

	switch authorizer.Type {
	case "TOKEN":
		if authorizer.IdentitySource != "" {
			headerName := extractHeaderFromIdentitySource(authorizer.IdentitySource)
			token := headers[headerName]
			if token == "" {
				for _, vals := range multiValueHeaders {
					if len(vals) > 0 {
						token = vals[0]
						break
					}
				}
			}
			if token == "" {
				for _, v := range headers {
					token = v
					break
				}
			}
			if token != "" {
				result["authorization"] = map[string]interface{}{
					"principalId": []string{"test-user"},
				}
				result["policy"] = buildTestPolicy(authorizer, apiId)
			} else {
				result["clientStatus"] = 403
				result["log"] = "Unauthorized: token not found in identity source"
			}
		}
	case "REQUEST":
		result["authorization"] = map[string]interface{}{
			"principalId": []string{"test-user"},
		}
		result["policy"] = buildTestPolicy(authorizer, apiId)
	case "AWS_IAM":
		result["authorization"] = map[string]interface{}{
			"principalId": []string{"test-user"},
		}
		result["policy"] = buildTestPolicy(authorizer, apiId)
		result["claims"] = map[string]string{
			"sub":   "test-user-id",
			"email": "test@example.com",
		}
	default:
		result["clientStatus"] = 502
		result["log"] = "Unsupported authorizer type: " + authorizer.Type
	}

	if len(headers) > 0 {
		result["headers"] = headers
	}
	if len(stageVariables) > 0 {
		result["stageVariables"] = stageVariables
	}
	if len(additionalContext) > 0 {
		result["additionalContext"] = additionalContext
	}

	return result, nil
}

func buildTestPolicy(authorizer *store.Authorizer, apiId string) string {
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":   "Allow",
				"Action":   "execute-api:Invoke",
				"Resource": "arn:aws:execute-api:*:*:" + apiId + "/*",
			},
		},
	}
	b, _ := json.Marshal(policy)
	return string(b)
}

func extractHeaderFromIdentitySource(identitySource string) string {
	identitySource = strings.TrimSpace(identitySource)
	if strings.HasPrefix(identitySource, "method.request.header.") {
		return strings.TrimPrefix(identitySource, "method.request.header.")
	}
	return identitySource
}
