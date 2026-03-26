package secretsmanager

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/secretsmanager"
	secretsmanagerconnect "vorpalstacks/internal/pb/aws/secretsmanager/secretsmanagerconnect"
)

// AdminHandler provides Secrets Manager service administration functionality.
// It implements the SecretsManagerServiceHandler interface for gRPC-Web communication.
type AdminHandler struct{}

var _ secretsmanagerconnect.SecretsManagerServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Secrets Manager AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListSecrets lists secrets in Secrets Manager.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSecrets(ctx context.Context, req *connect.Request[pb.ListSecretsRequest]) (*connect.Response[pb.ListSecretsResponse], error) {
	return connect.NewResponse(&pb.ListSecretsResponse{
		Secretlist: []*pb.SecretListEntry{},
	}), nil
}

// BatchGetSecretValue retrieves multiple secrets in a single request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) BatchGetSecretValue(ctx context.Context, req *connect.Request[pb.BatchGetSecretValueRequest]) (*connect.Response[pb.BatchGetSecretValueResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.BatchGetSecretValue is not implemented"))
}

// CancelRotateSecret cancels a secret rotation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CancelRotateSecret(ctx context.Context, req *connect.Request[pb.CancelRotateSecretRequest]) (*connect.Response[pb.CancelRotateSecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.CancelRotateSecret is not implemented"))
}

// CreateSecret creates a new secret in Secrets Manager.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateSecret(ctx context.Context, req *connect.Request[pb.CreateSecretRequest]) (*connect.Response[pb.CreateSecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.CreateSecret is not implemented"))
}

// DeleteResourcePolicy deletes the resource policy from a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteResourcePolicy(ctx context.Context, req *connect.Request[pb.DeleteResourcePolicyRequest]) (*connect.Response[pb.DeleteResourcePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.DeleteResourcePolicy is not implemented"))
}

// DeleteSecret deletes a secret from Secrets Manager.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteSecret(ctx context.Context, req *connect.Request[pb.DeleteSecretRequest]) (*connect.Response[pb.DeleteSecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.DeleteSecret is not implemented"))
}

// DescribeSecret returns metadata about a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeSecret(ctx context.Context, req *connect.Request[pb.DescribeSecretRequest]) (*connect.Response[pb.DescribeSecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.DescribeSecret is not implemented"))
}

// GetRandomPassword generates a random password.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRandomPassword(ctx context.Context, req *connect.Request[pb.GetRandomPasswordRequest]) (*connect.Response[pb.GetRandomPasswordResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.GetRandomPassword is not implemented"))
}

// GetResourcePolicy returns the resource policy for a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetResourcePolicy(ctx context.Context, req *connect.Request[pb.GetResourcePolicyRequest]) (*connect.Response[pb.GetResourcePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.GetResourcePolicy is not implemented"))
}

// GetSecretValue retrieves the secret value for a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSecretValue(ctx context.Context, req *connect.Request[pb.GetSecretValueRequest]) (*connect.Response[pb.GetSecretValueResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.GetSecretValue is not implemented"))
}

// ListSecretVersionIds lists versions of a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListSecretVersionIds(ctx context.Context, req *connect.Request[pb.ListSecretVersionIdsRequest]) (*connect.Response[pb.ListSecretVersionIdsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.ListSecretVersionIds is not implemented"))
}

// PutResourcePolicy attaches a resource policy to a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutResourcePolicy(ctx context.Context, req *connect.Request[pb.PutResourcePolicyRequest]) (*connect.Response[pb.PutResourcePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.PutResourcePolicy is not implemented"))
}

// PutSecretValue stores a new secret value for a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutSecretValue(ctx context.Context, req *connect.Request[pb.PutSecretValueRequest]) (*connect.Response[pb.PutSecretValueResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.PutSecretValue is not implemented"))
}

// RemoveRegionsFromReplication removes regions from secret replication.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RemoveRegionsFromReplication(ctx context.Context, req *connect.Request[pb.RemoveRegionsFromReplicationRequest]) (*connect.Response[pb.RemoveRegionsFromReplicationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.RemoveRegionsFromReplication is not implemented"))
}

// ReplicateSecretToRegions replicates a secret to specified regions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ReplicateSecretToRegions(ctx context.Context, req *connect.Request[pb.ReplicateSecretToRegionsRequest]) (*connect.Response[pb.ReplicateSecretToRegionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.ReplicateSecretToRegions is not implemented"))
}

// RestoreSecret restores a previously deleted secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RestoreSecret(ctx context.Context, req *connect.Request[pb.RestoreSecretRequest]) (*connect.Response[pb.RestoreSecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.RestoreSecret is not implemented"))
}

// RotateSecret starts or modifies secret rotation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RotateSecret(ctx context.Context, req *connect.Request[pb.RotateSecretRequest]) (*connect.Response[pb.RotateSecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.RotateSecret is not implemented"))
}

// StopReplicationToReplica removes replication from a replica secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) StopReplicationToReplica(ctx context.Context, req *connect.Request[pb.StopReplicationToReplicaRequest]) (*connect.Response[pb.StopReplicationToReplicaResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.StopReplicationToReplica is not implemented"))
}

// TagResource adds tags to a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.TagResource is not implemented"))
}

// UntagResource removes tags from a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.UntagResource is not implemented"))
}

// UpdateSecret updates the secret value for a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateSecret(ctx context.Context, req *connect.Request[pb.UpdateSecretRequest]) (*connect.Response[pb.UpdateSecretResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.UpdateSecret is not implemented"))
}

// UpdateSecretVersionStage modifies the version stage for a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateSecretVersionStage(ctx context.Context, req *connect.Request[pb.UpdateSecretVersionStageRequest]) (*connect.Response[pb.UpdateSecretVersionStageResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.UpdateSecretVersionStage is not implemented"))
}

// ValidateResourcePolicy validates a resource policy for a secret.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ValidateResourcePolicy(ctx context.Context, req *connect.Request[pb.ValidateResourcePolicyRequest]) (*connect.Response[pb.ValidateResourcePolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("secretsmanager.SecretsManagerService.ValidateResourcePolicy is not implemented"))
}
