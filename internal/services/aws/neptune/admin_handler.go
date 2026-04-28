package neptune

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/neptune"
	neptuneconnect "vorpalstacks/internal/pb/aws/neptune/neptuneconnect"
	storeneptune "vorpalstacks/internal/store/aws/neptune"
)

// AdminHandler implements the Neptune Management API gRPC-Web admin console
// handler. It exposes List/Describe operations for the Flutter management UI,
// delegating data access to the NeptuneStore via the shared NeptuneService
// per-region cache.
type AdminHandler struct {
	neptuneconnect.UnimplementedNeptuneServiceHandler
	service   *NeptuneService
	accountId string
}

var _ neptuneconnect.NeptuneServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Neptune admin console handler backed by the
// given service instance, ensuring the same per-region cached stores are used
// as the HTTP API handlers.
func NewAdminHandler(svc *NeptuneService, accountId string) *AdminHandler {
	return &AdminHandler{service: svc, accountId: accountId}
}

func (h *AdminHandler) getStore(header http.Header) (*storeneptune.NeptuneStore, error) {
	region := svccommon.GetRegionFromHeader(header)
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	s, ok := store.(*storeneptune.NeptuneStore)
	if !ok {
		return nil, fmt.Errorf("unexpected store type: %T", store)
	}
	return s, nil
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
				Engineversion:          "1.3.2.0",
				Dbparametergroupfamily: "neptune1",
			},
			{
				Engine:                 "neptune",
				Engineversion:          "1.3.1.0",
				Dbparametergroupfamily: "neptune1",
			},
			{
				Engine:                 "neptune",
				Engineversion:          "1.2.1.0",
				Dbparametergroupfamily: "neptune1",
			},
		},
	}), nil
}

// DescribeEventCategories returns the event categories for Neptune source types.
func (h *AdminHandler) DescribeEventCategories(ctx context.Context, req *connect.Request[pb.DescribeEventCategoriesMessage]) (*connect.Response[pb.EventCategoriesMessage], error) {
	return connect.NewResponse(&pb.EventCategoriesMessage{
		Eventcategoriesmaplist: []*pb.EventCategoriesMap{
			{Sourcetype: "db-cluster", Eventcategories: []string{"creation", "deletion", "failover", "failure", "maintenance", "notification", "read replica", "recovery", "restoration", "backup"}},
			{Sourcetype: "db-instance", Eventcategories: []string{"creation", "deletion", "failure", "maintenance", "notification", "recovery"}},
			{Sourcetype: "db-snapshot", Eventcategories: []string{"creation", "deletion", "restoration"}},
			{Sourcetype: "db-parameter-group", Eventcategories: []string{"creation", "modification", "deletion"}},
		},
	}), nil
}

// DescribeEvents returns Neptune events from the per-region store.
func (h *AdminHandler) DescribeEvents(ctx context.Context, req *connect.Request[pb.DescribeEventsMessage]) (*connect.Response[pb.EventsMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, err
	}
	result, err := store.ListEvents(storeneptune.EventListOptions{})
	if err != nil {
		return nil, err
	}
	events := make([]*pb.Event, 0, len(result.Events))
	for _, evt := range result.Events {
		var sourceType pb.SourceType
		switch evt.SourceType {
		case "db-cluster":
			sourceType = pb.SourceType_SOURCE_TYPE_DB_CLUSTER
		case "db-instance":
			sourceType = pb.SourceType_SOURCE_TYPE_DB_INSTANCE
		case "db-snapshot":
			sourceType = pb.SourceType_SOURCE_TYPE_DB_CLUSTER_SNAPSHOT
		case "db-parameter-group":
			sourceType = pb.SourceType_SOURCE_TYPE_DB_PARAMETER_GROUP
		default:
			sourceType = pb.SourceType_SOURCE_TYPE_DB_CLUSTER
		}
		events = append(events, &pb.Event{
			Date:             evt.Date.UTC().Format("2006-01-02T15:04:05.000Z"),
			Message:          evt.Message,
			Sourcearn:        evt.SourceArn,
			Sourceidentifier: evt.SourceIdentifier,
			Sourcetype:       sourceType,
			Eventcategories:  evt.EventCategories,
		})
	}
	return connect.NewResponse(&pb.EventsMessage{
		Events: events,
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
// parameter group, including system defaults and any user modifications.
func (h *AdminHandler) DescribeDBClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeDBClusterParametersMessage]) (*connect.Response[pb.DBClusterParameterGroupDetails], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, err
	}
	pg, err := store.GetClusterParameterGroup(req.Msg.Dbclusterparametergroupname)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	defaultParams := []struct{ name, value, desc, source, apply, dtype, modifiable string }{
		{"neptune_query_timeout", "120000", "Query execution timeout in milliseconds", "system", "dynamic", "integer", "true"},
		{"neptune_enable_audit_log", "0", "Enable audit logging", "system", "static", "boolean", "true"},
	}
	userMods := make(map[string]storeneptune.Parameter, len(pg.Parameters))
	for _, p := range pg.Parameters {
		userMods[p.ParameterName] = p
	}

	pbParams := make([]*pb.Parameter, 0, len(defaultParams))
	for _, dp := range defaultParams {
		if mod, ok := userMods[dp.name]; ok {
			pbParams = append(pbParams, &pb.Parameter{
				Parametername:  mod.ParameterName,
				Parametervalue: mod.ParameterValue,
				Description:    mod.Description,
				Source:         mod.Source,
				Applytype:      mod.ApplyType,
				Datatype:       mod.DataType,
				Ismodifiable:   mod.IsModifiable,
			})
			delete(userMods, dp.name)
		} else {
			pbParams = append(pbParams, &pb.Parameter{
				Parametername:  dp.name,
				Parametervalue: dp.value,
				Description:    dp.desc,
				Source:         dp.source,
				Applytype:      dp.apply,
				Datatype:       dp.dtype,
				Ismodifiable:   dp.modifiable == "true",
			})
		}
	}
	for _, p := range userMods {
		pbParams = append(pbParams, &pb.Parameter{
			Parametername:  p.ParameterName,
			Parametervalue: p.ParameterValue,
			Description:    p.Description,
			Source:         p.Source,
			Applytype:      p.ApplyType,
			Datatype:       p.DataType,
			Ismodifiable:   p.IsModifiable,
		})
	}

	return connect.NewResponse(&pb.DBClusterParameterGroupDetails{
		Parameters: pbParams,
		Marker:     "",
	}), nil
}

// DescribeDBParameters returns the parameters of a DB parameter group,
// including system defaults and any user modifications.
func (h *AdminHandler) DescribeDBParameters(ctx context.Context, req *connect.Request[pb.DescribeDBParametersMessage]) (*connect.Response[pb.DBParameterGroupDetails], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, err
	}
	pg, err := store.GetParameterGroup(req.Msg.Dbparametergroupname)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	defaultParams := []struct{ name, value, desc, source, apply, dtype, modifiable string }{
		{"neptune_query_timeout", "120000", "Query execution timeout", "system", "dynamic", "integer", "true"},
	}
	userMods := make(map[string]storeneptune.Parameter, len(pg.Parameters))
	for _, p := range pg.Parameters {
		userMods[p.ParameterName] = p
	}

	pbParams := make([]*pb.Parameter, 0, len(defaultParams))
	for _, dp := range defaultParams {
		if mod, ok := userMods[dp.name]; ok {
			pbParams = append(pbParams, &pb.Parameter{
				Parametername:  mod.ParameterName,
				Parametervalue: mod.ParameterValue,
				Description:    mod.Description,
				Source:         mod.Source,
				Applytype:      mod.ApplyType,
				Datatype:       mod.DataType,
				Ismodifiable:   mod.IsModifiable,
			})
			delete(userMods, dp.name)
		} else {
			pbParams = append(pbParams, &pb.Parameter{
				Parametername:  dp.name,
				Parametervalue: dp.value,
				Description:    dp.desc,
				Source:         dp.source,
				Applytype:      dp.apply,
				Datatype:       dp.dtype,
				Ismodifiable:   dp.modifiable == "true",
			})
		}
	}
	for _, p := range userMods {
		pbParams = append(pbParams, &pb.Parameter{
			Parametername:  p.ParameterName,
			Parametervalue: p.ParameterValue,
			Description:    p.Description,
			Source:         p.Source,
			Applytype:      p.ApplyType,
			Datatype:       p.DataType,
			Ismodifiable:   p.IsModifiable,
		})
	}

	return connect.NewResponse(&pb.DBParameterGroupDetails{
		Parameters: pbParams,
		Marker:     "",
	}), nil
}

// DescribeEngineDefaultClusterParameters returns the default engine
// parameters for a cluster parameter group family.
func (h *AdminHandler) DescribeEngineDefaultClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultClusterParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultClusterParametersResult], error) {
	pbParams := []*pb.Parameter{
		{Parametername: "neptune_query_timeout", Parametervalue: "120000", Description: "Query execution timeout in milliseconds", Source: "system", Applytype: "dynamic", Datatype: "integer", Ismodifiable: true},
		{Parametername: "neptune_enable_audit_log", Parametervalue: "0", Description: "Enable audit logging", Source: "system", Applytype: "static", Datatype: "boolean", Ismodifiable: true},
	}
	return connect.NewResponse(&pb.DescribeEngineDefaultClusterParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: "neptune1",
			Marker:                 "",
			Parameters:             pbParams,
		},
	}), nil
}

// DescribeEngineDefaultParameters returns the default engine parameters for
// a DB parameter group family.
func (h *AdminHandler) DescribeEngineDefaultParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultParametersResult], error) {
	pbParams := []*pb.Parameter{
		{Parametername: "neptune_query_timeout", Parametervalue: "120000", Description: "Query execution timeout", Source: "system", Applytype: "dynamic", Datatype: "integer", Ismodifiable: true},
	}
	return connect.NewResponse(&pb.DescribeEngineDefaultParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: "neptune1",
			Marker:                 "",
			Parameters:             pbParams,
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
		Enabledcloudwatchlogsexports:     c.EnabledCloudwatchLogsExports,
		Iamdatabaseauthenticationenabled: c.IAMDatabaseAuthenticationEnabled,
		Dbclusterarn:                     c.DBClusterArn,
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
		Dbinstancestatus:                 i.DBInstanceStatus,
		Availabilityzone:                 i.AvailabilityZone,
		Preferredmaintenancewindow:       i.PreferredMaintenanceWindow,
		Enabledcloudwatchlogsexports:     i.EnabledCloudwatchLogsExports,
		Iamdatabaseauthenticationenabled: i.IAMDatabaseAuthenticationEnabled,
		Publiclyaccessible:               i.PubliclyAccessible,
		Autominorversionupgrade:          i.AutoMinorVersionUpgrade,
		Copytagstosnapshot:               i.CopyTagsToSnapshot,
		Dbinstancearn:                    i.DBInstanceArn,
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

// NewConnectHandler creates a gRPC-Web connect handler for the Neptune admin console.
func NewConnectHandler(svc *NeptuneService, accountID string) (string, http.Handler) {
	return neptuneconnect.NewNeptuneServiceHandler(NewAdminHandler(svc, accountID))
}
