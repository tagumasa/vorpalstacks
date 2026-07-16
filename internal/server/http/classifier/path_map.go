package classifier

import (
	"strings"

	"vorpalstacks/internal/common/request"
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
	if isRDSDataPath(path) {
		return "rdsdata"
	}
	if request.IsNeptuneGraphPath(path) {
		return "neptunegraph"
	}
	if isIoTEventsPath(path) {
		return "iotevents"
	}
	if isIoTPath(path) {
		return "iot"
	}
	return ""
}

func isRDSDataPath(path string) bool {
	switch path {
	case "/Execute", "/BatchExecute", "/ExecuteSql",
		"/BeginTransaction", "/CommitTransaction", "/RollbackTransaction":
		return true
	default:
		return false
	}
}

func rdsDataOperationFromPath(path string) string {
	switch path {
	case "/Execute":
		return "ExecuteStatement"
	case "/BatchExecute":
		return "BatchExecuteStatement"
	case "/ExecuteSql":
		return "ExecuteSql"
	case "/BeginTransaction":
		return "BeginTransaction"
	case "/CommitTransaction":
		return "CommitTransaction"
	case "/RollbackTransaction":
		return "RollbackTransaction"
	default:
		return ""
	}
}

// iotPathPrefixes lists the URL path prefixes used by the AWS IoT Core
// REST-JSON control plane API. These are derived from the Smithy model
// (third_party/api-models-aws/models/iot/service/2015-05-28/iot-2015-05-28.json)
// @http uri traits. When a path is missing from this list, the AWS SDK v2
// (restJson1 protocol, no X-Amz-Target header) cannot route the request,
// causing the dispatcher to fall back to mock generation and return '{}'.
// R18-C1: added 41 missing prefixes to restore routing for ~140 operations.
var iotPathPrefixes = []string{
	// Original prefixes (pre-R18)
	"/things", "/thing-groups", "/thing-types", "/billing-groups",
	"/api/things/shadow",
	"/certificates", "/certificate", "/keys-and-certificate",
	"/policies", "/policy-principals", "/principal-policies",
	"/principals", "/target-policies", "/attached-policies",
	"/rules", "/jobs", "/endpoint", "/role-aliases",
	"/tags", "/untag", "/authorizers", "/authorizer/",
	"/provisioning-templates", "/provisioning-template",
	"/domainConfigurations", "/domainConfiguration",
	"/indexing", "/active-violations", "/violation-events",
	"/behavior-model-training", "/security-profiles",
	"/security-profile-behaviors",
	"/security-profiles-for-target",
	"/detector-models", "/detector-model/",
	"/inputs", "/input/",
	"/messages", "/destinations", "/effective-policies",
	// R18-C1: 41 missing prefixes (extracted from Smithy @http uri traits)
	// CA Certificate lifecycle (12 ops)
	"/cacertificate", "/cacertificates",
	"/accept-certificate-transfer", "/cancel-certificate-transfer",
	"/reject-certificate-transfer", "/transfer-certificate",
	"/certificates-by-ca", "/certificates-out-going",
	"/registrationcode",
	// Audit subsystem (~26 ops)
	"/audit",
	// Detect / Mitigation (~10 ops)
	"/detect", "/mitigationactions",
	// Custom Metrics / Dimensions (~10 ops)
	"/custom-metric", "/custom-metrics",
	"/dimensions",
	// Fleet Metrics (5 ops)
	"/fleet-metric", "/fleet-metrics",
	// Job Templates (6 ops)
	"/job-templates", "/managed-job-templates",
	// Logging / Encryption / Event Config (12 ops)
	"/loggingOptions", "/encryption-configuration",
	"/event-configurations", "/v2LoggingLevel", "/v2LoggingOptions",
	// Indices / Search (5 ops)
	"/indices",
	// Streams (5 ops)
	"/streams",
	// OTA Updates (4 ops)
	"/otaUpdates",
	// Thing Registration Tasks (5 ops)
	"/thing-registration-tasks",
	// Violations (1 op + PutVerificationState)
	"/violations",
	// Dynamic Thing Groups (3 ops)
	"/dynamic-thing-groups",
	// Authorizer extensions (4 ops; /authorizer singular)
	"/authorizer",
	// Default authorizer ops (3 ops)
	"/default-authorizer",
	// Policy targets (1 op)
	"/policy-targets",
	// TestAuthorization (1 op)
	"/test-authorization",
	// TopicRuleDestination confirmation (1 op)
	"/confirmdestination",
	// Tier 1 deferred features (paths still needed for routing classification
	// even if handlers are stubs — otherwise the stubs return '{}')
	"/certificate-providers",
	"/command-executions", "/commands",
	"/package-configuration", "/packages",
}

func isIoTPath(path string) bool {
	for _, prefix := range iotPathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// ioteventsPathPrefixes lists URL path prefixes for AWS IoT Events control
// plane. Derived from the Smithy model at
// third_party/api-models-aws/models/iot-events/service/2018-07-27/iot-events-2018-07-27.json.
// R18-C1: populated (previously empty) so IoT Events operations route
// correctly via restJson1 protocol.
var ioteventsPathPrefixes = []string{
	"/alarm-models",
	"/detector-models", "/detector-model/",
	"/inputs", "/input/",
	"/input-routings",
	"/logging",
	"/analysis",
	"/tags",
}

func isIoTEventsPath(path string) bool {
	for _, prefix := range ioteventsPathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
