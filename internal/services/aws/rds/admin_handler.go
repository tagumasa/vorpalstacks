// Package rds provides the RDS admin console handler that serves management
// RPCs for both Neptune and MySQL database engines via the gRPC-Web admin
// interface. The handler is a thin gRPC-Web wrapper that delegates all
// business logic to *Core methods on RDSService, keeping this file free of
// store-package imports per AGENTS.md #29.
package rds

import (
	"context"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	pb "vorpalstacks/internal/pb/aws/rds"
	rdsconnect "vorpalstacks/internal/pb/aws/rds/rdsconnect"
)

// AdminHandler implements the RDS Management API gRPC-Web admin console
// handler. Each RPC method resolves the region store, builds a DTO input,
// delegates to the corresponding *Core method on RDSService, and wraps the
// result. Error mapping uses svcerrors.AWSErrorToGRPC exclusively.
type AdminHandler struct {
	rdsconnect.UnimplementedRDSServiceHandler
	service *RDSService
}

var _ rdsconnect.RDSServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(service *RDSService) *AdminHandler {
	return &AdminHandler{service: service}
}

// getStore resolves the region from gRPC-Web request headers and returns
// the common RDS store wrapper. The return type *rdsStores is a local
// wrapper defined in service.go; admin_handler.go does not import the
// store package.
func (h *AdminHandler) getStore(header http.Header) (*rdsStores, error) {
	region := svccommon.GetRegionFromHeader(header)
	return h.service.GetStoreForRegion(region)
}

// NewConnectHandler creates a gRPC-Web connect handler for the RDS admin
// console.
func NewConnectHandler(service *RDSService) (string, http.Handler) {
	return rdsconnect.NewRDSServiceHandler(NewAdminHandler(service))
}

// =====================================================================
// DBCluster RPCs
// =====================================================================

func (h *AdminHandler) DescribeDBClusters(ctx context.Context, req *connect.Request[pb.DescribeDBClustersMessage]) (*connect.Response[pb.DBClusterMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBClustersCore(stores, DescribeDBClustersInput{
		DBClusterIdentifier: req.Msg.Dbclusteridentifier,
		Filters:             req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) CreateDBCluster(ctx context.Context, req *connect.Request[pb.CreateDBClusterMessage]) (*connect.Response[pb.CreateDBClusterResult], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.createDBClusterCore(stores, CreateDBClusterInput{
		DBClusterIdentifier:          req.Msg.Dbclusteridentifier,
		Engine:                       req.Msg.Engine,
		EngineVersion:                req.Msg.Engineversion,
		DatabaseName:                 req.Msg.Databasename,
		MasterUsername:               req.Msg.Masterusername,
		Port:                         req.Msg.GetPort(),
		BackupRetentionPeriod:        req.Msg.GetBackupretentionperiod(),
		AvailabilityZones:            req.Msg.Availabilityzones,
		DBSubnetGroupName:            req.Msg.Dbsubnetgroupname,
		DBClusterParameterGroupName:  req.Msg.Dbclusterparametergroupname,
		StorageEncrypted:             req.Msg.GetStorageencrypted(),
		CopyTagsToSnapshot:           req.Msg.GetCopytagstosnapshot(),
		DeletionProtection:           req.Msg.GetDeletionprotection(),
		EnabledCloudwatchLogsExports: req.Msg.Enablecloudwatchlogsexports,
		IAMDatabaseAuthentication:    req.Msg.GetEnableiamdatabaseauthentication(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DeleteDBCluster(ctx context.Context, req *connect.Request[pb.DeleteDBClusterMessage]) (*connect.Response[pb.DeleteDBClusterResult], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.deleteDBClusterCore(stores, DeleteDBClusterInput{
		DBClusterIdentifier:       req.Msg.Dbclusteridentifier,
		SkipFinalSnapshot:         req.Msg.GetSkipfinalsnapshot(),
		FinalDBSnapshotIdentifier: req.Msg.Finaldbsnapshotidentifier,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeDBClusterSnapshots(ctx context.Context, req *connect.Request[pb.DescribeDBClusterSnapshotsMessage]) (*connect.Response[pb.DBClusterSnapshotMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBClusterSnapshotsCore(stores, DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: req.Msg.Dbclustersnapshotidentifier,
		DBClusterIdentifier:         req.Msg.Dbclusteridentifier,
		SnapshotType:                req.Msg.Snapshottype,
		Filters:                     req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeDBClusterSnapshotAttributes(ctx context.Context, req *connect.Request[pb.DescribeDBClusterSnapshotAttributesMessage]) (*connect.Response[pb.DescribeDBClusterSnapshotAttributesResult], error) {
	return connect.NewResponse(&pb.DescribeDBClusterSnapshotAttributesResult{
		Dbclustersnapshotattributesresult: &pb.DBClusterSnapshotAttributesResult{
			Dbclustersnapshotattributes: []*pb.DBClusterSnapshotAttribute{},
		},
	}), nil
}

func (h *AdminHandler) DescribeDBClusterEndpoints(ctx context.Context, req *connect.Request[pb.DescribeDBClusterEndpointsMessage]) (*connect.Response[pb.DBClusterEndpointMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBClusterEndpointsCore(stores, DescribeDBClusterEndpointsInput{
		DBClusterIdentifier:         req.Msg.Dbclusteridentifier,
		DBClusterEndpointIdentifier: req.Msg.Dbclusterendpointidentifier,
		Filters:                     req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeDBClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeDBClusterParametersMessage]) (*connect.Response[pb.DBClusterParameterGroupDetails], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBClusterParametersCore(stores, DescribeDBClusterParametersInput{
		DBClusterParameterGroupName: req.Msg.Dbclusterparametergroupname,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeEngineDefaultClusterParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultClusterParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultClusterParametersResult], error) {
	family := req.Msg.Dbparametergroupfamily
	pbParams := make([]*pb.Parameter, 0, 4)
	for _, dp := range defaultClusterParamsForFamily(family) {
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
	sortParameters(pbParams)

	return connect.NewResponse(&pb.DescribeEngineDefaultClusterParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: family,
			Parameters:             pbParams,
		},
	}), nil
}

// =====================================================================
// DBInstance RPCs
// =====================================================================

func (h *AdminHandler) DescribeDBInstances(ctx context.Context, req *connect.Request[pb.DescribeDBInstancesMessage]) (*connect.Response[pb.DBInstanceMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBInstancesCore(stores, DescribeDBInstancesInput{
		DBInstanceIdentifier: req.Msg.Dbinstanceidentifier,
		Filters:              req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) CreateDBInstance(ctx context.Context, req *connect.Request[pb.CreateDBInstanceMessage]) (*connect.Response[pb.CreateDBInstanceResult], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.createDBInstanceCore(stores, CreateDBInstanceInput{
		DBInstanceIdentifier:               req.Msg.Dbinstanceidentifier,
		DBClusterIdentifier:                req.Msg.Dbclusteridentifier,
		Engine:                             req.Msg.Engine,
		EngineVersion:                      req.Msg.Engineversion,
		DBInstanceClass:                    req.Msg.Dbinstanceclass,
		AvailabilityZone:                   req.Msg.Availabilityzone,
		PreferredMaintenanceWindow:         req.Msg.Preferredmaintenancewindow,
		PreferredBackupWindow:              req.Msg.Preferredbackupwindow,
		DBParameterGroupName:               req.Msg.Dbparametergroupname,
		DBSubnetGroupName:                  req.Msg.Dbsubnetgroupname,
		PubliclyAccessible:                 req.Msg.GetPubliclyaccessible(),
		AutoMinorVersionUpgrade:            req.Msg.GetAutominorversionupgrade(),
		AllocatedStorage:                   req.Msg.GetAllocatedstorage(),
		MasterUsername:                     req.Msg.Masterusername,
		StorageType:                        req.Msg.Storagetype,
		BackupRetentionPeriod:              req.Msg.GetBackupretentionperiod(),
		LicenseModel:                       req.Msg.Licensemodel,
		StorageEncrypted:                   req.Msg.GetStorageencrypted(),
		KmsKeyId:                           req.Msg.Kmskeyid,
		DeletionProtection:                 req.Msg.GetDeletionprotection(),
		MultiAZ:                            req.Msg.GetMultiaz(),
		Port:                               req.Msg.GetPort(),
		OptionGroupName:                    req.Msg.Optiongroupname,
		Iops:                               req.Msg.GetIops(),
		MaxAllocatedStorage:                req.Msg.GetMaxallocatedstorage(),
		StorageThroughput:                  req.Msg.GetStoragethroughput(),
		MonitoringInterval:                 req.Msg.GetMonitoringinterval(),
		MonitoringRoleArn:                  req.Msg.Monitoringrolearn,
		EnablePerformanceInsights:          req.Msg.GetEnableperformanceinsights(),
		PerformanceInsightsKMSKeyId:        req.Msg.Performanceinsightskmskeyid,
		PerformanceInsightsRetentionPeriod: req.Msg.GetPerformanceinsightsretentionperiod(),
		CACertificateIdentifier:            req.Msg.Cacertificateidentifier,
		CopyTagsToSnapshot:                 req.Msg.GetCopytagstosnapshot(),
		EnabledCloudwatchLogsExports:       req.Msg.Enablecloudwatchlogsexports,
		EnableIAMDatabaseAuthentication:    req.Msg.GetEnableiamdatabaseauthentication(),
		VpcSecurityGroupIds:                req.Msg.Vpcsecuritygroupids,
		DBSecurityGroups:                   req.Msg.Dbsecuritygroups,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DeleteDBInstance(ctx context.Context, req *connect.Request[pb.DeleteDBInstanceMessage]) (*connect.Response[pb.DeleteDBInstanceResult], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.deleteDBInstanceCore(stores, DeleteDBInstanceInput{
		DBInstanceIdentifier:      req.Msg.Dbinstanceidentifier,
		SkipFinalSnapshot:         req.Msg.GetSkipfinalsnapshot(),
		FinalDBSnapshotIdentifier: req.Msg.Finaldbsnapshotidentifier,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeValidDBInstanceModifications(ctx context.Context, req *connect.Request[pb.DescribeValidDBInstanceModificationsMessage]) (*connect.Response[pb.DescribeValidDBInstanceModificationsResult], error) {
	return connect.NewResponse(&pb.DescribeValidDBInstanceModificationsResult{
		Validdbinstancemodificationsmessage: &pb.ValidDBInstanceModificationsMessage{
			Storage: []*pb.ValidStorageOptions{},
		},
	}), nil
}

func (h *AdminHandler) CreateDBSnapshot(ctx context.Context, req *connect.Request[pb.CreateDBSnapshotMessage]) (*connect.Response[pb.CreateDBSnapshotResult], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.createDBSnapshotCore(stores, CreateDBSnapshotInput{
		DBInstanceIdentifier: req.Msg.Dbinstanceidentifier,
		DBSnapshotIdentifier: req.Msg.Dbsnapshotidentifier,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeDBSnapshots(ctx context.Context, req *connect.Request[pb.DescribeDBSnapshotsMessage]) (*connect.Response[pb.DBSnapshotMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBSnapshotsCore(stores, DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: req.Msg.Dbsnapshotidentifier,
		DBInstanceIdentifier: req.Msg.Dbinstanceidentifier,
		SnapshotType:         req.Msg.Snapshottype,
		Filters:              req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

// =====================================================================
// Parameter Group RPCs
// =====================================================================

func (h *AdminHandler) DescribeDBClusterParameterGroups(ctx context.Context, req *connect.Request[pb.DescribeDBClusterParameterGroupsMessage]) (*connect.Response[pb.DBClusterParameterGroupsMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBClusterParameterGroupsCore(stores, DescribeDBClusterParameterGroupsInput{
		DBClusterParameterGroupName: req.Msg.Dbclusterparametergroupname,
		Filters:                     req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeDBParameterGroups(ctx context.Context, req *connect.Request[pb.DescribeDBParameterGroupsMessage]) (*connect.Response[pb.DBParameterGroupsMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBParameterGroupsCore(stores, DescribeDBParameterGroupsInput{
		DBParameterGroupName: req.Msg.Dbparametergroupname,
		Filters:              req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeDBParameters(ctx context.Context, req *connect.Request[pb.DescribeDBParametersMessage]) (*connect.Response[pb.DBParameterGroupDetails], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBParametersCore(stores, DescribeDBParametersInput{
		DBParameterGroupName: req.Msg.Dbparametergroupname,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeEngineDefaultParameters(ctx context.Context, req *connect.Request[pb.DescribeEngineDefaultParametersMessage]) (*connect.Response[pb.DescribeEngineDefaultParametersResult], error) {
	family := req.Msg.Dbparametergroupfamily
	pbParams := make([]*pb.Parameter, 0, 4)
	for _, dp := range defaultInstanceParamsForFamily(family) {
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
	sortParameters(pbParams)

	return connect.NewResponse(&pb.DescribeEngineDefaultParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: family,
			Parameters:             pbParams,
		},
	}), nil
}

// =====================================================================
// Misc RPCs (SubnetGroups, GlobalClusters, Events, Tags)
// =====================================================================

func (h *AdminHandler) DescribeDBSubnetGroups(ctx context.Context, req *connect.Request[pb.DescribeDBSubnetGroupsMessage]) (*connect.Response[pb.DBSubnetGroupMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeDBSubnetGroupsCore(stores, DescribeDBSubnetGroupsInput{
		DBSubnetGroupName: req.Msg.Dbsubnetgroupname,
		Filters:           req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeGlobalClusters(ctx context.Context, req *connect.Request[pb.DescribeGlobalClustersMessage]) (*connect.Response[pb.GlobalClustersMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeGlobalClustersCore(stores, DescribeGlobalClustersInput{
		GlobalClusterIdentifier: req.Msg.Globalclusteridentifier,
		Filters:                 req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeEventSubscriptions(ctx context.Context, req *connect.Request[pb.DescribeEventSubscriptionsMessage]) (*connect.Response[pb.EventSubscriptionsMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeEventSubscriptionsCore(stores, DescribeEventSubscriptionsInput{
		SubscriptionName: req.Msg.Subscriptionname,
		Filters:          req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) DescribeEvents(ctx context.Context, req *connect.Request[pb.DescribeEventsMessage]) (*connect.Response[pb.EventsMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.describeEventsCore(stores, DescribeEventsInput{
		SourceType:       req.Msg.Sourcetype,
		SourceIdentifier: req.Msg.Sourceidentifier,
		StartTime:        req.Msg.Starttime,
		EndTime:          req.Msg.Endtime,
		Duration:         req.Msg.GetDuration(),
		EventCategories:  req.Msg.Eventcategories,
		Marker:           req.Msg.Marker,
		MaxRecords:       req.Msg.GetMaxrecords(),
		Filters:          req.Msg.Filters,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceMessage]) (*connect.Response[pb.TagListMessage], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.listTagsForResourceCore(stores, ListTagsForResourceInput{
		ResourceName: req.Msg.Resourcename,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) AddTagsToResource(ctx context.Context, req *connect.Request[pb.AddTagsToResourceMessage]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.addTagsToResourceCore(stores, AddTagsToResourceInput{
		ResourceName: req.Msg.Resourcename,
		Tags:         req.Msg.Tags,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

func (h *AdminHandler) RemoveTagsFromResource(ctx context.Context, req *connect.Request[pb.RemoveTagsFromResourceMessage]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	result, err := h.service.removeTagsFromResourceCore(stores, RemoveTagsFromResourceInput{
		ResourceName: req.Msg.Resourcename,
		TagKeys:      req.Msg.Tagkeys,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}
	return connect.NewResponse(result), nil
}

// =====================================================================
// Static RPCs (no store access required)
// =====================================================================

func (h *AdminHandler) DescribeDBEngineVersions(ctx context.Context, req *connect.Request[pb.DescribeDBEngineVersionsMessage]) (*connect.Response[pb.DBEngineVersionMessage], error) {
	versions := allEngineVersions()

	if engine := req.Msg.Engine; engine != "" {
		filtered := make([]*pb.DBEngineVersion, 0)
		for _, v := range versions {
			if v.Engine == engine {
				filtered = append(filtered, v)
			}
		}
		versions = filtered
	}

	if ev := req.Msg.Engineversion; ev != "" {
		filtered := make([]*pb.DBEngineVersion, 0)
		for _, v := range versions {
			if v.Engineversion == ev {
				filtered = append(filtered, v)
			}
		}
		versions = filtered
	}

	return connect.NewResponse(&pb.DBEngineVersionMessage{
		Dbengineversions: versions,
	}), nil
}

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

func (h *AdminHandler) DescribePendingMaintenanceActions(ctx context.Context, req *connect.Request[pb.DescribePendingMaintenanceActionsMessage]) (*connect.Response[pb.PendingMaintenanceActionsMessage], error) {
	return connect.NewResponse(&pb.PendingMaintenanceActionsMessage{
		Pendingmaintenanceactions: []*pb.ResourcePendingMaintenanceActions{},
	}), nil
}

func (h *AdminHandler) DescribeOrderableDBInstanceOptions(ctx context.Context, req *connect.Request[pb.DescribeOrderableDBInstanceOptionsMessage]) (*connect.Response[pb.OrderableDBInstanceOptionsMessage], error) {
	options := []*pb.OrderableDBInstanceOption{
		{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.t3.medium", Licensemodel: "bring-your-own-license", Vpc: proto.Bool(true)},
		{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.r5.large", Licensemodel: "bring-your-own-license", Vpc: proto.Bool(true)},
		{Engine: "neptune", Engineversion: "1.4.0.1", Dbinstanceclass: "db.r5.xlarge", Licensemodel: "bring-your-own-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.t3.micro", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.t3.small", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.t3.medium", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.r5.large", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
		{Engine: "mysql", Engineversion: "8.0.40", Dbinstanceclass: "db.r5.xlarge", Licensemodel: "general-public-license", Vpc: proto.Bool(true)},
	}

	if engine := req.Msg.Engine; engine != "" {
		filtered := make([]*pb.OrderableDBInstanceOption, 0)
		for _, o := range options {
			if o.Engine == engine {
				filtered = append(filtered, o)
			}
		}
		options = filtered
	}

	return connect.NewResponse(&pb.OrderableDBInstanceOptionsMessage{
		Orderabledbinstanceoptions: options,
	}), nil
}

// =====================================================================
// Parameter family defaults (shared by admin_handler.go static RPCs and
// Core methods in dbcluster_core.go / dbparams_core.go)
// =====================================================================

func defaultClusterParamsForFamily(family string) []struct{ name, value, desc, source, apply, dtype, modifiable string } {
	switch {
	case strings.HasPrefix(family, "neptune"):
		return []struct{ name, value, desc, source, apply, dtype, modifiable string }{
			{"neptune_query_timeout", "120000", "Query execution timeout in milliseconds", "system", "dynamic", "integer", "true"},
			{"neptune_enable_audit_log", "0", "Enable audit logging", "system", "static", "boolean", "true"},
			{"neptune_streams", "0", "Enable Neptune streams", "system", "dynamic", "boolean", "true"},
			{"neptune_replica_primary_endpoint", "", "Override endpoint for replica primary", "system", "dynamic", "string", "true"},
		}
	case strings.HasPrefix(family, "mysql5.7"):
		return []struct{ name, value, desc, source, apply, dtype, modifiable string }{
			{"auto_increment_increment", "1", "Increments between auto-generated IDs", "system", "dynamic", "integer", "true"},
			{"auto_increment_offset", "1", "Offset for auto-generated IDs", "system", "dynamic", "integer", "true"},
			{"character_set_server", "utf8mb4", "Default server character set", "system", "dynamic", "string", "true"},
			{"collation_server", "utf8mb4_general_ci", "Server default collation", "system", "dynamic", "string", "true"},
			{"max_connections", "150", "Maximum simultaneous connections", "system", "dynamic", "integer", "true"},
			{"innodb_buffer_pool_size", "{DBInstanceClassMemory*3/4}", "Size of the InnoDB buffer pool", "system", "dynamic", "integer", "true"},
			{"innodb_log_file_size", "134217728", "Size of each InnoDB log file in bytes", "system", "dynamic", "integer", "true"},
			{"long_query_time", "10", "Threshold in seconds for slow_query_log", "system", "dynamic", "float", "true"},
			{"slow_query_log", "0", "Enable slow query log", "system", "dynamic", "boolean", "true"},
			{"transaction_isolation", "REPEATABLE-READ", "Default transaction isolation level", "system", "dynamic", "string", "true"},
		}
	case strings.HasPrefix(family, "mysql8.0"):
		return []struct{ name, value, desc, source, apply, dtype, modifiable string }{
			{"auto_increment_increment", "1", "Increments between auto-generated IDs", "system", "dynamic", "integer", "true"},
			{"auto_increment_offset", "1", "Offset for auto-generated IDs", "system", "dynamic", "integer", "true"},
			{"character_set_server", "utf8mb4", "Default server character set", "system", "dynamic", "string", "true"},
			{"collation_server", "utf8mb4_0900_ai_ci", "Server default collation", "system", "dynamic", "string", "true"},
			{"max_connections", "150", "Maximum simultaneous connections", "system", "dynamic", "integer", "true"},
			{"innodb_buffer_pool_size", "{DBInstanceClassMemory*3/4}", "Size of the InnoDB buffer pool", "system", "dynamic", "integer", "true"},
			{"innodb_log_file_size", "134217728", "Size of each InnoDB log file in bytes", "system", "dynamic", "integer", "true"},
			{"long_query_time", "10", "Threshold in seconds for slow_query_log", "system", "dynamic", "float", "true"},
			{"slow_query_log", "0", "Enable slow query log", "system", "dynamic", "boolean", "true"},
			{"transaction_isolation", "REPEATABLE-READ", "Default transaction isolation level", "system", "dynamic", "string", "true"},
			{"binlog_expire_logs_seconds", "0", "Binlog expiration in seconds", "system", "dynamic", "integer", "true"},
		}
	case strings.HasPrefix(family, "mysql8.4"):
		return []struct{ name, value, desc, source, apply, dtype, modifiable string }{
			{"auto_increment_increment", "1", "Increments between auto-generated IDs", "system", "dynamic", "integer", "true"},
			{"character_set_server", "utf8mb4", "Default server character set", "system", "dynamic", "string", "true"},
			{"collation_server", "utf8mb4_0900_ai_ci", "Server default collation", "system", "dynamic", "string", "true"},
			{"max_connections", "150", "Maximum simultaneous connections", "system", "dynamic", "integer", "true"},
			{"innodb_buffer_pool_size", "{DBInstanceClassMemory*3/4}", "Size of the InnoDB buffer pool", "system", "dynamic", "integer", "true"},
			{"long_query_time", "10", "Threshold in seconds for slow_query_log", "system", "dynamic", "float", "true"},
			{"slow_query_log", "0", "Enable slow query log", "system", "dynamic", "boolean", "true"},
			{"transaction_isolation", "REPEATABLE-READ", "Default transaction isolation level", "system", "dynamic", "string", "true"},
			{"binlog_expire_logs_seconds", "0", "Binlog expiration in seconds", "system", "dynamic", "integer", "true"},
			{"innodb_redo_log_capacity", "8589934592", "InnoDB redo log capacity in bytes", "system", "dynamic", "integer", "true"},
		}
	case strings.HasPrefix(family, "postgres"), strings.HasPrefix(family, "aurora-postgresql"):
		return []struct{ name, value, desc, source, apply, dtype, modifiable string }{
			{"max_connections", "100", "Maximum simultaneous connections", "system", "dynamic", "integer", "true"},
			{"shared_buffers", "128", "Shared buffer pool size in 8KB pages", "system", "static", "integer", "false"},
			{"work_mem", "4096", "Memory in KB used by internal sort operations", "system", "dynamic", "integer", "true"},
			{"maintenance_work_mem", "65536", "Memory in KB used by maintenance operations", "system", "dynamic", "integer", "true"},
			{"log_min_duration_statement", "-1", "Log statements slower than this many ms", "system", "dynamic", "integer", "true"},
		}
	default:
		return nil
	}
}

func defaultInstanceParamsForFamily(family string) []struct{ name, value, desc, source, apply, dtype, modifiable string } {
	switch {
	case strings.HasPrefix(family, "neptune"):
		return []struct{ name, value, desc, source, apply, dtype, modifiable string }{
			{"neptune_query_timeout", "120000", "Query execution timeout", "system", "dynamic", "integer", "true"},
		}
	case strings.HasPrefix(family, "mysql5.7"), strings.HasPrefix(family, "mysql8.0"), strings.HasPrefix(family, "mysql8.4"):
		return defaultClusterParamsForFamily(family)
	case strings.HasPrefix(family, "postgres"), strings.HasPrefix(family, "aurora-postgresql"):
		return defaultClusterParamsForFamily(family)
	default:
		return nil
	}
}
