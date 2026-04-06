package classifier

import (
	"strings"

	"vorpalstacks/internal/services/aws/common/request"
)

// lambdaRestPrefixes lists known Lambda REST API date-based path prefixes.
var lambdaRestPrefixes = []string{
	"/2014-11-13/",
	"/2015-03-31/",
	"/2016-08-19/",
	"/2017-03-31/",
	"/2017-10-31/",
	"/2018-10-31/",
	"/2019-09-25/",
	"/2019-09-30/",
	"/2020-01-01/",
	"/2020-06-30/",
	"/2021-01-01/",
	"/2021-10-01/",
	"/2021-10-31/",
	"/2021-11-15/",
	"/2022-01-01/",
	"/2022-07-01/",
	"/2023-01-01/",
	"/2023-07-01/",
	"/2024-01-01/",
	"/2025-01-01/",
}

// isLambdaRestPath reports whether the path matches a known Lambda REST API prefix.
func isLambdaRestPath(path string) bool {
	for _, prefix := range lambdaRestPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// isApiGatewayRuntimePath reports whether the path is an API Gateway runtime
// invocation path (e.g. /restapis/{id}/{stage}/_user_request_/...).
func isApiGatewayRuntimePath(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 4 && parts[0] == "restapis" && parts[3] == "_user_request_" {
		return true
	}
	return false
}

// isApiGatewayPath reports whether the path belongs to the API Gateway
// control plane. Runtime invocation paths are excluded.
func isApiGatewayPath(path string) bool {
	if isApiGatewayRuntimePath(path) {
		return false
	}
	return strings.HasPrefix(path, "/restapis") ||
		strings.HasPrefix(path, "/apikeys") ||
		strings.HasPrefix(path, "/usageplans") ||
		strings.HasPrefix(path, "/domainnames") ||
		strings.HasPrefix(path, "/vpclinks") ||
		strings.HasPrefix(path, "/apis") ||
		strings.HasPrefix(path, "/authorizers") ||
		(strings.HasPrefix(path, "/tags/") && strings.Contains(path, "apigateway"))
}

// isNeptunedataPath reports whether the path matches a Neptune Data API endpoint
// (gremlin, opencypher, sparql, rdf, ml, etc.). Delegates to the canonical
// implementation in the request package to avoid duplication.
func isNeptunedataPath(path string) bool {
	return request.IsNeptunedataPath(path)
}

// isAppSyncPath reports whether the path matches an AppSync API endpoint.
func isAppSyncPath(path string) bool {
	return strings.HasPrefix(path, "/v1/apis") || strings.HasPrefix(path, "/v2/apis") ||
		strings.HasPrefix(path, "/v1/domainnames") || strings.HasPrefix(path, "/v1/tags") ||
		strings.HasPrefix(path, "/v1/mergedApis") || strings.HasPrefix(path, "/v1/sourceApis") ||
		strings.HasPrefix(path, "/v1/datasources/introspections") || strings.HasPrefix(path, "/v1/dataplane-")
}

// lookupServiceByPath maps a URL path to an internal service name based on
// known path prefixes. Returns an empty string if no service matches.
func lookupServiceByPath(path string) string {
	if isLambdaRestPath(path) {
		return "lambda"
	}
	if isApiGatewayPath(path) {
		return "apigateway"
	}
	if strings.HasPrefix(path, "/schedule-groups") ||
		strings.HasPrefix(path, "/schedules") ||
		strings.HasPrefix(path, "/tags/arn:aws:scheduler") {
		return "scheduler"
	}
	if strings.HasPrefix(path, "/v2/email/") {
		return "email"
	}
	if strings.HasPrefix(path, "/2013-04-01/") {
		return "route53"
	}
	if strings.HasPrefix(path, "/2020-05-31/") {
		return "cloudfront"
	}
	if isNeptunedataPath(path) {
		return "neptunedata"
	}
	if strings.HasPrefix(path, "/service/GraniteServiceVersion20100801/") {
		return "monitoring"
	}
	if isAppSyncPath(path) {
		return "appsync"
	}
	return ""
}
