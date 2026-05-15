package neptune

import (
	"errors"

	rds "vorpalstacks/internal/store/aws/rds"
)

var (
	ErrDBClusterNotFound                    = rds.ErrDBClusterNotFound
	ErrDBInstanceNotFound                   = rds.ErrDBInstanceNotFound
	ErrDBClusterSnapshotNotFound            = rds.ErrDBClusterSnapshotNotFound
	ErrDBClusterParameterGroupNotFound      = rds.ErrDBClusterParameterGroupNotFound
	ErrDBParameterGroupNotFound             = rds.ErrDBParameterGroupNotFound
	ErrDBSubnetGroupNotFound                = rds.ErrDBSubnetGroupNotFound
	ErrGlobalClusterNotFound                = rds.ErrGlobalClusterNotFound
	ErrEventSubscriptionNotFound            = rds.ErrEventSubscriptionNotFound
	ErrEventNotFound                        = rds.ErrEventNotFound
	ErrDBClusterAlreadyExists               = rds.ErrDBClusterAlreadyExists
	ErrDBInstanceAlreadyExists              = rds.ErrDBInstanceAlreadyExists
	ErrDBClusterSnapshotAlreadyExists       = rds.ErrDBClusterSnapshotAlreadyExists
	ErrDBClusterParameterGroupAlreadyExists = rds.ErrDBClusterParameterGroupAlreadyExists
	ErrDBParameterGroupAlreadyExists        = rds.ErrDBParameterGroupAlreadyExists
	ErrDBSubnetGroupAlreadyExists           = rds.ErrDBSubnetGroupAlreadyExists
	ErrGlobalClusterAlreadyExists           = rds.ErrGlobalClusterAlreadyExists
	ErrEventSubscriptionAlreadyExists       = rds.ErrEventSubscriptionAlreadyExists
	ErrEventAlreadyExists                   = rds.ErrEventAlreadyExists
	ErrInvalidParameterGroupState           = rds.ErrInvalidParameterGroupState

	ErrDBClusterEndpointNotFound      = errors.New("neptune: DBClusterEndpoint not found")
	ErrDBClusterEndpointAlreadyExists = errors.New("neptune: DBClusterEndpoint already exists")
)
