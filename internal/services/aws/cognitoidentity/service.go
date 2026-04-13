// Package cognitoidentity provides AWS Cognito Identity service operations for vorpalstacks.
package cognitoidentity

import (
	"sync"

	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
)

// CognitoIdentityService provides Cognito Identity service operations.
type CognitoIdentityService struct {
	accountID string
	region    string
	stores    sync.Map // region → cognitoidentitystore.CognitoIdentityStoreInterface
}

// NewCognitoIdentityService creates a new Cognito Identity service.
func NewCognitoIdentityService(accountID, region string) *CognitoIdentityService {
	return &CognitoIdentityService{
		accountID: accountID,
		region:    region,
	}
}

// store returns the Cognito Identity store for the given request context.
func (s *CognitoIdentityService) store(reqCtx *request.RequestContext) (cognitoidentitystore.CognitoIdentityStoreInterface, error) {
	if store := reqCtx.GetCognitoIdentityStore(); store != nil {
		return store, nil
	}
	region := reqCtx.GetRegion()
	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(cognitoidentitystore.CognitoIdentityStoreInterface); ok {
			return typed, nil
		}
	}
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	store := cognitoidentitystore.NewCognitoIdentityStore(storage, s.accountID, region)
	if actual, loaded := s.stores.LoadOrStore(region, store); loaded {
		if typed, ok := actual.(cognitoidentitystore.CognitoIdentityStoreInterface); ok {
			return typed, nil
		}
	}
	return store, nil
}

// RegisterHandlers registers the Cognito Identity handlers with the dispatcher.
func (s *CognitoIdentityService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("cognito-identity", "CreateIdentityPool", s.CreateIdentityPool)
	d.RegisterHandlerForService("cognito-identity", "DeleteIdentityPool", s.DeleteIdentityPool)
	d.RegisterHandlerForService("cognito-identity", "DescribeIdentityPool", s.DescribeIdentityPool)
	d.RegisterHandlerForService("cognito-identity", "ListIdentityPools", s.ListIdentityPools)
	d.RegisterHandlerForService("cognito-identity", "UpdateIdentityPool", s.UpdateIdentityPool)
	d.RegisterHandlerForService("cognito-identity", "GetIdentityPoolRoles", s.GetIdentityPoolRoles)
	d.RegisterHandlerForService("cognito-identity", "SetIdentityPoolRoles", s.SetIdentityPoolRoles)
	d.RegisterHandlerForService("cognito-identity", "GetId", s.GetId)
	d.RegisterHandlerForService("cognito-identity", "GetCredentialsForIdentity", s.GetCredentialsForIdentity)
	d.RegisterHandlerForService("cognito-identity", "DescribeIdentity", s.DescribeIdentity)
	d.RegisterHandlerForService("cognito-identity", "TagResource", s.TagResource)
	d.RegisterHandlerForService("cognito-identity", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("cognito-identity", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("cognito-identity", "DeleteIdentities", s.DeleteIdentities)
	d.RegisterHandlerForService("cognito-identity", "ListIdentities", s.ListIdentities)
	d.RegisterHandlerForService("cognito-identity", "GetOpenIdToken", s.GetOpenIdToken)
	d.RegisterHandlerForService("cognito-identity", "GetOpenIdTokenForDeveloperIdentity", s.GetOpenIdTokenForDeveloperIdentity)
	d.RegisterHandlerForService("cognito-identity", "GetPrincipalTagAttributeMap", s.GetPrincipalTagAttributeMap)
	d.RegisterHandlerForService("cognito-identity", "SetPrincipalTagAttributeMap", s.SetPrincipalTagAttributeMap)
	d.RegisterHandlerForService("cognito-identity", "LookupDeveloperIdentity", s.LookupDeveloperIdentity)
	d.RegisterHandlerForService("cognito-identity", "MergeDeveloperIdentities", s.MergeDeveloperIdentities)
	d.RegisterHandlerForService("cognito-identity", "UnlinkDeveloperIdentity", s.UnlinkDeveloperIdentity)
	d.RegisterHandlerForService("cognito-identity", "UnlinkIdentity", s.UnlinkIdentity)
}
