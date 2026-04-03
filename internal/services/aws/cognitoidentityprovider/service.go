// Package cognito implements AWS Cognito service handlers for user pools,
// users, groups, and authentication operations.
package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"net/http"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/server/eventbus"
	"vorpalstacks/internal/services/aws/common"
	"vorpalstacks/internal/services/aws/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// CognitoService provides operations for AWS Cognito User Pools.
type CognitoService struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	region         string
	bus            eventbus.Bus
	lambdaInvoker  common.LambdaInvoker
}

// NewCognitoService creates a new Cognito User Pools service instance.
func NewCognitoService(store storage.BasicStorage, accountID, region string) *CognitoService {
	return &CognitoService{
		accountID: accountID,
		region:    region,
	}
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
		if eb, ok := bus.(*eventbus.EventBus); ok {
			_, _ = eventbus.SubscribeTyped[*eventbus.CognitoTriggerEvent](eb, s.handleCognitoTrigger,
				eventbus.WithCallerPrincipal("cognito-idp.amazonaws.com"),
			)
		}
	}
}

// SetLambdaInjector sets the Lambda invoker used for Cognito trigger
// Lambda functions. This must be called before any trigger can fire.
func (s *CognitoService) SetLambdaInvoker(invoker common.LambdaInvoker) {
	s.lambdaInvoker = invoker
}

// JWKSHandler serves the JSON Web Key Set for a Cognito User Pool.
// If no userPoolId query parameter is provided, the first available pool is used.
func (s *CognitoService) JWKSHandler(w http.ResponseWriter, r *http.Request) {
	if s.storageManager == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"keys": []interface{}{}})
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"keys": []interface{}{}})
		return
	}
	jwks, err := s.GetJWKS(reqCtx, userPoolID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"keys": []interface{}{}})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}

// store returns the Cognito store for the given request context.
func (s *CognitoService) store(reqCtx *request.RequestContext) (cognitostore.CognitoStoreInterface, error) {
	if store := reqCtx.GetCognitoStore(); store != nil {
		return store, nil
	}
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	return cognitostore.NewCognitoStore(storage, reqCtx.GetAccountID(), reqCtx.GetRegion()), nil
}

// RegisterHandlers registers the Cognito handlers with the dispatcher.
func (s *CognitoService) RegisterHandlers(d *dispatcher.Dispatcher) {
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
}
