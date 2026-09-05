// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	storecommon "vorpalstacks/internal/store/aws/common"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// IAMService provides IAM operations for managing users, groups, roles, and policies.
type IAMService struct {
	accountID             string
	stores                sync.Map // global — single cached instance
	storageManager        *storage.RegionStorageManager
	cloudTrailInvoker     invokers.CloudTrailInvoker
	reportWg              sync.WaitGroup
	slRoleDeletionWg      sync.WaitGroup
	credentialReportMu    sync.RWMutex
	credentialReportState string
	credentialReportData  string
	credentialReportTime  time.Time
}

// NewIAMService creates a new IAM service instance for the given account.
func NewIAMService(accountID string) *IAMService {
	return &IAMService{
		accountID: accountID,
	}
}

// SetCloudTrailInvoker injects the CloudTrail invoker for last-accessed analysis.
func (s *IAMService) SetCloudTrailInvoker(invoker invokers.CloudTrailInvoker) {
	s.cloudTrailInvoker = invoker
}

// SetStorageManager injects the region storage manager for lazy store creation.
func (s *IAMService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

func (s *IAMService) store(reqCtx *request.RequestContext) (*iamstore.IAMStore, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, "global", func() (*iamstore.IAMStore, error) {
		storage, err := reqCtx.GetGlobalStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get global storage: %w", err)
		}
		return iamstore.GetOrCreateGlobalStore(storage, s.accountID), nil
	})
}

// GetStoreForRegion returns the cached IAM store.
// IAM is a global service — the region parameter is ignored.
func (s *IAMService) GetStoreForRegion(_ string) (*iamstore.IAMStore, error) {
	if v, ok := s.stores.Load("global"); ok {
		return v.(*iamstore.IAMStore), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("IAM storage manager not initialised")
	}
	st, err := s.storageManager.GetGlobalStorage()
	if err != nil {
		return nil, err
	}
	store := iamstore.GetOrCreateGlobalStore(st, s.accountID)
	actual, _ := s.stores.LoadOrStore("global", store)
	return actual.(*iamstore.IAMStore), nil
}

// iamHandlerMethods maps every IAM operation to its handler method. The
// dispatcher registration iterates it, and the action registry guard test
// compares its keys with the registry entries in both directions.
var iamHandlerMethods = map[string]func(*IAMService, context.Context, *request.RequestContext, *request.ParsedRequest) (interface{}, error){
	"GetUser":                                   (*IAMService).GetUser,
	"CreateUser":                                (*IAMService).CreateUser,
	"DeleteUser":                                (*IAMService).DeleteUser,
	"UpdateUser":                                (*IAMService).UpdateUser,
	"ListUsers":                                 (*IAMService).ListUsers,
	"TagUser":                                   (*IAMService).TagUser,
	"UntagUser":                                 (*IAMService).UntagUser,
	"ListUserTags":                              (*IAMService).ListUserTags,
	"PutUserPermissionsBoundary":                (*IAMService).PutUserPermissionsBoundary,
	"DeleteUserPermissionsBoundary":             (*IAMService).DeleteUserPermissionsBoundary,
	"GetLoginProfile":                           (*IAMService).GetLoginProfile,
	"CreateLoginProfile":                        (*IAMService).CreateLoginProfile,
	"DeleteLoginProfile":                        (*IAMService).DeleteLoginProfile,
	"UpdateLoginProfile":                        (*IAMService).UpdateLoginProfile,
	"ChangePassword":                            (*IAMService).ChangePassword,
	"CreateAccessKey":                           (*IAMService).CreateAccessKey,
	"DeleteAccessKey":                           (*IAMService).DeleteAccessKey,
	"ListAccessKeys":                            (*IAMService).ListAccessKeys,
	"GetAccessKeyLastUsed":                      (*IAMService).GetAccessKeyLastUsed,
	"UpdateAccessKey":                           (*IAMService).UpdateAccessKey,
	"CreateGroup":                               (*IAMService).CreateGroup,
	"GetGroup":                                  (*IAMService).GetGroup,
	"UpdateGroup":                               (*IAMService).UpdateGroup,
	"DeleteGroup":                               (*IAMService).DeleteGroup,
	"ListGroups":                                (*IAMService).ListGroups,
	"ListGroupsForUser":                         (*IAMService).ListGroupsForUser,
	"AddUserToGroup":                            (*IAMService).AddUserToGroup,
	"RemoveUserFromGroup":                       (*IAMService).RemoveUserFromGroup,
	"CreateRole":                                (*IAMService).CreateRole,
	"GetRole":                                   (*IAMService).GetRole,
	"UpdateRole":                                (*IAMService).UpdateRole,
	"UpdateRoleDescription":                     (*IAMService).UpdateRoleDescription,
	"DeleteRole":                                (*IAMService).DeleteRole,
	"ListRoles":                                 (*IAMService).ListRoles,
	"UpdateAssumeRolePolicy":                    (*IAMService).UpdateAssumeRolePolicy,
	"TagRole":                                   (*IAMService).TagRole,
	"UntagRole":                                 (*IAMService).UntagRole,
	"ListRoleTags":                              (*IAMService).ListRoleTags,
	"PutRolePermissionsBoundary":                (*IAMService).PutRolePermissionsBoundary,
	"DeleteRolePermissionsBoundary":             (*IAMService).DeleteRolePermissionsBoundary,
	"ListInstanceProfilesForRole":               (*IAMService).ListInstanceProfilesForRole,
	"CreateInstanceProfile":                     (*IAMService).CreateInstanceProfile,
	"GetInstanceProfile":                        (*IAMService).GetInstanceProfile,
	"DeleteInstanceProfile":                     (*IAMService).DeleteInstanceProfile,
	"ListInstanceProfiles":                      (*IAMService).ListInstanceProfiles,
	"AddRoleToInstanceProfile":                  (*IAMService).AddRoleToInstanceProfile,
	"RemoveRoleFromInstanceProfile":             (*IAMService).RemoveRoleFromInstanceProfile,
	"ListInstanceProfileTags":                   (*IAMService).ListInstanceProfileTags,
	"TagInstanceProfile":                        (*IAMService).TagInstanceProfile,
	"UntagInstanceProfile":                      (*IAMService).UntagInstanceProfile,
	"CreatePolicy":                              (*IAMService).CreatePolicy,
	"GetPolicy":                                 (*IAMService).GetPolicy,
	"DeletePolicy":                              (*IAMService).DeletePolicy,
	"ListPolicies":                              (*IAMService).ListPolicies,
	"CreatePolicyVersion":                       (*IAMService).CreatePolicyVersion,
	"GetPolicyVersion":                          (*IAMService).GetPolicyVersion,
	"DeletePolicyVersion":                       (*IAMService).DeletePolicyVersion,
	"ListPolicyVersions":                        (*IAMService).ListPolicyVersions,
	"SetDefaultPolicyVersion":                   (*IAMService).SetDefaultPolicyVersion,
	"TagPolicy":                                 (*IAMService).TagPolicy,
	"UntagPolicy":                               (*IAMService).UntagPolicy,
	"ListPolicyTags":                            (*IAMService).ListPolicyTags,
	"AttachUserPolicy":                          (*IAMService).AttachUserPolicy,
	"DetachUserPolicy":                          (*IAMService).DetachUserPolicy,
	"ListAttachedUserPolicies":                  (*IAMService).ListAttachedUserPolicies,
	"AttachGroupPolicy":                         (*IAMService).AttachGroupPolicy,
	"DetachGroupPolicy":                         (*IAMService).DetachGroupPolicy,
	"ListAttachedGroupPolicies":                 (*IAMService).ListAttachedGroupPolicies,
	"AttachRolePolicy":                          (*IAMService).AttachRolePolicy,
	"DetachRolePolicy":                          (*IAMService).DetachRolePolicy,
	"ListAttachedRolePolicies":                  (*IAMService).ListAttachedRolePolicies,
	"PutUserPolicy":                             (*IAMService).PutUserPolicy,
	"GetUserPolicy":                             (*IAMService).GetUserPolicy,
	"DeleteUserPolicy":                          (*IAMService).DeleteUserPolicy,
	"ListUserPolicies":                          (*IAMService).ListUserPolicies,
	"PutGroupPolicy":                            (*IAMService).PutGroupPolicy,
	"GetGroupPolicy":                            (*IAMService).GetGroupPolicy,
	"DeleteGroupPolicy":                         (*IAMService).DeleteGroupPolicy,
	"ListGroupPolicies":                         (*IAMService).ListGroupPolicies,
	"PutRolePolicy":                             (*IAMService).PutRolePolicy,
	"GetRolePolicy":                             (*IAMService).GetRolePolicy,
	"DeleteRolePolicy":                          (*IAMService).DeleteRolePolicy,
	"ListRolePolicies":                          (*IAMService).ListRolePolicies,
	"CreateVirtualMFADevice":                    (*IAMService).CreateVirtualMFADevice,
	"DeleteVirtualMFADevice":                    (*IAMService).DeleteVirtualMFADevice,
	"EnableMFADevice":                           (*IAMService).EnableMFADevice,
	"DeactivateMFADevice":                       (*IAMService).DeactivateMFADevice,
	"ListMFADevices":                            (*IAMService).ListMFADevices,
	"ListVirtualMFADevices":                     (*IAMService).ListVirtualMFADevices,
	"ResyncMFADevice":                           (*IAMService).ResyncMFADevice,
	"TagMFADevice":                              (*IAMService).TagMFADevice,
	"UntagMFADevice":                            (*IAMService).UntagMFADevice,
	"ListMFADeviceTags":                         (*IAMService).ListMFADeviceTags,
	"GetAccountPasswordPolicy":                  (*IAMService).GetAccountPasswordPolicy,
	"UpdateAccountPasswordPolicy":               (*IAMService).UpdateAccountPasswordPolicy,
	"DeleteAccountPasswordPolicy":               (*IAMService).DeleteAccountPasswordPolicy,
	"GetAccountSummary":                         (*IAMService).GetAccountSummary,
	"CreateAccountAlias":                        (*IAMService).CreateAccountAlias,
	"DeleteAccountAlias":                        (*IAMService).DeleteAccountAlias,
	"ListAccountAliases":                        (*IAMService).ListAccountAliases,
	"UploadServerCertificate":                   (*IAMService).UploadServerCertificate,
	"GetServerCertificate":                      (*IAMService).GetServerCertificate,
	"UpdateServerCertificate":                   (*IAMService).UpdateServerCertificate,
	"DeleteServerCertificate":                   (*IAMService).DeleteServerCertificate,
	"ListServerCertificates":                    (*IAMService).ListServerCertificates,
	"TagServerCertificate":                      (*IAMService).TagServerCertificate,
	"UntagServerCertificate":                    (*IAMService).UntagServerCertificate,
	"ListServerCertificateTags":                 (*IAMService).ListServerCertificateTags,
	"UploadSigningCertificate":                  (*IAMService).UploadSigningCertificate,
	"ListSigningCertificates":                   (*IAMService).ListSigningCertificates,
	"UpdateSigningCertificate":                  (*IAMService).UpdateSigningCertificate,
	"DeleteSigningCertificate":                  (*IAMService).DeleteSigningCertificate,
	"UploadSSHPublicKey":                        (*IAMService).UploadSSHPublicKey,
	"GetSSHPublicKey":                           (*IAMService).GetSSHPublicKey,
	"UpdateSSHPublicKey":                        (*IAMService).UpdateSSHPublicKey,
	"ListSSHPublicKeys":                         (*IAMService).ListSSHPublicKeys,
	"DeleteSSHPublicKey":                        (*IAMService).DeleteSSHPublicKey,
	"CreateServiceSpecificCredential":           (*IAMService).CreateServiceSpecificCredential,
	"DeleteServiceSpecificCredential":           (*IAMService).DeleteServiceSpecificCredential,
	"ListServiceSpecificCredentials":            (*IAMService).ListServiceSpecificCredentials,
	"ResetServiceSpecificCredential":            (*IAMService).ResetServiceSpecificCredential,
	"UpdateServiceSpecificCredential":           (*IAMService).UpdateServiceSpecificCredential,
	"CreateSAMLProvider":                        (*IAMService).CreateSAMLProvider,
	"GetSAMLProvider":                           (*IAMService).GetSAMLProvider,
	"ListSAMLProviders":                         (*IAMService).ListSAMLProviders,
	"UpdateSAMLProvider":                        (*IAMService).UpdateSAMLProvider,
	"DeleteSAMLProvider":                        (*IAMService).DeleteSAMLProvider,
	"TagSAMLProvider":                           (*IAMService).TagSAMLProvider,
	"UntagSAMLProvider":                         (*IAMService).UntagSAMLProvider,
	"ListSAMLProviderTags":                      (*IAMService).ListSAMLProviderTags,
	"CreateOpenIDConnectProvider":               (*IAMService).CreateOpenIDConnectProvider,
	"GetOpenIDConnectProvider":                  (*IAMService).GetOpenIDConnectProvider,
	"ListOpenIDConnectProviders":                (*IAMService).ListOpenIDConnectProviders,
	"UpdateOpenIDConnectProviderThumbprint":     (*IAMService).UpdateOpenIDConnectProviderThumbprint,
	"AddClientIDToOpenIDConnectProvider":        (*IAMService).AddClientIDToOpenIDConnectProvider,
	"RemoveClientIDFromOpenIDConnectProvider":   (*IAMService).RemoveClientIDFromOpenIDConnectProvider,
	"DeleteOpenIDConnectProvider":               (*IAMService).DeleteOpenIDConnectProvider,
	"TagOpenIDConnectProvider":                  (*IAMService).TagOpenIDConnectProvider,
	"UntagOpenIDConnectProvider":                (*IAMService).UntagOpenIDConnectProvider,
	"ListOpenIDConnectProviderTags":             (*IAMService).ListOpenIDConnectProviderTags,
	"GetMFADevice":                              (*IAMService).GetMFADevice,
	"EnableOutboundWebIdentityFederation":       (*IAMService).EnableOutboundWebIdentityFederation,
	"DisableOutboundWebIdentityFederation":      (*IAMService).DisableOutboundWebIdentityFederation,
	"GetOutboundWebIdentityFederationInfo":      (*IAMService).GetOutboundWebIdentityFederationInfo,
	"SetSecurityTokenServicePreferences":        (*IAMService).SetSecurityTokenServicePreferences,
	"GetAccountAuthorizationDetails":            (*IAMService).GetAccountAuthorizationDetails,
	"ListEntitiesForPolicy":                     (*IAMService).ListEntitiesForPolicy,
	"GenerateCredentialReport":                  (*IAMService).GenerateCredentialReport,
	"GetCredentialReport":                       (*IAMService).GetCredentialReport,
	"CreateServiceLinkedRole":                   (*IAMService).CreateServiceLinkedRole,
	"DeleteServiceLinkedRole":                   (*IAMService).DeleteServiceLinkedRole,
	"GetServiceLinkedRoleDeletionStatus":        (*IAMService).GetServiceLinkedRoleDeletionStatus,
	"GenerateServiceLastAccessedDetails":        (*IAMService).GenerateServiceLastAccessedDetails,
	"GetServiceLastAccessedDetails":             (*IAMService).GetServiceLastAccessedDetails,
	"GetServiceLastAccessedDetailsWithEntities": (*IAMService).GetServiceLastAccessedDetailsWithEntities,
	"SimulatePrincipalPolicy":                   (*IAMService).SimulatePrincipalPolicy,
	"ListPoliciesGrantingServiceAccess":         (*IAMService).ListPoliciesGrantingServiceAccess,
}

// RegisterHandlers registers all IAM operation handlers with the dispatcher.
func (s *IAMService) RegisterHandlers(d handler.Registrar) {
	for op, method := range iamHandlerMethods {
		d.RegisterHandlerForService("iam", op, func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
			return method(s, ctx, reqCtx, req)
		})
	}
}

// RegisteredIAMOperations returns the sorted names of all operations with
// registered handlers.
func RegisteredIAMOperations() []string {
	ops := make([]string, 0, len(iamHandlerMethods))
	for op := range iamHandlerMethods {
		ops = append(ops, op)
	}
	sort.Strings(ops)
	return ops
}

// WaitForReport blocks until any in-flight credential report generation
// goroutines have finished. Call during shutdown.
func (s *IAMService) WaitForReport() {
	s.reportWg.Wait()
}

// WaitForSLRoleDeletions blocks until any in-flight service-linked role
// deletion goroutines have finished. Call during shutdown so that the
// background cleanup is not aborted mid-way, which would leave the role
// partially detached (instance profiles still attached, managed policies
// still linked, etc.). Recovery on next startup is possible but only
// marks the task as FAILED; it does not resume cleanup.
func (s *IAMService) WaitForSLRoleDeletions() {
	s.slRoleDeletionWg.Wait()
}
