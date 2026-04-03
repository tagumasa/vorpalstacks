package sts

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/config"
	pb "vorpalstacks/internal/pb/aws/sts"
	stsconnect "vorpalstacks/internal/pb/aws/sts/stsconnect"
)

type AdminHandler struct {
	stsconnect.UnimplementedSTSServiceHandler
	accountID string
}

var _ stsconnect.STSServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(accountID string) *AdminHandler {
	return &AdminHandler{accountID: accountID}
}

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
