package request

import (
	"net/http"
	"strings"
)

type iotRESTParser struct{}

func pathHasRoutePrefix(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func (p *iotRESTParser) MatchPath(path string) bool {
	// Covers all IoT REST-JSON paths from the Smithy model.
	// MatchPath gates ExtractPathParams, so it must return true for any path
	// that carries URI-bound parameters. Operation name extraction is handled
	// by ExtractOperation independently.
	return pathHasRoutePrefix(path, "/things") ||
		pathHasRoutePrefix(path, "/topics") ||
		pathHasRoutePrefix(path, "/api/things/shadow") ||
		pathHasRoutePrefix(path, "/thing-groups") ||
		pathHasRoutePrefix(path, "/thing-types") ||
		pathHasRoutePrefix(path, "/thing-registration-tasks") ||
		pathHasRoutePrefix(path, "/billing-groups") ||
		pathHasRoutePrefix(path, "/dynamic-thing-groups") ||
		pathHasRoutePrefix(path, "/certificates") ||
		pathHasRoutePrefix(path, "/certificates-by-ca") ||
		pathHasRoutePrefix(path, "/certificates-out-going") ||
		path == "/certificate/register" ||
		path == "/certificate/register-no-ca" ||
		path == "/keys-and-certificate" ||
		pathHasRoutePrefix(path, "/cacertificate") ||
		pathHasRoutePrefix(path, "/cacertificates") ||
		pathHasRoutePrefix(path, "/accept-certificate-transfer") ||
		pathHasRoutePrefix(path, "/cancel-certificate-transfer") ||
		pathHasRoutePrefix(path, "/reject-certificate-transfer") ||
		pathHasRoutePrefix(path, "/transfer-certificate") ||
		path == "/registrationcode" ||
		pathHasRoutePrefix(path, "/policies") ||
		pathHasRoutePrefix(path, "/policy-principals") ||
		pathHasRoutePrefix(path, "/principal-policies") ||
		path == "/principals/things" ||
		path == "/principals/things-v2" ||
		pathHasRoutePrefix(path, "/target-policies") ||
		pathHasRoutePrefix(path, "/attached-policies") ||
		pathHasRoutePrefix(path, "/policy-targets") ||
		pathHasRoutePrefix(path, "/rules") ||
		pathHasRoutePrefix(path, "/jobs") ||
		pathHasRoutePrefix(path, "/job-templates") ||
		pathHasRoutePrefix(path, "/managed-job-templates") ||
		path == "/endpoint" ||
		pathHasRoutePrefix(path, "/role-aliases") ||
		path == "/tags" ||
		path == "/untag" ||
		pathHasRoutePrefix(path, "/authorizers") ||
		pathHasRoutePrefix(path, "/authorizer") ||
		path == "/default-authorizer" ||
		pathHasRoutePrefix(path, "/provisioning-templates") ||
		pathHasRoutePrefix(path, "/provisioning-template") ||
		pathHasRoutePrefix(path, "/domainConfigurations") ||
		pathHasRoutePrefix(path, "/domainConfiguration") ||
		pathHasRoutePrefix(path, "/indexing") ||
		pathHasRoutePrefix(path, "/indices") ||
		pathHasRoutePrefix(path, "/active-violations") ||
		pathHasRoutePrefix(path, "/violation-events") ||
		pathHasRoutePrefix(path, "/violations") ||
		path == "/behavior-model-training/summaries" ||
		pathHasRoutePrefix(path, "/security-profiles") ||
		path == "/security-profile-behaviors/validate" ||
		pathHasRoutePrefix(path, "/security-profiles-for-target") ||
		pathHasRoutePrefix(path, "/audit") ||
		pathHasRoutePrefix(path, "/detect") ||
		pathHasRoutePrefix(path, "/mitigationactions") ||
		pathHasRoutePrefix(path, "/custom-metric") ||
		pathHasRoutePrefix(path, "/custom-metrics") ||
		pathHasRoutePrefix(path, "/dimensions") ||
		pathHasRoutePrefix(path, "/fleet-metric") ||
		pathHasRoutePrefix(path, "/fleet-metrics") ||
		pathHasRoutePrefix(path, "/metric-values") ||
		pathHasRoutePrefix(path, "/streams") ||
		pathHasRoutePrefix(path, "/otaUpdates") ||
		pathHasRoutePrefix(path, "/loggingOptions") ||
		pathHasRoutePrefix(path, "/v2LoggingLevel") ||
		pathHasRoutePrefix(path, "/v2LoggingOptions") ||
		pathHasRoutePrefix(path, "/encryption-configuration") ||
		pathHasRoutePrefix(path, "/event-configurations") ||
		path == "/test-authorization" ||
		pathHasRoutePrefix(path, "/confirmdestination") ||
		pathHasRoutePrefix(path, "/destinations") ||
		pathHasRoutePrefix(path, "/packages") ||
		path == "/package-configuration" ||
		pathHasRoutePrefix(path, "/commands") ||
		pathHasRoutePrefix(path, "/command-executions") ||
		pathHasRoutePrefix(path, "/certificate-providers") ||
		path == "/effective-policies" ||
		path == "/messages"
}

func (p *iotRESTParser) ExtractOperation(r *http.Request) string {
	// Normalise trailing slashes so that "/authorizers/" and "/authorizers"
	// are treated identically. The AWS CLI v2 emits a trailing slash for
	// some collection resources (e.g. ListAuthorizers sends
	// "/authorizers/"), which would otherwise miss exact-match cases.
	path := strings.TrimRight(r.URL.Path, "/")
	if path == "" {
		path = "/"
	}
	method := r.Method
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")

	switch {
	case strings.HasPrefix(path, "/topics/") && method == http.MethodPost:
		return "Publish"
	case strings.HasPrefix(path, "/api/things/shadow/ListNamedShadowsForThing/") && method == http.MethodGet:
		return "ListNamedShadowsForThing"
	case path == "/things" && method == http.MethodGet:
		return "ListThings"
	case strings.HasPrefix(path, "/things/") && len(parts) >= 2:
		switch {
		case len(parts) == 3 && parts[2] == "shadow" && method == http.MethodGet:
			return "GetThingShadow"
		case len(parts) == 3 && parts[2] == "shadow" && method == http.MethodPost:
			return "UpdateThingShadow"
		case len(parts) == 3 && parts[2] == "shadow" && method == http.MethodDelete:
			return "DeleteThingShadow"
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
		case len(parts) == 3 && parts[2] == "principals-v2" && method == http.MethodGet:
			return "ListThingPrincipalsV2"
		case len(parts) == 3 && parts[2] == "connectivity-data" && method == http.MethodPost:
			return "GetThingConnectivityData"
		case len(parts) >= 3 && parts[2] == "thing-groups" && method == http.MethodGet:
			return "ListThingGroupsForThing"
		case len(parts) == 3 && parts[2] == "jobs" && method == http.MethodGet:
			return "ListJobExecutionsForThing"
		case len(parts) == 4 && parts[2] == "jobs" && method == http.MethodGet:
			return "DescribeJobExecution"
		case len(parts) == 5 && parts[2] == "jobs" && parts[4] == "cancel" && method == http.MethodPut:
			return "CancelJobExecution"
		case len(parts) == 6 && parts[2] == "jobs" && parts[4] == "executionNumber" && method == http.MethodDelete:
			return "DeleteJobExecution"
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
	case path == "/certificates" && method == http.MethodPost:
		return "CreateCertificateFromCsr"
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
		case len(parts) == 3 && parts[2] == "version" && method == http.MethodPost:
			return "CreatePolicyVersion"
		case len(parts) >= 4 && parts[2] == "version" && method == http.MethodPatch:
			return "SetDefaultPolicyVersion"
		case len(parts) >= 4 && parts[2] == "version" && method == http.MethodDelete:
			return "DeletePolicyVersion"
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
		case len(parts) == 3 && parts[2] == "enable" && (method == http.MethodGet || method == http.MethodPost):
			return "EnableTopicRule"
		case len(parts) == 3 && parts[2] == "disable" && (method == http.MethodGet || method == http.MethodPost):
			return "DisableTopicRule"
		}

	case path == "/jobs" && method == http.MethodPost:
		return "CreateJob"
	case path == "/jobs" && method == http.MethodGet:
		return "ListJobs"
	case strings.HasPrefix(path, "/jobs/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPut:
			return "CreateJob"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeJob"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteJob"
		case len(parts) == 3 && parts[2] == "cancel" && (method == http.MethodPost || method == http.MethodPut):
			return "CancelJob"
		case len(parts) == 3 && parts[2] == "job-document" && method == http.MethodGet:
			return "GetJobDocument"
		case len(parts) == 3 && parts[2] == "targets" && method == http.MethodPost:
			return "AssociateTargetsWithJob"
		case len(parts) == 3 && parts[2] == "things" && method == http.MethodGet:
			return "ListJobExecutionsForJob"
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
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateRoleAlias"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeRoleAlias"
		case len(parts) == 2 && (method == http.MethodPatch || method == http.MethodPut):
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
	case strings.HasPrefix(path, "/billing-groups/") && len(parts) >= 2 && parts[1] != "addThingToBillingGroup" && parts[1] != "removeThingFromBillingGroup":
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateBillingGroup"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeBillingGroup"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateBillingGroup"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteBillingGroup"
		case len(parts) == 3 && parts[2] == "things" && method == http.MethodGet:
			return "ListThingsInBillingGroup"
		}

	case path == "/authorizers" && method == http.MethodPost:
		return "CreateAuthorizer"
	case path == "/authorizers" && method == http.MethodGet:
		return "ListAuthorizers"
	case strings.HasPrefix(path, "/authorizer/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateAuthorizer"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeAuthorizer"
		case len(parts) == 2 && (method == http.MethodPatch || method == http.MethodPut):
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
	case strings.HasPrefix(path, "/domainConfigurations/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateDomainConfiguration"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeDomainConfiguration"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateDomainConfiguration"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteDomainConfiguration"
		}

	case (path == "/indexing/config" || path == "/indexing") && method == http.MethodGet:
		return "GetIndexingConfiguration"
	case (path == "/indexing/config" || path == "/indexing") && (method == http.MethodPatch || method == http.MethodPost):
		return "UpdateIndexingConfiguration"

	case path == "/active-violations" && method == http.MethodGet:
		return "ListActiveViolations"
	case path == "/violation-events" && method == http.MethodGet:
		return "ListViolationEvents"
	case path == "/behavior-model-training/summaries" && method == http.MethodGet:
		return "GetBehaviorModelTrainingSummaries"

	case path == "/security-profile-behaviors/validate" && method == http.MethodPost:
		return "ValidateSecurityProfileBehaviors"
	case path == "/security-profiles" && method == http.MethodGet:
		return "ListSecurityProfiles"
	case strings.HasPrefix(path, "/security-profiles/") && len(parts) >= 2:
		switch {
		case len(parts) == 2 && method == http.MethodPost:
			return "CreateSecurityProfile"
		case len(parts) == 2 && method == http.MethodGet:
			return "DescribeSecurityProfile"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdateSecurityProfile"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeleteSecurityProfile"
		case len(parts) == 3 && parts[2] == "targets" && method == http.MethodPut:
			return "AttachSecurityProfile"
		case len(parts) == 3 && parts[2] == "targets" && method == http.MethodDelete:
			return "DetachSecurityProfile"
		case len(parts) == 3 && parts[2] == "targets" && method == http.MethodGet:
			return "ListTargetsForSecurityProfile"
		}

	case path == "/provisioning-templates" && method == http.MethodGet:
		return "ListProvisioningTemplates"
	case path == "/provisioning-templates" && method == http.MethodPost:
		return "CreateProvisioningTemplate"
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
		case len(parts) == 3 && parts[2] == "versions" && method == http.MethodPost:
			return "CreateProvisioningTemplateVersion"
		case len(parts) == 4 && parts[2] == "versions" && method == http.MethodGet:
			return "DescribeProvisioningTemplateVersion"
		case len(parts) == 4 && parts[2] == "versions" && method == http.MethodDelete:
			return "DeleteProvisioningTemplateVersion"
		case len(parts) == 3 && parts[2] == "provisioning-claim" && method == http.MethodPost:
			return "CreateProvisioningClaim"
		}

	case path == "/messages" && method == http.MethodPost:
		return "BatchPutMessage"

	case path == "/effective-policies" && method == http.MethodPost:
		return "GetEffectivePolicies"
	// Auto-generated cases for missing operations.
	// Derived from Smithy model @http traits.

	case strings.HasPrefix(path, "/accept-certificate-transfer/") && method == http.MethodPatch:
		return "AcceptCertificateTransfer"
	case path == "/audit/configuration" && method == http.MethodDelete:
		return "DeleteAccountAuditConfiguration"
	case path == "/audit/configuration" && method == http.MethodGet:
		return "DescribeAccountAuditConfiguration"
	case path == "/audit/configuration" && method == http.MethodPatch:
		return "UpdateAccountAuditConfiguration"
	case path == "/audit/findings" && method == http.MethodPost:
		return "ListAuditFindings"
	case strings.HasPrefix(path, "/audit/findings/") && method == http.MethodGet:
		return "DescribeAuditFinding"
	case path == "/audit/mitigationactions/executions" && method == http.MethodGet:
		return "ListAuditMitigationActionsExecutions"
	case path == "/audit/mitigationactions/tasks" && method == http.MethodGet:
		return "ListAuditMitigationActionsTasks"
	case strings.HasPrefix(path, "/audit/mitigationactions/tasks/") && method == http.MethodGet:
		return "DescribeAuditMitigationActionsTask"
	case strings.HasPrefix(path, "/audit/mitigationactions/tasks/") && method == http.MethodPost:
		return "StartAuditMitigationActionsTask"
	case strings.HasPrefix(path, "/audit/mitigationactions/tasks/") && method == http.MethodPut:
		return "CancelAuditMitigationActionsTask"
	case path == "/audit/relatedResources" && method == http.MethodGet:
		return "ListRelatedResourcesForAuditFinding"
	case path == "/audit/scheduledaudits" && method == http.MethodGet:
		return "ListScheduledAudits"
	case strings.HasPrefix(path, "/audit/scheduledaudits/") && method == http.MethodDelete:
		return "DeleteScheduledAudit"
	case strings.HasPrefix(path, "/audit/scheduledaudits/") && method == http.MethodGet:
		return "DescribeScheduledAudit"
	case strings.HasPrefix(path, "/audit/scheduledaudits/") && method == http.MethodPatch:
		return "UpdateScheduledAudit"
	case strings.HasPrefix(path, "/audit/scheduledaudits/") && method == http.MethodPost:
		return "CreateScheduledAudit"
	case path == "/audit/suppressions/create" && method == http.MethodPost:
		return "CreateAuditSuppression"
	case path == "/audit/suppressions/delete" && method == http.MethodPost:
		return "DeleteAuditSuppression"
	case path == "/audit/suppressions/describe" && method == http.MethodPost:
		return "DescribeAuditSuppression"
	case path == "/audit/suppressions/list" && method == http.MethodPost:
		return "ListAuditSuppressions"
	case path == "/audit/suppressions/update" && method == http.MethodPatch:
		return "UpdateAuditSuppression"
	case path == "/audit/tasks" && method == http.MethodGet:
		return "ListAuditTasks"
	case path == "/audit/tasks" && method == http.MethodPost:
		return "StartOnDemandAuditTask"
	case strings.HasPrefix(path, "/audit/tasks/") && method == http.MethodGet:
		return "DescribeAuditTask"
	case strings.HasPrefix(path, "/audit/tasks/") && method == http.MethodPut:
		return "CancelAuditTask"
	case strings.HasPrefix(path, "/authorizer/") && method == http.MethodPost:
		return "TestInvokeAuthorizer"
	case path == "/billing-groups/addThingToBillingGroup" && method == http.MethodPut:
		return "AddThingToBillingGroup"
	case path == "/billing-groups/removeThingFromBillingGroup" && method == http.MethodPut:
		return "RemoveThingFromBillingGroup"
	case strings.HasPrefix(path, "/billing-groups/") && method == http.MethodGet:
		return "ListThingsInBillingGroup"
	case path == "/cacertificate" && method == http.MethodPost:
		return "RegisterCACertificate"
	case strings.HasPrefix(path, "/cacertificate/") && method == http.MethodDelete:
		return "DeleteCACertificate"
	case strings.HasPrefix(path, "/cacertificate/") && method == http.MethodGet:
		return "DescribeCACertificate"
	case strings.HasPrefix(path, "/cacertificate/") && method == http.MethodPut:
		return "UpdateCACertificate"
	case path == "/cacertificates" && method == http.MethodGet:
		return "ListCACertificates"
	case strings.HasPrefix(path, "/cancel-certificate-transfer/") && method == http.MethodPatch:
		return "CancelCertificateTransfer"
	case path == "/certificate-providers" && method == http.MethodGet:
		return "ListCertificateProviders"
	case strings.HasPrefix(path, "/certificate-providers/") && method == http.MethodDelete:
		return "DeleteCertificateProvider"
	case strings.HasPrefix(path, "/certificate-providers/") && method == http.MethodGet:
		return "DescribeCertificateProvider"
	case strings.HasPrefix(path, "/certificate-providers/") && method == http.MethodPost:
		return "CreateCertificateProvider"
	case strings.HasPrefix(path, "/certificate-providers/") && method == http.MethodPut:
		return "UpdateCertificateProvider"
	case path == "/certificate/register-no-ca" && method == http.MethodPost:
		return "RegisterCertificateWithoutCA"
	case strings.HasPrefix(path, "/certificates-by-ca/") && method == http.MethodGet:
		return "ListCertificatesByCA"
	case path == "/certificates-out-going" && method == http.MethodGet:
		return "ListOutgoingCertificates"
	case path == "/command-executions" && method == http.MethodPost:
		return "ListCommandExecutions"
	case strings.HasPrefix(path, "/command-executions/") && method == http.MethodDelete:
		return "DeleteCommandExecution"
	case strings.HasPrefix(path, "/command-executions/") && method == http.MethodGet:
		return "GetCommandExecution"
	case path == "/commands" && method == http.MethodGet:
		return "ListCommands"
	case strings.HasPrefix(path, "/commands/") && method == http.MethodDelete:
		return "DeleteCommand"
	case strings.HasPrefix(path, "/commands/") && method == http.MethodGet:
		return "GetCommand"
	case strings.HasPrefix(path, "/commands/") && method == http.MethodPatch:
		return "UpdateCommand"
	case strings.HasPrefix(path, "/commands/") && method == http.MethodPut:
		return "CreateCommand"
	case strings.HasPrefix(path, "/confirmdestination/") && method == http.MethodGet:
		return "ConfirmTopicRuleDestination"
	case strings.HasPrefix(path, "/custom-metric/") && method == http.MethodDelete:
		return "DeleteCustomMetric"
	case strings.HasPrefix(path, "/custom-metric/") && method == http.MethodGet:
		return "DescribeCustomMetric"
	case strings.HasPrefix(path, "/custom-metric/") && method == http.MethodPatch:
		return "UpdateCustomMetric"
	case strings.HasPrefix(path, "/custom-metric/") && method == http.MethodPost:
		return "CreateCustomMetric"
	case path == "/custom-metrics" && method == http.MethodGet:
		return "ListCustomMetrics"
	case path == "/default-authorizer" && method == http.MethodDelete:
		return "ClearDefaultAuthorizer"
	case path == "/default-authorizer" && method == http.MethodGet:
		return "DescribeDefaultAuthorizer"
	case path == "/default-authorizer" && method == http.MethodPost:
		return "SetDefaultAuthorizer"
	case path == "/destinations" && method == http.MethodGet:
		return "ListTopicRuleDestinations"
	case path == "/destinations" && method == http.MethodPatch:
		return "UpdateTopicRuleDestination"
	case path == "/destinations" && method == http.MethodPost:
		return "CreateTopicRuleDestination"
	case strings.HasPrefix(path, "/destinations/") && method == http.MethodDelete:
		return "DeleteTopicRuleDestination"
	case strings.HasPrefix(path, "/destinations/") && method == http.MethodGet:
		return "GetTopicRuleDestination"
	case path == "/detect/mitigationactions/executions" && method == http.MethodGet:
		return "ListDetectMitigationActionsExecutions"
	case path == "/detect/mitigationactions/tasks" && method == http.MethodGet:
		return "ListDetectMitigationActionsTasks"
	case strings.HasPrefix(path, "/detect/mitigationactions/tasks/") && method == http.MethodGet:
		return "DescribeDetectMitigationActionsTask"
	case strings.HasPrefix(path, "/detect/mitigationactions/tasks/") && method == http.MethodPut:
		// /detect/mitigationactions/tasks/{taskId} -> Start;
		// /detect/mitigationactions/tasks/{taskId}/cancel -> Cancel.
		if strings.HasSuffix(path, "/cancel") {
			return "CancelDetectMitigationActionsTask"
		}
		return "StartDetectMitigationActionsTask"
	case path == "/dimensions" && method == http.MethodGet:
		return "ListDimensions"
	case strings.HasPrefix(path, "/dimensions/") && method == http.MethodDelete:
		return "DeleteDimension"
	case strings.HasPrefix(path, "/dimensions/") && method == http.MethodGet:
		return "DescribeDimension"
	case strings.HasPrefix(path, "/dimensions/") && method == http.MethodPatch:
		return "UpdateDimension"
	case strings.HasPrefix(path, "/dimensions/") && method == http.MethodPost:
		return "CreateDimension"
	case strings.HasPrefix(path, "/dynamic-thing-groups/") && method == http.MethodDelete:
		return "DeleteDynamicThingGroup"
	case strings.HasPrefix(path, "/dynamic-thing-groups/") && method == http.MethodPatch:
		return "UpdateDynamicThingGroup"
	case strings.HasPrefix(path, "/dynamic-thing-groups/") && method == http.MethodPost:
		return "CreateDynamicThingGroup"
	case path == "/encryption-configuration" && method == http.MethodGet:
		return "DescribeEncryptionConfiguration"
	case path == "/encryption-configuration" && method == http.MethodPatch:
		return "UpdateEncryptionConfiguration"
	case path == "/event-configurations" && method == http.MethodGet:
		return "DescribeEventConfigurations"
	case path == "/event-configurations" && method == http.MethodPatch:
		return "UpdateEventConfigurations"
	case strings.HasPrefix(path, "/fleet-metric/") && method == http.MethodDelete:
		return "DeleteFleetMetric"
	case strings.HasPrefix(path, "/fleet-metric/") && method == http.MethodGet:
		return "DescribeFleetMetric"
	case strings.HasPrefix(path, "/fleet-metric/") && method == http.MethodPatch:
		return "UpdateFleetMetric"
	case strings.HasPrefix(path, "/fleet-metric/") && method == http.MethodPut:
		return "CreateFleetMetric"
	case path == "/fleet-metrics" && method == http.MethodGet:
		return "ListFleetMetrics"
	case path == "/indices" && method == http.MethodGet:
		return "ListIndices"
	case path == "/indices/buckets" && method == http.MethodPost:
		return "GetBucketsAggregation"
	case path == "/indices/cardinality" && method == http.MethodPost:
		return "GetCardinality"
	case path == "/indices/percentiles" && method == http.MethodPost:
		return "GetPercentiles"
	case path == "/indices/search" && method == http.MethodPost:
		return "SearchIndex"
	case path == "/indices/statistics" && method == http.MethodPost:
		return "GetStatistics"
	case strings.HasPrefix(path, "/indices/") && method == http.MethodGet:
		return "DescribeIndex"
	case path == "/job-templates" && method == http.MethodGet:
		return "ListJobTemplates"
	case strings.HasPrefix(path, "/job-templates/") && method == http.MethodDelete:
		return "DeleteJobTemplate"
	case strings.HasPrefix(path, "/job-templates/") && method == http.MethodGet:
		return "DescribeJobTemplate"
	case strings.HasPrefix(path, "/job-templates/") && method == http.MethodPut:
		return "CreateJobTemplate"
	case path == "/loggingOptions" && method == http.MethodGet:
		return "GetLoggingOptions"
	case path == "/loggingOptions" && method == http.MethodPost:
		return "SetLoggingOptions"
	case path == "/managed-job-templates" && method == http.MethodGet:
		return "ListManagedJobTemplates"
	case strings.HasPrefix(path, "/managed-job-templates/") && method == http.MethodGet:
		return "DescribeManagedJobTemplate"
	case path == "/metric-values" && method == http.MethodGet:
		return "ListMetricValues"
	case path == "/mitigationactions/actions" && method == http.MethodGet:
		return "ListMitigationActions"
	case strings.HasPrefix(path, "/mitigationactions/actions/") && method == http.MethodDelete:
		return "DeleteMitigationAction"
	case strings.HasPrefix(path, "/mitigationactions/actions/") && method == http.MethodGet:
		return "DescribeMitigationAction"
	case strings.HasPrefix(path, "/mitigationactions/actions/") && method == http.MethodPatch:
		return "UpdateMitigationAction"
	case strings.HasPrefix(path, "/mitigationactions/actions/") && method == http.MethodPost:
		return "CreateMitigationAction"
	case path == "/otaUpdates" && method == http.MethodGet:
		return "ListOTAUpdates"
	case strings.HasPrefix(path, "/otaUpdates/") && method == http.MethodDelete:
		return "DeleteOTAUpdate"
	case strings.HasPrefix(path, "/otaUpdates/") && method == http.MethodGet:
		return "GetOTAUpdate"
	case strings.HasPrefix(path, "/otaUpdates/") && method == http.MethodPost:
		return "CreateOTAUpdate"
	case path == "/package-configuration" && method == http.MethodGet:
		return "GetPackageConfiguration"
	case path == "/package-configuration" && method == http.MethodPatch:
		return "UpdatePackageConfiguration"
	case path == "/packages" && method == http.MethodGet:
		return "ListPackages"
	case strings.HasPrefix(path, "/packages/") && len(parts) >= 2:
		switch {
		case len(parts) == 4 && parts[2] == "versions" && parts[3] == "sbom" && method == http.MethodDelete:
			return "DisassociateSbomFromPackageVersion"
		case len(parts) == 4 && parts[2] == "versions" && parts[3] == "sbom" && method == http.MethodPut:
			return "AssociateSbomWithPackageVersion"
		case len(parts) == 5 && parts[2] == "versions" && parts[4] == "sbom-validation-results" && method == http.MethodGet:
			return "ListSbomValidationResults"
		case len(parts) == 4 && parts[2] == "versions" && method == http.MethodDelete:
			return "DeletePackageVersion"
		case len(parts) == 4 && parts[2] == "versions" && method == http.MethodGet:
			return "GetPackageVersion"
		case len(parts) == 4 && parts[2] == "versions" && method == http.MethodPatch:
			return "UpdatePackageVersion"
		case len(parts) == 4 && parts[2] == "versions" && method == http.MethodPut:
			return "CreatePackageVersion"
		case len(parts) == 3 && parts[2] == "versions" && method == http.MethodGet:
			return "ListPackageVersions"
		case len(parts) == 2 && method == http.MethodDelete:
			return "DeletePackage"
		case len(parts) == 2 && method == http.MethodGet:
			return "GetPackage"
		case len(parts) == 2 && method == http.MethodPatch:
			return "UpdatePackage"
		case len(parts) == 2 && method == http.MethodPut:
			return "CreatePackage"
		}
	case strings.HasPrefix(path, "/policies/") && method == http.MethodPost:
		return "CreatePolicyVersion"
	case strings.HasPrefix(path, "/policies/") && method == http.MethodDelete:
		return "DeletePolicyVersion"
	case strings.HasPrefix(path, "/policies/") && method == http.MethodPatch:
		return "SetDefaultPolicyVersion"
	case strings.HasPrefix(path, "/policy-targets/") && method == http.MethodPost:
		return "ListTargetsForPolicy"
	case strings.HasPrefix(path, "/principal-policies/") && method == http.MethodDelete:
		return "DetachPrincipalPolicy"
	case strings.HasPrefix(path, "/principal-policies/") && method == http.MethodPut:
		return "AttachPrincipalPolicy"
	case path == "/principals/things-v2" && method == http.MethodGet:
		return "ListPrincipalThingsV2"
	case strings.HasPrefix(path, "/provisioning-templates/") && method == http.MethodDelete:
		return "DeleteProvisioningTemplateVersion"
	case strings.HasPrefix(path, "/provisioning-templates/") && method == http.MethodGet:
		return "DescribeProvisioningTemplateVersion"
	case path == "/registrationcode" && method == http.MethodDelete:
		return "DeleteRegistrationCode"
	case path == "/registrationcode" && method == http.MethodGet:
		return "GetRegistrationCode"
	case strings.HasPrefix(path, "/reject-certificate-transfer/") && method == http.MethodPatch:
		return "RejectCertificateTransfer"
	case strings.HasPrefix(path, "/rules/") && method == http.MethodGet:
		return "GetTopicRule"
	case path == "/security-profiles-for-target" && method == http.MethodGet:
		return "ListSecurityProfilesForTarget"
	case strings.HasPrefix(path, "/security-profiles/") && method == http.MethodDelete:
		return "DetachSecurityProfile"
	case strings.HasPrefix(path, "/security-profiles/") && method == http.MethodGet:
		return "ListTargetsForSecurityProfile"
	case strings.HasPrefix(path, "/security-profiles/") && method == http.MethodPut:
		return "AttachSecurityProfile"
	case path == "/streams" && method == http.MethodGet:
		return "ListStreams"
	case strings.HasPrefix(path, "/streams/") && method == http.MethodDelete:
		return "DeleteStream"
	case strings.HasPrefix(path, "/streams/") && method == http.MethodGet:
		return "DescribeStream"
	case strings.HasPrefix(path, "/streams/") && method == http.MethodPost:
		return "CreateStream"
	case strings.HasPrefix(path, "/streams/") && method == http.MethodPut:
		return "UpdateStream"
	case path == "/test-authorization" && method == http.MethodPost:
		return "TestAuthorization"
	case path == "/thing-groups/updateThingGroupsForThing" && method == http.MethodPut:
		return "UpdateThingGroupsForThing"
	case path == "/thing-registration-tasks" && method == http.MethodGet:
		return "ListThingRegistrationTasks"
	case path == "/thing-registration-tasks" && method == http.MethodPost:
		return "StartThingRegistrationTask"
	case strings.HasPrefix(path, "/thing-registration-tasks/") && method == http.MethodGet:
		return "DescribeThingRegistrationTask"
	case strings.HasPrefix(path, "/thing-registration-tasks/") && method == http.MethodPut:
		return "StopThingRegistrationTask"
	case strings.HasPrefix(path, "/thing-registration-tasks/") && method == http.MethodGet:
		return "ListThingRegistrationTaskReports"
	case path == "/things" && method == http.MethodPost:
		return "RegisterThing"
	case strings.HasPrefix(path, "/things/") && method == http.MethodPost:
		return "GetThingConnectivityData"
	case strings.HasPrefix(path, "/things/") && method == http.MethodGet:
		return "ListJobExecutionsForThing"
	case strings.HasPrefix(path, "/things/") && method == http.MethodGet:
		return "DescribeJobExecution"
	case strings.HasPrefix(path, "/things/") && method == http.MethodPut:
		return "CancelJobExecution"
	case strings.HasPrefix(path, "/things/") && method == http.MethodDelete:
		return "DeleteJobExecution"
	case strings.HasPrefix(path, "/things/") && method == http.MethodGet:
		return "ListThingPrincipalsV2"
	case strings.HasPrefix(path, "/transfer-certificate/") && method == http.MethodPatch:
		return "TransferCertificate"
	case path == "/v2LoggingLevel" && method == http.MethodDelete:
		return "DeleteV2LoggingLevel"
	case path == "/v2LoggingLevel" && method == http.MethodGet:
		return "ListV2LoggingLevels"
	case path == "/v2LoggingLevel" && method == http.MethodPost:
		return "SetV2LoggingLevel"
	case path == "/v2LoggingOptions" && method == http.MethodGet:
		return "GetV2LoggingOptions"
	case path == "/v2LoggingOptions" && method == http.MethodPost:
		return "SetV2LoggingOptions"
	case strings.HasPrefix(path, "/violations/verification-state/") && method == http.MethodPost:
		return "PutVerificationStateOnViolation"
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
		if len(parts) == 3 && parts[2] == "shadow" {
			// Thing Shadow data-plane: operation is derived in ExtractOperation;
			// named shadows arrive via the "name" query parameter.
			if name := r.URL.Query().Get("name"); name != "" {
				params["shadowName"] = name
			}
		} else if len(parts) >= 3 && parts[2] == "principals" {
			// AttachThingPrincipal/DetachThingPrincipal send the principal in the
			// X-Amzn-Principal header (per the AWS SDK restJson1 binding).
			if p := r.Header.Get("X-Amzn-Principal"); p != "" {
				params["principal"] = p
			}
		} else if len(parts) >= 4 && parts[2] == "jobs" {
			params["jobId"] = parts[3]
			if len(parts) >= 6 && parts[4] == "executionNumber" {
				params["executionNumber"] = parts[5]
			}
		} else if method == http.MethodPost {
			params["_operation"] = "CreateThing"
		} else if method == http.MethodPatch {
			params["_operation"] = "UpdateThing"
		} else if method == http.MethodGet {
			params["_operation"] = "DescribeThing"
		} else if method == http.MethodDelete {
			params["_operation"] = "DeleteThing"
		}
	case strings.HasPrefix(path, "/topics/") && len(parts) >= 2:
		params["topic"] = strings.Join(parts[1:], "/")
	case strings.HasPrefix(path, "/api/things/shadow/ListNamedShadowsForThing/") && len(parts) >= 5:
		params["thingName"] = parts[4]
	case strings.HasPrefix(path, "/certificates/") && len(parts) >= 2:
		params["certificateId"] = parts[1]
	case strings.HasPrefix(path, "/policies/") && len(parts) >= 2:
		params["policyName"] = parts[1]
		// GetPolicyVersion/DeletePolicyVersion: /policies/{name}/version/{versionId}
		if len(parts) >= 4 && parts[2] == "version" {
			params["policyVersionId"] = parts[3]
		}
		if method == http.MethodPost {
			params["_operation"] = "CreatePolicy"
		}
	case strings.HasPrefix(path, "/destinations/") && len(parts) >= 2:
		// {arn+} is a greedy label; capture the entire remainder of the path.
		params["arn"] = strings.TrimPrefix(path, "/destinations/")
	case strings.HasPrefix(path, "/confirmdestination/") && len(parts) >= 2:
		// {confirmationToken+} is greedy; capture the entire remainder.
		params["confirmationToken"] = strings.TrimPrefix(path, "/confirmdestination/")
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
		if parts[1] != "addThingToBillingGroup" && parts[1] != "removeThingFromBillingGroup" {
			params["billingGroupName"] = parts[1]
		}
	case strings.HasPrefix(path, "/authorizers/") && len(parts) >= 2:
		params["authorizerName"] = parts[1]
	case strings.HasPrefix(path, "/authorizer/") && len(parts) >= 2:
		params["authorizerName"] = parts[1]
	case strings.HasPrefix(path, "/provisioning-templates/") && len(parts) >= 2:
		params["templateName"] = parts[1]
		if len(parts) >= 4 && parts[2] == "versions" {
			params["versionId"] = parts[3]
		}
	case strings.HasPrefix(path, "/domainConfigurations/") && len(parts) >= 2:
		params["domainConfigurationName"] = parts[1]
	case strings.HasPrefix(path, "/security-profiles/") && len(parts) >= 2:
		params["securityProfileName"] = parts[1]
	case strings.HasPrefix(path, "/policy-principals/") && len(parts) >= 2:
		params["policyName"] = parts[1]
	case strings.HasPrefix(path, "/target-policies/") && len(parts) >= 2:
		params["policyName"] = parts[1]
	case strings.HasPrefix(path, "/attached-policies/"):
		params["target"] = strings.TrimPrefix(path, "/attached-policies/")
	case strings.HasPrefix(path, "/principal-policies/") && len(parts) >= 2:
		params["policyName"] = parts[1]
		// AttachPrincipalPolicy/DetachPrincipalPolicy send the principal via
		// the X-Amzn-Iot-Principal header (note: differs from AttachThingPrincipal
		// which uses X-Amzn-Principal without the "Iot" segment).
		if p := r.Header.Get("X-Amzn-Iot-Principal"); p != "" {
			params["principal"] = p
		}

	// Path parameter extraction for newly-routed operations.
	case strings.HasPrefix(path, "/cacertificate/") && len(parts) >= 2:
		params["certificateId"] = parts[1]
	case strings.HasPrefix(path, "/job-templates/") && len(parts) >= 2:
		params["jobTemplateId"] = parts[1]
	case strings.HasPrefix(path, "/managed-job-templates/") && len(parts) >= 2:
		params["templateName"] = parts[1]
	case strings.HasPrefix(path, "/otaUpdates/") && len(parts) >= 2:
		params["otaUpdateId"] = parts[1]
	case strings.HasPrefix(path, "/dynamic-thing-groups/") && len(parts) >= 2:
		params["thingGroupName"] = parts[1]
	case strings.HasPrefix(path, "/custom-metric/") && len(parts) >= 2:
		params["metricName"] = parts[1]
	case strings.HasPrefix(path, "/fleet-metric/") && len(parts) >= 2:
		params["metricName"] = parts[1]
	case strings.HasPrefix(path, "/dimensions/") && len(parts) >= 2:
		params["name"] = parts[1]
	case strings.HasPrefix(path, "/mitigationactions/actions/") && len(parts) >= 3:
		params["actionName"] = parts[2]
	case strings.HasPrefix(path, "/audit/scheduledaudits/") && len(parts) >= 3:
		params["scheduledAuditName"] = parts[2]
	case strings.HasPrefix(path, "/audit/findings/") && len(parts) >= 3:
		params["findingId"] = parts[2]
	case strings.HasPrefix(path, "/audit/tasks/") && len(parts) >= 3:
		params["taskId"] = parts[2]
	case strings.HasPrefix(path, "/audit/mitigationactions/tasks/") && len(parts) >= 4:
		params["taskId"] = parts[3]
	case path == "/principals/things-v2":
		if p := r.Header.Get("X-Amzn-Principal"); p != "" {
			params["principal"] = p
		}
	case strings.HasPrefix(path, "/streams/") && len(parts) >= 2:
		params["streamId"] = parts[1]
	case strings.HasPrefix(path, "/indices/") && len(parts) >= 2:
		params["indexName"] = parts[1]
	case strings.HasPrefix(path, "/thing-registration-tasks/") && len(parts) >= 3:
		params["taskId"] = parts[2]
	case strings.HasPrefix(path, "/violations/") && len(parts) >= 2:
		params["violationId"] = parts[1]
	case strings.HasPrefix(path, "/certificates-by-ca/") && len(parts) >= 3:
		params["caCertificateId"] = parts[2]
	case strings.HasPrefix(path, "/accept-certificate-transfer/") && len(parts) >= 2:
		params["certificateId"] = parts[1]
	case strings.HasPrefix(path, "/cancel-certificate-transfer/") && len(parts) >= 2:
		params["certificateId"] = parts[1]
	case strings.HasPrefix(path, "/reject-certificate-transfer/") && len(parts) >= 2:
		params["certificateId"] = parts[1]
	case strings.HasPrefix(path, "/transfer-certificate/") && len(parts) >= 2:
		params["certificateId"] = parts[1]
	case strings.HasPrefix(path, "/certificate-providers/") && len(parts) >= 2:
		params["certificateProviderName"] = parts[1]
	case strings.HasPrefix(path, "/packages/") && len(parts) >= 2:
		params["packageName"] = parts[1]
		if len(parts) >= 4 && parts[2] == "versions" {
			params["versionName"] = parts[3]
		}
	case strings.HasPrefix(path, "/commands/") && len(parts) >= 2:
		params["commandId"] = parts[1]
	case strings.HasPrefix(path, "/command-executions/") && len(parts) >= 2:
		params["executionId"] = parts[1]
	case strings.HasPrefix(path, "/authorizer/") && len(parts) >= 2:
		params["authorizerName"] = parts[1]
	case strings.HasPrefix(path, "/policy-targets/") && len(parts) >= 2:
		params["policyName"] = parts[1]
	case strings.HasPrefix(path, "/detect/mitigationactions/tasks/") && len(parts) >= 4:
		params["taskId"] = parts[3]
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
