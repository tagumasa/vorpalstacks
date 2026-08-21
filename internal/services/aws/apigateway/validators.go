package apigateway

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// stageNamePattern matches AWS API Gateway stage name rules: alphanumeric,
// underscore, and hyphen only, max 128 characters.
var stageNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,128}$`)

// fqdnPattern matches fully-qualified domain names (simplified RFC 1123).
var fqdnPattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

// validHTTPMethods is the set of accepted HTTP methods for API Gateway methods.
var validHTTPMethods = map[string]bool{
	"GET":     true,
	"PUT":     true,
	"POST":    true,
	"DELETE":  true,
	"PATCH":   true,
	"HEAD":    true,
	"OPTIONS": true,
	"ANY":     true,
}

// validAuthorizationTypes is the set of accepted authorization types.
var validAuthorizationTypes = map[string]bool{
	"NONE":               true,
	"AWS_IAM":            true,
	"CUSTOM":             true,
	"COGNITO_USER_POOLS": true,
}

// validCacheClusterSizes is the set of accepted cache cluster size values.
var validCacheClusterSizes = map[string]bool{
	"0.5":  true,
	"1.6":  true,
	"6.1":  true,
	"13.5": true,
	"28.4": true,
	"58.2": true,
	"118":  true,
	"237":  true,
}

// validLoggingLevels is the set of accepted logging level values.
var validLoggingLevels = map[string]bool{
	"OFF":   true,
	"ERROR": true,
	"INFO":  true,
}

// validIntegrationTypes is the set of accepted integration types.
var validIntegrationTypes = map[string]bool{
	"HTTP":       true,
	"HTTP_PROXY": true,
	"AWS":        true,
	"AWS_PROXY":  true,
	"MOCK":       true,
}

// validateStageName returns true if the stage name conforms to AWS API
// Gateway naming rules.
func validateStageName(name string) bool {
	return stageNamePattern.MatchString(name)
}

// validateFQDN returns true if the string is a syntactically valid
// fully-qualified domain name.
func validateFQDN(domain string) bool {
	if len(domain) == 0 || len(domain) > 255 {
		return false
	}
	return fqdnPattern.MatchString(domain)
}

// validateHTTPMethod returns true if the method is a recognised HTTP verb
// or the ANY wildcard.
func validateHTTPMethod(method string) bool {
	return validHTTPMethods[strings.ToUpper(method)]
}

// validateAuthorizationType returns true if the auth type is recognised.
func validateAuthorizationType(authType string) bool {
	return validAuthorizationTypes[authType]
}

// validateCacheClusterSize returns true if the size is a recognised value.
func validateCacheClusterSize(size string) bool {
	return size == "" || validCacheClusterSizes[size]
}

// validateLoggingLevel returns true if the level is recognised.
func validateLoggingLevel(level string) bool {
	return validLoggingLevels[level]
}

// validEndpointAccessModes is the set of accepted endpoint access mode values.
var validEndpointAccessModes = map[string]bool{
	"BASIC":  true,
	"STRICT": true,
}

// validateEndpointAccessMode returns true if the value is empty (optional)
// or a recognised endpoint access mode.
func validateEndpointAccessMode(mode string) bool {
	return mode == "" || validEndpointAccessModes[mode]
}

// validApiKeySources is the set of accepted API key source values.
var validApiKeySources = map[string]bool{
	"HEADER":     true,
	"AUTHORIZER": true,
}

// validateApiKeySource returns true if the value is empty (optional) or
// a recognised API key source.
func validateApiKeySource(source string) bool {
	return source == "" || validApiKeySources[source]
}

// validateSecurityPolicy returns true if the value is empty (optional) or
// starts with a recognised security policy prefix.
func validateSecurityPolicy(policy string) bool {
	if policy == "" {
		return true
	}
	return policy == "TLS_1_0" || policy == "TLS_1_2" ||
		strings.HasPrefix(policy, "SecurityPolicy_")
}

// validateIntegrationType returns true if the integration type is recognised.
func validateIntegrationType(t string) bool {
	return validIntegrationTypes[t]
}

// modelNamePattern matches AWS API Gateway model name rules. The API
// reference states "Must be alphanumeric", but the actual service also
// accepts hyphens and underscores — these are universally used in
// practice and accepted by real AWS API Gateway.
var modelNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// validateModelName returns true if the model name conforms to AWS API
// Gateway naming rules (alphanumeric, hyphen, underscore).
func validateModelName(name string) bool {
	return modelNamePattern.MatchString(name)
}

// maxModelSchemaSize is the maximum allowed size for a model schema in bytes.
// AWS API Gateway documentation states: "The maximum size of the model is 400 KB."
const maxModelSchemaSize = 400 * 1024

// validateModelSchemaSize returns true if the schema does not exceed the
// 400 KB limit imposed by AWS API Gateway.
func validateModelSchemaSize(schema string) bool {
	return len(schema) <= maxModelSchemaSize
}

// basePathPattern matches valid API Gateway basePath values: URL-safe
// path characters (alphanumeric, hyphen, underscore, period, forward
// slash).
var basePathPattern = regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)

// validateBasePath returns true if the basePath is empty, "(none)", or
// contains only URL-safe path characters accepted by AWS API Gateway.
func validateBasePath(basePath string) bool {
	if basePath == "" || basePath == "(none)" {
		return true
	}
	return basePathPattern.MatchString(basePath)
}

// stageKeyPattern matches the API Gateway stage key format: a non-empty
// restApiId, a forward slash, and a non-empty stageName. Both segments
// allow alphanumeric, hyphen, and underscore.
var stageKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+$`)

// validateStageKey returns true if the stage key conforms to the
// "restApiId/stageName" format expected by AWS API Gateway.
func validateStageKey(sk string) bool {
	return stageKeyPattern.MatchString(sk)
}

// pathParamPattern matches valid path parameter syntax: {name} or {name+}
// (greedy). A pathPart that does not use braces is always valid.
var pathParamPattern = regexp.MustCompile(`^\{[a-zA-Z0-9_-]+\+\}$|^\{[a-zA-Z0-9_-]+\}$`)

// validatePathPart returns true if the path segment is syntactically valid.
// Segments with braces must form exactly one well-formed path parameter.
func validatePathPart(part string) bool {
	hasOpen := strings.Contains(part, "{")
	hasClose := strings.Contains(part, "}")
	if !hasOpen && !hasClose {
		return true
	}
	if hasOpen != hasClose {
		return false
	}
	return pathParamPattern.MatchString(part)
}

// validAuthorizerTypes is the set of accepted authorizer types.
var validAuthorizerTypes = map[string]bool{
	"TOKEN":              true,
	"REQUEST":            true,
	"COGNITO_USER_POOLS": true,
}

// validConnectionTypes is the set of accepted connection types.
var validConnectionTypes = map[string]bool{
	"INTERNET": true,
	"VPC_LINK": true,
}

// validContentHandlingStrategies is the set of accepted content handling values.
var validContentHandlingStrategies = map[string]bool{
	"CONVERT_TO_BINARY": true,
	"CONVERT_TO_TEXT":   true,
}

// validResponseTransferModes is the set of accepted response transfer mode
// values per the AWS PutIntegration API specification.
var validResponseTransferModes = map[string]bool{
	"BUFFERED": true,
	"STREAM":   true,
}

// validateResponseTransferMode returns true if the response transfer mode is
// recognised. Empty string is accepted (field not provided).
func validateResponseTransferMode(mode string) bool {
	if mode == "" {
		return true
	}
	return validResponseTransferModes[mode]
}

// validPassthroughBehaviors is the set of accepted passthrough behaviour values.
var validPassthroughBehaviors = map[string]bool{
	"WHEN_NO_MATCH":     true,
	"WHEN_NO_TEMPLATES": true,
	"NEVER":             true,
}

// validQuotaPeriods is the set of accepted quota period values.
var validQuotaPeriods = map[string]bool{
	"DAY":   true,
	"WEEK":  true,
	"MONTH": true,
}

// validRoutingModes is the set of accepted routing mode values.
var validRoutingModes = map[string]bool{
	"ROUTING_RULE_THEN_BASE_PATH_MAPPING": true,
	"BASE_PATH_MAPPING_ONLY":              true,
	"ROUTING_RULE_ONLY":                   true,
}

// usageDatePattern matches the AWS GetUsage date format: YYYY-MM-DD.
var usageDatePattern = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`)

// validateUsageDateFormat returns true if the date string conforms to the
// YYYY-MM-DD format and represents a valid calendar date.
func validateUsageDateFormat(date string) bool {
	if !usageDatePattern.MatchString(date) {
		return false
	}
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// validateAuthorizerType returns true if the authorizer type is recognised.
// Empty string is accepted (no change to existing value in patch context).
func validateAuthorizerType(t string) bool {
	return t == "" || validAuthorizerTypes[t]
}

// validateConnectionType returns true if the connection type is recognised.
func validateConnectionType(ct string) bool {
	return ct == "" || validConnectionTypes[ct]
}

// validateContentHandling returns true if the content handling strategy is
// recognised.
func validateContentHandling(ch string) bool {
	return ch == "" || validContentHandlingStrategies[ch]
}

// validatePassthroughBehavior returns true if the passthrough behaviour is
// recognised.
func validatePassthroughBehavior(pb string) bool {
	return pb == "" || validPassthroughBehaviors[pb]
}

// validateQuotaPeriod returns true if the quota period is recognised.
func validateQuotaPeriod(period string) bool {
	return validQuotaPeriods[period]
}

// validateRoutingMode returns true if the routing mode is recognised.
func validateRoutingMode(rm string) bool {
	return rm == "" || validRoutingModes[rm]
}

// validatePolicyJSON returns true if the policy string is empty or contains
// a valid JSON object. AWS resource policies must be JSON objects, not
// scalars, arrays, or null.
func validatePolicyJSON(policy string) bool {
	if policy == "" {
		return true
	}
	var obj map[string]json.RawMessage
	return json.Unmarshal([]byte(policy), &obj) == nil
}

// maxThrottleBurstLimit is the upper bound for throttle burst limit values.
const maxThrottleBurstLimit = 10000

// maxThrottleRateLimit is the upper bound for throttle rate limit values.
const maxThrottleRateLimit = 10000.0

// validateThrottleBurstLimit returns true if the burst limit is within the
// accepted range [0, 10000].
func validateThrottleBurstLimit(v int64) bool {
	return v >= 0 && v <= maxThrottleBurstLimit
}

// validateThrottleRateLimit returns true if the rate limit is within the
// accepted range [0, 10000].
func validateThrottleRateLimit(v float64) bool {
	return v >= 0 && v <= maxThrottleRateLimit
}

// maxMinimumCompressionSize is the upper bound for minimumCompressionSize.
const maxMinimumCompressionSize = 10485760

// validateMinimumCompressionSize returns true if the value is within the
// accepted range [0, 10485760].
func validateMinimumCompressionSize(v int32) bool {
	return v >= 0 && v <= maxMinimumCompressionSize
}

// validateTimeoutInMillis returns true if the value is within the accepted
// range [50, 30000] for integration timeout.
func validateTimeoutInMillis(v int32) bool {
	return v >= 50 && v <= 30000
}

// validateCacheTtlInSeconds returns true if the value is within the
// accepted range [0, 86400] for method-setting cache TTL.
func validateCacheTtlInSeconds(v int32) bool {
	return v >= 0 && v <= 86400
}

// validatePercentTraffic returns true if the value is within the accepted
// range [0, 100] for canary percent traffic.
func validatePercentTraffic(v float64) bool {
	return v >= 0 && v <= 100
}

// validateAuthorizerTtl returns true if the value is within the accepted
// range [0, 3600] for authorizer result TTL.
func validateAuthorizerTtl(v int32) bool {
	return v >= 0 && v <= 3600
}

// validateUsagePlanNameLen returns true if the name length is within the
// accepted range [1, 255].
func validateUsagePlanNameLen(name string) bool {
	return len(name) >= 1 && len(name) <= 255
}

// maxMethodSettingThrottleBurstLimit is the upper bound for per-method
// throttle burst limit in stage method settings (distinct from the
// usage-plan-level limit of 10000).
const maxMethodSettingThrottleBurstLimit = 100000

// maxMethodSettingThrottleRateLimit is the upper bound for per-method
// throttle rate limit in stage method settings.
const maxMethodSettingThrottleRateLimit = 100000.0

// validateMethodSettingThrottleBurstLimit returns true if the per-method
// burst limit is within the accepted range [0, 100000].
func validateMethodSettingThrottleBurstLimit(v int64) bool {
	return v >= 0 && v <= maxMethodSettingThrottleBurstLimit
}

// validateMethodSettingThrottleRateLimit returns true if the per-method
// rate limit is within the accepted range [0, 100000].
func validateMethodSettingThrottleRateLimit(v float64) bool {
	return v >= 0 && v <= maxMethodSettingThrottleRateLimit
}

// validateAccessLogDestinationArn returns true if the ARN references either
// a CloudWatch Logs log group or a Kinesis Data Firehose delivery stream.
// AWS requires Firehose stream names used for access logs to begin with
// "amazon-apigateway-".
func validateAccessLogDestinationArn(arn string) bool {
	_, service, region, _, resource := svcarn.SplitARN(arn)
	if service == "" || region == "" {
		return false
	}
	switch service {
	case "logs":
		// CloudWatch Logs log group ARNs place the group name after a
		// colon: arn:<partition>:logs:<region>:<account>:log-group:<name>[:*]
		segs := strings.SplitN(resource, ":", 3)
		return len(segs) >= 2 && segs[0] == "log-group" && segs[1] != ""
	case "firehose":
		// Firehose ARNs use a slash before the stream name:
		// arn:<partition>:firehose:<region>:<account>:deliverystream/<name>
		if !strings.HasPrefix(resource, "deliverystream/") || len(resource) <= len("deliverystream/") {
			return false
		}
		streamName := strings.TrimPrefix(resource, "deliverystream/")
		return strings.HasPrefix(streamName, "amazon-apigateway-")
	default:
		return false
	}
}
