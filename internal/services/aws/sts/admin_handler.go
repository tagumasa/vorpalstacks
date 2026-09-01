package sts

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/sts"
	stsconnect "vorpalstacks/internal/pb/aws/sts/stsconnect"
	"vorpalstacks/internal/server/grpcweb"
)

// AdminHandler implements the STS admin console gRPC-Web handler.
type AdminHandler struct {
	stsconnect.UnimplementedSTSServiceHandler
	accountID string
}

var _ stsconnect.STSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new STS admin handler with the given account ID.
func NewAdminHandler(accountID string) *AdminHandler {
	return &AdminHandler{accountID: accountID}
}

// GetCallerIdentity returns the account ID, user ID, and ARN of the calling
// principal.  The caller identity is extracted from the JWT claims set by
// the gRPC-Web auth interceptor (grpcweb.ClaimsFromContext).  When claims
// are absent (e.g. the auth interceptor is not wired) the handler falls
// back to returning the root account identity.
func (h *AdminHandler) GetCallerIdentity(ctx context.Context, req *connect.Request[pb.GetCallerIdentityRequest]) (*connect.Response[pb.GetCallerIdentityResponse], error) {
	username, arnClaim := "", ""
	if claims, ok := grpcweb.ClaimsFromContext(ctx); ok && claims != nil {
		username = claims.Username
		arnClaim = claims.GetCustomClaimString("arn")
	}
	account, arn, userid := consoleCallerIdentityCore(h.accountID, username, arnClaim)
	return connect.NewResponse(&pb.GetCallerIdentityResponse{
		Account: account,
		Arn:     arn,
		Userid:  userid,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Sts admin console.
func NewConnectHandler(accountID string) (string, http.Handler) {
	return stsconnect.NewSTSServiceHandler(NewAdminHandler(accountID))
}
