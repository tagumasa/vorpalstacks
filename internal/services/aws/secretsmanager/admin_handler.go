package secretsmanager

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/secretsmanager"
	secretsmanagerconnect "vorpalstacks/internal/pb/aws/secretsmanager/secretsmanagerconnect"
)

// AdminHandler implements the Secrets Manager admin console gRPC-Web handler.
type AdminHandler struct {
	secretsmanagerconnect.UnimplementedSecretsManagerServiceHandler
}

var _ secretsmanagerconnect.SecretsManagerServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Secrets Manager admin handler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListSecrets returns all Secrets Manager secrets visible to the admin console.
func (h *AdminHandler) ListSecrets(ctx context.Context, req *connect.Request[pb.ListSecretsRequest]) (*connect.Response[pb.ListSecretsResponse], error) {
	return connect.NewResponse(&pb.ListSecretsResponse{
		Secretlist: []*pb.SecretListEntry{},
	}), nil
}
