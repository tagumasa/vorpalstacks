package sts

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"vorpalstacks/internal/config"
	pb "vorpalstacks/internal/pb/aws/sts"
	stsconnect "vorpalstacks/internal/pb/aws/sts/stsconnect"
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

// GetCallerIdentity returns the account ID, user ID, and ARN of the calling principal.
func (h *AdminHandler) GetCallerIdentity(ctx context.Context, req *connect.Request[pb.GetCallerIdentityRequest]) (*connect.Response[pb.GetCallerIdentityResponse], error) {
	accountID := h.accountID
	if accountID == "" {
		accountID = config.AWSAccountID()
	}
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
