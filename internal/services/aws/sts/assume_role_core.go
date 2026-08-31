package sts

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
	stsstore "vorpalstacks/internal/store/aws/sts"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

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
	// The IAM role store normalises MaxSessionDuration to 3600 when it
	// is zero or unset. A zero value here means the role data is
	// malformed — fail-closed by treating it as the AWS default (3600)
	// rather than silently skipping enforcement.
	if roleMaxSessionDuration <= 0 {
		roleMaxSessionDuration = 3600
	}
	if durationSeconds > roleMaxSessionDuration {
		return ErrInvalidDuration
	}
	return nil
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

// rootTaskPolicyDocuments maps each AWS root-task managed policy to the
// policy document that scopes the AssumeRoot session. The action sets are
// the ones AWS documents for each task policy ("Perform a privileged task on
// an AWS Organizations member account", IAM User Guide): audit reads for
// IAMAuditRootUserCredentials, credential deletion for
// IAMDeleteRootUserCredentials, password recovery for
// IAMCreateRootUserPassword, bucket-policy replacement for
// S3UnlockBucketPolicy and queue-policy replacement for SQSUnlockQueuePolicy.
// The document is attached as the session's inline policy so the authoriser
// evaluates every request of the session against the task scope instead of
// granting unrestricted root access.
var rootTaskPolicyDocuments = map[string]string{
	"IAMAuditRootUserCredentials": `{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "iam:GetUser",
      "iam:GetLoginProfile",
      "iam:ListAccessKeys",
      "iam:ListSigningCertificates",
      "iam:ListMFADevices",
      "iam:GetAccessKeyLastUsed"
    ],
    "Resource": "*"
  }]
}`,
	"IAMCreateRootUserPassword": `{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "iam:GetLoginProfile",
      "iam:CreateLoginProfile"
    ],
    "Resource": "*"
  }]
}`,
	"IAMDeleteRootUserCredentials": `{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "iam:DeleteLoginProfile",
      "iam:DeleteAccessKey",
      "iam:DeleteSigningCertificate",
      "iam:DeactivateMFADevice"
    ],
    "Resource": "*"
  }]
}`,
	"S3UnlockBucketPolicy": `{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "s3:ListBuckets",
      "s3:GetBucketPolicy",
      "s3:PutBucketPolicy",
      "s3:DeleteBucketPolicy"
    ],
    "Resource": "*"
  }]
}`,
	"SQSUnlockQueuePolicy": `{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "sqs:ListQueues",
      "sqs:GetQueueUrl",
      "sqs:GetQueueAttributes",
      "sqs:SetQueueAttributes"
    ],
    "Resource": "*"
  }]
}`,
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

// AssumeRoleResult carries the serialisation inputs of the AssumeRole
// response: the created session's credentials, the resolved role identity
// and the effective session name.
type AssumeRoleResult struct {
	Credentials      CredentialsResult
	RoleID           string
	RoleName         string
	RoleSessionName  string
	SourceIdentity   string
	PackedPolicySize int32
}

// assumeRoleCore owns the AssumeRole flow: parameter validation in the
// AWS-documented order, caller resolution, transitive session-tag merging,
// trust-policy evaluation with ExternalId / MFA / ProvidedContexts, the
// packed-policy-size gate and session persistence.
func (s *STSService) assumeRoleCore(reqCtx *request.RequestContext, in WireInput) (*AssumeRoleResult, error) {
	roleArn := request.GetStringParam(in.Parameters, "RoleArn")
	roleSessionName := request.GetStringParam(in.Parameters, "RoleSessionName")
	durationSeconds := request.GetIntParam(in.Parameters, "DurationSeconds")
	sessionPolicy := request.GetStringParam(in.Parameters, "Policy")
	sourceIdentity := request.GetStringParam(in.Parameters, "SourceIdentity")
	externalId := request.GetStringParam(in.Parameters, "ExternalId")
	serialNumber := request.GetStringParam(in.Parameters, "SerialNumber")
	tokenCode := request.GetStringParam(in.Parameters, "TokenCode")

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if roleArn == "" {
		return nil, ErrInvalidRoleArn
	}
	if err := validateRoleArn(roleArn); err != nil {
		return nil, err
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

	if err := validateUnrestrictedSessionPolicy(sessionPolicy); err != nil {
		return nil, err
	}

	sessionTags, err := extractSessionTags(in.Parameters)
	if err != nil {
		return nil, err
	}

	transitiveKeys, err := extractTransitiveTagKeys(in.Parameters)
	if err != nil {
		return nil, err
	}

	callerArn, callerName, err := s.resolveCallerArnOrReject(reqCtx, in.AccessKeyID)
	if err != nil {
		return nil, err
	}

	// Role chaining: when the caller uses temporary credentials, forward
	// transitive session tags from the caller's session. Transitive tags
	// take precedence over new tags with the same key to prevent privilege
	// escalation.
	callerSession := s.resolveCallerSession(reqCtx, in.SecurityToken)
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
				mergedTransitiveKeys[tk] = true
			}
		}
	}
	// Validate that every request-sent transitive tag key has a
	// corresponding value in the merged tags (either from the current
	// request's Tags or inherited from the caller session). A
	// transitive key without a value creates an inconsistent session:
	// TransitiveTagKeys references a key that Tags does not contain,
	// making the transitive flag a no-op in subsequent role chains.
	for _, tk := range transitiveKeys {
		if _, ok := mergedTags[tk]; !ok {
			return nil, ErrTransitiveKeyWithoutTag
		}
	}
	finalTransitiveKeys := make([]string, 0, len(mergedTransitiveKeys))
	for k := range mergedTransitiveKeys {
		finalTransitiveKeys = append(finalTransitiveKeys, k)
	}

	providedContexts, err := extractProvidedContexts(in.Parameters)
	if err != nil {
		return nil, err
	}

	actx := &assumeContext{
		ExternalId:       externalId,
		SerialNumber:     serialNumber,
		TokenCode:        tokenCode,
		CallerName:       callerName,
		ProvidedContexts: providedContexts,
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

	policyArns, err := extractPolicyArns(in.Parameters)
	if err != nil {
		return nil, err
	}

	packedPolicySize := computePackedPolicySize(sessionPolicy, policyArns, finalTransitiveKeys, mergedTags)
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

	return &AssumeRoleResult{
		Credentials:      credentialsOf(session),
		RoleID:           role.ID,
		RoleName:         role.RoleName,
		RoleSessionName:  roleSessionName,
		SourceIdentity:   session.SourceIdentity,
		PackedPolicySize: packedPolicySize,
	}, nil
}

// AssumeRoleSAMLResult carries the serialisation inputs of the
// AssumeRoleWithSAML response: credentials, role identity and the
// SAML assertion fields (subject / issuer / audience).
type AssumeRoleSAMLResult struct {
	Credentials      CredentialsResult
	RoleID           string
	RoleName         string
	RoleSessionName  string
	Subject          string
	Issuer           string
	Audience         string
	PackedPolicySize int32
}

// assumeRoleWithSAMLCore owns the AssumeRoleWithSAML flow. VorpalStacks
// cannot validate SAML assertions against external IdPs, so the operation is
// only available in TEST_MODE; the Core rejects production callers first,
// then validates parameters, resolves the role against the SAML principal
// and creates the session.
func (s *STSService) assumeRoleWithSAMLCore(reqCtx *request.RequestContext, in WireInput) (*AssumeRoleSAMLResult, error) {
	if os.Getenv("TEST_MODE") != "true" {
		return nil, ErrIDPCommunicationError
	}

	roleArn := request.GetStringParam(in.Parameters, "RoleArn")
	principalArn := request.GetStringParam(in.Parameters, "PrincipalArn")
	samlAssertion := request.GetStringParam(in.Parameters, "SAMLAssertion")
	roleSessionName := request.GetStringParam(in.Parameters, "RoleSessionName")
	durationSeconds := request.GetIntParam(in.Parameters, "DurationSeconds")
	sessionPolicy := request.GetStringParam(in.Parameters, "Policy")

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
	if err := validateRoleArn(roleArn); err != nil {
		return nil, err
	}

	if principalArn == "" {
		return nil, ErrInvalidPrincipalArn
	}
	if err := validatePrincipalArn(principalArn); err != nil {
		return nil, err
	}

	if samlAssertion == "" {
		return nil, ErrInvalidSAMLAssertion
	}
	if err := validateSAMLAssertion(samlAssertion); err != nil {
		return nil, err
	}

	if err := validateSessionPolicy(sessionPolicy); err != nil {
		return nil, err
	}

	if _, err := base64.StdEncoding.DecodeString(samlAssertion); err != nil {
		if _, err := base64.URLEncoding.DecodeString(samlAssertion); err != nil {
			return nil, ErrInvalidSAMLAssertion
		}
	}

	// Soft SAML expiry check: when the assertion is parseable XML with a
	// Conditions/NotOnOrAfter attribute, reject expired assertions with
	// ExpiredTokenException per the Smithy error mapping. Dummy tokens
	// used in SDK tests are not affected.
	if isSAMLAssertionExpired(samlAssertion) {
		return nil, ErrExpiredToken
	}

	role, err := s.resolveRoleForAssume(reqCtx, roleArn, principalArn, "sts:AssumeRoleWithSAML", nil)
	if err != nil {
		return nil, err
	}
	// Enforce role.MaxSessionDuration after role resolution.
	err = enforceRoleMaxSessionDuration(validDuration, role.MaxSessionDuration)
	if err != nil {
		return nil, err
	}
	policyArns, err := extractPolicyArns(in.Parameters)
	if err != nil {
		return nil, err
	}
	packedPolicySize := computePackedPolicySize(sessionPolicy, policyArns, nil, nil)
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
		Policy:          sessionPolicy,
		PolicyArns:      policyArns,
	})
	if err != nil {
		return nil, err
	}

	// Extract Issuer / Subject (NameID) / Audience from the SAML assertion
	// when it is parseable XML. Falls back to legacy placeholder values for
	// non-XML tokens (the SDK test fixture is a base64 plain-text dummy).
	samlIssuer, samlSubject, samlAudience, _ := decodeSAMLFields(samlAssertion)
	if samlIssuer == "" {
		samlIssuer = "VorpalStacks"
	}
	if samlSubject == "" {
		samlSubject = principalArn
	}
	if samlAudience == "" {
		samlAudience = "STS"
	}

	return &AssumeRoleSAMLResult{
		Credentials:      credentialsOf(session),
		RoleID:           role.ID,
		RoleName:         role.RoleName,
		RoleSessionName:  roleSessionName,
		Subject:          samlSubject,
		Issuer:           samlIssuer,
		Audience:         samlAudience,
		PackedPolicySize: packedPolicySize,
	}, nil
}

// AssumeRoleWebIdentityResult carries the serialisation inputs of the
// AssumeRoleWithWebIdentity response: credentials, role identity and the
// OIDC-derived subject / audience / provider fields.
type AssumeRoleWebIdentityResult struct {
	Credentials      CredentialsResult
	RoleID           string
	RoleName         string
	RoleSessionName  string
	SourceIdentity   string
	Subject          string
	Audience         string
	Provider         string
	PackedPolicySize int32
}

// assumeRoleWithWebIdentityCore owns the AssumeRoleWithWebIdentity flow.
// VorpalStacks cannot validate web identity tokens against external IdPs, so
// the operation is only available in TEST_MODE; the Core rejects production
// callers first, then validates parameters, resolves the role against the
// federated principal and creates the session.
func (s *STSService) assumeRoleWithWebIdentityCore(reqCtx *request.RequestContext, in WireInput) (*AssumeRoleWebIdentityResult, error) {
	if os.Getenv("TEST_MODE") != "true" {
		return nil, ErrIDPCommunicationError
	}

	roleArn := request.GetStringParam(in.Parameters, "RoleArn")
	roleSessionName := request.GetStringParam(in.Parameters, "RoleSessionName")
	webIdentityToken := request.GetStringParam(in.Parameters, "WebIdentityToken")
	providerId := request.GetStringParam(in.Parameters, "ProviderId")
	durationSeconds := request.GetIntParam(in.Parameters, "DurationSeconds")
	sessionPolicy := request.GetStringParam(in.Parameters, "Policy")
	sourceIdentity := request.GetStringParam(in.Parameters, "SourceIdentity")

	validDuration, err := validateDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if roleArn == "" {
		return nil, ErrInvalidRoleArn
	}
	if err := validateRoleArn(roleArn); err != nil {
		return nil, err
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
	// clientTokenType Smithy trait: length 4-20000 counted in Unicode
	// characters (no pattern).
	if n := utf8.RuneCountInString(webIdentityToken); n < 4 || n > 20000 {
		return nil, ErrInvalidWebIdentityToken
	}

	if err := validateSessionPolicy(sessionPolicy); err != nil {
		return nil, err
	}
	if err := validateProviderID(providerId); err != nil {
		return nil, err
	}

	// Soft JWT expiry check: when the web identity token is a parseable
	// JWT with an exp claim, reject expired tokens with
	// ExpiredTokenException per the Smithy error mapping. Tokens that
	// are not parseable JWTs (dummy tokens in TEST_MODE) are accepted.
	if isJWTExpired(webIdentityToken) {
		return nil, ErrExpiredToken
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

	// Enforce role.MaxSessionDuration after role resolution.
	err = enforceRoleMaxSessionDuration(validDuration, role.MaxSessionDuration)
	if err != nil {
		return nil, err
	}

	policyArns, err := extractPolicyArns(in.Parameters)
	if err != nil {
		return nil, err
	}
	packedPolicySize := computePackedPolicySize(sessionPolicy, policyArns, nil, nil)
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
		Policy:          sessionPolicy,
		PolicyArns:      policyArns,
	})
	if err != nil {
		return nil, err
	}

	// SubjectFromWebIdentityToken should be the OIDC sub claim of the
	// caller-supplied token.  Audience should be the aud claim.  Both
	// fall back to legacy defaults when the token is not a parseable JWT
	// (e.g. the SDK test dummy 'dummy-web-identity-token').
	subject := extractJWTClaim(webIdentityToken, "sub")
	if subject == "" {
		subject = roleSessionName
	}
	subject = fitWebIdentitySubject(subject)
	audience := extractJWTClaim(webIdentityToken, "aud")
	if audience == "" {
		audience = "sts.amazonaws.com"
	}

	// Per the Smithy AssumeRoleWithWebIdentityResponse.Provider
	// documentation: "For OpenID Connect ID tokens, this contains the
	// value of the iss field.  For OAuth 2.0 access tokens, this
	// contains the value of the ProviderId parameter."  When the token
	// is a parseable JWT with an iss claim, prefer it; otherwise fall
	// back to the caller-supplied ProviderId.
	provider := providerId
	if iss := extractJWTClaim(webIdentityToken, "iss"); iss != "" {
		provider = iss
	}

	return &AssumeRoleWebIdentityResult{
		Credentials:      credentialsOf(session),
		RoleID:           role.ID,
		RoleName:         role.RoleName,
		RoleSessionName:  roleSessionName,
		SourceIdentity:   session.SourceIdentity,
		Subject:          subject,
		Audience:         audience,
		Provider:         provider,
		PackedPolicySize: packedPolicySize,
	}, nil
}

// assumeRootCore owns the AssumeRoot flow. AWS requires the caller to be an
// IAM user or role in the Organizations management account (or an IAM
// delegated administrator) with an explicit sts:AssumeRoot grant, and
// explicitly forbids calling AssumeRoot with root user credentials.
// VorpalStacks does not implement Organizations (see docs/services.md "No
// organisations integration"), so the caller-side check reduces to the
// standard IAM policy evaluation for sts:AssumeRoot plus the root-caller
// rejection below. The session itself is scoped by the task policy: per the
// AWS AssumeRoot contract, TaskPolicyArn restricts the temporary credentials
// to the privileged task's action set instead of granting unrestricted root.
func (s *STSService) assumeRootCore(reqCtx *request.RequestContext, in WireInput) (*CredentialsResult, error) {
	// AWS: "You cannot use root user credentials to call sts:AssumeRoot."
	if reqCtx.Principal == iam.RootUserName {
		return nil, ErrAccessDenied
	}

	durationSeconds := request.GetIntParam(in.Parameters, "DurationSeconds")
	targetPrincipal := request.GetStringParam(in.Parameters, "TargetPrincipal")
	taskPolicyArn := request.GetStringParam(in.Parameters, "TaskPolicyArn.arn")

	validDuration, err := validateRootDurationSeconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	if targetPrincipal == "" {
		return nil, ErrTargetPrincipalRequired
	}
	// TargetPrincipalType Smithy trait: length 12-2048 counted in Unicode
	// characters (no pattern). Accepts account ID (12 digits) or principal
	// ARN.
	if n := utf8.RuneCountInString(targetPrincipal); n < 12 || n > 2048 {
		return nil, ErrInvalidTargetPrincipal
	}

	if taskPolicyArn == "" {
		return nil, ErrTaskPolicyArnRequired
	}
	policyName := extractPolicyNameFromArn(taskPolicyArn)
	if !allowedRootTaskPolicyNames[policyName] {
		return nil, ErrInvalidTaskPolicyArn
	}
	taskPolicy, ok := rootTaskPolicyDocuments[policyName]
	if !ok {
		return nil, ErrInvalidTaskPolicyArn
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Store TargetPrincipal as PrincipalName (the member account being
	// accessed) with an empty RoleSessionName — root sessions do not use
	// the assumed-role session name slot. The task policy document rides
	// along as the session's inline policy so the authoriser evaluates
	// every request against the task scope.
	session, err := store.Create(stsstore.CreateSessionParams{
		PrincipalType:   "Root",
		PrincipalName:   targetPrincipal,
		PrincipalArn:    arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root(),
		DurationSeconds: validDuration,
		Policy:          taskPolicy,
	})
	if err != nil {
		return nil, err
	}

	creds := credentialsOf(session)
	return &creds, nil
}
