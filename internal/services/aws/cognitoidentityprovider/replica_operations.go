package cognitoidentityprovider

import (
	"context"

	"vorpalstacks/internal/common/request"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// Handlers for the user-pool replica family; validation and store access
// live in replica_core.go.

// CreateUserPoolReplica creates a cross-region replica of a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_CreateUserPoolReplica.html
func (s *CognitoService) CreateUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.createUserPoolReplicaCore(reqCtx, CreateUserPoolReplicaInput{
		UserPoolID: req.GetParam("UserPoolId"),
		RegionName: req.GetParam("RegionName"),
		Params:     req.Parameters,
	})
}

// ListUserPoolReplicas lists replicas for a user pool.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_ListUserPoolReplicas.html
func (s *CognitoService) ListUserPoolReplicas(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.listUserPoolReplicasCore(reqCtx, ListUserPoolReplicasInput{
		UserPoolID: req.GetParam("UserPoolId"),
		NextToken:  req.GetParam("NextToken"),
	})
}

// DeleteUserPoolReplica deletes a cross-region replica.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_DeleteUserPoolReplica.html
func (s *CognitoService) DeleteUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.deleteUserPoolReplicaCore(reqCtx, DeleteUserPoolReplicaInput{
		UserPoolID: req.GetParam("UserPoolId"),
		RegionName: req.GetParam("RegionName"),
	})
}

// UpdateUserPoolReplica updates a cross-region replica.
// https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_UpdateUserPoolReplica.html
func (s *CognitoService) UpdateUserPoolReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateUserPoolReplicaCore(reqCtx, UpdateUserPoolReplicaInput{
		UserPoolID: req.GetParam("UserPoolId"),
		RegionName: req.GetParam("RegionName"),
		Status:     req.GetParam("Status"),
	})
}

func formatUserPoolReplica(r *cognitostore.UserPoolReplica) map[string]interface{} {
	result := map[string]interface{}{
		"RegionName": r.RegionName,
		"Status":     r.Status,
		"Role":       r.Role,
	}
	if r.UserPoolArn != "" {
		result["UserPoolArn"] = r.UserPoolArn
	}
	return result
}
