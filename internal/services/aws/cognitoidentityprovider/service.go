// Package cognito implements AWS Cognito service handlers for user pools,
// users, groups, and authentication operations.
package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// CognitoService provides operations for AWS Cognito User Pools.
type CognitoService struct {
	storageManager      *storage.RegionStorageManager
	accountID           string
	region              string
	bus                 eventbus.Bus
	stores              sync.Map // region → cognitostore.CognitoStoreInterface
	authCodes           sync.Map // code string → authCodeEntry
	authCodeCleanupOnce sync.Once
	bgCtx               context.Context
	bgCancel            context.CancelFunc
	bgWg                sync.WaitGroup
	// importCredentials supplies the SigV4 keys used to presign the CSV
	// upload URL handed out by CreateUserImportJob.
	importCredentials auth.CredentialsProvider
}

// NewCognitoService creates a new Cognito User Pools service instance.
func NewCognitoService(accountID, region string) *CognitoService {
	ctx, cancel := context.WithCancel(context.Background())
	return &CognitoService{
		accountID: accountID,
		region:    region,
		bgCtx:     ctx,
		bgCancel:  cancel,
	}
}

// SetImportCredentialsProvider sets the credentials used to presign the
// user import CSV upload URL.
func (s *CognitoService) SetImportCredentialsProvider(provider auth.CredentialsProvider) {
	s.importCredentials = provider
}

// Close stops background workers (user import jobs) and waits for them to
// finish.
func (s *CognitoService) Close() {
	if s.bgCancel != nil {
		s.bgCancel()
	}
	s.bgWg.Wait()
}

// SetStorageManager injects the storage manager, required for the JWKS handler.
func (s *CognitoService) SetStorageManager(sm *storage.RegionStorageManager) {
	s.storageManager = sm
}

// SetEventBus registers the Cognito trigger handler on the event bus.
// The handler invokes the Lambda function specified in the trigger event
// and returns the Lambda response payload.
func (s *CognitoService) SetEventBus(bus eventbus.Bus) {
	s.bus = bus
	if bus != nil {
		_, _ = eventbus.SubscribeTyped[*eventbus.CognitoTriggerEvent](bus, s.handleCognitoTrigger,
			eventbus.WithCallerPrincipal("cognito-idp.amazonaws.com"),
		)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, `{"error":"internal_error"}`, http.StatusInternalServerError)
	}
}

var emptyJWKS = map[string]interface{}{"keys": []interface{}{}}

// JWKSHandler serves the JSON Web Key Set for a Cognito User Pool.
// If no userPoolId query parameter is provided, the first available pool is used.
func (s *CognitoService) JWKSHandler(w http.ResponseWriter, r *http.Request) {
	if s.storageManager == nil {
		writeJSON(w, emptyJWKS)
		return
	}
	ctx := context.Background()
	reqCtx := request.NewRequestContext(ctx, s.storageManager, s.accountID, s.region)
	userPoolID := r.URL.Query().Get("userPoolId")
	if userPoolID == "" {
		pools, _ := s.ListUserPoolsRaw(reqCtx)
		if len(pools) > 0 {
			userPoolID = pools[0].ID
		}
	}
	if userPoolID == "" {
		writeJSON(w, emptyJWKS)
		return
	}
	jwks, err := s.GetJWKS(reqCtx, userPoolID)
	if err != nil {
		writeJSON(w, emptyJWKS)
		return
	}
	writeJSON(w, jwks)
}

func (s *CognitoService) store(reqCtx *request.RequestContext) (cognitostore.CognitoStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, reqCtx.GetRegion(), func() (cognitostore.CognitoStoreInterface, error) {
		storage, err := reqCtx.GetStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		return cognitostore.NewCognitoStore(storage, s.accountID, reqCtx.GetRegion()), nil
	})
}

// GetStoreForRegion returns the cached Cognito store for the given region,
// creating a new store instance if not already cached.
func (s *CognitoService) GetStoreForRegion(region string) (cognitostore.CognitoStoreInterface, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(cognitostore.CognitoStoreInterface), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("cognito idp storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	store := cognitostore.NewCognitoStore(st, s.accountID, region)
	actual, _ := s.stores.LoadOrStore(region, store)
	return actual.(cognitostore.CognitoStoreInterface), nil
}

// RegisterHandlers registers the Cognito handlers with the dispatcher.
func (s *CognitoService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("cognito-idp", "CreateUserPool", s.CreateUserPool)
	d.RegisterHandlerForService("cognito-idp", "DeleteUserPool", s.DeleteUserPool)
	d.RegisterHandlerForService("cognito-idp", "DescribeUserPool", s.DescribeUserPool)
	d.RegisterHandlerForService("cognito-idp", "ListUserPools", s.ListUserPools)
	d.RegisterHandlerForService("cognito-idp", "UpdateUserPool", s.UpdateUserPool)

	d.RegisterHandlerForService("cognito-idp", "CreateUserPoolClient", s.CreateUserPoolClient)
	d.RegisterHandlerForService("cognito-idp", "DeleteUserPoolClient", s.DeleteUserPoolClient)
	d.RegisterHandlerForService("cognito-idp", "DescribeUserPoolClient", s.DescribeUserPoolClient)
	d.RegisterHandlerForService("cognito-idp", "ListUserPoolClients", s.ListUserPoolClients)
	d.RegisterHandlerForService("cognito-idp", "UpdateUserPoolClient", s.UpdateUserPoolClient)

	d.RegisterHandlerForService("cognito-idp", "AdminCreateUser", s.AdminCreateUser)
	d.RegisterHandlerForService("cognito-idp", "AdminDeleteUser", s.AdminDeleteUser)
	d.RegisterHandlerForService("cognito-idp", "AdminDeleteUserAttributes", s.AdminDeleteUserAttributes)
	d.RegisterHandlerForService("cognito-idp", "AdminDisableUser", s.AdminDisableUser)
	d.RegisterHandlerForService("cognito-idp", "AdminEnableUser", s.AdminEnableUser)
	d.RegisterHandlerForService("cognito-idp", "AdminGetUser", s.AdminGetUser)
	d.RegisterHandlerForService("cognito-idp", "AdminResetUserPassword", s.AdminResetUserPassword)
	d.RegisterHandlerForService("cognito-idp", "AdminSetUserPassword", s.AdminSetUserPassword)
	d.RegisterHandlerForService("cognito-idp", "AdminUpdateUserAttributes", s.AdminUpdateUserAttributes)
	d.RegisterHandlerForService("cognito-idp", "AdminUserGlobalSignOut", s.AdminUserGlobalSignOut)
	d.RegisterHandlerForService("cognito-idp", "DeleteUser", s.DeleteUser)
	d.RegisterHandlerForService("cognito-idp", "DeleteUserAttributes", s.DeleteUserAttributes)
	d.RegisterHandlerForService("cognito-idp", "GetUser", s.GetUser)
	d.RegisterHandlerForService("cognito-idp", "ListUsers", s.ListUsers)
	d.RegisterHandlerForService("cognito-idp", "UpdateUserAttributes", s.UpdateUserAttributes)

	d.RegisterHandlerForService("cognito-idp", "CreateGroup", s.CreateGroup)
	d.RegisterHandlerForService("cognito-idp", "DeleteGroup", s.DeleteGroup)
	d.RegisterHandlerForService("cognito-idp", "GetGroup", s.GetGroup)
	d.RegisterHandlerForService("cognito-idp", "ListGroups", s.ListGroups)
	d.RegisterHandlerForService("cognito-idp", "UpdateGroup", s.UpdateGroup)
	d.RegisterHandlerForService("cognito-idp", "AdminAddUserToGroup", s.AdminAddUserToGroup)
	d.RegisterHandlerForService("cognito-idp", "AdminRemoveUserFromGroup", s.AdminRemoveUserFromGroup)
	d.RegisterHandlerForService("cognito-idp", "ListUsersInGroup", s.ListUsersInGroup)
	d.RegisterHandlerForService("cognito-idp", "AdminListGroupsForUser", s.AdminListGroupsForUser)

	d.RegisterHandlerForService("cognito-idp", "SignUp", s.SignUp)
	d.RegisterHandlerForService("cognito-idp", "ConfirmSignUp", s.ConfirmSignUp)
	d.RegisterHandlerForService("cognito-idp", "AdminConfirmSignUp", s.AdminConfirmSignUp)
	d.RegisterHandlerForService("cognito-idp", "InitiateAuth", s.InitiateAuth)
	d.RegisterHandlerForService("cognito-idp", "AdminInitiateAuth", s.AdminInitiateAuth)
	d.RegisterHandlerForService("cognito-idp", "RespondToAuthChallenge", s.RespondToAuthChallenge)
	d.RegisterHandlerForService("cognito-idp", "AdminRespondToAuthChallenge", s.AdminRespondToAuthChallenge)
	d.RegisterHandlerForService("cognito-idp", "SignOut", s.SignOut)
	d.RegisterHandlerForService("cognito-idp", "GlobalSignOut", s.GlobalSignOut)
	d.RegisterHandlerForService("cognito-idp", "ChangePassword", s.ChangePassword)
	d.RegisterHandlerForService("cognito-idp", "ForgotPassword", s.ForgotPassword)
	d.RegisterHandlerForService("cognito-idp", "ConfirmForgotPassword", s.ConfirmForgotPassword)

	d.RegisterHandlerForService("cognito-idp", "TagResource", s.TagResource)
	d.RegisterHandlerForService("cognito-idp", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("cognito-idp", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("cognito-idp", "GetUserPoolMfaConfig", s.GetUserPoolMfaConfig)
	d.RegisterHandlerForService("cognito-idp", "SetUserPoolMfaConfig", s.SetUserPoolMfaConfig)
	d.RegisterHandlerForService("cognito-idp", "AssociateSoftwareToken", s.AssociateSoftwareToken)
	d.RegisterHandlerForService("cognito-idp", "VerifySoftwareToken", s.VerifySoftwareToken)

	d.RegisterHandlerForService("cognito-idp", "CreateUserPoolDomain", s.CreateUserPoolDomain)
	d.RegisterHandlerForService("cognito-idp", "DescribeUserPoolDomain", s.DescribeUserPoolDomain)
	d.RegisterHandlerForService("cognito-idp", "DeleteUserPoolDomain", s.DeleteUserPoolDomain)
	d.RegisterHandlerForService("cognito-idp", "UpdateUserPoolDomain", s.UpdateUserPoolDomain)
	d.RegisterHandlerForService("cognito-idp", "CreateResourceServer", s.CreateResourceServer)
	d.RegisterHandlerForService("cognito-idp", "DescribeResourceServer", s.DescribeResourceServer)
	d.RegisterHandlerForService("cognito-idp", "UpdateResourceServer", s.UpdateResourceServer)
	d.RegisterHandlerForService("cognito-idp", "DeleteResourceServer", s.DeleteResourceServer)
	d.RegisterHandlerForService("cognito-idp", "ListResourceServers", s.ListResourceServers)
	d.RegisterHandlerForService("cognito-idp", "CreateIdentityProvider", s.CreateIdentityProvider)
	d.RegisterHandlerForService("cognito-idp", "DescribeIdentityProvider", s.DescribeIdentityProvider)
	d.RegisterHandlerForService("cognito-idp", "UpdateIdentityProvider", s.UpdateIdentityProvider)
	d.RegisterHandlerForService("cognito-idp", "DeleteIdentityProvider", s.DeleteIdentityProvider)
	d.RegisterHandlerForService("cognito-idp", "ListIdentityProviders", s.ListIdentityProviders)
	d.RegisterHandlerForService("cognito-idp", "GetCSVHeader", s.GetCSVHeader)
	d.RegisterHandlerForService("cognito-idp", "DescribeRiskConfiguration", s.DescribeRiskConfiguration)

	// Token & auth operations
	d.RegisterHandlerForService("cognito-idp", "RevokeToken", s.RevokeToken)
	d.RegisterHandlerForService("cognito-idp", "GetTokensFromRefreshToken", s.GetTokensFromRefreshToken)
	d.RegisterHandlerForService("cognito-idp", "GetUserAttributeVerificationCode", s.GetUserAttributeVerificationCode)
	d.RegisterHandlerForService("cognito-idp", "VerifyUserAttribute", s.VerifyUserAttribute)
	d.RegisterHandlerForService("cognito-idp", "ResendConfirmationCode", s.ResendConfirmationCode)
	d.RegisterHandlerForService("cognito-idp", "GetUserAuthFactors", s.GetUserAuthFactors)
	d.RegisterHandlerForService("cognito-idp", "AdminGetUserAuthFactors", s.AdminGetUserAuthFactors)

	// MFA & user settings
	d.RegisterHandlerForService("cognito-idp", "AdminSetUserMFAPreference", s.AdminSetUserMFAPreference)
	d.RegisterHandlerForService("cognito-idp", "SetUserMFAPreference", s.SetUserMFAPreference)
	d.RegisterHandlerForService("cognito-idp", "AdminSetUserSettings", s.AdminSetUserSettings)
	d.RegisterHandlerForService("cognito-idp", "SetUserSettings", s.SetUserSettings)

	// Device management
	d.RegisterHandlerForService("cognito-idp", "ConfirmDevice", s.ConfirmDevice)
	d.RegisterHandlerForService("cognito-idp", "GetDevice", s.GetDevice)
	d.RegisterHandlerForService("cognito-idp", "ForgetDevice", s.ForgetDevice)
	d.RegisterHandlerForService("cognito-idp", "ListDevices", s.ListDevices)
	d.RegisterHandlerForService("cognito-idp", "UpdateDeviceStatus", s.UpdateDeviceStatus)
	d.RegisterHandlerForService("cognito-idp", "AdminGetDevice", s.AdminGetDevice)
	d.RegisterHandlerForService("cognito-idp", "AdminForgetDevice", s.AdminForgetDevice)
	d.RegisterHandlerForService("cognito-idp", "AdminListDevices", s.AdminListDevices)
	d.RegisterHandlerForService("cognito-idp", "AdminUpdateDeviceStatus", s.AdminUpdateDeviceStatus)

	// Auth events
	d.RegisterHandlerForService("cognito-idp", "AdminListUserAuthEvents", s.AdminListUserAuthEvents)
	d.RegisterHandlerForService("cognito-idp", "AdminUpdateAuthEventFeedback", s.AdminUpdateAuthEventFeedback)
	d.RegisterHandlerForService("cognito-idp", "UpdateAuthEventFeedback", s.UpdateAuthEventFeedback)

	// Client secrets
	d.RegisterHandlerForService("cognito-idp", "AddUserPoolClientSecret", s.AddUserPoolClientSecret)
	d.RegisterHandlerForService("cognito-idp", "DeleteUserPoolClientSecret", s.DeleteUserPoolClientSecret)
	d.RegisterHandlerForService("cognito-idp", "ListUserPoolClientSecrets", s.ListUserPoolClientSecrets)

	// Log delivery
	d.RegisterHandlerForService("cognito-idp", "SetLogDeliveryConfiguration", s.SetLogDeliveryConfiguration)
	d.RegisterHandlerForService("cognito-idp", "GetLogDeliveryConfiguration", s.GetLogDeliveryConfiguration)

	// Risk configuration
	d.RegisterHandlerForService("cognito-idp", "SetRiskConfiguration", s.SetRiskConfiguration)

	// UI customisation
	d.RegisterHandlerForService("cognito-idp", "GetUICustomization", s.GetUICustomization)
	d.RegisterHandlerForService("cognito-idp", "SetUICustomization", s.SetUICustomization)

	// Provider user linking
	d.RegisterHandlerForService("cognito-idp", "AdminDisableProviderForUser", s.AdminDisableProviderForUser)
	d.RegisterHandlerForService("cognito-idp", "AdminLinkProviderForUser", s.AdminLinkProviderForUser)

	// Misc small operations
	d.RegisterHandlerForService("cognito-idp", "AddCustomAttributes", s.AddCustomAttributes)
	d.RegisterHandlerForService("cognito-idp", "GetIdentityProviderByIdentifier", s.GetIdentityProviderByIdentifier)
	d.RegisterHandlerForService("cognito-idp", "GetSigningCertificate", s.GetSigningCertificate)

	// Provisioned limits
	d.RegisterHandlerForService("cognito-idp", "GetProvisionedLimit", s.GetProvisionedLimit)
	d.RegisterHandlerForService("cognito-idp", "UpdateProvisionedLimit", s.UpdateProvisionedLimit)

	// User import
	d.RegisterHandlerForService("cognito-idp", "CreateUserImportJob", s.CreateUserImportJob)
	d.RegisterHandlerForService("cognito-idp", "DescribeUserImportJob", s.DescribeUserImportJob)
	d.RegisterHandlerForService("cognito-idp", "ListUserImportJobs", s.ListUserImportJobs)
	d.RegisterHandlerForService("cognito-idp", "StartUserImportJob", s.StartUserImportJob)
	d.RegisterHandlerForService("cognito-idp", "StopUserImportJob", s.StopUserImportJob)

	// WebAuthn
	d.RegisterHandlerForService("cognito-idp", "StartWebAuthnRegistration", s.StartWebAuthnRegistration)
	d.RegisterHandlerForService("cognito-idp", "CompleteWebAuthnRegistration", s.CompleteWebAuthnRegistration)
	d.RegisterHandlerForService("cognito-idp", "ListWebAuthnCredentials", s.ListWebAuthnCredentials)
	d.RegisterHandlerForService("cognito-idp", "DeleteWebAuthnCredential", s.DeleteWebAuthnCredential)

	// Managed login branding
	d.RegisterHandlerForService("cognito-idp", "CreateManagedLoginBranding", s.CreateManagedLoginBranding)
	d.RegisterHandlerForService("cognito-idp", "DescribeManagedLoginBranding", s.DescribeManagedLoginBranding)
	d.RegisterHandlerForService("cognito-idp", "DescribeManagedLoginBrandingByClient", s.DescribeManagedLoginBrandingByClient)
	d.RegisterHandlerForService("cognito-idp", "UpdateManagedLoginBranding", s.UpdateManagedLoginBranding)
	d.RegisterHandlerForService("cognito-idp", "DeleteManagedLoginBranding", s.DeleteManagedLoginBranding)

	// Terms
	d.RegisterHandlerForService("cognito-idp", "CreateTerms", s.CreateTerms)
	d.RegisterHandlerForService("cognito-idp", "DescribeTerms", s.DescribeTerms)
	d.RegisterHandlerForService("cognito-idp", "ListTerms", s.ListTerms)
	d.RegisterHandlerForService("cognito-idp", "UpdateTerms", s.UpdateTerms)
	d.RegisterHandlerForService("cognito-idp", "DeleteTerms", s.DeleteTerms)

	// Replicas
	d.RegisterHandlerForService("cognito-idp", "CreateUserPoolReplica", s.CreateUserPoolReplica)
	d.RegisterHandlerForService("cognito-idp", "ListUserPoolReplicas", s.ListUserPoolReplicas)
	d.RegisterHandlerForService("cognito-idp", "DeleteUserPoolReplica", s.DeleteUserPoolReplica)
	d.RegisterHandlerForService("cognito-idp", "UpdateUserPoolReplica", s.UpdateUserPoolReplica)
}
