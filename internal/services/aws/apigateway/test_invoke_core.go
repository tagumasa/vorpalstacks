package apigateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	integration "vorpalstacks/internal/services/aws/apigateway/runtime/integration"
	"vorpalstacks/internal/store/aws/apigateway"
)

// TestInvokePayload carries the parsed request-simulation members shared by
// the test-invoke operations.
type TestInvokePayload struct {
	Body                string
	PathWithQueryString string
	Headers             map[string]string
	MultiValueHeaders   map[string][]string
	StageVariables      map[string]string
	AdditionalContext   map[string]string
}

// testInvokeMethodCore simulates an API method invocation: it loads the
// method, dispatches through its configured integration executor and wraps
// the outcome in the TestInvokeMethodStatus response shape.
func (s *APIGatewayService) testInvokeMethodCore(
	ctx context.Context,
	stores *apiGatewayStores,
	apiId, resourceId, httpMethod string,
	p *TestInvokePayload,
) (map[string]interface{}, error) {
	startTime := time.Now()

	if apiId == "" || resourceId == "" {
		return nil, NewBadRequestException("restApiId and resourceId are required")
	}

	if httpMethod == "" {
		return nil, NewBadRequestException("httpMethod is required")
	}

	method, err := stores.restApis.GetMethod(apiId, resourceId, httpMethod)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	body := p.Body
	pathWithQueryString := p.PathWithQueryString
	headers := p.Headers
	multiValueHeaders := p.MultiValueHeaders
	stageVariables := p.StageVariables

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
			URI:                  mi.Uri,
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
		if s.testInvokeFactory != nil {
			created, createErr := s.testInvokeFactory.CreateExecutor(mi.Type)
			if createErr != nil {
				executor = nil
				responseStatus = 502
				responseBody = fmt.Sprintf(`{"message": "Integration execution failed: %v"}`, createErr)
				logEntries = append(logEntries, fmt.Sprintf("Execution failed: %v", createErr))
			} else {
				executor = created
			}
		} else {
			// InitRuntimeServer has not run (no bus wired): fall back to the
			// bus-less executors so the operation remains testable.
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
		}

		if executor != nil {
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
		"status":  responseStatus,
		"body":    responseBody,
		"log":     logStr,
		"latency": time.Since(startTime).Milliseconds(),
	}
	if len(headers) > 0 {
		result["headers"] = headers
	}
	if len(multiValueHeaders) > 0 {
		result["multiValueHeaders"] = multiValueHeaders
	}
	if pathWithQueryString != "" {
		result["pathWithQueryString"] = pathWithQueryString
	}
	if len(stageVariables) > 0 {
		result["stageVariables"] = stageVariables
	}

	return result, nil
}

// testInvokeAuthorizerCore simulates an authorizer invocation against the
// stub test authorizer, which mirrors the runtime identity-source and
// Authorization-header reading rules.
func (s *APIGatewayService) testInvokeAuthorizerCore(
	stores *apiGatewayStores,
	apiId, authorizerId string,
	p *TestInvokePayload,
) (map[string]interface{}, error) {
	if apiId == "" {
		return nil, NewBadRequestException("restApiId is required")
	}

	if authorizerId == "" {
		return nil, NewBadRequestException("authorizerId is required")
	}

	authorizer, err := stores.restApis.GetAuthorizer(apiId, authorizerId)
	if err != nil {
		return nil, toApiGatewayError(err)
	}

	headers := p.Headers
	stageVariables := p.StageVariables
	multiValueHeaders := p.MultiValueHeaders
	additionalContext := p.AdditionalContext

	result := map[string]interface{}{
		"clientStatus": 200,
		"log":          "TestInvokeAuthorizer completed successfully",
		"latency":      1,
	}

	switch authorizer.Type {
	case "TOKEN":
		if authorizer.IdentitySource != "" {
			headerName := extractHeaderFromIdentitySource(authorizer.IdentitySource)
			// Only the identity-source header supplies the token. Falling
			// back to arbitrary header values would accept tokens from
			// headers the authorizer is not configured to read, diverging
			// from the runtime behaviour under test.
			token := headers[headerName]
			if token == "" {
				if vals := multiValueHeaders[headerName]; len(vals) > 0 {
					token = vals[0]
				}
			}
			if token != "" {
				result["principalId"] = "test-user"
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
		result["principalId"] = "test-user"
		result["authorization"] = map[string]interface{}{
			"principalId": []string{"test-user"},
		}
		result["policy"] = buildTestPolicy(authorizer, apiId)
	case "COGNITO_USER_POOLS":
		// Verify that an Authorization header is present. Full JWT
		// validation against the Cognito user pool is not performed in
		// this test invocation path.
		authHeader := headers["Authorization"]
		if authHeader == "" {
			authHeader = headers["authorization"]
		}
		if authHeader == "" {
			result["clientStatus"] = 403
			result["log"] = "Unauthorized: missing Authorization header for COGNITO_USER_POOLS authorizer"
			break
		}
		result["principalId"] = "test-user"
		result["authorization"] = map[string]interface{}{
			"principalId": []string{"test-user"},
		}
		result["policy"] = buildTestPolicy(authorizer, apiId)
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

func buildTestPolicy(authorizer *apigateway.Authorizer, apiId string) string {
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
