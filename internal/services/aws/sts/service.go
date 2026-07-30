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

// extractPolicyArns parses PolicyArns.member.N.arn from the flat query-
// parameter map. AWS Smithy policyDescriptorType defines the ARN shape; the
// pack-size limit on the surrounding operation enforces an upper bound on
// the number of entries indirectly. The slice preserves caller order.
func extractPolicyArns(params map[string]interface{}) []string {
	var arns []string
	for i := 1; ; i++ {
		key := fmt.Sprintf("PolicyArns.member.%d.arn", i)
		arn := request.GetStringParam(params, key)
		if arn == "" {
			break
		}
		arns = append(arns, arn)
	}
	return arns
}

// extractProvidedContexts parses ProvidedContexts.member.N.ProviderArn and
// ProvidedContexts.member.N.ContextAssertion from the flat query-parameter
// map. Smithy ProvidedContextsListType limits the list to 1-5 entries; the
// trust policy evaluator inspects sts:ProvidedContextProviderArn and
// sts:ProvidedContextAssertion condition keys. Provided contexts are
// signed-and-encrypted by STS in real AWS; VorpalStacks does not verify the
// signature but exposes the values to the evaluation context for
// compatibility with policy templates that reference them.
func extractProvidedContexts(params map[string]interface{}) []ProvidedContextEntry {
	var entries []ProvidedContextEntry
	for i := 1; ; i++ {
		providerKey := fmt.Sprintf("ProvidedContexts.member.%d.ProviderArn", i)
		assertionKey := fmt.Sprintf("ProvidedContexts.member.%d.ContextAssertion", i)
		providerArn := request.GetStringParam(params, providerKey)
		contextAssertion := request.GetStringParam(params, assertionKey)
		if providerArn == "" && contextAssertion == "" {
			break
		}
		entries = append(entries, ProvidedContextEntry{
			ProviderArn:      providerArn,
			ContextAssertion: contextAssertion,
		})
	}
	return entries
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
		// M3: surface caller-supplied ProvidedContexts as
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
				arn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").STS().AssumedRole(roleName, session.RoleSessionName)
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
		ExternalId:       externalId,
		SerialNumber:     serialNumber,
		TokenCode:        tokenCode,
		CallerName:       callerName,
		ProvidedContexts: extractProvidedContexts(req.Parameters),
	}
	role, err := s.resolveRoleForAssume(reqCtx, roleArn, callerArn, "sts:AssumeRole", actx)
	if err != nil {
		return nil, err
	}

	// Enforce role.MaxSessionDuration after role resolution. Without this
	// gate a caller could request DurationSeconds beyond the role's
	// configured maximum (security control bypass).
	err = enforceRoleMaxSessionDuration(validDuration, role.MaxSessionDuration)
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
	policyArns := extractPolicyArns(req.Parameters)
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
		Policy:                 sessionPolicy,
		PolicyArns:             policyArns,
	})
	if err != nil {
		return nil, err
	}

	iamStore, err := s.iamStore(reqCtx)
	if err == nil {
		_ = iamStore.Roles().UpdateRoleLastUsed(role.RoleName, reqCtx.GetRegion())
	}

	return withSourceIdentity(map[string]interface{}{
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
	}, session.SourceIdentity), nil
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

	// AWS caps GetSessionToken duration at 3600 seconds for root users.
	isRoot := strings.HasSuffix(callerArn, ":root")
	if isRoot && validDuration > 3600 {
		return nil, ErrInvalidRootSessionDuration
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

	// AWS serialises the session policy, managed-policy ARNs, session
	// tags, transitive tag keys and source identity into a packed JSON
	// document whose total byte length must not exceed maxPolicySize.
	// The raw character counts alone underestimate the packed size
	// because they ignore JSON structural overhead (quotes, brackets,
	// key names, separators). We add per-element overhead so the
	// reported percentage more closely tracks real AWS.
	totalSize := len(policy)

	// Managed-policy ARNs — each serialised as a JSON string element
	// inside an array: ["arn:...","arn:..."]. Per-element overhead is
	// 3 bytes (two quotes + one comma/bracket).
	arnCount := 0
	for i := 1; ; i++ {
		arnKey := fmt.Sprintf("PolicyArns.member.%d.arn", i)
		arn := request.GetStringParam(params, arnKey)
		if arn == "" {
			break
		}
		totalSize += len(arn) + 3
		arnCount++
	}
	if arnCount > 0 {
		totalSize += 2 // array brackets
	}

	// Session tags — each serialised as {"Key":"...","Value":"..."}.
	// Structural overhead per tag: {"Key":"","Value":""} = 20 bytes.
	if len(tags) > 0 {
		totalSize += 2 // array brackets
	}
	for key, value := range tags {
		totalSize += len(key) + len(value) + 20
	}

	// Transitive tag keys — each serialised as a JSON string.
	transitiveCount := 0
	for i := 1; ; i++ {
		key := fmt.Sprintf("TransitiveTagKeys.member.%d", i)
		val := request.GetStringParam(params, key)
		if val == "" {
			break
		}
		totalSize += len(val) + 3
		transitiveCount++
	}
	if transitiveCount > 0 {
		totalSize += 2
	}

	if totalSize <= 0 {
		return 0
	}
	return int32((totalSize * 100) / maxPolicySize)
}

// withSourceIdentity returns resp with the SourceIdentity field set only when
// si is non-empty. The Smithy AssumeRole*Response shapes declare SourceIdentity
// as an optional member; AWS Query protocol responses omit the field entirely
// when the caller did not supply it, rather than serialising an empty string.
func withSourceIdentity(resp map[string]interface{}, si string) map[string]interface{} {
	if si != "" {
		resp["SourceIdentity"] = si
	}
	return resp
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

// enforceRoleMaxSessionDuration rejects a DurationSeconds value that
// exceeds the role's configured MaxSessionDuration. The caller must
// already have normalised 0 → DefaultDurationSeconds via
// validateDurationSeconds before calling. A roleMaxSessionDuration <= 0
// means the role store defaulted it to 3600 already, so no enforcement
// is needed.
func enforceRoleMaxSessionDuration(durationSeconds, roleMaxSessionDuration int) error {
	if roleMaxSessionDuration <= 0 {
		return nil
	}
	if durationSeconds > roleMaxSessionDuration {
		return ErrInvalidDuration
	}
	return nil
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
	// Enforce role.MaxSessionDuration after role resolution (H1).
	err = enforceRoleMaxSessionDuration(validDuration, role.MaxSessionDuration)
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
	policyArns := extractPolicyArns(req.Parameters)
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "SAML",
		PrincipalName:   principalArn,
		PrincipalArn:    roleArn,
		RoleArn:         roleArn,
		RoleSessionName: roleSessionName,
		DurationSeconds: validDuration,
		Policy:          sessionPolicy,
		PolicyArns:      policyArns,
	})
	if err != nil {
		return nil, err
	}

	// L5: extract Issuer / Subject (NameID) / Audience from the SAML assertion
	// when it is parseable XML. Falls back to legacy placeholder values for
	// non-XML tokens (the SDK test fixture is a base64 plain-text dummy).
	samlIssuer, samlSubject, samlAudience := decodeSAMLFields(samlAssertion)
	if samlIssuer == "" {
		samlIssuer = "VorpalStacks"
	}
	if samlSubject == "" {
		samlSubject = principalArn
	}
	if samlAudience == "" {
		samlAudience = "STS"
	}

	return withSourceIdentity(map[string]interface{}{
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
		"Subject":          samlSubject,
		"SubjectType":      "persistent",
		"Issuer":           samlIssuer,
		"NameQualifier":    "SAML",
		"Audience":         samlAudience,
		"PackedPolicySize": packedPolicySize,
		// SourceIdentity is optional in the Smithy
		// AssumeRoleWithSAMLResponse shape. AWS derives the value from the
		// SAML assertion's saml:AttributeStatement, which VorpalStacks
		// does not parse in TEST_MODE; withSourceIdentity omits the field
		// until a real SAML parser is introduced.
	}, ""), nil
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

	// Enforce role.MaxSessionDuration after role resolution (H1).
	err = enforceRoleMaxSessionDuration(validDuration, role.MaxSessionDuration)
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
	policyArns := extractPolicyArns(req.Parameters)
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "WebIdentity",
		PrincipalName:   roleSessionName,
		PrincipalArn:    roleArn,
		RoleArn:         roleArn,
		RoleSessionName: roleSessionName,
		SourceIdentity:  sourceIdentity,
		DurationSeconds: validDuration,
		Policy:          sessionPolicy,
		PolicyArns:      policyArns,
	})
	if err != nil {
		return nil, err
	}

	// M5: SubjectFromWebIdentityToken should be the OIDC sub claim of the
	// caller-supplied token. L4: Audience should be the aud claim. Both
	// fall back to legacy defaults when the token is not a parseable JWT
	// (e.g. the SDK test dummy 'dummy-web-identity-token').
	subject := extractJWTClaim(webIdentityToken, "sub")
	if subject == "" {
		subject = roleSessionName
	}
	audience := extractJWTClaim(webIdentityToken, "aud")
	if audience == "" {
		audience = "sts.amazonaws.com"
	}

	return withSourceIdentity(map[string]interface{}{
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
		"SubjectFromWebIdentityToken": subject,
		"Audience":                    audience,
		"PackedPolicySize":            packedPolicySize,
	}, session.SourceIdentity), nil
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

	// The Smithy GetAccessKeyInfoResponse shape declares only the Account
	// member; emitting an extra Status field would be non-compliant and
	// could confuse SDK clients that strict-deserialise responses.
	return map[string]interface{}{
		"Account": reqCtx.GetAccountID(),
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

	callerArn, _ := s.resolveCallerIdentity(reqCtx, req)
	if callerArn == "" {
		callerArn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root()
	}
	isRoot := strings.HasSuffix(callerArn, ":root")

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
	policyArns := extractPolicyArns(req.Parameters)
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "FederatedUser",
		PrincipalName:   name,
		PrincipalArn:    callerArn,
		RoleSessionName: name,
		DurationSeconds: validDuration,
		Tags:            fedTags,
		Policy:          policy,
		PolicyArns:      policyArns,
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

	// M4: parse caller-supplied Tags and forward them to the JWT signer so
	// external services consuming the token can read the caller-attached
	// session tags from the standard "tags" JWT claim.
	tags, err := extractSessionTags(req.Parameters)
	if err != nil {
		return nil, err
	}

	token, expiration, err := webIdentityTokenManagerInstance.generateToken(
		callerName, reqCtx.GetAccountID(), audiences, signingAlgorithm, validDuration, tags,
	)
	if err != nil {
		return nil, err
	}

	// Smithy GetWebIdentityToken can return JWTPayloadSizeExceededException
	// when the serialised JWT exceeds the platform limit. AWS enforces a
	// practical ceiling; we check the final token length as a guard.
	if len(token) > 32768 {
		return nil, ErrJWTPayloadSizeExceeded
	}

	return map[string]interface{}{
		"WebIdentityToken": token,
		"Expiration":       expiration.Format(timeutils.ISO8601SimpleFormat),
	}, nil
}
