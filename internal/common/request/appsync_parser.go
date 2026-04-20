package request

// AppSync request parser — dispatches to versioned and resource-specific sub-parsers.
// AWS AppSync uses REST-JSON 1.0 protocol with mixed HTTP methods across 74 control-plane
// operations plus the GraphQL execution data-plane endpoint.
// Path matching is case-sensitive as required by the AWS SDK wire format.

import (
	"net/http"
	"strings"
)

// isAppSyncPath reports whether the request path belongs to the AppSync service.
// Covers both v1 (/v1/apis, /v1/domainnames, /v1/tags, etc.) and v2 (/v2/apis) paths,
// as well as data-plane introspection and evaluation endpoints.
func isAppSyncPath(path string) bool {
	return strings.HasPrefix(path, "/v1/apis") || strings.HasPrefix(path, "/v2/apis") ||
		strings.HasPrefix(path, "/v1/domainnames") || strings.HasPrefix(path, "/v1/tags") ||
		strings.HasPrefix(path, "/v1/mergedApis") || strings.HasPrefix(path, "/v1/sourceApis") ||
		strings.HasPrefix(path, "/v1/datasources/introspections") || strings.HasPrefix(path, "/v1/dataplane-")
}

// extractAppSyncOperation determines the AppSync operation from the request path and HTTP method.
// Handles all 74 SDK operations plus the GraphQL execution sentinel.
// Delegates to versioned/resource-specific sub-parsers for each path prefix group.
func extractAppSyncOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if !isAppSyncPath(path) {
		return ""
	}

	switch {
	case strings.HasPrefix(path, "/v2/apis"):
		return extractAppSyncV2Operation(path, method)
	case strings.HasPrefix(path, "/v1/apis"):
		return extractAppSyncV1ApisOperation(path, method)
	case strings.HasPrefix(path, "/v1/domainnames"):
		return extractAppSyncDomainNamesOperation(path, method)
	case strings.HasPrefix(path, "/v1/tags"):
		return extractAppSyncTagsOperation(path, method)
	case strings.HasPrefix(path, "/v1/mergedApis"):
		return extractAppSyncMergedApisOperation(path, method)
	case strings.HasPrefix(path, "/v1/sourceApis"):
		return extractAppSyncSourceApisOperation(path, method)
	case strings.HasPrefix(path, "/v1/datasources/introspections"):
		return extractAppSyncDataSourceIntrospectionOperation(path, method)
	case strings.HasPrefix(path, "/v1/dataplane-evaluatecode"):
		if method == http.MethodPost {
			return "EvaluateCode"
		}
	case strings.HasPrefix(path, "/v1/dataplane-evaluatetemplate"):
		if method == http.MethodPost {
			return "EvaluateMappingTemplate"
		}
	}

	return ""
}

// extractAppSyncPathParams extracts URI-bound parameters from an AppSync request path
// into the params map. Parameter names match the AWS SDK wire format (camelCase).
// Delegates to versioned/resource-specific sub-parsers for each path prefix group.
func extractAppSyncPathParams(path string, params map[string]interface{}) {
	switch {
	case strings.HasPrefix(path, "/v2/apis"):
		extractAppSyncV2PathParams(path, params)
	case strings.HasPrefix(path, "/v1/apis"):
		extractAppSyncV1ApisPathParams(path, params)
	case strings.HasPrefix(path, "/v1/domainnames"):
		extractAppSyncDomainNamesPathParams(path, params)
	case strings.HasPrefix(path, "/v1/tags"):
		extractAppSyncTagsPathParams(path, params)
	case strings.HasPrefix(path, "/v1/mergedApis"):
		extractAppSyncMergedApisPathParams(path, params)
	case strings.HasPrefix(path, "/v1/sourceApis"):
		extractAppSyncSourceApisPathParams(path, params)
	case strings.HasPrefix(path, "/v1/datasources/introspections"):
		extractAppSyncIntrospectionPathParams(path, params)
	}
}

// appSyncRESTParser implements RESTServiceParser for AWS AppSync.
type appSyncRESTParser struct{}

// MatchPath returns true if the path belongs to AppSync.
func (p *appSyncRESTParser) MatchPath(path string) bool {
	return isAppSyncPath(path)
}

// ExtractOperation returns the AppSync operation name, or empty if the path does not match.
func (p *appSyncRESTParser) ExtractOperation(r *http.Request) string {
	return extractAppSyncOperation(r)
}

// ExtractPathParams extracts URI-bound parameters from the AppSync request path.
func (p *appSyncRESTParser) ExtractPathParams(r *http.Request, params map[string]interface{}) {
	extractAppSyncPathParams(r.URL.Path, params)
}
