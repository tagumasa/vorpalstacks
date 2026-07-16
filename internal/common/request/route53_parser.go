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
		} else if len(parts) >= 3 {
			if parts[2] == "rrset" {
				switch method {
				case "POST":
					return "ChangeResourceRecordSets"
				case "GET":
					return "ListResourceRecordSets"
				}
			}
			if parts[2] == "associatevpc" {
				return "AssociateVPCWithHostedZone"
			}
			if parts[2] == "disassociatevpc" {
				return "DisassociateVPCFromHostedZone"
			}
		}
		if len(parts) >= 2 {
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
		if len(parts) >= 2 && parts[1] != "" {
			resourceType := parts[1]
			if _, ok := params["ResourceType"]; !ok {
				params["ResourceType"] = resourceType
			}
			remaining := strings.Join(parts[2:], "/")
			if remaining != "" {
				if _, ok := params["ResourceId"]; !ok {
					params["ResourceId"] = remaining
				}
			}
		}
	}
}

// route53RESTParser implements RESTServiceParser for Amazon Route 53.
type route53RESTParser struct{}

// MatchPath returns true if the path belongs to Route 53.
func (p *route53RESTParser) MatchPath(path string) bool {
	return strings.HasPrefix(path, "/2013-04-01/")
}

// ExtractOperation returns the Route 53 operation name, or empty if the path does not match.
func (p *route53RESTParser) ExtractOperation(r *http.Request) string {
	return extractRoute53Operation(r)
}

// ExtractPathParams extracts URI-bound parameters from the Route 53 request path.
func (p *route53RESTParser) ExtractPathParams(r *http.Request, params map[string]interface{}) {
	extractRoute53PathParams(r.URL.Path, params)
}
