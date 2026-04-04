package request

// AppSync miscellaneous parsers — tags, merged APIs, source APIs, data source
// introspection, and data-plane evaluation endpoints.
// These path groups are smaller and do not warrant individual files.

import (
	"net/http"
	"strings"
)

// --- Tags ---

func extractAppSyncTagsOperation(path, method string) string {
	rel := strings.TrimPrefix(path, "/v1/tags")
	if rel == "" || rel == "/" {
		return ""
	}

	switch method {
	case http.MethodPost:
		return "TagResource"
	case http.MethodDelete:
		return "UntagResource"
	case http.MethodGet:
		return "ListTagsForResource"
	}

	return ""
}

func extractAppSyncTagsPathParams(path string, params map[string]interface{}) {
	rel := strings.TrimPrefix(path, "/v1/tags/")
	if rel != "" && rel != path {
		params["resourceArn"] = rel
	}
}

// --- Merged APIs (source API associations from the merged side) ---
// Path format: /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations[/{associationId}][/{sub-resource}]
// After trimming the prefix, parts[0] is mergedApiIdentifier, parts[1] is "sourceApiAssociations",
// parts[2] is associationId (if present), and parts[3] is the sub-resource (e.g. "types", "merge").

func extractAppSyncMergedApisOperation(path, method string) string {
	rel := strings.TrimPrefix(path, "/v1/mergedApis")
	if rel == "" || rel == "/" {
		return ""
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) < 2 || parts[1] != "sourceApiAssociations" {
		return ""
	}

	if len(parts) == 2 {
		if method == http.MethodPost {
			return "AssociateSourceGraphqlApi"
		}
		return ""
	}

	if len(parts) == 3 {
		switch method {
		case http.MethodGet:
			return "GetSourceApiAssociation"
		case http.MethodPost:
			return "UpdateSourceApiAssociation"
		case http.MethodDelete:
			return "DisassociateSourceGraphqlApi"
		}
		return ""
	}

	if len(parts) == 4 && parts[3] == "types" {
		if method == http.MethodGet {
			return "ListTypesByAssociation"
		}
	}

	if len(parts) == 4 && parts[3] == "merge" {
		if method == http.MethodPost {
			return "StartSchemaMerge"
		}
	}

	return ""
}

func extractAppSyncMergedApisPathParams(path string, params map[string]interface{}) {
	rel := strings.TrimPrefix(path, "/v1/mergedApis")
	if rel == "" || rel == "/" {
		return
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) >= 1 && parts[0] != "" {
		params["mergedApiIdentifier"] = parts[0]
		if len(parts) >= 3 && parts[1] == "sourceApiAssociations" && parts[2] != "" {
			params["associationId"] = parts[2]
		}
	}
}

// --- Source APIs (merged API associations from the source side) ---
// Path format: /v1/sourceApis/{sourceApiIdentifier}/mergedApiAssociations[/{associationId}]
// After trimming the prefix, parts[0] is sourceApiIdentifier, parts[1] is "mergedApiAssociations",
// and parts[2] is associationId (if present).

func extractAppSyncSourceApisOperation(path, method string) string {
	rel := strings.TrimPrefix(path, "/v1/sourceApis")
	if rel == "" || rel == "/" {
		return ""
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) < 2 || parts[1] != "mergedApiAssociations" {
		return ""
	}

	if len(parts) == 2 {
		if method == http.MethodPost {
			return "AssociateMergedGraphqlApi"
		}
		return ""
	}

	if len(parts) == 3 {
		if method == http.MethodDelete {
			return "DisassociateMergedGraphqlApi"
		}
	}

	return ""
}

func extractAppSyncSourceApisPathParams(path string, params map[string]interface{}) {
	rel := strings.TrimPrefix(path, "/v1/sourceApis")
	if rel == "" || rel == "/" {
		return
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) >= 1 && parts[0] != "" {
		params["sourceApiIdentifier"] = parts[0]
		if len(parts) >= 3 && parts[1] == "mergedApiAssociations" && parts[2] != "" {
			params["associationId"] = parts[2]
		}
	}
}

// --- Data source introspection ---

func extractAppSyncDataSourceIntrospectionOperation(path, method string) string {
	rel := strings.TrimPrefix(path, "/v1/datasources/introspections")
	if rel == "" || rel == "/" {
		if method == http.MethodPost {
			return "StartDataSourceIntrospection"
		}
		return ""
	}

	if method == http.MethodGet {
		return "GetDataSourceIntrospection"
	}

	return ""
}

func extractAppSyncIntrospectionPathParams(path string, params map[string]interface{}) {
	rel := strings.TrimPrefix(path, "/v1/datasources/introspections/")
	if rel != "" && rel != path {
		params["introspectionId"] = rel
	}
}
