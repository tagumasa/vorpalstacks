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

	"vorpalstacks/internal/common/defaults"
)

const maxRequestBodySize = 64 * 1024 * 1024 // 64 MiB

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

// GetRegion returns the request region, defaulting to defaults.DefaultRegion if unset.
func (r *ParsedRequest) GetRegion() string {
	if r.Region != "" {
		return r.Region
	}
	return defaults.DefaultRegion
}

// GetParam returns a parameter value, trying the original key first
// then the lower-first variant.
func (r *ParsedRequest) GetParam(key string) string {
	return GetParamLowerFirst(r.Parameters, key)
}

// ParseAWSRequest parses an AWS HTTP request into a ParsedRequest.
// Body parsing is driven by Content-Type; operation name and path
// parameters are extracted via registered REST service parsers.
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
		bodyBytes, err = io.ReadAll(io.LimitReader(r.Body, maxRequestBodySize))
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
		req.Region = defaults.DefaultRegion
	}
	req.AccessKeyID = ExtractAccessKeyIDFromAuth(authHeader)

	for _, p := range restParsers {
		if p.MatchPath(r.URL.Path) {
			p.ExtractPathParams(r, req.Parameters)
		}
	}

	return req, nil
}

// extractOperation determines the operation name from the request.
// Universal protocol headers (X-Amz-Target, Action) take precedence,
// followed by registered REST service parsers.
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

	for _, p := range restParsers {
		if op := p.ExtractOperation(r); op != "" {
			if op == "CreateDistribution" && r.Method == "POST" {
				if inspectCloudFrontBodyForTags(r, bodyBytes) {
					return "CreateDistributionWithTags"
				}
			}
			return op
		}
	}

	return ""
}

// inspectCloudFrontBodyForTags checks whether a CloudFront CreateDistribution
// request body contains a DistributionConfigWithTags wrapper element.
func inspectCloudFrontBodyForTags(r *http.Request, bodyBytes []byte) bool {
	if len(bodyBytes) > 0 {
		if bytes.Contains(bodyBytes, []byte("DistributionConfigWithTags")) {
			return true
		}
	}
	if r.Body != nil {
		buf, err := io.ReadAll(io.LimitReader(r.Body, maxRequestBodySize))
		if err == nil {
			r.Body = io.NopCloser(bytes.NewReader(buf))
			return bytes.Contains(buf, []byte("DistributionConfigWithTags"))
		}
	}
	return false
}

// convertCBORMapToStringMap recursively converts CBOR-decoded maps
// with interface{} keys to string-keyed maps.
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

// ExtractRESTOperation determines the operation for REST protocol services.
// Delegates to extractOperation which checks all registered REST parsers.
func ExtractRESTOperation(r *http.Request, bodyBytes []byte) string {
	return extractOperation(r, bodyBytes)
}

// ExtractRESTPathParams extracts path parameters for REST protocol services
// by iterating over all registered parsers and invoking those whose path
// predicate matches the request URL.
func ExtractRESTPathParams(r *http.Request, params map[string]interface{}) {
	for _, p := range restParsers {
		if p.MatchPath(r.URL.Path) {
			p.ExtractPathParams(r, params)
		}
	}
}

// ParseXMLBody parses XML body bytes into the given params map.
func ParseXMLBody(bodyBytes []byte, params map[string]interface{}) {
	parseXMLBody(bodyBytes, params)
}

// toString converts common types to their string representation.
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
