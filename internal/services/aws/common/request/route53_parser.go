package request

import (
	"net/http"
	"strings"
)

func extractRoute53Operation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if !strings.HasPrefix(path, "/2013-04-01/") {
		return ""
	}

	path = strings.TrimPrefix(path, "/2013-04-01/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 {
		return ""
	}

	switch parts[0] {
	case "hostedzone":
		if len(parts) == 1 || parts[1] == "" {
			switch method {
			case "POST":
				return "CreateHostedZone"
			case "GET":
				return "ListHostedZones"
			}
		} else if len(parts) >= 2 {
			if len(parts) >= 3 && parts[2] == "rrset" {
				switch method {
				case "POST":
					return "ChangeResourceRecordSets"
				case "GET":
					return "ListResourceRecordSets"
				}
			}
			switch method {
			case "GET":
				return "GetHostedZone"
			case "DELETE":
				return "DeleteHostedZone"
			case "POST":
				return "UpdateHostedZoneComment"
			}
		}
	case "healthcheck":
		if len(parts) == 1 || parts[1] == "" {
			switch method {
			case "POST":
				return "CreateHealthCheck"
			case "GET":
				return "ListHealthChecks"
			}
		} else if len(parts) >= 2 {
			switch method {
			case "GET":
				return "GetHealthCheck"
			case "DELETE":
				return "DeleteHealthCheck"
			case "POST":
				return "UpdateHealthCheck"
			}
		}
	case "change":
		if len(parts) >= 2 && method == "GET" {
			return "GetChange"
		}
	case "hostedzonesbyname":
		if method == "GET" {
			return "ListHostedZonesByName"
		}
	case "delegationset":
		if len(parts) == 1 || parts[1] == "" {
			if method == "GET" {
				return "ListReusableDelegationSets"
			}
		}
	case "tags":
		if len(parts) >= 2 {
			switch method {
			case "POST":
				return "ChangeTagsForResource"
			case "GET":
				return "ListTagsForResource"
			}
		}
	}

	return ""
}

func extractRoute53PathParams(path string, params map[string]interface{}) {
	path = strings.TrimPrefix(path, "/2013-04-01/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		return
	}

	switch parts[0] {
	case "hostedzone":
		if len(parts) >= 2 && parts[1] != "" {
			if _, ok := params["Id"]; !ok {
				params["Id"] = "/hostedzone/" + parts[1]
			}
			if _, ok := params["HostedZoneId"]; !ok {
				params["HostedZoneId"] = "/hostedzone/" + parts[1]
			}
		}
	case "healthcheck":
		if len(parts) >= 2 && parts[1] != "" {
			if _, ok := params["HealthCheckId"]; !ok {
				params["HealthCheckId"] = parts[1]
			}
		}
	case "change":
		if len(parts) >= 2 && parts[1] != "" {
			if _, ok := params["Id"]; !ok {
				params["Id"] = "/change/" + parts[1]
			}
		}
	case "tags":
		if len(parts) >= 3 && parts[1] != "" && parts[2] != "" {
			if _, ok := params["ResourceType"]; !ok {
				params["ResourceType"] = parts[1]
			}
			if _, ok := params["ResourceId"]; !ok {
				params["ResourceId"] = parts[2]
			}
		}
	}
}
