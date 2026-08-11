package rds

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/aws/rds"
	storerds "vorpalstacks/internal/store/aws/rds"
	"vorpalstacks/internal/utils/timeutils"
)

// ---------------------------------------------------------------------------
// Input DTOs
// ---------------------------------------------------------------------------

type DescribeDBClustersInput struct {
	DBClusterIdentifier string
	Filters             []*pb.Filter
	Marker              string
	MaxRecords          int32
}

type CreateDBClusterInput struct {
	DBClusterIdentifier          string
	Engine                       string
	EngineVersion                string
	DatabaseName                 string
	MasterUsername               string
	Port                         int32
	BackupRetentionPeriod        int32
	AvailabilityZones            []string
	DBSubnetGroupName            string
	DBClusterParameterGroupName  string
	StorageEncrypted             bool
	CopyTagsToSnapshot           bool
	DeletionProtection           bool
	EnabledCloudwatchLogsExports []string
	IAMDatabaseAuthentication    bool
}

type DeleteDBClusterInput struct {
	DBClusterIdentifier       string
	SkipFinalSnapshot         bool
	FinalDBSnapshotIdentifier string
}

type DescribeDBClusterSnapshotsInput struct {
	DBClusterSnapshotIdentifier string
	DBClusterIdentifier         string
	SnapshotType                string
	Filters                     []*pb.Filter
	Marker                      string
	MaxRecords                  int32
}

type DescribeDBClusterEndpointsInput struct {
	DBClusterIdentifier         string
	DBClusterEndpointIdentifier string
	Filters                     []*pb.Filter
	Marker                      string
	MaxRecords                  int32
}

type DescribeDBClusterParametersInput struct {
	DBClusterParameterGroupName string
}

// ---------------------------------------------------------------------------
// Core methods
// ---------------------------------------------------------------------------

func (s *RDSService) describeDBClustersCore(stores *rdsStores, in DescribeDBClustersInput) (*pb.DBClusterMessage, error) {
	clusters, nextMarker, err := QueryClusters(stores.store, in)
	if err != nil {
		return nil, err
	}
	pbClusters := make([]*pb.DBCluster, 0, len(clusters))
	for _, c := range clusters {
		pbClusters = append(pbClusters, clusterToPb(c, s.accountId))
	}
	return &pb.DBClusterMessage{Dbclusters: pbClusters, Marker: nextMarker}, nil
}

func (s *RDSService) createDBClusterCore(stores *rdsStores, in CreateDBClusterInput) (*pb.CreateDBClusterResult, error) {
	params := CreateClusterParams{
		DBClusterIdentifier:          in.DBClusterIdentifier,
		Engine:                       in.Engine,
		EngineVersion:                in.EngineVersion,
		DatabaseName:                 in.DatabaseName,
		MasterUsername:               in.MasterUsername,
		Port:                         int(in.Port),
		BackupRetentionPeriod:        int(in.BackupRetentionPeriod),
		AvailabilityZones:            in.AvailabilityZones,
		DBSubnetGroupName:            in.DBSubnetGroupName,
		DBClusterParameterGroupName:  in.DBClusterParameterGroupName,
		StorageEncrypted:             in.StorageEncrypted,
		CopyTagsToSnapshot:           in.CopyTagsToSnapshot,
		DeletionProtection:           in.DeletionProtection,
		IAMDatabaseAuthentication:    in.IAMDatabaseAuthentication,
		EnabledCloudwatchLogsExports: in.EnabledCloudwatchLogsExports,
		AccountID:                    s.accountId,
		Region:                       stores.region,
	}
	if err := ValidateCreateClusterParams(stores.store, params); err != nil {
		return nil, newValidationError("%v", err)
	}
	if _, err := s.engines(params.Engine); err != nil {
		return nil, newValidationError("engine %q is not supported on this platform: %v", params.Engine, err)
	}

	cluster := BuildCluster(params)

	if err := stores.store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	engineStarted := false
	if eng, engErr := s.engines(cluster.Engine); engErr == nil {
		port, openErr := eng.Open(stores.region, cluster.DBClusterIdentifier)
		if openErr != nil {
			logs.Warn("rds-admin: failed to start engine for cluster",
				logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(openErr))
		} else {
			cluster.Endpoint = &storerds.Endpoint{
				Address: fmt.Sprintf("%s.%s.%s.rds.amazonaws.com", cluster.DBClusterIdentifier, s.accountId, stores.region),
				Port:    port,
			}
			if err := stores.store.UpdateCluster(cluster); err != nil {
				logs.Warn("rds-admin: rolling back cluster engine after UpdateCluster failure",
					logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
				if closeErr := eng.Close(cluster.DBClusterIdentifier); closeErr != nil {
					logs.Warn("rds-admin: cluster engine rollback Close also failed",
						logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(closeErr))
				}
				cluster.Endpoint = nil
				cluster.Status = "failed"
				if persistErr := stores.store.UpdateCluster(cluster); persistErr != nil {
					logs.Warn("rds-admin: failed to persist cluster failed status",
						logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(persistErr))
				}
				return nil, translateStoreError(err)
			}
			engineStarted = true
		}
	}

	if engineStarted {
		cluster.Status = "available"
	} else {
		cluster.Status = "failed"
	}
	if err := stores.store.UpdateCluster(cluster); err != nil {
		if engineStarted {
			if eng, engErr := s.engines(cluster.Engine); engErr == nil {
				if closeErr := eng.Close(cluster.DBClusterIdentifier); closeErr != nil {
					logs.Warn("rds-admin: cleanup engine.Close after final cluster UpdateCluster failure",
						logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(closeErr))
				}
			}
		}
		return nil, translateStoreError(err)
	}

	return &pb.CreateDBClusterResult{
		Dbcluster: clusterToPb(cluster, s.accountId),
	}, nil
}

func (s *RDSService) deleteDBClusterCore(stores *rdsStores, in DeleteDBClusterInput) (*pb.DeleteDBClusterResult, error) {
	cluster, err := ValidateDeleteClusterParams(stores.store, DeleteClusterParams{
		DBClusterIdentifier:       in.DBClusterIdentifier,
		SkipFinalSnapshot:         in.SkipFinalSnapshot,
		FinalDBSnapshotIdentifier: in.FinalDBSnapshotIdentifier,
		AccountID:                 s.accountId,
		Region:                    stores.region,
	})
	if err != nil {
		return nil, translateStoreError(err)
	}

	if !in.SkipFinalSnapshot {
		if in.FinalDBSnapshotIdentifier == "" {
			return nil, newValidationError("FinalDBSnapshotIdentifier is required when SkipFinalSnapshot is false")
		}
		finalSnap := BuildFinalSnapshot(cluster, in.FinalDBSnapshotIdentifier, s.accountId, stores.region)
		if err := stores.store.CreateSnapshot(finalSnap); err != nil {
			return nil, translateStoreError(err)
		}
	}

	cluster.Status = "deleting"
	if err := stores.store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	if eng, engErr := s.engines(cluster.Engine); engErr == nil {
		if closeErr := eng.Close(in.DBClusterIdentifier); closeErr != nil {
			logs.Warn("rds-admin: engine.Close failed during DeleteDBCluster",
				logs.String("cluster", in.DBClusterIdentifier), logs.Err(closeErr))
		}
	}

	if err := stores.store.DeleteCluster(in.DBClusterIdentifier); err != nil {
		return nil, translateStoreError(err)
	}

	return &pb.DeleteDBClusterResult{
		Dbcluster: clusterToPb(cluster, s.accountId),
	}, nil
}

func (s *RDSService) describeDBClusterSnapshotsCore(stores *rdsStores, in DescribeDBClusterSnapshotsInput) (*pb.DBClusterSnapshotMessage, error) {
	snapshots, nextMarker, err := QueryClusterSnapshots(stores.store, in)
	if err != nil {
		return nil, err
	}
	pbSnapshots := make([]*pb.DBClusterSnapshot, 0, len(snapshots))
	for _, snap := range snapshots {
		pbSnapshots = append(pbSnapshots, snapshotToPb(snap, s.accountId))
	}
	return &pb.DBClusterSnapshotMessage{Dbclustersnapshots: pbSnapshots, Marker: nextMarker}, nil
}

func (s *RDSService) describeDBClusterEndpointsCore(stores *rdsStores, in DescribeDBClusterEndpointsInput) (*pb.DBClusterEndpointMessage, error) {
	endpoints, err := stores.store.ListClusterEndpoints(in.DBClusterIdentifier)
	if err != nil {
		return nil, translateStoreError(err)
	}

	pbEndpoints := make([]*pb.DBClusterEndpoint, 0, len(endpoints))
	for _, ep := range endpoints {
		if in.DBClusterEndpointIdentifier != "" && ep.DBClusterEndpointIdentifier != in.DBClusterEndpointIdentifier {
			continue
		}
		if !applyRDSFilters(in.Filters, clusterEndpointFilterGetter(ep)) {
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

	pbEndpoints, nextMarker := paginateRDSItems(pbEndpoints, in.Marker, in.MaxRecords, func(e *pb.DBClusterEndpoint) string {
		return e.Dbclusterendpointidentifier
	})
	return &pb.DBClusterEndpointMessage{Dbclusterendpoints: pbEndpoints, Marker: nextMarker}, nil
}

func (s *RDSService) describeDBClusterParametersCore(stores *rdsStores, in DescribeDBClusterParametersInput) (*pb.DBClusterParameterGroupDetails, error) {
	pg, err := stores.store.GetClusterParameterGroup(in.DBClusterParameterGroupName)
	if err != nil {
		return nil, translateStoreError(err)
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
	sortParameters(pbParams)

	return &pb.DBClusterParameterGroupDetails{
		Parameters: pbParams,
		Marker:     "",
	}, nil
}

// ---------------------------------------------------------------------------
// Conversion helpers (store → protobuf)
// ---------------------------------------------------------------------------

func clusterToPb(c *storerds.DBCluster, accountId string) *pb.DBCluster {
	p := &pb.DBCluster{
		Dbclusteridentifier:              c.DBClusterIdentifier,
		Engine:                           c.Engine,
		Engineversion:                    c.EngineVersion,
		Status:                           c.Status,
		Masterusername:                   c.MasterUsername,
		Databasename:                     c.DatabaseName,
		Port:                             proto.Int32(int32(c.Port)),
		Backupretentionperiod:            proto.Int32(int32(c.BackupRetentionPeriod)),
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
		endpointStr := c.Endpoint.Address
		if c.Endpoint.Port > 0 {
			endpointStr = fmt.Sprintf("%s:%d", c.Endpoint.Address, c.Endpoint.Port)
		}
		p.Endpoint = endpointStr
		p.Readerendpoint = endpointStr
	}
	return p
}

func snapshotToPb(s *storerds.DBClusterSnapshot, accountId string) *pb.DBClusterSnapshot {
	p := &pb.DBClusterSnapshot{
		Dbclustersnapshotidentifier: s.DBClusterSnapshotIdentifier,
		Dbclusteridentifier:         s.DBClusterIdentifier,
		Engine:                      s.Engine,
		Engineversion:               s.EngineVersion,
		Status:                      s.Status,
		Port:                        proto.Int32(int32(s.Port)),
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

// ---------------------------------------------------------------------------
// Filter getters
// ---------------------------------------------------------------------------

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
