package sts

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/config"
	pb "vorpalstacks/internal/pb/aws/sts"
	stsconnect "vorpalstacks/internal/pb/aws/sts/stsconnect"
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
		Arn:     fmt.Sprintf("arn:aws:iam::%s:root", accountID),
		Userid:  accountID,
	}), nil
}
