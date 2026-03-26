// Package request provides HTTP request parsing and context handling for vorpalstacks.
package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fxamacker/cbor/v2"
)

// ParsedRequest represents a parsed AWS HTTP request.
type ParsedRequest struct {
	Operation   string
	Parameters  map[string]interface{}
	Headers     http.Header
	QueryParams url.Values
	PathParams  map[string]string
	Body        []byte
	Region      string
	AccessKeyID string
}

// GetRegion returns the request region, defaulting to DefaultRegion if not set.
func (r *ParsedRequest) GetRegion() string {
	if r.Region != "" {
		return r.Region
	}
	return DefaultRegion
}

// ParseAWSRequest parses an AWS HTTP request into a ParsedRequest.
func ParseAWSRequest(r *http.Request) (*ParsedRequest, error) {
	req := &ParsedRequest{
		Headers:     r.Header,
		QueryParams: r.URL.Query(),
		PathParams:  make(map[string]string),
		Parameters:  make(map[string]interface{}),
	}

	var bodyBytes []byte
	if r.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}
	req.Body = bodyBytes

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "application/x-amz-json") {
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &req.Parameters); err != nil {
				return nil, err
			}
		}
	} else if strings.Contains(contentType, "application/cbor") {
		if len(bodyBytes) > 0 {
			var rawParams interface{}
			if err := cbor.Unmarshal(bodyBytes, &rawParams); err != nil {
				return nil, err
			}
			if m, ok := convertCBORMapToStringMap(rawParams).(map[string]interface{}); ok {
				req.Parameters = m
			} else {
				return nil, fmt.Errorf("CBOR body is not a map: %T", rawParams)
			}
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if len(bodyBytes) > 0 {
			values, err := url.ParseQuery(string(bodyBytes))
			if err == nil {
				for k, v := range values {
					if len(v) == 1 {
						req.Parameters[k] = v[0]
					} else {
						req.Parameters[k] = v
					}
				}
			}
		}
	} else if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") || strings.HasPrefix(r.URL.Path, "/2013-04-01/") || strings.HasPrefix(r.URL.Path, "/2020-05-31/") {
		if len(bodyBytes) > 0 {
			parseXMLBody(bodyBytes, req.Parameters)
		}
	} else if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &req.Parameters); err == nil {
		} else if values, err := url.ParseQuery(string(bodyBytes)); err == nil {
			for k, v := range values {
				if len(v) == 1 {
					req.Parameters[k] = v[0]
				} else {
					req.Parameters[k] = v
				}
			}
		}
	}

	for k, v := range req.QueryParams {
		if len(v) == 1 {
			req.Parameters[k] = v[0]
		} else {
			req.Parameters[k] = v
		}
	}

	req.Operation = extractOperation(r, bodyBytes)

	authHeader := r.Header.Get("Authorization")
	req.Region = ExtractRegionFromAuth(authHeader)
	if req.Region == "" {
		req.Region = DefaultRegion
	}
	req.AccessKeyID = ExtractAccessKeyIDFromAuth(authHeader)

	if strings.HasPrefix(r.URL.Path, "/20") {
		extractLambdaPathParams(r.URL.Path, req.Parameters)
	}

	if isApiGatewayPath(r.URL.Path) {
		extractApiGatewayPathParams(r.URL.Path, req.Parameters)
	}

	if strings.HasPrefix(r.URL.Path, "/schedule-groups") || strings.HasPrefix(r.URL.Path, "/schedules") {
		extractSchedulerPathParams(r.URL.Path, r.Method, req.Parameters)
	}

	if strings.HasPrefix(r.URL.Path, "/v2/email/") {
		extractSESv2PathParams(r.URL.Path, req.Parameters)
	}

	if strings.HasPrefix(r.URL.Path, "/2013-04-01/") {
		extractRoute53PathParams(r.URL.Path, req.Parameters)
	}

	if strings.HasPrefix(r.URL.Path, "/2020-05-31/") {
		extractCloudFrontPathParams(r.URL.Path, req.Parameters)
	}

	return req, nil
}

func extractOperation(r *http.Request, bodyBytes []byte) string {
	xAmzTarget := r.Header.Get("X-Amz-Target")
	if xAmzTarget != "" {
		parts := strings.Split(xAmzTarget, ".")
		if len(parts) > 1 {
			return parts[len(parts)-1]
		}
		return xAmzTarget
	}

	if strings.HasPrefix(r.URL.Path, "/service/") {
		parts := strings.Split(r.URL.Path, "/")
		for i, p := range parts {
			if p == "operation" && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}

	action := r.URL.Query().Get("Action")
	if action != "" {
		return action
	}

	if len(bodyBytes) > 0 {
		values, err := url.ParseQuery(string(bodyBytes))
		if err == nil {
			if action := values.Get("Action"); action != "" {
				return action
			}
		}
	}

	if op := extractLambdaOperation(r); op != "" {
		return op
	}

	if op := extractApiGatewayOperation(r); op != "" {
		return op
	}

	if op := extractSchedulerOperation(r); op != "" {
		return op
	}

	if op := extractSESv2Operation(r); op != "" {
		return op
	}

	if op := extractRoute53Operation(r); op != "" {
		return op
	}

	if op := extractCloudFrontOperation(r); op != "" {
		if op == "CreateDistribution" && r.Method == "POST" && r.Body != nil {
			buf, err := io.ReadAll(r.Body)
			if err == nil {
				r.Body = io.NopCloser(bytes.NewReader(buf))
				if bytes.Contains(buf, []byte("DistributionConfigWithTags")) {
					return "CreateDistributionWithTags"
				}
			}
		}
		return op
	}

	return ""
}

func convertCBORMapToStringMap(v interface{}) interface{} {
	switch val := v.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			keyStr, ok := k.(string)
			if !ok {
				keyStr = toString(k)
			}
			result[keyStr] = convertCBORMapToStringMap(v)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = convertCBORMapToStringMap(item)
		}
		return result
	default:
		return v
	}
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		return ""
	}
}
