// Package neptune provides the storage layer for AWS Neptune Management API resources.
package neptune

import "fmt"

// StoreError wraps a storage operation error with contextual information about
// which operation failed.
type StoreError struct {
	Operation string
	Err       error
}

// NewStoreError creates a new StoreError wrapping the given operation name and
// underlying error. Returns nil if the provided error is nil.
func NewStoreError(op string, err error) error {
	if err == nil {
		return nil
	}
	return &StoreError{Operation: op, Err: err}
}

func (e *StoreError) Error() string {
	return fmt.Sprintf("neptune store: %s: %v", e.Operation, e.Err)
}

func (e *StoreError) Unwrap() error {
	return e.Err
}

var (
	// ErrDBClusterNotFound is returned when a DB cluster cannot be found by its identifier.
	ErrDBClusterNotFound = fmt.Errorf("neptune: DBCluster not found")
	// ErrDBInstanceNotFound is returned when a DB instance cannot be found by its identifier.
	ErrDBInstanceNotFound = fmt.Errorf("neptune: DBInstance not found")
	// ErrDBClusterSnapshotNotFound is returned when a cluster snapshot cannot be found.
	ErrDBClusterSnapshotNotFound = fmt.Errorf("neptune: DBClusterSnapshot not found")
	// ErrDBClusterParameterGroupNotFound is returned when a cluster parameter group cannot be found.
	ErrDBClusterParameterGroupNotFound = fmt.Errorf("neptune: DBClusterParameterGroup not found")
	// ErrDBParameterGroupNotFound is returned when a DB parameter group cannot be found.
	ErrDBParameterGroupNotFound = fmt.Errorf("neptune: DBParameterGroup not found")
	// ErrDBSubnetGroupNotFound is returned when a DB subnet group cannot be found.
	ErrDBSubnetGroupNotFound = fmt.Errorf("neptune: DBSubnetGroup not found")
	// ErrGlobalClusterNotFound is returned when a global cluster cannot be found.
	ErrGlobalClusterNotFound = fmt.Errorf("neptune: GlobalCluster not found")
	// ErrEventSubscriptionNotFound is returned when an event subscription cannot be found.
	ErrEventSubscriptionNotFound = fmt.Errorf("neptune: EventSubscription not found")
	// ErrDBClusterEndpointNotFound is returned when a cluster endpoint cannot be found.
	ErrDBClusterEndpointNotFound            = fmt.Errorf("neptune: DBClusterEndpoint not found")
	ErrDBClusterAlreadyExists               = fmt.Errorf("neptune: DBCluster already exists")
	ErrDBInstanceAlreadyExists              = fmt.Errorf("neptune: DBInstance already exists")
	ErrDBClusterSnapshotAlreadyExists       = fmt.Errorf("neptune: DBClusterSnapshot already exists")
	ErrDBClusterParameterGroupAlreadyExists = fmt.Errorf("neptune: DBClusterParameterGroup already exists")
	ErrDBParameterGroupAlreadyExists        = fmt.Errorf("neptune: DBParameterGroup already exists")
	ErrDBSubnetGroupAlreadyExists           = fmt.Errorf("neptune: DBSubnetGroup already exists")
	ErrGlobalClusterAlreadyExists           = fmt.Errorf("neptune: GlobalCluster already exists")
	ErrEventSubscriptionAlreadyExists       = fmt.Errorf("neptune: EventSubscription already exists")
	ErrDBClusterEndpointAlreadyExists       = fmt.Errorf("neptune: DBClusterEndpoint already exists")
	ErrInvalidParameterGroupState           = fmt.Errorf("neptune: InvalidDBClusterParameterGroupState")
)
