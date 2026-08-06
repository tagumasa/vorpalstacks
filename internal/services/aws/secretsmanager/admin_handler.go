package secretsmanager

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	svcerrors "vorpalstacks/internal/common/errors"
	pb "vorpalstacks/internal/pb/aws/secretsmanager"
	secretsmanagerconnect "vorpalstacks/internal/pb/aws/secretsmanager/secretsmanagerconnect"
)

// AdminHandler implements the Secrets Manager admin console gRPC-Web handler.
// It is a thin adapter that converts proto requests to transport-agnostic
// Input structs, delegates to core functions, and converts results back to
// proto responses. No store package is imported directly (AGENTS.md #29).
type AdminHandler struct {
	secretsmanagerconnect.UnimplementedSecretsManagerServiceHandler
	service *SecretsManagerService
}

var _ secretsmanagerconnect.SecretsManagerServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Secrets Manager admin handler backed by the
// given service instance.
func NewAdminHandler(svc *SecretsManagerService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListSecrets returns all Secrets Manager secrets visible to the admin console.
func (h *AdminHandler) ListSecrets(ctx context.Context, req *connect.Request[pb.ListSecretsRequest]) (*connect.Response[pb.ListSecretsResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.listSecretsCore(ctx, store, pbToListSecretsInput(req.Msg))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	entries := make([]*pb.SecretListEntry, 0, len(result.Secrets))
	for _, s := range result.Secrets {
		entries = append(entries, secretToPbEntry(s))
	}

	return connect.NewResponse(&pb.ListSecretsResponse{
		Secretlist: entries,
		Nexttoken:  result.NextToken,
	}), nil
}

// CreateSecret creates a new secret via the admin console.
func (h *AdminHandler) CreateSecret(ctx context.Context, req *connect.Request[pb.CreateSecretRequest]) (*connect.Response[pb.CreateSecretResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.createSecretCore(ctx, store, pbToCreateSecretInput(req.Msg))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(createSecretResultToPb(result)), nil
}

// DeleteSecret deletes a secret via the admin console.
func (h *AdminHandler) DeleteSecret(ctx context.Context, req *connect.Request[pb.DeleteSecretRequest]) (*connect.Response[pb.DeleteSecretResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.deleteSecretCore(ctx, store, pbToDeleteSecretInput(req.Msg))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(deleteSecretResultToPb(result)), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Secrets
// Manager admin console.
func NewConnectHandler(svc *SecretsManagerService) (string, http.Handler) {
	return secretsmanagerconnect.NewSecretsManagerServiceHandler(NewAdminHandler(svc))
}
