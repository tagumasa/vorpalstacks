// Package integration provides API Gateway integration functionality for vorpalstacks.
package integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPExecutor handles HTTP integrations for API Gateway.
type HTTPExecutor struct {
	client *http.Client
}

// NewHTTPExecutor creates a new HTTP executor instance with a shared
// http.Client for connection pooling.
func NewHTTPExecutor() *HTTPExecutor {
	return &HTTPExecutor{
		client: &http.Client{},
	}
}

// Execute sends an HTTP request to the specified URI and returns the response.
// For non-proxy HTTP integrations, request mapping templates, integration
// response selection, and response templates are applied.
func (e *HTTPExecutor) Execute(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if req.URI == "" {
		return nil, &IntegrationError{
			Message:  "URI is required for HTTP integration",
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
		}
	}

	isProxy := req.IntegrationType == "HTTP_PROXY"

	if !isProxy {
		req = applyRequestParameterMapping(req)
		processed, pErr := processRequestBody(req)
		if pErr != nil {
			return nil, pErr
		}
		req.Body = processed
	}

	targetURL := substituteStageVariables(req.URI, req.StageVariables)
	if strings.Contains(targetURL, "{") {
		for key, value := range req.PathParams {
			targetURL = strings.ReplaceAll(targetURL, "{"+key+"}", url.PathEscape(value))
		}
	}

	timeout := 29 * time.Second
	if req.TimeoutInMillis > 0 {
		timeout = time.Duration(req.TimeoutInMillis) * time.Millisecond
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, targetURL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, &IntegrationError{
			Message:  "Failed to create HTTP request: " + err.Error(),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}
	for key, values := range req.MultiValueHeaders {
		for _, v := range values {
			httpReq.Header.Add(key, v)
		}
	}

	resp, err := e.client.Do(httpReq)
	if err != nil {
		httpCode := http.StatusBadGateway
		errType := "BadGatewayException"
		if errors.Is(err, context.DeadlineExceeded) {
			httpCode = http.StatusGatewayTimeout
			errType = "GatewayTimeoutException"
		}
		return nil, &IntegrationError{
			Message:  "HTTP request failed: " + err.Error(),
			Type:     errType,
			HTTPCode: httpCode,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &IntegrationError{
			Message:  "Failed to read response body: " + err.Error(),
			Type:     "InternalServerError",
			HTTPCode: http.StatusInternalServerError,
		}
	}

	headers := make(map[string]string)
	multiValueHeaders := make(map[string][]string)
	for key, values := range resp.Header {
		if len(values) == 1 {
			headers[key] = values[0]
		}
		if len(values) > 1 {
			multiValueHeaders[key] = values
		}
	}

	// For non-proxy HTTP, apply integration response selection and templates.
	if !isProxy && req.IntegrationResponses != nil {
		// For HTTP integrations, AWS matches selection pattern against the
		// HTTP status code, not the response body.
		statusStr := fmt.Sprintf("%d", resp.StatusCode)
		respConfig := matchIntegrationResponse(req.IntegrationResponses, statusStr, resp.StatusCode)
		if respConfig != nil {
			if respConfig.ResponseTemplates != nil {
				contentType := selectResponseContentType(respConfig.ResponseTemplates, headers, req.Headers)
				if tmpl, ok := respConfig.ResponseTemplates[contentType]; ok && tmpl != "" {
					transformed, tErr := applyResponseTemplate(tmpl, body, req)
					if tErr != nil {
						return nil, &IntegrationError{
							Message:  fmt.Sprintf("Failed to apply response template: %v", tErr),
							Type:     "InternalServerError",
							HTTPCode: 500,
						}
					}
					body = transformed
				}
			}

			if respConfig.StatusCode != "" {
				resp.StatusCode, _ = parseStatusCode(respConfig.StatusCode)
			}

			headers = applyResponseParameterMapping(respConfig.ResponseHeaders, headers, string(body))
			if respConfig.ContentHandling != "" {
				body = applyContentHandlingResponse(body, respConfig.ContentHandling)
			}
		}
	}

	return &IntegrationResponse{
		StatusCode:        resp.StatusCode,
		Headers:           headers,
		MultiValueHeaders: multiValueHeaders,
		Body:              body,
		IsBase64Encoded:   false,
	}, nil
}
