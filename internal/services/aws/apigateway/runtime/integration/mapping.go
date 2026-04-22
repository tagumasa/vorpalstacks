package integration

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	requestParamPattern  = regexp.MustCompile(`^method\.request\.(header|querystring|path)\.(.+)$`)
	responseParamPattern = regexp.MustCompile(`^(?:integration\.response\.(?:header|body)|method\.response\.header)\.(.+)$`)
)

// applyRequestParameterMapping applies integration-level request parameter
// mappings. The format is "method.request.{source}.{name}" -> "{target}".
// Source can be "header", "querystring", or "path". The target specifies
// where the mapped value should be placed in the integration request.
func applyRequestParameterMapping(req *IntegrationRequest) *IntegrationRequest {
	if len(req.RequestParameters) == 0 {
		return req
	}

	mappedHeaders := make(map[string]string)
	mappedQueryParams := make(map[string]string)
	mappedPathParams := make(map[string]string)

	for mappingKey, mappingValue := range req.RequestParameters {
		if mappingValue == "" || mappingValue == "false" {
			continue
		}

		matches := requestParamPattern.FindStringSubmatch(mappingKey)
		if matches == nil {
			continue
		}

		sourceType := matches[1]
		sourceName := matches[2]
		targetName := mappingValue

		var sourceValue string
		switch sourceType {
		case "header":
			sourceValue = req.Headers[sourceName]
		case "querystring":
			sourceValue = req.QueryParams[sourceName]
		case "path":
			sourceValue = req.PathParams[sourceName]
		}

		if sourceValue == "" {
			continue
		}

		switch {
		case strings.HasPrefix(targetName, "integration.header."):
			mappedHeaders[strings.TrimPrefix(targetName, "integration.header.")] = sourceValue
		case strings.HasPrefix(targetName, "integration.querystring."):
			mappedQueryParams[strings.TrimPrefix(targetName, "integration.querystring.")] = sourceValue
		case strings.HasPrefix(targetName, "integration.path."):
			mappedPathParams[strings.TrimPrefix(targetName, "integration.path.")] = sourceValue
		default:
			mappedHeaders[targetName] = sourceValue
		}
	}

	if len(mappedHeaders) > 0 {
		newHeaders := make(map[string]string, len(req.Headers)+len(mappedHeaders))
		for k, v := range req.Headers {
			newHeaders[k] = v
		}
		for k, v := range mappedHeaders {
			newHeaders[k] = v
		}
		req.Headers = newHeaders
	}

	if len(mappedQueryParams) > 0 {
		newQueryParams := make(map[string]string, len(req.QueryParams)+len(mappedQueryParams))
		for k, v := range req.QueryParams {
			newQueryParams[k] = v
		}
		for k, v := range mappedQueryParams {
			newQueryParams[k] = v
		}
		req.QueryParams = newQueryParams
	}

	if len(mappedPathParams) > 0 {
		newPathParams := make(map[string]string, len(req.PathParams)+len(mappedPathParams))
		for k, v := range req.PathParams {
			newPathParams[k] = v
		}
		for k, v := range mappedPathParams {
			newPathParams[k] = v
		}
		req.PathParams = newPathParams
	}

	return req
}

// applyResponseParameterMapping applies integration response parameter mappings.
// The format is "integration.response.header.{name}" -> "method.response.header.{name}"
// or "integration.response.body.{path}" -> "method.response.header.{name}".
// Static values (those not matching a pattern) are set directly as header values.
func applyResponseParameterMapping(responseParams map[string]string, currentHeaders map[string]string, responseBody string) map[string]string {
	if len(responseParams) == 0 {
		return currentHeaders
	}

	result := make(map[string]string, len(currentHeaders)+len(responseParams))
	for k, v := range currentHeaders {
		result[k] = v
	}

	for mappingKey, mappingValue := range responseParams {
		if mappingValue == "" || mappingValue == "false" {
			delete(result, mappingKey)
			continue
		}

		matches := responseParamPattern.FindStringSubmatch(mappingKey)
		if matches == nil {
			continue
		}

		source := mappingKey
		if strings.HasPrefix(source, "integration.response.body.") {
			jsonPath := strings.TrimPrefix(source, "integration.response.body.")
			extracted := extractValueFromJSON(responseBody, jsonPath)
			if extracted != "" {
				result[mappingValue] = extracted
			}
		} else {
			sourceHeader := matches[1]
			if val, ok := result[sourceHeader]; ok {
				result[mappingValue] = val
			}
		}
	}

	return result
}

// extractValueFromJSON extracts a simple dotted-path value from a JSON string.
// This handles flat JSON paths like "foo.bar.baz" without full JSONPath support.
func extractValueFromJSON(body string, path string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(body), &obj); err != nil {
		return ""
	}

	parts := strings.Split(path, ".")
	current := obj
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return ""
			}
		default:
			return ""
		}
	}

	switch v := current.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return ""
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(data)
	}
}
