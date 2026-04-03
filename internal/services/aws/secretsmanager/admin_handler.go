package secretsmanager

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/secretsmanager"
	secretsmanagerconnect "vorpalstacks/internal/pb/aws/secretsmanager/secretsmanagerconnect"
)

type AdminHandler struct {
	secretsmanagerconnect.UnimplementedSecretsManagerServiceHandler
}

var _ secretsmanagerconnect.SecretsManagerServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) ListSecrets(ctx context.Context, req *connect.Request[pb.ListSecretsRequest]) (*connect.Response[pb.ListSecretsResponse], error) {
	return connect.NewResponse(&pb.ListSecretsResponse{
		Secretlist: []*pb.SecretListEntry{},
	}), nil
}
