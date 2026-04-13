// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

var iamstoreKey struct{}

// IAMService provides IAM operations for managing users, groups, roles, and policies.
type IAMService struct {
	accountID string
	stores    sync.Map // global — single cached instance
}

// NewIAMService creates a new IAM service instance for the given account.
func NewIAMService(accountID string) *IAMService {
	return &IAMService{
		accountID: accountID,
	}
}

func (s *IAMService) store(reqCtx *request.RequestContext) (*iamstore.IAMStore, error) {
	if cached, ok := s.stores.Load(iamstoreKey); ok {
		if typed, ok := cached.(*iamstore.IAMStore); ok {
			return typed, nil
		}
	}
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage: %w", err)
	}
	store := iamstore.NewIAMStore(storage, s.accountID)
	if actual, loaded := s.stores.LoadOrStore(iamstoreKey, store); loaded {
		if typed, ok := actual.(*iamstore.IAMStore); ok {
			return typed, nil
		}
	}
	return store, nil
}

// RegisterHandlers registers all IAM operation handlers with the dispatcher.
func (s *IAMService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("iam", "GetUser", s.GetUser)
	d.RegisterHandlerForService("iam", "CreateUser", s.CreateUser)
	d.RegisterHandlerForService("iam", "DeleteUser", s.DeleteUser)
	d.RegisterHandlerForService("iam", "UpdateUser", s.UpdateUser)
	d.RegisterHandlerForService("iam", "ListUsers", s.ListUsers)
	d.RegisterHandlerForService("iam", "TagUser", s.TagUser)
	d.RegisterHandlerForService("iam", "UntagUser", s.UntagUser)
	d.RegisterHandlerForService("iam", "ListUserTags", s.ListUserTags)
	d.RegisterHandlerForService("iam", "PutUserPermissionsBoundary", s.PutUserPermissionsBoundary)
	d.RegisterHandlerForService("iam", "DeleteUserPermissionsBoundary", s.DeleteUserPermissionsBoundary)
	d.RegisterHandlerForService("iam", "GetLoginProfile", s.GetLoginProfile)
	d.RegisterHandlerForService("iam", "CreateLoginProfile", s.CreateLoginProfile)
	d.RegisterHandlerForService("iam", "DeleteLoginProfile", s.DeleteLoginProfile)
	d.RegisterHandlerForService("iam", "UpdateLoginProfile", s.UpdateLoginProfile)
	d.RegisterHandlerForService("iam", "ChangePassword", s.ChangePassword)
	d.RegisterHandlerForService("iam", "CreateAccessKey", s.CreateAccessKey)
	d.RegisterHandlerForService("iam", "DeleteAccessKey", s.DeleteAccessKey)
	d.RegisterHandlerForService("iam", "ListAccessKeys", s.ListAccessKeys)
	d.RegisterHandlerForService("iam", "GetAccessKeyLastUsed", s.GetAccessKeyLastUsed)
	d.RegisterHandlerForService("iam", "UpdateAccessKey", s.UpdateAccessKey)
	d.RegisterHandlerForService("iam", "CreateGroup", s.CreateGroup)
	d.RegisterHandlerForService("iam", "GetGroup", s.GetGroup)
	d.RegisterHandlerForService("iam", "UpdateGroup", s.UpdateGroup)
	d.RegisterHandlerForService("iam", "DeleteGroup", s.DeleteGroup)
	d.RegisterHandlerForService("iam", "ListGroups", s.ListGroups)
	d.RegisterHandlerForService("iam", "ListGroupsForUser", s.ListGroupsForUser)
	d.RegisterHandlerForService("iam", "AddUserToGroup", s.AddUserToGroup)
	d.RegisterHandlerForService("iam", "RemoveUserFromGroup", s.RemoveUserFromGroup)
	d.RegisterHandlerForService("iam", "TagGroup", s.TagGroup)
	d.RegisterHandlerForService("iam", "UntagGroup", s.UntagGroup)
	d.RegisterHandlerForService("iam", "ListGroupTags", s.ListGroupTags)
	d.RegisterHandlerForService("iam", "CreateRole", s.CreateRole)
	d.RegisterHandlerForService("iam", "GetRole", s.GetRole)
	d.RegisterHandlerForService("iam", "UpdateRole", s.UpdateRole)
	d.RegisterHandlerForService("iam", "UpdateRoleDescription", s.UpdateRoleDescription)
	d.RegisterHandlerForService("iam", "DeleteRole", s.DeleteRole)
	d.RegisterHandlerForService("iam", "ListRoles", s.ListRoles)
	d.RegisterHandlerForService("iam", "UpdateAssumeRolePolicy", s.UpdateAssumeRolePolicy)
	d.RegisterHandlerForService("iam", "TagRole", s.TagRole)
	d.RegisterHandlerForService("iam", "UntagRole", s.UntagRole)
	d.RegisterHandlerForService("iam", "ListRoleTags", s.ListRoleTags)
	d.RegisterHandlerForService("iam", "PutRolePermissionsBoundary", s.PutRolePermissionsBoundary)
	d.RegisterHandlerForService("iam", "DeleteRolePermissionsBoundary", s.DeleteRolePermissionsBoundary)
	d.RegisterHandlerForService("iam", "ListInstanceProfilesForRole", s.ListInstanceProfilesForRole)
	d.RegisterHandlerForService("iam", "CreateInstanceProfile", s.CreateInstanceProfile)
	d.RegisterHandlerForService("iam", "GetInstanceProfile", s.GetInstanceProfile)
	d.RegisterHandlerForService("iam", "DeleteInstanceProfile", s.DeleteInstanceProfile)
	d.RegisterHandlerForService("iam", "ListInstanceProfiles", s.ListInstanceProfiles)
	d.RegisterHandlerForService("iam", "AddRoleToInstanceProfile", s.AddRoleToInstanceProfile)
	d.RegisterHandlerForService("iam", "RemoveRoleFromInstanceProfile", s.RemoveRoleFromInstanceProfile)
	d.RegisterHandlerForService("iam", "ListInstanceProfileTags", s.ListInstanceProfileTags)
	d.RegisterHandlerForService("iam", "TagInstanceProfile", s.TagInstanceProfile)
	d.RegisterHandlerForService("iam", "UntagInstanceProfile", s.UntagInstanceProfile)
	d.RegisterHandlerForService("iam", "CreatePolicy", s.CreatePolicy)
	d.RegisterHandlerForService("iam", "GetPolicy", s.GetPolicy)
	d.RegisterHandlerForService("iam", "DeletePolicy", s.DeletePolicy)
	d.RegisterHandlerForService("iam", "ListPolicies", s.ListPolicies)
	d.RegisterHandlerForService("iam", "CreatePolicyVersion", s.CreatePolicyVersion)
	d.RegisterHandlerForService("iam", "GetPolicyVersion", s.GetPolicyVersion)
	d.RegisterHandlerForService("iam", "DeletePolicyVersion", s.DeletePolicyVersion)
	d.RegisterHandlerForService("iam", "ListPolicyVersions", s.ListPolicyVersions)
	d.RegisterHandlerForService("iam", "SetDefaultPolicyVersion", s.SetDefaultPolicyVersion)
	d.RegisterHandlerForService("iam", "TagPolicy", s.TagPolicy)
	d.RegisterHandlerForService("iam", "UntagPolicy", s.UntagPolicy)
	d.RegisterHandlerForService("iam", "ListPolicyTags", s.ListPolicyTags)
	d.RegisterHandlerForService("iam", "AttachUserPolicy", s.AttachUserPolicy)
	d.RegisterHandlerForService("iam", "DetachUserPolicy", s.DetachUserPolicy)
	d.RegisterHandlerForService("iam", "ListAttachedUserPolicies", s.ListAttachedUserPolicies)
	d.RegisterHandlerForService("iam", "AttachGroupPolicy", s.AttachGroupPolicy)
	d.RegisterHandlerForService("iam", "DetachGroupPolicy", s.DetachGroupPolicy)
	d.RegisterHandlerForService("iam", "ListAttachedGroupPolicies", s.ListAttachedGroupPolicies)
	d.RegisterHandlerForService("iam", "AttachRolePolicy", s.AttachRolePolicy)
	d.RegisterHandlerForService("iam", "DetachRolePolicy", s.DetachRolePolicy)
	d.RegisterHandlerForService("iam", "ListAttachedRolePolicies", s.ListAttachedRolePolicies)
	d.RegisterHandlerForService("iam", "PutUserPolicy", s.PutUserPolicy)
	d.RegisterHandlerForService("iam", "GetUserPolicy", s.GetUserPolicy)
	d.RegisterHandlerForService("iam", "DeleteUserPolicy", s.DeleteUserPolicy)
	d.RegisterHandlerForService("iam", "ListUserPolicies", s.ListUserPolicies)
	d.RegisterHandlerForService("iam", "PutGroupPolicy", s.PutGroupPolicy)
	d.RegisterHandlerForService("iam", "GetGroupPolicy", s.GetGroupPolicy)
	d.RegisterHandlerForService("iam", "DeleteGroupPolicy", s.DeleteGroupPolicy)
	d.RegisterHandlerForService("iam", "ListGroupPolicies", s.ListGroupPolicies)
	d.RegisterHandlerForService("iam", "PutRolePolicy", s.PutRolePolicy)
	d.RegisterHandlerForService("iam", "GetRolePolicy", s.GetRolePolicy)
	d.RegisterHandlerForService("iam", "DeleteRolePolicy", s.DeleteRolePolicy)
	d.RegisterHandlerForService("iam", "ListRolePolicies", s.ListRolePolicies)
	d.RegisterHandlerForService("iam", "CreateVirtualMFADevice", s.CreateVirtualMFADevice)
	d.RegisterHandlerForService("iam", "DeleteVirtualMFADevice", s.DeleteVirtualMFADevice)
	d.RegisterHandlerForService("iam", "EnableMFADevice", s.EnableMFADevice)
	d.RegisterHandlerForService("iam", "DeactivateMFADevice", s.DeactivateMFADevice)
	d.RegisterHandlerForService("iam", "ListMFADevices", s.ListMFADevices)
	d.RegisterHandlerForService("iam", "ListVirtualMFADevices", s.ListVirtualMFADevices)
	d.RegisterHandlerForService("iam", "ResyncMFADevice", s.ResyncMFADevice)
	d.RegisterHandlerForService("iam", "TagMFADevice", s.TagMFADevice)
	d.RegisterHandlerForService("iam", "UntagMFADevice", s.UntagMFADevice)
	d.RegisterHandlerForService("iam", "ListMFADeviceTags", s.ListMFADeviceTags)
	d.RegisterHandlerForService("iam", "GetAccountPasswordPolicy", s.GetAccountPasswordPolicy)
	d.RegisterHandlerForService("iam", "UpdateAccountPasswordPolicy", s.UpdateAccountPasswordPolicy)
	d.RegisterHandlerForService("iam", "DeleteAccountPasswordPolicy", s.DeleteAccountPasswordPolicy)
	d.RegisterHandlerForService("iam", "GetAccountSummary", s.GetAccountSummary)
	d.RegisterHandlerForService("iam", "CreateAccountAlias", s.CreateAccountAlias)
	d.RegisterHandlerForService("iam", "DeleteAccountAlias", s.DeleteAccountAlias)
	d.RegisterHandlerForService("iam", "ListAccountAliases", s.ListAccountAliases)
	d.RegisterHandlerForService("iam", "UploadServerCertificate", s.UploadServerCertificate)
	d.RegisterHandlerForService("iam", "GetServerCertificate", s.GetServerCertificate)
	d.RegisterHandlerForService("iam", "UpdateServerCertificate", s.UpdateServerCertificate)
	d.RegisterHandlerForService("iam", "DeleteServerCertificate", s.DeleteServerCertificate)
	d.RegisterHandlerForService("iam", "ListServerCertificates", s.ListServerCertificates)
	d.RegisterHandlerForService("iam", "TagServerCertificate", s.TagServerCertificate)
	d.RegisterHandlerForService("iam", "UntagServerCertificate", s.UntagServerCertificate)
	d.RegisterHandlerForService("iam", "ListServerCertificateTags", s.ListServerCertificateTags)
	d.RegisterHandlerForService("iam", "UploadSigningCertificate", s.UploadSigningCertificate)
	d.RegisterHandlerForService("iam", "ListSigningCertificates", s.ListSigningCertificates)
	d.RegisterHandlerForService("iam", "UpdateSigningCertificate", s.UpdateSigningCertificate)
	d.RegisterHandlerForService("iam", "DeleteSigningCertificate", s.DeleteSigningCertificate)
	d.RegisterHandlerForService("iam", "UploadSSHPublicKey", s.UploadSSHPublicKey)
	d.RegisterHandlerForService("iam", "GetSSHPublicKey", s.GetSSHPublicKey)
	d.RegisterHandlerForService("iam", "UpdateSSHPublicKey", s.UpdateSSHPublicKey)
	d.RegisterHandlerForService("iam", "ListSSHPublicKeys", s.ListSSHPublicKeys)
	d.RegisterHandlerForService("iam", "DeleteSSHPublicKey", s.DeleteSSHPublicKey)
	d.RegisterHandlerForService("iam", "CreateServiceSpecificCredential", s.CreateServiceSpecificCredential)
	d.RegisterHandlerForService("iam", "DeleteServiceSpecificCredential", s.DeleteServiceSpecificCredential)
	d.RegisterHandlerForService("iam", "ListServiceSpecificCredentials", s.ListServiceSpecificCredentials)
	d.RegisterHandlerForService("iam", "ResetServiceSpecificCredential", s.ResetServiceSpecificCredential)
	d.RegisterHandlerForService("iam", "UpdateServiceSpecificCredential", s.UpdateServiceSpecificCredential)
	d.RegisterHandlerForService("iam", "CreateSAMLProvider", s.CreateSAMLProvider)
	d.RegisterHandlerForService("iam", "GetSAMLProvider", s.GetSAMLProvider)
	d.RegisterHandlerForService("iam", "ListSAMLProviders", s.ListSAMLProviders)
	d.RegisterHandlerForService("iam", "UpdateSAMLProvider", s.UpdateSAMLProvider)
	d.RegisterHandlerForService("iam", "DeleteSAMLProvider", s.DeleteSAMLProvider)
	d.RegisterHandlerForService("iam", "TagSAMLProvider", s.TagSAMLProvider)
	d.RegisterHandlerForService("iam", "UntagSAMLProvider", s.UntagSAMLProvider)
	d.RegisterHandlerForService("iam", "ListSAMLProviderTags", s.ListSAMLProviderTags)
	d.RegisterHandlerForService("iam", "CreateOpenIDConnectProvider", s.CreateOpenIDConnectProvider)
	d.RegisterHandlerForService("iam", "GetOpenIDConnectProvider", s.GetOpenIDConnectProvider)
	d.RegisterHandlerForService("iam", "ListOpenIDConnectProviders", s.ListOpenIDConnectProviders)
	d.RegisterHandlerForService("iam", "UpdateOpenIDConnectProviderThumbprint", s.UpdateOpenIDConnectProviderThumbprint)
	d.RegisterHandlerForService("iam", "AddClientIDToOpenIDConnectProvider", s.AddClientIDToOpenIDConnectProvider)
	d.RegisterHandlerForService("iam", "RemoveClientIDFromOpenIDConnectProvider", s.RemoveClientIDFromOpenIDConnectProvider)
	d.RegisterHandlerForService("iam", "DeleteOpenIDConnectProvider", s.DeleteOpenIDConnectProvider)
	d.RegisterHandlerForService("iam", "TagOpenIDConnectProvider", s.TagOpenIDConnectProvider)
	d.RegisterHandlerForService("iam", "UntagOpenIDConnectProvider", s.UntagOpenIDConnectProvider)
	d.RegisterHandlerForService("iam", "ListOpenIDConnectProviderTags", s.ListOpenIDConnectProviderTags)
	d.RegisterHandlerForService("iam", "GetMFADevice", s.GetMFADevice)
	d.RegisterHandlerForService("iam", "EnableOutboundWebIdentityFederation", s.EnableOutboundWebIdentityFederation)
	d.RegisterHandlerForService("iam", "DisableOutboundWebIdentityFederation", s.DisableOutboundWebIdentityFederation)
	d.RegisterHandlerForService("iam", "GetOutboundWebIdentityFederationInfo", s.GetOutboundWebIdentityFederationInfo)
	d.RegisterHandlerForService("iam", "SetSecurityTokenServicePreferences", s.SetSecurityTokenServicePreferences)
	d.RegisterHandlerForService("iam", "GetAccountAuthorizationDetails", s.GetAccountAuthorizationDetails)
	d.RegisterHandlerForService("iam", "ListEntitiesForPolicy", s.ListEntitiesForPolicy)
	d.RegisterHandlerForService("iam", "GenerateCredentialReport", s.GenerateCredentialReport)
	d.RegisterHandlerForService("iam", "GetCredentialReport", s.GetCredentialReport)
	d.RegisterHandlerForService("iam", "CreateServiceLinkedRole", s.CreateServiceLinkedRole)
	d.RegisterHandlerForService("iam", "DeleteServiceLinkedRole", s.DeleteServiceLinkedRole)
	d.RegisterHandlerForService("iam", "GetServiceLinkedRoleDeletionStatus", s.GetServiceLinkedRoleDeletionStatus)
	d.RegisterHandlerForService("iam", "GenerateServiceLastAccessedDetails", s.GenerateServiceLastAccessedDetails)
	d.RegisterHandlerForService("iam", "GetServiceLastAccessedDetails", s.GetServiceLastAccessedDetails)
	d.RegisterHandlerForService("iam", "GetServiceLastAccessedDetailsWithEntities", s.GetServiceLastAccessedDetailsWithEntities)
	d.RegisterHandlerForService("iam", "SimulatePrincipalPolicy", s.SimulatePrincipalPolicy)
}
