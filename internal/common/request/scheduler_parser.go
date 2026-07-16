package request

import (
	"net/http"
	"strings"
)

func extractSchedulerOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method

	if !strings.HasPrefix(path, "/schedule-groups") && !strings.HasPrefix(path, "/schedules") {
		return ""
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) >= 1 && parts[0] == "schedule-groups" {
		if len(parts) >= 2 && parts[1] != "" {
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
			if len(parts) >= 3 && parts[2] == "tags" {
				switch method {
				case http.MethodPost:
					return "TagResource"
				case http.MethodDelete:
					return "UntagResource"
				case http.MethodGet:
					return "ListTagsForResource"
				}
			} else {
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
			}
		} else if method == http.MethodGet {
			return "ListSchedules"
		}
	}

	return ""
}

func extractSchedulerPathParams(path string, method string, params map[string]interface{}) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) >= 2 && parts[0] == "schedule-groups" {
		if parts[1] != "" {
			params["Name"] = parts[1]
			if method == http.MethodPost {
				params["_operation"] = "CreateScheduleGroup"
			} else if method == http.MethodGet {
				params["_operation"] = "GetScheduleGroup"
			} else if method == http.MethodDelete {
				params["_operation"] = "DeleteScheduleGroup"
			}
		} else if method == http.MethodGet {
			params["_operation"] = "ListScheduleGroups"
		}
	}

	if len(parts) >= 2 && parts[0] == "schedules" {
		if parts[1] != "" {
			params["Name"] = parts[1]
			if len(parts) >= 3 && parts[2] == "tags" {
				if method == http.MethodPost {
					params["_operation"] = "TagResource"
				} else if method == http.MethodDelete {
					params["_operation"] = "UntagResource"
				} else if method == http.MethodGet {
					params["_operation"] = "ListTagsForResource"
				}
			} else {
				if method == http.MethodPost {
					params["_operation"] = "CreateSchedule"
				} else if method == http.MethodGet {
					params["_operation"] = "GetSchedule"
				} else if method == http.MethodPut {
					params["_operation"] = "UpdateSchedule"
				} else if method == http.MethodDelete {
					params["_operation"] = "DeleteSchedule"
				}
			}
		}
	}
}

// schedulerRESTParser implements RESTServiceParser for Amazon EventBridge Scheduler.
type schedulerRESTParser struct{}

// MatchPath returns true if the path belongs to EventBridge Scheduler.
func (p *schedulerRESTParser) MatchPath(path string) bool {
	return strings.HasPrefix(path, "/schedule-groups") || strings.HasPrefix(path, "/schedules")
}

// ExtractOperation returns the Scheduler operation name, or empty if the path does not match.
func (p *schedulerRESTParser) ExtractOperation(r *http.Request) string {
	return extractSchedulerOperation(r)
}

// ExtractPathParams extracts URI-bound parameters from the Scheduler request path.
func (p *schedulerRESTParser) ExtractPathParams(r *http.Request, params map[string]interface{}) {
	extractSchedulerPathParams(r.URL.Path, r.Method, params)
}
