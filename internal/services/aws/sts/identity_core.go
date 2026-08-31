package sts

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"
	"unicode/utf8"

	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
	stsstore "vorpalstacks/internal/store/aws/sts"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// WireInput carries the raw wire request material every STS operation Core
// needs: the flat query-parameter map and the caller credential headers
// (already extracted by the handler). Core functions parse and validate
// members in the AWS-documented order and own every store call; the handler
// only serialises the result.
type WireInput struct {
	Parameters    map[string]interface{}
	AccessKeyID   string // X-Amz-Access-Key, falling back to the Authorization header
	SecurityToken string // X-Amz-Security-Token
}

// CredentialsResult mirrors the wire Credentials member shared by every
// session-issuing STS operation. The handler serialises the expiration via
// timeutils.ISO8601SimpleFormat.
type CredentialsResult struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

// CallerIdentityResult is the GetCallerIdentity response data: the final
// user ID and ARN after session, access-key and root fallback resolution.
type CallerIdentityResult struct {
	UserId string
	Arn    string
}

// credentialsOf copies a created session's credential material into the
// serialisation-facing result shape.
func credentialsOf(session *stsstore.Session) CredentialsResult {
	return CredentialsResult{
		AccessKeyId:     session.AccessKeyId,
		SecretAccessKey: session.SecretAccessKey,
		SessionToken:    session.SessionToken,
		Expiration:      session.Expiration,
	}
}

// iamStore returns the IAM store for the caller's account. STS is an IAM
// sub-service and resolves roles, access keys and MFA devices directly
// against it.
func (s *STSService) iamStore(reqCtx *request.RequestContext) (iamstore.IAMStoreInterface, error) {
	storage, err := reqCtx.GetGlobalStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to get global storage: %w", err)
	}
	return iamstore.GetOrCreateGlobalStore(storage, reqCtx.GetAccountID()), nil
}

// resolveCallerSession looks up the caller's STS session when the request
// includes a security token (temporary credentials). Returns nil for
// permanent IAM credentials or when no session is found.
func (s *STSService) resolveCallerSession(reqCtx *request.RequestContext, securityToken string) *stsstore.Session {
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

// resolveCallerIdentity resolves the caller's ARN and name from the access
// key ID extracted from the wire headers. The root configured key maps to
// the root principal by construction; permanent IAM keys resolve through
// the IAM store and temporary credentials through the STS session store.
func (s *STSService) resolveCallerIdentity(reqCtx *request.RequestContext, accessKeyId string) (arn, name string) {
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
func (s *STSService) resolveCallerArnOrReject(reqCtx *request.RequestContext, accessKeyId string) (string, string, error) {
	arn, name := s.resolveCallerIdentity(reqCtx, accessKeyId)
	if arn == "" {
		if os.Getenv("TEST_MODE") == "true" {
			return arnutil.NewARNBuilder(reqCtx.GetAccountID(), "").IAM().Root(), iam.RootUserName, nil
		}
		return "", "", ErrInvalidClientTokenId
	}
	return arn, name, nil
}

// getCallerIdentityCore resolves the calling principal for
// GetCallerIdentity. The security token is consulted first (temporary
// credentials); an unresolvable caller falls back to the root principal —
// the bootstrap-friendly behaviour the handler documents.
func (s *STSService) getCallerIdentityCore(reqCtx *request.RequestContext, in WireInput) (CallerIdentityResult, error) {
	securityToken := in.SecurityToken

	var userId, arn string

	if securityToken != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return CallerIdentityResult{}, err
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
		callerArn, callerName := s.resolveCallerIdentity(reqCtx, in.AccessKeyID)
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

	return CallerIdentityResult{UserId: userId, Arn: arn}, nil
}

// decodeAuthorizationMessageCore validates and decodes the EncodedMessage
// parameter: non-empty, at most 10240 Unicode characters, and Base64
// (standard or URL alphabet) encoded.
func (s *STSService) decodeAuthorizationMessageCore(in WireInput) (string, error) {
	encodedMessage := request.GetStringParam(in.Parameters, "EncodedMessage")

	if encodedMessage == "" {
		return "", ErrInvalidEncodedMessage
	}
	// encodedMessageType Smithy trait: length 1-10240 counted in Unicode
	// characters (no pattern).
	if utf8.RuneCountInString(encodedMessage) > 10240 {
		return "", ErrInvalidEncodedMessage
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedMessage)
	if err != nil {
		decodedBytes, err = base64.URLEncoding.DecodeString(encodedMessage)
		if err != nil {
			return "", ErrInvalidEncodedMessage
		}
	}

	return string(decodedBytes), nil
}

// getAccessKeyInfoCore verifies the AccessKeyId parameter exists as either
// a permanent IAM key or a temporary STS session key. AWS returns
// InvalidClientTokenId for non-existent access key IDs. Infrastructure
// failures (storage unavailable) surface as InternalError, not auth errors.
func (s *STSService) getAccessKeyInfoCore(reqCtx *request.RequestContext, in WireInput) error {
	accessKeyId := request.GetStringParam(in.Parameters, "AccessKeyId")

	if accessKeyId == "" {
		return ErrInvalidAccessKeyId
	}
	// accessKeyIdType Smithy trait: length 16-128, pattern [\w]*.
	if len(accessKeyId) < 16 || len(accessKeyId) > 128 {
		return ErrInvalidAccessKeyId
	}

	iamStore, err := s.iamStore(reqCtx)
	if err != nil {
		return ErrInternalError
	}
	if _, err := iamStore.AccessKeys().Get(accessKeyId); err == nil {
		return nil
	}
	sessionStore, err := s.store(reqCtx)
	if err != nil {
		return ErrInternalError
	}
	if _, err := sessionStore.GetByAccessKeyId(accessKeyId); err == nil {
		return nil
	}
	return ErrInvalidClientTokenId
}
