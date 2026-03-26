package sts

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/config"
	pb "vorpalstacks/internal/pb/aws/sts"
	stsconnect "vorpalstacks/internal/pb/aws/sts/stsconnect"
)

// AdminHandler handles STS gRPC requests for the admin interface.
type AdminHandler struct {
	accountID string
}

var _ stsconnect.STSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new STS admin handler with the given account ID.
func NewAdminHandler(accountID string) *AdminHandler {
	return &AdminHandler{accountID: accountID}
}

// GetCallerIdentity returns details about the IAM identity making the request.
func (h *AdminHandler) GetCallerIdentity(ctx context.Context, req *connect.Request[pb.GetCallerIdentityRequest]) (*connect.Response[pb.GetCallerIdentityResponse], error) {
	accountID := h.accountID
	if accountID == "" {
		accountID = config.AWSAccountID()
	}
	return connect.NewResponse(&pb.GetCallerIdentityResponse{
		Account: accountID,
		Arn:     fmt.Sprintf("arn:aws:iam::%s:root", accountID),
		Userid:  accountID,
	}), nil
}

// AssumeRole returns temporary security credentials for an IAM role.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssumeRole(ctx context.Context, req *connect.Request[pb.AssumeRoleRequest]) (*connect.Response[pb.AssumeRoleResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.AssumeRole is not implemented"))
}

// AssumeRoleWithSAML returns temporary credentials for an IAM role using SAML assertion.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssumeRoleWithSAML(ctx context.Context, req *connect.Request[pb.AssumeRoleWithSAMLRequest]) (*connect.Response[pb.AssumeRoleWithSAMLResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.AssumeRoleWithSAML is not implemented"))
}

// AssumeRoleWithWebIdentity returns temporary credentials for an IAM role using web identity.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssumeRoleWithWebIdentity(ctx context.Context, req *connect.Request[pb.AssumeRoleWithWebIdentityRequest]) (*connect.Response[pb.AssumeRoleWithWebIdentityResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.AssumeRoleWithWebIdentity is not implemented"))
}

// AssumeRoot returns temporary credentials for the root account.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssumeRoot(ctx context.Context, req *connect.Request[pb.AssumeRootRequest]) (*connect.Response[pb.AssumeRootResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.AssumeRoot is not implemented"))
}

// DecodeAuthorizationMessage decodes an encoded authorization message.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DecodeAuthorizationMessage(ctx context.Context, req *connect.Request[pb.DecodeAuthorizationMessageRequest]) (*connect.Response[pb.DecodeAuthorizationMessageResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.DecodeAuthorizationMessage is not implemented"))
}

// GetAccessKeyInfo returns the account identifier for an access key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAccessKeyInfo(ctx context.Context, req *connect.Request[pb.GetAccessKeyInfoRequest]) (*connect.Response[pb.GetAccessKeyInfoResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.GetAccessKeyInfo is not implemented"))
}

// GetDelegatedAccessToken returns a delegated access token.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDelegatedAccessToken(ctx context.Context, req *connect.Request[pb.GetDelegatedAccessTokenRequest]) (*connect.Response[pb.GetDelegatedAccessTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.GetDelegatedAccessToken is not implemented"))
}

// GetFederationToken returns temporary security credentials for a federated user.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetFederationToken(ctx context.Context, req *connect.Request[pb.GetFederationTokenRequest]) (*connect.Response[pb.GetFederationTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.GetFederationToken is not implemented"))
}

// GetSessionToken returns temporary credentials for an IAM user or root account.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSessionToken(ctx context.Context, req *connect.Request[pb.GetSessionTokenRequest]) (*connect.Response[pb.GetSessionTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.GetSessionToken is not implemented"))
}

// GetWebIdentityToken returns a web identity token.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetWebIdentityToken(ctx context.Context, req *connect.Request[pb.GetWebIdentityTokenRequest]) (*connect.Response[pb.GetWebIdentityTokenResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("sts.STSService.GetWebIdentityToken is not implemented"))
}
