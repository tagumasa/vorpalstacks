package appsync

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ============================================================================
// Enum validation maps (sourced from Smithy model 2017-07-25)
// ============================================================================

// AuthenticationType enum: smithy.api#AuthenticationType
var validAuthenticationTypes = map[string]bool{
	"API_KEY":                   true,
	"AWS_IAM":                   true,
	"AMAZON_COGNITO_USER_POOLS": true,
	"OPENID_CONNECT":            true,
	"AWS_LAMBDA":                true,
}

// DataSourceType enum: smithy.api#DataSourceType
var validDataSourceTypes = map[string]bool{
	"AWS_LAMBDA":                true,
	"AMAZON_DYNAMODB":           true,
	"AMAZON_ELASTICSEARCH":      true,
	"NONE":                      true,
	"HTTP":                      true,
	"RELATIONAL_DATABASE":       true,
	"AMAZON_OPENSEARCH_SERVICE": true,
	"AMAZON_EVENTBRIDGE":        true,
	"AMAZON_BEDROCK_RUNTIME":    true,
}

// ResolverKind enum: smithy.api#ResolverKind
var validResolverKinds = map[string]bool{
	"UNIT":     true,
	"PIPELINE": true,
}

// ApiCachingBehavior enum: smithy.api#ApiCachingBehavior
var validApiCachingBehaviors = map[string]bool{
	"FULL_REQUEST_CACHING":    true,
	"PER_RESOLVER_CACHING":    true,
	"OPERATION_LEVEL_CACHING": true,
}

// ApiCacheType enum: smithy.api#ApiCacheType
var validApiCacheTypes = map[string]bool{
	"T2_SMALL":   true,
	"T2_MEDIUM":  true,
	"R4_LARGE":   true,
	"R4_XLARGE":  true,
	"R4_2XLARGE": true,
	"R4_4XLARGE": true,
	"R4_8XLARGE": true,
	"SMALL":      true,
	"MEDIUM":     true,
	"LARGE":      true,
	"XLARGE":     true,
	"LARGE_2X":   true,
	"LARGE_4X":   true,
	"LARGE_8X":   true,
	"LARGE_12X":  true,
}

// ConflictDetectionType enum: smithy.api#ConflictDetectionType
var validConflictDetectionTypes = map[string]bool{
	"VERSION": true,
	"NONE":    true,
}

// ConflictHandlerType enum: smithy.api#ConflictHandlerType
var validConflictHandlerTypes = map[string]bool{
	"OPTIMISTIC_CONCURRENCY": true,
	"LAMBDA":                 true,
	"AUTOMERGE":              true,
	"NONE":                   true,
}

// TypeDefinitionFormat enum: smithy.api#TypeDefinitionFormat
var validTypeFormats = map[string]bool{
	"SDL":  true,
	"JSON": true,
}

// GraphQLApiType enum: smithy.api#GraphQLApiType
var validApiTypes = map[string]bool{
	"GRAPHQL": true,
	"MERGED":  true,
}

// GraphQLApiVisibility enum: smithy.api#GraphQLApiVisibility
var validVisibilities = map[string]bool{
	"GLOBAL":  true,
	"PRIVATE": true,
}

// GraphQLApiIntrospectionConfig enum: smithy.api#GraphQLApiIntrospectionConfig
var validIntrospectionConfigs = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
}

// FieldLogLevel enum: smithy.api#FieldLogLevel
var validFieldLogLevels = map[string]bool{
	"NONE":  true,
	"ERROR": true,
	"ALL":   true,
	"INFO":  true,
	"DEBUG": true,
}

// RelationalDatabaseSourceType enum: smithy.api#RelationalDatabaseSourceType
var validRelationalDatabaseSourceTypes = map[string]bool{
	"RDS_HTTP_ENDPOINT": true,
}

// EnabledDisabled enum covers all *MetricsConfig / *HealthMetricsConfig shapes.
var validEnabledDisabled = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
}

// ============================================================================
// Runtime validation
// ============================================================================

// validAppSyncRuntimeVersions maps runtime name to valid versions.
// Smithy RuntimeName enum: APPSYNC_JS only.
var validAppSyncRuntimeVersions = map[string]map[string]bool{
	"APPSYNC_JS": {"1.0.0": true},
}

// validateAppSyncRuntime validates the runtime name and version against the
// Smithy RuntimeName enum and known version constraints.
func validateAppSyncRuntime(rt *appsyncstore.AppSyncRuntime) error {
	if rt == nil {
		return nil
	}
	versions, ok := validAppSyncRuntimeVersions[rt.Name]
	if !ok {
		return NewBadRequestException(fmt.Sprintf("Unsupported runtime name: %s", rt.Name))
	}
	if rt.RuntimeVersion == "" {
		return NewBadRequestException("runtimeVersion is required when runtime is specified")
	}
	if !versions[rt.RuntimeVersion] {
		return NewBadRequestException(fmt.Sprintf("Unsupported runtimeVersion %s for runtime %s", rt.RuntimeVersion, rt.Name))
	}
	return nil
}

// ============================================================================
// Caching config validation
// ============================================================================

// validateCachingConfig validates the TTL range per Smithy TTL shape
// (smithy.api#range: min 0, max 3600).
func validateCachingConfig(cc *appsyncstore.CachingConfig) error {
	if cc == nil {
		return nil
	}
	// CachingConfig.ttl targets the Long shape which carries no @range trait.
	// AWS documentation states "Valid values are 1–3,600 seconds".
	// The @default:0 on Long only denotes the server-side default when the
	// field is omitted; it does not make 0 a valid caller-supplied value.
	if cc.Ttl < 1 || cc.Ttl > 3600 {
		return NewBadRequestException("cachingConfig.ttl must be between 1 and 3600 seconds")
	}
	return nil
}

// ============================================================================
// API key expiry validation
// ============================================================================

// MinApiKeyValidityDays and MaxApiKeyValidityDays bound the API key expiry
// window: the expiration must be set to a value between 1 and 365 days from
// creation (CreateApiKey) or from update (UpdateApiKey).
const (
	MinApiKeyValidityDays = 1
	MaxApiKeyValidityDays = 365
)

// validateApiKeyExpiry validates that the expiry timestamp falls within the
// AWS-mandated range of MinApiKeyValidityDays to MaxApiKeyValidityDays from
// now.
func validateApiKeyExpiry(expires int64) error {
	now := time.Now().Unix()
	minExpiry := now + int64(MinApiKeyValidityDays)*24*3600
	maxExpiry := now + int64(MaxApiKeyValidityDays)*24*3600
	if expires < minExpiry || expires > maxExpiry {
		return ErrApiKeyValidityOutOfBoundsException
	}
	return nil
}

// ============================================================================
// Sync config validation.
// ============================================================================

// validateSyncConfig validates ConflictDetection and ConflictHandler enum values
// per Smithy ConflictDetectionType and ConflictHandlerType enums.
func validateSyncConfig(sc *appsyncstore.SyncConfig) error {
	if sc == nil {
		return nil
	}
	if sc.ConflictDetection != "" && !validConflictDetectionTypes[sc.ConflictDetection] {
		return NewBadRequestException(fmt.Sprintf("Invalid conflictDetection: %s. Valid values: VERSION, NONE", sc.ConflictDetection))
	}
	if sc.ConflictHandler != "" && !validConflictHandlerTypes[sc.ConflictHandler] {
		return NewBadRequestException(fmt.Sprintf("Invalid conflictHandler: %s. Valid values: OPTIMISTIC_CONCURRENCY, LAMBDA, AUTOMERGE, NONE", sc.ConflictHandler))
	}
	return nil
}

// ============================================================================
// ApiCache TTL validation
// ============================================================================

// validateApiCacheTtl validates the ApiCache TTL range.
// Smithy CreateApiCacheRequest.ttl documentation: "Valid values are 1–3,600 seconds."
func validateApiCacheTtl(ttl int64) error {
	if ttl < 1 || ttl > 3600 {
		return NewBadRequestException("ttl must be between 1 and 3600 seconds")
	}
	return nil
}

// ============================================================================
// Generic enum validators
// ============================================================================

// validateAuthenticationType checks the AuthenticationType enum value.
func validateAuthenticationType(authType string) bool {
	return validAuthenticationTypes[authType]
}

// validateDataSourceType checks the DataSourceType enum value.
func validateDataSourceType(dsType string) bool {
	return validDataSourceTypes[dsType]
}

// validateResolverKind checks the ResolverKind enum value.
func validateResolverKind(kind string) bool {
	return validResolverKinds[kind]
}

// validateApiCachingBehavior checks the ApiCachingBehavior enum value.
func validateApiCachingBehavior(behavior string) bool {
	return validApiCachingBehaviors[behavior]
}

// validateApiCacheType checks the ApiCacheType enum value.
func validateApiCacheType(t string) bool {
	return validApiCacheTypes[t]
}

// validateTypeFormat checks the TypeDefinitionFormat enum value.
func validateTypeFormat(format string) bool {
	return validTypeFormats[format]
}

// validateApiType checks the GraphQLApiType enum value.
func validateApiType(t string) bool {
	return validApiTypes[t]
}

// validateVisibility checks the GraphQLApiVisibility enum value.
func validateVisibility(v string) bool {
	return validVisibilities[v]
}

// validateIntrospectionConfig checks the GraphQLApiIntrospectionConfig enum value.
func validateIntrospectionConfig(c string) bool {
	return validIntrospectionConfigs[c]
}

// validateFieldLogLevel checks the FieldLogLevel enum value.
func validateFieldLogLevel(level string) bool {
	return validFieldLogLevels[level]
}

// validateRelationalDatabaseSourceType checks the RelationalDatabaseSourceType enum value.
func validateRelationalDatabaseSourceType(t string) bool {
	return validRelationalDatabaseSourceTypes[t]
}

// validateEnabledDisabled checks the ENABLED/DISABLED enum shared by all
// *MetricsConfig and *HealthMetricsConfig shapes.
func validateEnabledDisabled(v string) bool {
	return validEnabledDisabled[v]
}

// ============================================================================
// Pattern validators (Smithy [pattern] + [length] traits)
// Source: third_party/api-models-aws/models/appsync/service/2017-07-25/appsync-2017-07-25.json
// ============================================================================

// Compiled patterns from Smithy string shapes.
var (
	apiNamePattern        = regexp.MustCompile(`^[A-Za-z0-9_\- ]+$`)
	resourceNamePattern   = regexp.MustCompile(`^[_A-Za-z][_0-9A-Za-z]*$`)
	namespacePattern      = regexp.MustCompile(`^([A-Za-z0-9](?:[A-Za-z0-9\-]{0,48}[A-Za-z0-9])?)$`)
	domainNamePattern     = regexp.MustCompile(`^(\*[\w\d-]*\.)?([\w\d-]+\.)+[\w\d-]+$`)
	certificateArnPattern = regexp.MustCompile(`^arn:[a-z-]*:(acm|iam):[a-z0-9-]*:\d{12}:(certificate|server-certificate)/[0-9A-Za-z_/-]*$`)
	envVarKeyPattern      = regexp.MustCompile(`^[A-Za-z]+\w*$`)
)

// validateApiName validates the ApiName shape: ^[A-Za-z0-9_\-\ ]+$, length 1-50.
func validateApiName(name string) error {
	if len(name) < 1 || len(name) > 50 {
		return NewBadRequestException("name must be between 1 and 50 characters")
	}
	if !apiNamePattern.MatchString(name) {
		return NewBadRequestException("name contains invalid characters; must match ^[A-Za-z0-9_\\- ]+$")
	}
	return nil
}

// validateResourceName validates the ResourceName shape:
// ^[_A-Za-z][_0-9A-Za-z]*$, length 1-65536.
func validateResourceName(name string) error {
	if len(name) < 1 || len(name) > 65536 {
		return NewBadRequestException("resource name must be between 1 and 65536 characters")
	}
	if !resourceNamePattern.MatchString(name) {
		return NewBadRequestException("resource name contains invalid characters; must match ^[_A-Za-z][_0-9A-Za-z]*$")
	}
	return nil
}

// validateNamespace validates the Namespace shape used by ChannelNamespace:
// ^([A-Za-z0-9](?:[A-Za-z0-9\-]{0,48}[A-Za-z0-9])?)$, length 1-50.
func validateNamespace(name string) error {
	if len(name) < 1 || len(name) > 50 {
		return NewBadRequestException("namespace name must be between 1 and 50 characters")
	}
	if !namespacePattern.MatchString(name) {
		return NewBadRequestException("namespace name contains invalid characters; must match ^([A-Za-z0-9](?:[A-Za-z0-9\\-]{0,48}[A-Za-z0-9])?)$")
	}
	return nil
}

// validateDomainName validates the DomainName shape:
// ^(\*[\w\d-]*\.)?([\w\d-]+\.)+[\w\d-]+$, length 1-253.
func validateDomainName(name string) error {
	if len(name) < 1 || len(name) > 253 {
		return NewBadRequestException("domainName must be between 1 and 253 characters")
	}
	if !domainNamePattern.MatchString(name) {
		return NewBadRequestException("domainName has an invalid format")
	}
	return nil
}

// validateCertificateArn validates the CertificateArn shape:
// ^arn:[a-z-]*:(acm|iam):..., length 20-2048.
func validateCertificateArn(arn string) error {
	if len(arn) < 20 || len(arn) > 2048 {
		return NewBadRequestException("certificateArn must be between 20 and 2048 characters")
	}
	if !certificateArnPattern.MatchString(arn) {
		return NewBadRequestException("certificateArn has an invalid format")
	}
	return nil
}

// validateEnvVarKey validates the EnvironmentVariableKey shape:
// ^[A-Za-z]+\w*, length 2-64.
func validateEnvVarKey(key string) error {
	if len(key) < 2 || len(key) > 64 {
		return NewBadRequestException("environment variable key must be between 2 and 64 characters")
	}
	if !envVarKeyPattern.MatchString(key) {
		return NewBadRequestException("environment variable key contains invalid characters; must match ^[A-Za-z]+\\w*$")
	}
	return nil
}

// validateEnvVarValue validates the EnvironmentVariableValue shape: length
// 0-512 counted in Unicode characters (no pattern; multibyte environment
// variable values are valid input).
func validateEnvVarValue(val string) error {
	if utf8.RuneCountInString(val) > 512 {
		return NewBadRequestException("environment variable value must not exceed 512 characters")
	}
	return nil
}

// ============================================================================
// Length validators (Smithy [length] trait only)
// ============================================================================

// maxDescriptionLength is the Description shape @length maximum, counted
// in Unicode characters like every @length trait; the shape's "^.*$"
// pattern admits multibyte descriptions.
const maxDescriptionLength = 255

// validateDescription validates the Description shape: length 0-255
// counted in Unicode characters.
func validateDescription(desc string) error {
	if utf8.RuneCountInString(desc) > maxDescriptionLength {
		return NewBadRequestException("description must not exceed 255 characters")
	}
	return nil
}

// validateCode validates the Code shape: length 1-32768 counted in Unicode
// characters (no pattern; JavaScript code may contain multibyte string
// literals).
func validateCode(code string) error {
	if n := utf8.RuneCountInString(code); n < 1 || n > 32768 {
		return NewBadRequestException("code must be between 1 and 32768 characters")
	}
	return nil
}

// validateMappingTemplate validates the MappingTemplate shape: length
// 1-65536 counted in Unicode characters (the "^.*$" pattern admits
// multibyte Velocity templates).
func validateMappingTemplate(tmpl string) error {
	if n := utf8.RuneCountInString(tmpl); n < 1 || n > 65536 {
		return NewBadRequestException("mapping template must be between 1 and 65536 characters")
	}
	return nil
}

// validateContext validates the Context shape: length 2-28000 counted in
// Unicode characters (the "^[\s\S]*$" pattern admits multibyte).
func validateContext(ctx string) error {
	if n := utf8.RuneCountInString(ctx); n < 2 || n > 28000 {
		return NewBadRequestException("context must be between 2 and 28000 characters")
	}
	return nil
}

// validateTemplate validates the Template shape: length 2-65536 counted in
// Unicode characters (the "^[\s\S]*$" pattern admits multibyte).
func validateTemplate(tmpl string) error {
	if n := utf8.RuneCountInString(tmpl); n < 2 || n > 65536 {
		return NewBadRequestException("template must be between 2 and 65536 characters")
	}
	return nil
}

// ============================================================================
// Range validators (Smithy [range] trait)
// ============================================================================

// validateQueryDepthLimit validates the QueryDepthLimit shape: range 0-75.
func validateQueryDepthLimit(v int32) error {
	if v < 0 || v > 75 {
		return NewBadRequestException("queryDepthLimit must be between 0 and 75")
	}
	return nil
}

// validateResolverCountLimit validates the ResolverCountLimit shape: range 0-10000.
func validateResolverCountLimit(v int32) error {
	if v < 0 || v > 10000 {
		return NewBadRequestException("resolverCountLimit must be between 0 and 10000")
	}
	return nil
}

// validateMaxBatchSize validates the MaxBatchSize shape: range 0-2000.
func validateMaxBatchSize(v int32) error {
	if v < 0 || v > 2000 {
		return NewBadRequestException("maxBatchSize must be between 0 and 2000")
	}
	return nil
}

// validateLambdaAuthorizerTtl validates the TTL shape (shared by
// LambdaAuthorizerConfig.authorizerResultTtlInSeconds and ApiCache.ttl):
// range 0-3600.
func validateLambdaAuthorizerTtl(v int32) error {
	if v < 0 || v > 3600 {
		return NewBadRequestException("authorizerResultTtlInSeconds must be between 0 and 3600")
	}
	return nil
}

// validateEnhancedMetricsConfig validates the EnhancedMetricsConfig enum fields.
// All three fields use the ENABLED/DISABLED enum per Smithy model.
func validateEnhancedMetricsConfig(ec *appsyncstore.EnhancedMetricsConfig) error {
	if ec == nil {
		return nil
	}
	if ec.DataSourceLevelMetricsBehavior != "" && !validateEnabledDisabled(ec.DataSourceLevelMetricsBehavior) {
		return NewBadRequestException(fmt.Sprintf("Invalid dataSourceLevelMetricsBehavior: %s. Valid values: ENABLED, DISABLED", ec.DataSourceLevelMetricsBehavior))
	}
	if ec.OperationLevelMetricsConfig != "" && !validateEnabledDisabled(ec.OperationLevelMetricsConfig) {
		return NewBadRequestException(fmt.Sprintf("Invalid operationLevelMetricsConfig: %s. Valid values: ENABLED, DISABLED", ec.OperationLevelMetricsConfig))
	}
	if ec.ResolverLevelMetricsBehavior != "" && !validateEnabledDisabled(ec.ResolverLevelMetricsBehavior) {
		return NewBadRequestException(fmt.Sprintf("Invalid resolverLevelMetricsBehavior: %s. Valid values: ENABLED, DISABLED", ec.ResolverLevelMetricsBehavior))
	}
	return nil
}

// validateAuthorizationType validates the HTTP data source AuthorizationType
// enum. Per Smithy, the only valid value is AWS_IAM.
func validateAuthorizationType(t string) bool {
	return t == "AWS_IAM"
}

// validateHandlerBehavior validates the Event API handler Behavior enum.
// Per Smithy, valid values are CODE and DIRECT.
var validHandlerBehaviors = map[string]bool{
	"CODE":   true,
	"DIRECT": true,
}

func validateHandlerBehavior(b string) bool {
	return validHandlerBehaviors[b]
}

// validateOpenIDConnectTTL validates the AuthTTL and IatTTL fields of the
// OpenIDConnectConfig shape. Smithy range: min 0, max 3600.
func validateOpenIDConnectTTL(ttl int64) error {
	if ttl < 0 || ttl > 3600 {
		return NewBadRequestException("OpenIDConnect TTL values must be between 0 and 3600")
	}
	return nil
}

// validateDeltaSyncTtl validates the DeltaSyncConfig TTL fields.
// Smithy range: BaseTableTTL 1-43200, DeltaSyncTableTTL 1-43200.
func validateDeltaSyncTtl(baseTTL, deltaTTL int64) error {
	if baseTTL < 1 || baseTTL > 43200 {
		return NewBadRequestException("deltaSyncConfig.baseTableTTL must be between 1 and 43200")
	}
	if deltaTTL < 1 || deltaTTL > 43200 {
		return NewBadRequestException("deltaSyncConfig.deltaSyncTableTTL must be between 1 and 43200")
	}
	return nil
}

// validateLambdaArn validates that the string is a well-formed Lambda ARN.
func validateLambdaArn(arn string) error {
	if arn == "" {
		return nil
	}
	partition, service, _, _, resource := svcarn.SplitARN(arn)
	if (partition != "aws" && partition != "aws-cn" && partition != "aws-us-gov") || service != "lambda" {
		return NewBadRequestException(fmt.Sprintf("Invalid Lambda ARN: %s", arn))
	}
	if !strings.Contains(resource, ":") {
		return NewBadRequestException(fmt.Sprintf("Invalid Lambda ARN format: %s", arn))
	}
	return nil
}

// validateRdsIdentifier validates the RDS HTTP endpoint config identifiers.
func validateRdsIdentifier(cfg *appsyncstore.RdsHttpEndpointConfig) error {
	if cfg == nil {
		return nil
	}
	if cfg.DatabaseName != "" && len(cfg.DatabaseName) > 100 {
		return NewBadRequestException("databaseName must not exceed 100 characters")
	}
	if cfg.DbClusterIdentifier != "" && len(cfg.DbClusterIdentifier) > 255 {
		return NewBadRequestException("dbClusterIdentifier must not exceed 255 characters")
	}
	if cfg.Schema != "" && len(cfg.Schema) > 100 {
		return NewBadRequestException("schema must not exceed 100 characters")
	}
	return nil
}

// validateEnvironmentVariableMapSize validates the EnvironmentVariableMap
// shape length constraint: min 0, max 50.
func validateEnvironmentVariableMapSize(m map[string]string) error {
	if len(m) > 50 {
		return NewBadRequestException("environmentVariables must not contain more than 50 entries")
	}
	return nil
}
