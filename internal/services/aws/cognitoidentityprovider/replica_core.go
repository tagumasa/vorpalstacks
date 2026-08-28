package cognitoidentityprovider

import (
	"time"

	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CreateUserPoolReplicaInput carries the wire parameters of
// CreateUserPoolReplica. Params holds the raw request parameter map for the
// UserPoolTags member.
type CreateUserPoolReplicaInput struct {
	UserPoolID string
	RegionName string
	Params     map[string]interface{}
}

// ListUserPoolReplicasInput carries the wire parameters of
// ListUserPoolReplicas.
type ListUserPoolReplicasInput struct {
	UserPoolID string
	NextToken  string
}

// DeleteUserPoolReplicaInput carries the wire parameters of
// DeleteUserPoolReplica.
type DeleteUserPoolReplicaInput struct {
	UserPoolID string
	RegionName string
}

// UpdateUserPoolReplicaInput carries the wire parameters of
// UpdateUserPoolReplica.
type UpdateUserPoolReplicaInput struct {
	UserPoolID string
	RegionName string
	Status     string
}

// createUserPoolReplicaCore creates a cross-region replica of a user pool.
func (s *CognitoService) createUserPoolReplicaCore(reqCtx *request.RequestContext, in CreateUserPoolReplicaInput) (interface{}, error) {
	if in.UserPoolID == "" || in.RegionName == "" {
		return nil, ErrInvalidParameter
	}
	if err := validateRegionName(in.RegionName); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	pool, err := store.GetUserPool(in.UserPoolID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	// The replica lives in the region named by the request, so its pool ARN
	// carries that region, not the service's constructor region.
	replicaArn := svcarn.NewARNBuilder(s.accountID, in.RegionName).Build("cognito-idp", "userpool/"+pool.ID)
	// ReplicaStatusType reports a replica available for both end-user and
	// administrator operations as ACTIVE; ReplicaRoleType assigns SECONDARY
	// to the replica (the source user pool keeps the PRIMARY role).
	replica := &cognitostore.UserPoolReplica{
		UserPoolID:   in.UserPoolID,
		RegionName:   in.RegionName,
		Status:       "ACTIVE",
		Role:         "SECONDARY",
		UserPoolArn:  replicaArn,
		CreationDate: time.Now().UTC(),
	}

	if err := store.SaveUserPoolReplica(replica); err != nil {
		return nil, ErrInternalError
	}

	parsedTags := tagutil.ParseTagsWithQueryFallback(in.Params, "UserPoolTags")
	if len(parsedTags) > 0 {
		replica.Tags = parsedTags
		tagMap := tagutil.ToMap(parsedTags)
		if err := store.Tag(replicaArn, tagMap); err != nil {
			// Rollback: remove the saved replica to avoid an untagged orphan.
			_ = store.DeleteUserPoolReplica(in.UserPoolID, in.RegionName)
			return nil, ErrInternalError
		}
	}

	return map[string]interface{}{"UserPoolReplica": formatUserPoolReplica(replica)}, nil
}

// listUserPoolReplicasCore lists replicas for a user pool.
func (s *CognitoService) listUserPoolReplicasCore(reqCtx *request.RequestContext, in ListUserPoolReplicasInput) (interface{}, error) {
	if in.UserPoolID == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListUserPoolReplicasPaginated(in.UserPoolID, storecommon.ListOptions{
		MaxItems: listLimitMax,
		Marker:   in.NextToken,
	})
	if err != nil {
		return nil, ErrInternalError
	}

	formatted := make([]map[string]interface{}, 0, len(result.Items))
	for _, r := range result.Items {
		formatted = append(formatted, formatUserPoolReplica(r))
	}

	resp := map[string]interface{}{"UserPoolReplicas": formatted}
	if result.IsTruncated && result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// deleteUserPoolReplicaCore deletes a cross-region replica.
func (s *CognitoService) deleteUserPoolReplicaCore(reqCtx *request.RequestContext, in DeleteUserPoolReplicaInput) (interface{}, error) {
	if in.UserPoolID == "" || in.RegionName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replica, err := store.GetUserPoolReplica(in.UserPoolID, in.RegionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if err := store.DeleteUserPoolReplica(in.UserPoolID, in.RegionName); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"UserPoolReplica": formatUserPoolReplica(replica)}, nil
}

// updateUserPoolReplicaCore updates a cross-region replica.
func (s *CognitoService) updateUserPoolReplicaCore(reqCtx *request.RequestContext, in UpdateUserPoolReplicaInput) (interface{}, error) {
	if in.UserPoolID == "" || in.RegionName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	replica, err := store.GetUserPoolReplica(in.UserPoolID, in.RegionName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	// UpdateUserPoolReplicaRequest marks Status required: the operation is
	// defined only as "the status to set for the replica" (ACTIVE or
	// INACTIVE), so an omitted member cannot identify an update.
	if in.Status == "" {
		return nil, ErrInvalidParameter
	}
	if !validateUpdateReplicaStatus(in.Status) {
		return nil, ErrInvalidParameter
	}
	replica.Status = in.Status

	if err := store.SaveUserPoolReplica(replica); err != nil {
		return nil, ErrInternalError
	}

	return map[string]interface{}{"UserPoolReplica": formatUserPoolReplica(replica)}, nil
}
