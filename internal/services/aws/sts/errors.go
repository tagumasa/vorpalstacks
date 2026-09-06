// Package sts provides STS (Security Token Service) operations for vorpalstacks.
package sts

import (
	"net/http"

	"vorpalstacks/internal/common/errors"
)

// Duration constants for STS sessions.
const (
	// MinDurationSeconds is the minimum allowed session duration in seconds.
	MinDurationSeconds = 900
	// MaxDurationSeconds is the maximum allowed session duration in seconds for standard operations.
	MaxDurationSeconds = 43200
	// MaxDurationSecondsExtended is the maximum allowed session duration in seconds for extended operations.
	MaxDurationSecondsExtended = 129600
	// DefaultDurationSeconds is the default session duration in seconds.
	DefaultDurationSeconds = 3600
	// MaxRootDurationSeconds is the maximum allowed session duration for
	// AssumeRoot per the Smithy RootDurationSecondsType trait.
	MaxRootDurationSeconds = 900
	// DefaultRootDurationSeconds is the default session duration for AssumeRoot.
	DefaultRootDurationSeconds = 900
	// DefaultFederationDurationSeconds is the default session duration for
	// GetFederationToken (AWS spec: 43200 seconds / 12 hours).
	DefaultFederationDurationSeconds = 43200
)

var (
	// ErrInvalidRoleArn is returned when the specified role ARN is empty
	// or does not correspond to a known role.
	ErrInvalidRoleArn = errors.NewAWSError("InvalidRoleArn", "The specified role ARN is invalid.", http.StatusBadRequest)
	// ErrInvalidRoleArnFormat is returned when RoleArn fails the Smithy
	// arnType constraint (length 20-2048, arnType pattern).  AWS returns
	// ValidationError for constraint violations before attempting role
	// resolution.
	ErrInvalidRoleArnFormat = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'roleArn' failed to satisfy constraint: Member must have length greater than or equal to 20", http.StatusBadRequest)
	// ErrInvalidPrincipalArnFormat is returned when PrincipalArn in
	// AssumeRoleWithSAML fails the Smithy arnType constraint.
	ErrInvalidPrincipalArnFormat = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'principalArn' failed to satisfy constraint: Member must have length greater than or equal to 20", http.StatusBadRequest)
	// ErrInvalidExternalId is returned when the ExternalId parameter fails AWS
	// spec validation (length 2-1224, pattern [\w+=,.@:\/-]*).
	ErrInvalidExternalId = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'externalId' failed to satisfy constraint: Member must satisfy regular expression pattern: [\\w+=,.@:\\/-]*", http.StatusBadRequest)
	// ErrInvalidPrincipalArn is returned when the specified principal ARN is invalid.
	ErrInvalidPrincipalArn = errors.NewAWSError("InvalidPrincipalArn", "The specified principal ARN is invalid.", http.StatusBadRequest)
	// ErrInvalidSAMLAssertion is returned when the SAML assertion is invalid.
	// The Smithy aws.protocols#awsQueryError trait for AssumeRoleWithSAML's
	// identity-token validation maps to "InvalidIdentityToken" (the same
	// code used for invalid web identity tokens). Reusing this code keeps
	// AWS SDK clients from having to special-case SAML vs WebIdentity
	// response failures.
	ErrInvalidSAMLAssertion = errors.NewAWSError("InvalidIdentityToken", "The SAML assertion is invalid.", http.StatusBadRequest)
	// ErrInvalidWebIdentityToken is returned when the web identity token is invalid.
	ErrInvalidWebIdentityToken = errors.NewAWSError("InvalidIdentityToken", "The web identity token is invalid.", http.StatusBadRequest)
	// ErrInvalidProviderID is returned when ProviderId fails the Smithy
	// urlType length constraint (4-2048).  AWS returns ValidationError
	// for this parameter, distinct from InvalidIdentityToken which is
	// reserved for token validation failures.
	ErrInvalidProviderID = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'providerId' failed to satisfy constraint: Member must have length greater than or equal to 4", http.StatusBadRequest)
	// ErrInvalidEncodedMessage is returned when the encoded message is invalid.
	// The Smithy aws.protocols#awsQueryError trait value for the
	// InvalidAuthorizationMessageException shape is "InvalidAuthorizationMessageException";
	// emitting that exact code keeps AWS SDK clients from receiving an
	// unrecognised error type when decoding fails.
	ErrInvalidEncodedMessage = errors.NewAWSError("InvalidAuthorizationMessageException", "The encoded message is invalid.", http.StatusBadRequest)
	// ErrInvalidAccessKeyId is returned when the access key ID is invalid.
	ErrInvalidAccessKeyId = errors.NewAWSError("InvalidAccessKeyId", "The access key ID is invalid.", http.StatusBadRequest)
	// ErrInvalidClientTokenId is returned by GetAccessKeyInfo when the
	// supplied AccessKeyId does not correspond to any IAM permanent key
	// or STS temporary session in the store.
	ErrInvalidClientTokenId = errors.NewAWSError("InvalidClientTokenId", "The security token included in the request is invalid.", http.StatusForbidden)
	// ErrInvalidFederationName is returned when the federation name is invalid.
	ErrInvalidFederationName = errors.NewAWSError("InvalidInput", "The federation name is invalid.", http.StatusBadRequest)
	// ErrInvalidTradeInToken is returned when the trade-in token is invalid.
	ErrInvalidTradeInToken = errors.NewAWSError("InvalidToken", "The trade-in token is invalid.", http.StatusBadRequest)
	// ErrInvalidDuration is returned when the duration is outside the allowed range.
	ErrInvalidDuration = errors.NewAWSError("InvalidDuration", "The duration seconds must be between 900 and 43200 seconds.", http.StatusBadRequest)
	// ErrInvalidDurationExtended is returned when the duration is outside the extended allowed range.
	ErrInvalidDurationExtended = errors.NewAWSError("InvalidDuration", "The duration seconds must be between 900 and 129600 seconds.", http.StatusBadRequest)
	// ErrNoSuchRole is returned when the specified role cannot be found.
	ErrNoSuchRole = errors.NewAWSError("NoSuchEntity", "The role with the specified ARN cannot be found.", http.StatusNotFound)
	// ErrMalformedPolicyDocument is returned when the policy document is invalid.
	ErrMalformedPolicyDocument = errors.NewAWSError("MalformedPolicyDocument", "The policy document is invalid.", http.StatusBadRequest)
	// ErrInvalidPolicyArn is returned when a PolicyDescriptorType.arn
	// fails the Smithy arnType constraint (length 20-2048, arnType
	// pattern).  Used only for PolicyArns.member.N.arn validation.
	ErrInvalidPolicyArn = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'policyArn' failed to satisfy constraint: Member must satisfy regular expression pattern: arnType", http.StatusBadRequest)
	// ErrInvalidProviderArn is returned when a ProvidedContext
	// ProviderArn fails the Smithy arnType constraint.  Separate from
	// ErrInvalidPolicyArn so the error message references the correct
	// field name ("ProviderArn", not "policyArn").
	ErrInvalidProviderArn = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'providerArn' failed to satisfy constraint: Member must satisfy regular expression pattern: arnType", http.StatusBadRequest)
	// ErrTooManyPolicyArns is returned when the PolicyArns list exceeds
	// the AWS-documented maximum of 10 managed policy ARNs.
	ErrTooManyPolicyArns = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'policyArns' failed to satisfy constraint: Member must have length less than or equal to 10", http.StatusBadRequest)
	// ErrDuplicatePolicyArn is returned when the same policy ARN appears
	// more than once in the PolicyArns parameter. Duplicate managed
	// policies are semantically redundant; the dedup check mirrors the
	// ErrDuplicateSessionTagKey guard used for session tags.
	ErrDuplicatePolicyArn = errors.NewAWSError("ValidationError", "Duplicate policy ARNs are not allowed.", http.StatusBadRequest)
	// ErrInvalidContextAssertion is returned when a ProvidedContext
	// ContextAssertion fails the Smithy contextAssertionType length
	// constraint (4-2048).
	ErrInvalidContextAssertion = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'contextAssertion' failed to satisfy constraint: Member must have length 4-2048", http.StatusBadRequest)
	// ErrTooManyProvidedContexts is returned when the ProvidedContexts
	// list exceeds the Smithy ProvidedContextsListType maximum of 5.
	ErrTooManyProvidedContexts = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'providedContexts' failed to satisfy constraint: Member must have length less than or equal to 5", http.StatusBadRequest)
	// ErrInvalidParameter is returned when a parameter value is invalid.
	ErrInvalidParameter = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'roleSessionName' failed to satisfy constraint: Member must satisfy regular expression pattern: [\\w+=,.@-]*", http.StatusBadRequest)
	// ErrAccessDenied is returned when the caller is not authorized to assume the role.
	ErrAccessDenied = errors.NewAWSError("AccessDenied", "Not authorized to perform sts:AssumeRole", http.StatusForbidden)
	// ErrExpiredTradeInToken is returned when the trade-in token has expired.
	ErrExpiredTradeInToken = errors.NewAWSError("ExpiredTradeInToken", "The trade-in token provided has expired.", http.StatusBadRequest)
	// ErrIDPCommunicationError is returned when external IdP validation is unavailable.
	// VorpalStacks cannot reach external IdPs; SAML and WebIdentity flows are test-only.
	ErrIDPCommunicationError = errors.NewAWSError("IDPCommunicationError", "The identity provider is not available. SAML and WebIdentity federation require external IdP connectivity.", http.StatusBadRequest)
	// ErrInvalidRootDuration is returned when AssumeRoot DurationSeconds falls
	// outside the AWS RootDurationSecondsType range (0-900).
	ErrInvalidRootDuration = errors.NewAWSError("ValidationError", "DurationSeconds must be between 0 and 900 for AssumeRoot.", http.StatusBadRequest)
	// ErrTargetPrincipalRequired is returned when AssumeRoot is called without
	// the required TargetPrincipal parameter.
	ErrTargetPrincipalRequired = errors.NewAWSError("ValidationError", "TargetPrincipal is required.", http.StatusBadRequest)
	// ErrInvalidTargetPrincipal is returned when TargetPrincipal fails the
	// Smithy TargetPrincipalType length constraint (12-2048).
	ErrInvalidTargetPrincipal = errors.NewAWSError("ValidationError", "TargetPrincipal must be between 12 and 2048 characters.", http.StatusBadRequest)
	// ErrTaskPolicyArnRequired is returned when AssumeRoot is called without
	// the required TaskPolicyArn parameter.
	ErrTaskPolicyArnRequired = errors.NewAWSError("ValidationError", "TaskPolicyArn is required.", http.StatusBadRequest)
	// ErrInvalidTaskPolicyArn is returned when TaskPolicyArn does not reference
	// one of the AWS root-task managed policies.
	ErrInvalidTaskPolicyArn = errors.NewAWSError("ValidationError", "TaskPolicyArn.arn must reference one of the root-task managed policies.", http.StatusBadRequest)
	// ErrInvalidMFARequirements is returned when only one of SerialNumber or
	// TokenCode is supplied. Both must be present together or both absent.
	ErrInvalidMFARequirements = errors.NewAWSError("ValidationError", "SerialNumber and TokenCode must both be supplied or both omitted.", http.StatusBadRequest)
	// ErrInvalidMFASerialNumber is returned when SerialNumber fails the
	// Smithy serialNumberType constraint (length 9-256, pattern
	// [\w+=/:,.@-]*).
	ErrInvalidMFASerialNumber = errors.NewAWSError("ValidationError", "The SerialNumber parameter is invalid.", http.StatusBadRequest)
	// ErrInvalidMFATokenCode is returned when TokenCode is not exactly six
	// digits per the Smithy tokenCodeType trait.
	ErrInvalidMFATokenCode = errors.NewAWSError("ValidationError", "The TokenCode parameter must be exactly six digits.", http.StatusBadRequest)
	// ErrMFADeviceNotFound is returned when the SerialNumber does not match a
	// registered IAM MFA device for the caller.
	ErrMFADeviceNotFound = errors.NewAWSError("InvalidClientTokenId", "The MFA device specified by SerialNumber does not exist or is not assigned to the calling user.", http.StatusBadRequest)
	// ErrInvalidMFACode is returned when the TokenCode does not match the
	// expected TOTP value for the specified MFA device.
	ErrInvalidMFACode = errors.NewAWSError("ValidationError", "The TokenCode provided does not match the expected value for the MFA device.", http.StatusBadRequest)
	// ErrInvalidWebIdentityDuration is returned when GetWebIdentityToken
	// DurationSeconds falls outside the Smithy
	// webIdentityTokenDurationSecondsType range (60-3600).
	ErrInvalidWebIdentityDuration = errors.NewAWSError("ValidationError", "DurationSeconds must be between 60 and 3600 for GetWebIdentityToken.", http.StatusBadRequest)
	// ErrAudienceRequired is returned when GetWebIdentityToken is called
	// without the required Audience parameter.
	ErrAudienceRequired = errors.NewAWSError("ValidationError", "Audience is required.", http.StatusBadRequest)
	// ErrTooManyAudiences is returned when GetWebIdentityToken Audience list
	// exceeds the Smithy maximum of 10 entries.
	ErrTooManyAudiences = errors.NewAWSError("ValidationError", "Too many audiences: maximum 10 allowed.", http.StatusBadRequest)
	// ErrInvalidSessionTag is returned when a session tag key or value fails
	// the Smithy tagKeyType or tagValueType trait validation.
	ErrInvalidSessionTag = errors.NewAWSError("ValidationError", "Invalid session tag key or value.", http.StatusBadRequest)
	// ErrTooManySessionTags is returned when more than 50 session tags are
	// provided, exceeding the Smithy tagListType length constraint.
	ErrTooManySessionTags = errors.NewAWSError("ValidationError", "Too many session tags: maximum 50 allowed.", http.StatusBadRequest)
	// ErrDuplicateSessionTagKey is returned when the same tag key appears more
	// than once in the Tags parameter. AWS rejects duplicate session tag keys.
	ErrDuplicateSessionTagKey = errors.NewAWSError("ValidationError", "Duplicate session tag keys are not allowed.", http.StatusBadRequest)
	// ErrTransitiveKeyWithoutTag is returned when a TransitiveTagKeys
	// entry has no corresponding value in the session tags (neither in
	// the request Tags nor inherited from the caller session). AWS docs:
	// "If you set a tag key as transitive, the corresponding key and
	// value passes to subsequent sessions."
	ErrTransitiveKeyWithoutTag = errors.NewAWSError("ValidationError", "Transitive tag key has no corresponding session tag.", http.StatusBadRequest)
	// ErrInvalidSigningAlgorithm is returned when GetWebIdentityToken
	// SigningAlgorithm is not RS256 or ES384.
	ErrInvalidSigningAlgorithm = errors.NewAWSError("ValidationError", "SigningAlgorithm must be RS256 or ES384.", http.StatusBadRequest)
	// ErrInvalidFederationRootDuration is returned when GetFederationToken is
	// called by the root user with a DurationSeconds exceeding the 3600-second
	// cap imposed by AWS.
	ErrInvalidFederationRootDuration = errors.NewAWSError("ValidationError", "DurationSeconds must not exceed 3600 when called with root user credentials.", http.StatusBadRequest)
	// ErrInvalidRootSessionDuration is returned when GetSessionToken is called
	// by the root user with a DurationSeconds exceeding the 3600-second cap
	// imposed by AWS.
	ErrInvalidRootSessionDuration = errors.NewAWSError("ValidationError", "DurationSeconds must not exceed 3600 when called with root user credentials.", http.StatusBadRequest)
	// ErrPackedPolicyTooLarge is returned when the combined packed size of
	// session policies and session tags exceeds the AWS-allowed limit.
	ErrPackedPolicyTooLarge = errors.NewAWSError("PackedPolicyTooLarge", "The total packed size of the session policies and session tags is too large.", http.StatusBadRequest)
	// ErrInvalidSourceIdentity is returned when SourceIdentity fails the
	// Smithy sourceIdentityType constraint (length 2-64, pattern
	// [\w+=,.@-]*) or begins with the reserved "aws:" prefix.
	ErrInvalidSourceIdentity = errors.NewAWSError("ValidationError", "1 validation error detected: Value at 'sourceIdentity' failed to satisfy constraint: Member must satisfy regular expression pattern: [\\w+=,.@-]*", http.StatusBadRequest)
	// ErrJWTPayloadSizeExceeded is returned by GetWebIdentityToken when
	// the generated JWT payload exceeds the AWS-allowed size limit.
	// Smithy JWTPayloadSizeExceededException, awsQueryError code
	// "JWTPayloadSizeExceededException", httpResponseCode 400.
	ErrJWTPayloadSizeExceeded = errors.NewAWSError("JWTPayloadSizeExceededException", "The payload size of the JWT token exceeds the allowed size.", http.StatusBadRequest)
	// ErrInternalError is returned when the service encounters an
	// infrastructure failure (e.g. storage unavailable). This must
	// not be masked as an authentication error. InternalFailure is
	// the only 500 code the STS contract documents (Common Error
	// Types).
	ErrInternalError = errors.NewAWSError("InternalFailure", "The request can't be processed right now because of an internal server issue.", http.StatusInternalServerError)
	// ErrExpiredToken is returned when the caller's temporary credentials
	// or the supplied web identity / SAML token has expired. Smithy
	// awsQueryError code "ExpiredTokenException", httpResponseCode 400.
	ErrExpiredToken = errors.NewAWSError("ExpiredTokenException", "The security token included in the request is expired.", http.StatusBadRequest)
	// ErrSessionDurationEscalation is returned by GetWebIdentityToken when
	// the requested DurationSeconds would extend beyond the caller's own
	// session expiration. Smithy awsQueryError code
	// "SessionDurationEscalationException", httpResponseCode 403.
	ErrSessionDurationEscalation = errors.NewAWSError("SessionDurationEscalationException", "The requested token duration would extend the session beyond its original expiration time.", http.StatusForbidden)
)
