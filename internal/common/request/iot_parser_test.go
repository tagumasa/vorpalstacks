package request

import (
	"net/http"
	"testing"
)

func TestMatchPathWordBoundary(t *testing.T) {
	parser := &iotRESTParser{}

	falsePaths := []string{
		"/thingsabc",
		"/thing-groups-archive",
		"/thing-types-foo",
		"/billing-groupsx",
		"/certificates-foo",
		"/keys-and-certificate-extended",
		"/policiesabc",
		"/policy-principals-extra",
		"/principal-policies-other",
		"/principals-thing",
		"/target-policies-more",
		"/attached-policies-extra",
		"/rulesengine",
		"/jobslist",
		"/endpoint-xyz",
		"/role-aliases-admin",
		"/tags-anything",
		"/untag-resource",
		"/authorizers-other",
		"/provisioning-templates-extended",
		"/provisioning-template-extended",
		"/domainConfigurations-other",
		"/domainConfiguration-other",
		"/indexing-configuration",
		"/active-violations-list",
		"/violation-events-list",
		"/behavior-model-training-sum",
		"/security-profiles-other",
		"/detector-models-other",
		"/detector-model-other",
		"/inputs-data",
		"/input-data",
		"/messages-queue",
		"/destinations-other",
		"/effective-policies-other",
	}

	for _, p := range falsePaths {
		t.Run("false/"+p, func(t *testing.T) {
			if parser.MatchPath(p) {
				t.Errorf("MatchPath(%q) = true, want false", p)
			}
		})
	}

	truePaths := []string{
		"/things",
		"/things/abc",
		"/thing-groups",
		"/thing-groups/myGroup",
		"/thing-types",
		"/thing-types/myType",
		"/billing-groups",
		"/billing-groups/myBillingGroup",
		"/certificates",
		"/certificates/cert-123",
		"/certificate/register",
		"/keys-and-certificate",
		"/policies",
		"/policies/myPolicy",
		"/policy-principals/abc",
		"/principal-policies/abc",
		"/principals/things",
		"/target-policies/abc",
		"/attached-policies/abc",
		"/rules",
		"/rules/myRule",
		"/jobs",
		"/jobs/myJob",
		"/endpoint",
		"/role-aliases",
		"/role-aliases/myAlias",
		"/tags",
		"/untag",
		"/authorizers",
		"/authorizer/abc",
		"/provisioning-templates",
		"/provisioning-templates/myTemplate",
		"/provisioning-template/abc",
		"/domainConfigurations",
		"/domainConfiguration/myDomain",
		"/indexing",
		"/active-violations",
		"/violation-events",
		"/behavior-model-training/summaries",
		"/security-profiles",
		"/security-profiles/myProfile",
		"/messages",
		"/destinations",
		"/effective-policies",
	}

	for _, p := range truePaths {
		t.Run("true/"+p, func(t *testing.T) {
			if !parser.MatchPath(p) {
				t.Errorf("MatchPath(%q) = false, want true", p)
			}
		})
	}
}

func TestExtractOperationStillWorks(t *testing.T) {
	parser := &iotRESTParser{}

	tests := []struct {
		method string
		path   string
		want   string
	}{
		{http.MethodGet, "/things", "ListThings"},
		{http.MethodGet, "/things/myThing", "DescribeThing"},
		{http.MethodGet, "/certificates", "ListCertificates"},
		{http.MethodGet, "/policies", "ListPolicies"},
		{http.MethodGet, "/rules", "ListTopicRules"},
		{http.MethodGet, "/jobs", "ListJobs"},
		{http.MethodGet, "/endpoint", "DescribeEndpoint"},
		{http.MethodGet, "/role-aliases", "ListRoleAliases"},
		{http.MethodGet, "/thing-types", "ListThingTypes"},
		{http.MethodGet, "/thing-groups", "ListThingGroups"},
		{http.MethodGet, "/billing-groups", "ListBillingGroups"},
		{http.MethodGet, "/tags", "ListTagsForResource"},
		{http.MethodPost, "/tags", "TagResource"},
		{http.MethodPost, "/untag", "UntagResource"},
		{http.MethodGet, "/indexing", "GetIndexingConfiguration"},
		{http.MethodGet, "/active-violations", "ListActiveViolations"},
		{http.MethodGet, "/violation-events", "ListViolationEvents"},
		{http.MethodGet, "/security-profiles", "ListSecurityProfiles"},
		{http.MethodPost, "/security-profile-behaviors/validate", "ValidateSecurityProfileBehaviors"},
		{http.MethodPost, "/messages", "BatchPutMessage"},
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

func TestPathHasRoutePrefix(t *testing.T) {
	tests := []struct {
		path   string
		prefix string
		want   bool
	}{
		{"/things", "/things", true},
		{"/things/abc", "/things", true},
		{"/thingsabc", "/things", false},
		{"/things-abc", "/things", false},
		{"/a/b", "/a", true},
		{"/abc", "/a", false},
		{"/", "/", true},
	}

	for _, tt := range tests {
		t.Run(tt.path+"~"+tt.prefix, func(t *testing.T) {
			got := pathHasRoutePrefix(tt.path, tt.prefix)
			if got != tt.want {
				t.Errorf("pathHasRoutePrefix(%q, %q) = %v, want %v", tt.path, tt.prefix, got, tt.want)
			}
		})
	}
}
