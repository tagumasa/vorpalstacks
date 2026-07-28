package rds

import (
	"fmt"
	"strings"
	"time"

	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
)

type Endpoint struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

type DBCluster struct {
	DBClusterIdentifier              string                            `json:"DBClusterIdentifier"`
	Engine                           string                            `json:"Engine"`
	EngineVersion                    string                            `json:"EngineVersion,omitempty"`
	Status                           string                            `json:"Status"`
	MasterUsername                   string                            `json:"MasterUsername,omitempty"`
	DatabaseName                     string                            `json:"DatabaseName,omitempty"`
	Port                             int                               `json:"Port"`
	BackupRetentionPeriod            int                               `json:"BackupRetentionPeriod"`
	PreferredBackupWindow            string                            `json:"PreferredBackupWindow,omitempty"`
	PreferredMaintenanceWindow       string                            `json:"PreferredMaintenanceWindow,omitempty"`
	AvailabilityZones                []string                          `json:"AvailabilityZones,omitempty"`
	MultiAZ                          bool                              `json:"MultiAZ"`
	VpcSecurityGroupIds              []string                          `json:"VpcSecurityGroupIds,omitempty"`
	DBSubnetGroupName                string                            `json:"DBSubnetGroup,omitempty"`
	DBClusterParameterGroupName      string                            `json:"DBClusterParameterGroup,omitempty"`
	StorageEncrypted                 bool                              `json:"StorageEncrypted"`
	KmsKeyId                         string                            `json:"KmsKeyId,omitempty"`
	CopyTagsToSnapshot               bool                              `json:"CopyTagsToSnapshot"`
	DeletionProtection               bool                              `json:"DeletionProtection"`
	EnabledCloudwatchLogsExports     []string                          `json:"EnabledCloudwatchLogsExports,omitempty"`
	IAMDatabaseAuthenticationEnabled bool                              `json:"IAMDatabaseAuthenticationEnabled"`
	ClusterCreateTime                *time.Time                        `json:"ClusterCreateTime,omitempty"`
	EarliestRestorableTime           *time.Time                        `json:"EarliestRestorableTime,omitempty"`
	LatestRestorableTime             *time.Time                        `json:"LatestRestorableTime,omitempty"`
	AssociatedRoles                  []DBClusterRole                   `json:"AssociatedRoles,omitempty"`
	ReplicationSourceIdentifier      string                            `json:"ReplicationSourceIdentifier,omitempty"`
	GlobalClusterIdentifier          string                            `json:"GlobalClusterIdentifier,omitempty"`
	StorageType                      string                            `json:"StorageType,omitempty"`
	ServerlessV2ScalingConfiguration *ServerlessV2ScalingConfiguration `json:"ServerlessV2ScalingConfiguration,omitempty"`
	AccountID                        string                            `json:"AccountId,omitempty"`
	Region                           string                            `json:"Region,omitempty"`
	DBClusterArn                     string                            `json:"DBClusterArn,omitempty"`
	Endpoint                         *Endpoint                         `json:"Endpoint,omitempty"`

	// MasterUserPasswordHash stores the bcrypt hash of the master user
	// password. It is persisted but never serialised to API responses
	// (write-only per AWS spec). json tag omits it from protocol output
	// but keeps it in Pebble persistence (H1 fix).
	MasterUserPasswordHash string `json:"-"`

	// AWS-standard DBCluster output fields previously dropped (M5 fix).
	// These are populated at create time and surfaced through
	// DescribeDBClusters so that SDK clients see a complete response.
	AllocatedStorage       int32                         `json:"AllocatedStorage,omitempty"`
	DBClusterMembers       []DBClusterMember             `json:"DBClusterMembers,omitempty"`
	ReaderEndpoint         *Endpoint                     `json:"ReaderEndpoint,omitempty"`
	DbClusterResourceId    string                        `json:"DbClusterResourceId,omitempty"`
	NetworkType            string                        `json:"NetworkType,omitempty"`
	PercentProgress        string                        `json:"PercentProgress,omitempty"`
	PendingModifiedValues  *ClusterPendingModifiedValues `json:"PendingModifiedValues,omitempty"`
	HostedZoneId           string                        `json:"HostedZoneId,omitempty"`
	ReadReplicaIdentifiers []string                      `json:"ReadReplicaIdentifiers,omitempty"`
	AutomaticRestartTime   *time.Time                    `json:"AutomaticRestartTime,omitempty"`
}

type ServerlessV2ScalingConfiguration struct {
	MinCapacity float64 `json:"MinCapacity"`
	MaxCapacity float64 `json:"MaxCapacity"`
}

type DBClusterRole struct {
	RoleArn     string `json:"RoleArn"`
	FeatureName string `json:"FeatureName,omitempty"`
	Status      string `json:"Status"`
}

// DBClusterMember represents a DB instance that belongs to a DB cluster.
// Surfaced in DBCluster output as the DBClusterMembers list (M5 fix).
type DBClusterMember struct {
	DBInstanceIdentifier          string `json:"DBInstanceIdentifier"`
	IsClusterWriter               bool   `json:"IsClusterWriter"`
	DBClusterParameterGroupStatus string `json:"DBClusterParameterGroupStatus,omitempty"`
	PromotionTier                 int32  `json:"PromotionTier,omitempty"`
}

// ClusterPendingModifiedValues holds pending changes that have not yet been
// applied to the DB cluster. Populated when ApplyImmediately=false during
// ModifyDBCluster (M5 fix).
type ClusterPendingModifiedValues struct {
	DBClusterIdentifier              string `json:"DBClusterIdentifier,omitempty"`
	IAMDatabaseAuthenticationEnabled *bool  `json:"IAMDatabaseAuthenticationEnabled,omitempty"`
	EngineVersion                    string `json:"EngineVersion,omitempty"`
	BackupRetentionPeriod            *int   `json:"BackupRetentionPeriod,omitempty"`
	StorageType                      string `json:"StorageType,omitempty"`
	AllocatedStorage                 *int32 `json:"AllocatedStorage,omitempty"`
	Iops                             *int32 `json:"Iops,omitempty"`
	NetworkType                      string `json:"NetworkType,omitempty"`
}

type DBInstance struct {
	DBInstanceIdentifier             string     `json:"DBInstanceIdentifier"`
	DBClusterIdentifier              string     `json:"DBClusterIdentifier"`
	Engine                           string     `json:"Engine"`
	EngineVersion                    string     `json:"EngineVersion,omitempty"`
	DBInstanceClass                  string     `json:"DBInstanceClass"`
	DBInstanceStatus                 string     `json:"DBInstanceStatus"`
	AvailabilityZone                 string     `json:"AvailabilityZone,omitempty"`
	PreferredMaintenanceWindow       string     `json:"PreferredMaintenanceWindow,omitempty"`
	DBParameterGroupName             string     `json:"DBParameterGroupName,omitempty"`
	DBSecurityGroups                 []string   `json:"DBSecurityGroups,omitempty"`
	VpcSecurityGroupIds              []string   `json:"VpcSecurityGroupIds,omitempty"`
	DBSubnetGroupName                string     `json:"DBSubnetGroupName,omitempty"`
	EnabledCloudwatchLogsExports     []string   `json:"EnabledCloudwatchLogsExports,omitempty"`
	IAMDatabaseAuthenticationEnabled bool       `json:"IAMDatabaseAuthenticationEnabled"`
	PubliclyAccessible               bool       `json:"PubliclyAccessible"`
	AutoMinorVersionUpgrade          bool       `json:"AutoMinorVersionUpgrade"`
	InstanceCreateTime               *time.Time `json:"InstanceCreateTime,omitempty"`
	CopyTagsToSnapshot               bool       `json:"CopyTagsToSnapshot"`
	AccountID                        string     `json:"AccountId,omitempty"`
	Region                           string     `json:"Region,omitempty"`
	DBInstanceArn                    string     `json:"DBInstanceArn,omitempty"`
	Endpoint                         *Endpoint  `json:"Endpoint,omitempty"`

	// AWS-standard DBInstance fields. These are captured from
	// CreateDBInstance parameters, persisted, and surfaced through
	// DescribeDBInstance / RestoreDBInstanceFromDBSnapshot. They were
	// previously dropped on the floor (RDS-5/RDS-20).
	AllocatedStorage                   int32      `json:"AllocatedStorage,omitempty"`
	MasterUsername                     string     `json:"MasterUsername,omitempty"`
	StorageType                        string     `json:"StorageType,omitempty"`
	BackupRetentionPeriod              int32      `json:"BackupRetentionPeriod,omitempty"`
	LicenseModel                       string     `json:"LicenseModel,omitempty"`
	StorageEncrypted                   bool       `json:"StorageEncrypted"`
	KmsKeyId                           string     `json:"KmsKeyId,omitempty"`
	DeletionProtection                 bool       `json:"DeletionProtection"`
	MultiAZ                            bool       `json:"MultiAZ"`
	SecondaryAvailabilityZone          string     `json:"SecondaryAvailabilityZone,omitempty"`
	Port                               int32      `json:"Port,omitempty"`
	OptionGroupName                    string     `json:"OptionGroupName,omitempty"`
	VpcId                              string     `json:"VpcId,omitempty"`
	Iops                               int32      `json:"Iops,omitempty"`
	MaxAllocatedStorage                int32      `json:"MaxAllocatedStorage,omitempty"`
	StorageThroughput                  int32      `json:"StorageThroughput,omitempty"`
	MonitoringInterval                 int32      `json:"MonitoringInterval,omitempty"`
	EnhancedMonitoringResourceArn      string     `json:"EnhancedMonitoringResourceArn,omitempty"`
	PerformanceInsightsEnabled         bool       `json:"PerformanceInsightsEnabled"`
	PerformanceInsightsKMSKeyId        string     `json:"PerformanceInsightsKMSKeyId,omitempty"`
	PerformanceInsightsRetentionPeriod int32      `json:"PerformanceInsightsRetentionPeriod,omitempty"`
	CACertificateIdentifier            string     `json:"CACertificateIdentifier,omitempty"`
	DbiResourceId                      string     `json:"DbiResourceId,omitempty"`
	LatestRestorableTime               *time.Time `json:"LatestRestorableTime,omitempty"`
	PreferredBackupWindow              string     `json:"PreferredBackupWindow,omitempty"`

	// MasterUserPasswordHash stores the bcrypt hash of the master user
	// password for the instance. Write-only: persisted but never surfaced
	// in API responses (H1 fix).
	MasterUserPasswordHash string `json:"-"`
}

type DBInstanceSnapshot struct {
	DBSnapshotIdentifier   string      `json:"DBSnapshotIdentifier"`
	DBInstanceIdentifier   string      `json:"DBInstanceIdentifier"`
	SnapshotCreateTime     *time.Time  `json:"SnapshotCreateTime,omitempty"`
	Engine                 string      `json:"Engine"`
	EngineVersion          string      `json:"EngineVersion,omitempty"`
	SnapshotType           string      `json:"SnapshotType,omitempty"`
	Status                 string      `json:"Status"`
	AllocatedStorage       int32       `json:"AllocatedStorage"`
	StorageType            string      `json:"StorageType,omitempty"`
	Port                   int32       `json:"Port"`
	AvailabilityZone       string      `json:"AvailabilityZone,omitempty"`
	VpcId                  string      `json:"VpcId,omitempty"`
	InstanceCreateTime     *time.Time  `json:"InstanceCreateTime,omitempty"`
	MasterUsername         string      `json:"MasterUsername,omitempty"`
	LicenseModel           string      `json:"LicenseModel,omitempty"`
	StorageEncrypted       bool        `json:"StorageEncrypted"`
	KmsKeyId               string      `json:"KmsKeyId,omitempty"`
	DBSnapshotArn          string      `json:"DBSnapshotArn,omitempty"`
	RestoreAttributeValues []string    `json:"RestoreAttributeValues,omitempty"`
	IAMDatabaseAuthEnabled bool        `json:"IAMDatabaseAuthenticationEnabled"`
	OptionGroupName        string      `json:"OptionGroupName,omitempty"`
	TagList                []types.Tag `json:"TagList,omitempty"`
	AccountID              string      `json:"AccountId,omitempty"`
	Region                 string      `json:"Region,omitempty"`
}

type DBClusterSnapshot struct {
	DBClusterSnapshotIdentifier string     `json:"DBClusterSnapshotIdentifier"`
	DBClusterIdentifier         string     `json:"DBClusterIdentifier"`
	SnapshotCreateTime          *time.Time `json:"SnapshotCreateTime,omitempty"`
	Engine                      string     `json:"Engine"`
	EngineVersion               string     `json:"EngineVersion,omitempty"`
	SnapshotType                string     `json:"SnapshotType,omitempty"`
	Status                      string     `json:"Status"`
	Port                        int        `json:"Port"`
	VpcId                       string     `json:"VpcId,omitempty"`
	ClusterCreateTime           *time.Time `json:"ClusterCreateTime,omitempty"`
	StorageEncrypted            bool       `json:"StorageEncrypted"`
	KmsKeyId                    string     `json:"KmsKeyId,omitempty"`
	DBSnapshotArn               string     `json:"DBClusterSnapshotArn,omitempty"`
	RestoreAttributeValues      []string   `json:"RestoreAttributeValues,omitempty"`
	AccountID                   string     `json:"AccountId,omitempty"`
	Region                      string     `json:"Region,omitempty"`

	// AWS-standard DBClusterSnapshot fields captured from the source
	// cluster at snapshot time. Previously dropped (RDS-2).
	MasterUsername                   string `json:"MasterUsername,omitempty"`
	AllocatedStorage                 int32  `json:"AllocatedStorage,omitempty"`
	StorageType                      string `json:"StorageType,omitempty"`
	LicenseModel                     string `json:"LicenseModel,omitempty"`
	IAMDatabaseAuthenticationEnabled bool   `json:"IAMDatabaseAuthenticationEnabled"`
}

type DBClusterParameterGroup struct {
	DBClusterParameterGroupName string      `json:"DBClusterParameterGroupName"`
	DBParameterGroupFamily      string      `json:"DBParameterGroupFamily"`
	Description                 string      `json:"Description"`
	ARN                         string      `json:"DBClusterParameterGroupArn,omitempty"`
	Parameters                  []Parameter `json:"Parameters,omitempty"`
}

type DBParameterGroup struct {
	DBParameterGroupName   string      `json:"DBParameterGroupName"`
	DBParameterGroupFamily string      `json:"DBParameterGroupFamily"`
	Description            string      `json:"Description"`
	ARN                    string      `json:"DBParameterGroupArn,omitempty"`
	Parameters             []Parameter `json:"Parameters,omitempty"`
}

type DBSubnetGroup struct {
	DBSubnetGroupName        string   `json:"DBSubnetGroupName"`
	DBSubnetGroupDescription string   `json:"DBSubnetGroupDescription"`
	VpcId                    string   `json:"VpcId"`
	SubnetGroupStatus        string   `json:"SubnetGroupStatus"`
	Subnets                  []Subnet `json:"Subnets"`
	ARN                      string   `json:"DBSubnetGroupArn,omitempty"`
}

type Subnet struct {
	SubnetIdentifier       string `json:"SubnetIdentifier"`
	SubnetAvailabilityZone string `json:"SubnetAvailabilityZone"`
	SubnetOutpost          string `json:"SubnetOutpost,omitempty"`
	SubnetStatus           string `json:"SubnetStatus"`
}

type GlobalCluster struct {
	GlobalClusterIdentifier string                `json:"GlobalClusterIdentifier"`
	GlobalClusterResourceId string                `json:"GlobalClusterResourceId"`
	GlobalClusterArn        string                `json:"GlobalClusterArn,omitempty"`
	Engine                  string                `json:"Engine"`
	EngineVersion           string                `json:"EngineVersion,omitempty"`
	Status                  string                `json:"Status"`
	StorageEncrypted        bool                  `json:"StorageEncrypted"`
	DeletionProtection      bool                  `json:"DeletionProtection"`
	GlobalClusterMembers    []GlobalClusterMember `json:"GlobalClusterMembers,omitempty"`
	AccountID               string                `json:"AccountId,omitempty"`
	Region                  string                `json:"Region,omitempty"`
}

type GlobalClusterMember struct {
	DBClusterArn            string   `json:"DBClusterArn,omitempty"`
	Readers                 []string `json:"Readers,omitempty"`
	IsWriter                bool     `json:"IsWriter"`
	GlobalClusterIdentifier string   `json:"GlobalClusterIdentifier,omitempty"`
}

type DBClusterEndpoint struct {
	DBClusterEndpointIdentifier string   `json:"DBClusterEndpointIdentifier"`
	DBClusterIdentifier         string   `json:"DBClusterIdentifier"`
	Endpoint                    string   `json:"Endpoint"`
	Status                      string   `json:"Status"`
	EndpointType                string   `json:"EndpointType"`
	ExcludedMembers             []string `json:"ExcludedMembers,omitempty"`
	StaticMembers               []string `json:"StaticMembers,omitempty"`
	DBClusterEndpointArn        string   `json:"DBClusterEndpointArn,omitempty"`
}

type EventSubscription struct {
	CustSubscriptionId       string     `json:"CustSubscriptionId"`
	SnsTopicArn              string     `json:"SnsTopicArn"`
	Status                   string     `json:"Status"`
	SubscriptionCreationTime *time.Time `json:"SubscriptionCreationTime,omitempty"`
	SourceType               string     `json:"SourceType,omitempty"`
	SourceIdsList            []string   `json:"SourceIdsList,omitempty"`
	EventCategoriesList      []string   `json:"EventCategoriesList,omitempty"`
	Enabled                  bool       `json:"Enabled"`
	CustSubscriptionArn      string     `json:"CustSubscriptionArn,omitempty"`
}

type Parameter struct {
	ParameterName        string `json:"ParameterName"`
	ParameterValue       string `json:"ParameterValue,omitempty"`
	Description          string `json:"Description,omitempty"`
	Source               string `json:"Source,omitempty"`
	ApplyType            string `json:"ApplyType,omitempty"`
	DataType             string `json:"DataType,omitempty"`
	AllowedValues        string `json:"AllowedValues,omitempty"`
	IsModifiable         bool   `json:"IsModifiable"`
	MinimumEngineVersion string `json:"MinimumEngineVersion,omitempty"`
	ApplyMethod          string `json:"ApplyMethod,omitempty"`
}

func ClusterParameterGroupARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("rds", "cluster-pg:"+name)
}

func ParameterGroupARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("rds", "pg:"+name)
}

func SubnetGroupARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("rds", "subgrp:"+name)
}

func EventSubscriptionARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("rds", "es:"+name)
}

func ResourceTagKey(accountID, region, resourceArn string) string {
	return fmt.Sprintf("rds:tags/%s/%s/%s", accountID, region, resourceArn)
}

func NormaliseResourceName(arn string) string {
	arn = strings.TrimPrefix(arn, "arn:aws:rds:")
	parts := strings.SplitN(arn, ":", 4)
	if len(parts) >= 4 {
		return parts[3]
	}
	return arn
}
