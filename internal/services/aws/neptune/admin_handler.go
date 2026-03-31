package neptune

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/rds"
	rdsconnect "vorpalstacks/internal/pb/aws/rds/rdsconnect"
	storeneptune "vorpalstacks/internal/store/aws/neptune"
)

type AdminHandler struct {
	rdsconnect.UnimplementedNeptuneServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ rdsconnect.NeptuneServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(sm *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{storageManager: sm, accountId: accountId}
}

func (h *AdminHandler) getStore(req connect.Request[pb.CreateDBClusterMessage]) (*storeneptune.NeptuneStore, error) {
	region := req.Header().Get("X-Aws-Region")
	if region == "" {
		region = "us-east-1"
	}
	rs, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return storeneptune.NewNeptuneStore(rs), nil
}

func (h *AdminHandler) DescribeDBClusters(ctx context.Context, req *connect.Request[pb.DescribeDBClustersMessage]) (*connect.Response[pb.DBClusterMessage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for DescribeDBClusters"))
}

func (h *AdminHandler) DescribeDBInstances(ctx context.Context, req *connect.Request[pb.DescribeDBInstancesMessage]) (*connect.Response[pb.DBInstanceMessage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for DescribeDBInstances"))
}

func nop(name string) error {
	return connect.NewError(connect.CodeUnimplemented, fmt.Errorf("use AWS REST API for %s", name))
}

func (h *AdminHandler) AddRoleToDBCluster(ctx context.Context, req *connect.Request[pb.AddRoleToDBClusterMessage]) (*connect.Response[common.Empty], error) {
	return nil, nop("AddRoleToDBCluster")
}

func (h *AdminHandler) AddSourceIdentifierToSubscription(ctx context.Context, req *connect.Request[pb.AddSourceIdentifierToSubscriptionMessage]) (*connect.Response[pb.AddSourceIdentifierToSubscriptionResult], error) {
	return nil, nop("AddSourceIdentifierToSubscription")
}

func (h *AdminHandler) AddTagsToResource(ctx context.Context, req *connect.Request[pb.AddTagsToResourceMessage]) (*connect.Response[common.Empty], error) {
	return nil, nop("AddTagsToResource")
}

func (h *AdminHandler) ApplyPendingMaintenanceAction(ctx context.Context, req *connect.Request[pb.ApplyPendingMaintenanceActionMessage]) (*connect.Response[pb.ApplyPendingMaintenanceActionResult], error) {
	return nil, nop("ApplyPendingMaintenanceAction")
}

func (h *AdminHandler) CopyDBClusterParameterGroup(ctx context.Context, req *connect.Request[pb.CopyDBClusterParameterGroupMessage]) (*connect.Response[pb.CopyDBClusterParameterGroupResult], error) {
	return nil, nop("CopyDBClusterParameterGroup")
}

func (h *AdminHandler) CopyDBClusterSnapshot(ctx context.Context, req *connect.Request[pb.CopyDBClusterSnapshotMessage]) (*connect.Response[pb.CopyDBClusterSnapshotResult], error) {
	return nil, nop("CopyDBClusterSnapshot")
}

func (h *AdminHandler) CopyDBParameterGroup(ctx context.Context, req *connect.Request[pb.CopyDBParameterGroupMessage]) (*connect.Response[pb.CopyDBParameterGroupResult], error) {
	return nil, nop("CopyDBParameterGroup")
}

func (h *AdminHandler) CreateDBCluster(ctx context.Context, req *connect.Request[pb.CreateDBClusterMessage]) (*connect.Response[pb.CreateDBClusterResult], error) {
	return nil, nop("CreateDBCluster")
}

func (h *AdminHandler) CreateDBClusterEndpoint(ctx context.Context, req *connect.Request[pb.CreateDBClusterEndpointMessage]) (*connect.Response[pb.CreateDBClusterEndpointOutput], error) {
	return nil, nop("CreateDBClusterEndpoint")
}

func (h *AdminHandler) CreateDBClusterParameterGroup(ctx context.Context, req *connect.Request[pb.CreateDBClusterParameterGroupMessage]) (*connect.Response[pb.CreateDBClusterParameterGroupResult], error) {
	return nil, nop("CreateDBClusterParameterGroup")
}

func (h *AdminHandler) CreateDBClusterSnapshot(ctx context.Context, req *connect.Request[pb.CreateDBClusterSnapshotMessage]) (*connect.Response[pb.CreateDBClusterSnapshotResult], error) {
	return nil, nop("CreateDBClusterSnapshot")
}

func (h *AdminHandler) CreateDBInstance(ctx context.Context, req *connect.Request[pb.CreateDBInstanceMessage]) (*connect.Response[pb.CreateDBInstanceResult], error) {
	return nil, nop("CreateDBInstance")
}

func (h *AdminHandler) CreateDBParameterGroup(ctx context.Context, req *connect.Request[pb.CreateDBParameterGroupMessage]) (*connect.Response[pb.CreateDBParameterGroupResult], error) {
	return nil, nop("CreateDBParameterGroup")
}

func (h *AdminHandler) CreateDBSubnetGroup(ctx context.Context, req *connect.Request[pb.CreateDBSubnetGroupMessage]) (*connect.Response[pb.CreateDBSubnetGroupResult], error) {
	return nil, nop("CreateDBSubnetGroup")
}

func (h *AdminHandler) CreateEventSubscription(ctx context.Context, req *connect.Request[pb.CreateEventSubscriptionMessage]) (*connect.Response[pb.CreateEventSubscriptionResult], error) {
	return nil, nop("CreateEventSubscription")
}

func (h *AdminHandler) CreateGlobalCluster(ctx context.Context, req *connect.Request[pb.CreateGlobalClusterMessage]) (*connect.Response[pb.CreateGlobalClusterResult], error) {
	return nil, nop("CreateGlobalCluster")
}

func (h *AdminHandler) DeleteDBCluster(ctx context.Context, req *connect.Request[pb.DeleteDBClusterMessage]) (*connect.Response[pb.DeleteDBClusterResult], error) {
	return nil, nop("DeleteDBCluster")
}

func (h *AdminHandler) DeleteDBClusterEndpoint(ctx context.Context, req *connect.Request[pb.DeleteDBClusterEndpointMessage]) (*connect.Response[pb.DeleteDBClusterEndpointOutput], error) {
	return nil, nop("DeleteDBClusterEndpoint")
}

func (h *AdminHandler) DeleteDBClusterParameterGroup(ctx context.Context, req *connect.Request[pb.DeleteDBClusterParameterGroupMessage]) (*connect.Response[common.Empty], error) {
	return nil, nop("DeleteDBClusterParameterGroup")
}

func (h *AdminHandler) DeleteDBClusterSnapshot(ctx context.Context, req *connect.Request[pb.DeleteDBClusterSnapshotMessage]) (*connect.Response[pb.DeleteDBClusterSnapshotResult], error) {
	return nil, nop("DeleteDBClusterSnapshot")
}

func (h *AdminHandler) DeleteDBInstance(ctx context.Context, req *connect.Request[pb.DeleteDBInstanceMessage]) (*connect.Response[pb.DeleteDBInstanceResult], error) {
	return nil, nop("DeleteDBInstance")
}

func (h *AdminHandler) DeleteDBParameterGroup(ctx context.Context, req *connect.Request[pb.DeleteDBParameterGroupMessage]) (*connect.Response[common.Empty], error) {
	return nil, nop("DeleteDBParameterGroup")
}

func (h *AdminHandler) DeleteDBSubnetGroup(ctx context.Context, req *connect.Request[pb.DeleteDBSubnetGroupMessage]) (*connect.Response[common.Empty], error) {
	return nil, nop("DeleteDBSubnetGroup")
}

func (h *AdminHandler) DeleteEventSubscription(ctx context.Context, req *connect.Request[pb.DeleteEventSubscriptionMessage]) (*connect.Response[pb.DeleteEventSubscriptionResult], error) {
	return nil, nop("DeleteEventSubscription")
}

func (h *AdminHandler) DeleteGlobalCluster(ctx context.Context, req *connect.Request[pb.DeleteGlobalClusterMessage]) (*connect.Response[pb.DeleteGlobalClusterResult], error) {
	return nil, nop("DeleteGlobalCluster")
}

func (h *AdminHandler) DescribeDBClusterEndpoints(ctx context.Context, req *connect.Request[pb.DescribeDBClusterEndpointsMessage]) (*connect.Response[pb.DBClusterEndpointMessage], error) {
	return nil, nop("DescribeDBClusterEndpoints")
}

func (h *AdminHandler) DescribeDBClusterParameterGroups(ctx context.Context, req *connect.Request[pb.DescribeDBClusterParameterGroupsMessage]) (*connect.Response[pb.DBClusterParameterGroupsMessage], error) {
	return nil, nop("DescribeDBClusterParameterGroups")
}

func (h *AdminHandler) DescribeDBClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeDBClusterParametersMessage]) (*connect.Response[pb.DBClusterParameterGroupDetails], error) {
	return nil, nop("DescribeDBClusterParameters")
}

func (h *AdminHandler) DescribeDBClusterSnapshotAttributes(ctx context.Context, req *connect.Request[pb.DescribeDBClusterSnapshotAttributesMessage]) (*connect.Response[pb.DescribeDBClusterSnapshotAttributesResult], error) {
	return nil, nop("DescribeDBClusterSnapshotAttributes")
}

func (h *AdminHandler) DescribeDBClusterSnapshots(ctx context.Context, req *connect.Request[pb.DescribeDBClusterSnapshotsMessage]) (*connect.Response[pb.DBClusterSnapshotMessage], error) {
	return nil, nop("DescribeDBClusterSnapshots")
}

func (h *AdminHandler) DescribeDBEngineVersions(ctx context.Context, req *connect.Request[pb.DescribeDBEngineVersionsMessage]) (*connect.Response[pb.DBEngineVersionMessage], error) {
	return nil, nop("DescribeDBEngineVersions")
}

func (h *AdminHandler) DescribeDBParameterGroups(ctx context.Context, req *connect.Request[pb.DescribeDBParameterGroupsMessage]) (*connect.Response[pb.DBParameterGroupsMessage], error) {
	return nil, nop("DescribeDBParameterGroups")
}

func (h *AdminHandler) DescribeDBParameters(ctx context.Context, req *connect.Request[pb.DescribeDBParametersMessage]) (*connect.Response[pb.DBParameterGroupDetails], error) {
	return nil, nop("DescribeDBParameters")
}

func (h *AdminHandler) DescribeDBSubnetGroups(ctx context.Context, req *connect.Request[pb.DescribeDBSubnetGroupsMessage]) (*connect.Response[pb.DBSubnetGroupMessage], error) {
	return nil, nop("DescribeDBSubnetGroups")
}

func (h *AdminHandler) DescribeEngineDefaultClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultClusterParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultClusterParametersResult], error) {
	return nil, nop("DescribeEngineDefaultClusterParameters")
}

func (h *AdminHandler) DescribeEngineDefaultParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultParametersResult], error) {
	return nil, nop("DescribeEngineDefaultParameters")
}

func (h *AdminHandler) DescribeEventCategories(ctx context.Context, req *connect.Request[pb.DescribeEventCategoriesMessage]) (*connect.Response[pb.EventCategoriesMessage], error) {
	return nil, nop("DescribeEventCategories")
}

func (h *AdminHandler) DescribeEventSubscriptions(ctx context.Context, req *connect.Request[pb.DescribeEventSubscriptionsMessage]) (*connect.Response[pb.EventSubscriptionsMessage], error) {
	return nil, nop("DescribeEventSubscriptions")
}

func (h *AdminHandler) DescribeEvents(ctx context.Context, req *connect.Request[pb.DescribeEventsMessage]) (*connect.Response[pb.EventsMessage], error) {
	return nil, nop("DescribeEvents")
}

func (h *AdminHandler) DescribeGlobalClusters(ctx context.Context, req *connect.Request[pb.DescribeGlobalClustersMessage]) (*connect.Response[pb.GlobalClustersMessage], error) {
	return nil, nop("DescribeGlobalClusters")
}

func (h *AdminHandler) DescribeOrderableDBInstanceOptions(ctx context.Context, req *connect.Request[pb.DescribeOrderableDBInstanceOptionsMessage]) (*connect.Response[pb.OrderableDBInstanceOptionsMessage], error) {
	return nil, nop("DescribeOrderableDBInstanceOptions")
}

func (h *AdminHandler) DescribePendingMaintenanceActions(ctx context.Context, req *connect.Request[pb.DescribePendingMaintenanceActionsMessage]) (*connect.Response[pb.PendingMaintenanceActionsMessage], error) {
	return nil, nop("DescribePendingMaintenanceActions")
}

func (h *AdminHandler) DescribeValidDBInstanceModifications(ctx context.Context, req *connect.Request[pb.DescribeValidDBInstanceModificationsMessage]) (*connect.Response[pb.DescribeValidDBInstanceModificationsResult], error) {
	return nil, nop("DescribeValidDBInstanceModifications")
}

func (h *AdminHandler) FailoverDBCluster(ctx context.Context, req *connect.Request[pb.FailoverDBClusterMessage]) (*connect.Response[pb.FailoverDBClusterResult], error) {
	return nil, nop("FailoverDBCluster")
}

func (h *AdminHandler) FailoverGlobalCluster(ctx context.Context, req *connect.Request[pb.FailoverGlobalClusterMessage]) (*connect.Response[pb.FailoverGlobalClusterResult], error) {
	return nil, nop("FailoverGlobalCluster")
}

func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceMessage]) (*connect.Response[pb.TagListMessage], error) {
	return nil, nop("ListTagsForResource")
}

func (h *AdminHandler) ModifyDBCluster(ctx context.Context, req *connect.Request[pb.ModifyDBClusterMessage]) (*connect.Response[pb.ModifyDBClusterResult], error) {
	return nil, nop("ModifyDBCluster")
}

func (h *AdminHandler) ModifyDBClusterEndpoint(ctx context.Context, req *connect.Request[pb.ModifyDBClusterEndpointMessage]) (*connect.Response[pb.ModifyDBClusterEndpointOutput], error) {
	return nil, nop("ModifyDBClusterEndpoint")
}

func (h *AdminHandler) ModifyDBClusterParameterGroup(ctx context.Context, req *connect.Request[pb.ModifyDBClusterParameterGroupMessage]) (*connect.Response[pb.DBClusterParameterGroupNameMessage], error) {
	return nil, nop("ModifyDBClusterParameterGroup")
}

func (h *AdminHandler) ModifyDBClusterSnapshotAttribute(ctx context.Context, req *connect.Request[pb.ModifyDBClusterSnapshotAttributeMessage]) (*connect.Response[pb.ModifyDBClusterSnapshotAttributeResult], error) {
	return nil, nop("ModifyDBClusterSnapshotAttribute")
}

func (h *AdminHandler) ModifyDBInstance(ctx context.Context, req *connect.Request[pb.ModifyDBInstanceMessage]) (*connect.Response[pb.ModifyDBInstanceResult], error) {
	return nil, nop("ModifyDBInstance")
}

func (h *AdminHandler) ModifyDBParameterGroup(ctx context.Context, req *connect.Request[pb.ModifyDBParameterGroupMessage]) (*connect.Response[pb.DBParameterGroupNameMessage], error) {
	return nil, nop("ModifyDBParameterGroup")
}

func (h *AdminHandler) ModifyDBSubnetGroup(ctx context.Context, req *connect.Request[pb.ModifyDBSubnetGroupMessage]) (*connect.Response[pb.ModifyDBSubnetGroupResult], error) {
	return nil, nop("ModifyDBSubnetGroup")
}

func (h *AdminHandler) ModifyEventSubscription(ctx context.Context, req *connect.Request[pb.ModifyEventSubscriptionMessage]) (*connect.Response[pb.ModifyEventSubscriptionResult], error) {
	return nil, nop("ModifyEventSubscription")
}

func (h *AdminHandler) ModifyGlobalCluster(ctx context.Context, req *connect.Request[pb.ModifyGlobalClusterMessage]) (*connect.Response[pb.ModifyGlobalClusterResult], error) {
	return nil, nop("ModifyGlobalCluster")
}

func (h *AdminHandler) PromoteReadReplicaDBCluster(ctx context.Context, req *connect.Request[pb.PromoteReadReplicaDBClusterMessage]) (*connect.Response[pb.PromoteReadReplicaDBClusterResult], error) {
	return nil, nop("PromoteReadReplicaDBCluster")
}

func (h *AdminHandler) RebootDBInstance(ctx context.Context, req *connect.Request[pb.RebootDBInstanceMessage]) (*connect.Response[pb.RebootDBInstanceResult], error) {
	return nil, nop("RebootDBInstance")
}

func (h *AdminHandler) RemoveFromGlobalCluster(ctx context.Context, req *connect.Request[pb.RemoveFromGlobalClusterMessage]) (*connect.Response[pb.RemoveFromGlobalClusterResult], error) {
	return nil, nop("RemoveFromGlobalCluster")
}

func (h *AdminHandler) RemoveRoleFromDBCluster(ctx context.Context, req *connect.Request[pb.RemoveRoleFromDBClusterMessage]) (*connect.Response[common.Empty], error) {
	return nil, nop("RemoveRoleFromDBCluster")
}

func (h *AdminHandler) RemoveSourceIdentifierFromSubscription(ctx context.Context, req *connect.Request[pb.RemoveSourceIdentifierFromSubscriptionMessage]) (*connect.Response[pb.RemoveSourceIdentifierFromSubscriptionResult], error) {
	return nil, nop("RemoveSourceIdentifierFromSubscription")
}

func (h *AdminHandler) RemoveTagsFromResource(ctx context.Context, req *connect.Request[pb.RemoveTagsFromResourceMessage]) (*connect.Response[common.Empty], error) {
	return nil, nop("RemoveTagsFromResource")
}

func (h *AdminHandler) ResetDBClusterParameterGroup(ctx context.Context, req *connect.Request[pb.ResetDBClusterParameterGroupMessage]) (*connect.Response[pb.DBClusterParameterGroupNameMessage], error) {
	return nil, nop("ResetDBClusterParameterGroup")
}

func (h *AdminHandler) ResetDBParameterGroup(ctx context.Context, req *connect.Request[pb.ResetDBParameterGroupMessage]) (*connect.Response[pb.DBParameterGroupNameMessage], error) {
	return nil, nop("ResetDBParameterGroup")
}

func (h *AdminHandler) RestoreDBClusterFromSnapshot(ctx context.Context, req *connect.Request[pb.RestoreDBClusterFromSnapshotMessage]) (*connect.Response[pb.RestoreDBClusterFromSnapshotResult], error) {
	return nil, nop("RestoreDBClusterFromSnapshot")
}

func (h *AdminHandler) RestoreDBClusterToPointInTime(ctx context.Context, req *connect.Request[pb.RestoreDBClusterToPointInTimeMessage]) (*connect.Response[pb.RestoreDBClusterToPointInTimeResult], error) {
	return nil, nop("RestoreDBClusterToPointInTime")
}

func (h *AdminHandler) StartDBCluster(ctx context.Context, req *connect.Request[pb.StartDBClusterMessage]) (*connect.Response[pb.StartDBClusterResult], error) {
	return nil, nop("StartDBCluster")
}

func (h *AdminHandler) StopDBCluster(ctx context.Context, req *connect.Request[pb.StopDBClusterMessage]) (*connect.Response[pb.StopDBClusterResult], error) {
	return nil, nop("StopDBCluster")
}

func (h *AdminHandler) SwitchoverGlobalCluster(ctx context.Context, req *connect.Request[pb.SwitchoverGlobalClusterMessage]) (*connect.Response[pb.SwitchoverGlobalClusterResult], error) {
	return nil, nop("SwitchoverGlobalCluster")
}
