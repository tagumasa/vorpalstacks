package rds

import (
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	storerds "vorpalstacks/internal/store/aws/rds"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// CreateClusterParams captures the common fields for cluster creation
// shared between the admin handler Core and the Neptune HTTP handler.
type CreateClusterParams struct {
	DBClusterIdentifier          string
	Engine                       string
	EngineVersion                string
	DatabaseName                 string
	MasterUsername               string
	Port                         int
	BackupRetentionPeriod        int
	AvailabilityZones            []string
	DBSubnetGroupName            string
	DBClusterParameterGroupName  string
	StorageEncrypted             bool
	CopyTagsToSnapshot           bool
	DeletionProtection           bool
	IAMDatabaseAuthentication    bool
	EnabledCloudwatchLogsExports []string
	AccountID                    string
	Region                       string
}

// ValidateCreateClusterParams validates all common cluster creation inputs
// and checks that referenced subnet groups and parameter groups exist.
// Returns a plain error; callers wrap with their own error type.
func ValidateCreateClusterParams(store storerds.StoreInterface, p CreateClusterParams) error {
	if p.DBClusterIdentifier == "" {
		return fmt.Errorf("DBClusterIdentifier is required")
	}
	if err := ValidateDBClusterIdentifier(p.DBClusterIdentifier); err != nil {
		return err
	}
	if p.Engine == "" {
		return fmt.Errorf("Engine is required")
	}
	if err := ValidateEngine(p.Engine); err != nil {
		return err
	}
	engineVersion := p.EngineVersion
	if engineVersion == "" {
		engineVersion = DefaultEngineVersion(p.Engine)
	}
	if err := ValidateEngineVersion(p.Engine, engineVersion); err != nil {
		return err
	}
	if p.DatabaseName != "" {
		if err := ValidateDatabaseName(p.DatabaseName); err != nil {
			return err
		}
	}
	if err := ValidatePort(int32(p.Port)); err != nil {
		return err
	}
	if err := ValidateBackupRetentionPeriod(int32(p.BackupRetentionPeriod)); err != nil {
		return err
	}
	if p.DBSubnetGroupName != "" {
		if _, err := store.GetSubnetGroup(p.DBSubnetGroupName); err != nil {
			return err
		}
	}
	if p.DBClusterParameterGroupName != "" {
		if _, err := store.GetClusterParameterGroup(p.DBClusterParameterGroupName); err != nil {
			return err
		}
	}
	return nil
}

// ResolveEngineVersion returns the effective engine version, applying the
// default when none is specified.
func ResolveEngineVersion(engine, requested string) string {
	if requested == "" {
		return DefaultEngineVersion(engine)
	}
	return requested
}

// BuildCluster constructs a *storerds.DBCluster from the given params,
// setting common fields shared between admin and HTTP layers. The
// returned cluster has Status="creating"; callers set additional
// Neptune-specific fields (PreferredBackupWindow, KmsKeyId, etc.) as needed.
func BuildCluster(p CreateClusterParams) *storerds.DBCluster {
	now := time.Now()
	engineVersion := ResolveEngineVersion(p.Engine, p.EngineVersion)
	return &storerds.DBCluster{
		DBClusterIdentifier:              p.DBClusterIdentifier,
		Engine:                           p.Engine,
		EngineVersion:                    engineVersion,
		Status:                           "creating",
		MasterUsername:                   p.MasterUsername,
		DatabaseName:                     p.DatabaseName,
		Port:                             p.Port,
		BackupRetentionPeriod:            p.BackupRetentionPeriod,
		AvailabilityZones:                p.AvailabilityZones,
		DBSubnetGroupName:                p.DBSubnetGroupName,
		DBClusterParameterGroupName:      p.DBClusterParameterGroupName,
		StorageEncrypted:                 p.StorageEncrypted,
		CopyTagsToSnapshot:               p.CopyTagsToSnapshot,
		DeletionProtection:               p.DeletionProtection,
		IAMDatabaseAuthenticationEnabled: p.IAMDatabaseAuthentication,
		EnabledCloudwatchLogsExports:     p.EnabledCloudwatchLogsExports,
		ClusterCreateTime:                &now,
		AccountID:                        p.AccountID,
		Region:                           p.Region,
		DBClusterArn:                     arnutil.NewARNBuilder(p.AccountID, p.Region).RDS().Cluster(p.DBClusterIdentifier),
	}
}

// DeleteClusterParams captures the common fields for cluster deletion.
type DeleteClusterParams struct {
	DBClusterIdentifier       string
	SkipFinalSnapshot         bool
	FinalDBSnapshotIdentifier string
	AccountID                 string
	Region                    string
}

// ValidateDeleteClusterParams validates cluster deletion inputs and checks
// deletion protection. It does NOT check SkipFinalSnapshot/FinalSnapshot
// combination logic (callers handle that separately for protocol-specific
// error codes).
func ValidateDeleteClusterParams(store storerds.StoreInterface, p DeleteClusterParams) (*storerds.DBCluster, error) {
	if p.DBClusterIdentifier == "" {
		return nil, awserrors.NewAWSError("InvalidParameterValue", "DBClusterIdentifier is required", http.StatusBadRequest)
	}
	cluster, err := store.GetCluster(p.DBClusterIdentifier)
	if err != nil {
		return nil, err
	}
	if cluster.DeletionProtection {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", "cannot delete cluster when DeletionProtection is enabled", http.StatusBadRequest)
	}
	return cluster, nil
}

// BuildFinalSnapshot creates a manual cluster snapshot for the final
// snapshot taken before cluster deletion.
func BuildFinalSnapshot(cluster *storerds.DBCluster, snapshotID, accountID, region string) *storerds.DBClusterSnapshot {
	now := time.Now()
	return &storerds.DBClusterSnapshot{
		DBClusterSnapshotIdentifier:      snapshotID,
		DBClusterIdentifier:              cluster.DBClusterIdentifier,
		SnapshotCreateTime:               &now,
		Engine:                           cluster.Engine,
		EngineVersion:                    cluster.EngineVersion,
		SnapshotType:                     "manual",
		Status:                           "available",
		Port:                             cluster.Port,
		ClusterCreateTime:                cluster.ClusterCreateTime,
		StorageEncrypted:                 cluster.StorageEncrypted,
		KmsKeyId:                         cluster.KmsKeyId,
		DBSnapshotArn:                    arnutil.NewARNBuilder(accountID, region).RDS().ClusterSnapshot(snapshotID),
		AccountID:                        accountID,
		Region:                           region,
		MasterUsername:                   cluster.MasterUsername,
		StorageType:                      cluster.StorageType,
		IAMDatabaseAuthenticationEnabled: cluster.IAMDatabaseAuthenticationEnabled,
	}
}

// DeleteInstanceParams captures the common fields for instance deletion.
type DeleteInstanceParams struct {
	DBInstanceIdentifier      string
	SkipFinalSnapshot         bool
	FinalDBSnapshotIdentifier string
	AccountID                 string
	Region                    string
}

// ValidateDeleteInstanceParams validates instance deletion inputs and checks
// deletion protection. Returns the instance if validation passes.
func ValidateDeleteInstanceParams(store storerds.StoreInterface, p DeleteInstanceParams) (*storerds.DBInstance, error) {
	if p.DBInstanceIdentifier == "" {
		return nil, awserrors.NewAWSError("InvalidParameterValue", "DBInstanceIdentifier is required", http.StatusBadRequest)
	}
	instance, err := store.GetInstance(p.DBInstanceIdentifier)
	if err != nil {
		return nil, err
	}
	if instance.DeletionProtection {
		return nil, awserrors.NewAWSError("InvalidDBInstanceStateFault", "cannot delete instance when DeletionProtection is enabled", http.StatusBadRequest)
	}
	return instance, nil
}

// BuildInstanceFinalSnapshot creates a manual instance snapshot for the
// final snapshot taken before instance deletion.
func BuildInstanceFinalSnapshot(instance *storerds.DBInstance, snapshotID, accountID, region string) *storerds.DBInstanceSnapshot {
	now := time.Now()
	snap := &storerds.DBInstanceSnapshot{
		DBSnapshotIdentifier:   snapshotID,
		DBInstanceIdentifier:   instance.DBInstanceIdentifier,
		InstanceCreateTime:     instance.InstanceCreateTime,
		Engine:                 instance.Engine,
		EngineVersion:          instance.EngineVersion,
		SnapshotType:           "manual",
		Status:                 "available",
		AvailabilityZone:       instance.AvailabilityZone,
		DBSnapshotArn:          arnutil.NewARNBuilder(accountID, region).RDS().Snapshot(snapshotID),
		IAMDatabaseAuthEnabled: instance.IAMDatabaseAuthenticationEnabled,
		AccountID:              accountID,
		Region:                 region,
		AllocatedStorage:       instance.AllocatedStorage,
		MasterUsername:         instance.MasterUsername,
		StorageType:            instance.StorageType,
		LicenseModel:           instance.LicenseModel,
		StorageEncrypted:       instance.StorageEncrypted,
		KmsKeyId:               instance.KmsKeyId,
		OptionGroupName:        instance.OptionGroupName,
		VpcId:                  instance.VpcId,
		SnapshotCreateTime:     &now,
	}
	if instance.Port > 0 {
		snap.Port = instance.Port
	} else if instance.Endpoint != nil {
		snap.Port = int32(instance.Endpoint.Port)
	}
	return snap
}

// CreateInstanceParams captures the common fields for instance creation
// shared between the admin handler Core and the Neptune HTTP handler.
type CreateInstanceParams struct {
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
	AccountID                          string
	Region                             string
}

// ValidateCreateInstanceParams validates all common instance creation
// inputs and checks that referenced subnet groups and parameter groups
// exist. Returns a plain error; callers wrap with their own error type.
func ValidateCreateInstanceParams(store storerds.StoreInterface, p CreateInstanceParams) error {
	if p.DBInstanceIdentifier == "" {
		return fmt.Errorf("DBInstanceIdentifier is required")
	}
	if err := ValidateDBInstanceIdentifier(p.DBInstanceIdentifier); err != nil {
		return err
	}
	if err := ValidateDBInstanceClass(p.DBInstanceClass); err != nil {
		return err
	}
	if p.Engine == "" {
		return fmt.Errorf("Engine is required")
	}
	if err := ValidateEngine(p.Engine); err != nil {
		return err
	}
	engineVersion := p.EngineVersion
	if engineVersion == "" {
		engineVersion = DefaultEngineVersion(p.Engine)
	}
	if err := ValidateEngineVersion(p.Engine, engineVersion); err != nil {
		return err
	}
	if err := ValidatePort(p.Port); err != nil {
		return err
	}
	if err := ValidateMonitoringInterval(p.MonitoringInterval); err != nil {
		return err
	}
	if err := ValidateStorageType(p.StorageType, p.Engine); err != nil {
		return err
	}
	if err := ValidateBackupRetentionPeriod(p.BackupRetentionPeriod); err != nil {
		return err
	}
	if err := ValidateAllocatedStorage(p.AllocatedStorage, p.Engine); err != nil {
		return err
	}
	if p.DBSubnetGroupName != "" {
		if _, err := store.GetSubnetGroup(p.DBSubnetGroupName); err != nil {
			return err
		}
	}
	if p.DBParameterGroupName != "" {
		if _, err := store.GetParameterGroup(p.DBParameterGroupName); err != nil {
			return err
		}
	}
	return nil
}

// BuildInstance constructs a *storerds.DBInstance from the given params,
// setting common fields shared between admin and HTTP layers. The
// returned instance has DBInstanceStatus="creating".
func BuildInstance(p CreateInstanceParams) *storerds.DBInstance {
	now := time.Now()
	engineVersion := ResolveEngineVersion(p.Engine, p.EngineVersion)
	return &storerds.DBInstance{
		DBInstanceIdentifier:               p.DBInstanceIdentifier,
		DBClusterIdentifier:                p.DBClusterIdentifier,
		Engine:                             p.Engine,
		EngineVersion:                      engineVersion,
		DBInstanceClass:                    p.DBInstanceClass,
		DBInstanceStatus:                   "creating",
		AvailabilityZone:                   p.AvailabilityZone,
		PreferredMaintenanceWindow:         p.PreferredMaintenanceWindow,
		PreferredBackupWindow:              p.PreferredBackupWindow,
		DBParameterGroupName:               p.DBParameterGroupName,
		DBSubnetGroupName:                  p.DBSubnetGroupName,
		PubliclyAccessible:                 p.PubliclyAccessible,
		AutoMinorVersionUpgrade:            p.AutoMinorVersionUpgrade,
		InstanceCreateTime:                 &now,
		AccountID:                          p.AccountID,
		Region:                             p.Region,
		DBInstanceArn:                      arnutil.NewARNBuilder(p.AccountID, p.Region).RDS().DBInstance(p.DBInstanceIdentifier),
		DbiResourceId:                      generateDbiResourceId(),
		AllocatedStorage:                   p.AllocatedStorage,
		MasterUsername:                     p.MasterUsername,
		StorageType:                        p.StorageType,
		BackupRetentionPeriod:              p.BackupRetentionPeriod,
		LicenseModel:                       p.LicenseModel,
		StorageEncrypted:                   p.StorageEncrypted,
		KmsKeyId:                           p.KmsKeyId,
		DeletionProtection:                 p.DeletionProtection,
		MultiAZ:                            p.MultiAZ,
		Port:                               p.Port,
		OptionGroupName:                    p.OptionGroupName,
		Iops:                               p.Iops,
		MaxAllocatedStorage:                p.MaxAllocatedStorage,
		StorageThroughput:                  p.StorageThroughput,
		MonitoringInterval:                 p.MonitoringInterval,
		EnhancedMonitoringResourceArn:      p.MonitoringRoleArn,
		PerformanceInsightsEnabled:         p.EnablePerformanceInsights,
		PerformanceInsightsKMSKeyId:        p.PerformanceInsightsKMSKeyId,
		PerformanceInsightsRetentionPeriod: p.PerformanceInsightsRetentionPeriod,
		CACertificateIdentifier:            p.CACertificateIdentifier,
		CopyTagsToSnapshot:                 p.CopyTagsToSnapshot,
		EnabledCloudwatchLogsExports:       p.EnabledCloudwatchLogsExports,
		IAMDatabaseAuthenticationEnabled:   p.EnableIAMDatabaseAuthentication,
		VpcSecurityGroupIds:                p.VpcSecurityGroupIds,
		DBSecurityGroups:                   p.DBSecurityGroups,
	}
}
