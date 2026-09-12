package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
)

// cognitoIDPServiceName is the dispatcher service namespace of the Cognito
// user pools API.
const cognitoIDPServiceName = "cognito-idp"

// cognitoOperationMethod is the method-expression form of an operation
// handler: the unbound service method, applied to the service instance at
// registration time.
type cognitoOperationMethod func(s *CognitoService, ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error)

// userPoolsAPIOperation is one row of the user pools API registration
// table. wafInspected marks the operations whose requests AWS WAF inspects
// when a WebACL is associated with the addressed user pool — the operations
// that do not require authentication with AWS credentials. The table is the
// single source of truth: RegisterHandlers wraps exactly the marked rows,
// and wafInspectedOperations is derived from the same marks, so the
// registration and the enforcement set cannot drift apart.
type userPoolsAPIOperation struct {
	operation    string
	handler      cognitoOperationMethod
	wafInspected bool
}

// bind applies the unbound operation method to the service instance.
func (op userPoolsAPIOperation) bind(s *CognitoService) cognitoOperationHandler {
	return func(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
		return op.handler(s, ctx, reqCtx, req)
	}
}

// userPoolsAPIOperations is the complete registration table of the user
// pools API. Each row registers its operation exactly once; the dispatcher
// registry is last-write-wins, so a duplicated row would silently drop the
// earlier registration together with any WAF wrapper it carried.
var userPoolsAPIOperations = []userPoolsAPIOperation{
	{operation: "CreateUserPool", handler: (*CognitoService).CreateUserPool},
	{operation: "DeleteUserPool", handler: (*CognitoService).DeleteUserPool},
	{operation: "DescribeUserPool", handler: (*CognitoService).DescribeUserPool},
	{operation: "ListUserPools", handler: (*CognitoService).ListUserPools},
	{operation: "UpdateUserPool", handler: (*CognitoService).UpdateUserPool},

	{operation: "CreateUserPoolClient", handler: (*CognitoService).CreateUserPoolClient},
	{operation: "DeleteUserPoolClient", handler: (*CognitoService).DeleteUserPoolClient},
	{operation: "DescribeUserPoolClient", handler: (*CognitoService).DescribeUserPoolClient},
	{operation: "ListUserPoolClients", handler: (*CognitoService).ListUserPoolClients},
	{operation: "UpdateUserPoolClient", handler: (*CognitoService).UpdateUserPoolClient},

	{operation: "AdminCreateUser", handler: (*CognitoService).AdminCreateUser},
	{operation: "AdminDeleteUser", handler: (*CognitoService).AdminDeleteUser},
	{operation: "AdminDeleteUserAttributes", handler: (*CognitoService).AdminDeleteUserAttributes},
	{operation: "AdminDisableUser", handler: (*CognitoService).AdminDisableUser},
	{operation: "AdminEnableUser", handler: (*CognitoService).AdminEnableUser},
	{operation: "AdminGetUser", handler: (*CognitoService).AdminGetUser},
	{operation: "AdminResetUserPassword", handler: (*CognitoService).AdminResetUserPassword},
	{operation: "AdminSetUserPassword", handler: (*CognitoService).AdminSetUserPassword},
	{operation: "AdminUpdateUserAttributes", handler: (*CognitoService).AdminUpdateUserAttributes},
	{operation: "AdminUserGlobalSignOut", handler: (*CognitoService).AdminUserGlobalSignOut},
	{operation: "ListUsers", handler: (*CognitoService).ListUsers},

	{operation: "CreateGroup", handler: (*CognitoService).CreateGroup},
	{operation: "DeleteGroup", handler: (*CognitoService).DeleteGroup},
	{operation: "GetGroup", handler: (*CognitoService).GetGroup},
	{operation: "ListGroups", handler: (*CognitoService).ListGroups},
	{operation: "UpdateGroup", handler: (*CognitoService).UpdateGroup},
	{operation: "AdminAddUserToGroup", handler: (*CognitoService).AdminAddUserToGroup},
	{operation: "AdminRemoveUserFromGroup", handler: (*CognitoService).AdminRemoveUserFromGroup},
	{operation: "ListUsersInGroup", handler: (*CognitoService).ListUsersInGroup},
	{operation: "AdminListGroupsForUser", handler: (*CognitoService).AdminListGroupsForUser},

	// The public (credential-free) operations carry the WAF enforcement
	// wrapper: when a WebACL is associated with the addressed user pool,
	// AWS WAF inspects these requests.
	{operation: "SignUp", handler: (*CognitoService).SignUp, wafInspected: true},
	{operation: "ConfirmSignUp", handler: (*CognitoService).ConfirmSignUp, wafInspected: true},
	{operation: "InitiateAuth", handler: (*CognitoService).InitiateAuth, wafInspected: true},
	{operation: "RespondToAuthChallenge", handler: (*CognitoService).RespondToAuthChallenge, wafInspected: true},
	{operation: "GlobalSignOut", handler: (*CognitoService).GlobalSignOut, wafInspected: true},
	{operation: "ChangePassword", handler: (*CognitoService).ChangePassword, wafInspected: true},
	{operation: "ForgotPassword", handler: (*CognitoService).ForgotPassword, wafInspected: true},
	{operation: "ConfirmForgotPassword", handler: (*CognitoService).ConfirmForgotPassword, wafInspected: true},
	{operation: "GetUser", handler: (*CognitoService).GetUser, wafInspected: true},
	{operation: "DeleteUser", handler: (*CognitoService).DeleteUser, wafInspected: true},
	{operation: "DeleteUserAttributes", handler: (*CognitoService).DeleteUserAttributes, wafInspected: true},
	{operation: "UpdateUserAttributes", handler: (*CognitoService).UpdateUserAttributes, wafInspected: true},
	{operation: "AssociateSoftwareToken", handler: (*CognitoService).AssociateSoftwareToken, wafInspected: true},
	{operation: "VerifySoftwareToken", handler: (*CognitoService).VerifySoftwareToken, wafInspected: true},
	{operation: "AdminDeleteSoftwareToken", handler: (*CognitoService).AdminDeleteSoftwareToken},
	{operation: "AdminConfirmSignUp", handler: (*CognitoService).AdminConfirmSignUp},
	{operation: "AdminInitiateAuth", handler: (*CognitoService).AdminInitiateAuth},
	{operation: "AdminRespondToAuthChallenge", handler: (*CognitoService).AdminRespondToAuthChallenge},

	{operation: "TagResource", handler: (*CognitoService).TagResource},
	{operation: "UntagResource", handler: (*CognitoService).UntagResource},
	{operation: "ListTagsForResource", handler: (*CognitoService).ListTagsForResource},
	{operation: "GetUserPoolMfaConfig", handler: (*CognitoService).GetUserPoolMfaConfig},
	{operation: "SetUserPoolMfaConfig", handler: (*CognitoService).SetUserPoolMfaConfig},

	{operation: "CreateUserPoolDomain", handler: (*CognitoService).CreateUserPoolDomain},
	{operation: "DescribeUserPoolDomain", handler: (*CognitoService).DescribeUserPoolDomain},
	{operation: "DeleteUserPoolDomain", handler: (*CognitoService).DeleteUserPoolDomain},
	{operation: "UpdateUserPoolDomain", handler: (*CognitoService).UpdateUserPoolDomain},
	{operation: "CreateResourceServer", handler: (*CognitoService).CreateResourceServer},
	{operation: "DescribeResourceServer", handler: (*CognitoService).DescribeResourceServer},
	{operation: "UpdateResourceServer", handler: (*CognitoService).UpdateResourceServer},
	{operation: "DeleteResourceServer", handler: (*CognitoService).DeleteResourceServer},
	{operation: "ListResourceServers", handler: (*CognitoService).ListResourceServers},
	{operation: "CreateIdentityProvider", handler: (*CognitoService).CreateIdentityProvider},
	{operation: "DescribeIdentityProvider", handler: (*CognitoService).DescribeIdentityProvider},
	{operation: "UpdateIdentityProvider", handler: (*CognitoService).UpdateIdentityProvider},
	{operation: "DeleteIdentityProvider", handler: (*CognitoService).DeleteIdentityProvider},
	{operation: "ListIdentityProviders", handler: (*CognitoService).ListIdentityProviders},
	{operation: "GetCSVHeader", handler: (*CognitoService).GetCSVHeader},
	{operation: "DescribeRiskConfiguration", handler: (*CognitoService).DescribeRiskConfiguration},

	// Token & auth operations
	{operation: "RevokeToken", handler: (*CognitoService).RevokeToken, wafInspected: true},
	{operation: "GetTokensFromRefreshToken", handler: (*CognitoService).GetTokensFromRefreshToken},
	{operation: "GetUserAttributeVerificationCode", handler: (*CognitoService).GetUserAttributeVerificationCode, wafInspected: true},
	{operation: "VerifyUserAttribute", handler: (*CognitoService).VerifyUserAttribute, wafInspected: true},
	{operation: "ResendConfirmationCode", handler: (*CognitoService).ResendConfirmationCode, wafInspected: true},
	{operation: "GetUserAuthFactors", handler: (*CognitoService).GetUserAuthFactors},
	{operation: "AdminGetUserAuthFactors", handler: (*CognitoService).AdminGetUserAuthFactors},

	// MFA & user settings
	{operation: "AdminSetUserMFAPreference", handler: (*CognitoService).AdminSetUserMFAPreference},
	{operation: "SetUserMFAPreference", handler: (*CognitoService).SetUserMFAPreference, wafInspected: true},
	{operation: "AdminSetUserSettings", handler: (*CognitoService).AdminSetUserSettings},
	{operation: "SetUserSettings", handler: (*CognitoService).SetUserSettings},

	// Device management
	{operation: "ConfirmDevice", handler: (*CognitoService).ConfirmDevice},
	{operation: "GetDevice", handler: (*CognitoService).GetDevice},
	{operation: "ForgetDevice", handler: (*CognitoService).ForgetDevice},
	{operation: "ListDevices", handler: (*CognitoService).ListDevices},
	{operation: "UpdateDeviceStatus", handler: (*CognitoService).UpdateDeviceStatus},
	{operation: "AdminGetDevice", handler: (*CognitoService).AdminGetDevice},
	{operation: "AdminForgetDevice", handler: (*CognitoService).AdminForgetDevice},
	{operation: "AdminListDevices", handler: (*CognitoService).AdminListDevices},
	{operation: "AdminUpdateDeviceStatus", handler: (*CognitoService).AdminUpdateDeviceStatus},

	// Auth events
	{operation: "AdminListUserAuthEvents", handler: (*CognitoService).AdminListUserAuthEvents},
	{operation: "AdminUpdateAuthEventFeedback", handler: (*CognitoService).AdminUpdateAuthEventFeedback},
	{operation: "UpdateAuthEventFeedback", handler: (*CognitoService).UpdateAuthEventFeedback},

	// Client secrets
	{operation: "AddUserPoolClientSecret", handler: (*CognitoService).AddUserPoolClientSecret},
	{operation: "DeleteUserPoolClientSecret", handler: (*CognitoService).DeleteUserPoolClientSecret},
	{operation: "ListUserPoolClientSecrets", handler: (*CognitoService).ListUserPoolClientSecrets},

	// Log delivery
	{operation: "SetLogDeliveryConfiguration", handler: (*CognitoService).SetLogDeliveryConfiguration},
	{operation: "GetLogDeliveryConfiguration", handler: (*CognitoService).GetLogDeliveryConfiguration},

	// Risk configuration
	{operation: "SetRiskConfiguration", handler: (*CognitoService).SetRiskConfiguration},

	// UI customisation
	{operation: "GetUICustomization", handler: (*CognitoService).GetUICustomization},
	{operation: "SetUICustomization", handler: (*CognitoService).SetUICustomization},

	// Provider user linking
	{operation: "AdminDisableProviderForUser", handler: (*CognitoService).AdminDisableProviderForUser},
	{operation: "AdminLinkProviderForUser", handler: (*CognitoService).AdminLinkProviderForUser},

	// Misc small operations
	{operation: "AddCustomAttributes", handler: (*CognitoService).AddCustomAttributes},
	{operation: "GetIdentityProviderByIdentifier", handler: (*CognitoService).GetIdentityProviderByIdentifier},
	{operation: "GetSigningCertificate", handler: (*CognitoService).GetSigningCertificate},

	// Provisioned limits
	{operation: "GetProvisionedLimit", handler: (*CognitoService).GetProvisionedLimit},
	{operation: "UpdateProvisionedLimit", handler: (*CognitoService).UpdateProvisionedLimit},

	// User import
	{operation: "CreateUserImportJob", handler: (*CognitoService).CreateUserImportJob},
	{operation: "DescribeUserImportJob", handler: (*CognitoService).DescribeUserImportJob},
	{operation: "ListUserImportJobs", handler: (*CognitoService).ListUserImportJobs},
	{operation: "StartUserImportJob", handler: (*CognitoService).StartUserImportJob},
	{operation: "StopUserImportJob", handler: (*CognitoService).StopUserImportJob},

	// WebAuthn
	{operation: "StartWebAuthnRegistration", handler: (*CognitoService).StartWebAuthnRegistration},
	{operation: "CompleteWebAuthnRegistration", handler: (*CognitoService).CompleteWebAuthnRegistration},
	{operation: "ListWebAuthnCredentials", handler: (*CognitoService).ListWebAuthnCredentials},
	{operation: "DeleteWebAuthnCredential", handler: (*CognitoService).DeleteWebAuthnCredential},

	// Managed login branding
	{operation: "CreateManagedLoginBranding", handler: (*CognitoService).CreateManagedLoginBranding},
	{operation: "DescribeManagedLoginBranding", handler: (*CognitoService).DescribeManagedLoginBranding},
	{operation: "DescribeManagedLoginBrandingByClient", handler: (*CognitoService).DescribeManagedLoginBrandingByClient},
	{operation: "UpdateManagedLoginBranding", handler: (*CognitoService).UpdateManagedLoginBranding},
	{operation: "DeleteManagedLoginBranding", handler: (*CognitoService).DeleteManagedLoginBranding},

	// Terms
	{operation: "CreateTerms", handler: (*CognitoService).CreateTerms},
	{operation: "DescribeTerms", handler: (*CognitoService).DescribeTerms},
	{operation: "DescribeTermsByClient", handler: (*CognitoService).DescribeTermsByClient},
	{operation: "ListTerms", handler: (*CognitoService).ListTerms},
	{operation: "UpdateTerms", handler: (*CognitoService).UpdateTerms},
	{operation: "DeleteTerms", handler: (*CognitoService).DeleteTerms},

	// Machine-to-machine authorisation. GetClientToken authorises with the
	// app client's secret rather than AWS credentials, so its requests take
	// the WAF inspection wrapper like the other credential-free operations.
	{operation: "GetClientToken", handler: (*CognitoService).GetClientToken, wafInspected: true},

	// Replicas
	{operation: "CreateUserPoolReplica", handler: (*CognitoService).CreateUserPoolReplica},
	{operation: "ListUserPoolReplicas", handler: (*CognitoService).ListUserPoolReplicas},
	{operation: "DeleteUserPoolReplica", handler: (*CognitoService).DeleteUserPoolReplica},
	{operation: "UpdateUserPoolReplica", handler: (*CognitoService).UpdateUserPoolReplica},
}

// RegisterHandlers registers every user pools API operation with the
// dispatcher from the single registration table, wrapping exactly the
// wafInspected rows with the WebACL enforcement.
func (s *CognitoService) RegisterHandlers(d handler.Registrar) {
	for _, op := range userPoolsAPIOperations {
		h := op.bind(s)
		if op.wafInspected {
			h = s.withWAFEnforcement(op.operation, h)
		}
		d.RegisterHandlerForService(cognitoIDPServiceName, op.operation, h)
	}
}

// buildWafInspectedOperations derives the WAF-inspected operation set from
// the registration table: an operation is inspected exactly when its table
// row is marked wafInspected, so the set always matches the registrations
// that carry the enforcement wrapper.
func buildWafInspectedOperations() map[string]bool {
	inspected := make(map[string]bool)
	for _, op := range userPoolsAPIOperations {
		if op.wafInspected {
			inspected[op.operation] = true
		}
	}
	return inspected
}
