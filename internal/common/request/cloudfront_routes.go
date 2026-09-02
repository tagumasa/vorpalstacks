package request

import (
	"net/url"
	"strings"
)

// cloudFrontRoute is one row of the CloudFront REST routing table. It mirrors
// the smithy.api#http binding of the operation in the CloudFront API model:
// method and segments select the operation, {Name} placeholders capture URI
// labels into the parameter named Name, and payloadRoot names the XML root
// element of the response (empty for operations with no response body).
// query narrows a row to requests carrying a given query parameter; an empty
// value means the parameter only has to be present.
type cloudFrontRoute struct {
	op          string
	method      string
	segments    []string
	query       map[string]string
	payloadRoot string
}

// cloudFrontRoutes lists every CloudFront operation the service registers.
// Rows carrying query constraints must precede the unconstrained rows they
// overlap with (see the tagging and distribution collection rows). A guard
// test cross-checks this table against both the handler registry and the
// API model so a registered operation can never be left unrouted.
var cloudFrontRoutes = []cloudFrontRoute{
	// Distribution
	{op: "CreateDistributionWithTags", method: "POST", segments: []string{"distribution"}, query: map[string]string{"WithTags": ""}, payloadRoot: "Distribution"},
	{op: "CreateDistribution", method: "POST", segments: []string{"distribution"}, payloadRoot: "Distribution"},
	{op: "ListDistributions", method: "GET", segments: []string{"distribution"}, payloadRoot: "DistributionList"},
	{op: "GetDistributionConfig", method: "GET", segments: []string{"distribution", "{Id}", "config"}, payloadRoot: "DistributionConfig"},
	{op: "UpdateDistribution", method: "PUT", segments: []string{"distribution", "{Id}", "config"}, payloadRoot: "Distribution"},
	{op: "GetInvalidation", method: "GET", segments: []string{"distribution", "{Id}", "invalidation", "{invalidationId}"}, payloadRoot: "Invalidation"},
	{op: "ListInvalidations", method: "GET", segments: []string{"distribution", "{Id}", "invalidation"}, payloadRoot: "InvalidationList"},
	{op: "CreateInvalidation", method: "POST", segments: []string{"distribution", "{Id}", "invalidation"}, payloadRoot: "Invalidation"},
	{op: "AssociateDistributionWebACL", method: "PUT", segments: []string{"distribution", "{Id}", "associate-web-acl"}, payloadRoot: "AssociateDistributionWebACLResult"},
	{op: "DisassociateDistributionWebACL", method: "PUT", segments: []string{"distribution", "{Id}", "disassociate-web-acl"}, payloadRoot: "DisassociateDistributionWebACLResult"},
	{op: "CopyDistribution", method: "POST", segments: []string{"distribution", "{Id}", "copy"}, payloadRoot: "Distribution"},
	{op: "GetDistribution", method: "GET", segments: []string{"distribution", "{Id}"}, payloadRoot: "Distribution"},
	{op: "DeleteDistribution", method: "DELETE", segments: []string{"distribution", "{Id}"}, payloadRoot: ""},

	// Distribution lookups by related resource
	{op: "ListDistributionsByWebACLId", method: "GET", segments: []string{"distributionsByWebACLId", "{WebACLId}"}, payloadRoot: "DistributionList"},
	{op: "ListDistributionsByCachePolicyId", method: "GET", segments: []string{"distributionsByCachePolicyId", "{CachePolicyId}"}, payloadRoot: "DistributionIdList"},
	{op: "ListDistributionsByKeyGroup", method: "GET", segments: []string{"distributionsByKeyGroupId", "{KeyGroupId}"}, payloadRoot: "DistributionIdList"},
	{op: "ListDistributionsByOriginRequestPolicyId", method: "GET", segments: []string{"distributionsByOriginRequestPolicyId", "{OriginRequestPolicyId}"}, payloadRoot: "DistributionIdList"},
	{op: "ListDistributionsByResponseHeadersPolicyId", method: "GET", segments: []string{"distributionsByResponseHeadersPolicyId", "{ResponseHeadersPolicyId}"}, payloadRoot: "DistributionIdList"},

	// Tagging. Tag and Untag requests arrive as POST with an Operation
	// query parameter per the API model; ListTagsForResource is plain GET.
	{op: "TagResource", method: "POST", segments: []string{"tagging"}, query: map[string]string{"Operation": "Tag"}, payloadRoot: ""},
	{op: "UntagResource", method: "POST", segments: []string{"tagging"}, query: map[string]string{"Operation": "Untag"}, payloadRoot: ""},
	{op: "ListTagsForResource", method: "GET", segments: []string{"tagging"}, payloadRoot: "Tags"},

	// Cache policy
	{op: "CreateCachePolicy", method: "POST", segments: []string{"cache-policy"}, payloadRoot: "CachePolicy"},
	{op: "ListCachePolicies", method: "GET", segments: []string{"cache-policy"}, payloadRoot: "CachePolicyList"},
	{op: "GetCachePolicyConfig", method: "GET", segments: []string{"cache-policy", "{Id}", "config"}, payloadRoot: "CachePolicyConfig"},
	{op: "GetCachePolicy", method: "GET", segments: []string{"cache-policy", "{Id}"}, payloadRoot: "CachePolicy"},
	{op: "UpdateCachePolicy", method: "PUT", segments: []string{"cache-policy", "{Id}"}, payloadRoot: "CachePolicy"},
	{op: "DeleteCachePolicy", method: "DELETE", segments: []string{"cache-policy", "{Id}"}, payloadRoot: ""},

	// Origin request policy
	{op: "CreateOriginRequestPolicy", method: "POST", segments: []string{"origin-request-policy"}, payloadRoot: "OriginRequestPolicy"},
	{op: "ListOriginRequestPolicies", method: "GET", segments: []string{"origin-request-policy"}, payloadRoot: "OriginRequestPolicyList"},
	{op: "GetOriginRequestPolicyConfig", method: "GET", segments: []string{"origin-request-policy", "{Id}", "config"}, payloadRoot: "OriginRequestPolicyConfig"},
	{op: "GetOriginRequestPolicy", method: "GET", segments: []string{"origin-request-policy", "{Id}"}, payloadRoot: "OriginRequestPolicy"},
	{op: "UpdateOriginRequestPolicy", method: "PUT", segments: []string{"origin-request-policy", "{Id}"}, payloadRoot: "OriginRequestPolicy"},
	{op: "DeleteOriginRequestPolicy", method: "DELETE", segments: []string{"origin-request-policy", "{Id}"}, payloadRoot: ""},

	// Origin access control
	{op: "CreateOriginAccessControl", method: "POST", segments: []string{"origin-access-control"}, payloadRoot: "OriginAccessControl"},
	{op: "ListOriginAccessControls", method: "GET", segments: []string{"origin-access-control"}, payloadRoot: "OriginAccessControlList"},
	{op: "GetOriginAccessControlConfig", method: "GET", segments: []string{"origin-access-control", "{Id}", "config"}, payloadRoot: "OriginAccessControlConfig"},
	{op: "UpdateOriginAccessControl", method: "PUT", segments: []string{"origin-access-control", "{Id}", "config"}, payloadRoot: "OriginAccessControl"},
	{op: "GetOriginAccessControl", method: "GET", segments: []string{"origin-access-control", "{Id}"}, payloadRoot: "OriginAccessControl"},
	{op: "DeleteOriginAccessControl", method: "DELETE", segments: []string{"origin-access-control", "{Id}"}, payloadRoot: ""},

	// Key group
	{op: "CreateKeyGroup", method: "POST", segments: []string{"key-group"}, payloadRoot: "KeyGroup"},
	{op: "ListKeyGroups", method: "GET", segments: []string{"key-group"}, payloadRoot: "KeyGroupList"},
	{op: "GetKeyGroupConfig", method: "GET", segments: []string{"key-group", "{Id}", "config"}, payloadRoot: "KeyGroupConfig"},
	{op: "GetKeyGroup", method: "GET", segments: []string{"key-group", "{Id}"}, payloadRoot: "KeyGroup"},
	{op: "UpdateKeyGroup", method: "PUT", segments: []string{"key-group", "{Id}"}, payloadRoot: "KeyGroup"},
	{op: "DeleteKeyGroup", method: "DELETE", segments: []string{"key-group", "{Id}"}, payloadRoot: ""},

	// Response headers policy
	{op: "CreateResponseHeadersPolicy", method: "POST", segments: []string{"response-headers-policy"}, payloadRoot: "ResponseHeadersPolicy"},
	{op: "ListResponseHeadersPolicies", method: "GET", segments: []string{"response-headers-policy"}, payloadRoot: "ResponseHeadersPolicyList"},
	{op: "GetResponseHeadersPolicyConfig", method: "GET", segments: []string{"response-headers-policy", "{Id}", "config"}, payloadRoot: "ResponseHeadersPolicyConfig"},
	{op: "GetResponseHeadersPolicy", method: "GET", segments: []string{"response-headers-policy", "{Id}"}, payloadRoot: "ResponseHeadersPolicy"},
	{op: "UpdateResponseHeadersPolicy", method: "PUT", segments: []string{"response-headers-policy", "{Id}"}, payloadRoot: "ResponseHeadersPolicy"},
	{op: "DeleteResponseHeadersPolicy", method: "DELETE", segments: []string{"response-headers-policy", "{Id}"}, payloadRoot: ""},

	// Public key
	{op: "CreatePublicKey", method: "POST", segments: []string{"public-key"}, payloadRoot: "PublicKey"},
	{op: "ListPublicKeys", method: "GET", segments: []string{"public-key"}, payloadRoot: "PublicKeyList"},
	{op: "GetPublicKeyConfig", method: "GET", segments: []string{"public-key", "{Id}", "config"}, payloadRoot: "PublicKeyConfig"},
	{op: "UpdatePublicKey", method: "PUT", segments: []string{"public-key", "{Id}", "config"}, payloadRoot: "PublicKey"},
	{op: "GetPublicKey", method: "GET", segments: []string{"public-key", "{Id}"}, payloadRoot: "PublicKey"},
	{op: "DeletePublicKey", method: "DELETE", segments: []string{"public-key", "{Id}"}, payloadRoot: ""},

	// Continuous deployment policy
	{op: "CreateContinuousDeploymentPolicy", method: "POST", segments: []string{"continuous-deployment-policy"}, payloadRoot: "ContinuousDeploymentPolicy"},
	{op: "ListContinuousDeploymentPolicies", method: "GET", segments: []string{"continuous-deployment-policy"}, payloadRoot: "ContinuousDeploymentPolicyList"},
	{op: "GetContinuousDeploymentPolicyConfig", method: "GET", segments: []string{"continuous-deployment-policy", "{Id}", "config"}, payloadRoot: "ContinuousDeploymentPolicyConfig"},
	{op: "GetContinuousDeploymentPolicy", method: "GET", segments: []string{"continuous-deployment-policy", "{Id}"}, payloadRoot: "ContinuousDeploymentPolicy"},
	{op: "UpdateContinuousDeploymentPolicy", method: "PUT", segments: []string{"continuous-deployment-policy", "{Id}"}, payloadRoot: "ContinuousDeploymentPolicy"},
	{op: "DeleteContinuousDeploymentPolicy", method: "DELETE", segments: []string{"continuous-deployment-policy", "{Id}"}, payloadRoot: ""},
}

// matchCloudFrontRoute resolves a CloudFront request against the routing
// table. It returns the operation name and the captured URI labels keyed by
// parameter name; both are zero values when no route matches. The path must
// be the escaped request path so percent-encoded slashes inside URI labels
// survive segment splitting.
func matchCloudFrontRoute(method, escapedPath string, query url.Values) (string, map[string]string) {
	trimmed := strings.TrimPrefix(escapedPath, "/2020-05-31/")
	parts := strings.Split(trimmed, "/")
	if len(parts) > 1 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	for i := range cloudFrontRoutes {
		route := &cloudFrontRoutes[i]
		if route.method != method || len(route.segments) != len(parts) {
			continue
		}
		if !cloudFrontQueryMatches(route.query, query) {
			continue
		}
		labels := make(map[string]string)
		matched := true
		for j, seg := range route.segments {
			if name, ok := cloudFrontLabelName(seg); ok {
				if parts[j] == "" {
					matched = false
					break
				}
				labels[name] = cloudFrontPathValue(parts[j])
			} else if seg != parts[j] {
				matched = false
				break
			}
		}
		if matched {
			return route.op, labels
		}
	}
	return "", nil
}

// cloudFrontQueryMatches reports whether the request query satisfies a
// route's query constraints.
func cloudFrontQueryMatches(constraints map[string]string, query url.Values) bool {
	for key, value := range constraints {
		if value == "" {
			if !query.Has(key) {
				return false
			}
		} else if query.Get(key) != value {
			return false
		}
	}
	return true
}

// cloudFrontLabelName reports whether a template segment is a {Name}
// placeholder and returns the label name.
func cloudFrontLabelName(seg string) (string, bool) {
	if len(seg) > 1 && seg[0] == '{' && seg[len(seg)-1] == '}' {
		return seg[1 : len(seg)-1], true
	}
	return "", false
}

// CloudFrontPayloadRoot returns the XML payload root element configured for
// a CloudFront operation and whether the operation appears in the routing
// table. An empty root with ok=true marks an operation whose response has no
// body.
func CloudFrontPayloadRoot(opName string) (string, bool) {
	for i := range cloudFrontRoutes {
		if cloudFrontRoutes[i].op == opName {
			return cloudFrontRoutes[i].payloadRoot, true
		}
	}
	return "", false
}

// CloudFrontRouteInfo is the exported view of one routing table row, for
// tests that verify the table against the API model and the handler
// registry.
type CloudFrontRouteInfo struct {
	Op          string
	Method      string
	Segments    []string
	Query       map[string]string
	PayloadRoot string
}

// CloudFrontRouteTable returns a copy of the CloudFront routing table.
func CloudFrontRouteTable() []CloudFrontRouteInfo {
	rows := make([]CloudFrontRouteInfo, 0, len(cloudFrontRoutes))
	for i := range cloudFrontRoutes {
		route := cloudFrontRoutes[i]
		rows = append(rows, CloudFrontRouteInfo{
			Op:          route.op,
			Method:      route.method,
			Segments:    append([]string(nil), route.segments...),
			Query:       route.query,
			PayloadRoot: route.payloadRoot,
		})
	}
	return rows
}
