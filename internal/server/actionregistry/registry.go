package actionregistry

import "sync"

// ActionRegistry maps AWS actions to their service names.
type ActionRegistry struct {
	mu      sync.RWMutex
	actions map[string]string
}

var globalActionRegistry = NewActionRegistry()

// NewActionRegistry creates a new action registry with default mappings.
func NewActionRegistry() *ActionRegistry {
	r := &ActionRegistry{
		actions: make(map[string]string),
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
	}
}

// Lookup returns the service name for a given action.
func (r *ActionRegistry) Lookup(action string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.actions[action]
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
		"DecodeAuthorizationMessage",
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
		"TagResource", "UntagResource", "ListTagsForResource",
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
		"TagResource", "UntagResource", "ListTagsForResource",
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
		"AddPermission", "RemovePermission", "GetPolicy",
		"TagResource", "UntagResource", "ListTags",
		"PutFunctionConcurrency", "GetFunctionConcurrency", "DeleteFunctionConcurrency",
		"PutProvisionedConcurrencyConcurrency", "GetProvisionedConcurrencyConfig", "DeleteProvisionedConcurrencyConcurrency", "ListProvisionedConcurrencyConfigs",
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
}

// LookupServiceByAction returns the service name for a given action using the global registry.
func LookupServiceByAction(action string) string {
	return globalActionRegistry.Lookup(action)
}
