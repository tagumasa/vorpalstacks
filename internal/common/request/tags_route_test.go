package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// The bare /tags/{arn} prefix is shared by API Gateway, EventBridge Scheduler
// and Neptune Graph; each service must extract operations only for ARNs whose
// service field is its own.
func TestTagsRouteOperationOwnership(t *testing.T) {
	cases := []struct {
		method string
		path   string
		op     string
	}{
		{http.MethodGet, "/tags/arn:aws:apigateway:us-east-1::/restapis/abc123", "GetTags"},
		{http.MethodPut, "/tags/arn:aws:apigateway:us-east-1::/restapis/abc123", "TagResource"},
		// TagResource is bound to PUT in the model; POST is no operation.
		{http.MethodPost, "/tags/arn:aws:apigateway:us-east-1::/restapis/abc123", ""},
		{http.MethodDelete, "/tags/arn:aws:apigateway:us-east-1::/restapis/abc123", "UntagResource"},
		{http.MethodGet, "/tags/arn:aws:scheduler:us-east-1:111122223333:schedule-group/default", "ListTagsForResource"},
		{http.MethodPost, "/tags/arn:aws:scheduler:us-east-1:111122223333:schedule-group/default", "TagResource"},
		{http.MethodDelete, "/tags/arn:aws:scheduler:us-east-1:111122223333:schedule-group/default", "UntagResource"},
		{http.MethodGet, "/tags/arn:aws:scheduler:us-east-1:111122223333:schedule/probe-sched", "ListTagsForResource"},
		{http.MethodGet, "/tags/arn:aws:neptune-graph:us-east-1:111122223333:graph/g-001", "ListTagsForResource"},
		{http.MethodPost, "/tags/arn:aws:neptune-graph:us-east-1:111122223333:graph/g-001", "TagResource"},
	}
	for _, tc := range cases {
		r := httptest.NewRequest(tc.method, tc.path, nil)
		if got := extractOperation(r, nil); got != tc.op {
			t.Errorf("extractOperation(%s %s) = %q, want %q", tc.method, tc.path, got, tc.op)
		}
	}
}

// An unparseable ARN after /tags/ is claimed by no service.
func TestTagsRouteUnparseableARNUnclaimed(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/tags/not-an-arn", nil)
	if got := extractOperation(r, nil); got != "" {
		t.Errorf("extractOperation = %q, want empty for an unparseable ARN", got)
	}
}

// The scheduler parser supplies the resourceArn path parameter itself; the
// API Gateway parser no longer sets it for foreign ARNs.
func TestSchedulerTagsPathParams(t *testing.T) {
	arn := "arn:aws:scheduler:us-east-1:111122223333:schedule-group/default"
	r := httptest.NewRequest(http.MethodGet, "/tags/"+arn, nil)

	params := map[string]interface{}{}
	for _, p := range restParsers {
		if p.MatchPath(r.URL.Path) {
			p.ExtractPathParams(r, params)
		}
	}
	if got, _ := params["resourceArn"].(string); got != arn {
		t.Errorf("resourceArn = %v, want %q", params["resourceArn"], arn)
	}
}

// API Gateway tag routes keep their resourceArn parameter.
func TestApiGatewayTagsPathParams(t *testing.T) {
	arn := "arn:aws:apigateway:us-east-1::/restapis/abc123"
	r := httptest.NewRequest(http.MethodGet, "/tags/"+arn, nil)

	params := map[string]interface{}{}
	for _, p := range restParsers {
		if p.MatchPath(r.URL.Path) {
			p.ExtractPathParams(r, params)
		}
	}
	if got, _ := params["resourceArn"].(string); got != arn {
		t.Errorf("resourceArn = %v, want %q", params["resourceArn"], arn)
	}
}
