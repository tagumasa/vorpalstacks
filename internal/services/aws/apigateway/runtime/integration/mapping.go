package integration

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	requestParamPattern = regexp.MustCompile(`^method\.request\.(header|querystring|path)\.(.+)$`)
)

// processRequestBody applies the VTL request template for the content type
// if one is available. If no template matches, the passthrough behaviour
// determines whether to reject the request or pass the body through unchanged.
//
// AWS passthrough behaviour semantics:
//   - NEVER: a matching template is always required; reject otherwise.
//   - WHEN_NO_TEMPLATES: pass through only when no templates are defined;
//     if templates exist but none match the content type, reject.
//   - WHEN_NO_MATCH (default): pass through when no template matches.
func processRequestBody(req *IntegrationRequest) ([]byte, *IntegrationError) {
	if len(req.Body) == 0 {
		return req.Body, nil
	}

	req.Body = applyContentHandlingRequest(req.Body, req.IntegrationContentHandling)

	// Per AWS docs: "You don't provide a mapping template for such a
	// conversion [CONVERT_TO_BINARY]." The base64-decoded body is passed
	// through directly without template processing.
	if req.IntegrationContentHandling == "CONVERT_TO_BINARY" {
		return req.Body, nil
	}

	contentType := req.Headers["Content-Type"]
	if contentType == "" {
		contentType = "application/json"
	}

	tmpl := req.RequestTemplates[contentType]
	if tmpl == "" {
		tmpl = req.RequestTemplates["application/json"]
	}

	if tmpl != "" {
		transformed, err := applyRequestTemplate(tmpl, req)
		if err != nil {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Failed to apply request template: %v", err),
				Type:     "InternalServerError",
				HTTPCode: 500,
			}
		}
		return transformed, nil
	}

	switch req.PassthroughBehavior {
	case "NEVER":
		return nil, &IntegrationError{
			Message:  "No matching template found and passthrough behaviour is NEVER",
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	case "WHEN_NO_TEMPLATES":
		if len(req.RequestTemplates) > 0 {
			return nil, &IntegrationError{
				Message:  "No matching template found and passthrough behaviour is WHEN_NO_TEMPLATES",
				Type:     "BadRequestException",
				HTTPCode: 400,
			}
		}
	}

	return req.Body, nil
}

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
		case strings.HasPrefix(targetName, "integration.request.header."):
			mappedHeaders[strings.TrimPrefix(targetName, "integration.request.header.")] = sourceValue
		case strings.HasPrefix(targetName, "integration.request.querystring."):
			mappedQueryParams[strings.TrimPrefix(targetName, "integration.request.querystring.")] = sourceValue
		case strings.HasPrefix(targetName, "integration.request.path."):
			mappedPathParams[strings.TrimPrefix(targetName, "integration.request.path.")] = sourceValue
		case strings.HasPrefix(targetName, "integration.header."):
			mappedHeaders[strings.TrimPrefix(targetName, "integration.header.")] = sourceValue
		case strings.HasPrefix(targetName, "integration.querystring."):
			mappedQueryParams[strings.TrimPrefix(targetName, "integration.querystring.")] = sourceValue
		case strings.HasPrefix(targetName, "integration.path."):
			mappedPathParams[strings.TrimPrefix(targetName, "integration.path.")] = sourceValue
		default:
			continue
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
// Per AWS spec (Smithy model), the KEY is a method response header parameter
// name (method.response.header.{name}) and the VALUE is the source:
//   - integration.response.header.{name} — from backend response header
//   - integration.response.body.{jsonPath} — from backend response body
//   - 'static value' (single-quoted) — literal string
//
// A value of "false" or "" removes the destination header.
func applyResponseParameterMapping(responseParams map[string]string, currentHeaders map[string]string, responseBody string) map[string]string {
	if len(responseParams) == 0 {
		return currentHeaders
	}

	result := make(map[string]string, len(currentHeaders)+len(responseParams))
	for k, v := range currentHeaders {
		result[k] = v
	}

	for mappingKey, mappingValue := range responseParams {
		// Extract the destination header name from the key.
		destHeader := extractResponseHeaderName(mappingKey)
		if destHeader == "" {
			continue
		}

		if mappingValue == "" || mappingValue == "false" {
			delete(result, destHeader)
			continue
		}

		switch {
		case strings.HasPrefix(mappingValue, "integration.response.header."):
			sourceHeader := strings.TrimPrefix(mappingValue, "integration.response.header.")
			if val, ok := result[sourceHeader]; ok {
				result[destHeader] = val
			}
		case strings.HasPrefix(mappingValue, "integration.response.body."):
			jsonPath := strings.TrimPrefix(mappingValue, "integration.response.body.")
			extracted := extractValueFromJSON(responseBody, jsonPath)
			if extracted != "" {
				result[destHeader] = extracted
			}
		case len(mappingValue) >= 2 && mappingValue[0] == '\'' && mappingValue[len(mappingValue)-1] == '\'':
			// Static value enclosed in single quotes (e.g., 'application/json')
			result[destHeader] = mappingValue[1 : len(mappingValue)-1]
		default:
			result[destHeader] = mappingValue
		}
	}

	return result
}

// extractResponseHeaderName extracts the header name from a response parameter
// key. Accepts both method.response.header.{name} and
// integration.response.header.{name} patterns.
func extractResponseHeaderName(key string) string {
	const methodRespPrefix = "method.response.header."
	const integrationRespPrefix = "integration.response.header."
	if strings.HasPrefix(key, methodRespPrefix) {
		return strings.TrimPrefix(key, methodRespPrefix)
	}
	if strings.HasPrefix(key, integrationRespPrefix) {
		return strings.TrimPrefix(key, integrationRespPrefix)
	}
	return ""
}

// extractValueFromJSON extracts a simple dotted-path value from a JSON string.
// Supports flat JSON paths like "foo.bar.baz" and array indices like
// "items[0].name" or "data[2]" without full JSONPath support.
func extractValueFromJSON(body string, path string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(body), &obj); err != nil {
		return ""
	}

	parts := strings.Split(path, ".")
	current := obj
	for _, part := range parts {
		// Split array index suffix: "items[0]" → key="items", idx=0
		key, idx, hasIdx := splitArrayIndex(part)
		if key != "" {
			m, ok := current.(map[string]interface{})
			if !ok {
				return ""
			}
			current, ok = m[key]
			if !ok {
				return ""
			}
		}
		if hasIdx {
			arr, ok := current.([]interface{})
			if !ok || idx < 0 || idx >= len(arr) {
				return ""
			}
			current = arr[idx]
		}
	}

	switch v := current.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
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

// splitArrayIndex parses a path segment that may contain an array index suffix.
// Examples: "items" → ("items", 0, false); "items[0]" → ("items", 0, true);
// "[2]" → ("", 2, true).
func splitArrayIndex(part string) (key string, idx int, hasIdx bool) {
	bracketStart := strings.IndexByte(part, '[')
	if bracketStart < 0 {
		return part, 0, false
	}
	key = part[:bracketStart]
	rest := part[bracketStart+1:]
	bracketEnd := strings.IndexByte(rest, ']')
	if bracketEnd < 0 {
		return part, 0, false
	}
	n, err := strconv.Atoi(rest[:bracketEnd])
	if err != nil || n < 0 {
		return part, 0, false
	}
	return key, n, true
}

// isBinaryContentType checks whether the given Content-Type matches any of the
// configured binary media types on the REST API. Supports exact match,
// type/* wildcards (e.g., image/*), and the global */* wildcard.
func isBinaryContentType(contentType string, binaryMediaTypes []string) bool {
	if len(binaryMediaTypes) == 0 || contentType == "" {
		return false
	}
	for _, bmt := range binaryMediaTypes {
		if bmt == "*/*" {
			return true
		}
		// Handle type/* wildcard (e.g., image/* matches image/png)
		if strings.HasSuffix(bmt, "/*") {
			prefix := strings.TrimSuffix(bmt, "/*")
			if strings.HasPrefix(contentType, prefix+"/") {
				return true
			}
			continue
		}
		if strings.HasPrefix(contentType, bmt) {
			return true
		}
	}
	return false
}

// applyContentHandlingRequest applies the integration-level ContentHandling
// setting to the request body. CONVERT_TO_BINARY base64-decodes the body so
// that base64-encoded binary payloads are restored to their binary form before
// being sent to the backend.
func applyContentHandlingRequest(body []byte, handling string) []byte {
	if handling == "CONVERT_TO_BINARY" && len(body) > 0 {
		decoded, err := base64.StdEncoding.DecodeString(string(body))
		if err == nil {
			return decoded
		}
	}
	return body
}

// applyContentHandlingResponse applies the response-level ContentHandling
// setting. CONVERT_TO_TEXT base64-encodes a binary response body so that it
// can be consumed by mapping templates or returned as text.
func applyContentHandlingResponse(body []byte, handling string) []byte {
	if handling == "CONVERT_TO_TEXT" && len(body) > 0 {
		return []byte(base64.StdEncoding.EncodeToString(body))
	}
	return body
}

// selectResponseContentType determines the Content-Type to use for
// integration response template selection. Priority:
//  1. Backend response Content-Type header
//  2. Client Accept header
//  3. application/json (AWS default)
func selectResponseContentType(templates map[string]string, responseHeaders, requestHeaders map[string]string) string {
	if ct := responseHeaders["Content-Type"]; ct != "" {
		if _, ok := templates[ct]; ok {
			return ct
		}
	}
	if accept := requestHeaders["Accept"]; accept != "" {
		if _, ok := templates[accept]; ok {
			return accept
		}
	}
	return "application/json"
}
