package neptune

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

type StoreError = common.StoreError

var (
	// ErrDBClusterNotFound is returned when a DB cluster cannot be found by its identifier.
	ErrDBClusterNotFound = errors.New("neptune: DBCluster not found")
	// ErrDBInstanceNotFound is returned when a DB instance cannot be found by its identifier.
	ErrDBInstanceNotFound = errors.New("neptune: DBInstance not found")
	// ErrDBClusterSnapshotNotFound is returned when a cluster snapshot cannot be found.
	ErrDBClusterSnapshotNotFound = errors.New("neptune: DBClusterSnapshot not found")
	// ErrDBClusterParameterGroupNotFound is returned when a cluster parameter group cannot be found.
	ErrDBClusterParameterGroupNotFound = errors.New("neptune: DBClusterParameterGroup not found")
	// ErrDBParameterGroupNotFound is returned when a DB parameter group cannot be found.
	ErrDBParameterGroupNotFound = errors.New("neptune: DBParameterGroup not found")
	// ErrDBSubnetGroupNotFound is returned when a DB subnet group cannot be found.
	ErrDBSubnetGroupNotFound = errors.New("neptune: DBSubnetGroup not found")
	// ErrGlobalClusterNotFound is returned when a global cluster cannot be found.
	ErrGlobalClusterNotFound = errors.New("neptune: GlobalCluster not found")
	// ErrEventSubscriptionNotFound is returned when an event subscription cannot be found.
	ErrEventSubscriptionNotFound = errors.New("neptune: EventSubscription not found")
	// ErrDBClusterEndpointNotFound is returned when a cluster endpoint cannot be found.
	ErrDBClusterEndpointNotFound = errors.New("neptune: DBClusterEndpoint not found")
	// ErrDBClusterAlreadyExists is returned when a DB cluster with the given identifier already exists.
	ErrDBClusterAlreadyExists = errors.New("neptune: DBCluster already exists")
	// ErrDBInstanceAlreadyExists is returned when a DB instance with the given identifier already exists.
	ErrDBInstanceAlreadyExists = errors.New("neptune: DBInstance already exists")
	// ErrDBClusterSnapshotAlreadyExists is returned when a cluster snapshot with the given identifier already exists.
	ErrDBClusterSnapshotAlreadyExists = errors.New("neptune: DBClusterSnapshot already exists")
	// ErrDBClusterParameterGroupAlreadyExists is returned when a cluster parameter group with the given name already exists.
	ErrDBClusterParameterGroupAlreadyExists = errors.New("neptune: DBClusterParameterGroup already exists")
	// ErrDBParameterGroupAlreadyExists is returned when a DB parameter group with the given name already exists.
	ErrDBParameterGroupAlreadyExists = errors.New("neptune: DBParameterGroup already exists")
	// ErrDBSubnetGroupAlreadyExists is returned when a DB subnet group with the given name already exists.
	ErrDBSubnetGroupAlreadyExists = errors.New("neptune: DBSubnetGroup already exists")
	// ErrGlobalClusterAlreadyExists is returned when a global cluster with the given identifier already exists.
	ErrGlobalClusterAlreadyExists = errors.New("neptune: GlobalCluster already exists")
	// ErrEventSubscriptionAlreadyExists is returned when an event subscription with the given name already exists.
	ErrEventSubscriptionAlreadyExists = errors.New("neptune: EventSubscription already exists")
	// ErrDBClusterEndpointAlreadyExists is returned when a cluster endpoint with the given identifier already exists.
	ErrDBClusterEndpointAlreadyExists = errors.New("neptune: DBClusterEndpoint already exists")
	// ErrInvalidParameterGroupState is returned when a parameter group is not in a valid state for the requested operation.
	ErrInvalidParameterGroupState = errors.New("neptune: InvalidDBClusterParameterGroupState")
)

// IsNotFound checks if the error indicates a Neptune resource was not found.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrDBClusterNotFound) ||
		errors.Is(err, ErrDBInstanceNotFound) ||
		errors.Is(err, ErrDBClusterSnapshotNotFound) ||
		errors.Is(err, ErrDBClusterParameterGroupNotFound) ||
		errors.Is(err, ErrDBParameterGroupNotFound) ||
		errors.Is(err, ErrDBSubnetGroupNotFound) ||
		errors.Is(err, ErrGlobalClusterNotFound) ||
		errors.Is(err, ErrEventSubscriptionNotFound) ||
		errors.Is(err, ErrDBClusterEndpointNotFound) ||
		common.IsNotFound(err)
}

// IsAlreadyExists checks if the error indicates a Neptune resource already exists.
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrDBClusterAlreadyExists) ||
		errors.Is(err, ErrDBInstanceAlreadyExists) ||
		errors.Is(err, ErrDBClusterSnapshotAlreadyExists) ||
		errors.Is(err, ErrDBClusterParameterGroupAlreadyExists) ||
		errors.Is(err, ErrDBParameterGroupAlreadyExists) ||
		errors.Is(err, ErrDBSubnetGroupAlreadyExists) ||
		errors.Is(err, ErrGlobalClusterAlreadyExists) ||
		errors.Is(err, ErrEventSubscriptionAlreadyExists) ||
		errors.Is(err, ErrDBClusterEndpointAlreadyExists) ||
		common.IsAlreadyExists(err)
}
