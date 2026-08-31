package sts

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// callerAccessKeyID extracts the caller's access key ID from the wire
// headers: the X-Amz-Access-Key header when present, otherwise the access
// key signed in the Authorization header.
func callerAccessKeyID(req *request.ParsedRequest) string {
	accessKeyId := req.Headers.Get("X-Amz-Access-Key")
	if accessKeyId == "" {
		authHeader := req.Headers.Get("Authorization")
		if authHeader != "" {
			accessKeyId = request.ExtractAccessKeyIDFromAuth(authHeader)
		}
	}
	return accessKeyId
}

// GetCallerIdentity returns details about the IAM user or role whose credentials are used to call the operation.
//
// When the caller cannot be resolved from the security token or access key
// (e.g. during InitialSetup before any IAM user exists), the response falls
// back to the root principal. This is an intentional design choice for
// VorpalStacks to support bootstrap/setup flows; AWS would return 403 for
// unauthenticated requests.
func (s *STSService) GetCallerIdentity(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.getCallerIdentityCore(reqCtx, WireInput{
		Parameters:    req.Parameters,
		AccessKeyID:   callerAccessKeyID(req),
		SecurityToken: req.Headers.Get("X-Amz-Security-Token"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"UserId":  result.UserId,
		"Account": reqCtx.GetAccountID(),
		"Arn":     result.Arn,
	}, nil
}

// DecodeAuthorizationMessage decodes additional information about the authorization status of a request from an encoded message.
func (s *STSService) DecodeAuthorizationMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	decodedMessage, err := s.decodeAuthorizationMessageCore(WireInput{Parameters: req.Parameters})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DecodedMessage": decodedMessage,
	}, nil
}

// GetAccessKeyInfo returns information about the access key in the request.
func (s *STSService) GetAccessKeyInfo(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.getAccessKeyInfoCore(reqCtx, WireInput{Parameters: req.Parameters}); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Account": reqCtx.GetAccountID(),
	}, nil
}
