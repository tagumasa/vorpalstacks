package request

import (
	"net/http"
	"testing"
)

// TestSchedulerExtractOperation pins the modelled Scheduler URI surface: the
// collection and item routes resolve to their operations, and the tag trio
// binds only to /tags/{ResourceArn+}. The model defines no sub-resource under
// /schedules/{Name}, so any deeper path stays unrouted instead of reaching a
// handler that can only fail on the missing members.
func TestSchedulerExtractOperation(t *testing.T) {
	parser := &schedulerRESTParser{}

	groupArn := "arn:aws:scheduler:us-east-1:111122223333:schedule-group/default"
	tests := []struct {
		method string
		path   string
		want   string
	}{
		{http.MethodPost, "/schedules/my-schedule", "CreateSchedule"},
		{http.MethodGet, "/schedules", "ListSchedules"},
		{http.MethodGet, "/schedules/my-schedule", "GetSchedule"},
		{http.MethodPut, "/schedules/my-schedule", "UpdateSchedule"},
		{http.MethodDelete, "/schedules/my-schedule", "DeleteSchedule"},
		{http.MethodPost, "/schedule-groups/my-group", "CreateScheduleGroup"},
		{http.MethodGet, "/schedule-groups/my-group", "GetScheduleGroup"},
		{http.MethodDelete, "/schedule-groups/my-group", "DeleteScheduleGroup"},
		{http.MethodGet, "/schedule-groups", "ListScheduleGroups"},
		{http.MethodPost, "/tags/" + groupArn, "TagResource"},
		{http.MethodDelete, "/tags/" + groupArn, "UntagResource"},
		{http.MethodGet, "/tags/" + groupArn, "ListTagsForResource"},

		// No modelled sub-resource exists under /schedules/{Name}, and the
		// collection itself carries no POST route.
		{http.MethodPost, "/schedules", ""},
		{http.MethodPost, "/schedules/my-schedule/tags", ""},
		{http.MethodGet, "/schedules/my-schedule/tags", ""},
		{http.MethodDelete, "/schedules/my-schedule/tags", ""},
		{http.MethodGet, "/schedules/my-schedule/other", ""},

		// No modelled sub-resource exists under /schedule-groups/{Name}
		// either — deeper paths stay unrouted, not routed to the group
		// operations on the truncated name.
		{http.MethodGet, "/schedule-groups/my-group/tags", ""},
		{http.MethodPost, "/schedule-groups/my-group/other", ""},
		{http.MethodDelete, "/schedule-groups/my-group/a/b", ""},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			got := parser.ExtractOperation(req)
			if got != tt.want {
				t.Errorf("ExtractOperation(%s %s) = %q, want %q", tt.method, tt.path, got, tt.want)
			}
		})
	}
}

// TestSchedulerExtractPathParams pins the parameter binding's depth guard:
// a two-segment resource path binds its Name, and the unrouted deeper paths
// bind nothing — the parameter extractor mirrors the operation extractor's
// refusal instead of binding a name the router will never use.
func TestSchedulerExtractPathParams(t *testing.T) {
	binds := []struct {
		path string
		want string
	}{
		{"/schedules/my-schedule", "my-schedule"},
		{"/schedule-groups/my-group", "my-group"},
	}
	for _, tt := range binds {
		params := map[string]interface{}{}
		extractSchedulerPathParams(tt.path, params)
		if params["Name"] != tt.want {
			t.Errorf("extractSchedulerPathParams(%q) bound Name = %v, want %q", tt.path, params["Name"], tt.want)
		}
	}

	for _, path := range []string{
		"/schedules/my-schedule/tags",
		"/schedules/my-schedule/other",
		"/schedule-groups/my-group/tags",
		"/schedule-groups/my-group/a/b",
	} {
		params := map[string]interface{}{}
		extractSchedulerPathParams(path, params)
		if _, ok := params["Name"]; ok {
			t.Errorf("extractSchedulerPathParams(%q) bound Name on an unrouted path", path)
		}
	}
}
