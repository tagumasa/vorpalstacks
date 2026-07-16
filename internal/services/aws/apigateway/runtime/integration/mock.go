// Package integration provides API Gateway integration functionality for vorpalstacks.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vorpalstacks/pkg/vtl"
)

const (
	defaultMockStatusCode = 200
	defaultMockBody       = `{"message": "Mock integration response"}`
)

// MockExecutor handles mock integrations for API Gateway.
//
// In AWS API Gateway, a mock integration returns a configured response without
// calling any backend. The response body is defined by the integration request
// templates (VTL), and the status code / headers / response templates come from
// the integration responses configuration.
type MockExecutor struct{}

// NewMockExecutor creates a new mock executor instance.
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{}
}

// Execute returns a mock response, using the IntegrationRequest configuration
// to determine status code, headers, and body.
//
// Resolution order:
//   - Body:   req.RequestTemplates["application/json"] (VTL template), or
//     req.RequestTemplates["*/*"] as fallback, or default JSON.
//   - Status: matched integration response StatusCode, or default 200.
//   - Headers: matched integration response ResponseHeaders, or default Content-Type.
//   - Response template: if the matched integration response has a response
//     template, it is applied to the body after initial template expansion.
func (e *MockExecutor) Execute(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	bodyTemplate := req.RequestTemplates["application/json"]
	if bodyTemplate == "" {
		bodyTemplate = req.RequestTemplates["*/*"]
	}

	body, err := e.renderBody(bodyTemplate, req)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Failed to render mock body template: %v", err),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	statusCode := defaultMockStatusCode
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if req.IntegrationResponses != nil {
		respConfig := matchIntegrationResponse(req.IntegrationResponses, string(body), statusCode)
		if respConfig != nil {
			if respConfig.StatusCode != "" {
				if parsed, ok := parseStatusCode(respConfig.StatusCode); ok {
					statusCode = parsed
				}
			}

			if respConfig.ResponseHeaders != nil {
				for k, v := range respConfig.ResponseHeaders {
					headers[k] = v
				}
			}

			if respConfig.ResponseTemplates != nil {
				if tmpl, ok := respConfig.ResponseTemplates["application/json"]; ok && tmpl != "" {
					transformed, err := applyResponseTemplate(tmpl, body, req)
					if err != nil {
						return nil, &IntegrationError{
							Message:  fmt.Sprintf("Failed to apply mock response template: %v", err),
							Type:     "InternalServerError",
							HTTPCode: http.StatusInternalServerError,
						}
					}
					body = transformed
				}
			}
		}
	}

	return &IntegrationResponse{
		StatusCode:      statusCode,
		Headers:         headers,
		Body:            body,
		IsBase64Encoded: false,
	}, nil
}

func (e *MockExecutor) renderBody(tmpl string, req *IntegrationRequest) ([]byte, error) {
	if tmpl == "" {
		return []byte(defaultMockBody), nil
	}

	engine := vtl.NewEngine()

	engine.GetContext().Body = string(req.Body)
	if len(req.Body) > 0 {
		var bodyObj interface{}
		if err := json.Unmarshal(req.Body, &bodyObj); err == nil {
			engine.GetContext().JSONBody = bodyObj
		}
	}

	params := make(map[string]string)
	for k, v := range req.PathParams {
		params[k] = v
	}
	for k, v := range req.QueryParams {
		params[k] = v
	}
	for k, v := range req.Headers {
		params[k] = v
	}
	engine.GetContext().Params = params

	engine.GetContext().Context = map[string]interface{}{
		"stage":     req.StageName,
		"apiId":     req.RestApiId,
		"requestId": fmt.Sprintf("%x", time.Now().UnixNano()),
		"path":      req.Path,
		"method":    req.Method,
	}

	engine.GetContext().StageVars = req.StageVariables

	result, err := engine.Transform(tmpl)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func parseStatusCode(s string) (int, bool) {
	var code int
	if _, err := fmt.Sscanf(s, "%d", &code); err != nil {
		return 0, false
	}
	if code < 100 || code > 599 {
		return 0, false
	}
	return code, true
}
