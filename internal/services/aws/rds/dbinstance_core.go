package rds

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/core/logs"
	pb "vorpalstacks/internal/pb/aws/rds"
	storerds "vorpalstacks/internal/store/aws/rds"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// ---------------------------------------------------------------------------
// Input DTOs
// ---------------------------------------------------------------------------

type DescribeDBInstancesInput struct {
	DBInstanceIdentifier string
	Filters              []*pb.Filter
	Marker               string
	MaxRecords           int32
}

type CreateDBInstanceInput struct {
	DBInstanceIdentifier               string
	DBClusterIdentifier                string
	Engine                             string
	EngineVersion                      string
	DBInstanceClass                    string
	AvailabilityZone                   string
	PreferredMaintenanceWindow         string
	PreferredBackupWindow              string
	DBParameterGroupName               string
	DBSubnetGroupName                  string
	PubliclyAccessible                 bool
	AutoMinorVersionUpgrade            bool
	AllocatedStorage                   int32
	MasterUsername                     string
	StorageType                        string
	BackupRetentionPeriod              int32
	LicenseModel                       string
	StorageEncrypted                   bool
	KmsKeyId                           string
	DeletionProtection                 bool
	MultiAZ                            bool
	Port                               int32
	OptionGroupName                    string
	Iops                               int32
	MaxAllocatedStorage                int32
	StorageThroughput                  int32
	MonitoringInterval                 int32
	MonitoringRoleArn                  string
	EnablePerformanceInsights          bool
	PerformanceInsightsKMSKeyId        string
	PerformanceInsightsRetentionPeriod int32
	CACertificateIdentifier            string
	CopyTagsToSnapshot                 bool
	EnabledCloudwatchLogsExports       []string
	EnableIAMDatabaseAuthentication    bool
	VpcSecurityGroupIds                []string
	DBSecurityGroups                   []string
}

type DeleteDBInstanceInput struct {
	DBInstanceIdentifier      string
	SkipFinalSnapshot         bool
	FinalDBSnapshotIdentifier string
}

type CreateDBSnapshotInput struct {
	DBInstanceIdentifier string
	DBSnapshotIdentifier string
}

type DescribeDBSnapshotsInput struct {
	DBSnapshotIdentifier string
	DBInstanceIdentifier string
	SnapshotType         string
	Filters              []*pb.Filter
	Marker               string
	MaxRecords           int32
}

// ---------------------------------------------------------------------------
// Core methods
// ---------------------------------------------------------------------------

func (s *RDSService) describeDBInstancesCore(stores *rdsStores, in DescribeDBInstancesInput) (*pb.DBInstanceMessage, error) {
	instances, nextMarker, err := QueryInstances(stores.store, in)
	if err != nil {
		return nil, err
	}
	pbInstances := make([]*pb.DBInstance, 0, len(instances))
	for _, i := range instances {
		pbInstances = append(pbInstances, instanceToPb(i, s.accountId))
	}
	return &pb.DBInstanceMessage{Dbinstances: pbInstances, Marker: nextMarker}, nil
}

func (s *RDSService) createDBInstanceCore(stores *rdsStores, in CreateDBInstanceInput) (*pb.CreateDBInstanceResult, error) {
	params := CreateInstanceParams{
		DBInstanceIdentifier:               in.DBInstanceIdentifier,
		DBClusterIdentifier:                in.DBClusterIdentifier,
		Engine:                             in.Engine,
		EngineVersion:                      in.EngineVersion,
		DBInstanceClass:                    in.DBInstanceClass,
		AvailabilityZone:                   in.AvailabilityZone,
		PreferredMaintenanceWindow:         in.PreferredMaintenanceWindow,
		PreferredBackupWindow:              in.PreferredBackupWindow,
		DBParameterGroupName:               in.DBParameterGroupName,
		DBSubnetGroupName:                  in.DBSubnetGroupName,
		PubliclyAccessible:                 in.PubliclyAccessible,
		AutoMinorVersionUpgrade:            in.AutoMinorVersionUpgrade,
		AllocatedStorage:                   in.AllocatedStorage,
		MasterUsername:                     in.MasterUsername,
		StorageType:                        in.StorageType,
		BackupRetentionPeriod:              in.BackupRetentionPeriod,
		LicenseModel:                       in.LicenseModel,
		StorageEncrypted:                   in.StorageEncrypted,
		KmsKeyId:                           in.KmsKeyId,
		DeletionProtection:                 in.DeletionProtection,
		MultiAZ:                            in.MultiAZ,
		Port:                               in.Port,
		OptionGroupName:                    in.OptionGroupName,
		Iops:                               in.Iops,
		MaxAllocatedStorage:                in.MaxAllocatedStorage,
		StorageThroughput:                  in.StorageThroughput,
		MonitoringInterval:                 in.MonitoringInterval,
		MonitoringRoleArn:                  in.MonitoringRoleArn,
		EnablePerformanceInsights:          in.EnablePerformanceInsights,
		PerformanceInsightsKMSKeyId:        in.PerformanceInsightsKMSKeyId,
		PerformanceInsightsRetentionPeriod: in.PerformanceInsightsRetentionPeriod,
		CACertificateIdentifier:            in.CACertificateIdentifier,
		CopyTagsToSnapshot:                 in.CopyTagsToSnapshot,
		EnabledCloudwatchLogsExports:       in.EnabledCloudwatchLogsExports,
		EnableIAMDatabaseAuthentication:    in.EnableIAMDatabaseAuthentication,
		VpcSecurityGroupIds:                in.VpcSecurityGroupIds,
		DBSecurityGroups:                   in.DBSecurityGroups,
		AccountID:                          s.accountId,
		Region:                             stores.region,
	}
	if err := ValidateCreateInstanceParams(stores.store, params); err != nil {
		return nil, newValidationError("%v", err)
	}
	if _, err := s.engines(params.Engine); err != nil {
		return nil, newValidationError("engine %q is not supported on this platform: %v", params.Engine, err)
	}

	instance := BuildInstance(params)

	if err := stores.store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	engineStarted := false
	if eng, engErr := s.engines(instance.Engine); engErr == nil {
		port, openErr := eng.Open(stores.region, instance.DBInstanceIdentifier)
		if openErr != nil {
			logs.Warn("rds-admin: failed to start engine for instance",
				logs.String("instance", instance.DBInstanceIdentifier), logs.Err(openErr))
		} else {
			instance.Endpoint = &storerds.Endpoint{
				Address: fmt.Sprintf("%s.%s.%s.rds.amazonaws.com", instance.DBInstanceIdentifier, s.accountId, stores.region),
				Port:    port,
			}
			if err := stores.store.UpdateInstance(instance); err != nil {
				logs.Warn("rds-admin: rolling back engine after UpdateInstance failure",
					logs.String("instance", instance.DBInstanceIdentifier), logs.Err(err))
				if closeErr := eng.Close(instance.DBInstanceIdentifier); closeErr != nil {
					logs.Warn("rds-admin: engine rollback Close also failed",
						logs.String("instance", instance.DBInstanceIdentifier), logs.Err(closeErr))
				}
				instance.Endpoint = nil
				instance.DBInstanceStatus = "failed"
				if persistErr := stores.store.UpdateInstance(instance); persistErr != nil {
					logs.Warn("rds-admin: failed to persist failed status",
						logs.String("instance", instance.DBInstanceIdentifier), logs.Err(persistErr))
				}
				return nil, translateStoreError(err)
			}
			engineStarted = true
		}
	}

	if engineStarted {
		instance.DBInstanceStatus = "available"
	} else {
		instance.DBInstanceStatus = "failed"
	}
	if err := stores.store.UpdateInstance(instance); err != nil {
		if engineStarted {
			if eng, engErr := s.engines(instance.Engine); engErr == nil {
				if closeErr := eng.Close(instance.DBInstanceIdentifier); closeErr != nil {
					logs.Warn("rds-admin: cleanup engine.Close after final UpdateInstance failure",
						logs.String("instance", instance.DBInstanceIdentifier), logs.Err(closeErr))
				}
			}
		}
		return nil, translateStoreError(err)
	}

	return &pb.CreateDBInstanceResult{
		Dbinstance: instanceToPb(instance, s.accountId),
	}, nil
}

func (s *RDSService) deleteDBInstanceCore(stores *rdsStores, in DeleteDBInstanceInput) (*pb.DeleteDBInstanceResult, error) {
	instance, err := ValidateDeleteInstanceParams(stores.store, DeleteInstanceParams{
		DBInstanceIdentifier:      in.DBInstanceIdentifier,
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
		if err := ValidateDBSnapshotIdentifier(in.FinalDBSnapshotIdentifier); err != nil {
			return nil, newValidationError("%v", err)
		}
		finalSnap := BuildInstanceFinalSnapshot(instance, in.FinalDBSnapshotIdentifier, s.accountId, stores.region)
		if err := stores.store.CreateInstanceSnapshot(finalSnap); err != nil {
			return nil, translateStoreError(err)
		}
		if s.snapOp != nil && instance.Engine == "mysql" {
			if err := s.snapOp.SnapshotData(in.DBInstanceIdentifier, in.FinalDBSnapshotIdentifier); err != nil {
				_ = stores.store.DeleteInstanceSnapshot(in.FinalDBSnapshotIdentifier)
				return nil, newInternalError("final snapshot data capture failed for instance %q: %v", in.DBInstanceIdentifier, err)
			}
		}
	}

	instance.DBInstanceStatus = "deleting"
	if err := stores.store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	if eng, engErr := s.engines(instance.Engine); engErr == nil {
		if closeErr := eng.Close(in.DBInstanceIdentifier); closeErr != nil {
			logs.Warn("rds-admin: engine.Close failed during DeleteDBInstance",
				logs.String("instance", in.DBInstanceIdentifier), logs.Err(closeErr))
		}
	}

	if err := stores.store.DeleteInstance(in.DBInstanceIdentifier); err != nil {
		return nil, translateStoreError(err)
	}

	return &pb.DeleteDBInstanceResult{
		Dbinstance: instanceToPb(instance, s.accountId),
	}, nil
}

func (s *RDSService) createDBSnapshotCore(stores *rdsStores, in CreateDBSnapshotInput) (*pb.CreateDBSnapshotResult, error) {
	instanceID := in.DBInstanceIdentifier
	snapshotID := in.DBSnapshotIdentifier
	if snapshotID == "" || instanceID == "" {
		return nil, newValidationError("DBSnapshotIdentifier and DBInstanceIdentifier are required")
	}
	if err := ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return nil, newValidationError("%v", err)
	}

	instance, err := stores.store.GetInstance(instanceID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	switch instance.DBInstanceStatus {
	case "deleting", "failed", "inaccessible-encryption-credentials":
		return nil, newFailedPreconditionError("cannot create snapshot for instance %q in status %q", instanceID, instance.DBInstanceStatus)
	}

	now := time.Now()
	arn := arnutil.NewARNBuilder(s.accountId, stores.region).RDS().Snapshot(snapshotID)
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
		AccountID:              s.accountId,
		Region:                 stores.region,
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
		snap.Port = instance.Port
	} else if instance.Endpoint != nil {
		snap.Port = int32(instance.Endpoint.Port)
	}

	if err := stores.store.CreateInstanceSnapshot(snap); err != nil {
		return nil, translateStoreError(err)
	}

	if s.snapOp != nil && instance.Engine == "mysql" {
		if err := s.snapOp.SnapshotData(instanceID, snapshotID); err != nil {
			_ = stores.store.DeleteInstanceSnapshot(snapshotID)
			return nil, newInternalError("snapshot data capture failed for instance %q: %v", instanceID, err)
		}
	}

	if instance.CopyTagsToSnapshot && instance.DBInstanceArn != "" {
		if tags, terr := stores.store.GetTags(instance.DBInstanceArn); terr == nil && len(tags) > 0 {
			snap.TagList = tags
			if uerr := stores.store.UpdateInstanceSnapshot(snap); uerr != nil {
				logs.Warn("rds-admin: failed to persist CopyTagsToSnapshot on snapshot",
					logs.String("snapshot", snapshotID), logs.Err(uerr))
			}
		}
	}

	return &pb.CreateDBSnapshotResult{
		Dbsnapshot: dbSnapshotToPb(snap, s.accountId),
	}, nil
}

func (s *RDSService) describeDBSnapshotsCore(stores *rdsStores, in DescribeDBSnapshotsInput) (*pb.DBSnapshotMessage, error) {
	snapshots, nextMarker, err := QuerySnapshots(stores.store, in)
	if err != nil {
		return nil, err
	}
	pbSnapshots := make([]*pb.DBSnapshot, 0, len(snapshots))
	for _, snap := range snapshots {
		pbSnapshots = append(pbSnapshots, dbSnapshotToPb(snap, s.accountId))
	}
	return &pb.DBSnapshotMessage{Dbsnapshots: pbSnapshots, Marker: nextMarker}, nil
}

// ---------------------------------------------------------------------------
// Conversion helpers (store -> protobuf)
// ---------------------------------------------------------------------------

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
		Allocatedstorage:                   proto.Int32(i.AllocatedStorage),
		Masterusername:                     i.MasterUsername,
		Storagetype:                        i.StorageType,
		Backupretentionperiod:              proto.Int32(i.BackupRetentionPeriod),
		Licensemodel:                       i.LicenseModel,
		Storageencrypted:                   proto.Bool(i.StorageEncrypted),
		Kmskeyid:                           i.KmsKeyId,
		Deletionprotection:                 proto.Bool(i.DeletionProtection),
		Multiaz:                            proto.Bool(i.MultiAZ),
		Secondaryavailabilityzone:          i.SecondaryAvailabilityZone,
		Iops:                               proto.Int32(i.Iops),
		Maxallocatedstorage:                proto.Int32(i.MaxAllocatedStorage),
		Storagethroughput:                  proto.Int32(i.StorageThroughput),
		Monitoringinterval:                 proto.Int32(i.MonitoringInterval),
		Enhancedmonitoringresourcearn:      i.EnhancedMonitoringResourceArn,
		Performanceinsightsenabled:         proto.Bool(i.PerformanceInsightsEnabled),
		Performanceinsightskmskeyid:        i.PerformanceInsightsKMSKeyId,
		Performanceinsightsretentionperiod: proto.Int32(i.PerformanceInsightsRetentionPeriod),
		Cacertificateidentifier:            i.CACertificateIdentifier,
		Dbiresourceid:                      i.DbiResourceId,
		Dbinstanceport:                     proto.Int32(i.Port),
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
			Port:    proto.Int32(int32(i.Endpoint.Port)),
		}
	}
	return p
}

func dbSnapshotToPb(s *storerds.DBInstanceSnapshot, accountId string) *pb.DBSnapshot {
	p := &pb.DBSnapshot{
		Dbsnapshotidentifier:             s.DBSnapshotIdentifier,
		Dbinstanceidentifier:             s.DBInstanceIdentifier,
		Engine:                           s.Engine,
		Engineversion:                    s.EngineVersion,
		Snapshottype:                     s.SnapshotType,
		Status:                           s.Status,
		Allocatedstorage:                 proto.Int32(int32(s.AllocatedStorage)),
		Storagetype:                      s.StorageType,
		Port:                             proto.Int32(int32(s.Port)),
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

func optionGroupMembershipsToPb(name string) []*pb.OptionGroupMembership {
	if name == "" {
		return nil
	}
	return []*pb.OptionGroupMembership{
		{Optiongroupname: name, Status: "in-sync"},
	}
}

// generateDbiResourceId allocates an AWS-shaped DB instance resource id
// ('db-' + 26 hex characters).
func generateDbiResourceId() string {
	b := make([]byte, 13)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand.Read failed for DbiResourceId: %v", err))
	}
	return "db-" + hex.EncodeToString(b)
}

// ---------------------------------------------------------------------------
// Filter getters
// ---------------------------------------------------------------------------

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
