package neptune

import (
	rds "vorpalstacks/internal/store/aws/rds"
)

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

type Endpoint = rds.Endpoint
type DBCluster = rds.DBCluster
type DBClusterRole = rds.DBClusterRole
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
