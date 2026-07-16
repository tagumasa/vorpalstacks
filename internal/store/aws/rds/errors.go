package rds

import "errors"

var (
	ErrDBClusterNotFound                    = errors.New("rds: DBCluster not found")
	ErrDBInstanceNotFound                   = errors.New("rds: DBInstance not found")
	ErrDBClusterSnapshotNotFound            = errors.New("rds: DBClusterSnapshot not found")
	ErrDBSnapshotNotFound                   = errors.New("rds: DBSnapshot not found")
	ErrDBClusterParameterGroupNotFound      = errors.New("rds: DBClusterParameterGroup not found")
	ErrDBParameterGroupNotFound             = errors.New("rds: DBParameterGroup not found")
	ErrDBSubnetGroupNotFound                = errors.New("rds: DBSubnetGroup not found")
	ErrGlobalClusterNotFound                = errors.New("rds: GlobalCluster not found")
	ErrEventSubscriptionNotFound            = errors.New("rds: EventSubscription not found")
	ErrEventNotFound                        = errors.New("rds: Event not found")
	ErrDBClusterAlreadyExists               = errors.New("rds: DBCluster already exists")
	ErrDBInstanceAlreadyExists              = errors.New("rds: DBInstance already exists")
	ErrDBClusterSnapshotAlreadyExists       = errors.New("rds: DBClusterSnapshot already exists")
	ErrDBSnapshotAlreadyExists              = errors.New("rds: DBSnapshot already exists")
	ErrDBClusterParameterGroupAlreadyExists = errors.New("rds: DBClusterParameterGroup already exists")
	ErrDBParameterGroupAlreadyExists        = errors.New("rds: DBParameterGroup already exists")
	ErrDBSubnetGroupAlreadyExists           = errors.New("rds: DBSubnetGroup already exists")
	ErrGlobalClusterAlreadyExists           = errors.New("rds: GlobalCluster already exists")
	ErrEventSubscriptionAlreadyExists       = errors.New("rds: EventSubscription already exists")
	ErrEventAlreadyExists                   = errors.New("rds: Event already exists")
	ErrInvalidParameterGroupState           = errors.New("rds: InvalidDBClusterParameterGroupState")
)
