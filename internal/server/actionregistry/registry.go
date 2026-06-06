package actionregistry

import "sync"

// ActionRegistry maps AWS actions to their service names. Multiple
// services may register the same action (e.g. TagResource for both
// SNS and Neptune). Lookup returns the last registered service, while
// LookupAll returns all candidates for disambiguation by the classifier.
type ActionRegistry struct {
	mu      sync.RWMutex
	actions map[string]string
	all     map[string][]string
}

var globalActionRegistry = NewActionRegistry()

// NewActionRegistry creates a new action registry with default mappings.
func NewActionRegistry() *ActionRegistry {
	r := &ActionRegistry{
		actions: make(map[string]string),
		all:     make(map[string][]string),
	}
	r.initDefaults()
	return r
}

// GetActionRegistry returns the global action registry.
func GetActionRegistry() *ActionRegistry {
	return globalActionRegistry
}

// Register registers actions for a service.
func (r *ActionRegistry) Register(service string, actions []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, action := range actions {
		r.actions[action] = service
		r.all[action] = append(r.all[action], service)
	}
}

// Lookup returns the last registered service name for a given action.
// For actions registered by a single service this is deterministic.
func (r *ActionRegistry) Lookup(action string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.actions[action]
}

// LookupAll returns all service names that registered the given action,
// in registration order. The classifier uses this to disambiguate when
// multiple services share an action name.
func (r *ActionRegistry) LookupAll(action string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	candidates := r.all[action]
	if candidates == nil {
		return nil
	}
	result := make([]string, len(candidates))
	copy(result, candidates)
	return result
}

func (r *ActionRegistry) initDefaults() {
	r.Register("iam", []string{
		"CreateUser", "GetUser", "UpdateUser", "DeleteUser", "ListUsers",
		"CreateGroup", "GetGroup", "UpdateGroup", "DeleteGroup", "ListGroups",
		"CreateRole", "GetRole", "UpdateRole", "DeleteRole", "ListRoles",
		"CreatePolicy", "GetPolicy", "DeletePolicy", "ListPolicies",
		"CreatePolicyVersion", "GetPolicyVersion", "DeletePolicyVersion", "ListPolicyVersions", "SetDefaultPolicyVersion",
		"AttachUserPolicy", "DetachUserPolicy", "ListAttachedUserPolicies",
		"AttachGroupPolicy", "DetachGroupPolicy", "ListAttachedGroupPolicies",
		"AttachRolePolicy", "DetachRolePolicy", "ListAttachedRolePolicies",
		"PutUserPolicy", "GetUserPolicy", "DeleteUserPolicy", "ListUserPolicies",
		"PutGroupPolicy", "GetGroupPolicy", "DeleteGroupPolicy", "ListGroupPolicies",
		"PutRolePolicy", "GetRolePolicy", "DeleteRolePolicy", "ListRolePolicies",
		"AddUserToGroup", "RemoveUserFromGroup", "ListGroupsForUser",
		"CreateAccessKey", "DeleteAccessKey", "ListAccessKeys", "GetAccessKeyLastUsed", "UpdateAccessKey",
		"CreateLoginProfile", "DeleteLoginProfile", "GetLoginProfile", "UpdateLoginProfile", "ChangePassword",
		"TagUser", "UntagUser", "ListUserTags",
		"TagGroup", "UntagGroup", "ListGroupTags",
		"TagRole", "UntagRole", "ListRoleTags",
		"TagPolicy", "UntagPolicy", "ListPolicyTags",
		"TagInstanceProfile", "UntagInstanceProfile", "ListInstanceProfileTags",
		"CreateInstanceProfile", "GetInstanceProfile", "DeleteInstanceProfile", "ListInstanceProfiles",
		"AddRoleToInstanceProfile", "RemoveRoleFromInstanceProfile", "ListInstanceProfilesForRole",
		"PutUserPermissionsBoundary", "DeleteUserPermissionsBoundary",
		"PutRolePermissionsBoundary", "DeleteRolePermissionsBoundary",
		"UpdateAssumeRolePolicy", "UpdateRoleDescription",
		"CreateVirtualMFADevice", "DeleteVirtualMFADevice", "EnableMFADevice", "DeactivateMFADevice",
		"ListMFADevices", "ListVirtualMFADevices", "GetMFADevice", "ResyncMFADevice",
		"TagMFADevice", "UntagMFADevice", "ListMFADeviceTags",
		"GetAccountSummary", "GetAccountPasswordPolicy", "UpdateAccountPasswordPolicy", "DeleteAccountPasswordPolicy",
		"CreateAccountAlias", "DeleteAccountAlias", "ListAccountAliases",
		"UploadServerCertificate", "GetServerCertificate", "UpdateServerCertificate", "DeleteServerCertificate", "ListServerCertificates",
		"TagServerCertificate", "UntagServerCertificate", "ListServerCertificateTags",
		"UploadSigningCertificate", "ListSigningCertificates", "UpdateSigningCertificate", "DeleteSigningCertificate",
		"UploadSSHPublicKey", "GetSSHPublicKey", "UpdateSSHPublicKey", "ListSSHPublicKeys", "DeleteSSHPublicKey",
		"CreateServiceSpecificCredential", "DeleteServiceSpecificCredential", "ListServiceSpecificCredentials",
		"ResetServiceSpecificCredential", "UpdateServiceSpecificCredential",
		"CreateSAMLProvider", "GetSAMLProvider", "ListSAMLProviders", "UpdateSAMLProvider", "DeleteSAMLProvider",
		"TagSAMLProvider", "UntagSAMLProvider", "ListSAMLProviderTags",
		"CreateOpenIDConnectProvider", "GetOpenIDConnectProvider", "ListOpenIDConnectProviders",
		"UpdateOpenIDConnectProviderThumbprint", "DeleteOpenIDConnectProvider",
		"AddClientIDToOpenIDConnectProvider", "RemoveClientIDFromOpenIDConnectProvider",
		"TagOpenIDConnectProvider", "UntagOpenIDConnectProvider", "ListOpenIDConnectProviderTags",
		"ListEntitiesForPolicy", "ListPoliciesGrantingServiceAccess",
		"GenerateCredentialReport", "GetCredentialReport",
		"GetAccountAuthorizationDetails",
		"CreateServiceLinkedRole", "DeleteServiceLinkedRole", "GetServiceLinkedRoleDeletionStatus",
		"EnableOutboundWebIdentityFederation", "DisableOutboundWebIdentityFederation",
		"GetOutboundWebIdentityFederationInfo", "SetSecurityTokenServicePreferences",
	})

	r.Register("sns", []string{
		"CreateTopic", "DeleteTopic", "GetTopicAttributes", "SetTopicAttributes", "ListTopics",
		"Subscribe", "Unsubscribe", "ConfirmSubscription",
		"GetSubscriptionAttributes", "SetSubscriptionAttributes",
		"ListSubscriptions", "ListSubscriptionsByTopic",
		"Publish", "PublishBatch",
		"CreatePlatformApplication", "DeletePlatformApplication",
		"GetPlatformApplicationAttributes", "SetPlatformApplicationAttributes",
		"ListPlatformApplications",
		"CreatePlatformEndpoint", "DeleteEndpoint",
		"GetEndpointAttributes", "SetEndpointAttributes",
		"ListEndpointsByPlatformApplication",
		"TagResource", "UntagResource", "ListTagsForResource",
		"AddPermission", "RemovePermission",
		"GetDataProtectionPolicy", "PutDataProtectionPolicy",
	})

	r.Register("sqs", []string{
		"CreateQueue", "DeleteQueue", "GetQueueUrl", "ListQueues",
		"GetQueueAttributes", "SetQueueAttributes", "PurgeQueue",
		"SendMessage", "SendMessageBatch", "ReceiveMessage",
		"DeleteMessage", "DeleteMessageBatch",
		"ChangeMessageVisibility", "ChangeMessageVisibilityBatch",
		"AddPermission", "RemovePermission",
		"TagQueue", "UntagQueue", "ListQueueTags",
		"ListDeadLetterSourceQueues",
	})

	r.Register("sts", []string{
		"GetCallerIdentity", "GetSessionToken", "GetFederationToken",
		"AssumeRole", "AssumeRoleWithSAML", "AssumeRoleWithWebIdentity",
		"AssumeRoot", "DecodeAuthorizationMessage", "GetAccessKeyInfo",
		"GetDelegatedAccessToken", "GetWebIdentityToken",
	})

	r.Register("events", []string{
		"CreateEventBus", "DeleteEventBus", "DescribeEventBus",
		"ListEventBuses", "UpdateEventBus",
		"PutRule", "DeleteRule", "DescribeRule",
		"ListRules", "EnableRule", "DisableRule",
		"PutTargets", "RemoveTargets", "ListTargetsByRule",
		"ListRuleNamesByTarget",
		"PutEvents",
		"CreateArchive", "DeleteArchive", "DescribeArchive",
		"ListArchives", "UpdateArchive",
		"CreateConnection", "DeleteConnection", "DescribeConnection",
		"UpdateConnection", "DeauthorizeConnection", "ListConnections",
		"CreateApiDestination", "DeleteApiDestination", "DescribeApiDestination",
		"UpdateApiDestination", "ListApiDestinations",
		"StartReplay", "DescribeReplay", "ListReplays", "CancelReplay",
		"TestEventPattern",
		// TagResource, UntagResource, ListTagsForResource intentionally
		// excluded: EventBridge uses JSON-RPC (X-Amz-Target), never
		// action-based routing. These would overwrite SNS Query entries.
	})

	r.Register("states", []string{
		"CreateStateMachine", "DeleteStateMachine", "DescribeStateMachine",
		"ListStateMachines", "UpdateStateMachine",
		"StartExecution", "StopExecution", "DescribeExecution",
		"ListExecutions", "GetExecutionHistory",
		"CreateActivity", "DeleteActivity", "DescribeActivity",
		"ListActivities",
		"RedriveExecution", "TestState",
	})

	r.Register("acm", []string{
		"RequestCertificate", "GetCertificate", "DeleteCertificate",
		"ListCertificates", "DescribeCertificate",
		"ResendValidationEmail",
		"AddTagsToCertificate", "RemoveTagsFromCertificate", "ListTagsForCertificate",
		"ImportCertificate", "ExportCertificate",
		"GetAccountConfiguration", "PutAccountConfiguration",
		"UpdateCertificateOptions", "RenewCertificate", "RevokeCertificate",
	})

	r.Register("cloudtrail", []string{
		"CreateTrail", "DeleteTrail", "UpdateTrail", "GetTrail",
		"DescribeTrails", "GetTrailStatus", "ListTrails",
		"StartLogging", "StopLogging",
		"GetEventSelectors", "PutEventSelectors",
		"GetInsightSelectors", "PutInsightSelectors",
		"AddTags", "RemoveTags", "ListTags",
		// TagResource, UntagResource, ListTagsForResource intentionally
		// excluded: CloudTrail uses JSON-RPC (X-Amz-Target), never
		// action-based routing. These would overwrite SNS Query entries.
		"LookupEvents", "ListPublicKeys",
		"GetResourcePolicy", "PutResourcePolicy", "DeleteResourcePolicy",
	})

	r.Register("lambda", []string{
		"CreateFunction", "DeleteFunction", "GetFunction", "GetFunctionConfiguration", "ListFunctions",
		"UpdateFunctionCode", "UpdateFunctionConfiguration",
		"Invoke", "InvokeAsync", "InvokeWithResponseStream",
		"PublishVersion", "ListVersionsByFunction",
		"CreateAlias", "DeleteAlias", "GetAlias", "UpdateAlias", "ListAliases",
		"PublishLayerVersion", "DeleteLayerVersion", "GetLayerVersion", "ListLayers", "ListLayerVersions",
		"CreateEventSourceMapping", "DeleteEventSourceMapping", "GetEventSourceMapping", "UpdateEventSourceMapping", "ListEventSourceMappings",
		// Note: AddPermission, RemovePermission, TagResource, UntagResource,
		// and ListTags are NOT registered here because Lambda uses REST-JSON
		// paths (e.g. /2015-03-31/functions/.../tags). These actions are
		// always routed via path lookup or signing service, never via
		// ActionRegistry. Registering them would overwrite SNS's entries for
		// the same action names, causing misrouting for SNS Query protocol
		// requests (see classifier.serviceFromAction).
		"GetPolicy",
		"PutFunctionConcurrency", "GetFunctionConcurrency", "DeleteFunctionConcurrency",
		"PutProvisionedConcurrencyConfig", "GetProvisionedConcurrencyConfig", "DeleteProvisionedConcurrencyConfig", "ListProvisionedConcurrencyConfigs",
		"PutFunctionEventInvokeConfig", "GetFunctionEventInvokeConfig", "DeleteFunctionEventInvokeConfig", "ListFunctionEventInvokeConfigs",
		"CreateFunctionUrlConfig", "DeleteFunctionUrlConfig", "GetFunctionUrlConfig", "UpdateFunctionUrlConfig", "ListFunctionUrlConfigs",
		"GetAccountSettings",
	})

	r.Register("neptune", []string{
		"CreateDBCluster", "DeleteDBCluster", "ModifyDBCluster", "DescribeDBClusters",
		"StartDBCluster", "StopDBCluster", "FailoverDBCluster",
		"CreateDBClusterEndpoint", "DescribeDBClusterEndpoints", "ModifyDBClusterEndpoint", "DeleteDBClusterEndpoint",
		"CreateDBInstance", "DeleteDBInstance", "ModifyDBInstance", "DescribeDBInstances", "RebootDBInstance",
		"CreateDBClusterSnapshot", "DeleteDBClusterSnapshot", "DescribeDBClusterSnapshots",
		"CopyDBClusterSnapshot", "DescribeDBClusterSnapshotAttributes", "ModifyDBClusterSnapshotAttribute",
		"CreateDBClusterParameterGroup", "DeleteDBClusterParameterGroup",
		"DescribeDBClusterParameterGroups", "DescribeDBClusterParameters", "ModifyDBClusterParameterGroup",
		"CreateDBParameterGroup", "DeleteDBParameterGroup",
		"DescribeDBParameterGroups", "DescribeDBParameters", "ModifyDBParameterGroup",
		"ResetDBClusterParameterGroup", "ResetDBParameterGroup",
		"CopyDBClusterParameterGroup", "CopyDBParameterGroup",
		"DescribeEngineDefaultClusterParameters", "DescribeEngineDefaultParameters",
		"CreateGlobalCluster", "DeleteGlobalCluster", "DescribeGlobalClusters",
		"ModifyGlobalCluster", "FailoverGlobalCluster", "SwitchoverGlobalCluster", "RemoveFromGlobalCluster",
		"CreateDBSubnetGroup", "DeleteDBSubnetGroup", "DescribeDBSubnetGroups", "ModifyDBSubnetGroup",
		"CreateEventSubscription", "DeleteEventSubscription", "DescribeEventSubscriptions",
		"ModifyEventSubscription", "AddSourceIdentifierToSubscription", "RemoveSourceIdentifierFromSubscription",
		"AddTagsToResource", "ListTagsForResource", "RemoveTagsFromResource",
		"DescribeEvents", "DescribeEventCategories", "DescribePendingMaintenanceActions", "ApplyPendingMaintenanceAction",
		"DescribeDBEngineVersions", "DescribeOrderableDBInstanceOptions", "DescribeValidDBInstanceModifications",
		"RestoreDBClusterFromSnapshot", "RestoreDBClusterToPointInTime", "PromoteReadReplicaDBCluster",
		"AddRoleToDBCluster", "RemoveRoleFromDBCluster",
	})

	r.Register("ec2", []string{
		"CreateVpc", "DescribeVpcs", "DeleteVpc",
		"CreateSubnet", "DescribeSubnets", "DeleteSubnet",
		"CreateSecurityGroup", "DescribeSecurityGroups", "DeleteSecurityGroup",
		"AuthorizeSecurityGroupIngress", "AuthorizeSecurityGroupEgress",
		"RevokeSecurityGroupIngress", "RevokeSecurityGroupEgress",
	})
}

// LookupServiceByAction returns the service name for a given action using the global registry.
func LookupServiceByAction(action string) string {
	return globalActionRegistry.Lookup(action)
}

// LookupAllServicesByAction returns all service names that registered the given action.
func LookupAllServicesByAction(action string) []string {
	return globalActionRegistry.LookupAll(action)
}
