package neptune

import (
	"fmt"
	"strings"
	"time"

	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// DBCluster represents a Neptune database cluster, including its configuration,
// status, and associated metadata.
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
}

// ServerlessV2ScalingConfiguration defines the minimum and maximum capacity for
// Neptune Serverless v2 scaling.
type ServerlessV2ScalingConfiguration struct {
	MinCapacity float64 `json:"MinCapacity"`
	MaxCapacity float64 `json:"MaxCapacity"`
}

// DBClusterRole represents an IAM role associated with a DB cluster.
type DBClusterRole struct {
	RoleArn     string `json:"RoleArn"`
	FeatureName string `json:"FeatureName,omitempty"`
	Status      string `json:"Status"`
}

// DBInstance represents a Neptune database instance within a cluster.
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
}

// DBClusterSnapshot represents a point-in-time backup of a Neptune DB cluster.
type DBClusterSnapshot struct {
	DBClusterSnapshotIdentifier string     `json:"DBClusterSnapshotIdentifier"`
	DBClusterIdentifier         string     `json:"DBClusterIdentifier"`
	SnapshotCreateTime          *time.Time `json:"SnapshotCreateTime,omitempty"`
	Engine                      string     `json:"Engine"`
	EngineVersion               string     `json:"EngineVersion,omitempty"`
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
}

// DBClusterParameterGroup defines a set of engine configuration parameters
// applied at the cluster level.
type DBClusterParameterGroup struct {
	DBClusterParameterGroupName string      `json:"DBClusterParameterGroupName"`
	DBParameterGroupFamily      string      `json:"DBParameterGroupFamily"`
	Description                 string      `json:"Description"`
	ARN                         string      `json:"DBClusterParameterGroupArn,omitempty"`
	Parameters                  []Parameter `json:"Parameters,omitempty"`
}

// DBParameterGroup defines a set of engine configuration parameters applied
// at the instance level.
type DBParameterGroup struct {
	DBParameterGroupName   string      `json:"DBParameterGroupName"`
	DBParameterGroupFamily string      `json:"DBParameterGroupFamily"`
	Description            string      `json:"Description"`
	ARN                    string      `json:"DBParameterGroupArn,omitempty"`
	Parameters             []Parameter `json:"Parameters,omitempty"`
}

// DBSubnetGroup specifies the VPC subnets used by DB instances within a cluster.
type DBSubnetGroup struct {
	DBSubnetGroupName        string   `json:"DBSubnetGroupName"`
	DBSubnetGroupDescription string   `json:"DBSubnetGroupDescription"`
	VpcId                    string   `json:"VpcId"`
	SubnetGroupStatus        string   `json:"SubnetGroupStatus"`
	Subnets                  []Subnet `json:"Subnets"`
	ARN                      string   `json:"DBSubnetGroupArn,omitempty"`
}

// Subnet represents a single VPC subnet belonging to a DB subnet group.
type Subnet struct {
	SubnetIdentifier       string `json:"SubnetIdentifier"`
	SubnetAvailabilityZone string `json:"SubnetAvailabilityZone"`
	SubnetOutpost          string `json:"SubnetOutpost,omitempty"`
	SubnetStatus           string `json:"SubnetStatus"`
}

// GlobalCluster represents a Neptune global database spanning multiple regions.
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

// GlobalClusterMember represents a regional DB cluster that participates in a
// global cluster.
type GlobalClusterMember struct {
	DBClusterArn            string   `json:"DBClusterArn,omitempty"`
	Readers                 []string `json:"Readers,omitempty"`
	IsWriter                bool     `json:"IsWriter"`
	GlobalClusterIdentifier string   `json:"GlobalClusterIdentifier,omitempty"`
}

// EventSubscription defines an SNS-based subscription for Neptune database events.
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

// Parameter represents a single engine configuration parameter with its
// metadata and constraints.
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

// DBClusterEndpoint represents a custom endpoint for connecting to a Neptune DB
// cluster with optional member filtering.
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

// ClusterParameterGroupARN returns the ARN for a DB cluster parameter group.
func ClusterParameterGroupARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("rds", "cluster-pg:"+name)
}

// ParameterGroupARN returns the ARN for a DB parameter group.
func ParameterGroupARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("rds", "pg:"+name)
}

// SubnetGroupARN returns the ARN for a DB subnet group.
func SubnetGroupARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("rds", "subgrp:"+name)
}

// EventSubscriptionARN returns the ARN for an event subscription.
func EventSubscriptionARN(accountID, region, name string) string {
	return arnutil.NewARNBuilder(accountID, region).Build("rds", "es:"+name)
}

// ResourceTagKey returns the storage key used to persist tags for a given resource ARN.
func ResourceTagKey(accountID, region, resourceArn string) string {
	return fmt.Sprintf("rds:tags/%s/%s/%s", accountID, region, resourceArn)
}

// NormaliseResourceName strips the ARN prefix and returns the resource name
// portion from a given ARN string.
func NormaliseResourceName(arn string) string {
	arn = strings.TrimPrefix(arn, "arn:aws:rds:")
	parts := strings.SplitN(arn, ":", 4)
	if len(parts) >= 4 {
		return parts[3]
	}
	return arn
}
