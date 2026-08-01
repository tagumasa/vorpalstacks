package appsync

import (
	"fmt"
	"time"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
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
// Runtime validation (moved from helpers.go)
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
// Caching config validation (moved from helpers.go)
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
// API key expiry validation (moved from api_key_operations.go)
// ============================================================================

// validateApiKeyExpiry validates that the expiry timestamp falls within the
// AWS-mandated range of 1 day to 365 days from now.
func validateApiKeyExpiry(expires int64) error {
	now := time.Now().Unix()
	minExpiry := now + 86400    // 1 day
	maxExpiry := now + 31536000 // 365 days
	if expires < minExpiry || expires > maxExpiry {
		return ErrApiKeyValidityOutOfBoundsException
	}
	return nil
}

// ============================================================================
// Sync config validation (new — addresses audit M5)
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
