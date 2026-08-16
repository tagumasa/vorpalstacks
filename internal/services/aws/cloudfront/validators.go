package cloudfront

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
)

// ---------------------------------------------------------------------------
// Smithy-derived patterns
// ---------------------------------------------------------------------------

// resourceIdPattern validates CloudFront resource IDs per the Smithy
// ResourceId shape: @length(min:1, max:64).
var resourceIdPattern = regexp.MustCompile(`^.{1,64}$`)

// tagKeyPattern validates tag keys per the Smithy TagKey shape.
var tagKeyPattern = regexp.MustCompile(`^([\p{L}\p{Z}\p{N}_.:/=+\-@]{1,128})$`)

// tagValuePattern validates tag values per the Smithy TagValue shape.
var tagValuePattern = regexp.MustCompile(`^([\p{L}\p{Z}\p{N}_.:/=+\-@]{0,256})$`)

// resourceArnPattern validates CloudFront ARNs per the Smithy ResourceARN
// shape: @pattern("^arn:aws(-cn)?:cloudfront::[0-9]+:").
var resourceArnPattern = regexp.MustCompile(`^arn:aws(-cn)?:cloudfront::[0-9]+:`)

// ---------------------------------------------------------------------------
// Smithy-derived enum sets
// ---------------------------------------------------------------------------

var validPriceClasses = map[string]bool{
	"PriceClass_100": true,
	"PriceClass_200": true,
	"PriceClass_All": true,
}

var validHttpVersions = map[string]bool{
	"http1.1":   true,
	"http2":     true,
	"http2and3": true,
}

var validOriginProtocolPolicies = map[string]bool{
	"http-only":    true,
	"match-viewer": true,
	"https-only":   true,
}

var validItemSelections = map[string]bool{
	"none":      true,
	"whitelist": true,
	"all":       true,
}

// CachePolicy enums — distinct from ItemSelection per the Smithy model.
var validCachePolicyQueryStringBehaviors = map[string]bool{
	"none":      true,
	"whitelist": true,
	"allExcept": true,
	"all":       true,
}

var validCachePolicyCookieBehaviors = map[string]bool{
	"none":      true,
	"whitelist": true,
	"allExcept": true,
	"all":       true,
}

var validCachePolicyHeaderBehaviors = map[string]bool{
	"none":      true,
	"whitelist": true,
}

// OriginRequestPolicy enums — distinct from CachePolicy and ItemSelection.
var validORPCookieBehaviors = map[string]bool{
	"none":      true,
	"whitelist": true,
	"all":       true,
	"allExcept": true,
}

var validORPHeaderBehaviors = map[string]bool{
	"none":                            true,
	"whitelist":                       true,
	"allViewer":                       true,
	"allViewerAndWhitelistCloudFront": true,
	"allExcept":                       true,
}

var validORPQueryStringBehaviors = map[string]bool{
	"none":      true,
	"whitelist": true,
	"all":       true,
	"allExcept": true,
}

var validOriginAccessControlOriginTypes = map[string]bool{
	"s3":             true,
	"mediastore":     true,
	"mediapackagev2": true,
	"lambda":         true,
	"smartplayer":    true,
}

var validSigningBehaviors = map[string]bool{
	"always":      true,
	"never":       true,
	"no-override": true,
}

var validSigningProtocols = map[string]bool{
	"sigv4": true,
}

// ---------------------------------------------------------------------------
// Validators — bool-returning per S8 policy
// ---------------------------------------------------------------------------

func isValidPriceClass(v string) bool {
	return validPriceClasses[v]
}

func isValidHttpVersion(v string) bool {
	return v == "" || validHttpVersions[v]
}

func isValidOriginProtocolPolicy(v string) bool {
	return validOriginProtocolPolicies[v]
}

func isValidItemSelection(v string) bool {
	return validItemSelections[v]
}

func isValidCachePolicyQueryStringBehavior(v string) bool {
	return validCachePolicyQueryStringBehaviors[v]
}

func isValidCachePolicyCookieBehavior(v string) bool {
	return validCachePolicyCookieBehaviors[v]
}

func isValidCachePolicyHeaderBehavior(v string) bool {
	return validCachePolicyHeaderBehaviors[v]
}

func isValidORPCookieBehavior(v string) bool {
	return validORPCookieBehaviors[v]
}

func isValidORPHeaderBehavior(v string) bool {
	return validORPHeaderBehaviors[v]
}

func isValidORPQueryStringBehavior(v string) bool {
	return validORPQueryStringBehaviors[v]
}

// validateBehavior returns an InvalidArgument error when the value is
// non-empty and not in the allowed set. Empty values are accepted because
// the behaviour fields are optional at the Smithy structural level.
func validateBehavior(field, value string, valid func(string) bool) error {
	if value != "" && !valid(value) {
		return invalidArgument(fmt.Sprintf("Invalid %s: %q", field, value))
	}
	return nil
}

func isValidResourceId(id string) bool {
	return resourceIdPattern.MatchString(id)
}

func isValidTagKey(k string) bool {
	return tagKeyPattern.MatchString(k)
}

func isValidTagValue(v string) bool {
	return tagValuePattern.MatchString(v)
}

func isValidResourceArn(arn string) bool {
	return resourceArnPattern.MatchString(arn)
}

func isValidOriginAccessControlOriginType(t string) bool {
	return validOriginAccessControlOriginTypes[t]
}

func originAccessControlOriginTypeValues() string {
	keys := make([]string, 0, len(validOriginAccessControlOriginTypes))
	for k := range validOriginAccessControlOriginTypes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func isValidSigningBehavior(b string) bool {
	return validSigningBehaviors[b]
}

func signingBehaviorValues() string {
	keys := make([]string, 0, len(validSigningBehaviors))
	for k := range validSigningBehaviors {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func isValidSigningProtocol(p string) bool {
	return validSigningProtocols[p]
}

// isValidPolicyListType reports whether the value is one of the managed
// and custom filters accepted by the policy list operations.
func isValidPolicyListType(t string) bool {
	return t == "managed" || t == "custom"
}

// policyMatchesListType reports whether a policy with the given managed
// flag is included by the managed|custom list filter. An empty filter
// includes every policy.
func policyMatchesListType(isManaged bool, listType string) bool {
	if listType == "managed" {
		return isManaged
	}
	if listType == "custom" {
		return !isManaged
	}
	return true
}

func signingProtocolValues() string {
	keys := make([]string, 0, len(validSigningProtocols))
	for k := range validSigningProtocols {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

// ---------------------------------------------------------------------------
// Error helpers
// ---------------------------------------------------------------------------

func invalidArgument(msg string) error {
	return awserrors.NewAWSError("InvalidArgument", msg, 400)
}
