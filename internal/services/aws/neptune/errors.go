package neptune

import (
	"errors"
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

// translateStoreError converts a Neptune store-level sentinel error into the
// corresponding AWS-compatible *AWSError so the dispatcher can serialise it
// correctly.  Unknown errors are returned unchanged and will fall through to
// the dispatcher's generic InternalFailure handler.
func translateStoreError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, neptunestore.ErrDBClusterNotFound):
		return awserrors.NewAWSError("DBClusterNotFoundFault", "DBCluster not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrDBInstanceNotFound):
		return awserrors.NewAWSError("DBInstanceNotFoundFault", "DBInstance not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrDBClusterSnapshotNotFound):
		return awserrors.NewAWSError("DBClusterSnapshotNotFoundFault", "DBClusterSnapshot not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrDBClusterParameterGroupNotFound):
		return awserrors.NewAWSError("DBParameterGroupNotFoundFault", "DBParameterGroup not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrDBParameterGroupNotFound):
		return awserrors.NewAWSError("DBParameterGroupNotFoundFault", "DBParameterGroup not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrDBSubnetGroupNotFound):
		return awserrors.NewAWSError("DBSubnetGroupNotFoundFault", "DBSubnetGroup not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrGlobalClusterNotFound):
		return awserrors.NewAWSError("GlobalClusterNotFoundFault", "GlobalCluster not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrEventSubscriptionNotFound):
		return awserrors.NewAWSError("SubscriptionNotFoundFault", "EventSubscription not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrDBClusterEndpointNotFound):
		return awserrors.NewAWSError("DBClusterEndpointNotFoundFault", "DBClusterEndpoint not found: "+extractStoreMsg(err), http.StatusNotFound)
	case errors.Is(err, neptunestore.ErrDBClusterAlreadyExists):
		return awserrors.NewAWSError("DBClusterAlreadyExistsFault", "DBCluster already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrDBInstanceAlreadyExists):
		return awserrors.NewAWSError("DBInstanceAlreadyExistsFault", "DBInstance already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrDBClusterSnapshotAlreadyExists):
		return awserrors.NewAWSError("DBClusterSnapshotAlreadyExistsFault", "DBClusterSnapshot already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrDBClusterParameterGroupAlreadyExists):
		return awserrors.NewAWSError("DBParameterGroupAlreadyExistsFault", "DBClusterParameterGroup already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrDBParameterGroupAlreadyExists):
		return awserrors.NewAWSError("DBParameterGroupAlreadyExistsFault", "DBParameterGroup already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrDBSubnetGroupAlreadyExists):
		return awserrors.NewAWSError("DBSubnetGroupAlreadyExistsFault", "DBSubnetGroup already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrGlobalClusterAlreadyExists):
		return awserrors.NewAWSError("GlobalClusterAlreadyExistsFault", "GlobalCluster already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrEventSubscriptionAlreadyExists):
		return awserrors.NewAWSError("SubscriptionAlreadyExistFault", "EventSubscription already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrDBClusterEndpointAlreadyExists):
		return awserrors.NewAWSError("DBClusterEndpointAlreadyExistsFault", "DBClusterEndpoint already exists", http.StatusConflict)
	case errors.Is(err, neptunestore.ErrInvalidParameterGroupState):
		return awserrors.NewAWSError("InvalidParameterGroupState", "Invalid DB cluster parameter group state", http.StatusBadRequest)
	}
	return err
}

// extractStoreMsg returns the error message string, handling both plain
// sentinel errors and wrapped errors from neptune.NewStoreError.
func extractStoreMsg(err error) string {
	return err.Error()
}
