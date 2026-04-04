package neptune

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/neptune"
	neptuneconnect "vorpalstacks/internal/pb/aws/neptune/neptuneconnect"
	svccommon "vorpalstacks/internal/services/aws/common"
	storeneptune "vorpalstacks/internal/store/aws/neptune"
)

// AdminHandler implements the Neptune Management API gRPC-Web admin console
// handler. It exposes List/Describe operations for the Flutter management UI,
// delegating data access to the NeptuneStore via RegionStorageManager.
type AdminHandler struct {
	neptuneconnect.UnimplementedNeptuneServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ neptuneconnect.NeptuneServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Neptune admin console handler.
func NewAdminHandler(sm *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{storageManager: sm, accountId: accountId}
}

func (h *AdminHandler) getStore(header http.Header) (*storeneptune.NeptuneStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	rs, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return storeneptune.NewNeptuneStore(rs), nil
}

// DescribeDBClusters returns information about DB clusters, optionally filtered
// by DBClusterIdentifier.
func (h *AdminHandler) DescribeDBClusters(ctx context.Context, req *connect.Request[pb.DescribeDBClustersMessage]) (*connect.Response[pb.DBClusterMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	clusters, err := store.ListClusters()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbClusters := make([]*pb.DBCluster, 0, len(clusters))
	for _, c := range clusters {
		if req.Msg.Dbclusteridentifier != "" && c.DBClusterIdentifier != req.Msg.Dbclusteridentifier {
			continue
		}
		pbClusters = append(pbClusters, clusterToPb(c, h.accountId))
	}

	return connect.NewResponse(&pb.DBClusterMessage{
		Dbclusters: pbClusters,
	}), nil
}

// DescribeDBInstances returns information about DB instances, optionally filtered
// by DBInstanceIdentifier.
func (h *AdminHandler) DescribeDBInstances(ctx context.Context, req *connect.Request[pb.DescribeDBInstancesMessage]) (*connect.Response[pb.DBInstanceMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	instances, err := store.ListInstances()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbInstances := make([]*pb.DBInstance, 0, len(instances))
	for _, i := range instances {
		if req.Msg.Dbinstanceidentifier != "" && i.DBInstanceIdentifier != req.Msg.Dbinstanceidentifier {
			continue
		}
		pbInstances = append(pbInstances, instanceToPb(i, h.accountId))
	}

	return connect.NewResponse(&pb.DBInstanceMessage{
		Dbinstances: pbInstances,
	}), nil
}

// DescribeDBClusterSnapshots returns information about DB cluster snapshots,
// optionally filtered by snapshot or cluster identifier.
func (h *AdminHandler) DescribeDBClusterSnapshots(ctx context.Context, req *connect.Request[pb.DescribeDBClusterSnapshotsMessage]) (*connect.Response[pb.DBClusterSnapshotMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	snapshots, err := store.ListSnapshots()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbSnapshots := make([]*pb.DBClusterSnapshot, 0, len(snapshots))
	for _, s := range snapshots {
		if req.Msg.Dbclustersnapshotidentifier != "" && s.DBClusterSnapshotIdentifier != req.Msg.Dbclustersnapshotidentifier {
			continue
		}
		if req.Msg.Dbclusteridentifier != "" && s.DBClusterIdentifier != req.Msg.Dbclusteridentifier {
			continue
		}
		pbSnapshots = append(pbSnapshots, snapshotToPb(s, h.accountId))
	}

	return connect.NewResponse(&pb.DBClusterSnapshotMessage{
		Dbclustersnapshots: pbSnapshots,
	}), nil
}

// DescribeDBClusterParameterGroups returns information about DB cluster
// parameter groups, optionally filtered by name.
func (h *AdminHandler) DescribeDBClusterParameterGroups(ctx context.Context, req *connect.Request[pb.DescribeDBClusterParameterGroupsMessage]) (*connect.Response[pb.DBClusterParameterGroupsMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	groups, err := store.ListClusterParameterGroups()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbGroups := make([]*pb.DBClusterParameterGroup, 0, len(groups))
	for _, g := range groups {
		if req.Msg.Dbclusterparametergroupname != "" && g.DBClusterParameterGroupName != req.Msg.Dbclusterparametergroupname {
			continue
		}
		pbGroups = append(pbGroups, clusterParamGroupToPb(g))
	}

	return connect.NewResponse(&pb.DBClusterParameterGroupsMessage{
		Dbclusterparametergroups: pbGroups,
	}), nil
}

// DescribeDBParameterGroups returns information about DB parameter groups,
// optionally filtered by name.
func (h *AdminHandler) DescribeDBParameterGroups(ctx context.Context, req *connect.Request[pb.DescribeDBParameterGroupsMessage]) (*connect.Response[pb.DBParameterGroupsMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	groups, err := store.ListParameterGroups()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbGroups := make([]*pb.DBParameterGroup, 0, len(groups))
	for _, g := range groups {
		if req.Msg.Dbparametergroupname != "" && g.DBParameterGroupName != req.Msg.Dbparametergroupname {
			continue
		}
		pbGroups = append(pbGroups, paramGroupToPb(g))
	}

	return connect.NewResponse(&pb.DBParameterGroupsMessage{
		Dbparametergroups: pbGroups,
	}), nil
}

// DescribeDBSubnetGroups returns information about DB subnet groups,
// optionally filtered by name.
func (h *AdminHandler) DescribeDBSubnetGroups(ctx context.Context, req *connect.Request[pb.DescribeDBSubnetGroupsMessage]) (*connect.Response[pb.DBSubnetGroupMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	groups, err := store.ListSubnetGroups()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbGroups := make([]*pb.DBSubnetGroup, 0, len(groups))
	for _, g := range groups {
		if req.Msg.Dbsubnetgroupname != "" && g.DBSubnetGroupName != req.Msg.Dbsubnetgroupname {
			continue
		}
		pbGroups = append(pbGroups, subnetGroupToPb(g))
	}

	return connect.NewResponse(&pb.DBSubnetGroupMessage{
		Dbsubnetgroups: pbGroups,
	}), nil
}

// DescribeGlobalClusters returns information about global clusters,
// optionally filtered by global cluster identifier.
func (h *AdminHandler) DescribeGlobalClusters(ctx context.Context, req *connect.Request[pb.DescribeGlobalClustersMessage]) (*connect.Response[pb.GlobalClustersMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	clusters, err := store.ListGlobalClusters()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbClusters := make([]*pb.GlobalCluster, 0, len(clusters))
	for _, c := range clusters {
		if req.Msg.Globalclusteridentifier != "" && c.GlobalClusterIdentifier != req.Msg.Globalclusteridentifier {
			continue
		}
		pbClusters = append(pbClusters, globalClusterToPb(c))
	}

	return connect.NewResponse(&pb.GlobalClustersMessage{
		Globalclusters: pbClusters,
	}), nil
}

// DescribeEventSubscriptions returns information about event subscriptions,
// optionally filtered by subscription name.
func (h *AdminHandler) DescribeEventSubscriptions(ctx context.Context, req *connect.Request[pb.DescribeEventSubscriptionsMessage]) (*connect.Response[pb.EventSubscriptionsMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	subs, err := store.ListEventSubscriptions()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbSubs := make([]*pb.EventSubscription, 0, len(subs))
	for _, s := range subs {
		if req.Msg.Subscriptionname != "" && s.CustSubscriptionId != req.Msg.Subscriptionname {
			continue
		}
		pbSubs = append(pbSubs, eventSubscriptionToPb(s))
	}

	return connect.NewResponse(&pb.EventSubscriptionsMessage{
		Eventsubscriptionslist: pbSubs,
	}), nil
}

// DescribeDBClusterEndpoints returns information about cluster endpoints,
// optionally filtered by cluster or endpoint identifier.
func (h *AdminHandler) DescribeDBClusterEndpoints(ctx context.Context, req *connect.Request[pb.DescribeDBClusterEndpointsMessage]) (*connect.Response[pb.DBClusterEndpointMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	clusterID := req.Msg.Dbclusteridentifier
	endpoints, err := store.ListClusterEndpoints(clusterID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbEndpoints := make([]*pb.DBClusterEndpoint, 0, len(endpoints))
	for _, ep := range endpoints {
		if req.Msg.Dbclusterendpointidentifier != "" && ep.DBClusterEndpointIdentifier != req.Msg.Dbclusterendpointidentifier {
			continue
		}
		pbEndpoints = append(pbEndpoints, clusterEndpointToPb(ep))
	}

	return connect.NewResponse(&pb.DBClusterEndpointMessage{
		Dbclusterendpoints: pbEndpoints,
	}), nil
}

// ListTagsForResource returns the tags associated with a Neptune resource.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceMessage]) (*connect.Response[pb.TagListMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	tags, err := store.GetTags(req.Msg.Resourcename)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	pbTags := make([]*pb.Tag, len(tags))
	for i, t := range tags {
		pbTags[i] = &pb.Tag{Key: t.Key, Value: t.Value}
	}

	return connect.NewResponse(&pb.TagListMessage{
		Taglist: pbTags,
	}), nil
}

// DescribeDBEngineVersions returns the available Neptune engine versions.
func (h *AdminHandler) DescribeDBEngineVersions(ctx context.Context, req *connect.Request[pb.DescribeDBEngineVersionsMessage]) (*connect.Response[pb.DBEngineVersionMessage], error) {
	return connect.NewResponse(&pb.DBEngineVersionMessage{
		Dbengineversions: []*pb.DBEngineVersion{
			{
				Engine:                 "neptune",
				Engineversion:          "1.4.0.1",
				Dbparametergroupfamily: "neptune1.4",
			},
			{
				Engine:                 "neptune",
				Engineversion:          "1.3.2.1",
				Dbparametergroupfamily: "neptune1.3",
			},
			{
				Engine:                 "neptune",
				Engineversion:          "1.2.1.0",
				Dbparametergroupfamily: "neptune1.2",
			},
		},
	}), nil
}

// DescribeEventCategories returns the event categories for Neptune source types.
func (h *AdminHandler) DescribeEventCategories(ctx context.Context, req *connect.Request[pb.DescribeEventCategoriesMessage]) (*connect.Response[pb.EventCategoriesMessage], error) {
	return connect.NewResponse(&pb.EventCategoriesMessage{
		Eventcategoriesmaplist: []*pb.EventCategoriesMap{
			{Sourcetype: "db-cluster", Eventcategories: []string{"creation", "failure", "failover", "maintenance", "notification", "recovery"}},
			{Sourcetype: "db-instance", Eventcategories: []string{"creation", "failure", "maintenance", "notification", "recovery"}},
			{Sourcetype: "db-snapshot", Eventcategories: []string{"creation", "deletion"}},
			{Sourcetype: "db-parameter-group", Eventcategories: []string{"creation", "modification", "deletion"}},
		},
	}), nil
}

// DescribeEvents returns Neptune events. Currently returns an empty list as
// event history is not persisted.
func (h *AdminHandler) DescribeEvents(ctx context.Context, req *connect.Request[pb.DescribeEventsMessage]) (*connect.Response[pb.EventsMessage], error) {
	return connect.NewResponse(&pb.EventsMessage{
		Events: []*pb.Event{},
	}), nil
}

// DescribePendingMaintenanceActions returns pending maintenance actions.
// Currently returns an empty list as maintenance scheduling is not implemented.
func (h *AdminHandler) DescribePendingMaintenanceActions(ctx context.Context, req *connect.Request[pb.DescribePendingMaintenanceActionsMessage]) (*connect.Response[pb.PendingMaintenanceActionsMessage], error) {
	return connect.NewResponse(&pb.PendingMaintenanceActionsMessage{
		Pendingmaintenanceactions: []*pb.ResourcePendingMaintenanceActions{},
	}), nil
}

// DescribeOrderableDBInstanceOptions returns the available DB instance classes
// for Neptune.
func (h *AdminHandler) DescribeOrderableDBInstanceOptions(ctx context.Context, req *connect.Request[pb.DescribeOrderableDBInstanceOptionsMessage]) (*connect.Response[pb.OrderableDBInstanceOptionsMessage], error) {
	return connect.NewResponse(&pb.OrderableDBInstanceOptionsMessage{
		Orderabledbinstanceoptions: []*pb.OrderableDBInstanceOption{
			{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.t3.medium", Licensemodel: "bring-your-own-license", Vpc: true},
			{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.r5.large", Licensemodel: "bring-your-own-license", Vpc: true},
			{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.r5.xlarge", Licensemodel: "bring-your-own-license", Vpc: true},
		},
	}), nil
}

// DescribeValidDBInstanceModifications returns valid modifications for a DB
// instance. Currently returns an empty list.
func (h *AdminHandler) DescribeValidDBInstanceModifications(ctx context.Context, req *connect.Request[pb.DescribeValidDBInstanceModificationsMessage]) (*connect.Response[pb.DescribeValidDBInstanceModificationsResult], error) {
	return connect.NewResponse(&pb.DescribeValidDBInstanceModificationsResult{
		Validdbinstancemodificationsmessage: &pb.ValidDBInstanceModificationsMessage{
			Storage: []*pb.ValidStorageOptions{},
		},
	}), nil
}

// DescribeDBClusterParameters returns the parameters of a DB cluster
// parameter group. Currently returns an empty parameter list.
func (h *AdminHandler) DescribeDBClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeDBClusterParametersMessage]) (*connect.Response[pb.DBClusterParameterGroupDetails], error) {
	return connect.NewResponse(&pb.DBClusterParameterGroupDetails{
		Parameters: []*pb.Parameter{},
		Marker:     "",
	}), nil
}

// DescribeDBParameters returns the parameters of a DB parameter group.
// Currently returns an empty parameter list.
func (h *AdminHandler) DescribeDBParameters(ctx context.Context, req *connect.Request[pb.DescribeDBParametersMessage]) (*connect.Response[pb.DBParameterGroupDetails], error) {
	return connect.NewResponse(&pb.DBParameterGroupDetails{
		Parameters: []*pb.Parameter{},
		Marker:     "",
	}), nil
}

// DescribeEngineDefaultClusterParameters returns the default engine
// parameters for a cluster parameter group family.
func (h *AdminHandler) DescribeEngineDefaultClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultClusterParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultClusterParametersResult], error) {
	return connect.NewResponse(&pb.DescribeEngineDefaultClusterParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: "neptune1.4",
			Marker:                 "",
			Parameters:             []*pb.Parameter{},
		},
	}), nil
}

// DescribeEngineDefaultParameters returns the default engine parameters for
// a DB parameter group family.
func (h *AdminHandler) DescribeEngineDefaultParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultParametersResult], error) {
	return connect.NewResponse(&pb.DescribeEngineDefaultParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: "neptune1.4",
			Marker:                 "",
			Parameters:             []*pb.Parameter{},
		},
	}), nil
}

// DescribeDBClusterSnapshotAttributes returns the attributes of a DB cluster
// snapshot. Currently returns an empty attribute list.
func (h *AdminHandler) DescribeDBClusterSnapshotAttributes(ctx context.Context, req *connect.Request[pb.DescribeDBClusterSnapshotAttributesMessage]) (*connect.Response[pb.DescribeDBClusterSnapshotAttributesResult], error) {
	return connect.NewResponse(&pb.DescribeDBClusterSnapshotAttributesResult{
		Dbclustersnapshotattributesresult: &pb.DBClusterSnapshotAttributesResult{
			Dbclustersnapshotattributes: []*pb.DBClusterSnapshotAttribute{},
		},
	}), nil
}

// clusterToPb converts a domain DBCluster to the AWS API protobuf DBCluster.
func clusterToPb(c *storeneptune.DBCluster, accountId string) *pb.DBCluster {
	p := &pb.DBCluster{
		Dbclusteridentifier:              c.DBClusterIdentifier,
		Engine:                           c.Engine,
		Engineversion:                    c.EngineVersion,
		Status:                           c.Status,
		Masterusername:                   c.MasterUsername,
		Databasename:                     c.DatabaseName,
		Port:                             int32(c.Port),
		Backupretentionperiod:            int32(c.BackupRetentionPeriod),
		Preferredbackupwindow:            c.PreferredBackupWindow,
		Preferredmaintenancewindow:       c.PreferredMaintenanceWindow,
		Multiaz:                          c.MultiAZ,
		Dbclusterparametergroup:          c.DBClusterParameterGroupName,
		Dbsubnetgroup:                    c.DBSubnetGroupName,
		Storageencrypted:                 c.StorageEncrypted,
		Kmskeyid:                         c.KmsKeyId,
		Copytagstosnapshot:               c.CopyTagsToSnapshot,
		Deletionprotection:               c.DeletionProtection,
		Enabledcloudwatchlogsexports:     c.EnableCloudwatchLogsExports,
		Iamdatabaseauthenticationenabled: c.EnableIAMDatabaseAuthentication,
		Dbclusterarn:                     c.ARN(accountId, c.Region),
		Replicationsourceidentifier:      c.ReplicationSourceIdentifier,
		Globalclusteridentifier:          c.GlobalClusterIdentifier,
		Storagetype:                      c.StorageType,
		Availabilityzones:                c.AvailabilityZones,
	}
	if c.ClusterCreateTime != nil {
		p.Clustercreatetime = c.ClusterCreateTime.Format("2006-01-02T15:04:05.000Z")
	}
	if c.EarliestRestorableTime != nil {
		p.Earliestrestorabletime = c.EarliestRestorableTime.Format("2006-01-02T15:04:05.000Z")
	}
	if c.LatestRestorableTime != nil {
		p.Latestrestorabletime = c.LatestRestorableTime.Format("2006-01-02T15:04:05.000Z")
	}
	if c.ServerlessV2ScalingConfiguration != nil {
		p.Serverlessv2Scalingconfiguration = &pb.ServerlessV2ScalingConfigurationInfo{
			Mincapacity: c.ServerlessV2ScalingConfiguration.MinCapacity,
			Maxcapacity: c.ServerlessV2ScalingConfiguration.MaxCapacity,
		}
	}
	for _, r := range c.AssociatedRoles {
		p.Associatedroles = append(p.Associatedroles, &pb.DBClusterRole{
			Rolearn:     r.RoleArn,
			Featurename: r.FeatureName,
			Status:      r.Status,
		})
	}
	return p
}

// instanceToPb converts a domain DBInstance to the AWS API protobuf DBInstance.
func instanceToPb(i *storeneptune.DBInstance, accountId string) *pb.DBInstance {
	p := &pb.DBInstance{
		Dbinstanceidentifier:             i.DBInstanceIdentifier,
		Dbclusteridentifier:              i.DBClusterIdentifier,
		Engine:                           i.Engine,
		Engineversion:                    i.EngineVersion,
		Dbinstanceclass:                  i.DBInstanceClass,
		Dbinstancestatus:                 i.Status,
		Availabilityzone:                 i.AvailabilityZone,
		Preferredmaintenancewindow:       i.PreferredMaintenanceWindow,
		Enabledcloudwatchlogsexports:     i.EnableCloudwatchLogsExports,
		Iamdatabaseauthenticationenabled: i.EnableIAMDatabaseAuthentication,
		Publiclyaccessible:               i.PubliclyAccessible,
		Autominorversionupgrade:          i.AutoMinorVersionUpgrade,
		Copytagstosnapshot:               i.CopyTagsToSnapshot,
		Dbinstancearn:                    i.ARN(accountId, i.Region),
	}
	if i.InstanceCreateTime != nil {
		p.Instancecreatetime = i.InstanceCreateTime.Format("2006-01-02T15:04:05.000Z")
	}
	return p
}

// snapshotToPb converts a domain DBClusterSnapshot to the AWS API protobuf type.
func snapshotToPb(s *storeneptune.DBClusterSnapshot, accountId string) *pb.DBClusterSnapshot {
	p := &pb.DBClusterSnapshot{
		Dbclustersnapshotidentifier: s.DBClusterSnapshotIdentifier,
		Dbclusteridentifier:         s.DBClusterIdentifier,
		Engine:                      s.Engine,
		Engineversion:               s.EngineVersion,
		Status:                      s.Status,
		Port:                        int32(s.Port),
		Vpcid:                       s.VpcId,
		Storageencrypted:            s.StorageEncrypted,
		Kmskeyid:                    s.KmsKeyId,
		Dbclustersnapshotarn:        s.DBSnapshotArn,
	}
	if s.SnapshotCreateTime != nil {
		p.Snapshotcreatetime = s.SnapshotCreateTime.Format("2006-01-02T15:04:05.000Z")
	}
	if s.ClusterCreateTime != nil {
		p.Clustercreatetime = s.ClusterCreateTime.Format("2006-01-02T15:04:05.000Z")
	}
	return p
}

// clusterParamGroupToPb converts a domain DBClusterParameterGroup to the
// AWS API protobuf type.
func clusterParamGroupToPb(g *storeneptune.DBClusterParameterGroup) *pb.DBClusterParameterGroup {
	return &pb.DBClusterParameterGroup{
		Dbclusterparametergroupname: g.DBClusterParameterGroupName,
		Dbparametergroupfamily:      g.DBParameterGroupFamily,
		Description:                 g.Description,
		Dbclusterparametergrouparn:  g.ARN,
	}
}

// paramGroupToPb converts a domain DBParameterGroup to the AWS API protobuf type.
func paramGroupToPb(g *storeneptune.DBParameterGroup) *pb.DBParameterGroup {
	return &pb.DBParameterGroup{
		Dbparametergroupname:   g.DBParameterGroupName,
		Dbparametergroupfamily: g.DBParameterGroupFamily,
		Description:            g.Description,
		Dbparametergrouparn:    g.ARN,
	}
}

// subnetGroupToPb converts a domain DBSubnetGroup to the AWS API protobuf type.
func subnetGroupToPb(g *storeneptune.DBSubnetGroup) *pb.DBSubnetGroup {
	p := &pb.DBSubnetGroup{
		Dbsubnetgroupname:        g.DBSubnetGroupName,
		Dbsubnetgroupdescription: g.DBSubnetGroupDescription,
		Vpcid:                    g.VpcId,
		Subnetgroupstatus:        g.SubnetGroupStatus,
		Dbsubnetgrouparn:         g.ARN,
	}
	for _, s := range g.Subnets {
		p.Subnets = append(p.Subnets, &pb.Subnet{
			Subnetidentifier:       s.SubnetIdentifier,
			Subnetavailabilityzone: &pb.AvailabilityZone{Name: s.SubnetAvailabilityZone},
			Subnetstatus:           s.SubnetStatus,
		})
	}
	return p
}

// globalClusterToPb converts a domain GlobalCluster to the AWS API protobuf type.
func globalClusterToPb(c *storeneptune.GlobalCluster) *pb.GlobalCluster {
	p := &pb.GlobalCluster{
		Globalclusteridentifier: c.GlobalClusterIdentifier,
		Globalclusterresourceid: c.GlobalClusterResourceId,
		Globalclusterarn:        c.GlobalClusterArn,
		Engine:                  c.Engine,
		Engineversion:           c.EngineVersion,
		Status:                  c.Status,
		Storageencrypted:        c.StorageEncrypted,
		Deletionprotection:      c.DeletionProtection,
	}
	for _, m := range c.GlobalClusterMembers {
		p.Globalclustermembers = append(p.Globalclustermembers, &pb.GlobalClusterMember{
			Dbclusterarn: m.DBClusterArn,
			Iswriter:     m.IsWriter,
			Readers:      m.Readers,
		})
	}
	return p
}

// eventSubscriptionToPb converts a domain EventSubscription to the AWS API
// protobuf type.
func eventSubscriptionToPb(s *storeneptune.EventSubscription) *pb.EventSubscription {
	p := &pb.EventSubscription{
		Custsubscriptionid:   s.CustSubscriptionId,
		Snstopicarn:          s.SnsTopicArn,
		Status:               s.Status,
		Sourcetype:           s.SourceType,
		Sourceidslist:        s.SourceIdsList,
		Eventcategorieslist:  s.EventCategoriesList,
		Enabled:              s.Enabled,
		Eventsubscriptionarn: s.CustSubscriptionArn,
	}
	return p
}

// clusterEndpointToPb converts a domain DBClusterEndpoint to the AWS API
// protobuf type.
func clusterEndpointToPb(ep *storeneptune.DBClusterEndpoint) *pb.DBClusterEndpoint {
	return &pb.DBClusterEndpoint{
		Dbclusterendpointidentifier: ep.DBClusterEndpointIdentifier,
		Dbclusteridentifier:         ep.DBClusterIdentifier,
		Endpoint:                    ep.Endpoint,
		Status:                      ep.Status,
		Endpointtype:                ep.EndpointType,
		Excludedmembers:             ep.ExcludedMembers,
		Staticmembers:               ep.StaticMembers,
		Dbclusterendpointarn:        ep.DBClusterEndpointArn,
	}
}
