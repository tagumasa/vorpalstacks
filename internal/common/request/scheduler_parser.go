package request

import (
	"net/http"
	"strings"
)

func extractSchedulerOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if TagsRouteService(path) == "scheduler" {
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

	if !strings.HasPrefix(path, "/schedule-groups") && !strings.HasPrefix(path, "/schedules") {
		return ""
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) >= 1 && parts[0] == "schedule-groups" {
		if len(parts) >= 2 && parts[1] != "" {
			// The model defines no sub-resource under
			// /schedule-groups/{Name} either — deeper paths stay unrouted
			// like the schedules branch below.
			if len(parts) >= 3 {
				return ""
			}
			switch method {
			case http.MethodPost:
				return "CreateScheduleGroup"
			case http.MethodGet:
				return "GetScheduleGroup"
			case http.MethodDelete:
				return "DeleteScheduleGroup"
			}
		} else if method == http.MethodGet {
			return "ListScheduleGroups"
		}
	}

	if len(parts) >= 1 && parts[0] == "schedules" {
		if len(parts) >= 2 && parts[1] != "" {
			if len(parts) >= 3 {
				// The model defines no sub-resource under /schedules/{Name};
				// the tag trio binds only to /tags/{ResourceArn+}.
				return ""
			}
			switch method {
			case http.MethodPost:
				return "CreateSchedule"
			case http.MethodGet:
				return "GetSchedule"
			case http.MethodPut:
				return "UpdateSchedule"
			case http.MethodDelete:
				return "DeleteSchedule"
			}
		} else if method == http.MethodGet {
			return "ListSchedules"
		}
	}

	return ""
}

func extractSchedulerPathParams(path string, params map[string]interface{}) {
	if TagsRouteService(path) == "scheduler" {
		if _, ok := params["resourceArn"]; !ok {
			params["resourceArn"] = strings.TrimPrefix(path, "/tags/")
		}
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")

	// The depth guards mirror extractSchedulerOperation: the model defines
	// no sub-resource under /schedule-groups/{Name} or /schedules/{Name},
	// so deeper paths stay unrouted and bind no path parameter.
	if len(parts) >= 2 && parts[0] == "schedule-groups" {
		if len(parts) == 2 && parts[1] != "" {
			params["Name"] = parts[1]
		}
	}

	if len(parts) >= 2 && parts[0] == "schedules" {
		if len(parts) == 2 && parts[1] != "" {
			params["Name"] = parts[1]
		}
	}
}

// schedulerRESTParser implements RESTServiceParser for Amazon EventBridge Scheduler.
type schedulerRESTParser struct{}

// MatchPath returns true if the path belongs to EventBridge Scheduler.
func (p *schedulerRESTParser) MatchPath(path string) bool {
	return strings.HasPrefix(path, "/schedule-groups") || strings.HasPrefix(path, "/schedules") ||
		TagsRouteService(path) == "scheduler"
}

// ExtractOperation returns the Scheduler operation name, or empty if the path does not match.
func (p *schedulerRESTParser) ExtractOperation(r *http.Request) string {
	return extractSchedulerOperation(r)
}

// ExtractPathParams extracts URI-bound parameters from the Scheduler request path.
func (p *schedulerRESTParser) ExtractPathParams(r *http.Request, params map[string]interface{}) {
	extractSchedulerPathParams(r.URL.Path, params)
}
