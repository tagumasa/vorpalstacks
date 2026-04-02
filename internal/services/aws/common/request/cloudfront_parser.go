package request

import (
	"net/http"
	"strings"
)

func extractCloudFrontOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if !strings.HasPrefix(path, "/2020-05-31/") {
		return ""
	}

	path = strings.TrimPrefix(path, "/2020-05-31/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 {
		return ""
	}

	switch parts[0] {
	case "distribution":
		if len(parts) == 1 || parts[1] == "" {
			switch method {
			case "POST":
				if r.URL.Query().Has("WithTags") {
					return "CreateDistributionWithTags"
				}
				return "CreateDistribution"
			case "GET":
				return "ListDistributions"
			}
		} else if len(parts) >= 2 {
			if len(parts) >= 3 && parts[2] == "config" {
				switch method {
				case "GET":
					return "GetDistributionConfig"
				}
			} else if len(parts) >= 3 && parts[2] == "invalidation" {
				switch method {
				case "POST":
					return "CreateInvalidation"
				case "GET":
					if len(parts) >= 4 && parts[3] != "" {
						return "GetInvalidation"
					}
					return "ListInvalidations"
				}
			}
			switch method {
			case "GET":
				return "GetDistribution"
			case "PUT":
				return "UpdateDistribution"
			case "DELETE":
				return "DeleteDistribution"
			}
		}
	case "tagging":
		op := r.URL.Query().Get("Operation")
		switch op {
		case "Tag":
			return "TagResource"
		case "Untag":
			return "UntagResource"
		default:
			switch method {
			case "POST":
				return "TagResource"
			case "DELETE":
				return "UntagResource"
			case "GET":
				return "ListTagsForResource"
			}
		}
	case "cache-policy":
		if len(parts) == 1 || parts[1] == "" {
			switch method {
			case "POST":
				return "CreateCachePolicy"
			case "GET":
				return "ListCachePolicies"
			}
		} else if len(parts) >= 2 {
			if len(parts) >= 3 && parts[2] == "config" {
				switch method {
				case "GET":
					return "GetCachePolicyConfig"
				}
			}
			switch method {
			case "GET":
				return "GetCachePolicy"
			case "PUT":
				return "UpdateCachePolicy"
			case "DELETE":
				return "DeleteCachePolicy"
			}
		}
	case "origin-request-policy":
		if len(parts) == 1 || parts[1] == "" {
			switch method {
			case "POST":
				return "CreateOriginRequestPolicy"
			case "GET":
				return "ListOriginRequestPolicies"
			}
		} else if len(parts) >= 2 {
			if len(parts) >= 3 && parts[2] == "config" {
				switch method {
				case "GET":
					return "GetOriginRequestPolicyConfig"
				}
			}
			switch method {
			case "GET":
				return "GetOriginRequestPolicy"
			case "PUT":
				return "UpdateOriginRequestPolicy"
			case "DELETE":
				return "DeleteOriginRequestPolicy"
			}
		}
	case "origin-access-control":
		if len(parts) == 1 || parts[1] == "" {
			switch method {
			case "POST":
				return "CreateOriginAccessControl"
			case "GET":
				return "ListOriginAccessControls"
			}
		} else if len(parts) >= 2 {
			if len(parts) >= 3 && parts[2] == "config" {
				switch method {
				case "PUT":
					return "UpdateOriginAccessControl"
				case "GET":
					return "GetOriginAccessControlConfig"
				}
			} else {
				switch method {
				case "GET":
					return "GetOriginAccessControl"
				case "DELETE":
					return "DeleteOriginAccessControl"
				}
			}
		}
	case "response-headers-policy":
		if len(parts) == 1 || parts[1] == "" {
			switch method {
			case "POST":
				return "CreateResponseHeadersPolicy"
			case "GET":
				return "ListResponseHeadersPolicies"
			}
		} else if len(parts) >= 2 {
			if len(parts) >= 3 && parts[2] == "config" {
				switch method {
				case "GET":
					return "GetResponseHeadersPolicyConfig"
				}
			} else {
				switch method {
				case "GET":
					return "GetResponseHeadersPolicy"
				case "PUT":
					return "UpdateResponseHeadersPolicy"
				case "DELETE":
					return "DeleteResponseHeadersPolicy"
				}
			}
		}
	}

	return ""
}

func extractCloudFrontPathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/2020-05-31/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		return
	}

	switch parts[0] {
	case "distribution":
		if len(parts) >= 2 && parts[1] != "" {
			if _, ok := params["Id"]; !ok {
				params["Id"] = parts[1]
			}
		}
		if len(parts) >= 4 && parts[2] == "invalidation" && parts[3] != "" {
			if _, ok := params["invalidationId"]; !ok {
				params["invalidationId"] = parts[3]
			}
		}
	case "cache-policy":
		if len(parts) >= 2 && parts[1] != "" {
			if _, ok := params["Id"]; !ok {
				params["Id"] = parts[1]
			}
		}
	case "origin-request-policy":
		if len(parts) >= 2 && parts[1] != "" {
			if _, ok := params["Id"]; !ok {
				params["Id"] = parts[1]
			}
		}
	case "origin-access-control":
		if len(parts) >= 2 && parts[1] != "" {
			if _, ok := params["Id"]; !ok {
				params["Id"] = parts[1]
			}
		}
	case "response-headers-policy":
		if len(parts) >= 2 && parts[1] != "" {
			if _, ok := params["Id"]; !ok {
				params["Id"] = parts[1]
			}
		}
	case "tagging":
		if resource := params["Resource"]; resource == nil {
			if queryResource, ok := params["Resource"].(string); ok && queryResource != "" {
				params["Resource"] = queryResource
			}
			if queryArn, ok := params["ResourceARN"].(string); ok && queryArn != "" {
				params["ResourceARN"] = queryArn
			}
		}
	}
}
