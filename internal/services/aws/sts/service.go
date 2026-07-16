// Package sts provides STS (Security Token Service) operations for vorpalstacks.
//
// STS is an IAM sub-service and directly accesses the IAM store
// (internal/store/aws/iam) for identity and role resolution.
// This is an intentional architectural decision: STS fundamentally
// depends on IAM roles and access keys, and synchronous store access
// is required for trust policy evaluation and caller identity resolution.
package sts

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
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
	"vorpalstacks/internal/utils/timeutils"
)

// externalIdPattern mirrors the Smithy externalIdType trait: [\w+=,.@:\/-]*.
// Combined with the length check (2-1224) it enforces the AWS AssumeRole
// ExternalId parameter constraints.
var externalIdPattern = regexp.MustCompile(`^[\w+=,.@:\/-]+$`)

// sessionTagKeyPattern mirrors the Smithy tagKeyType pattern.
var sessionTagKeyPattern = regexp.MustCompile(`^[\p{L}\p{Z}\p{N}_.:/=+\-@]+$`)

// sessionTagValuePattern mirrors the Smithy tagValueType pattern.
var sessionTagValuePattern = regexp.MustCompile(`^[\p{L}\p{Z}\p{N}_.:/=+\-@]*$`)

// extractSessionTags parses Tags.member.N.Key / Tags.member.N.Value pairs from
// the flat query-parameter map and validates them against the Smithy
// tagKeyType (1-128 chars) and tagValueType (0-256 chars) traits. At most 50
// tags are accepted per the tagListType length constraint.
func extractSessionTags(params map[string]interface{}) (map[string]string, error) {
	tags := make(map[string]string)
	for i := 1; ; i++ {
		key := request.GetStringParam(params, fmt.Sprintf("Tags.member.%d.Key", i))
		if key == "" {
			break
		}
		if i > 50 {
			return nil, ErrTooManySessionTags
		}
		if len(key) < 1 || len(key) > 128 || !sessionTagKeyPattern.MatchString(key) {
			return nil, ErrInvalidSessionTag
		}
		value := request.GetStringParam(params, fmt.Sprintf("Tags.member.%d.Value", i))
		if len(value) > 256 || !sessionTagValuePattern.MatchString(value) {
			return nil, ErrInvalidSessionTag
		}
		if _, exists := tags[key]; exists {
			return nil, ErrDuplicateSessionTagKey
		}
		tags[key] = value
	}
	if len(tags) == 0 {
		return nil, nil
	}
	return tags, nil
}

// extractTransitiveTagKeys parses TransitiveTagKeys.member.N from the flat
// query-parameter map. Each key is validated against the Smithy tagKeyType
// trait (1-128 chars). At most 50 keys are accepted. Duplicate keys are
// rejected.
func extractTransitiveTagKeys(params map[string]interface{}) ([]string, error) {
	var keys []string
	seen := make(map[string]bool)
	for i := 1; ; i++ {
		key := request.GetStringParam(params, fmt.Sprintf("TransitiveTagKeys.member.%d", i))
		if key == "" {
			break
		}
		if i > 50 {
			return nil, ErrTooManySessionTags
		}
		if len(key) < 1 || len(key) > 128 || !sessionTagKeyPattern.MatchString(key) {
			return nil, ErrInvalidSessionTag
		}
		if seen[key] {
			return nil, ErrDuplicateSessionTagKey
		}
		seen[key] = true
		keys = append(keys, key)
	}
	return keys, nil
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
}

// STSService provides AWS Security Token Service operations.
type STSService struct {
	stores sync.Map // caches STS SessionStore per region
}

// NewSTSService creates a new STS service instance.
func NewSTSService() *STSService {
	return &STSService{}
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
		// Device exists and is assigned but has no TOTP seed (e.g. a hardware
		// device we cannot verify server-side). Accept the code as-is provided
		// it is well-formed (validated upstream).
		return true, nil
	}
	if !verifyTOTP(device.Base32StringSeed, actx.TokenCode, time.Now().UTC()) {
		return false, ErrInvalidMFACode
	}
	return true, nil
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

// GetCallerIdentity returns details about the IAM user or role whose credentials are used to call the operation.
//
// When the caller cannot be resolved from the security token or access key
// (e.g. during InitialSetup before any IAM user exists), the response falls
// back to the root principal. This is an intentional design choice for
// VorpalStacks to support bootstrap/setup flows; AWS would return 403 for
// unauthenticated requests.
func (s *STSService) GetCallerIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	securityToken := req.Headers.Get("X-Amz-Security-Token")

	var userId, arn string

	if securityToken != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		session, err := store.Get(securityToken)
		if err == nil && session != nil {
			userId = session.AccessKeyId
			switch session.PrincipalType {
			case "AssumedRole":
				roleName := arnutil.ExtractRoleNameFromARN(session.RoleArn)
				arn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(roleName, session.RoleSessionName)
			case "SAML":
				roleName := arnutil.ExtractRoleNameFromARN(session.RoleArn)
				arn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(roleName, session.PrincipalName)
			case "WebIdentity":
				roleName := arnutil.ExtractRoleNameFromARN(session.RoleArn)
				arn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(roleName, session.RoleSessionName)
			case "FederatedUser":
				arn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().FederatedUser(session.PrincipalName)
			case "Root":
				arn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
			default:
				arn = session.PrincipalArn
			}
		}
	}

	if userId == "" && arn == "" {
		callerArn, callerName := s.resolveCallerIdentity(reqCtx, req)
		if callerArn != "" {
			arn = callerArn
			userId = callerName
		}
	}

	if userId == "" {
		userId = reqCtx.GetAccountID()
	}

	if arn == "" {
		arn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}

	return map[string]interface{}{
		"UserId":  userId,
		"Account": reqCtx.GetAccountID(),
		"Arn":     arn,
	}, nil
}

// AssumeRole returns a set of temporary security credentials that you can use to access AWS resources.
func (s *STSService) AssumeRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleArn := request.GetStringParam(req.Parameters, "RoleArn")
	roleSessionName := request.GetStringParam(req.Parameters, "RoleSessionName")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	sessionPolicy := request.GetStringParam(req.Parameters, "Policy")
	sourceIdentity := request.GetStringParam(req.Parameters, "SourceIdentity")
	externalId := request.GetStringParam(req.Parameters, "ExternalId")
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	tokenCode := request.GetStringParam(req.Parameters, "TokenCode")

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if roleArn == "" {
		return nil, ErrInvalidRoleArn
	}

	if err := validateRoleSessionName(roleSessionName); err != nil {
		return nil, err
	}

	if err := validateSourceIdentity(sourceIdentity); err != nil {
		return nil, err
	}

	if externalId != "" {
		if len(externalId) < 2 || len(externalId) > 1224 || !externalIdPattern.MatchString(externalId) {
			return nil, ErrInvalidExternalId
		}
	}

	if err := validateMFACredentials(serialNumber, tokenCode); err != nil {
		return nil, err
	}

	if sessionPolicy != "" {
		var js interface{}
		if err := json.Unmarshal([]byte(sessionPolicy), &js); err != nil {
			return nil, ErrMalformedPolicyDocument
		}
	}

	sessionTags, err := extractSessionTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	transitiveKeys, err := extractTransitiveTagKeys(req.Parameters)
	if err != nil {
		return nil, err
	}

	callerArn, callerName := s.resolveCallerIdentity(reqCtx, req)
	if callerArn == "" {
		callerArn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}

	// Role chaining: when the caller uses temporary credentials, forward
	// transitive session tags from the caller's session. Transitive tags
	// take precedence over new tags with the same key to prevent privilege
	// escalation.
	callerSession := s.resolveCallerSession(reqCtx, req)
	mergedTags := make(map[string]string)
	for k, v := range sessionTags {
		mergedTags[k] = v
	}
	mergedTransitiveKeys := make(map[string]bool)
	for _, tk := range transitiveKeys {
		mergedTransitiveKeys[tk] = true
	}
	if callerSession != nil {
		for _, tk := range callerSession.TransitiveTagKeys {
			if v, ok := callerSession.Tags[tk]; ok {
				mergedTags[tk] = v
			}
			mergedTransitiveKeys[tk] = true
		}
	}
	finalTransitiveKeys := make([]string, 0, len(mergedTransitiveKeys))
	for k := range mergedTransitiveKeys {
		finalTransitiveKeys = append(finalTransitiveKeys, k)
	}

	actx := &assumeContext{
		ExternalId:   externalId,
		SerialNumber: serialNumber,
		TokenCode:    tokenCode,
		CallerName:   callerName,
	}
	role, err := s.resolveRoleForAssume(reqCtx, roleArn, callerArn, "sts:AssumeRole", actx)
	if err != nil {
		return nil, err
	}

	packedPolicySize := computePackedPolicySize(sessionPolicy, req.Parameters, mergedTags)
	if packedPolicySize > 100 {
		return nil, ErrPackedPolicyTooLarge
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:          "AssumedRole",
		PrincipalName:          roleSessionName,
		PrincipalArn:           roleArn,
		RoleArn:                roleArn,
		RoleSessionName:        roleSessionName,
		SourceIdentity:         sourceIdentity,
		DurationSeconds:        validDuration,
		Tags:                   mergedTags,
		TransitiveTagKeys:      finalTransitiveKeys,
		MultiFactorAuthPresent: actx.MFAPresent,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": role.ID + ":" + roleSessionName,
			"Arn":           arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(role.RoleName, roleSessionName),
		},
		"PackedPolicySize": packedPolicySize,
		"SourceIdentity":   session.SourceIdentity,
	}, nil
}

// GetSessionToken returns a set of temporary credentials for an AWS account or IAM user.
func (s *STSService) GetSessionToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	serialNumber := request.GetStringParam(req.Parameters, "SerialNumber")
	tokenCode := request.GetStringParam(req.Parameters, "TokenCode")

	validDuration, err := validateDurationSecondsExtended(durationSeconds)
	if err != nil {
		return nil, err
	}

	accessKeyId := req.Headers.Get("X-Amz-Access-Key")
	if accessKeyId == "" {
		accessKeyId = request.GetStringParam(req.Parameters, "AccessKeyId")
	}

	var callerArn, callerName string

	if accessKeyId != "" {
		callerArn, callerName = s.resolveCallerIdentity(reqCtx, req)
	}

	if callerArn == "" {
		callerName = reqCtx.GetAccountID()
		callerArn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}

	if err := validateMFACredentials(serialNumber, tokenCode); err != nil {
		return nil, err
	}

	mfaPresent := false
	if serialNumber != "" && tokenCode != "" {
		iamStore, err := s.iamStore(reqCtx)
		if err != nil {
			return nil, err
		}
		actx := &assumeContext{
			SerialNumber: serialNumber,
			TokenCode:    tokenCode,
			CallerName:   callerName,
		}
		if _, err := s.verifyCallerMFA(reqCtx, iamStore, actx); err != nil {
			return nil, err
		}
		mfaPresent = true
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:          "User",
		PrincipalName:          callerName,
		PrincipalArn:           callerArn,
		DurationSeconds:        validDuration,
		MultiFactorAuthPresent: mfaPresent,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
	}, nil
}

func computePackedPolicySize(policy string, params map[string]interface{}, tags map[string]string) int32 {
	const maxPolicySize = 2048
	totalSize := len(policy)
	for i := 1; ; i++ {
		arnKey := fmt.Sprintf("PolicyArns.member.%d.arn", i)
		arn := request.GetStringParam(params, arnKey)
		if arn == "" {
			break
		}
		totalSize += len(arn)
	}
	for key, value := range tags {
		totalSize += len(key) + len(value)
	}
	if totalSize == 0 {
		return 0
	}
	return int32((totalSize * 100) / maxPolicySize)
}

func validateDurationSeconds(durationSeconds int) (int, error) {
	if durationSeconds == 0 {
		return DefaultDurationSeconds, nil
	}
	if durationSeconds < MinDurationSeconds || durationSeconds > MaxDurationSeconds {
		return 0, ErrInvalidDuration
	}
	return durationSeconds, nil
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
			return session.PrincipalArn, session.PrincipalName
		}
	}
	return "", ""
}

func validateDurationSecondsExtended(durationSeconds int) (int, error) {
	if durationSeconds == 0 {
		return DefaultDurationSeconds, nil
	}
	if durationSeconds < MinDurationSeconds || durationSeconds > MaxDurationSecondsExtended {
		return 0, ErrInvalidDurationExtended
	}
	return durationSeconds, nil
}

// validateRootDurationSeconds validates the DurationSeconds parameter for
// AssumeRoot per the Smithy RootDurationSecondsType trait (range 0-900,
// default 900). Duration 0 means "use the 900 s default", not
// immediate expiry.
func validateRootDurationSeconds(durationSeconds int) (int, error) {
	if durationSeconds == 0 {
		return DefaultRootDurationSeconds, nil
	}
	if durationSeconds < 0 || durationSeconds > MaxRootDurationSeconds {
		return 0, ErrInvalidRootDuration
	}
	return durationSeconds, nil
}

// validateFederationDurationSeconds validates the DurationSeconds parameter
// for GetFederationToken. The Smithy durationSecondsType trait allows
// 900-129600 with a default of 43200 (12 hours). When the caller is the root
// user, AWS caps the session at 3600 seconds (1 hour); the isRoot parameter
// enforces this cap.
func validateFederationDurationSeconds(durationSeconds int, isRoot bool) (int, error) {
	if durationSeconds == 0 {
		if isRoot {
			return 3600, nil
		}
		return DefaultFederationDurationSeconds, nil
	}
	if durationSeconds < MinDurationSeconds || durationSeconds > MaxDurationSecondsExtended {
		return 0, ErrInvalidDurationExtended
	}
	if isRoot && durationSeconds > 3600 {
		return 0, ErrInvalidFederationRootDuration
	}
	return durationSeconds, nil
}

// federationNamePattern mirrors the Smithy userNameType trait used by
// GetFederationToken's Name parameter: [\w+=,.@-]* with length 2-32.
var federationNamePattern = regexp.MustCompile(`^[\w+=,.@-]+$`)

// sessionNamePattern mirrors the Smithy roleSessionNameType and
// sourceIdentityType traits: [\w+=,.@-]* with length 2-64.
var sessionNamePattern = regexp.MustCompile(`^[\w+=,.@-]{2,64}$`)

// validateRoleSessionName checks the RoleSessionName parameter against the
// Smithy roleSessionNameType trait (pattern [\w+=,.@-]*, length 2-64).
func validateRoleSessionName(name string) error {
	if name == "" {
		return ErrInvalidParameter
	}
	if !sessionNamePattern.MatchString(name) {
		return ErrInvalidParameter
	}
	return nil
}

// validateSourceIdentity checks the SourceIdentity parameter against the
// Smithy sourceIdentityType trait (pattern [\w+=,.@-]*, length 2-64) and
// rejects values beginning with the reserved "aws:" prefix.
func validateSourceIdentity(si string) error {
	if si == "" {
		return nil
	}
	if strings.HasPrefix(si, "aws:") {
		return ErrInvalidSourceIdentity
	}
	if !sessionNamePattern.MatchString(si) {
		return ErrInvalidSourceIdentity
	}
	return nil
}

// allowedRootTaskPolicyNames enumerates the AWS root-task managed policies
// that AssumeRoot accepts as TaskPolicyArn. AWS docs reference:
// https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoot.html
var allowedRootTaskPolicyNames = map[string]bool{
	"IAMAuditRootUserCredentials":  true,
	"IAMCreateRootUserPassword":    true,
	"IAMDeleteRootUserCredentials": true,
	"S3UnlockBucketPolicy":         true,
	"SQSUnlockQueuePolicy":         true,
}

// extractPolicyNameFromArn returns the trailing path segment of a policy ARN
// (e.g. "IAMAuditRootUserCredentials" from
// "arn:aws:iam::aws:policy/root-task/IAMAuditRootUserCredentials"). Bare policy
// names without an ARN prefix are returned as-is for compatibility with SDK
// clients that send only the policy name.
func extractPolicyNameFromArn(arn string) string {
	if idx := strings.LastIndex(arn, "/"); idx >= 0 {
		return arn[idx+1:]
	}
	return arn
}

// AssumeRoleWithSAML returns a set of temporary security credentials for users who have been authenticated via a SAML authentication response.
// VorpalStacks cannot validate SAML assertions against external IdPs, so this is only available in TEST_MODE.
func (s *STSService) AssumeRoleWithSAML(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if os.Getenv("TEST_MODE") != "true" {
		return nil, ErrIDPCommunicationError
	}

	roleArn := request.GetStringParam(req.Parameters, "RoleArn")
	principalArn := request.GetStringParam(req.Parameters, "PrincipalArn")
	samlAssertion := request.GetStringParam(req.Parameters, "SAMLAssertion")
	roleSessionName := request.GetStringParam(req.Parameters, "RoleSessionName")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	sessionPolicy := request.GetStringParam(req.Parameters, "Policy")

	if roleSessionName == "" {
		roleSessionName = "SAML"
	} else if err := validateRoleSessionName(roleSessionName); err != nil {
		return nil, err
	}

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if roleArn == "" {
		return nil, ErrInvalidRoleArn
	}

	if principalArn == "" {
		return nil, ErrInvalidPrincipalArn
	}

	if samlAssertion == "" {
		return nil, ErrInvalidSAMLAssertion
	}

	if sessionPolicy != "" {
		var js interface{}
		if err := json.Unmarshal([]byte(sessionPolicy), &js); err != nil {
			return nil, ErrMalformedPolicyDocument
		}
	}

	if _, err := base64.StdEncoding.DecodeString(samlAssertion); err != nil {
		if _, err := base64.URLEncoding.DecodeString(samlAssertion); err != nil {
			return nil, ErrInvalidSAMLAssertion
		}
	}

	role, err := s.resolveRoleForAssume(reqCtx, roleArn, principalArn, "sts:AssumeRoleWithSAML", nil)
	if err != nil {
		return nil, err
	}
	packedPolicySize := computePackedPolicySize(sessionPolicy, req.Parameters, nil)
	if packedPolicySize > 100 {
		return nil, ErrPackedPolicyTooLarge
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "SAML",
		PrincipalName:   principalArn,
		PrincipalArn:    roleArn,
		RoleArn:         roleArn,
		RoleSessionName: roleSessionName,
		DurationSeconds: validDuration,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": role.ID + ":" + roleSessionName,
			"Arn":           arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(role.RoleName, roleSessionName),
		},
		"Subject":          principalArn,
		"SubjectType":      "persistent",
		"Issuer":           "VorpalStacks",
		"NameQualifier":    "SAML",
		"Audience":         "STS",
		"PackedPolicySize": packedPolicySize,
	}, nil
}

// AssumeRoleWithWebIdentity returns a set of temporary security credentials for users who have been authenticated in a mobile or web application with a web identity provider.
// VorpalStacks cannot validate web identity tokens against external IdPs, so this is only available in TEST_MODE.
func (s *STSService) AssumeRoleWithWebIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if os.Getenv("TEST_MODE") != "true" {
		return nil, ErrIDPCommunicationError
	}

	roleArn := request.GetStringParam(req.Parameters, "RoleArn")
	roleSessionName := request.GetStringParam(req.Parameters, "RoleSessionName")
	webIdentityToken := request.GetStringParam(req.Parameters, "WebIdentityToken")
	providerId := request.GetStringParam(req.Parameters, "ProviderId")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	sessionPolicy := request.GetStringParam(req.Parameters, "Policy")
	sourceIdentity := request.GetStringParam(req.Parameters, "SourceIdentity")

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if roleArn == "" {
		return nil, ErrInvalidRoleArn
	}

	if roleSessionName == "" {
		roleSessionName = "web-identity-session"
	} else if err := validateRoleSessionName(roleSessionName); err != nil {
		return nil, err
	}

	if err := validateSourceIdentity(sourceIdentity); err != nil {
		return nil, err
	}

	if webIdentityToken == "" {
		return nil, ErrInvalidWebIdentityToken
	}
	// clientTokenType Smithy trait: length 4-20000.
	if len(webIdentityToken) < 4 || len(webIdentityToken) > 20000 {
		return nil, ErrInvalidWebIdentityToken
	}

	if sessionPolicy != "" {
		var js interface{}
		if err := json.Unmarshal([]byte(sessionPolicy), &js); err != nil {
			return nil, ErrMalformedPolicyDocument
		}
	}

	// ProviderId accepts a URL (e.g. "www.amazon.com") or a full ARN per
	// the AWS spec. When it is already an ARN, use it verbatim; otherwise
	// construct the OIDC provider ARN from the URL.
	federatedPrincipal := ""
	if providerId != "" {
		if strings.HasPrefix(providerId, "arn:") {
			federatedPrincipal = providerId
		} else {
			federatedPrincipal = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().OIDCProvider(providerId)
		}
	}

	role, err := s.resolveRoleForAssume(reqCtx, roleArn, federatedPrincipal, "sts:AssumeRoleWithWebIdentity", nil)
	if err != nil {
		return nil, err
	}

	packedPolicySize := computePackedPolicySize(sessionPolicy, req.Parameters, nil)
	if packedPolicySize > 100 {
		return nil, ErrPackedPolicyTooLarge
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "WebIdentity",
		PrincipalName:   roleSessionName,
		PrincipalArn:    roleArn,
		RoleArn:         roleArn,
		RoleSessionName: roleSessionName,
		SourceIdentity:  sourceIdentity,
		DurationSeconds: validDuration,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"AssumedRoleUser": map[string]interface{}{
			"AssumedRoleId": role.ID + ":" + roleSessionName,
			"Arn":           arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(role.RoleName, roleSessionName),
		},
		"Provider":                    providerId,
		"SubjectFromWebIdentityToken": roleSessionName,
		"Audience":                    "sts.amazonaws.com",
		"PackedPolicySize":            packedPolicySize,
		"SourceIdentity":              session.SourceIdentity,
	}, nil
}

// AssumeRoot returns a set of temporary security credentials for performing
// privileged tasks on a member account. AWS requires the caller to be an
// Organizations management account or IAM delegated administrator; because
// VorpalStacks does not implement Organizations (see docs/services.md
// "No organisations integration") the caller authorisation check is omitted
// by design. Parameter validation (TargetPrincipal, TaskPolicyArn,
// DurationSeconds) is still enforced so that AWS SDK clients receive the
// correct validation errors and the session is scoped to a known task policy.
func (s *STSService) AssumeRoot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	targetPrincipal := request.GetStringParam(req.Parameters, "TargetPrincipal")
	taskPolicyArn := request.GetStringParam(req.Parameters, "TaskPolicyArn.arn")

	validDuration, err := validateRootDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if targetPrincipal == "" {
		return nil, ErrTargetPrincipalRequired
	}
	// TargetPrincipalType Smithy trait: length 12-2048. Accepts account ID
	// (12 digits) or principal ARN.
	if len(targetPrincipal) < 12 || len(targetPrincipal) > 2048 {
		return nil, ErrInvalidTargetPrincipal
	}

	if taskPolicyArn == "" {
		return nil, ErrTaskPolicyArnRequired
	}
	policyName := extractPolicyNameFromArn(taskPolicyArn)
	if !allowedRootTaskPolicyNames[policyName] {
		return nil, ErrInvalidTaskPolicyArn
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Store TargetPrincipal as PrincipalName (the member account being
	// accessed) with an empty RoleSessionName — root sessions do not use
	// the assumed-role session name slot.
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "Root",
		PrincipalName:   targetPrincipal,
		PrincipalArn:    arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root(),
		DurationSeconds: validDuration,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
	}, nil
}

// DecodeAuthorizationMessage decodes additional information about the authorization status of a request from an encoded message.
func (s *STSService) DecodeAuthorizationMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	encodedMessage := request.GetStringParam(req.Parameters, "EncodedMessage")

	if encodedMessage == "" {
		return nil, ErrInvalidEncodedMessage
	}
	// encodedMessageType Smithy trait: length 1-10240.
	if len(encodedMessage) > 10240 {
		return nil, ErrInvalidEncodedMessage
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedMessage)
	if err != nil {
		decodedBytes, err = base64.URLEncoding.DecodeString(encodedMessage)
		if err != nil {
			return nil, ErrInvalidEncodedMessage
		}
	}

	return map[string]interface{}{
		"DecodedMessage": string(decodedBytes),
	}, nil
}

// GetAccessKeyInfo returns information about the access key in the request.
func (s *STSService) GetAccessKeyInfo(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	accessKeyId := request.GetStringParam(req.Parameters, "AccessKeyId")

	if accessKeyId == "" {
		return nil, ErrInvalidAccessKeyId
	}
	// accessKeyIdType Smithy trait: length 16-128, pattern [\w]*.
	if len(accessKeyId) < 16 || len(accessKeyId) > 128 {
		return nil, ErrInvalidAccessKeyId
	}

	return map[string]interface{}{
		"Account": reqCtx.GetAccountID(),
		"Status":  "Active",
	}, nil
}

// GetFederationToken returns a set of temporary security credentials for a federated user.
func (s *STSService) GetFederationToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	policy := request.GetStringParam(req.Parameters, "Policy")
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")

	if name == "" {
		return nil, ErrInvalidFederationName
	}
	// userNameType Smithy trait: length 2-32, pattern [\w+=,.@-]*.
	if len(name) < 2 || len(name) > 32 || !federationNamePattern.MatchString(name) {
		return nil, ErrInvalidFederationName
	}

	callerArn, callerName := s.resolveCallerIdentity(reqCtx, req)
	if callerArn == "" {
		callerArn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}
	isRoot := callerName == iam.RootUserName

	validDuration, err := validateFederationDurationSeconds(durationSeconds, isRoot)
	if err != nil {
		return nil, err
	}

	if policy != "" {
		var js interface{}
		if err := json.Unmarshal([]byte(policy), &js); err != nil {
			return nil, ErrMalformedPolicyDocument
		}
	}

	fedTags, err := extractSessionTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	packedPolicySize := computePackedPolicySize(policy, req.Parameters, fedTags)
	if packedPolicySize > 100 {
		return nil, ErrPackedPolicyTooLarge
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "FederatedUser",
		PrincipalName:   name,
		PrincipalArn:    callerArn,
		RoleSessionName: name,
		DurationSeconds: validDuration,
		Tags:            fedTags,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"FederatedUser": map[string]interface{}{
			"FederatedUserId": reqCtx.GetAccountID() + ":" + name,
			"Arn":             arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().FederatedUser(name),
		},
		"PackedPolicySize": packedPolicySize,
	}, nil
}

// GetDelegatedAccessToken returns a set of temporary security credentials that represent an IAM identity centre user.
func (s *STSService) GetDelegatedAccessToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tradeInToken := request.GetStringParam(req.Parameters, "TradeInToken")

	if tradeInToken == "" {
		return nil, ErrInvalidTradeInToken
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	principalArn, err := store.RedeemDelegationToken(tradeInToken)
	if err != nil {
		if errors.Is(err, stsstore.ErrDelegationTokenExpired) {
			return nil, ErrExpiredTradeInToken
		}
		return nil, ErrInvalidTradeInToken
	}

	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "DelegatedAccess",
		PrincipalName:   principalArn,
		PrincipalArn:    principalArn,
		DurationSeconds: 3600,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Credentials": map[string]interface{}{
			"AccessKeyId":     session.AccessKeyId,
			"SecretAccessKey": session.SecretAccessKey,
			"SessionToken":    session.SessionToken,
			"Expiration":      session.Expiration.Format(timeutils.ISO8601SimpleFormat),
		},
		"AssumedPrincipal": principalArn,
		"PackedPolicySize": 0,
	}, nil
}

// GetWebIdentityToken returns a signed JSON Web Token (JWT) representing the
// calling AWS identity. The returned JWT can be used to authenticate with
// external services that support OIDC discovery. The token is signed using
// the caller-specified algorithm (RS256 or ES384).
//
// AWS spec: https://docs.aws.amazon.com/STS/latest/APIReference/API_GetWebIdentityToken.html
func (s *STSService) GetWebIdentityToken(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	durationSeconds := request.GetIntParam(req.Parameters, "DurationSeconds")
	signingAlgorithm := request.GetStringParam(req.Parameters, "SigningAlgorithm")

	validDuration, err := validateWebIdentityDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if signingAlgorithm != "RS256" && signingAlgorithm != "ES384" {
		return nil, ErrInvalidSigningAlgorithm
	}

	var audiences []string
	for i := 1; ; i++ {
		key := fmt.Sprintf("Audience.member.%d", i)
		val := request.GetStringParam(req.Parameters, key)
		if val == "" {
			break
		}
		audiences = append(audiences, val)
	}
	if len(audiences) == 0 {
		return nil, ErrAudienceRequired
	}
	if len(audiences) > 10 {
		return nil, ErrTooManyAudiences
	}

	callerArn, callerName := s.resolveCallerIdentity(reqCtx, req)
	if callerArn == "" {
		callerArn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
		callerName = reqCtx.GetAccountID()
	}

	token, expiration, err := webIdentityTokenManagerInstance.generateToken(
		callerName, reqCtx.GetAccountID(), audiences, signingAlgorithm, validDuration,
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"WebIdentityToken": token,
		"Expiration":       expiration.Format(timeutils.ISO8601SimpleFormat),
	}, nil
}
