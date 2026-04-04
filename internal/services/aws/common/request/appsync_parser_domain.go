package request

// AppSync domain name parser — handles /v1/domainnames/... paths.
// Covers domain name CRUD and API association (associate/dissociate a GraphQL API
// to a custom domain). Each domain name has at most one API association.

import (
	"net/http"
	"strings"
)

// extractAppSyncDomainNamesOperation handles domain name paths (/v1/domainnames/...).
func extractAppSyncDomainNamesOperation(path, method string) string {
	rel := strings.TrimPrefix(path, "/v1/domainnames")
	if rel == "" || rel == "/" {
		switch method {
		case http.MethodGet:
			return "ListDomainNames"
		case http.MethodPost:
			return "CreateDomainName"
		}
		return ""
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) == 1 {
		switch method {
		case http.MethodGet:
			return "GetDomainName"
		case http.MethodPost:
			return "UpdateDomainName"
		case http.MethodDelete:
			return "DeleteDomainName"
		}
		return ""
	}

	if len(parts) >= 2 && parts[1] == "apiassociation" {
		switch method {
		case http.MethodGet:
			return "GetApiAssociation"
		case http.MethodPost:
			return "AssociateApi"
		case http.MethodDelete:
			return "DisassociateApi"
		}
	}

	return ""
}

// extractAppSyncDomainNamesPathParams extracts path parameters for domain name requests.
func extractAppSyncDomainNamesPathParams(path string, params map[string]interface{}) {
	rel := strings.TrimPrefix(path, "/v1/domainnames")
	if rel == "" || rel == "/" {
		return
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) >= 1 && parts[0] != "" {
		params["domainName"] = parts[0]
	}
}
