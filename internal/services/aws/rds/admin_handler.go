// Package rds provides the RDS admin console handler that serves management
// RPCs for both Neptune and MySQL database engines via the gRPC-Web admin
// interface. It delegates data access to the common RDS store layer.
package rds

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
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

	// snapOp, when non-nil, allows CreateDBSnapshot and DeleteDBInstance
	// to capture row-level data for snapshots. When nil, snapshots store
	// metadata only (legacy behaviour).
	snapOp SnapshotOperator
}

var _ rdsconnect.RDSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new RDS admin console handler backed by the given
// store provider and engine provider.
func NewAdminHandler(stores StoreProvider, engines EngineProvider, accountId string) *AdminHandler {
	return &AdminHandler{stores: stores, engines: engines, accountId: accountId}
}

// SetSnapshotOperator wires the snapshot data-copier. Called by
// optional.go after the vmysql service has been initialised. Without
// this call, snapshots store metadata only.
func (h *AdminHandler) SetSnapshotOperator(op SnapshotOperator) {
	h.snapOp = op
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
		if !applyRDSFilters(req.Msg.Filters, clusterFilterGetter(c)) {
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
		if !applyRDSFilters(req.Msg.Filters, instanceFilterGetter(i)) {
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
	if err := ValidateDBInstanceIdentifier(id); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if err := ValidateDBInstanceClass(req.Msg.Dbinstanceclass); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	engine := req.Msg.Engine
	if engine == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Engine is required"))
	}
	if err := ValidateEngine(engine); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	engineVersion := req.Msg.Engineversion
	if engineVersion == "" {
		engineVersion = DefaultEngineVersion(engine)
	}
	if err := ValidateEngineVersion(engine, engineVersion); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	// Reject engines without a registered EngineProvider up front. AWS
	// accepts a long list of engine values but this platform only wires
	// concrete backing engines for a subset; the previous flow persisted
	// the row with status="failed" and leaked the storage entry, which
	// silently degraded the AWS UX (a real AWS CreateDBInstance with an
	// unsupported value returns InvalidParameterValue immediately).
	if _, err := h.engines(engine); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("engine %q is not supported on this platform: %v", engine, err))
	}

	// Validate referenced resource existence before persisting. AWS
	// returns DBSubnetGroupNotFoundFault / DBParameterGroupNotFoundFault
	// when the named groups do not exist at create time.
	if name := req.Msg.Dbsubnetgroupname; name != "" {
		if _, err := store.GetSubnetGroup(name); err != nil {
			return nil, connect.NewError(connect.CodeNotFound,
				fmt.Errorf("DB Subnet Group %q not found", name))
		}
	}
	if name := req.Msg.Dbparametergroupname; name != "" {
		if _, err := store.GetParameterGroup(name); err != nil {
			return nil, connect.NewError(connect.CodeNotFound,
				fmt.Errorf("DB Parameter Group %q not found", name))
		}
	}

	now := time.Now()
	instance := &storerds.DBInstance{
		DBInstanceIdentifier:       id,
		DBClusterIdentifier:        req.Msg.Dbclusteridentifier,
		Engine:                     engine,
		EngineVersion:              engineVersion,
		DBInstanceClass:            req.Msg.Dbinstanceclass,
		DBInstanceStatus:           "creating",
		AvailabilityZone:           req.Msg.Availabilityzone,
		PreferredMaintenanceWindow: req.Msg.Preferredmaintenancewindow,
		PreferredBackupWindow:      req.Msg.Preferredbackupwindow,
		DBParameterGroupName:       req.Msg.Dbparametergroupname,
		DBSubnetGroupName:          req.Msg.Dbsubnetgroupname,
		PubliclyAccessible:         req.Msg.GetPubliclyaccessible(),
		AutoMinorVersionUpgrade:    req.Msg.GetAutominorversionupgrade(),
		InstanceCreateTime:         &now,
		AccountID:                  h.accountId,
		Region:                     region,
		DBInstanceArn:              arnutil.NewARNBuilder(h.accountId, region).RDS().DBInstance(id),
		// DbiResourceId is the AWS Region-unique immutable identifier for
		// the DB instance. AWS allocates 'db-' + 26 base32 characters; we
		// generate a 26-character hex token which is functionally
		// equivalent for our local platform and is surfaced in
		// DescribeDBInstances, CloudTrail-style logging, and snapshot
		// lineage. Previously this field was plumbed end-to-end but never
		// assigned, leaving DescribeDBInstances to return an empty value.
		DbiResourceId: generateDbiResourceId(),

		// AWS-standard DBInstance parameters (RDS-5/RDS-20). These were
		// previously dropped even though the UI sends allocatedstorage and
		// masterusername.
		AllocatedStorage:                   req.Msg.Allocatedstorage,
		MasterUsername:                     req.Msg.Masterusername,
		StorageType:                        req.Msg.Storagetype,
		BackupRetentionPeriod:              req.Msg.Backupretentionperiod,
		LicenseModel:                       req.Msg.Licensemodel,
		StorageEncrypted:                   req.Msg.GetStorageencrypted(),
		KmsKeyId:                           req.Msg.Kmskeyid,
		DeletionProtection:                 req.Msg.GetDeletionprotection(),
		MultiAZ:                            req.Msg.GetMultiaz(),
		Port:                               req.Msg.Port,
		OptionGroupName:                    req.Msg.Optiongroupname,
		Iops:                               req.Msg.Iops,
		MaxAllocatedStorage:                req.Msg.Maxallocatedstorage,
		StorageThroughput:                  req.Msg.Storagethroughput,
		MonitoringInterval:                 req.Msg.Monitoringinterval,
		EnhancedMonitoringResourceArn:      req.Msg.Monitoringrolearn,
		PerformanceInsightsEnabled:         req.Msg.GetEnableperformanceinsights(),
		PerformanceInsightsKMSKeyId:        req.Msg.Performanceinsightskmskeyid,
		PerformanceInsightsRetentionPeriod: req.Msg.Performanceinsightsretentionperiod,
		CACertificateIdentifier:            req.Msg.Cacertificateidentifier,
		CopyTagsToSnapshot:                 req.Msg.GetCopytagstosnapshot(),
		EnabledCloudwatchLogsExports:       req.Msg.Enablecloudwatchlogsexports,
		IAMDatabaseAuthenticationEnabled:   req.Msg.GetEnableiamdatabaseauthentication(),
		VpcSecurityGroupIds:                req.Msg.Vpcsecuritygroupids,
		DBSecurityGroups:                   req.Msg.Dbsecuritygroups,
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Try to start the engine. If engine.Open succeeds but persisting the
	// resulting endpoint fails, we must roll back the engine so the
	// allocated port is released and the system is left in a consistent
	// state. The instance row stays in the store with status="failed" so
	// the operator can see what happened and explicitly delete it.
	engineStarted := false
	if eng, engErr := h.engines(engine); engErr == nil {
		port, openErr := eng.Open(region, id)
		if openErr != nil {
			logs.Warn("rds-admin: failed to start engine for instance",
				logs.String("instance", id), logs.Err(openErr))
		} else {
			instance.Endpoint = &storerds.Endpoint{
				Address: fmt.Sprintf("%s.%s.%s.rds.amazonaws.com", id, h.accountId, region),
				Port:    port,
			}
			if err := store.UpdateInstance(instance); err != nil {
				logs.Warn("rds-admin: rolling back engine after UpdateInstance failure",
					logs.String("instance", id), logs.Err(err))
				if closeErr := eng.Close(id); closeErr != nil {
					logs.Warn("rds-admin: engine rollback Close also failed",
						logs.String("instance", id), logs.Err(closeErr))
				}
				instance.Endpoint = nil
				instance.DBInstanceStatus = "failed"
				// Best-effort: persist the failed status so the operator
				// sees the failure rather than the original 'creating'.
				if persistErr := store.UpdateInstance(instance); persistErr != nil {
					logs.Warn("rds-admin: failed to persist failed status",
						logs.String("instance", id), logs.Err(persistErr))
				}
				return nil, svcerrors.StoreErrorToGRPC(err)
			}
			engineStarted = true
		}
	}

	if engineStarted {
		instance.DBInstanceStatus = "available"
	} else {
		instance.DBInstanceStatus = "failed"
	}
	if err := store.UpdateInstance(instance); err != nil {
		// Final status update failed; if the engine is still open, release
		// the port so we don't leak it while leaving the row inconsistent.
		if engineStarted {
			if eng, engErr := h.engines(engine); engErr == nil {
				if closeErr := eng.Close(id); closeErr != nil {
					logs.Warn("rds-admin: cleanup engine.Close after final UpdateInstance failure",
						logs.String("instance", id), logs.Err(closeErr))
				}
			}
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDBInstanceResult{
		Dbinstance: instanceToPb(instance, h.accountId),
	}), nil
}

// DeleteDBInstance deletes a DB instance and stops its engine if running.
// If SkipFinalSnapshot is false (the AWS default) the caller must supply
// FinalDBSnapshotIdentifier; otherwise the request is rejected.
func (h *AdminHandler) DeleteDBInstance(ctx context.Context, req *connect.Request[pb.DeleteDBInstanceMessage]) (*connect.Response[pb.DeleteDBInstanceResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	id := req.Msg.Dbinstanceidentifier
	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("DBInstanceIdentifier is required"))
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// M6: DeletionProtection must be checked before any destructive
	// action, matching DeleteDBCluster (line ~559) and AWS behaviour.
	if instance.DeletionProtection {
		return nil, connect.NewError(connect.CodeFailedPrecondition,
			fmt.Errorf("cannot delete instance when DeletionProtection is enabled"))
	}

	// AWS default: SkipFinalSnapshot=false means a final snapshot is
	// required. Reject the request if the caller did not provide a
	// FinalDBSnapshotIdentifier in that case, mirroring AWS behaviour.
	if !req.Msg.GetSkipfinalsnapshot() {
		if req.Msg.Finaldbsnapshotidentifier == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument,
				fmt.Errorf("FinalDBSnapshotIdentifier is required when SkipFinalSnapshot is false"))
		}
		if err := ValidateDBSnapshotIdentifier(req.Msg.Finaldbsnapshotidentifier); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		finalSnap := &storerds.DBInstanceSnapshot{
			DBSnapshotIdentifier:   req.Msg.Finaldbsnapshotidentifier,
			DBInstanceIdentifier:   id,
			SnapshotCreateTime:     nil,
			InstanceCreateTime:     instance.InstanceCreateTime,
			Engine:                 instance.Engine,
			EngineVersion:          instance.EngineVersion,
			SnapshotType:           "manual",
			Status:                 "available",
			AvailabilityZone:       instance.AvailabilityZone,
			DBSnapshotArn:          arnutil.NewARNBuilder(h.accountId, region).RDS().Snapshot(req.Msg.Finaldbsnapshotidentifier),
			IAMDatabaseAuthEnabled: instance.IAMDatabaseAuthenticationEnabled,
			AccountID:              h.accountId,
			Region:                 region,
			AllocatedStorage:       instance.AllocatedStorage,
			MasterUsername:         instance.MasterUsername,
			StorageType:            instance.StorageType,
			LicenseModel:           instance.LicenseModel,
			StorageEncrypted:       instance.StorageEncrypted,
			KmsKeyId:               instance.KmsKeyId,
			OptionGroupName:        instance.OptionGroupName,
			VpcId:                  instance.VpcId,
		}
		if instance.Port > 0 {
			finalSnap.Port = instance.Port
		} else if instance.Endpoint != nil {
			finalSnap.Port = int32(instance.Endpoint.Port)
		}
		now := time.Now()
		finalSnap.SnapshotCreateTime = &now
		if err := store.CreateInstanceSnapshot(finalSnap); err != nil {
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
		// Capture row-level data before the engine is torn down.
		// SnapshotData copies every database / table / row / index
		// into a 'snap_<snapshotID>' key prefix in the same Pebble
		// bucket.
		if h.snapOp != nil && instance.Engine == "mysql" {
			if err := h.snapOp.SnapshotData(id, req.Msg.Finaldbsnapshotidentifier); err != nil {
				_ = store.DeleteInstanceSnapshot(req.Msg.Finaldbsnapshotidentifier)
				return nil, connect.NewError(connect.CodeInternal,
					fmt.Errorf("final snapshot data capture failed for instance %q: %v", id, err))
			}
		}
	}

	if eng, engErr := h.engines(instance.Engine); engErr == nil {
		if closeErr := eng.Close(id); closeErr != nil {
			return nil, svcerrors.StoreErrorToGRPC(fmt.Errorf("failed to stop engine for instance %q: %w", id, closeErr))
		}
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
	if err := ValidateDBClusterIdentifier(id); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	// Per AWS RDS API spec, Engine is required for CreateDBCluster. Valid
	// values are aurora-mysql, aurora-postgresql, mysql, postgres, neptune.
	// Defaulting silently masked missing-engine requests as aurora-mysql
	// clusters that could never actually start (no engine provider registered
	// for aurora-mysql on this platform).
	engine := req.Msg.Engine
	if engine == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("Engine is required"))
	}
	if err := ValidateEngine(engine); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	engineVersion := req.Msg.Engineversion
	if engineVersion == "" {
		engineVersion = DefaultEngineVersion(engine)
	}
	if err := ValidateEngineVersion(engine, engineVersion); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if req.Msg.Databasename != "" {
		if err := ValidateDatabaseName(req.Msg.Databasename); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	}
	// Reject engines without a registered EngineProvider up front. See
	// CreateDBInstance for the rationale.
	if _, err := h.engines(engine); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("engine %q is not supported on this platform: %v", engine, err))
	}

	// Validate referenced resource existence before persisting. AWS
	// returns DBSubnetGroupNotFoundFault /
	// DBClusterParameterGroupNotFoundFault when the named groups do not
	// exist at create time.
	if name := req.Msg.Dbsubnetgroupname; name != "" {
		if _, err := store.GetSubnetGroup(name); err != nil {
			return nil, connect.NewError(connect.CodeNotFound,
				fmt.Errorf("DB Subnet Group %q not found", name))
		}
	}
	if name := req.Msg.Dbclusterparametergroupname; name != "" {
		if _, err := store.GetClusterParameterGroup(name); err != nil {
			return nil, connect.NewError(connect.CodeNotFound,
				fmt.Errorf("DB Cluster Parameter Group %q not found", name))
		}
	}

	now := time.Now()
	cluster := &storerds.DBCluster{
		DBClusterIdentifier:              id,
		Engine:                           engine,
		EngineVersion:                    engineVersion,
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

	// Allocate an engine port for the cluster so clients have a usable
	// Endpoint. This mirrors Neptune's CreateDBCluster behaviour. If the
	// engine.Open fails the cluster stays in 'creating' status (the
	// admin_handler does not have a background reconciler to drive async
	// transitions); if engine.Open succeeds but the subsequent
	// UpdateCluster fails we roll the port back so it is not leaked.
	engineStarted := false
	if eng, engErr := h.engines(engine); engErr == nil {
		port, openErr := eng.Open(region, id)
		if openErr != nil {
			logs.Warn("rds-admin: failed to start engine for cluster",
				logs.String("cluster", id), logs.Err(openErr))
		} else {
			cluster.Endpoint = &storerds.Endpoint{
				Address: fmt.Sprintf("%s.%s.%s.rds.amazonaws.com", id, h.accountId, region),
				Port:    port,
			}
			if err := store.UpdateCluster(cluster); err != nil {
				logs.Warn("rds-admin: rolling back cluster engine after UpdateCluster failure",
					logs.String("cluster", id), logs.Err(err))
				if closeErr := eng.Close(id); closeErr != nil {
					logs.Warn("rds-admin: cluster engine rollback Close also failed",
						logs.String("cluster", id), logs.Err(closeErr))
				}
				cluster.Endpoint = nil
				cluster.Status = "failed"
				if persistErr := store.UpdateCluster(cluster); persistErr != nil {
					logs.Warn("rds-admin: failed to persist cluster failed status",
						logs.String("cluster", id), logs.Err(persistErr))
				}
				return nil, svcerrors.StoreErrorToGRPC(err)
			}
			engineStarted = true
		}
	}

	if engineStarted {
		cluster.Status = "available"
	} else {
		cluster.Status = "failed"
	}
	if err := store.UpdateCluster(cluster); err != nil {
		if engineStarted {
			if eng, engErr := h.engines(engine); engErr == nil {
				if closeErr := eng.Close(id); closeErr != nil {
					logs.Warn("rds-admin: cleanup engine.Close after final cluster UpdateCluster failure",
						logs.String("cluster", id), logs.Err(closeErr))
				}
			}
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDBClusterResult{
		Dbcluster: clusterToPb(cluster, h.accountId),
	}), nil
}

// DeleteDBCluster deletes a DB cluster.
// If SkipFinalSnapshot is false (the AWS default) the caller must supply
// FinalDBSnapshotIdentifier; otherwise the request is rejected.
func (h *AdminHandler) DeleteDBCluster(ctx context.Context, req *connect.Request[pb.DeleteDBClusterMessage]) (*connect.Response[pb.DeleteDBClusterResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	region := svccommon.GetRegionFromHeader(req.Header())
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

	// AWS default: SkipFinalSnapshot=false means a final snapshot is
	// required. Reject the request if the caller did not provide a
	// FinalDBSnapshotIdentifier in that case.
	if !req.Msg.GetSkipfinalsnapshot() {
		if req.Msg.Finaldbsnapshotidentifier == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument,
				fmt.Errorf("FinalDBSnapshotIdentifier is required when SkipFinalSnapshot is false"))
		}
		now := time.Now()
		finalSnap := &storerds.DBClusterSnapshot{
			DBClusterSnapshotIdentifier:      req.Msg.Finaldbsnapshotidentifier,
			DBClusterIdentifier:              id,
			SnapshotCreateTime:               &now,
			Engine:                           cluster.Engine,
			EngineVersion:                    cluster.EngineVersion,
			SnapshotType:                     "manual",
			Status:                           "available",
			Port:                             cluster.Port,
			VpcId:                            "",
			ClusterCreateTime:                cluster.ClusterCreateTime,
			StorageEncrypted:                 cluster.StorageEncrypted,
			KmsKeyId:                         cluster.KmsKeyId,
			DBSnapshotArn:                    arnutil.NewARNBuilder(h.accountId, region).RDS().ClusterSnapshot(req.Msg.Finaldbsnapshotidentifier),
			AccountID:                        h.accountId,
			Region:                           region,
			MasterUsername:                   cluster.MasterUsername,
			AllocatedStorage:                 0,
			StorageType:                      cluster.StorageType,
			LicenseModel:                     "",
			IAMDatabaseAuthenticationEnabled: cluster.IAMDatabaseAuthenticationEnabled,
		}
		if err := store.CreateSnapshot(finalSnap); err != nil {
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
	}

	cluster.Status = "deleting"
	if err := store.UpdateCluster(cluster); err != nil {
		// Persisting the "deleting" status is best-effort — if Pebble is
		// unavailable the deletion itself will also fail on the next call,
		// so log the warning and surface the underlying error to the caller
		// rather than silently proceeding with an inconsistent view.
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Release the engine backing the cluster, mirroring DeleteDBInstance.
	// CreateDBCluster opened the engine via eng.Open(region, id); without
	// a matching Close the allocated port and goroutine state leak.
	if eng, engErr := h.engines(cluster.Engine); engErr == nil {
		if closeErr := eng.Close(id); closeErr != nil {
			logs.Warn("rds-admin: engine.Close failed during DeleteDBCluster",
				logs.String("cluster", id), logs.Err(closeErr))
		}
	}

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
		if req.Msg.Snapshottype != "" && s.SnapshotType != req.Msg.Snapshottype {
			continue
		}
		if !applyRDSFilters(req.Msg.Filters, clusterSnapshotFilterGetter(s)) {
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
		if !applyRDSFilters(req.Msg.Filters, clusterParamGroupFilterGetter(g)) {
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
		if !applyRDSFilters(req.Msg.Filters, paramGroupFilterGetter(g)) {
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
		if !applyRDSFilters(req.Msg.Filters, subnetGroupFilterGetter(g)) {
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
		if !applyRDSFilters(req.Msg.Filters, globalClusterFilterGetter(c)) {
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
		if !applyRDSFilters(req.Msg.Filters, eventSubscriptionFilterGetter(s)) {
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
// Neptune and MySQL engines. The full version list lives in validators.go
// (supportedMysqlVersions, supportedNeptuneVersions) so the same data drives
// DescribeDBEngineVersions and ValidateEngineVersion.
func (h *AdminHandler) DescribeDBEngineVersions(ctx context.Context, req *connect.Request[pb.DescribeDBEngineVersionsMessage]) (*connect.Response[pb.DBEngineVersionMessage], error) {
	versions := allEngineVersions()

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

	// Filter by engine version if requested
	if req.Msg.Engineversion != "" {
		filtered := make([]*pb.DBEngineVersion, 0)
		for _, v := range versions {
			if v.Engineversion == req.Msg.Engineversion {
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

// DescribeEvents returns RDS events from the per-region store. All AWS-spec
// filter parameters (SourceType, SourceIdentifier, StartTime, EndTime,
// Duration, EventCategories, Marker, MaxRecords) are honoured; an unknown
// SourceType value is rejected with InvalidArgument rather than silently
// collapsing to a default.
func (h *AdminHandler) DescribeEvents(ctx context.Context, req *connect.Request[pb.DescribeEventsMessage]) (*connect.Response[pb.EventsMessage], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// AWS RDS DescribeEvents does not officially support the Filters
	// parameter ("This parameter isn't currently supported."). However,
	// the proto3 enum mapping for SourceType assigns
	// BLUE_GREEN_DEPLOYMENT the zero value, making it indistinguishable
	// from "unset" (see sourceTypeToString comment above). To work
	// around this without a proto-breaking change, the admin console
	// client may pass a Filter with name "source-type" and a string
	// value to explicitly select a source type that cannot be expressed
	// through the proto enum's zero-value field. This is our own
	// extension; the AWS API ignores Filters entirely.
	var sourceTypeStr string
	for _, f := range req.Msg.Filters {
		if strings.EqualFold(f.Name, "source-type") && len(f.Values) > 0 {
			sourceTypeStr = f.Values[0]
		} else {
			logs.Warn("rds-admin: DescribeEvents ignores unsupported Filter",
				logs.String("filter", f.Name))
		}
	}
	if sourceTypeStr == "" {
		sourceTypeStr, err = sourceTypeToString(req.Msg.Sourcetype)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	}

	var startTime time.Time
	if st := req.Msg.Starttime; st != "" {
		startTime, err = time.Parse(time.RFC3339, st)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid StartTime: %v", err))
		}
	}
	var endTime time.Time
	if et := req.Msg.Endtime; et != "" {
		endTime, err = time.Parse(time.RFC3339, et)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid EndTime: %v", err))
		}
	}
	// AWS RDS Duration is in minutes and pairs with StartTime to derive an
	// implied end-of-window. An explicit EndTime overrides Duration.
	if req.Msg.Duration > 0 && !startTime.IsZero() && endTime.IsZero() {
		endTime = startTime.Add(time.Duration(req.Msg.Duration) * time.Minute)
	}

	// AWS RDS DescribeEvents MaxRecords constraints: min 20, max 100,
	// default 100. Clamp explicitly so callers passing 1 or 1000 still
	// receive a valid page rather than the raw value.
	maxRecords := int(req.Msg.Maxrecords)
	if maxRecords == 0 {
		maxRecords = 100
	} else if maxRecords < 20 {
		maxRecords = 20
	} else if maxRecords > 100 {
		maxRecords = 100
	}

	opts := storerds.EventListOptions{
		SourceType:       sourceTypeStr,
		SourceIdentifier: req.Msg.Sourceidentifier,
		StartTime:        startTime,
		EndTime:          endTime,
		EventCategories:  req.Msg.Eventcategories,
		Marker:           req.Msg.Marker,
		MaxRecords:       maxRecords,
	}

	result, err := store.ListEvents(opts)
	if err != nil {
		if errors.Is(err, storerds.ErrInvalidEventMarker) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid Marker: %s", req.Msg.Marker))
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	events := make([]*pb.Event, 0, len(result.Events))
	for _, evt := range result.Events {
		// stringToSourceType now returns an error on unknown values
		// rather than silently remapping to BLUE_GREEN_DEPLOYMENT.
		// RecordEvent rejects unknowns at write time so a healthy system
		// never hits this branch, but a corrupted legacy row could; log
		// and skip the row rather than fail the whole response.
		st, stErr := stringToSourceType(evt.SourceType)
		if stErr != nil {
			logs.Warn("rds-admin: skipping event with unknown SourceType",
				logs.String("event_id", evt.EventID),
				logs.String("source_type", evt.SourceType),
				logs.Err(stErr))
			continue
		}
		events = append(events, &pb.Event{
			Date:             evt.Date.UTC().Format(timeutils.ISO8601UTCFormat),
			Message:          evt.Message,
			Sourcearn:        evt.SourceArn,
			Sourceidentifier: evt.SourceIdentifier,
			Sourcetype:       st,
			Eventcategories:  evt.EventCategories,
		})
	}

	resp := &pb.EventsMessage{Events: events}
	if result.IsTruncated && result.Marker != "" {
		resp.Marker = result.Marker
	}
	return connect.NewResponse(resp), nil
}

// generateDbiResourceId allocates an AWS-shaped DB instance resource id
// ('db-' + 26 hex characters). AWS itself uses base32 but the wire format
// is opaque to clients and SDK tests do not assert the alphabet; hex is
// sufficient and avoids pulling in a base32 implementation with the
// required alphabet restrictions.
func generateDbiResourceId() string {
	b := make([]byte, 13) // 13 bytes -> 26 hex chars
	if _, err := rand.Read(b); err != nil {
		// crypto/rand.Read failing indicates a broken system RNG; fail
		// loud rather than fabricate a deterministic id.
		panic(fmt.Sprintf("crypto/rand.Read failed for DbiResourceId: %v", err))
	}
	return "db-" + hex.EncodeToString(b)
}

// sourceTypeToString converts the proto SourceType enum to the lowercase
// AWS wire string ("db-instance", "db-cluster", ...) used by the store filter.
//
// Proto3 zero-value ambiguity: SOURCE_TYPE_BLUE_GREEN_DEPLOYMENT == 0 is the
// enum's zero value. The proto3 wire format cannot distinguish "unset" from
// "explicitly zero", so any caller that omits SourceType (the common case
// for "list all events") delivers 0 to the server. We therefore treat 0 as
// the wildcard sentinel and return "" — matching the pre-L-3 behaviour and
// the AWS DescribeEvents API contract (omit SourceType => events of every
// source type).
//
// Workaround: DescribeEvents checks the Filters parameter for a "source-type"
// filter before calling this function. The admin console client can pass a
// Filter with name "source-type" and value "blue-green-deployment" to select
// blue-green events explicitly, bypassing the proto3 zero-value ambiguity
// without requiring a proto-breaking renumber.
func sourceTypeToString(st pb.SourceType) (string, error) {
	switch st {
	case pb.SourceType_SOURCE_TYPE_DB_PARAMETER_GROUP:
		return "db-parameter-group", nil
	case pb.SourceType_SOURCE_TYPE_DB_SHARD_GROUP:
		return "db-shard-group", nil
	case pb.SourceType_SOURCE_TYPE_CUSTOM_ENGINE_VERSION:
		return "custom-engine-version", nil
	case pb.SourceType_SOURCE_TYPE_DB_PROXY:
		return "db-proxy", nil
	case pb.SourceType_SOURCE_TYPE_DB_INSTANCE:
		return "db-instance", nil
	case pb.SourceType_SOURCE_TYPE_ZERO_ETL:
		return "zero-etl", nil
	case pb.SourceType_SOURCE_TYPE_DB_CLUSTER:
		return "db-cluster", nil
	case pb.SourceType_SOURCE_TYPE_DB_SECURITY_GROUP:
		return "db-security-group", nil
	case pb.SourceType_SOURCE_TYPE_DB_CLUSTER_SNAPSHOT:
		return "db-cluster-snapshot", nil
	case pb.SourceType_SOURCE_TYPE_DB_SNAPSHOT:
		return "db-snapshot", nil
	default:
		// Includes SOURCE_TYPE_BLUE_GREEN_DEPLOYMENT (= 0), which is the
		// proto3 zero value. Callers cannot distinguish "unset" from
		// "explicitly blue-green-deployment" on the wire, so treat both
		// as the wildcard and return "". Any other unmapped value is a
		// programming error and is surfaced.
		if st == 0 {
			return "", nil
		}
		return "", fmt.Errorf("unsupported SourceType: %d", st)
	}
}

// stringToSourceType is the inverse of sourceTypeToString. Unknown
// strings return an error rather than being silently remapped to a real
// category — previously unknowns were mapped to SOURCE_TYPE_BLUE_GREEN_DEPLOYMENT,
// the proto3 zero value, which fabricated audit data because the zero
// value is a *real* AWS source type and not a sentinel.
//
// RecordEvent rejects unknown values at write time, so a healthy system
// should never hit this error at read time. DescribeEvents logs and
// skips any event that does, which keeps the response honest without
// failing the whole request for a single corrupted row.
func stringToSourceType(s string) (pb.SourceType, error) {
	switch s {
	case "blue-green-deployment":
		return pb.SourceType_SOURCE_TYPE_BLUE_GREEN_DEPLOYMENT, nil
	case "db-parameter-group":
		return pb.SourceType_SOURCE_TYPE_DB_PARAMETER_GROUP, nil
	case "db-shard-group":
		return pb.SourceType_SOURCE_TYPE_DB_SHARD_GROUP, nil
	case "custom-engine-version":
		return pb.SourceType_SOURCE_TYPE_CUSTOM_ENGINE_VERSION, nil
	case "db-proxy":
		return pb.SourceType_SOURCE_TYPE_DB_PROXY, nil
	case "db-instance":
		return pb.SourceType_SOURCE_TYPE_DB_INSTANCE, nil
	case "zero-etl":
		return pb.SourceType_SOURCE_TYPE_ZERO_ETL, nil
	case "db-cluster":
		return pb.SourceType_SOURCE_TYPE_DB_CLUSTER, nil
	case "db-security-group":
		return pb.SourceType_SOURCE_TYPE_DB_SECURITY_GROUP, nil
	case "db-cluster-snapshot":
		return pb.SourceType_SOURCE_TYPE_DB_CLUSTER_SNAPSHOT, nil
	case "db-snapshot":
		return pb.SourceType_SOURCE_TYPE_DB_SNAPSHOT, nil
	default:
		return 0, fmt.Errorf("unsupported SourceType string: %q", s)
	}
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
// endpoint identifier. Filters (when supplied) are applied in addition to
// the explicit identifier-based filters, mirroring the AWS RDS API which
// supports both shapes simultaneously.
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
		if !applyRDSFilters(req.Msg.Filters, clusterEndpointFilterGetter(ep)) {
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

// clusterEndpointFilterGetter supports the AWS RDS filter names for
// DBClusterEndpoint resources: db-cluster-endpoint-id and db-cluster-id.
func clusterEndpointFilterGetter(ep *storerds.DBClusterEndpoint) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-cluster-endpoint-id":
			return ep.DBClusterEndpointIdentifier, true
		case "db-cluster-id":
			return ep.DBClusterIdentifier, true
		default:
			return "", false
		}
	}
}

// defaultClusterParamsForFamily returns the system-default parameters
// appropriate for a given DB parameter group family. The admin handler
// previously hard-coded Neptune defaults even when describing MySQL
// parameter groups, returning nonsensical `neptune_query_timeout` values
// for clusters that could never honour them.
//
// The parameter lists below are intentionally a representative subset
// of AWS's published defaults: the most-tuned parameters that AWS RDS
// customers actually inspect via DescribeDBClusterParameters. The
// complete AWS list per family runs to hundreds of entries; expanding
// to match AWS 1-for-1 is tracked as a separate data-entry task and is
// not required for correctness of the API response.
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

// defaultInstanceParamsForFamily mirrors defaultClusterParamsForFamily but
// for the per-instance DB parameter group.
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

	defaultParams := defaultClusterParamsForFamily(pg.DBParameterGroupFamily)
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

	defaultParams := defaultInstanceParamsForFamily(pg.DBParameterGroupFamily)
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
	sort.Slice(pbParams, func(i, j int) bool { return pbParams[i].Parametername < pbParams[j].Parametername })

	return connect.NewResponse(&pb.DescribeEngineDefaultClusterParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: family,
			Parameters:             pbParams,
		},
	}), nil
}

// DescribeEngineDefaultParameters returns the default engine parameters for a
// DB parameter group family.
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
	sort.Slice(pbParams, func(i, j int) bool { return pbParams[i].Parametername < pbParams[j].Parametername })

	return connect.NewResponse(&pb.DescribeEngineDefaultParametersResult{
		Enginedefaults: &pb.EngineDefaults{
			Dbparametergroupfamily: family,
			Parameters:             pbParams,
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
	if c.Endpoint != nil {
		// AWS DBCluster.Endpoint is a host:port writer DNS string; the
		// admin handler exposes the same shape so the UI can render a
		// connectable address. ReaderEndpoint mirrors the writer endpoint
		// for single-node clusters created via the admin console.
		endpointStr := c.Endpoint.Address
		if c.Endpoint.Port > 0 {
			endpointStr = fmt.Sprintf("%s:%d", c.Endpoint.Address, c.Endpoint.Port)
		}
		p.Endpoint = endpointStr
		p.Readerendpoint = endpointStr
	}
	return p
}

// instanceToPb converts a domain DBInstance to the RDS API protobuf DBInstance.
func instanceToPb(i *storerds.DBInstance, accountId string) *pb.DBInstance {
	p := &pb.DBInstance{
		Dbinstanceidentifier:               i.DBInstanceIdentifier,
		Dbclusteridentifier:                i.DBClusterIdentifier,
		Engine:                             i.Engine,
		Engineversion:                      i.EngineVersion,
		Dbinstanceclass:                    i.DBInstanceClass,
		Dbinstancestatus:                   i.DBInstanceStatus,
		Availabilityzone:                   i.AvailabilityZone,
		Preferredmaintenancewindow:         i.PreferredMaintenanceWindow,
		Preferredbackupwindow:              i.PreferredBackupWindow,
		Enabledcloudwatchlogsexports:       i.EnabledCloudwatchLogsExports,
		Iamdatabaseauthenticationenabled:   proto.Bool(i.IAMDatabaseAuthenticationEnabled),
		Publiclyaccessible:                 proto.Bool(i.PubliclyAccessible),
		Autominorversionupgrade:            proto.Bool(i.AutoMinorVersionUpgrade),
		Copytagstosnapshot:                 proto.Bool(i.CopyTagsToSnapshot),
		Dbinstancearn:                      i.DBInstanceArn,
		Allocatedstorage:                   i.AllocatedStorage,
		Masterusername:                     i.MasterUsername,
		Storagetype:                        i.StorageType,
		Backupretentionperiod:              i.BackupRetentionPeriod,
		Licensemodel:                       i.LicenseModel,
		Storageencrypted:                   proto.Bool(i.StorageEncrypted),
		Kmskeyid:                           i.KmsKeyId,
		Deletionprotection:                 proto.Bool(i.DeletionProtection),
		Multiaz:                            proto.Bool(i.MultiAZ),
		Secondaryavailabilityzone:          i.SecondaryAvailabilityZone,
		Iops:                               i.Iops,
		Maxallocatedstorage:                i.MaxAllocatedStorage,
		Storagethroughput:                  i.StorageThroughput,
		Monitoringinterval:                 i.MonitoringInterval,
		Enhancedmonitoringresourcearn:      i.EnhancedMonitoringResourceArn,
		Performanceinsightsenabled:         proto.Bool(i.PerformanceInsightsEnabled),
		Performanceinsightskmskeyid:        i.PerformanceInsightsKMSKeyId,
		Performanceinsightsretentionperiod: i.PerformanceInsightsRetentionPeriod,
		Cacertificateidentifier:            i.CACertificateIdentifier,
		Dbiresourceid:                      i.DbiResourceId,
		Dbinstanceport:                     i.Port,
		Vpcsecuritygroups:                  vpcSecurityGroupsToPb(i.VpcSecurityGroupIds),
		Optiongroupmemberships:             optionGroupMembershipsToPb(i.OptionGroupName),
	}
	if i.InstanceCreateTime != nil {
		p.Instancecreatetime = i.InstanceCreateTime.Format(timeutils.ISO8601UTCFormat)
	}
	if i.LatestRestorableTime != nil {
		p.Latestrestorabletime = i.LatestRestorableTime.Format(timeutils.ISO8601UTCFormat)
	}
	if i.Endpoint != nil {
		p.Endpoint = &pb.Endpoint{
			Address: i.Endpoint.Address,
			Port:    int32(i.Endpoint.Port),
		}
	}
	return p
}

// vpcSecurityGroupsToPb converts a list of VPC security group IDs to the
// RDS API protobuf VpcSecurityGroupMembership shape. Status is reported as
// "active" to match the Neptune handler's convention.
func vpcSecurityGroupsToPb(ids []string) []*pb.VpcSecurityGroupMembership {
	if len(ids) == 0 {
		return nil
	}
	out := make([]*pb.VpcSecurityGroupMembership, 0, len(ids))
	for _, id := range ids {
		out = append(out, &pb.VpcSecurityGroupMembership{
			Vpcsecuritygroupid: id,
			Status:             "active",
		})
	}
	return out
}

// optionGroupMembershipsToPb wraps a single option-group name into the
// AWS-standard OptionGroupMembership list shape with "in-sync" status.
func optionGroupMembershipsToPb(name string) []*pb.OptionGroupMembership {
	if name == "" {
		return nil
	}
	return []*pb.OptionGroupMembership{
		{Optiongroupname: name, Status: "in-sync"},
	}
}

// applyRDSFilters reports whether a candidate resource matches every filter
// in the supplied list. A resource matches a single filter when the value
// returned by getter(name) is equal (case-insensitive) to at least one of
// the filter's Values; an empty filter list matches everything.
//
// Semantics mirror AWS RDS: OR within a single filter's Values, AND across
// multiple filters. Unknown filter names cause the resource to be excluded
// rather than silently matching.
func applyRDSFilters(filters []*pb.Filter, getter func(name string) (string, bool)) bool {
	if len(filters) == 0 {
		return true
	}
	for _, f := range filters {
		if f == nil {
			continue
		}
		v, ok := getter(f.Name)
		if !ok {
			return false
		}
		matched := false
		for _, want := range f.Values {
			if strings.EqualFold(v, want) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// clusterFilterGetter returns the value for the supplied AWS RDS filter name
// against a DBCluster. Supported filter names: db-cluster-id, engine,
// db-cluster-parameter-group-name, db-subnet-group-name.
func clusterFilterGetter(c *storerds.DBCluster) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-cluster-id":
			return c.DBClusterIdentifier, true
		case "engine":
			return c.Engine, true
		case "db-cluster-parameter-group-name":
			return c.DBClusterParameterGroupName, true
		case "db-subnet-group-name":
			return c.DBSubnetGroupName, true
		default:
			return "", false
		}
	}
}

// instanceFilterGetter returns the value for the supplied AWS RDS filter name
// against a DBInstance. Supported filter names: db-instance-id, engine,
// db-parameter-group-name, db-cluster-id.
func instanceFilterGetter(i *storerds.DBInstance) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-instance-id":
			return i.DBInstanceIdentifier, true
		case "engine":
			return i.Engine, true
		case "db-parameter-group-name":
			return i.DBParameterGroupName, true
		case "db-cluster-id":
			return i.DBClusterIdentifier, true
		default:
			return "", false
		}
	}
}

// clusterSnapshotFilterGetter returns the value for the supplied AWS RDS
// filter name against a DBClusterSnapshot. Supported filter names:
// db-cluster-snapshot-id, db-cluster-id, engine, snapshot-type.
func clusterSnapshotFilterGetter(s *storerds.DBClusterSnapshot) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-cluster-snapshot-id":
			return s.DBClusterSnapshotIdentifier, true
		case "db-cluster-id":
			return s.DBClusterIdentifier, true
		case "engine":
			return s.Engine, true
		case "snapshot-type":
			return s.SnapshotType, true
		default:
			return "", false
		}
	}
}

// instanceSnapshotFilterGetter returns the value for the supplied AWS RDS
// filter name against a DBInstanceSnapshot. Supported filter names:
// db-snapshot-id, db-instance-id, engine, snapshot-type.
func instanceSnapshotFilterGetter(s *storerds.DBInstanceSnapshot) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-snapshot-id":
			return s.DBSnapshotIdentifier, true
		case "db-instance-id":
			return s.DBInstanceIdentifier, true
		case "engine":
			return s.Engine, true
		case "snapshot-type":
			return s.SnapshotType, true
		default:
			return "", false
		}
	}
}

// clusterParamGroupFilterGetter supports db-cluster-parameter-group-name.
func clusterParamGroupFilterGetter(g *storerds.DBClusterParameterGroup) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-cluster-parameter-group-name":
			return g.DBClusterParameterGroupName, true
		default:
			return "", false
		}
	}
}

// paramGroupFilterGetter supports db-parameter-group-name.
func paramGroupFilterGetter(g *storerds.DBParameterGroup) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-parameter-group-name":
			return g.DBParameterGroupName, true
		default:
			return "", false
		}
	}
}

// subnetGroupFilterGetter supports db-subnet-group-name.
func subnetGroupFilterGetter(g *storerds.DBSubnetGroup) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "db-subnet-group-name":
			return g.DBSubnetGroupName, true
		default:
			return "", false
		}
	}
}

// globalClusterFilterGetter supports global-cluster-id and engine.
func globalClusterFilterGetter(c *storerds.GlobalCluster) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "global-cluster-id":
			return c.GlobalClusterIdentifier, true
		case "engine":
			return c.Engine, true
		default:
			return "", false
		}
	}
}

// eventSubscriptionFilterGetter supports event-subscription-id.
func eventSubscriptionFilterGetter(s *storerds.EventSubscription) func(string) (string, bool) {
	return func(name string) (string, bool) {
		switch name {
		case "event-subscription-id":
			return s.CustSubscriptionId, true
		default:
			return "", false
		}
	}
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
// snapOp may be nil; when non-nil it enables row-level snapshot data capture
// for CreateDBSnapshot and DeleteDBInstance final snapshots.
func NewConnectHandler(stores StoreProvider, engines EngineProvider, accountID string, snapOp SnapshotOperator) (string, http.Handler) {
	h := NewAdminHandler(stores, engines, accountID)
	h.SetSnapshotOperator(snapOp)
	return rdsconnect.NewRDSServiceHandler(h)
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
	if err := ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	instance, err := store.GetInstance(instanceID)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Reject snapshots of instances in terminal or in-flight states.
	// AWS returns InvalidDBInstanceStateFault when the source instance
	// is not in a state that allows snapshot creation.
	switch instance.DBInstanceStatus {
	case "deleting", "failed", "inaccessible-encryption-credentials":
		return nil, connect.NewError(connect.CodeFailedPrecondition,
			fmt.Errorf("cannot create snapshot for instance %q in status %q", instanceID, instance.DBInstanceStatus))
	}

	now := time.Now()
	arn := arnutil.NewARNBuilder(h.accountId, region).RDS().Snapshot(snapshotID)
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

		// AWS-standard DBSnapshot fields captured from the source instance.
		// Previously only Engine/EngineVersion/AZ/IAM were copied, so
		// RestoreDBInstanceFromDBSnapshot could not reconstruct storage
		// shape, master username, port, encryption key, etc.
		AllocatedStorage: instance.AllocatedStorage,
		MasterUsername:   instance.MasterUsername,
		StorageType:      instance.StorageType,
		LicenseModel:     instance.LicenseModel,
		StorageEncrypted: instance.StorageEncrypted,
		KmsKeyId:         instance.KmsKeyId,
		OptionGroupName:  instance.OptionGroupName,
		VpcId:            instance.VpcId,
	}
	// Port: prefer the explicit Port field; fall back to Endpoint.Port for
	// instances that only carry the listener endpoint.
	if instance.Port > 0 {
		snap.Port = instance.Port
	} else if instance.Endpoint != nil {
		snap.Port = int32(instance.Endpoint.Port)
	}

	if err := store.CreateInstanceSnapshot(snap); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Capture row-level data for MySQL instances so that
	// RestoreDBInstanceFromDBSnapshot can recover user tables. Without
	// this call the snapshot is metadata-only and restore produces an
	// empty instance. Fail-closed: if data capture fails, remove the
	// metadata snapshot and return an error so the caller knows the
	// snapshot is not usable.
	if h.snapOp != nil && instance.Engine == "mysql" {
		if err := h.snapOp.SnapshotData(instanceID, snapshotID); err != nil {
			_ = store.DeleteInstanceSnapshot(snapshotID)
			return nil, connect.NewError(connect.CodeInternal,
				fmt.Errorf("snapshot data capture failed for instance %q: %v", instanceID, err))
		}
	}

	// AWS: when the source instance has CopyTagsToSnapshot=true, the
	// instance's tags are copied to the snapshot. Use the same Pebble tag
	// bucket so the snapshot's TagList reflects the source instance.
	if instance.CopyTagsToSnapshot && instance.DBInstanceArn != "" {
		if tags, terr := store.GetTags(instance.DBInstanceArn); terr == nil && len(tags) > 0 {
			snap.TagList = tags
			// Re-put the snapshot to persist the tag list we just attached.
			if uerr := store.UpdateInstanceSnapshot(snap); uerr != nil {
				logs.Warn("rds-admin: failed to persist CopyTagsToSnapshot on snapshot",
					logs.String("snapshot", snapshotID), logs.Err(uerr))
			}
		}
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
		if !applyRDSFilters(req.Msg.Filters, instanceSnapshotFilterGetter(s)) {
			continue
		}
		pbSnapshots = append(pbSnapshots, dbSnapshotToPb(s, h.accountId))
	}

	return connect.NewResponse(&pb.DBSnapshotMessage{
		Dbsnapshots: pbSnapshots,
	}), nil
}
