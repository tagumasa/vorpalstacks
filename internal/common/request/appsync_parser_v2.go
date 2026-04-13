package request

// AppSync v2 Event API parser — handles /v2/apis/... paths.
// The v2 API (AppSync Events) covers API and ChannelNamespace CRUD operations.
// Each API contains channel namespaces that group event channels.

import (
	"net/http"
	"strings"
)

// extractAppSyncV2Operation handles v2 Event API paths (/v2/apis/...).
func extractAppSyncV2Operation(path, method string) string {
	rel := strings.TrimPrefix(path, "/v2/apis")
	if rel == "" || rel == "/" {
		switch method {
		case http.MethodGet:
			return "ListApis"
		case http.MethodPost:
			return "CreateApi"
		}
		return ""
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) == 1 {
		switch method {
		case http.MethodGet:
			return "GetApi"
		case http.MethodPost:
			return "UpdateApi"
		case http.MethodDelete:
			return "DeleteApi"
		}
	}

	if len(parts) >= 2 && parts[1] == "channelNamespaces" {
		if len(parts) == 2 {
			switch method {
			case http.MethodGet:
				return "ListChannelNamespaces"
			case http.MethodPost:
				return "CreateChannelNamespace"
			}
		} else {
			switch method {
			case http.MethodGet:
				return "GetChannelNamespace"
			case http.MethodPost:
				return "UpdateChannelNamespace"
			case http.MethodDelete:
				return "DeleteChannelNamespace"
			}
		}
	}

	return ""
}

// extractAppSyncV2PathParams extracts path parameters for v2 Event API requests.
func extractAppSyncV2PathParams(path string, params map[string]interface{}) {
	rel := strings.TrimPrefix(path, "/v2/apis")
	if rel == "" || rel == "/" {
		return
	}

	rel = strings.TrimPrefix(rel, "/")
	parts := strings.Split(rel, "/")

	if len(parts) >= 1 && parts[0] != "" {
		params["apiId"] = parts[0]
		if len(parts) >= 3 && parts[1] == "channelNamespaces" && parts[2] != "" {
			params["name"] = parts[2]
		}
	}
}
