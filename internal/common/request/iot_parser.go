package request

import (
	"net/http"
	"strings"
)

type iotRESTParser struct{}

func (p *iotRESTParser) MatchPath(path string) bool {
	return strings.HasPrefix(path, "/things") ||
		strings.HasPrefix(path, "/thing-groups") ||
		strings.HasPrefix(path, "/thing-types") ||
		strings.HasPrefix(path, "/billing-groups") ||
		strings.HasPrefix(path, "/certificates") ||
		strings.HasPrefix(path, "/certificate") ||
		strings.HasPrefix(path, "/keys-and-certificate") ||
		strings.HasPrefix(path, "/policies") ||
		strings.HasPrefix(path, "/policy-principals") ||
		strings.HasPrefix(path, "/principal-policies") ||
		strings.HasPrefix(path, "/principals") ||
		strings.HasPrefix(path, "/target-policies") ||
		strings.HasPrefix(path, "/attached-policies") ||
		strings.HasPrefix(path, "/rules") ||
		strings.HasPrefix(path, "/jobs") ||
		strings.HasPrefix(path, "/endpoint") ||
		strings.HasPrefix(path, "/role-aliases") ||
		strings.HasPrefix(path, "/tags") ||
		strings.HasPrefix(path, "/untag") ||
		strings.HasPrefix(path, "/authorizers") ||
		strings.HasPrefix(path, "/authorizer/") ||
		strings.HasPrefix(path, "/provisioning-templates") ||
		strings.HasPrefix(path, "/provisioning-template") ||
		strings.HasPrefix(path, "/domainConfigurations") ||
		strings.HasPrefix(path, "/domainConfiguration") ||
		strings.HasPrefix(path, "/indexing") ||
		strings.HasPrefix(path, "/active-violations") ||
		strings.HasPrefix(path, "/violation-events") ||
		strings.HasPrefix(path, "/behavior-model-training") ||
		strings.HasPrefix(path, "/security-profiles") ||
		strings.HasPrefix(path, "/detector-models") ||
		strings.HasPrefix(path, "/detector-model") ||
		strings.HasPrefix(path, "/inputs") ||
		strings.HasPrefix(path, "/input/") ||
		strings.HasPrefix(path, "/messages") ||
		strings.HasPrefix(path, "/destinations") ||
		strings.HasPrefix(path, "/effective-policies")
}

func (p *iotRESTParser) ExtractOperation(r *http.Request) string {
	path := r.URL.Path
	method := r.Method
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")

	switch {
	case path == "/things" && method == http.MethodGet:
		return "ListThings"
	case strings.HasPrefix(path, "/things/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateThing"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateThing"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteThing"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeThing"
		case len(parts) >= 3 && parts[2] == "principals" && method == http.MethodGet:
			return "ListThingPrincipals"
		case len(parts) >= 3 && parts[2] == "principals" && (method == http.MethodPut || method == http.MethodPost):
			return "AttachThingPrincipal"
		case len(parts) >= 3 && parts[2] == "principals" && method == http.MethodDelete:
			return "DetachThingPrincipal"
		case len(parts) >= 3 && parts[2] == "thing-groups" && method == http.MethodGet:
			return "ListThingGroupsForThing"
		case len(parts) >= 3 && parts[2] == "jobs" && len(parts) >= 4 && method == http.MethodGet:
			return "DescribeJob"
		case len(parts) >= 3 && parts[2] == "jobs" && len(parts) >= 4 && method == http.MethodPost:
			return "UpdateJob"
		}

	case path == "/keys-and-certificate" && method == http.MethodPost:
		return "CreateKeysAndCertificate"
	case path == "/certificate/register" && method == http.MethodPost:
		return "RegisterCertificate"
	case strings.HasPrefix(path, "/certificates/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeCertificate"
		case len(parts) == 2 && (method == http.MethodPatch || method == http.MethodPut):
			return "UpdateCertificate"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteCertificate"
		}
	case path == "/certificates" && method == http.MethodGet:
		return "ListCertificates"
	case strings.HasPrefix(path, "/certificates/") && len(parts) >= 3 && parts[2] == "csr" && method == http.MethodPost:
		return "CreateCertificateFromCsr"

	case path == "/policies" && method == http.MethodGet:
		return "ListPolicies"
	case strings.HasPrefix(path, "/policies/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreatePolicy"
		case len(parts) == 2 && method == http.MethodGet:
			return "GetPolicy"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeletePolicy"
		case len(parts) >= 3 && parts[2] == "version" && method == http.MethodGet:
			if len(parts) >= 4 {
				return "GetPolicyVersion"
			}
			return "ListPolicyVersions"
		}

	case strings.HasPrefix(path, "/target-policies/") && len(parts) >= 2:
		switch {
		case method == http.MethodGet:
			return "ListPolicyPrincipals"
		case method == http.MethodPut:
			return "AttachPolicy"
		case method == http.MethodPost:
			return "DetachPolicy"
		}
	case strings.HasPrefix(path, "/attached-policies/") && method == http.MethodPost:
		return "ListAttachedPolicies"
	case path == "/principal-policies" && method == http.MethodGet:
		return "ListPrincipalPolicies"
	case strings.HasPrefix(path, "/principal-policies/") && method == http.MethodGet:
		return "ListPrincipalPolicies"
	case path == "/principals/things" && method == http.MethodGet:
		return "ListPrincipalThings"

	case path == "/rules" && method == http.MethodGet:
		return "ListTopicRules"
	case strings.HasPrefix(path, "/rules/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateTopicRule"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeTopicRule"
		case len(parts) == 2 && method == http.MethodPatch:
			return "ReplaceTopicRule"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteTopicRule"
		case len(parts) == 3 && parts[2] == "enable" && method == http.MethodGet:
			return "EnableTopicRule"
		case len(parts) == 3 && parts[2] == "disable" && method == http.MethodGet:
			return "DisableTopicRule"
		}

	case path == "/jobs" && method == http.MethodPost:
		return "CreateJob"
	case path == "/jobs" && method == http.MethodGet:
		return "ListJobs"
	case strings.HasPrefix(path, "/jobs/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeJob"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteJob"
		case len(parts) == 3 && parts[2] == "cancel" && method == http.MethodPost:
			return "CancelJob"
		case len(parts) == 3 && parts[2] == "job-document" && method == http.MethodGet:
			return "GetJobDocument"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateJob"
		}

	case path == "/endpoint" && method == http.MethodGet:
		return "DescribeEndpoint"

	case path == "/role-aliases" && method == http.MethodGet:
		return "ListRoleAliases"
	case path == "/role-aliases" && method == http.MethodPost:
		return "CreateRoleAlias"
	case strings.HasPrefix(path, "/role-aliases/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeRoleAlias"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateRoleAlias"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteRoleAlias"
		}

	case path == "/thing-types" && method == http.MethodGet:
		return "ListThingTypes"
	case strings.HasPrefix(path, "/thing-types/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateThingType"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeThingType"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateThingType"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteThingType"
		case len(parts) == 3 && parts[2] == "deprecate" && method == http.MethodPost:
			return "DeprecateThingType"
		}

	case path == "/thing-groups" && method == http.MethodGet:
		return "ListThingGroups"
	case path == "/thing-groups/addThingToThingGroup" && method == http.MethodPut:
		return "AddThingToThingGroup"
	case path == "/thing-groups/removeThingFromThingGroup" && method == http.MethodPut:
		return "RemoveThingFromThingGroup"
	case strings.HasPrefix(path, "/thing-groups/") && len(parts) >= 2 && parts[1] != "addThingToThingGroup" && parts[1] != "removeThingFromThingGroup":
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateThingGroup"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeThingGroup"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateThingGroup"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteThingGroup"
		case len(parts) == 3 && parts[2] == "things" && method == http.MethodGet:
			return "ListThingsInThingGroup"
		}

	case path == "/billing-groups" && method == http.MethodGet:
		return "ListBillingGroups"
	case strings.HasPrefix(path, "/billing-groups/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateBillingGroup"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeBillingGroup"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateBillingGroup"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteBillingGroup"
		}

	case path == "/authorizers" && method == http.MethodPost:
		return "CreateAuthorizer"
	case path == "/authorizers" && method == http.MethodGet:
		return "ListAuthorizers"
	case strings.HasPrefix(path, "/authorizer/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeAuthorizer"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateAuthorizer"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteAuthorizer"
		}

	case path == "/tags" && method == http.MethodGet:
		return "ListTagsForResource"
	case path == "/tags" && method == http.MethodPost:
		return "TagResource"
	case path == "/untag" && method == http.MethodPost:
		return "UntagResource"

	case path == "/domainConfigurations" && method == http.MethodGet:
		return "ListDomainConfigurations"
	case path == "/domainConfigurations" && method == http.MethodPost:
		return "CreateDomainConfiguration"
	case strings.HasPrefix(path, "/domainConfiguration/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeDomainConfiguration"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateDomainConfiguration"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteDomainConfiguration"
		}

	case path == "/indexing" && method == http.MethodGet:
		return "GetIndexingConfiguration"
	case path == "/indexing" && method == http.MethodPatch:
		return "UpdateIndexingConfiguration"

	case path == "/active-violations" && method == http.MethodGet:
		return "ListActiveViolations"
	case path == "/violation-events" && method == http.MethodGet:
		return "ListViolationEvents"
	case path == "/behavior-model-training/summaries" && method == http.MethodGet:
		return "GetBehaviorModelTrainingSummaries"

	case strings.HasPrefix(path, "/provisioning-templates") && method == http.MethodGet:
		return "ListProvisioningTemplates"
	case strings.HasPrefix(path, "/provisioning-templates/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateProvisioningTemplate"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeProvisioningTemplate"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateProvisioningTemplate"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteProvisioningTemplate"
		case len(parts) == 3 && parts[2] == "versions" && method == http.MethodGet:
			return "ListProvisioningTemplateVersions"
		}

	case path == "/messages" && method == http.MethodPost:
		return "BatchPutMessage"

	case path == "/effective-policies" && method == http.MethodPost:
		return "GetEffectivePolicies"
	}

	return ""
}

func (p *iotRESTParser) ExtractPathParams(r *http.Request, params map[string]interface{}) {
	path := r.URL.Path
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")
	method := r.Method

	switch {
	case strings.HasPrefix(path, "/things/") && len(parts) >= 2:
		params["thingName"] = parts[1]
		if method == http.MethodPost {
			params["_operation"] = "CreateThing"
		} else if method == http.MethodPatch {
			params["_operation"] = "UpdateThing"
		} else if method == http.MethodGet {
			params["_operation"] = "DescribeThing"
		} else if method == http.MethodDelete {
			params["_operation"] = "DeleteThing"
		}
	case strings.HasPrefix(path, "/certificates/") && len(parts) >= 2:
		params["certificateId"] = parts[1]
	case strings.HasPrefix(path, "/policies/") && len(parts) >= 2:
		params["policyName"] = parts[1]
		if method == http.MethodPost {
			params["_operation"] = "CreatePolicy"
		}
	case strings.HasPrefix(path, "/rules/") && len(parts) >= 2:
		params["ruleName"] = parts[1]
	case strings.HasPrefix(path, "/jobs/") && len(parts) >= 2:
		params["jobId"] = parts[1]
	case strings.HasPrefix(path, "/role-aliases/") && len(parts) >= 2:
		params["roleAlias"] = parts[1]
	case strings.HasPrefix(path, "/thing-types/") && len(parts) >= 2:
		params["thingTypeName"] = parts[1]
	case strings.HasPrefix(path, "/thing-groups/") && len(parts) >= 2:
		if parts[1] != "addThingToThingGroup" && parts[1] != "removeThingFromThingGroup" {
			params["thingGroupName"] = parts[1]
		}
	case strings.HasPrefix(path, "/billing-groups/") && len(parts) >= 2:
		params["billingGroupName"] = parts[1]
	case strings.HasPrefix(path, "/authorizer/") && len(parts) >= 2:
		params["authorizerName"] = parts[1]
	case strings.HasPrefix(path, "/provisioning-templates/") && len(parts) >= 2:
		params["templateName"] = parts[1]
	case strings.HasPrefix(path, "/domainConfiguration/") && len(parts) >= 2:
		params["domainConfigurationName"] = parts[1]
	case strings.HasPrefix(path, "/policy-principals/") && len(parts) >= 2:
		params["policyName"] = parts[1]
	case strings.HasPrefix(path, "/target-policies/") && len(parts) >= 2:
		params["policyName"] = parts[1]
	case strings.HasPrefix(path, "/attached-policies/"):
		params["target"] = strings.TrimPrefix(path, "/attached-policies/")
	case strings.HasPrefix(path, "/principal-policies/") && len(parts) >= 2:
		params["policyName"] = parts[1]
	}

	// Handle tag operations with resourceArn from body
	if path == "/tags" || path == "/untag" {
		params["_operation"] = "TagResource"
		if path == "/untag" {
			params["_operation"] = "UntagResource"
		}
	}

	// Extract query params for list operations
	if path == "/things" && method == http.MethodGet {
		if tn := r.URL.Query().Get("thingTypeName"); tn != "" {
			params["thingTypeName"] = tn
			params["_operation"] = "ListThingsForThingType"
		}
	}
	if path == "/endpoint" && method == http.MethodGet {
		if et := r.URL.Query().Get("endpointType"); et != "" {
			params["endpointType"] = et
		}
	}
}
