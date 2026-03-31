// Package neptune implements the AWS Neptune Management API service handlers.
package neptune

import (
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

// NeptuneService handles incoming Neptune Management API requests.
type NeptuneService struct{}

// NewNeptuneService creates a new NeptuneService. The store, accountID, and
// region arguments are accepted for initialisation compatibility but are not
// currently retained.
func NewNeptuneService(store storage.BasicStorage, accountID, region string) *NeptuneService {
	return &NeptuneService{}
}

// store retrieves or creates the NeptuneStoreInterface for the current request
// context, using the pre-initialised store when available.
func (s *NeptuneService) store(reqCtx *request.RequestContext) (neptunestore.NeptuneStoreInterface, error) {
	if ns := reqCtx.GetNeptuneStore(); ns != nil {
		return ns, nil
	}
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	return neptunestore.NewNeptuneStore(storage), nil
}

// RegisterHandlers registers all Neptune Management API action handlers with
// the given dispatcher.
func (s *NeptuneService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("neptune", "CreateDBCluster", s.CreateDBCluster)
	d.RegisterHandlerForService("neptune", "DeleteDBCluster", s.DeleteDBCluster)
	d.RegisterHandlerForService("neptune", "ModifyDBCluster", s.ModifyDBCluster)
	d.RegisterHandlerForService("neptune", "DescribeDBClusters", s.DescribeDBClusters)
	d.RegisterHandlerForService("neptune", "StartDBCluster", s.StartDBCluster)
	d.RegisterHandlerForService("neptune", "StopDBCluster", s.StopDBCluster)
	d.RegisterHandlerForService("neptune", "FailoverDBCluster", s.FailoverDBCluster)
	d.RegisterHandlerForService("neptune", "CreateDBClusterEndpoint", s.CreateDBClusterEndpoint)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterEndpoints", s.DescribeDBClusterEndpoints)
	d.RegisterHandlerForService("neptune", "ModifyDBClusterEndpoint", s.ModifyDBClusterEndpoint)
	d.RegisterHandlerForService("neptune", "DeleteDBClusterEndpoint", s.DeleteDBClusterEndpoint)
	d.RegisterHandlerForService("neptune", "CreateDBInstance", s.CreateDBInstance)
	d.RegisterHandlerForService("neptune", "DeleteDBInstance", s.DeleteDBInstance)
	d.RegisterHandlerForService("neptune", "ModifyDBInstance", s.ModifyDBInstance)
	d.RegisterHandlerForService("neptune", "DescribeDBInstances", s.DescribeDBInstances)
	d.RegisterHandlerForService("neptune", "RebootDBInstance", s.RebootDBInstance)
	d.RegisterHandlerForService("neptune", "CreateDBClusterSnapshot", s.CreateDBClusterSnapshot)
	d.RegisterHandlerForService("neptune", "DeleteDBClusterSnapshot", s.DeleteDBClusterSnapshot)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterSnapshots", s.DescribeDBClusterSnapshots)
	d.RegisterHandlerForService("neptune", "CopyDBClusterSnapshot", s.CopyDBClusterSnapshot)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterSnapshotAttributes", s.DescribeDBClusterSnapshotAttributes)
	d.RegisterHandlerForService("neptune", "ModifyDBClusterSnapshotAttribute", s.ModifyDBClusterSnapshotAttribute)
	d.RegisterHandlerForService("neptune", "CreateDBClusterParameterGroup", s.CreateDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "DeleteDBClusterParameterGroup", s.DeleteDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterParameterGroups", s.DescribeDBClusterParameterGroups)
	d.RegisterHandlerForService("neptune", "DescribeDBClusterParameters", s.DescribeDBClusterParameters)
	d.RegisterHandlerForService("neptune", "ModifyDBClusterParameterGroup", s.ModifyDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "CreateDBParameterGroup", s.CreateDBParameterGroup)
	d.RegisterHandlerForService("neptune", "DeleteDBParameterGroup", s.DeleteDBParameterGroup)
	d.RegisterHandlerForService("neptune", "DescribeDBParameterGroups", s.DescribeDBParameterGroups)
	d.RegisterHandlerForService("neptune", "DescribeDBParameters", s.DescribeDBParameters)
	d.RegisterHandlerForService("neptune", "ModifyDBParameterGroup", s.ModifyDBParameterGroup)
	d.RegisterHandlerForService("neptune", "ResetDBClusterParameterGroup", s.ResetDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "ResetDBParameterGroup", s.ResetDBParameterGroup)
	d.RegisterHandlerForService("neptune", "CopyDBClusterParameterGroup", s.CopyDBClusterParameterGroup)
	d.RegisterHandlerForService("neptune", "CopyDBParameterGroup", s.CopyDBParameterGroup)
	d.RegisterHandlerForService("neptune", "DescribeEngineDefaultClusterParameters", s.DescribeEngineDefaultClusterParameters)
	d.RegisterHandlerForService("neptune", "DescribeEngineDefaultParameters", s.DescribeEngineDefaultParameters)
	d.RegisterHandlerForService("neptune", "CreateGlobalCluster", s.CreateGlobalCluster)
	d.RegisterHandlerForService("neptune", "DeleteGlobalCluster", s.DeleteGlobalCluster)
	d.RegisterHandlerForService("neptune", "DescribeGlobalClusters", s.DescribeGlobalClusters)
	d.RegisterHandlerForService("neptune", "ModifyGlobalCluster", s.ModifyGlobalCluster)
	d.RegisterHandlerForService("neptune", "FailoverGlobalCluster", s.FailoverGlobalCluster)
	d.RegisterHandlerForService("neptune", "SwitchoverGlobalCluster", s.SwitchoverGlobalCluster)
	d.RegisterHandlerForService("neptune", "RemoveFromGlobalCluster", s.RemoveFromGlobalCluster)
	d.RegisterHandlerForService("neptune", "CreateDBSubnetGroup", s.CreateDBSubnetGroup)
	d.RegisterHandlerForService("neptune", "DeleteDBSubnetGroup", s.DeleteDBSubnetGroup)
	d.RegisterHandlerForService("neptune", "DescribeDBSubnetGroups", s.DescribeDBSubnetGroups)
	d.RegisterHandlerForService("neptune", "ModifyDBSubnetGroup", s.ModifyDBSubnetGroup)
	d.RegisterHandlerForService("neptune", "CreateEventSubscription", s.CreateEventSubscription)
	d.RegisterHandlerForService("neptune", "DeleteEventSubscription", s.DeleteEventSubscription)
	d.RegisterHandlerForService("neptune", "DescribeEventSubscriptions", s.DescribeEventSubscriptions)
	d.RegisterHandlerForService("neptune", "ModifyEventSubscription", s.ModifyEventSubscription)
	d.RegisterHandlerForService("neptune", "AddSourceIdentifierToSubscription", s.AddSourceIdentifierToSubscription)
	d.RegisterHandlerForService("neptune", "RemoveSourceIdentifierFromSubscription", s.RemoveSourceIdentifierFromSubscription)
	d.RegisterHandlerForService("neptune", "AddTagsToResource", s.AddTagsToResource)
	d.RegisterHandlerForService("neptune", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("neptune", "RemoveTagsFromResource", s.RemoveTagsFromResource)
	d.RegisterHandlerForService("neptune", "DescribeEvents", s.DescribeEvents)
	d.RegisterHandlerForService("neptune", "DescribeEventCategories", s.DescribeEventCategories)
	d.RegisterHandlerForService("neptune", "DescribePendingMaintenanceActions", s.DescribePendingMaintenanceActions)
	d.RegisterHandlerForService("neptune", "ApplyPendingMaintenanceAction", s.ApplyPendingMaintenanceAction)
	d.RegisterHandlerForService("neptune", "DescribeDBEngineVersions", s.DescribeDBEngineVersions)
	d.RegisterHandlerForService("neptune", "DescribeOrderableDBInstanceOptions", s.DescribeOrderableDBInstanceOptions)
	d.RegisterHandlerForService("neptune", "DescribeValidDBInstanceModifications", s.DescribeValidDBInstanceModifications)
	d.RegisterHandlerForService("neptune", "RestoreDBClusterFromSnapshot", s.RestoreDBClusterFromSnapshot)
	d.RegisterHandlerForService("neptune", "RestoreDBClusterToPointInTime", s.RestoreDBClusterToPointInTime)
	d.RegisterHandlerForService("neptune", "PromoteReadReplicaDBCluster", s.PromoteReadReplicaDBCluster)
	d.RegisterHandlerForService("neptune", "AddRoleToDBCluster", s.AddRoleToDBCluster)
	d.RegisterHandlerForService("neptune", "RemoveRoleFromDBCluster", s.RemoveRoleFromDBCluster)
}
