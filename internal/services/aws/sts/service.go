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
	"os"
	"sync"
	"time"

	"vorpalstacks/internal/common/handler"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	iamstore "vorpalstacks/internal/store/aws/iam"
	stsstore "vorpalstacks/internal/store/aws/sts"
	arnutil "vorpalstacks/internal/utils/aws/arn"
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

func (s *STSService) iamStore(reqCtx *request.RequestContext) (iamstore.IAMStoreInterface, error) {
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage: %w", err)
	}
	return iamstore.GetOrCreateGlobalStore(storage, reqCtx.GetAccountID()), nil
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

// resolveCallerSession looks up the caller's STS session when the request
// includes a security token (temporary credentials). Returns nil for
// permanent IAM credentials or when no session is found.
func (s *STSService) resolveCallerSession(reqCtx *request.RequestContext, req *request.ParsedRequest) *stsstore.Session {
	securityToken := req.Headers.Get("X-Amz-Security-Token")
	if securityToken == "" {
		return nil
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil
	}
	session, err := store.Get(securityToken)
	if err != nil {
		return nil
	}
	return session
}

// assumeContext bundles optional AssumeRole-family parameters that influence
// trust policy evaluation (ExternalId, MFA). nil is equivalent to an empty
// context and is used by callers that provide none of these parameters.
type assumeContext struct {
	ExternalId   string
	SerialNumber string
	TokenCode    string
	CallerName   string // IAM user name of the caller, for MFA device ownership check
	MFAPresent   bool   // set to true by resolveRoleForAssume after successful MFA verification
	// ProvidedContexts carries the parsed ProvidedContexts.member.N entries
	// from the AssumeRole request. resolveRoleForAssume injects each entry's
	// ProviderArn / ContextAssertion into evalCtx.Variables so that trust
	// policy conditions referencing sts:ProvidedContext* can match.
	ProvidedContexts []ProvidedContextEntry
}

// ProvidedContextEntry is the in-memory representation of one
// ProvidedContext from the AssumeRole request. It mirrors the Smithy
// ProvidedContext shape (ProviderArn + ContextAssertion) without coupling
// the service layer to the Smithy-generated package.
type ProvidedContextEntry struct {
	ProviderArn      string
	ContextAssertion string
}

// resolveRoleForAssume resolves and validates a role for STS Assume operations.
// It fetches the role by name, parses its trust policy, and evaluates whether
// the given principal is allowed to perform the specified action.
//
// actx carries optional caller-provided parameters (ExternalId, MFA). When
// non-nil, the parameters are validated and injected into the evaluation
// context so trust policies referencing sts:ExternalId,
// aws:MultiFactorAuthPresent, or aws:MultiFactorAuthAge match correctly.
// SAML and WebIdentity flows pass nil as the AWS API exposes no ExternalId or
// MFA parameters on those operations.
func (s *STSService) resolveRoleForAssume(reqCtx *request.RequestContext, roleArn, principalArn, action string, actx *assumeContext) (*iamstore.Role, error) {
	roleName := arnutil.ExtractRoleNameFromARN(roleArn)
	if roleName == "" {
		return nil, ErrInvalidRoleArn
	}

	iamStore, err := s.iamStore(reqCtx)
	if err != nil {
		return nil, err
	}

	role, err := iamStore.Roles().Get(roleName)
	if err != nil {
		return nil, ErrNoSuchRole
	}

	trustPolicyDoc, err := iamStore.Roles().GetAssumeRolePolicyDocument(roleName)
	if err != nil {
		return nil, ErrNoSuchRole
	}

	parsedPolicy, err := policy.ParseDocument(trustPolicyDoc)
	if err != nil {
		return nil, ErrInvalidRoleArn
	}

	evalCtx := iam.BuildEvaluationContext(reqCtx.GetAccountID(), principalArn)
	if actx != nil {
		if actx.ExternalId != "" {
			if evalCtx.Variables == nil {
				evalCtx.Variables = make(map[string]string)
			}
			evalCtx.Variables["sts:ExternalId"] = actx.ExternalId
		}
		if actx.SerialNumber != "" && actx.TokenCode != "" {
			mfaValid, err := s.verifyCallerMFA(reqCtx, iamStore, actx)
			if err != nil {
				return nil, err
			}
			if mfaValid {
				actx.MFAPresent = true
				evalCtx.MultiFactorAuthPresent = true
				if evalCtx.Variables == nil {
					evalCtx.Variables = make(map[string]string)
				}
				evalCtx.Variables["aws:MultiFactorAuthAge"] = "0"
				evalCtx.TokenIssueTime = time.Now().UTC()
			}
		}
		// Surface caller-supplied ProvidedContexts as
		// sts:ProvidedContextProviderArn.N and
		// sts:ProvidedContextAssertion.N variables on the trust-policy
		// evaluation context so policy templates that reference the
		// EC2 IMDS context-provider feature can match.
		for i, pc := range actx.ProvidedContexts {
			if pc.ProviderArn != "" || pc.ContextAssertion != "" {
				if evalCtx.Variables == nil {
					evalCtx.Variables = make(map[string]string)
				}
				keyProvider := fmt.Sprintf("sts:ProvidedContextProviderArn.%d", i+1)
				keyAssertion := fmt.Sprintf("sts:ProvidedContextAssertion.%d", i+1)
				evalCtx.Variables[keyProvider] = pc.ProviderArn
				evalCtx.Variables[keyAssertion] = pc.ContextAssertion
			}
		}
	}
	if err := iam.EvaluateTrustPolicyForAction(parsedPolicy, principalArn, action, evalCtx); err != nil {
		return nil, ErrAccessDenied
	}

	return role, nil
}

// verifyCallerMFA looks up the MFA device by SerialNumber, confirms it is
// assigned to the calling user, and validates the TokenCode via TOTP.
// Returns (true, nil) when the MFA device exists and the code is valid.
func (s *STSService) verifyCallerMFA(reqCtx *request.RequestContext, iamStore iamstore.IAMStoreInterface, actx *assumeContext) (bool, error) {
	device, err := iamStore.MFADevices().Get(actx.SerialNumber)
	if err != nil || device == nil {
		return false, ErrMFADeviceNotFound
	}
	if device.UserAssignment == nil || device.UserAssignment.UserName != actx.CallerName {
		return false, ErrMFADeviceNotFound
	}
	if device.Base32StringSeed == "" {
		// A device without a TOTP seed cannot be verified server-side.
		// VorpalStacks does not implement U2F/FIDO challenge-response, so
		// we reject the code rather than silently accepting any value.
		return false, ErrInvalidMFACode
	}
	if !verifyTOTP(device.Base32StringSeed, actx.TokenCode, time.Now().UTC()) {
		return false, ErrInvalidMFACode
	}
	return true, nil
}

func (s *STSService) resolveCallerIdentity(reqCtx *request.RequestContext, req *request.ParsedRequest) (arn, name string) {
	accessKeyId := req.Headers.Get("X-Amz-Access-Key")
	if accessKeyId == "" {
		authHeader := req.Headers.Get("Authorization")
		if authHeader != "" {
			accessKeyId = request.ExtractAccessKeyIDFromAuth(authHeader)
		}
	}
	if accessKeyId == "" {
		return "", ""
	}
	// The signature middleware's static verifier accepts exactly one
	// long-term credential: the configured root key. A request carrying
	// that key is the root principal by construction, independent of
	// whether the key is also registered in the IAM access-key store.
	if rootKey := os.Getenv("AWS_ACCESS_KEY_ID"); rootKey != "" && accessKeyId == rootKey {
		return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root(), iam.RootUserName
	}
	iamStore, err := s.iamStore(reqCtx)
	if err == nil {
		accessKey, err := iamStore.AccessKeys().Get(accessKeyId)
		if err == nil && accessKey != nil {
			// Root user access keys use the special RootUserName constant.
			if accessKey.UserName == iam.RootUserName {
				return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root(), iam.RootUserName
			}
			return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().User(accessKey.UserName), accessKey.UserName
		}
	}
	// Fall back to STS sessions (temporary credentials, e.g. role chaining).
	// When the caller uses a session token, the access key ID is temporary
	// and won't be found in IAM permanent keys.
	sessionStore, err := s.store(reqCtx)
	if err == nil {
		session, err := sessionStore.GetByAccessKeyId(accessKeyId)
		if err == nil && session != nil {
			// For assumed-role sessions the caller's ARN is the
			// assumed-role ARN (arn:aws:sts::account:assumed-role/
			// RoleName/SessionName), NOT the role ARN stored in
			// session.PrincipalArn.  Returning the role ARN here would
			// cause trust-policy condition keys (aws:PrincipalArn) to
			// mismatch in chained AssumeRole calls.  The stored
			// PrincipalArn is kept as the role ARN for the authorization
			// layer (authorizer.go) which expects that format.
			switch session.PrincipalType {
			case "AssumedRole", "SAML", "WebIdentity":
				roleName := arnutil.ExtractRoleNameFromARN(session.RoleArn)
				arn := arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(roleName, session.RoleSessionName)
				return arn, session.PrincipalName
			default:
				return session.PrincipalArn, session.PrincipalName
			}
		}
	}
	return "", ""
}

// resolveCallerArnOrReject resolves the caller's ARN and name, rejecting the
// request when the caller cannot be identified. Session-issuing operations
// must never default an unresolvable caller to the account root ARN: a
// session minted with a root principal ARN bypasses policy evaluation in the
// authorizer, so an identification failure has to fail closed. The single
// exception is TEST_MODE, which the documented SDK-test procedure runs with
// signature verification disabled and dummy credentials no store knows;
// there it mirrors GetCallerIdentity's bootstrap fallback and treats the
// caller as root. Production never takes that branch: with verification
// enabled the signature middleware rejects unknown keys before any handler
// runs.
func (s *STSService) resolveCallerArnOrReject(reqCtx *request.RequestContext, req *request.ParsedRequest) (string, string, error) {
	arn, name := s.resolveCallerIdentity(reqCtx, req)
	if arn == "" {
		if os.Getenv("TEST_MODE") == "true" {
			return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root(), iam.RootUserName, nil
		}
		return "", "", ErrInvalidClientTokenId
	}
	return arn, name, nil
}
