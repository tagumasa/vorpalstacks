// Package cognito implements AWS Cognito service handlers for user pools,
// users, groups, and authentication operations.
package cognitoidentityprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"vorpalstacks/internal/common/auth"
	"vorpalstacks/internal/common/iam"
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
	bus                 eventbus.ServiceBus
	stores              sync.Map // region → cognitostore.CognitoStoreInterface
	authCodes           sync.Map // code string → authCodeEntry
	authCodeCleanupOnce sync.Once
	bgCtx               context.Context
	bgCancel            context.CancelFunc
	bgWg                sync.WaitGroup
	// roleProvider supplies IAM role policies for the group-role trust
	// validation performed inside the group Cores.
	roleProvider iam.RolePolicyProvider
	// importCredentials supplies the SigV4 keys used to presign the CSV
	// upload URL handed out by CreateUserImportJob.
	importCredentials auth.CredentialsProvider
	// waf holds the injected WAF request-inspection entry point.
	waf wafEnforcement
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

// SetRoleProvider injects the IAM role policy provider so that the group
// Cores can validate group role trust policies without a request context.
func (s *CognitoService) SetRoleProvider(rp iam.RolePolicyProvider) {
	s.roleProvider = rp
}

// iamValidator builds the IAM validator used by the group Cores. A nil
// result (no role provider injected) leaves role trust validation skipped,
// the same nil handling the scheduler Core applies.
func (s *CognitoService) iamValidator() *iam.IAMValidator {
	if s.roleProvider == nil {
		return nil
	}
	return iam.NewIAMValidator(s.roleProvider, s.accountID)
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
func (s *CognitoService) SetEventBus(bus eventbus.ServiceBus) error {
	s.bus = bus
	if bus != nil {
		if _, err := eventbus.SubscribeTyped[*eventbus.CognitoTriggerEvent](bus, s.handleCognitoTrigger); err != nil {
			return fmt.Errorf("cognito-idp: subscribe CognitoTriggerEvent: %w", err)
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, `{"error":"internal_error"}`, http.StatusInternalServerError)
	}
}

var emptyJWKS = map[string]interface{}{"keys": []interface{}{}}

// JWKSHandler serves the JSON Web Key Set for a Cognito User Pool. AWS
// serves a pool's keys from a pool-specific address, so the pool is part of
// the request, never inferred: a request without a userPoolId is rejected
// rather than answered with an arbitrary pool's keys.
func (s *CognitoService) JWKSHandler(w http.ResponseWriter, r *http.Request) {
	if s.storageManager == nil {
		writeJSON(w, emptyJWKS)
		return
	}
	ctx := context.Background()
	reqCtx := request.NewRequestContext(ctx, s.storageManager, s.accountID, s.region)
	userPoolID := r.URL.Query().Get("userPoolId")
	if userPoolID == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "userPoolId is required"})
		return
	}
	jwks, err := s.getJWKSCore(reqCtx.GetRegion(), userPoolID)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			w.WriteHeader(http.StatusNotFound)
			writeJSON(w, map[string]string{"error": "user pool not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "internal_error"})
		return
	}
	writeJSON(w, jwks)
}

func (s *CognitoService) store(reqCtx *request.RequestContext) (cognitostore.CognitoStoreInterface, error) {
	return s.GetStoreForRegion(reqCtx.GetRegion())
}

// GetStoreForRegion returns the cached Cognito store for the given region,
// creating a new store instance if not already cached. It is the single
// construction path: the request plane arrives through store(reqCtx), the
// background workers and the hosted UI arrive here directly, and the
// GetOrCreateStoreE race handling (loser store closed) covers both. The
// storage manager is consulted only on a cache miss — a pre-seeded store
// (tests, cross-service wiring) resolves without one.
func (s *CognitoService) GetStoreForRegion(region string) (cognitostore.CognitoStoreInterface, error) {
	return storecommon.GetOrCreateStoreE(&s.stores, region, func() (cognitostore.CognitoStoreInterface, error) {
		if s.storageManager == nil {
			return nil, fmt.Errorf("cognito idp storage manager not initialised")
		}
		storage, err := s.storageManager.GetStorage(region)
		if err != nil {
			return nil, fmt.Errorf("failed to get storage: %w", err)
		}
		st := cognitostore.NewCognitoStore(storage, s.accountID, region)
		s.finaliseStaleImportJobs(region, st)
		return st, nil
	})
}

// RegisterHandlers lives in operation_registration.go, driven by the
// single registration table that also derives the WAF-inspected set.
