// Package integration provides API Gateway integration functionality for vorpalstacks.
package integration

import (
	"bytes"
	"context"
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

// NewHTTPExecutor creates a new HTTP executor instance.
func NewHTTPExecutor() *HTTPExecutor {
	return &HTTPExecutor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute sends an HTTP request to the specified URI and returns the response.
func (e *HTTPExecutor) Execute(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if req.URI == "" {
		return nil, &IntegrationError{
			Message:  "URI is required for HTTP integration",
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
		}
	}

	targetURL := req.URI
	if strings.Contains(targetURL, "{") {
		for key, value := range req.PathParams {
			targetURL = strings.ReplaceAll(targetURL, "{"+key+"}", url.PathEscape(value))
		}
	}

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

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return nil, &IntegrationError{
			Message:  "HTTP request failed: " + err.Error(),
			Type:     "BadGatewayException",
			HTTPCode: http.StatusBadGateway,
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
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return &IntegrationResponse{
		StatusCode:      resp.StatusCode,
		Headers:         headers,
		Body:            body,
		IsBase64Encoded: false,
	}, nil
}
