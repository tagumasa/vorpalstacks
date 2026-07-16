// Package rds provides the RDS admin console handler that serves management
// RPCs for both Neptune and MySQL database engines via the gRPC-Web admin
// interface. It delegates data access to the common RDS store layer.
package rds

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/rds"
	rdsconnect "vorpalstacks/internal/pb/aws/rds/rdsconnect"
	storerds "vorpalstacks/internal/store/aws/rds"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// StoreProvider returns the RDS store for a given region. This decouples the
// admin handler from any specific service (NeptuneService, etc.), allowing it
// to serve data from all RDS engines through the common store interface.
type StoreProvider func(region string) (storerds.StoreInterface, error)

// EngineProvider supplies per-engine lifecycle managers keyed by engine type.
type EngineProvider func(engineType string) (Engine, error)

// AdminHandler implements the RDS Management API gRPC-Web admin console
// handler. It exposes Describe/Create/Delete operations for the management UI,
// serving both Neptune and MySQL database engines through the common RDS store
// interface.
type AdminHandler struct {
	rdsconnect.UnimplementedRDSServiceHandler
	stores    StoreProvider
	engines   EngineProvider
	accountId string
}

var _ rdsconnect.RDSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new RDS admin console handler backed by the given
// store provider and engine provider.
func NewAdminHandler(stores StoreProvider, engines EngineProvider, accountId string) *AdminHandler {
	return &AdminHandler{stores: stores, engines: engines, accountId: accountId}
}

func (h *AdminHandler) getStore(header http.Header) (storerds.StoreInterface, error) {
	region := svccommon.GetRegionFromHeader(header)
	return h.stores(region)
}

// DescribeDBClusters returns information about DB clusters, optionally filtered
// by DBClusterIdentifier.
func (h *AdminHandler) DescribeDBClusters(ctx context.Context, req *connect.Request[pb.DescribeDBClustersMessage]) (*connect.Response[pb.DBClusterMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	clusters, err := store.ListClusters()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	instances, err := store.ListInstances()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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

// CreateDBInstance creates a new DB instance. For MySQL engine instances the
// vmysql engine is started on a dynamically allocated port.
func (h *AdminHandler) CreateDBInstance(ctx context.Context, req *connect.Request[pb.CreateDBInstanceMessage]) (*connect.Response[pb.CreateDBInstanceResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	id := req.Msg.Dbinstanceidentifier
	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DBInstanceIdentifier is required"))
	}
	engine := req.Msg.Engine
	if engine == "" {
		engine = "mysql"
	}

	now := time.Now()
	instance := &storerds.DBInstance{
		DBInstanceIdentifier:       id,
		DBClusterIdentifier:        req.Msg.Dbclusteridentifier,
		Engine:                     engine,
		EngineVersion:              req.Msg.Engineversion,
		DBInstanceClass:            req.Msg.Dbinstanceclass,
		DBInstanceStatus:           "creating",
		AvailabilityZone:           req.Msg.Availabilityzone,
		PreferredMaintenanceWindow: req.Msg.Preferredmaintenancewindow,
		DBParameterGroupName:       req.Msg.Dbparametergroupname,
		DBSubnetGroupName:          req.Msg.Dbsubnetgroupname,
		PubliclyAccessible:         req.Msg.GetPubliclyaccessible(),
		AutoMinorVersionUpgrade:    req.Msg.GetAutominorversionupgrade(),
		InstanceCreateTime:         &now,
		AccountID:                  h.accountId,
		Region:                     region,
		DBInstanceArn:              arnutil.NewARNBuilder(h.accountId, region).RDS().DBInstance(id),
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	engineStarted := false
	if eng, engErr := h.engines(engine); engErr == nil {
		if port, openErr := eng.Open(region, id); openErr != nil {
			logs.Warn("rds-admin: failed to start engine for instance",
				logs.String("instance", id), logs.Err(openErr))
		} else {
			engineStarted = true
			instance.Endpoint = &storerds.Endpoint{
				Address: fmt.Sprintf("%s.%s.%s.rds.amazonaws.com", id, h.accountId, region),
				Port:    port,
			}
			if err := store.UpdateInstance(instance); err != nil {
				logs.Warn("rds-admin: failed to persist instance endpoint",
					logs.String("instance", id), logs.Err(err))
			}
		}
	}

	if engineStarted {
		instance.DBInstanceStatus = "available"
	} else {
		instance.DBInstanceStatus = "failed"
	}
	if err := store.UpdateInstance(instance); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDBInstanceResult{
		Dbinstance: instanceToPb(instance, h.accountId),
	}), nil
}

// DeleteDBInstance deletes a DB instance and stops its engine if running.
func (h *AdminHandler) DeleteDBInstance(ctx context.Context, req *connect.Request[pb.DeleteDBInstanceMessage]) (*connect.Response[pb.DeleteDBInstanceResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	id := req.Msg.Dbinstanceidentifier
	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DBInstanceIdentifier is required"))
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if eng, engErr := h.engines(instance.Engine); engErr == nil {
		eng.Close(id)
	}

	instance.DBInstanceStatus = "deleting"
	if err := store.DeleteInstance(id); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteDBInstanceResult{
		Dbinstance: instanceToPb(instance, h.accountId),
	}), nil
}

// CreateDBCluster creates a new DB cluster.
func (h *AdminHandler) CreateDBCluster(ctx context.Context, req *connect.Request[pb.CreateDBClusterMessage]) (*connect.Response[pb.CreateDBClusterResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	id := req.Msg.Dbclusteridentifier
	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DBClusterIdentifier is required"))
	}
	engine := req.Msg.Engine
	if engine == "" {
		engine = "aurora-mysql"
	}

	now := time.Now()
	cluster := &storerds.DBCluster{
		DBClusterIdentifier:              id,
		Engine:                           engine,
		EngineVersion:                    req.Msg.Engineversion,
		Status:                           "creating",
		MasterUsername:                   req.Msg.Masterusername,
		DatabaseName:                     req.Msg.Databasename,
		Port:                             int(req.Msg.Port),
		BackupRetentionPeriod:            int(req.Msg.Backupretentionperiod),
		AvailabilityZones:                req.Msg.Availabilityzones,
		DBSubnetGroupName:                req.Msg.Dbsubnetgroupname,
		DBClusterParameterGroupName:      req.Msg.Dbclusterparametergroupname,
		StorageEncrypted:                 req.Msg.GetStorageencrypted(),
		CopyTagsToSnapshot:               req.Msg.GetCopytagstosnapshot(),
		DeletionProtection:               req.Msg.GetDeletionprotection(),
		EnabledCloudwatchLogsExports:     req.Msg.Enablecloudwatchlogsexports,
		IAMDatabaseAuthenticationEnabled: req.Msg.GetEnableiamdatabaseauthentication(),
		ClusterCreateTime:                &now,
		AccountID:                        h.accountId,
		Region:                           region,
		DBClusterArn:                     arnutil.NewARNBuilder(h.accountId, region).RDS().Cluster(id),
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDBClusterResult{
		Dbcluster: clusterToPb(cluster, h.accountId),
	}), nil
}

// DeleteDBCluster deletes a DB cluster.
func (h *AdminHandler) DeleteDBCluster(ctx context.Context, req *connect.Request[pb.DeleteDBClusterMessage]) (*connect.Response[pb.DeleteDBClusterResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	id := req.Msg.Dbclusteridentifier
	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DBClusterIdentifier is required"))
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if cluster.DeletionProtection {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("cannot delete cluster when DeletionProtection is enabled"))
	}

	cluster.Status = "deleting"
	_ = store.UpdateCluster(cluster)

	if err := store.DeleteCluster(id); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteDBClusterResult{
		Dbcluster: clusterToPb(cluster, h.accountId),
	}), nil
}

// DescribeDBClusterSnapshots returns information about DB cluster snapshots,
// optionally filtered by snapshot or cluster identifier.
func (h *AdminHandler) DescribeDBClusterSnapshots(ctx context.Context, req *connect.Request[pb.DescribeDBClusterSnapshotsMessage]) (*connect.Response[pb.DBClusterSnapshotMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	snapshots, err := store.ListSnapshots()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	groups, err := store.ListClusterParameterGroups()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	groups, err := store.ListParameterGroups()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	groups, err := store.ListSubnetGroups()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	clusters, err := store.ListGlobalClusters()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	subs, err := store.ListEventSubscriptions()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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

// ListTagsForResource returns the tags associated with an RDS resource.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceMessage]) (*connect.Response[pb.TagListMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	tags, err := store.GetTags(req.Msg.Resourcename)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	pbTags := make([]*pb.Tag, len(tags))
	for i, t := range tags {
		pbTags[i] = &pb.Tag{Key: t.Key, Value: t.Value}
	}

	return connect.NewResponse(&pb.TagListMessage{
		Taglist: pbTags,
	}), nil
}

// AddTagsToResource adds metadata tags to an RDS resource.
func (h *AdminHandler) AddTagsToResource(ctx context.Context, req *connect.Request[pb.AddTagsToResourceMessage]) (*connect.Response[pbcommon.Empty], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	tags := make([]types.Tag, len(req.Msg.Tags))
	for i, t := range req.Msg.Tags {
		tags[i] = types.Tag{Key: t.Key, Value: t.Value}
	}

	if err := store.AddTags(req.Msg.Resourcename, tags); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// RemoveTagsFromResource removes metadata tags from an RDS resource.
func (h *AdminHandler) RemoveTagsFromResource(ctx context.Context, req *connect.Request[pb.RemoveTagsFromResourceMessage]) (*connect.Response[pbcommon.Empty], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.RemoveTags(req.Msg.Resourcename, req.Msg.Tagkeys); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DescribeDBEngineVersions returns the available RDS engine versions for both
// Neptune and MySQL engines.
func (h *AdminHandler) DescribeDBEngineVersions(ctx context.Context, req *connect.Request[pb.DescribeDBEngineVersionsMessage]) (*connect.Response[pb.DBEngineVersionMessage], error) {
	versions := []*pb.DBEngineVersion{
		// Neptune engine versions
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
		// MySQL engine versions
		{
			Engine:                     "mysql",
			Engineversion:              "8.0.40",
			Dbparametergroupfamily:     "mysql8.0",
			Dbenginedescription:        "MySQL 8.0",
			Dbengineversiondescription: "MySQL 8.0.40",
		},
		{
			Engine:                     "mysql",
			Engineversion:              "8.0.39",
			Dbparametergroupfamily:     "mysql8.0",
			Dbenginedescription:        "MySQL 8.0",
			Dbengineversiondescription: "MySQL 8.0.39",
		},
		{
			Engine:                     "mysql",
			Engineversion:              "8.4.3",
			Dbparametergroupfamily:     "mysql8.4",
			Dbenginedescription:        "MySQL 8.4",
			Dbengineversiondescription: "MySQL 8.4.3",
		},
	}

	// Filter by engine if requested
	if req.Msg.Engine != "" {
		filtered := make([]*pb.DBEngineVersion, 0)
		for _, v := range versions {
			if v.Engine == req.Msg.Engine {
				filtered = append(filtered, v)
			}
		}
		versions = filtered
	}

	return connect.NewResponse(&pb.DBEngineVersionMessage{
		Dbengineversions: versions,
	}), nil
}

// DescribeEventCategories returns the event categories for RDS source types.
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

// DescribeEvents returns RDS events from the per-region store.
func (h *AdminHandler) DescribeEvents(ctx context.Context, req *connect.Request[pb.DescribeEventsMessage]) (*connect.Response[pb.EventsMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	result, err := store.ListEvents(storerds.EventListOptions{})
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
			Date:             evt.Date.UTC().Format(timeutils.ISO8601UTCFormat),
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
// for RDS engines.
func (h *AdminHandler) DescribeOrderableDBInstanceOptions(ctx context.Context, req *connect.Request[pb.DescribeOrderableDBInstanceOptionsMessage]) (*connect.Response[pb.OrderableDBInstanceOptionsMessage], error) {
	options := []*pb.OrderableDBInstanceOption{
		// Neptune options
		{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.t3.medium", Licensemodel: "bring-your-own-license", Vpc: proto.Bool(true)},
		{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.r5.large", Licensemodel: "bring-your-own-license", Vpc: proto.Bool(true)},
		{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.r5.xlarge", Licensemodel: "bring-your-own-license", Vpc: proto.Bool(true)},
		// MySQL options
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.t3.micro", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.t3.small", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.t3.medium", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.r5.large", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.r5.xlarge", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
	}

	// Filter by engine if requested
	if req.Msg.Engine != "" {
		filtered := make([]*pb.OrderableDBInstanceOption, 0)
		for _, o := range options {
			if o.Engine == req.Msg.Engine {
				filtered = append(filtered, o)
			}
		}
		options = filtered
	}

	return connect.NewResponse(&pb.OrderableDBInstanceOptionsMessage{
		Orderabledbinstanceoptions: options,
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

// DescribeDBClusterSnapshotAttributes returns the attributes of a DB cluster
// snapshot. Currently returns an empty attribute list.
func (h *AdminHandler) DescribeDBClusterSnapshotAttributes(ctx context.Context, req *connect.Request[pb.DescribeDBClusterSnapshotAttributesMessage]) (*connect.Response[pb.DescribeDBClusterSnapshotAttributesResult], error) {
	return connect.NewResponse(&pb.DescribeDBClusterSnapshotAttributesResult{
		Dbclustersnapshotattributesresult: &pb.DBClusterSnapshotAttributesResult{
			Dbclustersnapshotattributes: []*pb.DBClusterSnapshotAttribute{},
		},
	}), nil
}

// DescribeDBClusterEndpoints returns cluster endpoints filtered by cluster or
// endpoint identifier.
func (h *AdminHandler) DescribeDBClusterEndpoints(ctx context.Context, req *connect.Request[pb.DescribeDBClusterEndpointsMessage]) (*connect.Response[pb.DBClusterEndpointMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	clusterID := req.Msg.Dbclusteridentifier
	endpoints, err := store.ListClusterEndpoints(clusterID)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	pbEndpoints := make([]*pb.DBClusterEndpoint, 0, len(endpoints))
	for _, ep := range endpoints {
		if req.Msg.Dbclusterendpointidentifier != "" && ep.DBClusterEndpointIdentifier != req.Msg.Dbclusterendpointidentifier {
			continue
		}
		pbEndpoints = append(pbEndpoints, &pb.DBClusterEndpoint{
			Dbclusterendpointidentifier: ep.DBClusterEndpointIdentifier,
			Dbclusteridentifier:         ep.DBClusterIdentifier,
			Endpoint:                    ep.Endpoint,
			Status:                      ep.Status,
			Endpointtype:                ep.EndpointType,
			Excludedmembers:             ep.ExcludedMembers,
			Staticmembers:               ep.StaticMembers,
			Dbclusterendpointarn:        ep.DBClusterEndpointArn,
		})
	}

	return connect.NewResponse(&pb.DBClusterEndpointMessage{
		Dbclusterendpoints: pbEndpoints,
	}), nil
}

// DescribeDBClusterParameters returns the parameters of a DB cluster parameter
// group, including system defaults and user modifications.
func (h *AdminHandler) DescribeDBClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeDBClusterParametersMessage]) (*connect.Response[pb.DBClusterParameterGroupDetails], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	pg, err := store.GetClusterParameterGroup(req.Msg.Dbclusterparametergroupname)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	defaultParams := []struct{ name, value, desc, source, apply, dtype, modifiable string }{
		{"neptune_query_timeout", "120000", "Query execution timeout in milliseconds", "system", "dynamic", "integer", "true"},
		{"neptune_enable_audit_log", "0", "Enable audit logging", "system", "static", "boolean", "true"},
	}
	userMods := make(map[string]storerds.Parameter, len(pg.Parameters))
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
				Ismodifiable:   proto.Bool(mod.IsModifiable),
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
				Ismodifiable:   proto.Bool(dp.modifiable == "true"),
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
			Ismodifiable:   proto.Bool(p.IsModifiable),
		})
	}
	sort.Slice(pbParams, func(i, j int) bool { return pbParams[i].Parametername < pbParams[j].Parametername })

	return connect.NewResponse(&pb.DBClusterParameterGroupDetails{
		Parameters: pbParams,
		Marker:     "",
	}), nil
}

// DescribeDBParameters returns the parameters of a DB parameter group.
func (h *AdminHandler) DescribeDBParameters(ctx context.Context, req *connect.Request[pb.DescribeDBParametersMessage]) (*connect.Response[pb.DBParameterGroupDetails], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	pg, err := store.GetParameterGroup(req.Msg.Dbparametergroupname)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	defaultParams := []struct{ name, value, desc, source, apply, dtype, modifiable string }{
		{"neptune_query_timeout", "120000", "Query execution timeout", "system", "dynamic", "integer", "true"},
	}
	userMods := make(map[string]storerds.Parameter, len(pg.Parameters))
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
				Ismodifiable:   proto.Bool(mod.IsModifiable),
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
				Ismodifiable:   proto.Bool(dp.modifiable == "true"),
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
			Ismodifiable:   proto.Bool(p.IsModifiable),
		})
	}
	sort.Slice(pbParams, func(i, j int) bool { return pbParams[i].Parametername < pbParams[j].Parametername })

	return connect.NewResponse(&pb.DBParameterGroupDetails{
		Parameters: pbParams,
		Marker:     "",
	}), nil
}

// DescribeEngineDefaultClusterParameters returns the default engine parameters
// for a cluster parameter group family.
func (h *AdminHandler) DescribeEngineDefaultClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultClusterParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultClusterParametersResult], error) {
	return connect.NewResponse(&pb.DescribeEngineDefaultClusterParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: "neptune1",
			Parameters: []*pb.Parameter{
				{Parametername: "neptune_query_timeout", Parametervalue: "120000", Description: "Query execution timeout in milliseconds", Source: "system", Applytype: "dynamic", Datatype: "integer", Ismodifiable: proto.Bool(true)},
				{Parametername: "neptune_enable_audit_log", Parametervalue: "0", Description: "Enable audit logging", Source: "system", Applytype: "static", Datatype: "boolean", Ismodifiable: proto.Bool(true)},
			},
		},
	}), nil
}

// DescribeEngineDefaultParameters returns the default engine parameters for a
// DB parameter group family.
func (h *AdminHandler) DescribeEngineDefaultParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultParametersResult], error) {
	return connect.NewResponse(&pb.DescribeEngineDefaultParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: "neptune1",
			Parameters: []*pb.Parameter{
				{Parametername: "neptune_query_timeout", Parametervalue: "120000", Description: "Query execution timeout", Source: "system", Applytype: "dynamic", Datatype: "integer", Ismodifiable: proto.Bool(true)},
			},
		},
	}), nil
}

// --- Domain-to-proto conversion helpers ---

// clusterToPb converts a domain DBCluster to the RDS API protobuf DBCluster.
func clusterToPb(c *storerds.DBCluster, accountId string) *pb.DBCluster {
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
		Multiaz:                          proto.Bool(c.MultiAZ),
		Dbclusterparametergroup:          c.DBClusterParameterGroupName,
		Dbsubnetgroup:                    c.DBSubnetGroupName,
		Storageencrypted:                 proto.Bool(c.StorageEncrypted),
		Kmskeyid:                         c.KmsKeyId,
		Copytagstosnapshot:               proto.Bool(c.CopyTagsToSnapshot),
		Deletionprotection:               proto.Bool(c.DeletionProtection),
		Enabledcloudwatchlogsexports:     c.EnabledCloudwatchLogsExports,
		Iamdatabaseauthenticationenabled: proto.Bool(c.IAMDatabaseAuthenticationEnabled),
		Dbclusterarn:                     c.DBClusterArn,
		Replicationsourceidentifier:      c.ReplicationSourceIdentifier,
		Globalclusteridentifier:          c.GlobalClusterIdentifier,
		Storagetype:                      c.StorageType,
		Availabilityzones:                c.AvailabilityZones,
	}
	if c.ClusterCreateTime != nil {
		p.Clustercreatetime = c.ClusterCreateTime.Format(timeutils.ISO8601UTCFormat)
	}
	if c.EarliestRestorableTime != nil {
		p.Earliestrestorabletime = c.EarliestRestorableTime.Format(timeutils.ISO8601UTCFormat)
	}
	if c.LatestRestorableTime != nil {
		p.Latestrestorabletime = c.LatestRestorableTime.Format(timeutils.ISO8601UTCFormat)
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

// instanceToPb converts a domain DBInstance to the RDS API protobuf DBInstance.
func instanceToPb(i *storerds.DBInstance, accountId string) *pb.DBInstance {
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
		Iamdatabaseauthenticationenabled: proto.Bool(i.IAMDatabaseAuthenticationEnabled),
		Publiclyaccessible:               proto.Bool(i.PubliclyAccessible),
		Autominorversionupgrade:          proto.Bool(i.AutoMinorVersionUpgrade),
		Copytagstosnapshot:               proto.Bool(i.CopyTagsToSnapshot),
		Dbinstancearn:                    i.DBInstanceArn,
	}
	if i.InstanceCreateTime != nil {
		p.Instancecreatetime = i.InstanceCreateTime.Format(timeutils.ISO8601UTCFormat)
	}
	if i.Endpoint != nil {
		p.Endpoint = &pb.Endpoint{
			Address: i.Endpoint.Address,
			Port:    int32(i.Endpoint.Port),
		}
	}
	return p
}

// snapshotToPb converts a domain DBClusterSnapshot to the RDS API protobuf type.
func snapshotToPb(s *storerds.DBClusterSnapshot, accountId string) *pb.DBClusterSnapshot {
	p := &pb.DBClusterSnapshot{
		Dbclustersnapshotidentifier: s.DBClusterSnapshotIdentifier,
		Dbclusteridentifier:         s.DBClusterIdentifier,
		Engine:                      s.Engine,
		Engineversion:               s.EngineVersion,
		Status:                      s.Status,
		Port:                        int32(s.Port),
		Vpcid:                       s.VpcId,
		Storageencrypted:            proto.Bool(s.StorageEncrypted),
		Kmskeyid:                    s.KmsKeyId,
		Dbclustersnapshotarn:        s.DBSnapshotArn,
	}
	if s.SnapshotCreateTime != nil {
		p.Snapshotcreatetime = s.SnapshotCreateTime.Format(timeutils.ISO8601UTCFormat)
	}
	if s.ClusterCreateTime != nil {
		p.Clustercreatetime = s.ClusterCreateTime.Format(timeutils.ISO8601UTCFormat)
	}
	return p
}

// dbSnapshotToPb converts a domain DBInstanceSnapshot to the RDS API protobuf type.
func dbSnapshotToPb(s *storerds.DBInstanceSnapshot, accountId string) *pb.DBSnapshot {
	p := &pb.DBSnapshot{
		Dbsnapshotidentifier:             s.DBSnapshotIdentifier,
		Dbinstanceidentifier:             s.DBInstanceIdentifier,
		Engine:                           s.Engine,
		Engineversion:                    s.EngineVersion,
		Snapshottype:                     s.SnapshotType,
		Status:                           s.Status,
		Allocatedstorage:                 int32(s.AllocatedStorage),
		Storagetype:                      s.StorageType,
		Port:                             int32(s.Port),
		Availabilityzone:                 s.AvailabilityZone,
		Vpcid:                            s.VpcId,
		Masterusername:                   s.MasterUsername,
		Licensemodel:                     s.LicenseModel,
		Encrypted:                        proto.Bool(s.StorageEncrypted),
		Kmskeyid:                         s.KmsKeyId,
		Dbsnapshotarn:                    s.DBSnapshotArn,
		Iamdatabaseauthenticationenabled: proto.Bool(s.IAMDatabaseAuthEnabled),
		Optiongroupname:                  s.OptionGroupName,
	}
	if s.SnapshotCreateTime != nil {
		p.Snapshotcreatetime = s.SnapshotCreateTime.Format(timeutils.ISO8601UTCFormat)
	}
	if s.InstanceCreateTime != nil {
		p.Instancecreatetime = s.InstanceCreateTime.Format(timeutils.ISO8601UTCFormat)
	}
	if len(s.TagList) > 0 {
		p.Taglist = make([]*pb.Tag, 0, len(s.TagList))
		for _, t := range s.TagList {
			p.Taglist = append(p.Taglist, &pb.Tag{Key: t.Key, Value: t.Value})
		}
	}
	return p
}

// clusterParamGroupToPb converts a domain DBClusterParameterGroup to the
// RDS API protobuf type.
func clusterParamGroupToPb(g *storerds.DBClusterParameterGroup) *pb.DBClusterParameterGroup {
	return &pb.DBClusterParameterGroup{
		Dbclusterparametergroupname: g.DBClusterParameterGroupName,
		Dbparametergroupfamily:      g.DBParameterGroupFamily,
		Description:                 g.Description,
		Dbclusterparametergrouparn:  g.ARN,
	}
}

// paramGroupToPb converts a domain DBParameterGroup to the RDS API protobuf type.
func paramGroupToPb(g *storerds.DBParameterGroup) *pb.DBParameterGroup {
	return &pb.DBParameterGroup{
		Dbparametergroupname:   g.DBParameterGroupName,
		Dbparametergroupfamily: g.DBParameterGroupFamily,
		Description:            g.Description,
		Dbparametergrouparn:    g.ARN,
	}
}

// subnetGroupToPb converts a domain DBSubnetGroup to the RDS API protobuf type.
func subnetGroupToPb(g *storerds.DBSubnetGroup) *pb.DBSubnetGroup {
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

// globalClusterToPb converts a domain GlobalCluster to the RDS API protobuf type.
func globalClusterToPb(c *storerds.GlobalCluster) *pb.GlobalCluster {
	p := &pb.GlobalCluster{
		Globalclusteridentifier: c.GlobalClusterIdentifier,
		Globalclusterresourceid: c.GlobalClusterResourceId,
		Globalclusterarn:        c.GlobalClusterArn,
		Engine:                  c.Engine,
		Engineversion:           c.EngineVersion,
		Status:                  c.Status,
		Storageencrypted:        proto.Bool(c.StorageEncrypted),
		Deletionprotection:      proto.Bool(c.DeletionProtection),
	}
	for _, m := range c.GlobalClusterMembers {
		p.Globalclustermembers = append(p.Globalclustermembers, &pb.GlobalClusterMember{
			Dbclusterarn: m.DBClusterArn,
			Iswriter:     proto.Bool(m.IsWriter),
			Readers:      m.Readers,
		})
	}
	return p
}

// eventSubscriptionToPb converts a domain EventSubscription to the RDS API
// protobuf type.
func eventSubscriptionToPb(s *storerds.EventSubscription) *pb.EventSubscription {
	p := &pb.EventSubscription{
		Custsubscriptionid:   s.CustSubscriptionId,
		Snstopicarn:          s.SnsTopicArn,
		Status:               s.Status,
		Sourcetype:           s.SourceType,
		Sourceidslist:        s.SourceIdsList,
		Eventcategorieslist:  s.EventCategoriesList,
		Enabled:              proto.Bool(s.Enabled),
		Eventsubscriptionarn: s.CustSubscriptionArn,
	}
	return p
}

// NewConnectHandler creates a gRPC-Web connect handler for the RDS admin console.
func NewConnectHandler(stores StoreProvider, engines EngineProvider, accountID string) (string, http.Handler) {
	return rdsconnect.NewRDSServiceHandler(NewAdminHandler(stores, engines, accountID))
}

// CreateDBSnapshot creates a manual snapshot of a DB instance.
func (h *AdminHandler) CreateDBSnapshot(ctx context.Context, req *connect.Request[pb.CreateDBSnapshotMessage]) (*connect.Response[pb.CreateDBSnapshotResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	instanceID := req.Msg.Dbinstanceidentifier
	snapshotID := req.Msg.Dbsnapshotidentifier
	if snapshotID == "" || instanceID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DBSnapshotIdentifier and DBInstanceIdentifier are required"))
	}

	instance, err := store.GetInstance(instanceID)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	now := time.Now()
	arn := fmt.Sprintf("arn:aws:rds:%s:%s:snapshot:%s", region, h.accountId, snapshotID)
	snap := &storerds.DBInstanceSnapshot{
		DBSnapshotIdentifier:   snapshotID,
		DBInstanceIdentifier:   instanceID,
		SnapshotCreateTime:     &now,
		InstanceCreateTime:     instance.InstanceCreateTime,
		Engine:                 instance.Engine,
		EngineVersion:          instance.EngineVersion,
		SnapshotType:           "manual",
		Status:                 "available",
		AvailabilityZone:       instance.AvailabilityZone,
		DBSnapshotArn:          arn,
		IAMDatabaseAuthEnabled: instance.IAMDatabaseAuthenticationEnabled,
		AccountID:              h.accountId,
		Region:                 region,
	}

	if err := store.CreateInstanceSnapshot(snap); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDBSnapshotResult{
		Dbsnapshot: dbSnapshotToPb(snap, h.accountId),
	}), nil
}

// DescribeDBSnapshots lists DB instance snapshots with optional filtering.
func (h *AdminHandler) DescribeDBSnapshots(ctx context.Context, req *connect.Request[pb.DescribeDBSnapshotsMessage]) (*connect.Response[pb.DBSnapshotMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	snapshots, err := store.ListInstanceSnapshots()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	pbSnapshots := make([]*pb.DBSnapshot, 0, len(snapshots))
	for _, s := range snapshots {
		if req.Msg.Dbsnapshotidentifier != "" && s.DBSnapshotIdentifier != req.Msg.Dbsnapshotidentifier {
			continue
		}
		if req.Msg.Dbinstanceidentifier != "" && s.DBInstanceIdentifier != req.Msg.Dbinstanceidentifier {
			continue
		}
		if req.Msg.Snapshottype != "" && s.SnapshotType != req.Msg.Snapshottype {
			continue
		}
		pbSnapshots = append(pbSnapshots, dbSnapshotToPb(s, h.accountId))
	}

	return connect.NewResponse(&pb.DBSnapshotMessage{
		Dbsnapshots: pbSnapshots,
	}), nil
}
