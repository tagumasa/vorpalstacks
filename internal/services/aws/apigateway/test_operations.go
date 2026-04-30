package apigateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	integration "vorpalstacks/internal/services/aws/apigateway/runtime/integration"
	store "vorpalstacks/internal/store/aws/apigateway"
)

// TestInvokeMethod simulates an API method invocation for testing purposes.
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

	responseStatus := 200
	var responseBody string
	var logEntries []string

	if method.MethodIntegration != nil {
		mi := method.MethodIntegration

		reqTemplates := make(map[string]string)
		for k, v := range mi.RequestTemplates {
			reqTemplates[k] = v
		}

		intResponses := make(map[string]*integration.IntegrationResponseConfig)
		for code, ir := range mi.IntegrationResponses {
			respHeaders := make(map[string]string)
			for k, v := range ir.ResponseParameters {
				respHeaders[k] = v
			}
			respTemplates := make(map[string]string)
			for k, v := range ir.ResponseTemplates {
				respTemplates[k] = v
			}
			intResponses[code] = &integration.IntegrationResponseConfig{
				StatusCode:        ir.StatusCode,
				SelectionPattern:  ir.SelectionPattern,
				ResponseHeaders:   respHeaders,
				ResponseTemplates: respTemplates,
			}
		}

		intReq := &integration.IntegrationRequest{
			Method:               httpMethod,
			Headers:              headers,
			Body:                 []byte(body),
			PathParams:           make(map[string]string),
			QueryParams:          make(map[string]string),
			Path:                 pathWithQueryString,
			StageVariables:       stageVariables,
			IntegrationType:      mi.Type,
			RequestTemplates:     reqTemplates,
			IntegrationResponses: intResponses,
			RestApiId:            apiId,
			StageName:            "test-invoke-stage",
		}

		var executor integration.Executor
		switch mi.Type {
		case "MOCK":
			executor = integration.NewMockExecutor()
		case "HTTP", "HTTP_PROXY":
			executor = integration.NewHTTPExecutor()
		case "AWS", "AWS_PROXY":
			executor = integration.NewAWSExecutor(nil, s.accountID, s.region)
		default:
			executor = integration.NewMockExecutor()
		}

		resp, execErr := executor.Execute(ctx, intReq)
		if execErr != nil {
			responseStatus = 502
			responseBody = fmt.Sprintf(`{"message": "Integration execution failed: %v"}`, execErr)
			logEntries = append(logEntries, fmt.Sprintf("Execution failed: %v", execErr))
		} else {
			responseStatus = resp.StatusCode
			responseBody = string(resp.Body)
			logEntries = append(logEntries, "Execution completed successfully")
		}
	} else {
		responseStatus = 502
		responseBody = `{"message": "No integration configured"}`
		logEntries = append(logEntries, "No integration configured")
	}

	logStr := "TestInvokeMethod completed successfully"
	if len(logEntries) > 0 {
		logStr = strings.Join(logEntries, "; ")
	}

	result := map[string]interface{}{
		"status": responseStatus,
		"body":   responseBody,
		"log":    logStr,
	}
	if len(headers) > 0 {
		result["headers"] = headers
	}
	if pathWithQueryString != "" {
		result["pathWithQueryString"] = pathWithQueryString
	}
	if len(stageVariables) > 0 {
		result["stageVariables"] = stageVariables
	}

	return result, nil
}

// TestInvokeAuthorizer simulates an authoriser invocation for testing purposes.
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
