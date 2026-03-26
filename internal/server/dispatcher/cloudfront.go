package dispatcher

import (
	"net/http"
)

var cloudFrontPayloadOperations = map[string]string{
	"CreateDistribution":                       "Distribution",
	"GetDistribution":                          "Distribution",
	"GetDistributionConfig":                    "DistributionConfig",
	"UpdateDistribution":                       "Distribution",
	"CreateDistributionWithTags":               "Distribution",
	"ListDistributions":                        "DistributionList",
	"ListDistributionsByCachePolicyId":         "DistributionIdList",
	"ListDistributionsByKeyGroup":              "DistributionIdList",
	"ListDistributionsByOriginRequestPolicyId": "DistributionIdList",
	"ListDistributionsByRealtimeLogConfig":     "DistributionIdList",
	"ListDistributionsByAnycastIpListId":       "DistributionList",
	"CreateCachePolicy":                        "CachePolicy",
	"GetCachePolicy":                           "CachePolicy",
	"UpdateCachePolicy":                        "CachePolicy",
	"ListCachePolicies":                        "CachePolicyList",
	"CreateOriginRequestPolicy":                "OriginRequestPolicy",
	"GetOriginRequestPolicy":                   "OriginRequestPolicy",
	"UpdateOriginRequestPolicy":                "OriginRequestPolicy",
	"ListOriginRequestPolicies":                "OriginRequestPolicyList",
	"CreateCloudFrontOriginAccessIdentity":     "CloudFrontOriginAccessIdentity",
	"GetCloudFrontOriginAccessIdentity":        "CloudFrontOriginAccessIdentity",
	"UpdateCloudFrontOriginAccessIdentity":     "CloudFrontOriginAccessIdentity",
	"ListCloudFrontOriginAccessIdentities":     "CloudFrontOriginAccessIdentityList",
	"CreateInvalidation":                       "Invalidation",
	"GetInvalidation":                          "Invalidation",
	"ListInvalidations":                        "InvalidationList",
	"CreateFunction":                           "FunctionSummary",
	"GetFunction":                              "FunctionSummary",
	"UpdateFunction":                           "FunctionSummary",
	"ListFunctions":                            "FunctionList",
	"CreateKeyGroup":                           "KeyGroup",
	"GetKeyGroup":                              "KeyGroup",
	"UpdateKeyGroup":                           "KeyGroup",
	"ListKeyGroups":                            "KeyGroupList",
	"CreatePublicKey":                          "PublicKey",
	"GetPublicKey":                             "PublicKey",
	"UpdatePublicKey":                          "PublicKey",
	"ListPublicKeys":                           "PublicKeyList",
	"CreateRealtimeLogConfig":                  "RealtimeLogConfig",
	"GetRealtimeLogConfig":                     "RealtimeLogConfig",
	"UpdateRealtimeLogConfig":                  "RealtimeLogConfig",
	"ListRealtimeLogConfigs":                   "RealtimeLogConfigList",
	"CreateStreamingDistribution":              "StreamingDistribution",
	"GetStreamingDistribution":                 "StreamingDistribution",
	"UpdateStreamingDistribution":              "StreamingDistribution",
	"ListStreamingDistributions":               "StreamingDistributionList",
	"ListTagsForResource":                      "Tags",
	"CreateResponseHeadersPolicy":              "ResponseHeadersPolicy",
	"GetResponseHeadersPolicy":                 "ResponseHeadersPolicy",
	"GetResponseHeadersPolicyConfig":           "ResponseHeadersPolicyConfig",
	"UpdateResponseHeadersPolicy":              "ResponseHeadersPolicy",
	"DeleteResponseHeadersPolicy":              "",
	"ListResponseHeadersPolicies":              "ResponseHeadersPolicyList",
	"CreateOriginAccessControl":                "OriginAccessControl",
	"GetOriginAccessControl":                   "OriginAccessControl",
	"UpdateOriginAccessControl":                "OriginAccessControl",
	"ListOriginAccessControls":                 "OriginAccessControlList",
}

func isCloudFrontPayloadOperation(opName string) bool {
	_, ok := cloudFrontPayloadOperations[opName]
	return ok
}

func getCloudFrontPayloadRoot(opName string) string {
	if root, ok := cloudFrontPayloadOperations[opName]; ok {
		return root
	}
	return opName
}

func extractCloudFrontETag(w http.ResponseWriter, response interface{}, payloadRoot string) {
	m, ok := response.(map[string]interface{})
	if !ok {
		return
	}

	var etag string

	if inner, exists := m[payloadRoot]; exists {
		if innerMap, ok := inner.(map[string]interface{}); ok {
			if e, ok := innerMap["ETag"].(string); ok && e != "" {
				etag = e
			}
		}
	}

	if etag == "" {
		if e, ok := m["ETag"].(string); ok && e != "" {
			etag = e
		}
	}

	if etag != "" {
		w.Header().Set("ETag", etag)
	}
}
