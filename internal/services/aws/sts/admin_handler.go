package sts

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	commoniam "vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/config"
	pb "vorpalstacks/internal/pb/aws/sts"
	stsconnect "vorpalstacks/internal/pb/aws/sts/stsconnect"
	"vorpalstacks/internal/server/grpcweb"
	arnutil "vorpalstacks/internal/utils/aws/arn"
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
	accountID := h.accountID
	if accountID == "" {
		accountID = config.AWSAccountID()
	}

	if claims, ok := grpcweb.ClaimsFromContext(ctx); ok && claims != nil {
		arn := claims.GetCustomClaimString("arn")
		userID := claims.Username
		if claims.Username == commoniam.RootUserName {
			// AWS returns UserId = accountID for root callers.
			userID = accountID
			if arn == "" {
				arn = arnutil.NewARNBuilder(accountID, "").IAM().Root()
			}
		} else if arn == "" {
			arn = arnutil.NewARNBuilder(accountID, "").IAM().User(claims.Username)
		}
		return connect.NewResponse(&pb.GetCallerIdentityResponse{
			Account: accountID,
			Arn:     arn,
			Userid:  userID,
		}), nil
	}

	// Fallback: no JWT claims (auth interceptor not active).
	return connect.NewResponse(&pb.GetCallerIdentityResponse{
		Account: accountID,
		Arn:     arnutil.NewARNBuilder(accountID, "").IAM().Root(),
		Userid:  accountID,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Sts admin console.
func NewConnectHandler(accountID string) (string, http.Handler) {
	return stsconnect.NewSTSServiceHandler(NewAdminHandler(accountID))
}
