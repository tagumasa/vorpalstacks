package sts

import (
	"context"
	"encoding/base64"
	"unicode/utf8"

	"vorpalstacks/internal/common/request"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

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

// DecodeAuthorizationMessage decodes additional information about the authorization status of a request from an encoded message.
func (s *STSService) DecodeAuthorizationMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	encodedMessage := request.GetStringParam(req.Parameters, "EncodedMessage")

	if encodedMessage == "" {
		return nil, ErrInvalidEncodedMessage
	}
	// encodedMessageType Smithy trait: length 1-10240 counted in Unicode
	// characters (no pattern).
	if utf8.RuneCountInString(encodedMessage) > 10240 {
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

	// Verify the access key exists as either a permanent IAM key or a
	// temporary STS session key. AWS returns InvalidClientTokenId for
	// non-existent access key IDs. Infrastructure failures (storage
	// unavailable) must surface as InternalError, not auth errors.
	iamStore, err := s.iamStore(reqCtx)
	if err != nil {
		return nil, ErrInternalError
	}
	if _, err := iamStore.AccessKeys().Get(accessKeyId); err == nil {
		return map[string]interface{}{
			"Account": reqCtx.GetAccountID(),
		}, nil
	}
	sessionStore, err := s.store(reqCtx)
	if err != nil {
		return nil, ErrInternalError
	}
	if _, err := sessionStore.GetByAccessKeyId(accessKeyId); err == nil {
		return map[string]interface{}{
			"Account": reqCtx.GetAccountID(),
		}, nil
	}
	return nil, ErrInvalidClientTokenId
}
