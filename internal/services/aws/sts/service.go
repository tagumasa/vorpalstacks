// Package sts provides STS (Security Token Service) operations for vorpalstacks.
//
// STS is an IAM sub-service and directly accesses the IAM store
// (internal/store/aws/iam) for identity and role resolution.
// This is an intentional architectural decision: STS fundamentally
// depends on IAM roles and access keys, and synchronous store access
// is required for trust policy evaluation and caller identity resolution.
package sts

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	stsstore "vorpalstacks/internal/store/aws/sts"
)

// STSService provides AWS Security Token Service operations.
type STSService struct {
	stores sync.Map // caches STS SessionStore per region
}

// Close stops all background goroutines in cached SessionStores.
func (s *STSService) Close() {
	s.stores.Range(func(_, v any) bool {
		if c, ok := v.(interface{ Close() }); ok {
			c.Close()
		}
		return true
	})
}

// NewSTSService creates a new STS service instance.
func NewSTSService() *STSService {
	return &STSService{}
}

func (s *STSService) store(reqCtx *request.RequestContext) (stsstore.SessionStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, "global", func() (stsstore.SessionStoreInterface, error) {
		storage, err := reqCtx.GetGlobalStorage()
		if err != nil {
			return nil, fmt.Errorf("failed to get global storage: %w", err)
		}
		return stsstore.NewSessionStore(storage, reqCtx.GetRegion()), nil
	})
}

// RegisterHandlers registers all STS operation handlers with the dispatcher.
func (s *STSService) RegisterHandlers(d handler.Registrar) {
	d.RegisterHandlerForService("sts", "GetCallerIdentity", s.GetCallerIdentity)
	d.RegisterHandlerForService("sts", "AssumeRole", s.AssumeRole)
	d.RegisterHandlerForService("sts", "GetSessionToken", s.GetSessionToken)
	d.RegisterHandlerForService("sts", "AssumeRoleWithSAML", s.AssumeRoleWithSAML)
	d.RegisterHandlerForService("sts", "AssumeRoleWithWebIdentity", s.AssumeRoleWithWebIdentity)
	d.RegisterHandlerForService("sts", "AssumeRoot", s.AssumeRoot)
	d.RegisterHandlerForService("sts", "DecodeAuthorizationMessage", s.DecodeAuthorizationMessage)
	d.RegisterHandlerForService("sts", "GetAccessKeyInfo", s.GetAccessKeyInfo)
	d.RegisterHandlerForService("sts", "GetFederationToken", s.GetFederationToken)
	d.RegisterHandlerForService("sts", "GetDelegatedAccessToken", s.GetDelegatedAccessToken)
	d.RegisterHandlerForService("sts", "GetWebIdentityToken", s.GetWebIdentityToken)
}
