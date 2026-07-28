package neptune

import (
	rds "vorpalstacks/internal/store/aws/rds"
)

type DBClusterEndpoint = rds.DBClusterEndpoint

type Endpoint = rds.Endpoint
type DBCluster = rds.DBCluster
type DBClusterRole = rds.DBClusterRole
type DBClusterMember = rds.DBClusterMember
type ClusterPendingModifiedValues = rds.ClusterPendingModifiedValues
type ServerlessV2ScalingConfiguration = rds.ServerlessV2ScalingConfiguration
type DBInstance = rds.DBInstance
type DBClusterSnapshot = rds.DBClusterSnapshot
type DBClusterParameterGroup = rds.DBClusterParameterGroup
type DBParameterGroup = rds.DBParameterGroup
type DBSubnetGroup = rds.DBSubnetGroup
type Subnet = rds.Subnet
type GlobalCluster = rds.GlobalCluster
type GlobalClusterMember = rds.GlobalClusterMember
type EventSubscription = rds.EventSubscription
type Parameter = rds.Parameter

var ClusterParameterGroupARN = rds.ClusterParameterGroupARN
var ParameterGroupARN = rds.ParameterGroupARN
var SubnetGroupARN = rds.SubnetGroupARN
var EventSubscriptionARN = rds.EventSubscriptionARN
var ResourceTagKey = rds.ResourceTagKey
var NormaliseResourceName = rds.NormaliseResourceName
